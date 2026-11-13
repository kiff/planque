# Zoho CRM

Checked: 2026-11-13 · Sources: the connection service's Zoho tool
definitions as planque reads them (110 tools, [generated](generated/zoho.md)),
and Zoho CRM's API reference.

## How it connects

- **A sign-in to Zoho**, held by the connection service, choosing the data
  centre the owner's account lives in (EU accounts are on `zoho.eu`). The
  tools are Zoho CRM's: leads, contacts, accounts, deals.
- The catalog lists the tools and puts the ones its access doesn't cover
  under `unusable`. Offer only what's under `tools`.
- **One approval for the sign-in, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `ZOHO_SEARCH_ZOHO_RECORDS` | Finds records in a module by email, phone or criteria | read | **record ids** |
| `ZOHO_GET_ZOHO_RECORDS` | Lists a module's records | read | record ids |
| `ZOHO_GET_RECORD` | Reads one record | read | |
| `ZOHO_GET_MODULE_FIELDS` | Lists a module's fields, with their API names | read | field API names |
| `ZOHO_CREATE_LEAD` / `ZOHO_CREATE_CONTACT` / `ZOHO_CREATE_DEAL` | Creates a record | write | **the record's id** |
| `ZOHO_UPDATE_ZOHO_RECORD` | Changes fields of a record | write | |
| `ZOHO_CREATE_NOTE` | Adds a note to a record | write | the note's id |
| `ZOHO_ADD_TAGS_TO_RECORDS` | Tags records | write | |

**Search before creating, by email.** A search with no match is an empty
answer, not an error; only then create. Keep the id it answers.

**Fields have API names** (`Last_Name`, not "Last name"), which an owner's
custom fields change. Read the module's fields once and keep them.

## Don't offer

- Deleting records, notes, tags or drafts.
- `ZOHO_CONVERT_ZOHO_LEAD` unless the job is about converting: it turns a lead
  into a contact, an account and a deal at once.
- Workflow rules, functions, blueprints and assignment rules: they change
  what happens to every record.
- Email drafts sent from Zoho, unless the job is about emailing.

## Access and plans

- **The free edition has no API access** for most modules; paid editions do.
- **Calls are counted per day,** by edition and number of licences (from a
  few thousand). A job reading every record daily should read only changes.
- **Rate:** Zoho also limits calls at once; pages answer up to 200 records.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused (403, no permission) | The owner's profile can't do this, or the edition has no API | "Your Zoho edition or role doesn't allow this." |
| Wrong data centre (401 `invalid_oauth_token`) | Signed in to another region | "Zoho needs you to sign in again, choosing your region (EU)." |
| Refused (400, invalid data or mandatory field missing) | A field's API name or value is wrong | Read the fields and try once; else "Zoho didn't accept the record." |
| The day's calls are used (429) | Over the daily allowance | "Zoho's daily call allowance is used up. The taf continues tomorrow." |
| A duplicate (`DUPLICATE_DATA`) | A record with that email exists | Use the existing record's id; don't tell the owner. |

After any refusal the taf stops that step and says so.
