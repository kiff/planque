# Dynamics 365: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

119 tools: 47 read, 72 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `DYNAMICS365_ADD_CONTACT_TO_ACCOUNT` | Add Contact to Account | write | `account_id`, `contact_id` |  | `user_impersonation` |
| `DYNAMICS365_ASSIGN_USER_SECURITY_ROLE` | Assign User Security Role | write | `role_id`, `user_id` |  | `user_impersonation` |
| `DYNAMICS365_ASSOCIATE_CONTACT_ACCOUNT` | Associate Contact Account | write | `account_id`, `contact_id` |  | `user_impersonation` |
| `DYNAMICS365_CLEAR_ACCOUNT_DESCRIPTION` | Clear Account Description | write | `account_id` |  | `user_impersonation` |
| `DYNAMICS365_CREATE_ANNOTATION` | Create Annotation | write | `subject`; *document_body, file_name, mime_type, note_text, parent_id, parent_type* |  | `user_impersonation` |
| `DYNAMICS365_CREATE_APPOINTMENT` | Create appointment | write | `properties` |  | `user_impersonation` |
| `DYNAMICS365_CREATE_EMAIL_ACTIVITY` | Create Email Activity | write | *description, is_outgoing, priority, subject* |  | `user_impersonation` |
| `DYNAMICS365_CREATE_INVOICE_LINE_ITEM` | Create invoice line item | write | `properties` |  | `user_impersonation` |
| `DYNAMICS365_CREATE_OPPORTUNITY_LINE_ITEM` | Create opportunity line item | write | `properties` |  | `user_impersonation` |
| `DYNAMICS365_CREATE_PHONE_CALL` | Create Phone Call | write | `subject`; *description, direction, phone_number, scheduled_end, scheduled_start* |  | `user_impersonation` |
| `DYNAMICS365_CREATE_PRICE_LIST` | Create price list | write | `properties` |  | `user_impersonation` |
| `DYNAMICS365_CREATE_PRODUCT` | Create Product | write | `default_uom_id`, `default_uom_schedule_id`, `name`, `product_number`, `quantity_decimal`; *description, price, select* |  | `user_impersonation` |
| `DYNAMICS365_CREATE_PRODUCT_PRICE_LIST_ITEM` | Create Product Price List Item | write | `amount`, `price_list_id`, `product_id`, `unit_id`; *percentage, pricing_method, quantity_selling_option, rounding_option_amount, rounding_option_code, rounding_policy* |  | `user_impersonation` |
| `DYNAMICS365_CREATE_QUOTE` | Create quote | write | `properties` |  | `user_impersonation` |
| `DYNAMICS365_CREATE_QUOTE_LINE_ITEM` | Create quote line item | write | `properties` |  | `user_impersonation` |
| `DYNAMICS365_CREATE_SALES_ORDER_LINE_ITEM` | Create sales order line item | write | `properties` |  | `user_impersonation` |
| `DYNAMICS365_CREATE_SOCIAL_ACTIVITY` | Create Social Activity | write | *description, is_outgoing, message_type, post_id, post_url, priority, subject, thread_id* |  | `user_impersonation` |
| `DYNAMICS365_CREATE_SYSTEM_USER` | Create system user | write | `properties` |  | `user_impersonation` |
| `DYNAMICS365_CREATE_TASK` | Create Task | write | `subject`; *description, percent_complete, priority_code, scheduled_end, scheduled_start, select* |  | `user_impersonation` |
| `DYNAMICS365_CREATE_TRANSACTION_CURRENCY` | Create transaction currency | write | `properties` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_ACCOUNT` | Delete account | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_APPOINTMENT` | Delete appointment | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_CASE` | Delete case | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_CONTACT` | Delete contact | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_EMAIL_ACTIVITY` | Delete email activity | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_INVOICE` | Delete invoice | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_INVOICE_LINE_ITEM` | Delete invoice line item | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_LEAD` | Delete lead | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_NOTE` | Delete note | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_OPPORTUNITY` | Delete opportunity | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_OPPORTUNITY_LINE_ITEM` | Delete opportunity line item | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_PRICE_LIST` | Delete price list | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_PRODUCT` | Delete product | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_QUOTE` | Delete quote | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_QUOTE_LINE_ITEM` | Delete quote line item | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_SALES_ORDER` | Delete sales order | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_SALES_ORDER_LINE_ITEM` | Delete sales order line item | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_SYSTEM_USER` | Delete system user | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_TASK` | Delete task | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DELETE_TRANSACTION_CURRENCY` | Delete transaction currency | write | `record_id` |  | `user_impersonation` |
| `DYNAMICS365_DISASSOCIATE_CONTACT_ACCOUNT` | Disassociate Contact Account | write | `contact_id` |  | `user_impersonation` |
| `DYNAMICS365_DYNAMICSCRM_CREATE_ACCOUNT` | Create Account | write | *accountcategorycode, creditonhold, description, name, revenue, user_id* |  | `https://{{subdomain}}.api.crm.dynamics.com/.default` |
| `DYNAMICS365_DYNAMICSCRM_CREATE_CASE` | Create Case | write | *account_id, caseorigincode, contact_id, description, prioritycode, title, user_id* |  | `https://{{subdomain}}.api.crm.dynamics.com/.default` |
| `DYNAMICS365_DYNAMICSCRM_CREATE_CONTACT` | Create Contact | write | *address1_city, address1_country, address1_postalcode, address1_stateorprovince, emailaddress1, firstname, jobtitle, lastname and 3 more* |  | `https://{{subdomain}}.api.crm.dynamics.com/.default`, `user_impersonation` |
| `DYNAMICS365_DYNAMICSCRM_CREATE_INVOICE` | Create Invoice | write | *account_id, currency_id, description, name, price_level_id, user_id* |  | `https://{{subdomain}}.api.crm.dynamics.com/.default` |
| `DYNAMICS365_DYNAMICSCRM_CREATE_LEAD` | Create Lead | write | *companyname, emailaddress1, firstname, lastname, subject, telephone1, user_id* |  | `https://{{subdomain}}.api.crm.dynamics.com/.default`, `user_impersonation` |
| `DYNAMICS365_DYNAMICSCRM_CREATE_OPPORTUNITY` | Create Opportunity | write | `name`; *customer_account_id, customer_contact_id, description, estimatedclosedate, estimatedvalue* |  | `https://{{subdomain}}.api.crm.dynamics.com/.default` |
| `DYNAMICS365_DYNAMICSCRM_CREATE_SALES_ORDER` | Create Sales Order | write | *account_id, currency_id, description, name, price_level_id, user_id* |  | `https://{{subdomain}}.api.crm.dynamics.com/.default` |
| `DYNAMICS365_DYNAMICSCRM_GET_ALL_LEADS` | Dynamicscrm get all leads | read | *filter, orderby, select, skip, top* |  | `user_impersonation` |
| `DYNAMICS365_DYNAMICSCRM_GET_A_INVOICE` | Dynamicscrm get a invoice | read | `invoice_id` |  | `https://{{subdomain}}.api.crm.dynamics.com/.default` |
| `DYNAMICS365_DYNAMICSCRM_GET_A_LEAD` | Dynamicscrm get a lead | read | `lead_id` |  | `https://{{subdomain}}.crm.dynamics.com/.default` |
| `DYNAMICS365_DYNAMICSCRM_UPDATE_CASE` | Update Case | write | `case_id`; *description, prioritycode, statuscode, title, user_id* |  | `https://{{subdomain}}.crm.dynamics.com/.default` |
| `DYNAMICS365_DYNAMICSCRM_UPDATE_INVOICE` | Update Invoice | write | `invoice_id`; *description, duedate, name, totalamount, user_id* |  | `https://{{subdomain}}.api.crm.dynamics.com/user_impersonation` |
| `DYNAMICS365_DYNAMICSCRM_UPDATE_LEAD` | Update Lead | write | `lead_id`; *companyname, emailaddress1, firstname, lastname, subject, telephone1, user_id* |  | `user_impersonation` |
| `DYNAMICS365_DYNAMICSCRM_UPDATE_OPPORTUNITY` | Update Opportunity | write | `opportunity_id`; *description, estimatedclosedate, estimatedvalue, name, opportunityratingcode, salesstagecode, user_id* |  | `user_impersonation` |
| `DYNAMICS365_DYNAMICSCRM_UPDATE_SALES_ORDER` | Update Sales Order | write | `salesorder_id`; *description, discountamount, freightamount, name, user_id* |  | `https://{{subdomain}}.api.crm.dynamics.com/.default` |
| `DYNAMICS365_EXECUTE_BATCH` | Execute Batch | write | `body`, `boundary`; *continue_on_error* |  | `user_impersonation` |
| `DYNAMICS365_GET_ACCOUNT` | Get account | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_ALL_INVOICES_ACTION` | Get all invoices action | read | *filter, orderby, select, skip, top* |  | `user_impersonation` |
| `DYNAMICS365_GET_APPOINTMENT` | Get appointment | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_ATTRIBUTE_METADATA` | Get Attribute Metadata | read | `attribute_metadata_id`, `entity_metadata_id`; *attribute_type, expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_CASE` | Get case | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_CONTACT` | Get contact | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_EMAIL_ACTIVITY` | Get email activity | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_INVOICE_LINE_ITEM` | Get invoice line item | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_NOTE` | Get note | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_OPPORTUNITY` | Get opportunity | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_OPPORTUNITY_LINE_ITEM` | Get opportunity line item | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_PRICE_LIST` | Get price list | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_PRODUCT` | Get product | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_QUOTE` | Get quote | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_QUOTE_LINE_ITEM` | Get quote line item | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_SALES_ORDER` | Get sales order | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_SALES_ORDER_LINE_ITEM` | Get sales order line item | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_SYSTEM_USER` | Get system user | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_TASK` | Get task | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_GET_TRANSACTION_CURRENCY` | Get transaction currency | read | `record_id`; *expand, select* |  | `user_impersonation` |
| `DYNAMICS365_LIST_APPOINTMENTS` | List appointments | read | *count, expand, filter, next_link, orderby, page_size, select, top* |  | `user_impersonation` |
| `DYNAMICS365_LIST_ENTITY_DEFINITIONS` | List Entity Definitions | read | *filter, label_languages, select* |  | `user_impersonation` |
| `DYNAMICS365_LIST_ENTITY_SETS` | List entity sets | read |  |  | `user_impersonation` |
| `DYNAMICS365_LIST_INVOICE_LINE_ITEMS` | List invoice line items | read | *count, expand, filter, next_link, orderby, page_size, select, top* |  | `user_impersonation` |
| `DYNAMICS365_LIST_NOTES` | List notes | read | *count, expand, filter, next_link, orderby, page_size, select, top* |  | `user_impersonation` |
| `DYNAMICS365_LIST_OPPORTUNITY_LINE_ITEMS` | List opportunity line items | read | *count, expand, filter, next_link, orderby, page_size, select, top* |  | `user_impersonation` |
| `DYNAMICS365_LIST_PRICE_LISTS` | List price lists | read | *count, expand, filter, next_link, orderby, page_size, select, top* |  | `user_impersonation` |
| `DYNAMICS365_LIST_PRODUCTS` | List products | read | *count, expand, filter, next_link, orderby, page_size, select, top* |  | `user_impersonation` |
| `DYNAMICS365_LIST_QUOTES` | List quotes | read | *count, expand, filter, next_link, orderby, page_size, select, top* |  | `user_impersonation` |
| `DYNAMICS365_LIST_QUOTE_LINE_ITEMS` | List quote line items | read | *count, expand, filter, next_link, orderby, page_size, select, top* |  | `user_impersonation` |
| `DYNAMICS365_LIST_SALES_ORDERS` | List sales orders | read | *count, expand, filter, next_link, orderby, page_size, select, top* |  | `user_impersonation` |
| `DYNAMICS365_LIST_SYSTEM_USERS` | List system users | read | *count, expand, filter, next_link, orderby, page_size, select, top* |  | `user_impersonation` |
| `DYNAMICS365_LIST_TRANSACTION_CURRENCIES` | List transaction currencies | read | *count, expand, filter, next_link, orderby, page_size, select, top* |  | `user_impersonation` |
| `DYNAMICS365_QUERY_ACCOUNTS` | Query accounts | read | *apply, expand, filter, include_count, orderby, select, top* |  | `user_impersonation` |
| `DYNAMICS365_QUERY_CAMPAIGNS` | Query campaigns | read | *apply, expand, filter, include_count, orderby, select, top* |  | `user_impersonation` |
| `DYNAMICS365_QUERY_CASES` | Query cases | read | *filter, next_link, orderby, select, top* |  | `user_impersonation` |
| `DYNAMICS365_QUERY_CONTACTS` | Query contacts | read | *apply, expand, filter, include_count, orderby, select, top* |  | `user_impersonation` |
| `DYNAMICS365_QUERY_EMAIL_ACTIVITIES` | Query email activities | read | *filter, next_link, orderby, select, top* |  | `user_impersonation` |
| `DYNAMICS365_QUERY_ENTITY_DATA` | Query Entity Data | read | `entity_set`; *filter, orderby, select, skiptoken, top* |  | `user_impersonation` |
| `DYNAMICS365_QUERY_OPPORTUNITIES` | Query opportunities | read | *filter, next_link, orderby, select, top* |  | `user_impersonation` |
| `DYNAMICS365_QUERY_PHONE_CALLS` | Query phone calls | read | *filter, next_link, orderby, select, top* |  | `user_impersonation` |
| `DYNAMICS365_QUERY_SOCIAL_ACTIVITIES` | Query social activities | read | *filter, next_link, orderby, select, top* |  | `user_impersonation` |
| `DYNAMICS365_QUERY_TASKS` | Query tasks | read | *filter, next_link, orderby, select, top* |  | `user_impersonation` |
| `DYNAMICS365_REMOVE_CONTACT_FROM_ACCOUNT` | Remove Contact from Account | write | `account_id`, `contact_id` |  | `user_impersonation` |
| `DYNAMICS365_RESOLVE_CASE` | Resolve Case | write | `case_id`, `subject`; *description, status* |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_ACCOUNT` | Update account | write | `properties`, `record_id` |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_ACCOUNT_NAME` | Update Account Name | write | `account_id`, `name` |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_APPOINTMENT` | Update appointment | write | `properties`, `record_id` |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_CONTACT` | Update contact | write | `properties`, `record_id` |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_EMAIL_ACTIVITY` | Update email activity | write | `properties`, `record_id` |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_INVOICE_LINE_ITEM` | Update invoice line item | write | `properties`, `record_id` |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_NOTE` | Update note | write | `properties`, `record_id` |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_OPPORTUNITY_LINE_ITEM` | Update opportunity line item | write | `properties`, `record_id` |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_PRICE_LIST` | Update price list | write | `properties`, `record_id` |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_PRODUCT` | Update product | write | `properties`, `record_id` |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_QUOTE` | Update quote | write | `properties`, `record_id` |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_QUOTE_LINE_ITEM` | Update quote line item | write | `properties`, `record_id` |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_SALES_ORDER_LINE_ITEM` | Update sales order line item | write | `properties`, `record_id` |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_SYSTEM_USER` | Update system user | write | `properties`, `record_id` |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_TASK` | Update task | write | `properties`, `record_id` |  | `user_impersonation` |
| `DYNAMICS365_UPDATE_TRANSACTION_CURRENCY` | Update transaction currency | write | `properties`, `record_id` |  | `user_impersonation` |
| `DYNAMICS365_WHO_AM_I` | Get Caller Identity | read |  |  | `user_impersonation` |
