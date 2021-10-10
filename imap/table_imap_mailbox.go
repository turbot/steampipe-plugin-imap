package imap

import (
	"context"

	"github.com/emersion/go-imap"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableIMAPMailbox(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "imap_mailbox",
		Description: "Mailboxes in IMAP.",
		List: &plugin.ListConfig{
			Hydrate:    tableIMAPMailboxList,
			KeyColumns: plugin.OptionalColumns([]string{"name"}),
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the mailbox, e.g. 'INBOX', '[Gmail]/Drafts'."},
			// Other columns
			{Name: "attributes", Type: proto.ColumnType_JSON, Description: "Attributes set on the mailbox."},
			{Name: "delimiter", Type: proto.ColumnType_STRING, Description: "The server's path separator."},
			{Name: "flags", Type: proto.ColumnType_JSON, Hydrate: tableIMAPMailboxGet, Description: "The mailbox flags."},
			{Name: "permanent_flags", Type: proto.ColumnType_JSON, Hydrate: tableIMAPMailboxGet, Description: "The mailbox permanent flags."},
			{Name: "messages", Type: proto.ColumnType_INT, Hydrate: tableIMAPMailboxGet, Description: "The number of messages in this mailbox."},
			{Name: "recent", Type: proto.ColumnType_INT, Hydrate: tableIMAPMailboxGet, Description: "The number of messages not seen since the last time the mailbox was opened."},
			{Name: "unseen", Type: proto.ColumnType_INT, Hydrate: tableIMAPMailboxGet, Description: "The number of unread messages."},
			{Name: "read_only", Type: proto.ColumnType_BOOL, Hydrate: tableIMAPMailboxGet, Description: "True if the mailbox is open in read-only mode."},
		},
	}
}

func tableIMAPMailboxList(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	c, err := login(ctx, d)
	if err != nil {
		return nil, err
	}
	defer c.Logout()

	// Default to all mailboxes, but limit to the requested mailbox if given in a qual
	name := "*"
	if d.KeyColumnQuals["name"] != nil {
		name = d.KeyColumnQuals["name"].GetStringValue()
	}

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", name, mailboxes)
	}()

	for m := range mailboxes {
		d.StreamListItem(ctx, m)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return nil, nil
}

func tableIMAPMailboxGet(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	c, err := login(ctx, d)
	if err != nil {
		return nil, err
	}
	defer c.Logout()

	name := h.Item.(*imap.MailboxInfo).Name

	mboxDetail, err := c.Select(name, true)
	if err != nil {
		// Return an empty status instead of nil, so the attributes can be hydrated
		return &imap.MailboxStatus{}, nil
	}

	return mboxDetail, nil
}
