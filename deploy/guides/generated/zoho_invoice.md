# Zoho Invoice: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

137 tools: 55 read, 82 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `ZOHO_INVOICE_ADD_CREDITNOTE_TO_INVOICES` | Add Credit Note to Invoices | write | `creditnote_id`, `invoices`, `organization_id` |  | `ZohoInvoice.creditnotes.CREATE`, `ZohoInvoice.creditnotes.UPDATE` |
| `ZOHO_INVOICE_ADD_INVOICE_COMMENT` | Add Invoice Comment | write | `description`, `invoice_id`, `organization_id`; *show_comment_to_client* |  | `ZohoInvoice.invoices.CREATE`, `ZohoInvoice.fullaccess.all` |
| `ZOHO_INVOICE_APPLY_CREDITS_TO_INVOICE` | Apply Credits to Invoice | write | `apply_creditnotes`, `invoice_id`, `organization_id` |  | `ZohoInvoice.invoices.CREATE` |
| `ZOHO_INVOICE_CANCEL_WRITE_OFF_INVOICE` | Cancel Write Off Invoice | write | `invoice_id`, `organization_id` |  | `ZohoInvoice.invoices.CREATE` |
| `ZOHO_INVOICE_CLONE_PROJECT` | Clone Zoho Invoice Project | write | `organization_id`, `project_id`, `project_name`; *description* |  | `ZohoInvoice.projects.CREATE` |
| `ZOHO_INVOICE_CREATE_ADDITIONAL_ADDRESS` | Create Additional Address | write | `contact_id`, `organization_id`; *address, attention, city, country, fax, phone, state, state_code and 2 more* |  | `ZohoInvoice.contacts.CREATE` |
| `ZOHO_INVOICE_CREATE_CONTACT` | Create Contact | write | `contact_name`, `organization_id`; *billing_address, company_name, contact_persons, country_code, currency_id, custom_fields, facebook, gst_no and 17 more* |  | `ZohoInvoice.contacts.CREATE` |
| `ZOHO_INVOICE_CREATE_CONTACT_PERSON` | Create Contact Person | write | `contact_id`, `first_name`, `organization_id`; *department, designation, email, enable_portal, last_name, mobile, phone, salutation and 1 more* |  | `ZohoInvoice.contacts.CREATE` |
| `ZOHO_INVOICE_CREATE_CREDIT_NOTE` | Create Credit Note | write | `customer_id`, `date`, `line_items`, `organization_id`; *creditnote_number, notes, reference_number, terms* |  | `ZohoInvoice.creditnotes.CREATE`, `ZohoInvoice.fullaccess.all` |
| `ZOHO_INVOICE_CREATE_CREDIT_NOTE_COMMENT` | Create Credit Note Comment | write | `creditnote_id`, `description`, `organization_id` |  | `ZohoInvoice.creditnotes.CREATE` |
| `ZOHO_INVOICE_CREATE_CURRENCY` | Create Currency | write | `currency_code`, `currency_format`, `currency_symbol`, `organization_id`; *price_precision* |  | `ZohoInvoice.settings.CREATE` |
| `ZOHO_INVOICE_CREATE_CUSTOMER_PAYMENT` | Create Customer Payment | write | `amount`, `customer_id`, `date`, `invoices`, `organization_id`, `payment_mode`; *bank_charges, description, reference_number* |  | `ZohoInvoice.customerpayments.CREATE`, `ZohoInvoice.customerpayments.READ`, `ZohoInvoice.customerpayments.UPDATE`, `ZohoInvoice.customerpayments.DELETE` |
| `ZOHO_INVOICE_CREATE_EMPLOYEE` | Create Employee | write | `email`, `name`, `organization_id` |  | `ZohoInvoice.expenses.CREATE` |
| `ZOHO_INVOICE_CREATE_ESTIMATE` | Create Estimate | write | `customer_id`, `line_items`, `organization_id`; *adjustment, adjustment_description, contact_persons, custom_fields, date, discount, discount_type, estimate_number and 16 more* |  | `ZohoInvoice.estimates.Create` |
| `ZOHO_INVOICE_CREATE_ESTIMATE_COMMENT` | Create Estimate Comment | write | `description`, `estimate_id`, `organization_id` |  | `ZohoInvoice.estimates.CREATE`, `ZohoInvoice.estimates.Create`, `ZohoInvoice.fullaccess.all` |
| `ZOHO_INVOICE_CREATE_EXCHANGE_RATE` | Create Exchange Rate | write | `currency_id`, `effective_date`, `organization_id`, `rate` |  | `ZohoInvoice.settings.CREATE` |
| `ZOHO_INVOICE_CREATE_EXPENSE_CATEGORY` | Create Expense Category | write | `account_name`, `organization_id`; *description* |  | `ZohoInvoice.settings.CREATE` |
| `ZOHO_INVOICE_CREATE_INVOICE` | Create Invoice | write | `customer_id`, `date`, `line_items`, `organization_id`; *adjustment, adjustment_description, custom_fields, discount, due_date, invoice_number, notes, payment_terms and 6 more* |  | `ZohoInvoice.invoices.Create` |
| `ZOHO_INVOICE_CREATE_ITEM` | Create Item | write | `name`, `organization_id`, `rate`; *custom_fields, description, hsn_or_sac, is_taxable, item_tax_preferences, product_type, sat_item_key_code, sku and 3 more* |  | `ZohoInvoice.settings.CREATE` |
| `ZOHO_INVOICE_CREATE_PROJECT_COMMENT` | Create Project Comment | write | `description`, `organization_id`, `project_id` |  | `ZohoInvoice.projects.CREATE` |
| `ZOHO_INVOICE_CREATE_RECURRING_INVOICE` | Create Recurring Invoice | write | `customer_id`, `organization_id`, `recurrence_frequency`, `recurrence_name`; *allow_partial_payments, avatax_exempt_no, avatax_use_code, cfdi_usage, contact_persons, custom_fields, email, end_date and 12 more* |  | `ZohoInvoice.invoices.CREATE`, `ZohoInvoice.fullaccess.all` |
| `ZOHO_INVOICE_CREATE_REFUND_CREDIT_NOTE_REFUNDS` | Create Refund Credit Note Refunds | write | `amount`, `creditnote_id`, `date`, `organization_id`, `refund_mode`; *description, reference_number* |  | `ZohoInvoice.creditnotes.CREATE` |
| `ZOHO_INVOICE_CREATE_TASK` | Create Task | write | `organization_id`, `project_id`, `task_name`; *budget_hours, description, rate* |  | `ZohoInvoice.projects.CREATE` |
| `ZOHO_INVOICE_CREATE_TAX` | Create Tax | write | `organization_id`, `tax_name`, `tax_percentage`, `tax_type`; *country_code, purchase_tax_account_id, tax_account_id, tax_authority_id, tax_factor, tax_specific_type* |  | `ZohoInvoice.settings.CREATE` |
| `ZOHO_INVOICE_CREATE_TAX_GROUP` | Create Tax Group | write | `organization_id`, `tax_group_name`, `taxes` |  | `ZohoInvoice.settings.CREATE`, `ZohoInvoice.fullaccess.all` |
| `ZOHO_INVOICE_CREATE_TIME_ENTRY` | Create Time Entry | write | `log_date`, `organization_id`, `project_id`; *begin_time, end_time, is_billable, log_time, notes, task_id, user_id* |  | `ZohoInvoice.projects.CREATE` |
| `ZOHO_INVOICE_CREATE_USER` | Create Zoho Invoice User | write | `email`, `name`, `organization_id`; *user_role* |  | `ZohoInvoice.settings.CREATE` |
| `ZOHO_INVOICE_DELETE_ADDITIONAL_ADDRESS` | Delete Additional Address | write | `address_id`, `contact_id`, `organization_id` |  | `ZohoInvoice.contacts.DELETE` |
| `ZOHO_INVOICE_DELETE_CONTACT` | Delete Contact | write | `contact_id`, `organization_id` |  | `ZohoInvoice.contacts.DELETE`, `ZohoInvoice.contacts.ALL` |
| `ZOHO_INVOICE_DELETE_CONTACT_PERSON` | Delete Contact Person | write | `contact_person_id`, `organization_id` |  | `ZohoInvoice.contacts.DELETE` |
| `ZOHO_INVOICE_DELETE_CREDIT_NOTES_APPLIED_TO_INVOICE` | Delete Credit Notes Applied to Invoice | write | `creditnote_id`, `creditnote_invoice_id`, `organization_id` |  | `ZohoInvoice.creditnotes.DELETE` |
| `ZOHO_INVOICE_DELETE_CURRENCY` | Delete Currency | write | `currency_id`, `organization_id` |  | `ZohoInvoice.settings.DELETE` |
| `ZOHO_INVOICE_DELETE_CUSTOMER_PAYMENT` | Delete Customer Payment | write | `organization_id`, `payment_id` |  | `ZohoInvoice.customerpayments.DELETE` |
| `ZOHO_INVOICE_DELETE_EMPLOYEE` | Delete Employee | write | `employee_id`, `organization_id` |  | `ZohoInvoice.expenses.DELETE` |
| `ZOHO_INVOICE_DELETE_ESTIMATES` | Delete Estimates | write | `estimate_ids`, `organization_id` |  | `ZohoInvoice.estimates.DELETE` |
| `ZOHO_INVOICE_DELETE_ESTIMATE_COMMENT` | Delete Estimate Comment | write | `comment_id`, `estimate_id`, `organization_id` |  | `ZohoInvoice.estimates.DELETE` |
| `ZOHO_INVOICE_DELETE_EXPENSE` | Delete Expense | write | `expense_id`, `organization_id` |  | `ZohoInvoice.expenses.DELETE` |
| `ZOHO_INVOICE_DELETE_EXPENSE_CATEGORY` | Delete Expense Category | write | `category_id`, `organization_id` |  | `ZohoInvoice.settings.DELETE` |
| `ZOHO_INVOICE_DELETE_INVOICE` | Delete Invoice | write | `invoice_id`, `organization_id` |  | `ZohoInvoice.invoices.DELETE` |
| `ZOHO_INVOICE_DELETE_INVOICE_ATTACHMENT` | Delete Invoice Attachment | write | `invoice_id`, `organization_id` |  | `ZohoInvoice.invoices.DELETE` |
| `ZOHO_INVOICE_DELETE_INVOICE_COMMENT` | Delete Invoice Comment | write | `comment_id`, `invoice_id`, `organization_id` |  | `ZohoInvoice.invoices.DELETE` |
| `ZOHO_INVOICE_DELETE_INVOICE_EXPENSE_RECEIPT` | Delete Invoice Expense Receipt | write | `expense_id`, `organization_id` |  | `ZohoInvoice.invoices.DELETE` |
| `ZOHO_INVOICE_DELETE_ITEM` | Delete Item | write | `item_id`, `organization_id` |  | `ZohoInvoice.settings.DELETE` |
| `ZOHO_INVOICE_DELETE_PROJECT` | Delete Project | write | `organization_id`, `project_id` |  | `ZohoInvoice.projects.DELETE` |
| `ZOHO_INVOICE_DELETE_PROJECT_COMMENT` | Delete Project Comment | write | `comment_id`, `organization_id`, `project_id` |  | `ZohoInvoice.projects.DELETE` |
| `ZOHO_INVOICE_DELETE_SETTING_TAXES` | Delete Tax | write | `organization_id`, `tax_id` |  | `ZohoInvoice.settings.DELETE` |
| `ZOHO_INVOICE_DELETE_TASK` | Delete Task | write | `organization_id`, `project_id`, `task_id` |  | `ZohoInvoice.projects.DELETE` |
| `ZOHO_INVOICE_DELETE_TIME_ENTRY` | Delete Time Entry | write | `organization_id`, `time_entry_id` |  | `ZohoInvoice.projects.DELETE` |
| `ZOHO_INVOICE_DELETE_USER` | Delete User | write | `organization_id`, `user_id` |  | `ZohoInvoice.settings.DELETE` |
| `ZOHO_INVOICE_DISABLE_CONTACT_PAYMENT_REMINDERS` | Disable Contact Payment Reminders | write | `contact_id`, `organization_id` |  | `ZohoInvoice.contacts.CREATE` |
| `ZOHO_INVOICE_DISABLE_INVOICE_PAYMENT_REMINDER` | Disable Invoice Payment Reminder | write | `invoice_id`, `organization_id` |  | `ZohoInvoice.invoices.CREATE` |
| `ZOHO_INVOICE_EMAIL_CONTACT_STATEMENT` | Email Contact Statement | write | `contact_id`, `organization_id`; *end_date, start_date* |  | `ZohoInvoice.contacts.CREATE`, `ZohoInvoice.fullaccess.all` |
| `ZOHO_INVOICE_EMAIL_ESTIMATE` | Email Estimate | write | `estimate_id`, `organization_id`, `to_mail_ids`; *body, cc_mail_ids, send_from_org_email_id, subject* |  | `ZohoInvoice.estimates.CREATE`, `ZohoInvoice.estimates.Create` |
| `ZOHO_INVOICE_EMAIL_INVOICE` | Email Invoice | write | `invoice_id`, `organization_id`, `to_mail_ids`; *body, cc_mail_ids, send_attachment, send_customer_statement, send_from_org_email_id, subject* |  | `ZohoInvoice.invoices.CREATE` |
| `ZOHO_INVOICE_EMAIL_MULTIPLE_ESTIMATES` | Email Multiple Estimates | write | `estimate_ids`, `organization_id` |  | `ZohoInvoice.estimates.CREATE` |
| `ZOHO_INVOICE_ENABLE_INVOICE_PAYMENT_REMINDER` | Enable Invoice Payment Reminder | write | `invoice_id`, `organization_id` |  | `ZohoInvoice.invoices.CREATE`, `ZohoInvoice.invoices.ALL`, `ZohoInvoice.fullaccess.all` |
| `ZOHO_INVOICE_ENABLE_PAYMENT_REMINDERS` | Enable Payment Reminders | write | `contact_id`, `organization_id` |  | `ZohoInvoice.contacts.CREATE`, `ZohoInvoice.fullaccess.all` |
| `ZOHO_INVOICE_ENABLE_PORTAL_ACCESS` | Enable Portal Access | write | `contact_id`, `contact_persons`, `organization_id` |  | `ZohoInvoice.contacts.CREATE` |
| `ZOHO_INVOICE_EXPORT_INVOICES_BULK` | Bulk Export Invoices | read | `invoice_ids`, `organization_id` |  | `ZohoInvoice.invoices.READ` |
| `ZOHO_INVOICE_GET_ALL_TASKS` | Get All Tasks | read | `organization_id`, `project_id`; *page, per_page* |  | `ZohoInvoice.projects.READ` |
| `ZOHO_INVOICE_GET_CLIENT_REVIEW` | Get Client Review | read | `comment_id`, `organization_id`; *contact_id* |  | `ZohoInvoice.contacts.READ` |
| `ZOHO_INVOICE_GET_CONTACT` | Get Contact | read | `contact_id`, `organization_id` |  | `ZohoInvoice.contacts.READ` |
| `ZOHO_INVOICE_GET_CONTACT_ADDRESSES` | Get Contact Addresses | read | `contact_id`, `organization_id` |  | `ZohoInvoice.contacts.READ` |
| `ZOHO_INVOICE_GET_CONTACT_PERSON` | Get Contact Person | read | `contact_id`, `contact_person_id`, `organization_id` |  | `ZohoInvoice.contacts.READ` |
| `ZOHO_INVOICE_GET_CREDIT_NOTE` | Get Credit Note | read | `creditnote_id`, `organization_id`; *accept, print* |  | `ZohoInvoice.creditnotes.READ` |
| `ZOHO_INVOICE_GET_CREDIT_NOTE_EMAIL_HISTORY` | Get Credit Note Email History | read | `creditnote_id`, `organization_id` |  | `ZohoInvoice.creditnotes.READ` |
| `ZOHO_INVOICE_GET_CREDIT_NOTE_REFUND` | Get Credit Note Refund | read | `creditnote_id`, `creditnote_refund_id`, `organization_id` |  | `ZohoInvoice.creditnotes.READ` |
| `ZOHO_INVOICE_GET_EXPENSE` | Get Expense | read | `expense_id`, `organization_id` |  | `ZohoInvoice.expenses.READ` |
| `ZOHO_INVOICE_GET_INVOICE` | Get Invoice | read | `invoice_id`, `organization_id` |  | `ZohoInvoice.invoices.READ` |
| `ZOHO_INVOICE_GET_INVOICE_ATTACHMENT` | Get Invoice Attachment | read | `invoice_id`, `organization_id`; *preview* |  | `ZohoInvoice.invoices.READ` |
| `ZOHO_INVOICE_GET_INVOICE_EMAIL_CONTENT` | Get Invoice Email Content | read | `invoice_id`, `organization_id`; *email_template_id* |  | `ZohoInvoice.invoices.READ`, `ZohoInvoice.fullaccess.all` |
| `ZOHO_INVOICE_GET_ITEM` | Get Zoho Invoice Item | read | `item_id`, `organization_id` |  | `ZohoInvoice.settings.READ` |
| `ZOHO_INVOICE_GET_PAYMENT_REMINDER_MAIL_CONTENT` | Get Payment Reminder Mail Content | read | `invoice_id`, `organization_id` |  | `ZohoInvoice.invoices.READ` |
| `ZOHO_INVOICE_GET_PRICE_LIST` | Get Price List | read | `organization_id`, `pricebook_id` |  | `ZohoInvoice.settings.READ` |
| `ZOHO_INVOICE_GET_PROJECT` | Get Project | read | `organization_id`, `project_id` |  | `ZohoInvoice.projects.READ` |
| `ZOHO_INVOICE_GET_PROJECT_USER` | Get Project User | read | `organization_id`, `project_id`, `user_id` |  | `ZohoInvoice.projects.READ` |
| `ZOHO_INVOICE_GET_RECURRING_INVOICE` | Get Recurring Invoice | read | `organization_id`, `recurring_invoice_id` |  | `ZohoInvoice.invoices.READ` |
| `ZOHO_INVOICE_GET_STATEMENT_MAIL_CONTENT` | Get Statement Mail Content | read | `contact_id`, `organization_id`; *end_date, start_date* |  | `ZohoInvoice.contacts.READ` |
| `ZOHO_INVOICE_GET_TASK` | Get Task | read | `organization_id`, `project_id`, `task_id` |  | `ZohoInvoice.projects.READ` |
| `ZOHO_INVOICE_GET_TAX` | Get Tax | read | `tax_id`; *organization_id* |  | `ZohoInvoice.settings.READ`, `ZohoInvoice.fullaccess.all` |
| `ZOHO_INVOICE_GET_TAX_GROUP` | Get Tax Group | read | `organization_id`, `tax_group_id` |  | `ZohoInvoice.settings.READ` |
| `ZOHO_INVOICE_GET_TIME_ENTRY` | Get Time Entry | read | `organization_id`, `time_entry_id` |  | `ZohoInvoice.projects.READ` |
| `ZOHO_INVOICE_INACTIVATE_PROJECT` | Inactivate Project | write | `organization_id`, `project_id` |  | `ZohoInvoice.projects.CREATE` |
| `ZOHO_INVOICE_LIST_CHILD_EXPENSES_CREATED` | List Child Expenses Created | read | `recurring_expense_id`; *organization_id, page, per_page, sort_column* |  | `ZohoInvoice.expenses.READ` |
| `ZOHO_INVOICE_LIST_CLIENT_REVIEWS` | List Client Reviews | read | `organization_id`; *contact_id, page, per_page* |  | `ZohoInvoice.contacts.READ` |
| `ZOHO_INVOICE_LIST_CONTACTS` | List Contacts | read | `organization_id`; *address, address_contains, address_startswith, company_name, company_name_contains, company_name_startswith, contact_name, contact_name_contains and 20 more* |  | `ZohoInvoice.contacts.READ` |
| `ZOHO_INVOICE_LIST_CONTACT_COMMENTS` | List Contact Comments | read | `contact_id`, `organization_id`; *page, per_page* |  | `ZohoInvoice.contacts.READ` |
| `ZOHO_INVOICE_LIST_CONTACT_REFUNDS` | List Contact Refunds | read | `contact_id`, `organization_id`; *page, per_page* |  | `ZohoInvoice.contacts.READ` |
| `ZOHO_INVOICE_LIST_CREDIT_NOTES` | List Credit Notes | read | `organization_id`; *creditnote_number, customer_id, customer_name, date, filter_by, item_description, item_id, item_name and 7 more* |  | `ZohoInvoice.creditnotes.READ` |
| `ZOHO_INVOICE_LIST_CURRENCIES` | List Currencies | read | `organization_id`; *page, per_page* |  | `ZohoInvoice.settings.READ` |
| `ZOHO_INVOICE_LIST_CUSTOMER_PAYMENT_REFUNDS` | List Customer Payment Refunds | read | `customer_payment_id`, `organization_id`; *page, per_page* |  | `ZohoInvoice.customerpayments.READ` |
| `ZOHO_INVOICE_LIST_EMPLOYEES` | List Employees | read | `organization_id`; *page, per_page* |  | `ZohoInvoice.expenses.READ` |
| `ZOHO_INVOICE_LIST_ESTIMATES` | List Estimates | read | `organization_id`; *page, per_page, zcrm_potential_id* |  | `ZohoInvoice.estimates.READ` |
| `ZOHO_INVOICE_LIST_EXPENSES` | List Expenses | read | *account_name, customer_id, customer_name, date_end, date_start, filter_by, page, per_page and 3 more* |  | `ZohoInvoice.expenses.READ`, `ZohoInvoice.expenses.CREATE`, `ZohoInvoice.expenses.UPDATE`, `ZohoInvoice.expenses.DELETE`, `ZohoInvoice.expenses.ALL` |
| `ZOHO_INVOICE_LIST_EXPENSE_CATEGORIES` | List Expense Categories | read | `organization_id`; *filter_by, page, per_page, sort_column* |  | `ZohoInvoice.settings.READ` |
| `ZOHO_INVOICE_LIST_EXPENSE_COMMENTS` | List Expense Comments | read | `expense_id`, `organization_id` |  | `ZohoInvoice.expenses.READ` |
| `ZOHO_INVOICE_LIST_INVOICES` | List Invoices | read | `organization_id`; *balance, custom_field, custom_field_contains, custom_field_startswith, customer_id, customer_name, date, date_after and 27 more* |  | `ZohoInvoice.invoices.CREATE`, `ZohoInvoice.invoices.UPDATE`, `ZohoInvoice.invoices.READ`, `ZohoInvoice.invoices.DELETE`, `ZohoInvoice.invoices.ALL` |
| `ZOHO_INVOICE_LIST_INVOICES_CREDITED` | List Invoices Credited | read | `creditnote_id`, `organization_id` |  | `ZohoInvoice.creditnotes.READ` |
| `ZOHO_INVOICE_LIST_INVOICE_COMMENTS` | List Invoice Comments | read | `invoice_id`, `organization_id` |  | `ZohoInvoice.invoices.READ` |
| `ZOHO_INVOICE_LIST_INVOICE_PAYMENTS` | List Invoice Payments | read | `invoice_id`, `organization_id` |  | `ZohoInvoice.invoices.READ` |
| `ZOHO_INVOICE_LIST_ITEMS` | List Items | read | `organization_id`; *filter_by, page, per_page, search_text* |  | `ZohoInvoice.settings.READ` |
| `ZOHO_INVOICE_LIST_ORGANIZATIONS` | List Organizations | read |  |  | `ZohoInvoice.settings.READ` |
| `ZOHO_INVOICE_LIST_PAYMENTS` | List Payments | read | *customer_id, date_end, date_start, filter_by, invoice_id, mode, page, payment_number and 4 more* |  | `ZohoInvoice.customerpayments.READ`, `ZohoInvoice.customerpayments.CREATE`, `ZohoInvoice.customerpayments.UPDATE`, `ZohoInvoice.customerpayments.DELETE`, `ZohoInvoice.customerpayments.ALL` |
| `ZOHO_INVOICE_LIST_PRICE_LISTS` | List Price Lists | read | `organization_id`; *page, per_page* |  | `ZohoInvoice.settings.READ` |
| `ZOHO_INVOICE_LIST_PROJECTS` | List Projects | read | `organization_id`; *customer_id, filter_by, page, per_page, sort_column* |  | `ZohoInvoice.projects.READ` |
| `ZOHO_INVOICE_LIST_PROJECT_COMMENTS` | List Project Comments | read | `organization_id`, `project_id` |  | `ZohoInvoice.projects.READ` |
| `ZOHO_INVOICE_LIST_PROJECT_INVOICES` | List Project Invoices | read | `organization_id`, `project_id`; *page, per_page, sort_column* |  | `ZohoInvoice.projects.READ` |
| `ZOHO_INVOICE_LIST_PROJECT_USERS` | List Project Users | read | `organization_id`, `project_id` |  | `ZohoInvoice.projects.READ` |
| `ZOHO_INVOICE_LIST_RECURRING_INVOICES` | List Recurring Invoices | read | *filter_by, organization_id, page, per_page, search_text, sort_column* |  | `ZohoInvoice.invoices.READ`, `ZohoInvoice.fullaccess.all` |
| `ZOHO_INVOICE_LIST_RETAINER_INVOICES` | List Retainer Invoices | read | *filter_by, organization_id, page, per_page, print, sort_column, sort_order* |  | `ZohoInvoice.invoices.READ` |
| `ZOHO_INVOICE_LIST_RETAINER_INVOICE_TEMPLATES` | List Retainer Invoice Templates | read | `organization_id` |  | `ZohoInvoice.invoices.READ` |
| `ZOHO_INVOICE_LIST_USERS` | List Users | read | `organization_id`; *filter_by, page, per_page, sort_column* |  | `ZohoInvoice.settings.READ` |
| `ZOHO_INVOICE_MARK_CONTACT_ACTIVE` | Mark Contact as Active | write | `contact_id`, `organization_id` |  | `ZohoInvoice.contacts.CREATE` |
| `ZOHO_INVOICE_MARK_CONTACT_INACTIVE` | Mark Contact as Inactive | write | `contact_id`, `organization_id` |  | `ZohoInvoice.contacts.CREATE` |
| `ZOHO_INVOICE_MARK_CONTACT_PERSON_AS_PRIMARY` | Mark Contact Person as Primary | write | `contact_person_id`, `organization_id` |  | `ZohoInvoice.contacts.CREATE` |
| `ZOHO_INVOICE_MARK_ESTIMATE_AS_DECLINED` | Mark Estimate as Declined | write | `estimate_id`, `organization_id` |  | `ZohoInvoice.estimates.CREATE` |
| `ZOHO_INVOICE_MARK_EXPENSE_CATEGORY_ACTIVE` | Mark Expense Category as Active | write | `category_id`, `organization_id` |  | `ZohoInvoice.settings.CREATE` |
| `ZOHO_INVOICE_MARK_INVOICE_AS_SENT` | Mark Invoice as Sent | write | `invoice_id`, `organization_id` |  | `ZohoInvoice.invoices.CREATE` |
| `ZOHO_INVOICE_MARK_INVOICE_AS_VOID` | Mark Invoice as Void | write | `invoice_id`, `organization_id` |  | `ZohoInvoice.invoices.CREATE` |
| `ZOHO_INVOICE_MARK_ITEM_AS_INACTIVE` | Mark Item as Inactive | write | `item_id`, `organization_id` |  | `ZohoInvoice.settings.CREATE`, `ZohoInvoice.settings.UPDATE`, `ZohoInvoice.settings.READ`, `ZohoInvoice.settings.DELETE` |
| `ZOHO_INVOICE_MARK_RETAINER_INVOICE_AS_SENT` | Mark Retainer Invoice as Sent | write | `organization_id`, `retainerinvoice_id` |  | `ZohoInvoice.invoices.CREATE` |
| `ZOHO_INVOICE_PRINT_ESTIMATES` | Bulk Print Estimates | read | `estimate_ids`, `organization_id` |  | `ZohoInvoice.estimates.READ` |
| `ZOHO_INVOICE_PRINT_INVOICES` | Bulk Print Invoices | read | `invoice_ids`, `organization_id` |  | `ZohoInvoice.invoices.READ` |
| `ZOHO_INVOICE_REFUND_CUSTOMER_PAYMENT` | Refund Customer Payment | write | `amount`, `customer_payment_id`, `date`, `organization_id`; *description, exchange_rate, payment_form, reference_number, refund_mode* |  | `ZohoInvoice.customerpayments.CREATE` |
| `ZOHO_INVOICE_RESUME_RECURRING_INVOICE` | Resume Recurring Invoice | write | `recurring_invoice_id`; *organization_id* |  | `ZohoInvoice.invoices.CREATE` |
| `ZOHO_INVOICE_SEND_BULK_INVOICE_REMINDER` | Send Bulk Invoice Reminder | write | `invoice_ids`, `organization_id` |  | `ZohoInvoice.invoices.CREATE` |
| `ZOHO_INVOICE_SEND_CONTACT_EMAIL` | Send Contact Email | write | `body`, `contact_id`, `organization_id`, `subject`, `to_mail_ids`; *send_customer_statement* |  | `ZohoInvoice.contacts.CREATE` |
| `ZOHO_INVOICE_START_TIMER` | Start Timer | write | `organization_id`, `time_entry_id` |  | `ZohoInvoice.projects.CREATE` |
| `ZOHO_INVOICE_STOP_RECURRING_INVOICE` | Stop Recurring Invoice | write | `recurring_invoice_id`; *organization_id* |  | `ZohoInvoice.invoices.CREATE`, `ZohoInvoice.fullaccess.all` |
| `ZOHO_INVOICE_UPDATE_ADDITIONAL_ADDRESS` | Update Additional Address | write | `address_id`, `contact_id`, `organization_id`; *address, attention, city, country, fax, phone, state, state_code and 2 more* |  | `ZohoInvoice.contacts.UPDATE` |
| `ZOHO_INVOICE_UPDATE_CONTACT` | Update Contact | write | `contact_id`, `organization_id`; *billing_address, company_name, contact_name, contact_persons, country_code, currency_id, custom_fields, facebook and 18 more* |  | `ZohoInvoice.contacts.UPDATE` |
| `ZOHO_INVOICE_UPDATE_CONTACT_PERSON` | Update Contact Person | write | `contact_id`, `contact_person_id`, `first_name`, `organization_id`; *department, designation, email, enable_portal, last_name, mobile, phone, salutation and 1 more* |  | `ZohoInvoice.contacts.UPDATE` |
| `ZOHO_INVOICE_UPDATE_CREDIT_NOTE` | Update Credit Note | write | `creditnote_id`, `creditnote_number`, `customer_id`, `date`, `line_items`, `organization_id`; *ignore_auto_number_generation, notes, reference_number, terms* |  | `ZohoInvoice.creditnotes.UPDATE`, `ZohoInvoice.fullaccess.all` |
| `ZOHO_INVOICE_UPDATE_CREDIT_NOTE_TEMPLATE` | Update Credit Note Template | write | `creditnote_id`, `organization_id`, `template_id` |  | `ZohoInvoice.creditnotes.UPDATE` |
| `ZOHO_INVOICE_UPDATE_CUSTOMER_PAYMENT_REFUND` | Update Customer Payment Refund | write | `customer_payment_id`, `organization_id`, `refund_id`; *amount, date, description, exchange_rate, payment_form, reference_number, refund_mode* |  | `ZohoInvoice.customerpayments.UPDATE` |
| `ZOHO_INVOICE_UPDATE_ESTIMATE_SHIPPING_ADDRESS` | Update Estimate Shipping Address | write | `address`, `city`, `country`, `estimate_id`, `organization_id`, `state`, `zip`; *attention, fax, phone, street2* |  | `ZohoInvoice.estimates.UPDATE` |
| `ZOHO_INVOICE_WRITE_OFF_INVOICE` | Write Off Invoice | write | `invoice_id`, `organization_id` |  | `ZohoInvoice.invoices.CREATE` |
