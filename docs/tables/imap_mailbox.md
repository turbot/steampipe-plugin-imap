---
title: "Steampipe Table: imap_mailbox - Query IMAP Mailboxes using SQL"
description: "Allows users to query IMAP Mailboxes, specifically the mailbox details, providing insights into email organization and potential anomalies."
---

# Table: imap_mailbox - Query IMAP Mailboxes using SQL

IMAP (Internet Message Access Protocol) is a standard email protocol that stores email messages on a mail server and allows the end user to view and manipulate the messages as though they were stored locally on the end user's computing device(s). The IMAP protocol is designed to enable users to view, organize and sort their email messages in various ways. It is particularly useful for people who need to access their email on multiple devices or from different locations.

## Table Usage Guide

The `imap_mailbox` table provides insights into mailboxes within an IMAP server. As a system administrator, explore mailbox-specific details through this table, including mailbox size, message count, and associated metadata. Utilize it to uncover information about mailboxes, such as those with high message count, the size of each mailbox, and the organization of emails.

## Examples

### List all mailboxs
Discover the segments that encompass all your mailboxes, providing a comprehensive view of your email landscape. This is useful for gaining an overall understanding of your email organization and management.

```sql
select
  *
from
  imap_mailbox
```

### Get a specific mailbox
Explore the specific mailbox that has been marked as important or noteworthy in your email system. This is useful to quickly identify and access important communications without having to manually search through all your emails.

```sql
select
  *
from
  imap_mailbox
where
  name = '[Gmail]/Starred'
```

### Mailboxes by message count
Explore which mailboxes contain the highest number of messages to better manage storage and prioritize clean-up efforts.

```sql
select
  name,
  messages
from
  imap_mailbox
order by
  messages desc
```

### Mailboxes with unseen messages
Identify mailboxes that contain unread messages to prioritize checking and responding to these communications.

```sql
select
  name,
  unseen
from
  imap_mailbox
where
  unseen > 0
order by
  unseen desc
```

### Get all mailboxes with the Important attribute
Explore which mailboxes have been marked as important. This is useful for prioritizing email management and focusing on high-priority communications.

```sql
select
  *
from
  imap_mailbox
where
  attributes ? '\Important'
```