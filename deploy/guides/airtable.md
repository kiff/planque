# Airtable

Checked: 2026-11-13 · Sources: the connection service's published Airtable
tool list, Airtable's API reference (rate limits, errors, access scopes,
records) and its help pages on monthly call limits. Not yet run on a phone.
The tool names below are checked against the generated part once the keyed
run lands.

## How it connects

- **A sign-in to Airtable**, held by the connection service; the taf never
  sees a password or token. With a known sign-in, the catalog puts the tools
  its access doesn't cover under `unusable`, with why.
- **The owner picks the bases at sign-in.** The taf reaches only the bases
  and workspaces the owner chose then, and never more than the owner can do
  in them. A base added later needs the owner to add it to the sign-in.
- **A held sign-in is reused** for a new taf when it covers the tools and
  the bases asked. Ask for every tool the job needs in the first request.
- **A base holds tables; a table holds records.** Most jobs read the base's
  tables and fields first, then keep the table's id.

## Tools a job usually needs

Read means it changes nothing in Airtable. Returns are what a taf keeps to
use next time; the catalog's `returns` gives their exact paths.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `AIRTABLE_LIST_BASES` | Lists the bases this sign-in can reach | read | base ids and names |
| `AIRTABLE_GET_BASE_SCHEMA` | Lists a base's tables and their fields | read | table ids, field ids and types |
| `AIRTABLE_LIST_RECORDS` | Lists a table's records, filtered and sorted, a page at a time | read | record ids, their fields, a page offset |
| `AIRTABLE_GET_RECORD` | Reads one record | read | the record and its fields |
| `AIRTABLE_LIST_COMMENTS` | Reads the comments on a record | read | comment ids, text |
| `AIRTABLE_CREATE_RECORDS` | Adds one or more records to a table | write | the records' ids |
| `AIRTABLE_UPDATE_RECORD` | Changes some fields of one record, leaving the rest | write | the record's id |
| `AIRTABLE_UPDATE_MULTIPLE_RECORDS` | Changes several records; can also add the ones a key field doesn't match | write | the records' ids |
| `AIRTABLE_CREATE_COMMENT` | Adds a comment to a record | write | the comment's id |

**What to keep:** the base id, the table id, and the id of every record the
taf creates. Next run, update that record by its id. Prefer ids to names for
tables and fields: an owner can rename them, and the ids don't change.

**Start read-only.** "Tell me which leads have no next step" is reads only.
Offer adding or changing records only when the owner asked for it.

**Update, don't replace.** The tools that replace a whole record clear every
field the taf didn't send. Use the ones above, which change only the fields
named.

## Don't offer

- **Deleting tools** (records, comments): a taf has no way to bring back
  what it deleted.
- **Structure tools** (creating a base or table, adding or changing fields
  or tables): they change the base for everyone who uses it. If a job needs
  a new field, the owner adds it.
- **Replacing tools** (`AIRTABLE_UPDATE_RECORD_PUT`,
  `AIRTABLE_UPDATE_MULTIPLE_RECORDS_PUT`): they clear the fields not sent.
- **Record from free text** (`AIRTABLE_CREATE_RECORD_FROM_NATURAL_LANGUAGE`):
  the taf can't check which values go where before it writes. Use
  `AIRTABLE_CREATE_RECORDS` with named fields.
- **Event triggers** (base or view changed): they aren't tools a taf calls.

## Access and plans

- **Reading records needs record read access** plus reading the base's
  structure, which lists its tables and fields. Writing needs record write
  access. A sign-in made for reading doesn't let the taf write.
- **The owner's role in the base counts too.** Reading needs at least
  read-only access to the base; adding or changing records needs editor
  access. A table or field can also be locked to some people; then a write
  to it is refused even with the right sign-in.
- **Monthly call caps, per workspace:** about 1,000 calls a month on the
  free plan, and about 100,000 on the Team plan. Business and above have no
  monthly cap. On the free plan, going over starts a one-time 30-day grace
  period. On Team, going over slows calls down.
- **Rate:** about 5 calls a second per base, whatever the plan. Over it,
  Airtable refuses for 30 seconds.
- **Batches:** a call adds or changes up to 10 records. A job adding 200
  rows makes 20 calls, which counts against the monthly cap.
- **Big answers:** lists answer at most 100 records a page. A complex filter
  over a large table can time out; narrow it or ask for fewer per page.
- **An enterprise account's admin** can block outside apps. Then the list of
  bases comes back empty, with no error.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access | The sign-in didn't give this tool's access: reading only, say, and the tool adds records | "Airtable is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Refused, or "not found", on a base or table | The base wasn't chosen at sign-in, the owner can't see it, or the name is wrong | "I can't reach that base. Add it to the connection in Airtable, or check you have access to it." |
| Refused writing one table or field | It's locked to some people, or the owner is read-only in the base | "That table or field is locked for your account in Airtable. Someone who owns the base can unlock it." |
| The list of bases comes back empty | An enterprise admin blocks outside apps | "Your Airtable admin blocks outside apps. Ask them to allow it, or connect another account." |
| Too many calls, right away | Over 5 calls a second on one base | Wait 30 seconds and continue; tell the owner only if it keeps happening. |
| Too many calls, all day | The workspace's monthly calls are used up | "Your Airtable plan's monthly allowance is used up. The taf will continue next month, or sooner on a bigger plan." |
| Not accepted | A value doesn't fit its field (text in a number field, an unknown choice) | Fix it and try once; else "Airtable didn't accept that change. Nothing was written." |
| The sign-in no longer works | Access removed, or it expired | "Airtable needs you to sign in again. Nothing was changed." |
| Airtable's own error | A problem on Airtable's side | "Airtable had a problem. The taf will try again later. Nothing was changed." |

After any refusal the taf stops that step and says so. A refused read never
leads to a create.
