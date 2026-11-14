# Pipedrive

Checked: 2026-11-13 · Sources: the connection service's Pipedrive tool
definitions as planque reads them, and Pipedrive's API documentation (rate
limits, status codes, sign-in, and its published API description).

## How it connects

- **A sign-in to Pipedrive**, held by the connection service; the taf never
  sees a password or token. The owner allows a set of access areas (deals,
  persons, activities and so on), and calls outside them are refused.
  planque's catalog lists the tools and, with a known sign-in, puts the ones
  its access doesn't cover under `unusable`, with why.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.
- **A sign-in unused for 60 days ends.** Then the owner signs in again.

## Tools a job usually needs

Read means it changes nothing in Pipedrive. The catalog gives no `returns`
paths for Pipedrive today; the ids below are those Pipedrive's documented
answers carry.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `PIPEDRIVE_GET_CURRENT_USER_DATA` | Reads the signed-in user | read | the user's id |
| `PIPEDRIVE_SEARCH_PERSONS` | Finds people by name, email or phone | read | person ids, their organization's id |
| `PIPEDRIVE_SEARCH_ORGANIZATIONS` | Finds companies by name | read | organization ids |
| `PIPEDRIVE_SEARCH_DEALS` | Finds deals by title or text | read | deal ids, their stage, person and organization |
| `PIPEDRIVE_GET_DEAL` | Reads one deal | read | the deal's id, stage, pipeline, person and organization ids |
| `PIPEDRIVE_GET_ALL_DEALS_BETA` | Lists deals, by stage, status or owner, a page at a time | read | deal ids, a cursor for the next page |
| `PIPEDRIVE_GET_ALL_ACTIVITIES_BETA` | Lists calls, meetings and to-dos, by date or done | read | activity ids with their deal and person ids, a cursor |
| `PIPEDRIVE_GET_ALL_LEADS` | Lists leads not yet turned into deals | read | lead ids, their person and organization ids |
| `PIPEDRIVE_ADD_A_PERSON` | Adds a person | write | the person's id |
| `PIPEDRIVE_ADD_LEAD` | Adds a lead, linked to a person or organization | write | the lead's id |
| `PIPEDRIVE_ADD_A_DEAL` | Adds a deal | write | the deal's id |
| `PIPEDRIVE_UPDATE_A_DEAL` | Changes a deal: stage, value, expected close, won or lost | write | the deal, with its new stage |
| `PIPEDRIVE_ADD_AN_ACTIVITY` | Adds a follow-up call, meeting or to-do | write | the activity's id |
| `PIPEDRIVE_ADD_NOTE` | Adds a note to a deal, person, organization or lead | write | the note's id |

**What to keep:** each person's, deal's and lead's id, to update that same
record next time; the signed-in user's id, to list only their own deals.
Never look a record up again by name: use the id the taf kept.

**Search first.** Before adding a person, search by their email. Add them
only when nothing matches, so no one is entered twice. A search needs at
least two characters.

**Start read-only.** "Tell me which deals have no next step" needs only the
read tools. Adding a follow-up activity or a note is a small write; moving
a deal's stage or marking it won or lost changes the owner's pipeline, so
offer it only when asked, and each change waits for the owner's yes.

## Don't offer

- Deleting (`PIPEDRIVE_DELETE_PERSON`, `PIPEDRIVE_DELETE_A_PERSON`,
  `PIPEDRIVE_DELETE_DEAL` and the others): records are removed for everyone.
- Merging (`PIPEDRIVE_MERGE_PERSONS`, `PIPEDRIVE_MERGE_DEALS`,
  `PIPEDRIVE_MERGE_ORGANIZATIONS` and their two-record versions): a merge
  can't be undone from here.
- Setup and admin: pipelines and stages (`PIPEDRIVE_ADD_PIPELINE`,
  `PIPEDRIVE_ADD_STAGE`), custom fields (`PIPEDRIVE_ADD_PERSON_FIELD`),
  saved filters, activity types and lead labels. They change Pipedrive for
  the whole team.
- Turning deals into leads or back (`PIPEDRIVE_CONVERT_DEAL_TO_LEAD`,
  `PIPEDRIVE_CONVERT_LEAD_TO_DEAL`) unless the owner asks.

## Access and plans

- **A daily budget, per company:** Pipedrive gives each company a number of
  call tokens a day, set by its plan and its number of seats, and shared by
  every app. Each call costs tokens; a search costs more than a read of one
  record. When the day's budget is spent, calls are refused until the next
  day. Pipedrive's admins get an email at 75% and 100%.
- **Bursts:** calls are also limited over any two seconds, more on higher
  plans. Searches are limited to 10 every two seconds on every plan.
- **A job that runs every hour** should list what changed since the last
  run, not read everything again.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused, not allowed (403) | The sign-in didn't give this area, or the account hit a usage limit | "Pipedrive is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Too many calls (429) | Over the two-second limit, or the day's budget is spent | "Pipedrive asked us to slow down. The taf will try again shortly, or tomorrow if today's allowance is used up." |
| Account not open (402) | The trial ended, or payment details are missing | "Your Pipedrive account isn't active right now. Nothing was changed." |
| The sign-in no longer works (401) | The sign-in ended, or was removed | "Pipedrive needs you to sign in again. Nothing was changed." |
| Not found (404) or gone (410) | The record was deleted or merged | "That record isn't in Pipedrive any more. The taf will skip it." |
| Pipedrive is unavailable (503) | Planned maintenance | "Pipedrive is down for maintenance. The taf will try again later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write: a taf that can't find a person never adds them again.
