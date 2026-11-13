# Pipedrive

Checked: 2026-11-13 · Sources: the connection service's Pipedrive tool
definitions as planque reads them (403 tools), Pipedrive's API reference,
its rate-limit, status-code and app-scope pages.

## How it connects

- **A sign-in to Pipedrive**, held by the connection service; the taf never
  sees a password or token. planque's catalog puts the tools the sign-in's
  access doesn't cover under `unusable`, with why.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  A sign-in made for reading deals doesn't let the taf change them, so ask
  for every tool the job needs in the first request.
- **The taf sees what the signed-in user sees.** Pipedrive's visibility
  settings still apply: deals or contacts hidden from that user stay hidden.

## Tools a job usually needs

There are hundreds; a sales job needs a handful. Read means it changes
nothing in the account.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `PIPEDRIVE_SEARCH_DEALS` | Finds deals by a word in the title or notes | read | deal ids, titles |
| `PIPEDRIVE_SEARCH_PERSONS` | Finds people by name, email or phone | read | person ids, names, their organisation's id |
| `PIPEDRIVE_LIST_DEALS` | Lists deals by owner, stage, status or ids, a page at a time | read | deal ids, stage ids, a cursor |
| `PIPEDRIVE_GET_PERSON` / `PIPEDRIVE_GET_ORGANIZATION` | Reads one person or organisation | read | the id, owner, contact details |
| `PIPEDRIVE_LIST_DEAL_ACTIVITIES` | Lists the calls, meetings and tasks on a deal | read | activity ids, due dates, done or not |
| `PIPEDRIVE_GET_ALL_NOTES` | Lists notes on a deal, person or organisation | read | note ids, content |
| `PIPEDRIVE_GET_ALL_PIPELINES` / `PIPEDRIVE_GET_ALL_STAGES` | Reads the pipelines and their stages | read | pipeline ids, stage ids, names |
| `PIPEDRIVE_ADD_NOTE` | Adds a note to a deal, person, organisation or lead | write | the note id |
| `PIPEDRIVE_ADD_AN_ACTIVITY` | Plans a call, meeting or task, linked to a deal or person | write | the activity id |
| `PIPEDRIVE_UPDATE_ACTIVITY` | Changes an activity or marks it done | write | the activity id |
| `PIPEDRIVE_ADD_A_PERSON` | Adds a person | write | the person id |
| `PIPEDRIVE_ADD_A_DEAL` | Opens a deal | write | the deal id |
| `PIPEDRIVE_UPDATE_DEAL` | Changes a deal: stage, value, close date, won or lost | write | the deal id, its stage id |

**What to keep:** the ids the taf works on (deal, person, organisation,
activity, note) and the pipeline and stage ids the owner's process uses.
Never search for something the taf made: use the id it kept. Look up
pipelines and stages once and keep them; don't guess a stage from its name.

**Start read-only.** "Tell me which deals have no next step" is a list tool
plus `PIPEDRIVE_LIST_DEAL_ACTIVITIES`. Adding notes and activities comes
next. Moving a deal's stage or marking it won or lost shows on the team's
reports: offer `PIPEDRIVE_UPDATE_DEAL` only when the owner asked, and show
each change first.

## Don't offer

- Deleting tools (`PIPEDRIVE_DELETE_*`), including bulk deletes, and
  `PIPEDRIVE_UPDATE_DEAL_V2`, which can mark a deal deleted.
- Merging tools (`PIPEDRIVE_MERGE_*`): two records become one, which can't be
  undone from here.
- Admin tools: pipelines and stages (adding, changing, deleting), activity
  types, fields, filters shared with the team, users, teams, roles and their
  pipeline visibility, and webhooks.
- Mail, call-log, video-call and messaging tools, products, subscriptions,
  installments and goals, unless the job is about exactly that and the owner
  asked.

## Access and plans

- **Every plan can use the API**, with limits that grow with the plan.
- **Daily budget per company:** Pipedrive counts tokens, not calls. Each day
  a company gets 30,000 tokens times its plan's multiplier (Lite 1, Growth 2,
  Premium 5, Ultimate 7) times its number of seats. Reading one record costs
  2, a list 20, an update 10, a search 40. The budget is shared by every app
  the company uses and resets at midnight in Pipedrive's server time.
- **Bursts:** per user, over a rolling two seconds, from 80 calls on Lite to
  480 on Ultimate for a signed-in app. Searches are capped at 10 every two
  seconds on every plan, so search once and keep the ids.
- **A lapsed account** (trial over, payment missing) is refused with a
  payment notice until it is active again.
- **Pipelines, stages, fields and users** need a Pipedrive admin's sign-in;
  they aren't offered anyway.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access | The sign-in didn't allow this tool: it was made for reading, say | "Pipedrive is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Refused (403) | The user isn't allowed this action, or reached a usage limit for that kind of record | "Pipedrive says your user can't do this. A Pipedrive admin can change your rights, or the taf can skip it." |
| Payment required (402) | The company account isn't active: trial over or payment missing | "Your Pipedrive account isn't active, so nothing can be read or changed. Once it's active the taf carries on." |
| Too many calls (429), in a burst | Too many calls in two seconds | Wait a few seconds and try again; tell the owner only if it keeps failing. |
| Too many calls (429), all day | The company's daily budget is used | "Your company's Pipedrive allowance for today is used up. The taf will continue after it resets tonight." |
| The sign-in no longer works (401) | The owner removed the app, or the sign-in expired | "Pipedrive needs you to sign in again. Nothing was changed." |
| Not found (404) | The record was deleted or merged, or the user can't see it | "That record isn't in Pipedrive any more. The taf will skip it." |
| Bad request (400) | A field the taf sent isn't valid: a stage from another pipeline, say | Fix it and try once; else "Pipedrive didn't accept that change. Nothing was written." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
