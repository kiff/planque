# Zoho Books: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

270 tools: 103 read, 167 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `ZOHO_BOOKS_ACCEPT_ESTIMATE` | Accept Estimate | write | `estimate_id`, `organization_id` |  | `ZohoBooks.estimates.CREATE` |
| `ZOHO_BOOKS_ACTIVATE_BANK_ACCOUNT` | Activate Bank Account | write | `account_id`, `organization_id` |  | `ZohoBooks.banking.CREATE` |
| `ZOHO_BOOKS_ADD_BILL_ATTACHMENT` | Add Bill Attachment | write | `attachment`, `bill_id`, `organization_id` |  | `ZohoBooks.bills.UPDATE`, `ZohoBooks.bills.ALL` |
| `ZOHO_BOOKS_ADD_BILL_COMMENT` | Add Bill Comment | write | `bill_id`, `description`, `organization_id` |  | `ZohoBooks.bills.CREATE`, `ZohoBooks.bills.UPDATE`, `ZohoBooks.bills.ALL` |
| `ZOHO_BOOKS_ADD_CONTACT_ADDRESS` | Add Contact Address | write | `contact_id`, `organization_id`; *address, attention, city, country, fax, phone, state, state_code and 2 more* |  | `ZohoBooks.contacts.CREATE`, `ZohoBooks.contacts.UPDATE`, `ZohoBooks.contacts.ALL` |
| `ZOHO_BOOKS_ADD_INVOICE_COMMENT` | Add Invoice Comment | write | `description`, `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.CREATE`, `ZohoBooks.invoices.UPDATE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_ADD_PROJECT_COMMENT` | Add Project Comment | write | `description`, `organization_id`, `project_id` |  | `ZohoBooks.projects.CREATE` |
| `ZOHO_BOOKS_ADD_PURCHASE_ORDER_COMMENT` | Add Purchase Order Comment | write | `description`, `organization_id`, `purchaseorder_id` |  | `ZohoBooks.purchaseorders.CREATE` |
| `ZOHO_BOOKS_ADD_SALES_ORDER_COMMENT` | Add Sales Order Comment | write | `description`, `organization_id`, `salesorder_id` |  | `ZohoBooks.salesorders.CREATE`, `ZohoBooks.salesorders.ALL` |
| `ZOHO_BOOKS_ADD_VENDOR_CREDIT_COMMENT` | Add Vendor Credit Comment | write | `description`, `organization_id`, `vendor_credit_id` |  | `ZohoBooks.debitnotes.CREATE` |
| `ZOHO_BOOKS_APPLY_CREDITS_TO_BILL` | Apply Credits to Bill | write | `bill_id`, `organization_id`, `vendor_credits` |  | `ZohoBooks.bills.CREATE`, `ZohoBooks.bills.ALL` |
| `ZOHO_BOOKS_APPLY_CREDITS_TO_INVOICE` | Apply Credits to Invoice | write | `creditnotes`, `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.CREATE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_APPLY_CREDIT_NOTE_TO_INVOICE` | Apply Credit Note to Invoice | write | `creditnote_id`, `invoices`, `organization_id` |  | `ZohoBooks.creditnotes.CREATE`, `ZohoBooks.creditnotes.ALL` |
| `ZOHO_BOOKS_BULK_EXPORT_INVOICES_PDF` | Bulk Export Invoices PDF | read | `invoice_ids`, `organization_id` |  | `ZohoBooks.invoices.READ` |
| `ZOHO_BOOKS_BULK_PRINT_INVOICES` | Bulk Print Invoices | read | `invoice_ids`, `organization_id` |  | `ZohoBooks.invoices.READ` |
| `ZOHO_BOOKS_CATEGORIZE_AS_CUSTOMER_PAYMENT_REFUND` | Categorize As Customer Payment Refund | write | `amount`, `date`, `from_account_id`, `organization_id`, `payment_id`, `transaction_id` |  | `ZohoBooks.banking.CREATE`, `ZohoBooks.banking.ALL` |
| `ZOHO_BOOKS_CATEGORIZE_UNCATEGORIZED_TRANSACTION` | Categorize Uncategorized Transaction | write | `organization_id`, `transaction_id`, `transaction_type`; *amount, bank_charges, currency_id, custom_fields, customer_id, date, description, documents and 14 more* |  | `ZohoBooks.banking.UPDATE`, `ZohoBooks.banking.ALL` |
| `ZOHO_BOOKS_CLONE_PROJECT` | Clone Project | write | `organization_id`, `project_id`, `project_name`; *description* |  | `ZohoBooks.projects.CREATE`, `ZohoBooks.projects.ALL` |
| `ZOHO_BOOKS_CONVERT_PURCHASE_ORDER_TO_BILL` | Convert Purchase Order To Bill | read | `organization_id`, `purchaseorder_ids` |  | `ZohoBooks.bills.READ`, `ZohoBooks.bills.ALL` |
| `ZOHO_BOOKS_CREATE_BANK_ACCOUNT` | Create Bank Account | write | `bankaccount`, `organization_id` |  | `ZohoBooks.banking.CREATE` |
| `ZOHO_BOOKS_CREATE_BANK_TRANSACTION` | Create Bank Transaction | write | `organization_id`, `transaction_type`; *amount, bank_charges, currency_id, custom_fields, customer_id, date, description, documents and 13 more* |  | `ZohoBooks.banking.CREATE`, `ZohoBooks.banking.UPDATE`, `ZohoBooks.banking.READ`, `ZohoBooks.banking.DELETE`, `ZohoBooks.banking.ALL` |
| `ZOHO_BOOKS_CREATE_CHART_OF_ACCOUNT` | Create Chart Of Account | write | `account_name`, `organization_id`; *account_code, account_type, can_show_in_ze, currency_id, description, include_in_vat_return, parent_account_id, show_on_dashboard* |  | `ZohoBooks.accountants.CREATE`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_CREATE_CONTACT` | Create Contact | write | `contact_name`, `organization_id`; *avatax_exempt_no, avatax_use_code, billing_address, company_name, contact_number, contact_persons, contact_type, country_code and 40 more* |  | `ZohoBooks.contacts.CREATE` |
| `ZOHO_BOOKS_CREATE_CONTACT_PERSON` | Create Contact Person | write | `contact_id`, `first_name`, `organization_id`; *department, designation, email, enable_portal, last_name, mobile, phone, salutation and 1 more* |  | `ZohoBooks.contacts.CREATE`, `ZohoBooks.contacts.ALL` |
| `ZOHO_BOOKS_CREATE_CURRENCY` | Create Currency | write | `currency_code`, `currency_format`, `organization_id`; *currency_symbol, price_precision* |  | `ZohoBooks.settings.CREATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_CREATE_CUSTOMER_PAYMENT_REFUND` | Create Customer Payment Refund | write | `amount`, `customer_payment_id`, `date`, `from_account_id`, `organization_id`; *description, exchange_rate, payment_form, reference_number, refund_mode* |  | `ZohoBooks.customerpayments.CREATE`, `ZohoBooks.customerpayments.ALL` |
| `ZOHO_BOOKS_CREATE_EMPLOYEE` | Create Employee | write | `email`, `employee_name`, `organization_id`; *cost_rate, role_id* |  | `ZohoBooks.expenses.CREATE`, `ZohoBooks.expenses.ALL` |
| `ZOHO_BOOKS_CREATE_EXCHANGE_RATE` | Create Exchange Rate | write | `currency_id`, `effective_date`, `organization_id`, `rate` |  | `ZohoBooks.settings.CREATE` |
| `ZOHO_BOOKS_CREATE_EXPENSE` | Create Expense | write | `account_id`, `amount`, `date`, `organization_id`, `paid_through_account_id`; *currency_id, customer_id, description, exchange_rate, is_billable, is_inclusive_tax, is_personal, project_id and 3 more* |  | `ZohoBooks.expenses.CREATE`, `ZohoBooks.expenses.ALL` |
| `ZOHO_BOOKS_CREATE_INVOICE_FROM_SALESORDER` | Create Invoice From Sales Order | write | `organization_id`, `salesorder_id` |  | `ZohoBooks.invoices.CREATE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_CREATE_ITEM` | Create Item | write | `name`, `organization_id`, `rate`; *asset_account_id, description, opening_stock, opening_stock_rate, purchase_account_id, purchase_rate, reorder_level, sales_account_id and 5 more* |  |  |
| `ZOHO_BOOKS_CREATE_JOURNAL` | Create Journal | write | `journal_date`, `line_items`, `organization_id`; *currency_id, custom_fields, exchange_rate, journal_type, notes, reference_number, tags* |  | `ZohoBooks.accountants.CREATE`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_CREATE_LOCATION` | Create Location | write | `address`, `location_name`, `organization_id`; *associated_series_ids, auto_number_generation_id, email, is_all_users_selected, parent_location_id, phone, tax_settings_id, type and 1 more* |  | `ZohoBooks.settings.CREATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_CREATE_PROJECT` | Create Project | write | `billing_type`, `customer_id`, `organization_id`, `project_name`, `user_id`; *budget_amount, budget_hours, budget_type, cost_budget_amount, description, rate* |  | `ZohoBooks.projects.CREATE`, `ZohoBooks.projects.ALL` |
| `ZOHO_BOOKS_CREATE_RECURRING_BILL` | Create Recurring Bill | write | `organization_id`, `recurrence_frequency`, `recurrence_name`, `repeat_every`, `start_date`, `vendor_id`; *currency_id, custom_fields, destination_of_supply, end_date, gst_no, gst_treatment, is_inclusive_tax, line_items and 9 more* |  | `ZohoBooks.bills.CREATE`, `ZohoBooks.bills.ALL` |
| `ZOHO_BOOKS_CREATE_RECURRING_EXPENSE` | Create Recurring Expense | write | `account_id`, `amount`, `organization_id`, `recurrence_frequency`, `recurrence_name`, `repeat_every`, `start_date`; *currency_id, customer_id, is_billable, is_inclusive_tax, paid_through_account_id, project_id, tax_id* |  | `ZohoBooks.expenses.CREATE` |
| `ZOHO_BOOKS_CREATE_RECURRING_INVOICE` | Create Recurring Invoice | write | `customer_id`, `organization_id`, `recurrence_frequency`, `recurrence_name`; *billing_address, contact_persons_associated, currency_id, custom_fields, discount, discount_type, email, end_date and 16 more* |  | `ZohoBooks.invoices.CREATE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_CREATE_REPORTING_TAG` | Create Reporting Tag | write | `entities`, `multi_preference_entities`, `organization_id`, `tag_name`; *description, is_mandatory* |  | `ZohoBooks.settings.CREATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_CREATE_RULE` | Create Bank Account Rule | write | `apply_to`, `criteria_type`, `criterion`, `organization_id`, `rule_name`, `target_account_id`; *account_id, customer_id, record_as, reference_number, tax_id* |  | `ZohoBooks.banking.CREATE`, `ZohoBooks.banking.ALL` |
| `ZOHO_BOOKS_CREATE_SALES_RECEIPT` | Create Sales Receipt | write | `customer_id`, `line_items`, `organization_id`, `payment_mode`; *can_send_in_mail, date, ignore_auto_number_generation, notes, receipt_number, tags, terms* |  | `ZohoBooks.invoices.CREATE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_CREATE_TIME_ENTRY` | Create Time Entry | write | `log_date`, `organization_id`, `project_id`, `task_id`, `user_id`; *begin_time, end_time, is_billable, log_time, notes* |  | `ZohoBooks.projects.CREATE` |
| `ZOHO_BOOKS_CREATE_USER` | Create User | write | `email`, `name`, `organization_id`, `role_id` |  | `ZohoBooks.settings.CREATE` |
| `ZOHO_BOOKS_CREATE_VENDOR_CREDIT` | Create Vendor Credit | write | `date`, `line_items`, `organization_id`, `vendor_id`; *bill_id, ignore_auto_number_generation, notes, reference_number, vendor_credit_number* |  | `ZohoBooks.debitnotes.CREATE`, `ZohoBooks.debitnotes.ALL` |
| `ZOHO_BOOKS_CREATE_VENDOR_PAYMENT` | Create Vendor Payment | write | `amount`, `date`, `organization_id`, `paid_through_account_id`, `vendor_id`; *bills, description, exchange_rate, payment_mode, reference_number* |  | `ZohoBooks.vendorpayments.CREATE`, `ZohoBooks.vendorpayments.ALL` |
| `ZOHO_BOOKS_DEACTIVATE_BANK_ACCOUNT` | Deactivate Bank Account | write | `account_id`, `organization_id` |  | `ZohoBooks.banking.CREATE` |
| `ZOHO_BOOKS_DELETE_BANK_ACCOUNT` | Delete Bank Account | write | `account_id`, `organization_id` |  | `ZohoBooks.banking.DELETE`, `ZohoBooks.banking.ALL` |
| `ZOHO_BOOKS_DELETE_BANK_TRANSACTION` | Delete Bank Transaction | write | `bank_transaction_id`, `organization_id` |  | `ZohoBooks.banking.DELETE`, `ZohoBooks.banking.ALL` |
| `ZOHO_BOOKS_DELETE_BILL` | Delete Bill | write | `bill_id`, `organization_id` |  | `ZohoBooks.bills.DELETE`, `ZohoBooks.bills.ALL` |
| `ZOHO_BOOKS_DELETE_BILL_ATTACHMENT` | Delete Bill Attachment | write | `bill_id`, `organization_id` |  | `ZohoBooks.bills.DELETE`, `ZohoBooks.bills.ALL` |
| `ZOHO_BOOKS_DELETE_BILL_COMMENT` | Delete Bill Comment | write | `bill_id`, `comment_id`, `organization_id` |  | `ZohoBooks.bills.DELETE`, `ZohoBooks.bills.ALL` |
| `ZOHO_BOOKS_DELETE_BILL_PAYMENT` | Delete Bill Payment | write | `bill_id`, `bill_payment_id`, `organization_id` |  | `ZohoBooks.bills.DELETE` |
| `ZOHO_BOOKS_DELETE_BULK_CUSTOMER_PAYMENTS` | Delete Bulk Customer Payments | write | `organization_id`, `payment_id` |  | `ZohoBooks.customerpayments.DELETE`, `ZohoBooks.customerpayments.ALL` |
| `ZOHO_BOOKS_DELETE_BULK_VENDOR_PAYMENTS` | Delete Bulk Vendor Payments | write | `organization_id`, `vendorpayment_id` |  | `ZohoBooks.vendorpayments.DELETE`, `ZohoBooks.vendorpayments.ALL` |
| `ZOHO_BOOKS_DELETE_CHART_OF_ACCOUNT` | Delete Chart of Account | write | `account_id`, `organization_id` |  | `ZohoBooks.accountants.DELETE`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_DELETE_CHART_OF_ACCOUNT_TRANSACTION` | Delete Chart Of Account Transaction | write | `organization_id`, `transaction_id` |  | `ZohoBooks.accountants.DELETE` |
| `ZOHO_BOOKS_DELETE_CONTACT` | Delete Contact | write | `contact_id`, `organization_id` |  | `ZohoBooks.contacts.DELETE`, `ZohoBooks.contacts.ALL` |
| `ZOHO_BOOKS_DELETE_CONTACT_ADDRESS` | Delete Contact Address | write | `address_id`, `contact_id`, `organization_id` |  | `ZohoBooks.contacts.DELETE`, `ZohoBooks.contacts.ALL` |
| `ZOHO_BOOKS_DELETE_CONTACT_PERSON` | Delete Contact Person | write | `contact_person_id`, `organization_id` |  | `ZohoBooks.contacts.DELETE` |
| `ZOHO_BOOKS_DELETE_CREDIT_NOTE` | Delete Credit Note | write | `creditnote_id`, `organization_id` |  | `ZohoBooks.creditnotes.DELETE`, `ZohoBooks.creditnotes.ALL` |
| `ZOHO_BOOKS_DELETE_CREDIT_NOTE_COMMENT` | Delete Credit Note Comment | write | `comment_id`, `creditnote_id`, `organization_id` |  | `ZohoBooks.creditnotes.DELETE`, `ZohoBooks.creditnotes.ALL` |
| `ZOHO_BOOKS_DELETE_CREDIT_NOTE_REFUND` | Delete Credit Note Refund | write | `creditnote_id`, `creditnote_refund_id`, `organization_id` |  | `ZohoBooks.creditnotes.DELETE`, `ZohoBooks.creditnotes.ALL` |
| `ZOHO_BOOKS_DELETE_CURRENCY` | Delete Currency | write | `currency_id`, `organization_id` |  | `ZohoBooks.settings.DELETE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_DELETE_CUSTOMER_PAYMENT` | Delete Customer Payment | write | `organization_id`, `payment_id` |  | `ZohoBooks.customerpayments.DELETE`, `ZohoBooks.customerpayments.ALL` |
| `ZOHO_BOOKS_DELETE_CUSTOMER_PAYMENT_REFUND` | Delete Customer Payment Refund | write | `customer_payment_id`, `organization_id`, `refund_id` |  | `ZohoBooks.customerpayments.DELETE`, `ZohoBooks.customerpayments.ALL` |
| `ZOHO_BOOKS_DELETE_EMPLOYEE` | Delete Employee | write | `employee_id`, `organization_id` |  | `ZohoBooks.expenses.DELETE` |
| `ZOHO_BOOKS_DELETE_ESTIMATE` | Delete Estimate | write | `estimate_id`, `organization_id` |  | `ZohoBooks.estimates.DELETE`, `ZohoBooks.estimates.ALL` |
| `ZOHO_BOOKS_DELETE_EXCHANGE_RATE` | Delete Exchange Rate | write | `currency_id`, `exchange_rate_id`, `organization_id` |  | `ZohoBooks.settings.DELETE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_DELETE_EXPENSE` | Delete Expense | write | `expense_id`, `organization_id` |  | `ZohoBooks.expenses.DELETE`, `ZohoBooks.expenses.ALL` |
| `ZOHO_BOOKS_DELETE_EXPENSE_RECEIPT` | Delete Expense Receipt | write | `expense_id`, `organization_id` |  | `ZohoBooks.expenses.DELETE`, `ZohoBooks.expenses.ALL` |
| `ZOHO_BOOKS_DELETE_INVOICE` | Delete Invoice | write | `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.DELETE` |
| `ZOHO_BOOKS_DELETE_INVOICE_ATTACHMENT` | Delete Invoice Attachment | write | `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.DELETE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_DELETE_INVOICE_COMMENT` | Delete Invoice Comment | write | `comment_id`, `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.DELETE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_DELETE_INVOICE_DOCUMENT` | Delete Invoice Document | write | `document_id`, `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.DELETE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_DELETE_INVOICE_PAYMENT` | Delete Invoice Payment | write | `invoice_id`, `invoice_payment_id`, `organization_id` |  | `ZohoBooks.invoices.DELETE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_DELETE_ITEM` | Delete Item | write | `item_id`, `organization_id` |  | `ZohoBooks.settings.CREATE`, `ZohoBooks.settings.UPDATE`, `ZohoBooks.settings.READ`, `ZohoBooks.settings.DELETE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_DELETE_JOURNALS` | Delete Journals | write | `journal_ids`, `organization_id` |  | `ZohoBooks.accountants.DELETE`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_DELETE_JOURNAL_COMMENT_V4` | Delete Journal Comment V4 | write | `comment_id`, `journal_id`, `organization_id`; *entity* |  | `ZohoBooks.accountants.DELETE`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_DELETE_JOURNAL_DOCUMENT` | Delete Journal Document | write | `document_id`, `journal_id`, `organization_id` |  | `ZohoBooks.accountants.DELETE`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_DELETE_JOURNAL_V4` | Delete Journal (v4 Beta) | write | `journal_id`, `organization_id` |  | `ZohoBooks.accountants.DELETE`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_DELETE_LOCATION` | Delete Location | write | `location_id`, `organization_id` |  | `ZohoBooks.settings.DELETE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_DELETE_PROJECT_COMMENT_V4` | Delete Project Comment (v4 Beta) | write | `comment_id`, `organization_id`, `project_id` |  | `ZohoBooks.projects.DELETE`, `ZohoBooks.projects.ALL` |
| `ZOHO_BOOKS_DELETE_PROJECT_TASK` | Delete Project Task | write | `organization_id`, `project_id`, `task_id` |  | `ZohoBooks.projects.DELETE`, `ZohoBooks.projects.ALL` |
| `ZOHO_BOOKS_DELETE_PROJECT_V4` | Delete Project V4 | write | `organization_id`, `project_id` |  | `ZohoBooks.projects.DELETE`, `ZohoBooks.projects.ALL` |
| `ZOHO_BOOKS_DELETE_PURCHASE_ORDER` | Delete Purchase Order | write | `organization_id`, `purchase_order_id` |  | `ZohoBooks.purchaseorders.DELETE`, `ZohoBooks.purchaseorders.ALL`, `ZohoBooks.fullaccess.all` |
| `ZOHO_BOOKS_DELETE_PURCHASE_ORDER_ATTACHMENT` | Delete Purchase Order Attachment | write | `organization_id`, `purchaseorder_id` |  | `ZohoBooks.purchaseorders.DELETE`, `ZohoBooks.purchaseorders.ALL` |
| `ZOHO_BOOKS_DELETE_PURCHASE_ORDER_COMMENT` | Delete Purchase Order Comment | write | `comment_id`, `organization_id`, `purchaseorder_id` |  | `ZohoBooks.purchaseorders.DELETE` |
| `ZOHO_BOOKS_DELETE_RECURRING_BILL` | Delete Recurring Bill | write | `organization_id`, `recurring_bill_id` |  | `ZohoBooks.bills.DELETE` |
| `ZOHO_BOOKS_DELETE_RECURRING_INVOICE` | Delete Recurring Invoice | write | `organization_id`, `recurring_invoice_id` |  | `ZohoBooks.invoices.DELETE` |
| `ZOHO_BOOKS_DELETE_RULE` | Delete Bank Account Rule | write | `organization_id`, `rule_id` |  | `ZohoBooks.banking.DELETE`, `ZohoBooks.banking.ALL` |
| `ZOHO_BOOKS_DELETE_SALES_ORDER_COMMENT` | Delete Sales Order Comment | write | `comment_id`, `organization_id`, `salesorder_id` |  | `ZohoBooks.salesorders.DELETE`, `ZohoBooks.salesorders.ALL` |
| `ZOHO_BOOKS_DELETE_SALES_RECEIPT` | Delete Sales Receipt | write | `organization_id`, `sales_receipt_id` |  | `ZohoBooks.invoices.DELETE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_DELETE_TIME_ENTRY` | Delete Time Entry | write | `organization_id`, `time_entry_id` |  | `ZohoBooks.projects.DELETE`, `ZohoBooks.projects.ALL`, `ZohoBooks.fullaccess.all` |
| `ZOHO_BOOKS_DELETE_VENDOR_CREDIT` | Delete Vendor Credit | write | `organization_id`, `vendor_credit_id` |  | `ZohoBooks.debitnotes.DELETE`, `ZohoBooks.debitnotes.ALL` |
| `ZOHO_BOOKS_DELETE_VENDOR_CREDIT_REFUND` | Delete Vendor Credit Refund | write | `organization_id`, `vendor_credit_id`, `vendor_credit_refund_id` |  | `ZohoBooks.debitnotes.DELETE` |
| `ZOHO_BOOKS_DELETE_VENDOR_PAYMENT` | Delete Vendor Payment | write | `organization_id`, `payment_id` |  | `ZohoBooks.vendorpayments.DELETE`, `ZohoBooks.vendorpayments.ALL` |
| `ZOHO_BOOKS_DISABLE_CONTACT_PAYMENT_REMINDER` | Disable Contact Payment Reminder | write | `contact_id`, `organization_id` |  | `ZohoBooks.contacts.CREATE`, `ZohoBooks.contacts.ALL` |
| `ZOHO_BOOKS_DISABLE_INVOICE_PAYMENT_REMINDER` | Disable Invoice Payment Reminder | write | `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.CREATE` |
| `ZOHO_BOOKS_EMAIL_CONTACT_STATEMENT` | Email Contact Statement | write | `contact_id`, `organization_id`; *body, cc_mail_ids, end_date, send_attachment, send_from_org_email_id, start_date, subject, to_mail_ids* |  | `ZohoBooks.contacts.CREATE`, `ZohoBooks.contacts.ALL` |
| `ZOHO_BOOKS_EMAIL_ESTIMATE` | Email Estimate | write | `estimate_id`, `organization_id`, `to_mail_ids`; *attachments, body, cc_mail_ids, from_address_id, send_attachment, send_from_org_email_id, subject* |  | `ZohoBooks.estimates.CREATE`, `ZohoBooks.estimates.ALL` |
| `ZOHO_BOOKS_EMAIL_INVOICE` | Email Invoice | write | `invoice_id`, `organization_id`, `to_mail_ids`; *attachments, body, cc_mail_ids, from_address_id, send_attachment, send_customer_statement, send_from_org_email_id, subject* |  | `ZohoBooks.invoices.CREATE` |
| `ZOHO_BOOKS_ENABLE_CONTACT_PAYMENT_REMINDER` | Enable Contact Payment Reminder | write | `contact_id`, `organization_id` |  | `ZohoBooks.contacts.UPDATE`, `ZohoBooks.contacts.ALL` |
| `ZOHO_BOOKS_ENABLE_CONTACT_PORTAL` | Enable Contact Portal | write | `contact_id`, `organization_id` |  | `ZohoBooks.contacts.CREATE` |
| `ZOHO_BOOKS_ENABLE_INVOICE_PAYMENT_REMINDER` | Enable Invoice Payment Reminder | write | `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.UPDATE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_EXCLUDE_BANK_TRANSACTION` | Exclude Bank Transaction | write | `organization_id`, `transaction_id`; *account_id* |  | `ZohoBooks.banking.CREATE` |
| `ZOHO_BOOKS_EXPORT_BULK_ESTIMATES_PDF` | Bulk Export Estimates PDF | read | `estimate_ids`, `organization_id` |  | `ZohoBooks.estimates.READ` |
| `ZOHO_BOOKS_EXPORT_SALES_ORDERS_PDF` | Bulk Export Sales Orders PDF | read | `organization_id`, `salesorder_ids` |  | `ZohoBooks.salesorders.READ`, `ZohoBooks.salesorders.ALL` |
| `ZOHO_BOOKS_GET_ACCOUNT_TRANSACTION_REPORT` | Get Account Transaction Report | read | `from_date`, `organization_id`, `to_date`; *filter_by, page, per_page* |  | `ZohoBooks.accountants.READ`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_GET_ALL_TAG_OPTIONS` | Get All Tag Options | read | `organization_id`, `tag_id`; *show_untagged* |  | `ZohoBooks.settings.READ`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_GET_BANK_ACCOUNT` | Get Bank Account | read | `account_id`, `organization_id` |  | `ZohoBooks.banking.READ` |
| `ZOHO_BOOKS_GET_BASE_CURRENCY_ADJUSTMENT` | Get Base Currency Adjustment | read | `base_currency_adjustment_id`, `organization_id` |  | `ZohoBooks.accountants.READ`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_GET_BILL` | Get Bill | read | `bill_id`, `organization_id` |  | `ZohoBooks.bills.READ` |
| `ZOHO_BOOKS_GET_BILL_ATTACHMENT` | Get Bill Attachment | read | `bill_id`, `organization_id`; *preview* |  | `ZohoBooks.bills.READ`, `ZohoBooks.bills.ALL` |
| `ZOHO_BOOKS_GET_CHART_OF_ACCOUNT` | Get Chart Of Account | read | `account_id`, `organization_id` |  | `ZohoBooks.accountants.READ`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_GET_CONTACT` | Get Contact | read | `contact_id`, `organization_id` |  |  |
| `ZOHO_BOOKS_GET_CONTACT_ADDRESS` | Get Contact Address | read | `contact_id`, `organization_id` |  | `ZohoBooks.contacts.READ`, `ZohoBooks.contacts.ALL` |
| `ZOHO_BOOKS_GET_CONTACT_PERSON` | Get Contact Person | read | `contact_id`, `contact_person_id`, `organization_id` |  | `ZohoBooks.contacts.READ`, `ZohoBooks.contacts.ALL` |
| `ZOHO_BOOKS_GET_CONTACT_STATEMENT_MAIL` | Get Contact Statement Mail | read | `contact_id`, `organization_id`; *end_date, start_date* |  | `ZohoBooks.contacts.READ`, `ZohoBooks.contacts.ALL` |
| `ZOHO_BOOKS_GET_CREDIT_NOTE` | Get Credit Note | read | `creditnote_id`, `organization_id`; *accept, print* |  | `ZohoBooks.creditnotes.READ`, `ZohoBooks.creditnotes.ALL` |
| `ZOHO_BOOKS_GET_CURRENCY` | Get Currency | read | `currency_id`, `organization_id` |  | `ZohoBooks.settings.READ`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_GET_DETAILED_GENERAL_LEDGER` | Get Detailed General Ledger | read | `from_date`, `organization_id`, `to_date`; *page, per_page* |  | `ZohoBooks.accountants.READ`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_GET_ESTIMATE` | Get Estimate | read | `estimate_id`, `organization_id` |  | `ZohoBooks.estimates.READ` |
| `ZOHO_BOOKS_GET_ESTIMATE_EMAIL` | Get Estimate Email | read | `estimate_id`, `organization_id`; *email_template_id* |  | `ZohoBooks.estimates.READ`, `ZohoBooks.estimates.ALL` |
| `ZOHO_BOOKS_GET_EXPENSE` | Get Expense | read | `expense_id`, `organization_id` |  | `ZohoBooks.expenses.READ`, `ZohoBooks.expenses.ALL` |
| `ZOHO_BOOKS_GET_FINANCIAL_STATEMENT` | Get Financial Statement | read | `from_date`, `organization_id`, `statement_type`, `to_date`; *filter_by* |  | `ZohoBooks.accountants.READ`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_GET_INVOICE` | Get Invoice | read | `invoice_id`, `organization_id`; *accept, print* |  | `ZohoBooks.invoices.READ`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_GET_INVOICE_ATTACHMENT` | Get Invoice Attachment | read | `invoice_id`, `organization_id`; *preview* |  | `ZohoBooks.invoices.READ`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_GET_INVOICE_EMAIL` | Get Invoice Email | read | `invoice_id`, `organization_id`; *email_template_id* |  | `ZohoBooks.invoices.READ` |
| `ZOHO_BOOKS_GET_ITEM` | Get Item | read | `item_id`, `organization_id` |  | `ZohoBooks.settings.CREATE` |
| `ZOHO_BOOKS_GET_JOURNALS_LIST_V4` | Get Journals List (v4 Beta) | read | `organization_id`; *date_end, date_start, filter_by, page, per_page, search_text, sort_column, status* |  | `ZohoBooks.accountants.READ`, `ZohoBooks.accountants.ALL`, `ZohoBooks.fullaccess.all` |
| `ZOHO_BOOKS_GET_JOURNAL_CREDITS_LIST_V4` | Get Journal Credits List (v4 Beta) | read | `journal_id`, `organization_id` |  | `ZohoBooks.accountants.READ`, `ZohoBooks.accountants.ALL`, `ZohoBooks.fullaccess.all` |
| `ZOHO_BOOKS_GET_JOURNAL_DETAILS_V4` | Get Journal Details V4 | read | `journal_id`, `organization-id`; *accept, include, language_code, print* |  | `ZohoBooks.accountants.READ`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_GET_LAST_IMPORTED_STATEMENT` | Get Last Imported Statement | read | `account_id`, `organization_id` |  | `ZohoBooks.banking.READ`, `ZohoBooks.banking.ALL` |
| `ZOHO_BOOKS_GET_MATCHING_BANK_TRANSACTIONS` | Get Matching Bank Transactions | read | `organization_id`, `transaction_id`; *amount_end, amount_start, contact, date_after, date_before, page, per_page, reference_number and 1 more* |  | `ZohoBooks.banking.READ` |
| `ZOHO_BOOKS_GET_OPENING_BALANCE` | Get Opening Balance | read | `organization_id` |  | `ZohoBooks.settings.READ`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_GET_ORGANIZATION` | Get Organization | read | `organization_id` |  | `ZohoBooks.settings.READ`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_GET_PAYMENT_REMINDER_MAIL_CONTENT_FOR_INVOICE` | Get Payment Reminder Mail Content | read | `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.READ` |
| `ZOHO_BOOKS_GET_PROJECT` | Get Project | read | `organization_id`, `project_id` |  | `ZohoBooks.projects.READ` |
| `ZOHO_BOOKS_GET_PROJECT_USER` | Get Project User | read | `organization_id`, `project_id`, `user_id` |  | `ZohoBooks.projects.READ`, `ZohoBooks.projects.ALL` |
| `ZOHO_BOOKS_GET_PURCHASE_ORDER` | Get Purchase Order | read | `organization_id`, `purchase_order_id` |  | `ZohoBooks.purchaseorders.READ`, `ZohoBooks.purchaseorders.ALL` |
| `ZOHO_BOOKS_GET_RECURRING_INVOICE` | Get Recurring Invoice | read | `organization_id`, `recurring_invoice_id` |  | `ZohoBooks.invoices.READ`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_GET_SALES_ORDER` | Get Sales Order | read | `organization_id`, `salesorder_id`; *accept, print* |  | `ZohoBooks.salesorders.READ` |
| `ZOHO_BOOKS_GET_SALES_ORDER_ATTACHMENT` | Get Sales Order Attachment | read | `organization_id`, `salesorder_id`; *inline, preview* |  | `ZohoBooks.salesorders.READ`, `ZohoBooks.salesorders.ALL` |
| `ZOHO_BOOKS_GET_SALES_ORDER_EMAIL_CONTENT` | Get Sales Order Email Content | read | `organization_id`, `salesorder_id`; *email_template_id* |  | `ZohoBooks.salesorders.READ` |
| `ZOHO_BOOKS_GET_TEMPLATE` | Get Template | read | `organization_id`, `template_id` |  | `ZohoBooks.settings.READ`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_GET_USER` | Get User | read | `organization_id`, `user_id` |  | `ZohoBooks.settings.CREATE`, `ZohoBooks.settings.UPDATE`, `ZohoBooks.settings.READ`, `ZohoBooks.settings.DELETE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_GET_VENDOR_CREDIT_REFUND` | Get Vendor Credit Refund | read | `organization_id`, `vendor_credit_id`, `vendor_credit_refund_id` |  | `ZohoBooks.debitnotes.READ`, `ZohoBooks.debitnotes.ALL` |
| `ZOHO_BOOKS_IMPORT_BANK_STATEMENTS` | Import Bank Statements | write | `account_id`, `organization_id`, `transactions`; *end_date, start_date* |  | `ZohoBooks.banking.CREATE`, `ZohoBooks.banking.UPDATE`, `ZohoBooks.banking.READ`, `ZohoBooks.banking.DELETE`, `ZohoBooks.banking.ALL` |
| `ZOHO_BOOKS_LIST_BANK_ACCOUNTS` | List Bank Accounts | read | `organization_id`; *filter_by, page, per_page, sort_column* |  | `ZohoBooks.banking.READ`, `ZohoBooks.banking.CREATE`, `ZohoBooks.banking.UPDATE`, `ZohoBooks.banking.DELETE`, `ZohoBooks.banking.ALL` |
| `ZOHO_BOOKS_LIST_BANK_ACCOUNT_STATEMENTS` | List Bank Account Statements | read | `account_id`, `organization_id`; *filter_by, page, per_page* |  | `ZohoBooks.banking.READ`, `ZohoBooks.banking.ALL` |
| `ZOHO_BOOKS_LIST_BANK_RULES` | List Bank Rules | read | `account_id`, `organization_id` |  | `ZohoBooks.banking.READ`, `ZohoBooks.fullaccess.all` |
| `ZOHO_BOOKS_LIST_BANK_TRANSACTIONS` | List Bank Transactions | read | `organization_id`; *account_id, amount_end, amount_start, date_end, date_start, filter_by, page, per_page and 5 more* |  | `ZohoBooks.banking.CREATE` |
| `ZOHO_BOOKS_LIST_BASE_CURRENCY_ADJUSTMENTS` | List Base Currency Adjustments | read | `organization_id`; *filter_by, last_modified_time, page, per_page, sort_column* |  | `ZohoBooks.accountants.READ` |
| `ZOHO_BOOKS_LIST_BILLS` | List Bills | read | `organization_id`; *bill_number, bill_number_contains, bill_number_startswith, date_after, date_before, date_end, date_start, description_contains and 17 more* |  | `ZohoBooks.bills.READ`, `ZohoBooks.bills.UPDATE`, `ZohoBooks.bills.DELETE`, `ZohoBooks.bills.CREATE`, `ZohoBooks.bills.ALL` |
| `ZOHO_BOOKS_LIST_BILL_COMMENTS_AND_HISTORY` | List Bill Comments and History | read | `bill_id`, `organization_id` |  | `ZohoBooks.bills.READ`, `ZohoBooks.bills.ALL` |
| `ZOHO_BOOKS_LIST_BILL_PAYMENTS` | List Bill Payments | read | `bill_id`, `organization_id` |  | `ZohoBooks.bills.READ`, `ZohoBooks.bills.ALL` |
| `ZOHO_BOOKS_LIST_CHART_OF_ACCOUNTS` | List Chart Of Accounts | read | `organization_id`; *filter_by, page, per_page, showbalance, sort_column* |  | `ZohoBooks.accountants.READ` |
| `ZOHO_BOOKS_LIST_CHART_OF_ACCOUNT_TRANSACTIONS` | List Chart of Account Transactions | read | `account_id`, `organization_id`; *amount_greater_equals, amount_greater_than, amount_less_equals, amount_less_than, date_after, date_before, date_end, date_start and 5 more* |  | `ZohoBooks.accountants.READ`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_LIST_CONTACTS` | List Contacts | read | `organization_id`; *address_contains, address_startswith, company_name_contains, company_name_startswith, contact_name_contains, contact_name_startswith, contact_type, email_contains and 15 more* |  | `ZohoBooks.contacts.READ` |
| `ZOHO_BOOKS_LIST_CONTACT_COMMENTS` | List Contact Comments | read | `contact_id`, `organization_id`; *page, per_page* |  | `ZohoBooks.contacts.READ`, `ZohoBooks.contacts.ALL` |
| `ZOHO_BOOKS_LIST_CONTACT_PERSONS` | List Contact Persons | read | `organization_id`; *page, per_page* |  | `ZohoBooks.contacts.READ`, `ZohoBooks.contacts.ALL` |
| `ZOHO_BOOKS_LIST_CONTACT_REFUNDS` | List Contact Refunds | read | `contact_id`, `organization_id`; *page, per_page* |  | `ZohoBooks.contacts.READ`, `ZohoBooks.creditnotes.READ` |
| `ZOHO_BOOKS_LIST_CREDIT_NOTES` | List Credit Notes | read | `organization_id`; *creditnote_number, customer_id, customer_name, date, filter_by, item_description, item_id, item_name and 9 more* |  | `ZohoBooks.creditnotes.READ` |
| `ZOHO_BOOKS_LIST_CREDIT_NOTE_REFUNDS` | List Credit Note Refunds | read | `organization_id`; *customer_id, page, per_page, sort_column* |  | `ZohoBooks.creditnotes.READ`, `ZohoBooks.creditnotes.ALL` |
| `ZOHO_BOOKS_LIST_CURRENCIES` | List Currencies | read | `organization_id`; *filter_by, page, per_page* |  | `ZohoBooks.settings.READ` |
| `ZOHO_BOOKS_LIST_CUSTOMER_PAYMENTS` | List Customer Payments | read | `organization_id`; *amount, customer_id, customer_name, date, filter_by, notes, page, payment_mode and 4 more* |  | `ZohoBooks.customerpayments.READ`, `ZohoBooks.customerpayments.ALL` |
| `ZOHO_BOOKS_LIST_CUSTOMER_PAYMENT_REFUNDS` | List Customer Payment Refunds | read | `customer_payment_id`, `organization_id`; *page, per_page* |  | `ZohoBooks.customerpayments.READ`, `ZohoBooks.customerpayments.ALL` |
| `ZOHO_BOOKS_LIST_EMPLOYEES` | List Employees | read | `organization_id`; *page, per_page* |  | `ZohoBooks.expenses.READ`, `ZohoBooks.expenses.ALL` |
| `ZOHO_BOOKS_LIST_ESTIMATES` | List Estimates | read | `organization_id`; *customer_id, customer_name, customer_name_contains, date, date_end, date_start, estimate_number, filter_by and 8 more* |  | `ZohoBooks.estimates.READ`, `ZohoBooks.estimates.ALL`, `ZohoBooks.fullaccess.all` |
| `ZOHO_BOOKS_LIST_ESTIMATE_COMMENTS_AND_HISTORY` | List Estimate Comments & History | read | `estimate_id`, `organization_id` |  | `ZohoBooks.estimates.READ`, `ZohoBooks.estimates.ALL` |
| `ZOHO_BOOKS_LIST_ESTIMATE_TEMPLATES` | List Estimate Templates | read | `organization_id` |  | `ZohoBooks.estimates.READ` |
| `ZOHO_BOOKS_LIST_EXPENSES` | List Expenses | read | `organization_id`; *account_name, account_name_contains, account_name_startswith, amount, amount_greater_equals, amount_greater_than, amount_less_equals, amount_less_than and 27 more* |  | `ZohoBooks.expenses.READ`, `ZohoBooks.expenses.ALL` |
| `ZOHO_BOOKS_LIST_EXPENSE_COMMENTS` | List Expense Comments | read | `expense_id`, `organization_id` |  | `ZohoBooks.expenses.READ`, `ZohoBooks.expenses.ALL` |
| `ZOHO_BOOKS_LIST_FIXED_ASSETS` | List Fixed Assets | read | `organization_id`; *filter_by, page, per_page, sort_column, sort_order* |  | `ZohoBooks.fixedasset.READ`, `ZohoBooks.fixedasset.ALL` |
| `ZOHO_BOOKS_LIST_FIXED_ASSET_TYPES` | List Fixed Asset Types | read | `organization_id`; *page, per_page* |  | `ZohoBooks.fixedasset.READ`, `ZohoBooks.fixedasset.ALL` |
| `ZOHO_BOOKS_LIST_INVOICES` | List Invoices | read | `organization_id`; *customer_id, customer_name, date_end, date_start, due_date_end, due_date_start, filter_by, invoice_number and 9 more* |  | `ZohoBooks.invoices.READ`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_LIST_INVOICE_COMMENTS_AND_HISTORY` | List Invoice Comments and History | read | `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.READ`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_LIST_INVOICE_CREDITS_APPLIED` | List Invoice Credits Applied | read | `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.READ`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_LIST_INVOICE_PAYMENTS` | List Invoice Payments | read | `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.READ`, `ZohoBooks.customerpayments.READ`, `ZohoBooks.invoices.DELETE` |
| `ZOHO_BOOKS_LIST_INVOICE_TEMPLATES` | List Invoice Templates | read | `organization_id` |  | `ZohoBooks.invoices.READ`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_LIST_ITEMS` | List Items | read | `organization_id`; *account_id, description, description_contains, description_startswith, filter_by, is_taxable, name, name_contains and 16 more* |  | `ZohoBooks.settings.READ`, `ZohoBooks.settings.CREATE`, `ZohoBooks.settings.UPDATE`, `ZohoBooks.settings.DELETE` |
| `ZOHO_BOOKS_LIST_ITEM_DETAILS` | List Item Details | read | `item_ids`, `organization_id` |  | `ZohoBooks.settings.READ`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_LIST_JOURNAL_TEMPLATES` | List Journal Templates | read | `organization_id`; *page, per_page* |  | `ZohoBooks.accountants.READ`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_LIST_LOCATIONS` | List Locations | read | `organization_id` |  | `ZohoBooks.settings.READ` |
| `ZOHO_BOOKS_LIST_ORGANIZATIONS` | List Organizations | read |  |  | `ZohoBooks.settings.READ` |
| `ZOHO_BOOKS_LIST_PROJECTS` | List Projects | read | `organization_id`; *customer_id, filter_by, page, per_page, sort_column* |  | `ZohoBooks.projects.READ`, `ZohoBooks.projects.ALL` |
| `ZOHO_BOOKS_LIST_PROJECT_USERS` | List Project Users | read | `organization_id`, `project_id` |  | `ZohoBooks.projects.READ`, `ZohoBooks.projects.ALL` |
| `ZOHO_BOOKS_LIST_PURCHASE_ORDERS` | List Purchase Orders | read | `organization_id`; *custom_field, custom_field_contains, custom_field_startswith, date, filter_by, item_description, item_description_contains, item_description_startswith and 24 more* |  | `ZohoBooks.purchaseorders.READ`, `ZohoBooks.purchaseorders.ALL` |
| `ZOHO_BOOKS_LIST_RECURRING_BILL_HISTORY` | List Recurring Bill History | read | `organization_id`, `recurring_bill_id` |  | `ZohoBooks.bills.READ`, `ZohoBooks.bills.ALL` |
| `ZOHO_BOOKS_LIST_RECURRING_INVOICE_HISTORY` | List Recurring Invoice History | read | `organization_id`, `recurring_invoice_id` |  | `ZohoBooks.invoices.READ` |
| `ZOHO_BOOKS_LIST_REPORTING_TAGS` | List Reporting Tags | read | `organization_id` |  | `ZohoBooks.settings.READ` |
| `ZOHO_BOOKS_LIST_RETAINER_INVOICES` | List Retainer Invoices | read | `organization_id`; *filter_by, page, per_page, print, sort_column, sort_order* |  | `ZohoBooks.invoices.READ`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_LIST_SALES_ORDERS` | List Sales Orders | read | `organization_id`; *filter_by, page, per_page, search_text, sort_column* |  | `ZohoBooks.salesorders.READ`, `ZohoBooks.salesorders.CREATE`, `ZohoBooks.salesorders.UPDATE`, `ZohoBooks.salesorders.DELETE`, `ZohoBooks.salesorders.ALL` |
| `ZOHO_BOOKS_LIST_SALES_ORDER_COMMENTS_AND_HISTORY` | List Sales Order Comments & History | read | `organization_id`, `salesorder_id` |  | `ZohoBooks.salesorders.READ`, `ZohoBooks.salesorders.ALL` |
| `ZOHO_BOOKS_LIST_SALES_RECEIPTS` | List Sales Receipts | read | `organization_id`; *customer_id, date, date_after, date_before, date_end, date_start, filter_by, item_name and 13 more* |  | `ZohoBooks.invoices.READ`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_LIST_TASKS` | List Tasks | read | `organization_id`, `project_id`; *page, per_page* |  | `ZohoBooks.projects.READ`, `ZohoBooks.projects.ALL` |
| `ZOHO_BOOKS_LIST_TAXES` | List Taxes | read | `organization_id`; *page, per_page* |  | `ZohoBooks.settings.READ`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_LIST_TAX_EXEMPTIONS` | List Tax Exemptions | read | `organization_id` |  | `ZohoBooks.settings.READ`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_LIST_USERS` | List Users | read | `organization_id`; *filter_by, page, per_page, sort_column* |  |  |
| `ZOHO_BOOKS_LIST_VENDOR_CREDITS` | List Vendor Credits | read | `organization_id`; *custom_field, custom_field_contains, custom_field_startswith, customer_id, customer_name, customer_name_contains, customer_name_startswith, date and 36 more* |  | `ZohoBooks.debitnotes.READ`, `ZohoBooks.debitnotes.ALL` |
| `ZOHO_BOOKS_LIST_VENDOR_CREDIT_REFUNDS` | List Vendor Credit Refunds | read | `organization_id`; *customer_id, last_modified_time, page, per_page, sort_column* |  | `ZohoBooks.debitnotes.READ`, `ZohoBooks.debitnotes.ALL` |
| `ZOHO_BOOKS_LIST_VENDOR_CREDIT_REFUNDS2` | List Vendor Credit Refunds | read | `organization_id`, `vendor_credit_id`; *page, per_page* |  | `ZohoBooks.debitnotes.READ`, `ZohoBooks.debitnotes.ALL` |
| `ZOHO_BOOKS_LIST_VENDOR_PAYMENTS` | List Vendor Payments | read | `organization_id`; *amount, amount_greater_equals, amount_greater_than, amount_less_equals, amount_less_than, bill_id, date, date_after and 28 more* |  | `ZohoBooks.vendorpayments.READ`, `ZohoBooks.vendorpayments.ALL` |
| `ZOHO_BOOKS_LIST_VENDOR_PAYMENT_REFUNDS` | List Vendor Payment Refunds | read | `organization_id`, `payment_id`; *page, per_page* |  | `ZohoBooks.vendorpayments.READ`, `ZohoBooks.vendorpayments.ALL` |
| `ZOHO_BOOKS_MARK_BILL_OPEN` | Mark Bill Open | write | `bill_id`, `organization_id` |  | `ZohoBooks.bills.CREATE`, `ZohoBooks.bills.UPDATE`, `ZohoBooks.bills.ALL` |
| `ZOHO_BOOKS_MARK_BILL_VOID` | Mark Bill Void | write | `bill_id`, `organization_id` |  | `ZohoBooks.bills.CREATE`, `ZohoBooks.bills.ALL` |
| `ZOHO_BOOKS_MARK_CONTACT_AS_ACTIVE` | Mark Contact as Active | write | `contact_id`, `organization_id` |  | `ZohoBooks.contacts.CREATE` |
| `ZOHO_BOOKS_MARK_CONTACT_AS_INACTIVE` | Mark Contact as Inactive | write | `contact_id`, `organization_id` |  | `ZohoBooks.contacts.CREATE` |
| `ZOHO_BOOKS_MARK_ESTIMATE_AS_SENT` | Mark Estimate As Sent | write | `estimate_id`, `organization_id` |  | `ZohoBooks.estimates.UPDATE`, `ZohoBooks.estimates.ALL` |
| `ZOHO_BOOKS_MARK_INVOICE_AS_SENT` | Mark Invoice As Sent | write | `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.UPDATE` |
| `ZOHO_BOOKS_MARK_ITEM_AS_ACTIVE` | Mark Item as Active | write | `item_id`, `organization_id` |  | `ZohoBooks.settings.UPDATE`, `ZohoBooks.settings.CREATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_MARK_ITEM_INACTIVE` | Mark Item as Inactive | write | `item_id`, `organization_id` |  | `ZohoBooks.settings.UPDATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_MARK_LOCATION_ACTIVE` | Mark Location Active | write | `location_id`, `organization_id` |  | `ZohoBooks.settings.CREATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_MARK_LOCATION_INACTIVE` | Mark Location as Inactive | write | `location_id`, `organization_id` |  | `ZohoBooks.settings.CREATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_MARK_LOCATION_PRIMARY` | Mark Location Primary | write | `location_id`, `organization_id` |  | `ZohoBooks.settings.CREATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_MARK_PROJECT_ACTIVE` | Mark Project Active | write | `organization_id`, `project_id` |  | `ZohoBooks.projects.CREATE`, `ZohoBooks.projects.ALL` |
| `ZOHO_BOOKS_MARK_REPORTING_TAG_AS_ACTIVE` | Mark Reporting Tag as Active | write | `organization_id`, `tag_id` |  | `ZohoBooks.settings.UPDATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_MARK_REPORTING_TAG_AS_INACTIVE` | Mark Reporting Tag as Inactive | write | `organization_id`, `tag_id` |  | `ZohoBooks.settings.UPDATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_MARK_REPORTING_TAG_DEFAULT_OPTION` | Mark Reporting Tag Default Option | write | `default_option_id`, `organization_id`, `tag_id` |  | `ZohoBooks.settings.READ` |
| `ZOHO_BOOKS_MARK_RETAINER_INVOICE_SENT` | Mark Retainer Invoice Sent | write | `organization_id`, `retainerinvoice_id` |  | `ZohoBooks.invoices.CREATE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_MARK_SALES_ORDER_VOID` | Mark Sales Order Void | write | `organization_id`, `salesorder_id` |  | `ZohoBooks.salesorders.CREATE`, `ZohoBooks.salesorders.ALL` |
| `ZOHO_BOOKS_MARK_TAG_OPTION_ACTIVE` | Mark Tag Option Active | write | `option_id`, `organization_id`, `tag_id` |  | `ZohoBooks.settings.UPDATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_MARK_TAG_OPTION_AS_INACTIVE` | Mark Tag Option as Inactive | write | `option_id`, `organization_id`, `tag_id` |  | `ZohoBooks.settings.UPDATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_MARK_USER_INACTIVE` | Mark User as Inactive | write | `organization_id`, `user_id` |  | `ZohoBooks.settings.CREATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_MARK_VENDOR_CREDIT_VOID` | Mark Vendor Credit Void | write | `organization_id`, `vendor_credit_id` |  | `ZohoBooks.debitnotes.CREATE` |
| `ZOHO_BOOKS_OPEN_SALES_ORDER` | Open Sales Order | write | `organization_id`, `salesorder_id` |  |  |
| `ZOHO_BOOKS_PRINT_BULK_ESTIMATES` | Bulk Print Estimates | read | `estimate_ids`, `organization_id` |  | `ZohoBooks.estimates.READ`, `ZohoBooks.estimates.ALL` |
| `ZOHO_BOOKS_PRINT_BULK_SALES_ORDERS` | Bulk Print Sales Orders | read | `organization_id`, `salesorder_ids` |  | `ZohoBooks.salesorders.READ`, `ZohoBooks.salesorders.ALL` |
| `ZOHO_BOOKS_REMIND_CUSTOMER_INVOICE` | Send Payment Reminder | write | `invoice_id`, `organization_id`; *body, cc_mail_ids, send_customer_statement, send_from_org_email_id, subject, to_mail_ids* |  | `ZohoBooks.invoices.CREATE` |
| `ZOHO_BOOKS_RESUME_RECURRING_INVOICE` | Resume Recurring Invoice | write | `organization_id`, `recurring_invoice_id` |  | `ZohoBooks.invoices.CREATE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_SEND_BULK_INVOICE_REMINDER` | Send Bulk Invoice Reminder | write | `invoice_ids`, `organization_id` |  | `ZohoBooks.invoices.CREATE` |
| `ZOHO_BOOKS_SEND_CONTACT_EMAIL` | Send Contact Email | write | `body`, `contact_id`, `organization_id`, `subject`, `to_mail_ids`; *cc_mail_ids, send_attachment, send_customer_statement, send_from_org_email_id* |  | `ZohoBooks.contacts.CREATE` |
| `ZOHO_BOOKS_STOP_RECURRING_INVOICE` | Stop Recurring Invoice | write | `organization_id`, `recurring_invoice_id` |  | `ZohoBooks.invoices.CREATE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_UNTRACK_CONTACT1099` | Untrack Contact 1099 | write | `contact_id`, `organization_id` |  | `ZohoBooks.contacts.UPDATE`, `ZohoBooks.contacts.ALL` |
| `ZOHO_BOOKS_UPDATE_BANK_TRANSACTION` | Update Bank Transaction | write | `bank_transaction_id`, `organization_id`, `transaction_type`; *amount, bank_charges, currency_id, custom_fields, customer_id, date, description, documents and 14 more* |  | `ZohoBooks.banking.UPDATE`, `ZohoBooks.banking.ALL`, `ZohoBooks.fullaccess.all` |
| `ZOHO_BOOKS_UPDATE_BILL` | Update Vendor Bill | write | `bill_id`, `bill_number`, `organization_id`, `vendor_id`; *adjustment, adjustment_description, approvers, attachment, currency_id, custom_fields, date, destination_of_supply and 22 more* |  | `ZohoBooks.bills.READ`, `ZohoBooks.bills.UPDATE` |
| `ZOHO_BOOKS_UPDATE_CONTACT` | Update Contact | write | `contact_id`, `contact_name`, `organization_id`; *avatax_exempt_no, avatax_use_code, billing_address, company_name, contact_number, contact_persons, contact_type, country_code and 37 more* |  | `ZohoBooks.contacts.UPDATE`, `ZohoBooks.contacts.ALL` |
| `ZOHO_BOOKS_UPDATE_CONTACT_PERSON` | Update Contact Person | write | `contact_id`, `contact_person_id`, `first_name`, `organization_id`; *department, designation, email, enable_portal, last_name, mobile, phone, salutation and 1 more* |  | `ZohoBooks.contacts.UPDATE`, `ZohoBooks.contacts.ALL` |
| `ZOHO_BOOKS_UPDATE_CREDIT_NOTE_REFUND` | Update Credit Note Refund | write | `creditnote_id`, `creditnote_refund_id`, `organization_id`; *amount, date, description, from_account_id, reference_number, refund_mode* |  | `ZohoBooks.creditnotes.UPDATE`, `ZohoBooks.creditnotes.ALL` |
| `ZOHO_BOOKS_UPDATE_CURRENCY` | Update Currency | write | `currency_code`, `currency_format`, `currency_id`, `organization_id`; *currency_symbol, price_precision* |  | `ZohoBooks.settings.UPDATE` |
| `ZOHO_BOOKS_UPDATE_CUSTOM_FIELDS_IN_ITEM` | Update Custom Fields in Item | write | `custom_fields`, `item_id`, `organization_id` |  | `ZohoBooks.settings.UPDATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_UPDATE_ESTIMATE` | Update Estimate | write | `estimate_id`, `organization_id`; *adjustment, adjustment_description, contact_persons_associated, currency_id, custom_body, custom_fields, custom_subject, customer_id and 22 more* |  | `ZohoBooks.estimates.UPDATE`, `ZohoBooks.estimates.ALL` |
| `ZOHO_BOOKS_UPDATE_ESTIMATE_BILLING_ADDRESS` | Update Estimate Billing Address | write | `estimate_id`, `organization_id`; *address, attention, city, country, fax, phone, state, street2 and 1 more* |  | `ZohoBooks.estimates.UPDATE`, `ZohoBooks.estimates.ALL` |
| `ZOHO_BOOKS_UPDATE_ESTIMATE_SHIPPING_ADDRESS` | Update Estimate Shipping Address | write | `estimate_id`, `organization_id`; *address, attention, city, country, fax, phone, state, state_code and 2 more* |  | `ZohoBooks.estimates.UPDATE`, `ZohoBooks.estimates.ALL` |
| `ZOHO_BOOKS_UPDATE_ESTIMATE_TEMPLATE` | Update Estimate Template | write | `estimate_id`, `organization_id`, `template_id` |  | `ZohoBooks.estimates.UPDATE`, `ZohoBooks.estimates.ALL` |
| `ZOHO_BOOKS_UPDATE_INVOICE` | Update Invoice | write | `customer_id`, `invoice_id`, `line_items`, `organization_id`; *currency_id, date, discount, discount_type, due_date, exchange_rate, ignore_auto_number_generation, invoice_number and 6 more* |  | `ZohoBooks.invoices.UPDATE` |
| `ZOHO_BOOKS_UPDATE_INVOICE_ATTACHMENT_PREFERENCE` | Update Invoice Attachment Preference | write | `can_send_in_mail`, `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.UPDATE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_UPDATE_INVOICE_BILLING_ADDRESS` | Update Invoice Billing Address | write | `invoice_id`, `organization_id`; *address, city, country, fax, state, street2, zip* |  | `ZohoBooks.invoices.UPDATE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_UPDATE_INVOICE_SHIPPING_ADDRESS` | Update Invoice Shipping Address | write | `invoice_id`, `organization_id`; *address, city, country, fax, state, street2, zip* |  | `ZohoBooks.invoices.UPDATE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_UPDATE_INVOICE_TEMPLATE` | Update Invoice Template | write | `invoice_id`, `organization_id`, `template_id` |  | `ZohoBooks.invoices.UPDATE`, `ZohoBooks.invoices.ALL` |
| `ZOHO_BOOKS_UPDATE_ITEM` | Update Item | write | `item_id`, `name`, `organization_id`, `rate`; *account_id, avatax_tax_code, avatax_use_code, custom_fields, description, hsn_or_sac, inventory_account_id, is_taxable and 18 more* |  | `ZohoBooks.settings.UPDATE` |
| `ZOHO_BOOKS_UPDATE_JOURNAL` | Update Journal | write | `journal_id`, `organization_id`; *currency_id, custom_fields, exchange_rate, journal_date, journal_type, line_items, notes, reference_number and 1 more* |  | `ZohoBooks.accountants.UPDATE`, `ZohoBooks.accountants.ALL` |
| `ZOHO_BOOKS_UPDATE_LOCATION` | Update Location | write | `address`, `location_id`, `location_name`, `organization_id`; *associated_series_ids, auto_number_generation_id, email, is_all_users_selected, parent_location_id, phone, type, user_ids* |  | `ZohoBooks.settings.UPDATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_UPDATE_ORGANIZATION` | Update Organization | write | `organization_id`; *contact_name, date_format, email, field_separator, fiscal_year_start_month, is_logo_uploaded, language_code, name and 2 more* |  | `ZohoBooks.settings.UPDATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_UPDATE_PROJECT` | Update Project | write | `billing_type`, `customer_id`, `organization_id`, `project_id`, `project_name`, `user_id`; *budget_amount, budget_hours, budget_type, cost_budget_amount, currency_id, description, rate, tasks and 1 more* |  | `ZohoBooks.projects.UPDATE`, `ZohoBooks.projects.ALL`, `ZohoBooks.fullaccess.all` |
| `ZOHO_BOOKS_UPDATE_PROJECT_USER` | Update Project User | write | `organization_id`, `project_id`, `user_id`; *budget_hours, rate, user_role* |  | `ZohoBooks.projects.UPDATE`, `ZohoBooks.projects.ALL` |
| `ZOHO_BOOKS_UPDATE_PURCHASE_ORDER_BILLING_ADDRESS` | Update Purchase Order Billing Address | write | `address`, `city`, `country`, `organization_id`, `purchaseorder_id`, `state`, `zip`; *attention, fax, phone, street2* |  | `ZohoBooks.purchaseorders.UPDATE` |
| `ZOHO_BOOKS_UPDATE_PURCHASE_ORDER_COMMENT` | Update Purchase Order Comment | write | `comment_id`, `description`, `organization_id`, `purchaseorder_id` |  | `ZohoBooks.purchaseorders.UPDATE`, `ZohoBooks.purchaseorders.ALL` |
| `ZOHO_BOOKS_UPDATE_RECURRING_BILL` | Update Recurring Bill | write | `organization_id`, `recurring_bill_id`; *currency_id, custom_fields, destination_of_supply, discount, end_date, gst_no, gst_treatment, is_inclusive_tax and 15 more* |  | `ZohoBooks.bills.UPDATE` |
| `ZOHO_BOOKS_UPDATE_RECURRING_INVOICE_TEMPLATE` | Update Recurring Invoice Template | write | `organization_id`, `recurring_invoice_id`, `template_id` |  | `ZohoBooks.invoices.UPDATE` |
| `ZOHO_BOOKS_UPDATE_REPORTING_TAG` | Update Reporting Tag | write | `organization_id`, `tag_id`, `tag_name`; *description, entities, is_mandatory, multi_preference_entities* |  | `ZohoBooks.settings.UPDATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_UPDATE_REPORTING_TAG_OPTIONS` | Update Reporting Tag Options | write | `options`, `organization_id`, `tag_id` |  | `ZohoBooks.settings.UPDATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_UPDATE_SALES_ORDER` | Update Sales Order | write | `organization_id`, `salesorder_id`; *contact_persons_associated, currency_id, custom_fields, customer_id, date, line_items, shipment_date, tax_exemption_code and 1 more* |  | `ZohoBooks.salesorders.UPDATE` |
| `ZOHO_BOOKS_UPDATE_SALES_ORDER_ATTACHMENT_PREFERENCE` | Update Sales Order Attachment Preference | write | `can_send_in_mail`, `organization_id`, `salesorder_id` |  | `ZohoBooks.salesorders.UPDATE`, `ZohoBooks.salesorders.ALL` |
| `ZOHO_BOOKS_UPDATE_SALES_ORDER_BILLING_ADDRESS` | Update Sales Order Billing Address | write | `organization_id`, `salesorder_id`; *address, attention, city, country, fax, state, street2, zip* |  | `ZohoBooks.salesorders.UPDATE`, `ZohoBooks.salesorders.ALL` |
| `ZOHO_BOOKS_UPDATE_SALES_ORDER_SHIPPING_ADDRESS` | Update Sales Order Shipping Address | write | `organization_id`, `salesorder_id`; *address, attention, city, country, fax, state, street2, zip* |  | `ZohoBooks.salesorders.UPDATE`, `ZohoBooks.salesorders.ALL` |
| `ZOHO_BOOKS_UPDATE_SALES_ORDER_TEMPLATE` | Update Sales Order Template | write | `organization_id`, `salesorder_id`, `template_id` |  | `ZohoBooks.salesorders.UPDATE`, `ZohoBooks.salesorders.ALL` |
| `ZOHO_BOOKS_UPDATE_USER` | Update User | write | `email`, `name`, `organization_id`, `user_id`; *role_id* |  | `ZohoBooks.settings.UPDATE`, `ZohoBooks.settings.ALL` |
| `ZOHO_BOOKS_UPDATE_VENDOR_PAYMENT_REFUND` | Update Vendor Payment Refund | write | `amount`, `date`, `organization_id`, `payment_id`, `vendorpayment_refund_id`; *description, exchange_rate, reference_number, refund_mode* |  | `ZohoBooks.vendorpayments.UPDATE`, `ZohoBooks.vendorpayments.ALL` |
| `ZOHO_BOOKS_VOID_INVOICE` | Void Invoice | write | `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.CREATE` |
| `ZOHO_BOOKS_WRITE_OFF_INVOICE` | Write Off Invoice | write | `invoice_id`, `organization_id` |  | `ZohoBooks.invoices.CREATE`, `ZohoBooks.invoices.ALL` |
