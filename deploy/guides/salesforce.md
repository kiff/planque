# Salesforce

Checked: 2026-11-13 · Sources: the connection service's Salesforce tool
definitions as planque reads them (288 tools), and Salesforce's developer
docs: REST API query and error references, supported editions, and API
request limits and allocations.

## How it connects

- **A sign-in to the owner's Salesforce org**, held by the connection
  service; the taf never sees a password or token. The taf acts as the user
  who signed in, and sees only the records that user may see. With a known
  sign-in, the catalog puts the tools its access doesn't cover under
  `unusable`, with why.
- **Ask for every tool in the first request.** A held sign-in is reused for
  a new taf only when it covers the tools asked.
- **There are hundreds of tools; a job needs a handful.** Offer the ones
  below, by name.

## Tools a job usually needs

Read changes nothing in the org. Returns are what Salesforce's reference
documents: every record has an 18-character `Id`; a create answers with the
new `Id`; a query answers `records`, `totalSize`, `done` and, when more
remain, a `nextRecordsUrl`.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `SALESFORCE_RUN_SOQL_QUERY` | Runs a query over any object, up to 2,000 records a page | read | record ids, `nextRecordsUrl` |
| `SALESFORCE_QUERY_MORE` | Fetches the next page of a query | read | record ids, `nextRecordsUrl` |
| `SALESFORCE_SEARCH_ACCOUNTS` / `SALESFORCE_SEARCH_CONTACTS` | Finds accounts or contacts by name, email, and so on | read | account or contact ids |
| `SALESFORCE_SEARCH_LEADS` / `SALESFORCE_SEARCH_OPPORTUNITIES` | Finds leads or opportunities | read | lead or opportunity ids |
| `SALESFORCE_GET_ACCOUNT` / `SALESFORCE_GET_CONTACT` / `SALESFORCE_GET_OPPORTUNITY` / `SALESFORCE_GET_LEAD` | Reads one record | read | the record and its `Id` |
| `SALESFORCE_GET_ALL_FIELDS_FOR_OBJECT` | Lists an object's fields, including custom ones | read | field names and types |
| `SALESFORCE_LIST_REPORTS` / `SALESFORCE_RUN_REPORT` | Lists saved reports, or runs one | read | report ids, the report's rows |
| `SALESFORCE_CREATE_LEAD` | Adds a lead | write | the new lead `Id` |
| `SALESFORCE_CREATE_TASK` | Adds a follow-up task to a record | write | the new task `Id` |
| `SALESFORCE_LOG_CALL` / `SALESFORCE_CREATE_NOTE` | Records a call or a note on a record | write | the new task or note `Id` |
| `SALESFORCE_UPDATE_OPPORTUNITY` / `SALESFORCE_UPDATE_CONTACT` | Changes fields of one record | write | nothing beyond success |

**What to keep:** the `Id` of each record the taf made or follows. Never
search for a record the taf made: use the id it kept. Read
`SALESFORCE_GET_ALL_FIELDS_FOR_OBJECT` once before writing: orgs add
required custom fields.

**Start read-only.** "Brief me on this week's pipeline" is queries and
reports only. Offer a creating or updating tool only when the owner asked;
each change is shown to them before it is written. Sending email
(`SALESFORCE_SEND_EMAIL`, `SALESFORCE_SEND_EMAIL_FROM_TEMPLATE`) goes to
customers: only on request, and each send waits for their yes.

## Don't offer

