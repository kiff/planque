# Supabase

Checked: 2026-11-14 · Sources: the connection service's Supabase tool
definitions as planque reads them (230 tools, [generated](generated/supabase.md)),
and Supabase's management API reference.

## How it connects

- **A sign-in to the owner's Supabase account**, held by the connection
  service. The owner picks an organisation at sign-in, and the sign-in gets
  **full management access to its projects** on the owner's behalf: offer
  it only to an owner who said they use Supabase, and name the one project
  the job uses.
- **One approval for the sign-in, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `SUPABASE_LIST_ALL_PROJECTS` | Lists the owner's projects | read | **the project's ref** |
| `SUPABASE_LIST_TABLES` | Lists a project's tables | read | table names |
| `SUPABASE_GET_TABLE_SCHEMAS` | Reads tables' columns and types | read | not checked here |
| `SUPABASE_SELECT_FROM_TABLE` | Reads rows, with filters, order and a limit | read | row ids |
| `SUPABASE_COUNT_TABLE_ROWS` | Counts rows matching a filter | read | not checked here |
| `SUPABASE_RUN_READ_ONLY_QUERY` | Runs a SQL query that can only read | read | not checked here |
| `SUPABASE_INSERT_TABLE_ROWS` | Adds rows | write | the new rows' ids |
| `SUPABASE_UPDATE_TABLE_ROWS` | Changes rows matching a filter | write | not checked here |

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

- **Free projects pause after a week of low activity;** a paid plan's
  projects don't. The owner can restore a paused project from Supabase's
  dashboard, for up to a year after it paused.
- **The management API allows 120 calls a minute** per user, counted
  separately for each project, organisation and endpoint. Over it, calls are
  refused until the minute ends.
- **Row-level rules** in the owner's database apply or not depending on the
  role the tools use: the service role bypasses them.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused | This sign-in isn't allowed to do that in Supabase | "Supabase doesn't let this connection do that in this project." |
| The project doesn't answer | A free project paused after a week of low activity | "Your Supabase project is paused. Restore it in Supabase, then I'll try again." |
| A query is refused | A table or column that doesn't exist, or SQL that can't run | "The query didn't run: the table may have changed. Want me to read it again?" |
| Too many calls (429) | Over 120 calls a minute on one endpoint | "Supabase asked us to slow down. The taf will continue in a minute." |
| The sign-in no longer works | The sign-in isn't accepted any more | "Supabase needs you to sign in again." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
