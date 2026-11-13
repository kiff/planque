# Clay

Checked: 2026-11-13 · Sources: [the generated part](generated/clay.md)
(15 tools), and Clay's developer docs at developers.clay.com: the Public API
reference, authentication, errors, rate limits, pagination, searches,
routines, batch runs and tables pages, plus Clay's credits page.

## How it connects

- **A key the owner pastes**, held by the connection service; the taf never
  sees it. The owner makes it in Clay under API and CLI, API keys. A key
  belongs to one Clay user and acts with that user's access to one
  workspace.
- **One approval for the key, one for the tools and the Card.** Ask for every
  tool the job needs in the first request.
- **Check the key first** with `CLAY_GET_CURRENT_USER`: it names the key's
  user and workspace.

## Tools a job usually needs

Read means it changes nothing in the workspace. Returns are the ids a taf
keeps; they come from Clay's API reference, since the catalog lists none.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `CLAY_GET_CURRENT_USER` | Says whose key this is | read | the user, the workspace |
| `CLAY_GET_CREDIT_BALANCE` | Reads the workspace's credits left, and the key owner's monthly limit; costs nothing | read | balances |
| `CLAY_GET_SEARCH_QUERY_REFERENCE` | Reads how to write a people or company search | read | the reference text |
| `CLAY_SEARCH_PEOPLE_OR_COMPANIES` | Starts a search of Clay's people or company database | read | a search id, records, whether more remain |
| `CLAY_CONTINUE_SEARCH` | Fetches the next page of a search | read | records, whether more remain |
| `CLAY_RUN_ROUTINE` | Runs one of the workspace's functions on 1 to 100 items; it runs in the background | write, spends credits | a run id, its status |
| `CLAY_GET_ROUTINE_RUN_RESULTS` | Reads a run's progress, then each item's result | read | each item's own id, its status, a cursor |
| `CLAY_QUERY_TABLE_ROWS` | Reads rows from a Clay table whose id is known (Enterprise only) | read | rows, field names, a cursor |
| `CLAY_GET_WORKFLOW_RUNS_QUERY_REFERENCE` / `CLAY_QUERY_WORKFLOW_RUNS` | Reads how to query workflow runs, then finds runs | read | run records, a cursor |

**What to keep:** a search id, to fetch more of the same search; a run id,
to read its results later; and, for each item sent to a function, the id
the taf gave it, since results are matched by that id and not by order.
Never run a search or function again to find something already kept.

**A run takes time.** Until it's done, its results show progress only.
Wait between reads. A finished run can hold failed items: check each
item's status.

**Start read-only.** Searching, reading tables and checking the balance
change nothing. Offer `CLAY_RUN_ROUTINE` only when the owner asked for
enrichment, name the function and how many items it will run on, and check
the balance with `CLAY_GET_CREDIT_BALANCE` first.

## Don't offer

- `CLAY_PRIVACY_DELETE_PEOPLE`: it queues people for deletion from the
  workspace's Audiences. It's for privacy requests the owner handles
  themselves, not for a taf.
- Batch tools (`CLAY_GET_ROUTINE_BATCH_UPLOAD_URL`,
  `CLAY_START_ROUTINE_BATCH`, `CLAY_GET_ROUTINE_BATCH_RESULTS`): they need a
  file uploaded outside the tools and can spend many credits at once. A job
  of more than 100 items runs `CLAY_RUN_ROUTINE` in pieces, with the
  owner's yes.

## Access and plans

- **Search results are capped by plan,** across all API use of the
  workspace:
  - Free: 50 a request, 50 a search, 100 every 30 days.
  - Trial: 50 a request, 50 a search, 10,000 for the trial.
  - Paid: 500 a request, 1,000,000 every 30 days.
  - Older Starter, Explorer and Pro plans: 500 a request, 1,000,000 a year,
    reset on 1 January.
  - Enterprise: 500 a request, 10,000,000 every 30 days.

  Counting searches and job searches aren't offered through the API.
- **Tables are Enterprise only,** and there's no way to list them: the owner
  gives the table id, the part of the table's Clay link after `/tables/`.
- **Credits:** Clay counts Data Credits (bought data and outside AI) and
  Actions (each enrichment step per record). Running out of Data Credits
  can be fixed with a top-up; running out of Actions needs a higher plan.
- **Rate:** each workspace has a request limit Clay doesn't publish. When it's
  reached, Clay says how many seconds to wait.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| "Authentication failed" (401 or 403) | The key is wrong, was deleted, or its user lost access | "Clay didn't accept the key. Could you make a new one in Clay's API keys page and paste it again?" |
| Search refused for a limit (402) | The plan's per-search or per-period result cap is reached; the message says which, and when it resets | "Clay's search limit for your plan is reached. It frees up on the date Clay gives; a higher plan raises it." |
| Search refused as invalid (400) | The search was written wrongly | Read the search reference, fix it and try once; else "Clay didn't understand that search. Can you describe who you're looking for another way?" |
| Not found (404) | A table, function or run id that doesn't exist or isn't in this workspace | "I can't find that in your Clay workspace. Could you check the link or name?" |
| Refused (409 or 422) | The request doesn't fit the item's current state | Fix it and try once; else "Clay didn't accept that request. Nothing was run." |
| Too many calls (429) | Over the workspace's rate | "Clay asked us to slow down. The taf will continue in a moment." |
| Still running (202) | The run or job hasn't finished | Wait and read again later; tell the owner only if it never finishes. |
| An item failed in a finished run | That item's lookup failed at the data source | "Clay couldn't enrich some of the records. The rest are done." |
| Out of credits or Actions | The workspace's balance is spent | "Your Clay workspace is out of credits. Add more in Clay, then the taf can finish." |
| Server error (5xx) | A passing problem at Clay | Retry a read later; tell the owner only if it keeps failing. |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
