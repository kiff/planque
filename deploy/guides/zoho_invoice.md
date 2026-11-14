# Zoho Invoice

Checked: 2026-11-14 · Sources: the connection service's Zoho Invoice tool
definitions as planque reads them (137 tools), and Zoho Invoice's API
documentation (introduction, sign-in, errors, paging, invoices, contacts,
customer payments, estimates, items).

## How it connects

- **A sign-in to Zoho**, held by the connection service; the taf never sees
  a password or token. planque's catalog lists the tools and, with a known
  sign-in, puts the ones its access doesn't cover under `unusable`, with why.
- **Zoho keeps each account in one region** (Europe, the United States,
  India and others), and an account is reached only through its own
  region's addresses.
- **Every call names one organization.** The taf reads the list of the
  person's organizations first and keeps the one the owner means.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.

## Tools a job usually needs

Read means it changes nothing. The catalog's `returns` gives no paths for
Zoho Invoice today; the ids below are those Zoho's answers carry.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `ZOHO_INVOICE_LIST_ORGANIZATIONS` | Lists the person's organizations | read | organization ids |
| `ZOHO_INVOICE_LIST_INVOICES` | Lists invoices, filtered by customer, status, balance or dates, a page at a time | read | invoice ids, numbers, balances, due dates |
| `ZOHO_INVOICE_GET_INVOICE` | Reads one invoice with its lines | read | the invoice id, number, status, balance and link |
| `ZOHO_INVOICE_LIST_INVOICE_PAYMENTS` | Lists payments recorded on one invoice | read | payment ids |
| `ZOHO_INVOICE_LIST_PAYMENTS` | Lists payments received, by customer or date | read | payment ids and numbers |
| `ZOHO_INVOICE_LIST_CONTACTS` / `ZOHO_INVOICE_GET_CONTACT` | Lists or reads customers | read | contact ids, names, emails |
| `ZOHO_INVOICE_LIST_ESTIMATES` | Lists quotes | read | estimate ids and numbers |
| `ZOHO_INVOICE_LIST_ITEMS` | Lists the products and services sold | read | item ids |
| `ZOHO_INVOICE_CREATE_CONTACT` | Adds a customer | write | the new contact id |
| `ZOHO_INVOICE_CREATE_ESTIMATE` | Writes a quote for a customer; sends nothing | write | the estimate id and number |
| `ZOHO_INVOICE_CREATE_INVOICE` | Writes an invoice for a customer; sends nothing | write | the invoice id, number and link |
| `ZOHO_INVOICE_CREATE_CUSTOMER_PAYMENT` | Records a payment received against one or more invoices | write | the payment id |
| `ZOHO_INVOICE_EMAIL_INVOICE` | Emails an invoice to the customer | write, outside | a confirmation only |
| `ZOHO_INVOICE_EMAIL_ESTIMATE` | Emails a quote to the customer | write, outside | a confirmation only |

**What to keep:** the organization id, needed on every call; each invoice's
and estimate's id, to read, email or record payment against that same one
later; each customer's contact id. Never look an invoice up by name when the
taf kept its id.

**Money and outside steps wait for the owner.** Creating an invoice or
recording a payment changes the books, and emailing an invoice or quote
reaches the customer. Offer those only when the owner asked for them, and
each one waits for their yes. A good first job is read-only: "which
invoices are overdue, and by how much".

**Reminders:** this tool set sends payment reminders only in bulk
(`ZOHO_INVOICE_SEND_BULK_INVOICE_REMINDER`), so a reminder job sends one
invoice at a time with `ZOHO_INVOICE_EMAIL_INVOICE`, each on the owner's
yes. There's no tool to change an invoice once made.

## Don't offer

- Deleting anything (invoices, payments, contacts, estimates): it changes
  the books' history.
- Voiding or writing off an invoice (`ZOHO_INVOICE_MARK_INVOICE_AS_VOID`,
  `ZOHO_INVOICE_WRITE_OFF_INVOICE`), refunds
  (`ZOHO_INVOICE_REFUND_CUSTOMER_PAYMENT`,
  `ZOHO_INVOICE_CREATE_REFUND_CREDIT_NOTE_REFUNDS`), and credit notes: only
  the owner decides those.
- Organization and user admin (`ZOHO_INVOICE_CREATE_USER`), currencies,
  exchange rates and taxes: they change the whole organization.
- Bulk reminders and multiple-estimate emails: a mistake reaches many
  customers at once.
- Free-text emails and statements to a customer
  (`ZOHO_INVOICE_SEND_CONTACT_EMAIL`, `ZOHO_INVOICE_EMAIL_CONTACT_STATEMENT`),
  and marking an invoice sent without sending it, unless the owner asks for
  exactly that.
- Recurring invoices (create, stop, resume): they bill customers on a
  schedule without a yes each time.
- Opening the customer portal (`ZOHO_INVOICE_ENABLE_PORTAL_ACCESS`): it
  gives an outside person a login.

## Access and plans

- **Per minute:** at most 100 calls per organization.
- **Per day:** 1,000 calls on the Free plan, counted per organization.
  Zoho's page gives no other plan's figure here.
- **At once:** 5 calls on Free, and about 10 on paid plans.
- **A daily job that reads every invoice** can use up the day. Filter by
  status or date, and page through lists rather than reading all each run.
- **Paging:** lists say whether more pages follow.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused, not allowed (403) | The person's role in the organization doesn't allow this, or they aren't in that organization | "Zoho Invoice didn't let the taf do this. Your role there may not allow it; an admin of that organization can change it." |
| Refused for missing access | The sign-in didn't give this tool's access | "Zoho Invoice is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| The sign-in no longer works (401) | The sign-in was removed or is no longer valid | "Zoho Invoice needs you to sign in again. Nothing was changed." |
| Too many calls (429) | Over the minute's limit, too many at once, or the day's allowance | "Zoho Invoice asked us to slow down, or today's allowance is used up. The taf will try again later." |
| Not found (404) | Deleted, or in another organization | "That invoice isn't in Zoho Invoice any more. The taf will skip it." |
| Rejected as malformed (400) | A required field is missing or a value doesn't fit | "Zoho Invoice didn't accept that. The taf has stopped and will show you what it tried." |
| Server error (500) | Trouble on Zoho's side | "Zoho Invoice isn't answering properly right now. The taf will try again later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write, and an invoice is never created or emailed twice to make
up for an unclear answer.
