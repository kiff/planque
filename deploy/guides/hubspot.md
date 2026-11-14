# HubSpot

Checked: 2026-11-14 · Sources: the connection service's HubSpot tool
definitions as planque reads them (352 tools), and HubSpot's developer
documentation (usage guidelines, CRM search, contacts, notes, error handling).

## How it connects

- **A sign-in to HubSpot**, held by the connection service; the taf never
  sees a password or token. The owner picks the HubSpot account to connect
  on HubSpot's own page. planque's catalog lists the tools and, with a known
  sign-in, puts the ones its access doesn't cover under `unusable`, with why.
- **Access is per kind of record.** Contacts, companies and deals each have
  their own read and write access; owners can only be read. A sign-in that reads deals
  doesn't necessarily read contacts. Ask for every tool the job needs in the
  first request, so nothing is asked twice.
- **A held sign-in is reused** for a new taf when it covers the tools asked.

## Tools a job usually needs

Read means it changes nothing in the account. The catalog's `returns` gives
exact paths where it has them; where it has none yet, the answer still
carries the record's id, which is what a taf keeps.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `HUBSPOT_SEARCH_CONTACTS_BY_CRITERIA` | Finds contacts by email, name, or any property | read | contact ids, a paging cursor |
| `HUBSPOT_READ_CONTACT` | Reads one contact, with chosen properties and linked records | read | the contact and its id |
| `HUBSPOT_SEARCH_COMPANIES` | Finds companies by domain, name, or any property | read | company ids, a paging cursor |
| `HUBSPOT_GET_COMPANY` | Reads one company | read | the company and its id |
| `HUBSPOT_SEARCH_DEALS` | Finds deals by stage, amount, close date or owner | read | deal ids, a paging cursor |
| `HUBSPOT_GET_DEAL` | Reads one deal | read | the deal and its id |
| `HUBSPOT_RETRIEVE_ALL_PIPELINES_FOR_SPECIFIED_OBJECT_TYPE` | Lists the deal pipelines and their stages | read | pipeline and stage ids |
| `HUBSPOT_RETRIEVE_OWNERS` | Lists the people records can be assigned to | read | owner ids |
| `HUBSPOT_LIST_CONTACT_NOTES` | Lists the notes on one contact | read | note ids |
| `HUBSPOT_CREATE_CONTACT` | Adds a contact | write | the new contact's id |
| `HUBSPOT_UPDATE_CONTACT` | Changes properties on one contact | write | the contact's id |
| `HUBSPOT_UPDATE_DEAL` | Moves a deal's stage or changes its properties | write | the deal's id |
| `HUBSPOT_CREATE_NOTE` | Adds a note to the timeline of a contact, company or deal | write | the new note's id |
| `HUBSPOT_CREATE_TASK` | Adds a to-do for a person on the team, linked to a record | write | the new task's id |

**What to keep:** each record's id (contact, company, deal), to read or
update that same record next run; the ids of notes and tasks the taf made,
so it never adds the same one twice. Look up a pipeline's stage ids once and
keep them: a deal's stage is set by id, not by its name.

**Search before creating a contact.** HubSpot treats an email address as a
contact's unique key, so search by email first and update the contact found.
Notes and tasks need a time (when they happened, or when they're due).

**Start read-only.** "Tell me which deals went quiet" is search and read
tools only. Offer a write only for what the owner asked to change.

## Don't offer

- Archiving, merging and deleting (`HUBSPOT_ARCHIVE_*`,
  `HUBSPOT_MERGE_*`, `HUBSPOT_DELETE_*`,
  `HUBSPOT_PERMANENTLY_DELETE_CONTACT_VIA_GDPR`, `HUBSPOT_REMOVE_DEAL`):
  merges and permanent deletes can't be undone.
- Sending, publishing or scheduling (`HUBSPOT_PUBLISH_*`,
  `HUBSPOT_SCHEDULE_*`) and changing email opt-ins
  (`HUBSPOT_UPDATE_EMAIL_SUBSCRIPTIONS`): they reach the owner's customers.
- Imports (`HUBSPOT_START_IMPORT`, `HUBSPOT_IMPORT_FILE_FROM_URL`), users
  (`HUBSPOT_REMOVE_SETTINGS_USER`), and calling or video settings.
- Enrolling contacts in sequences or workflows
  (`HUBSPOT_ENROLL_CONTACT_IN_SEQUENCE`, `HUBSPOT_ENROLL_CONTACT_IN_WORKFLOW`):
  they set off emails and automations the taf doesn't control.
- Pipeline and schema changes (creating, replacing or deleting pipelines,
  stages, properties or association labels): they change the account for the
  whole team.
- The "from natural language" create tools and the legacy contact-profile
  tools: use the plain create and read tools above.
- Batch creates and updates, unless the job really handles many records at
  once; one at a time is easier for the owner to check.

## Access and plans

- **Free HubSpot accounts can use the API.** Which features need a paid
  plan isn't checked here.
- **Rate:** an app connected to an account gets about 110 calls every ten
  seconds. Search is stricter, at five calls a second per account, and
  pages of at most 200 results, up to 10,000 for one search. A bigger job
  narrows its search (by date, say) rather than paging further.
- **Batches** take at most 100 records at a time.
- **Busy accounts:** when many records change at once, HubSpot can lock the
  account for about two seconds. Wait and try again.
- **A record the owner deletes in HubSpot** goes to its recycle bin, where
  the owner can restore it.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403) | The sign-in didn't give this kind of record, or this write | "HubSpot is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Too many calls (429) | Over the account's short-term limit | "HubSpot asked us to slow down. The taf will try again in a minute." |
| Account locked briefly (423) | Many records changed at once | "HubSpot is busy with a lot of changes. The taf will try again in a moment." |
| The sign-in no longer works (401) | The owner removed the app, or the sign-in expired | "HubSpot needs you to sign in again. Nothing was changed." |
| A contact with that email already exists | Email is a contact's unique key | "That person is already in HubSpot. The taf will update their record instead of adding a new one." |
| A record isn't found (404) | It was deleted or merged in HubSpot | "That record isn't in HubSpot any more. The taf will skip it." |
| A search asks past 10,000 results (400) | Too many matches for one search | "That search finds too much at once. The taf will look at a shorter period." |
| HubSpot is unavailable (5xx) | A short outage on HubSpot's side | "HubSpot isn't answering right now. The taf will try again later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a create or a write.
