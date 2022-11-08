package imap

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/emersion/go-imap/client"

	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func login(ctx context.Context, d *plugin.QueryData) (*client.Client, error) {

	port := 993
	tlsEnabled := true
	insecureSkipVerify := false

	// Check env var settings
	host := os.Getenv("IMAP_HOST")
	login := os.Getenv("IMAP_LOGIN")
	password := os.Getenv("IMAP_PASSWORD")

	if portString, ok := os.LookupEnv("IMAP_PORT"); ok {
		p, err := strconv.Atoi(portString)
		if err == nil {
			port = p
		}
	}

	// Prefer config settings
	imapConfig := GetConfig(d.Connection)
	if imapConfig.Host != nil {
		host = *imapConfig.Host
	}
	if imapConfig.Port != nil {
		port = *imapConfig.Port
	}
	if imapConfig.Login != nil {
		login = *imapConfig.Login
	}
	if imapConfig.Password != nil {
		password = *imapConfig.Password
	}
	if imapConfig.TLSEnabled != nil {
		tlsEnabled = *imapConfig.TLSEnabled
	}
	if imapConfig.InsecureSkipVerify != nil {
		insecureSkipVerify = *imapConfig.InsecureSkipVerify
	}

	// Error if the minimum config is not set
	if host == "" {
		return nil, errors.New("host must be configured")
	}
	if login == "" {
		return nil, errors.New("login must be configured")
	}
	if password == "" {
		return nil, errors.New("password must be configured")
	}

	// Connect to server
	hostPort := fmt.Sprintf("%s:%d", host, port)
	var c *client.Client
	var err error
	if tlsEnabled {
		c, err = client.DialTLS(hostPort, &tls.Config{InsecureSkipVerify: insecureSkipVerify})
	} else {
		c, err = client.Dial(hostPort)
	}
	if err != nil {
		plugin.Logger(ctx).Error("connection_error", "host", host, "port", port, "hostPort", hostPort, "tlsEnabled", tlsEnabled, "login", login, "err", err)
		return nil, err
	}

	// Login
	if err := c.Login(login, password); err != nil {
		plugin.Logger(ctx).Error("connection_error", "host", host, "port", port, "hostPort", hostPort, "tlsEnabled", tlsEnabled, "login", login, "err", err)
		return nil, err
	}

	return c, nil
}
