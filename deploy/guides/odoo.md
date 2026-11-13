# Odoo

Checked: 2026-11-13 · Sources: the connection service's Odoo tool
definitions as planque reads them (12 tools), Odoo's external API reference
and document Extract API reference (Odoo 19 documentation).

## How it connects

- **A key the owner pastes**, an API key made in Odoo under Preferences,
  Account Security, held by the connection service; the taf never sees a
  password. The key acts as the Odoo user who made it, with that user's
  access rights and record rules, no more.
- **Keys have a set lifetime.** The owner picks a duration when making the
  key, and Odoo allows at most three months. A job that runs for longer needs
  a new key before the old one ends. Suggest a dedicated Odoo user with only
  the rights the job needs.
- **The document-reading tools need a second token**, the account token of
  the owner's Odoo credit account for document digitization. It is not the
  API key. Ask for it only when the job reads invoices, receipts, bank
  statements or CVs.
- **Ask for every tool the job needs at once,** so nothing is asked twice.

## Tools a job usually needs

Read means it changes nothing in Odoo and spends nothing.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `ODOO_LIST_DATABASES` | Lists the databases on the owner's Odoo server | read | database names |
| `ODOO_CREATE_CONTACT` | Adds a contact (a person or a company) | write | the new contact's id |
| `ODOO_UPDATE_RECORD` | Changes fields on records of one model | write | true when saved |
| `ODOO_PARSE_INVOICE` | Sends one vendor bill or invoice (PDF, PNG or JPG) to be read | write, spends a credit | a document token |
| `ODOO_GET_INVOICE_RESULT` | Fetches what was read from that bill: supplier, dates, totals, lines | read | the status and the read fields |
| `ODOO_PARSE_EXPENSE` / `ODOO_GET_EXPENSE_RESULT` | The same for a receipt | write, spends a credit / read | a document token / the read fields |
| `ODOO_PARSE_BANK_STATEMENT` / `ODOO_GET_BANK_STATEMENT_RESULT` | The same for a bank statement | write, spends a credit / read | a document token / the statement lines |
| `ODOO_PARSE_APPLICANT` / `ODOO_GET_APPLICANT_RESULT` | The same for a CV | write, spends a credit / read | a document token / the applicant's fields |

**What to keep:** a new contact's id, to update that same contact later
instead of searching again; a parse's document token and the version number
sent with it, since the result must be asked for with the same version.
Results take a little while: ask again after a pause, a few times at most.

**No tool here reads records.** There is no search or read tool for
contacts, invoices or orders, short of the raw call below. Don't offer a job
that must first look something up in Odoo; say it can't be done with this
connection yet.

**Start read-only.** "Read this bill and tell me what's on it" is a parse
tool plus its result tool; it writes nothing in Odoo, but each parse costs a
credit, so say so. Offer `ODOO_CREATE_CONTACT` or `ODOO_UPDATE_RECORD` only
when the owner asked for changes, and show each change first.
`ODOO_UPDATE_RECORD` can change any model: offer it only for the one model
the job names (contacts, say), one record at a time.

## Don't offer

- `ODOO_CALL_ODOO_JSONRPC`: a raw call that can read, change or delete
  anything the key's user can. It is the bulk and raw-request tool.
- `ODOO_UPDATE_RECORD` on users, companies, access rights, settings or
  accounting entries: it changes who can do what, or the books.

## Access and plans

- **Odoo Online needs the Custom plan** for the external API; One App Free
  and Standard don't include it. A self-hosted server depends on
  its own setup.
- **Every call is checked** against the key user's rights. A user who can't
  see invoices in Odoo can't through the taf either.
- **Document reading uses prepaid credits**, one per document read
  successfully. With no credit left, parsing stops until the owner buys more.
  Files must be PDF, PNG or JPG; password-protected PDFs, very small files
  and documents with too many pages are refused.
- **The older remote-call endpoints are going away:** Odoo plans to remove
  them in Odoo Online 21.1 (winter 2027) and Odoo 22 (fall 2028). Check this
  guide again before then.
- **Rate:** Odoo's reference states no general call limit.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused, key not valid | The key was revoked, mistyped, or its duration ended | "Your Odoo key no longer works. Make a new one in Odoo under Preferences, Account Security, and paste it here. Nothing was changed." |
| Refused, access error | The key's user lacks the right for this record or model | "Your Odoo user isn't allowed to do this. An Odoo admin can give that user the right, or the taf can skip this step." |
| External API not available | The Odoo Online plan doesn't include it | "Your Odoo plan doesn't allow outside connections. Odoo's Custom plan does; until then this job can't run." |
| Parse answers `error_no_credit` | The document-reading credit is used up | "Odoo's document reading has no credit left. Add credit in Odoo and the taf will carry on." |
| Parse answers `error_unsupported_format`, or a result answers `error_password_protected` or `error_too_many_pages` | The file can't be read as sent | "Odoo couldn't read that file. Send it as a plain PDF or a photo, without a password." |
| Result answers `error_document_not_found` | The token is wrong or the version differs from the parse | Ask once more with the kept token and version; else "Odoo lost that document. Send it again?" |
| Parse or result answers `error_maintenance` | Odoo's reading service is down for a while | "Odoo's document reading is paused for maintenance. The taf will try again later." |
| A record isn't found | Deleted, or hidden from the key's user | "That record isn't in Odoo any more, or you can't see it. The taf will skip it." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
