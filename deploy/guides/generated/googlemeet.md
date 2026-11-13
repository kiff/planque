# Google Meet: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

25 tools: 18 read, 7 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `GOOGLEMEET_BATCH_UPDATE_MEMBERS` | Batch Update Google Meet Members | write | `parent`, `requests`; *updateMask* |  | `meetings.space.created` |
| `GOOGLEMEET_CREATE_MEET` | Create Google Meet Space | write | *config* |  | `meetings.space.created` |
| `GOOGLEMEET_CREATE_MEMBER` | Create Member | write | `email`, `parent`; *fields, name, role* |  | `meetings.space.created` |
| `GOOGLEMEET_DELETE_MEMBER` | Delete Google Meet Space Member | write | `name` |  | `meetings.space.created` |
| `GOOGLEMEET_END_ACTIVE_CONFERENCE` | End active conference | write | `space_name` |  | `meetings.space.created` |
| `GOOGLEMEET_GET_CONFERENCE_RECORD_BY_NAME` | Get conference record by name | read | `name` |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLEMEET_GET_MEET` | Get Meet details | read | `space_name` |  | `meetings.space.readonly`, `meetings.space.created`, `meetings.space.settings` |
| `GOOGLEMEET_GET_MEMBER` | Get Member | read | `name`; *fields* |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLEMEET_GET_PARTICIPANT` | Get Participant | read | `name` |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLEMEET_GET_PARTICIPANT_SESSION` | Get Participant Details | read | `name` |  | `meetings.space.readonly`, `meetings.space.created` |
| `GOOGLEMEET_GET_RECORDING` | Get Recording | read | `name` |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLEMEET_GET_RECORDINGS_BY_CONFERENCE_RECORD_ID` | Get recordings by conference record ID | read | `conference_record_id` |  | `meetings.space.readonly`, `meetings.space.created` |
| `GOOGLEMEET_GET_SMART_NOTE` | Get Smart Note | read | `name` |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLEMEET_GET_TRANSCRIPT` | Get Transcript | read | `name` |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLEMEET_GET_TRANSCRIPTS_BY_CONFERENCE_RECORD_ID` | Get transcripts by conference record ID | read | `conference_record_id`; *page_size, page_token* |  | `meetings.space.readonly`, `meetings.space.created` |
| `GOOGLEMEET_GET_TRANSCRIPT_ENTRY` | Get Transcript Entry | read | `name` |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLEMEET_LIST_CONFERENCE_RECORDS` | List Conference Records | read | *filter, page_size, page_token* |  | `meetings.space.readonly`, `meetings.space.created` |
| `GOOGLEMEET_LIST_MEMBERS` | List Members | read | `parent`; *fields, page_size, page_token* |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLEMEET_LIST_PARTICIPANTS` | List Participants | read | `parent`; *filter, page_size, page_token* |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLEMEET_LIST_PARTICIPANT_SESSIONS` | List Participant Sessions | read | `conference_record_id`, `participant_id`; *filter, page_size, page_token* |  | `meetings.space.readonly`, `meetings.space.created` |
| `GOOGLEMEET_LIST_RECORDINGS` | List Recordings | read | `parent`; *page_size, page_token* |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLEMEET_LIST_SMART_NOTES` | List Smart Notes | read | `parent`; *page_size, page_token* |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLEMEET_LIST_TRANSCRIPT_ENTRIES` | List Transcript Entries | read | `conference_record_id`, `transcript_id`; *page_size, page_token* |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLEMEET_UPDATE_MEMBER` | Update Google Meet Member | write | `name`; *email, role, update_mask* |  | `meetings.space.created` |
| `GOOGLEMEET_UPDATE_SPACE` | Update Google Meet Space | write | `name`; *config, updateMask* |  | `meetings.space.created`, `meetings.space.settings` |
