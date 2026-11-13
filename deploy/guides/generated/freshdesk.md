# Freshdesk: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner pastes a key from the vendor's settings). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

184 tools: 85 read, 99 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `FRESHDESK_ADD_NOTE_TO_TICKET` | Add Note to Ticket | write | `body`, `ticket_id`; *notify_emails, private, user_id* |  |  |
| `FRESHDESK_ADD_TICKET_USER_ACCESS` | Add Ticket User Access | write | `ticket_id`, `user_ids` |  |  |
| `FRESHDESK_ADD_WATCHER` | Add Watcher to Ticket | write | `ticket_id` |  |  |
| `FRESHDESK_BULK_UNWATCH_TICKETS` | Bulk Unwatch Tickets | write | `ids` |  |  |
| `FRESHDESK_BULK_UPDATE_TICKETS` | Bulk Update Tickets | write | `ids`, `properties` |  |  |
| `FRESHDESK_CANCEL_CONTACT_IMPORT` | Cancel Contact Import | write | `import_id` |  |  |
| `FRESHDESK_CREATE_ADMIN` | Create Admin Group | write | `name`; *agent_ids, allow_agents_to_change_availability, auto_ticket_assign, business_hour_id, description, escalate_to, unassigned_for* |  |  |
| `FRESHDESK_CREATE_ADMIN_GROUP` | Create Admin Group | write | `name`, `type`; *agent_ids, description, escalate_to, unassigned_for* |  |  |
| `FRESHDESK_CREATE_ADMIN_TICKET_FIELD` | Create Admin Ticket Field | write | `label`, `label_for_customers`, `type`; *choices, customers_can_edit, displayed_to_customers, position, required_for_agents, required_for_closure, required_for_customers* |  |  |
| `FRESHDESK_CREATE_ADMIN_TICKET_FIELD_SECTION` | Create Admin Ticket Field Section | write | `choice_ids`, `id`, `label` |  |  |
| `FRESHDESK_CREATE_AGENTS` | Create Multiple Agents | write | `agents` |  |  |
| `FRESHDESK_CREATE_AUTOMATION_RULE` | Create Automation Rule | write | `actions`, `automation_type_id`, `conditions`, `name`; *active, description, events, operator, performer, position* |  |  |
| `FRESHDESK_CREATE_CANNED_RESPONSE` | Create Canned Response | write | `content_html`, `title`, `visibility`; *folder_id, group_ids* |  |  |
| `FRESHDESK_CREATE_CANNED_RESPONSE_BULK` | Bulk Create Canned Responses | write | `canned_responses` |  |  |
| `FRESHDESK_CREATE_CANNED_RESPONSE_FOLDER` | Create Canned Response Folder | write | `name` |  |  |
| `FRESHDESK_CREATE_COMPANIES` | Create Company | write | `name`; *account_tier, custom_fields, description, domains, health_score, industry, note, renewal_date* |  |  |
| `FRESHDESK_CREATE_CONTACT` | Create Contact | write | `custom_fields`, `name`; *address, avatar, company_id, description, email, job_title, language, mobile and 8 more* |  |  |
| `FRESHDESK_CREATE_CONTACT_FIELDS` | Create Contact Field | write | `label`, `label_for_customers`, `type`; *agents_can_edit, choices, customers_can_edit, displayed_for_agents, displayed_for_customers, editable_in_signup, position, quick_add_for_agent and 2 more* |  |  |
| `FRESHDESK_CREATE_DISCUSSION_CATEGORY` | Create Discussion Category | write | `description`, `name` |  |  |
| `FRESHDESK_CREATE_DISCUSSION_FORUM_TOPIC` | Create Discussion Forum Topic | write | `forum_id`, `message`, `title`; *locked, sticky* |  |  |
| `FRESHDESK_CREATE_DISCUSSION_TOPIC_COMMENT` | Create Discussion Topic Comment | write | `body`, `topic_id` |  |  |
| `FRESHDESK_CREATE_EMAIL` | Create Email Mailbox | write | `group_id`, `mailbox_type`, `name`, `product_id`, `support_email`; *default_reply_email* |  |  |
| `FRESHDESK_CREATE_FORUM` | Create Forum | write | `category_id`, `name`; *description, forum_type, forum_visibility* |  |  |
| `FRESHDESK_CREATE_SLA_POLICIES` | Create SLA Policy | write | `applicable_to`, `name`, `sla_target`; *active, description* |  |  |
| `FRESHDESK_CREATE_SOLUTION_ARTICLE` | Create Solution Article | write | `description`, `folder_id`, `title`; *agent_id, seo_data, status, tags* |  |  |
| `FRESHDESK_CREATE_SOLUTION_CATEGORY` | Create Solution Category | write | `name`; *description, visible_in_portals* |  |  |
| `FRESHDESK_CREATE_TICKET` | Create Ticket | write | `description`, `subject`; *attachments, cc_emails, company_id, custom_fields, due_by, email, email_config_id, facebook_id and 16 more* |  |  |
| `FRESHDESK_CREATE_TICKET_FORMS` | Create Ticket Form | write | `title`; *description, fields, portal_ids* |  |  |
| `FRESHDESK_CREATE_TICKET_OUTBOUND_EMAIL` | Create Ticket via Outbound Email | write | `description`, `email`, `email_config_id`, `priority`, `subject` |  |  |
| `FRESHDESK_CREATE_TICKET_TIME_ENTRY` | Create Ticket Time Entry | write | `ticket_id`, `time_spent`; *agent_id, billable, executed_at, note, start_time* |  |  |
| `FRESHDESK_CREATE_TRANSLATED_SOLUTION_ARTICLE` | Create Translated Solution Article | write | `article_id`, `description`, `language_code`, `status`, `title`; *seo_data, tags* |  |  |
| `FRESHDESK_CREATE_TRANSLATED_SOLUTION_CATEGORY` | Create Translated Solution Category | write | `id`, `language_code`, `name`; *description* |  |  |
| `FRESHDESK_CREATE_TRANSLATED_SOLUTION_FOLDER` | Create Translated Solution Folder | write | `folder_id`, `language_code`, `name`; *description* |  |  |
| `FRESHDESK_CURRENTLY_AUTHENTICATED_AGENT` | Currently Authenticated Agent | read |  |  |  |
| `FRESHDESK_DELETE_ADMIN_GROUP` | Delete Admin Group | write | `id` |  |  |
| `FRESHDESK_DELETE_AGENT` | Delete Agent | write | `agent_id` |  |  |
| `FRESHDESK_DELETE_AUTOMATION_RULE` | Delete Automation Rule | write | `automation_type_id`, `rule_id` |  |  |
| `FRESHDESK_DELETE_COMPANY` | Delete Company | write | `company_id` |  |  |
| `FRESHDESK_DELETE_COMPANY_FIELD` | Delete Company Field | write | `id` |  |  |
| `FRESHDESK_DELETE_CONTACT` | Delete Contact | write | `contact_id` |  |  |
| `FRESHDESK_DELETE_CONTACT_FIELD` | Delete Contact Field | write | `contact_field_id` |  |  |
| `FRESHDESK_DELETE_CONVERSATION` | Delete Conversation | write | `conversation_id` |  |  |
| `FRESHDESK_DELETE_DISCUSSIONS_FORUMS_FOLLOW` | Unmonitor Forum | write | `id`, `user_id` |  |  |
| `FRESHDESK_DELETE_DISCUSSIONS_TOPICS_FOLLOW` | Unmonitor Topic | write | `id`, `user_id` |  |  |
| `FRESHDESK_DELETE_DISCUSSION_CATEGORY` | Delete Discussion Category | write | `category_id` |  |  |
| `FRESHDESK_DELETE_DISCUSSION_COMMENT` | Delete Discussion Comment | write | `id` |  |  |
| `FRESHDESK_DELETE_DISCUSSION_FORUM` | Delete Discussion Forum | write | `id` |  |  |
| `FRESHDESK_DELETE_DISCUSSION_TOPIC` | Delete Discussion Topic | write | `id` |  |  |
| `FRESHDESK_DELETE_EMAIL_MAILBOXES` | Delete Email Mailbox | write | `id` |  |  |
| `FRESHDESK_DELETE_GROUP` | Delete Group | write | `group_id` |  |  |
| `FRESHDESK_DELETE_MULTIPLE_TICKETS` | Delete Multiple Tickets | write | `ids` |  |  |
| `FRESHDESK_DELETE_SECTION` | Delete Section | write | `section_id`, `ticket_field_id` |  |  |
| `FRESHDESK_DELETE_SKILL` | Delete Skill | write | `id` |  |  |
| `FRESHDESK_DELETE_SOLUTION_ARTICLE` | Delete Solution Article | write | `article_id` |  |  |
| `FRESHDESK_DELETE_SOLUTION_CATEGORY` | Delete Solution Category | write | `category_id` |  |  |
| `FRESHDESK_DELETE_SOLUTION_FOLDER` | Delete Solution Folder | write | `folder_id` |  |  |
| `FRESHDESK_DELETE_TICKET` | Delete Ticket | write | `ticket_id` |  |  |
| `FRESHDESK_DELETE_TICKET_FIELD` | Delete Ticket Field | write | `ticket_field_id` |  |  |
| `FRESHDESK_DELETE_TICKET_FORM` | Delete Ticket Form | write | `ticket_form_id` |  |  |
| `FRESHDESK_DELETE_TICKET_FORMS` | Delete Ticket Form Field | write | `field_id`, `form_id` |  |  |
| `FRESHDESK_DELETE_TICKET_SUMMARY` | Delete Ticket Summary | write | `ticket_id` |  |  |
| `FRESHDESK_DELETE_TICKET_USER_ACCESS` | Delete Ticket User Access | write | `ticket_id` |  |  |
| `FRESHDESK_DELETE_TIME_ENTRY` | Delete Time Entry | write | `time_entry_id` |  |  |
| `FRESHDESK_EXPORT_CONTACTS` | Export Contacts | read | `fields` |  |  |
| `FRESHDESK_GET_ACCOUNT` | Get Account | read |  |  |  |
| `FRESHDESK_GET_AGENT` | Get Agent | read | `agent_id` |  |  |
| `FRESHDESK_GET_AGENTS` | Search Agents Autocomplete | read | `term` |  |  |
| `FRESHDESK_GET_AGENT_AVAILABILITY` | Get Agent Availability | read | `agent_id` |  |  |
| `FRESHDESK_GET_BUSINESS_HOURS` | Get Business Hours | read |  |  |  |
| `FRESHDESK_GET_CANNED_RESPONSE_FOLDER` | Get Canned Response Folder | read | `folder_id` |  |  |
| `FRESHDESK_GET_CANNED_RESPONSE_FOLDERS` | Get Canned Response Folders | read |  |  |  |
| `FRESHDESK_GET_COMPANIES` | Get Companies | read | *page, per_page* |  |  |
| `FRESHDESK_GET_COMPANY` | Get Company | read | `company_id` |  |  |
| `FRESHDESK_GET_COMPANY_FIELDS` | Get Company Fields | read |  |  |  |
| `FRESHDESK_GET_CONTACT` | Get Contact | read | `contact_id` |  |  |
| `FRESHDESK_GET_CONTACTS` | Get Contacts | read | *company_id, email, mobile, page, per_page, phone, state, updated_since* |  |  |
| `FRESHDESK_GET_CONTACT_FIELDS` | Get Contact Fields | read |  |  |  |
| `FRESHDESK_GET_DISCUSSION_CATEGORY` | Get Discussion Category | read | `id` |  |  |
| `FRESHDESK_GET_DISCUSSION_TOPIC` | View Discussion Topic | read | `id` |  |  |
| `FRESHDESK_GET_EMAIL_MAILBOX` | Get Email Mailbox | read | `id` |  |  |
| `FRESHDESK_GET_FOLDER_RESPONSES` | Get Folder Responses | read | `folder_id` |  |  |
| `FRESHDESK_GET_IMPORTED_CONTACTS` | Get Imported Contacts | read | *status* |  |  |
| `FRESHDESK_GET_JOB` | Get Job | read | `id` |  |  |
| `FRESHDESK_GET_SCENARIO_AUTOMATIONS_JSON` | List All Scenario Automations | read |  |  |  |
| `FRESHDESK_GET_SEARCH` | Filter Tickets | read | `query`; *page* |  |  |
| `FRESHDESK_GET_TICKETS` | Get Tickets | read | *agent_id, created_since, filter, group_id, page, per_page, priority, requester_id and 3 more* |  |  |
| `FRESHDESK_GET_TICKET_TIME_ENTRIES` | Get Ticket Time Entries | read | `ticket_id`; *page* |  |  |
| `FRESHDESK_GET_TICKET_USER_ACCESS` | Get Ticket User Access | read | `ticket_id` |  |  |
| `FRESHDESK_GET_TRANSLATED_SOLUTION_CATEGORY` | Get Translated Solution Category | read | `id`, `language_code` |  |  |
| `FRESHDESK_HARD_DELETE_CONTACT` | Hard Delete Contact | write | `contact_id`; *force* |  |  |
| `FRESHDESK_IMPORT_CONTACT` | Import Contact | write | `company_name_column_index`, `email_column_index`, `name_column_index`; *file, job_title_column_index, raw_content, raw_filename, raw_mimetype* |  |  |
| `FRESHDESK_LIST_ADMIN_GROUPS` | List Admin Groups | read | *page, per_page* |  |  |
| `FRESHDESK_LIST_ALL_AGENTS_IN_A_GROUP` | List All Agents in a Group | read | `group_id` |  |  |
| `FRESHDESK_LIST_ALL_SECTIONS_FOR_TICKET_FIELD` | List All Sections for a Ticket Field | read | `ticket_field_id` |  |  |
| `FRESHDESK_LIST_ALL_SKILLS` | List All Skills | read |  |  |  |
| `FRESHDESK_LIST_ALL_TICKET_CONVERSATIONS` | List All Ticket Conversations | read | `ticket_id`; *page, per_page* |  |  |
| `FRESHDESK_LIST_AUTOMATION_RULES` | List Automation Rules | read | `automation_type_id` |  |  |
| `FRESHDESK_LIST_DISCUSSIONS` | List All Forum Categories | read | *page* |  |  |
| `FRESHDESK_LIST_DISCUSSIONS_TOPICS` | List All Topics in a Forum | read | `forum_id`; *page, per_page* |  |  |
| `FRESHDESK_LIST_EMAIL` | List Email Mailboxes | read | *active, group_id, order_by, order_type, page, per_page, product_id, support_email* |  |  |
| `FRESHDESK_LIST_EMAIL_CONFIGS` | List All Email Configs | read |  |  |  |
| `FRESHDESK_LIST_FORUMS_IN_CATEGORY` | List Forums in Category | read | `category_id`; *page, per_page* |  |  |
| `FRESHDESK_LIST_GROUPS` | List Groups | read | *page, per_page* |  |  |
| `FRESHDESK_LIST_MONITORED_TOPICS` | List Monitored Topics | read | *page, user_id* |  |  |
| `FRESHDESK_LIST_PARTICIPATED_TOPICS` | List Participated Topics | read | *page, user_id* |  |  |
| `FRESHDESK_LIST_PRODUCTS` | List All Products | read |  |  |  |
| `FRESHDESK_LIST_ROLES` | List All Roles | read | *page, per_page* |  |  |
| `FRESHDESK_LIST_SATISFACTION_RATINGS` | List Satisfaction Ratings | read | *created_since, page, per_page* |  |  |
| `FRESHDESK_LIST_SLA_POLICIES` | List All SLA Policies | read |  |  |  |
| `FRESHDESK_LIST_SOLUTIONS_CATEGORY_FOLDERS_TRANSLATED` | List Translated Category Folders | read | `category_id`, `language_code`; *page* |  |  |
| `FRESHDESK_LIST_SOLUTION_ARTICLES` | List Solution Articles | read | `folder_id`; *page* |  |  |
| `FRESHDESK_LIST_SOLUTION_CATEGORIES` | List Solution Categories | read | *page* |  |  |
| `FRESHDESK_LIST_SOLUTION_FOLDERS` | List Solution Folders | read | `category_id`; *page* |  |  |
| `FRESHDESK_LIST_SOLUTION_SUBFOLDERS` | List Solution Subfolders | read | `folder_id`; *page* |  |  |
| `FRESHDESK_LIST_SPECIFIC_SECTION_DETAILS` | List a Specific Section Details | read | `section_id`, `ticket_field_id` |  |  |
| `FRESHDESK_LIST_SURVEYS` | List All Surveys | read | *state* |  |  |
| `FRESHDESK_LIST_TICKETS_WATCHERS` | List Ticket Watchers | read | `id` |  |  |
| `FRESHDESK_LIST_TICKET_FIELDS` | List All Ticket Fields | read |  |  |  |
| `FRESHDESK_LIST_TICKET_FORMS` | List Ticket Forms | read | *include, page, per_page* |  |  |
| `FRESHDESK_LIST_TICKET_SATISFACTION_RATINGS` | List Ticket Satisfaction Ratings | read | `ticket_id`; *page, per_page* |  |  |
| `FRESHDESK_LIST_TIME_ENTRIES` | List Time Entries | read | *agent_id, billable, company_id, executed_after, executed_before, page, per_page* |  |  |
| `FRESHDESK_LIST_TOPIC_COMMENTS2` | List Topic Comments (Paginated) | read | `id`; *page, per_page* |  |  |
| `FRESHDESK_LIST_TRANSLATED_SOLUTION_ARTICLES` | List Translated Solution Articles | read | `id`, `language_code`; *page* |  |  |
| `FRESHDESK_LIST_TRANSLATED_SUBFOLDERS` | List Translated Subfolders | read | `folder_id`, `language_code`; *page* |  |  |
| `FRESHDESK_MERGE_TICKETS` | Merge Tickets | write | `primary_id`, `ticket_ids`; *convert_recipients_to_cc, note_in_primary, note_in_secondary* |  |  |
| `FRESHDESK_PATCH_ADMIN_GROUPS_AGENTS` | Update Group Agents | write | `agents`, `id` |  |  |
| `FRESHDESK_REMOVE_WATCHER` | Remove Watcher from Ticket | write | `ticket_id` |  |  |
| `FRESHDESK_REPLY_TO_FORWARD_TICKET` | Reply to Forward Ticket | write | `body`, `ticket_id`, `to_emails`; *bcc_emails, cc_emails, from_email* |  |  |
| `FRESHDESK_REPLY_TO_TICKET` | Reply to Ticket | write | `body`, `ticket_id`; *bcc_emails, cc_emails, from_email, to_emails, user_id* |  |  |
| `FRESHDESK_SEARCH_AGENTS` | Search Agents | read | *email, mobile, page, per_page, phone, state* |  |  |
| `FRESHDESK_SEARCH_COMPANIES` | Search Companies | read | `query`; *page* |  |  |
| `FRESHDESK_SEARCH_COMPANY` | Search Companies Autocomplete | read | `name` |  |  |
| `FRESHDESK_SEARCH_CONTACTS` | Search Contacts | read | `query`; *page* |  |  |
| `FRESHDESK_SEARCH_CONTACTS_AUTOCOMPLETE` | Search Contacts Autocomplete | read | `term` |  |  |
| `FRESHDESK_SEARCH_SOLUTION_ARTICLES` | Search Solution Articles | read | `term`; *page* |  |  |
| `FRESHDESK_TOGGLE_TIMER` | Toggle Timer | write | `time_entry_id` |  |  |
| `FRESHDESK_UPDATE_ADMIN_GROUP` | Update Admin Group | write | `id`; *agent_ids, business_calendar_id, description, escalate_to, name, unassigned_for* |  |  |
| `FRESHDESK_UPDATE_ADMIN_TICKET_FIELD` | Update Admin Ticket Field | write | `id`; *customers_can_edit, displayed_to_customers, label, label_for_customers, position, required_for_agents, required_for_closure, required_for_customers* |  |  |
| `FRESHDESK_UPDATE_ADMIN_TICKET_FIELD_SECTION` | Update Admin Ticket Field Section | write | `id`, `section_id`; *choice_ids, label, ticket_field_ids* |  |  |
| `FRESHDESK_UPDATE_AGENT` | Update Agent | write | `id`; *agent_type, email, focus_mode, group_ids, language, occasional, role_ids, signature and 3 more* |  |  |
| `FRESHDESK_UPDATE_AGENT_AVAILABILITY` | Update Agent Availability | write | `agent_id`, `channel_availability` |  |  |
| `FRESHDESK_UPDATE_AUTOMATIONS` | Update Automation Rule | write | `automation_type_id`, `id`; *actions, active, conditions, description, name, position* |  |  |
| `FRESHDESK_UPDATE_CANNED_RESPONSE` | Update Canned Response | write | `id`; *content_html, folder_id, group_ids, title, visibility* |  |  |
| `FRESHDESK_UPDATE_CANNED_RESPONSE_FOLDER` | Update Canned Response Folder | write | `folder_id`, `name` |  |  |
| `FRESHDESK_UPDATE_COMPANIES` | Update Company | write | `company_id`; *account_tier, custom_fields, description, domains, health_score, industry, name, note and 1 more* |  |  |
| `FRESHDESK_UPDATE_CONTACT` | Update Contact | write | `id`; *address, company_id, custom_fields, description, email, job_title, language, mobile and 8 more* |  |  |
| `FRESHDESK_UPDATE_CONTACT_FIELD` | Update Contact Field | write | `contact_field_id`; *customers_can_edit, displayed_for_customers, label, label_for_customers, required_for_agents, required_for_customers* |  |  |
| `FRESHDESK_UPDATE_CONTACT_MAKE_AGENT` | Make Agent | write | `id`, `role_ids`, `ticket_scope`; *occasional* |  |  |
| `FRESHDESK_UPDATE_CONVERSATIONS` | Update Conversation | write | `body`, `id` |  |  |
| `FRESHDESK_UPDATE_DISCUSSION_CATEGORY` | Update Discussion Category | write | `id`; *description, name* |  |  |
| `FRESHDESK_UPDATE_DISCUSSION_COMMENT` | Update Discussion Comment | write | `body`, `id` |  |  |
| `FRESHDESK_UPDATE_DISCUSSION_FORUM` | Update Discussion Forum | write | `id`; *company_ids, description, forum_category_id, forum_type, forum_visibility, name* |  |  |
| `FRESHDESK_UPDATE_DISCUSSION_TOPIC` | Update Discussion Topic | write | `id`; *forum_id, locked, message, stamp_type, sticky, title* |  |  |
| `FRESHDESK_UPDATE_EMAIL_SETTINGS` | Update Email Settings | write | *allow_agent_to_initiate_conversation, allow_wildcard_ticket_create, auto_response_detector_toggle, create_requester_using_reply_to, email_subject_match, extended_quoted_text, multiple_to, original_sender_as_requester_for_forward and 4 more* |  |  |
| `FRESHDESK_UPDATE_NOTIFICATION_EMAIL_BCC` | Update Automatic BCC Emails | write | `emails` |  |  |
| `FRESHDESK_UPDATE_SLA_POLICIES` | Update SLA Policy | write | `id`; *active, applicable_to, description, name, sla_target* |  |  |
| `FRESHDESK_UPDATE_SOLUTIONS` | Update Solution Article | write | `id`; *agent_id, description, language_code, seo_data, status, tags, title* |  |  |
| `FRESHDESK_UPDATE_SOLUTION_CATEGORY` | Update Solution Category | write | `category_id`; *description, name* |  |  |
| `FRESHDESK_UPDATE_SOLUTION_FOLDER` | Update Solution Folder | write | `folder_id`; *company_ids, company_segment_ids, contact_segment_ids, description, language_code, name, parent_folder_id, visibility* |  |  |
| `FRESHDESK_UPDATE_TICKET` | Update Ticket | write | `ticket_id`; *attachments, company_id, custom_fields, description, due_by, email, email_config_id, facebook_id and 19 more* |  |  |
| `FRESHDESK_UPDATE_TICKET_BULK_WATCH` | Bulk Add Watcher to Tickets | write | `ids` |  |  |
| `FRESHDESK_UPDATE_TICKET_FORMS` | Update Ticket Form | write | `ticket_form_id`; *description, fields, include, title* |  |  |
| `FRESHDESK_UPDATE_TICKET_SUMMARY` | Update Ticket Summary | write | `body`, `ticket_id` |  |  |
| `FRESHDESK_UPDATE_TICKET_USER_ACCESS` | Update Ticket User Access | write | `id`, `user_ids` |  |  |
| `FRESHDESK_UPDATE_TIME_ENTRY` | Update Time Entry | write | `id`; *agent_id, billable, executed_at, note, start_time, time_spent, timer_running* |  |  |
| `FRESHDESK_VIEW_AUTOMATIONS` | View Automation Rule | read | `automation_type_id`, `id` |  |  |
| `FRESHDESK_VIEW_A_CANNED_RESPONSE` | View a Canned Response | read | `canned_response_id` |  |  |
| `FRESHDESK_VIEW_BUSINESS_HOUR` | View a Business Hour | read | `business_hours_id` |  |  |
| `FRESHDESK_VIEW_COMPANY_FIELD` | View a Company Field | read | `company_field_id` |  |  |
| `FRESHDESK_VIEW_CONTACT_FIELD` | View a Contact Field | read | `contact_field_id` |  |  |
| `FRESHDESK_VIEW_EMAIL` | View Email Mailbox Settings | read |  |  |  |
| `FRESHDESK_VIEW_EMAIL_CONFIGS` | View Email Config | read | `id` |  |  |
| `FRESHDESK_VIEW_GROUP` | View Group | read | `group_id` |  |  |
| `FRESHDESK_VIEW_NOTIFICATIONS` | View Automatic BCC Emails | read |  |  |  |
| `FRESHDESK_VIEW_PRODUCTS` | View Product | read | `id` |  |  |
| `FRESHDESK_VIEW_ROLES` | View a Role | read | `id` |  |  |
| `FRESHDESK_VIEW_SETTINGS` | View Helpdesk Settings | read |  |  |  |
| `FRESHDESK_VIEW_SKILL` | View a Skill | read | `id` |  |  |
| `FRESHDESK_VIEW_SOLUTIONS` | View Solution Article | read | `article_id`; *language_code* |  |  |
| `FRESHDESK_VIEW_SOLUTION_CATEGORY` | View Solution Category | read | `category_id` |  |  |
| `FRESHDESK_VIEW_SOLUTION_FOLDER` | View Solution Folder | read | `folder_id`; *language_code* |  |  |
| `FRESHDESK_VIEW_TICKET` | View Ticket | read | `ticket_id`; *include* |  |  |
| `FRESHDESK_VIEW_TICKET_FIELD` | View a Ticket Field | read | `ticket_field_id` |  |  |
| `FRESHDESK_VIEW_TICKET_FORM_FIELD` | View a Ticket Form's Field | read | `field_id`, `form_id` |  |  |
