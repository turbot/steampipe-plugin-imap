---
title: "Steampipe Table: imap_message - Query IMAP Messages using SQL"
description: "Allows users to query IMAP Messages, providing detailed information about each message available in the user's mailbox."
---

# Table: imap_message - Query IMAP Messages using SQL

Internet Message Access Protocol (IMAP) is an Internet standard protocol used by email clients to retrieve messages from a mail server over a TCP/IP connection. It is a method of accessing electronic mail or bulletin board messages that are kept on a (possibly shared) mail server. IMAP allows for managing and manipulating the email on the server without downloading the messages to the local device.

## Table Usage Guide

The `imap_message` table provides insights into messages within an IMAP server. As a system administrator, explore message-specific details through this table, including sender, recipient, subject, and associated metadata. Utilize it to uncover information about messages, such as those with specific subjects, from certain senders, or sent at particular dates.

## Examples

### List messages from the default mailbox (e.g. INBOX)
Explore the contents of your default mailbox, such as the INBOX, to gain insights into all your email messages. This can be particularly useful for analyzing overall email activity or identifying specific messages.

```sql
select
  *
from
  imap_message
```

### List messages from a specific mailbox
Determine the areas in which starred emails reside within a specific Gmail mailbox. This can be useful to quickly pinpoint important communications that have been marked for follow-up.

```sql
select
  *
from
  imap_message
where
  mailbox = '[Gmail]/Starred'
```

### Find messages greater than 1MB in size
Explore which emails have a large size, potentially indicating attachments or extensive content. This can help manage storage space and identify important communications that may require more attention due to their size.

```sql
select
  timestamp,
  from_email,
  subject,
  size
from
  imap_message
where
  size > 1000000
order by
  size desc
```

### List messages received between 7 and 14 days ago
Discover the segments of your email inbox that contain messages received between a week and two weeks ago. This can help you manage your inbox by focusing on specific time frames.

```sql
select
  timestamp,
  from_email,
  subject
from
  imap_message
where
  timestamp > current_timestamp - interval '14 days'
  and timestamp < current_timestamp - interval '7 days'
order by
  timestamp
```

### Find messages from a given address
Discover the segments that contain emails sent from a specific address. This is useful for tracking communication history and identifying the content of the messages from a particular sender.

```sql
select
  timestamp,
  from_email,
  subject
from
  imap_message
where
  from_email = 'jim@dundermifflin.com'
```

### Search drafts for messages with a keyword
Explore drafts for messages containing a specific keyword. This can help in quickly identifying and reviewing relevant draft messages without having to manually search through each one.

```sql
select
  timestamp,
  from_email,
  subject
from
  imap_message
where
  mailbox = '[Gmail]/Drafts'
  and query = 'keyword'
```

### List all attachments on Starred messages
Explore the attachments linked to your most important emails. This query is useful for identifying and reviewing all attachments connected to your starred messages, providing a quick way to assess important documents or files.

```sql
select
  m.timestamp,
  m.from_email,
  m.subject,
  a ->> 'file_name' as attachment_filename,
  a ->> 'content_type' as attachment_content_type
from
  imap_message as m,
  jsonb_array_elements(attachments) as a
where
  mailbox = '[Gmail]/Starred'
```