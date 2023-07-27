---
organization: Turbot
category: ["internet"]
icon_url: "/images/plugins/turbot/imap.svg"
brand_color: "#003FFF"
display_name: "IMAP"
short_name: "imap"
description: "Steampipe plugin to query mailboxes and messages using IMAP."
og_description: "Query IMAP with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/turbot/imap-social-graphic.png"
---

# IMAP + Steampipe

IMAP is a protocol for email access and management.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

List messages from your inbox:

```sql
select
  timestamp,
  from_email,
  subject
from
  imap_message;
```

```
+---------------------+---------------------------+---------------------+
| timestamp           | from_email                | subject             |
+---------------------+---------------------------+---------------------+
| 2021-06-04 10:04:09 | michael@dundermifflin.com | FW: Joke            |
| 2021-09-20 10:10:01 | dwight@dundermifflin.com  | Where's my stapler? |
| 2021-09-21 10:13:01 | pam@dundermifflin.com     | Tonight...          |
+---------------------+---------------------------+---------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/imap/tables)**

## Get started

### Install

Download and install the latest IMAP plugin:

```bash
steampipe plugin install imap
```

### Configuration

Installing the latest imap plugin will create a config file (`~/.steampipe/config/imap.spc`) with a single connection named `imap`.

Here is an example configuration for Gmail, which requires [allowing your user to use less secure apps](https://support.google.com/a/answer/6260879?hl=en).

```hcl
connection "imap" {
  plugin   = "imap"
  host     = "imap.gmail.com"
  port     = 993
  login    = "michael@dundermifflin.com"
  password = "Great Scott!"
}
```

- `host` - Hostname of the IMAP server. Required. Can also be set with the `IMAP_HOST` environment variable.
- `login` - Login name, usually the email address. Required. Can also be set with the `IMAP_LOGIN` environment variable.
- `password` - Password. Required. Can also be set with the `IMAP_PASSWORD` environment variable.
- `port` - Port to connect on the host, usually 143 for IMAP and 993 for IMAPS. Valid values are 143, 993, or a value between 1024 and 65535. Default 993. Can also be set with the `IMAP_PORT` environment variable.
- `tls_enabled` - If true, use TLS to connecto the host. Default true.
- `insecure_skip_verify` - If true, skip certificate verification. Default false.
- `mailbox` - The mailbox to query for messages if not specifically given in the query. Default is INBOX.

By default, variables in the configuration file will take precedence over any configured environment variables.

## Get involved

- Open source: https://github.com/turbot/steampipe-plugin-imap
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
