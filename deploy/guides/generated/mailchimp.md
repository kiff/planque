# Mailchimp: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

275 tools: 144 read, 131 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `MAILCHIMP_ADD_AUTOMATION` | Add automation | write | `recipients__list__id`; *recipients__store__id, settings__from__name, settings__reply__to, trigger__settings__workflow__type* |  |  |
| `MAILCHIMP_ADD_BATCH_WEBHOOK` | Add batch webhook | write | `url`; *enabled* |  |  |
| `MAILCHIMP_ADD_CAMPAIGN` | Add campaign | write | `type`; *content_type, recipients__list__id, recipients__segment__opts__conditions, recipients__segment__opts__match, recipients__segment__opts__saved__segment__id, rss__opts__constrain__rss__img, rss__opts__feed__url, rss__opts__frequency and 45 more* |  |  |
| `MAILCHIMP_ADD_CAMPAIGN_FEEDBACK` | Add Campaign Feedback | write | `campaign_id`, `message`; *block_id, is_complete* |  |  |
| `MAILCHIMP_ADD_CAMPAIGN_FOLDER` | Add campaign folder | write | `name` |  |  |
| `MAILCHIMP_ADD_CART` | Add cart | write | `currency_code`, `customer__email__address`, `customer__id`, `customer__opt__in__status`, `id`, `lines`, `order_total`, `store_id`; *campaign_id, checkout_url, customer__address__address1, customer__address__address2, customer__address__city, customer__address__country, customer__address__country__code, customer__address__postal__code and 6 more* |  |  |
| `MAILCHIMP_ADD_CART_LINE_ITEM` | Add cart line item | write | `cart_id`, `id`, `price`, `product_id`, `product_variant_id`, `quantity`, `store_id` |  |  |
| `MAILCHIMP_ADD_CONNECTED_SITE` | Add connected site | write | `domain`, `foreign_id` |  |  |
| `MAILCHIMP_ADD_CONTACT_TO_AUDIENCE` | Add contact to audience | write | `audience_id`; *data_mode, email_channel, language, merge_field_validation_mode, merge_fields, sms_channel, tags* |  |  |
| `MAILCHIMP_ADD_DOMAIN_TO_ACCOUNT` | Add domain to account | write | `verification_email` |  |  |
| `MAILCHIMP_ADD_EVENT` | Add event | write | `list_id`, `name`, `subscriber_hash`; *is_syncing, occurred_at, properties* |  |  |
| `MAILCHIMP_ADD_EXPORT` | Add export | read | `include_stages`; *since_timestamp* |  |  |
| `MAILCHIMP_ADD_FILE` | Upload File to File Manager | write | `file_data`, `name`; *folder_id* |  |  |
| `MAILCHIMP_ADD_FOLDER` | Add folder | write | `name`; *created_at, created_by, file_count* |  |  |
| `MAILCHIMP_ADD_INTEREST_CATEGORY` | Add interest category | write | `list_id`, `title`, `type`; *display_order* |  |  |
| `MAILCHIMP_ADD_INTEREST_IN_CATEGORY` | Add interest in category | write | `interest_category_id`, `list_id`, `name`; *display_order* |  |  |
| `MAILCHIMP_ADD_LANDING_PAGE` | Add landing page | write | *description, list_id, name, store_id, template_id, title, tracking__enable__restricted__data__processing, tracking__track__with__mailchimp and 2 more* |  |  |
| `MAILCHIMP_ADD_LIST` | Add list | write | `campaign__defaults__from__email`, `campaign__defaults__from__name`, `campaign__defaults__language`, `campaign__defaults__subject`, `contact__address1`, `contact__city`, `contact__company`, `contact__country`, `name`, `permission_reminder`; *contact__address2, contact__phone, contact__state, contact__zip, double_optin, email_type_option, marketing_permissions, notify_on_subscribe and 2 more* |  |  |
| `MAILCHIMP_ADD_LIST_MEMBER` | Add list member | write | `email_address`, `list_id`, `status`; *email_type, interests, ip_opt, ip_signup, language, location, marketing_permissions, merge_fields and 5 more* |  |  |
| `MAILCHIMP_ADD_MEMBER_NOTE` | Add member note | write | `list_id`, `note`, `subscriber_hash` |  |  |
| `MAILCHIMP_ADD_MEMBER_TO_SEGMENT` | Add member to segment | write | `email_address`, `list_id`, `segment_id` |  |  |
| `MAILCHIMP_ADD_MERGE_FIELD` | Add merge field | write | `list_id`, `name`, `type`; *default_value, display_order, help_text, options__choices, options__date__format, options__default__country, options__phone__format, options__size and 3 more* |  |  |
| `MAILCHIMP_ADD_ORDER_LINE_ITEM` | Add order line item | write | `id`, `order_id`, `price`, `product_id`, `product_variant_id`, `quantity`, `store_id`; *discount* |  |  |
| `MAILCHIMP_ADD_OR_REMOVE_MEMBER_TAGS` | Add or remove member tags | write | `list_id`, `subscriber_hash`, `tags`; *is_syncing* |  |  |
| `MAILCHIMP_ADD_OR_UPDATE_CUSTOMER` | Add or update customer | write | `customer_id`, `email_address`, `id`, `opt_in_status`, `store_id`; *address__address1, address__address2, address__city, address__country, address__country__code, address__postal__code, address__province, address__province__code and 3 more* |  |  |
| `MAILCHIMP_ADD_OR_UPDATE_LIST_MEMBER` | Add or update list member | write | `email_address`, `list_id`, `status_if_new`, `subscriber_hash`; *email_type, interests, ip_opt, ip_signup, language, location__latitude, location__longitude, marketing_permissions and 6 more* |  |  |
| `MAILCHIMP_ADD_OR_UPDATE_PRODUCT_VARIANT` | Add or update product variant | write | `id`, `product_id`, `store_id`, `title`, `variant_id`; *backorders, image_url, inventory_quantity, price, sku, url, visibility* |  |  |
| `MAILCHIMP_ADD_PRODUCT` | Add product | write | `id`, `store_id`, `title`, `variants`; *description, handle, image_url, images, published_at_foreign, type, url, vendor* |  |  |
| `MAILCHIMP_ADD_PRODUCT_IMAGE` | Add product image | write | `id`, `product_id`, `store_id`, `url`; *variant_ids* |  |  |
| `MAILCHIMP_ADD_PROMO_CODE` | Add promo code | write | `code`, `id`, `promo_rule_id`, `redemption_url`, `store_id`; *created_at_foreign, enabled, updated_at_foreign, usage_count* |  |  |
| `MAILCHIMP_ADD_PROMO_RULE` | Add promo rule | write | `amount`, `description`, `id`, `store_id`, `target`, `type`; *created_at_foreign, enabled, ends_at, starts_at, title, updated_at_foreign* |  |  |
| `MAILCHIMP_ADD_SEGMENT` | Add segment | write | `list_id`, `name`; *options__conditions, options__match, static_segment* |  |  |
| `MAILCHIMP_ADD_STORE` | Add store | write | `currency_code`, `id`, `list_id`, `name`; *address__address1, address__address2, address__city, address__country, address__country__code, address__latitude, address__longitude, address__postal__code and 10 more* |  |  |
| `MAILCHIMP_ADD_SUBSCRIBER_TO_WORKFLOW_EMAIL` | Add subscriber to workflow email | write | `email_address`, `workflow_email_id`, `workflow_id` |  |  |
| `MAILCHIMP_ADD_TEMPLATE` | Add template | write | `html`, `name`; *folder_id* |  |  |
| `MAILCHIMP_ADD_TEMPLATE_FOLDER` | Add Template Folder | write | `name` |  |  |
| `MAILCHIMP_ADD_WEBHOOK` | Add webhook | write | `list_id`, `url`; *events__campaign, events__cleaned, events__profile, events__subscribe, events__unsubscribe, events__upemail, sources__admin, sources__api and 1 more* |  |  |
| `MAILCHIMP_ARCHIVE_AUTOMATION` | Archive automation | write | `workflow_id` |  |  |
| `MAILCHIMP_ARCHIVE_CONTACT` | Archive contact | write | `audience_id`, `contact_id` |  |  |
| `MAILCHIMP_ARCHIVE_LIST_MEMBER` | Archive list member | write | `list_id`, `subscriber_hash` |  |  |
| `MAILCHIMP_BATCH_ADD_OR_REMOVE_MEMBERS` | Batch add or remove members | write | `list_id`, `segment_id`; *members_to_add, members_to_remove* |  |  |
| `MAILCHIMP_BATCH_SUBSCRIBE_OR_UNSUBSCRIBE` | Batch subscribe or unsubscribe | write | `list_id`, `members`; *skip_duplicate_check, skip_merge_validation, sync_tags, update_existing* |  |  |
| `MAILCHIMP_CAMPAIGN_ABUSE_REPORT_DETAILS` | Get Campaign Abuse Report Details | read | `campaign_id`, `report_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_CAMPAIGN_STATISTICS_FEEDBACK` | Campaign statistics feedback | read | `campaign_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_CANCEL_CAMPAIGN` | Cancel campaign | write | `campaign_id` |  |  |
| `MAILCHIMP_CREATE_A_SURVEY_CAMPAIGN` | Create a survey campaign | write | `list_id`, `survey_id` |  |  |
| `MAILCHIMP_CUSTOMER_JOURNEYS_API_TRIGGER_FOR_A_CONTACT` | Customer journeys api trigger for a contact | write | `email_address`, `journey_id`, `step_id` |  |  |
| `MAILCHIMP_CUSTOMIZE_SIGNUP_FORM` | Customize signup form | write | `list_id`; *contents, header__image__align, header__image__alt, header__image__border__color, header__image__border__style, header__image__border__width, header__image__height, header__image__link and 5 more* |  |  |
| `MAILCHIMP_DELETE_BATCH_REQUEST` | Delete batch request | write | `batch_id` |  |  |
| `MAILCHIMP_DELETE_BATCH_WEBHOOK` | Delete batch webhook | write | `batch_webhook_id` |  |  |
| `MAILCHIMP_DELETE_CAMPAIGN` | Delete campaign | write | `campaign_id` |  |  |
| `MAILCHIMP_DELETE_CAMPAIGN_FEEDBACK_MESSAGE` | Delete campaign feedback message | write | `campaign_id`, `feedback_id` |  |  |
| `MAILCHIMP_DELETE_CAMPAIGN_FOLDER` | Delete campaign folder | write | `folder_id` |  |  |
| `MAILCHIMP_DELETE_CART` | Delete cart | write | `cart_id`, `store_id` |  |  |
| `MAILCHIMP_DELETE_CART_LINE_ITEM` | Delete cart line item | write | `cart_id`, `line_id`, `store_id` |  |  |
| `MAILCHIMP_DELETE_CONNECTED_SITE` | Delete connected site | write | `connected_site_id` |  |  |
| `MAILCHIMP_DELETE_CUSTOMER` | Delete customer | write | `customer_id`, `store_id` |  |  |
| `MAILCHIMP_DELETE_DOMAIN` | Delete domain | write | `domain_name` |  |  |
| `MAILCHIMP_DELETE_FILE` | Delete file | write | `file_id` |  |  |
| `MAILCHIMP_DELETE_FOLDER` | Delete folder | write | `folder_id` |  |  |
| `MAILCHIMP_DELETE_INTEREST_CATEGORY` | Delete interest category | write | `interest_category_id`, `list_id` |  |  |
| `MAILCHIMP_DELETE_INTEREST_IN_CATEGORY` | Delete interest in category | write | `interest_category_id`, `interest_id`, `list_id` |  |  |
| `MAILCHIMP_DELETE_LANDING_PAGE` | Delete landing page | write | `page_id` |  |  |
| `MAILCHIMP_DELETE_LIST` | Delete list | write | `list_id` |  |  |
| `MAILCHIMP_DELETE_LIST_MEMBER` | Delete list member | write | `list_id`, `subscriber_hash` |  |  |
| `MAILCHIMP_DELETE_MERGE_FIELD` | Delete merge field | write | `list_id`, `merge_id` |  |  |
| `MAILCHIMP_DELETE_NOTE` | Delete note | write | `list_id`, `note_id`, `subscriber_hash` |  |  |
| `MAILCHIMP_DELETE_ORDER` | Delete order | write | `order_id`, `store_id` |  |  |
| `MAILCHIMP_DELETE_ORDER_LINE_ITEM` | Delete order line item | write | `line_id`, `order_id`, `store_id` |  |  |
| `MAILCHIMP_DELETE_PRODUCT` | Delete product | write | `product_id`, `store_id` |  |  |
| `MAILCHIMP_DELETE_PRODUCT_IMAGE` | Delete product image | write | `image_id`, `product_id`, `store_id` |  |  |
| `MAILCHIMP_DELETE_PRODUCT_VARIANT` | Delete product variant | write | `product_id`, `store_id`, `variant_id` |  |  |
| `MAILCHIMP_DELETE_PROMO_CODE` | Delete promo code | write | `promo_code_id`, `promo_rule_id`, `store_id` |  |  |
| `MAILCHIMP_DELETE_PROMO_RULE` | Delete promo rule | write | `promo_rule_id`, `store_id` |  |  |
| `MAILCHIMP_DELETE_SEGMENT` | Delete segment | write | `list_id`, `segment_id` |  |  |
| `MAILCHIMP_DELETE_STORE` | Delete store | write | `store_id` |  |  |
| `MAILCHIMP_DELETE_TEMPLATE` | Delete template | write | `template_id` |  |  |
| `MAILCHIMP_DELETE_TEMPLATE_FOLDER` | Delete template folder | write | `folder_id` |  |  |
| `MAILCHIMP_DELETE_WEBHOOK` | Delete webhook | write | `list_id`, `webhook_id` |  |  |
| `MAILCHIMP_DELETE_WORKFLOW_EMAIL` | Delete workflow email | write | `workflow_email_id`, `workflow_id` |  |  |
| `MAILCHIMP_FORGET_CONTACT` | Forget contact | write | `audience_id`, `contact_id` |  |  |
| `MAILCHIMP_GET_ABUSE_REPORT` | Get abuse report | read | `list_id`, `report_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_ACCOUNT_EXPORT_INFO` | Get account export info | read | `export_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_AUDIENCES_CONTACTS` | Get audience contacts | read | `audience_id`; *count, created_before, created_since, cursor, exclude_fields, fields, sort_dir, sort_field and 2 more* |  |  |
| `MAILCHIMP_GET_AUDIENCES_CONTACTS_DETAIL` | Get audience contact detail | read | `audience_id`, `contact_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_AUTHORIZED_APP_INFO` | Get authorized app info | read | `app_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_AUTOMATED_EMAIL_SUBSCRIBER` | Get automated email subscriber | read | `subscriber_hash`, `workflow_email_id`, `workflow_id` |  |  |
| `MAILCHIMP_GET_AUTOMATION_INFO` | Get automation info | read | `workflow_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_BATCH_OPERATION_STATUS` | Get batch operation status | read | `batch_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_BATCH_WEBHOOK_INFO` | Get batch webhook info | read | `batch_webhook_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_CAMPAIGN_CONTENT` | Get campaign content | read | `campaign_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_CAMPAIGN_FEEDBACK_MESSAGE` | Get campaign feedback message | read | `campaign_id`, `feedback_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_CAMPAIGN_FOLDER` | Get campaign folder | read | `folder_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_CAMPAIGN_INFO` | Get campaign info | read | `campaign_id`; *exclude_fields, fields, include_resend_shortcut_eligibility* |  |  |
| `MAILCHIMP_GET_CAMPAIGN_LINK_DETAILS` | Get campaign link details | read | `campaign_id`, `link_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_CAMPAIGN_RECIPIENT_INFO` | Get campaign recipient info | read | `campaign_id`, `subscriber_hash`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_CAMPAIGN_REPORT` | Get campaign report | read | `campaign_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_CAMPAIGN_SEND_CHECKLIST` | Get campaign send checklist | read | `campaign_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_CART_INFO` | Get cart info | read | `cart_id`, `store_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_CART_LINE_ITEM` | Get cart line item | read | `cart_id`, `line_id`, `store_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_CLICKED_LINK_SUBSCRIBER` | Get clicked link subscriber | read | `campaign_id`, `link_id`, `subscriber_hash`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_CONNECTED_SITE` | Get connected site | read | `connected_site_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_CUSTOMER_INFO` | Get customer info | read | `customer_id`, `store_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_DOMAIN_INFO` | Get domain info | read | `domain_name` |  |  |
| `MAILCHIMP_GET_FACEBOOK_AD_INFO` | Get facebook ad info | read | `outreach_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_FACEBOOK_AD_REPORT` | Get facebook ad report | read | `outreach_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_FILE` | Get file | read | `file_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_FOLDER` | Get folder | read | `folder_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_GROWTH_HISTORY_BY_MONTH` | Get growth history by month | read | `list_id`, `month`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_INFORMATION_ABOUT_ALL_SURVEYS_FOR_A_LIST` | Get information about all surveys for a list | read | `list_id` |  |  |
| `MAILCHIMP_GET_INTEREST_CATEGORY_INFO` | Get interest category info | read | `interest_category_id`, `list_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_INTEREST_IN_CATEGORY` | Get interest in category | read | `interest_category_id`, `interest_id`, `list_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_LANDING_PAGE_CONTENT` | Get landing page content | read | `page_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_LANDING_PAGE_INFO` | Get landing page info | read | `page_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_LANDING_PAGE_REPORT` | Get landing page report | read | `outreach_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_LATEST_CHIMP_CHATTER` | Get latest chimp chatter | read | *count, offset* |  |  |
| `MAILCHIMP_GET_LISTS_INFO` | Get lists info | read | *before_campaign_last_sent, before_date_created, count, email, exclude_fields, fields, has_ecommerce_store, include_total_contacts and 5 more* |  |  |
| `MAILCHIMP_GET_LIST_INFO` | Get list info | read | `list_id`; *exclude_fields, fields, include_total_contacts* |  |  |
| `MAILCHIMP_GET_MEMBER_INFO` | Get member info | read | `list_id`, `subscriber_hash`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_MEMBER_NOTE` | Get member note | read | `list_id`, `note_id`, `subscriber_hash`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_MERGE_FIELD` | Get merge field | read | `list_id`, `merge_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_MESSAGE` | Get message | read | `conversation_id`, `message_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_OPENED_CAMPAIGN_SUBSCRIBER` | Get opened campaign subscriber | read | `campaign_id`, `subscriber_hash`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_ORDER_INFO` | Get order info | read | `order_id`, `store_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_ORDER_LINE_ITEM` | Get order line item | read | `line_id`, `order_id`, `store_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_PRODUCT_IMAGE_INFO` | Get product image info | read | `image_id`, `product_id`, `store_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_PRODUCT_INFO` | Get product info | read | `product_id`, `store_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_PRODUCT_VARIANT_INFO` | Get product variant info | read | `product_id`, `store_id`, `variant_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_PROMO_CODE` | Get promo code | read | `promo_code_id`, `promo_rule_id`, `store_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_PROMO_RULE` | Get promo rule | read | `promo_rule_id`, `store_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_SEGMENT_INFO` | Get segment info | read | `list_id`, `segment_id`; *exclude_fields, fields, include_cleaned, include_transactional, include_unsubscribed* |  |  |
| `MAILCHIMP_GET_STORE_INFO` | Get store info | read | `store_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_SUBSCRIBER_EMAIL_ACTIVITY` | Get subscriber email activity | read | `campaign_id`, `subscriber_hash`; *exclude_fields, fields, since* |  |  |
| `MAILCHIMP_GET_SUBSCRIBER_REMOVED_FROM_WORKFLOW` | Get subscriber removed from workflow | read | `subscriber_hash`, `workflow_id` |  |  |
| `MAILCHIMP_GET_SURVEY` | Get survey | read | `list_id`, `survey_id` |  |  |
| `MAILCHIMP_GET_SURVEY_QUESTION_REPORT` | Get survey question report | read | `question_id`, `survey_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_SURVEY_REPORT` | Get survey report | read | `survey_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_SURVEY_RESPONSE` | Get survey response | read | `response_id`, `survey_id` |  |  |
| `MAILCHIMP_GET_TEMPLATE_FOLDER` | Get Template Folder | read | `folder_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_TEMPLATE_INFO` | Get template info | read | `template_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_UNSUBSCRIBED_MEMBER` | Get unsubscribed member | read | `campaign_id`, `subscriber_hash`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_GET_WEBHOOK_INFO` | Get webhook info | read | `list_id`, `webhook_id` |  |  |
| `MAILCHIMP_GET_WORKFLOW_EMAIL_INFO` | Get workflow email info | read | `workflow_email_id`, `workflow_id` |  |  |
| `MAILCHIMP_LIST_ABUSE_REPORTS` | List abuse reports | read | `list_id`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_ACCOUNT_EXPORTS` | List account exports | read | *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_ACCOUNT_ORDERS` | List account orders | read | *campaign_id, count, customer_id, exclude_fields, fields, has_outreach, offset, outreach_id* |  |  |
| `MAILCHIMP_LIST_ANSWERS_FOR_QUESTION` | List answers for question | read | `question_id`, `survey_id`; *exclude_fields, fields, respondent_familiarity_is* |  |  |
| `MAILCHIMP_LIST_API_ROOT_RESOURCES` | Get API Root and Account Info | read | *exclude_fields, fields* |  |  |
| `MAILCHIMP_LIST_AUTHORIZED_APPS` | List authorized apps | read | *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_AUTOMATED_EMAILS` | List automated emails | read | `workflow_id` |  |  |
| `MAILCHIMP_LIST_AUTOMATED_EMAIL_SUBSCRIBERS` | List automated email subscribers | read | `workflow_email_id`, `workflow_id` |  |  |
| `MAILCHIMP_LIST_AUTOMATIONS` | List automations | read | *before_create_time, before_start_time, count, exclude_fields, fields, offset, since_create_time, since_start_time and 1 more* |  |  |
| `MAILCHIMP_LIST_BATCH_REQUESTS` | List batch requests | read | *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_BATCH_WEBHOOKS` | List batch webhooks | read | *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_CAMPAIGNS` | List campaigns | read | *before_create_time, before_send_time, count, exclude_fields, fields, folder_id, include_resend_shortcut_eligibility, list_id and 8 more* |  |  |
| `MAILCHIMP_LIST_CAMPAIGN_CLICK_DETAILS` | List campaign click details | read | `campaign_id`; *count, exclude_fields, fields, filter_bots, offset, sort_dir, sort_field* |  |  |
| `MAILCHIMP_LIST_CAMPAIGN_FEEDBACK` | List campaign feedback | read | `campaign_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_LIST_CAMPAIGN_FOLDERS` | List campaign folders | read | *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_CAMPAIGN_OPEN_DETAILS` | List campaign open details | read | `campaign_id`; *count, exclude_fields, fields, offset, since, sort_dir, sort_field* |  |  |
| `MAILCHIMP_LIST_CAMPAIGN_PRODUCT_ACTIVITY` | List campaign product activity | read | `campaign_id`; *count, exclude_fields, fields, offset, sort_field* |  |  |
| `MAILCHIMP_LIST_CAMPAIGN_RECIPIENTS` | List campaign recipients | read | `campaign_id`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_CAMPAIGN_REPORTS` | List campaign reports | read | *before_send_time, count, exclude_fields, fields, offset, since_send_time, type* |  |  |
| `MAILCHIMP_LIST_CARTS` | List carts | read | `store_id`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_CART_LINE_ITEMS` | List cart line items | read | `cart_id`, `store_id`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_CHILD_CAMPAIGN_REPORTS` | List child campaign reports | read | `campaign_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_LIST_CLICKED_LINK_SUBSCRIBERS` | List clicked link subscribers | read | `campaign_id`, `link_id`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_CONNECTED_SITES` | List connected sites | read | *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_CUSTOMERS` | List customers | read | `store_id`; *count, email_address, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_CUSTOMER_JOURNEYS` | List customer journeys | read | *count* |  |  |
| `MAILCHIMP_LIST_DOMAIN_PERFORMANCE_STATS` | List domain performance stats | read | `campaign_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_LIST_EEP_URL_ACTIVITY` | List EepURL Activity | read | `campaign_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_LIST_EMAIL_ACTIVITY` | List email activity | read | `campaign_id`; *count, exclude_fields, fields, offset, since* |  |  |
| `MAILCHIMP_LIST_FACEBOOK_ADS` | List facebook ads | read | *count, exclude_fields, fields, offset, sort_dir, sort_field* |  |  |
| `MAILCHIMP_LIST_FACEBOOK_ADS_REPORTS` | List facebook ads reports | read | *count, exclude_fields, fields, offset, sort_dir, sort_field* |  |  |
| `MAILCHIMP_LIST_FACEBOOK_ECOMMERCE_REPORT` | List facebook ecommerce report | read | `outreach_id`; *count, exclude_fields, fields, offset, sort_field* |  |  |
| `MAILCHIMP_LIST_FOLDERS` | List folders | read | *before_created_at, count, created_by, exclude_fields, fields, offset, since_created_at* |  |  |
| `MAILCHIMP_LIST_GROWTH_HISTORY_DATA` | List growth history data | read | `list_id`; *count, exclude_fields, fields, offset, sort_dir, sort_field* |  |  |
| `MAILCHIMP_LIST_INTERESTS_IN_CATEGORY` | List interests in category | read | `interest_category_id`, `list_id`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_INTEREST_CATEGORIES` | List interest categories | read | `list_id`; *count, exclude_fields, fields, offset, type* |  |  |
| `MAILCHIMP_LIST_LANDING_PAGES` | List landing pages | read | *count, exclude_fields, fields, sort_dir, sort_field* |  |  |
| `MAILCHIMP_LIST_LANDING_PAGES_REPORTS` | List landing pages reports | read | *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_LOCATIONS` | List locations | read | `list_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_LIST_MEMBERS_INFO` | List members info | read | `list_id`; *before_last_changed, before_timestamp_opt, count, email_type, exclude_fields, fields, interest_category_id, interest_ids and 11 more* |  |  |
| `MAILCHIMP_LIST_MEMBERS_IN_SEGMENT` | List members in segment | read | `list_id`, `segment_id`; *count, exclude_fields, fields, include_cleaned, include_transactional, include_unsubscribed, offset* |  |  |
| `MAILCHIMP_LIST_MEMBER_EVENTS` | List member events | read | `list_id`, `subscriber_hash`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_MEMBER_GOAL_EVENTS` | List member goal events | read | `list_id`, `subscriber_hash`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_LIST_MEMBER_TAGS` | List member tags | read | `list_id`, `subscriber_hash`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_MERGE_FIELDS` | List merge fields | read | `list_id`; *count, exclude_fields, fields, offset, required, type* |  |  |
| `MAILCHIMP_LIST_ORDERS` | List orders | read | `store_id`; *campaign_id, count, customer_id, exclude_fields, fields, has_outreach, offset, outreach_id* |  |  |
| `MAILCHIMP_LIST_ORDER_LINE_ITEMS` | List order line items | read | `order_id`, `store_id`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_PRODUCT` | List product | read | `store_id`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_PRODUCT_IMAGES` | List product images | read | `product_id`, `store_id`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_PRODUCT_VARIANTS` | List product variants | read | `product_id`, `store_id`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_PROMO_CODES` | List promo codes | read | `promo_rule_id`, `store_id`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_PROMO_RULES` | List promo rules | read | `store_id`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_RECENT_ACTIVITY` | List recent activity | read | `list_id`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_RECENT_MEMBER_NOTES` | List recent member notes | read | `list_id`, `subscriber_hash`; *count, exclude_fields, fields, offset, sort_dir, sort_field* |  |  |
| `MAILCHIMP_LIST_SEGMENTS` | List segments | read | `list_id`; *before_created_at, before_updated_at, count, exclude_fields, fields, include_cleaned, include_transactional, include_unsubscribed and 4 more* |  |  |
| `MAILCHIMP_LIST_SENDING_DOMAINS` | List sending domains | read |  |  |  |
| `MAILCHIMP_LIST_SIGNUP_FORMS` | List signup forms | read | `list_id` |  |  |
| `MAILCHIMP_LIST_STORED_FILES` | List stored files | read | *before_created_at, count, created_by, exclude_fields, fields, offset, since_created_at, sort_dir and 2 more* |  |  |
| `MAILCHIMP_LIST_STORES` | List stores | read | *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_SUBSCRIBERS_REMOVED_FROM_WORKFLOW` | List subscribers removed from workflow | read | `workflow_id` |  |  |
| `MAILCHIMP_LIST_SURVEY_QUESTION_REPORTS` | List survey question reports | read | `survey_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_LIST_SURVEY_REPORTS` | List survey reports | read | *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_SURVEY_RESPONSES` | List survey responses | read | `survey_id`; *answered_question, chose_answer, count, exclude_fields, fields, offset, respondent_familiarity_is* |  |  |
| `MAILCHIMP_LIST_TEMPLATES` | List templates | read | *before_date_created, category, content_type, count, created_by, exclude_fields, fields, folder_id and 5 more* |  |  |
| `MAILCHIMP_LIST_TEMPLATE_FOLDERS` | List template folders | read | *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_TOP_EMAIL_CLIENTS` | List top email clients | read | `list_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_LIST_TOP_OPEN_ACTIVITIES` | List top open activities | read | `campaign_id`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_UNSUBSCRIBED_MEMBERS` | List unsubscribed members | read | `campaign_id`; *count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_LIST_WEBHOOKS` | List webhooks | read | `list_id` |  |  |
| `MAILCHIMP_PAUSE_AUTOMATED_EMAIL` | Pause automated email | write | `workflow_email_id`, `workflow_id` |  |  |
| `MAILCHIMP_PAUSE_AUTOMATION_EMAILS` | Pause automation emails | write | `workflow_id` |  |  |
| `MAILCHIMP_PAUSE_RSS_CAMPAIGN` | Pause RSS Campaign | write | `campaign_id` |  |  |
| `MAILCHIMP_PING` | Ping | read |  |  |  |
| `MAILCHIMP_PUBLISH_A_SURVEY` | Publish a survey | write | `list_id`, `survey_id` |  |  |
| `MAILCHIMP_PUBLISH_LANDING_PAGE` | Publish landing page | write | `page_id` |  |  |
| `MAILCHIMP_REMOVE_LIST_MEMBER_FROM_SEGMENT` | Remove list member from segment | write | `list_id`, `segment_id`, `subscriber_hash` |  |  |
| `MAILCHIMP_REMOVE_SUBSCRIBER_FROM_WORKFLOW` | Remove subscriber from workflow | write | `email_address`, `workflow_id` |  |  |
| `MAILCHIMP_REPLICATE_CAMPAIGN` | Replicate campaign | write | `campaign_id` |  |  |
| `MAILCHIMP_RESEND_CAMPAIGN` | Resend campaign | write | `campaign_id`; *shortcut_type* |  |  |
| `MAILCHIMP_RESUME_RSS_CAMPAIGN` | Resume rss campaign | write | `campaign_id` |  |  |
| `MAILCHIMP_RETRIEVE_CAMPAIGN_ABUSE_COMPLAINTS` | Retrieve campaign abuse complaints | read | `campaign_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_RETRIEVE_FOLDER_CONTENTS` | Retrieve folder contents | read | `folder_id`; *before_created_at, count, created_by, exclude_fields, fields, offset, since_created_at, sort_dir and 2 more* |  |  |
| `MAILCHIMP_SCHEDULE_CAMPAIGN` | Schedule campaign | write | `campaign_id`, `schedule_time`; *batch__delivery__batch__count, batch__delivery__batch__delay, timewarp* |  |  |
| `MAILCHIMP_SEARCH_CAMPAIGNS` | Search campaigns | read | `query`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_SEARCH_FOR_TAGS_ON_A_LIST_BY_NAME` | Search for tags on a list by name | read | `list_id`; *name* |  |  |
| `MAILCHIMP_SEARCH_MEMBERS` | Search members | read | `query`; *exclude_fields, fields, list_id* |  |  |
| `MAILCHIMP_SEND_CAMPAIGN` | Send campaign | write | `campaign_id` |  |  |
| `MAILCHIMP_SEND_TEST_EMAIL` | Send test email | write | `campaign_id`, `test_emails`; *send_type* |  |  |
| `MAILCHIMP_SET_CAMPAIGN_CONTENT` | Set campaign content | write | `campaign_id`; *archive__archive__content, archive__archive__type, html, plain_text, template__id, url, variate_contents* |  |  |
| `MAILCHIMP_START_AUTOMATED_EMAIL` | Start automated email | write | `workflow_email_id`, `workflow_id` |  |  |
| `MAILCHIMP_START_AUTOMATION_EMAILS` | Start automation emails | write | `workflow_id` |  |  |
| `MAILCHIMP_START_BATCH_OPERATION` | Start batch operation | write | `operations` |  |  |
| `MAILCHIMP_UNPUBLISH_A_SURVEY` | Unpublish a survey | write | `list_id`, `survey_id` |  |  |
| `MAILCHIMP_UNPUBLISH_LANDING_PAGE` | Unpublish landing page | write | `page_id` |  |  |
| `MAILCHIMP_UNSCHEDULE_CAMPAIGN` | Unschedule campaign | write | `campaign_id` |  |  |
| `MAILCHIMP_UPDATE_AUDIENCES_CONTACTS` | Update contact in audience | write | `audience_id`, `contact_id`; *data_mode, email_channel, language, merge_field_validation_mode, merge_fields, sms_channel, tags* |  |  |
| `MAILCHIMP_UPDATE_BATCH_WEBHOOK` | Update batch webhook | write | `batch_webhook_id`; *enabled, url* |  |  |
| `MAILCHIMP_UPDATE_CAMPAIGN_FEEDBACK_MESSAGE` | Update campaign feedback message | write | `campaign_id`, `feedback_id`; *block_id, is_complete, message* |  |  |
| `MAILCHIMP_UPDATE_CAMPAIGN_FOLDER` | Update campaign folder | write | `folder_id`, `name` |  |  |
| `MAILCHIMP_UPDATE_CAMPAIGN_SETTINGS` | Update campaign settings | write | `campaign_id`; *recipients__list__id, recipients__segment__opts__conditions, recipients__segment__opts__match, recipients__segment__opts__prebuilt__segment__id, recipients__segment__opts__saved__segment__id, rss__opts__constrain__rss__img, rss__opts__feed__url, rss__opts__frequency and 45 more* |  |  |
| `MAILCHIMP_UPDATE_CART` | Update cart | write | `cart_id`, `store_id`; *campaign_id, checkout_url, currency_code, customer__address__address1, customer__address__address2, customer__address__city, customer__address__country, customer__address__country__code and 10 more* |  |  |
| `MAILCHIMP_UPDATE_CART_LINE_ITEM` | Update cart line item | write | `cart_id`, `line_id`, `store_id`; *price, product_id, product_variant_id, quantity* |  |  |
| `MAILCHIMP_UPDATE_CUSTOMER` | Update customer | write | `customer_id`, `store_id`; *address__address1, address__address2, address__city, address__country, address__country__code, address__postal__code, address__province, address__province__code and 4 more* |  |  |
| `MAILCHIMP_UPDATE_FILE` | Update file | write | `file_id`; *folder_id, name* |  |  |
| `MAILCHIMP_UPDATE_FOLDER` | Update folder | write | `folder_id`, `name`; *created_at, created_by, file_count* |  |  |
| `MAILCHIMP_UPDATE_INTEREST_CATEGORY` | Update interest category | write | `interest_category_id`, `list_id`, `title`; *display_order, type* |  |  |
| `MAILCHIMP_UPDATE_INTEREST_IN_CATEGORY` | Update interest in category | write | `interest_category_id`, `interest_id`, `list_id`, `name`; *display_order* |  |  |
| `MAILCHIMP_UPDATE_LANDING_PAGE` | Update landing page | write | `page_id`; *description, list_id, name, store_id, title, tracking__enable__restricted__data__processing, tracking__track__with__mailchimp* |  |  |
| `MAILCHIMP_UPDATE_LISTS` | Update lists | write | `list_id`; *campaign__defaults__from__email, campaign__defaults__from__name, campaign__defaults__language, campaign__defaults__subject, contact__address1, contact__address2, contact__city, contact__company and 12 more* |  |  |
| `MAILCHIMP_UPDATE_LISTS_SEGMENTS` | Update lists segments | write | `list_id`, `name`, `segment_id`; *options__conditions, options__match, static_segment* |  |  |
| `MAILCHIMP_UPDATE_LIST_MEMBER` | Update list member | write | `list_id`, `subscriber_hash`; *email_address, email_type, interests, ip_opt, ip_signup, language, location__latitude, location__longitude and 7 more* |  |  |
| `MAILCHIMP_UPDATE_MERGE_FIELD` | Update merge field | write | `list_id`, `merge_id`; *default_value, display_order, help_text, name, options__choices, options__date__format, options__default__country, options__phone__format and 3 more* |  |  |
| `MAILCHIMP_UPDATE_NOTE` | Update note | write | `list_id`, `note_id`, `subscriber_hash`; *note* |  |  |
| `MAILCHIMP_UPDATE_ORDER` | Update order | write | `order_id`, `store_id`; *billing__address__address1, billing__address__address2, billing__address__city, billing__address__company, billing__address__country, billing__address__country__code, billing__address__latitude, billing__address__longitude and 50 more* |  |  |
| `MAILCHIMP_UPDATE_ORDER_LINE_ITEM` | Update order line item | write | `line_id`, `order_id`, `store_id`; *discount, price, product_id, product_variant_id, quantity* |  |  |
| `MAILCHIMP_UPDATE_PRODUCT` | Update product | write | `product_id`, `store_id`; *description, handle, image_url, images, published_at_foreign, title, type, url and 2 more* |  |  |
| `MAILCHIMP_UPDATE_PRODUCT_IMAGE` | Update product image | write | `image_id`, `product_id`, `store_id`; *id, url, variant_ids* |  |  |
| `MAILCHIMP_UPDATE_PRODUCT_VARIANT` | Update product variant | write | `product_id`, `store_id`, `variant_id`; *backorders, image_url, inventory_quantity, price, sku, title, url, visibility* |  |  |
| `MAILCHIMP_UPDATE_PROMO_CODE` | Update promo code | write | `promo_code_id`, `promo_rule_id`, `store_id`; *code, created_at_foreign, enabled, redemption_url, updated_at_foreign, usage_count* |  |  |
| `MAILCHIMP_UPDATE_PROMO_RULE` | Update promo rule | write | `promo_rule_id`, `store_id`; *amount, created_at_foreign, description, enabled, ends_at, starts_at, target, title and 2 more* |  |  |
| `MAILCHIMP_UPDATE_STORE` | Update store | write | `store_id`; *address__address1, address__address2, address__city, address__country, address__country__code, address__latitude, address__longitude, address__postal__code and 12 more* |  |  |
| `MAILCHIMP_UPDATE_TEMPLATE` | Update template | write | `html`, `name`, `template_id`; *folder_id* |  |  |
| `MAILCHIMP_UPDATE_TEMPLATE_FOLDER` | Update Template Folder | write | `folder_id`, `name` |  |  |
| `MAILCHIMP_UPDATE_WEBHOOK` | Update webhook | write | `list_id`, `webhook_id`; *events__campaign, events__cleaned, events__profile, events__subscribe, events__unsubscribe, events__upemail, sources__admin, sources__api and 2 more* |  |  |
| `MAILCHIMP_UPDATE_WORKFLOW_EMAIL` | Update workflow email | write | `workflow_email_id`, `workflow_id`; *delay__action, delay__amount, delay__direction, delay__type, settings__from__name, settings__preview__text, settings__reply__to, settings__subject__line and 1 more* |  |  |
| `MAILCHIMP_UPSERT_ECOMMERCE_STORES_ORDERS` | Upsert ecommerce stores orders | write | `customer__id`, `id`, `order_id`, `store_id`; *billing__address__address1, billing__address__address2, billing__address__city, billing__address__company, billing__address__country, billing__address__country__code, billing__address__latitude, billing__address__longitude and 52 more* |  |  |
| `MAILCHIMP_UPSERT_ECOMMERCE_STORES_PRODUCTS` | Upsert ecommerce stores products | write | `id`, `product_id`, `store_id`, `title`, `variants`; *description, handle, image_url, images, published_at_foreign, type, url, vendor* |  |  |
| `MAILCHIMP_VERIFY_CONNECTED_SITE_SCRIPT` | Verify connected site script | read | `connected_site_id` |  |  |
| `MAILCHIMP_VERIFY_DOMAIN` | Verify domain | write | `code`, `domain_name` |  |  |
| `MAILCHIMP_VIEW_DEFAULT_CONTENT` | View default content | read | `template_id`; *exclude_fields, fields* |  |  |
| `MAILCHIMP_VIEW_RECENT_ACTIVITY` | View recent activity | read | `list_id`, `subscriber_hash`; *activity_filters, count, exclude_fields, fields, offset* |  |  |
| `MAILCHIMP_VIEW_RECENT_ACTIVITY50` | View recent activity 50 | read | `list_id`, `subscriber_hash`; *action, exclude_fields, fields* |  |  |
