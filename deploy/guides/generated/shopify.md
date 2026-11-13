# Shopify: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

386 tools: 199 read, 187 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `SHOPIFY_ACCEPT_FULFILLMENT_CANCELLATION_REQUEST` | Accept Fulfillment Cancellation Request | write | `fulfillment_order_id`; *message* |  | `write_assigned_fulfillment_orders` |
| `SHOPIFY_ACCEPT_FULFILLMENT_REQUEST` | Accept Fulfillment Request | write | `fulfillment_order_id`; *message* |  | `write_assigned_fulfillment_orders` |
| `SHOPIFY_ADD_PRODUCT_TO_CUSTOM_COLLECTION` | Add product to custom collection | write | `collection_id`, `product_id` |  | `read_products`, `write_products` |
| `SHOPIFY_ADJUST_INVENTORY_LEVEL` | Adjust inventory levels | write | `available_adjustment`, `inventory_item_id`, `location_id` |  | `write_inventory` |
| `SHOPIFY_APPLY_FULFILLMENT_HOLD` | Apply fulfillment hold | write | `fulfillment_order_id`, `reason`; *notify_merchant, reason_notes* |  | `write_assigned_fulfillment_orders` |
| `SHOPIFY_APPROVE_COMMENT` | Approve a comment | write | `comment_id` |  | `write_content` |
| `SHOPIFY_BULK_CREATE_PRODUCTS` | Bulk create products via GraphQL | write | `products`; *max_poll_attempts, poll_interval_seconds* |  | `write_products` |
| `SHOPIFY_BULK_DELETE_METAFIELDS` | Delete metafields in bulk | write | `metafields` |  | `write_products`, `write_customers`, `write_orders`, `write_draft_orders`, `write_locations`, `write_inventory`, `write_content`, `write_themes`, `write_markets`, `write_price_rules`, `write_discounts`, `write_gift_cards`, `write_fulfillments`, `write_shipping`, `write_translations`, `write_metaobjects` |
| `SHOPIFY_BULK_DELETE_METAOBJECTS` | Bulk delete metaobjects | write | `where` |  | `write_metaobjects` |
| `SHOPIFY_BULK_QUERY_OPERATION` | Run Shopify bulk query operation | write | `query`; *client_identifier, poll_interval_seconds, timeout_seconds* |  | `read_products`, `read_orders`, `read_customers` |
| `SHOPIFY_BULK_UPDATE_PRODUCT_STATUS` | Bulk update Shopify product status | write | `status`; *collection_id, product_ids* |  | `write_products` |
| `SHOPIFY_CALCULATE_REFUND` | Calculate a refund | read | `order_id`; *currency, refund_line_items, shipping* |  | `read_orders` |
| `SHOPIFY_CANCEL_APP_SUBSCRIPTION` | Cancel App Subscription | write | `id`; *prorate* |  | `write_own_subscription_contracts` |
| `SHOPIFY_CANCEL_BULK_OPERATION` | Cancel Shopify bulk operation | write | `id` |  | `read_products`, `write_products` |
| `SHOPIFY_CANCEL_FULFILLMENT` | Cancel a fulfillment | write | `fulfillment_id` |  | `write_orders` |
| `SHOPIFY_CANCEL_FULFILLMENT_ORDER` | Cancel a Fulfillment Order | write | `fulfillment_order_id` |  | `read_assigned_fulfillment_orders` |
| `SHOPIFY_CANCEL_ORDER` | Cancel an order | write | `order_id`; *amount, currency, email, reason, refund, restock* |  | `write_orders`, `write_marketplace_orders` |
| `SHOPIFY_CANCEL_RECURRING_APPLICATION_CHARGE` | Cancel Recurring Application Charge | write | `recurring_application_charge_id` |  |  |
| `SHOPIFY_CLOSE_FULFILLMENT_ORDER` | Close Fulfillment Order | write | `fulfillment_order_id`; *message* |  | `write_assigned_fulfillment_orders` |
| `SHOPIFY_CLOSE_ORDER` | Close an order | write | `order_id` |  | `write_orders` |
| `SHOPIFY_COMPLETE_DRAFT_ORDER` | Complete a draft order | write | `draft_order_id`; *payment_pending* |  | `write_draft_orders` |
| `SHOPIFY_CONNECT_INVENTORY_LEVEL` | Connect inventory item to location | write | `inventory_item_id`, `location_id`; *relocate_if_necessary* |  | `write_inventory` |
| `SHOPIFY_COUNT_ARTICLES` | Count articles in blog | read | `blog_id`; *created_at_max, created_at_min, published_at_max, published_at_min, published_status, updated_at_max, updated_at_min* |  | `read_content` |
| `SHOPIFY_COUNT_BLOGS` | Get blogs count | read |  |  | `read_content`, `write_content` |
| `SHOPIFY_COUNT_DRAFT_ORDERS` | Get draft orders count | read | *status, updated_at_max, updated_at_min* |  | `read_draft_orders` |
| `SHOPIFY_COUNT_FULFILLMENTS` | Count fulfillments for order | read | `order_id`; *created_at_max, created_at_min, updated_at_max, updated_at_min* |  | `read_orders`, `write_orders` |
| `SHOPIFY_COUNT_METAFIELDS` | Count shop metafields | read |  |  | `read_products` |
| `SHOPIFY_COUNT_PAGES` | Retrieve page count | read | *created_at_max, created_at_min, published_at_max, published_at_min, published_status, title, updated_at_max, updated_at_min* |  | `read_content` |
| `SHOPIFY_COUNT_PRICE_RULES` | Count price rules | read |  |  | `read_price_rules` |
| `SHOPIFY_COUNT_PRODUCTS` | Retrieves a count of products | read | *collection_id, created_at_max, created_at_min, product_type, published_at_max, published_at_min, published_status, status and 3 more* |  | `read_products` |
| `SHOPIFY_COUNT_PRODUCT_IMAGES` | Count product images | read | `product_id` |  | `read_products` |
| `SHOPIFY_COUNT_PRODUCT_LISTINGS` | Count product listings | read |  |  | `read_product_listings`, `write_product_listings` |
| `SHOPIFY_COUNT_PRODUCT_VARIANTS` | Count product variants | read | `product_id` |  | `read_products` |
| `SHOPIFY_CREATES_OR_UPDATES_AN_ASSET_FOR_A_THEME` | Create or Update Theme Asset | write | `key`, `theme_id`; *attachment, source_key, src, value* |  | `write_themes` |
| `SHOPIFY_CREATE_APP_SUBSCRIPTION` | Create app subscription | write | `lineItems`, `name`, `returnUrl`; *replacementBehavior, test, trialDays* |  |  |
| `SHOPIFY_CREATE_ARTICLE` | Create Article | write | `blog_id`, `title`; *author, body_html, image, metafields, published, published_at, summary_html, tags* |  | `write_content` |
| `SHOPIFY_CREATE_ARTICLE_COMMENT` | Create Comment for Article | write | `article_id`, `author`, `blog_id`, `body`, `email`; *ip, user_agent* |  | `write_content` |
| `SHOPIFY_CREATE_A_METAFIELD` | Create a metafield | write | `key`, `namespace`, `resource`, `type`, `value`; *description, resource_id* |  | `write_products` |
| `SHOPIFY_CREATE_BLOG` | Create a new blog | write | `title`; *commentable, handle, metafields, template_suffix* |  | `write_content` |
| `SHOPIFY_CREATE_CARRIER_SERVICE` | Create carrier service | write | `callback_url`, `name`; *active, carrier_service_type, format, service_discovery* |  | `write_shipping` |
| `SHOPIFY_CREATE_CHECKOUT_PAYMENT` | Create checkout payment | write | `amount`, `request_details`, `session_id`, `token`, `unique_token` |  |  |
| `SHOPIFY_CREATE_CUSTOMER` | Create Customer | write | *addresses, email, first_name, last_name, note, password, password_confirmation, phone and 4 more* |  |  |
| `SHOPIFY_CREATE_CUSTOMER_ACCOUNT_ACTIVATION_URL` | Generate Customer Account Activation URL | write | `customer_id` |  | `write_customers` |
| `SHOPIFY_CREATE_CUSTOMER_ADDRESS` | Create customer address | write | `address`, `customer_id` |  | `write_products` |
| `SHOPIFY_CREATE_CUSTOM_COLLECTION` | Create a custom collection | write | `title`; *body_html, collects, handle, image, metafields, published, published_scope, sort_order and 1 more* |  | `read_custom_collections`, `write_custom_collections` |
| `SHOPIFY_CREATE_DELEGATE_ACCESS_TOKEN` | Create delegate access token | write | `delegate_access_scope`; *expires_in* |  |  |
| `SHOPIFY_CREATE_DISCOUNT` | Create Shopify discount code | write | `code`, `starts_at`, `title`, `value_type`; *amount, applies_on_each_item, applies_once_per_customer, collection_ids, combines_with, customer_ids, customer_selection, ends_at and 9 more* |  | `write_discounts` |
| `SHOPIFY_CREATE_DISCOUNT_CODE` | Create Discount Code | write | `code`, `price_rule_id` |  | `write_discounts` |
| `SHOPIFY_CREATE_DISCOUNT_CODE_BATCH` | Create Discount Codes Batch | write | `discount_codes`, `price_rule_id` |  | `write_price_rules` |
| `SHOPIFY_CREATE_DRAFT_ORDER` | Create a new draft order | write | `line_items`; *applied_discount, billing_address, currency, customer_id, email, note, note_attributes, shipping_address and 4 more* |  | `write_shopify_payments_dispute_evidences` |
| `SHOPIFY_CREATE_EVENT_BRIDGE_WEBHOOK_SUBSCRIPTION` | Create EventBridge webhook subscription | write | `topic`, `webhook_subscription` |  | `write_webhooks` |
| `SHOPIFY_CREATE_FULFILLMENT` | Create Fulfillment | write | `line_items_by_fulfillment_order`; *message, notify_customer, origin_address, tracking_info* |  | `read_orders` |
| `SHOPIFY_CREATE_FULFILLMENT_EVENT` | Create Fulfillment Event | write | `fulfillment_id`, `order_id`, `status`; *address1, city, country, happened_at, latitude, longitude, message, province and 1 more* |  | `write_orders` |
| `SHOPIFY_CREATE_FULFILLMENT_SERVICE` | Create Fulfillment Service | write | `callback_url`, `name`; *fulfillment_orders_opt_in, inventory_management, requires_shipping_method, tracking_support* |  | `write_fulfillments` |
| `SHOPIFY_CREATE_GIFT_CARD` | Create Gift Card | write | `gift_card` |  | `write_gift_cards` |
| `SHOPIFY_CREATE_MARKETING_ENGAGEMENTS` | Create Marketing Event Engagements | write | `engagements`, `marketing_event_id` |  | `write_marketing_events` |
| `SHOPIFY_CREATE_MARKETING_EVENT` | Create Marketing Event | write | `event_type`, `marketing_channel`, `started_at`; *budget, budget_type, currency, description, ended_at, manage_url, marketed_resources, paid and 7 more* |  | `read_locations` |
| `SHOPIFY_CREATE_METAFIELD` | Create blog metafield | write | `blog_id`, `key`, `namespace`, `type`, `value`; *description* |  | `read_products` |
| `SHOPIFY_CREATE_MOBILE_PLATFORM_APPLICATION` | Create Mobile Platform Application | write | `application_id`, `platform`; *app_clip_application_id, enabled_app_clips, enabled_shared_webcredentials, enabled_universal_or_app_links, sha256_cert_fingerprints* |  | `mobile_platform_applications` |
| `SHOPIFY_CREATE_ONE_TIME_APPLICATION_CHARGE` | Create One-Time Application Charge | write | `name`, `price`, `return_url`; *test* |  |  |
| `SHOPIFY_CREATE_ORDER` | Create an order | write | `line_items`; *billing_address, currency, customer, discount_codes, email, financial_status, fulfillment_status, inventory_behaviour and 9 more* |  | `read_orders`, `write_orders` |
| `SHOPIFY_CREATE_ORDER_TRANSACTION` | Create transaction for order | write | `kind`, `order_id`; *amount, authorization, currency, gateway, parent_id, source_name, test* |  | `write_orders`, `marketplace_orders` |
| `SHOPIFY_CREATE_PAGE` | Create a new page | write | `title`; *author, body_html, handle, metafields, published, published_at, template_suffix* |  | `write_content`, `read_content` |
| `SHOPIFY_CREATE_PRICE_RULE` | Create Price Rule | write | `allocation_method`, `customer_selection`, `starts_at`, `target_selection`, `target_type`, `title`, `value`, `value_type`; *allocation_limit, customer_segment_prerequisite_ids, ends_at, entitled_collection_ids, entitled_country_ids, entitled_product_ids, entitled_variant_ids, once_per_customer and 9 more* |  | `write_price_rules` |
| `SHOPIFY_CREATE_PRODUCT_METAFIELD` | Create product metafield | write | `key`, `namespace`, `product_id`, `value`; *description, type* |  | `write_products` |
| `SHOPIFY_CREATE_PRODUCT_RESOURCE_FEEDBACK` | Create Product Resource Feedback | write | `feedback_generated_at`, `product_id`, `resource_updated_at`, `state`; *messages* |  | `write_resource_feedbacks` |
| `SHOPIFY_CREATE_PRODUCT_VARIANT` | Create a new product variant | write | `product_id`, `variant` |  | `write_product_listings` |
| `SHOPIFY_CREATE_PRODUCT_WITH_VARIANTS` | Create a Shopify product with variants | write | *description_html, handle, images, product_options, product_type, status, tags, title and 2 more* |  | `write_products` |
| `SHOPIFY_CREATE_PUB_SUB_WEBHOOK_SUBSCRIPTION` | Create Pub/Sub webhook subscription | write | `topic`, `webhook_subscription` |  | `write_webhooks` |
| `SHOPIFY_CREATE_RECURRING_APPLICATION_CHARGE` | Create recurring application charge | write | `name`, `price`, `return_url`; *capped_amount, terms, test, trial_days* |  |  |
| `SHOPIFY_CREATE_REDIRECT` | Create Redirect | write | `path`, `target` |  | `write_online_store_navigation` |
| `SHOPIFY_CREATE_REFUND` | Create order refund | write | `order_id`; *currency, discrepancy_reason, note, notify, refund_line_items, shipping, transactions* |  | `write_orders` |
| `SHOPIFY_CREATE_RESOURCE_FEEDBACK` | Create Resource Feedback | write | `feedback_generated_at`, `state`; *messages* |  | `write_resource_feedbacks` |
| `SHOPIFY_CREATE_SCRIPT_TAG` | Create Script Tag | write | `event`, `src`; *cache, display_scope* |  | `write_script_tags` |
| `SHOPIFY_CREATE_SMART_COLLECTIONS` | Create Smart Collection | write | `rules`, `title`; *body_html, disjunctive, image, published, sort_order, template_suffix* |  | `write_products` |
| `SHOPIFY_CREATE_STOREFRONT_ACCESS_TOKEN` | Create Storefront Access Token | write | `title` |  | `unauthenticated_read_product_listings`, `unauthenticated_write_checkouts`, `unauthenticated_read_checkouts`, `unauthenticated_read_customers`, `unauthenticated_write_customers`, `unauthenticated_read_customer_tags`, `unauthenticated_read_content`, `unauthenticated_read_metaobjects`, `unauthenticated_read_product_inventory`, `unauthenticated_read_product_pickup_locations`, `unauthenticated_read_product_tags`, `unauthenticated_read_selling_plans` |
| `SHOPIFY_CREATE_THEME` | Create Theme | write | `name`; *role, src* |  | `write_themes` |
| `SHOPIFY_CREATE_USAGE_CHARGE` | Create Usage Charge | write | `description`, `price`, `recurring_application_charge_id` |  |  |
| `SHOPIFY_CREATE_WEBHOOK_SUBSCRIPTION` | Create webhook subscription | write | `topic`, `webhookSubscription` |  | `write_webhooks` |
| `SHOPIFY_CUSTOMIZE_RECURRING_APPLICATION_CHARGE` | Customize Recurring Application Charge | write | `capped_amount`, `recurring_application_charge_id` |  |  |
| `SHOPIFY_DELETES_A_CUSTOM_COLLECTION` | Delete a custom collection | write | `custom_collection_id` |  | `write_products` |
| `SHOPIFY_DELETES_A_DISCOUNT_CODE` | Delete a discount code | write | `discount_code_id`, `price_rule_id` |  | `write_orders` |
| `SHOPIFY_DELETE_ARTICLE` | Delete Article | write | `article_id`, `blog_id` |  | `write_content` |
| `SHOPIFY_DELETE_BLOG` | Delete a blog | write | `blog_id` |  | `write_content` |
| `SHOPIFY_DELETE_CARRIER_SERVICE` | Delete Carrier Service | write | `carrier_service_id` |  | `write_shipping` |
| `SHOPIFY_DELETE_COLLECT` | Remove product from collection | write | `collect_id` |  | `write_products` |
| `SHOPIFY_DELETE_COMMENT` | Remove a comment | write | `comment_id` |  | `write_content` |
| `SHOPIFY_DELETE_CUSTOMER` | Delete a customer | write | `customer_id` |  | `write_customers` |
| `SHOPIFY_DELETE_CUSTOMER_ADDRESS` | Delete customer address | write | `address_id`, `customer_id` |  | `write_products` |
| `SHOPIFY_DELETE_DISPUTE_EVIDENCE_FILE` | Delete dispute evidence file | write | `dispute_file_upload_id`, `dispute_id` |  | `write_shopify_payments_dispute_file_uploads` |
| `SHOPIFY_DELETE_DRAFT_ORDER` | Delete a draft order | write | `draft_order_id` |  | `write_draft_orders` |
| `SHOPIFY_DELETE_FULFILLMENT_SERVICE` | Remove fulfillment service | write | `fulfillment_service_id` |  | `write_fulfillments` |
| `SHOPIFY_DELETE_INVENTORY_LEVEL` | Delete inventory levels | write | `inventory_item_id`, `location_id` |  | `write_inventory` |
| `SHOPIFY_DELETE_MARKETING_EVENT` | Delete a marketing event | write | `marketing_event_id` |  | `write_marketing_events` |
| `SHOPIFY_DELETE_METAFIELD` | Delete a metafield by its ID | write | `metafield_id`, `resource`, `resource_id` |  | `write_products` |
| `SHOPIFY_DELETE_METAOBJECT` | Delete metaobject | write | `id` |  | `write_metaobjects` |
| `SHOPIFY_DELETE_METAOBJECT_DEFINITION` | Delete metaobject definition | write | `id` |  | `write_metaobject_definitions` |
| `SHOPIFY_DELETE_MOBILE_PLATFORM_APPLICATION` | Delete Mobile Platform Application | write | `mobile_platform_application_id` |  | `mobile_platform_applications` |
| `SHOPIFY_DELETE_ORDER` | Delete an order | write | `order_id` |  | `write_orders` |
| `SHOPIFY_DELETE_PAGE` | Delete a page | write | `page_id` |  | `write_content` |
| `SHOPIFY_DELETE_PRICE_RULE` | Delete price rule | write | `price_rule_id` |  | `write_price_rules` |
| `SHOPIFY_DELETE_PRODUCT` | Delete a product | write | `product_id` |  | `write_products` |
| `SHOPIFY_DELETE_PRODUCT_IMAGE` | Delete product image | write | `image_id`, `product_id` |  | `write_products` |
| `SHOPIFY_DELETE_PRODUCT_METAFIELD` | Delete product metafield | write | `metafield_id`, `product_id` |  | `write_products` |
| `SHOPIFY_DELETE_PRODUCT_VARIANT` | Remove product variant | write | `product_id`, `variant_id` |  | `write_products` |
| `SHOPIFY_DELETE_REDIRECT` | Delete redirect | write | `redirect_id` |  | `write_content` |
| `SHOPIFY_DELETE_SAVED_SEARCH` | Delete a saved search | write | `id` |  | `write_customers`, `write_products`, `write_orders`, `write_draft_orders`, `write_discounts`, `write_files`, `write_online_store_navigation` |
| `SHOPIFY_DELETE_SMART_COLLECTION` | Delete smart collection | write | `smart_collection_id` |  | `write_products` |
| `SHOPIFY_DELETE_STOREFRONT_ACCESS_TOKEN` | Delete Storefront Access Token | write | `storefront_access_token_id` |  | `unauthenticated_read_product_listings`, `unauthenticated_write_checkouts`, `unauthenticated_read_checkouts`, `unauthenticated_read_customers`, `unauthenticated_write_customers`, `unauthenticated_read_customer_tags`, `unauthenticated_read_content`, `unauthenticated_read_product_tags`, `unauthenticated_read_product_inventory`, `unauthenticated_read_product_pickup_locations`, `unauthenticated_read_selling_plans`, `unauthenticated_read_metaobjects` |
| `SHOPIFY_DELETE_THEME` | Delete theme | write | `theme_id` |  | `write_themes` |
| `SHOPIFY_DELETE_THEME_ASSET` | Delete theme asset | write | `asset_key`, `theme_id` |  | `write_themes`, `read_themes` |
| `SHOPIFY_DELETE_WEBHOOK_SUBSCRIPTION` | Delete webhook subscription | write | `id` |  | `write_webhooks` |
| `SHOPIFY_DELETE_WEB_PRESENCE` | Delete web presence | write | `id` |  | `read_markets`, `write_markets` |
| `SHOPIFY_DISABLE_GIFT_CARD` | Disable a Gift Card | write | `gift_card_id` |  | `write_gift_cards` |
| `SHOPIFY_ENABLE_STANDARD_METAOBJECT_DEFINITION` | Enable standard metaobject definition | write | `type` |  | `write_metaobject_definitions`, `read_metaobject_definitions` |
| `SHOPIFY_GET_ABANDONMENT` | Get abandonment | read | `id` |  | `read_orders` |
| `SHOPIFY_GET_ABANDONMENT_BY_ABANDONED_CHECKOUT_ID` | Get abandonment by abandoned checkout ID | read | `abandoned_checkout_id` |  | `read_orders`, `write_orders` |
| `SHOPIFY_GET_ACCESS_SCOPES` | Get a list of access scopes | read |  |  |  |
| `SHOPIFY_GET_APP` | Get App | read | *id* |  |  |
| `SHOPIFY_GET_APPLICATION_CHARGES` | Get Application Charges | read | *fields, since_id* |  |  |
| `SHOPIFY_GET_APPLICATION_CHARGE_BY_ID` | Get Application Charge By ID | read | `id`; *fields* |  |  |
| `SHOPIFY_GET_APPLICATION_CREDIT` | Get application credit | read | `application_credit_id`; *fields* |  |  |
| `SHOPIFY_GET_APPLICATION_CREDITS` | Get application credits | read | *fields* |  |  |
| `SHOPIFY_GET_APP_BY_HANDLE` | Get app by handle | read | `handle` |  | `read_apps` |
| `SHOPIFY_GET_APP_INSTALLATION` | Get app installation details | read | *id* |  | `read_apps` |
| `SHOPIFY_GET_ARTICLE` | Get Article | read | `article_id`, `blog_id`; *fields* |  | `read_content` |
| `SHOPIFY_GET_BATCH_DISCOUNT_CODES` | Get Batch Discount Codes | read | `batch_id`, `price_rule_id` |  | `read_price_rules` |
| `SHOPIFY_GET_BLOG` | Get Single Blog | read | `blog_id`; *fields* |  | `read_content` |
| `SHOPIFY_GET_BLOG_ARTICLE_TAGS` | Get blog article tags | read | `blog_id`; *limit, popular* |  | `read_content` |
| `SHOPIFY_GET_BULK_OPERATION` | Get bulk operation status | read | `id` |  | `read_products`, `read_orders`, `read_customers` |
| `SHOPIFY_GET_CARRIER_SERVICE` | Get carrier service | read | `carrier_service_id` |  | `read_shipping`, `write_shipping` |
| `SHOPIFY_GET_CARRIER_SERVICES` | Get carrier services | read |  |  | `read_shipping`, `write_shipping` |
| `SHOPIFY_GET_CATALOGS_COUNT` | Get catalogs count | read | *query, type* |  | `read_products`, `read_product_listings` |
| `SHOPIFY_GET_CHECKOUTS_COUNT` | Get checkouts count | read | *created_at_max, created_at_min, since_id, status, updated_at_max, updated_at_min* |  | `read_checkouts` |
| `SHOPIFY_GET_CHECKOUT_PAYMENT` | Get checkout payment | read | `payment_id`, `token` |  | `read_checkouts`, `write_checkouts` |
| `SHOPIFY_GET_CHECKOUT_PROFILE` | Get checkout profile | read | `id` |  | `read_checkouts` |
| `SHOPIFY_GET_COLLECTION_BY_ID` | Get collection by ID | read | `collection_id`; *fields* |  | `read_checkouts` |
| `SHOPIFY_GET_COLLECTION_BY_IDENTIFIER` | Get collection by identifier | read | `identifier` |  | `read_products` |
| `SHOPIFY_GET_COLLECTION_LISTING` | Get collection listing | read | `collection_listing_id` |  | `read_product_listings`, `write_product_listings` |
| `SHOPIFY_GET_COLLECTS` | Get collects | read | *collection_id, fields, limit, product_id, since_id* |  | `read_products`, `write_products` |
| `SHOPIFY_GET_COLLECTS_COUNT` | Get collects count | read | *collection_id, product_id* |  | `read_products` |
| `SHOPIFY_GET_COLLECT_BY_ID` | Get collect by ID | read | `collect_id`; *fields* |  | `read_products` |
| `SHOPIFY_GET_COMMENT` | Get Comment | read | `comment_id`; *fields* |  | `read_content`, `write_content` |
| `SHOPIFY_GET_COMMENTS_COUNT` | Get count of comments | read | *article_id, blog_id, created_at_max, created_at_min, published_at_max, published_at_min, published_status, status and 2 more* |  | `read_content` |
| `SHOPIFY_GET_CONSENT_POLICY` | Get consent policy | read | *consent_required, country_code, data_sale_opt_out_required, id, region_code* |  | `read_privacy_settings` |
| `SHOPIFY_GET_COUNTRIES` | Get Countries | read | *fields, since_id* |  | `read_shipping` |
| `SHOPIFY_GET_COUNTRIES_PARAM_COUNTRY_ID_PROVINCES_COUNT` | Get provinces count for country | read | `country_id` |  | `read_shipping` |
| `SHOPIFY_GET_COUNTRIES_PARAM_COUNTRY_ID_PROVINCES_PARAM_PROVI` | Get Province for Country | read | `country_id`, `province_id`; *fields* |  | `read_shipping` |
| `SHOPIFY_GET_COUNTRY` | Get country | read | `country_id`; *fields* |  | `read_shipping` |
| `SHOPIFY_GET_COUNTRY_PROVINCES` | Get provinces for country | read | `country_id`; *fields, since_id* |  | `read_shipping` |
| `SHOPIFY_GET_CURRENCIES` | Get Currencies | read |  |  | `read_shop` |
| `SHOPIFY_GET_CURRENT_USER` | Get current user | read |  |  | `read_users` |
| `SHOPIFY_GET_CUSTOMER` | Retrieve a single customer | read | `customer_id`; *fields* |  | `read_customers` |
| `SHOPIFY_GET_CUSTOMERS_COUNT` | Get customers count | read | *created_at_max, created_at_min, updated_at_max, updated_at_min* |  | `read_customers` |
| `SHOPIFY_GET_CUSTOMERS_SEARCH` | Search Customers | read | `query`; *fields, limit, order* |  | `read_customers` |
| `SHOPIFY_GET_CUSTOMER_ADDRESS` | Get Customer Address | read | `address_id`, `customer_id` |  | `read_shop` |
| `SHOPIFY_GET_CUSTOMER_ADDRESSES` | Get Customer Addresses | read | `customer_id`; *limit* |  | `read_customers` |
| `SHOPIFY_GET_CUSTOMER_ORDERS` | Get customer orders | read | `customer_id`; *limit, page_info, since_id, status* |  | `read_orders` |
| `SHOPIFY_GET_CUSTOM_COLLECTIONS` | Get custom collections | write | *fields, handle, ids, limit, product_id, published_at_max, published_at_min, published_status and 4 more* |  | `read_products` |
| `SHOPIFY_GET_CUSTOM_COLLECTIONS_COUNT` | Get custom collections count | write | *product_id, published_at_max, published_at_min, published_status, title, updated_at_max, updated_at_min* |  |  |
| `SHOPIFY_GET_CUSTOM_COLLECTION_BY_ID` | Get custom collection by ID | write | `custom_collection_id`; *fields* |  | `read_products` |
| `SHOPIFY_GET_DISCOUNT_CODE` | Get Discount Code | read | `discount_code_id`, `price_rule_id` |  | `read_price_rules` |
| `SHOPIFY_GET_DISCOUNT_CODES_COUNT` | Get discount codes count | read | *times_used, times_used_max, times_used_min* |  | `read_price_rules` |
| `SHOPIFY_GET_DISCOUNT_CODES_LOOKUP` | Get Discount Code Location | write | `code` |  | `read_price_rules` |
| `SHOPIFY_GET_DISCOUNT_CODE_BATCH_JOB` | Retrieve Discount Code Creation Job | read | `batch_id`, `price_rule_id` |  | `read_products` |
| `SHOPIFY_GET_DISPUTE_EVIDENCE` | Get dispute evidence | read | `dispute_id` |  | `read_shopify_payments_dispute_evidences`, `write_shopify_payments_dispute_evidences` |
| `SHOPIFY_GET_DRAFT_ORDER` | Receive a single draft order | read | `draft_order_id`; *fields* |  | `read_draft_orders` |
| `SHOPIFY_GET_EVENT` | Get event by ID | read | `event_id`; *fields* |  |  |
| `SHOPIFY_GET_FULFILLMENT` | Get fulfillments by fulfillment order | read | `fulfillment_order_id` |  | `read_assigned_fulfillment_orders`, `read_merchant_managed_fulfillment_orders`, `read_third_party_fulfillment_orders` |
| `SHOPIFY_GET_FULFILLMENT_EVENTS` | Get fulfillment events | read | `fulfillment_id`, `order_id` |  |  |
| `SHOPIFY_GET_FULFILLMENT_ORDERS_FOR_ORDER` | Get fulfillment orders for order | read | `order_id`; *include_financial_summaries, include_order_reference_fields* |  | `write_assigned_fulfillment_orders` |
| `SHOPIFY_GET_FULFILLMENT_ORDER_LOCATIONS_FOR_MOVE` | Get fulfillment order locations for move | read | `fulfillment_order_id` |  | `read_merchant_managed_fulfillment_orders` |
| `SHOPIFY_GET_FULFILLMENT_SERVICE` | Get Fulfillment Service | read | `fulfillment_service_id` |  | `read_fulfillments`, `write_fulfillments` |
| `SHOPIFY_GET_FULFILLMENT_SERVICES` | Get Fulfillment Services | read | *scope* |  | `write_fulfillments` |
| `SHOPIFY_GET_GIFT_CARDS` | Retrieves a list of gift cards | read | *fields, limit, since_id, status* |  | `read_gift_cards`, `write_gift_cards` |
| `SHOPIFY_GET_GIFT_CARDS_COUNT` | Retrieves a count of gift cards | read | *status* |  | `read_gift_cards` |
| `SHOPIFY_GET_GIFT_CARD_BY_ID` | Get Gift Card By ID | read | `gift_card_id` |  | `read_gift_cards` |
| `SHOPIFY_GET_INVENTORY_ITEM` | Get Inventory Item | read | `inventory_item_id` |  | `read_inventory` |
| `SHOPIFY_GET_INVENTORY_ITEMS` | Get Inventory Items | read | `ids`; *limit* |  | `read_inventory` |
| `SHOPIFY_GET_INVENTORY_LEVELS` | Get Inventory Levels | read | *inventory_item_ids, limit, location_ids, updated_at_min* |  | `read_inventory` |
| `SHOPIFY_GET_INVENTORY_LEVELS_FOR_LOCATION` | Get Inventory Levels For Location | read | `location_id`; *inventory_item_ids, limit, page_info, updated_at_min* |  | `read_inventory` |
| `SHOPIFY_GET_LOCATION` | Get location | read | `location_id` |  | `read_locations` |
| `SHOPIFY_GET_LOCATIONS_COUNT` | Get locations count | read |  |  | `read_locations`, `read_inventory`, `read_markets_home` |
| `SHOPIFY_GET_MARKETING_EVENT` | Retrieve a single marketing event | read | `marketing_event_id` |  | `read_marketing_events` |
| `SHOPIFY_GET_MARKETING_EVENTS_COUNT` | Get marketing events count | read |  |  | `read_marketing_events` |
| `SHOPIFY_GET_METAFIELD` | Get blog metafields | read | `blog_id`; *created_at_max, created_at_min, fields, key, limit, metafield_type, namespace, since_id and 2 more* |  | `read_content` |
| `SHOPIFY_GET_METAFIELDS` | Get metafields | read | *created_at_max, created_at_min, fields, key, limit, metafield_owner_id, metafield_owner_resource, namespace and 3 more* |  | `read_products` |
| `SHOPIFY_GET_METAFIELD_BY_ID` | Get metafield by ID | read | `metafield_id`; *fields* |  | `read_products` |
| `SHOPIFY_GET_MOBILE_PLATFORM_APPLICATION` | Get mobile platform application | read | `mobile_platform_application_id` |  | `mobile_platform_applications` |
| `SHOPIFY_GET_ORDER` | Get order | read | `order_id`; *fields* |  | `read_orders`, `read_all_orders` |
| `SHOPIFY_GET_ORDERS_WITH_FILTERS` | Get orders with filters | read | *attribution_app_id, created_at_max, created_at_min, fields, financial_status, fulfillment_status, ids, limit and 8 more* |  | `read_orders` |
| `SHOPIFY_GET_ORDER_FULFILLMENTS` | Retrieve order fulfillments | read | `order_id`; *created_at_max, created_at_min, fields, limit, since_id, updated_at_max, updated_at_min* |  | `write_orders` |
| `SHOPIFY_GET_ORDER_REFUNDS` | Get refunds for order | read | `order_id`; *fields, in_shop_currency, limit* |  | `read_orders` |
| `SHOPIFY_GET_ORDER_REFUND_BY_ID` | Get specific refund for order | read | `order_id`, `refund_id`; *fields, in_shop_currency* |  | `read_orders` |
| `SHOPIFY_GET_ORDER_TRANSACTIONS_COUNT` | Get order transactions count | read | `order_id` |  | `read_orders` |
| `SHOPIFY_GET_PAGE_BY_ID` | Get page by ID | read | `page_id`; *fields* |  | `read_content` |
| `SHOPIFY_GET_PAYMENTS_BALANCE` | Get payments balance | read |  |  | `read_shopify_payments_payouts`, `read_shopify_payments_accounts` |
| `SHOPIFY_GET_PAYMENTS_PAYOUT` | Get payments payout | read | `payout_id` |  | `read_shopify_payments_payouts` |
| `SHOPIFY_GET_PAYMENT_DISPUTE` | Get payment dispute | read | `dispute_id` |  | `read_shopify_payments_disputes` |
| `SHOPIFY_GET_POLICIES` | Get Policies | read |  |  | `read_orders` |
| `SHOPIFY_GET_PRICE_RULE_BATCH_INFO` | Check Batch Listing Support for Price Rule | read | `price_rule_id` |  | `read_price_rules` |
| `SHOPIFY_GET_PRICE_RULE_DISCOUNT_CODES` | Get Discount Codes for Price Rule | read | `price_rule_id` |  | `read_price_rules` |
| `SHOPIFY_GET_PRODUCT` | Get product | read | `product_id` |  | `read_products` |
| `SHOPIFY_GET_PRODUCTS_PAGINATED` | Get products (paginated) | read | *fields, limit, page_info* |  | `read_products` |
| `SHOPIFY_GET_PRODUCT_IMAGE` | Get product image | read | `image_id`, `product_id` |  | `read_products`, `write_products` |
| `SHOPIFY_GET_PRODUCT_LISTING` | Get product listing | read | `product_listing_id` |  | `read_product_listings`, `write_product_listings` |
| `SHOPIFY_GET_PRODUCT_METAFIELD` | Get product metafield | read | `metafield_id`, `product_id` |  | `read_products` |
| `SHOPIFY_GET_PRODUCT_METAFIELDS` | Get product metafields | read | `product_id`; *key, namespace* |  | `read_products` |
| `SHOPIFY_GET_PRODUCT_VARIANT` | Receive a single product variant | read | `variant_id`; *fields* |  | `read_products` |
| `SHOPIFY_GET_RECURRING_APPLICATION_CHARGE` | Get recurring application charge | read | `recurring_application_charge_id`; *fields* |  |  |
| `SHOPIFY_GET_RECURRING_APPLICATION_CHARGES` | Get recurring application charges | read | *fields, since_id* |  |  |
| `SHOPIFY_GET_REDIRECTS` | Get redirects | read | *fields, limit, path, since_id, target* |  | `read_content` |
| `SHOPIFY_GET_REDIRECTS_COUNT` | Get redirects count | read | *path, target* |  | `read_content` |
| `SHOPIFY_GET_REDIRECT_BY_ID` | Get redirect by ID | read | `redirect_id`; *fields* |  | `read_content` |
| `SHOPIFY_GET_SCRIPTTAG` | Get Script Tag | read | `script_tag_id`; *fields* |  | `read_script_tags` |
| `SHOPIFY_GET_SEGMENTS_COUNT` | Get segments count | read |  |  | `read_customers` |
| `SHOPIFY_GET_SHIPPING_ZONES` | Get Shipping Zones | read | *fields* |  |  |
| `SHOPIFY_GET_SHOP_BILLING_ADDRESS` | Get shop billing address | read |  |  | `read_shop_data` |
| `SHOPIFY_GET_SHOP_INFO` | Get Shopify shop information | read |  |  |  |
| `SHOPIFY_GET_SMART_COLLECTIONS` | Get smart collections | read | *fields, handle, ids, limit, product_id, published_at_max, published_at_min, published_status and 4 more* |  | `read_products` |
| `SHOPIFY_GET_SMART_COLLECTIONS_COUNT` | Get smart collections count | read | *product_id, published_at_max, published_at_min, published_status, title, updated_at_max, updated_at_min* |  |  |
| `SHOPIFY_GET_SMART_COLLECTION_BY_ID` | Get smart collection by ID | read | `smart_collection_id`; *fields* |  | `read_products` |
| `SHOPIFY_GET_STOREFRONT_ACCESS_TOKENS` | Get storefront access tokens | read |  |  | `unauthenticated_read_checkouts` |
| `SHOPIFY_GET_TENDER_TRANSACTIONS` | Get tender transactions | read | *limit, order, processed_at, processed_at_max, processed_at_min, since_id* |  | `unauthenticated_read_checkouts` |
| `SHOPIFY_GET_THEME` | Retrieve a single theme | read | `theme_id`; *fields* |  | `read_themes` |
| `SHOPIFY_GET_THEMES` | Get themes | read | *fields* |  | `read_themes` |
| `SHOPIFY_GET_THEMES_PARAM_THEME_ID_ASSETS` | Get Theme Asset | read | `asset_key`, `theme_id`; *fields* |  | `read_themes` |
| `SHOPIFY_GET_TRANSACTION` | Retrieve a specific transaction | read | `order_id`, `transaction_id`; *fields, in_shop_currency* |  | `read_orders` |
| `SHOPIFY_GET_USAGE_CHARGE` | Get usage charge | read | `recurring_application_charge_id`, `usage_charge_id`; *fields* |  |  |
| `SHOPIFY_GET_USER` | Get user | read | `user_id` |  | `read_users` |
| `SHOPIFY_GET_WEBHOOK` | Get Webhook | read | `webhook_id`; *fields* |  |  |
| `SHOPIFY_GET_WEBHOOK_SUBSCRIPTIONS_COUNT` | Get webhook subscriptions count | read | *query* |  | `read_webhooks` |
| `SHOPIFY_GRAPH_QL_ADMIN_EXECUTE` | Shopify GraphQL Admin execute | write | `document`; *api_version, operation_name, raise_on_graphql_errors, raise_on_user_errors, variables* |  |  |
| `SHOPIFY_GRAPH_QL_COLLECTIONS` | Shopify GraphQL collections and publication-ready collection reads | write | `operation`; *after, api_version, before, custom_document, fields, first, last, operation_name and 6 more* |  | `read_products`, `write_products` |
| `SHOPIFY_GRAPH_QL_CUSTOMERS` | Shopify GraphQL customers and customer search | write | `operation`; *after, api_version, before, custom_document, fields, first, last, operation_name and 6 more* |  | `read_customers`, `write_customers` |
| `SHOPIFY_GRAPH_QL_DISCOUNTS` | Shopify GraphQL discounts and discount search | write | `operation`; *after, api_version, before, custom_document, fields, first, last, operation_name and 6 more* |  | `read_discounts`, `write_discounts` |
| `SHOPIFY_GRAPH_QL_DRAFT_ORDERS` | Shopify GraphQL draft orders | write | `operation`; *after, api_version, before, custom_document, fields, first, last, operation_name and 6 more* |  | `read_draft_orders`, `write_draft_orders` |
| `SHOPIFY_GRAPH_QL_INVENTORY` | Shopify GraphQL inventory items and levels | write | `operation`; *after, api_version, before, custom_document, fields, first, last, operation_name and 6 more* |  | `read_inventory`, `write_inventory`, `read_products` |
| `SHOPIFY_GRAPH_QL_METAFIELDS_METAOBJECTS` | Shopify GraphQL metafields and metaobjects | write | `operation`; *after, api_version, custom_document, fields, first, operation_name, raise_on_graphql_errors, variables* |  | `read_products`, `read_customers`, `read_orders`, `read_metaobjects`, `write_metaobjects` |
| `SHOPIFY_GRAPH_QL_MUTATION` | Execute safe Shopify GraphQL mutation | write | `mutation`; *operation_name, variables* |  |  |
| `SHOPIFY_GRAPH_QL_ORDERS` | Shopify GraphQL orders and order search | write | `operation`; *after, api_version, before, custom_document, fields, first, last, operation_name and 6 more* |  | `read_orders`, `read_all_orders`, `write_orders` |
| `SHOPIFY_GRAPH_QL_PRODUCTS` | Shopify GraphQL products, variants, media, and publishing | write | `operation`; *after, api_version, before, fields, first, last, query, raise_on_graphql_errors and 4 more* |  | `read_products`, `write_products` |
| `SHOPIFY_GRAPH_QL_QUERY` | Execute Shopify GraphQL query | write | `query`; *variables* |  | `read_products`, `write_products` |
| `SHOPIFY_GRAPH_QL_READ_QUERY` | Execute safe read-only Shopify GraphQL query | write | `query`; *operation_name, variables* |  |  |
| `SHOPIFY_GRAPH_QL_SCHEMA` | Inspect Shopify GraphQL schema | write | `operation`; *field_name, search, type_name* |  |  |
| `SHOPIFY_GRAPH_QL_SHOP` | Shopify GraphQL shop and app metadata | write | *api_version, custom_document, fields, raise_on_graphql_errors, variables* |  | `read_shop`, `read_products`, `read_orders` |
| `SHOPIFY_GRAPH_QL_WEBHOOKS` | Shopify GraphQL webhook subscriptions | write | `operation`; *after, api_version, before, custom_document, fields, first, last, operation_name and 6 more* |  | `read_webhooks`, `write_webhooks` |
| `SHOPIFY_GRAPH_QL_WRITE_OPERATIONS` | Shopify GraphQL write operations | write | `operation`; *api_version, custom_document, operation_name, raise_on_graphql_errors, raise_on_user_errors, variables* |  | `write_products`, `write_orders`, `write_customers`, `write_inventory`, `write_discounts`, `write_draft_orders`, `write_webhooks`, `write_metaobjects` |
| `SHOPIFY_LIST_ARTICLE_AUTHORS` | List article authors | read |  |  | `read_checkouts` |
| `SHOPIFY_LIST_ARTICLE_TAGS` | List article tags | read | *limit, popular* |  | `read_content` |
| `SHOPIFY_LIST_ASSIGNED_FULFILLMENT_ORDERS` | List assigned fulfillment orders | read | *assignment_status, location_ids* |  | `read_assigned_fulfillment_orders`, `write_assigned_fulfillment_orders` |
| `SHOPIFY_LIST_AVAILABLE_LOCALES` | List available locales | read |  |  | `read_locales` |
| `SHOPIFY_LIST_BLOGS` | Retrieve a list of all blogs | read | *fields, handle, limit, since_id* |  | `read_content` |
| `SHOPIFY_LIST_BLOG_ARTICLES` | Retrieve list of articles from blog | read | `blog_id`; *author, created_at_max, created_at_min, fields, handle, limit, published_at_max, published_at_min and 2 more* |  | `read_content` |
| `SHOPIFY_LIST_BULK_OPERATIONS` | List bulk operations | read | *after, before, first, last, query, reverse, sort_key* |  | `read_bulk_operations` |
| `SHOPIFY_LIST_CHANNELS` | List channels | read | *after, before, first, last, reverse* |  | `read_publications` |
| `SHOPIFY_LIST_CHECKOUT_PROFILES` | List checkout profiles | read | *after, before, first, last, query, reverse, sortKey* |  | `read_checkout_branding_settings` |
| `SHOPIFY_LIST_COLLECTION_LISTINGS` | List Collection Listings | read | *limit, page_info* |  | `read_product_listings`, `write_product_listings` |
| `SHOPIFY_LIST_COLLECTION_LISTING_PRODUCT_IDS` | List collection listing product IDs | read | `collection_listing_id`; *limit, page_info* |  | `read_product_listings`, `write_product_listings` |
| `SHOPIFY_LIST_COLLECTION_PRODUCTS` | Retrieve products in collection | read | `collection_id`; *fields, limit, page_info, since_id* |  | `read_products` |
| `SHOPIFY_LIST_COMMENTS` | Retrieve a list of comments | read | *created_at_max, created_at_min, fields, limit, published_at_max, published_at_min, published_status, since_id and 3 more* |  | `read_content` |
| `SHOPIFY_LIST_CUSTOMERS` | Retrieves a list of customers | read | *created_at_max, created_at_min, fields, ids, limit, page_info, since_id, status and 2 more* |  | `write_customers` |
| `SHOPIFY_LIST_DEPRECATED_API_CALLS` | List Deprecated API Calls | read |  |  |  |
| `SHOPIFY_LIST_DRAFT_ORDERS` | Retrieve list of draft orders | read | *fields, ids, limit, since_id, status, updated_at_max, updated_at_min* |  | `read_draft_orders`, `write_draft_orders` |
| `SHOPIFY_LIST_LOCATIONS` | List locations | read | *limit* |  | `read_locations` |
| `SHOPIFY_LIST_LOCATION_INVENTORY_LEVELS` | List Location Inventory Levels | read | `location_id`; *limit, page_info, since_id* |  | `read_inventory`, `write_inventory` |
| `SHOPIFY_LIST_MARKETING_EVENTS` | Retrieve all marketing events | read | *limit* |  | `read_marketing_events` |
| `SHOPIFY_LIST_MOBILE_PLATFORM_APPLICATIONS` | List Mobile Platform Applications | read |  |  | `mobile_platform_applications` |
| `SHOPIFY_LIST_PAGES` | Retrieve list of pages | read | *created_at_max, created_at_min, fields, handle, limit, published_at_max, published_at_min, published_status and 4 more* |  | `read_content` |
| `SHOPIFY_LIST_PAYMENTS_BALANCE_TRANSACTIONS` | List Payments Balance Transactions | read | *last_id, limit, page_info, payout_id, payout_status, since_id, test* |  | `read_shopify_payments_payouts` |
| `SHOPIFY_LIST_PAYMENTS_PAYOUTS` | List Payments Payouts | read | *date, date_max, date_min, last_id, limit, page_info, since_id, status* |  | `read_shopify_payments_payouts` |
| `SHOPIFY_LIST_PAYMENT_DISPUTES` | List Payment Disputes | read | *initiated_at, last_id, limit, page_info, since_id, status* |  | `read_shopify_payments_disputes` |
| `SHOPIFY_LIST_PAYMENT_TERMS_TEMPLATES` | List Payment Terms Templates | read | *payment_terms_type* |  | `read_payment_terms` |
| `SHOPIFY_LIST_PRODUCT_IMAGES` | Receive A List Of All Product Images | read | `product_id`; *fields, limit, since_id* |  | `read_products` |
| `SHOPIFY_LIST_PRODUCT_LISTINGS` | List Product Listings | read | *collection_id, handle, limit, page_info, product_ids, updated_at_min* |  | `read_product_listings`, `write_product_listings` |
| `SHOPIFY_LIST_PRODUCT_LISTING_PRODUCT_IDS` | List Product Listing Product IDs | read | *limit, page_info* |  | `read_product_listings`, `write_product_listings` |
| `SHOPIFY_LIST_PRODUCT_RESOURCE_FEEDBACK` | List Product Resource Feedback | read | `product_id` |  | `read_resource_feedbacks`, `write_resource_feedbacks` |
| `SHOPIFY_LIST_PRODUCT_VARIANTS` | Get product variants | read | `product_id`; *fields, limit, presentment_currencies, since_id* |  | `write_products` |
| `SHOPIFY_LIST_RESOURCE_FEEDBACKS` | Get resource feedback | read |  |  | `resource_feedbacks` |
| `SHOPIFY_LIST_SCRIPT_TAGS` | Retrieve all script tags | read | *created_at_max, created_at_min, fields, limit, since_id, src, updated_at_max, updated_at_min* |  | `read_script_tags` |
| `SHOPIFY_LIST_TRANSACTIONS` | Retrieve order transactions | read | `order_id`; *fields, in_shop_currency, since_id* |  | `write_orders` |
| `SHOPIFY_LIST_USAGE_CHARGES` | List usage charges | read | `recurring_application_charge_id`; *fields* |  |  |
| `SHOPIFY_LIST_USERS` | List Users | read | *limit, page_info* |  | `read_users` |
| `SHOPIFY_MARK_COMMENT_AS_NOT_SPAM` | Mark Comment as Not Spam | write | `comment_id` |  | `read_content`, `write_content` |
| `SHOPIFY_MARK_COMMENT_AS_SPAM` | Mark Comment as Spam | write | `comment_id` |  | `read_content`, `write_content` |
| `SHOPIFY_MODIFY_AN_EXISTING_PRODUCT_VARIANT` | Modify an existing product variant | write | `variant_id`; *barcode, compare_at_price, image_id, inventory_policy, option1, option2, option3, price and 6 more* |  | `write_products` |
| `SHOPIFY_MOVE_FULFILLMENT_ORDER` | Move Fulfillment Order to New Location | write | `fulfillment_order_id`, `new_location_id`; *fulfillment_order_line_items* |  | `write_assigned_fulfillment_orders` |
| `SHOPIFY_OPEN_FULFILLMENT_ORDER` | Open Fulfillment Order | write | `fulfillment_order_id` |  | `write_merchant_managed_fulfillment_orders`, `write_third_party_fulfillment_orders` |
| `SHOPIFY_PIN_METAFIELD_DEFINITION` | Pin metafield definition | write | *definitionId, identifier* |  | `write_products`, `write_customers`, `write_orders`, `write_draft_orders`, `write_content`, `write_locations`, `write_markets`, `write_inventory`, `write_fulfillments`, `write_gift_cards`, `write_files`, `write_discounts`, `write_price_rules`, `write_shipping`, `write_themes`, `write_translations` |
| `SHOPIFY_PUBLISH_COLLECTION_LISTING` | Publish collection listing | write | `collection_listing`, `collection_listing_id` |  | `write_product_listings` |
| `SHOPIFY_PUBLISH_PRODUCT_LISTING` | Publish product listing | write | `product_listing`, `product_listing_id` |  | `write_product_listings` |
| `SHOPIFY_QUERY_APP_BY_KEY` | Query app by API key | read | `api_key` |  | `read_apps` |
| `SHOPIFY_QUERY_BUSINESS_ENTITIES` | Query business entities | read |  |  | `write_products` |
| `SHOPIFY_QUERY_BUSINESS_ENTITY` | Query business entity | read | *id* |  | `read_markets` |
| `SHOPIFY_QUERY_CATALOGS` | Query catalogs | read | *after, before, first, last, query, reverse, sort_key, type* |  | `read_products`, `read_product_listings` |
| `SHOPIFY_QUERY_CHANNEL` | Query channel by ID | read | `id` |  | `write_products` |
| `SHOPIFY_QUERY_CONSENT_POLICY_REGIONS` | Query consent policy regions | read |  |  | `read_privacy_settings` |
| `SHOPIFY_QUERY_CURRENT_APP_INSTALLATION` | Query current app installation | read |  |  | `read_apps` |
| `SHOPIFY_QUERY_CURRENT_BULK_OPERATION` | Query current Shopify bulk operation | read | *type* |  | `read_products`, `write_products` |
| `SHOPIFY_QUERY_CUSTOMER_ACCOUNT_PAGES` | Query customer account pages | read | *after, before, first, last, reverse* |  | `read_content`, `read_online_store_pages` |
| `SHOPIFY_QUERY_DELETION_EVENTS` | Query deletion events | read | *after, before, first, last, query, reverse, sort_key, subject_types* |  | `read_products` |
| `SHOPIFY_QUERY_DISPUTES` | Query disputes | read | *after, before, first, last, query, reverse* |  | `read_shopify_payments_disputes` |
| `SHOPIFY_QUERY_DOMAIN` | Query domain | read | `id` |  |  |
| `SHOPIFY_QUERY_EVENTS` | Query events | read | *after, before, first, last, query, reverse, sort_key* |  |  |
| `SHOPIFY_QUERY_EVENTS_COUNT` | Query events count | read | *query* |  |  |
| `SHOPIFY_QUERY_JOB` | Query job | read | `id` |  |  |
| `SHOPIFY_QUERY_METAFIELD_DEFINITION` | Query metafield definition by ID | read | `id` |  | `read_products`, `write_products`, `read_customers`, `write_customers`, `read_orders`, `write_orders`, `read_draft_orders`, `write_draft_orders`, `read_locations`, `write_locations`, `read_inventory`, `write_inventory`, `read_content`, `write_content` |
| `SHOPIFY_QUERY_METAFIELD_DEFINITIONS` | Query metafield definitions | read | `ownerType`; *after, first, namespace, pinnedStatus* |  | `read_products`, `write_products`, `read_customers`, `write_customers`, `read_orders`, `write_orders`, `read_collections`, `write_collections`, `read_locations`, `write_locations`, `read_companies`, `write_companies`, `read_metaobject_definitions`, `write_metaobject_definitions` |
| `SHOPIFY_QUERY_METAFIELD_DEFINITION_TYPES` | Query metafield definition types | read |  |  |  |
| `SHOPIFY_QUERY_NODE` | Query node by ID | read | `id`; *fields* |  |  |
| `SHOPIFY_QUERY_NODES` | Query nodes by IDs | read | `ids` |  |  |
| `SHOPIFY_QUERY_ONLINE_STORE` | Query online store | read |  |  | `read_online_store_pages` |
| `SHOPIFY_QUERY_PRODUCT_DUPLICATE_JOB` | Query product duplicate job | read | `id` |  | `read_products` |
| `SHOPIFY_QUERY_PUBLICATION` | Query publication by ID | read | `id` |  | `read_publications` |
| `SHOPIFY_QUERY_PUBLIC_API_VERSIONS` | Query public API versions | read |  |  |  |
| `SHOPIFY_QUERY_SERVER_PIXEL` | Query server pixel configuration | read |  |  | `read_pixels`, `read_server_pixels` |
| `SHOPIFY_QUERY_SHOP` | Query Shop (GraphQL) | read | *fields* |  | `read_shop` |
| `SHOPIFY_QUERY_SHOPIFY_FUNCTIONS` | Query Shopify Functions | read | *after, apiType, before, first, last, reverse, useCreationUi* |  | `read_functions` |
| `SHOPIFY_QUERY_SHOP_BILLING_PREFERENCES` | Query shop billing preferences | read |  |  | `read_shop_data` |
| `SHOPIFY_QUERY_SHOP_PAY_PAYMENT_REQUEST_RECEIPTS` | Query Shop Pay payment request receipts | read | *after, before, first, last, query, reverse, sort_key* |  | `read_payment_terms` |
| `SHOPIFY_QUERY_STANDARD_METAFIELD_DEFINITION_TEMPLATES` | Query standard metafield templates | read | *after, before, constraintStatus, constraintSubtype, excludeActivated, first, last, reverse* |  | `read_metafield_definitions` |
| `SHOPIFY_QUERY_TAXONOMY` | Query taxonomy | read | *after, before, childrenOf, descendantsOf, first, last, search, siblingsOf* |  | `read_products` |
| `SHOPIFY_QUERY_WEBHOOK_SUBSCRIPTION` | Query webhook subscription by ID | read | `id` |  | `read_webhooks` |
| `SHOPIFY_QUERY_WEBHOOK_SUBSCRIPTIONS` | Query webhook subscriptions | write | *after, first, query* |  | `read_webhooks` |
| `SHOPIFY_RECEIVE_A_SINGLE_FULFILLMENT` | Receive a single fulfillment | read | `fulfillment_id`, `order_id` |  | `read_orders` |
| `SHOPIFY_REJECT_FULFILLMENT_CANCELLATION_REQUEST` | Reject fulfillment cancellation request | write | `fulfillment_order_id`; *message* |  | `write_assigned_fulfillment_orders` |
| `SHOPIFY_REJECT_FULFILLMENT_REQUEST` | Reject Fulfillment Request | write | `fulfillment_order_id`; *line_items, message, reason* |  | `write_assigned_fulfillment_orders` |
| `SHOPIFY_RELEASE_FULFILLMENT_HOLD` | Release fulfillment hold | write | `fulfillment_order_id` |  | `write_merchant_managed_fulfillment_orders`, `write_third_party_fulfillment_orders`, `write_marketplace_fulfillment_orders` |
| `SHOPIFY_REOPEN_CLOSED_ORDER` | Reopen a closed order | write | `order_id` |  | `write_orders` |
| `SHOPIFY_RESCHEDULE_FULFILLMENT_ORDER` | Reschedule Fulfillment Order | write | `fulfillment_order`, `fulfillment_order_id` |  | `write_merchant_managed_fulfillment_orders`, `write_third_party_fulfillment_orders` |
| `SHOPIFY_RESTORE_COMMENT` | Restore Comment | write | `comment_id` |  | `write_content` |
| `SHOPIFY_RETRIEVES_A_LIST_OF_PRICE_RULES` | Get Price Rules List | read | *created_at_max, created_at_min, ends_at_max, ends_at_min, limit, since_id, starts_at_max, starts_at_min and 3 more* |  | `read_price_rules` |
| `SHOPIFY_RETRIEVE_LIST_METAFIELDS_RESOURCE_S_ENDPOINT` | Get resource metafields | read | `resource`, `resource_id`; *created_at_max, created_at_min, fields, key, limit, metafield_type, namespace, since_id and 2 more* |  | `read_products` |
| `SHOPIFY_RETRIEVE_SPECIFIC_FULFILLMENT_ORDER` | Retrieve specific fulfillment order | read | `fulfillment_order_id`; *include_financial_summaries, include_order_reference_fields* |  | `read_assigned_fulfillment_orders`, `read_merchant_managed_fulfillment_orders`, `read_third_party_fulfillment_orders`, `read_marketplace_fulfillment_orders` |
| `SHOPIFY_REVOKE_APP_ACCESS_SCOPES` | Revoke app access scopes | write | `scopes` |  |  |
| `SHOPIFY_RUN_ANALYTICS_QUERY` | Run Shopify analytics query | read | `query` |  | `read_reports` |
| `SHOPIFY_RUN_BULK_MUTATION_OPERATION` | Run bulk mutation operation | write | `mutation`, `staged_upload_path`; *client_identifier* |  | `write_products` |
| `SHOPIFY_RUN_BULK_OPERATION_QUERY` | Run Shopify bulk operation query | write | `query`; *groupObjects* |  | `read_products`, `read_orders`, `read_customers` |
| `SHOPIFY_SEARCHES_FOR_GIFT_CARDS` | Search Gift Cards | read | *created_at_max, created_at_min, fields, limit, order, page_info, query, updated_at_max and 1 more* |  | `read_gift_cards` |
| `SHOPIFY_SEARCH_PRODUCTS` | Search Shopify products | read | *after, first, query, reverse, sort_key* |  | `read_products` |
| `SHOPIFY_SEND_CUSTOMER_INVITE` | Send Customer Account Invite | write | `customer_id`; *bcc, custom_message, from, subject, to* |  | `write_customers` |
| `SHOPIFY_SEND_INVOICE` | Send an invoice | write | `draft_order_id`; *draft_order_invoice* |  | `write_draft_orders` |
| `SHOPIFY_SET_DEFAULT_CUSTOMER_ADDRESS` | Set Default Customer Address | write | `address_id`, `customer_id` |  | `write_customers` |
| `SHOPIFY_SET_FULFILLMENT_ORDERS_DEADLINE` | Set Fulfillment Orders Deadline | write | `fulfillment_deadline`, `fulfillment_order_ids` |  | `write_merchant_managed_fulfillment_orders`, `write_third_party_fulfillment_orders`, `write_orders` |
| `SHOPIFY_SET_INVENTORY_LEVEL` | Set inventory level | write | `available`, `inventory_item_id`, `location_id`; *disconnect_if_necessary* |  | `write_inventory` |
| `SHOPIFY_SET_METAFIELDS` | Set metafields | write | `metafields` |  | `write_products`, `write_customers`, `write_orders`, `write_draft_orders`, `write_content` |
| `SHOPIFY_SUBMIT_FULFILLMENT_CANCELLATION_REQUEST` | Submit Fulfillment Cancellation Request | write | `fulfillment_order_id`; *cancellation_request* |  | `write_third_party_fulfillment_orders` |
| `SHOPIFY_SUBMIT_FULFILLMENT_REQUEST` | Submit Fulfillment Request | write | `fulfillment_order_id`; *fulfillment_request* |  | `write_third_party_fulfillment_orders` |
| `SHOPIFY_TRIGGER_SHOPIFY_FLOW` | Trigger Shopify Flow | write | *handle, payload* |  | `write_flow_triggers` |
| `SHOPIFY_UNPUBLISH_COLLECTION_LISTING` | Unpublish collection listing | write | `collection_listing_id` |  | `write_product_listings` |
| `SHOPIFY_UNPUBLISH_PRODUCT_LISTING` | Unpublish product listing | write | `product_listing_id` |  | `write_product_listings` |
| `SHOPIFY_UPDATES_AN_EXISTING_DISCOUNT_CODE` | Updates An Existing Discount Code | write | `code`, `discount_code_id`, `price_rule_id` |  | `write_price_rules` |
| `SHOPIFY_UPDATE_APP_SUBSCRIPTION_LINE_ITEM` | Update app subscription line item capped amount | write | `cappedAmount`, `id` |  |  |
| `SHOPIFY_UPDATE_ARTICLE` | Update Article | write | `article_id`, `blog_id`; *author, body_html, image, metafields, published, published_at, summary_html, tags and 1 more* |  | `write_content` |
| `SHOPIFY_UPDATE_BLOG` | Modify an existing blog | write | `blog_id`; *commentable, feedburner, feedburner_location, handle, metafields, template_suffix, title* |  | `write_content`, `read_content` |
| `SHOPIFY_UPDATE_CARRIER_SERVICE` | Update carrier service | write | `carrier_service`, `carrier_service_id` |  | `write_shipping` |
| `SHOPIFY_UPDATE_COMMENT` | Update Comment | write | `comment_id`; *author, body, email, published_at* |  | `write_content` |
| `SHOPIFY_UPDATE_CUSTOMER` | Update Customer | write | `customer_id`; *addresses, email, email_marketing_consent, first_name, last_name, multipass_identifier, note, phone and 7 more* |  | `write_customers`, `read_customers` |
| `SHOPIFY_UPDATE_CUSTOMERS_PARAM_CUSTOMER_ID_ADDRESSES_SET` | Bulk delete customer addresses | write | `address_ids`, `customer_id`, `operation` |  |  |
| `SHOPIFY_UPDATE_CUSTOMER_ADDRESS` | Update Customer Address | write | `address_id`, `customer_id`; *address1, address2, city, company, country, first_name, last_name, phone and 2 more* |  | `read_products` |
| `SHOPIFY_UPDATE_CUSTOM_COLLECTION` | Update custom collection | write | `custom_collection_id`; *body_html, handle, image, published, published_at, published_scope, sort_order, template_suffix and 1 more* |  | `read_products`, `write_products` |
| `SHOPIFY_UPDATE_DISPUTE_EVIDENCE` | Update dispute evidence | write | `dispute_evidence`, `dispute_id` |  | `write_shopify_payments_dispute_evidences` |
| `SHOPIFY_UPDATE_DRAFT_ORDER` | Modify an existing draft order | write | `draft_order_id`; *allow_discount_codes_in_checkout, applied_discount, b2b, billing_address, customer, customer_id, email, line_items and 9 more* |  | `write_draft_orders` |
| `SHOPIFY_UPDATE_FULFILLMENT_SERVICE` | Update Fulfillment Service | write | `fulfillment_service`, `fulfillment_service_id` |  | `write_fulfillments` |
| `SHOPIFY_UPDATE_FULFILLMENT_TRACKING` | Update Fulfillment Tracking | write | `fulfillment`, `fulfillment_id` |  | `write_orders` |
| `SHOPIFY_UPDATE_GIFT_CARD` | Update an Existing Gift Card | write | `gift_card_id`; *customer_id, expires_on, note, template_suffix* |  | `write_gift_cards`, `read_gift_cards` |
| `SHOPIFY_UPDATE_INVENTORY_ITEM` | Update inventory item | write | `inventory_item`, `inventory_item_id` |  | `write_inventory` |
| `SHOPIFY_UPDATE_MARKETING_EVENT` | Update Marketing Event | write | `marketing_event_id`; *budget, budget_type, currency, description, ended_at, event_type, manage_url, marketing_channel and 9 more* |  | `write_marketing_events`, `read_marketing_events` |
| `SHOPIFY_UPDATE_METAFIELD_BY_ID` | Update metafield by ID | write | `metafield_id`; *description, type, value* |  | `write_products` |
| `SHOPIFY_UPDATE_METAFIELD_DEFINITION` | Update metafield definition | write | `key`, `namespace`, `ownerType`; *access, capabilities, description, name, pin, useAsCollectionCondition, validations* |  | `write_products`, `write_customers`, `write_draft_orders`, `write_orders`, `write_locations`, `write_companies`, `write_variants`, `write_collections` |
| `SHOPIFY_UPDATE_MOBILE_PLATFORM_APPLICATION` | Update Mobile Platform Application | write | `mobile_platform_application`, `mobile_platform_application_id` |  | `mobile_platform_applications` |
| `SHOPIFY_UPDATE_ORDER` | Update Order | write | `id`; *billing_address, buyer_accepts_marketing, email, metafields, note, note_attributes, phone, po_number and 5 more* |  | `write_orders` |
| `SHOPIFY_UPDATE_PAGE` | Update Page | write | `page_id`; *author, body_html, handle, published, template_suffix, title* |  | `write_content` |
| `SHOPIFY_UPDATE_PRICE_RULE` | Update Price Rule | write | `price_rule_id`; *allocation_limit, allocation_method, customer_segment_prerequisite_ids, customer_selection, ends_at, entitled_collection_ids, entitled_country_ids, entitled_product_ids and 18 more* |  | `read_price_rules` |
| `SHOPIFY_UPDATE_PRODUCT` | Update a Shopify product | write | `product_id`; *description_html, handle, images, product_type, status, tags, title, variants and 1 more* |  | `write_products` |
| `SHOPIFY_UPDATE_PRODUCT_IMAGE` | Update Product Image | write | `image`, `image_id`, `product_id` |  | `write_products` |
| `SHOPIFY_UPDATE_PRODUCT_METAFIELD` | Update product metafield | write | `metafield_id`, `product_id`; *description, type, value* |  | `write_products` |
| `SHOPIFY_UPDATE_PRODUCT_VARIANTS_BULK` | Bulk update product variants via GraphQL | write | `variants`; *max_poll_attempts, poll_interval_seconds* |  | `write_products` |
| `SHOPIFY_UPDATE_PUB_SUB_WEBHOOK_SUBSCRIPTION` | Update Pub/Sub webhook subscription | write | `id`; *format, pub_sub_project, pub_sub_topic* |  | `write_webhooks` |
| `SHOPIFY_UPDATE_REDIRECT` | Update redirect | write | `redirect_id`; *path, target* |  | `read_content`, `write_content` |
| `SHOPIFY_UPDATE_SCRIPT_TAG` | Update Script Tag | write | `script_tag_id`; *cache, display_scope, event, src* |  | `write_script_tags` |
| `SHOPIFY_UPDATE_SMART_COLLECTION` | Update Smart Collection | write | `smart_collection_id`; *body_html, disjunctive, handle, image, published, published_at, published_scope, rules and 3 more* |  | `write_products` |
| `SHOPIFY_UPDATE_SMART_COLLECTION_ORDER` | Update smart collection product order | write | `products`, `smart_collection_id` |  | `write_products` |
| `SHOPIFY_UPDATE_THEME` | Modify existing theme | write | `theme_id`; *name, role* |  | `write_themes` |
| `SHOPIFY_UPDATE_WEBHOOK_SUBSCRIPTION` | Update webhook subscription | write | `id`, `webhook_subscription` |  |  |
| `SHOPIFY_UPLOAD_DISPUTE_EVIDENCE_FILE` | Upload dispute evidence file | write | `dispute_file_upload`, `dispute_id` |  | `write_shopify_payments_dispute_file_uploads` |
