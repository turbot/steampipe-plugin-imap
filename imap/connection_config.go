package imap

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type imapConfig struct {
	Host               *string `hcl:"host"`
	Port               *int    `hcl:"port"`
	Login              *string `hcl:"login"`
	Password           *string `hcl:"password"`
	TLSEnabled         *bool   `hcl:"tls_enabled"`
	InsecureSkipVerify *bool   `hcl:"insecure_skip_verify"`
	Mailbox            *string `hcl:"mailbox"`
}

func ConfigInstance() interface{} {
	return &imapConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) imapConfig {
	if connection == nil || connection.Config == nil {
		return imapConfig{}
	}
	config, _ := connection.Config.(imapConfig)
	return config
}
