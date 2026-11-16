# Google BigQuery

Checked: 2026-11-14 · Sources: the connection service's BigQuery tool
definitions as planque reads them (63 tools, [generated](generated/googlebigquery.md)),
and Google's BigQuery reference.
Failure table checked against planque's answers on 2026-11-16.

## How it connects

- **A sign-in to Google**, held by the connection service. The grant's catalog lists
  the tools and, with a known sign-in, puts the ones its access doesn't cover
  under `unusable`. Offer only what's under `tools`.
- **The owner needs a Google Cloud project** with BigQuery in it, and a role
  that lets them run queries in it. Most people don't have one: offer
  BigQuery only to an owner who said they use it.

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `GOOGLEBIGQUERY_LIST_PROJECTS` | Lists the projects the owner can use | read | **the project id** |
| `GOOGLEBIGQUERY_LIST_DATASETS` | Lists a project's datasets | read | dataset ids |
| `GOOGLEBIGQUERY_LIST_TABLES` | Lists a dataset's tables | read | table ids |
| `GOOGLEBIGQUERY_GET_TABLE_SCHEMA` | Reads a table's columns and types | read | |
| `GOOGLEBIGQUERY_QUERY` | Runs a SQL query and answers its rows | **counted as a write** | the job's id |
| `GOOGLEBIGQUERY_GET_QUERY_RESULTS` | Reads the rows of a query that took long | read | |
| `GOOGLEBIGQUERY_LIST_TABLE_DATA` | Reads a table's rows without a query | read | |

**Why a query counts as a write:** a query can also change or delete data
(`DELETE`, `UPDATE`, `CREATE`), and it costs money. The Card shows the query
itself, and a job that only reads says so in its instructions: `SELECT` only.

**What to keep:** the project, dataset and table ids; the query, so the same
one runs every time.

## Don't offer

- Creating, changing, deleting or restoring datasets, tables, models and routines, and deleting job history.
- Publishing data to others (the data exchange, listing and query template tools), connections to outside data, and data policies.
- `GOOGLEBIGQUERY_INSERT_JOB`: it runs any kind of job, exports and copies included. A query uses `GOOGLEBIGQUERY_QUERY`.
- Reservations and capacity commitments
  (`GOOGLEBIGQUERY_CREATE_CAPACITY_COMMITMENT`,
  `GOOGLEBIGQUERY_CREATE_RESERVATION`): they commit the owner to spending.
- Access policies (`GOOGLEBIGQUERY_SET_ROUTINE_IAM_POLICY` and the other
  policy tools).
- `GOOGLEBIGQUERY_INSERT_ALL` and uploads, unless the job is about loading
  data the owner named.

## Access and plans

- **Queries are billed by the data they read,** after a free allowance of
  1 TB a month. A query over a big table without a date filter can cost real
  money every run. The job's query filters on dates, and names its columns
  instead of `SELECT *`.
- **Results:** a query answers a page of rows; long answers are read in
  pages with `GOOGLEBIGQUERY_GET_QUERY_RESULTS`.
- **The project's admin** decides who can query; the owner's role is what
  the taf can do.

## When it fails

The rows are Google BigQuery's own answers: a refusal reaches the taf as
"Google BigQuery refused it (it answered …): …", with Google BigQuery's
reason. planque's own refusals, and an expired sign-in, are in the part
below.

| What the taf gets | Why | Tell the owner |
|---|---|---|
| Refused (403, access denied) | The owner's role doesn't allow it, or BigQuery isn't on in the project | "Your Google Cloud role doesn't let me run that. A project admin can give you query access." |
| Refused (403, quota or billing) | The free allowance is used and billing isn't set up, or a cost cap was hit | "BigQuery's allowance for this project is used up. Billing needs to be set up to run more queries." |
| Refused (400), bad query | A SQL error, or a column that doesn't exist | "The query didn't run: a table or column has changed. Want me to look at the table again?" |
| Refused (404), not found | A dataset or table renamed or deleted | "I can't find that table any more." |
| Refused (429), too many calls, or a 403 that says the rate is exceeded | Over the project's rate | "BigQuery asked us to slow down. The taf will try again in a minute." |
| Google BigQuery's sign-in has expired (401) | Removed, or expired | "BigQuery needs you to sign in again." |
