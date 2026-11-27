# Qonto

Checked: 2026-11-27 · Sources: planque's reviewed entries `qonto-key@1`,
`qonto-key@2` (see [the generated parts](generated/qonto-key-2.md)) and
`qonto-sandbox@1`, Qonto's published invoicing API description, and
Qonto's developer pages (Base URLs & environments, Testing in the Sandbox).
Failure table checked against planque's answers on 2026-11-16.

## Say these five things before the sign-in card

Tell the owner all five, in your own words, in the message before the Sign in
card, and ask them to say yes. Don't leave one for later.

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
   exporting them; client invoices and credit notes. It can create client
   invoices and add clients. Through Qonto's own tools that write access also
   reaches further than a draft: finalising an invoice, sending it by email or
   as an e-invoice (an e-invoice can't be recalled), marking it paid, and
   deleting drafts or clients. Say so, and that the Card can hold those for
   their approval. It can't make payments, approve or refuse card
   and transfer requests, set trusted beneficiaries, or store their login.

## How it connects

- **Qonto's own sign-in is the first route.** Catalog entry `qonto`: the
  owner signs in on Qonto's own page and approves; no key is pasted and no
  taf sees the sign-in. Its scopes are reviewed. Planque never asks a bank for a scope that
  moves money, so no transfer is possible on it. It reads the business, its
  accounts, balances, IBANs, team members and transactions, clients, client and supplier invoices, products,
  attachments and e-invoicing, and writes clients and client invoices.
- **Before the owner signs in, tell them what Qonto needs.** Qonto doesn't
  list planque yet, so **an owner or admin of their Qonto organisation must
  first approve planque's address in Qonto:** Qonto AI, then Qonto MCP, then
  Advanced settings, add `https://planque.kiff.dev/oauth/callback`. Qonto's
  sign-in page says, until the address is approved, that it isn't on its list
  and that Qonto hasn't verified it. Once it is approved, the sign-in goes
  on to Qonto's consent screen, which lists what the connection can't do: no
  payments, and no approving cards or transfers. A member who isn't
  an owner or admin can't connect. Whether to approve it is their choice,
  and they can remove the address in the same screen at any time. If they'd
  rather not, nothing is connected.
- **Tell them about professional secrecy before they sign in.** Qonto's
  consent page also asks them to accept the lifting of Qonto's professional
  secrecy for the tool they chose: it lets Qonto pass planque the data that
  page lists. Whether to accept is their choice; without that tick, nothing
  is connected.
- **The sign-in is with the owner's own Qonto login.** A Google account that
  isn't the one their Qonto organisation uses (a work or test address) makes
  Qonto offer to create a company: never create one, close the page. On a
  phone, Qonto's page may show only "Continue with Apple": open the browser's
  "Desktop site" view, which shows Google and "another sign-in method" (email
  and password, for an account made that way).
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
