## v0.5.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters.
- Recompiled plugin with Go version `1.21`.

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
