# Holded

Checked: 2026-11-14 · Sources: the connection service's Holded tool
definitions as planque reads them (10 tools), and Holded's API reference
(invoices, contacts, error handling and API limits).

## How it connects

- **By the owner's API key**, made in Holded's developer settings and pasted
  once. The connection service holds it; the taf never sees it.
- **A key can be limited** to some kinds of record (contacts, sales
  invoices) and to reading or writing them. A key made for reading doesn't
  let the taf create or change anything.
- **With a known sign-in**, planque's catalog puts the tools the key can't
  run under `unusable`. Ask for every tool the job needs in the first
  request, so nothing is asked twice.

## Tools a job usually needs

Read means it changes nothing in Holded. The catalog gives no `returns`
paths for Holded today; the ids below are those Holded's answers carry.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `HOLDED_FIND_CONTACTS` | Finds contacts by exact email, phone or tax id, or by the start of a name, a page at a time | read | contact ids, a page cursor |
| `HOLDED_GET_CONTACT` | Reads one contact | read | the contact's id, name, email, tax id |
| `HOLDED_FIND_INVOICES` | Lists sales invoices by contact, date, due date or number, a page at a time | read | invoice ids, document numbers, status, a page cursor |
| `HOLDED_GET_INVOICE` | Reads one invoice: lines, totals, due date, payment status | read | the invoice's id, its document number, its contact id |
| `HOLDED_CREATE_CONTACT` | Adds a customer or supplier | write | the new contact's id |
| `HOLDED_UPDATE_CONTACT` | Changes a contact's details | write | the contact's id |
| `HOLDED_CREATE_INVOICE` | Writes a sales invoice for an existing contact | write | the new invoice's id |
| `HOLDED_UPDATE_INVOICE` | Changes an invoice | write | the invoice's id |

**What to keep:** each contact's id, since an invoice can only be made for
an existing contact; each invoice's id, to read or change that same invoice
later. Look a contact up by exact email or tax id before adding one, so the
taf doesn't make a second copy.

**A new invoice is a draft** unless the account approves new documents
automatically, and an approved invoice can't be deleted. "Prepare this
month's invoices for me to check" fits; anything that approves or sends
invoices is the owner's step.

**Updating an invoice replaces it whole** on Holded's side: a field left out
goes back to its default. Read the invoice just before changing it, and
check the answer after.

## Don't offer

- `HOLDED_DELETE_INVOICE`: deleting removes the invoice and any payments
  recorded on it for good, and its number isn't reused.
- `HOLDED_DELETE_CONTACT`: it removes a customer or supplier the owner's
  documents name.
- Setting your own invoice numbers: Holded warns that a number shaped like
  an active series breaks that series' order. Leave numbering to Holded.

## Access and plans

- **Calls are counted for the whole account,** across every key, by minute
  and by month. Holded's plan table: 60 calls a minute on its Starter,
  Growth and Basic plans, 120 on Standard and Advanced, 600 on Premium.
  Monthly budgets run from 500 calls (Starter) to 100,000 (Premium). Below
  Premium, calls are refused once the month's budget is used. Holded's
  error page gives a different per-minute figure (100 per key), so
  treat these numbers as a guide.
- **Lists come a page at a time:** up to 200 invoices or 100 contacts a page.
  A taf follows the cursor rather than reading everything each run.
- **A job that runs every day on a small plan** can use up the month's
  budget. Read only what changed since the last run.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403) | The key doesn't allow this kind of record, or only allows reading it | "Holded is only partly connected: this step needs access your key doesn't give. Want to add a key that allows it, or do without this step?" |
| The key doesn't work (401) | The key is missing, wrong or was removed in Holded | "Holded needs a new key. Nothing was changed." |
| Too many calls (429), this minute | Over the per-minute limit | "Holded asked us to slow down. The taf will try again in a minute." |
| Too many calls (429), this month | The plan's monthly budget is used up | "Holded's monthly limit for your plan is reached. The taf will pick up again next month, or sooner if you move to a bigger plan." |
| Rejected (400 or 422) | A value doesn't fit, or the invoice is approved and can't be deleted | "Holded didn't accept that change. The taf has stopped and will show you what it tried." |
| Not found (404) | The contact or invoice was deleted | "That isn't in Holded any more. The taf will skip it." |
| Holded unavailable (500) | Trouble on Holded's side | "Holded isn't answering right now. The taf will try again later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write: a taf that can't find a contact asks before adding one.
