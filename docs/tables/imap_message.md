# Table: imap_message

Read messages from a given mailbox.

Note: A `mailbox` must be provided in all queries to this table.

## Examples

### List messages from the Inbox

```sql
select
  *
from
  imap_message
where
  mailbox = 'INBOX'
```

### Find messages greater than 1MB in size

```sql
select
  timestamp,
  from_email,
  subject,
  size
from
  imap_message
where
  mailbox = 'INBOX'
  and size > 1000000
order by
  size desc
```

### Find messages from a given address

```sql
select
  timestamp,
  from_email,
  subject
from
  imap_message
where
  mailbox = 'INBOX'
  and from_email = 'jim@dundermifflin.com'
```

### Search drafts for messages with a keyword

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

```sql
select
  m.timestamp,
  m.from_email,
  m.subject,
  a ->> 'Filename' as attachment_filename,
  a ->> 'ContentType' as attachment_content_type
from
  imap_message as m,
  jsonb_array_elements(attachments) as a
where
  mailbox = '[Gmail]/Starred'
```
