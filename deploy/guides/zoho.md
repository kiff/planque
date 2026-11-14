# Zoho CRM

Checked: 2026-11-14 · Sources: the connection service's Zoho tool
definitions as planque reads them (110 tools, [generated](generated/zoho.md)),
and Zoho CRM's API documentation (limits, status codes, data centres,
search, records).

## How it connects

- **A sign-in to Zoho**, held by the connection service. Each account's data
  lives in one data centre (EU accounts are on `zoho.eu`), and every call
  goes to that centre's address. The tools are Zoho CRM's: leads, contacts,
  accounts, deals.
- The catalog lists the tools and puts the ones its access doesn't cover
  under `unusable`. Offer only what's under `tools`.
- **One approval for the sign-in, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `ZOHO_SEARCH_ZOHO_RECORDS` | Finds records in a module by email, phone or criteria | read | **record ids** |
| `ZOHO_GET_ZOHO_RECORDS` | Lists a module's records | read | record ids |
| `ZOHO_GET_RECORD` | Reads one record | read | the record, with its id |
| `ZOHO_GET_MODULE_FIELDS` | Lists a module's fields, with their API names | read | field API names |
| `ZOHO_CREATE_LEAD` / `ZOHO_CREATE_CONTACT` / `ZOHO_CREATE_DEAL` | Creates a record | write | **the record's id** |
| `ZOHO_UPDATE_ZOHO_RECORD` | Changes fields of a record | write | not checked here |
| `ZOHO_CREATE_NOTE` | Adds a note to a record | write | the note's id |
| `ZOHO_ADD_TAGS_TO_RECORDS` | Tags records | write | not checked here |

**Search before creating, by email.** A search with no match answers with
no content, not an error; only then create. Keep the id it answers. A record
created or changed a moment ago may not show in a search yet: use the id
the taf kept.

**Fields have API names** (`Last_Name`, not "Last name"), which an owner's
custom fields change. Read the module's fields once and keep them.

## Don't offer

- Deleting records, notes, tags or drafts.
- `ZOHO_CONVERT_ZOHO_LEAD` unless the job is about converting: it turns a lead
  into a contact, an account and a deal at once, and costs 5 calls' worth of
  credits.
- Workflow rules, functions, blueprints and assignment rules: they change
  what happens to every record.
- Email drafts sent from Zoho, unless the job is about emailing.

## Access and plans

- **Calls are counted in credits over a rolling 24 hours,** usually one a
  call: 5,000 on the free edition, and from 50,000 plus an amount per user
  licence on paid editions, up to each edition's cap. A trial counts like
  its paid edition. A job reading every record daily should read only
  changes.
- **Rate:** Zoho also limits calls at once. Lists answer up to 200 records a
  page, a search at most 2,000 in all, and one create takes up to 100
  records.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused (403, no permission) | The owner's account isn't allowed to do this | "Your Zoho account doesn't allow this step. Your Zoho admin can check its permissions." |
| The sign-in isn't accepted (401) | The sign-in is invalid or expired | "Zoho needs you to sign in again." |
| Refused (400, invalid data or mandatory field missing) | A field's API name or value is wrong | Read the fields and try once; else "Zoho didn't accept the record." |
| Too many calls (429) | Over the 24-hour credits, or too many calls at once | "Zoho's call allowance is used up for now. The taf will continue later." |
| A duplicate (`DUPLICATE_DATA`) | A field marked unique already has that value | Use the existing record's id; don't tell the owner. |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
