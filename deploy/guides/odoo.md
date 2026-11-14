# Odoo

Checked: 2026-11-13 · Sources: the connection service's Odoo tool
definitions as planque reads them (12 tools), and Odoo's External API and
Extract API documentation (version 18.0).

## How it connects

- **By the owner's API key**, pasted once and held by the connection
  service; the taf never sees it. An Odoo user makes a key in their own
  preferences, under account security. Odoo shows it only once, so the
  owner copies it straight into the connection.
- **The key acts as that Odoo user.** What a taf can read or change is what
  that user's access rights allow, no more.
- **Odoo Online plans:** Odoo says its external API is open only on its
  Custom plan, not on the free one-app plan or Standard. On another plan,
  contact and record tools are refused.
- **The document-reading tools use a separate account.** Reading invoices,
  expenses, bank statements and CVs runs on Odoo's paid digitization
  service, charged in credits to an in-app purchase account. Those tools
  need that account's token, not the API key.
- **A held connection is reused** for a new taf when it covers the tools
  asked. Ask for every tool the job needs in the first request.

## Tools a job usually needs

Read means it changes nothing in Odoo. The catalog gives no `returns` paths
for Odoo today; the cells below are what Odoo's documentation shows.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `ODOO_LIST_DATABASES` | Lists the databases on the Odoo server | read | not documented |
| `ODOO_CREATE_CONTACT` | Adds a person or company to the contacts | write | the new contact's id |
| `ODOO_UPDATE_RECORD` | Changes fields on given records of one kind | write | not documented |
| `ODOO_PARSE_INVOICE` | Sends one vendor invoice (PDF or image) to be read | write, costs a credit | a document token |
| `ODOO_GET_INVOICE_RESULT` | Fetches what was read from that invoice | read | total, dates, supplier, invoice lines |
| `ODOO_PARSE_EXPENSE` | Sends one receipt to be read | write, costs a credit | a document token |
| `ODOO_GET_EXPENSE_RESULT` | Fetches what was read from that receipt | read | total, date, currency |
| `ODOO_PARSE_BANK_STATEMENT` | Sends one bank statement to be read | write, costs a credit | a document token |
| `ODOO_GET_BANK_STATEMENT_RESULT` | Fetches what was read from that statement | read | balances, statement lines |
| `ODOO_PARSE_APPLICANT` | Sends one CV to be read | write, costs a credit | a document token |
| `ODOO_GET_APPLICANT_RESULT` | Fetches what was read from that CV | read | name, email, phone |

**What to keep:** each document token, to fetch its result on a later run
instead of sending the document again (each send costs a credit); a new
contact's id, to change that same contact later. Never look a contact up by
name to change it: use the id the taf kept.

**Reading a document takes two steps.** The taf sends it, keeps the token,
and checks back for the result. Odoo reads one document per send, as PDF,
PNG or JPG.

**Start read-only.** "Read this week's vendor invoices and list what's due"
needs the parse and result tools for invoices only. Changing records is
offered only when the owner asked, and each change is shown first and
waits for their yes.

## Don't offer

- **The raw call tool** (`ODOO_CALL_ODOO_JSONRPC`): it can run any method on
  any kind of record, deleting included, so a Card can't say what it will
  do.
- **`ODOO_UPDATE_RECORD` on accounting records** (posted invoices, payments,
  journal entries): it changes the books. Offer it only for the kind of
  record the owner named, such as contacts.
- **CV reading for anything but a recruiting job:** a CV holds personal
  details.

## Access and plans

- **Credits:** each successful document read costs one credit from the
  digitization account. With no credit left, sends are refused until the
  owner buys more in Odoo.
- **Rights:** a key from a user without access to an app (Accounting, say)
  can't read or change that app's records.
- **Online accounts:** a user with no local password can't use the
  external API until an administrator sets one.
- **Rate limits:** not given in the pages checked here.

## When it fails

The document-reading service answers most problems inside a normal answer,
with a status word, rather than as an error.

| What happens | Why | Tell the owner |
|---|---|---|
| Refused, not enough credit (`error_no_credit`) | The digitization account has run out of credits | "Odoo has no document-reading credits left. Buy more in Odoo, then the taf can carry on." |
| Unsupported file (`error_unsupported_format`) | The file isn't a PDF, PNG or JPG | "Odoo can't read that file type. Send it as a PDF or a photo." |
| Protected or too long (`error_password_protected`, `error_too_many_pages`) | The PDF has a password, or too many pages | "Odoo couldn't read that document: it's locked or too long. Send an unlocked copy, or fewer pages." |
| Result not found (`error_document_not_found`) | The token is wrong or the document is gone | "Odoo can't find that document any more. The taf will skip it." |
| Under maintenance (`error_maintenance`) | Odoo's reading service is down for a while | "Odoo's document reading is paused for maintenance. The taf will try again later." |
| Refused, not allowed | The key's user lacks rights for that kind of record, or the plan has no external API | "Odoo didn't let the taf do that with this login. Check the user's rights, or your Odoo plan." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write, and a failed read is never sent again without checking
for its result first.
