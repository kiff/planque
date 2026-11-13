# Attio

Checked: 2026-11-13 · Sources: the connection service's Attio tool
definitions as planque reads them (105 tools), and Attio's REST API
documentation (authentication, rate limits, pagination, records, notes,
tasks).

## How it connects

- **A sign-in to Attio**, held by the connection service; the taf never sees
  a password or token. The sign-in is for one Attio workspace, the one the
  owner picks on Attio's page. planque's catalog lists the tools and, with a
  known sign-in, puts the ones its access doesn't cover under `unusable`,
  with why.
- **Access is per area:** records, lists, notes, tasks, comments, meetings
  and workspace members each have their own read or read-and-write access.
  Most record tools also need to read the workspace's object setup. Ask for
  every tool the job needs in the first request.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
- **When the sign-in acts as one member** of the workspace, it has only that
  member's permissions: a record they can't see in Attio, the taf can't
  either.

## Tools a job usually needs

Read means it changes nothing in the workspace. The catalog's `returns`
gives exact paths where it has them; where it has none yet, Attio's answer
carries the record's id and a link to its page in Attio.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `ATTIO_SEARCH_RECORDS` | Finds people, companies or deals by a word, such as a name, email or domain | read | record ids |
| `ATTIO_QUERY_RECORDS` | Lists one kind of record with filters and sorting, a page at a time | read | record ids |
| `ATTIO_PEOPLE_GET_PERSON` | Reads one person | read | the person, its id, its link |
| `ATTIO_GET_COMPANY` | Reads one company | read | the company, its id, its link |
| `ATTIO_GET_DEAL_RECORD` | Reads one deal | read | the deal, its id, its link |
| `ATTIO_LIST_NOTES` | Lists notes, all or for one record | read | note ids |
| `ATTIO_GET_V2_TASKS` | Lists tasks, by assignee, record or done | read | task ids |
| `ATTIO_LIST_WORKSPACE_MEMBERS` | Lists the people in the workspace, to assign tasks | read | member ids |
| `ATTIO_ASSERT_PERSON` | Adds a person, or updates the one that matches (by email, say) | write | the person's id, its link |
| `ATTIO_ASSERT_COMPANY` | Adds a company, or updates the one that matches (by domain, say) | write | the company's id, its link |
| `ATTIO_UPDATE_PERSON` | Changes values on one person | write | the person's id |
| `ATTIO_UPDATE_DEAL_RECORD` | Changes a deal's stage, value or other fields | write | the deal's id |
| `ATTIO_CREATE_NOTE` | Adds a note to a record | write | the new note's id |
| `ATTIO_CREATE_TASK` | Adds a to-do, linked to records and assigned to members | write | the new task's id |
| `ATTIO_POST_V2_LISTS_LIST_ENTRIES` | Adds a record to a list (a pipeline or a board) | write | the entry's id |

**What to keep:** each record's id, to read or update that same record next
time; its link, to show the owner; the ids of notes and tasks the taf made,
so it never adds the same one twice.

**Match, don't duplicate.** The "assert" tools update the record whose
matching field (an email, a domain) has the given value, and add one only if
none does. Match on a field that is unique; if the value matches more than
one record, Attio refuses.

**Notes and tasks:** a note needs a title and belongs to one record. A task's
text is plain, up to 2,000 characters, and its records are linked separately
rather than written into the text.

**Start read-only.** "Summarise my open deals" is search and read tools
only. Offer a write only for what the owner asked to change.

## Don't offer

- Deleting tools (records, people, companies, deals, notes, tasks, list
  entries, comments): deleted records can't be brought back by the taf.
- Merging records (`ATTIO_MERGE_RECORDS`): it can't be undone.
- Overwriting tools (`ATTIO_PUT_V2_OBJECTS_OBJECT_RECORDS_RECORD_ID`,
  `ATTIO_PUT_V2_LISTS_LIST_ENTRIES_ENTRY_ID`): they clear any value not sent.
  Use the update tools, which change only what's given.
- Workspace setup (creating or changing objects, attributes, options,
  statuses or lists) and webhooks: they change Attio for the whole team.
- Call recordings: they hold people's voices, which a job rarely needs.

## Access and plans

- **Rate:** Attio allows about 100 reads and 25 writes a second. A write that also creates linked records counts each one.
  When over, Attio says when to try again, usually within a second.
- **Paging:** lists come a page at a time; a taf keeps asking until a page
  comes back short.
- **Plans:** we found no plan limit on the tools above in Attio's API
  documentation.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403) | The sign-in didn't give this area, or the person who signed in can't do this in Attio | "Attio is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Too many calls (429) | Over Attio's per-second limit | "Attio asked us to slow down. The taf will try again in a moment." |
| The sign-in no longer works | The owner removed the app, or the sign-in expired | "Attio needs you to sign in again. Nothing was changed." |
| A record already has that value (409, uniqueness) | Another record holds the same unique email or domain | "That's already in Attio. The taf will update the existing record instead." |
| Matches more than one record | The matching field isn't unique in this workspace | "More than one record in Attio matches. Which one should the taf use?" |
| Changed at the same moment (409) | Another change landed during the write | "Someone changed that record at the same time. The taf left it as it is and will try on its next run." |
| A value isn't one Attio knows (400) | A stage or option name that isn't set up | "Attio doesn't have that option. Which existing one should the taf use?" |
| Not found (404) | The record was deleted or merged, or the object doesn't exist | "That record isn't in Attio any more. The taf will skip it." |

After any refusal the taf stops that step and says so. A refused read never
leads to a create or a write.
