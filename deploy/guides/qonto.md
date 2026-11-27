# Qonto

Checked: 2026-11-27 · Sources: planque's reviewed entries `qonto-key@1`,
`qonto-key@2` (see [the generated parts](generated/qonto-key-2.md)) and
`qonto-sandbox@1`, Qonto's published invoicing API description, and
Qonto's developer pages (Base URLs & environments, Testing in the Sandbox).
Failure table checked against planque's answers on 2026-11-16.

## How it connects

- **By the business's API key**, reviewed in planque's catalog. The owner
  makes it in Qonto's settings (Integrations and partnerships, API key) and
  pastes it as `<organisation login>:<secret key>`. The key and the taf's
  Card come in one approval.
- **Two reviewed entries:** `qonto-key@1` reads only; `qonto-key@2` reads more
  and can draft an invoice. Ask for `@2` only when the job drafts.
- **Qonto's own server sign-in** isn't usable yet: it can't send the owner
  back to planque after signing in.
- **Qonto no longer issues the key to a new developer account:** its
  Developer Portal creates an OAuth app (a client id and secret) instead.
  The OAuth route is planque's own client, as for Google Sheets:
  `qonto-sandbox@1` is the first entry, on Qonto's **sandbox** (test data):
  the owner signs in on Qonto's page and approves, no key is pasted. The
  sandbox wants a staging token on every request, including the token
  call; planque keeps it with the client secret and sends it to Qonto's own
  hosts only. Reads only for now (the business and its accounts, and a
  bank account's transactions); a transfer request comes later.

## Tools

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `qonto_organization` | Reads the business and its bank accounts | read | account ids and IBANs |
| `qonto_list_clients` | Lists the clients | read | client ids |
| `qonto_list_client_invoices` | Lists invoices sent to clients | read | invoice ids, numbers, status |
| `qonto_get_client_invoice` | Reads one client invoice | read | the invoice, its link |
| `qonto_list_quotes` (`@2`) | Lists quotes | read | quote ids, status |
| `qonto_get_quote` (`@2`) | Reads one quote | read | the quote, its items |
| `qonto_list_supplier_invoices` (`@2`) | Lists invoices received | read | invoice ids, amounts, due dates |
| `qonto_draft_client_invoice` (`@2`) | Drafts an invoice from an accepted quote | write, draft only | the draft's id |

**What to keep:** a client's id and a quote's id; the draft's id, so the
owner can find it in Qonto.

**A draft is never sent.** The owner reviews and sends it in Qonto. The
draft's IBAN must be one of the business's own accounts: planque refuses
any other.

## Don't offer

- Transfers, cards, or anything that moves money: not in either entry, and
  a key that can do them shouldn't be pasted for a taf.
- Sending or cancelling invoices: the owner does that in Qonto.

## Access and plans

- **Drafting an invoice and listing clients work on every Qonto plan**, per
  Qonto's pages for them.
- **The key belongs to the business,** not a person: anyone with it reads
  the business's accounts.
- **Rate:** Qonto allows 1,000 calls per 10 seconds and 10,000 per 10
  minutes, counted per address, and slows a key down after many refused
  calls in an hour. A taf reading invoices daily stays far below it.

## When it fails

planque calls Qonto with the key itself, so Qonto's answers reach the taf
by its address: "thirdparty.qonto.com refused it (403): …". An answer of 500
or above reads "may or may not have done it": read the invoices before any
new call. planque's own refusals are in the part below.

| What the taf gets | Why | Tell the owner |
|---|---|---|
| Refused (401) | Wrong format, revoked, or the login and secret swapped | "Qonto didn't accept the key. Tap Connect on this taf and paste it again as your organisation login, a colon, then the secret key." |
| Refused (403) | This key isn't allowed to call that part of Qonto | "Qonto doesn't let this key do that step. The taf will do without it." |
| Refused, on the draft: IBAN | The IBAN isn't one of the business's accounts | "That IBAN isn't one of your Qonto accounts, so the draft wasn't made." |
| Refused, on the draft: quote not found | The quote doesn't exist, or belongs to another organisation | "That quote isn't in your Qonto account. Nothing was drafted." |
| Refused, on the draft: conflict | The invoice number is already used, or the invoice conflicts with another | "Qonto refused the draft because it clashes with an existing invoice. Nothing was drafted." |
| Refused (429), too many calls | Over Qonto's rate | "Qonto asked us to slow down. The taf will try again in a minute." |
