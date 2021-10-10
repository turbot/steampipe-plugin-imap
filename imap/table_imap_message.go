package imap

import (
	"context"
	"net/mail"
	"strings"

	"github.com/DusanKasan/parsemail"
	"github.com/emersion/go-imap"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableIMAPMessage(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "imap_message",
		Description: "Messagees in IMAP.",
		List: &plugin.ListConfig{
			Hydrate: tableIMAPMessageList,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "mailbox"},
				{Name: "size", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
				{Name: "timestamp", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
				{Name: "seq_num", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
				{Name: "from_email", Require: plugin.Optional},
				{Name: "sender", Require: plugin.Optional},
				{Name: "message_id", Require: plugin.Optional},
				{Name: "subject", Require: plugin.Optional},
				{Name: "query", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Envelope.Date"), Description: "Time when the message was sent."},
			{Name: "from_email", Type: proto.ColumnType_STRING, Hydrate: tableIMAPParsedMessage, Transform: transform.FromField("From").Transform(getFirstAddress), Description: "Email address, in lower case, of the first (and usually only) mailbox in the From header."},
			{Name: "subject", Type: proto.ColumnType_STRING, Hydrate: tableIMAPParsedMessage, Description: "Subject of the message."},
			// Other columns
			{Name: "attachments", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Transform: transform.FromField("Attachments").Transform(getAttachmentsWithoutData), Description: "Array of the names and content types of any attachments. The actual content of the attachment is not included."},
			{Name: "bcc", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Transform: transform.FromField("BCC"), Description: "The bcc field (where the 'BCC' means 'Blind Carbon Copy') contains addresses of recipients of the message whose addresses are not to be revealed to other recipients of the message."},
			{Name: "cc", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Transform: transform.FromField("CC"), Description: "The cc field (where the 'CC' means 'Carbon Copy' in the sense of making a copy on a typewriter using carbon paper) contains the addresses of others who are to receive the message, though the content of the message may not be directed at them."},
			{Name: "content_type", Type: proto.ColumnType_STRING, Hydrate: tableIMAPParsedMessage, Description: "Content type of the message."},
			{Name: "embedded_files", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Transform: transform.FromField("EmbeddedFiles").Transform(getEmbeddedFilesWithoutData), Description: "Array of content IDs and types for embedded files in the message. The actual content of the embedded file is not included."},
			{Name: "flags", Type: proto.ColumnType_JSON, Transform: transform.FromField("Flags"), Description: "Flags set on the message."},
			{Name: "from", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Transform: transform.FromField("From"), Description: "The from field specifies the author(s) of the message, that is, the mailbox(es) of the person(s) or system(s) responsible for the writing of the message."},
			{Name: "header", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Description: "Map of key value pairs for message header information."},
			{Name: "html_body", Type: proto.ColumnType_STRING, Hydrate: tableIMAPParsedMessage, Description: "HTML formatted body of the message."},
			{Name: "in_reply_to", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Description: "The message_id of the message to which this one is a reply (the parent message)."},
			{Name: "mailbox", Type: proto.ColumnType_STRING, Transform: transform.FromQual("mailbox"), Description: "Mailbox queried for messages."},
			{Name: "message_id", Type: proto.ColumnType_STRING, Hydrate: tableIMAPParsedMessage, Description: "Unique message identifier that refers to a particular version of a particular message."},
			{Name: "query", Type: proto.ColumnType_STRING, Transform: transform.FromQual("query"), Description: "Search query to match messages."},
			{Name: "references", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Description: "An array of message_id's for the parent and it's ancestors."},
			{Name: "reply_to", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Description: "An array of mailboxes that replies to this message should be sent to."},
			{Name: "sender", Type: proto.ColumnType_STRING, Hydrate: tableIMAPParsedMessage, Transform: transform.FromField("Sender.Address"), Description: "The mailbox of the agent responsible for the actual transmission of the message. For example, if a secretary were to send a message for another person, the mailbox of the secretary would appear in the sender field and the mailbox of the actual author would appear in the from field."},
			{Name: "seq_num", Type: proto.ColumnType_INT, Description: "Sequence number of the message."},
			{Name: "size", Type: proto.ColumnType_INT, Description: "Size in bytes of the message."},
			{Name: "text_body", Type: proto.ColumnType_STRING, Hydrate: tableIMAPParsedMessage, Description: "Text formatted body of the message."},
			{Name: "to", Type: proto.ColumnType_JSON, Hydrate: tableIMAPParsedMessage, Description: "List of mailboxes the message was sent to."},
		},
	}
}

func tableIMAPMessageList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	c, err := login(ctx, d)
	if err != nil {
		return nil, err
	}
	defer c.Logout()

	// Select the mailbox
	mailbox := d.KeyColumnQuals["mailbox"].GetStringValue()
	mbox, err := c.Select(mailbox, false)
	if err != nil {
		plugin.Logger(ctx).Error("imap_message.tableIMAPMessageList", "query_error", err, "mailbox", mailbox)
		return nil, err
	}

	// Convenience
	quals := d.Quals
	keyQuals := d.KeyColumnQuals

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
				break
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

	if keyQuals["sender"] != nil {
		criteria.Header.Add("Sender", keyQuals["sender"].GetStringValue())
	}

	if quals["timestamp"] != nil {
		for _, q := range quals["timestamp"].Quals {
			ts := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case ">":
				criteria.SentSince = ts
			case ">=":
				criteria.SentSince = ts
			case "=":
				criteria.SentSince = ts
				criteria.SentBefore = ts
			case "<=":
				criteria.SentBefore = ts
			case "<":
				criteria.SentBefore = ts
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

	plugin.Logger(ctx).Debug("imap_message.tableIMAPMessageList", "criteria", criteria)

	ids, err := c.Search(criteria)
	if err != nil {
		plugin.Logger(ctx).Error("imap_message.tableIMAPMessageList", "query_error", err, "criteria", criteria)
		return nil, err
	}

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
		d.StreamListItem(ctx, msg)
	}

	if err := <-done; err != nil {
		plugin.Logger(ctx).Error("imap_message.tableIMAPMessageList", "query_error", err, "mailbox", mailbox)
		return nil, err
	}

	return nil, nil
}

func tableIMAPParsedMessage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	msg := h.Item.(*imap.Message)
	r := msg.GetBody(&imap.BodySectionName{})
	email, err := parsemail.Parse(r)
	if err != nil {
		plugin.Logger(ctx).Error("imap_message.tableIMAPParsedMessage", "message_parse_error", err, "subject", msg.Envelope.Subject)
		return nil, err
	}
	return email, nil
}

type attachmentWithoutData struct {
	Filename    string
	ContentType string
}

type embeddedFileWithoutData struct {
	CID         string
	ContentType string
}

func getAttachmentsWithoutData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	items := d.Value.([]parsemail.Attachment)
	result := []attachmentWithoutData{}
	for _, i := range items {
		result = append(result, attachmentWithoutData{Filename: i.Filename, ContentType: i.ContentType})
	}
	return result, nil
}

func getEmbeddedFilesWithoutData(_ context.Context, d *transform.TransformData) (interface{}, error) {
	items := d.Value.([]parsemail.EmbeddedFile)
	result := []embeddedFileWithoutData{}
	for _, i := range items {
		result = append(result, embeddedFileWithoutData{CID: i.CID, ContentType: i.ContentType})
	}
	return result, nil
}

func getFirstAddress(_ context.Context, d *transform.TransformData) (interface{}, error) {
	items := d.Value.([]*mail.Address)
	if len(items) > 0 {
		return strings.ToLower(items[0].Address), nil
	}
	return nil, nil
}
