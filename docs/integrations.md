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
  and write tags. Since 16 Nov a usual sign-in keeps the app's full managed
  access, and a taf's Card decides what it may do; planque narrows a sign-in
  only where the provider's rules or a consent screen require it.
- **Proof:** a sign-in on a phone, one call, nothing kept: the taf removed and
  its grant disconnected afterwards.

| App | Sign-in | Asked today | Narrowest for its jobs | Catalog | Guide | Proof | Left |
|---|---|---|---|---|---|---|---|
| Shopify | ours (a development-dashboard app, legacy install flow) | 35: the earlier 36 minus `read_webhooks`, which Shopify rejects (restored 16 Nov to the full managed access: the Card decides what a taf may do) | `read_orders`, `read_customers`; `write_orders` only for `SHOPIFY_CREATE_REFUND` | `shopify`, 386 tools (199 read, 187 write) | [shopify](../deploy/guides/shopify.md) | not proved: the 16 Nov sign-in stopped at `invalid_scope: read_webhooks`, before the scopes were fixed | the read proof on a development store; then a refund held for approval, on the store's test gateway only, before the beta; the sign-in isn't narrowed, so every Shopify taf holds `write_orders`; two sign-ins (read-only, and with refunds) need planque to choose between them |
| Typeform | managed | 14: `accounts:read`, forms, images, themes, responses, webhooks and workspaces read and write, `offline` (restored 16 Nov to the full managed access: the Card decides what a taf may do); the grant is `usual_access` | `forms:read`, `responses:read` | `typeform`, 35 tools (16 read, 19 write) | [typeform](../deploy/guides/typeform.md) | 16 Nov, a US data-centre account: `list_forms`, 1 row, nothing written | an EU data-centre account (another API address); narrow the sign-in to the two read scopes |
| Google Sheets | managed, narrowed per taf | `spreadsheets` (the usual sign-in: `drive`, `spreadsheets`, `userinfo.email`, restored 16 Nov to the full managed access: the Card decides what a taf may do) | `spreadsheets.readonly` to read; `spreadsheets` to add rows | `googlesheets`, 50 tools (17 read, 33 write) | [googlesheets](../deploy/guides/googlesheets.md) | 16 Nov, a test account: `read_rows`, 31 rows, nothing written | the managed ceiling has no `spreadsheets.readonly` (asked of the platform, 16 Nov), so a read-only taf's sign-in is wider and the card says so |
| Google Drive | managed, narrowed per taf | `drive` (the usual sign-in: `drive`, `userinfo.email`) | `drive.readonly` to find files; `drive.file` to save receipts | `googledrive`, 92 tools (41 read, 51 write) | [googledrive](../deploy/guides/googledrive.md) | 16 Nov, a test account: `find_files`, 0 files, nothing written | the managed ceiling has no `drive.readonly` (asked of the platform, 16 Nov) |
| Google Calendar | managed, narrowed per taf | `calendar.events` (the usual sign-in: `calendar`, `calendar.events`, restored 16 Nov to the full managed access: the Card decides what a taf may do) | `calendar.events.readonly` to read; `calendar.events` to book | `googlecalendar`, 47 tools (19 read, 28 write) | [googlecalendar](../deploy/guides/googlecalendar.md) | 16 Nov, a test account: `list_week_events`, 0 events, nothing written | the managed ceiling has no read-only Calendar scope (asked of the platform, 16 Nov) |
| Gmail | managed; planque narrows per taf | the usual sign-in: `https://mail.google.com/` plus 10 profile and contacts scopes (restored 16 Nov to the full managed access: the Card decides what a taf may do); planque-narrowed: `https://mail.google.com/` | `gmail.readonly` to read; `gmail.send` to send | `gmail`, 61 tools (28 read, 33 write) | [gmail](../deploy/guides/gmail.md) | not proved for a job yet | the managed ceiling has no `gmail.readonly` (asked of the platform) |
| Outlook | managed | 12: mail, calendars, contacts, chat and mailbox settings, read and write, `User.Read`, `offline_access` (restored 16 Nov to the full managed access: the Card decides what a taf may do) | the same five | `outlook` | [outlook](../deploy/guides/outlook.md) | not proved yet | a phone proof |
| Intercom | managed | 31: read and write on users, companies, conversations, admins, tags, events, segments, counts, data attributes, tickets, teams, notes, news, articles and webhooks, `offline_access` (restored 16 Nov to the full managed access: the Card decides what a taf may do) | the same | `intercom`, 138 tools (68 read, 70 write) | [intercom](../deploy/guides/intercom.md) | not proved yet | a phone proof |
| Gorgias | managed | 31: `openid`, `email`, `profile`, `offline`, and read and write on account, users, customers, tickets, custom fields, events, integrations, jobs, macros, rules, satisfaction surveys, tags and apps, `statistics:read` (restored 16 Nov to the full managed access: the Card decides what a taf may do) | the same | `gorgias`, 51 tools (25 read, 26 write) | [gorgias](../deploy/guides/gorgias.md) | not proved yet | a phone proof |
| Zendesk | managed | 10: `read`, `write`, and organizations, tickets, users and webhooks read and write (restored 16 Nov to the full managed access: the Card decides what a taf may do) | the same | `zendesk`, 480 tools (264 read, 216 write) | [zendesk](../deploy/guides/zendesk.md) | not proved yet | a phone proof |
| QuickBooks | managed | 6: `com.intuit.quickbooks.accounting`, `openid`, `profile`, `email`, `phone`, `address` (restored 16 Nov to the full managed access: the Card decides what a taf may do) | the same | `quickbooks`, 114 tools (76 read, 38 write) | [quickbooks](../deploy/guides/quickbooks.md) | not possible: no account can be opened | — |
| Zoho Books | managed | `ZohoBooks.fullaccess.all`, the only scope the managed sign-in offers | read-only invoices and contacts | `zoho_books`, 270 tools (103 read, 167 write) | [zoho_books](../deploy/guides/zoho_books.md) | not proved: no account yet | read-only scopes (to be asked of the platform) |

Rows for the other apps in the platform's console follow from the sign-in
audit.
