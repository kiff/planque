# QuickBooks: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

114 tools: 76 read, 38 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `QUICKBOOKS_CAPTURE_CHARGE` | Capture Charge | write | `charge_id`; *amount, context* |  | `com.intuit.quickbooks.payment` |
| `QUICKBOOKS_CREATE_ACCOUNT` | Create Account | write | `account_type`, `name`; *account_number, account_sub_type, active, currency_ref, description, opening_balance, opening_balance_date, parent_ref and 1 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_BANK_ACCOUNT` | Create Bank Account | write | `account_number`, `account_type`, `customer_id`, `name`, `routing_number`; *phone* |  | `com.intuit.quickbooks.payment` |
| `QUICKBOOKS_CREATE_BILL` | Create Bill | write | `Line`, `VendorRef`; *APAccountRef, CurrencyRef, DepartmentRef, DocNumber, DueDate, ExchangeRate, GlobalTaxCalculation, PrivateNote and 3 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_BILL_PAYMENT` | Create Bill Payment | write | `Line`, `PayType`, `TotalAmt`, `VendorRef`; *CheckPayment, CreditCardPayment, DepartmentRef, PrivateNote, TxnDate, minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_CLASS` | Create Class | write | `name`; *Active, ParentRef, SubClass* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_CREDIT_MEMO` | Create Credit Memo | write | `customer_ref`, `line`; *apply_tax_after_discount, bill_addr, bill_email, class_ref, currency_ref, custom_field, customer_memo, department_ref and 12 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_CUSTOMER` | Create Customer | write | *Active, AlternatePhone, Balance, BillAddr, BillWithParent, CompanyName, CurrencyRef, Fax and 19 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_DEPOSIT` | Create Deposit | write | `deposit_to_account_ref`, `line`; *currency_ref, department_ref, exchange_rate, minorversion, private_note, txn_date* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_ECHECK_PAYMENT` | Create eCheck Payment | write | `amount`, `bankAccount`, `currency`; *context* |  | `com.intuit.quickbooks.payment` |
| `QUICKBOOKS_CREATE_EMPLOYEE` | Create Employee | write | `family_name`, `given_name`; *active, alternate_phone, bill_rate, billable_time, birth_date, company_name, cost_rate, custom_field and 22 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_ESTIMATE` | Create Estimate | write | `customer_id`, `lines`; *apply_tax_after_discount, bill_addr, bill_email, class_ref, currency_ref, custom_field, customer_memo, department_ref and 12 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_INVOICE` | Create Invoice | write | `customer_id`, `lines`; *allow_ipn_payment, allow_online_ach_payment, allow_online_credit_card_payment, allow_online_payment, apply_tax_after_discount, bill_addr, bill_email, bill_email_bcc and 26 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_ITEM` | Create Item | write | `name`, `type`; *active, asset_account_ref, description, expense_account_ref, income_account_ref, inv_start_date, minorversion, parent_ref and 6 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_JOURNAL_ENTRY` | Create Journal Entry | write | `Line`; *Adjustment, CurrencyRef, DocNumber, ExchangeRate, PrivateNote, TxnDate, minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_PAYMENT` | Create Payment | write | `customer_id`, `lines`, `total_amt`; *ar_account_ref_id, credit_card_payment, currency_ref_value, customer_name, deposit_to_account_ref_id, exchange_rate, minorversion, payment_method_ref_id and 4 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_PURCHASE` | Create Purchase | write | `AccountRef`, `Line`, `PaymentType`; *CurrencyRef, DepartmentRef, DocNumber, EntityRef, ExchangeRate, GlobalTaxCalculation, PaymentMethodRef, PrivateNote and 4 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_PURCHASE_ORDER` | Create Purchase Order | write | `ap_account_id`, `lines`, `total_amt`, `vendor_id`; *class_ref, currency_ref, custom_field, department_ref, doc_number, due_date, exchange_rate, expected_date and 13 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_REFUND_RECEIPT` | Create Refund Receipt | write | `CustomerRef`, `DepositToAccountRef`, `Line`; *CurrencyRef, CustomField, DocNumber, ExchangeRate, PaymentRefNum, PrivateNote, TxnDate, minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_SALES_RECEIPT` | Create Sales Receipt | write | `customer_ref`, `line`; *apply_tax_after_discount, bill_addr, bill_email, bill_email_bcc, bill_email_cc, class_ref, currency_ref, custom_field and 16 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_TAXSERVICE_TAXCODE` | Create Tax Code via TaxService | write | `TaxCode`, `TaxRateDetails` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_TAX_AGENCY` | Create Tax Agency | write | `DisplayName` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_TIME_ACTIVITY` | Create Time Activity | write | `name_of`; *billable_status, break_hours, break_minutes, class_ref, customer_ref, department_ref, description, employee_ref and 10 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_VENDOR` | Create Vendor | write | *acct_num, active, alt_contact_name, alternate_phone, ap_account_ref, balance, bill_addr, bill_rate and 36 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CREATE_VENDOR_CREDIT` | Create Vendor Credit | write | `Line`, `VendorRef`; *APAccountRef, CurrencyRef, DepartmentRef, DocNumber, ExchangeRate, GlobalTaxCalculation, Memo, PrivateNote and 2 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CUSTOMER_BALANCE_DETAIL` | Customer Balance Detail | read | *aging_method, arpaid, columns, custom1, customer_ids, department_ids, end_duedate, report_date and 5 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_CUSTOMER_BALANCE_REPORT` | Customer Balance Report | read | *accounting_method, arpaid, customer_ids, date_macro, department_ids, qzurl, report_date, sort_order and 1 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_DELETE_BANK_ACCOUNT` | Delete Bank Account | write | `bank_account_id`, `customer_id`; *request_id* |  | `com.intuit.quickbooks.payment` |
| `QUICKBOOKS_DELETE_ITEM` | Delete Item | write | `item_id`, `sync_token`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_EXECUTE_BATCH_OPERATION` | Execute Batch Operation | write | `BatchItemRequest`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_AGED_RECEIVABLES_REPORT` | Aged Receivables Report | read | *aging_method, aging_period, arpaid, customer_ids, date_macro, minorversion, num_periods, report_date* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_ATTACHABLE` | Get Attachable | read | `attachable_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_BALANCE_SHEET_REPORT` | Get Balance Sheet Report | read | *accounting_method, date_macro, department_ids, report_date, sort_order, summarize_column_by* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_BANK_ACCOUNT` | Get Bank Account | read | `bank_account_id`, `customer_id` |  | `com.intuit.quickbooks.payment` |
| `QUICKBOOKS_GET_BILL` | Get Bill | read | `bill_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_BILL_PAYMENT` | Get Bill Payment | read | `bill_payment_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_CHANGED_ENTITIES` | Get Changed Entities | read | `changedSince`, `entities` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_COMPANY_INFO` | Get Company Info | read | *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_CREDIT_MEMO` | Get Credit Memo | read | `credit_memo_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_CREDIT_MEMO_PDF` | Get Credit Memo PDF | read | `credit_memo_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_DEPARTMENT` | Get Department | read | `department_id` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_DEPOSIT` | Get Deposit | read | `deposit_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_ESTIMATE` | Get Estimate | read | `estimate_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_ESTIMATE_PDF` | Get Estimate PDF | read | `estimate_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_EXCHANGE_RATE` | Get Exchange Rate | read | `asofdate`, `sourcecurrencycode`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_GENERAL_LEDGER_REPORT` | Get General Ledger Report | read | *account_ids, account_type, accounting_method, class_ids, columns, customer_ids, date_macro, department_ids and 7 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_INVENTORY_VALUATION_SUMMARY` | Get Inventory Valuation Summary | read | *accounting_method, date_macro, department_ids, report_date, sort_order, summarize_column_by* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_INVOICE_PDF` | Get Invoice PDF | read | `invoice_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_ITEM` | Get Item | read | `item_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_JOURNAL_ENTRY` | Get Journal Entry | read | `journal_entry_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_PAYMENT` | Get Payment | read | `payment_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_PAYMENT_METHOD` | Get Payment Method | read | `payment_method_id` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_PAYMENT_PDF` | Get Payment PDF | read | `payment_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_PREFERENCES` | Get Company Preferences | read | *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_PROFIT_AND_LOSS_DETAIL_REPORT` | Get Profit and Loss Detail Report | read | *accounting_method, class_ids, customer_ids, date_macro, department_ids, end_date, report_date, sort_order and 3 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_PROFIT_AND_LOSS_REPORT` | Get Profit and Loss Report | read | *accounting_method, class_ids, customer_ids, date_macro, department_ids, end_date, report_date, sort_order and 3 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_PURCHASE` | Get Purchase | read | `purchase_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_PURCHASE_ORDER` | Get Purchase Order | read | `purchase_order_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_PURCHASE_ORDER_PDF` | Get Purchase Order PDF | read | `purchase_order_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_REFUND_RECEIPT` | Get Refund Receipt | read | `refund_receipt_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_REFUND_RECEIPT_PDF` | Get Refund Receipt PDF | read | `refund_receipt_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_REPORTS` | Get Reports | read | `report_name`; *accounting_method, cleared, columns, customer, date_macro, department_ids, end_date, group_by and 6 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_REPORT_ACCOUNT_LIST` | Get Report - Account List | read | *columns, sort_order* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_REPORT_AGED_PAYABLES` | Get Report - Aged Payables | read | *aging_method, aging_period, appaid, date_macro, minorversion, num_periods, report_date, vendor_ids* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_REPORT_AGED_PAYABLE_DETAIL` | Get Report - Aged Payable Detail | read | *accounting_method, aging_method, aging_period, appaid, columns, date_macro, department_ids, duedate_macro and 8 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_REPORT_AGED_RECEIVABLE_DETAIL` | Get Report - Aged Receivable Detail | read | *accounting_method, aging_method, aging_period, arpaid, columns, custom1, customer_ids, date_macro and 10 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_REPORT_CASH_FLOW` | Get Report - Cash Flow | read | *accounting_method, date_macro, department_ids, report_date, sort_order, summarize_column_by* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_REPORT_CLASS_SALES` | Get Report - Class Sales | read | *accounting_method, class_ids, date_macro, department_ids, report_date, sort_order, summarize_column_by* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_REPORT_CUSTOMER_INCOME` | Get Report - Customer Income | read | *accounting_method, customer_ids, date_macro, department_ids, end_date, report_date, sort_order, start_date and 1 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_REPORT_CUSTOMER_SALES` | Get Report - Customer Sales | read | *accounting_method, customer_ids, date_macro, department_ids, end_date, report_date, sort_order, start_date and 1 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_REPORT_DEPARTMENT_SALES` | Get Report - Department Sales | read | *accounting_method, date_macro, department_ids, end_date, sort_order, start_date, summarize_column_by* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_REPORT_ITEM_SALES` | Get Report - Item Sales | read | *accounting_method, date_macro, department_ids, item_ids, report_date, sort_order, summarize_column_by* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_REPORT_TRIAL_BALANCE` | Get Report - Trial Balance | read | *accounting_method, date_macro, department_ids, end_date, report_date, sort_order, start_date, summarize_column_by* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_SALESRECEIPT_PDF` | Get Sales Receipt PDF | read | `sales_receipt_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_SALES_RECEIPT` | Get Sales Receipt | read | `sales_receipt_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_TAX_RATE` | Get Tax Rate | read | `tax_rate_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_TERM` | Get Term | read | `term_id` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_TIME_ACTIVITY` | Get Time Activity | read | `time_activity_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_TRANSACTION_LIST_REPORT` | Get Transaction List Report | read | *account_ids, accounting_method, class_ids, cleared, columns, customer_ids, date_macro, department_ids and 6 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_TRANSFER` | Get Transfer | read | `transfer_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_VENDOR_CREDIT` | Get Vendor Credit | read | `vendor_credit_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_GET_VENDOR_EXPENSES_REPORT` | Get Vendor Expenses Report | read | *accounting_method, appaid, date_macro, department_ids, qzurl, report_date, sort_order, summarize_column_by and 1 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_LIST_CARDS` | List Cards | read | `customer_id` |  | `com.intuit.quickbooks.payment` |
| `QUICKBOOKS_LIST_DEPARTMENTS` | List Departments | read | *query* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_LIST_INVOICES` | List Invoices | read | *fields, max_results, start_position* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_QUERY_ACCOUNT` | Query Account Entities | read | `query` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_QUERY_ACCOUNTS` | Query Account Entities | read | `query` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_QUERY_ATTACHABLE` | Query Attachable Entities | read | *query* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_QUERY_BILLS` | Query Bills | read | `query` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_QUERY_CUSTOMERS` | Query Customer Entities | read | `query` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_QUERY_ENTITIES` | Query QuickBooks Entities | read | `query` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_QUERY_ESTIMATES` | Query Estimate Entities | read | `query` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_QUERY_INVOICES` | Query Invoices | read | *customer_id, end_date, max_results, query, start_date, start_position, status* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_QUERY_ITEMS` | Query Item Entities | read | `query` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_QUERY_PAYMENTS` | Query Payment Entities | read | `query` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_QUERY_VENDORS` | Query Vendor Entities | read | `query` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_READ_ACCOUNT` | Read Account | read | `account_id` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_READ_CLASS` | Read Class | read | `class_id`; *minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_READ_CUSTOMER` | Read Customer | read | `customer_id` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_READ_EMPLOYEE` | Read Employee | read | `employee_id` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_READ_INVOICE` | Read Invoice | read | *invoice_id, minorversion* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_READ_VENDOR` | Read Vendor | read | `vendor_id` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_SEND_CREDIT_MEMO` | Send Credit Memo | write | `credit_memo_id`; *minorversion, send_to* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_UPDATE_ATTACHABLE` | Update Attachable | write | `Id`, `SyncToken`; *AttachableRef, Category, Lat, Long, Note, PlaceName, Tag* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_UPDATE_COMPANY_INFO` | Update Company Info | write | `CompanyAddr`, `Id`, `SyncToken`; *CompanyName, CompanyStartDate, Country, CustomerCommunicationAddr, CustomerCommunicationEmailAddr, Email, FiscalYearStartMonth, LegalAddr and 5 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_UPDATE_DEPARTMENT` | Update Department | write | `Id`, `Name`, `SyncToken`; *Active, ParentRef, SubDepartment, sparse* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_UPDATE_FULL_INVOICE` | Update Full Invoice | write | `invoice`; *minorversion, requestid* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_UPDATE_PAYMENT_METHOD` | Update Payment Method | write | `Id`, `Name`, `SyncToken`; *Active, Type* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_UPDATE_PREFERENCES` | Update Preferences | write | `Id`, `SyncToken`; *AccountingInfoPrefs, CurrencyPrefs, EmailMessagesPrefs, OtherPrefs, ProductAndServicesPrefs, ReportPrefs, SalesFormsPrefs, TaxPrefs and 4 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_UPDATE_SPARSE_INVOICE` | Update Sparse Invoice | write | `invoice` |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_UPDATE_TERM` | Update Term | write | `Id`, `Name`, `SyncToken`, `Type`; *Active, DayOfMonthDue, DiscountDayOfMonth, DiscountDays, DiscountPercent, DueDays, DueNextMonthDays* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_UPDATE_TRANSFER` | Update Transfer | write | `Amount`, `FromAccountRef`, `Id`, `SyncToken`, `ToAccountRef`; *CurrencyRef, DepartmentRef, ExchangeRate, PrivateNote, TxnDate* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_VENDOR_BALANCE_DETAIL` | Vendor Balance Detail | read | *accounting_method, aging_method, aging_period, arpaid, columns, date_macro, department_ids, duedate_macro and 8 more* |  | `com.intuit.quickbooks.accounting` |
| `QUICKBOOKS_VENDOR_BALANCE_REPORT` | Vendor Balance Report | read | *accounting_method, appaid, date_macro, department_ids, qzurl, report_date, sort_order, summarize_column_by and 1 more* |  | `com.intuit.quickbooks.accounting` |
