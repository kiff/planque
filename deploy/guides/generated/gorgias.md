# Gorgias: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

51 tools: 25 read, 26 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `GORGIAS_ADD_TICKET_TAGS` | Add Ticket Tags | write | `tag_ids`, `ticket_id` |  |  |
| `GORGIAS_CREATE_ACCOUNT_SETTING` | Create Account Setting | write | `name`, `value`; *description* |  |  |
| `GORGIAS_CREATE_CUSTOMER` | Create Customer | write | `name`; *address, channels, data, email, external_id* |  |  |
| `GORGIAS_CREATE_MACRO` | Create Macro | write | `name`; *actions, external_id, intent, language* |  | `macros:write` |
| `GORGIAS_CREATE_RULE` | Create Rule | write | `code`, `name`; *code_ast, deactivated_datetime, description, event_types, priority* |  | `rules:write` |
| `GORGIAS_CREATE_TAG` | Create Tag | write | `name`; *decoration, description* |  | `tags:write` |
| `GORGIAS_CREATE_TEAM` | Create Team | write | `name`; *decoration, description, members* |  |  |
| `GORGIAS_CREATE_TICKET` | Create Ticket | write | `messages`, `via`; *assignee_team_id, assignee_user_id, channel, customer, external_id, language, meta, priority and 3 more* |  |  |
| `GORGIAS_CREATE_VIEW` | Create View | write | `slug`; *decoration, fields, filters, name, order_by, order_dir, shared_with_teams, shared_with_users and 2 more* |  | `tickets:write` |
| `GORGIAS_DELETE_CUSTOMER` | Delete Customer | write | `customer_id` |  |  |
| `GORGIAS_DELETE_CUSTOMERS` | Delete Customers | write | `customer_ids` |  |  |
| `GORGIAS_DELETE_CUSTOMER_FIELD_VALUE` | Delete Customer Field Value | write | `customer_id`, `field_id` |  |  |
| `GORGIAS_DELETE_TEAM` | Delete Team | write | `team_id` |  |  |
| `GORGIAS_DELETE_TICKET` | Delete Ticket | write | `ticket_id` |  |  |
| `GORGIAS_DELETE_TICKET_FIELD_VALUE` | Delete Ticket Field Value | write | `field_id`, `ticket_id` |  |  |
| `GORGIAS_GET_ACCOUNT` | Get Account | read |  |  |  |
| `GORGIAS_GET_CUSTOMER` | Get Customer | read | `customer_id` |  |  |
| `GORGIAS_GET_EVENT` | Get Event | read | `event_id` |  |  |
| `GORGIAS_GET_MACRO` | Get Macro | read | `macro_id` |  | `macros:read` |
| `GORGIAS_GET_RULE` | Get Rule | read | `rule_id` |  | `rules:read` |
| `GORGIAS_GET_TEAM` | Get Team | read | `team_id` |  |  |
| `GORGIAS_GET_TICKET` | Get Ticket | read | `ticket_id` |  |  |
| `GORGIAS_GET_TICKET_MESSAGE` | Get Ticket Message | read | `message_id`, `ticket_id` |  | `tickets:read` |
| `GORGIAS_LIST_ACCOUNT_SETTINGS` | List Account Settings | read | *limit, offset* |  |  |
| `GORGIAS_LIST_CUSTOMERS` | List Customers | read | *cursor, email, external_id, limit, name, order_by* |  |  |
| `GORGIAS_LIST_CUSTOMER_FIELD_VALUES` | List Customer Field Values | read | `customer_id` |  |  |
| `GORGIAS_LIST_EVENTS` | List Events | read | *cursor, limit, object_id, object_type, types* |  |  |
| `GORGIAS_LIST_INTEGRATIONS` | List Integrations | read | *cursor, limit, order_by, type* |  | `integrations:read` |
| `GORGIAS_LIST_MACROS` | List Macros | read | *archived, cursor, languages, limit, message_id, number_predictions, order_by, search and 2 more* |  | `macros:read` |
| `GORGIAS_LIST_MESSAGES` | List Messages | read | *cursor, limit, order_by, ticket_id* |  | `tickets:read` |
| `GORGIAS_LIST_RULES` | List Rules | read | *cursor, limit, order_by, search* |  | `rules:read` |
| `GORGIAS_LIST_TAGS` | List Tags | read | *cursor, limit, order_by, search* |  | `tags:read` |
| `GORGIAS_LIST_TEAMS` | List Teams | read | *cursor, limit* |  |  |
| `GORGIAS_LIST_TICKETS` | List Tickets | read | *cursor, customer_id, external_id, limit, order_by, rule_id, ticket_ids, trashed and 1 more* |  |  |
| `GORGIAS_LIST_TICKET_FIELD_VALUES` | List Ticket Field Values | read | `ticket_id` |  |  |
| `GORGIAS_LIST_TICKET_TAGS` | List Ticket Tags | read | `ticket_id` |  |  |
| `GORGIAS_LIST_USERS` | List Users | read | *available_first, cursor, email, external_id, limit, order_by, roles, search* |  | `users:read` |
| `GORGIAS_LIST_VIEWS` | List Views | read | *category, cursor, limit, order_by* |  | `tickets:read` |
| `GORGIAS_LIST_VIEW_ITEMS` | List View Items | read | `view_id`; *cursor, direction, ignored_item, limit, order_by* |  | `tickets:read` |
| `GORGIAS_MERGE_CUSTOMERS` | Merge Customers | write | `source_customer_id`, `target_customer_id` |  |  |
| `GORGIAS_QUERY_REPORTING_STATISTICS` | Query Reporting Statistics | read | `query`; *cursor, limit* |  | `statistics:read` |
| `GORGIAS_REMOVE_TICKET_TAGS` | Remove Ticket Tags | write | `tag_ids`, `ticket_id` |  |  |
| `GORGIAS_SEND_TICKET_MESSAGE` | Send Ticket Message | write | `channel`, `ticket_id`; *attachments, body_html, body_text, external_id, if_unmodified_since, integration_id, receiver, sender and 2 more* |  | `tickets:write` |
| `GORGIAS_SET_CUSTOMER_DATA` | Set Customer Data | write | `customer_id`, `data` |  |  |
| `GORGIAS_SET_TICKET_TAGS` | Set Ticket Tags | write | `tag_ids`, `ticket_id` |  |  |
| `GORGIAS_UPDATE_ACCOUNT_SETTING` | Update Account Setting | write | `setting_id`, `value`; *description* |  |  |
| `GORGIAS_UPDATE_CUSTOMER` | Update Customer | write | `customer_id`; *address, channels, data, email, external_id, name* |  |  |
| `GORGIAS_UPDATE_MACRO` | Update Macro | write | `macro_id`; *actions, external_id, intent, language, name* |  | `macros:write` |
| `GORGIAS_UPDATE_TAG` | Update Tag | write | `tag_id`; *decoration, description, name* |  | `tags:write` |
| `GORGIAS_UPDATE_TEAM` | Update Team | write | `team_id`; *decoration, description, members, name* |  |  |
| `GORGIAS_UPDATE_TICKET` | Update Ticket | write | `ticket_id`; *assignee_team_id, assignee_user_id, customer_id, external_id, language, meta, priority, status and 1 more* |  |  |
