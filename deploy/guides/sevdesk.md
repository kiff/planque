# Sevdesk

Checked: 2026-11-14 · Sources: the connection service's Sevdesk tool
definitions as planque reads them (134 tools), and Sevdesk's published API
description (its OpenAPI file and the guides inside it).

## How it connects

- **By the owner's API token**, pasted once and held by the connection
  service; the taf never sees it. Each Sevdesk administrator has one token,
  a 32-character code.
- **The token acts for that administrator's Sevdesk account:** contacts,
  invoices, vouchers and bank accounts. The API description shows no way
  to narrow a token, so the Card is what keeps a taf to its job.
- **A held connection is reused** for a new taf when it covers the tools
  asked. Ask for every tool the job needs in the first request.

## Tools a job usually needs

Read means it changes nothing in Sevdesk. The catalog gives no `returns`
paths for Sevdesk today; the cells below are what Sevdesk's API description
shows.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `SEVDESK_GET_CONTACTS` | Lists or searches contacts, a page at a time | read | contact ids, names, customer numbers |
| `SEVDESK_GET_CONTACT_BY_ID` | Reads one contact | read | the contact |
| `SEVDESK_GET_INVOICES` | Lists invoices, filtered by status, contact, number or date | read | invoice ids, numbers, status, totals |
| `SEVDESK_GET_INVOICE_BY_ID` | Reads one invoice | read | the invoice, its status and amounts |
| `SEVDESK_GET_INVOICE_POSITIONS` | Reads an invoice's lines | read | line ids, amounts |
| `SEVDESK_GET_INVOICE_IS_PARTIALLY_PAID` | Says whether an invoice is part paid | read | true or false |
| `SEVDESK_GET_VOUCHERS` | Lists receipts and bills | read | voucher ids, supplier, date |
| `SEVDESK_GET_CHECK_ACCOUNTS` | Lists payment and bank accounts | read | account ids |
| `SEVDESK_GET_CHECK_ACCOUNT_TRANSACTIONS` | Lists bank transactions | read | transaction ids, amounts, payer names, dates |
| `SEVDESK_CREATE_CONTACT` | Adds a customer or supplier | write | the new contact's id |
| `SEVDESK_SAVE_INVOICE` | Creates an invoice with its lines | write | the new invoice, its lines |
| `SEVDESK_SEND_INVOICE_VIA_EMAIL` | Emails an invoice to the customer | write, outside | the sent email's id |
| `SEVDESK_BOOK_INVOICE_AMOUNT` | Records a payment against an invoice | write, books | the change, with old and new status |
| `SEVDESK_SAVE_VOUCHER` | Records a receipt or bill, with its lines | write, books | the new voucher, its lines |
| `SEVDESK_UPLOAD_VOUCHER_FILE` | Uploads a receipt file to attach to a voucher | write | the file name to use |

**What to keep:** each invoice's id, to read, send or book that same
invoice later; a contact's id, to put on new invoices; the file name an
upload gives, to attach it in the next voucher. Never search for an
invoice the taf made: use the id it kept.

**Status:** invoices go from draft to open to paid. The list tools filter by
it, so "open invoices past due" is a read.

**Books and outside steps wait for the owner.** Creating an invoice,
emailing it, booking a payment and recording a voucher are offered only
when the owner asked, and each one waits for their yes. Once an invoice is
emailed it's marked sent, and Sevdesk doesn't let it change after that.

## Don't offer

- **Deleting anything** (the 14 delete tools): contacts, orders, credit
  notes, accounts and transactions.
- **Cancelling an invoice** (`SEVDESK_CANCEL_INVOICE`): it makes a
  cancellation invoice that is booked as paid at once.
- **Resetting** an invoice, credit note or voucher to open
  (`SEVDESK_RESET_INVOICE_TO_OPEN` and the like): it unlinks its bank
  transactions.
- **Locking records for good** (`SEVDESK_ENSHRINE_CHECK_ACCOUNT_TRANSACTION`):
  locked records can't be changed or reset.
- **Account set-up:** payment methods, bank accounts, clearing accounts,
  transaction rules, custom fields and export settings.
- **Marking as sent without sending** (`SEVDESK_SEND_INVOICE_BY`): it
  changes the invoice's status without the customer getting it.

## Access and plans

- **Lists come a page at a time,** with a page size and an offset.
- **VAT is checked strictly.** A voucher's tax rate must fit its booking
  account and tax rule; Sevdesk's receipt guidance tools say which fit.
  Otherwise the voucher is refused.
- **New vouchers start as draft or unpaid.** A voucher can't be created as
  paid; it is paid by booking an amount on it.
- **The invoice download** (`SEVDESK_DOWNLOAD_INVOICE_PDF`) counts as a
  write: its options decide whether downloading also marks the invoice
  sent.
- **Which plan includes the API, and its rate limits,** aren't given in the
  API description, so they're left out here.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The token doesn't work (401) | The token is wrong or no longer valid | "Sevdesk needs a new API token. Nothing was changed." |
| Refused as invalid (422) | A tax rate or booking account doesn't fit, or a status change isn't allowed | "Sevdesk didn't accept that entry: its tax or account doesn't fit. The taf has stopped and will show you what it tried." |
| Bad request (400) | A field is missing or wrong, or the invoice isn't there | "Sevdesk couldn't find or accept that invoice. The taf will skip it and tell you which." |
| Changing a sent or locked invoice | Sent and locked invoices can't change | "That invoice is already sent, so Sevdesk won't change it. Correct it in Sevdesk if it needs changing." |
| Server error (500) | Trouble on Sevdesk's side | "Sevdesk isn't answering properly right now. The taf will try again later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write, and a payment is never booked twice: the taf reads the
invoice's status first.
