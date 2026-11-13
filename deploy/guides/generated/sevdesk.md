# Sevdesk: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner pastes a key from the vendor's settings). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

134 tools: 64 read, 70 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `SEVDESK_BOOK_INVOICE_AMOUNT` | Book Invoice Amount | write | `amount`, `checkAccount`, `date`, `invoiceId`, `type`; *checkAccountTransaction, createFeed* |  |  |
| `SEVDESK_BOOK_VOUCHER_AMOUNT` | Book Voucher Amount | write | `amount`, `booking_type`, `check_account`, `date`, `voucher_id`; *check_account_transaction, create_feed* |  |  |
| `SEVDESK_CANCEL_INVOICE` | Cancel Invoice | write | `invoice_id` |  |  |
| `SEVDESK_CHECK_CUSTOMER_NUMBER_AVAILABILITY` | Check Customer Number Availability | read | *customerNumber* |  |  |
| `SEVDESK_CREATE_ACCOUNTING_CONTACT` | Create Accounting Contact | write | `contact`; *creditorNumber, debitorNumber* |  |  |
| `SEVDESK_CREATE_CHECK_ACCOUNT_TRANSACTION` | Create Check Account Transaction | write | `amount`, `checkAccount`, `payeePayerName`, `status`, `valueDate`; *entryDate, payeePayerAcctNo, payeePayerBankCode, paymtPurpose, sourceTransaction, targetTransaction* |  |  |
| `SEVDESK_CREATE_CLEARING_ACCOUNT` | Create Clearing Account | write | `name`; *accountingNumber* |  |  |
| `SEVDESK_CREATE_COMMUNICATION_WAY` | Create Communication Way | write | `key`, `type`, `value`; *contact, main* |  |  |
| `SEVDESK_CREATE_CONTACT` | Create Contact | write | `category`, `type`; *academicTitle, address, bankDetails, birthday, customerNumber, description, email, familyname and 11 more* |  |  |
| `SEVDESK_CREATE_CONTACT2` | Create Contact (REST API) | write | `category`; *academicTitle, bankAccount, bankNumber, birthday, buyerReference, customerNumber, defaultCashbackPercent, defaultCashbackTime and 17 more* |  |  |
| `SEVDESK_CREATE_CONTACT_ADDRESS` | Create Contact Address | write | `category`, `contact`, `country`; *city, name, name2, name3, name4, street, zip* |  |  |
| `SEVDESK_CREATE_CONTACT_CUSTOM_FIELD` | Create Contact Custom Field | write | `contact`, `contactCustomFieldSetting`, `value` |  |  |
| `SEVDESK_CREATE_CONTACT_CUSTOM_FIELD_SETTING` | Create Contact Custom Field Setting | write | `name`; *description, objectName* |  |  |
| `SEVDESK_CREATE_CREDIT_NOTE_FROM_INVOICE` | Create Credit Note from Invoice | write | `invoice` |  |  |
| `SEVDESK_CREATE_FILE_IMPORT_ACCOUNT` | Create File Import Account | write | `importType`, `name`; *accountingNumber, iban* |  |  |
| `SEVDESK_CREATE_PART` | Create Part | write | `name`, `partNumber`, `stock`, `stockEnabled`, `taxRate`, `unity`; *category, internalComment, price, priceGross, priceNet, pricePurchase, status, text* |  |  |
| `SEVDESK_CREATE_PAYMENT_METHOD` | Create Payment Method | write | `name`, `type`; *isActive* |  |  |
| `SEVDESK_CREATE_PRIVATE_TRANSACTION_RULE` | Create Private Transaction Rule | write | *counterpartName, objectName, paymentPurpose* |  |  |
| `SEVDESK_CREATE_TAG` | Create Tag | write | `name`, `object` |  |  |
| `SEVDESK_DELETE_ACCOUNTING_CONTACT_BY_ID` | Delete Accounting Contact By ID | write | `accountingContactId` |  |  |
| `SEVDESK_DELETE_CHECK_ACCOUNT_BY_ID` | Delete Check Account By ID | write | `checkAccountId` |  |  |
| `SEVDESK_DELETE_CHECK_ACCOUNT_TRANSACTION_BY_ID` | Delete Check Account Transaction | write | `checkAccountTransactionId` |  |  |
| `SEVDESK_DELETE_COMMUNICATION_WAY_BY_ID` | Delete Communication Way By ID | write | `communicationWayId` |  |  |
| `SEVDESK_DELETE_CONTACT` | Delete Contact | write | `id` |  |  |
| `SEVDESK_DELETE_CONTACT_ADDRESS_BY_ID` | Delete Contact Address By ID | write | `contactAddressId` |  |  |
| `SEVDESK_DELETE_CONTACT_CUSTOM_FIELD_BY_ID` | Delete Contact Custom Field By ID | write | `contactCustomFieldId` |  |  |
| `SEVDESK_DELETE_CONTACT_CUSTOM_FIELD_SETTING_BY_ID` | Delete Contact Custom Field Setting By ID | write | `contactCustomFieldSettingId` |  |  |
| `SEVDESK_DELETE_CREDIT_NOTE_BY_ID` | Delete Credit Note By ID | write | `creditNoteId` |  |  |
| `SEVDESK_DELETE_ORDER_BY_ID` | Delete Order By ID | write | `orderId` |  |  |
| `SEVDESK_DELETE_ORDER_POS_BY_ID` | Delete Order Position By ID | write | `orderPosId` |  |  |
| `SEVDESK_DELETE_PAYMENT_METHOD` | Delete Payment Method | write | `id` |  |  |
| `SEVDESK_DELETE_PRIVATE_TRANSACTION_RULE_BY_ID` | Delete Private Transaction Rule By ID | write | `id` |  |  |
| `SEVDESK_DELETE_TAG_BY_ID` | Delete Tag By ID | write | `tagId` |  |  |
| `SEVDESK_DOWNLOAD_DOCUMENT` | Download Document | read | `document_id` |  |  |
| `SEVDESK_DOWNLOAD_INVOICE_PDF` | Download Invoice PDF | write | `invoice_id`; *download, prevent_send_by* |  |  |
| `SEVDESK_ENSHRINE_CHECK_ACCOUNT_TRANSACTION` | Enshrine Check Account Transaction | write | `checkAccountTransactionId` |  |  |
| `SEVDESK_EXPORT_CREDIT_NOTE_CSV` | Export Credit Note CSV | read | *download, sevQuery_modelName, sevQuery_objectName* |  |  |
| `SEVDESK_FETCH_DICTIONARY_ENTRIES_BY_TYPE` | Fetch Dictionary Entries By Type | read | `objectName`; *subObjectName* |  |  |
| `SEVDESK_FIND_CONTACTS_BY_CUSTOM_FIELD_VALUE` | Find Contacts By Custom Field Value | write | `value`; *customFieldName, customFieldSetting_id, customFieldSetting_objectName* |  |  |
| `SEVDESK_GENERATE_PROGRESS_DOWNLOAD_HASH` | Generate Progress Download Hash | write | `jobId` |  |  |
| `SEVDESK_GET_ACCOUNTING_CONTACTS` | Get Accounting Contacts | read | *contact_id, contact_objectName* |  |  |
| `SEVDESK_GET_ACCOUNTING_CONTACT_BY_ID` | Get Accounting Contact By ID | read | `accountingContactId` |  |  |
| `SEVDESK_GET_ACCOUNT_INVOICES` | Get SevClient Account Invoices | read |  |  |  |
| `SEVDESK_GET_CHECK_ACCOUNTS` | Get Check Accounts | read | *embed, limit, offset* |  |  |
| `SEVDESK_GET_CHECK_ACCOUNT_BALANCE_AT_DATE` | Get Check Account Balance At Date | read | `checkAccountId`, `date` |  |  |
| `SEVDESK_GET_CHECK_ACCOUNT_BY_ID` | Get Check Account By ID | read | `checkAccountId` |  |  |
| `SEVDESK_GET_CHECK_ACCOUNT_TRANSACTIONS` | Get Check Account Transactions | read | *amount, checkAccount, direction, embed, limit, offset, status, voucherDate* |  |  |
| `SEVDESK_GET_CHECK_ACCOUNT_TRANSACTION_BY_ID` | Get Check Account Transaction By ID | read | `checkAccountTransactionId` |  |  |
| `SEVDESK_GET_CHECK_ACCOUNT_TRANSACTION_LOG` | Get Check Account Transaction Log | read | `check_account_transaction_id`; *embed* |  |  |
| `SEVDESK_GET_COMMUNICATIONWAY_BY_ID` | Get CommunicationWay By ID | read | `communicationWayId` |  |  |
| `SEVDESK_GET_COMMUNICATION_WAYS` | Get Communication Ways | read | *contact_id, contact_objectName, main, type* |  |  |
| `SEVDESK_GET_COMMUNICATION_WAY_KEYS` | Get Communication Way Keys | read |  |  |  |
| `SEVDESK_GET_CONTACTS` | Get Contacts | read | *category, customerNumber, embed, field, limit, offset, order, search and 4 more* |  |  |
| `SEVDESK_GET_CONTACT_ADDRESSES` | Get Contact Addresses | read | *embed, limit, offset* |  |  |
| `SEVDESK_GET_CONTACT_ADDRESS_BY_ID` | Get Contact Address By ID | read | `contactAddressId` |  |  |
| `SEVDESK_GET_CONTACT_BY_ID` | Get Contact By ID | read | `id` |  |  |
| `SEVDESK_GET_CONTACT_CUSTOM_FIELDS` | Get Contact Custom Fields | write | *countAll, embed, limit, offset* |  |  |
| `SEVDESK_GET_CONTACT_CUSTOM_FIELD_BY_ID` | Get Contact Custom Field By ID | write | `contactCustomFieldId` |  |  |
| `SEVDESK_GET_CONTACT_CUSTOM_FIELD_SETTINGS` | Get Contact Custom Field Settings | write | *countAll, embed, limit, offset* |  |  |
| `SEVDESK_GET_CONTACT_CUSTOM_FIELD_SETTING_BY_ID` | Get Contact Custom Field Setting By ID | write | `contactCustomFieldSettingId` |  |  |
| `SEVDESK_GET_CONTACT_CUSTOM_FIELD_SETTING_REFERENCE_COUNT` | Get Contact Custom Field Setting Reference Count | write | `contactCustomFieldSettingId` |  |  |
| `SEVDESK_GET_CONTACT_TABS_ITEM_COUNT` | Get Contact Tabs Item Count | read | `contact_id` |  |  |
| `SEVDESK_GET_CREDIT_NOTES` | Get Credit Notes | read | *contact_familyname, contact_id, contact_surname, creditNoteNumber, customer_customerNumber, customer_id, customer_name, embed and 5 more* |  |  |
| `SEVDESK_GET_CREDIT_NOTE_POSITIONS` | Get Credit Note Positions | read | *creditNote_id, creditNote_objectName* |  |  |
| `SEVDESK_GET_EXPORT_INVOICECSV` | Get Export Invoice CSV | read | *download, sevQuery_modelName, sevQuery_objectName* |  |  |
| `SEVDESK_GET_EXPORT_TRANSACTIONS_CSV` | Get Export Transactions CSV | read | *download, sevQuery_modelName, sevQuery_objectName* |  |  |
| `SEVDESK_GET_INVOICES` | Get Invoices | read | *contact_id, customerInternalNote, deliveryDate, dueDate, embed, invoiceDate, limit, number and 4 more* |  |  |
| `SEVDESK_GET_INVOICE_BY_ID` | Get Invoice By ID | read | `invoiceId` |  |  |
| `SEVDESK_GET_INVOICE_IS_PARTIALLY_PAID` | Get Invoice Is Partially Paid | read | `invoiceId` |  |  |
| `SEVDESK_GET_INVOICE_POSITIONS` | Get Invoice Positions | read | `invoiceId`; *embed, limit, offset* |  |  |
| `SEVDESK_GET_INVOICE_POSS` | Get Invoice Positions | read | *id, invoice_id, invoice_objectName, part_id, part_objectName* |  |  |
| `SEVDESK_GET_LETTERPAPERS_WITH_THUMB` | Get Letterpapers With Thumbnails | read |  |  |  |
| `SEVDESK_GET_NEXT_CUSTOMER_NUMBER` | Get Next Customer Number | read |  |  |  |
| `SEVDESK_GET_ORDERS` | Get Orders | read | *embed, filter_customer, filter_search, filter_status, limit, offset, order* |  |  |
| `SEVDESK_GET_ORDER_BY_ID` | Get Order By ID | read | `orderId` |  |  |
| `SEVDESK_GET_ORDER_DISCOUNTS` | Get Order Discounts | read | `orderId`; *embed, limit, offset* |  |  |
| `SEVDESK_GET_ORDER_POSITIONS` | Get Order Positions | read | *order_id, order_objectName* |  |  |
| `SEVDESK_GET_ORDER_POS_BY_ID` | Get Order Position By ID | read | `orderPosId` |  |  |
| `SEVDESK_GET_PARTS` | Get Parts | read | *archived, countAll, embed, limit, offset, properties* |  |  |
| `SEVDESK_GET_PART_BY_ID` | Get Part By ID | read | `part_id` |  |  |
| `SEVDESK_GET_PART_CONTACT_PRICES` | Get Part Contact Prices | read | *contact, embed, fields, includeDeleted, limit, offset, part, updatedAfter* |  |  |
| `SEVDESK_GET_PART_STOCK` | Get Part Stock | read | `partId` |  |  |
| `SEVDESK_GET_PART_UNITIES` | Get Part Unities | read | *embed, includeDeleted, limit, offset* |  |  |
| `SEVDESK_GET_PAYMENT_METHODS` | Get Payment Methods | read | *embed, limit, offset* |  |  |
| `SEVDESK_GET_PLACES` | Get Places | read | *countAll, depth, embed, limit, offset* |  |  |
| `SEVDESK_GET_PRIVATE_TRANSACTION_RULES` | Get Private Transaction Rules | read | *embed, limit, offset* |  |  |
| `SEVDESK_GET_PROGRESS` | Get Export Progress | read | `hash` |  |  |
| `SEVDESK_GET_RECEIPT_GUIDANCE_FOR_ACCOUNT_NUMBER` | Get Receipt Guidance For Account Number | read | `accountNumber` |  |  |
| `SEVDESK_GET_RECEIPT_GUIDANCE_FOR_ALL_ACCOUNTS` | Get Receipt Guidance For All Accounts | read |  |  |  |
| `SEVDESK_GET_RECEIPT_GUIDANCE_FOR_EXPENSE` | Get Receipt Guidance For Expense | read |  |  |  |
| `SEVDESK_GET_RECEIPT_GUIDANCE_FOR_REVENUE` | Get Receipt Guidance For Revenue | read |  |  |  |
| `SEVDESK_GET_RECEIPT_GUIDANCE_FOR_TAX_RULE` | Get Receipt Guidance For Tax Rule | read | `taxRule` |  |  |
| `SEVDESK_GET_REPORT_INVOICELIST` | Get Report Invoice List | read | `view`; *download, sevQuery_limit, sevQuery_modelName, sevQuery_objectName* |  |  |
| `SEVDESK_GET_REPORT_ORDERLIST` | Get Report Orderlist | read | `view`; *download, sevQuery_limit, sevQuery_modelName, sevQuery_objectName* |  |  |
| `SEVDESK_GET_SEV_CLIENT_CONFIG` | Get SevClient Config | read | *embed* |  |  |
| `SEVDESK_GET_TAGS` | Get Tags | read | *id, name* |  |  |
| `SEVDESK_GET_TAG_BY_ID` | Get Tag By ID | read | `tagId` |  |  |
| `SEVDESK_GET_TAG_RELATIONS` | Get Tag Relations | read | *embed, limit, offset* |  |  |
| `SEVDESK_GET_TAX_RULES` | Get Tax Rules | read | *embed, includeDeleted, limit, offset* |  |  |
| `SEVDESK_GET_TOOLS_BOOKKEEPING_SYSTEM_VERSION` | Get Bookkeeping System Version | read |  |  |  |
| `SEVDESK_GET_VOUCHERS` | Get Vouchers | read | *accounting_type_id, accounting_type_object_name, commentLike, contact_id, contact_object_name, contact_or_object_id, contact_or_object_name, cost_centre_id and 36 more* |  |  |
| `SEVDESK_GET_VOUCHER_BY_ID` | Get Voucher By ID | read | `voucher_id` |  |  |
| `SEVDESK_GET_VOUCHER_DOCUMENT_IMAGE` | Get Voucher Document Image | read | `voucher_id` |  |  |
| `SEVDESK_GET_VOUCHER_POSITIONS` | Get Voucher Positions | read | *embed, voucher_id* |  |  |
| `SEVDESK_RENDER_INVOICE` | Render Invoice | read | `invoiceId`; *forceReload* |  |  |
| `SEVDESK_RESET_CREDIT_NOTE_TO_OPEN` | Reset Credit Note To Open | write | `creditNoteId` |  |  |
| `SEVDESK_RESET_INVOICE_TO_OPEN` | Reset Invoice To Open | write | `invoiceId` |  |  |
| `SEVDESK_RESET_VOUCHER_STATE` | Reset Voucher State | write | `state`, `voucher_id` |  |  |
| `SEVDESK_SAVE_INVOICE` | Save Invoice | write | `invoice`, `invoicePosSave`; *discountDelete, discountSave, filename, invoicePosDelete* |  |  |
| `SEVDESK_SAVE_ORDER` | Save Order | write | `order`; *orderPosSave* |  |  |
| `SEVDESK_SAVE_VOUCHER` | Save Voucher | write | `voucher`, `voucherPosSave`; *filename, voucherPosDelete* |  |  |
| `SEVDESK_SEND_CREDIT_NOTE_BY` | Send Credit Note By | write | `creditNoteId`, `sendDraft`, `sendType` |  |  |
| `SEVDESK_SEND_CREDIT_NOTE_BY_WITH_RENDER` | Send Credit Note By Printing | write | `creditNoteId`, `sendType` |  |  |
| `SEVDESK_SEND_INVOICE_BY` | Send Invoice By | write | `invoiceId`, `sendDraft`, `sendType` |  |  |
| `SEVDESK_SEND_INVOICE_VIA_EMAIL` | Send Invoice Via Email | write | `invoice_id`, `subject`, `text`, `to_email`; *additional_attachments, bcc_email, cc_email, copy, send_xml* |  |  |
| `SEVDESK_SEND_ORDER_BY` | Mark Order as Sent | write | `orderId`, `sendDraft`, `sendType` |  |  |
| `SEVDESK_UPDATE_ACCOUNTING_CONTACT_BY_ID` | Update Accounting Contact | write | `accountingContactId`; *contact, creditorNumber, debitorNumber* |  |  |
| `SEVDESK_UPDATE_CHECK_ACCOUNT` | Update Check Account | write | `id`; *accountHolder, active, bankAccountNumber, bankCode, bankName, bic, customer, default and 2 more* |  |  |
| `SEVDESK_UPDATE_CHECK_ACCOUNT_TRANSACTION_BY_ID` | Update Check Account Transaction By ID | write | `checkAccountTransactionId`; *amount, checkAccount, entryDate, payeePayerName, paymtPurpose, sourceTransaction, status, targetTransaction and 1 more* |  |  |
| `SEVDESK_UPDATE_COMMUNICATION_WAY_BY_ID` | Update Communication Way By ID | write | `communicationWayId`; *contact, key, main, type, value* |  |  |
| `SEVDESK_UPDATE_CONTACT` | Update Contact | write | `id`; *academicTitle, address, bankDetails, category, creditNumber, customerNumber, description, familyname and 8 more* |  |  |
| `SEVDESK_UPDATE_CONTACT_ADDRESS_BY_ID` | Update Contact Address By ID | write | `contactAddressId`; *category, city, contact, country, name, name2, name3, name4 and 2 more* |  |  |
| `SEVDESK_UPDATE_CONTACT_CUSTOM_FIELD_BY_ID` | Update Contact Custom Field | write | `contact`, `contactCustomFieldId`, `contactCustomFieldSetting`, `value`; *objectName* |  |  |
| `SEVDESK_UPDATE_CONTACT_CUSTOM_FIELD_SETTING_BY_ID` | Update Contact Custom Field Setting | write | `contactCustomFieldSettingId`; *description, name, objectName* |  |  |
| `SEVDESK_UPDATE_INVOICE_PARAMETER` | Update Invoice Parameter | write | `invoiceId`, `key`, `value` |  |  |
| `SEVDESK_UPDATE_INVOICE_POSITION` | Update Invoice Position | write | `invoice_position_id`; *discount, name, part, position_number, price, price_gross, price_tax, quantity and 3 more* |  |  |
| `SEVDESK_UPDATE_ORDER_BY_ID` | Update Order By ID | write | `orderId`; *address, addressCountry, contact, contactPerson, currency, customerInternalNote, deliveryTerms, footText and 17 more* |  |  |
| `SEVDESK_UPDATE_ORDER_POSITION` | Update Order Position | write | `orderPosId`; *discount, name, optional, order, part, positionNumber, price, priceGross and 5 more* |  |  |
| `SEVDESK_UPDATE_ORDER_TEMPLATE` | Update Order Template | write | `key`, `orderId`, `value` |  |  |
| `SEVDESK_UPDATE_PART_BY_ID` | Update Part By ID | write | `partId`; *category, internalComment, name, partNumber, price, priceGross, priceNet, pricePurchase and 6 more* |  |  |
| `SEVDESK_UPDATE_SEV_CLIENT_EXPORT_CONFIG` | Update SevClient Export Config | write | `accountant_client_number`, `accountant_number`, `accounting_year_begin`, `sev_client_id` |  |  |
| `SEVDESK_UPDATE_TAG_BY_ID` | Update Tag By ID | write | `name`, `tagId` |  |  |
| `SEVDESK_UPDATE_VOUCHER` | Update Voucher | write | `voucher_id`; *cost_centre, credit_debit, currency, delivery_date, delivery_date_until, description, document, pay_date and 10 more* |  |  |
| `SEVDESK_UPLOAD_VOUCHER_FILE` | Upload Voucher File | write | `file` |  |  |
