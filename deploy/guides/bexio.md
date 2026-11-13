# bexio

Checked: 2026-11-13 · Sources: the bexio tool definitions as planque reads
them (309 tools, generated/bexio.md), and bexio's API reference at
docs.bexio.com (personal access tokens, rate limiting, errors, pagination,
changelog).

## How it connects

- **A key the owner pastes**: a personal access token the owner makes in
  bexio's developer portal (developer.bexio.com/pat), held by the
  connection service. The taf never sees the key after it is pasted.
- **A key opens the whole company.** bexio gives a personal access token
  every default access, so the catalog's `unusable` list rarely narrows
  anything. The Card and the grant are what keep a taf to its tools: ask for
  only the tools the job needs, all in the first request.
- **The key lasts 60 days.** bexio stops accepting it 60 days after it was
  made; the owner makes a new one and pastes it again. A key the owner
  deletes can keep working for up to an hour.

## Tools a job usually needs

Read means it changes nothing in bexio. The generated part gives no
returns; the ids below are the ones bexio's API reference shows.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `BEXIO_LIST_CONTACTS` / `BEXIO_SEARCH_CONTACT` | Lists contacts, or finds them by name, email or number | read | contact ids, contact numbers |
| `BEXIO_GET_CONTACT` | Reads one contact | read | the contact id and its details |
| `BEXIO_LIST_INVOICES` / `BEXIO_SEARCH_INVOICES` | Lists invoices, or finds them by status, contact or date | read | invoice ids, document numbers |
| `BEXIO_GET_INVOICE` | Reads one invoice with its lines | read | the invoice id, document number, contact id, total |
| `BEXIO_GET_INVOICE_PDF` | Gets an invoice as a PDF | read | the file's name and content |
| `BEXIO_LIST_INVOICE_PAYMENTS` | Lists the payments booked on one invoice | read | payment ids, dates, amounts |
| `BEXIO_LIST_QUOTES` / `BEXIO_GET_QUOTE` | Reads quotes | read | quote ids, document numbers |
| `BEXIO_LIST_ITEMS` / `BEXIO_LIST_TAXES` | Reads products and tax rates, needed to write invoice lines | read | item ids, tax ids |
| `BEXIO_CREATE_CONTACT` | Adds a contact | write | the new contact id |
| `BEXIO_UPDATE_CONTACT` | Changes a contact's details | write | the contact id |
| `BEXIO_CREATE_QUOTE` | Writes a quote as a draft | write | the quote id, its document number |
| `BEXIO_CREATE_INVOICE` | Writes an invoice as a draft; sends nothing | write | the invoice id, its document number |
| `BEXIO_ISSUE_INVOICE` | Turns a draft invoice into an open one | write | whether it worked |
| `BEXIO_SEND_INVOICE` | Emails an invoice to the customer | write, outside | whether it was sent |

**What to keep:** a contact's id, and an invoice's or quote's id with its
document number. The id is what the next call needs; the document number is
what the owner and the customer see. Never search for a document the taf
made: use the id it kept.

**Start read-only.** "Draft the invoice and let me check it" is the read
tools plus `BEXIO_CREATE_INVOICE`. Issuing and sending are separate steps:
offer `BEXIO_ISSUE_INVOICE` or `BEXIO_SEND_INVOICE` only when the owner asked
for it, and say each one waits for their yes.

## Don't offer

- Deleting tools (`BEXIO_DELETE_INVOICE`, `BEXIO_DELETE_CONTACT`,
  `BEXIO_DELETE_PAYMENT` and the other `BEXIO_DELETE_*`), and cancelling or
  undoing tools (`BEXIO_CANCEL_INVOICE`, `BEXIO_CANCEL_PAYMENT`,
  `BEXIO_REVERT_ISSUE_INVOICE`). They change the books.
- Money and accounting: `BEXIO_CREATE_OUTGOING_PAYMENT`,
  `BEXIO_CREATE_PAYMENT`, `BEXIO_CREATE_MANUAL_ENTRY`,
  `BEXIO_CREATE_CALENDAR_YEAR`, `BEXIO_UPDATE_BILL_BOOKING_STATUS`,
  `BEXIO_UPDATE_EXPENSE_BOOKING_STATUS`.
- People and payroll: `BEXIO_CREATE_EMPLOYEE`, `BEXIO_UPDATE_EMPLOYEE`,
  `BEXIO_CREATE_FICTIONAL_USER`, `BEXIO_DOWNLOAD_PAYSTUB_PDF`,
  `BEXIO_GET_PERMISSIONS`.
- Settings the whole company shares: currencies, taxes, countries, units,
  titles and salutations (`BEXIO_CREATE_CURRENCY`, `BEXIO_UPDATE_COUNTRY`,
  `BEXIO_DELETE_TAX` and the like).
- Bulk writes: `BEXIO_CREATE_CONTACTS_BULK`.

## Access and plans

- **Purchase tools** (bills, expenses, outgoing payments) work only for
  companies that use bexio's newer purchase module. Elsewhere they fail; read
  the owner's setup before offering them.
- **Big lists:** a list answers 500 items by default and at most 2,000 a
  page; bills at most 500. Page through with an offset.
- **Rate:** bexio limits calls per period and answers "too many requests"
  past it, with a header saying how many seconds until the next period. Its
  docs don't publish the number, so a job that reads thousands of records
  pages and pauses.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The key isn't accepted (401) | The key expired after 60 days, or the owner deleted it | "bexio stopped accepting the key. Make a new one in bexio's developer portal and paste it here. Nothing was changed." |
| Refused (403) | The key's user can't use this part of bexio, or the company doesn't have the module | "bexio didn't let us do this with your account. Check that your user can use this part of bexio, or do without this step." |
| Not found (404) | The contact, invoice or quote was deleted, or the id is wrong | "That record isn't in bexio any more. The taf will skip it." |
| Not saved (422) | A required field is missing or not valid, such as a contact, tax or account that doesn't exist | Fix the input and try once; else "bexio didn't accept that document. Nothing was saved." |
| Too many requests (429) | Over bexio's rate limit | "bexio asked us to slow down. The taf will try again shortly." |
| Unavailable (503) | bexio is down or in maintenance | "bexio is not available right now. The taf will try again later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