- Deleting: `SALESFORCE_DELETE_*`.
- Bulk and raw-request tools: `SALESFORCE_EXECUTE_COMPOSITE`,
  `SALESFORCE_EXECUTE_COMPOSITE_BATCH`, `SALESFORCE_POST_COMPOSITE_GRAPH`,
  `SALESFORCE_POST_COMPOSITE_SOBJECTS`, `SALESFORCE_PATCH_COMPOSITE_SOBJECTS`,
  `SALESFORCE_UPDATE_S_OBJECT_COLLECTION`, `SALESFORCE_CREATE_SOBJECT_TREE`,
  `SALESFORCE_CREATE_BULK_INGEST_JOB`, `SALESFORCE_UPLOAD_JOB_DATA`,
  `SALESFORCE_MASS_TRANSFER_OWNERSHIP`, `SALESFORCE_SEND_MASS_EMAIL`, and
  the any-object writers `SALESFORCE_CREATE_A_RECORD`,
  `SALESFORCE_UPDATE_RECORD`, `SALESFORCE_UPSERT_SOBJECT_BY_EXTERNAL_ID`.
- Admin and org setup: `SALESFORCE_CREATE_CUSTOM_OBJECT`,
  `SALESFORCE_CREATE_CUSTOM_FIELD`, `SALESFORCE_UPDATE_CUSTOM_FIELD`,
  `SALESFORCE_CREATE_VALIDATION_RULE`, `SALESFORCE_START_METADATA_DEPLOYMENT`,
  the developer tools (`SALESFORCE_CREATE_TOOLING_RECORD`,
  `SALESFORCE_UPDATE_TOOLING_RECORD`, `SALESFORCE_CREATE_TOOLING_FLOW`,
  `SALESFORCE_UPDATE_TOOLING_FLOW`, `SALESFORCE_EXECUTE_TOOLING_COMPOSITE`,
  `SALESFORCE_RUN_TESTS_SYNCHRONOUS`, `SALESFORCE_RUN_TESTS_ASYNCHRONOUS`).
- People, passwords and sharing: `SALESFORCE_SET_USER_PASSWORD`,
  `SALESFORCE_SET_SELF_SERVICE_USER_PASSWORD`,
  `SALESFORCE_RESET_SELF_SERVICE_USER_PASSWORD`,
  `SALESFORCE_CREATE_FOLDER_SHARES`, `SALESFORCE_REPLACE_FOLDER_SHARES`,
  `SALESFORCE_UPDATE_FOLDER_SHARE`.
- Automation that acts for others: `SALESFORCE_PROCESS_APPROVALS`,
  `SALESFORCE_TRIGGER_PROCESS_RULES`, `SALESFORCE_APPLY_LEAD_ASSIGNMENT_RULES`,
  `SALESFORCE_EXECUTE_INVOCABLE_ACTION`.

## Access and plans

- **Editions:** Enterprise, Unlimited, Performance and Developer include API
  access. Professional has API access only with Salesforce's add-on, so check
  the edition first.
- **The user needs the "API Enabled" permission** in their profile or a
  permission set, set by the org's admin.
- **Daily calls are shared by the whole org,** over a rolling 24 hours. On
  Enterprise, Salesforce's own example is 100,000 plus 1,000 per user
  licence. Past a point, every call is refused until usage drops.
- **Long calls:** at most 25 calls running 20 seconds or more at once on a
  production org, 5 on a Developer or trial org.
- **Pages:** a query returns up to 2,000 records a page.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The sign-in no longer works (401) | The session expired or access was removed | "Salesforce needs you to sign in again. Nothing was changed." |
| Refused, API not enabled (403) | The edition lacks API access, or the user lacks the permission | "Your Salesforce plan or your user doesn't allow outside apps. Your Salesforce admin can turn this on." |
| Refused, request limit exceeded (403) | The org's daily calls are used up, or too many long calls at once | "Salesforce's limit for your org is reached. The taf will continue later; your admin can see usage in Setup." |
| Refused for missing rights (insufficient access) | The user can't see or edit that record | "Your Salesforce user can't change that record. Someone with access, or your admin, can share it." |
| A field or validation rule refuses the write (400) | A required field is missing, or an org rule blocks the value | "Salesforce didn't accept that change: a rule in your org blocks it. Nothing was written." |
| Not found (404), or the record was deleted | It was deleted or merged | "That record isn't in Salesforce any more. The taf will skip it." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
