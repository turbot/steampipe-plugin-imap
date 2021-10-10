package imap

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type imapConfig struct {
	Host               *string `cty:"host"`
	Port               *int    `cty:"port"`
	Login              *string `cty:"login"`
	Password           *string `cty:"password"`
	TLSEnabled         *bool   `cty:"tls_enabled"`
	InsecureSkipVerify *bool   `cty:"insecure_skip_verify"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"host": {
		Type: schema.TypeString,
	},
	"port": {
		Type: schema.TypeInt,
	},
	"login": {
		Type: schema.TypeString,
	},
	"password": {
		Type: schema.TypeString,
	},
	"tls_enabled": {
		Type: schema.TypeBool,
	},
	"insecure_skip_verify": {
		Type: schema.TypeBool,
	},
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
