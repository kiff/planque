# Moneybird: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

32 tools: 17 read, 15 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `MONEYBIRD_ADD_NOTE_TO_CONTACT` | Add Note to Contact | write | `administration_id`, `contact_id`, `note` |  |  |
| `MONEYBIRD_ARCHIVE_CONTACT` | Archive Contact | write | `administration_id`, `contact_id` |  |  |
| `MONEYBIRD_CREATE_CONTACT` | Create Moneybird Contact | write | `administration_id`; *address1, address2, bank_account, chamber_of_commerce, city, company_name, country, custom_fields_attributes and 20 more* |  |  |
| `MONEYBIRD_CREATE_CONTACT_PERSON` | Create Contact Person | write | `administration_id`, `contact_id`, `contact_person` |  |  |
| `MONEYBIRD_CREATE_EXTERNAL_SALES_INVOICE_PAYMENT` | Create External Sales Invoice Payment | write | `administration_id`, `external_sales_invoice_id`, `payment_date`, `price`; *financial_account_id, financial_mutation_id, invoice_id, ledger_account_id, manual_payment_action, price_base, transaction_identifier* |  | `sales_invoices`, `bank`, `documents` |
| `MONEYBIRD_CREATE_SALES_INVOICE` | Create Sales Invoice | write | `administration_id`, `sales_invoice` |  |  |
| `MONEYBIRD_DELETE_CONTACT` | Delete Contact | write | `administration_id`, `contact_id` |  |  |
| `MONEYBIRD_DELETE_CONTACT_NOTE` | Delete Contact Note | write | `administration_id`, `contact_id`, `note_id` |  |  |
| `MONEYBIRD_DELETE_CONTACT_PERSON` | Delete Contact Person | write | `administration_id`, `contact_id`, `contact_person_id` |  |  |
| `MONEYBIRD_DELETE_DOCUMENT` | Delete Document | write | `administration_id`, `document_id`, `document_type` |  | `documents` |
| `MONEYBIRD_DOWNLOAD_DOCUMENT_ATTACHMENT` | Download Document Attachment | read | `administration_id`, `attachment_id`, `document_id`, `document_type` |  | `documents` |
| `MONEYBIRD_FILTER_CONTACTS` | Filter Contacts | read | `administration_id`; *created_after, first_name, include_archived, last_name, per_page, updated_after* |  |  |
| `MONEYBIRD_GET_ADDITIONAL_CHARGES` | Get Additional Charges | read | `administration_id`, `contact_id`; *include_billed* |  |  |
| `MONEYBIRD_GET_CONTACT` | Get Contact | read | `administration_id`, `contact_id`; *include_archived* |  |  |
| `MONEYBIRD_GET_CONTACT_BY_CUSTOMER_ID` | Get Contact by Customer ID | read | `administration_id`, `customer_id` |  |  |
| `MONEYBIRD_GET_CONTACT_PERSON` | Get Contact Person | read | `administration_id`, `contact_id`, `contact_person_id` |  |  |
| `MONEYBIRD_GET_DOCUMENT` | Get Document | read | `administration_id`, `document_id`, `document_type` |  | `documents` |
| `MONEYBIRD_GET_EXTERNAL_SALES_INVOICE` | Get External Sales Invoice | read | `administration_id`, `external_sales_invoice_id` |  | `sales_invoices` |
| `MONEYBIRD_GET_SALES_INVOICE` | Get Sales Invoice | read | `administration_id`, `invoice_id` |  |  |
| `MONEYBIRD_GET_TAX_REPORT` | Get Tax Report | read | `administration_id`; *granularity, period* |  | `documents`, `sales_invoices` |
| `MONEYBIRD_LIST_ADMINISTRATIONS` | List Administrations | read |  |  |  |
| `MONEYBIRD_LIST_CONTACTS_SYNCHRONIZATION` | List Contacts Synchronization | read | `administration_id`; *filter, include_archived* |  |  |
| `MONEYBIRD_LIST_FINANCIAL_ACCOUNTS` | List Financial Accounts | read | `administration_id` |  | `settings` |
| `MONEYBIRD_LIST_LEDGER_ACCOUNTS` | List Ledger Accounts | read | `administration_id` |  | `settings`, `sales_invoices`, `documents`, `estimates` |
| `MONEYBIRD_LIST_PURCHASE_INVOICES` | List Purchase Invoices | read | `administration_id`; *filter, page, per_page, sort* |  | `documents` |
| `MONEYBIRD_LIST_SALES_INVOICES` | List Sales Invoices | read | `administration_id`; *filter, page, per_page* |  |  |
| `MONEYBIRD_LIST_TAX_RATES` | List Tax Rates | read | `administration_id`; *filter, page, per_page* |  | `settings`, `sales_invoices`, `documents`, `estimates` |
| `MONEYBIRD_REQUEST_CONTACT_PAYMENTS_MANDATE_URL` | Request Payments Mandate URL | write | `administration_id`, `contact_id`; *identity_id* |  |  |
| `MONEYBIRD_UPDATE_CONTACT` | Update Contact | write | `administration_id`, `contact`, `contact_id` |  |  |
| `MONEYBIRD_UPDATE_CONTACT_PERSON` | Update Contact Person | write | `administration_id`, `contact_id`, `contact_person_id`; *department, email, firstname, lastname, phone* |  |  |
| `MONEYBIRD_UPDATE_SALES_INVOICE` | Update Sales Invoice | write | `administration_id`, `sales_invoice`, `sales_invoice_id` |  |  |
| `MONEYBIRD_UPLOAD_DOCUMENT_ATTACHMENT` | Upload Document Attachment | write | `administration_id`, `document_id`, `document_type`, `file` |  | `documents` |
