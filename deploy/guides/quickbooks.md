# QuickBooks

Checked: 2026-11-13 · Sources: the connection service's QuickBooks tool
definitions as planque reads them (114 tools), and Intuit's QuickBooks Online
developer docs: call limits and throttles, scopes, the OAuth FAQ and the
error codes page.

## How it connects

- **A sign-in to QuickBooks Online on Intuit's page**, held by the connection
  service; the taf never sees a password or token. One sign-in reaches one
  QuickBooks company. With a known sign-in, the catalog puts the tools its
  access doesn't cover under `unusable`, with why.
- **The person who signs in must be an admin** of that QuickBooks company.
- **Two kinds of access.** Bookkeeping (invoices, bills, customers, reports)
  is one; card and bank payments is another. Ask for every tool the job needs
  in the first request, so the owner signs in once.
- **A sign-in left unused for 100 days stops working,** and the owner signs
  in again.

## Tools a job usually needs

Read changes nothing in the books. Returns are the ids Intuit's reference
gives for that call; a taf keeps the `Id`, and the `SyncToken` before any
update.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `QUICKBOOKS_QUERY_INVOICES` | Lists invoices by customer, date or status, a page at a time | read | invoice ids, doc numbers, balances |
| `QUICKBOOKS_READ_INVOICE` | Reads one invoice | read | the invoice, its `Id` and `SyncToken` |
| `QUICKBOOKS_GET_INVOICE_PDF` | Gets an invoice as a PDF | read | the PDF file |
| `QUICKBOOKS_QUERY_CUSTOMERS` / `QUICKBOOKS_READ_CUSTOMER` | Finds or reads customers | read | customer ids, display names |
| `QUICKBOOKS_QUERY_BILLS` / `QUICKBOOKS_GET_BILL` | Finds or reads supplier bills | read | bill ids, due dates, balances |
| `QUICKBOOKS_GET_CHANGED_ENTITIES` | Lists what changed since a date | read | ids of changed records, by type |
| `QUICKBOOKS_GET_PROFIT_AND_LOSS_REPORT` | Profit and loss for a period | read | the report's rows and totals |
| `QUICKBOOKS_GET_AGED_RECEIVABLES_REPORT` | Who owes what, by age | read | the report's rows |
| `QUICKBOOKS_GET_REPORT_AGED_PAYABLES` | What the company owes, by age | read | the report's rows |
| `QUICKBOOKS_CREATE_CUSTOMER` | Adds a customer | write | the customer `Id`, `SyncToken` |
| `QUICKBOOKS_CREATE_INVOICE` | Writes an invoice; it isn't emailed | write | the invoice `Id`, `SyncToken`, doc number |
| `QUICKBOOKS_UPDATE_SPARSE_INVOICE` | Changes only the given fields of an invoice | write | the invoice `Id`, new `SyncToken` |
| `QUICKBOOKS_CREATE_BILL` | Records a supplier bill | write | the bill `Id`, `SyncToken` |

**What to keep:** the `Id` of each record the taf made, and its latest
`SyncToken`: an update with an old one is refused. Never search for a record
the taf made: use the id it kept.

**Start read-only.** "Tell me who is late paying" or "a monthly summary" is
reads and reports only. Offer a creating tool only when the owner asked for
it, and say each new invoice or bill is shown to them before it is written.

## Don't offer

- Deleting: `QUICKBOOKS_DELETE_ITEM`, `QUICKBOOKS_DELETE_BANK_ACCOUNT`.
- Moving money or touching card and bank details: `QUICKBOOKS_CAPTURE_CHARGE`,
  `QUICKBOOKS_CREATE_ECHECK_PAYMENT`, `QUICKBOOKS_CREATE_BANK_ACCOUNT`,
  `QUICKBOOKS_GET_BANK_ACCOUNT`, `QUICKBOOKS_LIST_CARDS`.
- Company settings and the chart of accounts: `QUICKBOOKS_UPDATE_COMPANY_INFO`,
  `QUICKBOOKS_UPDATE_PREFERENCES`, `QUICKBOOKS_UPDATE_TERM`,
  `QUICKBOOKS_UPDATE_PAYMENT_METHOD`, `QUICKBOOKS_CREATE_ACCOUNT`,
  `QUICKBOOKS_CREATE_TAX_AGENCY`, `QUICKBOOKS_CREATE_TAXSERVICE_TAXCODE`.
- People: `QUICKBOOKS_CREATE_EMPLOYEE`.
- Entries an accountant should make: `QUICKBOOKS_CREATE_JOURNAL_ENTRY`,
  `QUICKBOOKS_UPDATE_TRANSFER`.
- Bulk and whole-record tools: `QUICKBOOKS_EXECUTE_BATCH_OPERATION`, and
  `QUICKBOOKS_UPDATE_FULL_INVOICE` (it blanks any field left out; use the
  sparse update).
- `QUICKBOOKS_SEND_CREDIT_MEMO` emails the customer: only when the owner
  asked for sending, and each send waits for their yes.

## Access and plans

- **Bookkeeping tools** need accounting access; the payment tools need
  payments access, which the bookkeeping sign-in doesn't give.
- **Rate, per company:** 500 calls a minute and 10 a second. Over that,
  QuickBooks refuses and asks to wait 60 seconds. Batches are capped at 40 a
  minute.
- **Big answers:** a query returns at most 1,000 records; page through more.
  A call running over 120 seconds times out.
- **Closed books:** if the owner has closed a period, entries dated in it
  are refused.
- **Subscription:** when a company's QuickBooks subscription has ended, was
  cancelled or has a billing problem, QuickBooks refuses the calls.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The sign-in no longer works (401) | Access was removed, the sign-in expired, or 100 days without use | "QuickBooks needs you to sign in again. Nothing was changed." |
| Refused for missing access (403) | The sign-in didn't give this tool's access, payments say, or the user isn't an admin | "QuickBooks is only partly connected: this step needs access you didn't give. An admin of the company can sign in again and allow it." |
| Too many calls (429) | Over the per-minute or per-second limit | "QuickBooks asked us to slow down. The taf will try again in a minute." |
| Stale object (code 5010) | The record changed since the taf read it | Read it again and retry once; tell the owner only if it fails again. |
| Duplicate name or number (codes 6240, 6140) | A customer name or document number already exists | "QuickBooks already has one with that name or number. Should I use the existing one, or pick another?" |
| Period closed (code 6200) | The books are closed for that date | "That date falls in a period you've closed in QuickBooks. Nothing was written. Use a later date, or reopen the period first." |
| Company status (code 6190) | The subscription ended or has a billing problem | "QuickBooks won't accept changes until the company's subscription is sorted out. Nothing was written." |
| Not found (code 610) | The record was deleted | "That record isn't in QuickBooks any more. The taf will skip it." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
