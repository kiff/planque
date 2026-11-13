# bexio: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner pastes a key from the vendor's settings). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

309 tools: 157 read, 152 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `BEXIO_ACCEPT_QUOTE` | Accept Quote | write | `quote_id` |  |  |
| `BEXIO_ARCHIVE_PROJECT` | Archive Project | write | `project_id` |  |  |
| `BEXIO_CANCEL_INVOICE` | Cancel Invoice | write | `invoice_id` |  |  |
| `BEXIO_CANCEL_PAYMENT` | Cancel Payment | write | `payment_id` |  |  |
| `BEXIO_CHECK_BILL_DOCUMENT_NUMBER` | Check Bill Document Number | read | `document_no` |  |  |
| `BEXIO_CHECK_EXPENSE_DOCUMENT_NUMBER` | Check Expense Document Number | read | `document_no` |  |  |
| `BEXIO_COPY_INVOICE` | Copy Invoice | write | `contact_id`, `invoice_id`; *contact_sub_id, is_valid_from, title* |  |  |
| `BEXIO_COPY_QUOTE` | Copy Quote | write | `contact_id`, `quote_id`; *contact_sub_id, is_valid_from, pr_project_id, title* |  |  |
| `BEXIO_CREATE_ADDITIONAL_ADDRESS` | Create Additional Address | write | `contact_id`; *address_addition, city, country_id, description, house_number, name, name_addition, postcode and 2 more* |  |  |
| `BEXIO_CREATE_BILL` | Create Bill | write | `address`, `attachment_ids`, `bill_date`, `contact_partner_id`, `currency_code`, `discounts`, `due_date`, `item_net`, `line_items`, `manual_amount`, `supplier_id`; *amount_calc, amount_man, base_currency_amount, exchange_rate, payment, purchase_order_id, qr_bill_information, title and 1 more* |  |  |
| `BEXIO_CREATE_BUSINESS_ACTIVITY` | Create Business Activity | write | `name`; *account_id, default_is_billable, default_price_per_hour* |  |  |
| `BEXIO_CREATE_CALENDAR_YEAR` | Create Calendar Year | write | *default_tax_expense_id, default_tax_income_id, is_annual_reporting, is_vat_subject, vat_accounting_method, vat_accounting_type, year* |  |  |
| `BEXIO_CREATE_COMMENT` | Create Comment | write | `document_id`, `kb_document_type`, `text`, `user_id`, `user_name`; *is_public, user_email* |  |  |
| `BEXIO_CREATE_CONTACT` | Create Contact | write | `contact_type_id`, `name_1`, `owner_id`, `user_id`; *address_addition, birthday, city, contact_branch_ids, contact_group_ids, contact_number, country_id, fax and 16 more* |  |  |
| `BEXIO_CREATE_CONTACTS_BULK` | Create Contacts Bulk | write | *items* |  |  |
| `BEXIO_CREATE_CONTACT_GROUP` | Create Contact Group | write | `name` |  |  |
| `BEXIO_CREATE_CONTACT_RELATION` | Create Contact Relation | write | `contact_id`, `contact_sub_id`; *description* |  |  |
| `BEXIO_CREATE_COUNTRY` | Create Country | write | `iso3166_alpha2`, `name`, `name_short` |  |  |
| `BEXIO_CREATE_CURRENCY` | Create Currency | write | `name`, `round_factor` |  |  |
| `BEXIO_CREATE_DEFAULT_POSITION` | Create Default Position | write | `document_id`, `kb_document_type`; *account_id, amount, amount_completed, amount_open, amount_reserved, discount_in_percent, is_optional, tax_id and 3 more* |  |  |
| `BEXIO_CREATE_DELIVERY_FROM_ORDER` | Create Delivery From Order | write | `order_id`; *positions* |  |  |
| `BEXIO_CREATE_DISCOUNT_POSITION` | Create Discount Position | write | `document_id`, `kb_document_type`; *is_percentual, text, value* |  |  |
| `BEXIO_CREATE_EMPLOYEE` | Create Employee | write | *address, ahv_number, annual_vacation_days, date_of_birth, email, first_name, gender, iban and 6 more* |  |  |
| `BEXIO_CREATE_EMPLOYEE_ABSENCE` | Create Employee Absence | write | `employee_id`, `reason`, `start_date`; *continued_pay, disability, end_date, half_day, paid_hours* |  |  |
| `BEXIO_CREATE_EXPENSE` | Create Expense | write | `amount`, `attachment_ids`, `currency_code`, `paid_on`; *address, bank_account_id, base_currency_amount, booking_account_id, exchange_rate, supplier_id, tax_id, title* |  |  |
| `BEXIO_CREATE_FICTIONAL_USER` | Create Fictional User | write | `email`, `firstname`, `lastname`, `salutation_type`; *title_id* |  |  |
| `BEXIO_CREATE_FILE` | Create File | write | `file` |  |  |
| `BEXIO_CREATE_INVOICE` | Create Invoice | write | *api_reference, bank_account_id, contact_address_manual, contact_id, contact_sub_id, currency_id, document_nr, footer and 15 more* |  |  |
| `BEXIO_CREATE_INVOICE_FROM_ORDER` | Create Invoice From Order | write | `order_id`; *positions* |  |  |
| `BEXIO_CREATE_INVOICE_FROM_QUOTE` | Create Invoice From Quote | write | `quote_id`; *positions* |  |  |
| `BEXIO_CREATE_INVOICE_PAYMENT` | Create Invoice Payment | write | `invoice_id`, `value`; *bank_account_id, date, payment_service_id* |  |  |
| `BEXIO_CREATE_INVOICE_REMINDER` | Create Invoice Reminder | write | `invoice_id` |  |  |
| `BEXIO_CREATE_ITEM` | Create Item | write | *account_id, article_group_id, article_type_id, contact_id, currency_id, deliverer_code, deliverer_description, deliverer_name and 24 more* |  |  |
| `BEXIO_CREATE_ITEM_POSITION` | Create Item Position | write | `document_id`, `kb_document_type`; *account_id, amount, amount_completed, amount_open, amount_reserved, article_id, discount_in_percent, is_optional and 4 more* |  |  |
| `BEXIO_CREATE_MANUAL_ENTRY` | Create Manual Entry | write | `date`, `entries`, `type`; *reference_nr* |  |  |
| `BEXIO_CREATE_MILESTONE` | Create Milestone | write | `name`, `project_id`; *comment, end_date, pr_parent_milestone_id* |  |  |
| `BEXIO_CREATE_NOTE` | Create Note | write | `event_start`, `subject`, `user_id`; *contact_id, entry_id, info, module_id, pr_project_id* |  |  |
| `BEXIO_CREATE_ORDER` | Create Order | write | *api_reference, bank_account_id, contact_address_manual, contact_id, contact_sub_id, currency_id, delivery_address_manual, delivery_address_type and 15 more* |  |  |
| `BEXIO_CREATE_ORDER_FROM_QUOTE` | Create Order From Quote | write | `quote_id`; *positions* |  |  |
| `BEXIO_CREATE_OUTGOING_PAYMENT` | Create Outgoing Payment | write | `amount`, `bill_id`, `currency_code`, `exchange_rate`, `execution_date`, `is_salary_payment`, `payment_type`, `sender_bank_account_id`; *booking_text, fee_type, message, note, receiver_account_no, receiver_bank_name, receiver_bank_no, receiver_bc_no and 18 more* |  |  |
| `BEXIO_CREATE_PAGEBREAK_POSITION` | Create Pagebreak Position | write | `document_id`, `kb_document_type`; *pagebreak* |  |  |
| `BEXIO_CREATE_PAYMENT` | Create Payment | write | `account_id`, `amount`, `currency`, `execution_date`, `is_salary`, `recipient`, `type`; *additional_information, allowance, is_editing_restricted, message, purchase_reference, qr_reference_number* |  |  |
| `BEXIO_CREATE_PROJECT` | Create Project | write | `contact_id`, `name`, `pr_project_type_id`, `pr_state_id`, `user_id`; *comment, contact_sub_id, document_nr, end_date, pr_budget_type_amount, pr_budget_type_id, pr_invoice_type_amount, pr_invoice_type_id and 1 more* |  |  |
| `BEXIO_CREATE_PURCHASE_ORDER` | Create Purchase Order | write | *api_reference, bank_account_id, contact_address_manual, contact_id, contact_sub_id, created_at, currency, currency_id and 32 more* |  |  |
| `BEXIO_CREATE_QUOTE` | Create Quote | write | *api_reference, bank_account_id, contact_address_manual, contact_id, contact_sub_id, currency_id, delivery_address_manual, delivery_address_type and 18 more* |  |  |
| `BEXIO_CREATE_SALUTATION` | Create Salutation | write | `name` |  |  |
| `BEXIO_CREATE_SUBPOSITION_POSITION` | Create Subposition Position | write | `document_id`, `kb_document_type`; *show_pos_nr, text* |  |  |
| `BEXIO_CREATE_SUBTOTAL_POSITION` | Create Subtotal Position | write | `document_id`, `kb_document_type`; *text* |  |  |
| `BEXIO_CREATE_TASK` | Create Task | write | `subject`, `user_id`; *communication_kind_id, contact_id, entry_id, finish_date, have_remember, info, module_id, pr_project_id and 5 more* |  |  |
| `BEXIO_CREATE_TEXT_POSITION` | Create Text Position | write | `document_id`, `kb_document_type`; *show_pos_nr, text* |  |  |
| `BEXIO_CREATE_TIMESHEET` | Create Timesheet | write | `allowable_bill`, `client_service_id`, `tracking`, `user_id`; *charge, contact_id, estimated_time, pr_milestone_id, pr_package_id, pr_project_id, status_id, sub_contact_id and 1 more* |  |  |
| `BEXIO_CREATE_TITLE` | Create Title | write | `name` |  |  |
| `BEXIO_CREATE_UNIT` | Create Unit | write | `name` |  |  |
| `BEXIO_CREATE_WORK_PACKAGE` | Create Work Package | write | `name`, `project_id`; *comment, estimated_time_in_hours, pr_milestone_id, spent_time_in_hours* |  |  |
| `BEXIO_DECLINE_QUOTE` | Decline Quote | write | `quote_id` |  |  |
| `BEXIO_DELETE_ADDITIONAL_ADDRESS` | Delete Additional Address | write | `additional_address_id`, `contact_id` |  |  |
| `BEXIO_DELETE_BILL` | Delete Bill | write | `id` |  |  |
| `BEXIO_DELETE_CONTACT` | Delete Contact | write | `contact_id` |  |  |
| `BEXIO_DELETE_CONTACT_GROUP` | Delete Contact Group | write | `contact_group_id` |  |  |
| `BEXIO_DELETE_CONTACT_RELATION` | Delete Contact Relation | write | `contact_relation_id` |  |  |
| `BEXIO_DELETE_COUNTRY` | Delete Country | write | `country_id` |  |  |
| `BEXIO_DELETE_CURRENCY` | Delete Currency | write | `currency_id` |  |  |
| `BEXIO_DELETE_DEFAULT_POSITION` | Delete Default Position | write | `document_id`, `kb_document_type`, `position_id` |  |  |
| `BEXIO_DELETE_DISCOUNT_POSITION` | Delete Discount Position | write | `document_id`, `kb_document_type`, `position_id` |  |  |
| `BEXIO_DELETE_EMPLOYEE_ABSENCE` | Delete Employee Absence | write | `absence_id`, `employee_id` |  |  |
| `BEXIO_DELETE_EXPENSE` | Delete Expense | write | `id` |  |  |
| `BEXIO_DELETE_FICTIONAL_USER` | Delete Fictional User | write | `fictional_user_id` |  |  |
| `BEXIO_DELETE_FILE` | Delete File | write | `file_id` |  |  |
| `BEXIO_DELETE_INVOICE` | Delete Invoice | write | `invoice_id` |  |  |
| `BEXIO_DELETE_INVOICE_PAYMENT` | Delete Invoice Payment | write | `invoice_id`, `payment_id` |  |  |
| `BEXIO_DELETE_INVOICE_REMINDER` | Delete Invoice Reminder | write | `invoice_id`, `reminder_id` |  |  |
| `BEXIO_DELETE_ITEM` | Delete Item | write | `article_id` |  |  |
| `BEXIO_DELETE_ITEM_POSITION` | Delete Item Position | write | `document_id`, `kb_document_type`, `position_id` |  |  |
| `BEXIO_DELETE_MANUAL_COMPOUND_ENTRY_FILE` | Delete Manual Compound Entry File | write | `file_id`, `manual_entry_id` |  |  |
| `BEXIO_DELETE_MANUAL_ENTRY` | Delete Manual Entry | write | `manual_entry_id` |  |  |
| `BEXIO_DELETE_MANUAL_ENTRY_FILE` | Delete Manual Entry File | write | `entry_id`, `file_id`, `manual_entry_id` |  |  |
| `BEXIO_DELETE_MILESTONE` | Delete Milestone | write | `milestone_id`, `project_id` |  |  |
| `BEXIO_DELETE_NOTE` | Delete Note | write | `note_id` |  |  |
| `BEXIO_DELETE_ORDER` | Delete Order | write | `order_id` |  |  |
| `BEXIO_DELETE_ORDER_REPETITION` | Delete Order Repetition | write | `order_id` |  |  |
| `BEXIO_DELETE_OUTGOING_PAYMENT` | Delete Outgoing Payment | write | `id` |  |  |
| `BEXIO_DELETE_PAGEBREAK_POSITION` | Delete Pagebreak Position | write | `document_id`, `kb_document_type`, `position_id` |  |  |
| `BEXIO_DELETE_PAYMENT` | Delete Payment | write | `payment_id` |  |  |
| `BEXIO_DELETE_PROJECT` | Delete Project | write | `project_id` |  |  |
| `BEXIO_DELETE_PURCHASE_ORDER` | Delete Purchase Order | write | `purchase_order_id` |  |  |
| `BEXIO_DELETE_QUOTE` | Delete Quote | write | `quote_id` |  |  |
| `BEXIO_DELETE_SALUTATION` | Delete Salutation | write | `salutation_id` |  |  |
| `BEXIO_DELETE_SUBPOSITION_POSITION` | Delete Subposition Position | write | `document_id`, `kb_document_type`, `position_id` |  |  |
| `BEXIO_DELETE_SUBTOTAL_POSITION` | Delete Subtotal Position | write | `document_id`, `kb_document_type`, `position_id` |  |  |
| `BEXIO_DELETE_TASK` | Delete Task | write | `task_id` |  |  |
| `BEXIO_DELETE_TAX` | Delete Tax | write | `tax_id` |  |  |
| `BEXIO_DELETE_TEXT_POSITION` | Delete Text Position | write | `document_id`, `kb_document_type`, `position_id` |  |  |
| `BEXIO_DELETE_TIMESHEET` | Delete Timesheet | write | `timesheet_id` |  |  |
| `BEXIO_DELETE_TITLE` | Delete Title | write | `title_id` |  |  |
| `BEXIO_DELETE_UNIT` | Delete Unit | write | `unit_id` |  |  |
| `BEXIO_DELETE_WORK_PACKAGE` | Delete Work Package | write | `package_id`, `project_id` |  |  |
| `BEXIO_DOWNLOAD_FILE` | Download File | read | `file_id` |  |  |
| `BEXIO_DOWNLOAD_PAYSTUB_PDF` | Download Paystub Pdf | read | `employee_id`, `month`, `year` |  |  |
| `BEXIO_DUPLICATE_BILL` | Duplicate Bill | write | `action`, `id` |  |  |
| `BEXIO_DUPLICATE_EXPENSE` | Duplicate Expense | write | `action`, `id` |  |  |
| `BEXIO_GET_ADDITIONAL_ADDRESS` | Get Additional Address | read | `additional_address_id`, `contact_id` |  |  |
| `BEXIO_GET_BANK_ACCOUNT` | Get Bank Account | read | `bank_account_id` |  |  |
| `BEXIO_GET_BILL` | Get Bill | read | `id` |  |  |
| `BEXIO_GET_BUSINESS_YEAR` | Get Business Year | read | `business_year_id` |  |  |
| `BEXIO_GET_CALENDAR_YEAR` | Get Calendar Year | read | `calendar_year_id` |  |  |
| `BEXIO_GET_COMMENT` | Get Comment | read | `comment_id`, `document_id`, `kb_document_type` |  |  |
| `BEXIO_GET_COMPANY_PROFILE` | Get Company Profile | read | `profile_id` |  |  |
| `BEXIO_GET_CONTACT` | Get Contact | read | `contact_id`; *show_archived* |  |  |
| `BEXIO_GET_CONTACT_GROUP` | Get Contact Group | read | `contact_group_id` |  |  |
| `BEXIO_GET_CONTACT_RELATION` | Get Contact Relation | read | `contact_relation_id` |  |  |
| `BEXIO_GET_COUNTRY` | Get Country | read | `country_id` |  |  |
| `BEXIO_GET_CURRENCY` | Get Currency | read | `currency_id` |  |  |
| `BEXIO_GET_CURRENT_USER` | Get Current User | read |  |  |  |
| `BEXIO_GET_DEFAULT_POSITION` | Get Default Position | read | `document_id`, `kb_document_type`, `position_id` |  |  |
| `BEXIO_GET_DELIVERY` | Get Delivery | read | `delivery_id` |  |  |
| `BEXIO_GET_DISCOUNT_POSITION` | Get Discount Position | read | `document_id`, `kb_document_type`, `position_id` |  |  |
| `BEXIO_GET_EMPLOYEE` | Get Employee | read | `date`, `employee_id` |  |  |
| `BEXIO_GET_EMPLOYEE_ABSENCE` | Get Employee Absence | read | `absence_id`, `employee_id` |  |  |
| `BEXIO_GET_EXPENSE` | Get Expense | read | `id` |  |  |
| `BEXIO_GET_FICTIONAL_USER` | Get Fictional User | read | `fictional_user_id` |  |  |
| `BEXIO_GET_FILE` | Get File | read | `file_id` |  |  |
| `BEXIO_GET_FILE_USAGE` | Get File Usage | read | `file_id` |  |  |
| `BEXIO_GET_INVOICE` | Get Invoice | read | `invoice_id` |  |  |
| `BEXIO_GET_INVOICE_PAYMENT` | Get Invoice Payment | read | `invoice_id`, `payment_id` |  |  |
| `BEXIO_GET_INVOICE_PDF` | Get Invoice Pdf | read | `invoice_id`; *logopaper* |  |  |
| `BEXIO_GET_INVOICE_REMINDER` | Get Invoice Reminder | read | `invoice_id`, `reminder_id` |  |  |
| `BEXIO_GET_INVOICE_REMINDER_PDF` | Get Invoice Reminder Pdf | read | `invoice_id`, `reminder_id`; *logopaper* |  |  |
| `BEXIO_GET_ITEM` | Get Item | read | `article_id` |  |  |
| `BEXIO_GET_ITEM_POSITION` | Get Item Position | read | `document_id`, `kb_document_type`, `position_id` |  |  |
| `BEXIO_GET_MANUAL_COMPOUND_ENTRY_FILE` | Get Manual Compound Entry File | read | `file_id`, `manual_entry_id` |  |  |
| `BEXIO_GET_MANUAL_ENTRY_FILE` | Get Manual Entry File | read | `entry_id`, `file_id`, `manual_entry_id` |  |  |
| `BEXIO_GET_MILESTONE` | Get Milestone | read | `milestone_id`, `project_id` |  |  |
| `BEXIO_GET_NEXT_REFERENCE_NUMBER` | Get Next Reference Number | read |  |  |  |
| `BEXIO_GET_NOTE` | Get Note | read | `note_id` |  |  |
| `BEXIO_GET_ORDER` | Get Order | read | `order_id` |  |  |
| `BEXIO_GET_ORDER_PDF` | Get Order Pdf | read | `order_id`; *logopaper* |  |  |
| `BEXIO_GET_ORDER_REPETITION` | Get Order Repetition | read | `order_id` |  |  |
| `BEXIO_GET_OUTGOING_PAYMENT` | Get Outgoing Payment | read | `id` |  |  |
| `BEXIO_GET_PAGEBREAK_POSITION` | Get Pagebreak Position | read | `document_id`, `kb_document_type`, `position_id` |  |  |
| `BEXIO_GET_PAYMENT` | Get Payment | read | `payment_id` |  |  |
| `BEXIO_GET_PERMISSIONS` | Get Permissions | read |  |  |  |
| `BEXIO_GET_PROJECT` | Get Project | read | `project_id` |  |  |
| `BEXIO_GET_PURCHASE_ORDER` | Get Purchase Order | read | `purchase_order_id` |  |  |
| `BEXIO_GET_QUOTE` | Get Quote | read | `quote_id` |  |  |
| `BEXIO_GET_QUOTE_PDF` | Get Quote Pdf | read | `quote_id`; *logopaper* |  |  |
| `BEXIO_GET_SALUTATION` | Get Salutation | read | `salutation_id` |  |  |
| `BEXIO_GET_SUBPOSITION_POSITION` | Get Subposition Position | read | `document_id`, `kb_document_type`, `position_id` |  |  |
| `BEXIO_GET_SUBTOTAL_POSITION` | Get Subtotal Position | read | `document_id`, `kb_document_type`, `position_id` |  |  |
| `BEXIO_GET_TASK` | Get Task | read | `task_id` |  |  |
| `BEXIO_GET_TAX` | Get Tax | read | `tax_id` |  |  |
| `BEXIO_GET_TEXT_POSITION` | Get Text Position | read | `document_id`, `kb_document_type`, `position_id` |  |  |
| `BEXIO_GET_TIMESHEET` | Get Timesheet | read | `timesheet_id` |  |  |
| `BEXIO_GET_TITLE` | Get Title | read | `title_id` |  |  |
| `BEXIO_GET_UNIT` | Get Unit | read | `unit_id` |  |  |
| `BEXIO_GET_USER` | Get User | read | `user_id` |  |  |
| `BEXIO_GET_VAT_PERIOD` | Get Vat Period | read | `vat_period_id` |  |  |
| `BEXIO_GET_WORK_PACKAGE` | Get Work Package | read | `package_id`, `project_id` |  |  |
| `BEXIO_ISSUE_DELIVERY` | Issue Delivery | write | `delivery_id` |  |  |
| `BEXIO_ISSUE_INVOICE` | Issue Invoice | write | `invoice_id` |  |  |
| `BEXIO_ISSUE_QUOTE` | Issue Quote | write | `quote_id` |  |  |
| `BEXIO_LIST_ACCOUNTS` | List Accounts | read | *limit, next_cursor* |  |  |
| `BEXIO_LIST_ACCOUNT_GROUPS` | List Account Groups | read | *limit, next_cursor* |  |  |
| `BEXIO_LIST_ADDITIONAL_ADDRESSES` | List Additional Addresses | read | `contact_id`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_BANK_ACCOUNTS` | List Bank Accounts | read | *limit, next_cursor* |  |  |
| `BEXIO_LIST_BILLS` | List Bills | read | *average_exchange_rate_enabled, bill_date_end, bill_date_start, currency_code, document_no, due_date_end, due_date_start, fields and 16 more* |  |  |
| `BEXIO_LIST_BUSINESS_ACTIVITIES` | List Business Activities | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_BUSINESS_YEARS` | List Business Years | read | *limit, next_cursor* |  |  |
| `BEXIO_LIST_CALENDAR_YEARS` | List Calendar Years | read | *limit, next_cursor* |  |  |
| `BEXIO_LIST_COMMENTS` | List Comments | read | `document_id`, `kb_document_type`; *limit, next_cursor* |  |  |
| `BEXIO_LIST_COMMUNICATION_TYPES` | List Communication Types | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_COMPANY_PROFILE` | List Company Profile | read |  |  |  |
| `BEXIO_LIST_CONTACTS` | List Contacts | read | *limit, next_cursor, order_by, show_archived* |  |  |
| `BEXIO_LIST_CONTACT_GROUPS` | List Contact Groups | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_CONTACT_RELATIONS` | List Contact Relations | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_CONTACT_SECTORS` | List Contact Sectors | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_COUNTRIES` | List Countries | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_CURRENCIES` | List Currencies | read | *date, embed, limit, next_cursor* |  |  |
| `BEXIO_LIST_CURRENCY_CODES` | List Currency Codes | read |  |  |  |
| `BEXIO_LIST_DEFAULT_POSITIONS` | List Default Positions | read | `document_id`, `kb_document_type`; *limit, next_cursor* |  |  |
| `BEXIO_LIST_DELIVERIES` | List Deliveries | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_DISCOUNT_POSITIONS` | List Discount Positions | read | `document_id`, `kb_document_type`; *limit, next_cursor* |  |  |
| `BEXIO_LIST_DOCUMENT_SETTINGS` | List Document Settings | read | *order_by* |  |  |
| `BEXIO_LIST_DOCUMENT_TEMPLATE` | List Document Template | read |  |  |  |
| `BEXIO_LIST_EMPLOYEES` | List Employees | read |  |  |  |
| `BEXIO_LIST_EMPLOYEE_ABSENCES` | List Employee Absences | read | `business_year`, `employee_id` |  |  |
| `BEXIO_LIST_EXCHANGE_RATES_FOR_CURRENCY` | List Exchange Rates For Currency | read | `currency_id`; *date* |  |  |
| `BEXIO_LIST_EXPENSES` | List Expenses | read | *created_at_end, created_at_start, currency_code, document_no, gross_max, gross_min, limit, net_max and 10 more* |  |  |
| `BEXIO_LIST_FICTIONAL_USERS` | List Fictional Users | read | *limit, next_cursor* |  |  |
| `BEXIO_LIST_FILES` | List Files | read | *archived_state, next_cursor, order_by* |  |  |
| `BEXIO_LIST_INVOICES` | List Invoices | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_INVOICE_PAYMENTS` | List Invoice Payments | read | `invoice_id`; *limit, next_cursor* |  |  |
| `BEXIO_LIST_INVOICE_REMINDERS` | List Invoice Reminders | read | `invoice_id` |  |  |
| `BEXIO_LIST_ITEMS` | List Items | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_ITEM_POSITIONS` | List Item Positions | read | `document_id`, `kb_document_type`; *limit, next_cursor* |  |  |
| `BEXIO_LIST_JOURNAL_ENTRIES` | List Journal Entries | read | *account_uuid, from, limit, next_cursor, to* |  |  |
| `BEXIO_LIST_LANGUAGES` | List Languages | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_MANUAL_COMPOUND_ENTRY_FILES` | List Manual Compound Entry Files | read | `manual_entry_id`; *limit, next_cursor* |  |  |
| `BEXIO_LIST_MANUAL_ENTRIES` | List Manual Entries | read | *limit, next_cursor* |  |  |
| `BEXIO_LIST_MANUAL_ENTRY_FILES` | List Manual Entry Files | read | `entry_id`, `manual_entry_id`; *limit, next_cursor* |  |  |
| `BEXIO_LIST_MILESTONES` | List Milestones | read | `project_id`; *limit, next_cursor* |  |  |
| `BEXIO_LIST_NOTES` | List Notes | read | *limit, next_cursor* |  |  |
| `BEXIO_LIST_ORDERS` | List Orders | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_OUTGOING_PAYMENTS` | List Outgoing Payments | read | `bill_id`; *limit, next_cursor, order, sort* |  |  |
| `BEXIO_LIST_PAGEBREAK_POSITIONS` | List Pagebreak Positions | read | `document_id`, `kb_document_type`; *limit, next_cursor* |  |  |
| `BEXIO_LIST_PAYMENTS` | List Payments | read | *filter_by, limit, next_cursor* |  |  |
| `BEXIO_LIST_PAYMENT_TYPES` | List Payment Types | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_PROJECTS` | List Projects | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_PROJECT_STATUS` | List Project Status | read |  |  |  |
| `BEXIO_LIST_PROJECT_TYPE` | List Project Type | read | *order_by* |  |  |
| `BEXIO_LIST_PURCHASE_ORDERS` | List Purchase Orders | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_QUOTES` | List Quotes | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_SALUTATIONS` | List Salutations | read | *limit, next_cursor* |  |  |
| `BEXIO_LIST_STOCK_AREAS` | List Stock Areas | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_STOCK_LOCATIONS` | List Stock Locations | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_SUBPOSITION_POSITIONS` | List Subposition Positions | read | `document_id`, `kb_document_type`; *limit, next_cursor* |  |  |
| `BEXIO_LIST_SUBTOTAL_POSITIONS` | List Subtotal Positions | read | `document_id`, `kb_document_type`; *limit, next_cursor* |  |  |
| `BEXIO_LIST_TASKS` | List Tasks | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_TASK_PRIORITIES` | List Task Priorities | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_TASK_STATUS` | List Task Status | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_TAXES` | List Taxes | read | *date, limit, next_cursor, scope, types* |  |  |
| `BEXIO_LIST_TEXT_POSITIONS` | List Text Positions | read | `document_id`, `kb_document_type`; *limit, next_cursor* |  |  |
| `BEXIO_LIST_TIMESHEETS` | List Timesheets | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_TIME_SHEET_STATUS` | List Time Sheet Status | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_TITLES` | List Titles | read | *limit, next_cursor, order_by* |  |  |
| `BEXIO_LIST_UNITS` | List Units | read | *limit, next_cursor* |  |  |
| `BEXIO_LIST_USERS` | List Users | read | *limit, next_cursor* |  |  |
| `BEXIO_LIST_VAT_PERIODS` | List Vat Periods | read | *limit, next_cursor* |  |  |
| `BEXIO_LIST_WORK_PACKAGES` | List Work Packages | read | `project_id`; *limit, next_cursor* |  |  |
| `BEXIO_MARK_AS_SENT_INVOICE` | Mark As Sent Invoice | write | `invoice_id` |  |  |
| `BEXIO_MARK_AS_SENT_INVOICE_REMINDER` | Mark As Sent Invoice Reminder | write | `invoice_id`, `reminder_id` |  |  |
| `BEXIO_MARK_AS_SENT_QUOTE` | Mark As Sent Quote | write | `quote_id` |  |  |
| `BEXIO_MARK_AS_UNSENT_INVOICE_REMINDER` | Mark As Unsent Invoice Reminder | write | `invoice_id`, `reminder_id` |  |  |
| `BEXIO_PREVIEW_FILE` | Preview File | read | `file_id` |  |  |
| `BEXIO_REISSUE_QUOTE` | Reissue Quote | write | `quote_id` |  |  |
| `BEXIO_RESTORE_CONTACT` | Restore Contact | write | `contact_id` |  |  |
| `BEXIO_REVERT_ISSUE_INVOICE` | Revert Issue Invoice | write | `invoice_id` |  |  |
| `BEXIO_REVERT_ISSUE_QUOTE` | Revert Issue Quote | write | `quote_id` |  |  |
| `BEXIO_SEARCH_ACCOUNTS` | Search Accounts | read | `filters`; *limit, next_cursor* |  |  |
| `BEXIO_SEARCH_ADDITIONAL_ADDRESSES` | Search Additional Addresses | read | `contact_id`, `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_BUSINESS_ACTIVITIES` | Search Business Activities | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_CALENDAR_YEARS` | Search Calendar Years | read | `filters`; *limit, next_cursor* |  |  |
| `BEXIO_SEARCH_COMMUNICATION_TYPES` | Search Communication Types | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_CONTACT` | Search Contact | read | `filters`; *limit, next_cursor, order_by, show_archived* |  |  |
| `BEXIO_SEARCH_CONTACT_GROUPS` | Search Contact Groups | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_CONTACT_RELATIONS` | Search Contact Relations | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_CONTACT_SECTORS` | Search Contact Sectors | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_COUNTRIES` | Search Countries | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_FILES` | Search Files | read | `filters`; *archived_state, limit, next_cursor* |  |  |
| `BEXIO_SEARCH_INVOICES` | Search Invoices | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_ITEMS` | Search Items | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_LANGUAGES` | Search Languages | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_NOTES` | Search Notes | read | `filters`; *limit, next_cursor* |  |  |
| `BEXIO_SEARCH_ORDERS` | Search Orders | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_PAYMENT_TYPES` | Search Payment Types | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_PROJECTS` | Search Projects | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_QUOTES` | Search Quotes | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_REMINDERS` | Search Reminders | read | `filters`, `invoice_id` |  |  |
| `BEXIO_SEARCH_SALUTATIONS` | Search Salutations | read | `filters`; *limit, next_cursor* |  |  |
| `BEXIO_SEARCH_STOCK_AREAS` | Search Stock Areas | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_STOCK_LOCATIONS` | Search Stock Locations | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_TASKS` | Search Tasks | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_TIMESHEETS` | Search Timesheets | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_TITLES` | Search Titles | read | `filters`; *limit, next_cursor, order_by* |  |  |
| `BEXIO_SEARCH_UNITS` | Search Units | read | `filters`; *limit, next_cursor* |  |  |
| `BEXIO_SEND_INVOICE` | Send Invoice | write | `invoice_id`, `message`, `recipient_email`, `subject`; *attach_pdf, mark_as_open* |  |  |
| `BEXIO_SEND_INVOICE_REMINDER` | Send Invoice Reminder | write | `invoice_id`, `message`, `recipient_email`, `reminder_id`, `subject` |  |  |
| `BEXIO_SEND_QUOTE` | Send Quote | write | `message`, `quote_id`, `recipient_email`, `subject`; *attach_pdf, mark_as_open* |  |  |
| `BEXIO_UNARCHIVE_PROJECT` | Unarchive Project | write | `project_id` |  |  |
| `BEXIO_UPDATE_ADDITIONAL_ADDRESS` | Update Additional Address | write | `additional_address_id`, `contact_id`; *address_addition, city, country_id, description, house_number, name, name_addition, postcode and 2 more* |  |  |
| `BEXIO_UPDATE_BILL` | Update Bill | write | `address`, `attachment_ids`, `bill_date`, `contact_partner_id`, `currency_code`, `discounts`, `due_date`, `id`, `item_net`, `line_items`, `manual_amount`, `split_into_line_items`, `supplier_id`; *amount_calc, amount_man, base_currency_amount, document_no, exchange_rate, payment, title, vendor_ref* |  |  |
| `BEXIO_UPDATE_BILL_BOOKING_STATUS` | Update Bill Booking Status | write | `id`, `status` |  |  |
| `BEXIO_UPDATE_CONTACT` | Update Contact | write | `contact_id`, `contact_type_id`, `name_1`, `owner_id`, `user_id`; *address_addition, birthday, city, contact_branch_ids, contact_group_ids, contact_number, country_id, fax and 16 more* |  |  |
| `BEXIO_UPDATE_CONTACT_GROUP` | Update Contact Group | write | `contact_group_id`, `name` |  |  |
| `BEXIO_UPDATE_CONTACT_RELATION` | Update Contact Relation | write | `contact_id`, `contact_relation_id`, `contact_sub_id`; *description* |  |  |
| `BEXIO_UPDATE_COUNTRY` | Update Country | write | `country_id`, `iso3166_alpha2`, `name`, `name_short` |  |  |
| `BEXIO_UPDATE_CURRENCY` | Update Currency | write | `currency_id`, `data` |  |  |
| `BEXIO_UPDATE_DEFAULT_POSITION` | Update Default Position | write | `document_id`, `kb_document_type`, `position_id`; *account_id, amount, amount_completed, amount_open, amount_reserved, discount_in_percent, is_optional, tax_id and 3 more* |  |  |
| `BEXIO_UPDATE_DISCOUNT_POSITION` | Update Discount Position | write | `document_id`, `kb_document_type`, `position_id`; *is_percentual, text, value* |  |  |
| `BEXIO_UPDATE_EMPLOYEE` | Update Employee | write | `employee_id`; *address, ahv_number, annual_vacation_days, date_of_birth, email, first_name, gender, iban and 6 more* |  |  |
| `BEXIO_UPDATE_EMPLOYEE_ABSENCE` | Update Employee Absence | write | `absence_id`, `continued_pay`, `disability`, `employee_id`, `end_date`, `half_day`, `paid_hours`, `reason`, `start_date` |  |  |
| `BEXIO_UPDATE_EXPENSE` | Update Expense | write | `amount`, `attachment_ids`, `currency_code`, `id`, `paid_on`; *address, bank_account_id, base_currency_amount, booking_account_id, document_no, exchange_rate, supplier_id, tax_id and 1 more* |  |  |
| `BEXIO_UPDATE_EXPENSE_BOOKING_STATUS` | Update Expense Booking Status | write | `id`, `status` |  |  |
| `BEXIO_UPDATE_FICTIONAL_USER` | Update Fictional User | write | `email`, `fictional_user_id`, `firstname`, `lastname`, `salutation_type`; *title_id* |  |  |
| `BEXIO_UPDATE_FILE` | Update File | write | `file_id`; *is_archived, name, source_type* |  |  |
| `BEXIO_UPDATE_INVOICE` | Update Invoice | write | `invoice_id`; *api_reference, bank_account_id, contact_address_manual, contact_id, contact_sub_id, currency_id, document_nr, footer and 14 more* |  |  |
| `BEXIO_UPDATE_ITEM` | Update Item | write | `article_id`; *account_id, article_group_id, contact_id, currency_id, deliverer_code, deliverer_description, deliverer_name, delivery_price and 22 more* |  |  |
| `BEXIO_UPDATE_ITEM_POSITION` | Update Item Position | write | `document_id`, `kb_document_type`, `position_id`; *account_id, amount, amount_completed, amount_open, amount_reserved, article_id, discount_in_percent, is_optional and 4 more* |  |  |
| `BEXIO_UPDATE_MANUAL_ENTRY` | Update Manual Entry | write | `date`, `entries`, `manual_entry_id`, `type`; *reference_nr* |  |  |
| `BEXIO_UPDATE_MILESTONE` | Update Milestone | write | `milestone_id`, `name`, `project_id`; *comment, end_date, pr_parent_milestone_id* |  |  |
| `BEXIO_UPDATE_NOTE` | Update Note | write | `event_start`, `note_id`, `subject`, `user_id`; *contact_id, entry_id, info, module_id, pr_project_id* |  |  |
| `BEXIO_UPDATE_ORDER` | Update Order | write | `order_id`; *api_reference, bank_account_id, contact_address_manual, contact_id, contact_sub_id, currency_id, delivery_address_manual, delivery_address_type and 14 more* |  |  |
| `BEXIO_UPDATE_ORDER_REPETITION` | Update Order Repetition | write | `order_id`; *end, repetition, start* |  |  |
| `BEXIO_UPDATE_OUTGOING_PAYMENT` | Update Outgoing Payment | write | `amount`, `execution_date`, `is_salary_payment`, `payment_id`; *fee_type, message, receiver_city, receiver_country_code, receiver_house_no, receiver_iban, receiver_name, receiver_postcode and 2 more* |  |  |
| `BEXIO_UPDATE_PAGEBREAK_POSITION` | Update Pagebreak Position | write | `document_id`, `kb_document_type`, `position_id`; *pagebreak* |  |  |
| `BEXIO_UPDATE_PAYMENT` | Update Payment | write | `data`, `payment_id` |  |  |
| `BEXIO_UPDATE_PROJECT` | Update Project | write | `contact_id`, `name`, `pr_project_type_id`, `pr_state_id`, `project_id`, `user_id`; *comment, contact_sub_id, document_nr, end_date, pr_budget_type_amount, pr_budget_type_id, pr_invoice_type_amount, pr_invoice_type_id and 1 more* |  |  |
| `BEXIO_UPDATE_PURCHASE_ORDER` | Update Purchase Order | write | `purchase_order_id`; *api_reference, bank_account_id, contact_address_manual, contact_id, contact_sub_id, created_at, currency, currency_id and 31 more* |  |  |
| `BEXIO_UPDATE_QUOTE` | Update Quote | write | `quote_id`; *api_reference, bank_account_id, contact_address_manual, contact_id, contact_sub_id, currency_id, delivery_address_manual, delivery_address_type and 17 more* |  |  |
| `BEXIO_UPDATE_SALUTATION` | Update Salutation | write | `name`, `salutation_id` |  |  |
| `BEXIO_UPDATE_SUBPOSITION_POSITION` | Update Subposition Position | write | `document_id`, `kb_document_type`, `position_id`; *show_pos_nr, text* |  |  |
| `BEXIO_UPDATE_SUBTOTAL_POSITION` | Update Subtotal Position | write | `document_id`, `kb_document_type`, `position_id`; *text* |  |  |
| `BEXIO_UPDATE_TASK` | Update Task | write | `subject`, `task_id`, `user_id`; *communication_kind_id, contact_id, entry_id, finish_date, have_remember, info, module_id, pr_project_id and 5 more* |  |  |
| `BEXIO_UPDATE_TEXT_POSITION` | Update Text Position | write | `document_id`, `kb_document_type`, `position_id`; *show_pos_nr, text* |  |  |
| `BEXIO_UPDATE_TIMESHEET` | Update Timesheet | write | `allowable_bill`, `client_service_id`, `timesheet_id`, `tracking`, `user_id`; *charge, contact_id, estimated_time, pr_milestone_id, pr_package_id, pr_project_id, status_id, sub_contact_id and 1 more* |  |  |
| `BEXIO_UPDATE_TITLE` | Update Title | write | `name`, `title_id` |  |  |
| `BEXIO_UPDATE_UNIT` | Update Unit | write | `name`, `unit_id` |  |  |
| `BEXIO_UPDATE_WORK_PACKAGE` | Update Work Package | write | `name`, `package_id`, `project_id`; *comment, estimated_time_in_hours, pr_milestone_id, spent_time_in_hours* |  |  |
| `BEXIO_UPLOAD_MANUAL_COMPOUND_ENTRY_FILE` | Upload Manual Compound Entry File | write | `file`, `manual_entry_id` |  |  |
| `BEXIO_UPLOAD_MANUAL_ENTRY_FILE` | Upload Manual Entry File | write | `entry_id`, `file`, `manual_entry_id` |  |  |
