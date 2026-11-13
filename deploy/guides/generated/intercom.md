# Intercom: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

138 tools: 68 read, 70 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `INTERCOM_ADD_SUBSCRIPTION_TO_A_CONTACT` | Add subscription to a contact | write | `consent_type`, `contact_id`, `id` |  | `contacts` |
| `INTERCOM_ADD_TAG_TO_CONTACT` | Add tag to contact | write | `contact_id`, `id` |  |  |
| `INTERCOM_ARCHIVE_CONTACT` | Archive contact | write | `contact_id` |  |  |
| `INTERCOM_ASSIGN_CONVERSATION` | Assign conversation | write | `conversation_id`; *admin_id, assignee_id, team_id* |  |  |
| `INTERCOM_ATTACH_CONTACT_TO_COMPANY` | Attach contact to company | write | `company_id`, `contact_id` |  |  |
| `INTERCOM_ATTACH_CONTACT_TO_CONVERSATION` | Attach contact to conversation | write | `admin_id`, `conversation_id`, `customer` |  |  |
| `INTERCOM_ATTACH_TAG_TO_CONVERSATION` | Attach tag to conversation | write | `admin_id`, `conversation_id`, `id` |  |  |
| `INTERCOM_ATTACH_TAG_TO_TICKET` | Attach tag to ticket | write | `admin_id`, `id`, `ticket_id` |  |  |
| `INTERCOM_BLOCK_CONTACT` | Block contact | write | `contact_id` |  |  |
| `INTERCOM_CANCEL_DATA_EXPORT` | Cancel data export | write | `job_identifier` |  |  |
| `INTERCOM_CLOSE_CONVERSATION` | Close conversation | write | `admin_id`, `conversation_id`; *body* |  |  |
| `INTERCOM_CREATE_AN_ARTICLE` | Create an article | write | `author_id`, `title`; *body, description, parent_id, parent_type, state, translated__content__ar__author__id, translated__content__ar__body, translated__content__ar__created__at and 331 more* |  | `Read and Write Articles` |
| `INTERCOM_CREATE_A_COLLECTION` | Create a collection | write | `name`; *description, help_center_id, parent_id, translated__content__ar__description, translated__content__ar__name, translated__content__ar__type, translated__content__bg__description, translated__content__bg__name and 107 more* |  | `read_and_write_articles` |
| `INTERCOM_CREATE_A_NOTE` | Create a note | write | `body`, `id`; *admin_id, contact_id* |  | `write_users_and_companies` |
| `INTERCOM_CREATE_CONTACT` | Create contact | write | *avatar, custom_attributes, email, external_id, last_seen_at, name, owner_id, phone and 3 more* |  |  |
| `INTERCOM_CREATE_CONTENT_IMPORT_SOURCE` | Create content import source | write | `sync_behavior`, `url`; *status* |  |  |
| `INTERCOM_CREATE_CONVERSATION` | Create conversation | write | `body`; *from_admin_id, from_contact_id, from_user_id, message_type, subject* |  |  |
| `INTERCOM_CREATE_DATA_ATTRIBUTE` | Create data attribute | write | `data_type`, `model`, `name`; *description, messenger_writable, options* |  |  |
| `INTERCOM_CREATE_DATA_EVENT` | Create data event | write | `created_at`, `event_name`; *email, id, metadata, user_id* |  |  |
| `INTERCOM_CREATE_DATA_EXPORT` | Create data export | read | `created_at_after`, `created_at_before` |  |  |
| `INTERCOM_CREATE_EXTERNAL_PAGE` | Create external page | write | `external_id`, `html`, `locale`, `source_id`, `title`, `url`; *ai_agent_availability, ai_copilot_availability* |  |  |
| `INTERCOM_CREATE_HELP_CENTER_SECTION` | Create help center section | write | `name`, `parent_id`; *translated_content* |  | `Read and Write Articles` |
| `INTERCOM_CREATE_INTERNAL_ARTICLE` | Create internal article | write | `author_id`, `owner_id`, `title`; *body, locale* |  |  |
| `INTERCOM_CREATE_OR_UPDATE_A_COMPANY` | Create or update a company | write | *company_id, custom_attributes, industry, monthly_spend, name, plan, remote_created_at, size and 1 more* |  | `Write users and companies` |
| `INTERCOM_CREATE_TAG` | Create or update tag | write | `name`; *companies, id, users* |  |  |
| `INTERCOM_CREATE_TICKET` | Create ticket | write | `contacts`, `ticket_type_id`; *assignment, company_id, conversation_to_link_id, created_at, ticket_attributes* |  |  |
| `INTERCOM_CREATE_TICKET_TYPE` | Create a ticket type | write | `category`, `name`; *description, icon, is_internal* |  |  |
| `INTERCOM_CREATE_TICKET_TYPE_ATTRIBUTE` | Create ticket type attribute | write | `data_type`, `description`, `name`, `ticket_type_id`; *allow_multiple_values, list_items, multiline, required_to_create, required_to_create_for_contacts, visible_on_create, visible_to_contacts* |  |  |
| `INTERCOM_DATA_EVENT_SUMMARIES` | Create event summaries | write | `event_summaries`, `user_id` |  |  |
| `INTERCOM_DELETE_AN_ARTICLE` | Delete an article | write | `id` |  | `Read and Write Articles` |
| `INTERCOM_DELETE_A_COLLECTION` | Delete a collection | write | `id` |  | `Read and Write Articles` |
| `INTERCOM_DELETE_A_COMPANY` | Delete a company | write | `id` |  | `Read and list users and companies` |
| `INTERCOM_DELETE_A_TAG_DELETE_TAG` | Delete a tag | write | `tag_id` |  |  |
| `INTERCOM_DELETE_A_VISITOR` | Delete a visitor | write | `id` |  |  |
| `INTERCOM_DELETE_CONTACT` | Delete a contact | write | `contact_id` |  |  |
| `INTERCOM_DELETE_CONTENT_IMPORT_SOURCE` | Delete content import source | write | `source_id` |  |  |
| `INTERCOM_DELETE_CONVERSATION` | Delete conversation | write | `conversation_id`; *retain_metrics* |  | `Write conversations` |
| `INTERCOM_DELETE_EXTERNAL_PAGE` | Delete external page | write | `page_id` |  |  |
| `INTERCOM_DELETE_INTERNAL_ARTICLE` | Delete internal article | write | `internal_article_id` |  |  |
| `INTERCOM_DELETE_TICKET` | Delete ticket | write | `ticket_id` |  |  |
| `INTERCOM_DETACH_A_CONTACT` | Detach a contact from tag | write | `contact_id`, `tag_id` |  |  |
| `INTERCOM_DETACH_CONTACT_FROM_COMPANY` | Detach contact from company | write | `company_id`, `contact_id` |  |  |
| `INTERCOM_DETACH_TAG_FROM_CONVERSATION` | Detach tag from conversation | write | `admin_id`, `conversation_id`, `tag_id` |  |  |
| `INTERCOM_DETACH_TAG_FROM_TICKET` | Detach tag from ticket | write | `admin_id`, `tag_id`, `ticket_id` |  |  |
| `INTERCOM_DOWNLOAD_DATA_EXPORT` | Download data export | read | `job_identifier` |  |  |
| `INTERCOM_ENQUEUE_CREATE_TICKET` | Enqueue create ticket | write | `contacts`, `ticket_type_id`; *assignment, company_id, conversation_to_link_id, created_at, skip_notifications, ticket_attributes* |  |  |
| `INTERCOM_EXPORT_WORKFLOW` | Export workflow | read | `ruleset_id` |  |  |
| `INTERCOM_FIND_TAG` | Find a tag | read | `tag_id` |  |  |
| `INTERCOM_GET_A_CONTACT` | Get a contact | read | `id` |  | `Read one user and one company` |
| `INTERCOM_GET_CONTENT_IMPORT_SOURCE` | Get content import source | read | `source_id` |  |  |
| `INTERCOM_GET_CONVERSATION` | Get conversation | read | `conversation_id`; *display_as* |  |  |
| `INTERCOM_GET_COUNTS` | Get entity counts | read | *type* |  |  |
| `INTERCOM_GET_CUSTOM_OBJECT_INSTANCE_BY_EXTERNAL_ID` | Get custom object instance by external ID | write | `custom_object_type_identifier`, `external_id` |  |  |
| `INTERCOM_GET_EXTERNAL_PAGE` | Get external page | read | `page_id` |  |  |
| `INTERCOM_GET_TICKET` | Get ticket | read | `ticket_id` |  |  |
| `INTERCOM_GET_TICKET_TYPE` | Get a ticket type | read | `ticket_type_id` |  |  |
| `INTERCOM_IDENTIFY_AN_ADMIN` | Identify an admin | read |  |  | `read_admins` |
| `INTERCOM_JOBS_STATUS` | Retrieve job status | read | `job_id` |  |  |
| `INTERCOM_LIST_ALL_ACTIVITY_LOGS` | List all activity logs | read | `created_at_after`; *created_at_before* |  | `Read admin activity logs` |
| `INTERCOM_LIST_ALL_ADMINS` | List all admins | read |  |  | `read_admins` |
| `INTERCOM_LIST_ALL_ARTICLES` | List all articles | read |  |  | `Read and List Articles` |
| `INTERCOM_LIST_ALL_COLLECTIONS` | List all collections | read |  |  | `read_and_list_articles` |
| `INTERCOM_LIST_ALL_COMPANIES` | List all companies | read | *order, page, per_page* |  | `Read and list users and companies` |
| `INTERCOM_LIST_ALL_HELP_CENTERS` | List all help centers | read |  |  | `read_and_list_articles` |
| `INTERCOM_LIST_ALL_MACROS` | List all macros | read | *per_page, starting_after, updated_since* |  |  |
| `INTERCOM_LIST_ALL_NOTES` | List all notes | read | `id` |  | `Read and list users and companies` |
| `INTERCOM_LIST_ATTACHED_COMPANIES_FOR_CONTACT` | List attached companies for contact | read | `id` |  | `read_and_list_users_and_companies` |
| `INTERCOM_LIST_ATTACHED_CONTACTS` | List attached contacts | read | `id` |  | `Read and list users and companies` |
| `INTERCOM_LIST_ATTACHED_SEGMENTS_FOR_COMPANIES` | List attached segments for companies | read | `id` |  | `Read and list users and companies` |
| `INTERCOM_LIST_ATTACHED_SEGMENTS_FOR_CONTACT` | List attached segments for contact | read | `contact_id` |  | `read_and_list_users_and_companies` |
| `INTERCOM_LIST_AWAY_STATUS_REASONS` | List away status reasons | read |  |  |  |
| `INTERCOM_LIST_CALLS` | List calls | read | *page, per_page* |  |  |
| `INTERCOM_LIST_CALLS_WITH_TRANSCRIPTS` | List calls with transcripts | read | `conversation_ids` |  |  |
| `INTERCOM_LIST_COMPANY_NOTES` | List company notes | read | `company_id` |  | `Read and list users and companies` |
| `INTERCOM_LIST_CONTACTS` | List all contacts | read | *per_page, starting_after* |  |  |
| `INTERCOM_LIST_CONTENT_IMPORT_SOURCES` | List content import sources | read | *per_page, starting_after* |  |  |
| `INTERCOM_LIST_CONVERSATIONS` | List conversations | read | *per_page, starting_after* |  |  |
| `INTERCOM_LIST_DATA_ATTRIBUTES` | List data attributes | read | *include_archived, model* |  | `Write data attributes` |
| `INTERCOM_LIST_DATA_EVENTS` | List data events | read | *email, intercom_user_id, summary, user_id* |  | `Read events` |
| `INTERCOM_LIST_EXTERNAL_PAGES` | List external pages | read | *per_page, starting_after* |  |  |
| `INTERCOM_LIST_HELP_CENTER_SECTIONS` | List help center sections | read | *per_page, starting_after* |  | `Read and List articles` |
| `INTERCOM_LIST_INTERNAL_ARTICLES` | List internal articles | read | *per_page, starting_after* |  | `Read and List articles` |
| `INTERCOM_LIST_NEWS_ITEMS` | List all news items | read |  |  |  |
| `INTERCOM_LIST_SEGMENTS` | List all segments | read | *include_count* |  |  |
| `INTERCOM_LIST_SUBSCRIPTIONS_FOR_A_CONTACT` | List subscriptions for a contact | read | `contact_id` |  | `read_and_list_users_and_companies` |
| `INTERCOM_LIST_SUBSCRIPTION_TYPES` | List subscription types | read |  |  |  |
| `INTERCOM_LIST_TAGS` | List all tags | read |  |  |  |
| `INTERCOM_LIST_TAGS_ATTACHED_TO_A_CONTACT` | List tags attached to a contact | read | `contact_id` |  | `tags` |
| `INTERCOM_LIST_TEAMS` | List all teams | read |  |  |  |
| `INTERCOM_LIST_TICKET_STATES` | List all ticket states | read |  |  |  |
| `INTERCOM_LIST_TICKET_TYPES` | List all ticket types | read |  |  |  |
| `INTERCOM_MERGE_A_LEAD_AND_A_USER` | Merge a lead and a user | write | *from, into* |  | `Write users and companies` |
| `INTERCOM_REDACT_CONVERSATION_CONTENT` | Redact conversation content | write | `conversation_id`, `conversation_part_id` |  | `Write conversations` |
| `INTERCOM_REGISTER_FIN_VOICE_CALL` | Register Fin Voice call | write | `call_id`, `phone_number`; *data, source* |  | `Create phone call redirects` |
| `INTERCOM_REMOVE_SUBSCRIPTION_FROM_A_CONTACT` | Remove subscription from a contact | write | `contact_id`, `id` |  | `write_users_and_companies` |
| `INTERCOM_REMOVE_TAG_FROM_A_CONTACT` | Remove tag from a contact | write | `contact_id`, `id` |  | `write_tags` |
| `INTERCOM_REOPEN_CONVERSATION` | Reopen conversation | write | `admin_id`, `conversation_id`; *body* |  |  |
| `INTERCOM_REPLY_TICKET` | Reply to ticket | write | `body`, `message_type`, `ticket_id`; *admin_id, attachment_urls, created_at, email, intercom_user_id, user_id* |  |  |
| `INTERCOM_REPLY_TO_CONVERSATION` | Reply to conversation | write | `admin_id`, `conversation_id`, `message_body`; *attachment_urls, message_type* |  |  |
| `INTERCOM_RETRIEVE_AN_ADMIN` | Retrieve an admin | read | `id` |  | `read:admins` |
| `INTERCOM_RETRIEVE_AN_ARTICLE` | Retrieve an article | read | `id` |  | `Read and List Articles` |
| `INTERCOM_RETRIEVE_A_COLLECTION` | Retrieve a collection | read | `id` |  | `read_and_list_articles` |
| `INTERCOM_RETRIEVE_A_COMPANY_BY_ID` | Retrieve a company by id | read | `id` |  | `Read and list users and companies` |
| `INTERCOM_RETRIEVE_A_HELP_CENTER` | Retrieve a help center | read | `id` |  | `read_and_list_articles` |
| `INTERCOM_RETRIEVE_A_JOB_STATUS` | Retrieve job status | read | `job_identifier` |  |  |
| `INTERCOM_RETRIEVE_A_MACRO` | Retrieve a macro | read | `id` |  |  |
| `INTERCOM_RETRIEVE_A_SEGMENT` | Retrieve a segment | read | `segment_id`; *include_count* |  |  |
| `INTERCOM_RETRIEVE_COMPANIES` | Retrieve companies | read | *company_id, name, page, per_page, segment_id, tag_id* |  | `read:users`, `read:companies` |
| `INTERCOM_RETRIEVE_INTERNAL_ARTICLE` | Retrieve internal article | read | `internal_article_id` |  |  |
| `INTERCOM_RETRIEVE_NOTE` | Retrieve note | read | `note_id` |  |  |
| `INTERCOM_RETRIEVE_TEAM` | Retrieve a team | read | `team_id` |  |  |
| `INTERCOM_RETRIEVE_VISITOR_WITH_USER_ID` | Retrieve visitor with user ID | read | `user_id` |  |  |
| `INTERCOM_SCROLL_OVER_ALL_COMPANIES` | Scroll over all companies | read | *scroll_param* |  | `Read and list users and companies` |
| `INTERCOM_SEARCH_CONTACTS` | Search contacts | read | `query`; *per_page, starting_after* |  |  |
| `INTERCOM_SEARCH_CONVERSATIONS` | Search conversations | read | `query`; *per_page, sort_field, sort_order, starting_after* |  |  |
| `INTERCOM_SEARCH_FOR_ARTICLES` | Search for articles | read | *help_center_id, highlight, phrase, state* |  | `read:articles` |
| `INTERCOM_SEARCH_INTERNAL_ARTICLES` | Search internal articles | read | *folder_id, per_page, starting_after* |  |  |
| `INTERCOM_SEARCH_TICKETS` | Search tickets | read | `query`; *pagination* |  |  |
| `INTERCOM_SEND_OUTBOUND_MESSAGE` | Send outbound message | write | `from_admin_id`, `message_type`, `to`; *bcc, body, brand_id, cc, components, create_conversation_without_contact_reply, created_at, locale and 2 more* |  | `Write conversations` |
| `INTERCOM_SET_ADMIN_TO_AWAY` | Set admin to away | write | `admin_id`; *away_mode_enabled, away_mode_reassign, away_status_reason_id* |  |  |
| `INTERCOM_SET_AN_ADMIN_TO_AWAY` | Set an admin to away | write | `id`; *away_mode_enabled, away_mode_reassign* |  | `update_admins` |
| `INTERCOM_SHOW_CALL` | Show call | read | `call_id` |  |  |
| `INTERCOM_SHOW_CALL_TRANSCRIPT` | Show call transcript | read | `call_id` |  |  |
| `INTERCOM_SHOW_CONTACT_BY_EXTERNAL_ID` | Show contact by external ID | read | `external_id` |  |  |
| `INTERCOM_UNARCHIVE_CONTACT` | Unarchive contact | write | `contact_id` |  |  |
| `INTERCOM_UPDATE_AN_ARTICLE` | Update an article | write | `id`; *author_id, body, description, parent_id, parent_type, state, title, translated__content__ar__author__id and 333 more* |  | `Read and Write Articles` |
| `INTERCOM_UPDATE_A_COLLECTION` | Update a collection | write | `id`; *description, name, parent_id, translated__content__ar__description, translated__content__ar__name, translated__content__ar__type, translated__content__bg__description, translated__content__bg__name and 107 more* |  | `read_and_write_articles` |
| `INTERCOM_UPDATE_A_COMPANY` | Update a company | write | `id`; *custom_attributes, industry, monthly_spend, name, plan, remote_created_at, size, website* |  | `Write users and companies` |
| `INTERCOM_UPDATE_A_CONTACT` | Update a contact | write | `id`; *avatar, companies, custom_attributes, email, external_id, last_seen_at, name, owner_id and 4 more* |  | `write_users_and_companies` |
| `INTERCOM_UPDATE_CONTACT` | Update contact | write | `contact_id`; *avatar, custom_attributes, email, external_id, last_seen_at, name, owner_id, phone and 3 more* |  |  |
| `INTERCOM_UPDATE_CONTENT_IMPORT_SOURCE` | Update content import source | write | `source_id`, `sync_behavior`, `url`; *status* |  |  |
| `INTERCOM_UPDATE_CONVERSATION` | Update conversation | write | `conversation_id`; *company_id, custom_attributes, display_as, read, title* |  | `Write conversations` |
| `INTERCOM_UPDATE_DATA_ATTRIBUTE` | Update data attribute | write | `data_attribute_id`; *archived, description, messenger_writable, options* |  |  |
| `INTERCOM_UPDATE_EXTERNAL_PAGE` | Update external page | write | `external_id`, `html`, `locale`, `page_id`, `source_id`, `title`, `url`; *fin_availability* |  |  |
| `INTERCOM_UPDATE_INTERNAL_ARTICLE` | Update internal article | write | `author_id`, `body`, `internal_article_id`, `owner_id`, `title` |  |  |
| `INTERCOM_UPDATE_TICKET` | Update ticket | write | `ticket_id`; *admin_id, assignee_id, company_id, is_shared, open, snoozed_until, ticket_attributes, ticket_state_id* |  |  |
| `INTERCOM_UPDATE_TICKET_TYPE` | Update a ticket type | write | `ticket_type_id`; *archived, category, description, icon, is_internal, name* |  |  |
| `INTERCOM_UPDATE_TICKET_TYPE_ATTRIBUTE` | Update ticket type attribute | write | `attribute_id`, `ticket_type_id`; *allow_multiple_values, description, list_items, multiline, name, required_to_create, required_to_create_for_contacts, visible_on_create and 1 more* |  |  |
