# Lexoffice: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner pastes a key from the vendor's settings). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

43 tools: 28 read, 15 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `LEXOFFICE_CREATE_ARTICLE` | Create Article | write | `price`, `title`, `type`, `unitName`; *articleNumber* |  |  |
| `LEXOFFICE_CREATE_CONTACT` | Create contact | write | `roles`; *addresses, company, emailAddresses, note, person, phoneNumbers, version, xRechnung* |  |  |
| `LEXOFFICE_CREATE_CREDIT_NOTE` | Create Credit Note | write | `address`, `lineItems`, `taxConditions`, `totalPrice`, `voucherDate`; *archived, finalize, introduction, language, printLayoutId, remark, title, totalDiscountAbsolute and 1 more* |  |  |
| `LEXOFFICE_CREATE_DELIVERY_NOTE` | Create Delivery Note | write | `address`, `lineItems`, `shippingConditions`, `taxConditions`, `voucherDate`; *archived, deliveryTerms, finalize, introduction, language, printLayoutId, remark, title* |  |  |
| `LEXOFFICE_CREATE_EVENT_SUBSCRIPTION` | Create Event Subscription | write | `callbackUrl`, `eventType` |  |  |
| `LEXOFFICE_CREATE_INVOICE` | Create Invoice | write | `address`, `line_items`, `shipping_conditions`, `tax_conditions`, `total_price`, `voucher_date`; *finalize, introduction, language, payment_conditions, print_layout_id, remark, title, x_rechnung* |  |  |
| `LEXOFFICE_CREATE_ORDER_CONFIRMATION` | Create Order Confirmation | write | `address`, `lineItems`, `shippingConditions`, `taxConditions`, `totalPrice`, `voucherDate`; *archived, deliveryTerms, introduction, language, printLayoutId, remark, title* |  |  |
| `LEXOFFICE_CREATE_QUOTATION` | Create Quotation | write | `address`, `expirationDate`, `lineItems`, `taxConditions`, `totalPrice`, `voucherDate`; *deliveryTerms, finalize, introduction, language, printLayoutId, remark, shippingConditions, title* |  |  |
| `LEXOFFICE_CREATE_VOUCHER` | Create Voucher | write | `taxType`, `type`, `voucherItems`; *contactId, dueDate, remark, shippingDate, totalGrossAmount, totalTaxAmount, useCollectiveContact, voucherDate and 2 more* |  |  |
| `LEXOFFICE_DELETE_ARTICLE` | Delete Article | write | `id` |  |  |
| `LEXOFFICE_DELETE_EVENT_SUBSCRIPTION` | Delete Event Subscription | write | `subscription_id` |  |  |
| `LEXOFFICE_DOWNLOAD_FILE` | Download File | read | `file_id`; *accept, file_name* |  |  |
| `LEXOFFICE_GET_ARTICLE` | Get Article | read | `id` |  |  |
| `LEXOFFICE_GET_CONTACT` | Get Contact | read | `id` |  |  |
| `LEXOFFICE_GET_CREDIT_NOTE` | Get Credit Note | read | `id` |  |  |
| `LEXOFFICE_GET_CREDIT_NOTE_DOCUMENT` | Get Credit Note Document | read | `id` |  |  |
| `LEXOFFICE_GET_DELIVERY_NOTE` | Get Delivery Note | read | `id` |  |  |
| `LEXOFFICE_GET_DUNNING` | Get Dunning | read | `dunningsid` |  |  |
| `LEXOFFICE_GET_DUNNING_DOCUMENT` | Get Dunning Document | read | `dunningsid` |  |  |
| `LEXOFFICE_GET_EVENT_SUBSCRIPTION` | Get Event Subscription | read | `subscription_id` |  |  |
| `LEXOFFICE_GET_INVOICE` | Get Invoice | read | `id` |  |  |
| `LEXOFFICE_GET_INVOICE_DOCUMENT` | Get Invoice Document | read | `id` |  |  |
| `LEXOFFICE_GET_ORDER_CONFIRMATION` | Get Order Confirmation | read | `id` |  |  |
| `LEXOFFICE_GET_ORDER_CONFIRMATION_DOCUMENT` | Render Order Confirmation Document | read | `id` |  |  |
| `LEXOFFICE_GET_PAYMENT` | Get Payment Information | read | `voucher_id` |  |  |
| `LEXOFFICE_GET_PROFILE` | Get Profile | read |  |  |  |
| `LEXOFFICE_GET_QUOTATION` | Get Quotation | read | `id` |  |  |
| `LEXOFFICE_GET_QUOTATION_DOCUMENT` | Get Quotation Document | read | `id` |  |  |
| `LEXOFFICE_GET_VOUCHER` | Get Voucher | read | `id` |  |  |
| `LEXOFFICE_LIST_ARTICLES` | List Articles | read | *article_number, gtin, page, size, sort, type* |  |  |
| `LEXOFFICE_LIST_CONTACTS` | List Contacts | read | *customer, email, name, number, page, size, sort, vendor* |  |  |
| `LEXOFFICE_LIST_COUNTRIES` | List Countries | read |  |  |  |
| `LEXOFFICE_LIST_EVENT_SUBSCRIPTIONS` | List Event Subscriptions | read |  |  |  |
| `LEXOFFICE_LIST_PAYMENT_CONDITIONS` | List Payment Conditions | read |  |  |  |
| `LEXOFFICE_LIST_POSTING_CATEGORIES` | List Posting Categories | read |  |  |  |
| `LEXOFFICE_LIST_PRINT_LAYOUTS` | List Print Layouts | read |  |  |  |
| `LEXOFFICE_LIST_RECURRING_TEMPLATES` | List Recurring Templates | read | *page, size, sort* |  |  |
| `LEXOFFICE_LIST_VOUCHERLIST` | List Voucherlist | read | *archived, contactId, createdDateFrom, createdDateTo, page, size, sort, updatedDateFrom and 6 more* |  |  |
| `LEXOFFICE_LIST_VOUCHERS` | List Vouchers | read | `voucherNumber`; *page, size* |  |  |
| `LEXOFFICE_UPDATE_ARTICLE` | Update Article | write | `id`, `price`, `type`, `unitName`, `version`; *articleNumber, gtin, note, title* |  |  |
| `LEXOFFICE_UPDATE_CONTACT` | Update lexoffice contact | write | `id`, `roles`, `version`; *addresses, company, emailAddresses, note, person, phoneNumbers* |  |  |
| `LEXOFFICE_UPDATE_VOUCHER` | Update Voucher | write | `current_version`, `voucher_id`; *clear_shipping_date, contact_id, contact_name, due_date, files, remark, shipping_date, tax_type and 8 more* |  |  |
| `LEXOFFICE_UPLOAD_VOUCHER_FILE` | Upload Voucher File | write | `file`, `id` |  |  |
