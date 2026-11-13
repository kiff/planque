# Salesforce: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

288 tools: 181 read, 107 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `SALESFORCE_ABORT_QUERY_JOB` | Abort query job | write | `job_id` |  | `full`, `api` |
| `SALESFORCE_ADD_CONTACT_TO_CAMPAIGN` | Add contact to campaign | write | `campaign_id`, `contact_id`; *custom_fields, status* |  | `full`, `api` |
| `SALESFORCE_ADD_LEAD_TO_CAMPAIGN` | Add lead to campaign | write | `campaign_id`, `lead_id`; *custom_fields, status* |  | `full`, `api` |
| `SALESFORCE_ADD_OPPORTUNITY_LINE_ITEM` | Add product to opportunity | write | `opportunity_id`, `pricebook_entry_id`, `quantity`; *custom_fields, description, discount, service_date, total_price, unit_price* |  | `full`, `api` |
| `SALESFORCE_APPLY_LEAD_ASSIGNMENT_RULES` | Apply lead assignment rules | write | `lead_id` |  | `full`, `api` |
| `SALESFORCE_ASSOCIATE_CONTACT_TO_ACCOUNT` | Associate contact to account | write | `account_id`, `contact_id`; *custom_fields* |  | `full`, `api` |
| `SALESFORCE_CLONE_OPPORTUNITY_WITH_PRODUCTS` | Clone opportunity with products | write | `opportunity_id`; *clone_products, close_date, new_name, stage_name* |  | `full`, `api` |
| `SALESFORCE_CLONE_RECORD` | Clone record | write | `object_type`, `record_id`; *field_updates* |  | `full`, `api` |
| `SALESFORCE_CLONE_REPORT` | Clone report | write | `name`, `report_id` |  | `full`, `api` |
| `SALESFORCE_CLOSE_OR_ABORT_JOB` | Close or abort a job | write | `job_id`, `state` |  | `full`, `api` |
| `SALESFORCE_COMPLETE_TASK` | Complete task | write | `task_id`; *completion_notes, custom_fields* |  | `full`, `api` |
| `SALESFORCE_CREATE_ACCOUNT` | Create account | write | `name`; *account_source, annual_revenue, billing_city, billing_country, billing_postal_code, billing_state, billing_street, custom_fields and 14 more* |  | `full`, `api` |
| `SALESFORCE_CREATE_ANALYTICS_NOTIFICATION` | Create analytics notification | write | *active, configuration, deactivateOnTrigger, name, recordId, runAs, schedule, source and 1 more* |  | `full`, `api` |
| `SALESFORCE_CREATE_A_RECORD` | Create a record | write | `api_name`, `fields`; *allow_save_on_duplicate* |  | `full`, `api` |
| `SALESFORCE_CREATE_BULK_INGEST_JOB` | Create bulk ingest job | write | `object_type`, `operation`; *assignment_rule_id, column_delimiter, content_type, external_id_field_name, line_ending* |  | `full`, `api` |
| `SALESFORCE_CREATE_BULK_QUERY_JOB` | Create bulk query job | read | `query`; *column_delimiter, content_type, line_ending, operation* |  | `full`, `api` |
| `SALESFORCE_CREATE_CAMPAIGN` | Create campaign | write | `name`; *actual_cost, budgeted_cost, custom_fields, description, end_date, expected_response, expected_revenue, is_active and 5 more* |  | `full`, `api` |
| `SALESFORCE_CREATE_CONTACT` | Create contact | write | `last_name`; *account_id, birthdate, custom_fields, department, email, first_name, lead_source, mailing_city and 7 more* |  | `full`, `api` |
| `SALESFORCE_CREATE_CUSTOM_FIELD` | Create custom field | write | `field_api_name`, `field_type`, `label`, `object_api_name`; *check_existing, default_value, description, external_id, inline_help_text, length, picklist_values, precision and 8 more* |  | `full`, `api` |
| `SALESFORCE_CREATE_CUSTOM_OBJECT` | Create custom object | write | `full_name`, `label`, `name_field`, `plural_label`; *deployment_status, description, enable_activities, enable_reports, sharing_model* |  | `full`, `api` |
| `SALESFORCE_CREATE_FOLDER` | Create folder | write | `label`, `type`; *name, parentId* |  | `full`, `api` |
| `SALESFORCE_CREATE_FOLDER_SHARES` | Create folder shares | write | `folder_id`, `shares` |  | `full`, `api` |
| `SALESFORCE_CREATE_KNOWLEDGE_ARTICLE_DRAFT` | Create Knowledge article draft | write | `article_id` |  | `full`, `api` |
| `SALESFORCE_CREATE_LEAD` | Create lead | write | `Company`, `LastName`; *AnnualRevenue, City, Country, CustomFields, Email, FirstName, Industry, LeadSource and 10 more* |  | `full`, `api` |
| `SALESFORCE_CREATE_LIST_VIEW` | Create list view | write | `object_api_name`; *display_columns, filter_logic_string, filtered_by_info, label, list_shares, list_view_api_name, scope, visibility* |  | `api`, `full` |
| `SALESFORCE_CREATE_NOTE` | Create note | write | `parent_id`, `title`; *body, custom_fields, is_private, owner_id* |  | `full`, `api` |
| `SALESFORCE_CREATE_OPPORTUNITY` | Create opportunity | write | `close_date`, `name`, `stage_name`; *account_id, amount, contact_id, custom_fields, description, lead_source, next_step, pricebook2_id and 2 more* |  | `full`, `api` |
| `SALESFORCE_CREATE_SOBJECT_TREE` | Create sObject tree | write | `records`, `sobject_api_name` |  | `full`, `api` |
| `SALESFORCE_CREATE_S_OBJECT_RECORD` | Create SObject record | write | `fields`, `sobject_type` |  | `full`, `api` |
| `SALESFORCE_CREATE_TASK` | Create task | write | `subject`; *activity_date, custom_fields, description, is_reminder_set, owner_id, priority, reminder_date_time, status and 2 more* |  | `full`, `api` |
| `SALESFORCE_CREATE_TOOLING_FLOW` | Create Tooling Flow | write | `full_name`, `metadata`; *is_template* |  | `full`, `api` |
| `SALESFORCE_CREATE_TOOLING_RECORD` | Create Tooling API record | write | `fields`, `object_name` |  | `full`, `api` |
| `SALESFORCE_CREATE_VALIDATION_RULE` | Create validation rule | write | `active`, `error_condition_formula`, `error_message`, `full_name`; *description, error_display_field* |  | `full`, `api` |
| `SALESFORCE_DELETE_ACCOUNT` | Delete account | write | `account_id` |  | `full`, `api` |
| `SALESFORCE_DELETE_ANALYTICS_NOTIFICATION` | Delete analytics notification | write | `notification_id` |  | `full`, `api` |
| `SALESFORCE_DELETE_CAMPAIGN` | Delete campaign | write | `campaign_id` |  | `full`, `api` |
| `SALESFORCE_DELETE_CONTACT` | Delete contact | write | `contact_id` |  | `full`, `api` |
| `SALESFORCE_DELETE_DASHBOARD` | Delete dashboard | write | `dashboard_id` |  | `full`, `api` |
| `SALESFORCE_DELETE_FILE` | Delete file | write | `file_id` |  | `full`, `api` |
| `SALESFORCE_DELETE_FOLDER` | Delete folder | write | `folder_id` |  | `full`, `api` |
| `SALESFORCE_DELETE_FOLDER_SHARE` | Delete folder share | write | `folder_id`, `share_id` |  | `full`, `api` |
| `SALESFORCE_DELETE_INGEST_JOB` | Delete ingest job | write | `job_id` |  | `full`, `api` |
| `SALESFORCE_DELETE_JOB_QUERY` | Delete job query | write | `job_id` |  | `full`, `api` |
| `SALESFORCE_DELETE_LEAD` | Delete lead | write | `lead_id` |  | `full`, `api` |
| `SALESFORCE_DELETE_NOTE` | Delete note | write | `note_id` |  | `full`, `api` |
| `SALESFORCE_DELETE_OPPORTUNITY` | Delete opportunity | write | `opportunity_id` |  | `full`, `api` |
| `SALESFORCE_DELETE_REPORT` | Delete report | write | `report_id` |  | `full`, `api` |
| `SALESFORCE_DELETE_SOBJECT` | Delete sObject record | write | `id`, `sobject_name` |  | `api`, `full` |
| `SALESFORCE_DELETE_SOBJECT_COLLECTIONS` | Delete multiple records (SObject Collections) | write | `ids`; *all_or_none* |  | `full`, `api` |
| `SALESFORCE_DELETE_SURVEY_TRANSLATIONS` | Delete survey translations | write | `mode`; *developer_name, language, translations* |  | `full`, `api` |
| `SALESFORCE_DELETE_TOOLING_RECORD` | Delete Tooling API record | write | `record_id`, `tooling_object` |  | `full`, `api` |
| `SALESFORCE_DESCRIBE_DASHBOARD` | Describe dashboard | read | `dashboard_id`; *load_component_properties* |  | `full`, `api` |
| `SALESFORCE_DESCRIBE_LIST_VIEW` | Describe list view | read | `list_view_id`, `sobject_name` |  | `full`, `api` |
| `SALESFORCE_DESCRIBE_REPORT_TYPE` | Describe report type | read | `report_type` |  | `full`, `api` |
| `SALESFORCE_DESCRIBE_TOOLING_OBJECT` | Describe Tooling object | read | `tooling_object` |  | `full`, `api` |
| `SALESFORCE_DISCOVER_UNIT_TESTS` | Discover unit tests | read | *category, namespacePrefix, nextRecord, pageSize, showAllMethods* |  | `full`, `api` |
| `SALESFORCE_DOWNLOAD_ATTACHMENT_BODY` | Download attachment body | read | `attachment_id` |  | `full`, `api` |
| `SALESFORCE_EXECUTE_COMPOSITE` | Execute composite | write | `composite_request`; *all_or_none, collate_subrequests* |  | `full`, `api` |
| `SALESFORCE_EXECUTE_COMPOSITE_BATCH` | Execute composite batch | write | `batch_requests`; *halt_on_error* |  | `full`, `api` |
| `SALESFORCE_EXECUTE_GLOBAL_QUICK_ACTION` | Execute global quick action | write | `action_name`, `record`; *context_id* |  | `full`, `api` |
| `SALESFORCE_EXECUTE_INVOCABLE_ACTION` | Execute invocable action | write | `action_name`, `action_type`, `inputs` |  | `full`, `api` |
| `SALESFORCE_EXECUTE_SOBJECT_QUICK_ACTION` | Execute sObject Quick Action | write | `action_name`, `record`, `sobject`; *context_id* |  | `api`, `full` |
| `SALESFORCE_EXECUTE_SOSL_SEARCH` | Execute SOSL search | write | `q` |  | `full`, `api` |
| `SALESFORCE_EXECUTE_TOOLING_COMPOSITE` | Execute Tooling composite | write | `composite_request`; *all_or_none* |  | `full`, `api` |
| `SALESFORCE_GET_ACCOUNT` | Get account | read | `account_id` |  | `full`, `api` |
| `SALESFORCE_GET_ALL_CUSTOM_OBJECTS` | Get all objects | write | *custom_only, page, page_size* |  | `full`, `api` |
| `SALESFORCE_GET_ALL_FIELDS_FOR_OBJECT` | Get all fields for an object | read | `object_name` |  | `full`, `api` |
| `SALESFORCE_GET_ALL_NAVIGATION_ITEMS` | Get all navigation items | read | *formFactor, navItemNames, page, pageSize* |  | `api`, `full` |
| `SALESFORCE_GET_ANALYTICS_NOTIFICATION` | Get analytics notification | read | `notification_id`, `source`; *owner_id, record_id* |  | `full`, `api` |
| `SALESFORCE_GET_ANALYTICS_NOTIFICATION_LIMITS` | Get Analytics notification limits | read | `source`; *record_id* |  | `full`, `api` |
| `SALESFORCE_GET_APEX_LOG_BODY` | Get Apex log body | read | `log_id`; *use_tooling_api* |  | `full`, `api` |
| `SALESFORCE_GET_API` | Get API resources by version | read | *version* |  | `full`, `api` |
| `SALESFORCE_GET_APP` | Get an app | read | `app_id`; *form_factor* |  | `full`, `api` |
| `SALESFORCE_GET_APPOINTMENT_CANDIDATES` | Get appointment candidates | read | `territoryIds`; *accountId, allowConcurrentScheduling, correlationId, endTime, engagementChannelTypeIds, filterByResources, resourceLimitApptDistribution, schedulingPolicyId and 3 more* |  | `full`, `api` |
| `SALESFORCE_GET_APPOINTMENT_SLOTS` | Get appointment slots | read | `requiredResourceIds`, `territoryIds`; *accountId, allowConcurrentScheduling, correlationId, endTime, engagementChannelTypeIds, primaryResourceId, schedulingPolicyId, startTime and 2 more* |  | `full`, `api` |
| `SALESFORCE_GET_APPS` | Get apps | read | *form_factor* |  | `api`, `full` |
| `SALESFORCE_GET_APP_MENU` | Get App Menu | read | *headers_only, menu* |  | `full`, `api` |
| `SALESFORCE_GET_A_BATCH_OF_RECORDS` | Get batch of UI API records | read | `record_ids`; *fields, optional_fields* |  | `full`, `api` |
| `SALESFORCE_GET_CAMPAIGN` | Get campaign | read | `campaign_id` |  | `full`, `api` |
| `SALESFORCE_GET_CHATTER_RESOURCES` | Get Chatter resources | read |  |  | `full`, `api`, `chatter_api` |
| `SALESFORCE_GET_CHILD_RECORDS` | Get child records | read | `record_id`, `relationship_name`; *fields, optional_fields, page_size, page_token* |  | `full`, `api` |
| `SALESFORCE_GET_COMPACT_LAYOUTS` | Get Compact Layouts | read | `object_list` |  | `full`, `api` |
| `SALESFORCE_GET_COMPOSITE_RESOURCES` | Get composite resources | read |  |  | `full`, `api` |
| `SALESFORCE_GET_COMPOSITE_SOBJECTS` | Get multiple records by IDs | read | `fields`, `ids`, `sobject_name` |  | `full`, `api` |
| `SALESFORCE_GET_CONSENT_ACTION` | Get consent preferences | read | `action`, `ids` |  | `full`, `api` |
| `SALESFORCE_GET_CONTACT` | Get contact | read | `contact_id`; *fields* |  | `full`, `api` |
| `SALESFORCE_GET_CONTACT_BY_ID` | Retrieve specific contact by id | read | `id`; *fields* |  | `full`, `api` |
| `SALESFORCE_GET_DASHBOARD` | Get dashboard | read | `dashboard_id` |  | `full`, `api` |
| `SALESFORCE_GET_DASHBOARD_COMPONENTS` | Get dashboard components | read | `component_ids`, `dashboard_id` |  | `full`, `api` |
| `SALESFORCE_GET_DASHBOARD_FILTER_OPTIONS` | Get dashboard filter options | read | `filter_columns`, `options` |  | `full`, `api` |
| `SALESFORCE_GET_DASHBOARD_STATUS` | Get dashboard status | read | `dashboard_id`; *filter1, filter2, filter3, running_user* |  | `full`, `api` |
| `SALESFORCE_GET_DATA_CATEGORIES` | Get data categories | read | *accept_language, category, group, operation, top_categories_only* |  | `full`, `api` |
| `SALESFORCE_GET_DELETED_RECORDS` | Get deleted records | read | `end`, `sobject`, `start` |  | `full`, `api` |
| `SALESFORCE_GET_EMBEDDED_SERVICE_CONFIGURATION` | Get Embedded Service configuration | read | `embedded_service_config_developer_name`; *headers_only* |  | `full`, `api` |
| `SALESFORCE_GET_ENHANCED_NOTE_CONTENT` | Get enhanced note content | read | `content_note_id` |  | `full`, `api` |
| `SALESFORCE_GET_FILE_CONTENT` | Get file content | read | `file_id`; *rendition_type, version_id* |  | `full`, `api`, `chatter_api` |
| `SALESFORCE_GET_FILE_INFORMATION` | Get file information | read | `file_id` |  | `full`, `api`, `chatter_api` |
| `SALESFORCE_GET_FILE_SHARES` | Get file shares | read | `file_id`; *page, page_size* |  | `full`, `api`, `chatter_api` |
| `SALESFORCE_GET_FOLDER` | Get folder | read | `folder_id` |  | `full`, `api` |
| `SALESFORCE_GET_FOLDER_SHARE` | Get folder share | read | `folder_id`, `share_id` |  | `full`, `api` |
| `SALESFORCE_GET_FOLDER_SHARES` | Get folder shares | read | `folder_id` |  | `full`, `api` |
| `SALESFORCE_GET_GLOBAL_ACTIONS` | Get global actions | read | *actionTypes, apiNames, formFactor* |  | `api`, `full` |
| `SALESFORCE_GET_INGEST_JOB_INFO` | Get ingest job info | read | `job_id` |  | `full`, `api` |
| `SALESFORCE_GET_INVALID_APEX_COMPILATION_RESULTS` | Get invalid Apex compilation results | read | *api_version* |  | `full`, `api` |
| `SALESFORCE_GET_JOB_FAILED_RECORD_RESULTS` | Get job failed record results | read | `job_id`; *locator, max_records* |  | `full`, `api` |
| `SALESFORCE_GET_JOB_SUCCESSFUL_RECORD_RESULTS` | Get job successful record results | read | `job_id`; *locator, max_records* |  | `full`, `api` |
| `SALESFORCE_GET_JOB_UNPROCESSED_RECORD_RESULTS` | Get job unprocessed record results | read | `job_id`; *locator, max_records* |  | `full`, `api` |
| `SALESFORCE_GET_LAST_SELECTED_APP` | Get last selected app | read | *formFactor* |  | `full`, `api` |
| `SALESFORCE_GET_LEAD` | Get lead | read | `lead_id` |  | `full`, `api` |
| `SALESFORCE_GET_LIST_VIEW_ACTIONS` | Get List View Header Actions | read | `list_view_id` |  | `full`, `api` |
| `SALESFORCE_GET_LIST_VIEW_METADATA_BATCH` | Get batch of list view metadata | read | `list_view_ids` |  | `full`, `api` |
| `SALESFORCE_GET_LIST_VIEW_METADATA_BY_NAME` | Get list view metadata | read | `list_view_api_name`, `sobject_api_name` |  | `api`, `full` |
| `SALESFORCE_GET_LIST_VIEW_RECORDS_BY_ID` | Get list view records by ID | read | `list_view_id`; *fields, optional_fields, page_size, page_token, sort_by* |  | `api`, `full` |
| `SALESFORCE_GET_LIST_VIEW_RECORDS_BY_NAME` | Get list view records by API name | read | `list_view_api_name`, `sobject_api_name`; *fields, optionalFields, pageSize, pageToken, searchTerm, sortBy, where* |  | `full`, `api` |
| `SALESFORCE_GET_LIST_VIEW_RESULTS` | Get list view results | read | `list_view_id`, `sobject` |  | `full`, `api` |
| `SALESFORCE_GET_LOOKUP_FIELD_SUGGESTIONS` | Get lookup field suggestions | read | `field_api_name`, `object_api_name`; *dependent_field_bindings, form_factor, page, page_size, q, search_type* |  | `full`, `api` |
| `SALESFORCE_GET_LOOKUP_SUGGESTIONS_CASE_CONTACT` | Get lookup field suggestions with POST | read | `field_api_name`, `object_api_name`, `source_record`; *dependent_field_bindings, page, page_size, q, search_type* |  | `full`, `api` |
| `SALESFORCE_GET_METADATA_DEPLOYMENT_STATUS` | Get Metadata Deployment Status | read | `deployment_request_id`; *include_details* |  | `full`, `api` |
| `SALESFORCE_GET_NAMED_LAYOUT` | Get named layout | read | `layout_name`, `sobject` |  | `full`, `api` |
| `SALESFORCE_GET_NOTE` | Get note | read | `note_id` |  | `full`, `api` |
| `SALESFORCE_GET_OBJECT_LIST_VIEWS` | Get list views for an object | read | `sobject_api_name` |  | `api`, `full` |
| `SALESFORCE_GET_OBJECT_UI_LAYOUT` | Get object UI layout | read | `object_api_name`; *form_factor, layout_type, mode, record_type_id* |  | `api`, `full` |
| `SALESFORCE_GET_OPPORTUNITY` | Get opportunity | read | `opportunity_id`; *fields* |  | `full`, `api` |
| `SALESFORCE_GET_ORG_LIMITS` | Get organization limits | read |  |  | `full`, `api` |
| `SALESFORCE_GET_PHOTO_ACTIONS` | Get photo actions | read | `record_ids`; *action_types, form_factor* |  | `full`, `api` |
| `SALESFORCE_GET_PICKLIST_VALUES_BY_RECORD_TYPE` | Get Values for All Picklist Fields | read | `record_type_id`, `sobject_api_name` |  | `api`, `full` |
| `SALESFORCE_GET_PLATFORM_EVENT_SCHEMA` | Get platform event schema | read | *event_name, payload_format, schema_id* |  | `full`, `api` |
| `SALESFORCE_GET_PROCESS_APPROVALS` | Get Process Approvals | read | *headers_only* |  | `full`, `api` |
| `SALESFORCE_GET_PROCESS_RULES` | Get process rules | read | *headers_only, scope, sobject_name, workflow_rule_id* |  | `full`, `api` |
| `SALESFORCE_GET_QUERY_JOB_INFO` | Get job info query | read | `job_id` |  | `full`, `api` |
| `SALESFORCE_GET_QUERY_JOB_RESULTS` | Get job query result | read | `job_id`; *locator, max_records* |  | `api`, `full` |
| `SALESFORCE_GET_QUERY_RESULT_PAGES` | Get query result pages | read | `job_id`; *locator* |  | `full`, `api` |
| `SALESFORCE_GET_QUICK_ACTIONS` | Get quick actions | read |  |  | `full`, `api` |
| `SALESFORCE_GET_RECORD_COUNTS` | Get record count for objects | read | `sobjects` |  | `full`, `api` |
| `SALESFORCE_GET_RECORD_EDIT_PAGE_ACTIONS` | Get record edit page actions | read | `record_ids`; *action_types, form_factor, retrieval_mode* |  | `full`, `api` |
| `SALESFORCE_GET_RELATED_LIST_ACTIONS` | Get related list actions | read | `record_ids`, `related_list_ids`; *action_types, form_factor, retrieval_mode* |  | `full`, `api` |
| `SALESFORCE_GET_RELATED_LIST_PREFERENCES_BATCH` | Get batch of related list user preferences | read | `related_list_ids` |  | `api`, `full` |
| `SALESFORCE_GET_RELATED_LIST_RECORDS_CONTACTS` | Get related list records with request body | read | `parent_record_id`, `related_list_id`; *fields, optional_fields, page_size, page_token, sort_by, where* |  | `full`, `api` |
| `SALESFORCE_GET_RELEVANT_ITEMS` | Get relevant items | read | *last_updated_id, mode, object_last_updated_ids, sobjects* |  | `full`, `api` |
| `SALESFORCE_GET_REPORT_FIELDS` | Get report fields | read | `report_id`; *intersect_with* |  | `full`, `api` |
| `SALESFORCE_GET_REPORT_INSTANCE` | Get report instance results | read | `instance_id`, `report_id` |  | `full`, `api` |
| `SALESFORCE_GET_RICH_TEXT_IMAGE` | Get rich text image | read | `content_reference_id`, `field_name`, `record_id`, `sobject` |  | `full`, `api` |
| `SALESFORCE_GET_SCHEDULING_RESOURCES` | Get scheduling resources | read |  |  | `full`, `api` |
| `SALESFORCE_GET_SEARCH_LAYOUT` | Get search result layouts | read | `object_names` |  | `full`, `api` |
| `SALESFORCE_GET_SEARCH_SCOPE` | Get search scope | read |  |  | `full`, `api` |
| `SALESFORCE_GET_SEARCH_SUGGESTIONS` | Search suggested queries | read | `channel`, `language`, `q`; *sobject* |  | `full`, `api` |
| `SALESFORCE_GET_SELF_SERVICE_USER_PASSWORD_STATUS` | Get self-service user password status | read | `self_service_user_id`; *headers_only* |  | `full`, `api` |
| `SALESFORCE_GET_SOBJECTS_SOBJECT_DESCRIBE_APPROVALLAYOUTS` | Get approval layouts for object | read | `sobject_name` |  | `full`, `api` |
| `SALESFORCE_GET_SOBJECT_BASIC_INFORMATION` | Get sObject basic information | read | `sobject` |  | `full`, `api` |
| `SALESFORCE_GET_SOBJECT_BY_EXTERNAL_ID` | Get sObject by external ID | read | `field_name`, `field_value`, `sobject`; *fields* |  | `full`, `api` |
| `SALESFORCE_GET_SOBJECT_COLLECTIONS` | Get multiple sObject records | read | `fields`, `ids`, `sobject_type`; *all_or_none* |  | `api`, `full` |
| `SALESFORCE_GET_SOBJECT_COMPACT_LAYOUTS` | Get sObject compact layouts | read | `sobject` |  | `full`, `api` |
| `SALESFORCE_GET_SOBJECT_LAYOUTS` | Get sObject layouts | read | `sobject`; *headers_only* |  | `full`, `api` |
| `SALESFORCE_GET_SOBJECT_LAYOUT_HEADERS` | Get sObject layout headers | read | `layout_type`, `sobject` |  | `full`, `api` |
| `SALESFORCE_GET_SOBJECT_LIST_VIEW` | Get sObject list view information | read | `list_view_id`, `sobject_name` |  | `full`, `api` |
| `SALESFORCE_GET_SOBJECT_LIST_VIEWS` | Get list views for sObject | read | `sobject_name` |  | `full`, `api` |
| `SALESFORCE_GET_SOBJECT_PLATFORMACTION` | Get PlatformAction SObject metadata | read |  |  | `full`, `api` |
| `SALESFORCE_GET_SOBJECT_QUICK_ACTIONS` | Get sObject quick actions | read | `sobject`; *headers_only* |  | `full`, `api` |
| `SALESFORCE_GET_SOBJECT_QUICK_ACTION_DEFAULT_VALUES` | Get sObject Quick Action Default Values | read | `action_name`, `sobject` |  | `full`, `api` |
| `SALESFORCE_GET_SOBJECT_RELATIONSHIP` | Get records using sObject relationships | read | `id`, `relationship_field_name`, `sobject`; *fields* |  | `full`, `api` |
| `SALESFORCE_GET_SUPPORT` | Get support knowledge | read |  |  | `full`, `api` |
| `SALESFORCE_GET_SUPPORTED_OBJECTS_DIRECTORY` | Get directory of supported objects | read | *q* |  | `full`, `api` |
| `SALESFORCE_GET_SUPPORT_KNOWLEDGE_ARTICLES` | Get support knowledge articles | read | *accept_language, categories, channel, page_number, page_size, publish_status, sort, topics* |  | `full`, `api` |
| `SALESFORCE_GET_SURVEY_TRANSLATIONS` | Get Survey Translations | read | `mode`; *developer_name, language, translations* |  | `full`, `api` |
| `SALESFORCE_GET_S_OBJECTS_DESCRIBE_LAYOUTS_RECORD_TYPE_ID` | Get layouts for object with record type | read | `record_type_id`, `sobject` |  | `full`, `api` |
| `SALESFORCE_GET_S_OBJECTS_UPDATED` | Get updated sObject records | read | `end`, `sobject`, `start` |  | `full`, `api` |
| `SALESFORCE_GET_S_OBJECT_QUICK_ACTION_DEFAULT_VALUES` | Get sObject Quick Action Default Values by ID | read | `action_name`, `context_id`, `sobject` |  | `full`, `api` |
| `SALESFORCE_GET_S_OBJECT_RECORD` | Get sObject record by ID | read | `record_id`, `sobject_api_name`; *fields* |  | `full`, `api` |
| `SALESFORCE_GET_TABS` | Get tabs | read | *headers_only* |  | `full`, `api` |
| `SALESFORCE_GET_THEME` | Get theme | read |  |  | `full`, `api` |
| `SALESFORCE_GET_TOOLING_COMPLETIONS` | Get Tooling completions | read | `type` |  | `full`, `api` |
| `SALESFORCE_GET_TOOLING_FLEXI_PAGE` | Get Tooling FlexiPage | read | `flexipage_id` |  | `full`, `api` |
| `SALESFORCE_GET_TOOLING_RECORD` | Get Tooling API record | read | `record_id`, `tooling_object` |  | `full`, `api` |
| `SALESFORCE_GET_UIAPI_ACTIONS_MRU_LIST_ACCOUNT` | Get MRU List View Actions | read | *action_types, form_factor, object_api_name* |  | `full`, `api` |
| `SALESFORCE_GET_UIAPI_LIST_INFO_ACCOUNT_ALL_ACCOUNTS` | Get Account AllAccounts list view metadata | read |  |  | `api`, `full` |
| `SALESFORCE_GET_UIAPI_LIST_INFO_ACCOUNT_RECENT` | Get List View Metadata | read | *list_view_api_name, object_api_name* |  | `api`, `full` |
| `SALESFORCE_GET_UIAPI_LIST_INFO_ACCOUNT_SEARCH_RESULT` | Get Account SearchResult list view metadata | read |  |  | `full`, `api` |
| `SALESFORCE_GET_UIAPI_RELATED_LIST_PREFERENCES` | Get related list user preferences | read | `object_api_name`, `related_list_id` |  | `full`, `api` |
| `SALESFORCE_GET_UI_API_ACTIONS_LOOKUP_ACCOUNT` | Get lookup field actions | read | *action_types, form_factor, object_api_name, sections* |  | `api`, `full` |
| `SALESFORCE_GET_UI_API_ACTIONS_RECORD_RELATED_LIST` | Get related list actions | read | `record_id`; *action_types, form_factor, related_list_ids, retrieval_mode* |  | `full`, `api` |
| `SALESFORCE_GET_UI_API_APPS_USER_NAV_ITEMS` | Get Personalized Navigation Items | read | `app_id`; *form_factor, nav_item_names, page, page_size* |  | `full`, `api` |
| `SALESFORCE_GET_UI_API_LIST_INFO_RECENT` | Get List Views for Object | read | *object_api_name, page_size, page_token, q, recent_lists_only* |  | `full`, `api` |
| `SALESFORCE_GET_UI_API_RECORD_UI` | Get record UI data and metadata | read | `record_ids`; *child_relationships, form_factor, layout_types, modes, optional_fields, page_size* |  | `api`, `full` |
| `SALESFORCE_GET_USER_INFO` | Get user info | read | *include_permissions, user_id* |  | `full`, `api` |
| `SALESFORCE_HEAD_ACTIONS_CUSTOM` | Get custom actions headers | write | *if_unmodified_since* |  | `full`, `api` |
| `SALESFORCE_HEAD_ACTIONS_STANDARD` | Get standard actions headers | read | *if_modified_since, if_unmodified_since* |  | `api`, `full` |
| `SALESFORCE_HEAD_APPMENU_SALESFORCE1` | Get AppMenu Salesforce1 headers | read | *if_modified_since, if_none_match* |  | `full`, `api` |
| `SALESFORCE_HEAD_PROCESS_RULES_S_OBJECT` | Get process rules headers | read | `sobject` |  | `full`, `api` |
| `SALESFORCE_HEAD_QUICK_ACTIONS` | Get Quick Actions headers | read |  |  | `full`, `api` |
| `SALESFORCE_HEAD_SOBJECTS_GLOBAL_DESCRIBE_LAYOUTS` | Get global publisher layouts headers | read | *if_modified_since, if_none_match* |  | `full`, `api` |
| `SALESFORCE_HEAD_SOBJECTS_QUICK_ACTION` | Get sObject quick action headers | read | `action_name`, `sobject_type`; *if_modified_since, if_none_match* |  | `full`, `api` |
| `SALESFORCE_HEAD_SOBJECTS_USER_PASSWORD` | Get User password headers | read | `user_id` |  | `full`, `api` |
| `SALESFORCE_HEAD_SOBJECT_QUICK_ACTION_DEFAULT_VALUES` | Get quick action default values headers by ID | read | `action_name`, `context_id`, `sobject_type`; *if_modified_since, if_none_match* |  | `full`, `api` |
| `SALESFORCE_LIST_ACCOUNTS` | List accounts | read | *query* |  | `full`, `api` |
| `SALESFORCE_LIST_ANALYTICS_NOTIFICATIONS` | List Analytics notifications | read | `source`; *owner_id, record_id* |  | `full`, `api` |
| `SALESFORCE_LIST_ANALYTICS_TEMPLATES` | List Analytics templates | read | *options, page, pageSize, type* |  | `full`, `wave_api` |
| `SALESFORCE_LIST_API_VERSIONS` | List API versions | read |  |  | `full`, `api` |
| `SALESFORCE_LIST_BULK_JOBS` | List bulk jobs | read | `job_category`; *concurrency_mode, is_pk_chunking_enabled, job_type, query_locator* |  | `full`, `api` |
| `SALESFORCE_LIST_CAMPAIGNS` | List campaigns | read | *query* |  | `full`, `api` |
| `SALESFORCE_LIST_COMMUNITIES` | List communities | read | *status* |  | `full`, `api` |
| `SALESFORCE_LIST_CONTACTS` | List contacts | read | *query* |  | `full`, `api` |
| `SALESFORCE_LIST_CUSTOM_INVOCABLE_ACTIONS` | Get custom invocable actions | write |  |  | `full`, `api` |
| `SALESFORCE_LIST_DASHBOARDS` | List dashboards | read |  |  | `full`, `api` |
| `SALESFORCE_LIST_EMAIL_TEMPLATES` | List email templates | read | *folder_name, include_body, is_active_only, limit, order_by, order_direction, search_term, template_type* |  | `full`, `api` |
| `SALESFORCE_LIST_FOLDERS` | List folders | read | *page, pageSize, type* |  | `full`, `api` |
| `SALESFORCE_LIST_FOLDER_CHILDREN` | List folder children | read | `folder_id`; *page, page_size* |  | `full`, `api` |
| `SALESFORCE_LIST_FOLDER_SHARE_RECIPIENTS` | List folder share recipients | read | `folder_id`; *limit, search_term, share_type* |  | `full`, `api` |
| `SALESFORCE_LIST_LEADS` | List leads | read | *query* |  | `full`, `api` |
| `SALESFORCE_LIST_NOTES` | List notes | read | *query* |  | `full`, `api` |
| `SALESFORCE_LIST_OPPORTUNITIES` | List opportunities | read | *query* |  | `full`, `api` |
| `SALESFORCE_LIST_PRICEBOOKS` | List pricebooks | read | *query* |  | `full`, `api` |
| `SALESFORCE_LIST_PRICEBOOK_ENTRIES` | List pricebook entries | read | *query* |  | `full`, `api` |
| `SALESFORCE_LIST_RECENT_LIST_VIEWS` | List recent list views | read | `sobject` |  | `full`, `api` |
| `SALESFORCE_LIST_RECENT_RECORDS` | List recent records | read | *limit* |  | `full`, `api` |
| `SALESFORCE_LIST_REPORTS` | List reports | read |  |  | `full`, `api` |
| `SALESFORCE_LIST_REPORT_FILTER_OPERATORS` | List report filter operators | read | *for_dashboards* |  | `full`, `api` |
| `SALESFORCE_LIST_REPORT_RUNS` | List report runs | read | `report_id` |  | `full`, `api` |
| `SALESFORCE_LIST_REPORT_TYPES` | List report types | read |  |  | `full`, `api` |
| `SALESFORCE_LIST_STANDARD_INVOCABLE_ACTIONS` | Get standard invocable actions | read |  |  | `full`, `api` |
| `SALESFORCE_LIST_TOOLING_OBJECTS` | List Tooling API objects | read |  |  | `full`, `api` |
| `SALESFORCE_LOG_CALL` | Log call | write | *call_date, call_disposition, call_duration_seconds, call_type, comments, custom_fields, subject, what_id and 1 more* |  | `full`, `api` |
| `SALESFORCE_LOG_EMAIL_ACTIVITY` | Log email activity | write | `from_address`, `related_to_id`, `subject`, `to_address`; *bcc_address, cc_address, custom_fields, html_body, is_client_managed, is_externally_visible, is_incoming, message_date and 3 more* |  | `full`, `api` |
| `SALESFORCE_MASS_TRANSFER_OWNERSHIP` | Mass transfer ownership | write | `new_owner_id`, `object_type`, `record_ids`; *trigger_assignment_rules* |  | `full`, `api` |
| `SALESFORCE_PARAMETERIZED_SEARCH` | Parameterized Search | read | `q`; *defaultSearchScope, fields, in, method, overallLimit, sobjects, spellCorrection, where* |  | `full`, `api` |
| `SALESFORCE_PATCH_COMPOSITE_SOBJECTS` | Upsert records using external ID | write | `external_id_field_name`, `records`, `sobject_name`; *all_or_none* |  | `full`, `api` |
| `SALESFORCE_POST_COMPOSITE_GRAPH` | Execute Composite Graph | write | `graphs` |  | `full`, `api` |
| `SALESFORCE_POST_COMPOSITE_SOBJECTS` | Create records using sObject Collections | write | `records`; *all_or_none* |  | `full`, `api` |
| `SALESFORCE_POST_PARAMETERIZED_SEARCH` | Parameterized Search | read | `q`; *defaultLimit, fields, in, overallLimit, sobjects, spellCorrection, where* |  | `full`, `api` |
| `SALESFORCE_PROCESS_APPROVALS` | Process Approvals | write | `requests` |  | `full`, `api` |
| `SALESFORCE_QUERY_ALL` | Query All (including deleted) | read | `query` |  | `full`, `api` |
| `SALESFORCE_QUERY_MORE` | Query more (next page of SOQL results) | read | `next_records_url` |  | `full`, `api` |
| `SALESFORCE_REFRESH_DASHBOARD` | Refresh dashboard | write | `dashboard_id`; *filter1, filter2, filter3, running_user* |  | `full`, `api` |
| `SALESFORCE_REMOVE_FROM_CAMPAIGN` | Remove from campaign | write | `campaign_id`; *campaign_member_id, member_id* |  | `full`, `api` |
| `SALESFORCE_REPLACE_FOLDER_SHARES` | Replace folder shares | write | `folder_id`, `shares` |  | `full`, `api` |
| `SALESFORCE_RESET_SELF_SERVICE_USER_PASSWORD` | Reset self-service user password | write | `self_service_user_id` |  | `full`, `api` |
| `SALESFORCE_RETRIEVE_LEAD_BY_ID` | Retrieve lead by id | read | `id`; *fields, filtered_fields_message* |  | `full`, `api` |
| `SALESFORCE_RETRIEVE_OPPORTUNITIES_DATA` | Retrieve opportunities data | read |  |  | `full`, `api` |
| `SALESFORCE_RUN_REPORT` | Run report | read | `report_id` |  | `full`, `api` |
| `SALESFORCE_RUN_REPORT_ASYNC` | Run report async | read | `report_id`; *include_details, report_metadata* |  | `full`, `api` |
| `SALESFORCE_RUN_REPORT_WITH_OVERRIDES` | Run report with overrides | read | `report_id`, `report_metadata`; *include_details* |  | `full`, `api` |
| `SALESFORCE_RUN_SOQL_QUERY` | Run SOQL query | read | `query` |  | `full`, `api` |
| `SALESFORCE_RUN_TESTS_ASYNCHRONOUS` | Run tests asynchronously | write | `tests`; *max_failed_tests, skip_code_coverage* |  | `full`, `api` |
| `SALESFORCE_RUN_TESTS_SYNCHRONOUS` | Run tests synchronously | write | *class_id, class_name, max_failed_tests, skip_code_coverage, test_methods* |  | `full`, `api` |
| `SALESFORCE_SAVE_REPORT` | Save report | write | `report_metadata`; *report_id* |  | `full`, `api` |
| `SALESFORCE_SEARCH_ACCOUNTS` | Search accounts | read | *billing_city, billing_country, billing_state, fields, industry, limit, name, phone and 2 more* |  | `full`, `api` |
| `SALESFORCE_SEARCH_CAMPAIGNS` | Search campaigns | read | *fields, is_active, limit, name, start_date_from, start_date_to, status, type* |  | `full`, `api` |
| `SALESFORCE_SEARCH_CONTACTS` | Search contacts | read | *account_name, email, fields, limit, name, phone, title* |  | `full`, `api` |
| `SALESFORCE_SEARCH_KNOWLEDGE_ARTICLES` | Search suggested article title matches | read | `language`, `publishStatus`, `q`; *channel, limit, topics, validationStatus* |  | `full`, `api` |
| `SALESFORCE_SEARCH_LEADS` | Search leads | read | *company, email, fields, lead_source, limit, name, phone, status and 1 more* |  | `full`, `api` |
| `SALESFORCE_SEARCH_NOTES` | Search notes | write | *body, created_date_from, created_date_to, fields, is_private, limit, owner_name, parent_name and 1 more* |  | `full`, `api` |
| `SALESFORCE_SEARCH_OPPORTUNITIES` | Search opportunities | read | *account_name, amount_max, amount_min, close_date_from, close_date_to, fields, is_closed, is_won and 4 more* |  | `full`, `api` |
| `SALESFORCE_SEARCH_TASKS` | Search tasks | read | *account_name, activity_date_from, activity_date_to, assigned_to_name, contact_name, fields, is_closed, limit and 3 more* |  | `full`, `api` |
| `SALESFORCE_SEARCH_TOOLING` | Search Tooling API records | read | `q` |  | `full`, `api` |
| `SALESFORCE_SEND_EMAIL` | Send email | write | `body`, `subject`, `to_addresses`; *attachment_ids, bcc_addresses, cc_addresses, is_html, log_email, org_wide_email_address_id, recipient_id, related_record_id and 2 more* |  | `full`, `api` |
| `SALESFORCE_SEND_EMAIL_FROM_TEMPLATE` | Send email from template | write | `recipient_id`, `template_id`; *add_threading_tokens, additional_to_addresses, attachment_ids, bcc_addresses, cc_addresses, log_email, related_record_id, sender_address and 1 more* |  | `full`, `api` |
| `SALESFORCE_SEND_MASS_EMAIL` | Send mass email | write | `recipient_ids`; *batch_size, body, is_html, log_emails, sender_address, sender_type, subject, template_id* |  | `full`, `api` |
| `SALESFORCE_SET_SELF_SERVICE_USER_PASSWORD` | Set self-service user password | write | `new_password`, `self_service_user_id` |  | `full`, `api` |
| `SALESFORCE_SET_USER_PASSWORD` | Set user password | write | `user_id`; *new_password* |  | `full`, `api` |
| `SALESFORCE_SOBJECT_ROWS_UPDATE` | Update sObject record | write | `fields`, `record_id`, `sobject_api_name` |  | `api`, `full` |
| `SALESFORCE_SOBJECT_USER_PASSWORD` | Get user password expiration status | read | `user_id` |  | `full`, `api` |
| `SALESFORCE_START_METADATA_DEPLOYMENT` | Start metadata deployment | write | *content_b64, deploy_options, file, filename* |  | `api`, `full` |
| `SALESFORCE_TOOLING_QUERY` | Execute Tooling API query | read | `query` |  | `full`, `api` |
| `SALESFORCE_TRIGGER_PROCESS_RULES` | Trigger process rules | write | `context_ids` |  | `full`, `api` |
| `SALESFORCE_UPDATE_ACCOUNT` | Update account | write | `account_id`; *account_source, annual_revenue, billing_city, billing_country, billing_postal_code, billing_state, billing_street, custom_fields and 15 more* |  | `full`, `api` |
| `SALESFORCE_UPDATE_ANALYTICS_NOTIFICATION` | Update analytics notification | write | `notification_id`; *active, configuration, deactivateOnTrigger, name, recordId, runAs, schedule, source and 1 more* |  | `full`, `api` |
| `SALESFORCE_UPDATE_CAMPAIGN` | Update campaign | write | `campaign_id`; *actual_cost, budgeted_cost, custom_fields, description, end_date, expected_response, expected_revenue, is_active and 6 more* |  | `full`, `api` |
| `SALESFORCE_UPDATE_CONTACT` | Update contact | write | `contact_id`; *account_id, birthdate, custom_fields, department, email, first_name, last_name, mailing_city and 7 more* |  | `full`, `api` |
| `SALESFORCE_UPDATE_CUSTOM_FIELD` | Update custom field metadata | write | `custom_field_id`, `metadata` |  | `full`, `api` |
| `SALESFORCE_UPDATE_DASHBOARD` | Update dashboard | write | `dashboard_id`, `dashboard_metadata` |  | `full`, `api` |
| `SALESFORCE_UPDATE_FAVORITE` | Update a favorite | write | `favorite_id`; *name, sort_order* |  | `full`, `api` |
| `SALESFORCE_UPDATE_FOLDER` | Update folder | write | `folder_id`; *label, name* |  | `full`, `api` |
| `SALESFORCE_UPDATE_FOLDER_SHARE` | Update folder share | write | `access_type`, `folder_id`, `share_id`, `share_type`, `share_with_id` |  | `full`, `api` |
| `SALESFORCE_UPDATE_KNOWLEDGE_ARTICLE_VERSION` | Update Knowledge article version | write | `version_id`; *archive_schedule_date, assignee_id, due_date, instruction, publish_schedule_date, publish_status, version_number* |  | `full`, `api` |
| `SALESFORCE_UPDATE_LEAD` | Update lead | write | `lead_id`; *annual_revenue, city, company, country, custom_fields, description, email, first_name and 12 more* |  | `full`, `api` |
| `SALESFORCE_UPDATE_LIST_VIEW_PREFERENCES` | Update list view preferences | write | `list_view_api_name`, `object_api_name`; *column_order, column_widths, column_wrap* |  | `full`, `api` |
| `SALESFORCE_UPDATE_NOTE` | Update note | write | `note_id`; *body, custom_fields, is_private, owner_id, title* |  | `full`, `api` |
| `SALESFORCE_UPDATE_OPPORTUNITY` | Update opportunity | write | `opportunity_id`; *account_id, amount, close_date, custom_fields, description, lead_source, name, next_step and 3 more* |  | `full`, `api` |
| `SALESFORCE_UPDATE_RECORD` | Update a record | write | `fields`, `record_id`; *allow_save_on_duplicate, api_name, if_unmodified_since, trigger_other_email, trigger_user_email, use_default_rule* |  | `full`, `api` |
| `SALESFORCE_UPDATE_RELATED_LIST_PREFERENCES` | Update Related List Preferences | write | `composite_layout_name`; *column_widths, column_wrap, ordered_by* |  | `full`, `api` |
| `SALESFORCE_UPDATE_S_OBJECT_COLLECTION` | Update sObject collection | write | `records`; *all_or_none* |  | `full`, `api` |
| `SALESFORCE_UPDATE_TASK` | Update task | write | `task_id`; *activity_date, custom_fields, description, is_reminder_set, priority, reminder_date_time, status, subject and 2 more* |  | `full`, `api` |
| `SALESFORCE_UPDATE_TOOLING_FLOW` | Update Tooling Flow | write | `flow_id`, `flow_metadata` |  | `full`, `api` |
| `SALESFORCE_UPDATE_TOOLING_RECORD` | Update Tooling API record | write | `fields`, `object_name`, `record_id` |  | `full`, `api` |
| `SALESFORCE_UPLOAD_FILE` | Upload file to Salesforce Files | write | *content_b64, desc, file, filename, first_publish_location_id, mimetype_override, title* |  | `api`, `chatter_api`, `full` |
| `SALESFORCE_UPLOAD_JOB_DATA` | Upload job data | write | `csv_data`, `job_id` |  | `full`, `api` |
| `SALESFORCE_UPSERT_SOBJECT_BY_EXTERNAL_ID` | Upsert sObject by External ID | write | `field_name`, `field_value`, `fields`, `sobject`; *update_only* |  | `full`, `api` |
| `SALESFORCE_UPSERT_SURVEY_TRANSLATIONS` | Upsert Survey Translations | write | `mode`; *developer_name, language, translations, value* |  | `full`, `api` |
| `SALESFORCE_WHO_AM_I` | Who Am I | read |  |  |  |
