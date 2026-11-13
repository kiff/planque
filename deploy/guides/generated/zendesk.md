# Zendesk: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

480 tools: 264 read, 216 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `ZENDESK_APPLY_ZENDESK_MACRO` | Apply Zendesk Macro Preview | read | `macro_id`, `ticket_id` |  | `tickets:read`, `read` |
| `ZENDESK_AUTOCOMPLETE_PROBLEMS` | Autocomplete Problems | read | `text` |  | `read` |
| `ZENDESK_AUTOCOMPLETE_USERS` | Autocomplete Users | read | *field_id, filter, name, phone, source* |  | `users:read`, `read` |
| `ZENDESK_CHECK_HOST_MAPPING_VALIDITY_FOR_EXISTING_BRAND` | Check Host Mapping Validity for Existing Brand | read | `brand_id` |  | `read` |
| `ZENDESK_COUNT_DELETED_USERS` | Count Deleted Users | read |  |  | `users:read`, `read` |
| `ZENDESK_COUNT_USER_CCD_TICKETS` | Count User CCD Tickets | read | `user_id` |  | `users:read`, `read` |
| `ZENDESK_COUNT_ZENDESK_CUSTOM_OBJECT_RECORDS` | Count Custom Object Records | write | `custom_object_key` |  | `read` |
| `ZENDESK_COUNT_ZENDESK_ORGANIZATIONS` | Count Zendesk Organizations | read |  |  | `organizations:read`, `read` |
| `ZENDESK_CREATE_ACCESS_RULE` | Create Access Rule | write | `access_rule`, `custom_object_key` |  | `write` |
| `ZENDESK_CREATE_APPS_INSTALLATION` | Install Zendesk App | write | `app_id`; *enabled, group_restrictions, recurring_payment, role_restrictions, settings, settings_objects* |  | `apps:write`, `write` |
| `ZENDESK_CREATE_APPS_NOTIFY` | Send Notification to App | write | `app_id`, `event`; *agent_id, body* |  | `apps:write`, `write` |
| `ZENDESK_CREATE_ARTICLE_IN_SECTION` | Create Article in Section | write | `section_id`; *author_id, body, comments_disabled, content_tag_ids, draft, label_names, locale, notify_subscribers and 7 more* |  | `hc:write`, `write` |
| `ZENDESK_CREATE_AUTOCOMPLETE_TAGS` | Create Autocomplete Tags | read | *name* |  | `write` |
| `ZENDESK_CREATE_AUTOMATION` | Create Automation | write | `conditions`, `title`; *actions, active, position* |  | `automations:write`, `write` |
| `ZENDESK_CREATE_BOOKMARK` | Create Bookmark | write | `ticket_id` |  | `tickets:write`, `write` |
| `ZENDESK_CREATE_COMMENT_REDACTION` | Redact Zendesk Comment | write | `text`, `ticket_comment_id`, `ticket_id` |  | `tickets:write`, `write` |
| `ZENDESK_CREATE_CUSTOM_OBJECT_BULK_JOB` | Create Custom Object Bulk Job | write | `action`, `custom_object_key`, `items` |  | `write` |
| `ZENDESK_CREATE_CUSTOM_OBJECT_RECORD` | Create Custom Object Record | write | `custom_object_key`, `name`; *custom_object_fields, external_id* |  | `write` |
| `ZENDESK_CREATE_CUSTOM_OBJECT_RECORD_ATTACHMENT` | Create Custom Object Record Attachment | write | `custom_object_key`, `file`, `record_id` |  | `write` |
| `ZENDESK_CREATE_CUSTOM_STATUSES` | Create Custom Ticket Status | write | `agent_label`, `description`, `end_user_label`, `status_category`; *active, end_user_description* |  | `write` |
| `ZENDESK_CREATE_DELETION_SCHEDULE` | Create Deletion Schedule | write | `object`, `title`; *active, conditions, description* |  | `write` |
| `ZENDESK_CREATE_DYNAMIC_CONTENT_ITEMS` | Create Zendesk Dynamic Content Item | write | `default_locale_id`, `name`, `variants` |  | `dynamic_content:write`, `write` |
| `ZENDESK_CREATE_DYNAMIC_CONTENT_ITEMS_VARIANTS` | Create Zendesk Dynamic Content Item Variant | write | `content`, `dynamic_content_item_id`, `locale_id`; *active, default* |  | `dynamic_content:write`, `write` |
| `ZENDESK_CREATE_END_USER_IDENTITY` | Create End User Identity | write | `type`, `user_id`, `value`; *primary* |  | `users:write`, `write` |
| `ZENDESK_CREATE_GROUP_MEMBERSHIPS` | Create Zendesk Group Membership | write | `group_id`, `user_id`; *default* |  | `write` |
| `ZENDESK_CREATE_IMPORTS_TICKETS_CREATE_MANY` | Bulk Import Tickets | write | `tickets`; *archive_immediately* |  | `tickets:write`, `write` |
| `ZENDESK_CREATE_MACRO` | Create Macro | write | `actions`, `title`; *active, attachments, description, position, restriction* |  | `macros:write`, `write` |
| `ZENDESK_CREATE_MANY_DYNAMIC_CONTENT_ITEMS_VARIANTS` | Create Many Dynamic Content Variants | write | `dynamic_content_item_id`, `variants` |  | `dynamic_content:write`, `write` |
| `ZENDESK_CREATE_MANY_TICKETS` | Create Many Zendesk Tickets | write | `tickets` |  | `tickets:write`, `write` |
| `ZENDESK_CREATE_MANY_USERS` | Create Many Zendesk Users | write | `users`; *skip_verify_email* |  | `users:write`, `write` |
| `ZENDESK_CREATE_OBJECT_TRIGGER` | Create Object Trigger | write | `actions`, `conditions`, `custom_object_key`, `title`; *active, description, position* |  | `write` |
| `ZENDESK_CREATE_ORGANIZATIONS_TAGS` | Set Organization Tags | write | `organization_id`, `tags` |  | `organizations:write`, `write` |
| `ZENDESK_CREATE_ORGANIZATION_FIELD` | Create Organization Field | write | `key`, `title`, `type`; *active, custom_field_options, description, position, regexp_for_validation, relationship_filter, relationship_target_type, tag* |  | `organizations:write`, `write` |
| `ZENDESK_CREATE_ORGANIZATION_MEMBERSHIPS` | Create Organization Membership | write | `organization_id`, `user_id`; *default* |  | `organizations:write`, `write` |
| `ZENDESK_CREATE_ORGANIZATION_MEMBERSHIPS_CREATE_MANY` | Create Many Organization Memberships | write | `organization_memberships` |  | `write` |
| `ZENDESK_CREATE_ORGANIZATION_MERGE` | Merge Zendesk Organizations | write | `organization_id`, `winner_id` |  | `organizations:write`, `write` |
| `ZENDESK_CREATE_ORGANIZATION_SUBSCRIPTION` | Create Organization Subscription | write | `organization_id`, `user_id` |  | `organizations:write`, `write` |
| `ZENDESK_CREATE_OR_UPDATE_MANY_USERS` | Create Or Update Many Users | write | `users` |  | `users:write`, `write` |
| `ZENDESK_CREATE_OR_UPDATE_ORGANIZATION` | Create Or Update Organization | write | `name`; *details, domain_names, external_id, group_id, notes, organization_fields, shared_comments, shared_tickets and 1 more* |  | `organizations:write`, `write` |
| `ZENDESK_CREATE_OR_UPDATE_USER` | Create or Update User | write | `user` |  | `users:write`, `write` |
| `ZENDESK_CREATE_PUSH_NOTIFICATION_DEVICES_DESTROY_MANY` | Bulk Unregister Push Notification Devices | write | `push_notification_devices` |  | `write` |
| `ZENDESK_CREATE_RESOURCE_COLLECTIONS` | Create Resource Collection | write | `actions`, `title` |  | `tickets:write`, `write` |
| `ZENDESK_CREATE_SIDE_CONVERSATION` | Create Side Conversation | write | `body`, `ticket_id`, `to`; *attachment_ids, external_ids, html_body, subject* |  | `tickets:write`, `write` |
| `ZENDESK_CREATE_SUPPORT_ADDRESS` | Create Support Address | write | `email`, `name`; *brand_id, default* |  | `write` |
| `ZENDESK_CREATE_SUSPENDED_TICKETS_EXPORT` | Export Suspended Tickets | read |  |  | `tickets:read`, `read` |
| `ZENDESK_CREATE_TICKET_CONTENT_PIN` | Create Ticket Content Pin | write | `content_id`, `content_type`, `ticket_id`; *locale* |  | `tickets:write`, `write` |
| `ZENDESK_CREATE_TICKET_FIELD` | Create Ticket Field | write | `title`, `type`; *active, agent_can_edit, agent_description, collapsed_for_agents, custom_field_options, description, editable_in_portal, position and 8 more* |  | `tickets:write`, `write` |
| `ZENDESK_CREATE_TICKET_FIELD_OPTION` | Create Ticket Field Option | write | `name`, `ticket_field_id`, `value`; *allow_solving, position* |  | `tickets:write`, `write` |
| `ZENDESK_CREATE_TRIGGER` | Create Trigger | write | `actions`, `conditions`, `title`; *active, category_id, description* |  | `triggers:write`, `write` |
| `ZENDESK_CREATE_TRIGGER_CATEGORIES_JOB` | Create Trigger Categories Batch Job | write | `trigger_category_ids`; *action* |  | `write` |
| `ZENDESK_CREATE_TRIGGER_CATEGORY` | Create Zendesk Trigger Category | write | `name`; *position* |  | `triggers:write`, `write` |
| `ZENDESK_CREATE_USERS_LOGOUT_MANY` | Logout Many Zendesk Users | write | `ids` |  | `users:write`, `write` |
| `ZENDESK_CREATE_USERS_ORGANIZATION_MEMBERSHIPS` | Create User Organization Membership | write | `organization_id`, `user_id`; *default* |  | `organizations:write`, `users:write`, `write` |
| `ZENDESK_CREATE_USERS_TAGS` | Create Users Tags | write | `tags`, `user_id` |  | `users:write`, `write` |
| `ZENDESK_CREATE_USER_FIELD_OPTION` | Create User Field Option | write | `name`, `user_field_id`, `value`; *allow_solving, position* |  | `users:write`, `write` |
| `ZENDESK_CREATE_USER_IDENTITY` | Create User Identity | write | `type`, `user_id`, `value`; *primary* |  | `users:write`, `write` |
| `ZENDESK_CREATE_USER_PROFILES_EVENTS` | Create User Profile Event | write | `event`, `profile` |  | `users:write`, `write` |
| `ZENDESK_CREATE_VIEW` | Create View | write | `conditions`, `title`; *active, description, output, restriction* |  | `write` |
| `ZENDESK_CREATE_VIEWS_PREVIEW` | Preview View | read | `view`; *columns, exclude, include, page, per_page* |  | `write` |
| `ZENDESK_CREATE_VIEWS_PREVIEW_COUNT` | Preview View Ticket Count | read | `all_conditions`; *any_conditions, output* |  | `write` |
| `ZENDESK_CREATE_WEBHOOK` | Create Webhook | write | `endpoint`, `http_method`, `name`, `request_format`, `status`; *authentication, custom_headers, description, subscriptions* |  | `webhooks:write`, `write` |
| `ZENDESK_CREATE_ZENDESK_ATTACHMENTS` | Upload Zendesk Attachment | write | `file` |  | `tickets:write`, `write` |
| `ZENDESK_CREATE_ZENDESK_CUSTOM_OBJECT` | Create Zendesk Custom Object | write | `custom_object` |  | `write` |
| `ZENDESK_CREATE_ZENDESK_CUSTOM_OBJECT_FIELD` | Create Zendesk Custom Object Field | write | `custom_object_key`, `key`, `title`, `type`; *active, custom_field_options, description, position, properties, regexp_for_validation, relationship_target_type, required and 1 more* |  | `write` |
| `ZENDESK_CREATE_ZENDESK_GROUP` | Create Zendesk Group | write | `name`; *description, is_public* |  | `write` |
| `ZENDESK_CREATE_ZENDESK_IMPORT_TICKET` | Import Zendesk Ticket | write | `description`, `subject`; *archive_immediately, assignee_id, collaborator_ids, collaborators, created_at, custom_fields, custom_status_id, due_at and 11 more* |  | `tickets:write`, `write` |
| `ZENDESK_CREATE_ZENDESK_ORGANIZATION` | Create Zendesk Organization | write | `name`; *details, domain_names, external_id, group_id, notes, shared_comments, shared_tickets, tags* |  | `organizations:write`, `write` |
| `ZENDESK_CREATE_ZENDESK_REQUEST` | Create Zendesk Request | write | `comment_body`, `subject`; *custom_fields, external_id, priority, requester_email, requester_id, requester_name, status, tags and 2 more* |  | `requests:write`, `write` |
| `ZENDESK_CREATE_ZENDESK_TARGET` | Create Zendesk Target | write | `target_type`, `title`; *account_name, active, api_id, api_token, attribute, email, from_field, group_id and 12 more* |  | `targets:write`, `write` |
| `ZENDESK_CREATE_ZENDESK_TICKET` | Create Zendesk Ticket | write | `description`, `subject`; *additional_collaborators, assignee_id, collaborator_ids, collaborators, custom_fields, due_at, email_ccs, external_id and 11 more* |  | `tickets:write`, `write` |
| `ZENDESK_CREATE_ZENDESK_TICKET_OR_VOICEMAIL_TICKET` | Create Ticket or Voicemail Ticket | write | `description`, `subject`; *assignee_id, custom_fields, display_to_agent, external_id, group_id, priority, requester_email, requester_id and 6 more* |  | `tickets:write`, `write` |
| `ZENDESK_CREATE_ZENDESK_USER` | Create Zendesk User | write | `user`; *skip_verify_email* |  | `users:write`, `write` |
| `ZENDESK_CREATE_ZENDESK_USER_FIELD` | Create Zendesk User Field | write | `key`, `title`, `type`; *active, custom_field_options, description, position, raw_description, raw_title, regexp_for_validation, relationship_filter and 2 more* |  | `users:write`, `write` |
| `ZENDESK_DELETE_ACCESS_RULE` | Delete Access Rule | write | `custom_object_key`, `id` |  | `write` |
| `ZENDESK_DELETE_AUTOMATION` | Delete Zendesk Automation | write | `automation_id` |  | `automations:write`, `write` |
| `ZENDESK_DELETE_BUSINESS_HOURS_SCHEDULE` | Delete Zendesk Business Hours Schedule | write | `schedule_id` |  | `write` |
| `ZENDESK_DELETE_CUSTOM_OBJECT_RECORD_ATTACHMENT` | Delete Custom Object Record Attachment | write | `custom_object_key`, `id`, `record_id` |  | `write` |
| `ZENDESK_DELETE_CUSTOM_OBJECT_RECORD_BY_EXTERNAL_ID_OR_NAME` | Delete Custom Object Record | write | `custom_object_key`; *external_id, name* |  | `write` |
| `ZENDESK_DELETE_DELETED_TICKET` | Delete Deleted Ticket Permanently | write | `ticket_id` |  | `tickets:write`, `write` |
| `ZENDESK_DELETE_DELETION_SCHEDULE` | Delete Zendesk Deletion Schedule | write | `deletion_schedule_id` |  | `write` |
| `ZENDESK_DELETE_GROUP_MEMBERSHIP` | Delete Zendesk Group Membership | write | `group_membership_id` |  | `write` |
| `ZENDESK_DELETE_GROUP_MEMBERSHIPS_DESTROY_MANY` | Bulk Delete Zendesk Group Memberships | write | *ids* |  | `write` |
| `ZENDESK_DELETE_MACRO` | Delete Zendesk Macro | write | `macro_id` |  | `macros:write`, `write` |
| `ZENDESK_DELETE_MACROS_DESTROY_MANY` | Bulk Delete Zendesk Macros | write | `ids` |  | `macros:write`, `write` |
| `ZENDESK_DELETE_MANY_AUTOMATIONS` | Bulk Delete Zendesk Automations | write | `ids` |  | `automations:write`, `write` |
| `ZENDESK_DELETE_MANY_OBJECT_TRIGGERS` | Bulk Delete Zendesk Custom Object Triggers | write | `custom_object_key`, `ids` |  | `write` |
| `ZENDESK_DELETE_MANY_ORGANIZATION_MEMBERSHIPS` | Bulk Delete Zendesk Organization Memberships | write | `ids` |  | `write` |
| `ZENDESK_DELETE_MANY_VIEWS` | Bulk Delete Zendesk Views | write | `ids` |  | `write` |
| `ZENDESK_DELETE_ORGANIZATIONS_TAGS` | Delete Organization Tags | write | `organization_id` |  | `organizations:write`, `write` |
| `ZENDESK_DELETE_ORGANIZATION_FIELD` | Delete Zendesk Organization Field | write | `organization_field_id` |  | `write` |
| `ZENDESK_DELETE_ORGANIZATION_MEMBERSHIP` | Delete Organization Membership | write | `organization_membership_id` |  | `write` |
| `ZENDESK_DELETE_ORGANIZATION_SUBSCRIPTION` | Delete Zendesk Organization Subscription | write | `organization_subscription_id` |  | `organizations:write`, `write` |
| `ZENDESK_DELETE_RESOURCE_COLLECTIONS` | Delete Resource Collection | write | `resource_collection_id` |  | `write` |
| `ZENDESK_DELETE_SUSPENDED_TICKETS_DESTROY_MANY` | Bulk Delete Zendesk Suspended Tickets | write | `ids` |  | `tickets:write`, `write` |
| `ZENDESK_DELETE_TARGET` | Delete Zendesk Target | write | `target_id` |  | `targets:write`, `write` |
| `ZENDESK_DELETE_TICKETS_TAGS` | Delete Ticket Tags | write | `ticket_id`; *tags* |  | `tickets:write`, `write` |
| `ZENDESK_DELETE_TICKET_FIELD` | Delete Zendesk Ticket Field | write | `ticket_field_id`; *creator* |  | `write` |
| `ZENDESK_DELETE_TICKET_FIELD_OPTION` | Delete Ticket Field Option | write | `ticket_field_id`, `ticket_field_option_id` |  | `write` |
| `ZENDESK_DELETE_TRIGGER` | Delete Zendesk Trigger | write | `trigger_id` |  | `triggers:write`, `write` |
| `ZENDESK_DELETE_TRIGGER_CATEGORY` | Delete Zendesk Trigger Category | write | `trigger_category_id` |  | `write` |
| `ZENDESK_DELETE_USERS_DESTROY_MANY` | Bulk Delete Zendesk Users | write | *external_ids, ids* |  | `users:write`, `write` |
| `ZENDESK_DELETE_USERS_GROUP_MEMBERSHIPS` | Delete Zendesk User's Group Membership | write | `group_membership_id`, `user_id` |  | `users:write`, `write` |
| `ZENDESK_DELETE_USERS_IDENTITY` | Delete User Identity | write | `user_id`, `user_identity_id` |  | `users:write`, `write` |
| `ZENDESK_DELETE_USERS_ME_LOGOUT` | Delete Authenticated Session | write |  |  | `users:write`, `write` |
| `ZENDESK_DELETE_USERS_ORGANIZATION_MEMBERSHIPS` | Delete User's Organization Membership | write | `organization_membership_id`, `user_id` |  | `users:write`, `write` |
| `ZENDESK_DELETE_USERS_SESSIONS` | Delete All User Sessions | write | `user_id` |  | `users:write`, `write` |
| `ZENDESK_DELETE_USERS_TAGS` | Delete User Tags | write | `user_id`; *tags* |  | `users:write`, `write` |
| `ZENDESK_DELETE_USER_FIELD` | Delete Zendesk User Field | write | `user_field_id` |  | `users:write`, `write` |
| `ZENDESK_DELETE_USER_FIELD_OPTION` | Delete User Field Option | write | `user_field_id`, `user_field_option_id` |  | `users:write`, `write` |
| `ZENDESK_DELETE_USER_PERMANENTLY` | Permanently Delete User | write | `deleted_user_id` |  | `users:write`, `write` |
| `ZENDESK_DELETE_USER_PROFILE` | Delete Zendesk User Profile | write | `profile_id` |  | `users:write`, `write` |
| `ZENDESK_DELETE_USER_SESSION` | Delete Zendesk User Session | write | `session_id`, `user_id` |  | `users:write`, `write` |
| `ZENDESK_DELETE_VIEW` | Delete Zendesk View | write | `view_id` |  | `write` |
| `ZENDESK_DELETE_WEBHOOK` | Delete Webhook | write | `webhook_id` |  | `webhooks:write`, `write` |
| `ZENDESK_DELETE_ZENDESK_APP_INSTALLATION` | Delete Zendesk App Installation | write | `app_installation_id` |  | `apps:write`, `write` |
| `ZENDESK_DELETE_ZENDESK_BOOKMARK` | Delete Zendesk Bookmark | write | `bookmark_id` |  | `write` |
| `ZENDESK_DELETE_ZENDESK_CUSTOM_OBJECT` | Delete Zendesk Custom Object | write | `custom_object_key` |  | `write` |
| `ZENDESK_DELETE_ZENDESK_CUSTOM_OBJECT_FIELD` | Delete Zendesk Custom Object Field | write | `custom_object_field_key_or_id`, `custom_object_key` |  | `write` |
| `ZENDESK_DELETE_ZENDESK_CUSTOM_OBJECT_RECORD` | Delete Zendesk Custom Object Record | write | `custom_object_key`, `custom_object_record_id` |  | `write` |
| `ZENDESK_DELETE_ZENDESK_CUSTOM_STATUS` | Delete Zendesk Custom Status | write | `custom_status_id` |  | `tickets:write`, `write` |
| `ZENDESK_DELETE_ZENDESK_DYNAMIC_CONTENT_ITEM` | Delete Zendesk Dynamic Content Item | write | `dynamic_content_item_id` |  | `dynamic_content:write`, `write` |
| `ZENDESK_DELETE_ZENDESK_DYNAMIC_CONTENT_ITEM_VARIANT` | Delete Zendesk Dynamic Content Item Variant | write | `dynamic_content_item_id`, `dynamic_content_variant_id` |  | `dynamic_content:write`, `write` |
| `ZENDESK_DELETE_ZENDESK_OBJECT_TRIGGER` | Delete Zendesk Custom Object Trigger | write | `custom_object_key`, `trigger_id` |  | `write` |
| `ZENDESK_DELETE_ZENDESK_ORGANIZATION` | Delete Zendesk Organization | write | `organization_id` |  | `organizations:write`, `write` |
| `ZENDESK_DELETE_ZENDESK_SUPPORT_ADDRESS` | Delete Zendesk Support Address | write | `support_address_id` |  | `write` |
| `ZENDESK_DELETE_ZENDESK_TICKET` | Delete Zendesk Ticket | write | `ticket_id` |  | `tickets:read`, `tickets:write`, `read`, `write` |
| `ZENDESK_DELETE_ZENDESK_USER` | Delete Zendesk User | write | `user_id` |  | `users:write`, `write` |
| `ZENDESK_DESTROY_MANY_DELETED_TICKETS` | Bulk Delete Soft-Deleted Zendesk Tickets | write | `ids` |  | `tickets:write`, `write` |
| `ZENDESK_DESTROY_MANY_ORGANIZATIONS` | Bulk Delete Zendesk Organizations | write | *external_ids, ids* |  | `organizations:write`, `write` |
| `ZENDESK_DETECT_BEST_LOCALE` | Detect Best Locale | read |  |  | `read` |
| `ZENDESK_DOWNLOAD_ATTACHMENT` | Download Attachment | read | `attachment_id` |  | `tickets:read`, `read` |
| `ZENDESK_DOWNLOAD_CUSTOM_OBJECT_RECORD_ATTACHMENT` | Download Custom Object Record Attachment | write | `custom_object_key`, `id`, `record_id`; *inline* |  | `read` |
| `ZENDESK_GET_ABOUT_ME` | Get Zendesk About Me | read |  |  | `users:read`, `read` |
| `ZENDESK_GET_ACCOUNT_SETTINGS` | Get Account Settings | read | *authenticity_token* |  | `read` |
| `ZENDESK_GET_ACTIVE_TRIGGERS` | List Active Triggers | read | *category_id, page, per_page, sort, sort_by, sort_order* |  | `triggers:read`, `read` |
| `ZENDESK_GET_ACTIVITIES_COUNT` | Count Activities | read |  |  | `tickets:read`, `read` |
| `ZENDESK_GET_ALL_ZENDESK_ORGANIZATIONS` | Get Zendesk Organizations | read |  |  | `organizations:read`, `read` |
| `ZENDESK_GET_APP` | Get App | read | `app_id` |  | `apps:read`, `read` |
| `ZENDESK_GET_APPS_INSTALLATIONS_REQUIREMENTS` | Get App Installation Requirements | read | `app_installation_id` |  | `apps:read`, `read` |
| `ZENDESK_GET_APPS_LOCATION_INSTALLATIONS` | List Apps Location Installations | read |  |  | `apps:read`, `read` |
| `ZENDESK_GET_APPS_PUBLIC_KEY` | Get App Public Key | read | `app_id` |  | `apps:read`, `read` |
| `ZENDESK_GET_APPS_PUBLIC_KEY_PEM` | Get App Public Key PEM | read | `app_id` |  | `apps:read`, `read` |
| `ZENDESK_GET_ARTICLE` | Get Article | read | `article_id`; *include* |  | `read` |
| `ZENDESK_GET_ARTICLE_ATTACHMENT` | Get Article Attachment | read | `article_attachment_id` |  | `read` |
| `ZENDESK_GET_ARTICLE_TRANSLATION` | Get Article Translation | read | `article_id`, `locale` |  | `read` |
| `ZENDESK_GET_ATTACHMENT` | Get Attachment | read | `attachment_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_AUTOCOMPLETE_TAGS` | Get Autocomplete Tags | read | *name* |  | `read` |
| `ZENDESK_GET_AUTOMATION` | Show Automation | read | `automation_id` |  | `automations:read`, `read` |
| `ZENDESK_GET_AUTOMATIONS` | List Automations | read | *active, include, page, per_page, sort* |  | `automations:read`, `read` |
| `ZENDESK_GET_AUTOMATIONS_SEARCH` | Search Automations | read | `query`; *active, include, sort_by, sort_order* |  | `automations:read`, `read` |
| `ZENDESK_GET_BRANDS` | List Brands | read | *page_after, page_size, per_page, sort* |  | `read` |
| `ZENDESK_GET_BRANDS2` | Get Zendesk Brand | read | `brand_id` |  | `read` |
| `ZENDESK_GET_BRANDS_AGENTS` | List Agents by Brand | read | `brand_id`; *page, per_page, sort* |  | `read` |
| `ZENDESK_GET_CURRENT_OAUTH_TOKEN` | Get Current OAuth Token | read |  |  | `read` |
| `ZENDESK_GET_CUSTOM_OBJECTS_LIMIT` | Get Custom Objects Limit | write |  |  | `read` |
| `ZENDESK_GET_CUSTOM_OBJECT_FIELDS_LIMIT` | Get Custom Object Fields Limit | write | `custom_object_key` |  | `read` |
| `ZENDESK_GET_CUSTOM_OBJECT_RECORDS_LIMIT` | Get Custom Object Records Limit | write |  |  | `read` |
| `ZENDESK_GET_CUSTOM_STATUS` | Get Custom Status | write | `custom_status_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_CUSTOM_STATUSES` | List Custom Ticket Statuses | write | *active, default, status_categories* |  | `tickets:read`, `read` |
| `ZENDESK_GET_DELETED_TICKETS` | List Deleted Tickets | read | *page, per_page, sort_by, sort_order, support_type_scope* |  | `tickets:read`, `read` |
| `ZENDESK_GET_DELETED_USER` | Get Deleted User | read | `deleted_user_id` |  | `users:read`, `read` |
| `ZENDESK_GET_DELETION_SCHEDULE` | Get Deletion Schedule | read | `deletion_schedule_id` |  | `read` |
| `ZENDESK_GET_DELETION_SCHEDULES` | List Deletion Schedules | read |  |  | `read` |
| `ZENDESK_GET_DYNAMIC_CONTENT_ITEMS` | List Dynamic Content Items | read | *page, per_page, sort* |  | `read` |
| `ZENDESK_GET_DYNAMIC_CONTENT_ITEMS_SHOW_MANY` | Show Many Dynamic Content Items | read | *identifiers, page, per_page* |  | `dynamic_content:read`, `read` |
| `ZENDESK_GET_DYNAMIC_CONTENT_ITEMS_VARIANTS` | List Dynamic Content Variants | read | `dynamic_content_item_id`; *page, per_page, sort* |  | `dynamic_content:read`, `read` |
| `ZENDESK_GET_DYNAMIC_CONTENT_ITEM_SHOW` | Show Dynamic Content Item | read | `dynamic_content_item_id` |  | `read` |
| `ZENDESK_GET_DYNAMIC_CONTENT_ITEM_VARIANT` | Show Dynamic Content Item Variant | read | `dynamic_content_item_id`, `dynamic_content_variant_id` |  | `read` |
| `ZENDESK_GET_EMAIL_NOTIFICATION` | Get Email Notification | read | `notification_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_EMAIL_NOTIFICATIONS` | Get Email Notifications | read | *filter_comment_id, filter_notification_id, filter_ticket_id, per_page, sort* |  | `tickets:read`, `read` |
| `ZENDESK_GET_EMAIL_NOTIFICATIONS_SHOW_MANY` | Show Many Email Notifications | read | *comment_ids, ids, ticket_ids* |  | `tickets:read`, `read` |
| `ZENDESK_GET_END_USER_IDENTITIES` | List End User Identities | read | `user_id`; *page, page_after, page_before, page_size, per_page, type_list* |  | `users:read`, `read` |
| `ZENDESK_GET_END_USER_IDENTITY` | Get End User Identity | read | `user_id`, `user_identity_id` |  | `users:read`, `read` |
| `ZENDESK_GET_GROUP` | Get Group | read | `group_id`; *include* |  | `read` |
| `ZENDESK_GET_GROUPS_ASSIGNABLE` | List Assignable Groups | read | *page, page_after, page_before, page_size, per_page* |  | `read` |
| `ZENDESK_GET_GROUPS_COUNT` | Count Groups | read |  |  | `read` |
| `ZENDESK_GET_GROUPS_USERS_COUNT` | Count Users by Group | read | `group_id`; *permission_set, role, role_list* |  | `read` |
| `ZENDESK_GET_GROUP_MEMBERSHIP` | Get Group Membership | read | `group_membership_id` |  | `read` |
| `ZENDESK_GET_GROUP_MEMBERSHIPS` | List Group Memberships | read | *include, page, page_after, page_before, page_size, per_page, sort* |  | `read` |
| `ZENDESK_GET_GROUP_MEMBERSHIPS_BY_GROUP` | List Memberships By Group | read | `group_id`; *include, page, page_after, page_before, page_size, per_page, sort* |  | `read` |
| `ZENDESK_GET_GROUP_USERS` | List Group Users | read | `group_id`; *external_id, page, per_page, permission_set, role, role_list* |  | `read` |
| `ZENDESK_GET_GUIDE_SURVEYS` | List Guide Surveys | read | `locale`; *filter_types, page_after, page_size* |  | `read` |
| `ZENDESK_GET_GUIDE_SURVEY_RESPONSES` | List Guide Survey Responses | read | *filter_created_at_end, filter_created_at_start, filter_responder_ids, filter_subject_zrns, page_after, page_before, page_size, sort* |  | `tickets:read`, `read` |
| `ZENDESK_GET_GUIDE_THEME_JOB` | Get Guide Theme Job | read | `job_id` |  | `themes:read`, `themes:write` |
| `ZENDESK_GET_INCREMENTAL_ORGANIZATIONS` | Get Incremental Organizations | read | `start_time`; *per_page* |  | `organizations:read`, `read` |
| `ZENDESK_GET_INCREMENTAL_ROUTING_ATTRIBUTES` | Incremental Routing Attributes Export | read | *cursor, per_page, start_time* |  | `read` |
| `ZENDESK_GET_INCREMENTAL_ROUTING_ATTRIBUTE_VALUES` | Get Incremental Routing Attribute Values | read | `start_time`; *cursor, per_page* |  | `read` |
| `ZENDESK_GET_INCREMENTAL_ROUTING_INSTANCE_VALUES` | Export Incremental Routing Instance Values | read | *cursor* |  | `read` |
| `ZENDESK_GET_INCREMENTAL_SAMPLE` | Get Incremental Sample Export | read | `incremental_resource`, `start_time`; *per_page* |  | `tickets:read`, `read` |
| `ZENDESK_GET_INCREMENTAL_TICKETS` | Get Incremental Tickets | read | `start_time`; *per_page, support_type_scope* |  | `tickets:read`, `read` |
| `ZENDESK_GET_INCREMENTAL_TICKETS_CURSOR` | Get Incremental Tickets (Cursor) | read | *cursor, start_time, support_type_scope* |  | `tickets:read`, `read` |
| `ZENDESK_GET_INCREMENTAL_TICKET_EVENTS` | Get Incremental Ticket Events | read | `start_time`; *include, support_type_scope* |  | `tickets:read`, `read` |
| `ZENDESK_GET_INCREMENTAL_TICKET_METRIC_EVENTS` | List Ticket Metric Events | read | `start_time`; *exclude_deleted, include_changes, page, per_page* |  | `tickets:read`, `read` |
| `ZENDESK_GET_INCREMENTAL_USERS` | Get Incremental Users | read | `start_time`; *per_page* |  | `users:read`, `read` |
| `ZENDESK_GET_INCREMENTAL_USERS_CURSOR` | Incremental User Export (Cursor) | read | *cursor, per_page, start_time* |  | `users:read`, `read` |
| `ZENDESK_GET_INCREMENTAL_VOICE_CALLS` | Get Incremental Voice Calls | read | `start_time` |  | `read` |
| `ZENDESK_GET_JOB_STATUSES` | List Job Statuses | read | *page_after, page_before, page_size* |  | `tickets:read`, `read` |
| `ZENDESK_GET_JOB_STATUSES2` | Get Job Status | read | `job_status_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_LOCALE` | Get Locale | read | `locale_id` |  | `read` |
| `ZENDESK_GET_LOCALES` | List Locales | read | *page, per_page* |  | `read` |
| `ZENDESK_GET_LOCALES_AGENT` | List Agent Locales | read |  |  | `read` |
| `ZENDESK_GET_LOCALES_CURRENT` | Get Current Locale | read |  |  | `read` |
| `ZENDESK_GET_LOCALIZED_ARTICLE` | Get Localized Article | read | `article_id`, `locale`; *include* |  | `read` |
| `ZENDESK_GET_MACRO2` | Show Macro | read | `macro_id`; *include* |  | `macros:read`, `read` |
| `ZENDESK_GET_MACROS` | List Macros | read | *access, active, category, group_id, include, only_viewable, page, per_page and 2 more* |  | `macros:read`, `read` |
| `ZENDESK_GET_MACROS_ACTIONS` | List Macro Actions | read |  |  | `macros:read`, `read` |
| `ZENDESK_GET_MACROS_ACTIVE` | List Active Macros | read | *access, category, group_id, include, page, page_after, page_before, page_size and 3 more* |  | `macros:read`, `read` |
| `ZENDESK_GET_MACROS_APPLY` | Show Changes to Ticket | read | `macro_id` |  | `macros:read`, `read` |
| `ZENDESK_GET_MACROS_SEARCH` | Search Macros | read | `query`; *access, active, category, group_id, include, only_viewable, sort_by, sort_order* |  | `macros:read`, `read` |
| `ZENDESK_GET_MACRO_ATTACHMENTS` | List Macro Attachments | read | `macro_id` |  | `macros:read`, `read` |
| `ZENDESK_GET_MACRO_CATEGORIES` | List Macro Categories | read |  |  | `macros:read`, `read` |
| `ZENDESK_GET_MACRO_DEFINITIONS` | List Macro Action Definitions | read |  |  | `macros:read`, `read` |
| `ZENDESK_GET_MACRO_REPLICA` | Show Macro Replica | read | `macro_id`, `ticket_id` |  | `macros:read`, `read` |
| `ZENDESK_GET_OBJECT_TRIGGER` | Show Object Trigger | read | `custom_object_key`, `trigger_id` |  | `read` |
| `ZENDESK_GET_OPEN_REQUESTS` | List Open Requests | read | *page, per_page, sort_by, sort_order* |  | `requests:read`, `read` |
| `ZENDESK_GET_ORGANIZATIONS_AUTOCOMPLETE` | Autocomplete Organizations | read | `name`; *field_id, source* |  | `organizations:read`, `read` |
| `ZENDESK_GET_ORGANIZATIONS_MERGES` | List Organization Merges | read | `organization_id`; *page, per_page* |  | `organizations:read`, `read` |
| `ZENDESK_GET_ORGANIZATIONS_SHOW_MANY` | Show Many Organizations | read | *external_ids, ids* |  | `organizations:read`, `read` |
| `ZENDESK_GET_ORGANIZATIONS_TICKETS_COUNT` | Count Organization Tickets | read | `organization_id` |  | `tickets:read`, `organizations:read`, `read` |
| `ZENDESK_GET_ORGANIZATION_FIELD2` | Get Organization Field | read | `organization_field_id`; *resolve_dc* |  | `organizations:read`, `read` |
| `ZENDESK_GET_ORGANIZATION_FIELDS` | List Organization Fields | read | *page, resolve_dc* |  | `organizations:read`, `read` |
| `ZENDESK_GET_ORGANIZATION_MEMBERSHIP2` | Get Organization Membership | read | `organization_membership_id` |  | `read` |
| `ZENDESK_GET_ORGANIZATION_RELATED` | Get Organization Related Info | read | `organization_id` |  | `organizations:read`, `read` |
| `ZENDESK_GET_ORGANIZATION_REQUESTS` | List Organization Requests | read | `organization_id`; *page, page_after, page_before, page_size, per_page, sort_by, sort_order* |  | `organizations:read`, `read` |
| `ZENDESK_GET_ORGANIZATION_SUBSCRIPTION` | Get Organization Subscription | read | `organization_subscription_id` |  | `organizations:read`, `read` |
| `ZENDESK_GET_ORGANIZATION_SUBSCRIPTIONS` | List Organization Subscriptions | read | *page, page_after, page_before, page_size, per_page* |  | `read` |
| `ZENDESK_GET_ORGANIZATION_TAGS` | Get Organization Tags | read | `organization_id` |  | `organizations:read`, `read` |
| `ZENDESK_GET_ORGANIZATION_TICKETS` | Get Organization Tickets | read | `organization_id`; *page, per_page, sort_by, sort_order* |  | `tickets:read`, `read` |
| `ZENDESK_GET_ORGANIZATION_USERS` | List Organization Users | read | `organization_id`; *external_id, page, per_page, permission_set, role, role_list, sort_by, sort_order* |  | `users:read`, `organizations:read`, `read` |
| `ZENDESK_GET_ORGANIZATION_USERS_COUNT` | Count Organization Users | read | `organization_id`; *permission_set, role* |  | `organizations:read`, `users:read`, `read` |
| `ZENDESK_GET_OWNED_APPS` | Get Owned Apps | read | *exclude, page, per_page* |  | `apps:read`, `read` |
| `ZENDESK_GET_PROBLEMS` | List Ticket Problems | read | *include_boundary_indicators, include_item_cursors, page, page_after, page_before, page_size, per_page, sort_by and 1 more* |  | `tickets:read`, `read` |
| `ZENDESK_GET_PUBLIC_LOCALES` | List Public Locales | read |  |  | `read` |
| `ZENDESK_GET_RECENT_TICKETS` | Get Recent Tickets | read | *page, per_page* |  | `tickets:read`, `read` |
| `ZENDESK_GET_REMOTE_AUTHENTICATIONS` | List Remote Authentications | read |  |  | `read` |
| `ZENDESK_GET_REQUESTS_CCDB` | List CCD Requests | read | *page, per_page, sort_by, sort_order* |  | `requests:read`, `read` |
| `ZENDESK_GET_REQUESTS_COMMENT` | Get Request Comment | read | `request_id`, `ticket_comment_id` |  | `requests:read`, `read` |
| `ZENDESK_GET_REQUESTS_SEARCH` | Search Requests | read | *page, per_page, query, sort_by, sort_order* |  | `requests:read`, `read` |
| `ZENDESK_GET_REQUESTS_SOLVED` | List Solved Requests | read | *page, per_page, sort_by, sort_order* |  | `requests:read`, `read` |
| `ZENDESK_GET_REQUEST_COMMENTS` | List Request Comments | read | `request_id`; *page, per_page, role, since* |  | `requests:read`, `read` |
| `ZENDESK_GET_ROUTING_AGENTS_INSTANCE_VALUES` | Get Agent Attribute Values | read | `agent_ids`; *page_size* |  | `read` |
| `ZENDESK_GET_SATISFACTION_RATINGS` | List Satisfaction Ratings | read | *page, per_page, sort* |  | `tickets:read`, `read` |
| `ZENDESK_GET_SATISFACTION_RATINGS_COUNT` | Count Satisfaction Ratings | read |  |  | `tickets:read`, `read` |
| `ZENDESK_GET_SEARCH_COUNT` | Get Search Results Count | read | `query` |  | `read` |
| `ZENDESK_GET_SEARCH_EXPORT` | Export Search Results | read | `query`; *filter_type, include, page_after, page_size* |  | `read` |
| `ZENDESK_GET_SECURITY_SETTINGS` | Get Security Settings | read |  |  | `read` |
| `ZENDESK_GET_SIDE_CONVERSATIONS_EVENTS` | List Side Conversation Events | read | `start_time`; *page, per_page* |  | `tickets:read`, `read` |
| `ZENDESK_GET_SUPPORT_ADDRESS` | Get Support Address | read | `support_address_id` |  | `read` |
| `ZENDESK_GET_TAGS` | List Tags | read | *page, per_page, sort* |  | `read` |
| `ZENDESK_GET_TAGS_COUNT` | Count Tags | read |  |  | `read` |
| `ZENDESK_GET_TARGET` | Show Target | read | `target_id` |  | `read` |
| `ZENDESK_GET_TARGETS` | List Targets | read | *page, per_page* |  | `targets:read`, `read` |
| `ZENDESK_GET_TARGET_FAILURE` | Show Target Failure | write | `target_failure_id` |  | `read` |
| `ZENDESK_GET_TARGET_FAILURES` | List Target Failures | read |  |  | `read` |
| `ZENDESK_GET_TICKETS_AUDITS` | List Ticket Audits | read | `ticket_id`; *filter_events, include, include_boundary_indicators, include_item_cursors, page, page_after, page_before, page_size and 2 more* |  | `tickets:read`, `read` |
| `ZENDESK_GET_TICKETS_AUDITS_COUNT` | Count Ticket Audits | read | `ticket_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_TICKETS_COMMENTS_COUNT` | Count Ticket Comments | read | `ticket_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_TICKETS_CONVERSATION_LOG` | Get Ticket Conversation Log | read | `ticket_id`; *page_after, page_before, page_size, sort* |  | `read` |
| `ZENDESK_GET_TICKETS_COUNT` | Count Zendesk Tickets | read |  |  | `tickets:read`, `read` |
| `ZENDESK_GET_TICKETS_METRICS` | Get Ticket Metrics | read | `ticket_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_TICKETS_SIDE_CONVERSATIONS_EVENTS` | List Side Conversation Events | read | `side_conversation_id`, `ticket_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_TICKET_AUDIT_BY_ID` | Get Ticket Audit | read | `ticket_audit_id`, `ticket_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_TICKET_COLLABORATORS` | List Ticket Collaborators | read | `ticket_id`; *page, per_page* |  | `tickets:read`, `read` |
| `ZENDESK_GET_TICKET_COMMENTS` | Get Ticket Comments | read | `ticket_id`; *include, include_inline_images, per_page, sort_order* |  | `read` |
| `ZENDESK_GET_TICKET_EMAIL_C_CS` | List Ticket Email CCs | read | `ticket_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_TICKET_FIELD2` | Show Ticket Field | read | `ticket_field_id`; *creator* |  | `read` |
| `ZENDESK_GET_TICKET_FIELDS` | List Ticket Fields | read | *creator, include_boundary_indicators, include_item_cursors, locale, page_after, page_before, page_size, sort* |  | `tickets:read`, `read` |
| `ZENDESK_GET_TICKET_FIELDS_COUNT` | Count Ticket Fields | read |  |  | `tickets:read`, `read` |
| `ZENDESK_GET_TICKET_FIELDS_SHOW_MANY` | Show Many Ticket Fields | read | *creator, exclude_sub_selection_options, ids, keys* |  | `read` |
| `ZENDESK_GET_TICKET_FOLLOWERS` | List Ticket Followers | read | `ticket_id`; *page, per_page* |  | `tickets:read`, `read` |
| `ZENDESK_GET_TICKET_FORMS` | List Ticket Forms | read | *active, associated_to_brand, end_user_visible, fallback_to_default, form_type, include_boundary_indicators, include_item_cursors, locale and 3 more* |  | `tickets:read`, `read` |
| `ZENDESK_GET_TICKET_FORMS_SHOW_MANY` | Show Many Ticket Forms | read | `ids`; *active, associated_to_brand, end_user_visible, fallback_to_default, include_boundary_indicators, include_item_cursors* |  | `read` |
| `ZENDESK_GET_TICKET_INCIDENTS` | List Ticket Incidents | read | `ticket_id`; *page, page_after, page_before, page_size, per_page* |  | `tickets:read`, `read` |
| `ZENDESK_GET_TICKET_METRIC` | Get Ticket Metric | read | `ticket_metric_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_TICKET_METRICS` | Get Ticket Metrics | read | *page_after, page_before, page_size, sort_by* |  | `read` |
| `ZENDESK_GET_TICKET_SKIPS` | List Ticket Skips | read | `ticket_id`; *sort_order* |  | `read` |
| `ZENDESK_GET_TICKET_SLA_POLICY_METRICS` | Get Ticket SLA Policy Metrics | read | `ticket_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_TRIGGER` | Show Trigger | read | `trigger_id` |  | `triggers:read`, `read` |
| `ZENDESK_GET_TRIGGERS` | List Ticket Triggers | read | *active, category_id, include, page, per_page, sort, sort_by, sort_order* |  | `triggers:read`, `read` |
| `ZENDESK_GET_TRIGGERS_DEFINITIONS` | Get Ticket Trigger Definitions | read |  |  | `triggers:read`, `read` |
| `ZENDESK_GET_TRIGGERS_SEARCH` | Search Ticket Triggers | read | *active, filter, include, query, sort, sort_by, sort_order* |  | `triggers:read`, `read` |
| `ZENDESK_GET_TRIGGER_CATEGORIES` | List Trigger Categories | read | *include, page, sort* |  | `read` |
| `ZENDESK_GET_TRIGGER_CATEGORIES2` | Show Trigger Category | read | `trigger_category_id` |  | `read` |
| `ZENDESK_GET_USER` | Get User | read | `user_id` |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_ASSIGNED_TICKETS` | Get User's Assigned Tickets | read | `user_id`; *page, per_page, sort_by, sort_order* |  | `tickets:read`, `read` |
| `ZENDESK_GET_USERS_AUTOCOMPLETE` | Autocomplete Users | read | *field_id, filter, include, name, per_page, phone, source* |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_BRAND_AGENTS` | List Brand Agent Memberships By User | read | `user_id`; *page, per_page, sort* |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_CCD_TICKETS` | Get User's CC'd Tickets | read | `user_id`; *page, per_page, sort_by, sort_order* |  | `tickets:read`, `read` |
| `ZENDESK_GET_USERS_COMPLIANCE_DELETION_STATUSES` | Get User Compliance Deletion Statuses | read | `user_id`; *application, page, per_page* |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_COUNT` | Count Users | read | *permission_set, role, role_list* |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_FOLLOWED_TICKETS` | Get User's Followed Tickets | read | `user_id`; *exclude_archived, page, per_page, sort_by, sort_order* |  | `tickets:read`, `read` |
| `ZENDESK_GET_USERS_GROUPS` | List User Groups | read | `user_id`; *exclude_deleted, include_boundary_indicators, include_item_cursors, page, per_page, sort* |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_GROUPS_COUNT` | Count User Groups | read | `user_id` |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_GROUP_MEMBERSHIPS` | List Memberships By User | read | `user_id`; *include, page, page_after, page_before, page_size, per_page, sort* |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_GROUP_MEMBERSHIPS_BY_ID` | Get User's Group Membership | read | `group_membership_id`, `user_id` |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_IDENTITIES2` | Show User Identity | read | `user_id`, `user_identity_id` |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_ME_SESSION` | Get Current Session | read | *page, per_page* |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_ME_SETTINGS` | Get Current User Settings | read |  |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_ORGANIZATIONS` | List User Organizations | read | `user_id`; *page, per_page, sort* |  | `organizations:read`, `users:read`, `read` |
| `ZENDESK_GET_USERS_ORGANIZATIONS_COUNT` | Get User's Organizations Count | read | `user_id` |  | `users:read`, `organizations:read`, `read` |
| `ZENDESK_GET_USERS_ORGANIZATION_MEMBERSHIPS` | List User Organization Memberships | read | `user_id`; *include* |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_ORGANIZATION_MEMBERSHIPS2` | Get User's Organization Membership | read | `organization_membership_id`, `user_id` |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_ORGANIZATION_SUBSCRIPTIONS` | List User's Organization Subscriptions | read | `user_id`; *page, page_after, page_before, page_size, per_page* |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_PASSWORD_REQUIREMENTS` | Get User Password Requirements | read | `user_id` |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_PROFILES` | Get Users Profiles | read | `user_id` |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_RELATED` | Get User Related Info | read | `user_id` |  | `users:read`, `read` |
| `ZENDESK_GET_USERS_REQUESTED_TICKETS` | Get User's Requested Tickets | read | `user_id`; *page, per_page, sort_by, sort_order* |  | `tickets:read`, `read` |
| `ZENDESK_GET_USERS_REQUESTS` | List User Requests | read | `user_id`; *sort_by, sort_order* |  | `requests:read`, `users:read`, `read` |
| `ZENDESK_GET_USERS_SKIPS` | List User Skips | read | `user_id`; *page, per_page, sort_order* |  | `tickets:read`, `read` |
| `ZENDESK_GET_USERS_TAGS` | Get User Tags | read | `user_id` |  | `users:read`, `read` |
| `ZENDESK_GET_USER_ENTITLEMENTS_FULL` | Get Full User Entitlements | read | `user_id` |  | `users:read`, `read` |
| `ZENDESK_GET_USER_EVENTS` | Get User Events | read | `user_id`; *filter_end_time, filter_source, filter_start_time, filter_type, page_size* |  | `read`, `users:read` |
| `ZENDESK_GET_USER_FIELDS` | List User Fields | read | *page_after, page_before, page_size, resolve_dc* |  | `users:read`, `read` |
| `ZENDESK_GET_USER_FIELDS_SHOW_MANY` | Get User Fields Show Many | read | *keys* |  | `users:read`, `read` |
| `ZENDESK_GET_USER_FIELD_BY_ID` | Get User Field By ID | read | `user_field_id` |  | `users:read`, `read` |
| `ZENDESK_GET_USER_FIELD_OPTIONS` | List User Field Options | read | `user_field_id`; *page, per_page* |  | `users:read`, `read` |
| `ZENDESK_GET_USER_PROFILES2` | Get User Profile by Profile ID | read | `profile_id` |  | `users:read`, `read` |
| `ZENDESK_GET_USER_PROFILES_EVENTS2` | Get User Profile Events by ID | read | `profile_id`; *filter_end_time, filter_source, filter_start_time, filter_type, page_size* |  | `users:read`, `read` |
| `ZENDESK_GET_USER_SESSIONS` | List User Sessions | read | `user_id`; *page_after, page_before, page_size, sort* |  | `users:read`, `read` |
| `ZENDESK_GET_USER_TICKETS_ASSIGNED_COUNT` | Count User Assigned Tickets | read | `user_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_VIEWS` | List Views | read | *access, active, group_id, page, per_page, sort, sort_by, sort_order* |  | `read` |
| `ZENDESK_GET_VIEWS_ACTIVE` | List Active Views | read | *access, group_id, page, per_page, sort_by, sort_order* |  | `read` |
| `ZENDESK_GET_VIEWS_COMPACT` | List Views - Compact | read |  |  | `read` |
| `ZENDESK_GET_VIEWS_COUNT` | Count Views | read |  |  | `read` |
| `ZENDESK_GET_VIEWS_COUNT_MANY` | Count Tickets in Views | read | `ids` |  | `read` |
| `ZENDESK_GET_VIEWS_EXECUTE` | Execute View | write | `view_id`; *exclude, group_by, include, page, sort_by, sort_order* |  | `read` |
| `ZENDESK_GET_VIEWS_EXPORT` | Export View | read | `view_id` |  | `read` |
| `ZENDESK_GET_VIEWS_SEARCH` | Search Views | read | `query`; *access, active, group_id, include, page, per_page, sort_by, sort_order* |  | `read` |
| `ZENDESK_GET_VIEWS_SHOW_MANY` | List Views By IDs | read | `ids`; *active* |  | `read` |
| `ZENDESK_GET_VIEWS_TICKETS` | List Tickets From View | read | `view_id`; *page_after, page_before, page_size, sort_by, sort_order* |  | `read` |
| `ZENDESK_GET_VIEW_DEFINITIONS` | List View Filter Definitions | read |  |  | `read` |
| `ZENDESK_GET_VIEW_TICKET_COUNT` | Count Tickets in View | read | `view_id` |  | `read` |
| `ZENDESK_GET_VOICE_AVAILABILITY` | Get Voice Availability | read | `agent_id` |  | `read` |
| `ZENDESK_GET_ZENDESK_ACTIVITIES` | List Activities | read | *include, page, per_page, since, sort* |  | `tickets:read`, `read` |
| `ZENDESK_GET_ZENDESK_CUSTOM_ROLES` | List Custom Roles | write |  |  | `read` |
| `ZENDESK_GET_ZENDESK_ORGANIZATION` | Get Zendesk Organization | read | `organization_id` |  | `organizations:read`, `read` |
| `ZENDESK_GET_ZENDESK_ORGANIZATIONS_SUBSCRIPTIONS` | List Organization Subscriptions | read | `organization_id`; *page, per_page* |  | `organizations:read`, `read` |
| `ZENDESK_GET_ZENDESK_ORGANIZATION_MERGE` | Get Zendesk Organization Merge | read | `organization_merge_id` |  | `organizations:read`, `read` |
| `ZENDESK_GET_ZENDESK_REQUESTS` | List Requests | read | *page, per_page, sort_by, sort_order* |  | `requests:read`, `read` |
| `ZENDESK_GET_ZENDESK_SESSIONS` | List Sessions | read | *page, sort* |  | `read` |
| `ZENDESK_GET_ZENDESK_SHARING_AGREEMENTS` | List Sharing Agreements | read |  |  | `read` |
| `ZENDESK_GET_ZENDESK_SUSPENDED_TICKETS` | List Suspended Tickets | read | *page, per_page, sort_by, sort_order* |  | `tickets:read`, `read` |
| `ZENDESK_GET_ZENDESK_TICKETS_SIDE_CONVERSATION_ID` | Get Side Conversation | read | `side_conversation_id`, `ticket_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_ZENDESK_TICKETS_TAGS` | Get Ticket Tags | read | `ticket_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_ZENDESK_TICKET_BY_ID` | Get Zendesk Ticket | read | `ticket_id`; *page, per_page* |  | `tickets:read`, `read` |
| `ZENDESK_GET_ZENDESK_TICKET_FIELDS_OPTIONS` | List Ticket Field Options | read | `ticket_field_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_ZENDESK_TICKET_FIELDS_OPTIONS2` | Show Ticket Field Option | read | `ticket_field_id`, `ticket_field_option_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_ZENDESK_TICKET_FORM` | Get Zendesk Ticket Form | read | `ticket_form_id` |  | `tickets:read`, `read` |
| `ZENDESK_GET_ZENDESK_TICKET_RELATED` | Get Zendesk Ticket Related | read | `ticket_id` |  | `tickets:read`, `read` |
| `ZENDESK_LIST_ACTIVE_AUTOMATIONS` | List Active Automations | read | *page, page_after, page_before, page_size, per_page* |  | `automations:read`, `read` |
| `ZENDESK_LIST_APP_INSTALLATIONS` | List App Installations | read | *include, page, per_page* |  | `apps:read`, `read` |
| `ZENDESK_LIST_APP_LOCATIONS` | List App Locations | read |  |  | `apps:read`, `read` |
| `ZENDESK_LIST_ARTICLE_ATTACHMENTS` | List Article Attachments | read | `article_id`; *page, page_after, page_before, page_size, per_page* |  | `read` |
| `ZENDESK_LIST_ARTICLE_TRANSLATIONS` | List Article Translations | read | `article_id`; *draft, include, locales, outdated, page, page_after, page_before, page_size and 2 more* |  | `read` |
| `ZENDESK_LIST_ASSIGNABLE_GROUP_MEMBERSHIPS` | List Assignable Group Memberships | read | *page, page_after, page_size, per_page* |  | `read` |
| `ZENDESK_LIST_ASSIGNABLE_GROUP_MEMBERSHIPS_BY_GROUP` | List Assignable Memberships By Group | read | `group_id`; *page, page_after, page_size, per_page* |  | `read` |
| `ZENDESK_LIST_BOOKMARKS` | List Bookmarks | read | *page, per_page* |  | `read` |
| `ZENDESK_LIST_BRAND_AGENTS` | List Brand Agent Memberships | read | *page, per_page, sort* |  | `read` |
| `ZENDESK_LIST_CUSTOM_OBJECT_RECORDS` | List Custom Object Records | write | `custom_object_key`; *filter_external_ids, filter_ids, page_after, page_before, page_size, sort* |  |  |
| `ZENDESK_LIST_CUSTOM_OBJECT_RECORD_ATTACHMENTS` | List Custom Object Record Attachments | write | `custom_object_key`, `record_id` |  | `read` |
| `ZENDESK_LIST_DELETED_USERS` | List Deleted Users | read | *page, per_page, sort* |  | `users:read`, `read` |
| `ZENDESK_LIST_GROUPS` | List Groups | read | *exclude_deleted, include, include_boundary_indicators, include_item_cursors, page, page_after, page_before, page_size and 2 more* |  | `read` |
| `ZENDESK_LIST_HELP_CENTER_CONTENT` | List Help Center Content | read | `content_type`, `locale`; *page, per_page, sort_by, sort_order* |  | `read` |
| `ZENDESK_LIST_INCREMENTAL_CUSTOM_OBJECT_RECORDS` | List Incremental Custom Object Records | write | `custom_object_key`; *cursor, filter_exclude_deleted, per_page, start_time* |  | `read` |
| `ZENDESK_LIST_MONITORED_TWITTER_HANDLES` | List Monitored X Handles | read | *page, per_page* |  | `read` |
| `ZENDESK_LIST_OBJECT_TRIGGERS` | List Object Triggers | read | `custom_object_key`; *active, page, per_page, sort_by, sort_order* |  | `read` |
| `ZENDESK_LIST_OBJECT_TRIGGERS_DEFINITIONS` | List Object Trigger Definitions | read | `custom_object_key` |  | `read` |
| `ZENDESK_LIST_ORGANIZATION_MEMBERSHIPS` | List Organization Memberships | read | *include, page, per_page, sort* |  | `organizations:read`, `read` |
| `ZENDESK_LIST_ORGANIZATION_MEMBERSHIPS_BY_ORGANIZATION` | List Organization Memberships by Organization | read | `organization_id`; *include, page, per_page, sort* |  | `organizations:read`, `read` |
| `ZENDESK_LIST_RESOURCE_COLLECTIONS` | List Resource Collections | read | *per_page* |  | `read` |
| `ZENDESK_LIST_ROUTING_ATTRIBUTES` | List Routing Attributes | read | *include* |  | `read` |
| `ZENDESK_LIST_ROUTING_ATTRIBUTE_VALUES` | List Routing Attribute Values | read | `attribute_id` |  | `read` |
| `ZENDESK_LIST_SKIPS` | List All Skips | read | *page, per_page, sort_order* |  | `tickets:read`, `read` |
| `ZENDESK_LIST_SUPPORT_ADDRESSES` | List Support Addresses | read | *include, page, per_page, sort* |  | `read` |
| `ZENDESK_LIST_TICKET_AUDITS` | List Ticket Audits | read | *page_after, page_before, page_size* |  | `tickets:read`, `read` |
| `ZENDESK_LIST_TICKET_CONTENT_PINS` | List Ticket Content Pins | read | `ticket_id` |  | `tickets:read`, `read` |
| `ZENDESK_LIST_USER_IDENTITIES` | List User Identities | read | `user_id`; *include_boundary_indicators, include_item_cursors, page, sort, type_list* |  | `users:read`, `read` |
| `ZENDESK_LIST_WEBHOOK_INVOCATIONS` | List Webhook Invocations | read | `webhook_id`; *from_ts, page_after, page_before, page_size, sort, status, to_ts* |  | `webhooks:read` |
| `ZENDESK_LIST_ZENDESK_ACCESS_RULES` | List Access Rules | read | `custom_object_key` |  | `read` |
| `ZENDESK_LIST_ZENDESK_CUSTOM_OBJECTS` | List Custom Objects | write | *include_ui_path* |  | `read` |
| `ZENDESK_LIST_ZENDESK_CUSTOM_OBJECT_FIELDS` | List Custom Object Fields | write | `custom_object_key`; *include_standard_fields, page, per_page* |  | `read` |
| `ZENDESK_LIST_ZENDESK_PERMISSION_POLICIES` | List Permission Policies | read | `custom_object_key` |  | `read` |
| `ZENDESK_LIST_ZENDESK_TICKETS` | List Tickets | read | *external_id, page, page_after, page_before, page_size, per_page, sort_by, sort_order* |  | `tickets:read`, `read` |
| `ZENDESK_LIST_ZENDESK_USERS` | List Users | read | *external_id, page, page_after, page_size, per_page, permission_set, role, role_list* |  | `users:read`, `read` |
| `ZENDESK_MAKE_COMMENT_PRIVATE` | Make Zendesk Ticket Comment Private | write | `ticket_comment_id`, `ticket_id` |  | `tickets:write`, `write` |
| `ZENDESK_MAKE_DEFAULT_ORGANIZATION_MEMBERSHIP` | Make Organization Membership Default | write | `organization_membership_id`, `user_id` |  | `users:write`, `write` |
| `ZENDESK_MAKE_END_USER_IDENTITY_PRIMARY` | Make End User Identity Primary | write | `user_id`, `user_identity_id` |  | `write` |
| `ZENDESK_MAKE_TICKET_COMMENT_PRIVATE` | Make Ticket Comment Private | write | `ticket_audit_id`, `ticket_id` |  | `tickets:write`, `write` |
| `ZENDESK_MAKE_USER_IDENTITY_PRIMARY` | Make User Identity Primary | write | `user_id`, `user_identity_id` |  | `users:write`, `write` |
| `ZENDESK_MARK_MANY_TICKETS_AS_SPAM` | Mark Multiple Tickets as Spam | write | `ids` |  | `tickets:write`, `write` |
| `ZENDESK_MERGE_USERS` | Merge Users | write | `user`, `user_id` |  | `users:write`, `write` |
| `ZENDESK_MERGE_ZENDESK_TICKETS` | Merge Zendesk Tickets | write | `source_ticket_ids`, `target_ticket_id`; *source_comment, source_comment_is_public, target_comment, target_comment_is_public* |  | `tickets:write`, `write` |
| `ZENDESK_OPEN_USERS_PROFILE_IN_AGENT_BROWSER` | Open User's Profile in Agent Browser | write | `agent_id`, `user_id` |  | `write` |
| `ZENDESK_RECOVER_MANY_SUSPENDED_TICKETS` | Recover Multiple Suspended Tickets | write | `ids` |  | `tickets:write`, `write` |
| `ZENDESK_RENEW_SESSION` | Renew Session | write |  |  | `users:read`, `read` |
| `ZENDESK_REORDER_CUSTOM_OBJECT_FIELDS` | Reorder Custom Object Fields | write | `custom_object_field_ids`, `custom_object_key` |  | `write` |
| `ZENDESK_REORDER_LOCATION_INSTALLATIONS` | Reorder App Location Installations | write | `installations`, `location_name` |  | `apps:write`, `write` |
| `ZENDESK_REORDER_ORGANIZATION_FIELDS` | Reorder Organization Fields | write | `organization_field_ids` |  | `organizations:write`, `write` |
| `ZENDESK_REORDER_TICKET_FIELDS` | Reorder Ticket Fields | write | `ticket_field_ids` |  | `tickets:write`, `write` |
| `ZENDESK_REORDER_TRIGGERS` | Reorder Ticket Triggers | write | `trigger_ids` |  | `triggers:write`, `write` |
| `ZENDESK_REORDER_USER_FIELDS` | Reorder User Fields | write | `user_field_ids` |  | `users:write`, `write` |
| `ZENDESK_REPLY_ZENDESK_TICKET` | Reply to Zendesk Ticket | write | `ticket_id`; *assignee_id, body, public* |  | `tickets:write`, `write` |
| `ZENDESK_REQUEST_END_USER_VERIFICATION` | Request End User Verification | write | `user_id`, `user_identity_id` |  | `write` |
| `ZENDESK_REQUEST_USER_CREATE` | Request User Create | write | `user` |  | `users:write`, `write` |
| `ZENDESK_REQUEST_USER_VERIFICATION` | Request User Verification | write | `user_id`, `user_identity_id` |  | `users:write`, `write` |
| `ZENDESK_RESTORE_MANY_TICKETS` | Restore Multiple Deleted Tickets | write | `ids` |  | `tickets:write`, `write` |
| `ZENDESK_RESTORE_ZENDESK_TICKET` | Restore Zendesk Ticket | write | `ticket_id` |  | `tickets:write`, `write` |
| `ZENDESK_SEARCH_CUSTOM_OBJECT_RECORDS` | Search Custom Object Records | write | `custom_object_key`; *filter, page_after, page_before, page_size, query, sort* |  | `read` |
| `ZENDESK_SEARCH_HELP_CENTER_ARTICLES` | Search Help Center Articles | read | `query`; *brand_id, category, label_names, locale, multibrand, page, per_page, section and 2 more* |  | `read` |
| `ZENDESK_SEARCH_OBJECT_TRIGGERS` | Search Object Triggers | read | `custom_object_key`; *active, filter, include, query, sort, sort_by, sort_order* |  | `read` |
| `ZENDESK_SEARCH_ZENDESK` | Search Zendesk | read | `query`; *include, page, per_page, sort_by, sort_order* |  | `read` |
| `ZENDESK_SEARCH_ZENDESK_ORGANIZATIONS` | Search Organizations | read | *external_id, name* |  | `organizations:read`, `read` |
| `ZENDESK_SEARCH_ZENDESK_USERS` | Search Zendesk Users | read | *external_id, include, page, per_page, query* |  | `users:read`, `read` |
| `ZENDESK_SET_DEFAULT_CUSTOM_STATUS` | Set Default Custom Ticket Status | write | `ids` |  | `tickets:write`, `write` |
| `ZENDESK_SET_GROUP_MEMBERSHIP_AS_DEFAULT` | Set Group Membership as Default | write | `group_membership_id`, `user_id` |  | `users:write`, `write` |
| `ZENDESK_SET_USERS_ORGANIZATION_AS_DEFAULT` | Set Organization as Default for User | write | `organization_id`, `user_id` |  | `users:write`, `organizations:write`, `write` |
| `ZENDESK_SHOW_ACCESS_RULE` | Show Access Rule | read | `custom_object_key`, `id` |  | `read` |
| `ZENDESK_SHOW_CUSTOM_OBJECT` | Show Custom Object | write | `custom_object_key`; *include_permissions_metadata, include_ui_path* |  | `read` |
| `ZENDESK_SHOW_CUSTOM_OBJECT_RECORD` | Show Custom Object Record | write | `custom_object_key`, `custom_object_record_id` |  | `read` |
| `ZENDESK_SHOW_MANY_JOB_STATUSES` | Show Many Job Statuses | read | `ids` |  | `read` |
| `ZENDESK_SHOW_MANY_TICKETS` | Show Many Tickets | read | `ids`; *include* |  | `read` |
| `ZENDESK_SHOW_MANY_USERS` | Show Many Users | read | *external_ids, ids, include, include_deleted* |  | `users:read`, `read` |
| `ZENDESK_SHOW_PERMISSION_POLICY` | Show Permission Policy | read | `custom_object_key`, `id` |  | `read` |
| `ZENDESK_SHOW_REQUEST` | Show Request | read | `request_id` |  | `requests:read`, `read` |
| `ZENDESK_TEST_WEBHOOK` | Test Webhook | write | `endpoint`, `http_method`, `request_format`; *custom_headers, payload, webhook_id* |  | `webhooks:write` |
| `ZENDESK_UNASSIGN_USER_ORGANIZATION` | Unassign User from Organization | write | `organization_id`, `user_id` |  | `organizations:write`, `users:write`, `write` |
| `ZENDESK_UPDATE_ACCESS_RULE` | Update Access Rule | write | `custom_object_key`, `id`; *conditions, description, title* |  | `write` |
| `ZENDESK_UPDATE_APPS_INSTALLATION` | Update Zendesk App Installation | write | `app_installation_id`; *enabled, group_restrictions, role_restrictions, settings, settings_objects* |  | `apps:write`, `write` |
| `ZENDESK_UPDATE_ARTICLE` | Update Article | write | `article_id`; *author_id, comments_disabled, content_tag_ids, label_names, permission_group_id, position, promoted, section_id and 2 more* |  | `write` |
| `ZENDESK_UPDATE_ARTICLE_TRANSLATION` | Update Article Translation | write | `article_id`, `locale`; *body, draft, new_locale, outdated, title* |  | `write` |
| `ZENDESK_UPDATE_AUTOMATION` | Update Automation | write | `automation_id`; *actions, active, conditions, position, title* |  | `automations:write`, `write` |
| `ZENDESK_UPDATE_CUSTOM_OBJECT` | Update Custom Object | write | `custom_object_key`; *allows_attachments, allows_photos, description, include_in_list_view, title, title_pluralized* |  | `write` |
| `ZENDESK_UPDATE_CUSTOM_OBJECT_RECORD` | Update Custom Object Record | write | `custom_object_key`, `custom_object_record_id`; *custom_object_fields, external_id, name* |  | `write` |
| `ZENDESK_UPDATE_CUSTOM_STATUS` | Update Custom Ticket Status | write | `custom_status_id`; *active, agent_label, description, end_user_description, end_user_label* |  | `tickets:write`, `write` |
| `ZENDESK_UPDATE_DELETION_SCHEDULE` | Update Deletion Schedule | write | `deletion_schedule_id`; *active, conditions, description, title* |  | `write` |
| `ZENDESK_UPDATE_DYNAMIC_CONTENT_ITEMS_VARIANT` | Update Dynamic Content Variant | write | `content`, `dynamic_content_item_id`, `dynamic_content_variant_id`, `locale_id`; *active, default* |  | `dynamic_content:write`, `write` |
| `ZENDESK_UPDATE_DYNAMIC_CONTENT_ITEMS_VARIANTS` | Update Many Dynamic Content Variants | write | `dynamic_content_item_id`, `variants` |  | `dynamic_content:write`, `write` |
| `ZENDESK_UPDATE_MACRO` | Update Macro | write | `actions`, `macro_id`, `title`; *active, description, restriction* |  | `macros:write`, `write` |
| `ZENDESK_UPDATE_MANY_AUTOMATIONS` | Update Many Automations | write | `automations` |  | `automations:write`, `write` |
| `ZENDESK_UPDATE_MANY_MACROS` | Update Many Macros | write | `macros` |  | `macros:write`, `write` |
| `ZENDESK_UPDATE_MANY_OBJECT_TRIGGERS` | Update Many Object Triggers | write | `custom_object_key`, `triggers` |  | `write` |
| `ZENDESK_UPDATE_MANY_ORGANIZATIONS` | Update Many Zendesk Organizations | write | `organizations`; *external_ids, ids* |  | `organizations:write`, `write` |
| `ZENDESK_UPDATE_MANY_TICKETS` | Update Many Tickets | write | `tickets`; *ids* |  | `tickets:write`, `write` |
| `ZENDESK_UPDATE_MANY_TRIGGERS` | Update Many Ticket Triggers | write | `triggers` |  | `triggers:write`, `write` |
| `ZENDESK_UPDATE_MANY_USERS` | Update Many Users | write | *external_ids, ids, user, users* |  | `users:write`, `write` |
| `ZENDESK_UPDATE_MANY_VIEWS` | Update Many Views | write | `views` |  | `write` |
| `ZENDESK_UPDATE_OBJECT_TRIGGER` | Update Object Trigger | write | `custom_object_key`, `trigger_id`; *actions, active, conditions, description, title* |  | `write` |
| `ZENDESK_UPDATE_ORGANIZATIONS_TAGS` | Add Organization Tags | write | `organization_id`, `tags` |  | `organizations:write`, `write` |
| `ZENDESK_UPDATE_ORGANIZATION_FIELD` | Update Organization Field | write | `key`, `organization_field_id`, `title`, `type`; *active, custom_field_options, description, position, regexp_for_validation, relationship_filter, relationship_target_type, tag* |  | `organizations:write`, `write` |
| `ZENDESK_UPDATE_PERMISSION_POLICY` | Update Permission Policy | write | `custom_object_key`, `id`; *records* |  | `write` |
| `ZENDESK_UPDATE_RESOURCE_COLLECTION` | Update Resource Collection | write | `resource_collection_id`; *name, resources* |  | `write` |
| `ZENDESK_UPDATE_SECTION` | Update Section | write | `locale`, `section_id`; *category_id, description, name, parent_section_id, position, sorting, theme_template* |  | `write` |
| `ZENDESK_UPDATE_SIDE_CONVERSATION` | Update Side Conversation | write | `side_conversation_id`, `ticket_id`; *state, subject* |  | `tickets:write`, `write` |
| `ZENDESK_UPDATE_SUPPORT_ADDRESS` | Update Support Address | write | `support_address_id`; *brand_id, default, name* |  | `write` |
| `ZENDESK_UPDATE_TARGET` | Update Zendesk Target | write | `target_id`; *active, data, title* |  | `targets:write`, `write` |
| `ZENDESK_UPDATE_TICKETS_TAGS` | Update Zendesk Ticket Tags | write | `tags`, `ticket_id` |  | `tickets:write`, `write` |
| `ZENDESK_UPDATE_TICKET_FIELD` | Update Ticket Field | write | `ticket_field_id`; *active, agent_can_edit, agent_description, collapsed_for_agents, creator, description, editable_in_portal, position and 8 more* |  | `write` |
| `ZENDESK_UPDATE_TICKET_FORM` | Update Ticket Form | write | `ticket_form_id`; *active, agent_conditions, default, display_name, end_user_conditions, end_user_visible, in_all_brands, name and 4 more* |  | `tickets:write`, `write` |
| `ZENDESK_UPDATE_TICKET_TRIGGER` | Update Ticket Trigger | write | `trigger_id`; *actions, active, category_id, conditions, description, position, title* |  | `triggers:write`, `write` |
| `ZENDESK_UPDATE_TRIGGER_CATEGORY` | Update Trigger Category | write | `trigger_category_id`; *name, position* |  | `write` |
| `ZENDESK_UPDATE_USER` | Update User | write | `user`, `user_id` |  | `users:write`, `write` |
| `ZENDESK_UPDATE_USERS_IDENTITY` | Update User Identity | write | `user_id`, `user_identity_id`; *verified* |  | `users:write`, `write` |
| `ZENDESK_UPDATE_USERS_ME_SETTINGS` | Update Current User Settings | write | `settings` |  | `users:write`, `write` |
| `ZENDESK_UPDATE_USERS_PROFILES` | Update Users Profiles | write | `identifier`, `profile`, `user_id` |  | `users:write`, `write` |
| `ZENDESK_UPDATE_USERS_TAGS` | Update User Tags | write | `tags`, `user_id` |  | `users:write`, `write` |
| `ZENDESK_UPDATE_USER_FIELD` | Update Zendesk User Field | write | `user_field_id`; *active, custom_field_options, description, position, raw_description, raw_title, regexp_for_validation, relationship_filter and 3 more* |  | `users:write`, `write` |
| `ZENDESK_UPDATE_USER_PROFILE` | Update Zendesk User Profile | write | `identifier`; *attributes, identifiers, name, source, type* |  | `users:write`, `write` |
| `ZENDESK_UPDATE_USER_PROFILE_BY_ID` | Update Zendesk User Profile By ID | write | `profile_id`; *attributes, identifiers, name, source, type* |  | `users:write`, `write` |
| `ZENDESK_UPDATE_VIEW` | Update View | write | `view_id`; *active, default, description, position, title* |  | `write` |
| `ZENDESK_UPDATE_ZENDESK_ACCOUNT_SETTINGS` | Update Zendesk Account Settings | write | `settings` |  | `write` |
| `ZENDESK_UPDATE_ZENDESK_ATTACHMENT` | Update Zendesk Attachment | write | `attachment_id`, `malware_access_override` |  | `write` |
| `ZENDESK_UPDATE_ZENDESK_CUSTOM_OBJECT_FIELD` | Update Custom Object Field | write | `custom_object_field_key_or_id`, `custom_object_key`; *active, description, position, regexp_for_validation, required, title* |  | `write` |
| `ZENDESK_UPDATE_ZENDESK_DYNAMIC_CONTENT_ITEM` | Update Dynamic Content Item | write | `dynamic_content_item_id`, `name` |  | `dynamic_content:write`, `write` |
| `ZENDESK_UPDATE_ZENDESK_ORGANIZATION` | Update Zendesk Organization | write | `data`, `organization_id` |  | `organizations:read`, `organizations:write`, `read`, `write` |
| `ZENDESK_UPDATE_ZENDESK_REQUESTS` | Update Request | write | `request_id`; *assignee_id, comment, custom_fields, group_id, priority, solved, status, tags* |  | `requests:write`, `write` |
| `ZENDESK_UPDATE_ZENDESK_TICKET` | Update Zendesk Ticket | write | `ticket_id`; *additional_collaborators, assignee_id, collaborator_ids, collaborators, comment_author_id, comment_body, comment_html_body, comment_public and 15 more* |  | `tickets:write`, `write` |
| `ZENDESK_UPSERT_CUSTOM_OBJECT_RECORD_BY_EXTERNAL_ID_OR_NAME` | Upsert Custom Object Record By External ID or Name | write | `custom_object_key`; *external_id, name, record_data* |  | `write` |
| `ZENDESK_VERIFY_SUBDOMAIN_AVAILABILITY` | Verify Subdomain Availability | read | `subdomain` |  | `read` |
| `ZENDESK_VERIFY_SUPPORT_ADDRESS` | Verify Support Address | write | `support_address_id` |  | `write` |
| `ZENDESK_VERIFY_USER_IDENTITY` | Verify User Identity | write | `user_id`, `user_identity_id` |  | `users:write`, `write` |
