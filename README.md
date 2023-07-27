![image](https://hub.steampipe.io/images/plugins/turbot/imap-social-graphic.png)

# IMAP Plugin for Steampipe

Use SQL to query mailboxes, messages and more using IMAP.

- **[Get started →](https://hub.steampipe.io/plugins/turbot/imap)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/imap/tables)
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/turbot/steampipe-plugin-imap/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install imap
```

Run a query:

```sql
select
  timestamp,
  from_email,
  subject
from
  imap_message
where
  mailbox = 'INBOX';
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

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-imap.git
cd steampipe-plugin-imap
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/imap.spc
```

Try it!

```
steampipe query
> .inspect imap
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-imap/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [IMAP Plugin](https://github.com/turbot/steampipe-plugin-imap/labels/help%20wanted)
