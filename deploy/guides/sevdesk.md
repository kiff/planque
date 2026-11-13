# Sevdesk

Checked: 2026-11-13 · Sources: the connection service's Sevdesk tool
definitions as planque reads them ([generated part](generated/sevdesk.md),
134 tools), and sevdesk's published API description (authentication,
pagination, status codes, invoice and voucher status, the 2024 bookkeeping
update).

## How it connects

- **A key the owner pastes**, held by the connection service; the taf never
  sees it after that. It is the API token of a sevdesk administrator, found
  in sevdesk's settings.
- **The key is the whole account.** Each administrator has one token, with
  no narrower access level, so the catalog can't put any Sevdesk tool under
  `unusable`. What the taf may do rests on the tools in the Card alone.
- **It doesn't expire.** It works as long as that sevdesk user exists. If
  the owner generates a new token in sevdesk, the old one stops and the key
  must be pasted again.
- **One approval for the key, the tools and the Card.** Ask for every tool
  the job needs in the first request.

## Tools a job usually needs

Read changes nothing in sevdesk. Returns are the ids sevdesk's reference
documents for each call.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `SEVDESK_GET_CONTACTS` | Lists or searches contacts (customers, suppliers) | read | contact ids, customer numbers |
| `SEVDESK_GET_CONTACT_BY_ID` | Reads one contact | read | the contact id |
| `SEVDESK_GET_INVOICES` | Lists invoices, filtered by status, number, contact or date | read | invoice ids, numbers, status |
| `SEVDESK_GET_INVOICE_BY_ID` | Reads one invoice | read | the invoice id, number, status |
| `SEVDESK_GET_INVOICE_POSITIONS` | Reads an invoice's lines | read | position ids |
| `SEVDESK_GET_INVOICE_IS_PARTIALLY_PAID` | Says whether an invoice is partly paid | read | yes or no |
| `SEVDESK_GET_VOUCHERS` / `SEVDESK_GET_VOUCHER_BY_ID` | Lists or reads receipts (vouchers) | read | voucher ids, status |
| `SEVDESK_GET_CHECK_ACCOUNTS` | Lists payment accounts (bank, cash) | read | account ids |
| `SEVDESK_GET_CHECK_ACCOUNT_TRANSACTIONS` | Lists transactions on payment accounts | read | transaction ids, status |
| `SEVDESK_GET_ORDERS` | Lists quotes and orders | read | order ids, numbers, status |
| `SEVDESK_GET_TOOLS_BOOKKEEPING_SYSTEM_VERSION` | Says which bookkeeping system the account uses | read | the version |
| `SEVDESK_GET_RECEIPT_GUIDANCE_FOR_ALL_ACCOUNTS` | Lists accounts and the tax rules and rates each allows | read | account numbers, tax rules |
| `SEVDESK_CREATE_CONTACT` | Adds a contact | write | the contact id |
| `SEVDESK_SAVE_INVOICE` | Creates an invoice with its lines | write | the invoice id, its positions |
| `SEVDESK_UPLOAD_VOUCHER_FILE` | Uploads a receipt file for a voucher | write | the file's temporary name |
| `SEVDESK_SAVE_VOUCHER` | Creates a voucher, with an uploaded file | write | the voucher id, its positions |
| `SEVDESK_SEND_INVOICE_VIA_EMAIL` | Emails an invoice to the customer and marks it sent | write, outside | the sent email's id |

**What to keep:** each contact's, invoice's and voucher's id, to read or
link it next time. Never search for something the taf made: use the id it
kept.

**Start with drafts.** Create invoices as drafts (status 100) and vouchers
as drafts (status 50), and let the owner check them in sevdesk. Offer
`SEVDESK_SEND_INVOICE_VIA_EMAIL` only when the owner asked for sending, and
say each send waits for their yes: once sent, an invoice must not be
changed.

## Don't offer

- Deleting tools (`SEVDESK_DELETE_*`): contacts, credit notes, orders,
  payment accounts and transactions.
- Locking and undoing bookkeeping: `SEVDESK_ENSHRINE_CHECK_ACCOUNT_TRANSACTION`
  (locks the transaction and every document linked to it, for good),
  `SEVDESK_CANCEL_INVOICE` (issues a cancellation invoice),
  `SEVDESK_RESET_INVOICE_TO_OPEN`, `SEVDESK_RESET_CREDIT_NOTE_TO_OPEN` and
  `SEVDESK_RESET_VOUCHER_STATE` (they unlink payments).
- Booking payments and bank lines (`SEVDESK_BOOK_INVOICE_AMOUNT`,
  `SEVDESK_BOOK_VOUCHER_AMOUNT`, `SEVDESK_CREATE_CHECK_ACCOUNT_TRANSACTION`,
  `SEVDESK_UPDATE_CHECK_ACCOUNT_TRANSACTION_BY_ID`): they change the books
  the owner's accountant relies on.
- Account settings: payment accounts, clearing and import accounts,
  payment methods, transaction rules, custom field settings, and
  `SEVDESK_UPDATE_SEV_CLIENT_EXPORT_CONFIG` (the accountant export setup).

## Access and plans

- **Status numbers:** an invoice is 100 draft, 200 open or due, 750 partly
  paid, 1000 paid. A voucher is 50 draft, 100 unpaid or due, 1000 paid.
  Invoices partly paid before mid-2024 still show 200: ask
  `SEVDESK_GET_INVOICE_IS_PARTIALLY_PAID`.
- **The 2024 bookkeeping update** replaced tax types with numbered tax
  rules. On an updated account, a line whose tax rate the chosen rule
  doesn't allow is refused. Check the version and the allowed rates first.
- **Pages:** lists take a limit (1 to 1,000) and an offset. sevdesk advises
  pages of 10 to 100.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The key no longer works (401) | The token was regenerated, or that sevdesk user was removed | "sevdesk no longer accepts the key. Paste the current API token from your sevdesk settings. Nothing was changed." |
| Refused (403) | The sevdesk user behind the key can't reach this | "sevdesk refused this for the user whose key you gave. Use an administrator's key, or do without this step." |
| Not found (404) | The contact, invoice or voucher was deleted, or the id is wrong | "That record isn't in sevdesk any more. The taf will skip it." |
| Refused as invalid (400) | A missing or wrong field | Fix it and try once; else "sevdesk didn't accept that entry. Nothing was saved." |
| Validation error (422) | A tax rate the tax rule doesn't allow, or a field that fails sevdesk's checks | "sevdesk refused the tax setting on this entry. Check the tax rule, then the taf can try again." |
| Server error (500) | A problem on sevdesk's side, or an input sevdesk didn't catch | Check whether it was saved, then try once more later; else "sevdesk had a problem with this step. The taf will try again later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
