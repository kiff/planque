# Stripe

Checked: 2026-11-14 · Sources: planque's direct Stripe connection and its
reviewed refund entry (`docs/managed-connectors.md`), the managed app's
invoice tools, Stripe's API reference and invoicing docs, and Stripe's
published server tools.

## How it connects

- **A sign-in on Stripe's own page**, to Stripe's own server: no key is pasted.
  The owner picks the Stripe account, and test or live mode, on that page.
- **One account per connection.** If the sign-in reaches exactly one account,
  planque sets it. With several, or none, planque refuses the tools and the
  owner signs in again choosing one account.
- **Stripe also appears as a managed app** in the catalog ([generated](generated/stripe.md),
  428 tools). The refunds job doesn't use it: it uses Stripe's own server,
  above, narrowed to the operations below.
- **Refunds are a reviewed entry** (`stripe-refunds@1`), switched on per
  deployment, and **test mode only** for now: a live object is refused.

## Tools a job usually needs

The direct connection narrows Stripe's general tools to named operations.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `stripe_api_details` | Describes one Stripe operation's fields, before a call | read | the fields |
| `stripe_api_read` (the account) | Reads the account the sign-in is in | read | the account id, its country, its mode |
| `stripe_api_read` (a charge) | Reads one payment | read | the charge id, amount, currency (its refunds only when asked for, or by listing refunds) |
| `stripe_api_read` (refunds) | Lists recent refunds, at most 100 | read | refund ids, amounts, the charge each is for |
| `stripe_api_write` (a refund) | Refunds a payment, whole or in part | write, money out | the refund id and its status |

**What to keep:** a refund's id, to report it and never refund twice; the
charge id the owner named. Look a payment up by its id, never by searching
for an amount.

**Amounts are in cents** (the currency's smallest unit): €12.50 is `1250`.
KIFF's Card caps refunds in euros; a refund in another currency is refused.

## Invoices: chasing overdue ones

Invoices go through the managed app's tools ([generated](generated/stripe.md)),
not the refund entry above.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `STRIPE_LIST_INVOICES` | Lists invoices, newest first, filtered to status `open` (and collection method `send_invoice`); 10 a page unless asked, at most 100 | read | invoice ids, due dates, amounts left to pay, the last id to page from |
| `STRIPE_GET_INVOICES_INVOICE` | Reads one invoice | read | its status, due date, amount left, customer, collection method |
| `STRIPE_SEND_INVOICE` | Emails the invoice to its customer now, outside Stripe's own schedule | write, reaches the customer | the invoice |

**Overdue** means `open` with a due date in the past.

**Sending only works for invoices sent for the customer to pay** (collection
method `send_invoice`), not ones Stripe charges to the customer's card
automatically: Stripe doesn't email those.

**Read the invoice again just before sending.** Stripe sends a paid invoice
too, with no mention of the payment, so send only while it's still `open`.
A draft gets finalized first, which locks its amount and customer: never
send a draft.

**Test mode sends no email**, though Stripe records the invoice as sent.

**What to keep**, in the taf's notes: each invoice's id and how many
reminders it has sent, so it never reminds twice. Offer sending only when
the owner asked for reminders, and say each one waits for their yes.

## Don't offer

- Creating customers, products, prices, payment links or subscriptions: none
  is in the reviewed entry, and each puts something in front of the owner's
  customers.
- Cancelling subscriptions, closing disputes, or payouts: they move money or
  end a customer's service, and need their own review.
- A refund in live mode, until the live entry is reviewed.
- `STRIPE_POST_V1_INVOICES_SEND_VOID_AN_INVOICE`: its name says "send", but
  it voids the invoice, and that can't be undone.
- The other void tools (`STRIPE_VOID_INVOICE`, `STRIPE_POST_INVOICES_INVOICE_VOID`,
  `STRIPE_POST_V1_INVOICES_VOID_AN_INVOICE`,
  `STRIPE_POST_V1_INVOICES_FINALIZE_VOID_AN_INVOICE`,
  `STRIPE_DELETE_V1_INVOICES_VOID_AN_INVOICE`, `STRIPE_VOID_CREDIT_NOTE`,
  `STRIPE_POST_BILLING_CREDIT_GRANTS_VOID_CREDIT_GRANT`): they cancel a bill
  or a credit for good.
- `STRIPE_DELETE_INVOICE`: it deletes a draft for good.
- `STRIPE_FINALIZE_INVOICE`: it turns a draft into a bill the customer owes.
- `STRIPE_PAY_INVOICE`: it charges the customer's saved payment method now.
- `STRIPE_MARK_INVOICE_UNCOLLECTIBLE`: it writes the debt off.

## Access and plans

- **No paid Stripe plan is needed** for the API. Stripe charges its own fees
  on payments; a refund doesn't give its fee back.
- **Test mode first:** every refund job starts on a test-mode sign-in, with
  test payments, before the owner is offered live.
- **Rate:** Stripe allows about 100 calls a second in live mode and 25 in
  test mode. A refund job never comes near it; a job reading thousands of
  payments should page.
- **A refund can't be undone by a taf.** Stripe allows cancelling only a
  few kinds of refund, mostly from its own dashboard, so treat every refund
  as final.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused: several accounts, or none | The sign-in reaches more than one Stripe account | "Your Stripe sign-in reaches several accounts. Sign in again and pick the one this taf should use." |
| Refused: a live object | The refund entry is test mode only | "This taf only works with test payments for now. Nothing was refunded." |
| Refused: another currency | The Card caps refunds in euros | "That payment isn't in euros, so this taf can't refund it. Refund it in Stripe directly." |
| Refused: more than was paid, or already refunded | Stripe checks what's left to refund | "That payment was already refunded, or the amount is more than what's left. Nothing new was refunded." |
| Not found | A wrong id, or a payment in the other mode | "I can't find that payment in this Stripe account. Check the id, or whether it's a test payment." |
| Too many calls (429) | Over Stripe's rate | "Stripe asked us to slow down. The taf will try again in a moment." |
| The sign-in no longer works (401) | The owner disconnected it in Stripe | "Stripe needs you to sign in again. Nothing was refunded." |
| An invoice isn't found | A wrong id, or an invoice in the other mode | "I can't find that invoice in this Stripe account. Check it, or whether it's a test invoice." |
| The invoice is no longer open | It was paid, voided or written off since the list | "That invoice is already paid or closed, so no reminder went out." |
| Charged automatically | Its collection method isn't `send_invoice`: don't send | "Stripe charges that customer's card automatically, so I didn't send a reminder for it." |
| Sent, but no email arrived | A test-mode invoice | "That's a test invoice, so Stripe didn't email anyone." |

After any refusal the taf stops that step and says so. A refused refund is
never tried again with a different amount.
