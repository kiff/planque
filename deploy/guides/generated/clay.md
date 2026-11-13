# Clay: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner pastes a key from the vendor's settings). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

15 tools: 11 read, 4 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `CLAY_CONTINUE_SEARCH` | Continue Search | read | `limit`, `search_id` |  |  |
| `CLAY_GET_CREDIT_BALANCE` | Get Credit Balance | read |  |  |  |
| `CLAY_GET_CURRENT_USER` | Get Current User | read |  |  |  |
| `CLAY_GET_PRIVACY_DELETE_STATUS` | Get Privacy Delete Status | read | `job_id` |  |  |
| `CLAY_GET_ROUTINE_BATCH_RESULTS` | Get Routine Batch Results | read | `routine_run_id` |  |  |
| `CLAY_GET_ROUTINE_BATCH_UPLOAD_URL` | Get Routine Batch Upload URL | write | `routine_id` |  |  |
| `CLAY_GET_ROUTINE_RUN_RESULTS` | Get Routine Run Results | read | `routine_run_id`; *cursor, limit* |  |  |
| `CLAY_GET_SEARCH_QUERY_REFERENCE` | Get Search Query Reference | read |  |  |  |
| `CLAY_GET_WORKFLOW_RUNS_QUERY_REFERENCE` | Get Workflow Runs Query Reference | read |  |  |  |
| `CLAY_PRIVACY_DELETE_PEOPLE` | Privacy Delete People | write | `person_identifiers`; *dry_run* |  |  |
| `CLAY_QUERY_TABLE_ROWS` | Query Table Rows | read | `query`; *cursor, limit* |  |  |
| `CLAY_QUERY_WORKFLOW_RUNS` | Query Workflow Runs | read | `query`; *cursor, limit* |  |  |
| `CLAY_RUN_ROUTINE` | Run Routine | write | `items`, `routine_id`; *poll_interval_seconds, wait_timeout_seconds, webhook_id* |  |  |
| `CLAY_SEARCH_PEOPLE_OR_COMPANIES` | Search People or Companies | read | `limit`, `query` |  |  |
| `CLAY_START_ROUTINE_BATCH` | Start Routine Batch | write | `file_id`, `routine_id`; *webhook_id* |  |  |
