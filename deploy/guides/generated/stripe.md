# Stripe: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

428 tools: 203 read, 225 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `STRIPE_ACCEPT_QUOTE` | Accept quote | write | `quote_id` |  | `read_write` |
| `STRIPE_ACTIVATE_BILLING_ALERT` | Activate billing alert | write | `alert_id` |  |  |
| `STRIPE_ADD_INVOICE_LINES` | Add lines to invoice | write | `id`, `lines`; *invoice_metadata* |  | `read_write` |
| `STRIPE_ADVANCE_TEST_CLOCK` | Advance test clock | write | `clock_id`, `frozen_time` |  | `read_write` |
| `STRIPE_APPLY_CUSTOMER_BALANCE` | Apply customer balance to payment intent | write | `intent`; *amount, currency* |  | `read_write` |
| `STRIPE_ARCHIVE_BILLING_ALERT` | Archive billing alert | write | `alert_id` |  | `read_write` |
| `STRIPE_ATTACH_CUSTOMER_SOURCE` | Attach source to customer | write | `customer_id`, `source`; *metadata* |  | `read_only` |
| `STRIPE_ATTACH_INVOICE_PAYMENT` | Attach payment to invoice | write | `invoice`; *payment_intent, payment_record* |  | `read_write` |
| `STRIPE_ATTACH_PAYMENT_METHOD` | Attach payment method to customer | write | `payment_method_id`; *customer, customer_account* |  | `read_write` |
| `STRIPE_ATTACH_PRODUCT_FEATURE` | Attach feature to product | write | `entitlement_feature`, `product_id` |  | `read_only` |
| `STRIPE_CANCEL_PAYMENT_INTENT` | Cancel payment intent | write | `payment_intent_id`; *cancellation_reason* |  |  |
| `STRIPE_CANCEL_QUOTE` | Cancel Quote | write | `quote_id` |  |  |
| `STRIPE_CANCEL_SETUP_INTENT` | Cancel setup intent | write | `setup_intent_id`; *cancellation_reason* |  |  |
| `STRIPE_CANCEL_SUBSCRIPTION` | Cancel subscription | write | `subscription_id`; *cancellation_details, invoice_now, prorate* |  |  |
| `STRIPE_CANCEL_SUBSCRIPTION_SCHEDULE` | Cancel subscription schedule | write | `schedule`; *invoice_now, prorate* |  |  |
| `STRIPE_CANCEL_TERMINAL_READER_ACTION` | Cancel Terminal Reader Action | write | `reader_id` |  |  |
| `STRIPE_CAPTURE_CHARGE` | Capture Charge | write | `charge`; *amount, application_fee_amount, receipt_email, statement_descriptor, statement_descriptor_suffix, transfer_data_amount, transfer_group* |  |  |
| `STRIPE_CAPTURE_PAYMENT_INTENT` | Capture payment intent | write | `intent`; *amount_to_capture, application_fee_amount, final_capture, metadata, statement_descriptor, statement_descriptor_suffix* |  |  |
| `STRIPE_CLOSE_DISPUTE` | Close Dispute | write | `dispute` |  |  |
| `STRIPE_COLLECT_TERMINAL_READER_INPUTS` | Set Terminal Reader Display | write | `inputs`, `reader_id`; *metadata* |  | `read_write` |
| `STRIPE_COLLECT_TERMINAL_READER_PAYMENT_METHOD` | Collect payment method on Terminal Reader | write | `payment_intent`, `reader`; *collect_config* |  | `read_write` |
| `STRIPE_CONFIRM_PAYMENT_INTENT` | Confirm payment intent | write | `payment_intent_id`; *capture_method, confirmation_token, error_on_requires_action, mandate, off_session, payment_method, payment_method_data, payment_method_options and 4 more* |  |  |
| `STRIPE_CONFIRM_SETUP_INTENT` | Confirm setup intent | write | `setup_intent_id`; *confirmation_token, payment_method, return_url, use_stripe_sdk* |  | `read_write` |
| `STRIPE_CONFIRM_TERMINAL_READER_PAYMENT_INTENT` | Confirm PaymentIntent on Terminal Reader | write | `payment_intent`, `reader_id`; *confirm_config* |  | `read_write` |
| `STRIPE_CREATE_ACCOUNTS_CAPABILITIES_CARD_PAYMENTS` | Create Card Payments Capability | write | `id`; *requested* |  | `read_write` |
| `STRIPE_CREATE_ACCOUNT_SESSION` | Create Account Session | write | `account`, `components`; *expand* |  | `connected_account_read` |
| `STRIPE_CREATE_APPLE_PAY_DOMAIN` | Create Apple Pay Domain | write | `domain_name`; *expand* |  | `read_write` |
| `STRIPE_CREATE_APPS_SECRET` | Create apps secret | write | `name`, `payload`, `scope_type`; *expires_at, scope_user* |  | `read_write` |
| `STRIPE_CREATE_BANK_ACCOUNT_TOKEN` | Create Bank Account Token | write | `bank_account`; *customer* |  |  |
| `STRIPE_CREATE_BILLING_ALERT` | Create billing alert | write | `alert_type`, `title`, `usage_threshold` |  | `read_write` |
| `STRIPE_CREATE_BILLING_CREDIT_GRANT` | Create billing credit grant | write | `amount`, `applicability_config`; *category, customer, customer_account, effective_at, expires_at, metadata, name, priority* |  | `read_write` |
| `STRIPE_CREATE_BILLING_METER` | Create billing meter | write | `default_aggregation`, `display_name`, `event_name`; *customer_mapping, event_time_window, value_settings* |  |  |
| `STRIPE_CREATE_BILLING_METER_EVENT` | Create Billing Meter Event V2 | write | `event_name`, `payload`; *identifier, timestamp* |  | `read_write` |
| `STRIPE_CREATE_BILLING_METER_EVENT_ADJUSTMENT` | Create Billing Meter Event Adjustment | write | `event_name`, `type`; *cancel* |  | `read_write` |
| `STRIPE_CREATE_BILLING_METER_EVENT_SESSION` | Create Billing Meter Event Session | write |  |  |  |
| `STRIPE_CREATE_BILLING_PORTAL_CONFIGURATION` | Create Billing Portal Configuration | write | `features`; *business_profile, default_return_url, login_page, metadata, name* |  |  |
| `STRIPE_CREATE_BILLING_PORTAL_SESSION` | Create Billing Portal Session | write | *configuration, customer, customer_account, flow_data, locale, on_behalf_of, return_url* |  | `read_write` |
| `STRIPE_CREATE_CHARGE_REFUND` | Create Charge Refund | write | `charge`; *amount, instructions_email, metadata, payment_intent, reason, refund_application_fee, reverse_transfer* |  |  |
| `STRIPE_CREATE_CHECKOUT_SESSION` | Create Checkout Session | write | `line_items`, `mode`, `success_url`; *allow_promotion_codes, cancel_url, client_reference_id, customer, metadata, payment_method_types* |  |  |
| `STRIPE_CREATE_COUPON` | Create Coupon | write | *amount_off, applies_to, currency, currency_options, duration, duration_in_months, id, max_redemptions and 4 more* |  | `read_write` |
| `STRIPE_CREATE_CREDIT_NOTE` | Create Credit Note | write | `invoice`; *amount, credit_amount, effective_at, email_type, lines, memo, metadata, out_of_band_amount and 4 more* |  | `read_write` |
| `STRIPE_CREATE_CUSTOMER` | Create Customer | write | *address, description, email, name, phone* |  |  |
| `STRIPE_CREATE_CUSTOMER_BALANCE_TRANSACTION` | Create Customer Balance Transaction | write | `amount`, `currency`, `customer`; *description, metadata* |  | `read_write` |
| `STRIPE_CREATE_CUSTOMER_BANK_ACCOUNT` | Create customer bank account | write | `customer`, `source`; *metadata* |  | `read_only` |
| `STRIPE_CREATE_CUSTOMER_CARD` | Create customer card | write | `customer`; *alipay_account, bank_account, card, expand, metadata, source* |  | `read_write` |
| `STRIPE_CREATE_CUSTOMER_SESSION` | Create Customer Session | write | `components`; *customer, customer_account* |  |  |
| `STRIPE_CREATE_CUSTOMER_SOURCE` | Create card or payment source | write | `customer`, `source`; *metadata* |  | `read_write` |
| `STRIPE_CREATE_CUSTOMER_SUBSCRIPTION` | Create customer subscription | write | `customer`; *add_invoice_items, application_fee_percent, automatic_tax, billing_cycle_anchor, billing_thresholds, cancel_at, cancel_at_period_end, collection_method and 17 more* |  | `read_write` |
| `STRIPE_CREATE_CUSTOMER_TAX_ID` | Create Customer Tax ID | write | `customer`, `type`, `value` |  | `read_write` |
| `STRIPE_CREATE_CVC_UPDATE_TOKEN` | Create CVC Update Token | write | `cvc` |  |  |
| `STRIPE_CREATE_ENTITLEMENT_FEATURE` | Create Feature Entitlement | write | `lookup_key`, `name`; *metadata* |  | `read_write` |
| `STRIPE_CREATE_EPHEMERAL_KEY` | Create Ephemeral Key | write | `stripe_version`; *customer, expand, issuing_card, nonce, verification_session* |  | `read_write` |
| `STRIPE_CREATE_FILE` | Create File | write | `file`, `purpose`; *file_link_data* |  | `read_write` |
| `STRIPE_CREATE_FILE_LINK` | Create File Link | write | `file`; *expires_at, metadata* |  | `read_write` |
| `STRIPE_CREATE_FINANCIAL_CONNECTIONS_SESSION` | Create Financial Connections Session | write | `account_holder`, `permissions`; *filters, prefetch, return_url* |  | `read_write` |
| `STRIPE_CREATE_FX_QUOTE` | Create FX Quote | write | `from_currencies`, `to_currency`; *lock_duration, usage_payment_destination, usage_payment_on_behalf_of, usage_transfer_destination, usage_type* |  |  |
| `STRIPE_CREATE_INVOICE` | Create an invoice | write | `customer`; *account_tax_ids, application_fee_amount, auto_advance, automatic_tax, automatically_finalizes_at, collection_method, currency, custom_fields and 21 more* |  |  |
| `STRIPE_CREATE_INVOICE_ITEM` | Create invoice item | write | *amount, currency, customer, customer_account, description, discountable, discounts, invoice and 10 more* |  | `read_only` |
| `STRIPE_CREATE_PAYMENT_INTENT` | Create payment intent | write | `amount`, `currency`; *application_fee_amount, automatic_payment_methods, confirm, customer, description, metadata, off_session, payment_method and 8 more* |  |  |
| `STRIPE_CREATE_PAYMENT_LINK` | Create Payment Link | write | `line_items`; *after_completion, allow_promotion_codes, application_fee_amount, application_fee_percent, automatic_tax, billing_address_collection, currency, customer_creation and 12 more* |  | `read_write` |
| `STRIPE_CREATE_PAYMENT_METHOD` | Create payment method | write | `type`; *acss_debit, allow_redisplay, au_becs_debit, bacs_debit, billing_details, boleto, card, klarna and 4 more* |  | `read_write` |
| `STRIPE_CREATE_PAYMENT_METHOD_CONFIGURATION` | Create payment method configuration | write | *acss_debit, affirm, afterpay_clearpay, alipay, alma, amazon_pay, apple_pay, apple_pay_later and 51 more* |  |  |
| `STRIPE_CREATE_PAYMENT_METHOD_DOMAIN` | Create payment method domain | write | `domain_name`; *enabled* |  | `read_write` |
| `STRIPE_CREATE_PII_TOKEN` | Create PII Token | write | `pii` |  |  |
| `STRIPE_CREATE_PLAN` | Create a plan | write | `currency`, `interval`; *active, amount, amount_decimal, billing_scheme, id, interval_count, metadata, meter and 7 more* |  |  |
| `STRIPE_CREATE_PREVIEW_INVOICE` | Create Preview Invoice | read | *automatic_tax, currency, customer, customer_account, customer_details, discounts, invoice_items, issuer and 5 more* |  | `read_write` |
| `STRIPE_CREATE_PRICE` | Create a price | write | `currency`; *active, billing_scheme, currency_options, custom_unit_amount, lookup_key, metadata, nickname, product and 9 more* |  |  |
| `STRIPE_CREATE_PRODUCT` | Create product | write | `name`; *active, custom_fields, default_price_data, description, id, images, marketing_features, metadata and 7 more* |  |  |
| `STRIPE_CREATE_PROMOTION_CODE` | Create Promotion Code | write | `coupon`; *active, code, customer, expires_at, max_redemptions, metadata, restrictions* |  | `read_write` |
| `STRIPE_CREATE_QUOTE` | Create a quote | write | *application_fee_amount, application_fee_percent, automatic_tax, collection_method, customer, customer_account, default_tax_rates, description and 12 more* |  | `read_write` |
| `STRIPE_CREATE_REFUND` | Create Refund | write | *amount, charge, metadata, payment_intent, reason, refund_application_fee, reverse_transfer* |  |  |
| `STRIPE_CREATE_REPORT_RUN` | Create a Report Run | write | `report_type`; *parameters* |  | `read_write` |
| `STRIPE_CREATE_REVERSAL_TAX_TRANSACTION` | Create Reversal Tax Transaction | write | `mode`, `original_transaction`, `reference`; *flat_amount, line_items, metadata, shipping_cost* |  | `read_write` |
| `STRIPE_CREATE_SETUP_INTENT` | Create Setup Intent | write | *attach_to_self, automatic_payment_methods, confirm, confirmation_token, customer, customer_account, description, excluded_payment_method_types and 13 more* |  |  |
| `STRIPE_CREATE_SHIPPING_RATE` | Create a shipping rate | write | `display_name`, `fixed_amount`; *delivery_estimate, metadata, tax_behavior, tax_code, type* |  | `read_write` |
| `STRIPE_CREATE_SOURCE` | Create Source | write | *amount, currency, flow, mandate, metadata, owner, receiver, redirect and 5 more* |  | `read_write` |
| `STRIPE_CREATE_SUBSCRIPTION` | Create subscription | write | `customer`, `items`; *add_invoice_items, application_fee_percent, automatic_tax, backdate_start_date, billing_cycle_anchor, billing_cycle_anchor_config, billing_thresholds, cancel_at and 23 more* |  |  |
| `STRIPE_CREATE_SUBSCRIPTION_ITEM` | Create subscription item | write | `subscription`; *billing_thresholds, discounts, metadata, payment_behavior, price, price_data, proration_behavior, proration_date and 2 more* |  | `read_write` |
| `STRIPE_CREATE_SUBSCRIPTION_SCHEDULE` | Create subscription schedule | write | *billing_mode, customer, customer_account, default_settings, end_behavior, from_subscription, metadata, phases and 1 more* |  | `read_write` |
| `STRIPE_CREATE_TAX_CALCULATION` | Create Tax Calculation | write | `currency`, `line_items`; *customer, customer_details, ship_from_details, shipping_cost, tax_date* |  | `read_write` |
| `STRIPE_CREATE_TAX_ID` | Create Tax ID | write | `type`, `value`; *owner* |  |  |
| `STRIPE_CREATE_TAX_RATE` | Create tax rate | write | `display_name`, `inclusive`, `percentage`; *active, country, description, jurisdiction, metadata, state, tax_type* |  | `read_write` |
| `STRIPE_CREATE_TAX_REGISTRATION` | Create Tax Registration | write | `active_from`, `country`, `country_options`; *expires_at* |  |  |
| `STRIPE_CREATE_TAX_TRANSACTION_FROM_CALCULATION` | Create Tax Transaction from Calculation | write | `calculation`, `reference`; *metadata, posted_at* |  | `read_write` |
| `STRIPE_CREATE_TERMINAL_CONFIGURATION` | Create Terminal Configuration | write | *bbpos_wisepad3, bbpos_wisepos_e, name, offline, reader_security, reboot_window, stripe_s700, tipping and 2 more* |  | `read_write` |
| `STRIPE_CREATE_TERMINAL_CONNECTION_TOKEN` | Create Terminal Connection Token | write | *location* |  | `read_write` |
| `STRIPE_CREATE_TERMINAL_LOCATION` | Create Terminal Location | write | `address`; *configuration_overrides, display_name, metadata, phone* |  |  |
| `STRIPE_CREATE_TERMINAL_ONBOARDING_LINK` | Create Terminal Onboarding Link | write | `link_options`, `link_type`; *expand, on_behalf_of* |  | `read_write` |
| `STRIPE_CREATE_TERMINAL_READER` | Create Terminal Reader | write | `location`, `registration_code`; *label, metadata* |  | `read_only` |
| `STRIPE_CREATE_TEST_CLOCK` | Create Test Clock | write | `frozen_time`; *name* |  |  |
| `STRIPE_CREATE_TEST_CONFIRMATION_TOKEN` | Create Test Confirmation Token (Test Mode Only) | write | *payment_method, payment_method_data, payment_method_options, return_url, setup_future_usage, shipping* |  |  |
| `STRIPE_CREATE_TOKEN` | Create CVC update token | write | `cvc_update` |  | `read_write` |
| `STRIPE_CREATE_TRANSFER` | Create Transfer | write | `currency`, `destination`; *amount, description, expand, metadata, source_transaction, source_type, transfer_group* |  | `read_write` |
| `STRIPE_CREATE_USAGE_RECORD` | Create usage record | write | `quantity`, `subscription_item`, `timestamp`; *action* |  | `read_write` |
| `STRIPE_DEACTIVATE_BILLING_ALERT` | Deactivate billing alert | write | `id` |  | `read_write` |
| `STRIPE_DEACTIVATE_BILLING_METER` | Deactivate Billing Meter | write | `id` |  | `read_write` |
| `STRIPE_DELETEV1_SUBSCRIPTIONS_LIST_SUBSCRIPTIONS` | Cancel subscription | write | `subscription_id`; *cancellation_details, invoice_now, prorate* |  | `read_write` |
| `STRIPE_DELETE_APPLE_PAY_DOMAIN` | Delete Apple Pay Domain | write | `domain` |  | `read_write` |
| `STRIPE_DELETE_APPS_SECRET` | Delete apps secret | write | `name`, `scope_type`; *scope_user* |  | `read_write` |
| `STRIPE_DELETE_COUPON` | Delete coupon | write | `coupon` |  | `read_write` |
| `STRIPE_DELETE_CUSTOMERS_CUSTOMER` | Delete customer | write | `customer` |  | `read_write` |
| `STRIPE_DELETE_CUSTOMER_BANK_ACCOUNT` | Delete customer bank account | write | `customer`, `id`; *expand* |  | `read_write` |
| `STRIPE_DELETE_CUSTOMER_DISCOUNT` | Delete customer discount | write | `customer` |  |  |
| `STRIPE_DELETE_CUSTOMER_SUBSCRIPTION` | Cancel customer subscription | write | `customer`, `subscription_exposed_id`; *invoice_now, prorate* |  |  |
| `STRIPE_DELETE_CUSTOMER_SUBSCRIPTION_DISCOUNT` | Delete customer subscription discount | write | `customer`, `subscription_exposed_id` |  |  |
| `STRIPE_DELETE_CUSTOMER_TAX_ID` | Delete customer tax ID | write | `customer`, `id` |  |  |
| `STRIPE_DELETE_EPHEMERAL_KEY` | Delete ephemeral key | write | `key` |  |  |
| `STRIPE_DELETE_INVOICE` | Delete draft invoice | write | `invoice` |  | `read_write` |
| `STRIPE_DELETE_INVOICEITEM` | Delete invoice item | write | `invoiceitem` |  | `read_write` |
| `STRIPE_DELETE_PLAN` | Delete plan | write | `plan` |  | `read_write` |
| `STRIPE_DELETE_PRODUCT` | Delete product | write | `product_id` |  | `read_write` |
| `STRIPE_DELETE_PRODUCT_FEATURE` | Delete product feature | write | `id`, `product` |  | `read_write` |
| `STRIPE_DELETE_SUBSCRIPTION_DISCOUNT` | Delete subscription discount | write | `subscription_id` |  |  |
| `STRIPE_DELETE_SUBSCRIPTION_ITEM` | Delete subscription item | write | `item`; *clear_usage, proration_behavior, proration_date* |  | `read_write` |
| `STRIPE_DELETE_TAX_ID` | Delete tax ID | write | `id` |  | `read_write` |
| `STRIPE_DELETE_TERMINAL_CONFIGURATION` | Delete Terminal Configuration | write | `configuration` |  |  |
| `STRIPE_DELETE_TERMINAL_LOCATION` | Delete Terminal Location | write | `location` |  | `read_write` |
| `STRIPE_DELETE_TERMINAL_READER` | Delete Terminal Reader | write | `reader` |  | `read_write` |
| `STRIPE_DELETE_TEST_CLOCK` | Delete test clock | write | `test_clock` |  | `read_write` |
| `STRIPE_DELETE_V1_INVOICES_VOID_AN_INVOICE` | Void an invoice | write | `invoice_id` |  | `read_write` |
| `STRIPE_DETACH_PAYMENT_FROM_INVOICE` | Detach payment from invoice | write | `id`; *payment_intent, payment_record* |  | `read_write` |
| `STRIPE_DETACH_PAYMENT_METHOD` | Detach payment method | write | `payment_method_id` |  | `read_write` |
| `STRIPE_DISABLE_EVENT_DESTINATION` | Disable Event Destination | write | `event_destination_id` |  | `EventDestinations:Write` |
| `STRIPE_DOWNLOAD_FILE_CONTENTS` | Download File Contents | read | `file_id` |  | `read_only` |
| `STRIPE_DOWNLOAD_QUOTE_PDF` | Download Quote PDF | read | `quote` |  | `read_only` |
| `STRIPE_EXPIRE_BILLING_CREDIT_GRANT` | Expire billing credit grant | write | `id` |  |  |
| `STRIPE_EXPIRE_CHECKOUT_SESSION` | Expire Checkout Session | write | `session_id` |  | `read_write` |
| `STRIPE_FINALIZE_INVOICE` | Finalize an invoice | write | `invoice_id`; *auto_advance, expand* |  | `read_write` |
| `STRIPE_FINALIZE_QUOTE` | Finalize quote | write | `quote`; *expand, expires_at* |  | `read_write` |
| `STRIPE_FIND_APPS_SECRET` | Find a secret by name and scope | read | `name`, `scope` |  | `read_write` |
| `STRIPE_FIND_TAX_ASSOCIATION` | Find tax association | read | `payment_intent` |  |  |
| `STRIPE_FUND_TEST_CASH_BALANCE` | Fund Test Mode Cash Balance | write | `amount`, `currency`, `customer_id`; *reference* |  | `read_write` |
| `STRIPE_GET_ACCOUNT` | Retrieve Account | read | *expand* |  |  |
| `STRIPE_GET_ACCOUNTS_ACCOUNT` | Retrieve Account | read | `account`; *expand* |  | `read_only` |
| `STRIPE_GET_ACCOUNT_CAPABILITIES` | List Account Capabilities | read | `account` |  | `read_only` |
| `STRIPE_GET_ACCOUNT_CAPABILITY` | Retrieve Account Capability | read | `account_id`, `capability_id`; *expand* |  | `read_only` |
| `STRIPE_GET_ACCOUNT_PERSON` | Retrieve Person from Account | read | `account`, `person` |  | `read_only` |
| `STRIPE_GET_ACCOUNT_PERSONS` | List all persons for an account | read | `account`; *ending_before, limit, relationship_authorizer, relationship_director, relationship_executive, relationship_legal_guardian, relationship_owner, relationship_representative and 1 more* |  | `read_only` |
| `STRIPE_GET_ACCOUNT_TOKEN` | Retrieve Account Token | read | `token_id`; *expand* |  | `api_key` |
| `STRIPE_GET_ACTIVE_ENTITLEMENT` | Retrieve active entitlement | read | `id`; *expand* |  | `read_only` |
| `STRIPE_GET_APPLE_PAY_DOMAIN` | Retrieve Apple Pay Domain | read | `domain` |  | `read_only` |
| `STRIPE_GET_BALANCE_SETTINGS` | Get balance settings overview | read | *stripe_account* |  | `read_write`, `read_only` |
| `STRIPE_GET_BALANCE_TRANSACTION` | Retrieve Balance Transaction | read | `id`; *expand* |  | `read_only` |
| `STRIPE_GET_BILLING_ALERT` | Retrieve Billing Alert | read | `id` |  | `read_only` |
| `STRIPE_GET_BILLING_CREDIT_BALANCE_SUMMARY` | Get Billing Credit Balance Summary | read | `filter`; *customer, customer_account* |  | `read_only` |
| `STRIPE_GET_BILLING_CREDIT_BALANCE_TRANSACTIONS_ID` | Retrieve credit balance transaction | read | `id` |  |  |
| `STRIPE_GET_BILLING_CREDIT_GRANTS` | List billing credit grants | read | *customer, customer_account, ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_GET_BILLING_CREDIT_GRANTS_ID` | Retrieve billing credit grant | read | `id` |  |  |
| `STRIPE_GET_BILLING_METERS` | List billing meters | read | *ending_before, limit, starting_after, status* |  | `read_only` |
| `STRIPE_GET_BILLING_METERS_ID` | Retrieve Billing Meter | read | `id` |  | `read_only` |
| `STRIPE_GET_BILLING_PORTAL_CONFIGURATIONS` | List portal configurations | read | *active, ending_before, is_default, limit, starting_after* |  |  |
| `STRIPE_GET_BILLING_PORTAL_CONFIGURATIONS_CONFIGURATION` | Get billing portal configuration | read | `configuration` |  | `read_only` |
| `STRIPE_GET_CHARGES_CHARGE_DISPUTE` | Retrieve Charge Dispute | read | `charge` |  |  |
| `STRIPE_GET_CHARGES_CHARGE_REFUNDS_REFUND` | Get Charge Refund Details | read | `charge`, `refund`; *expand* |  |  |
| `STRIPE_GET_CHECKOUT_SESSIONS_SESSION_LINE_ITEMS` | Get Checkout Session Line Items | read | `session`; *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_GET_CLIMATE_PRODUCTS_PRODUCT` | Retrieve climate product | read | `product` |  |  |
| `STRIPE_GET_CLIMATE_SUPPLIERS_SUPPLIER` | Retrieve climate supplier details | read | `supplier` |  |  |
| `STRIPE_GET_CONFIRMATION_TOKENS_CONFIRMATION_TOKEN` | Retrieve Confirmation Token | read | `confirmation_token`; *expand* |  |  |
| `STRIPE_GET_COUNTRY_SPECS_COUNTRY` | Retrieve Country Spec | read | `country` |  |  |
| `STRIPE_GET_COUPON` | Retrieve Coupon Details | read | `coupon_id` |  | `read_only` |
| `STRIPE_GET_CREDIT_NOTES_CREDIT_NOTE_LINES` | Get Credit Note Lines | read | `credit_note`; *ending_before, limit, starting_after* |  |  |
| `STRIPE_GET_CREDIT_NOTES_ID` | Retrieve a credit note | read | `id`; *expand* |  | `read_only` |
| `STRIPE_GET_CREDIT_NOTES_PREVIEW` | Preview Credit Note | read | `invoice`; *amount, credit_amount, effective_at, email_type, lines, memo, metadata, out_of_band_amount and 4 more* |  | `read_write` |
| `STRIPE_GET_CREDIT_NOTES_PREVIEW_LINES` | Get Credit Note Preview Lines | read | `invoice`; *amount, credit_amount, effective_at, email_type, ending_before, limit, lines, memo and 7 more* |  | `read_only` |
| `STRIPE_GET_CUSTOMERS_CASH_BALANCE_TRANSACTIONS` | Retrieve customer cash balance transaction | read | `customer`, `transaction`; *expand* |  |  |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_BALANCE_TRANSACTIONS` | List customer balance transactions | read | `customer`; *created, ending_before, invoice, limit, starting_after* |  | `read_only` |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_BALANCE_TXNS_TXN` | Retrieve Customer Balance Transaction | read | `customer`, `transaction` |  |  |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_BANK_ACCOUNTS` | List customer bank accounts | read | `customer_id`; *ending_before, limit, starting_after* |  |  |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_BANK_ACCOUNTS_ID` | Retrieve customer bank account | read | `customer`, `id`; *expand* |  | `read_only` |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_CARDS` | List all cards | read | `customer`; *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_CARDS_ID` | Retrieve customer card | read | `customer`, `id` |  |  |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_CASH_BALANCE` | Get customer cash balance | read | `customer`; *expand* |  | `customer_read`, `read_only`, `read_write` |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_CASH_BALANCE_TRANSACTIONS` | List cash balance transactions | read | `customer_id`; *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_DISCOUNT` | Get customer discount | read | `customer` |  |  |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_PAYMENT_METHODS` | List customer payment methods | read | `customer_id`; *allow_redisplay, ending_before, limit, starting_after, type* |  |  |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_PAYMENT_METHODS_PAYMENT_METHOD` | Retrieve customer's payment method | read | `customer`, `payment_method`; *expand* |  | `read_only` |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_SOURCES` | List customer sources | read | `customer`; *ending_before, limit, starting_after* |  |  |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_SUBS_SUB_EXPOSED_ID` | Retrieve customer subscription | read | `customer`, `subscription_exposed_id` |  | `read_only` |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_SUBS_SUB_EXPOSED_ID_DISCOUNT` | Get customer subscription discount | read | `customer`, `subscription_exposed_id`; *expand* |  |  |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_TAX_IDS` | List customer tax IDs | read | `customer`; *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_GET_CUSTOMERS_CUSTOMER_TAX_IDS_ID` | Retrieve a customer tax ID | read | `customer`, `id` |  | `read_only` |
| `STRIPE_GET_ENTITLEMENTS_ACTIVE_ENTITLEMENTS` | List active entitlements | read | `customer`; *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_GET_ENTITLEMENTS_FEATURES_ID` | Retrieve entitlements feature | read | `id` |  | `read_only` |
| `STRIPE_GET_EVENTS` | List events | read | *created, delivery_success, ending_before, limit, starting_after, type, types* |  | `read_only` |
| `STRIPE_GET_EVENTS_ID` | Retrieve event | read | `id` |  | `read_only` |
| `STRIPE_GET_FILES_FILE` | Retrieve File | read | `file` |  | `read_only` |
| `STRIPE_GET_FILE_LINKS_LINK` | Retrieve file link | read | `file_link_id`; *expand* |  | `read_only` |
| `STRIPE_GET_FINANCIAL_CONNECTIONS_SESSIONS_SESSION` | Retrieve Financial Connections Session | read | `session` |  | `read_only` |
| `STRIPE_GET_FX_QUOTE` | Retrieve FX Quote | read | `id`; *expand* |  |  |
| `STRIPE_GET_INVOICEITEMS_INVOICEITEM` | Retrieve Invoice Item | read | `invoiceitem`; *expand* |  | `read_only` |
| `STRIPE_GET_INVOICES_INVOICE` | Retrieve Invoice | read | `invoice`; *expand* |  | `read_only` |
| `STRIPE_GET_INVOICES_INVOICE_LINES` | Get Invoice Line Items | read | `invoice`; *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_GET_INVOICE_PAYMENTS_ID` | Retrieve an invoice payment | read | `id` |  | `read_only` |
| `STRIPE_GET_LINKED_ACCOUNT` | Retrieve Linked Account | read | `account`; *expand* |  |  |
| `STRIPE_GET_MANDATE` | Retrieve mandate | read | `mandate_id`; *expand* |  | `read_only` |
| `STRIPE_GET_MONTHLY_RECURRING_REVENUE` | Get monthly recurring revenue | read | `ends_at`, `granularity`, `starts_at`; *currency, filters, group_by, timezone* |  | `read_only` |
| `STRIPE_GET_PAYMENT_ATTEMPT_RECORDS_ID` | Retrieve Payment Attempt Record | read | `id`; *expand* |  | `read_only` |
| `STRIPE_GET_PAYMENT_INTENT_AMOUNT_DETAILS_LINE_ITEMS` | Get payment intent amount details line items | read | `payment_intent_id`; *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_GET_PAYMENT_LINK` | Retrieve payment link | read | `payment_link_id`; *expand* |  | `read_only` |
| `STRIPE_GET_PAYMENT_LINKS_PAYMENT_LINK_LINE_ITEMS` | Get payment link line items | read | `payment_link`; *ending_before, limit, starting_after* |  |  |
| `STRIPE_GET_PAYMENT_METHOD` | Retrieve payment method | read | `payment_method_id`; *expand* |  |  |
| `STRIPE_GET_PAYMENT_METHOD_CONFIGURATION` | Retrieve payment method configuration | read | `id` |  |  |
| `STRIPE_GET_PAYMENT_METHOD_DOMAIN` | Retrieve Payment Method Domain | read | `payment_method_domain` |  | `read_only` |
| `STRIPE_GET_PAYMENT_RECORD` | Retrieve Payment Record | read | `id`; *expand* |  | `read_only` |
| `STRIPE_GET_PAYOUT` | Retrieve a payout | read | `payout_id` |  | `read_only` |
| `STRIPE_GET_PLAN` | Retrieve plan | read | `plan_id` |  | `read_only` |
| `STRIPE_GET_PRICES_PRICE` | Retrieve a price | read | `price` |  |  |
| `STRIPE_GET_PRODUCTS_ID` | Retrieve a product | read | `id`; *expand* |  | `read_only` |
| `STRIPE_GET_PRODUCTS_PRODUCT_FEATURES_ID` | Retrieve product feature | read | `id`, `product` |  | `read_only` |
| `STRIPE_GET_PROMOTION_CODES` | List Stripe promotion codes | write | *active, code, coupon, created, customer, customer_account, ending_before, limit and 1 more* |  | `read_only` |
| `STRIPE_GET_QUOTES_QUOTE` | Get quote details | read | `quote`; *expand* |  | `read_write` |
| `STRIPE_GET_QUOTES_QUOTE_COMPUTED_UPFRONT_LINE_ITEMS` | Get Quote Computed Upfront Line Items | read | `quote`; *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_GET_QUOTES_QUOTE_LINE_ITEMS` | Get Quote Line Items | read | `quote`; *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_GET_RADAR_VALUE_LIST` | Get Radar Value List | read | `id` |  |  |
| `STRIPE_GET_RADAR_VALUE_LIST_ITEMS` | List Radar value list items | read | `value_list`; *created, ending_before, limit, starting_after, value* |  | `read_only` |
| `STRIPE_GET_REFUNDS_REFUND` | Retrieve a refund | read | `refund` |  | `read_write` |
| `STRIPE_GET_REPORTING_REPORT_RUNS_REPORT_RUN` | Retrieve a Report Run | read | `report_run` |  |  |
| `STRIPE_GET_REPORTING_REPORT_TYPE` | Retrieve a report type | read | `id` |  | `read_only` |
| `STRIPE_GET_SETUP_ATTEMPTS` | List setup attempts | read | `setup_intent`; *created, ending_before, limit, starting_after* |  | `read_only`, `read_write` |
| `STRIPE_GET_SETUP_INTENT` | Retrieve setup intent | read | `setup_intent_id`; *client_secret* |  | `read_only` |
| `STRIPE_GET_SHIPPING_RATE` | Retrieve Shipping Rate | read | `shipping_rate_id` |  | `read_only` |
| `STRIPE_GET_SUBSCRIPTION` | Retrieve subscription | read | `subscription_id` |  |  |
| `STRIPE_GET_SUBSCRIPTION_ITEM` | Retrieve subscription item | read | `subscription_item_id` |  | `read_only` |
| `STRIPE_GET_SUBSCRIPTION_SCHEDULES_SCHEDULE` | Retrieve subscription schedule | read | `schedule`; *expand* |  | `read_only` |
| `STRIPE_GET_TAX_CALCULATION` | Retrieve Tax Calculation | read | `id` |  |  |
| `STRIPE_GET_TAX_CALCULATIONS_CALCULATION_LINE_ITEMS` | Get Tax Calculation Line Items | read | `calculation`; *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_GET_TAX_CODE` | Retrieve tax code | read | `id` |  |  |
| `STRIPE_GET_TAX_ID` | Retrieve tax ID | read | `id` |  |  |
| `STRIPE_GET_TAX_RATE` | Retrieve a tax rate | read | `tax_rate` |  | `read_write` |
| `STRIPE_GET_TAX_REGISTRATION` | Retrieve tax registration | read | `id` |  |  |
| `STRIPE_GET_TAX_SETTINGS` | Get Tax Settings | read | *expand* |  |  |
| `STRIPE_GET_TAX_TRANSACTION` | Retrieve Tax Transaction | read | `transaction` |  |  |
| `STRIPE_GET_TAX_TRANSACTIONS_TRANSACTION_LINE_ITEMS` | Get tax transaction line items | read | `transaction`; *ending_before, expand, limit, starting_after* |  | `read_only` |
| `STRIPE_GET_TERMINAL_CONFIGURATION` | Retrieve terminal configuration | read | `configuration` |  | `read_write` |
| `STRIPE_GET_TERMINAL_LOCATION` | Retrieve Terminal Location | read | `location` |  | `read_only` |
| `STRIPE_GET_TERMINAL_READER` | Retrieve Terminal Reader | read | `reader`; *expand* |  | `read_only` |
| `STRIPE_GET_TEST_CLOCK` | Retrieve test clock | read | `test_clock` |  | `read_only` |
| `STRIPE_GET_TOKEN` | Retrieve Token | read | `token` |  | `read_only` |
| `STRIPE_GET_TRANSFER_REVERSAL` | Retrieve transfer reversal | read | `id`, `transfer`; *expand* |  | `read_only` |
| `STRIPE_GET_UPCOMING_INVOICE_LINE_ITEMS` | Retrieve Upcoming Invoice Line Items | read | *automatic_tax, coupon, currency, customer, discounts, ending_before, invoice_items, limit and 4 more* |  |  |
| `STRIPE_GET_V1_CUSTOMERS_SEARCH_CUSTOMERS` | Search Stripe customers | read | `query`; *limit, page* |  | `read_write`, `read_only`, `customer_read`, `customer_write` |
| `STRIPE_GET_V1_DISPUTES_DISPUTES` | Retrieve Dispute Details | read | `id`; *expand* |  | `read_only` |
| `STRIPE_GET_V1_INVOICE_PAYMENTS_OVERVIEW2` | List Invoice Payments | read | *ending_before, invoice, limit, starting_after, status* |  |  |
| `STRIPE_GET_V1_PAYMENT_METHODS_PAYMENT_METHODS2` | List payment methods | read | *allow_redisplay, customer, customer_account, ending_before, limit, starting_after, type* |  | `read_write` |
| `STRIPE_GET_V1_PAYMENT_METHOD_CONFIGS_LIST_PAYMENT_METHOD2` | List payment method configurations | read | *application, ending_before, limit, starting_after* |  |  |
| `STRIPE_GET_V1_PAYMENT_METHOD_DOMAINS` | List payment method domains | read | *domain_name, enabled, ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_GET_V1_RADAR_EARLY_FRAUD_WARNINGS_OVERVIEW` | List Radar Early Fraud Warnings | read | *charge, created, ending_before, limit, payment_intent, starting_after* |  | `read_only` |
| `STRIPE_GET_V1_RADAR_VALUE_LISTS_OVERVIEW_RADAR_VALUE_LISTS` | Get Radar value lists overview | read | *alias, contains, created, ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_GET_V1_REPORTING_REPORT_TYPES_OVERVIEW_REPORTING` | Get reporting report types overview | read |  |  | `read_only` |
| `STRIPE_GET_V1_SHIPPING_RATES_OVERVIEW2` | Get Shipping Rates Overview | read | *active, currency, ending_before, limit, starting_after* |  | `shipping_rate_read` |
| `STRIPE_GET_V1_SUBSCRIPTION_SCHEDULES_OVERVIEW2` | List subscription schedules | read | *canceled_at, completed_at, created, customer, customer_account, ending_before, limit, released_at and 2 more* |  | `read_only` |
| `STRIPE_GET_V1_TAX_CODES_TAX_CODE2` | Get tax codes | read | *ending_before, limit, starting_after* |  | `read_only`, `read_write` |
| `STRIPE_GET_V1_TAX_IDS_OVERVIEW2` | List tax IDs | read | *ending_before, limit, owner_account, owner_customer, owner_customer_account, owner_type, starting_after* |  | `read_only` |
| `STRIPE_GET_V1_TAX_RATES_TAX_RATE2` | Get tax rates | read | *active, created, ending_before, inclusive, limit, starting_after* |  | `tax_rate_read` |
| `STRIPE_GET_V2_CORE_EVENTS_RETRIEVE_AN_EVENT` | Retrieve V2 Core Event | read | `id` |  |  |
| `STRIPE_LIST_ALL_SUBSCRIPTION_ITEM_PERIOD_SUMMARIES` | List subscription item period summaries | read | `subscription_item`; *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_LIST_APPLE_PAY_DOMAINS` | List Apple Pay domains | read | *ending_before, limit, starting_after* |  |  |
| `STRIPE_LIST_APPLICATION_FEES` | List Application Fees | read | *charge, created, ending_before, limit, starting_after* |  | `read_write` |
| `STRIPE_LIST_APPS_SECRETS` | List apps secrets | read | `scope_type`; *ending_before, limit, scope_user, starting_after* |  | `read_only` |
| `STRIPE_LIST_BALANCE_TRANSACTIONS` | List Balance Transactions | read | *created, currency, ending_before, expand, limit, payout, source, starting_after and 1 more* |  |  |
| `STRIPE_LIST_BILLING_ALERTS` | List billing alerts | read | *alert_type, ending_before, limit, meter, starting_after* |  | `read_only` |
| `STRIPE_LIST_BILLING_CREDIT_BALANCE_TRANSACTIONS` | List credit balance transactions | read | *credit_grant, customer, customer_account, ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_LIST_CHARGES` | List Charges | read | *created, customer, ending_before, limit, payment_intent, starting_after, transfer_group* |  |  |
| `STRIPE_LIST_CHARGE_REFUNDS` | List Charge Refunds | read | `charge`; *ending_before, expand, limit, starting_after* |  |  |
| `STRIPE_LIST_CHECKOUT_SESSIONS` | List Checkout Sessions | read | *created, customer, customer_account, customer_details_email, ending_before, limit, payment_intent, payment_link and 3 more* |  |  |
| `STRIPE_LIST_CLIMATE_ORDERS` | List Climate Orders | read | *ending_before, limit, starting_after* |  |  |
| `STRIPE_LIST_CLIMATE_PRODUCTS` | List climate products | read | *ending_before, limit, starting_after* |  |  |
| `STRIPE_LIST_CLIMATE_SUPPLIERS` | List climate suppliers | read | *ending_before, limit, starting_after* |  |  |
| `STRIPE_LIST_COUNTRY_SPECS` | List country specs | read | *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_LIST_COUPONS` | List Stripe coupons | read | *ending_before, limit, starting_after* |  |  |
| `STRIPE_LIST_CREDIT_NOTES` | List Credit Notes | read | *created, customer, customer_account, ending_before, invoice, limit, starting_after* |  | `read_write` |
| `STRIPE_LIST_CUSTOMERS` | List customers | read | *created, email, ending_before, limit, starting_after, test_clock* |  |  |
| `STRIPE_LIST_CUSTOMER_SUBSCRIPTIONS` | List customer subscriptions | read | `customer`; *ending_before, expand, limit, starting_after* |  | `read_only` |
| `STRIPE_LIST_DISPUTES` | List Disputes | read | *charge, created, ending_before, limit, payment_intent, starting_after* |  |  |
| `STRIPE_LIST_ENTITLEMENTS_FEATURES` | List entitlements features | read | *archived, ending_before, limit, lookup_key, starting_after* |  | `read_only` |
| `STRIPE_LIST_FILES` | List files | read | *created, ending_before, limit, purpose, starting_after* |  | `read_only` |
| `STRIPE_LIST_FILE_LINKS` | List file links | read | *created, ending_before, expired, file, limit, starting_after* |  | `read_only` |
| `STRIPE_LIST_FINANCIAL_CONNECTIONS_ACCOUNTS` | List Financial Connections Accounts | read | *account_holder, ending_before, limit, session, starting_after* |  |  |
| `STRIPE_LIST_FORWARDING_REQUESTS` | List forwarding requests | read | *created, ending_before, limit, starting_after* |  |  |
| `STRIPE_LIST_FX_QUOTES` | List FX quotes | read | *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_LIST_INVOICEITEMS` | List Invoice Items | read | *created, customer, customer_account, ending_before, invoice, limit, pending, starting_after* |  |  |
| `STRIPE_LIST_INVOICES` | List Invoices | read | *collection_method, created, customer, ending_before, limit, starting_after, status, subscription* |  |  |
| `STRIPE_LIST_INVOICE_PAYMENTS` | List Invoice Payments | read | *created, ending_before, invoice, limit, starting_after, status* |  |  |
| `STRIPE_LIST_INVOICE_RENDERING_TEMPLATES` | List Invoice Rendering Templates | read | *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_LIST_METER_EVENT_SUMMARIES` | List meter event summaries | read | `customer`, `end_time`, `id`, `start_time`; *ending_before, limit, starting_after, value_grouping_window* |  | `read_only` |
| `STRIPE_LIST_PAYMENT_ATTEMPT_RECORDS` | List Payment Attempt Records | read | `payment_record`; *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_LIST_PAYMENT_INTENTS` | List payment intents | read | *created, customer, customer_account, ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_LIST_PAYMENT_LINKS` | List payment links | read | *active, ending_before, limit, starting_after* |  |  |
| `STRIPE_LIST_PAYOUTS` | List Payouts | read | *arrival_date, created, destination, ending_before, limit, starting_after, status* |  | `read_write` |
| `STRIPE_LIST_PLANS` | List plans | read | *active, created, ending_before, limit, product, starting_after* |  |  |
| `STRIPE_LIST_PRICES` | List prices | read | *active, created, currency, ending_before, expand, limit, lookup_keys, product and 3 more* |  |  |
| `STRIPE_LIST_PRODUCTS` | List products | read | *active, created, ending_before, ids, limit, shippable, starting_after, url* |  |  |
| `STRIPE_LIST_PRODUCTS_PRODUCT_FEATURES` | List product features | read | `product`; *ending_before, limit, starting_after* |  |  |
| `STRIPE_LIST_PROMOTION_CODES` | List Stripe promotion codes | write | *active, code, coupon, created, customer, customer_account, ending_before, limit and 1 more* |  | `read_write` |
| `STRIPE_LIST_QUOTES` | List Quotes | read | *customer, customer_account, ending_before, limit, starting_after, status, test_clock* |  |  |
| `STRIPE_LIST_RADAR_EARLY_FRAUD_WARNINGS` | Get Radar Early Fraud Warnings | read | *charge, created, ending_before, limit, payment_intent, starting_after* |  | `read_only` |
| `STRIPE_LIST_RADAR_REVIEWS` | List Radar Reviews | read | *created, ending_before, limit, starting_after* |  |  |
| `STRIPE_LIST_REFUNDS` | List Refunds | read | *charge, ending_before, limit, payment_intent, starting_after* |  |  |
| `STRIPE_LIST_REPORTING_REPORT_RUNS` | List Report Runs | read | *created, ending_before, limit, starting_after* |  | `report_runs_and_report_types_read` |
| `STRIPE_LIST_REVIEWS` | List open reviews | read | *created, ending_before, limit, starting_after* |  |  |
| `STRIPE_LIST_SETUP_ATTEMPTS` | List SetupAttempts | read | `setup_intent`; *created, ending_before, limit, starting_after* |  |  |
| `STRIPE_LIST_SETUP_INTENTS` | List setup intents | read | *attach_to_self, created, customer, ending_before, limit, payment_method, starting_after* |  | `read_only` |
| `STRIPE_LIST_SIGMA_SCHEDULED_QUERY_RUNS` | List Sigma scheduled query runs | read | *ending_before, expand, limit, starting_after* |  | `read_only` |
| `STRIPE_LIST_SUBSCRIPTIONS` | List subscriptions | read | *automatic_tax, collection_method, created, current_period_end, current_period_start, customer, ending_before, limit and 4 more* |  |  |
| `STRIPE_LIST_SUBSCRIPTION_ITEMS` | List subscription items | read | `subscription`; *ending_before, limit, starting_after* |  |  |
| `STRIPE_LIST_TAX_CODES` | List tax codes | read | *ending_before, limit, starting_after* |  |  |
| `STRIPE_LIST_TAX_RATES` | List tax rates | read | *active, created, ending_before, inclusive, limit, starting_after* |  |  |
| `STRIPE_LIST_TAX_REGISTRATIONS` | List tax registrations | read | *ending_before, limit, starting_after, status* |  |  |
| `STRIPE_LIST_TERMINAL_CONFIGURATIONS` | List terminal configurations | read | *ending_before, is_account_default, limit, starting_after* |  | `read_write` |
| `STRIPE_LIST_TERMINAL_LOCATIONS` | List terminal locations | read | *ending_before, limit, starting_after* |  | `read_write` |
| `STRIPE_LIST_TERMINAL_READERS` | List terminal readers | read | *device_type, ending_before, limit, location, serial_number, starting_after, status* |  | `read_only` |
| `STRIPE_LIST_TEST_CLOCKS` | List test clocks | read | *ending_before, limit, starting_after* |  | `read_only` |
| `STRIPE_LIST_TOPUPS` | List Top-ups | read | *amount, created, ending_before, limit, starting_after, status* |  | `read_only` |
| `STRIPE_LIST_TRANSFERS` | List transfers | read | *created, destination, ending_before, limit, starting_after, transfer_group* |  | `read_only` |
| `STRIPE_LIST_V2_CORE_EVENTS` | List v2 core events | read | *created, limit, object_id, page, types* |  | `read_write`, `read_only` |
| `STRIPE_LIST_V2_CORE_EVENT_DESTINATIONS` | List V2 Core Event Destinations | read | *include, limit, page* |  | `read_only` |
| `STRIPE_MARK_INVOICE_UNCOLLECTIBLE` | Mark invoice as uncollectible | write | `invoice` |  | `read_write` |
| `STRIPE_MIGRATE_SUBSCRIPTION` | Migrate subscription to flexible billing | write | `billing_mode`, `subscription` |  |  |
| `STRIPE_PAY_INVOICE` | Pay an invoice | write | `invoice`; *forgive, mandate, off_session, paid_out_of_band, payment_method, source* |  | `read_write` |
| `STRIPE_POST_BILLING_CREDIT_GRANTS_VOID_CREDIT_GRANT` | Void billing credit grant | write | `id` |  |  |
| `STRIPE_POST_CHARGES_CHARGE` | Update Charge | write | `charge`; *customer, description, fraud_details, metadata, receipt_email, shipping, transfer_group* |  |  |
| `STRIPE_POST_CUSTOMERS_CUSTOMER_SUBS_SUB_EXPOSED_ID` | Update customer subscription | write | `customer`, `subscription_exposed_id`; *automatic_tax, billing_cycle_anchor, cancel_at, cancel_at_period_end, collection_method, default_payment_method, default_tax_rates, description and 6 more* |  | `subscription_read`, `subscription_write` |
| `STRIPE_POST_CUSTOMERS_SOURCES_VERIFY_THE_BANK_ACCOUNT` | Verify customer bank account | write | `bank_account_id`, `customer_id`; *amounts* |  |  |
| `STRIPE_POST_FILE_LINKS_LINK` | Update file link | write | `link`; *expires_at, metadata* |  | `read_write` |
| `STRIPE_POST_INVOICES_INVOICE_VOID` | Void an invoice | write | `invoice` |  | `invoice_write` |
| `STRIPE_POST_PAYMENT_RECORDS_REPORT_PAYMENT` | Report Payment Record | write | `amount_requested`, `initiated_at`, `payment_method_details`, `processor_details`; *customer_details, customer_presence, description, failed, guaranteed, metadata, outcome, shipping_details* |  |  |
| `STRIPE_POST_SETUP_INTENTS_INTENT` | Update setup intent | write | `intent`; *attach_to_self, customer, customer_account, description, flow_directions, metadata, payment_method, payment_method_configuration and 1 more* |  |  |
| `STRIPE_POST_SUBSCRIPTIONS_SUBSCRIPTION_EXPOSED_ID` | Update subscription | write | `subscription_exposed_id`; *add_invoice_items, application_fee_percent, automatic_tax, billing_cycle_anchor, billing_thresholds, cancel_at, cancel_at_period_end, cancellation_details and 22 more* |  | `read_write` |
| `STRIPE_POST_SUBSCRIPTION_ITEMS_ITEM` | Update subscription item | write | `item`; *billing_thresholds, discounts, metadata, off_session, payment_behavior, price, price_data, proration_behavior and 3 more* |  | `read_write` |
| `STRIPE_POST_TERMINAL_CONFIGURATIONS_CONFIGURATION` | Update Terminal Configuration | write | `configuration`; *bbpos_wisepad3, bbpos_wisepos_e, name, offline, reader_security, reboot_window, stripe_s700, tipping and 2 more* |  |  |
| `STRIPE_POST_TERMINAL_READERS_READER_PROCESS_PAYMENT_INTENT` | Process payment on terminal reader | write | `payment_intent`, `reader`; *process_config* |  | `read_write` |
| `STRIPE_POST_TERMINAL_READERS_READER_SET_READER_DISPLAY` | Set Reader Display | write | `cart`, `reader`, `type` |  | `read_write` |
| `STRIPE_POST_TERMINAL_READERS_SET_READER_DISPLAY` | Set Terminal Reader Display | write | `cart`, `reader_id`, `type` |  |  |
| `STRIPE_POST_V1_ACCOUNTS_CAPABILITIES_CAPABILITIES` | Update Account Capability | write | `account`, `capability`; *experimental, requested* |  | `read_write` |
| `STRIPE_POST_V1_BALANCE_SETTINGS_OVERVIEW` | Update balance settings overview | write | *payments, stripe_account* |  | `read_write` |
| `STRIPE_POST_V1_BILLING_PORTAL_CONFIGS_OVERVIEW_BILLING` | Update Billing Portal Configuration | write | `configuration_id`; *active, business_profile, default_return_url, features, login_page, metadata, name* |  | `read_write` |
| `STRIPE_POST_V1_CHECKOUT_SESSIONS_RETRIEVE_A_CHECKOUT` | Update and retrieve Checkout Session | write | `id`; *expand, metadata* |  |  |
| `STRIPE_POST_V1_CREDIT_NOTES_CREDIT_NOTE2` | Update Credit Note | write | `credit_note_id`; *memo, metadata* |  | `read_only` |
| `STRIPE_POST_V1_CUSTOMERS_FUNDING_INSTRUCTIONS_FUND_A_TEST` | Create or Retrieve Customer Funding Instructions | write | `bank_transfer`, `currency`, `customer_id`, `funding_type` |  | `read_write` |
| `STRIPE_POST_V1_CUSTOMERS_SEARCH_CUSTOMERS` | Search Stripe customers | read | `query`; *limit, page* |  | `customer_read` |
| `STRIPE_POST_V1_INVOICES_FINALIZE_VOID_AN_INVOICE` | Void an invoice | write | `invoice_id` |  |  |
| `STRIPE_POST_V1_INVOICES_LINES_OVERVIEW` | Update Invoice Line Item | write | `invoice`, `line_item_id`; *amount, description, discountable, discounts, metadata, period, price_data, pricing and 3 more* |  | `read_only` |
| `STRIPE_POST_V1_INVOICES_REMOVE_LINES_OVERVIEW` | Remove lines from invoice | write | `invoice_id`, `lines`; *invoice_metadata* |  | `read_write` |
| `STRIPE_POST_V1_INVOICES_SEND_VOID_AN_INVOICE` | Void invoice | write | `invoice_id` |  | `invoice_write`, `read_write` |
| `STRIPE_POST_V1_INVOICES_VOID_AN_INVOICE` | Void an invoice | write | `invoice_id` |  |  |
| `STRIPE_POST_V1_PAYMENT_METHOD_CONFIGS_LIST_PAYMENT_METHOD2` | Update payment method configuration | write | `id`; *acss_debit, active, affirm, afterpay_clearpay, alipay, alma, amazon_pay, apple_pay and 48 more* |  | `read_write` |
| `STRIPE_POST_V1_PAYMENT_METHOD_DOMAINS_PAYMENT_METHOD` | Update Payment Method Domain | write | `id`; *enabled* |  | `read_write` |
| `STRIPE_POST_V1_PAYMENT_RECORDS_REPORT_ATTEMPT` | Report Payment Attempt | write | `initiated_at`, `outcome`, `payment_record_id`; *guaranteed* |  | `read_write` |
| `STRIPE_POST_V1_PAYMENT_RECORDS_REPORT_CANCELED` | Report Payment Attempt Canceled | write | `canceled_at`, `id` |  | `read_write` |
| `STRIPE_POST_V1_PAYMENT_RECORDS_REPORT_GUARANTEED` | Report payment attempt guaranteed overview | write | `guaranteed_at`, `id` |  | `read_write` |
| `STRIPE_POST_V1_PAYMENT_RECORDS_REPORT_INFO` | Report payment attempt information | write | `description`, `payment_record_id` |  |  |
| `STRIPE_POST_V1_PAYMENT_RECORDS_REPORT_REFUND_OVERVIEW` | Report Payment Record Refund | write | `id`, `outcome`, `processor_details`, `refunded`; *amount, initiated_at, metadata* |  | `read_write` |
| `STRIPE_POST_V1_PAYOUTS_RETRIEVE_A_PAYOUT` | Update payout metadata | write | `payout_id`; *metadata* |  | `read_only` |
| `STRIPE_POST_V1_PLANS_OVERVIEW` | Update Plan Overview | write | `id`; *active, metadata, nickname, product, trial_period_days* |  | `read_write` |
| `STRIPE_POST_V1_PROMOTION_CODES_PROMOTION_CODE` | Update Promotion Code | write | `promotion_code`; *active, metadata, restrictions* |  |  |
| `STRIPE_POST_V1_SHIPPING_RATES_OVERVIEW2` | Update Shipping Rate | write | `id`; *active, fixed_amount, metadata* |  | `read_write` |
| `STRIPE_POST_V1_SUBSCRIPTIONS_LIST_SUBSCRIPTIONS` | Update subscription | write | `subscription_id`; *billing_cycle_anchor, cancel_at_period_end, collection_method, days_until_due, default_payment_method, default_source, description, metadata and 2 more* |  | `read_write` |
| `STRIPE_POST_V1_SUBSCRIPTION_SCHEDULES_OVERVIEW2` | Update subscription schedule | write | `id`; *default_settings, end_behavior, metadata, phases, proration_behavior* |  | `read_write` |
| `STRIPE_POST_V1_SUBSCRIPTION_SCHEDULES_RELEASE_OVERVIEW` | Release subscription schedule | write | `schedule_id`; *preserve_cancel_date* |  | `read_write` |
| `STRIPE_POST_V1_TAX_RATES_TAX_RATE2` | Update Tax Rate | write | `tax_rate_id`; *active, country, description, display_name, jurisdiction, metadata, state, tax_type* |  | `read_write` |
| `STRIPE_POST_V2_BILLING_METER_EVENT_ADJUSTMENTS_METER_EVENT` | Create Meter Event Adjustment | write | `cancel`, `event_name`, `type` |  | `billing_meter_event_write` |
| `STRIPE_PRESENT_PAYMENT_METHOD` | Present payment method on reader | write | `reader_id`; *amount_tip, card, card_present, interac_present, type* |  |  |
| `STRIPE_PROCESS_SETUP_INTENT` | Process setup intent on terminal reader | write | `customer_consent_collected`, `reader_id`, `setup_intent`; *allow_redisplay, process_config* |  | `read_write` |
| `STRIPE_REACTIVATE_BILLING_METER` | Reactivate Billing Meter | write | `id` |  | `read_write` |
| `STRIPE_REPORT_PAYMENT_ATTEMPT_FAILED` | Report Payment Attempt Failed | write | `failed_at`, `id` |  | `read_write` |
| `STRIPE_RESUME_SUBSCRIPTION` | Resume subscription | write | `subscription_id`; *billing_cycle_anchor, proration_behavior, proration_date* |  | `read_write` |
| `STRIPE_RETRIEVE_BALANCE` | Retrieve Balance | read |  |  |  |
| `STRIPE_RETRIEVE_CHARGE` | Retrieve Charge Details | read | `charge_id` |  |  |
| `STRIPE_RETRIEVE_CHECKOUT_SESSION` | Retrieve Checkout Session | read | `session`; *expand* |  | `read_only` |
| `STRIPE_RETRIEVE_CUSTOMER` | Retrieve customer | read | `customer_id` |  |  |
| `STRIPE_RETRIEVE_PAYMENT_INTENT` | Retrieve payment intent | read | `payment_intent_id`; *client_secret* |  | `read_only` |
| `STRIPE_RETRIEVE_PROMOTION_CODE` | Retrieve Promotion Code | read | `promotion_code_id` |  | `read_only` |
| `STRIPE_RETRIEVE_UPCOMING_INVOICE` | Retrieve Upcoming Invoice | read | *automatic_tax, coupon, currency, customer, customer_details, discounts, invoice_items, issuer and 18 more* |  | `read_only` |
| `STRIPE_SEARCH_CHARGES` | Search Stripe charges | read | `query`; *limit, page* |  |  |
| `STRIPE_SEARCH_CUSTOMERS` | Search Stripe customers | read | `query`; *limit, page* |  |  |
| `STRIPE_SEARCH_INVOICES` | Search invoices | read | `query`; *limit, page* |  | `read_write` |
| `STRIPE_SEARCH_PAYMENT_INTENTS` | Search payment intents | read | `query`; *limit, page* |  | `read_write` |
| `STRIPE_SEARCH_PRICES` | Search Stripe prices | read | `query`; *limit, page* |  |  |
| `STRIPE_SEARCH_PRODUCTS` | Search Stripe products | read | `query`; *limit, page* |  | `read_write` |
| `STRIPE_SEARCH_SUBSCRIPTIONS` | Search subscriptions | read | `query`; *limit, page* |  | `read_write` |
| `STRIPE_SEND_INVOICE` | Send invoice for manual payment | write | `invoice` |  | `read_write` |
| `STRIPE_SET_READER_DISPLAY` | Set reader display | write | `cart`, `reader`; *type* |  |  |
| `STRIPE_SET_TERMINAL_READER_DISPLAY` | Set terminal reader display | write | `reader_id`, `type`; *cart* |  |  |
| `STRIPE_SUCCEED_TEST_TERMINAL_READER_INPUT_COLLECTION` | Simulate successful input collection | write | `reader`; *expand, skip_non_required_inputs* |  |  |
| `STRIPE_TIMEOUT_TEST_TERMINAL_READER_INPUT_COLLECTION` | Simulate Terminal Reader Input Collection Timeout | write | `reader`; *expand* |  |  |
| `STRIPE_UPDATE_APPLICATION_FEE_REFUND` | Update Application Fee Refund | write | `fee`, `id`; *expand, metadata* |  | `read_write` |
| `STRIPE_UPDATE_BILLING_CREDIT_GRANT` | Update billing credit grant | write | `id`; *expires_at, metadata* |  |  |
| `STRIPE_UPDATE_BILLING_METER` | Update Billing Meter | write | `id`; *display_name, expand* |  | `read_write` |
| `STRIPE_UPDATE_BILLING_PORTAL_CONFIGURATION` | Update Billing Portal Configuration | write | `configuration`; *active, business_profile, default_return_url, features, login_page, metadata, name* |  |  |
| `STRIPE_UPDATE_CASH_BALANCE` | Update Cash Balance Settings | write | `customer_id`; *settings* |  | `read_write` |
| `STRIPE_UPDATE_CHARGE` | Update Charge | write | `charge_id`; *customer, description, fraud_details, metadata, receipt_email, shipping, transfer_group* |  |  |
| `STRIPE_UPDATE_CHARGE_DISPUTE` | Update Charge Dispute | write | `charge`; *metadata* |  |  |
| `STRIPE_UPDATE_CHARGE_REFUND` | Update Charge Refund | write | `refund`; *metadata* |  |  |
| `STRIPE_UPDATE_CHECKOUT_SESSION` | Update Checkout Session | write | `session`; *collected_information, line_items, metadata, shipping_options* |  | `read_write` |
| `STRIPE_UPDATE_COUPON` | Update Coupon | write | `coupon`; *currency_options, metadata, name* |  | `read_write` |
| `STRIPE_UPDATE_CUSTOMER` | Update Customer | write | `customer_id`; *address, balance, coupon, default_payment_method, default_source, description, email, invoice_settings and 6 more* |  |  |
| `STRIPE_UPDATE_CUSTOMER_BALANCE_TXN` | Update Customer Balance Transaction | write | `customer_id`, `transaction_id`; *description, metadata* |  | `read_write` |
| `STRIPE_UPDATE_CUSTOMER_BANK_ACCOUNT` | Update customer bank account | write | `customer`, `id`; *account_holder_name, account_holder_type, metadata* |  |  |
| `STRIPE_UPDATE_CUSTOMER_BANK_ACCOUNT_SOURCE` | Update customer bank account | write | `customer_id`, `source_id`; *account_holder_name, account_holder_type, metadata* |  | `read_write` |
| `STRIPE_UPDATE_CUSTOMER_CARD` | Update customer card | write | `customer`, `id`; *address_city, address_country, address_line1, address_line2, address_state, address_zip, exp_month, exp_year and 2 more* |  | `read_write` |
| `STRIPE_UPDATE_CUSTOMER_SOURCE` | Update customer source | write | `customer`, `id`; *account_holder_name, account_holder_type, address_city, address_country, address_line1, address_line2, address_state, address_zip and 6 more* |  | `read_only` |
| `STRIPE_UPDATE_CUSTOMER_SUBSCRIPTION` | Update Customer Subscription | write | `customer`, `subscription_exposed_id`; *automatic_tax, billing_cycle_anchor, billing_thresholds, cancel_at, cancel_at_period_end, cancellation_details, collection_method, days_until_due and 15 more* |  | `read_write` |
| `STRIPE_UPDATE_DISPUTE` | Update Dispute | write | `dispute`; *evidence, metadata, submit* |  | `read_write` |
| `STRIPE_UPDATE_ENTITLEMENT_FEATURE` | Update entitlements feature | write | `id`; *active, metadata, name* |  | `read_write` |
| `STRIPE_UPDATE_EVENT_DESTINATION` | Update Event Destination | write | `event_destination_id`; *description, enabled_events, include, metadata, name, webhook_endpoint* |  | `EventDestinations:Write` |
| `STRIPE_UPDATE_EXTERNAL_ACCOUNT` | Update External Account | write | `account`, `id`; *account_holder_name, account_holder_type, account_type, address_city, address_country, address_line1, address_line2, address_state and 5 more* |  | `read_write`, `read_only` |
| `STRIPE_UPDATE_INVOICE` | Update Invoice | write | `invoice`; *account_tax_ids, application_fee_amount, auto_advance, automatic_tax, automatically_finalizes_at, collection_method, custom_fields, days_until_due and 18 more* |  | `read_write` |
| `STRIPE_UPDATE_INVOICEITEM` | Update invoice item | write | `invoiceitem_id`; *amount, description, discountable, discounts, metadata, period, price_data, pricing and 5 more* |  | `read_write` |
| `STRIPE_UPDATE_INVOICE_LINES` | Bulk update invoice line items | write | `invoice`, `lines`; *invoice_metadata* |  |  |
| `STRIPE_UPDATE_PAYMENT_INTENT` | Update Payment Intent | write | `payment_intent_id`; *amount, currency, customer, description, metadata, payment_method, receipt_email, setup_future_usage and 3 more* |  |  |
| `STRIPE_UPDATE_PAYMENT_LINK` | Update Payment Link | write | `payment_link`; *active, after_completion, allow_promotion_codes, automatic_tax, billing_address_collection, custom_fields, custom_text, customer_creation and 14 more* |  | `read_write` |
| `STRIPE_UPDATE_PAYMENT_METHOD` | Update Payment Method | write | `payment_method_id`; *allow_redisplay, billing_details, card, metadata, payto, us_bank_account* |  |  |
| `STRIPE_UPDATE_PRICE` | Update a price | write | `price`; *active, currency_options, lookup_key, metadata, nickname, tax_behavior, transfer_lookup_key* |  | `read_write` |
| `STRIPE_UPDATE_PRODUCT` | Update product | write | `id`; *active, default_price, description, images, marketing_features, metadata, name, package_dimensions and 5 more* |  | `read_write` |
| `STRIPE_UPDATE_PROMOTION_CODE` | Update Promotion Code | write | `id`; *active, metadata, restrictions* |  |  |
| `STRIPE_UPDATE_QUOTE` | Update Quote | write | `quote`; *application_fee_amount, application_fee_percent, automatic_tax, collection_method, customer, customer_account, default_tax_rates, description and 10 more* |  |  |
| `STRIPE_UPDATE_SUBSCRIPTION` | Update Subscription | write | `subscription_id`; *automatic_tax, billing_cycle_anchor, cancel_at_period_end, collection_method, default_payment_method, description, items, metadata and 2 more* |  |  |
| `STRIPE_UPDATE_SUBSCRIPTION_SCHEDULE` | Update Subscription Schedule | write | `schedule`; *default_settings, end_behavior, metadata, phases, proration_behavior* |  | `read_only` |
| `STRIPE_UPDATE_TAX_REGISTRATION` | Update Tax Registration | write | `id`; *active_from, expires_at* |  | `read_write` |
| `STRIPE_UPDATE_TAX_SETTINGS` | Update Tax Settings | write | *defaults, head_office* |  | `read_only` |
| `STRIPE_UPDATE_TERMINAL_LOCATION` | Update Terminal Location | write | `location`; *address, address_kana, address_kanji, configuration_overrides, display_name, display_name_kana, display_name_kanji, metadata and 1 more* |  |  |
| `STRIPE_UPDATE_TERMINAL_READER` | Update Terminal Reader | write | `reader`; *label, metadata* |  |  |
| `STRIPE_VALIDATE_PAYMENT_METHOD_DOMAIN` | Validate payment method domain | write | `payment_method_domain` |  |  |
| `STRIPE_VERIFY_MICRODEPOSITS` | Verify microdeposits on setup intent | write | `setup_intent_id`; *amounts, descriptor_code* |  | `read_write` |
| `STRIPE_VERIFY_MICRODEPOSITS_PAYMENT_INTENT` | Verify microdeposits on payment intent | write | `intent`; *amounts, descriptor_code* |  | `read_write` |
| `STRIPE_VOID_CREDIT_NOTE` | Void a credit note | write | `id` |  | `read_write` |
| `STRIPE_VOID_INVOICE` | Void an invoice | write | `invoice_id` |  |  |
| `STRIPE_WHO_AM_I` | Who Am I | read |  |  |  |
