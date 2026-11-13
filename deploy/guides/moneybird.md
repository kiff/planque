# Moneybird

Checked: 2026-11-13 · Sources: the connection service's Moneybird tool
definitions as planque reads them (32 tools), and Moneybird's API
documentation (introduction, authentication, contacts, sales invoices).

## How it connects

- **A sign-in to Moneybird**, held by the connection service; the taf never
  sees a password or token. Moneybird lets only an administration's owner
  approve the sign-in.
- **Every call names an administration** (one set of books). A person with
  several administrations has the taf list them first and keep the one the
  job is for.
- **Access is granted by area:** sales invoices, documents (purchase invoices
  and receipts), estimates, bank, time entries and settings. A sign-in made
  for one area doesn't reach the others. With a known sign-in, planque's
  catalog puts the tools it doesn't cover under `unusable`, with why.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.

## Tools a job usually needs

Read means it changes nothing in the books. Moneybird's catalog entries give
no `returns` paths today; the ids below are the ones Moneybird's answers
carry.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `MONEYBIRD_LIST_ADMINISTRATIONS` | Lists the sets of books the sign-in can reach | read | administration ids |
| `MONEYBIRD_FILTER_CONTACTS` | Finds customers and suppliers by name or by when they changed | read | contact ids |
| `MONEYBIRD_GET_CONTACT` | Reads one contact | read | the contact's id, customer number, email |
| `MONEYBIRD_GET_CONTACT_BY_CUSTOMER_ID` | Finds a contact by its customer number | read | the contact's id |
| `MONEYBIRD_LIST_SALES_INVOICES` | Lists sales invoices, filtered by state or period, a page at a time | read | invoice ids, states, links |
| `MONEYBIRD_GET_SALES_INVOICE` | Reads one sales invoice: lines, total, state | read | the invoice's id, its number, its link and payment link |
| `MONEYBIRD_LIST_PURCHASE_INVOICES` | Lists bills received | read | not documented |
| `MONEYBIRD_LIST_TAX_RATES` | Lists the VAT rates set up in the books | read | not documented |
| `MONEYBIRD_CREATE_CONTACT` | Adds a customer or supplier | write | the new contact's id |
| `MONEYBIRD_ADD_NOTE_TO_CONTACT` | Adds a note to a contact | write | the note's id |
| `MONEYBIRD_CREATE_SALES_INVOICE` | Creates a sales invoice; it is not sent by this tool | write | the invoice's id, its state |
| `MONEYBIRD_UPDATE_SALES_INVOICE` | Changes a sales invoice's details | write | the invoice's id |
| `MONEYBIRD_UPLOAD_DOCUMENT_ATTACHMENT` | Attaches a file to a document | write | not documented |

**What to keep:** the administration id (every call needs it); each
invoice's `id` the taf created, to read or update that same invoice next
time. An invoice's number stays empty while it's a draft, so look a draft
up by its `id`, never by number. Never search for an invoice the taf made:
use the id it kept.

**Money and outside steps wait for the owner.** Creating or changing an
invoice changes the books. Offer those only when the owner asked, and each
one waits for their yes. This set of tools has no "send invoice" tool. A
payment tool exists (`MONEYBIRD_CREATE_EXTERNAL_SALES_INVOICE_PAYMENT`,
recording a payment, with its date and amount, on an invoice made in
another system): offer it only on request,
with the owner's yes for each payment.

## Don't offer

- Deleting (`MONEYBIRD_DELETE_CONTACT`, `MONEYBIRD_DELETE_DOCUMENT`,
  `MONEYBIRD_DELETE_CONTACT_NOTE`, `MONEYBIRD_DELETE_CONTACT_PERSON`): books
  records shouldn't disappear. `MONEYBIRD_ARCHIVE_CONTACT` keeps the
  contact, only on request.
- `MONEYBIRD_REQUEST_CONTACT_PAYMENTS_MANDATE_URL`: it makes a link asking a
  customer to set up direct debit, with a small first payment. It needs
  Moneybird's own payments switched on, and reaches the customer's money.
- Tax reports (`MONEYBIRD_GET_TAX_REPORT`) as a basis for filing: offer
  them for reading only, and the owner checks before filing anything.

## Access and plans

- **Rate:** Moneybird allows 150 calls per 5 minutes from one address, and
  50 per 5 minutes for reports. Over that it says how long to wait. A job
  that reads many invoices should page and pause.
- **Paging:** 50 items a page by default, 100 at most.
- **Account limits:** when the Moneybird account reaches its limit, Moneybird
  refuses the call as payment required.
- **Test books:** a sandbox administration has every feature without a
  subscription, but its invoices are watermarked.
- **The sign-in** doesn't expire today, by Moneybird's docs; it ends when
  the access is revoked.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The sign-in no longer works (401) | Access was revoked, or the sign-in is invalid | "Moneybird needs you to sign in again. Nothing was changed in your books." |
| Account limit reached (402) | The Moneybird account hit its plan's limit | "Your Moneybird account has reached its limit, so the taf can't add more. Check your plan in Moneybird." |
| Refused, blocked (403) | Moneybird has blocked the calling address | "Moneybird is refusing the taf's requests right now. The taf has stopped; this needs Moneybird's support." |
| A step's area isn't covered | The sign-in didn't include that area (bank or settings, say) | "Moneybird is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Not found (404) | The contact or invoice was deleted, or a draft was looked up by number | "That record isn't in Moneybird any more. The taf will skip it." |
| Rejected (422) | A required field is missing or a value isn't valid | "Moneybird didn't accept that invoice: some details are missing or wrong. The taf has stopped and will show you what it tried." |
| Too many calls (429) | Over the 5-minute limit | "Moneybird asked us to slow down. The taf will try again in a few minutes." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write: a taf that can't find its invoice asks, rather than
creating another one.
