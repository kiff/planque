# Salesforce

Checked: 2026-11-13 · Sources: the connection service's Salesforce tool
definitions as planque reads them, and Salesforce's help articles on API
access by edition and on the daily request limit. Salesforce's developer
reference couldn't be read for this check, so what each tool returns is
marked "not checked".

## How it connects

- **A sign-in to Salesforce**, held by the connection service; the taf never
  sees a password or token. The sign-in acts as one Salesforce user, so it
  sees and changes only what that user may. planque's catalog lists the
  tools and, with a known sign-in, puts the ones its access doesn't cover
  under `unusable`, with why.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.
- **API access depends on the Salesforce edition.** Enterprise, Unlimited,
  Developer and Performance include it. Professional has it only if the
  company bought it as an add-on. Group and Essentials can't have it at
  all. Ask which edition the owner has before offering a Salesforce job.

## Tools a job usually needs

Read means it changes nothing in Salesforce. The catalog gives no `returns`
paths for Salesforce today.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `SALESFORCE_GET_USER_INFO` | Reads the signed-in user | read | not checked |
| `SALESFORCE_SEARCH_CONTACTS` | Finds contacts by name, email, phone or company | read | not checked |
| `SALESFORCE_SEARCH_ACCOUNTS` | Finds companies by name, city or industry | read | not checked |
| `SALESFORCE_SEARCH_LEADS` | Finds leads by name, email, company or status | read | not checked |
| `SALESFORCE_SEARCH_OPPORTUNITIES` | Finds deals by company, amount, close date, open or won | read | not checked |
| `SALESFORCE_GET_OPPORTUNITY` | Reads one deal | read | not checked |
| `SALESFORCE_SEARCH_TASKS` | Finds follow-ups by contact, company, due date or done | read | not checked |
| `SALESFORCE_CREATE_CONTACT` | Adds a contact | write | not checked |
| `SALESFORCE_CREATE_LEAD` | Adds a lead | write | not checked |
| `SALESFORCE_CREATE_TASK` | Adds a follow-up task | write | not checked |
| `SALESFORCE_LOG_CALL` | Records a call that happened | write | not checked |
| `SALESFORCE_CREATE_NOTE` | Adds a note to a record | write | not checked |
| `SALESFORCE_UPDATE_OPPORTUNITY` | Changes a deal: stage, amount, close date, next step | write | not checked |

**What to keep:** each record's Salesforce id from the answer, to update
that same record next time. Never look a record up again by its name: use
the id the taf kept.

**Search first.** Before adding a contact or a lead, search by email. Add
one only when nothing matches, so no one is entered twice.

**Start read-only.** "Tell me which open deals close this month" needs only
the read tools. Tasks, call logs and notes are small writes. Changing a
deal's stage or amount changes the team's forecast, so offer it only when
asked, and each change waits for the owner's yes.

## Don't offer

- Deleting (`SALESFORCE_DELETE_CONTACT`, `SALESFORCE_DELETE_LEAD`,
  `SALESFORCE_DELETE_ACCOUNT`, `SALESFORCE_DELETE_OPPORTUNITY` and others).
- Bulk changes: `SALESFORCE_MASS_TRANSFER_OWNERSHIP`,
  `SALESFORCE_UPDATE_S_OBJECT_COLLECTION`, and the bulk job tools. One call
  changes many records at once.
- Setup and admin: custom fields and objects (`SALESFORCE_CREATE_CUSTOM_FIELD`,
  `SALESFORCE_CREATE_CUSTOM_OBJECT`), validation rules, list views, folders
  and sharing, flows and tooling records. They change Salesforce for the
  whole company.
- Free-form queries (`SALESFORCE_RUN_SOQL_QUERY`, `SALESFORCE_QUERY_ALL`)
  when a search tool above covers the job: the search tools say plainly
  what they look for.
- Generic record tools (`SALESFORCE_CREATE_A_RECORD`,
  `SALESFORCE_UPDATE_RECORD`) when a named tool exists for that kind of
  record.

## Access and plans

- **Edition first:** see How it connects. On an edition without API access,
  every call is refused, whatever the sign-in allows.
- **A daily request allowance, per company,** shared by every connected app.
  It grows with the number of licenses. A company that goes far over it may
  have every app's calls refused for up to a day.
- **A job that runs often** should search for what changed since its last
  run, not read everything again.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused, not allowed | The signed-in user's permissions don't cover this record or action | "Salesforce doesn't let your account do this step. Your Salesforce admin can allow it, or the taf can do without it." |
| No API access | The company's Salesforce edition doesn't include it | "Your Salesforce edition doesn't let outside apps connect. Your Salesforce admin can tell you whether it can be added." |
| Request limit exceeded (REQUEST_LIMIT_EXCEEDED) | The company used up its daily allowance | "Salesforce's daily limit for your company is used up. The taf will try again tomorrow." |
| The sign-in no longer works | Access was removed, or the sign-in ended | "Salesforce needs you to sign in again. Nothing was changed." |
| A record isn't found | It was deleted or merged | "That record isn't in Salesforce any more. The taf will skip it." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write: a taf that can't find a contact never adds them again.
