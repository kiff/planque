# Supabase

Checked: 2026-11-13 · Sources: the connection service's Supabase tool
definitions as planque reads them (230 tools, [generated](generated/supabase.md)),
and Supabase's management API reference.

## How it connects

- **A sign-in to the owner's Supabase account**, held by the connection
  service. The sign-in reaches **every project in the owner's
  organisations**, with the owner's own rights: offer it only to an owner
  who said they use Supabase, and name the one project the job uses.
- **One approval for the sign-in, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `SUPABASE_LIST_ALL_PROJECTS` | Lists the owner's projects | read | **the project's ref** |
| `SUPABASE_LIST_TABLES` | Lists a project's tables | read | table names |
| `SUPABASE_GET_TABLE_SCHEMAS` | Reads tables' columns and types | read | |
| `SUPABASE_SELECT_FROM_TABLE` | Reads rows, with filters, order and a limit | read | row ids |
| `SUPABASE_COUNT_TABLE_ROWS` | Counts rows matching a filter | read | |
| `SUPABASE_RUN_READ_ONLY_QUERY` | Runs a SQL query that can only read | read | |
| `SUPABASE_INSERT_TABLE_ROWS` | Adds rows | write | the new rows' ids |
| `SUPABASE_UPDATE_TABLE_ROWS` | Changes rows matching a filter | write | |

**Prefer the read-only query** to the general one: it can't change anything,
whatever the query says. A job that writes names the table and the columns
on the Card.

**An update needs a filter.** `SUPABASE_UPDATE_TABLE_ROWS` with no filter
changes every row; the taf always names the row by its id.

**What to keep:** the project ref and table name; each row's id the taf
adds.

## Don't offer

- The general SQL tool (`SUPABASE_BETA_RUN_SQL_QUERY`), migrations
  (`SUPABASE_APPLY_A_MIGRATION`, `SUPABASE_UPSERT_MIGRATION`) and
  `SUPABASE_DELETE_TABLE_ROWS`: they can change or drop the owner's data.
- Project administration: creating or deleting projects, keys, secrets,
  network rules, replicas, custom domains, Postgres upgrades, add-ons (which
  cost money).
- Edge functions and auth settings.

## Access and plans

- **Free projects pause after a week without use;** a paused project answers
  nothing until the owner restores it in Supabase.
- **The management API is limited** to about 120 calls a minute per account.
- **Row-level rules** in the owner's database still apply to what the taf
  reads, depending on the role the tools use.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused (403) | The owner's role in that organisation doesn't allow it | "Your Supabase role doesn't let me do that in this project." |
| The project doesn't answer (503, paused) | A free project paused for inactivity | "Your Supabase project is paused. Restore it in Supabase, then I'll try again." |
| A query is refused (400) | A table or column that doesn't exist, or SQL that can't run | "The query didn't run: the table may have changed. Want me to read it again?" |
| Too many calls (429) | Over the account's rate | "Supabase asked us to slow down. The taf will continue in a minute." |
| The sign-in no longer works (401) | Removed, or expired | "Supabase needs you to sign in again." |

After any refusal the taf stops that step and says so.
