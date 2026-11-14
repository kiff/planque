# QuickBooks

Checked: 2026-11-13 · Sources: the connection service's QuickBooks tool
definitions as planque reads them (the generated part beside this file),
and QuickBooks Online's API documentation (limits and throttles,
authorization, error codes).

## How it connects

- **A sign-in to one QuickBooks company**, held by the connection service;
  the taf never sees a password or token. planque's catalog lists the tools
  and, with a known sign-in, puts the ones its access doesn't cover under
  `unusable`, with why.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.

## Tools a job usually needs

Read means it changes nothing in the books. The catalog gives no `returns`
paths for QuickBooks today, and QuickBooks' answer reference couldn't be
read for this guide, so the ids each answer carries are not checked here.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `QUICKBOOKS_GET_COMPANY_INFO` | Reads the company's name, address and settings | read | the company record |
| `QUICKBOOKS_QUERY_INVOICES` | Finds invoices by customer, date range or status | read | up to 1,000 invoices a query; ids not checked here |
| `QUICKBOOKS_READ_INVOICE` | Reads one invoice: lines, total, balance due, due date | read | not checked here |
| `QUICKBOOKS_GET_INVOICE_PDF` | Gets an invoice as a PDF | read | the file |
| `QUICKBOOKS_GET_AGED_RECEIVABLES_REPORT` | Shows what each customer owes, by how late | read | a report |
| `QUICKBOOKS_QUERY_CUSTOMERS` | Finds customers by name or other fields | read | not checked here |
| `QUICKBOOKS_READ_CUSTOMER` | Reads one customer, with contact details and balance | read | not checked here |
| `QUICKBOOKS_QUERY_PAYMENTS` | Finds payments received | read | not checked here |
| `QUICKBOOKS_QUERY_ITEMS` | Finds the products and services invoices use | read | not checked here |
| `QUICKBOOKS_GET_PROFIT_AND_LOSS_REPORT` | Income and expenses over a period | read | a report |
| `QUICKBOOKS_CREATE_CUSTOMER` | Adds a customer | write | not checked here |
| `QUICKBOOKS_CREATE_INVOICE` | Makes an invoice for a customer from item lines | write, the books | not checked here |
| `QUICKBOOKS_UPDATE_SPARSE_INVOICE` | Changes only the invoice fields sent | write, the books | not checked here |
| `QUICKBOOKS_CREATE_PAYMENT` | Records a customer's payment against invoices | write, the books | not checked here |

**What to keep:** the ids of the invoices, customers and payments the taf
made or acts on; read an invoice again just before changing it. Never look an invoice up by its number when the taf has its id.

**Money waits for the owner.** Reports and reminders built from reads work
on their own. Making or changing an invoice, adding a customer, or recording
a payment changes the books: offer them only when the owner asked, and
each one waits for their yes. **This tool set has no "send invoice" tool**,
so a taf can't email an invoice from QuickBooks; it can read the invoice
(or its PDF) and the owner's own mail tools take it from there.

## Don't offer

- Deleting (`QUICKBOOKS_DELETE_ITEM`,
  `QUICKBOOKS_DELETE_BANK_ACCOUNT`): they remove records others rely on.
- Replacing a whole invoice (`QUICKBOOKS_UPDATE_FULL_INVOICE`): fields not
  sent are cleared. Use `QUICKBOOKS_UPDATE_SPARSE_INVOICE`.
- Card and bank charges (`QUICKBOOKS_CAPTURE_CHARGE`,
  `QUICKBOOKS_CREATE_ECHECK_PAYMENT`, `QUICKBOOKS_CREATE_BANK_ACCOUNT`): they
  move real money.
- The ledger and company setup: journal entries, accounts, deposits,
  transfers, tax codes and agencies, terms, preferences, company info
  (`QUICKBOOKS_UPDATE_COMPANY_INFO`, `QUICKBOOKS_UPDATE_PREFERENCES`).
- Bills and vendor payments, unless the job is about paying suppliers and
  the owner asked.
- Batch calls (`QUICKBOOKS_EXECUTE_BATCH_OPERATION`): many writes the owner
  can't check one by one.
- Employees and time tracking.

## Access and plans

- **All these tools need accounting access**; the card and bank tools need a
  separate payments access.
- **Rate, per company:** 500 calls a minute, and 10 a second per company and
  app. Over the limit, QuickBooks refuses and asks for a minute's wait.
- **Queries answer at most 1,000 records;** page through more.
- **Batch calls** are capped at 40 a minute per company and app.
- **Plans:** we checked no QuickBooks plan limit on these tools.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Too many calls (429) | Over the per-minute or per-second limit for the company | "QuickBooks asked us to slow down. The taf will try again in a minute." |
| The sign-in no longer works (401) | Access was removed, or the sign-in ended | "QuickBooks needs you to sign in again. Nothing was changed in your books." |
| Refused for missing access (403) | The sign-in doesn't cover this tool | "QuickBooks is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Not found | The record was deleted, or made inactive | "That record isn't active in QuickBooks any more. The taf will skip it." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write, and a payment or invoice is never retried without the
owner's yes.
