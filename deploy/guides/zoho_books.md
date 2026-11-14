# Zoho Books

Checked: 2026-11-14 · Sources: the connection service's Zoho Books tool
definitions as planque reads them (270 tools), and Zoho Books' API
documentation (introduction, sign-in, errors, paging, invoices, contacts,
customer payments, estimates, items).

## How it connects

- **A sign-in to Zoho**, held by the connection service; the taf never sees
  a password or token. planque's catalog lists the tools and, with a known
  sign-in, puts the ones its access doesn't cover under `unusable`, with why.
- **Zoho keeps each account in one region** (Europe, the United States,
  India and others). An account held in Europe signs in, and is reached,
  only through Europe's addresses; another region's address doesn't find
  it.
- **Every call names one organization.** A person can belong to several, so
  the taf reads the list first and keeps the one the owner means.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.

## Tools a job usually needs

Read means it changes nothing in the books. The catalog's `returns` gives
no paths for Zoho Books today; the ids below are those Zoho's answers carry.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `ZOHO_BOOKS_LIST_ORGANIZATIONS` | Lists the organizations the person belongs to | read | organization ids |
| `ZOHO_BOOKS_LIST_INVOICES` | Lists invoices, filtered by customer, status or dates, a page at a time | read | invoice ids, numbers, balances, due dates, whether more pages follow |
| `ZOHO_BOOKS_GET_INVOICE` | Reads one invoice with its lines | read | the invoice id, number, status, balance and link |
| `ZOHO_BOOKS_LIST_INVOICE_PAYMENTS` | Lists payments recorded on one invoice | read | payment ids |
| `ZOHO_BOOKS_LIST_CUSTOMER_PAYMENTS` | Lists payments received | read | payment ids and numbers, the invoices they cover |
| `ZOHO_BOOKS_LIST_CONTACTS` | Lists customers and vendors, filtered by name or email | read | contact ids, names, emails |
| `ZOHO_BOOKS_GET_CONTACT` | Reads one customer or vendor | read | the contact id |
| `ZOHO_BOOKS_LIST_ESTIMATES` / `ZOHO_BOOKS_GET_ESTIMATE` | Lists or reads quotes | read | estimate ids and numbers |
| `ZOHO_BOOKS_LIST_ITEMS` | Lists the products and services sold | read | item ids |
| `ZOHO_BOOKS_CREATE_CONTACT` | Adds a customer or vendor | write | the new contact id |
| `ZOHO_BOOKS_UPDATE_INVOICE` | Changes an invoice's lines, dates or customer | write | the invoice id, number and status |
| `ZOHO_BOOKS_ADD_INVOICE_COMMENT` | Adds an internal comment to an invoice | write | not documented |
| `ZOHO_BOOKS_EMAIL_INVOICE` | Emails an invoice to the customer | write, outside | a confirmation only |
| `ZOHO_BOOKS_REMIND_CUSTOMER_INVOICE` | Emails a payment reminder for one invoice | write, outside | a confirmation only |

**What to keep:** the organization id, needed on every call; each invoice's
id, to read, remind or change that same invoice later; a customer's contact
id. Never look an invoice up by its customer's name when the taf kept its id.

**Money and outside steps wait for the owner.** Changing an invoice changes
the books, and emailing an invoice or a reminder reaches the customer.
Offer those only when the owner asked for them, and each one waits for
their yes. "Tell me which invoices are overdue" needs only the read tools.

**This tool set has no plain "create an invoice" or "record a payment"
tool.** It can make an invoice from a sales order or a recurring schedule,
and list payments, but a job that must raise new invoices or record
payments by hand can't be done with these tools.

## Don't offer

- Deleting anything (invoices, payments, contacts, attachments, comments):
  it changes the books' history.
- Voiding or writing off an invoice (`ZOHO_BOOKS_VOID_INVOICE`,
  `ZOHO_BOOKS_WRITE_OFF_INVOICE`), refunds, vendor payments, applying
  credits, and recurring invoices (create, stop, resume): only the owner or
  their accountant decides those.
- Organization and user admin (`ZOHO_BOOKS_UPDATE_ORGANIZATION`,
  `ZOHO_BOOKS_CREATE_USER`), currencies, taxes, the chart of accounts and
  journals: they change how the whole organization keeps its books.
- Banking tools (bank accounts, transactions, rules): they touch
  reconciliation.
- Bulk reminders (`ZOHO_BOOKS_SEND_BULK_INVOICE_REMINDER`) and contact
  emails: a mistake reaches many customers at once.
- Marking an invoice as sent without sending it
  (`ZOHO_BOOKS_MARK_INVOICE_AS_SENT`), unless the owner asks for exactly
  that.

## Access and plans

- **Calls per day depend on the plan:** 1,000 on Free, 2,000 on Standard,
  5,000 on Professional, and 10,000 on Premium and above, counted per
  organization.
- **Per minute:** at most 100 calls per organization. At most 5 calls at
  once on Free, and about 10 on paid plans.
- **A daily job that reads every invoice** can use up a Free plan's day.
  Filter by status or date, and page through lists rather than reading
  everything each run.
- **Paging:** lists say whether more pages follow; the taf asks for the
  next page until none is left.
- **The sign-in ends** if the person removes it from their Zoho account's
  connected apps, or changes their password and signs out everywhere.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access | The sign-in didn't give this tool's access, or the person's role in the organization doesn't allow it | "Zoho Books is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| The sign-in no longer works (401) | Removed from connected apps, or a password change signed it out | "Zoho Books needs you to sign in again. Nothing was changed." |
| Too many calls (429) | Over the minute's limit, too many at once, or the plan's daily allowance | "Zoho Books asked us to slow down, or today's allowance is used up. The taf will try again later." |
| Not found (404), or "invoice does not exist" | Deleted, or in another organization | "That invoice isn't in Zoho Books any more. The taf will skip it." |
| Rejected as malformed (400) | A required field is missing or a value doesn't fit | "Zoho Books didn't accept that change. The taf has stopped and will show you what it tried." |
| Zoho has an internal error (500) | Trouble on Zoho's side | "Zoho Books isn't answering properly right now. The taf will try again later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write, and an email to a customer is never sent twice to make up
for an unclear answer.
