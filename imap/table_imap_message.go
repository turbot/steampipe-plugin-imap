package imap

import (
	"context"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/emersion/go-imap"
	"github.com/jhillyerd/enmime"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableIMAPMessage(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "imap_message",
		Description: "Messages in IMAP.",
		List: &plugin.ListConfig{
			Hydrate: tableIMAPMessageList,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "mailbox", Require: plugin.Optional},
				{Name: "size", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
				{Name: "timestamp", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
				{Name: "seq_num", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
				{Name: "from_email", Require: plugin.Optional},
				{Name: "message_id", Require: plugin.Optional},
				{Name: "subject", Require: plugin.Optional},
				{Name: "query", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Hydrate: tableIMAPParsedMessage, Description: "Time when the message was sent."},
			{Name: "from_email", Type: proto.ColumnType_STRING, Hydrate: tableIMAPParsedMessage, Transform: transform.FromField("FromAddresses").Transform(getFirstAddress), Description: "Email address, in lower case, of the first (and usually only) mailbox in the From header."},
			{Name: "subject", Type: proto.ColumnType_STRING, Hydrate: tableIMAPParsedMessage, Transform: transform.FromField("Subject"), Description: "Subject of the message."},
			{Name: "message_id", Type: proto.ColumnType_STRING, Hydrate: tableIMAPParsedMessage, Description: "Unique message identifier that refers to a particular version of a particular message."},
			{Name: "to_addresses", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Description: "Array of To addresses."},
			{Name: "cc_addresses", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Description: "Array of CC addresses."},
			{Name: "bcc_addresses", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Description: "Array of BCC addresses."},
			{Name: "seq_num", Type: proto.ColumnType_INT, Transform: transform.FromField("Message.SeqNum"), Description: "Sequence number of the message."},
			{Name: "size", Type: proto.ColumnType_INT, Transform: transform.FromField("Message.Size"), Description: "Size in bytes of the message."},
			// Other columns
			{Name: "attachments", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Transform: transform.FromField("Envelope.Attachments").Transform(getAttachmentsWithoutData), Description: "All parts having a Content-Disposition of attachment."},
			{Name: "body_html", Type: proto.ColumnType_STRING, Hydrate: tableIMAPParsedMessage, Description: "HTML body of the message."},
			{Name: "body_text", Type: proto.ColumnType_STRING, Hydrate: tableIMAPParsedMessage, Description: "Text body of the message."},
			{Name: "embedded_files", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Transform: transform.FromField("Envelope.Inlines").Transform(getAttachmentsWithoutData), Description: "All parts having a Content-Disposition of inline."},
			{Name: "errors", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Transform: transform.FromField("Envelope.Errors"), Description: "Errors returned while parsing the email."},
			{Name: "flags", Type: proto.ColumnType_JSON, Transform: transform.FromField("Message.Flags"), Description: "Flags set on the message."},
			{Name: "from_addresses", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Description: "Array of From addresses."},
			{Name: "headers", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Transform: transform.FromField("Envelope.Root.Header"), Description: "Full set of headers defined in the message."},
			{Name: "in_reply_to", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Description: "Array of message IDs that this message is a reply to."},
			{Name: "mailbox", Type: proto.ColumnType_STRING, Description: "Mailbox queried for messages."},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromQual("query"), Description: "Search query to match messages."},
		},
	}
}

type msgWrapper struct {
	Message *imap.Message
	Mailbox string
}

type wrapper struct {
	Mailbox       string
	Envelope      *enmime.Envelope
	FromAddresses []*mail.Address
	ToAddresses   []*mail.Address
	CcAddresses   []*mail.Address
	BccAddresses  []*mail.Address
	Timestamp     time.Time
	InReplyTo     []string
	MessageID     string
	From          string
	Subject       string
	BodyText      string
	BodyHTML      string
}

func tableIMAPMessageList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	c, err := login(ctx, d)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = c.Logout()
	}()

	if err != nil {
		return nil, err
	}
	// Convenience
	quals := d.Quals
	keyQuals := d.KeyColumnQuals

	// Select the mailbox for queries:
	// 1. Check the mailbox qual
	// 2. Use the mailbox config setting
	// 3. Default to INBOX
	var mailbox string
	if keyQuals["mailbox"] != nil {
		mailbox = keyQuals["mailbox"].GetStringValue()
	} else if d.Connection != nil {
		imapConfig := GetConfig(d.Connection)
		if imapConfig.Mailbox != nil {
			mailbox = *imapConfig.Mailbox
		}
	}
	if mailbox == "" {
		mailbox = "INBOX"
	}

	mbox, err := c.Select(mailbox, false)
	if err != nil {
		plugin.Logger(ctx).Error("imap_message.tableIMAPMessageList", "query_error", err, "mailbox", mailbox)
		return nil, err
	}

	// Setup search criteria
	criteria := imap.NewSearchCriteria()

	// Limit by sequence number. Adding multiple ranges expands the range (or), so instead we collect the pieces
	from := uint32(1)
	to := mbox.Messages
	if quals["seq_num"] != nil {
		for _, q := range quals["seq_num"].Quals {
			sn := uint32(q.Value.GetInt64Value())
			switch q.Operator {
			case "=":
				from = sn
				to = sn
			case ">":
				if sn >= from {
					from = sn + 1
				}
			case ">=":
				if sn > from {
					from = sn
				}
			case "<":
				if sn <= to {
					to = sn - 1
				}
			case "<=":
				if sn < to {
					to = sn
				}
			}
		}
	}
	criteria.SeqNum = new(imap.SeqSet)
	criteria.SeqNum.AddRange(from, to)

	if keyQuals["query"] != nil {
		criteria.Text = append(criteria.Text, keyQuals["query"].GetStringValue())
	}

	if keyQuals["subject"] != nil {
		criteria.Header.Add("Subject", keyQuals["subject"].GetStringValue())
	}

	if keyQuals["message_id"] != nil {
		criteria.Header.Add("Message-Id", keyQuals["message_id"].GetStringValue())
	}

	if keyQuals["from_email"] != nil {
		criteria.Header.Add("From", keyQuals["from_email"].GetStringValue())
	}

	// Note: I don't know how these SentSince and SentBefore settings work,
	// because in testing they seem to be (kinda) on date boundaries, not on
	// timestamp boundaries. So, the only approach I could think to make it
	// accurate is to expand the search 24 hours each side and leave the Postgres
	// engine to do the final filtering.
	if quals["timestamp"] != nil {
		for _, q := range quals["timestamp"].Quals {
			ts := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case ">":
				criteria.SentSince = ts.Add(-24 * time.Hour)
			case ">=":
				criteria.SentSince = ts.Add(-24 * time.Hour)
			case "=":
				criteria.SentSince = ts.Add(-24 * time.Hour)
				criteria.SentBefore = ts.Add(24 * time.Hour)
			case "<=":
				criteria.SentBefore = ts.Add(24 * time.Hour)
			case "<":
				criteria.SentBefore = ts.Add(24 * time.Hour)
			}
		}
	}

	if quals["size"] != nil {
		for _, q := range quals["size"].Quals {
			size := q.Value.GetInt64Value()
			switch q.Operator {
			case "=":
				criteria.Larger = uint32(size) - 1
				criteria.Smaller = uint32(size) + 1
			case ">":
				criteria.Larger = uint32(size)
			case ">=":
				criteria.Larger = uint32(size) - 1
			case "<=":
				criteria.Smaller = uint32(size) + 1
			case "<":
				criteria.Smaller = uint32(size)
			}
		}
	}

	plugin.Logger(ctx).Warn("imap_message.tableIMAPMessageList", "criteria", criteria)

	ids, err := c.Search(criteria)
	if err != nil {
		plugin.Logger(ctx).Error("imap_message.tableIMAPMessageList", "query_error", err, "criteria", criteria)
		return nil, err
	}

	plugin.Logger(ctx).Warn("imap_message.tableIMAPMessageList", "ids", ids)

	if len(ids) == 0 {
		return nil, nil
	}

	searchSeqset := new(imap.SeqSet)
	limit := len(ids)
	if d.QueryContext.Limit != nil {
		i := int(*d.QueryContext.Limit)
		if i < limit {
			limit = i
		}
	}
	searchSeqset.AddNum(ids[:limit]...)

	section := &imap.BodySectionName{}
	fetchItems := append(imap.FetchFull.Expand(), section.FetchItem())

	messages := make(chan *imap.Message, limit)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(searchSeqset, fetchItems, messages)
	}()

	for msg := range messages {
		mw := msgWrapper{
			Message: msg,
			Mailbox: mailbox,
		}
		d.StreamListItem(ctx, mw)
	}

	if err := <-done; err != nil {
		plugin.Logger(ctx).Error("imap_message.tableIMAPMessageList", "query_error", err, "mailbox", mailbox)
		return nil, err
	}

	return nil, nil
}

