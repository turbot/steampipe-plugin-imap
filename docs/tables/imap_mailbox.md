# Table: imap_mailbox

Mailboxes accessible on the IMAP server.

## Examples

### List all mailboxs

```sql
select
  *
from
  imap_mailbox
```

### Get a specific mailbox

```sql
select
  *
from
  imap_mailbox
where
  name = '[Gmail]/Starred'
```

### Mailboxes by message count

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

```sql
select
  *
from
  imap_mailbox
where
  attributes ? '\Important'
```
