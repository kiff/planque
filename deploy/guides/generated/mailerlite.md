# MailerLite: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner pastes a key from the vendor's settings). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

94 tools: 40 read, 54 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `MAILERLITE_ADD_PRODUCT_TO_CATEGORY` | Add Product to Category | write | `category_id`, `product_id`, `shop_id` |  |  |
| `MAILERLITE_ADD_SUBSCRIBER_TO_GROUP` | Add Subscriber To Group | write | `group_id`, `subscriber_id` |  |  |
| `MAILERLITE_CANCEL_CAMPAIGN` | Cancel Campaign | write | `campaign_id` |  |  |
| `MAILERLITE_CREATE_AUTOMATION` | Create automation | write | `name`, `trigger`; *enabled* |  |  |
| `MAILERLITE_CREATE_AUTOMATION_DELAY_STEP` | Create Automation Delay Step | write | `automation_id`, `unit`, `value`; *branch, move_to_branch, parent_id* |  |  |
| `MAILERLITE_CREATE_CAMPAIGN` | Create Campaign | write | `emails`, `name`, `type`; *groups, language_id, segments, settings* |  |  |
| `MAILERLITE_CREATE_ECOMMERCE_CART_ITEM` | Create E-commerce Cart Item | write | `cart_id`, `ecommerce_product_id`, `shop_id`; *price, quantity, variant* |  |  |
| `MAILERLITE_CREATE_ECOMMERCE_CATEGORY` | Create E-commerce Category | write | `name`, `shop_id`; *exclude_from_automations* |  |  |
| `MAILERLITE_CREATE_ECOMMERCE_CUSTOMER` | Create/Update E-commerce Customer | write | `email`, `shop_id`; *accepts_marketing, create_subscriber, total_spent* |  |  |
| `MAILERLITE_CREATE_ECOMMERCE_ORDER` | Create E-commerce Order | write | `cart`, `customer`, `shop_id`; *status, total_price* |  |  |
| `MAILERLITE_CREATE_ECOMMERCE_PRODUCT` | Create E-commerce Product | write | `name`, `price`, `shop_id`; *description, exclude_from_automations, image, short_description, url* |  |  |
| `MAILERLITE_CREATE_ECOMMERCE_SHOP` | Create E-commerce Shop | write | `currency`, `name`, `url`; *access_data, enable_popups, enabled, group_id, platform* |  |  |
| `MAILERLITE_CREATE_FIELD` | Create Field | write | `name`, `type` |  |  |
| `MAILERLITE_CREATE_GROUP` | Create Group | write | `name` |  |  |
| `MAILERLITE_CREATE_SEGMENT` | Create Segment | write | `name` |  |  |
| `MAILERLITE_CREATE_SUBSCRIBER` | Create Subscriber | write | `email`; *fields, groups, ip_address, opted_in_at, optin_ip, resubscribe, status, subscribed_at and 1 more* |  |  |
| `MAILERLITE_CREATE_WEBHOOK` | Create Webhook | write | `events`, `url`; *batchable, enabled, name* |  |  |
| `MAILERLITE_DELETE_AUTOMATION` | Delete Automation | write | `automation_id` |  |  |
| `MAILERLITE_DELETE_CAMPAIGN` | Delete Campaign | write | `campaign_id` |  |  |
| `MAILERLITE_DELETE_ECOMMERCE_CART_ITEM` | Delete E-commerce Cart Item | write | `cart_id`, `item_id`, `shop_id` |  |  |
| `MAILERLITE_DELETE_ECOMMERCE_CATEGORY` | Delete E-commerce Category | write | `category_id`, `shop_id` |  |  |
| `MAILERLITE_DELETE_ECOMMERCE_CUSTOMER` | Delete E-commerce Customer | write | `customer_id`, `shop_id` |  |  |
| `MAILERLITE_DELETE_ECOMMERCE_ORDER` | Delete E-commerce Order | write | `order_id`, `shop_id` |  |  |
| `MAILERLITE_DELETE_ECOMMERCE_PRODUCT` | Delete E-commerce Product | write | `product_id`, `shop_id` |  |  |
| `MAILERLITE_DELETE_ECOMMERCE_SHOP` | Delete E-commerce Shop | write | `shop_id` |  |  |
| `MAILERLITE_DELETE_FIELD` | Delete Field | write | `field_id` |  |  |
| `MAILERLITE_DELETE_GROUP` | Delete Group | write | `group_id` |  |  |
| `MAILERLITE_DELETE_SEGMENT` | Delete Segment | write | `segment_id` |  |  |
| `MAILERLITE_DELETE_SUBSCRIBER` | Delete Subscriber | write | `id` |  |  |
| `MAILERLITE_DELETE_WEBHOOK` | Delete Webhook | write | `webhook_id` |  |  |
| `MAILERLITE_EXECUTE_BATCH_REQUEST` | Execute Batch Request | write | `requests` |  |  |
| `MAILERLITE_FORGET_SUBSCRIBER` | Forget Subscriber | write | `subscriber_id` |  |  |
| `MAILERLITE_GET_ACCOUNT_INFO` | Get Account Info | read |  |  |  |
| `MAILERLITE_GET_ACCOUNT_STATS` | Get Account Stats | read |  |  |  |
| `MAILERLITE_GET_AUTOMATION` | Get Automation | read | `automation_id`; *fields* |  |  |
| `MAILERLITE_GET_AUTOMATION_ACTIVITY` | Get Automation Activity | read | `automation_id`, `filter_status` |  |  |
| `MAILERLITE_GET_CAMPAIGN` | Get Campaign | read | `campaign_id` |  |  |
| `MAILERLITE_GET_CAMPAIGNS` | Get Campaigns | read | *limit, page, status, type* |  |  |
| `MAILERLITE_GET_CAMPAIGN_LANGUAGES` | Get Campaign Languages | read |  |  |  |
| `MAILERLITE_GET_CAMPAIGN_SUBSCRIBER_ACTIVITY` | Get Campaign Subscriber Activity | read | `campaign_id`; *activity_type, include, limit, page, search, sort* |  |  |
| `MAILERLITE_GET_ECOMMERCE_CART` | Get E-commerce Cart | read | `cart_id`, `shop_id` |  |  |
| `MAILERLITE_GET_ECOMMERCE_CART_ITEM` | Get E-commerce Cart Item | read | `cart_id`, `item_id`, `shop_id` |  |  |
| `MAILERLITE_GET_ECOMMERCE_CATEGORY` | Get E-commerce Category | read | `category_id`, `shop_id` |  |  |
| `MAILERLITE_GET_ECOMMERCE_CUSTOMER` | Get E-commerce Customer | read | `customer_id`, `shop_id` |  |  |
| `MAILERLITE_GET_ECOMMERCE_CUSTOMERS` | Get E-commerce Customers | read | `shop_id`; *limit, page* |  |  |
| `MAILERLITE_GET_ECOMMERCE_ORDER` | Get E-commerce Order | read | `order_id`, `shop_id` |  |  |
| `MAILERLITE_GET_ECOMMERCE_PRODUCT` | Get E-commerce Product | read | `product_id`, `shop_id` |  |  |
| `MAILERLITE_GET_ECOMMERCE_SHOP` | Get E-commerce Shop | read | `shop_id` |  |  |
| `MAILERLITE_GET_ECOMMERCE_SHOPS` | Get E-commerce Shops | read | *limit, page* |  |  |
| `MAILERLITE_GET_FIELDS` | Get Fields | read | *filter_keyword, filter_type, limit, page, sort* |  |  |
| `MAILERLITE_GET_GROUP` | Get Group | read | `group_id` |  |  |
| `MAILERLITE_GET_GROUPS` | Get Groups | read | *filter_name, limit, page, sort* |  |  |
| `MAILERLITE_GET_GROUP_SUBSCRIBERS` | Get Group Subscribers | read | `group_id`; *cursor, filter_status, include, limit* |  |  |
| `MAILERLITE_GET_IMPORT_STATUS` | Get Import Status | read | `import_id` |  |  |
| `MAILERLITE_GET_SEGMENT` | Get Segment | read | `segment_id` |  |  |
| `MAILERLITE_GET_SEGMENTS` | Get Segments | read | *limit, page* |  |  |
| `MAILERLITE_GET_SEGMENT_SUBSCRIBERS` | Get Segment Subscribers | read | `segment_id`; *cursor, filter_status, per_page* |  |  |
| `MAILERLITE_GET_SUBSCRIBER` | Get Subscriber | read | `subscriber_id` |  |  |
| `MAILERLITE_GET_SUBSCRIBERS` | Get Subscribers | read | *cursor, filter_status, include, limit* |  |  |
| `MAILERLITE_GET_SUBSCRIBERS_STATS` | Get Subscribers Stats | read | *x_locale* |  |  |
| `MAILERLITE_GET_SUBSCRIBER_ACTIVITY` | Get Subscriber Activity | read | `subscriber_id`; *cursor, filter_type, limit* |  |  |
| `MAILERLITE_GET_WEBHOOK` | Get Webhook | read | `webhook_id` |  |  |
| `MAILERLITE_GET_WEBHOOKS` | Get Webhooks | read | *limit, page* |  |  |
| `MAILERLITE_IMPORT_ECOMMERCE_CATEGORIES` | Import E-commerce Categories | write | `categories`, `shop_id` |  |  |
| `MAILERLITE_IMPORT_ECOMMERCE_ORDERS` | Import E-commerce Orders | write | `orders`, `shop_id` |  |  |
| `MAILERLITE_IMPORT_ECOMMERCE_PRODUCTS` | Import E-commerce Products | write | `products`, `shop_id` |  |  |
| `MAILERLITE_IMPORT_SUBSCRIBERS_TO_GROUP` | Import Subscribers to Group | write | `group_id`, `subscribers`; *autoresponders, resubscribe, return_status* |  |  |
| `MAILERLITE_LIST_AUTOMATIONS` | List Automations | read | *filter_enabled, filter_group, filter_name, limit, page* |  |  |
| `MAILERLITE_LIST_ECOMMERCE_CARTS` | List E-commerce Carts | read | `shop_id`; *limit, page* |  |  |
| `MAILERLITE_LIST_ECOMMERCE_CART_ITEMS` | List E-commerce Cart Items | read | `cart_id`, `shop_id`; *limit, page* |  |  |
| `MAILERLITE_LIST_ECOMMERCE_CATEGORIES` | List E-commerce Categories | read | `shop_id`; *limit, page* |  |  |
| `MAILERLITE_LIST_ECOMMERCE_CATEGORY_PRODUCTS` | List E-commerce Category Products | read | `category_id`, `shop_id`; *limit, page* |  |  |
| `MAILERLITE_LIST_ECOMMERCE_ORDERS` | List E-commerce Orders | read | `shop_id`; *limit, page* |  |  |
| `MAILERLITE_LIST_ECOMMERCE_PRODUCTS` | List E-commerce Products | read | `shop_id`; *limit, page* |  |  |
| `MAILERLITE_LIST_FORMS` | List Forms | read | `type`; *filter_name, limit, page, sort* |  |  |
| `MAILERLITE_LIST_TIMEZONES` | List Timezones | read |  |  |  |
| `MAILERLITE_REMOVE_PRODUCT_FROM_CATEGORY` | Remove Product From Category | write | `category_id`, `product_id`, `shop_id` |  |  |
| `MAILERLITE_REMOVE_SUBSCRIBER_FROM_GROUP` | Remove Subscriber From Group | write | `group_id`, `subscriber_id` |  |  |
| `MAILERLITE_SEND_OR_SCHEDULE_CAMPAIGN` | Send or Schedule Campaign | write | `campaign_id`; *delivery, resend, schedule, sending_time_test_schedule* |  |  |
| `MAILERLITE_SET_DOUBLE_OPTIN` | Set Double Opt-In | write | `double_opt_in` |  |  |
| `MAILERLITE_UPDATE_AUTOMATION_DELAY_STEP` | Update Automation Delay Step | write | `automation_id`, `step_id`, `unit`, `value`; *parent_id, reschedule_waiting_subscribers* |  |  |
| `MAILERLITE_UPDATE_CAMPAIGN` | Update Campaign | write | `campaign_id`, `emails`, `name`; *ab_settings, groups, language_id, resend_settings, segments* |  |  |
| `MAILERLITE_UPDATE_ECOMMERCE_CART` | Update E-commerce Cart | write | `cart_id`, `shop_id`; *cart_total, checkout_url* |  |  |
| `MAILERLITE_UPDATE_ECOMMERCE_CART_ITEM` | Update E-commerce Cart Item | write | `cart_id`, `ecommerce_product_id`, `item_id`, `shop_id`; *price, quantity, variant* |  |  |
| `MAILERLITE_UPDATE_ECOMMERCE_CATEGORY` | Update E-commerce Category | write | `category_id`, `name`, `shop_id`; *exclude_from_automations* |  |  |
| `MAILERLITE_UPDATE_ECOMMERCE_CUSTOMER` | Update E-commerce Customer | write | `customer_id`, `email`, `shop_id`; *accepts_marketing, create_subscriber, total_spent* |  |  |
| `MAILERLITE_UPDATE_ECOMMERCE_ORDER` | Update E-commerce Order | write | `order_id`, `shop_id`; *status, total_price* |  |  |
| `MAILERLITE_UPDATE_ECOMMERCE_PRODUCT` | Update E-commerce Product | write | `name`, `product_id`, `shop_id`; *description, exclude_from_automations, image, price, short_description, url* |  |  |
| `MAILERLITE_UPDATE_ECOMMERCE_SHOP` | Update E-commerce Shop | write | `currency`, `name`, `shop_id`, `url`; *access_data, enable_popups, enabled, group_id, platform* |  |  |
| `MAILERLITE_UPDATE_FIELD` | Update Field | write | `field_id`, `name` |  |  |
| `MAILERLITE_UPDATE_GROUP` | Update Group | write | `id`, `name` |  |  |
| `MAILERLITE_UPDATE_SEGMENT` | Update Segment | write | `id`, `name` |  |  |
| `MAILERLITE_UPDATE_SUBSCRIBER` | Update Subscriber | write | `id`; *fields, groups, ip_address, opted_in_at, optin_ip, status, subscribed_at, unsubscribed_at* |  |  |
| `MAILERLITE_UPDATE_WEBHOOK` | Update Webhook | write | `events`, `url`, `webhook_id`; *batchable, enabled, name* |  |  |