func tableIMAPParsedMessage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	mw := h.Item.(msgWrapper)

	// Convenience
	msg := mw.Message

	r := msg.GetBody(&imap.BodySectionName{})

	// Parse message body
	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		plugin.Logger(ctx).Error("imap_message.tableIMAPParsedMessage", "CANNOT READ ENVELOPE", msg.Envelope.Subject)
		return nil, nil
	}

	te := wrapper{
		Mailbox:   mw.Mailbox,
		Envelope:  env,
		From:      env.GetHeader("From"),
		MessageID: env.GetHeader("Message-Id"),
		InReplyTo: env.GetHeaderValues("In-Reply-To"),
	}

	// Sometimes the subject has invalid UTF-8
	subject := env.GetHeader("Subject")
	if utf8.ValidString(subject) {
		te.Subject = subject
	}

	from, err := env.AddressList("From")
	if err == nil {
		te.FromAddresses = from
	}
	to, err := env.AddressList("To")
	if err == nil {
		te.ToAddresses = to
	}
	cc, err := env.AddressList("CC")
	if err == nil {
		te.CcAddresses = cc
	}
	bcc, err := env.AddressList("BCC")
	if err == nil {
		te.BccAddresses = bcc
	}
	dateString := env.GetHeader("Date")
	if dateString != "" {
		ts, err := mail.ParseDate(dateString)
		if err == nil {
			te.Timestamp = ts
		}
	}

	// It's common for emails to have non-UTF strings. They cause the gRPC layer
	// to fail, so filter that invalid data out here.
	if utf8.ValidString(env.Text) {
		te.BodyText = env.Text
	}
	if utf8.ValidString(env.HTML) {
		te.BodyHTML = env.HTML
	}

	return te, nil
}

