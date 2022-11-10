connection "imap" {
  plugin = "imap"

  # Required: Hostname of the IMAP server, can also be set with the IMAP_HOST environment variable.
  # host = "imap.gmail.com"

  # Required: Login name, usually the email address, can also be set with the IMAP_LOGIN environment variable.
  # login = "michael@dundermifflin.com"

  # Required: Password, can also be set with the IMAP_PASSWORD environment variable.
  # password = "Great Scott!"

  # Optional: Port (default: 993), can also be set with the IMAP_PORT environment variable.
  # port = 993

  # Example Gmail configuration
  # host = "imap.gmail.com"
  # port = 993
  # tls_enabled = true
  # login = "michael@gmail.com"
  # password = "Password1234"
}
