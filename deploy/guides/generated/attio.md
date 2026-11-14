# Attio: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

105 tools: 52 read, 53 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `ATTIO_ASSERT_COMPANY` | Assert Company (Create or Update) | write | `matching_attribute`; *description, domain, domains, name, primary_location, team* |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_ASSERT_PERSON` | Assert Person Record | write | `matching_attribute`, `values` |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_ASSERT_USER_RECORD` | Assert User Record (Create or Update) | write | `matching_attribute`, `values` |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_ASSERT_WORKSPACE` | Assert Workspace (Create or Update) | write | `matching_attribute`; *avatar_url, company, name, users, workspace_id* |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_CREATE_ATTRIBUTE` | Create Attribute | write | `api_slug`, `config`, `identifier`, `target`, `title`, `type`; *default_value, description, is_multiselect, is_required, is_unique, relationship* |  | `object_configuration:read-write` |
| `ATTIO_CREATE_COMMENT` | Create Comment | write | `author`, `content`; *entry, format, record, thread_id* |  | `comment:read-write` |
| `ATTIO_CREATE_COMPANY` | Create Company | write | *description, domains, name, primary_location, team* |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_CREATE_DEAL_RECORD` | Create Deal Record | write | *associated_company_domain, associated_company_id, associated_people_emails, associated_people_ids, name, owner_email, owner_id, stage and 1 more* |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_CREATE_LIST` | Create List | write | `api_slug`, `name`, `parent_object`, `workspace_access`, `workspace_member_access` |  | `list_configuration:read-write` |
| `ATTIO_CREATE_MEETING` | Create Meeting | write | `description`, `end`, `is_all_day`, `participants`, `start`, `title`; *linked_records* |  | `meeting:read-write`, `record_permission:read` |
| `ATTIO_CREATE_NOTE` | Create Note | write | `content`, `parent_object`, `parent_record_id`, `title`; *created_at* |  | `note:read-write`, `object_configuration:read`, `record_permission:read` |
| `ATTIO_CREATE_OBJECT` | Create Object | write | `data` |  | `object_configuration:read-write` |
| `ATTIO_CREATE_PERSON` | Create Person | write | *company, description, email_addresses, job_title, linkedin, name, phone_numbers, primary_location* |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_CREATE_RECORD` | Create Record | write | `object_type`, `values` |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_CREATE_SELECT_OPTION` | Create Select Option | write | `attribute`, `data`, `identifier`, `target` |  | `object_configuration:read-write` |
| `ATTIO_CREATE_STATUS` | Create Status | write | `attribute`, `data`, `identifier`, `target` |  | `object_configuration:read-write` |
| `ATTIO_CREATE_TASK` | Create Task | write | `content`; *assignees, deadline_at, is_completed, linked_records* |  | `task:read-write`, `object_configuration:read`, `record_permission:read`, `user_management:read` |
| `ATTIO_CREATE_USER_RECORD` | Create User Record | write | `primary_email_address`, `user_id`, `workspace_id`; *person_id* |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_CREATE_WEBHOOK` | Create Webhook | write | `subscriptions`, `target_url` |  | `webhook:read-write` |
| `ATTIO_CREATE_WORKSPACE_RECORD` | Create Workspace Record | write | `workspace_id`; *avatar_url, company, name, users* |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_DELETE_COMMENT` | Delete Comment | write | `comment_id` |  | `comment:read-write` |
| `ATTIO_DELETE_COMPANY` | Delete Company | write | `record_id` |  | `object_configuration:read`, `record_permission:read-write` |
| `ATTIO_DELETE_DEAL` | Delete Deal | write | `record_id` |  | `object_configuration:read`, `record_permission:read-write` |
| `ATTIO_DELETE_ENTRY` | Delete List Entry | write | `entry_id`, `list` |  | `list_entry:read-write`, `list_configuration:read` |
| `ATTIO_DELETE_NOTE` | Delete Note | write | `note_id` |  | `note:read-write` |
| `ATTIO_DELETE_PERSON` | Delete Person | write | `record_id` |  | `object_configuration:read`, `record_permission:read-write` |
| `ATTIO_DELETE_RECORD` | Delete Record | write | `object_type`, `record_id` |  | `object_configuration:read`, `record_permission:read-write` |
| `ATTIO_DELETE_TASK` | Delete Task | write | `task_id` |  | `task:read-write` |
| `ATTIO_DELETE_USER` | Delete User | write | `record_id` |  | `object_configuration:read`, `record_permission:read-write` |
| `ATTIO_DELETE_WEBHOOK` | Delete Webhook | write | `webhook_id` |  | `webhook:read-write` |
| `ATTIO_DELETE_WORKSPACE_RECORD` | Delete Workspace Record | write | `record_id` |  | `object_configuration:read`, `record_permission:read-write` |
| `ATTIO_FIND_RECORD` | Find Record | read | `object_id`; *attributes, limit, offset, record_id* |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_GET_ATTRIBUTE` | Get Attribute Details | read | `attribute`, `identifier`, `target` |  | `object_configuration:read-write` |
| `ATTIO_GET_CALL_RECORDING` | Get Call Recording | read | `call_recording_id`, `meeting_id` |  | `meeting:read`, `call_recording:read` |
| `ATTIO_GET_COMMENT` | Get Comment | read | `comment_id` |  | `comment:read`, `object_configuration:read`, `record_permission:read` |
| `ATTIO_GET_COMPANY` | Get Company | read | `record_id` |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_GET_DEAL_RECORD` | Get Deal Record | read | `record_id` |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_GET_LIST` | Get List | read | `list` |  | `list_configuration:read` |
| `ATTIO_GET_LIST_ENTRY` | Get List Entry | read | `entry_id`, `list` |  | `list_entry:read`, `list_configuration:read` |
| `ATTIO_GET_MEETING` | Get Meeting | read | `meeting_id` |  | `meeting:read`, `record_permission:read` |
| `ATTIO_GET_NOTE` | Get Note | read | `note_id` |  | `note:read`, `object_configuration:read`, `record_permission:read` |
| `ATTIO_GET_OBJECT` | Get Object | read | `object` |  | `object_configuration:read` |
| `ATTIO_GET_RECORD_ENTRIES` | List Record Entries | read | `object`, `record_id`; *limit, offset* |  | `record_permission:read`, `object_configuration:read`, `list_entry:read` |
| `ATTIO_GET_SELF` | Get Current Token Info | read |  |  |  |
| `ATTIO_GET_TASK` | Get Task | read | `task_id` |  | `task:read`, `object_configuration:read`, `record_permission:read`, `user_management:read` |
| `ATTIO_GET_V2_TASKS` | List Tasks | read | *assignee, is_completed, limit, linked_object, linked_record_id, offset, sort* |  | `task:read`, `object_configuration:read`, `record_permission:read`, `user_management:read` |
| `ATTIO_GET_WEBHOOK` | Get Webhook | read | `webhook_id` |  | `webhook:read` |
| `ATTIO_GET_WORKSPACE_MEMBER` | Get Workspace Member | read | `workspace_member_id` |  | `user_management:read` |
| `ATTIO_GET_WORKSPACE_RECORD` | Get Workspace Record | read | `record_id` |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_LIST_ATTRIBUTES` | List Attributes | read | `identifier`, `target`; *limit, offset, show_archived* |  | `object_configuration:read` |
| `ATTIO_LIST_ATTRIBUTE_OPTIONS` | List Attribute Options | read | `attribute`, `identifier`, `target`; *limit, offset, show_archived* |  | `object_configuration:read` |
| `ATTIO_LIST_ATTRIBUTE_STATUSES` | List Attribute Statuses | read | `attribute`, `identifier`, `target`; *show_archived* |  | `object_configuration:read` |
| `ATTIO_LIST_CALL_RECORDINGS` | List Call Recordings | read | `meeting_id`; *cursor, limit* |  | `meeting:read`, `call_recording:read` |
| `ATTIO_LIST_COMPANIES` | List Companies | read | *filter, limit, offset, sorts* |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_LIST_COMPANY_ATTRIBUTE_VALUES` | List Company Attribute Values | read | `attribute`, `record_id`; *limit, offset, show_historic* |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_LIST_COMPANY_RECORD_ENTRIES` | List Company Record Entries | read | `record_id`; *limit, offset* |  | `record_permission:read`, `object_configuration:read`, `list_entry:read` |
| `ATTIO_LIST_DEAL_ENTRIES` | List Deal Entries | read | `record_id`; *limit, offset* |  | `record_permission:read`, `object_configuration:read`, `list_entry:read` |
| `ATTIO_LIST_DEAL_RECORDS` | List Deal Records | read | *filter, limit, offset, sorts* |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_LIST_DEAL_RECORD_ATTRIBUTE_VALUES` | List Deal Record Attribute Values | read | `attribute`, `record_id`; *limit, offset, show_historic* |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_LIST_LISTS` | List Lists | read |  |  | `list_configuration:read` |
| `ATTIO_LIST_LIST_ENTRY_ATTRIBUTE_VALUES` | List List Entry Attribute Values | read | `attribute`, `entry_id`, `list`; *limit, offset, show_historic* |  | `list_entry:read`, `list_configuration:read` |
| `ATTIO_LIST_MEETINGS` | List Meetings | read | *cursor, ends_from, limit, linked_object, linked_record_id, participants, sort, starts_before and 1 more* |  | `meeting:read`, `record_permission:read` |
| `ATTIO_LIST_NOTES` | List Notes | read | *limit, offset, parent_object, parent_record_id* |  | `note:read`, `object_configuration:read`, `record_permission:read` |
| `ATTIO_LIST_OBJECTS` | List Objects | read |  |  | `object_configuration:read` |
| `ATTIO_LIST_PEOPLE_ATTRIBUTE_VALUES` | List People Attribute Values | read | `attribute`, `record_id`; *limit, offset, show_historic* |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_LIST_PEOPLE_RECORD_ENTRIES` | List People Record Entries | read | `record_id`; *limit, offset* |  | `record_permission:read`, `object_configuration:read`, `list_entry:read` |
| `ATTIO_LIST_RECORDS` | List Records | read | `object_type`; *limit, offset* |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_LIST_RECORD_ATTRIBUTE_VALUES` | List Record Attribute Values | read | `attribute`, `object_type`, `record_id`; *limit, offset, show_historic* |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_LIST_THREADS` | List Threads | read | *entry_id, limit, list, object, offset, record_id* |  | `comment:read`, `object_configuration:read`, `record_permission:read` |
| `ATTIO_LIST_USER_RECORDS` | List User Records | read | *filter, limit, offset, sorts* |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_LIST_USER_RECORD_ENTRIES` | List User Record Entries | read | `record_id`; *limit, offset* |  | `record_permission:read`, `object_configuration:read`, `list_entry:read` |
| `ATTIO_LIST_VIEWS_FOR_LIST` | List Views for List | read | `list`; *cursor, limit, show_archived* |  | `list_configuration:read` |
| `ATTIO_LIST_WEBHOOKS` | List Webhooks | read | *limit, offset* |  | `webhook:read` |
| `ATTIO_LIST_WORKSPACE_MEMBERS` | List Workspace Members | read |  |  | `user_management:read` |
| `ATTIO_LIST_WORKSPACE_RECORDS` | List Workspace Records | read | *filter, limit, offset, sorts* |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_LIST_WORKSPACE_RECORD_ATTRIBUTE_VALUES` | List Workspace Record Attribute Values | read | `attribute`, `record_id`; *limit, offset, show_historic* |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_LIST_WORKSPACE_RECORD_ENTRIES` | List Workspace Record Entries | read | `record_id`; *limit, offset* |  | `record_permission:read`, `object_configuration:read`, `list_entry:read` |
| `ATTIO_MERGE_RECORDS` | Merge Records | write | `object`, `primary_record_id`, `secondary_record_id` |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_PATCH_RECORD` | Patch Record | write | `object`, `record_id`, `values` |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_PATCH_V2_LISTS_LIST_ENTRIES_ENTRY_ID` | Update List Entry by Entry ID | write | `entry_id`, `entry_values`, `list` |  | `list_entry:read-write`, `list_configuration:read` |
| `ATTIO_PEOPLE_GET_PERSON` | Get Person | read | `record_id` |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_PEOPLE_LIST_PERSONS` | List Person Records | read | *filter, limit, offset, sorts* |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_POST_V2_LISTS_LIST_ENTRIES` | Create List Entry | write | `list`, `parent_object`, `parent_record_id`; *entry_values* |  | `list_entry:read-write`, `list_configuration:read` |
| `ATTIO_POST_V2_LISTS_LIST_ENTRIES_QUERY` | List List Entries | read | `list`; *filter, limit, offset, sorts* |  | `list_entry:read`, `list_configuration:read` |
| `ATTIO_POST_V2_OBJECTS_OBJECT_RECORDS_QUERY` | List Records with Query | read | `object`; *filter, limit, offset, sorts* |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_PUT_V2_LISTS_LIST_ENTRIES` | Assert List Entry by Parent | write | `list`, `parent_object`, `parent_record_id`; *entry_values* |  | `list_entry:read-write`, `list_configuration:read` |
| `ATTIO_PUT_V2_LISTS_LIST_ENTRIES_ENTRY_ID` | Update List Entry by Entry ID (PUT) | write | `entry_id`, `entry_values`, `list` |  | `list_entry:read-write`, `list_configuration:read` |
| `ATTIO_PUT_V2_OBJECTS_OBJECT_RECORDS` | Assert Record (Create or Update) | write | `matching_attribute`, `object`, `values` |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_PUT_V2_OBJECTS_OBJECT_RECORDS_RECORD_ID` | Put Record (Overwrite) | write | `object`, `record_id`, `values` |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_QUERY_RECORDS` | Query Records | read | `object`; *filter, limit, offset, sorts* |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_SEARCH_RECORDS` | Search Records | read | `objects`, `query`, `request_as`; *limit* |  | `record_permission:read`, `object_configuration:read` |
| `ATTIO_UPDATE_ATTRIBUTE` | Update Attribute | write | `attribute`, `data`, `identifier`, `target` |  | `object_configuration:read-write` |
| `ATTIO_UPDATE_COMPANY` | Update Company | write | `record_id`, `values` |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_UPDATE_DEAL_RECORD` | Update Deal Record | write | `record_id`, `values` |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_UPDATE_LIST` | Update List | write | `list`; *api_slug, name, workspace_access, workspace_member_access* |  | `list_configuration:read-write` |
| `ATTIO_UPDATE_NOTE` | Update Note | write | `note_id`; *content, format, title* |  | `note:read-write`, `object_configuration:read`, `record_permission:read` |
| `ATTIO_UPDATE_OBJECT` | Update Object | write | `data`, `object_id` |  | `object_configuration:read-write` |
| `ATTIO_UPDATE_PERSON` | Update Person | write | `record_id`, `values` |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_UPDATE_RECORD` | Update Record | write | `object_type`, `record_id`, `values` |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_UPDATE_SELECT_OPTION` | Update Select Option | write | `attribute`, `data`, `identifier`, `option`, `target` |  | `object_configuration:read-write` |
| `ATTIO_UPDATE_STATUS` | Update Status | write | `attribute`, `data`, `identifier`, `status`, `target` |  | `object_configuration:read-write` |
| `ATTIO_UPDATE_TASK` | Update Task | write | `data`, `task_id` |  | `task:read-write`, `object_configuration:read`, `record_permission:read`, `user_management:read` |
| `ATTIO_UPDATE_USER_RECORD` | Update User Record | write | `record_id`, `values` |  | `record_permission:read-write`, `object_configuration:read` |
| `ATTIO_UPDATE_WEBHOOK` | Update Webhook | write | `webhook_id`; *subscriptions, target_url* |  | `webhook:read-write` |
| `ATTIO_UPDATE_WORKSPACE_RECORD` | Update Workspace Record | write | `record_id`, `values` |  | `record_permission:read-write`, `object_configuration:read` |
