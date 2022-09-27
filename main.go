package main

import (
	"github.com/turbot/steampipe-plugin-imap/imap"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: imap.Plugin})
}
