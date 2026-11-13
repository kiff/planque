# Pipedrive: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

403 tools: 206 read, 197 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `PIPEDRIVE_ADD_ACTIVITY_TYPE` | Add new activity type | write | `icon_key`, `name`; *color* |  | `admin` |
| `PIPEDRIVE_ADD_AN_ACTIVITY` | Add an activity | write | `subject`, `type`; *attendees, busy_flag, deal_id, done, due_date, due_time, duration, lead_id and 8 more* |  | `activities:full` |
| `PIPEDRIVE_ADD_AN_ORGANIZATION` | Add an organization | write | `name`; *owner_id, visible_to* |  | `organizations:add` |
| `PIPEDRIVE_ADD_A_DEAL` | Add a deal | write | `title`; *currency, expected_close_date, org_id, person_id, pipeline_id, probability, stage_id, status and 3 more* |  | `deals:add` |
| `PIPEDRIVE_ADD_A_PERSON` | Add a person | write | `name`; *add_time, email, label_ids, marketing_status, org_id, owner_id, phone, update_time and 1 more* |  | `persons:add` |
| `PIPEDRIVE_ADD_CALL_LOG` | Add a call log | write | `end_time`, `outcome`, `start_time`, `to_phone_number`; *activity_id, deal_id, duration, from_phone_number, lead_id, note, org_id, person_id and 2 more* |  | `callLogs:add` |
| `PIPEDRIVE_ADD_CALL_LOG_AUDIO_FILE` | Add Call Log Audio File | write | `file`, `id` |  | `phone-integration` |
| `PIPEDRIVE_ADD_CHANNEL` | Add a channel | write | `name`, `provider_channel_id`; *avatar_url, provider_type, template_support* |  | `messengers-integration` |
| `PIPEDRIVE_ADD_DEALFIELDS_OPTIONS` | Add deal field options | write | `field_code`, `options` |  | `admin` |
| `PIPEDRIVE_ADD_DEAL_DISCOUNT` | Add discount to deal | write | `amount`, `description`, `id`, `type` |  | `deals:read` |
| `PIPEDRIVE_ADD_DEAL_FIELD` | Add a new deal field | write | `field_type`, `name`; *add_visible_flag, options* |  | `dealFields:read`, `dealFields:write` |
| `PIPEDRIVE_ADD_DEAL_FOLLOWER` | Add a follower to a deal | write | `id`, `user_id` |  | `deals:full` |
| `PIPEDRIVE_ADD_DEAL_PARTICIPANT` | Add a participant to a deal | write | `id`, `person_id` |  | `deals:full` |
| `PIPEDRIVE_ADD_DEAL_PRODUCT` | Add product to deal | write | `id`, `item_price`, `product_id`, `quantity`; *billing_frequency, billing_frequency_cycles, billing_start_date, comments, discount, discount_percentage, discount_type, duration and 4 more* |  | `deals:read` |
| `PIPEDRIVE_ADD_DEAL_PRODUCTS` | Bulk add products to deal | write | `data`, `id` |  | `deals:full` |
| `PIPEDRIVE_ADD_FILE` | Add file | write | *activity_id, deal_id, file, lead_id, org_id, person_id, product_id* |  | `deals:full` |
| `PIPEDRIVE_ADD_FILTER` | Add a new filter | write | `conditions`, `name`, `type` |  | `deals:full`, `activities:full`, `contacts:full` |
| `PIPEDRIVE_ADD_GOAL` | Add a new goal | write | `assignee`, `duration`, `expected_outcome`, `interval`, `type`; *title* |  | `goals:full` |
| `PIPEDRIVE_ADD_INCOMING_MESSAGE` | Receives an incoming message | write | `channel_id`, `conversation_id`, `created_at`, `id`, `message`, `sender_id`, `status`; *attachments, conversation_link, reply_by* |  | `messengers-integration` |
| `PIPEDRIVE_ADD_INSTALLMENT_SUBSCRIPTION` | Add an installment subscription | write | `currency`, `deal_id`, `payments`; *update_deal_value* |  | `deals:full` |
| `PIPEDRIVE_ADD_LEAD` | Add a lead | write | `title`; *expected_close_date, label_ids, organization_id, owner_id, person_id, value__amount, value__currency, was_seen* |  | `leads:full` |
| `PIPEDRIVE_ADD_LEAD_LABEL` | Add a lead label | write | `color`, `name` |  | `leadLabels:add` |
| `PIPEDRIVE_ADD_NOTE` | Add Note | write | `content`; *add_time, deal_id, lead_id, org_id, person_id, pinned_to_deal_flag, pinned_to_lead_flag, pinned_to_organization_flag and 4 more* |  | `notes:add` |
| `PIPEDRIVE_ADD_NOTE_COMMENT` | Add comment to note | write | `content`, `id` |  | `deals:full` |
| `PIPEDRIVE_ADD_ORGANIZATION_FIELD` | Add a new organization field | write | `field_type`, `name`; *add_visible_flag, options* |  | `admin` |
| `PIPEDRIVE_ADD_ORGANIZATION_FIELD_OPTIONS` | Add Organization Field Options | write | `field_code`, `options` |  | `admin` |
| `PIPEDRIVE_ADD_ORGANIZATION_FOLLOWER` | Add organization follower | write | `id`, `user_id` |  | `contacts:full` |
| `PIPEDRIVE_ADD_OR_UPDATE_ROLE_SETTING` | Add or update role setting | write | `id`, `setting_key`, `value` |  | `admin` |
| `PIPEDRIVE_ADD_PERSONFIELDS_OPTIONS` | Add person field options | write | `field_code`, `options` |  | `admin` |
| `PIPEDRIVE_ADD_PERSON_FIELD` | Add a new person field | write | `field_type`, `name`; *add_visible_flag, api_token, options* |  | `contact-fields:full`, `admin` |
| `PIPEDRIVE_ADD_PERSON_FOLLOWER` | Add a follower to a person | write | `id`, `user_id` |  | `persons:read`, `persons:write` |
| `PIPEDRIVE_ADD_PERSON_PICTURE` | Add person picture | write | `id`; *crop_height, crop_width, crop_x, crop_y, file* |  | `contacts:full` |
| `PIPEDRIVE_ADD_PIPELINE` | Add a new pipeline | write | `name`; *is_deal_probability_enabled* |  | `deals:read` |
| `PIPEDRIVE_ADD_PRODUCT` | Add a product | write | `name`; *billing_frequency, billing_frequency_cycles, category, code, description, is_linkable, owner_id, prices and 3 more* |  | `products:full` |
| `PIPEDRIVE_ADD_PRODUCTFIELDS_OPTIONS` | Add product field options | write | `field_code`, `options` |  | `products:full` |
| `PIPEDRIVE_ADD_PRODUCT_FIELD` | Add a new product field | write | `field_type`, `name`; *options* |  | `product-fields:full`, `products:full` |
| `PIPEDRIVE_ADD_PRODUCT_FOLLOWER` | Add follower to product | write | `id`, `user_id` |  | `products:full` |
| `PIPEDRIVE_ADD_PRODUCT_IMAGE` | Add Product Image | write | `data`, `id` |  | `products:read` |
| `PIPEDRIVE_ADD_PRODUCT_VARIATION` | Create product variation | write | `id`, `name`; *prices* |  | `products:read` |
| `PIPEDRIVE_ADD_PROJECT` | Add a project | write | `board_id`, `phase_id`, `title`; *deal_ids, description, end_date, labels, org_id, owner_id, person_id, start_date and 2 more* |  | `projects:full` |
| `PIPEDRIVE_ADD_RECURRING_SUBSCRIPTION` | Add a recurring subscription | write | `cadence_type`, `currency`, `cycle_amount`, `deal_id`, `start_date`; *cycles_count, description, infinite, payments, update_deal_value* |  | `deals:full` |
| `PIPEDRIVE_ADD_ROLE` | Add a role | write | `name`; *parent_role_id* |  | `admin` |
| `PIPEDRIVE_ADD_ROLE_ASSIGNMENT` | Add role assignment | write | `id`, `user_id` |  | `admin` |
| `PIPEDRIVE_ADD_STAGE` | Add a new stage | write | `name`, `pipeline_id`; *deal_probability, rotten_days, rotten_flag* |  | `admin` |
| `PIPEDRIVE_ADD_TASK` | Add a task | write | `project_id`, `title`; *assignee_id, description, done, due_date, parent_task_id* |  | `projects:full` |
| `PIPEDRIVE_ADD_TEAM` | Add a new team | write | `manager_id`, `name`; *description, users* |  | `admin` |
| `PIPEDRIVE_ADD_TEAM_USERS` | Add users to a team | write | `id`, `users` |  | `admin` |
| `PIPEDRIVE_ADD_USER` | Add a new user | write | `email`; *access, active_flag* |  | `admin` |
| `PIPEDRIVE_ARCHIVE_PROJECT` | Archive a project | write | `id` |  | `projects:full` |
| `PIPEDRIVE_CANCEL_RECURRING_SUBSCRIPTION` | Cancel a recurring subscription | write | `id`; *end_date* |  | `deals:full` |
| `PIPEDRIVE_CONVERT_DEAL_TO_LEAD` | Convert deal to lead | write | `id` |  | `deals:full` |
| `PIPEDRIVE_CONVERT_LEAD_TO_DEAL` | Convert Lead to Deal | write | `id`; *pipeline_id, stage_id* |  | `leads:full` |
| `PIPEDRIVE_CREATE_ORGANIZATION_RELATIONSHIP` | Create an organization relationship | write | `rel_linked_org_id`, `rel_owner_org_id`, `type`; *org_id* |  | `contacts:full` |
| `PIPEDRIVE_CREATE_REMOTE_FILE_LINK` | Create a remote file and link it to an item | write | `file_type`, `item_id`, `item_type`, `remote_location`, `title` |  | `deals:full` |
| `PIPEDRIVE_CREATE_WEBHOOK` | Create a new webhook | write | `event_action`, `event_object`, `name`, `subscription_url`; *http_auth_password, http_auth_user, user_id, version* |  | `webhooks:full` |
| `PIPEDRIVE_CTIVITYIELDS_GET_ALL_ACTIVITYFIELDS` | Get all ActivityFields | read | *cursor, include_fields, limit* |  |  |
| `PIPEDRIVE_DEALS_UPDATE_DEAL` | Update deal product | write | `id`, `product_attachment_id`; *billing_frequency, billing_frequency_cycles, billing_start_date, comments, discount, discount_type, is_enabled, item_price and 5 more* |  | `deals:full` |
| `PIPEDRIVE_DELETE_ACTIVITIES_BULK` | Delete multiple activities in bulk | write | `ids` |  | `activities:full` |
| `PIPEDRIVE_DELETE_ACTIVITY` | Delete Activity | write | `id` |  | `activities:full` |
| `PIPEDRIVE_DELETE_ACTIVITY_TYPE` | Delete Activity Type | write | `id` |  | `admin` |
| `PIPEDRIVE_DELETE_ACTIVITY_TYPES_BULK` | Delete multiple activity types in bulk | write | `ids` |  | `admin` |
| `PIPEDRIVE_DELETE_AN_ACTIVITY` | Delete an activity | write | `id` |  | `activities:delete` |
| `PIPEDRIVE_DELETE_A_PERSON` | Delete a person | write | `id` |  | `persons:delete` |
| `PIPEDRIVE_DELETE_A_PERSON_FIELD` | Delete a person field | write | `id` |  | `contacts:read` |
| `PIPEDRIVE_DELETE_A_STAGE` | Delete a stage | write | `id` |  | `admin` |
| `PIPEDRIVE_DELETE_CALL_LOG` | Delete a call log | write | `id` |  | `callLogs:delete` |
| `PIPEDRIVE_DELETE_CHANNEL` | Delete a channel | write | `id` |  | `messengers-integration` |
| `PIPEDRIVE_DELETE_COMMENT` | Delete comment from note | write | `commentId`, `id` |  | `deals:full` |
| `PIPEDRIVE_DELETE_CONVERSATION` | Delete a conversation | write | `channel__id`, `conversation__id` |  | `messengers-integration` |
| `PIPEDRIVE_DELETE_DEAL` | Delete Deal | write | `id` |  | `deals:full` |
| `PIPEDRIVE_DELETE_DEALS_BULK` | Delete multiple deals in bulk | write | `ids` |  | `deals:delete` |
| `PIPEDRIVE_DELETE_DEAL_DISCOUNT` | Delete Deal Discount | write | `discount_id`, `id` |  | `deals:read` |
| `PIPEDRIVE_DELETE_DEAL_FIELD` | Delete a deal field | write | `id` |  | `admin` |
| `PIPEDRIVE_DELETE_DEAL_FIELDS_BULK` | Delete multiple deal fields in bulk | write | `ids` |  | `deal-fields:full`, `admin` |
| `PIPEDRIVE_DELETE_DEAL_FIELD_OPTIONS` | Delete deal field options | write | `field_code`, `options` |  | `admin` |
| `PIPEDRIVE_DELETE_DEAL_FOLLOWER` | Delete deal follower | write | `follower_id`, `id` |  | `deals:read` |
| `PIPEDRIVE_DELETE_DEAL_PARTICIPANT` | Delete a participant from a deal | write | `deal_participant_id`, `id` |  | `deals:delete` |
| `PIPEDRIVE_DELETE_DEAL_PRODUCT` | Delete Deal Product | write | `id`, `product_attachment_id` |  | `deals:full` |
| `PIPEDRIVE_DELETE_DEAL_PRODUCTS` | Delete Deals Products | write | `id`; *ids* |  | `deals:full` |
| `PIPEDRIVE_DELETE_FILE` | Delete a file | write | `id` |  | `deals:full` |
| `PIPEDRIVE_DELETE_FILTER` | Delete Pipedrive Filter | write | `id` |  | `deals:full`, `activities:full`, `contacts:full` |
| `PIPEDRIVE_DELETE_FILTERS_BULK` | Delete multiple filters in bulk | write | `ids` |  | `deals:full`, `activities:full`, `contacts:full` |
| `PIPEDRIVE_DELETE_GOAL` | Delete existing goal | write | `id` |  | `goals:full` |
| `PIPEDRIVE_DELETE_LEAD` | Delete a lead | write | `id` |  | `leads:full` |
| `PIPEDRIVE_DELETE_LEAD_LABEL` | Delete lead label | write | `id` |  | `leads:full` |
| `PIPEDRIVE_DELETE_MAIL_THREAD` | Delete mail thread | write | `id` |  | `mail:full` |
| `PIPEDRIVE_DELETE_MULTIPLE_PRODUCT_FIELDS` | Delete multiple product fields in bulk | write | `ids` |  | `product-fields:full`, `products:full` |
| `PIPEDRIVE_DELETE_MULTIPLE_STAGES` | Delete multiple stages in bulk | write | `ids` |  | `admin` |
| `PIPEDRIVE_DELETE_NOTE` | Delete a note | write | `id` |  | `deals:full`, `contacts:full` |
| `PIPEDRIVE_DELETE_ORGANIZATION` | Delete an organization | write | `id` |  | `contacts:full` |
| `PIPEDRIVE_DELETE_ORGANIZATIONS_BULK` | Delete multiple organizations in bulk | write | `ids` |  | `api_key`, `oauth2` |
| `PIPEDRIVE_DELETE_ORGANIZATION_FIELD` | Delete an organization field | write | `id` |  | `organizationFields:read`, `organizationFields:write` |
| `PIPEDRIVE_DELETE_ORGANIZATION_FIELDS_BULK` | Delete multiple organization fields in bulk | write | `ids` |  | `admin` |
| `PIPEDRIVE_DELETE_ORGANIZATION_FIELD_OPTIONS` | Delete organization field options | write | `field_code`, `options` |  | `admin` |
| `PIPEDRIVE_DELETE_ORGANIZATION_FOLLOWER` | Delete organization follower | write | `follower_id`, `id` |  | `contacts:read` |
| `PIPEDRIVE_DELETE_ORGANIZATION_RELATIONSHIP` | Delete an organization relationship | write | `id` |  | `contacts:full` |
| `PIPEDRIVE_DELETE_PERSON` | Delete a person | write | `id` |  | `contacts:full` |
| `PIPEDRIVE_DELETE_PERSONS_BULK` | Delete multiple persons in bulk | write | `ids` |  | `persons:delete` |
| `PIPEDRIVE_DELETE_PERSON_FIELD` | Delete a person field | write | `id` |  |  |
| `PIPEDRIVE_DELETE_PERSON_FIELDS_BULK` | Delete multiple person fields in bulk | write | `ids`; *api_token* |  | `admin` |
| `PIPEDRIVE_DELETE_PERSON_FIELD_OPTIONS` | Delete person field options | write | `field_code`, `options` |  | `admin` |
| `PIPEDRIVE_DELETE_PERSON_FOLLOWER` | Delete a follower from a person | write | `follower_id`, `id` |  | `persons:delete` |
| `PIPEDRIVE_DELETE_PERSON_PICTURE` | Delete person picture | write | `id` |  | `contacts:full` |
| `PIPEDRIVE_DELETE_PIPELINE` | Delete a pipeline | write | `id` |  | `admin` |
| `PIPEDRIVE_DELETE_PRODUCT` | Delete a product | write | `id` |  | `products:delete` |
| `PIPEDRIVE_DELETE_PRODUCT_FIELD` | Delete a product field | write | `id` |  | `productFields:delete` |
| `PIPEDRIVE_DELETE_PRODUCT_FIELD_OPTIONS` | Delete product field options | write | `field_code`, `options` |  | `products:full` |
| `PIPEDRIVE_DELETE_PRODUCT_FOLLOWER` | Delete a follower from a product | write | `follower_id`, `id` |  | `products:full` |
| `PIPEDRIVE_DELETE_PRODUCT_IMAGE` | Delete Product Image | write | `id`, `image_id` |  | `products:full` |
| `PIPEDRIVE_DELETE_PRODUCT_VARIATION` | Delete Product Variation | write | `id`, `product_variation_id` |  | `products:full` |
| `PIPEDRIVE_DELETE_PROJECT` | Delete a project | write | `id` |  | `projects:full` |
| `PIPEDRIVE_DELETE_ROLE` | Delete a role | write | `id` |  | `admin` |
| `PIPEDRIVE_DELETE_ROLE_ASSIGNMENT` | Delete a role assignment | write | `id`, `user_id` |  | `admin` |
| `PIPEDRIVE_DELETE_STAGE` | Delete a stage | write | `id` |  | `admin` |
| `PIPEDRIVE_DELETE_SUBSCRIPTION` | Delete a subscription | write | `id` |  | `deals:full` |
| `PIPEDRIVE_DELETE_TASK` | Delete a task | write | `id` |  | `projects:full` |
| `PIPEDRIVE_DELETE_TEAM_USERS` | Delete users from a team | write | `id`, `users` |  | `admin` |
| `PIPEDRIVE_DELETE_WEBHOOK` | Delete existing webhook | write | `id` |  | `admin` |
| `PIPEDRIVE_DOWNLOAD_FILE` | Download one file | read | `id` |  | `deals:read` |
| `PIPEDRIVE_DUPLICATE_DEAL` | Duplicate deal | write | `id` |  | `deals:full` |
| `PIPEDRIVE_DUPLICATE_PRODUCT` | Duplicate Product | write | `id` |  | `products:full` |
| `PIPEDRIVE_ENUMERATE_ACCESSIBLE_USERS_FOR_LEAD` | Enumerate accessible users for lead | read | `id` |  | `leads:read` |
| `PIPEDRIVE_FIND_GOALS` | Find goals | read | *assignee__id, assignee__type, expected__outcome__currency__id, expected__outcome__target, expected__outcome__tracking__metric, is_active, period__end, period__start and 5 more* |  | `goals:full` |
| `PIPEDRIVE_FIND_USERS_BY_NAME` | Find users by name | read | `term`; *search_by_email* |  | `users:read` |
| `PIPEDRIVE_GET_ACTIVITY` | Get details of an activity | read | `id` |  | `activities:read` |
| `PIPEDRIVE_GET_ACTIVITY_FIELD` | Get one activity field | read | `field_code`; *include_fields* |  | `activities:read`, `activities:full`, `admin` |
| `PIPEDRIVE_GET_ALL_ACTIVITIES_ASSIGNED_TO_A_PARTICULAR_USER` | Get all activities assigned to a particular user | read | *done, end_date, fields, filter_id, limit, response_detail, start, start_date and 2 more* |  | `activities:read` |
| `PIPEDRIVE_GET_ALL_ACTIVITIES_BETA` | Get all activities beta | read | *cursor, done, limit, since, type, until, user_id* |  | `activities:read` |
| `PIPEDRIVE_GET_ALL_ACTIVITY_FIELDS` | Get all activity fields | read |  |  | `activities:read` |
| `PIPEDRIVE_GET_ALL_ACTIVITY_TYPES` | Get all activity types | read |  |  | `activityTypes:read` |
| `PIPEDRIVE_GET_ALL_DEALS` | Get all deals | read | *filter_id, limit, owned_by_you, sort, stage_id, start, status, user_id* |  | `deals:read` |
| `PIPEDRIVE_GET_ALL_DEALS_BETA` | Get all deals beta | read | *cursor, limit, since, stage_id, status, until, user_id* |  | `deals:read` |
| `PIPEDRIVE_GET_ALL_DEAL_FIELDS` | Get all deal fields | read | *limit, start* |  | `deals:read` |
| `PIPEDRIVE_GET_ALL_FILES` | Get all files | read | *limit, sort, start* |  | `files:read` |
| `PIPEDRIVE_GET_ALL_FILTERS` | Get all filters | read | *type* |  | `filters:read` |
| `PIPEDRIVE_GET_ALL_FILTER_HELPERS` | Get all filter helpers | read |  |  | `filters:read` |
| `PIPEDRIVE_GET_ALL_LEADS` | Get all leads | read | *archived_status, filter_id, limit, organization_id, owner_id, person_id, sort, start* |  | `leads:read` |
| `PIPEDRIVE_GET_ALL_LEAD_LABELS` | Get all lead labels | read |  |  | `leadLabels:read` |
| `PIPEDRIVE_GET_ALL_LEAD_SOURCES` | Get all lead sources | read |  |  | `leads:read` |
| `PIPEDRIVE_GET_ALL_NOTES` | Get all notes | read | *deal_id, end_date, lead_id, limit, org_id, person_id, pinned_to_deal_flag, pinned_to_lead_flag and 6 more* |  | `notes:read` |
| `PIPEDRIVE_GET_ALL_NOTE_FIELDS` | Get all note fields | read |  |  | `noteFields:read` |
| `PIPEDRIVE_GET_ALL_ORGANIZATIONS` | Get all organizations | read | *filter_id, first_char, limit, sort, start, user_id* |  | `contacts:read`, `contacts:full` |
| `PIPEDRIVE_GET_ALL_ORGANIZATIONS_BETA` | Get all organizations beta | read | *cursor, first_char, limit, owner_id, since, until* |  | `contacts:read`, `contacts:full` |
| `PIPEDRIVE_GET_ALL_ORGANIZATION_FIELDS` | Get all organization fields | read | *limit, start* |  | `organizationFields:read` |
| `PIPEDRIVE_GET_ALL_PERMISSION_SETS` | Get all permission sets | read | *app* |  | `admin` |
| `PIPEDRIVE_GET_ALL_PERSONS` | Get all persons | read | *filter_id, first_char, limit, sort, start, user_id* |  | `contacts:read` |
| `PIPEDRIVE_GET_ALL_PERSONS_BETA` | Get all persons beta | read | *cursor, first_char, limit, owner_id, since, until* |  | `contacts:read` |
| `PIPEDRIVE_GET_ALL_PERSON_FIELDS` | Get all person fields | read | *limit, start* |  | `api_key`, `oauth2` |
| `PIPEDRIVE_GET_ALL_PIPELINES` | Get all pipelines | read |  |  | `deals:read`, `deals:full`, `admin` |
| `PIPEDRIVE_GET_ALL_PRODUCTS` | Get all products | read | *filter_id, first_char, get_summary, ids, limit, start, user_id* |  | `products:read` |
| `PIPEDRIVE_GET_ALL_PRODUCT_FIELDS` | Get all product fields | read | *limit, start* |  | `products:read` |
| `PIPEDRIVE_GET_ALL_PROJECTS` | Get all projects | read | *cursor, filter_id, include_archived, limit, phase_id, status* |  | `projects:read` |
| `PIPEDRIVE_GET_ALL_PROJECT_BOARDS` | Get all project boards | read |  |  | `projects:read` |
| `PIPEDRIVE_GET_ALL_PROJECT_TEMPLATES` | Get all project templates | read | *cursor, limit* |  | `projects:read` |
| `PIPEDRIVE_GET_ALL_RELATIONSHIPS_FOR_ORGANIZATION` | Get all relationships for organization | read | `org_id` |  | `contacts:read`, `contacts:full` |
| `PIPEDRIVE_GET_ALL_ROLES` | Get all roles | read | *limit, start* |  | `admin` |
| `PIPEDRIVE_GET_ALL_STAGES` | Get all stages | read | *limit, pipeline_id, start* |  | `deals:read`, `deals:full`, `admin` |
| `PIPEDRIVE_GET_ALL_SUPPORTED_CURRENCIES` | Get all supported currencies | read | *term* |  | `base` |
| `PIPEDRIVE_GET_ALL_TASKS` | Get all tasks | read | *assignee_id, cursor, done, limit, parent_task_id, project_id* |  | `projects:read` |
| `PIPEDRIVE_GET_ALL_TEAMS` | Get all teams | read | *order_by, skip_users* |  | `teams:read` |
| `PIPEDRIVE_GET_ALL_TEAMS_OF_A_USER` | Get all teams of a user | read | `id`; *order_by, skip_users* |  | `users:read` |
| `PIPEDRIVE_GET_ALL_USERS` | Get all users | read |  |  | `users:read` |
| `PIPEDRIVE_GET_ALL_USERS_IN_A_TEAM` | Get all users in a team | read | `id` |  | `users:read` |
| `PIPEDRIVE_GET_ALL_USER_CONNECTIONS` | Get all user connections | read |  |  | `base` |
| `PIPEDRIVE_GET_ALL_WEBHOOKS` | Get all webhooks | read |  |  | `webhooks:read` |
| `PIPEDRIVE_GET_ARCHIVED_DEALS` | Get archived deals | read | *cursor, custom_fields, filter_id, ids, include_fields, limit, org_id, owner_id and 8 more* |  | `deals:read`, `deals:full` |
| `PIPEDRIVE_GET_ARCHIVED_LEADS` | Get archived leads | read | *filter_id, limit, organization_id, owner_id, person_id, sort, start* |  | `leads:read`, `leads:full` |
| `PIPEDRIVE_GET_CALL_LOG` | Get call log details | read | `id` |  | `phone-integration` |
| `PIPEDRIVE_GET_COMMENT` | Get one comment | read | `commentId`, `id` |  | `deals:read`, `deals:full`, `contacts:read`, `contacts:full` |
| `PIPEDRIVE_GET_COMPANY_ADD_ONS` | Get all add ons for a single company | read |  |  | `users:read` |
| `PIPEDRIVE_GET_CONVERSATIONS` | Get conversations for a channel | read | `channel_id`; *after, conversations_limit, messages_limit* |  | `base` |
| `PIPEDRIVE_GET_CURRENT_USER_DATA` | Get current user data | read |  |  | `users:read` |
| `PIPEDRIVE_GET_DEAL` | Get Deal Details | read | `id`; *custom_fields, include_fields* |  | `deals:full` |
| `PIPEDRIVE_GET_DEALS_CONVERSION_RATES_IN_PIPELINE` | Get deals conversion rates in pipeline | read | `end_date`, `id`, `start_date`; *user_id* |  | `deals:read` |
| `PIPEDRIVE_GET_DEALS_CONVERT_STATUS` | Get deal conversion status | read | `conversion_id`, `id` |  | `deals:read` |
| `PIPEDRIVE_GET_DEALS_DISCOUNTS` | Get deals discounts | read | `id` |  | `deals:read` |
| `PIPEDRIVE_GET_DEALS_IN_A_PIPELINE` | Get deals in a pipeline | read | `id`; *everyone, filter_id, get_summary, limit, stage_id, start, totals_convert_currency, user_id* |  | `deals:read` |
| `PIPEDRIVE_GET_DEALS_IN_A_STAGE` | Get deals in a stage | read | `id`; *everyone, filter_id, limit, start, user_id* |  | `deals:read` |
| `PIPEDRIVE_GET_DEALS_PRODUCTS` | Get Deals Products | read | `deal_ids`; *cursor, limit, sort_by, sort_direction* |  | `deals:read` |
| `PIPEDRIVE_GET_DEALS_SUMMARY` | Get deals summary | read | *filter_id, stage_id, status, user_id* |  | `deals:read`, `deals:full` |
| `PIPEDRIVE_GET_DEALS_SUMMARY_ARCHIVED` | Get archived deals summary | read | *filter_id, pipeline_id, stage_id, status, user_id* |  | `deals:read`, `deals:full` |
| `PIPEDRIVE_GET_DEALS_TIMELINE` | Get deals timeline | read | `amount`, `field_key`, `interval`, `start_date`; *exclude_deals, filter_id, pipeline_id, totals_convert_currency, user_id* |  | `deals:read` |
| `PIPEDRIVE_GET_DEALS_WITH_PRODUCT` | Get deals where a product is attached to | read | `id`; *limit, start, status* |  | `deals:read` |
| `PIPEDRIVE_GET_DEAL_CHANGELOG` | Get Deal Changelog | read | `id`; *cursor, limit* |  | `deals:read`, `deals:full` |
| `PIPEDRIVE_GET_DEAL_FIELD` | Get one deal field | read | `id` |  | `deals:read` |
| `PIPEDRIVE_GET_DEAL_FILES` | Get Deal Files | read | `id`; *limit, sort, start* |  | `deals:read`, `deals:full` |
| `PIPEDRIVE_GET_DEAL_FOLLOWERS` | Get deal followers | read | `id`; *cursor, limit* |  | `deals:read`, `deals:full` |
| `PIPEDRIVE_GET_DEAL_FOLLOWERS_CHANGELOG` | Get deal followers changelog | read | `id`; *cursor, limit* |  | `deals:read` |
| `PIPEDRIVE_GET_DEAL_MAIL_MESSAGES` | Get Deal Mail Messages | read | `id`; *limit, start* |  | `mail:read`, `deals:read` |
| `PIPEDRIVE_GET_DEAL_PARTICIPANTS` | Get Deal Participants | read | `id`; *limit, start* |  | `deals:read`, `deals:full` |
| `PIPEDRIVE_GET_DEAL_PARTICIPANTS_CHANGELOG` | Get Deal Participants Changelog | read | `id`; *cursor, limit* |  | `deals:read`, `deals:full` |
| `PIPEDRIVE_GET_DEAL_PRODUCTS` | Get Deal Products | read | `id`; *cursor, limit, sort_by, sort_direction* |  | `products:read`, `products:full` |
| `PIPEDRIVE_GET_DEAL_SUBSCRIPTION` | Find subscription by deal | read | `dealId` |  | `deals:read` |
| `PIPEDRIVE_GET_DEAL_UPDATES` | Get deal updates | read | `id`; *all_changes, items, limit, start* |  | `recents:read` |
| `PIPEDRIVE_GET_DEAL_USERS` | Get Deal Permitted Users | read | `id` |  | `deals:read`, `deals:full` |
| `PIPEDRIVE_GET_DETAILS_OF_AN_ORGANIZATION` | Get details of an organization | read | `id` |  | `contacts:read` |
| `PIPEDRIVE_GET_FILE` | Get one file | read | `id` |  | `files:read` |
| `PIPEDRIVE_GET_FILTER` | Get one filter | read | `id` |  | `deals:read`, `deals:full`, `activities:read`, `activities:full`, `contacts:read`, `contacts:full` |
| `PIPEDRIVE_GET_GOAL_RESULT` | Get result of a goal | read | `id`, `period__end`, `period__start` |  | `goals:read`, `goals:full` |
| `PIPEDRIVE_GET_LEAD_CONVERSION_STATUS` | Get lead conversion status by ID | read | `conversion_id`, `id` |  | `leads:read` |
| `PIPEDRIVE_GET_LEAD_FIELDS` | Get all lead fields | read | *limit, start* |  | `leads:read` |
| `PIPEDRIVE_GET_LEAD_USERS` | Get Lead Permitted Users | read | `id` |  | `leads:read`, `leads:full` |
| `PIPEDRIVE_GET_MAIL_MESSAGE` | Get one mail message | read | `id`; *include_body* |  | `mail:read` |
| `PIPEDRIVE_GET_MAIL_THREAD` | Get one mail thread | read | `id` |  | `mailbox:read` |
| `PIPEDRIVE_GET_MAIL_THREADS` | Get mail threads | read | `folder`; *limit, start* |  | `mail:read`, `mail:full` |
| `PIPEDRIVE_GET_MAIL_THREAD_MESSAGES` | Get all mail messages of mail thread | read | `id` |  | `mailbox:read` |
| `PIPEDRIVE_GET_MARKETPLACE_CLIENT_ID` | Get Marketplace Client ID | read | *marketplace_client_id_hint, provider* |  |  |
| `PIPEDRIVE_GET_NOTE` | Get one note | read | `id` |  | `notes:read` |
| `PIPEDRIVE_GET_NOTE_COMMENTS` | Get all comments for a note | read | `id`; *limit, start* |  | `deals:full`, `deals:read`, `contacts:full`, `contacts:read` |
| `PIPEDRIVE_GET_ONE_LEAD` | Get one lead | read | `id` |  | `leads:read`, `leads:full` |
| `PIPEDRIVE_GET_ORGANIZATION` | Get Organization Details | read | `id`; *custom_fields, include_fields* |  | `contacts:read`, `contacts:full` |
| `PIPEDRIVE_GET_ORGANIZATION_CHANGELOG` | Get Organization Changelog | read | `id`; *cursor, limit* |  | `contacts:read` |
| `PIPEDRIVE_GET_ORGANIZATION_DEALS` | Get organization deals | read | `id`; *limit, only_primary_association, sort, start, status* |  | `deals:read`, `deals:full` |
| `PIPEDRIVE_GET_ORGANIZATION_FIELD` | Get one organization field | read | `id` |  | `contacts:read` |
| `PIPEDRIVE_GET_ORGANIZATION_FILES` | Get Organization Files | read | `id`; *limit, sort, start* |  | `recents:read` |
| `PIPEDRIVE_GET_ORGANIZATION_FOLLOWERS` | Get Organization Followers | read | `id`; *cursor, limit* |  | `contacts:full` |
| `PIPEDRIVE_GET_ORGANIZATION_FOLLOWERS_CHANGELOG` | Get Organization Followers Changelog | read | `id`; *cursor, limit* |  | `contacts:read`, `contacts:full` |
| `PIPEDRIVE_GET_ORGANIZATION_MAIL_MESSAGES` | Get Organization Mail Messages | read | `id`; *limit, start* |  | `mail:read`, `mail:full` |
| `PIPEDRIVE_GET_ORGANIZATION_PERSONS` | Get Organization Persons | read | `id`; *limit, start* |  | `contacts:read`, `contacts:full` |
| `PIPEDRIVE_GET_ORGANIZATION_RELATIONSHIP` | Get one organization relationship | read | `id`; *org_id* |  | `contacts:read`, `contacts:full` |
| `PIPEDRIVE_GET_ORGANIZATION_UPDATES` | Get Organization Updates | read | `id`; *all_changes, items, limit, start* |  | `contacts:read` |
| `PIPEDRIVE_GET_ORGANIZATION_USERS` | Get organization permitted users | read | `id` |  | `contacts:read`, `contacts:full` |
| `PIPEDRIVE_GET_PERMISSION_SET` | Get one permission set | read | `id` |  | `admin` |
| `PIPEDRIVE_GET_PERMISSION_SET_ASSIGNMENTS` | Get Permission Set Assignments | read | `id`; *limit, start* |  | `admin` |
| `PIPEDRIVE_GET_PERSON` | Get Person Details | read | `id`; *custom_fields, include_fields* |  | `contacts:read`, `contacts:full` |
| `PIPEDRIVE_GET_PERSON_CHANGELOG` | Get Person Changelog | read | `id`; *cursor, limit* |  | `contacts:read` |
| `PIPEDRIVE_GET_PERSON_DEALS` | Get person deals | read | `id`; *limit, sort, start, status* |  | `deals:read`, `deals:full` |
| `PIPEDRIVE_GET_PERSON_FIELD` | Get one person field | read | `id` |  | `admin` |
| `PIPEDRIVE_GET_PERSON_FILES` | Get Person Files | read | `id`; *limit, sort, start* |  | `recents:read` |
| `PIPEDRIVE_GET_PERSON_FOLLOWERS` | Get Person Followers | read | `id`; *cursor, limit* |  | `contacts:read`, `contacts:full` |
| `PIPEDRIVE_GET_PERSON_FOLLOWERS_CHANGELOG` | Get person followers changelog | read | `id`; *cursor, limit* |  | `contacts:read` |
| `PIPEDRIVE_GET_PERSON_MAIL_MESSAGES` | Get Person Mail Messages | read | `id`; *limit, start* |  | `mail:read`, `mail:full` |
| `PIPEDRIVE_GET_PERSON_PICTURE` | Get Person Picture | read | `id` |  | `contacts:read`, `contacts:full` |
| `PIPEDRIVE_GET_PERSON_PRODUCTS` | Get Person Products | read | `id`; *limit, start* |  | `contacts:full` |
| `PIPEDRIVE_GET_PERSON_UPDATES` | Get person updates | read | `id`; *all_changes, items, limit, start* |  | `recents:read` |
| `PIPEDRIVE_GET_PERSON_USERS` | Get Person Permitted Users | read | `id` |  | `contacts:read`, `contacts:full` |
| `PIPEDRIVE_GET_PIPELINE` | Get one pipeline | read | `id`; *totals_convert_currency* |  | `deals:read`, `deals:full`, `admin` |
| `PIPEDRIVE_GET_PIPELINE_CONVERSION_STATISTICS` | Get Pipeline Conversion Statistics | read | `end_date`, `id`, `start_date`; *user_id* |  | `deals:read`, `deals:full` |
| `PIPEDRIVE_GET_PIPELINE_MOVEMENT_STATISTICS` | Get Pipeline Movement Statistics | read | `end_date`, `id`, `start_date`; *user_id* |  | `deals:read`, `deals:full` |
| `PIPEDRIVE_GET_PRODUCT` | Get one product | read | `id` |  | `products:read` |
| `PIPEDRIVE_GET_PRODUCT_DEALS` | Get Product Deals | read | `id`; *limit, start, status* |  | `products:read`, `products:full` |
| `PIPEDRIVE_GET_PRODUCT_FIELD` | Get one product field | read | `id` |  | `products:read`, `product-fields:full`, `products:full` |
| `PIPEDRIVE_GET_PRODUCT_FILES` | Get Product Files | read | `id`; *limit, sort, start* |  | `deals:read` |
| `PIPEDRIVE_GET_PRODUCT_FOLLOWERS` | Get product followers | read | `id`; *limit, start* |  | `products:read`, `products:full` |
| `PIPEDRIVE_GET_PRODUCT_FOLLOWERS_CHANGELOG` | Get Product Followers Changelog | read | `id`; *cursor, limit* |  | `products:read`, `products:full` |
| `PIPEDRIVE_GET_PRODUCT_IMAGES` | Get product images | read | `id` |  | `products:read` |
| `PIPEDRIVE_GET_PRODUCT_USERS` | List permitted users for product | read | `id` |  | `products:read`, `products:full` |
| `PIPEDRIVE_GET_PRODUCT_VARIATIONS` | Get product variations | read | `id`; *cursor, limit* |  | `products:read` |
| `PIPEDRIVE_GET_PROJECT` | Get details of a project | read | `id` |  | `projects:read` |
| `PIPEDRIVE_GET_PROJECT_ACTIVITIES` | Returns project activities | read | `id` |  | `projects:read` |
| `PIPEDRIVE_GET_PROJECT_BOARD` | Get details of a board | read | `id` |  | `projects:read` |
| `PIPEDRIVE_GET_PROJECT_GROUPS` | Returns project groups | read | `id` |  | `projects:full` |
| `PIPEDRIVE_GET_PROJECT_PHASE` | Get details of a phase | read | `id` |  | `projects:read` |
| `PIPEDRIVE_GET_PROJECT_PHASES` | Get project phases | read | `board_id` |  | `projects:read` |
| `PIPEDRIVE_GET_PROJECT_PLAN` | Returns project plan | read | `id` |  | `projects:full` |
| `PIPEDRIVE_GET_PROJECT_TASKS` | Returns project tasks | read | `id` |  | `projects:read` |
| `PIPEDRIVE_GET_RECENTS` | Get recents | read | `since_timestamp`; *items, limit, start* |  | `recents:read` |
| `PIPEDRIVE_GET_ROLE` | Get one role | read | `id` |  | `admin` |
| `PIPEDRIVE_GET_ROLE_ASSIGNMENTS` | Get role assignments | read | `id`; *limit, start* |  | `admin` |
| `PIPEDRIVE_GET_ROLE_PIPELINES` | Get role pipelines | read | `id`; *visible* |  | `admin` |
| `PIPEDRIVE_GET_ROLE_SETTINGS` | Get role settings | read | `id` |  | `admin` |
| `PIPEDRIVE_GET_STAGE` | Get one stage | read | `id`; *everyone* |  | `stages:read` |
| `PIPEDRIVE_GET_SUBSCRIPTION` | Get details of a subscription | read | `id` |  | `deals:read` |
| `PIPEDRIVE_GET_SUBSCRIPTION_PAYMENTS` | Get all payments of a subscription | read | `id` |  | `deals:read`, `deals:full` |
| `PIPEDRIVE_GET_TASK` | Get details of a task | read | `id` |  | `projects:read` |
| `PIPEDRIVE_GET_TEAM` | Get a single team | read | `id`; *skip_users* |  | `teams:read` |
| `PIPEDRIVE_GET_TEMPLATE` | Get details of a template | read | `id` |  | `projects:read`, `projects:full` |
| `PIPEDRIVE_GET_TEMPLATES` | Get templates for a channel | read | `channel_id` |  | `base` |
| `PIPEDRIVE_GET_USER` | Get one user | read | `id` |  | `base` |
| `PIPEDRIVE_GET_USER_CALL_LOGS` | Get all call logs assigned to a particular user | read | *limit, start* |  | `callLogs:read` |
| `PIPEDRIVE_GET_USER_FOLLOWERS` | Get User Followers | read | `id`; *cursor, limit* |  | `users:read` |
| `PIPEDRIVE_GET_USER_PERMISSIONS` | List user permissions | read | `id` |  | `users:read` |
| `PIPEDRIVE_GET_USER_ROLE_ASSIGNMENTS` | Get User Role Assignments | read | `id`; *limit, start* |  | `users:read` |
| `PIPEDRIVE_GET_USER_ROLE_SETTINGS` | Get user role settings | read | `id` |  | `users:read` |
| `PIPEDRIVE_LINK_REMOTE_FILE` | Link a remote file to an item | write | `item_id`, `item_type`, `remote_id`, `remote_location` |  | `deals:full`, `activities:full`, `contacts:full` |
| `PIPEDRIVE_LINK_USER_VIDEO_INTEGRATION` | Link User Video Integration | write | `company_id`, `marketplace_client_id`, `user_id`, `user_provider_id` |  | `video-calls` |
| `PIPEDRIVE_LIST_ACTIVITIES` | List activities | read | *cursor, deal_id, done, filter_id, ids, include_fields, lead_id, limit and 7 more* |  | `activities:read` |
| `PIPEDRIVE_LIST_DEALS` | List deals | read | *cursor, custom_fields, filter_id, ids, include_fields, include_labels, include_option_labels, limit and 10 more* |  | `deals:read` |
| `PIPEDRIVE_LIST_DEAL_ACTIVITIES` | List activities associated with a deal | read | `id`; *done, exclude, limit, start* |  | `activities:read` |
| `PIPEDRIVE_LIST_DEAL_FILES` | List files attached to a deal | read | `id`; *limit, sort, start* |  | `deals:read`, `deals:full` |
| `PIPEDRIVE_LIST_DEAL_FOLLOWERS` | List followers of a deal | read | `id` |  | `deals:read` |
| `PIPEDRIVE_LIST_DEAL_MAIL_MESSAGES` | List mail messages associated with a deal | read | `id`; *limit, start* |  | `mail:read` |
| `PIPEDRIVE_LIST_DEAL_PERSONS` | List all persons associated with a deal | read | `id`; *limit, start* |  | `deals:read`, `persons:read` |
| `PIPEDRIVE_LIST_DEAL_PRODUCTS` | List products attached to a deal | read | `id`; *include_product_data, limit, start* |  | `products:read`, `products:full`, `deals:read`, `deals:full` |
| `PIPEDRIVE_LIST_DEAL_UPDATES` | List updates about a deal | read | `id`; *all_changes, items, limit, start* |  | `recents:read` |
| `PIPEDRIVE_LIST_DEAL_USERS` | List permitted users | read | `id` |  | `deals:read` |
| `PIPEDRIVE_LIST_FILES` | List all files | read | *limit, sort, start* |  | `deals:read`, `deals:full`, `activities:read`, `activities:full`, `contacts:read`, `contacts:full` |
| `PIPEDRIVE_LIST_ORGANIZATION_ACTIVITIES` | List activities associated with an organization | read | `id`; *done, exclude, limit, start* |  | `activities:read`, `activities:full` |
| `PIPEDRIVE_LIST_ORGANIZATION_DEALS` | List deals associated with an organization | read | `id`; *limit, only_primary_association, sort, start, status* |  | `deals:read` |
| `PIPEDRIVE_LIST_ORGANIZATION_FILES` | List files attached to an organization | read | `id`; *limit, sort, start* |  | `contacts:read` |
| `PIPEDRIVE_LIST_ORGANIZATION_FOLLOWERS` | List followers of an organization | read | `id` |  | `contacts:read`, `contacts:full` |
| `PIPEDRIVE_LIST_ORGANIZATION_MAIL_MESSAGES` | List mail messages associated with an organization | read | `id`; *limit, start* |  | `mail:read` |
| `PIPEDRIVE_LIST_ORGANIZATION_PERSONS` | List persons of an organization | read | `id`; *limit, start* |  | `organizations:read`, `persons:read` |
| `PIPEDRIVE_LIST_PARTICIPANTS_OF_A_DEAL` | List participants of a deal | read | `id`; *limit, start* |  | `deals:read` |
| `PIPEDRIVE_LIST_PERMISSION_SET_ASSIGNMENTS` | List permission set assignments | read | `id`; *limit, start* |  | `admin` |
| `PIPEDRIVE_LIST_PERSONS` | List persons | read | *cursor, custom_fields, deal_id, filter_id, ids, include_fields, include_labels, include_option_labels and 7 more* |  | `contacts:read` |
| `PIPEDRIVE_LIST_PERSON_ACTIVITIES` | List activities associated with a person | read | `id`; *done, exclude, limit, start* |  | `activities:read` |
| `PIPEDRIVE_LIST_PERSON_DEALS` | List deals associated with a person | read | `id`; *limit, sort, start, status* |  | `deals:read` |
| `PIPEDRIVE_LIST_PERSON_FILES` | List files attached to a person | read | `id`; *limit, sort, start* |  | `contacts:read`, `contacts:full` |
| `PIPEDRIVE_LIST_PERSON_FOLLOWERS` | List followers of a person | read | `id` |  | `contacts:read` |
| `PIPEDRIVE_LIST_PERSON_MAIL_MESSAGES` | List mail messages associated with a person | read | `id`; *limit, start* |  | `mail:read` |
| `PIPEDRIVE_LIST_PERSON_PRODUCTS` | List products associated with a person | read | `id`; *limit, start* |  | `contacts:read` |
| `PIPEDRIVE_LIST_PERSON_USERS` | List Person Access Users | read | `id` |  | `contacts:read`, `contacts:full` |
| `PIPEDRIVE_LIST_PRODUCT_FILES` | List files attached to a product | read | `id`; *limit, sort, start* |  | `products:read`, `products:full` |
| `PIPEDRIVE_LIST_PRODUCT_FOLLOWERS` | List followers of a product | read | `id`; *limit, start* |  | `products:read` |
| `PIPEDRIVE_LIST_PRODUCT_USERS` | List Product Permitted Users | read | `id` |  | `products:read`, `products:full` |
| `PIPEDRIVE_LIST_ROLE_ASSIGNMENTS` | List role assignments | read | `id`; *limit, start* |  | `admin` |
| `PIPEDRIVE_LIST_ROLE_PIPELINE_VISIBILITY` | List pipeline visibility for a role | read | `id`; *visible* |  | `admin` |
| `PIPEDRIVE_LIST_ROLE_SETTINGS` | List role settings | read | `id` |  | `admin` |
| `PIPEDRIVE_LIST_UPDATES_ABOUT_AN_ORGANIZATION` | List updates about an organization | read | `id`; *all_changes, items, limit, start* |  | `organizations:read` |
| `PIPEDRIVE_LIST_UPDATES_ABOUT_A_PERSON` | List updates about a person | read | `id`; *all_changes, items, limit, start* |  | `recents:read` |
| `PIPEDRIVE_LIST_UPDATES_ABOUT_PARTICIPANTS_OF_A_DEAL` | List updates about participants of a deal | read | `id`; *cursor, limit* |  | `deals:read` |
| `PIPEDRIVE_LIST_USER_FOLLOWERS` | List followers of a user | read | `id` |  | `users:read` |
| `PIPEDRIVE_LIST_USER_PERMISSIONS` | List user permissions | read | `id` |  | `users:read` |
| `PIPEDRIVE_LIST_USER_ROLE_ASSIGNMENTS` | List User Role Assignments | read | `id`; *limit, start* |  | `users:read` |
| `PIPEDRIVE_LIST_USER_ROLE_SETTINGS` | List user role settings | read | `id` |  | `users:read` |
| `PIPEDRIVE_LIST_USER_SETTINGS` | List settings of an authorized user | read |  |  | `base` |
| `PIPEDRIVE_MERGE_DEALS` | Merge Deals | write | `id`, `merge_with_id` |  | `mail:read` |
| `PIPEDRIVE_MERGE_ORGANIZATIONS` | Merge Organizations | write | `id`, `merge_with_id` |  | `mail:read` |
| `PIPEDRIVE_MERGE_PERSONS` | Merge two persons | write | `id`, `merge_with_id` |  | `mail:read` |
| `PIPEDRIVE_MERGE_TWO_DEALS` | Merge two deals | write | `id`, `merge_with_id` |  | `deals:full` |
| `PIPEDRIVE_MERGE_TWO_ORGANIZATIONS` | Merge two organizations | write | `id`, `merge_with_id` |  | `contacts:full` |
| `PIPEDRIVE_MERGE_TWO_PERSONS` | Merge two persons | write | `id`, `merge_with_id` |  | `Persons:merge` |
| `PIPEDRIVE_ORGANIZATION_ACCESSIBLE_USER_LIST` | Organization accessible user list | read | `id` |  | `contacts:full` |
| `PIPEDRIVE_PATCH_DEAL_FIELDS` | Patch deal fields | write | `field_code`; *description, field_name, important_fields, required_fields, ui_visibility* |  | `admin` |
| `PIPEDRIVE_PATCH_ORGANIZATION_FIELDS` | Update organization field (v2) | write | `field_code`; *description, field_name, important_fields, required_fields, ui_visibility* |  | `admin` |
| `PIPEDRIVE_PATCH_PERSON_FIELDS` | Update Person Field Configuration | write | `field_code`; *description, field_name, important_fields, required_fields, ui_visibility* |  | `admin` |
| `PIPEDRIVE_PRODUCTFIELDS_PATCH_PRODUCTFIELDS` | Update product field | write | `field_code`; *description, field_name, ui_visibility* |  | `products:full` |
| `PIPEDRIVE_REFRESH_TOKENS` | Refreshing the tokens | write | *client_id, client_secret, grant_type, refresh_token* |  | `base` |
| `PIPEDRIVE_REQUEST_AUTHORIZATION` | Requesting authorization | read | `client_id`, `redirect_uri`; *state* |  | `base` |
| `PIPEDRIVE_SEARCH_DEALS` | Search deals | read | `term`; *exact_match, fields, include_fields, limit, organization_id, person_id, start, status* |  | `deals:read`, `deals:full`, `search:read` |
| `PIPEDRIVE_SEARCH_ITEM_BY_FIELD` | Search Item By Field | read | `field_key`, `field_type`, `term`; *exact_match, limit, return_item_ids, start* |  | `search:read` |
| `PIPEDRIVE_SEARCH_LEADS` | Search leads | read | `term`; *exact_match, fields, include_fields, limit, organization_id, person_id, start* |  | `leads:read` |
| `PIPEDRIVE_SEARCH_MULTIPLE_ITEM_TYPES` | Perform a search from multiple item types | read | `term`; *exact_match, fields, include_fields, item_types, limit, search_for_related_items, start* |  | `deals:read`, `persons:read`, `organizations:read`, `leads:read`, `products:read`, `files:read`, `mail_attachments:read`, `projects:read` |
| `PIPEDRIVE_SEARCH_ORGANIZATIONS` | Search organizations | read | `term`; *exact_match, fields, limit, start* |  | `organizations:read` |
| `PIPEDRIVE_SEARCH_PERSONS` | Search persons | read | `term`; *exact_match, fields, include_fields, limit, organization_id, start* |  | `persons:read` |
| `PIPEDRIVE_SEARCH_PRODUCTS` | Search products | read | `term`; *exact_match, fields, include_fields, limit, start* |  | `products:read` |
| `PIPEDRIVE_TAGES_UPDATE_STAGE` | Update stage | write | `id`; *days_to_rotten, deal_probability, is_deal_rot_enabled, name, pipeline_id* |  | `admin` |
| `PIPEDRIVE_UNLINK_USER_FROM_VIDEO_CALL_INTEGRATION` | Unlinkuserfromvideocallintegration | write | `id` |  | `video-calls` |
| `PIPEDRIVE_UPDATE_ACTIVITY` | Update an activity | write | `id`; *attendees, busy, deal_id, done, due_date, due_time, duration, lead_id and 11 more* |  | `activities:full` |
| `PIPEDRIVE_UPDATE_ACTIVITY_IN_PROJECT_PLAN` | Update activity in project plan | write | `activityId`, `id`; *group_id, phase_id* |  | `projects:full` |
| `PIPEDRIVE_UPDATE_ACTIVITY_TYPE` | Update activity type | write | `id`; *color, icon_key, name, order_nr* |  | `admin` |
| `PIPEDRIVE_UPDATE_AN_ACTIVITY` | Update an activity | write | `id`; *deal_id, done, due_date, due_time, duration, location, note, org_id and 4 more* |  | `activities:read`, `activities:write` |
| `PIPEDRIVE_UPDATE_AN_ACTIVITY_TYPE` | Update an activity type | write | `id`; *color, icon_key, name, order_nr* |  | `admin` |
| `PIPEDRIVE_UPDATE_AN_ORGANIZATION` | Update an organization | write | `id`; *address, annual_revenue, employee_count, industry, label_ids, linkedin, name, owner_id and 2 more* |  | `organizations:read`, `organizations:write` |
| `PIPEDRIVE_UPDATE_AN_ORGANIZATION_FIELD` | Update an organization field | write | `id`; *add_visible_flag, name, options* |  | `admin` |
| `PIPEDRIVE_UPDATE_AN_ORGANIZATION_RELATIONSHIP` | Update an organization relationship | write | `id`; *org_id, rel_linked_org_id, rel_owner_org_id, type* |  | `contacts:full` |
| `PIPEDRIVE_UPDATE_A_COMMENT_RELATED_TO_A_NOTE` | Update a comment related to a note | write | `commentId`, `content`, `id` |  | `notes:read`, `notes:write` |
| `PIPEDRIVE_UPDATE_A_DEAL` | Update a deal | write | `id`; *currency, expected_close_date, lost_reason, org_id, person_id, pipeline_id, probability, stage_id and 5 more* |  | `deals:full` |
| `PIPEDRIVE_UPDATE_A_DEAL_FIELD` | Update a deal field | write | `id`; *add_visible_flag, name, options* |  | `admin` |
| `PIPEDRIVE_UPDATE_A_FILTER` | Update a filter | write | `id`; *conditions, name* |  | `deals:full`, `activities:full`, `contacts:full` |
| `PIPEDRIVE_UPDATE_A_LEAD` | Update a lead | write | `id`; *expected_close_date, is_archived, label_ids, organization_id, owner_id, person_id, title, value__amount and 2 more* |  | `leads:full` |
| `PIPEDRIVE_UPDATE_A_LEAD_LABEL` | Update a lead label | write | `id`; *color, name* |  | `leads:full` |
| `PIPEDRIVE_UPDATE_A_NOTE` | Update a note | write | `id`; *add_time, content, deal_id, lead_id, org_id, person_id, pinned_to_deal_flag, pinned_to_lead_flag and 5 more* |  | `notes:read`, `notes:write` |
| `PIPEDRIVE_UPDATE_A_PERSON` | Update a person | write | `id`; *add_time, email, label_ids, marketing_status, name, org_id, owner_id, phone and 2 more* |  | `contacts:full` |
| `PIPEDRIVE_UPDATE_A_PERSON_FIELD` | Update a person field | write | `id`; *add_visible_flag, name, options* |  | `contact-fields:full`, `admin` |
| `PIPEDRIVE_UPDATE_A_PIPELINE` | Update a pipeline | write | `id`; *is_deal_probability_enabled, name* |  | `admin` |
| `PIPEDRIVE_UPDATE_A_PRODUCT` | Update a product | write | `id`; *billing_frequency, billing_frequency_cycles, category, code, description, is_linkable, name, owner_id and 4 more* |  | `products:full` |
| `PIPEDRIVE_UPDATE_A_PRODUCT_FIELD` | Update a product field | write | `id`; *name, options* |  | `product-fields:full`, `products:full` |
| `PIPEDRIVE_UPDATE_COMMENT_FOR_NOTE` | Update comment for note | write | `commentId`, `content`, `id` |  | `deals:full` |
| `PIPEDRIVE_UPDATE_DEAL` | Update a deal | write | `id`; *archive_time, close_time, currency, expected_close_date, is_archived, label_ids, lost_reason, lost_time and 11 more* |  | `deals:full` |
| `PIPEDRIVE_UPDATE_DEALFIELDS_OPTIONS` | Update deal field options | write | `field_code`, `options` |  | `admin` |
| `PIPEDRIVE_UPDATE_DEAL_DISCOUNT` | Update deal discount | write | `discount_id`, `id`; *amount, description, type* |  | `deals:read` |
| `PIPEDRIVE_UPDATE_DEAL_FIELD` | Update a deal field | write | `id`; *add_visible_flag, name, options* |  | `admin` |
| `PIPEDRIVE_UPDATE_DEAL_PRODUCT` | Update Deal Product | write | `id`, `product_attachment_id`; *billing_frequency, billing_frequency_cycles, billing_start_date, comments, discount, discount_type, is_enabled, item_price and 3 more* |  | `deals:full` |
| `PIPEDRIVE_UPDATE_DEAL_PRODUCT_ATTACHMENT` | Update the product attached to a deal | write | `id`, `product_attachment_id`; *comments, discount, discount_type, duration, enabled_flag, item_price, product_id, product_variation_id and 3 more* |  | `deals:read`, `deals:write`, `products:read`, `products:write` |
| `PIPEDRIVE_UPDATE_DEAL_V2` | Update deal (v2 API) | write | `id`; *archive_time, close_time, currency, custom_fields, expected_close_date, is_archived, is_deleted, label_ids and 13 more* |  | `deals:full` |
| `PIPEDRIVE_UPDATE_FILE` | Update file details | write | `id`; *description, name* |  | `deals:full`, `activities:full`, `contacts:full` |
| `PIPEDRIVE_UPDATE_FILE_DETAILS` | Update file details | write | `id`; *description, name* |  | `deals:full`, `activities:full`, `contacts:full` |
| `PIPEDRIVE_UPDATE_FILTER` | Update filter | write | `id`; *name* |  | `filters:read`, `filters:write` |
| `PIPEDRIVE_UPDATE_GOAL` | Update existing goal | write | `id`; *assignee, duration, expected_outcome, interval, title, type* |  |  |
| `PIPEDRIVE_UPDATE_INSTALLMENT_SUBSCRIPTION` | Update an installment subscription | write | `id`, `payments`; *update_deal_value* |  | `deals:full` |
| `PIPEDRIVE_UPDATE_LEAD` | Update a lead | write | `id`; *channel, channel_id, expected_close_date, is_archived, label_ids, organization_id, owner_id, person_id and 4 more* |  | `leads:full` |
| `PIPEDRIVE_UPDATE_LEAD_LABEL` | Update Lead Label | write | `id`; *color, name* |  | `leads:full` |
| `PIPEDRIVE_UPDATE_MAIL_THREAD` | Update mail thread details | write | `id`; *archived_flag, deal_id, lead_id, read_flag, shared_flag* |  | `mail:full` |
| `PIPEDRIVE_UPDATE_NOTE` | Update Note | write | `id`; *add_time, content, deal_id, lead_id, org_id, person_id, pinned_to_deal_flag, pinned_to_lead_flag and 5 more* |  | `deals:full` |
| `PIPEDRIVE_UPDATE_ORGANIZATION` | Update Organization | write | `id`; *address, label_ids, name, owner_id, visible_to* |  | `contacts:full` |
| `PIPEDRIVE_UPDATE_ORGANIZATION_FIELD` | Update organization field | write | `id`; *add_visible_flag, name, options* |  | `admin` |
| `PIPEDRIVE_UPDATE_ORGANIZATION_FIELD_OPTIONS` | Update organization field options | write | `field_code`, `options` |  | `admin` |
| `PIPEDRIVE_UPDATE_ORGANIZATION_RELATIONSHIP` | Update organization relationship | write | `id`; *org_id, rel_linked_org_id, rel_owner_org_id, type* |  | `contacts:full` |
| `PIPEDRIVE_UPDATE_PERSON` | Update a person | write | `id`; *add_time, email, label, label_ids, marketing_status, name, org_id, owner_id and 3 more* |  | `contacts:full` |
| `PIPEDRIVE_UPDATE_PERSONFIELDS_OPTIONS` | Update person field options | write | `field_code`, `options` |  | `admin` |
| `PIPEDRIVE_UPDATE_PERSON_FIELD` | Update person field | write | `id`; *add_visible_flag, name, options* |  | `admin` |
| `PIPEDRIVE_UPDATE_PERSON_V2` | Update person (v2) | write | `id`; *add_time, emails, first_name, label_ids, last_name, marketing_status, name, org_id and 4 more* |  | `contacts:read` |
| `PIPEDRIVE_UPDATE_PIPELINE` | Update a pipeline | write | `id`; *active, deal_probability, name, order_nr* |  | `admin` |
| `PIPEDRIVE_UPDATE_PIPELINE_VISIBILITY_FOR_ROLE` | Update pipeline visibility for a role | write | `id`, `visible_pipeline_ids` |  | `admin` |
| `PIPEDRIVE_UPDATE_PRODUCT` | Update a product | write | `id`; *active_flag, category, code, description, name, owner_id, prices, selectable and 3 more* |  | `products:full` |
| `PIPEDRIVE_UPDATE_PRODUCTFIELDS_OPTIONS` | Update product field options | write | `field_code`, `options` |  | `products:full` |
| `PIPEDRIVE_UPDATE_PRODUCT_FIELD` | Update product field | write | `id`; *name, options* |  | `products:full` |
| `PIPEDRIVE_UPDATE_PRODUCT_IMAGE` | Update product image | write | `data`, `id` |  | `products:full` |
| `PIPEDRIVE_UPDATE_PRODUCT_VARIATION` | Update product variation | write | `id`, `product_variation_id`; *name, prices* |  | `products:full` |
| `PIPEDRIVE_UPDATE_PROJECT` | Update a project | write | `id`; *board_id, deal_ids, description, end_date, labels, org_id, owner_id, person_id and 4 more* |  | `projects:full` |
| `PIPEDRIVE_UPDATE_RECURRING_SUBSCRIPTION` | Update a recurring subscription | write | `effective_date`, `id`; *cycle_amount, description, payments, update_deal_value* |  | `deals:full` |
| `PIPEDRIVE_UPDATE_ROLE` | Update role details | write | `id`; *name, parent_role_id* |  | `admin` |
| `PIPEDRIVE_UPDATE_ROLE_PIPELINES` | Update role pipelines | write | `id`, `visible_pipeline_ids` |  | `admin` |
| `PIPEDRIVE_UPDATE_STAGE` | Update stage details | write | `id`; *deal_probability, name, order_nr, pipeline_id, rotten_days, rotten_flag* |  | `admin` |
| `PIPEDRIVE_UPDATE_STAGE_DETAILS` | Update stage details | write | `id`; *deal_probability, name, order_nr, pipeline_id, rotten_days, rotten_flag* |  | `stages:read`, `stages:write` |
| `PIPEDRIVE_UPDATE_TASK` | Update a task | write | `id`; *assignee_id, description, done, due_date, parent_task_id, project_id, title* |  | `projects:read` |
| `PIPEDRIVE_UPDATE_TASK_IN_PROJECT_PLAN` | Update task in project plan | write | `id`, `taskId`; *group_id, phase_id* |  | `projects:full` |
| `PIPEDRIVE_UPDATE_TEAM` | Update a team | write | `id`; *active_flag, deleted_flag, description, manager_id, name, users* |  | `teams:read`, `teams:write` |
| `PIPEDRIVE_UPDATE_USER` | Update user details | write | `active_flag`, `id` |  | `admin` |
| `PIPEDRIVE_UPDATE_USER_DETAILS` | Update user details | write | `active_flag`, `id` |  | `admin` |
| `PIPEDRIVE_WHO_AM_I` | Who Am I | read |  |  |  |
