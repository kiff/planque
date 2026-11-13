# Agent Mail: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner pastes a key from the vendor's settings). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

93 tools: 46 read, 47 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `AGENT_MAIL_BATCH_GET_MESSAGES` | Batch Get Messages | read | `inbox_id`, `message_ids` |  | `message_read` |
| `AGENT_MAIL_BATCH_UPDATE_MESSAGE_LABELS` | Batch Update Message Labels | write | `inbox_id`, `message_ids`; *add_labels, remove_labels* |  | `message_update` |
| `AGENT_MAIL_COMPLETE_FILE_UPLOAD` | Complete File Upload | write | `file_id`, `inbox_id` |  | `file_create` |
| `AGENT_MAIL_CREATE_AGENT_ARMOR_FEEDBACK` | Create Agent Armor Feedback | write | `feedback`, `inbox_id`, `message_id`; *notes* |  | `message_read` |
| `AGENT_MAIL_CREATE_CALENDAR_EVENT` | Create Calendar Event | write | `end`, `inbox_id`, `start`, `title`; *all_day, attendees, client_id, description, duration_mode, location, metadata, recurrence and 3 more* |  | `calendar_event_create` |
| `AGENT_MAIL_CREATE_DRAFT` | Create Draft | write | `inbox_id`; *attachments, bcc, cc, client_id, forward_of, html, in_reply_to, labels and 6 more* |  | `draft_create`, `message_read` |
| `AGENT_MAIL_CREATE_FILE` | Create File | write | `inbox_id`, `path`, `sha256`, `size`; *content_type, expires_at* |  | `file_create` |
| `AGENT_MAIL_CREATE_INBOX` | Create Inbox | write | *client_id, display_name, domain, username* |  |  |
| `AGENT_MAIL_CREATE_INBOX_API_KEY` | Create Inbox API Key | write | `inbox_id`; *expires_at, name, permissions* |  | `api_key_create` |
| `AGENT_MAIL_CREATE_INBOX_IN_POD` | Create Inbox in Pod | write | `pod_id`; *client_id, display_name, domain, metadata, status, username* |  | `inbox_create` |
| `AGENT_MAIL_CREATE_INBOX_POLLING_WEBHOOK` | Create Inbox Polling Webhook | write | `client_id`, `event_types`, `inbox_id` |  | `webhook_create` |
| `AGENT_MAIL_CREATE_INBOX_WEBHOOK` | Create Inbox Webhook | write | `event_types`, `inbox_id`, `url`; *client_id, headers* |  | `webhook_create` |
| `AGENT_MAIL_CREATE_LIST_ENTRY` | Create Allowlist or Blocklist Entry | write | `direction`, `entry`, `list_type`; *reason* |  | `list_entry_create` |
| `AGENT_MAIL_CREATE_POD` | Create Pod | write | *client_id, name* |  | `pod_create` |
| `AGENT_MAIL_CREATE_POLLING_WEBHOOK` | Create Polling Webhook | write | `client_id`, `event_types`; *inbox_ids, pod_ids* |  | `webhook_create` |
| `AGENT_MAIL_CREATE_WEBHOOK` | Create Webhook | write | `event_types`, `url`; *client_id, headers, inbox_ids, pod_ids* |  | `webhook_create` |
| `AGENT_MAIL_DELETE_CALENDAR_EVENT` | Delete Calendar Event | write | `event_id`, `inbox_id`; *idempotency_key, if_match, mode, send_invites* |  | `calendar_event_delete` |
| `AGENT_MAIL_DELETE_DRAFT` | Delete Draft | write | `draft_id`, `inbox_id` |  | `draft_delete` |
| `AGENT_MAIL_DELETE_FILE` | Delete File | write | `file_id`, `inbox_id` |  | `file_delete` |
| `AGENT_MAIL_DELETE_INBOX` | Delete Inbox | write | `inbox_id` |  | `inbox_delete` |
| `AGENT_MAIL_DELETE_INBOX_POLLING_WEBHOOK` | Delete Inbox Polling Webhook | write | `inbox_id`, `webhook_id` |  | `webhook_delete` |
| `AGENT_MAIL_DELETE_INBOX_WEBHOOK` | Delete Inbox Webhook | write | `inbox_id`, `webhook_id` |  | `webhook_delete` |
| `AGENT_MAIL_DELETE_LIST_ENTRY` | Delete Allowlist or Blocklist Entry | write | `direction`, `entry`, `list_type` |  | `list_entry_delete` |
| `AGENT_MAIL_DELETE_MESSAGE` | Delete Message | write | `inbox_id`, `message_id` |  | `message_delete` |
| `AGENT_MAIL_DELETE_POD` | Delete Pod | write | `pod_id` |  | `pod_delete` |
| `AGENT_MAIL_DELETE_POD_INBOX` | Delete Pod Inbox | write | `inbox_id`, `pod_id` |  | `inbox_delete` |
| `AGENT_MAIL_DELETE_POLLING_WEBHOOK` | Delete Polling Webhook | write | `webhook_id` |  | `webhook_delete` |
| `AGENT_MAIL_DELETE_THREAD` | Delete Thread | write | `inbox_id`, `thread_id` |  | `message_delete` |
| `AGENT_MAIL_DELETE_WEBHOOK` | Delete Webhook | write | `webhook_id` |  | `webhook_delete` |
| `AGENT_MAIL_FORWARD_MESSAGE` | Forward Message | write | `inbox_id`, `message_id`; *attachments, bcc, cc, headers, html, labels, reply_to, subject and 3 more* |  | `message_send` |
| `AGENT_MAIL_GET_AGENDA` | Get Agenda | read | `inbox_id`; *after, before, consistency, include_overlapping, limit, page_token* |  | `calendar_event_read` |
| `AGENT_MAIL_GET_API_KEY` | Get API Key | read | `api_key_id` |  | `api_key_read` |
| `AGENT_MAIL_GET_CALENDAR` | Get Calendar | read | `inbox_id`; *consistency* |  | `calendar_read` |
| `AGENT_MAIL_GET_CALENDAR_EVENT` | Get Calendar Event | read | `event_id`, `inbox_id`; *consistency* |  | `calendar_event_read` |
| `AGENT_MAIL_GET_CURRENT_IDENTITY` | Get Current Identity | read |  |  |  |
| `AGENT_MAIL_GET_DOMAIN` | Get Domain | read | `domain_id` |  | `domain_read` |
| `AGENT_MAIL_GET_DOMAIN_SETUP_LINK` | Get Domain Setup Link | read | `domain_id` |  | `domain_read` |
| `AGENT_MAIL_GET_DOMAIN_ZONE_FILE` | Get Domain Zone File | read | `domain_id` |  | `domain_read` |
| `AGENT_MAIL_GET_DRAFT` | Get Draft | read | `draft_id`, `inbox_id` |  | `draft_read` |
| `AGENT_MAIL_GET_DRAFT_ATTACHMENT` | Get Draft Attachment | read | `attachment_id`, `draft_id` |  | `draft_read` |
| `AGENT_MAIL_GET_FILE` | Get File | read | `file_id`, `inbox_id` |  | `file_read` |
| `AGENT_MAIL_GET_INBOX` | Get Inbox | read | `inbox_id` |  | `inbox_read` |
| `AGENT_MAIL_GET_INBOX_POLLING_WEBHOOK` | Get Inbox Polling Webhook | read | `inbox_id`, `webhook_id` |  | `webhook_read` |
| `AGENT_MAIL_GET_INBOX_WEBHOOK` | Get Inbox Webhook | read | `inbox_id`, `webhook_id` |  | `webhook_read` |
| `AGENT_MAIL_GET_INBOX_WEBHOOK_HEADERS` | Get Inbox Webhook Headers | read | `inbox_id`, `webhook_id` |  | `webhook_read` |
| `AGENT_MAIL_GET_MESSAGE` | Get Message | read | `inbox_id`, `message_id` |  |  |
| `AGENT_MAIL_GET_MESSAGE_ATTACHMENT` | Get Message Attachment | read | `attachment_id`, `inbox_id`, `message_id` |  | `message_read` |
| `AGENT_MAIL_GET_ORGANIZATION` | Get Organization | read |  |  |  |
| `AGENT_MAIL_GET_POD` | Get Pod | read | `pod_id` |  | `pod_read` |
| `AGENT_MAIL_GET_POD_INBOX` | Get Pod Inbox | read | `inbox_id`, `pod_id` |  | `inbox_read` |
| `AGENT_MAIL_GET_POLLING_WEBHOOK` | Get Polling Webhook | read | `webhook_id` |  | `webhook_read` |
| `AGENT_MAIL_GET_RAW_MESSAGE` | Get Raw Message | write | `inbox_id`, `message_id` |  | `message_read` |
| `AGENT_MAIL_GET_THREAD` | Get Thread | read | `inbox_id`, `thread_id`; *limit, page_token* |  | `message_read` |
| `AGENT_MAIL_GET_THREAD_ATTACHMENT` | Get Thread Attachment | read | `attachment_id`, `thread_id` |  | `message_read` |
| `AGENT_MAIL_GET_WEBHOOK` | Get Webhook | read | `webhook_id` |  | `webhook_read` |
| `AGENT_MAIL_GET_WEBHOOK_HEADERS` | Get Webhook Headers | read | `webhook_id` |  | `webhook_read` |
| `AGENT_MAIL_LIST_CALENDAR_EVENTS` | List Calendar Events | read | `inbox_id`; *limit, page_token* |  | `calendar_event_read` |
| `AGENT_MAIL_LIST_CALENDAR_EVENT_INSTANCES` | List Calendar Event Instances | read | `event_id`, `inbox_id`; *after, before, consistency, include_overlapping, limit, page_token* |  | `calendar_event_read` |
| `AGENT_MAIL_LIST_DOMAINS` | List Domains | read | *ascending, limit, page_token* |  | `domain_read` |
| `AGENT_MAIL_LIST_DRAFTS` | List Drafts | read | `inbox_id`; *after, ascending, before, labels, limit, page_token* |  | `draft_read` |
| `AGENT_MAIL_LIST_FILES` | List Files | read | `inbox_id`; *ascending, limit, page_token, path_prefix* |  | `file_read` |
| `AGENT_MAIL_LIST_INBOXES` | List Inboxes | read | *ascending, limit, page_token* |  |  |
| `AGENT_MAIL_LIST_INBOX_API_KEYS` | List Inbox API Keys | read | `inbox_id`; *limit, page_token* |  | `api_key_read` |
| `AGENT_MAIL_LIST_INBOX_EVENTS` | List Inbox Events | read | `inbox_id`; *ascending, limit, page_token* |  | `message_read` |
| `AGENT_MAIL_LIST_INBOX_POLLING_WEBHOOKS` | List Inbox Polling Webhooks | read | `inbox_id`; *ascending, limit, page_token* |  | `webhook_read` |
| `AGENT_MAIL_LIST_INBOX_WEBHOOKS` | List Inbox Webhooks | read | `inbox_id`; *ascending, limit, page_token* |  | `webhook_read` |
| `AGENT_MAIL_LIST_LIST_ENTRIES` | List Allowlist or Blocklist Entries | read | `direction`, `list_type`; *entry, limit, page_token* |  | `list_entry_read` |
| `AGENT_MAIL_LIST_MESSAGES` | List Messages | read | `inbox_id`; *after, ascending, before, include_spam, labels, limit, page_token* |  |  |
| `AGENT_MAIL_LIST_PODS` | List Pods | read | *ascending, limit, page_token* |  | `pod_read` |
| `AGENT_MAIL_LIST_POD_INBOXES` | List Pod Inboxes | read | `pod_id`; *ascending, limit, page_token, query* |  | `inbox_read` |
| `AGENT_MAIL_LIST_POLLING_WEBHOOKS` | List Polling Webhooks | read | *ascending, limit, page_token* |  | `webhook_read` |
| `AGENT_MAIL_LIST_THREADS` | List Threads | read | `inbox_id`; *after, ascending, before, include_blocked, include_spam, include_trash, include_unauthenticated, labels and 5 more* |  | `message_read` |
| `AGENT_MAIL_LIST_WEBHOOKS` | List Webhooks | read | *ascending, limit, page_token* |  | `webhook_read` |
| `AGENT_MAIL_QUERY_INBOX_EVENT_METRICS` | Query Inbox Event Metrics | read | `inbox_id`; *descending, end, event_types, limit, period, start* |  | `metrics_read` |
| `AGENT_MAIL_QUERY_METRICS` | Query Metrics | read | `metric_kind`; *descending, end, event_types, limit, period, start, usage_types* |  | `metrics_read` |
| `AGENT_MAIL_QUERY_RATE_METRICS` | Query Rate Metrics | read | *descending, end, limit, period, rate_types, start, window* |  | `metrics_read` |
| `AGENT_MAIL_RELEASE_HELD_MESSAGE` | Release Held Message | write | `inbox_id`, `message_id`, `reason` |  | `agent_armor_release`, `message_read` |
| `AGENT_MAIL_REPLY_TO_MESSAGE` | Reply To Message | write | `inbox_id`, `message_id`; *attachments, bcc, cc, headers, html, labels, reply_all, reply_to and 3 more* |  | `message_send` |
| `AGENT_MAIL_RESPOND_TO_CALENDAR_EVENT` | Respond to Calendar Event | write | `event_id`, `inbox_id`, `status`; *comment, if_match, send_reply* |  | `calendar_event_update` |
| `AGENT_MAIL_SEND_DRAFT` | Send Draft | write | `draft_id`, `inbox_id`; *add_labels, remove_labels* |  | `draft_send` |
| `AGENT_MAIL_SEND_EMAIL` | Send Email | write | `inbox_id`, `subject`, `to`; *bcc, cc, html, labels, reply_to, text* |  |  |
| `AGENT_MAIL_SEND_MESSAGE` | Send Message | write | `inbox_id`; *attachments, bcc, cc, headers, html, labels, reply_to, subject and 3 more* |  | `message_send` |
| `AGENT_MAIL_UPDATE_CALENDAR` | Update Calendar | write | `inbox_id`, `timezone`; *if_match* |  | `calendar_update` |
| `AGENT_MAIL_UPDATE_CALENDAR_EVENT` | Update Calendar Event | write | `event_id`, `inbox_id`; *all_day, attendees, description, duration_mode, end, if_match, location, metadata and 7 more* |  | `calendar_event_update` |
| `AGENT_MAIL_UPDATE_DRAFT` | Update Draft | write | `draft_id`, `inbox_id`; *add_attachments, add_labels, bcc, cc, html, remove_attachments, remove_labels, reply_to and 4 more* |  | `draft_update` |
| `AGENT_MAIL_UPDATE_INBOX` | Update Inbox | write | `inbox_id`; *display_name, metadata, status* |  | `inbox_update` |
| `AGENT_MAIL_UPDATE_INBOX_WEBHOOK` | Update Inbox Webhook | write | `inbox_id`, `webhook_id`; *enabled, event_types* |  | `webhook_update` |
| `AGENT_MAIL_UPDATE_INBOX_WEBHOOK_HEADERS` | Update Inbox Webhook Headers | write | `inbox_id`, `webhook_id`; *headers, remove_headers* |  | `webhook_update` |
| `AGENT_MAIL_UPDATE_MESSAGE_LABELS` | Update Message Labels | write | `inbox_id`, `message_id`; *add_labels, remove_labels* |  | `message_update` |
| `AGENT_MAIL_UPDATE_POD_INBOX` | Update Pod Inbox | write | `inbox_id`, `pod_id`; *display_name, metadata, status* |  | `inbox_update` |
| `AGENT_MAIL_UPDATE_THREAD_LABELS` | Update Thread Labels | write | `thread_id`; *add_labels, remove_labels* |  | `message_update` |
| `AGENT_MAIL_UPDATE_WEBHOOK` | Update Webhook | write | `webhook_id`; *add_inbox_ids, add_pod_ids, enabled, event_types, remove_inbox_ids, remove_pod_ids* |  | `webhook_update` |
| `AGENT_MAIL_UPDATE_WEBHOOK_HEADERS` | Update Webhook Headers | write | `webhook_id`; *headers, remove_headers* |  | `webhook_update` |
