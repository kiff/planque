# Calendly: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

53 tools: 34 read, 19 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `CALENDLY_CANCEL_SCHEDULED_EVENT` | Cancel scheduled event | write | `event_uuid`; *reason* |  | `default` |
| `CALENDLY_CREATE_EVENT_TYPE` | Create Event Type | write | `name`, `owner`; *active, color, description, duration, duration_options, locale, locations* |  | `default` |
| `CALENDLY_CREATE_ONE_OFF_EVENT_TYPE` | Create One-Off Event Type | write | `date_setting`, `duration`, `host`, `name`; *location, timezone* |  | `default`, `event_types:write` |
| `CALENDLY_CREATE_SCHEDULING_LINK` | Create scheduling link | write | `owner`; *max_event_count, owner_type* |  | `default`, `scheduling_links:write` |
| `CALENDLY_CREATE_SHARE` | Create share | write | `event_type`; *availability__rule__rules, availability__rule__timezone, duration, duration_options, end_date, hide_location, location_configurations, max_booking_time and 3 more* |  | `default`, `scheduling_links:write` |
| `CALENDLY_CREATE_SINGLE_USE_SCHEDULING_LINK` | Create single use scheduling link | write | `max_event_count`, `owner`, `owner_type` |  | `default`, `scheduling_links:write` |
| `CALENDLY_CREATE_WEBHOOKS` | Create webhook subscription | write | `events`, `organization`, `scope`, `url`; *group, signing_key, user* |  | `routing_forms:read` |
| `CALENDLY_DELETE_INVITEE_DATA` | Delete invitee data | write | `emails` |  | `default`, `data_compliance:write` |
| `CALENDLY_DELETE_INVITEE_NO_SHOW` | Delete invitee no show | write | `uuid` |  | `default`, `scheduled_events:write` |
| `CALENDLY_DELETE_ORGANIZATION_MEMBERSHIP` | Delete organization membership | write | `uuid` |  | `organizations:write` |
| `CALENDLY_DELETE_SCHEDULED_EVENT_DATA` | Delete scheduled event data | write | `end_time`, `start_time` |  | `default`, `data_compliance:write` |
| `CALENDLY_DELETE_WEBHOOK_SUBSCRIPTION` | Delete webhook subscription | write | `webhook_uuid` |  | `organization` |
| `CALENDLY_GET_EVENT` | Get event | read | `uuid` |  | `scheduled_events:read` |
| `CALENDLY_GET_EVENT_INVITEE` | Get event invitee | read | `event_uuid`, `invitee_uuid` |  | `default`, `scheduled_events:read` |
| `CALENDLY_GET_EVENT_TYPE` | Get event type | read | *uuid* |  | `user:read` |
| `CALENDLY_GET_EVENT_TYPE_AVAILABILITY` | Get event type availability | read | `event_type`; *user* |  | `default` |
| `CALENDLY_GET_GROUP` | Get group | read | `uuid` |  | `default` |
| `CALENDLY_GET_GROUP_RELATIONSHIP` | Get group relationship | read | `uuid` |  | `default`, `groups:read` |
| `CALENDLY_GET_INVITEE_NO_SHOW` | Get invitee no show | read | `uuid` |  | `default`, `scheduled_events:read` |
| `CALENDLY_GET_ORGANIZATION` | Get organization | read | `uuid` |  | `organizations:read` |
| `CALENDLY_GET_ORGANIZATION_INVITATION` | Get organization invitation | read | `org_uuid`, `uuid` |  | `default`, `organizations:read` |
| `CALENDLY_GET_ORGANIZATION_MEMBERSHIP` | Get organization membership | read | `uuid` |  | `default`, `organizations:read` |
| `CALENDLY_GET_ROUTING_FORM` | Get routing form | read | `uuid` |  | `default`, `routing_forms:read` |
| `CALENDLY_GET_ROUTING_FORM_SUBMISSION` | Get routing form submission | read | `uuid` |  | `routing_forms:read` |
| `CALENDLY_GET_SAMPLE_WEBHOOK_DATA` | Get sample webhook data | read | `event`, `organization`, `scope`; *user* |  | `default` |
| `CALENDLY_GET_USER` | Get user | read | *uuid* |  | `user` |
| `CALENDLY_GET_USER_AVAILABILITY_SCHEDULE` | Get user availability schedule | read | `uuid` |  | `default`, `availability:read` |
| `CALENDLY_GET_WEBHOOK_SUBSCRIPTION` | Get webhook subscription | read | `webhook_uuid` |  | `organization:read`, `user:read` |
| `CALENDLY_INVITEE_NO_SHOW` | Mark invitee as no-show | write | `invitee` |  | `scheduled_events:write` |
| `CALENDLY_LIST_ACTIVITY_LOG_ENTRIES` | List activity log entries | read | `organization`; *action, actor, count, max_occurred_at, min_occurred_at, namespace, page_token, search_term and 1 more* |  | `default`, `organizations:read` |
| `CALENDLY_LIST_EVENT_INVITEES` | List event invitees | read | `uuid`; *count, email, page_token, sort, status* |  | `default`, `scheduled_events:read` |
| `CALENDLY_LIST_EVENT_TYPES` | List Event Types | read | *active, admin_managed, count, organization, page_token, sort, user, user_availability_schedule* |  | `default` |
| `CALENDLY_LIST_EVENT_TYPE_AVAILABLE_TIMES` | List event type available times | read | `end_time`, `event_type`, `start_time` |  | `default`, `availability:read` |
| `CALENDLY_LIST_EVENT_TYPE_MEMBERSHIPS` | List event type hosts | read | `event_type`; *count, page_token* |  |  |
| `CALENDLY_LIST_GROUPS` | List groups | read | `organization`; *count, page_token* |  | `default`, `groups:read` |
| `CALENDLY_LIST_GROUP_RELATIONSHIPS` | List group relationships | read | *count, group, organization, owner, page_token* |  | `default`, `groups:read` |
| `CALENDLY_LIST_ORGANIZATION_INVITATIONS` | List organization invitations | read | `uuid`; *count, email, page_token, sort, status* |  | `default`, `organizations:read` |
| `CALENDLY_LIST_ORGANIZATION_MEMBERSHIPS` | List organization memberships | read | *count, email, organization, page_token, user* |  | `organization:read` |
| `CALENDLY_LIST_OUTGOING_COMMUNICATIONS` | List outgoing communications | read | `organization`; *count, max_created_at, min_created_at, page_token* |  | `default`, `organizations:read` |
| `CALENDLY_LIST_ROUTING_FORMS` | List routing forms | read | `organization`; *count, page_token, sort* |  | `default` |
| `CALENDLY_LIST_ROUTING_FORM_SUBMISSIONS` | List routing form submissions | read | `form`; *count, page_token, sort* |  | `routing_forms:read` |
| `CALENDLY_LIST_SCHEDULED_EVENTS` | List scheduled events | read | *count, group, invitee_email, max_start_time, min_start_time, organization, page_token, sort and 2 more* |  | `scheduled_events:read` |
| `CALENDLY_LIST_USER_AVAILABILITY_SCHEDULES` | List user availability schedules | read | `user` |  |  |
| `CALENDLY_LIST_USER_BUSY_TIMES` | List user busy times | read | `end_time`, `start_time`, `user` |  | `default` |
| `CALENDLY_LIST_USER_LOCATIONS` | List User Meeting Locations | read | `user` |  | `event_types:read` |
| `CALENDLY_LIST_WEBHOOK_SUBSCRIPTIONS` | List webhook subscriptions | read | `organization`, `scope`; *count, group, page_token, sort, user* |  | `organization`, `user`, `group` |
| `CALENDLY_ORGANIZATION_INVITATION` | Invite user to organization | write | `email`, `org_uuid` |  | `default` |
| `CALENDLY_POST_INVITEE` | Create Event Invitee | write | `event_type`, `invitee`, `start_time`; *event_guests, location, questions_and_answers, tracking* |  | `default` |
| `CALENDLY_REMOVE_USER_FROM_ORGANIZATION` | Remove user from organization | write | `uuid` |  | `organizations:write` |
| `CALENDLY_REVOKE_USER_S_ORGANIZATION_INVITATION` | Revoke a user's organization invitation | write | `org_uuid`, `uuid` |  |  |
| `CALENDLY_UPDATE_EVENT_TYPE` | Update Event Type | write | `uuid`; *active, color, description, duration, duration_options, locale, locations, name* |  | `event_types:write` |
| `CALENDLY_UPDATE_EVENT_TYPE_AVAILABILITY` | Update Event Type Availability | write | `availability_rule`, `event_type`; *availability_setting, user* |  | `default` |
| `CALENDLY_WHO_AM_I` | Who Am I | read |  |  |  |
