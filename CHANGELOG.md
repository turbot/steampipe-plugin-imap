## v0.6.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/downloads), as a [Postgres FDW](https://steampipe.io/docs/steampipe_postgres/overview), as a [SQLite extension](https://steampipe.io/docs//steampipe_sqlite/overview) and as a standalone [exporter](https://steampipe.io/docs/steampipe_export/overview). ([#40](https://github.com/turbot/steampipe-plugin-imap/pull/40))
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension. ([#40](https://github.com/turbot/steampipe-plugin-imap/pull/40))
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-aiven/blob/main/docs/LICENSE). ([#40](https://github.com/turbot/steampipe-plugin-imap/pull/40))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server encapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#39](https://github.com/turbot/steampipe-plugin-imap/pull/39))

## v0.5.1 [2023-10-05]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#32](https://github.com/turbot/steampipe-plugin-imap/pull/32))

## v0.5.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#30](https://github.com/turbot/steampipe-plugin-imap/pull/30))
- Recompiled plugin with Go version `1.21`. ([#30](https://github.com/turbot/steampipe-plugin-imap/pull/30))

## v0.4.0 [2023-03-23]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#21](https://github.com/turbot/steampipe-plugin-imap/pull/21))

## v0.3.0 [2022-11-10]

_Enhancements_

- Added environment variables for setting the `host`, `login`, `password`, and `port` config arguments. ([#18](https://github.com/turbot/steampipe-plugin-imap/pull/18)) (Thanks to [@graza-io](https://github.com/graza-io) for the additions!)

## v0.2.0 [2022-09-27]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#15](https://github.com/turbot/steampipe-plugin-imap/pull/15))
- Recompiled plugin with Go version `1.19`. ([#15](https://github.com/turbot/steampipe-plugin-imap/pull/15))

_Bug fixes_

- Fixed the `message_id` column of `imap_message` table to return `null` if it is invalid UTF-8. ([#5](https://github.com/turbot/steampipe-plugin-imap/pull/5))

## v0.1.1 [2022-05-23]

_Bug fixes_

- Fixed the Slack community links in README and docs/index.md files. ([#10](https://github.com/turbot/steampipe-plugin-imap/pull/10))

## v0.1.0 [2022-04-27]

_Enhancements_

- Added support for native Linux ARM and Mac M1 builds. ([#8](https://github.com/turbot/steampipe-plugin-imap/pull/8))
- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#7](https://github.com/turbot/steampipe-plugin-imap/pull/7))

## v0.0.1 [2022-02-16]

_What's new?_

- New tables added
  - [imap_mailbox](https://hub.steampipe.io/plugins/turbot/imap/tables/imap_mailbox)
  - [imap_message](https://hub.steampipe.io/plugins/turbot/imap/tables/imap_message)