type attachmentWithoutData struct {
	PartID string `json:"part_id,omitempty"`
	//Header            textproto.MIMEHeader `json:"header,omitempty"`
	Boundary          string            `json:"boundary,omitempty"`
	ContentID         string            `json:"content_id,omitempty"`
	ContentType       string            `json:"content_type,omitempty"`
	ContentTypeParams map[string]string `json:"content_type_params,omitempty"`
	Disposition       string            `json:"disposition,omitempty"`
	FileName          string            `json:"file_name,omitempty"`
	//FileModDate       time.Time         `json:"file_mod_date,omitempty"`
	Charset     string          `json:"charset,omitempty"`
	OrigCharset string          `json:"orig_charset,omitempty"`
	Errors      []*enmime.Error `json:"errors,omitempty"`
}

func getFirstAddress(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}
	items := d.Value.([]*mail.Address)
	if len(items) > 0 {
		return strings.ToLower(items[0].Address), nil
	}
	return nil, nil
}

func getAttachmentsWithoutData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}
	items := d.Value.([]*enmime.Part)
	result := []attachmentWithoutData{}
	for _, i := range items {
		result = append(result, attachmentWithoutData{
			PartID: i.PartID,
			//Header:            i.Header,
			Boundary:          i.Boundary,
			ContentID:         i.ContentID,
			ContentType:       i.ContentType,
			ContentTypeParams: i.ContentTypeParams,
			Disposition:       i.Disposition,
			FileName:          i.FileName,
			//FileModDate:       i.FileModDate,
			Charset:     i.Charset,
			OrigCharset: i.OrigCharset,
			Errors:      i.Errors,
		})
	}
	return result, nil
}
