# Table: imap_message

Query messages from a given mailbox.

All queries are against a single mailbox, chosen in this order of precedence:
1. A `where mailbox = 'INBOX'` qualifier in the query.
2. The `mailbox` config setting in `imap.spc`.
3. Default is `INBOX`.

## Examples

### List messages from the default mailbox (e.g. INBOX)

```sql
select
  *
from
  imap_message
```

### List messages from a specific mailbox

```sql
select
  *
from
  imap_message
where
  mailbox = '[Gmail]/Starred'
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
  size > 1000000
order by
  size desc
```

### List messages received between 7 and 14 days ago

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
  a ->> 'file_name' as attachment_filename,
  a ->> 'content_type' as attachment_content_type
from
  imap_message as m,
  jsonb_array_elements(attachments) as a
where
  mailbox = '[Gmail]/Starred'
```
