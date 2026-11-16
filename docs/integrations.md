# Integrations record

One row per app that taf jobs connect through the connection platform: how its sign-in is set up, what it asks the vendor for, what its
jobs need at least, and whether a real sign-in has proved it. A change to an
app's sign-in in the platform's console lands in the same pull request as
its row here.

No logins, passwords, tokens, client ids or config ids go in this file.

- **Sign-in:** *managed* is the platform's own OAuth client. *Ours* is a
  client registered under KIFF's own developer account and entered in the
  platform's console.
- **Asked today:** the scopes the sign-in requests, read in the console.
  planque narrows a managed Google sign-in per taf (`Scope`, see
  [managed connections](managed-connectors.md)) to the narrowest scope its
  ceiling offers. Where a grant is `usual_access`, every taf on that app gets
  the whole list.
- **Narrowest:** the least the app's current jobs need, from the guide's read
  and write tags.
- **Proof:** a sign-in on a phone, one call, nothing kept: the taf removed and
  its grant disconnected afterwards.

| App | Sign-in | Asked today | Narrowest for its jobs | Catalog | Guide | Proof | Left |
|---|---|---|---|---|---|---|---|
| Shopify | ours (a development-dashboard app, legacy install flow) | `read_orders`, `write_orders`, `read_customers` (16 Nov; before: 36 scopes, among them `read_webhooks`, which Shopify rejects) | `read_orders`, `read_customers`; `write_orders` only for `SHOPIFY_CREATE_REFUND` | `shopify`, 386 tools (199 read, 187 write) | [shopify](../deploy/guides/shopify.md) | not proved: the 16 Nov sign-in stopped at `invalid_scope: read_webhooks`, before the scopes were fixed | the read proof on a development store; then a refund held for approval, on the store's test gateway only, before the beta; the sign-in isn't narrowed, so every Shopify taf holds `write_orders`; two sign-ins (read-only, and with refunds) need planque to choose between them |
| Typeform | managed | 14: `accounts:read`, `forms:read`, `forms:write`, `images:read`, `images:write`, `themes:read`, `themes:write`, `responses:read`, `responses:write`, `webhooks:read`, `webhooks:write`, `workspaces:read`, `workspaces:write`, `offline`; the grant is `usual_access` | `forms:read`, `responses:read` | `typeform`, 35 tools (16 read, 19 write) | [typeform](../deploy/guides/typeform.md) | 16 Nov, a US data-centre account: `list_forms`, 1 row, nothing written | an EU data-centre account (another API address); narrow the sign-in to the two read scopes |
| Google Sheets | managed, narrowed per taf | `spreadsheets` (the usual sign-in: `drive`, `spreadsheets`, `userinfo.email`) | `spreadsheets.readonly` to read; `spreadsheets` to add rows | `googlesheets`, 50 tools (17 read, 33 write) | [googlesheets](../deploy/guides/googlesheets.md) | 16 Nov, a test account: `read_rows`, 31 rows, nothing written | the managed ceiling has no `spreadsheets.readonly` (asked of the platform, 16 Nov), so a read-only taf's sign-in is wider and the card says so |
| Google Drive | managed, narrowed per taf | `drive` (the usual sign-in: `drive`, `userinfo.email`) | `drive.readonly` to find files; `drive.file` to save receipts | `googledrive`, 92 tools (41 read, 51 write) | [googledrive](../deploy/guides/googledrive.md) | 16 Nov, a test account: `find_files`, 0 files, nothing written | the managed ceiling has no `drive.readonly` (asked of the platform, 16 Nov) |
| Google Calendar | managed, narrowed per taf | `calendar.events` | `calendar.events.readonly` to read; `calendar.events` to book | `googlecalendar`, 47 tools (19 read, 28 write) | [googlecalendar](../deploy/guides/googlecalendar.md) | 16 Nov, a test account: `list_week_events`, 0 events, nothing written | the managed ceiling has no read-only Calendar scope (asked of the platform, 16 Nov) |
| Gmail | managed; planque narrows per taf | the usual sign-in: `https://mail.google.com/`, `userinfo.email`, `userinfo.profile` (16 Nov; before: also contacts, other contacts, birthday, addresses, phone numbers, emails and language, which no job reads); planque-narrowed: `https://mail.google.com/` | `gmail.readonly` to read; `gmail.send` to send | `gmail`, 61 tools (28 read, 33 write) | [gmail](../deploy/guides/gmail.md) | not proved for a job yet | the managed ceiling has no `gmail.readonly` (asked of the platform) |
| Outlook | managed | `Mail.Read`, `Mail.Send`, `Calendars.ReadWrite`, `User.Read`, `offline_access` (16 Nov; before: 12, among them `Mail.ReadWrite`, `Chat.ReadWrite`, `Contacts.ReadWrite`, `MailboxSettings.ReadWrite`, `Calendars.ReadWrite.Shared`) | the same five | `outlook` | [outlook](../deploy/guides/outlook.md) | not proved yet | a phone proof |

Rows for the other apps in the platform's console follow from the sign-in
audit.
