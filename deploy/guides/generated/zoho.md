# Zoho: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

110 tools: 66 read, 44 write (a tool not marked read-only counts as a write). 3 more couldn't be read and are left out.

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `ZOHO_ADD_TAGS_TO_RECORDS` | Add Tags to Zoho CRM Records | write | `module_api_name`, `record_ids`, `tags`; *over_write* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.{module_name}.WRITE`, `ZohoCRM.modules.{module_name}.CREATE` |
| `ZOHO_ASSOCIATE_EMAIL` | Associate Email with Zoho CRM Record | write | `content`, `date_time`, `from_email`, `module_api_name`, `original_message_id`, `record_id`, `sent`, `subject`, `to`; *attachments, bcc, cc, linked_record, mail_format, owner* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.{module_name}.ALL`, `ZohoCRM.modules.emails.ALL` |
| `ZOHO_CONVERT_ZOHO_LEAD` | Convert Zoho CRM Lead | write | `lead_id`; *Deals, account_id, assign_to, carry_over_tags, contact_id, move_attachments_to, notify_lead_owner, notify_new_entity_owner and 1 more* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.leads.CREATE` |
| `ZOHO_COUNT_RECORDS` | Count Zoho CRM Records | read | `module_api_name`; *approved, converted, criteria, custom_view_id, email, phone, word* |  | `ZohoSearch.securesearch.READ`, `ZohoCRM.modules.{module_name}.READ`, `ZohoCRM.settings.layouts.READ`, `ZohoCRM.settings.layouts.ALL`, `ZohoCRM.settings.ALL` |
| `ZOHO_CREATE_AUTOMATION_TASK` | Create Zoho CRM Automation Task | write | `field_mappings`, `module`, `name`; *feature_type, related_module* |  | `ZohoCRM.settings.automation_actions.CREATE`, `ZohoCRM.settings.automation_actions.ALL` |
| `ZOHO_CREATE_CONTACT` | Create Contact in Zoho CRM | write | `last_name`; *account_name, apply_feature_execution, assistant, asst_phone, date_of_birth, department, description, email and 27 more* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.CREATE`, `ZohoCRM.modules.contacts.ALL`, `ZohoCRM.modules.contacts.WRITE`, `ZohoCRM.modules.contacts.CREATE` |
| `ZOHO_CREATE_DEAL` | Create Deal in Zoho CRM | write | `Deal_Name`, `Stage`; *Account_Name, Amount, Campaign_Source, Closing_Date, Contact_Name, Description, Lead_Source, Next_Step and 7 more* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.CREATE`, `ZohoCRM.modules.deals.ALL`, `ZohoCRM.modules.deals.WRITE`, `ZohoCRM.modules.deals.CREATE` |
| `ZOHO_CREATE_EMAIL_DRAFT` | Create Email Draft in Zoho CRM | write | `email_drafts`, `module_api_name`, `record_id` |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.CREATE` |
| `ZOHO_CREATE_EVENT` | Create Zoho Event | write | `End_DateTime`, `Event_Title`, `Start_DateTime`; *All_day, Check_In_Status, Description, Owner, Participants, Recurring_Activity, Remind_At, Venue and 5 more* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.CREATE`, `ZohoCRM.modules.events.ALL`, `ZohoCRM.modules.events.WRITE`, `ZohoCRM.modules.events.CREATE` |
| `ZOHO_CREATE_FUNCTION` | Create Zoho CRM Function | write | `api_name`, `category`, `name`, `runtime`; *code_file, deluge_code, description, publish_changelog* |  | `ZohoCRM.settings.functions.CREATE`, `ZohoCRM.settings.functions.ALL`, `ZohoCRM.settings.CREATE` |
| `ZOHO_CREATE_LEAD` | Create Lead in Zoho CRM | write | `Last_Name`; *Annual_Revenue, City, Company, Country, Description, Designation, Email, Fax and 20 more* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.CREATE`, `ZohoCRM.modules.leads.ALL`, `ZohoCRM.modules.leads.WRITE`, `ZohoCRM.modules.leads.CREATE` |
| `ZOHO_CREATE_NOTE` | Create Note in Zoho CRM | write | `Note_Content`, `module_api_name`, `record_id`; *Note_Title, is_shared_to_client* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.CREATE`, `ZohoCRM.modules.notes.CREATE` |
| `ZOHO_CREATE_WORKFLOW_RULE` | Create Zoho CRM Workflow Rule | write | `workflow_rule` |  | `ZohoCRM.settings.workflow_rules.CREATE`, `ZohoCRM.settings.workflow_rules.ALL` |
| `ZOHO_CREATE_ZOHO_RECORD` | Create Zoho CRM Record | write | `data`, `module_api_name`; *apply_feature_execution, lar_id, skip_feature_execution, trigger* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.CREATE` |
| `ZOHO_CREATE_ZOHO_TAG` | Create Zoho CRM Tag | write | `module_api_name`, `name`; *color_code* |  | `ZohoCRM.settings.ALL`, `ZohoCRM.settings.tags.ALL`, `ZohoCRM.settings.tags.WRITE`, `ZohoCRM.settings.tags.CREATE` |
| `ZOHO_DELETE_ACCOUNT` | Delete Zoho Account | write | `account_id`; *wf_trigger* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.DELETE`, `ZohoCRM.modules.accounts.ALL`, `ZohoCRM.modules.accounts.DELETE` |
| `ZOHO_DELETE_CONTACT` | Delete Zoho CRM Contact | write | `contact_id`; *wf_trigger* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.DELETE`, `ZohoCRM.modules.contacts.ALL`, `ZohoCRM.modules.contacts.DELETE` |
| `ZOHO_DELETE_DEAL` | Delete Zoho CRM Deal | write | `deal_id`; *wf_trigger* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.DELETE`, `ZohoCRM.modules.deals.ALL`, `ZohoCRM.modules.deals.DELETE` |
| `ZOHO_DELETE_EMAIL_DRAFT` | Delete Zoho CRM Email Draft | write | `draft_id`, `module_api_name`, `record_id` |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.{module_name}.ALL`, `ZohoCRM.modules.{module_name}.DELETE` |
| `ZOHO_DELETE_EVENT` | Delete Zoho CRM Event | write | `event_id`; *wf_trigger* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.events.DELETE` |
| `ZOHO_DELETE_LEAD` | Delete Zoho CRM Lead | write | `record_id`; *wf_trigger* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.leads.ALL`, `ZohoCRM.modules.leads.DELETE` |
| `ZOHO_DELETE_NOTE` | Delete Zoho CRM Note | write | `note_id` |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.notes.ALL`, `ZohoCRM.modules.notes.DELETE` |
| `ZOHO_DELETE_RECORDS` | Delete Zoho CRM Records | write | `module_api_name`, `record_ids`; *wf_trigger* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.{module_name}.ALL`, `ZohoCRM.modules.{module_name}.DELETE` |
| `ZOHO_DELETE_TAG` | Delete Zoho CRM Tag | write | `tag_id` |  | `ZohoCRM.settings.ALL`, `ZohoCRM.settings.tags.ALL`, `ZohoCRM.settings.tags.WRITE`, `ZohoCRM.settings.tags.DELETE` |
| `ZOHO_DELETE_WORKFLOW_RULE` | Delete Zoho CRM Workflow Rule | write | `workflow_rule_id` |  | `ZohoCRM.settings.workflow_rules.DELETE`, `ZohoCRM.settings.workflow_rules.ALL` |
| `ZOHO_DOWNLOAD_RECORD_ATTACHMENT` | Download Zoho CRM Record Attachment | read | `attachment_id`, `module_api_name`, `record_id` |  | `ZohoCRM.modules.READ`, `ZohoCRM.modules.{module_name}.READ`, `ZohoCRM.modules.ALL`, `ZohoCRM.modules.{module_name}.ALL`, `ZohoCRM.modules.attachments.READ`, `ZohoCRM.modules.attachments.ALL` |
| `ZOHO_DOWNLOAD_RECORD_EMAIL_ATTACHMENT` | Download Zoho CRM Record Email Attachment | read | `attachment_id`, `message_id`, `module_api_name`, `name`, `record_id`, `user_id` |  | `ZohoCRM.modules.READ`, `ZohoCRM.modules.{module_name}.READ`, `ZohoCRM.modules.emails.READ` |
| `ZOHO_DOWNLOAD_RECORD_PHOTO` | Download Zoho CRM Record Photo | read | `module_api_name`, `record_id` |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.{module_name}.READ` |
| `ZOHO_EXPORT_AUDIT_LOGS` | Export Zoho CRM Audit Logs | read | *criteria* |  | `ZohoCRM.settings.audit_logs.CREATE`, `ZohoCRM.settings.audit_logs.READ`, `ZohoFiles.files.READ` |
| `ZOHO_GET_ACCOUNT` | Get Zoho CRM Account | read | `account_id`; *fields* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.accounts.ALL`, `ZohoCRM.modules.accounts.READ` |
| `ZOHO_GET_ASSIGNMENT_RULE` | Get Zoho CRM Assignment Rule | read | `assignment_rule_id`, `module` |  | `ZohoCRM.settings.assignment_rules.READ` |
| `ZOHO_GET_AUDIT_LOG_EXPORT` | Get Zoho CRM Audit Log Export | read | `export_job_id` |  | `ZohoCRM.settings.audit_logs.READ` |
| `ZOHO_GET_AUTOMATION_TASK` | Get Zoho CRM Automation Task | read | `task_id` |  | `ZohoCRM.settings.automation_actions.READ`, `ZohoCRM.settings.automation_actions.ALL` |
| `ZOHO_GET_BLUEPRINT` | Get Zoho CRM Blueprint | read | `blueprint_id`; *include* |  | `ZohoCRM.settings.blueprint.READ`, `ZohoCRM.settings.blueprint.ALL` |
| `ZOHO_GET_BLUEPRINT_TRANSITION` | Get Zoho CRM Blueprint Transition | read | `transition_id` |  | `ZohoCRM.settings.transitions.read` |
| `ZOHO_GET_CALL` | Get Call | read | `record_id`; *fields* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.calls.ALL`, `ZohoCRM.modules.calls.READ` |
| `ZOHO_GET_CONTACT` | Get Zoho CRM Contact | read | `record_id`; *fields* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.contacts.ALL`, `ZohoCRM.modules.contacts.READ` |
| `ZOHO_GET_CUSTOM_VIEWS` | Get Zoho CRM Custom Views | write | `module_api_name`; *custom_view_id* |  | `ZohoCRM.settings.custom_views.READ`, `ZohoCRM.settings.custom_views.ALL`, `ZohoCRM.settings.ALL` |
| `ZOHO_GET_DEAL` | Get Zoho CRM Deal | read | `deal_id` |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.deals.ALL`, `ZohoCRM.modules.deals.READ` |
| `ZOHO_GET_EMAIL_DRAFTS` | Get Email Drafts for Zoho CRM Record | read | `module`, `record`; *draft_id, page, per_page* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ` |
| `ZOHO_GET_EVENT` | Get Zoho Event | read | `record_id`; *fields* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.events.ALL`, `ZohoCRM.modules.events.READ` |
| `ZOHO_GET_FROM_ADDRESSES` | Get From Addresses | read |  |  | `ZohoCRM.settings.ALL`, `ZohoCRM.settings.emails.READ` |
| `ZOHO_GET_FUNCTION` | Get Zoho CRM Function | read | `function_identifier` |  | `ZohoCRM.settings.functions.READ`, `ZohoCRM.settings.functions.ALL`, `ZohoCRM.settings.ALL` |
| `ZOHO_GET_KIOSK_STATES` | Get Zoho CRM Kiosk States | read | `kiosk_id`; *include_configured_properties, state_id* |  | `ZohoCRM.settings.kiosks.READ`, `ZohoCRM.settings.kiosks.ALL` |
| `ZOHO_GET_LAYOUTS` | Get Zoho CRM Layouts | read | `module_api_name`; *layout_id* |  | `ZohoCRM.settings.layouts.READ`, `ZohoCRM.settings.layouts.ALL`, `ZohoCRM.settings.ALL` |
| `ZOHO_GET_MODULE_FIELDS` | Get Zoho CRM Module Fields Metadata | read | `module_api_name`; *field_unique_id, include, type* |  | `ZohoCRM.settings.ALL`, `ZohoCRM.settings.fields.ALL`, `ZohoCRM.settings.fields.READ` |
| `ZOHO_GET_NOTE` | Get Zoho CRM Note | read | `note_id`; *fields* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.notes.ALL`, `ZohoCRM.modules.notes.READ` |
| `ZOHO_GET_ORGANIZATION_DETAILS` | Get Zoho CRM Organization Details | read |  |  | `ZohoCRM.org.READ`, `ZohoCRM.org.ALL` |
| `ZOHO_GET_PIPELINES` | Get Zoho CRM Pipelines | read | `layout_id` |  | `ZohoCRM.settings.pipeline.READ` |
| `ZOHO_GET_PRODUCT` | Get Zoho CRM Product | read | `record_id`; *fields* |  | `ZohoCRM.modules.products.READ`, `ZohoCRM.modules.products.ALL`, `ZohoCRM.modules.ALL` |
| `ZOHO_GET_PROFILE` | Get Zoho CRM Profile | read | `profile_id` |  | `ZohoCRM.settings.profiles.READ`, `ZohoCRM.settings.profiles.ALL` |
| `ZOHO_GET_RECORD` | Get Zoho CRM Record | read | `module_api_name`, `record_id`; *fields* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ` |
| `ZOHO_GET_RECORD_EMAIL` | Get Zoho CRM Record Email | read | `message_id`, `module_api_name`, `record_id`; *user_id* |  | `ZohoCRM.modules.READ`, `ZohoCRM.modules.{module_name}.READ`, `ZohoCRM.modules.emails.READ` |
| `ZOHO_GET_RECORD_EMAILS` | Get Zoho CRM Record Emails | read | `module_api_name`, `record_id`; *index, owner_id, type* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.emails.READ` |
| `ZOHO_GET_RECORD_TIMELINE` | Get Zoho CRM Record Timeline | read | `module_api_name`, `record_id`; *filters, include, include_inner_details, include_timeline_type, page_token, per_page, sort_by* |  | `ZohoCRM.modules.{module_name}.ALL`, `ZohoCRM.modules.{module_name}.READ` |
| `ZOHO_GET_RELATED_LISTS` | Get Zoho CRM Related Lists Metadata | read | `module_api_name`; *layout_id* |  | `ZohoCRM.settings.ALL`, `ZohoCRM.settings.related_lists.ALL`, `ZohoCRM.settings.related_lists.READ` |
| `ZOHO_GET_RELATED_RECORDS` | Get Zoho CRM Related Records | read | `fields`, `module_api_name`, `record_id`, `related_list_api_name`; *converted, ids, page, page_token, per_page, sort_by, sort_order* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ` |
| `ZOHO_GET_TASK` | Get Zoho CRM Task | read | `record_id`; *fields* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.tasks.ALL`, `ZohoCRM.modules.tasks.READ` |
| `ZOHO_GET_USER` | Get Zoho CRM User | read | `user_id` |  | `ZohoCRM.users.ALL`, `ZohoCRM.users.READ` |
| `ZOHO_GET_WEBFORM` | Get Zoho CRM Webform | read | `module_api_name`, `webform_id`; *from_time, include_active_user_details, to_time* |  | `ZohoCRM.settings.webforms.READ`, `ZohoCRM.settings.webforms.ALL` |
| `ZOHO_GET_WORKFLOW_RULE` | Get Zoho CRM Workflow Rule | read | `workflow_rule_id` |  | `ZohoCRM.settings.workflow_rules.READ`, `ZohoCRM.settings.workflow_rules.ALL` |
| `ZOHO_GET_ZOHO_RECORDS` | Get Zoho CRM Records | read | `fields`, `module_api_name`; *converted, cvid, ids, include_child, page, page_token, per_page, sort_by and 2 more* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ` |
| `ZOHO_GET_ZOHO_USERS` | Get Zoho CRM Users | read | *ids, if_modified_since, page, per_page, type, type__s* |  | `ZohoCRM.users.ALL`, `ZohoCRM.users.READ` |
| `ZOHO_LIST_ACCESS_SETTINGS` | List Zoho CRM Access Settings | read | `collection` |  | `ZohoCRM.settings.roles.READ`, `ZohoCRM.settings.roles.ALL`, `ZohoCRM.settings.profiles.READ`, `ZohoCRM.settings.profiles.ALL` |
| `ZOHO_LIST_ACCOUNTS` | List Zoho CRM Accounts | read | `fields`; *approved, cvid, ids, include_child, page, page_token, per_page, sort_by and 2 more* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.accounts.ALL`, `ZohoCRM.modules.accounts.READ` |
| `ZOHO_LIST_CADENCES` | List Zoho CRM Cadences | read | *active, fields, module_api_name, name, page, per_page* |  | `ZohoCRM.settings.cadences.READ`, `ZohoCRM.settings.cadences.ALL`, `ZohoCRM.settings.ALL` |
| `ZOHO_LIST_CALLS` | List Zoho CRM Calls | read | `fields`; *cvid, ids, page, page_token, per_page, sort_by, sort_order* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.calls.ALL`, `ZohoCRM.modules.calls.READ` |
| `ZOHO_LIST_CONNECTED_WORKFLOW_RULES` | List Zoho CRM Connected Workflow Rules | read | `connected_workflow_id`, `module_api_name` |  | `ZohoCRM.settings.connected_workflows.READ`, `ZohoCRM.settings.connected_workflows.ALL`, `ZohoCRM.settings.ALL` |
| `ZOHO_LIST_CONTACTS` | List Contacts | read | *cvid, fields, ids, include_child, page, page_token, per_page, sort_by and 2 more* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.contacts.ALL`, `ZohoCRM.modules.contacts.READ` |
| `ZOHO_LIST_DEALS` | List Deals | read | *approval_state, converted, cvid, fields, ids, include_child, page, page_token and 4 more* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.deals.ALL`, `ZohoCRM.modules.deals.READ` |
| `ZOHO_LIST_EVENTS` | List Zoho CRM Events | read | `fields`; *cvid, ids, page, page_token, per_page, sort_by, sort_order* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.events.ALL`, `ZohoCRM.modules.events.READ` |
| `ZOHO_LIST_FUNCTIONS` | List Zoho CRM Functions | read | *page, page_token, per_page* |  | `ZohoCRM.settings.functions.READ`, `ZohoCRM.settings.functions.ALL`, `ZohoCRM.settings.ALL` |
| `ZOHO_LIST_LEADS` | List Zoho CRM Leads | read | *converted, cvid, fields, ids, page, page_token, per_page, sort_by and 1 more* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.leads.ALL`, `ZohoCRM.modules.leads.READ` |
| `ZOHO_LIST_MODULES` | List Zoho CRM Modules | read | *status* |  | `ZohoCRM.settings.ALL`, `ZohoCRM.settings.modules.ALL`, `ZohoCRM.settings.modules.READ` |
| `ZOHO_LIST_NOTES` | List Zoho CRM Notes | read | `fields`; *page, per_page* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.notes.ALL`, `ZohoCRM.modules.notes.READ` |
| `ZOHO_LIST_RECORD_ATTACHMENTS` | List Attachments for Zoho CRM Record | read | `module_api_name`, `record_id`; *fields, page, per_page* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.attachments.READ` |
| `ZOHO_LIST_TAGS` | List Zoho CRM Tags | read | `module_api_name`; *my_tags* |  | `ZohoCRM.settings.tags.READ`, `ZohoCRM.settings.tags.ALL`, `ZohoCRM.settings.ALL` |
| `ZOHO_LIST_TASKS` | List Zoho CRM Tasks | read | *cvid, fields, ids, page, page_token, per_page, sort_by, sort_order* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.tasks.ALL`, `ZohoCRM.modules.tasks.READ` |
| `ZOHO_LIST_VARIABLES` | List Zoho CRM Variables | read |  |  | `ZohoCRM.settings.variables.READ`, `ZohoCRM.settings.variables.ALL` |
| `ZOHO_LIST_WORKFLOW_RULES` | List Zoho CRM Workflow Rules | read | *execute_on, include_inner_details, module, name_contains, page, per_page, sort_order, status_active and 1 more* |  | `ZohoCRM.settings.workflow_rules.READ`, `ZohoCRM.settings.workflow_rules.ALL` |
| `ZOHO_MERGE_LEADS` | Merge Zoho CRM Leads | write | `child_records`, `master_record_id`; *master_record_fields* |  | `ZohoCRM.modules.leads.ALL`, `ZohoCRM.modules.leads.WRITE` |
| `ZOHO_QUERY_RECORDS_WITH_COQL` | Query Zoho CRM Records with COQL | read | `select_query`; *include_meta* |  | `ZohoCRM.coql.READ`, `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.settings.fields.READ` |
| `ZOHO_SEARCH_ACCOUNTS` | Search Zoho CRM Accounts | read | *approved, criteria, fields, page, per_page, phone, word* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.accounts.ALL`, `ZohoCRM.modules.accounts.READ`, `ZohoSearch.securesearch.READ` |
| `ZOHO_SEARCH_CALLS` | Search Zoho CRM Calls | read | *approved, criteria, fields, page, per_page, word* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.calls.ALL`, `ZohoCRM.modules.calls.READ`, `ZohoSearch.securesearch.READ` |
| `ZOHO_SEARCH_CONTACTS` | Search Zoho CRM Contacts | read | *approved, criteria, email, fields, page, per_page, phone, word* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.contacts.ALL`, `ZohoCRM.modules.contacts.READ`, `ZohoSearch.securesearch.READ` |
| `ZOHO_SEARCH_DEALS` | Search Deals in Zoho CRM | read | *approved, criteria, fields, page, per_page, word* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.deals.ALL`, `ZohoCRM.modules.deals.READ`, `ZohoSearch.securesearch.READ` |
| `ZOHO_SEARCH_EVENTS` | Search Events | read | *approved, criteria, fields, page, per_page, word* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.events.ALL`, `ZohoCRM.modules.events.READ`, `ZohoSearch.securesearch.READ` |
| `ZOHO_SEARCH_LEADS` | Search Leads | read | *approved, converted, criteria, email, fields, page, per_page, phone and 1 more* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.leads.ALL`, `ZohoCRM.modules.leads.READ`, `ZohoSearch.securesearch.READ` |
| `ZOHO_SEARCH_NOTES` | Search Notes | read | *approved, criteria, fields, page, per_page, word* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.notes.ALL`, `ZohoCRM.modules.notes.READ`, `ZohoSearch.securesearch.READ` |
| `ZOHO_SEARCH_TASKS` | Search Zoho CRM Tasks | read | *approved, criteria, fields, page, per_page, word* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoCRM.modules.tasks.ALL`, `ZohoCRM.modules.tasks.READ`, `ZohoSearch.securesearch.READ` |
| `ZOHO_SEARCH_ZOHO_RECORDS` | Search Zoho CRM Records | read | `module_api_name`; *approved, converted, criteria, email, fields, page, per_page, phone and 1 more* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.READ`, `ZohoSearch.securesearch.READ` |
| `ZOHO_UPDATE_ACCOUNT` | Update Zoho CRM Account | write | `record_id`; *Account_Name, Account_Number, Account_Site, Account_Type, Annual_Revenue, Billing_City, Billing_Code, Billing_Country and 22 more* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.UPDATE`, `ZohoCRM.modules.accounts.ALL`, `ZohoCRM.modules.accounts.WRITE`, `ZohoCRM.modules.accounts.UPDATE` |
| `ZOHO_UPDATE_CALL` | Update Zoho CRM Call | write | `data`; *apply_feature_execution, skip_feature_execution, trigger* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.UPDATE`, `ZohoCRM.modules.calls.ALL`, `ZohoCRM.modules.calls.WRITE`, `ZohoCRM.modules.calls.UPDATE` |
| `ZOHO_UPDATE_CUSTOM_FIELDS` | Update Zoho CRM Custom Fields | write | `fields`, `module_api_name`; *field_id* |  | `ZohoCRM.settings.fields.UPDATE`, `ZohoCRM.settings.fields.ALL`, `ZohoCRM.settings.ALL` |
| `ZOHO_UPDATE_DEAL` | Update Zoho CRM Deal | write | `data`, `deal_id`; *apply_feature_execution, skip_feature_execution, trigger* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.UPDATE`, `ZohoCRM.modules.deals.ALL`, `ZohoCRM.modules.deals.WRITE`, `ZohoCRM.modules.deals.UPDATE` |
| `ZOHO_UPDATE_EMAIL_DRAFT` | Update Zoho CRM Email Draft | write | `email_drafts`, `module_api_name`, `record_id` |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.UPDATE` |
| `ZOHO_UPDATE_EVENT` | Update Zoho CRM Event | write | `data`; *apply_feature_execution, skip_feature_execution, trigger* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.UPDATE`, `ZohoCRM.modules.events.ALL`, `ZohoCRM.modules.events.WRITE`, `ZohoCRM.modules.events.UPDATE` |
| `ZOHO_UPDATE_FUNCTION` | Update Zoho CRM Function | write | `function_identifier`; *code_file, deluge_code, description, name, publish_changelog* |  | `ZohoCRM.settings.functions.UPDATE`, `ZohoCRM.settings.functions.ALL`, `ZohoCRM.settings.UPDATE` |
| `ZOHO_UPDATE_LAYOUT` | Update Zoho CRM Layout | write | `layout`, `layout_id`, `module_api_name` |  | `ZohoCRM.settings.layouts.UPDATE`, `ZohoCRM.settings.layouts.ALL`, `ZohoCRM.settings.ALL` |
| `ZOHO_UPDATE_LEAD` | Update Lead | write | `data`; *apply_feature_execution, skip_feature_execution, trigger* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.UPDATE`, `ZohoCRM.modules.leads.ALL`, `ZohoCRM.modules.leads.WRITE`, `ZohoCRM.modules.leads.UPDATE` |
| `ZOHO_UPDATE_NOTE` | Update Zoho CRM Note | write | `id`; *Note_Content, Note_Title, is_shared_to_client* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.UPDATE`, `ZohoCRM.modules.notes.ALL`, `ZohoCRM.modules.notes.WRITE`, `ZohoCRM.modules.notes.UPDATE` |
| `ZOHO_UPDATE_RELATED_RECORDS` | Update Related Records in Zoho CRM | write | `data`, `module_api_name`, `record_id`, `related_list_api_name` |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.UPDATE` |
| `ZOHO_UPDATE_USER` | Update Zoho CRM User | write | `changes`, `user_id` |  | `ZohoCRM.users.UPDATE`, `ZohoCRM.users.ALL` |
| `ZOHO_UPDATE_WEBFORM` | Update Zoho CRM Webform | write | `module_api_name`, `webform`, `webform_id` |  | `ZohoCRM.settings.webforms.UPDATE`, `ZohoCRM.settings.webforms.ALL` |
| `ZOHO_UPDATE_WORKFLOW_RULE` | Update Zoho CRM Workflow Rule | write | `rule_id`, `workflow_rule` |  | `ZohoCRM.settings.workflow_rules.UPDATE`, `ZohoCRM.settings.workflow_rules.ALL` |
| `ZOHO_UPDATE_ZOHO_RECORD` | Update Zoho CRM Record | write | `data`, `module_api_name`; *apply_feature_execution, skip_feature_execution, trigger* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.UPDATE` |
| `ZOHO_UPLOAD_ATTACHMENT` | Upload Attachment to Zoho CRM Record | write | `module_api_name`, `record_id`; *attachmentUrl, file, title* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.CREATE`, `ZohoCRM.modules.attachments.CREATE` |
| `ZOHO_UPLOAD_RECORD_PHOTO` | Upload Zoho CRM Record Photo | write | `file`, `module_api_name`, `record_id`; *restrict_triggers* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.{module_name}.ALL`, `ZohoCRM.modules.{module_name}.CREATE`, `ZohoCRM.modules.{module_name}.WRITE` |
| `ZOHO_UPSERT_RECORDS` | Upsert Zoho CRM Records | write | `data`, `module_api_name`; *apply_feature_execution, duplicate_check_fields, if_unmodified_since, skip_feature_execution, trigger* |  | `ZohoCRM.modules.ALL`, `ZohoCRM.modules.{module_name}.ALL`, `ZohoCRM.modules.{module_name}.WRITE`, `ZohoCRM.modules.{module_name}.CREATE` |
| `ZOHO_VALIDATE_CREDENTIAL` | Validate Zoho CRM Credentials | read |  |  | `ZohoCRM.users.ALL`, `ZohoCRM.users.READ` |
