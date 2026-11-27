# Qonto

Checked: 2026-11-27 · Sources: planque's reviewed entries `qonto-key@1`,
`qonto-key@2` (see [the generated parts](generated/qonto-key-2.md)) and
`qonto-sandbox@1`, Qonto's published invoicing API description, and
Qonto's developer pages (Base URLs & environments, Testing in the Sandbox).
Failure table checked against planque's answers on 2026-11-16.

## Put these five lines on the card before the sign-in

Put all five, as written here, in the text of the card that asks the owner to
go ahead with the sign-in, one a line, above its question. Don't leave one for
later, and don't reword them: a reworded line ("read-only") is refused.

1. **Qonto's own sign-in, no key.** They sign in on Qonto's own page; no one
   makes or pastes a key.
2. **An owner or admin approves planque's address first:** Qonto AI, then
   Qonto MCP, then Advanced settings, add `https://planque.kiff.dev/oauth/callback`.
3. **Use their own Qonto login, and "Desktop site" on a phone.** If Qonto's
   page offers to create a company, they close it; if it shows only Apple,
   they open the browser's "Desktop site" view.
4. **Professional secrecy.** Qonto's consent page asks them to accept the
   lifting of Qonto's professional secrecy for the tool they chose. Their
   choice; without the tick, nothing is connected.
5. **What it can read, as Qonto's page lists it:** the account balance, IBAN,
   transactions, the business and its team; supplier invoices; receipts and
   other attachments; client invoices and credit notes. It can draft a client
   invoice, which they send themselves, and add or update clients. It can't
   finalise, cancel, send or mark an invoice paid, or delete an invoice or a
   client. It can't make payments, approve or refuse card
   and transfer requests, set trusted beneficiaries, or store their login.

## How it connects

- **Qonto's own sign-in is the first route.** Catalog entry `qonto`: the
  owner signs in on Qonto's own page and approves; no key is pasted and no
  taf sees the sign-in. Its scopes are reviewed (the list above). Planque never
  asks for a scope that moves money, so no transfer is possible on it.
- **Before the sign-in, the five lines above are on the card.** The address approval is
  an owner's or admin's choice, can be removed in the same screen at any
  time, and a member who isn't an owner or admin can't connect. Until it is
  approved, Qonto's page says the address isn't on its list and isn't
  verified; once approved, the sign-in goes on to Qonto's consent screen. If
  they'd rather not, or won't tick the secrecy line, nothing is connected.
- **Own login.** A Google account that isn't the one their Qonto organisation
  uses makes Qonto offer to create a company: never create one, close the
  page. "Desktop site" shows Google and "another sign-in method" (email and
  password).
- **Do not ask for an API key.** Qonto no longer issues the login-and-secret
  key to a new developer account (its Developer Portal makes an OAuth app
  instead). The reviewed key entries `qonto-key@1` and `qonto-key@2` are only
  for an owner who already holds such a key (`<organisation login>:<secret
  key>`); never propose them to anyone else, and never tell an owner to make
  a key.
- **The sandbox twin, for test data:** `qonto-sandbox@1`, planque's own
  OAuth client on Qonto's sandbox, offered to test accounts only. It reads
  the business and its accounts, and a bank account's transactions. The
  sandbox wants a staging token on every request, including the token call;
  planque keeps it with the client secret and sends it to Qonto's own hosts
  only.

## Tools

These are the tools of the key entries `qonto-key@1` and `qonto-key@2`. The
sign-in route (`qonto`) uses Qonto's own MCP tools, so don't name one of
these for it: read the tools the connection lists.

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
