package imap

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-imap",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		ConnectionKeyColumns: []plugin.ConnectionKeyColumn{
			{
				Name:    "login",
				Hydrate: getLoginName,
			},
		},
		DefaultTransform: transform.FromGo(),
		TableMap: map[string]*plugin.Table{
			"imap_mailbox": tableIMAPMailbox(ctx),
			"imap_message": tableIMAPMessage(ctx),
		},
	}
	return p
}
