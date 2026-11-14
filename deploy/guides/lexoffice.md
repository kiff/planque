# Lexoffice

Checked: 2026-11-14 · Sources: the connection service's Lexoffice tool
definitions as planque reads them (43 tools), and the app's public API
documentation (now published under its new name, Lexware Office).

## How it connects

- **By the owner's API key**, made on the app's public API page in the
  owner's account and pasted once. The connection service holds it; the
  taf never sees it.
- **One key reaches one organization's books.** `LEXOFFICE_GET_PROFILE`
  names the company and its tax setup, so a taf can check it's in the right
  books before it writes.
- **With a known sign-in**, planque's catalog puts the tools the key can't
  run under `unusable`. Ask for every tool the job needs in the first
  request, so nothing is asked twice.

## Tools a job usually needs

Read means it changes nothing in the books. The catalog gives no `returns`
paths for Lexoffice today; the ids below are those the app's answers carry.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `LEXOFFICE_GET_PROFILE` | Reads the company the key belongs to and its tax setup | read | the organization id, company name, tax type |
| `LEXOFFICE_LIST_VOUCHERLIST` | Lists invoices, quotes, credit notes and receipts by type and status (open, overdue, paid…), a page at a time | read | voucher ids, numbers, due dates, open amounts, contact ids |
| `LEXOFFICE_GET_INVOICE` | Reads one invoice: lines, totals, status | read | the invoice id, its version |
| `LEXOFFICE_GET_PAYMENT` | Reads what's been paid on an invoice or receipt, and what's still open | read | the open amount, payment status, paid date |
| `LEXOFFICE_LIST_CONTACTS` | Finds customers and suppliers by name, email or number | read | contact ids |
| `LEXOFFICE_GET_CONTACT` | Reads one contact | read | the contact id, its version |
| `LEXOFFICE_GET_QUOTATION` | Reads one quote | read | the quote id, its version |
| `LEXOFFICE_CREATE_CONTACT` | Adds a customer or supplier | write | the new contact's id, its link, its version |
| `LEXOFFICE_UPDATE_CONTACT` | Changes a contact; needs its current version | write | the contact's id, its link, its new version |
| `LEXOFFICE_CREATE_QUOTATION` | Writes a quote, as a draft unless told to finalize | write | the new quote's id, its link, its version |
| `LEXOFFICE_CREATE_INVOICE` | Writes an invoice, as a draft unless told to finalize | write | the new invoice's id, its link, its version |
| `LEXOFFICE_CREATE_VOUCHER` | Books a receipt or incoming bill | write | the new voucher's id, its link, its version |
| `LEXOFFICE_UPLOAD_VOUCHER_FILE` | Attaches a scan or PDF to a booked voucher | write | not documented |

**What to keep:** each document's id, to read it or check its payment
later; each contact's id and version, since a change must send the current
version. Find vouchers by status in the voucher list ("overdue", "open")
rather than by searching names.

**Create as a draft.** Leave `finalize` off: a draft is something the owner
checks and finalizes in the app. A finalized invoice is open in the books
and the API can't change its status afterwards, so finalizing is the
owner's call, not the taf's.

## Don't offer

- Finalizing on create, unless the owner asked for it in so many words.
- `LEXOFFICE_DELETE_ARTICLE`: removes a product from the catalog for good.
- Event subscription tools: the taf reads on its schedule, it doesn't
  listen.
- Credit notes, delivery notes and order confirmations, unless the job is
  about them: each is a document the owner's customer receives.

## Access and plans

- **Rate:** 2 calls a second per client, across all tools. Over that, the
  call is refused and not done. Repeated bursts can get a client blocked
  for seconds to minutes, and it stays blocked until it sends fewer calls.
- **Every change carries a version.** If someone changed a record since the
  taf read it, the change is refused; read it again first.
- **The app has no delete** for invoices, quotes, vouchers or contacts
  through the API.
- **Plan problems are their own answer:** a call the owner's contract
  doesn't cover is refused as such. Which plans include API access isn't
  stated in the API docs, so it isn't given here. New accounts get every
  feature free for 30 days.
- **A call that runs over 30 seconds** is cut off.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The key doesn't work (401) | The key is missing or wrong, or the API client was switched off | "Lexoffice needs a new key. Nothing was changed." |
| Refused for the contract (402) | The owner's contract doesn't cover this action | "Your Lexoffice plan doesn't include this step. You can check your plan in the app, or the taf can do without it." |
| Refused for missing access (403) | The key's access doesn't cover this action | "Lexoffice is only partly connected: this step needs access your key doesn't give. Want to add a key that allows it, or do without this step?" |
| Not found (404) | The document or contact doesn't exist | "That isn't in Lexoffice any more. The taf will skip it." |
| Rejected (406) | A value failed the app's checks, or a file was asked for on a draft | "Lexoffice didn't accept that entry. The taf has stopped and will show you what it tried." |
| Conflict (409) | Someone changed the record since the taf read it, or the document is still a draft | "Someone changed this in Lexoffice in the meantime. The taf will read it again before trying once more." |
| Too many calls (429) | Over two calls a second | "Lexoffice asked us to slow down. The taf will try again in a moment." |
| Lexoffice unavailable (5xx) | Trouble on the app's side, or a call ran too long | "Lexoffice isn't answering right now. The taf will check whether anything was saved before trying again." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write: a taf that can't find a contact asks before adding one.
