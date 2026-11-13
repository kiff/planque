# WhatsApp: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

127 tools: 50 read, 77 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `WHATSAPP_ACCEPT_BUSINESS_SOLUTION` | Accept business solution | write | `partner_app_id`, `solution_id`; *log_session_id* |  | `whatsapp_business_management` |
| `WHATSAPP_ACCEPT_SOLUTION_DEACTIVATION` | Accept solution deactivation | write | `solution_id`; *fields* |  | `whatsapp_business_management` |
| `WHATSAPP_ADD_BUSINESS_PHONE_NUMBER` | Add business phone number | write | `business_id`, `phone_number`; *user_agent* |  | `business_management` |
| `WHATSAPP_APPROVE_GROUP_JOIN_REQUESTS` | Approve group join requests | write | `group_id`, `join_requests`; *messaging_product* |  | `whatsapp_business_messaging` |
| `WHATSAPP_ASSIGN_USER` | Assign user | write | `tasks`, `user`, `whatsapp_business_account_id` |  | `whatsapp_business_management` |
| `WHATSAPP_BLOCK_USERS` | Block WhatsApp Users | write | `block_users`, `phone_number_id` |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_CONFIGURE_CONVERSATIONAL_AUTOMATION` | Configure WhatsApp Conversational Automation | write | `phone_number_id`; *commands, enable_welcome_message, prompts* |  | `whatsapp_business_messaging`, `whatsapp_business_management` |
| `WHATSAPP_CREATE_BUSINESS_ACCOUNT` | Create business account | write | `business_id`, `name`; *business_type, currency, on_behalf_of_business_id, primary_funding_id, purchase_order_number, timezone_id* |  | `business_management` |
| `WHATSAPP_CREATE_BUSINESS_SOLUTION` | Create business solution | write | `application_id`, `owner_permissions`, `partner_app_id`, `partner_permissions`, `solution_name` |  | `whatsapp_business_management` |
| `WHATSAPP_CREATE_FLOW` | Create WhatsApp Flow | write | `categories`, `endpoint_uri`, `name`; *clone_flow_id, waba_id* |  | `whatsapp_business_messaging`, `whatsapp_business_management` |
| `WHATSAPP_CREATE_GROUP` | Create group | write | `messaging_product`, `phone_number_id`, `subject`; *description, join_approval_mode* |  | `whatsapp_business_messaging` |
| `WHATSAPP_CREATE_GROUP_INVITE_LINK` | Create group invite link | write | `group_id`, `messaging_product` |  | `whatsapp_business_messaging` |
| `WHATSAPP_CREATE_MAX_PRICE_AGREEMENTS` | Create WhatsApp Max Price Agreement | write | `business_id`, `currency`, `max_price_micros`, `message_type`, `signer_email`, `signer_name` |  | `business_management`, `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_CREATE_MESSAGE_TEMPLATE` | Create message template | write | `category`, `language`, `name`; *allow_category_change, components, library_template_body_inputs, library_template_button_inputs, library_template_name* |  | `whatsapp_business_management` |
| `WHATSAPP_CREATE_OBO_MOBILITY_INTENT` | Create obo mobility intent | write | `body`, `waba_id` |  | `whatsapp_business_management` |
| `WHATSAPP_CREATE_PHONE_NUMBER` | Create phone number | write | `phone_number`, `verified_name`, `waba_id`; *cc, migrate_phone_number, preverified_id* |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_CREATE_QR_CODE` | Create QR Code | write | `phone_number_id`, `prefilled_message`; *generate_qr_image* |  | `whatsapp_business_management` |
| `WHATSAPP_CREATE_SCHEDULE` | Create schedule | write | `request_body`; *waba_id* |  | `whatsapp_business_management` |
| `WHATSAPP_CREATE_UPLOAD_SESSION` | Create upload session for WhatsApp | write | `file_length`, `file_type`; *file_name* |  | `whatsapp_business_messaging` |
| `WHATSAPP_DELETE_FLOW` | Delete flow | write | `flow_id` |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_DELETE_GROUP` | Delete group | write | `group_id` |  | `whatsapp_business_messaging` |
| `WHATSAPP_DELETE_GROUP_INVITE_LINK` | Delete group invite link | write | `group_id`; *messaging_product* |  | `whatsapp_business_messaging` |
| `WHATSAPP_DELETE_MEDIA` | Delete WhatsApp Media | write | `media_id` |  | `whatsapp_business_messaging` |
| `WHATSAPP_DELETE_MESSAGE_TEMPLATE` | Delete message template | write | `name`; *hsm_id, waba_id* |  | `whatsapp_business_management` |
| `WHATSAPP_DELETE_PREVERIFIED_NUMBER` | Delete preverified number | write | `preverified_phone_number_id` |  | `whatsapp_business_management` |
| `WHATSAPP_DELETE_QR_CODE` | Delete QR Code | write | `phone_number_id`, `qr_code_id` |  | `whatsapp_business_management` |
| `WHATSAPP_DEPRECATE_FLOW` | Deprecate flow | write | `flow_id` |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_DEREGISTER_PHONE` | Deregister WhatsApp Business Phone Number | write | `phone_number_id` |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_DOWNLOAD_MEDIA` | Download media | read | `media_url` |  | `whatsapp_business_messaging` |
| `WHATSAPP_EDIT_MESSAGE_TEMPLATE` | Edit message template | write | `template_id`; *allow_category_change, category, components, cta_url_link_tracking_opted_out, display_format, is_primary_device_delivery_only, message_send_ttl_seconds, parameter_format and 1 more* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_ACCOUNT_NUMBER` | Get account number | read | `account_number_id`; *fields* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_ACTIVITIES` | Get WhatsApp Business Account Activities | read | *activity_type, after, before, fields, limit, since, until, waba_id* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_ASSIGNED_BUSINESS_ACCOUNTS` | Get assigned business accounts | read | `user_id`; *after, before, fields, limit* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_ASSIGNED_USERS` | Get assigned users | read | `business`, `waba_id`; *after, before, fields, limit* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_BUSINESS_ACCOUNT_DETAILS` | Get WhatsApp Business Account Details | read | *fields, waba_id* |  | `whatsapp_business_messaging` |
| `WHATSAPP_GET_BUSINESS_BOT` | Get business bot | read | `waba_bot_id`; *fields* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_BUSINESS_COMPLIANCE_INFO` | Get WhatsApp Business Compliance Info | read | `phone_number_id`; *fields* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_BUSINESS_ENCRYPTION` | Get Business Encryption | read | `phone_number_id`; *fields* |  | `whatsapp_business_messaging` |
| `WHATSAPP_GET_BUSINESS_PORTFOLIO` | Get business portfolio | read | `business_id`; *fields* |  | `business_management` |
| `WHATSAPP_GET_BUSINESS_PROFILE` | Get business profile | read | `phone_number_id`; *fields* |  | `whatsapp_business_messaging` |
| `WHATSAPP_GET_BUSINESS_PROFILE_NODE` | Get business profile node | read | `whatsapp_business_profile_id`; *fields* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_CALL_PERMISSIONS` | Get call permissions | read | `phone_number_id`; *recipient, user_wa_id* |  | `whatsapp_business_messaging` |
| `WHATSAPP_GET_CLIENT_BUSINESS_ACCOUNTS` | Get client WhatsApp Business Accounts | read | `business_id`; *after, before, business_type, fields, find, limit* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_COMMERCE_SETTINGS` | Get WhatsApp Commerce Settings | read | `phone_number_id` |  | `whatsapp_business_messaging` |
| `WHATSAPP_GET_CONNECTED_CLIENT_BUSINESSES` | Get connected client businesses | read | `application_id`; *after, before, fields, limit* |  | `business_management` |
| `WHATSAPP_GET_EXTENDED_CREDITS` | Get extended credits | read | `business_id`; *after, before, fields, limit, order_by_is_owned_credential* |  | `business_management` |
| `WHATSAPP_GET_FLOW` | Get flow | write | `flow_id`; *fields* |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_GET_FLOW_ASSETS` | Get Flow Assets | read | `flow_id` |  | `whatsapp_business_messaging` |
| `WHATSAPP_GET_GROUP` | Get group | read | `group_id`; *fields* |  | `whatsapp_business_messaging` |
| `WHATSAPP_GET_JOIN_REQUESTS` | Get Join Requests | read | `group_id`; *after, before, limit* |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_GET_MAX_PRICE_AGREEMENTS` | Get max price agreements | read | `business_id` |  | `business_management` |
| `WHATSAPP_GET_MEDIA_INFO` | Get media info | read | `media_id` |  | `whatsapp_business_messaging` |
| `WHATSAPP_GET_MESSAGE_HISTORY` | Get WhatsApp Message History | read | `phone_number_id`; *after, before, fields, limit, message_id* |  | `whatsapp_business_messaging`, `whatsapp_business_management` |
| `WHATSAPP_GET_MESSAGE_HISTORY_EVENTS` | Get message history events | read | `message_history_id`; *after, before, delivery_status, fields, limit* |  | `whatsapp_business_messaging` |
| `WHATSAPP_GET_MESSAGE_TEMPLATES` | Get message templates | read | *after, category, fields, language, limit, name_or_content, status, waba_id* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_MIGRATION_INTENT` | Get migration intent | read | `migration_intent_id`; *fields* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_OFFICIAL_BUSINESS_ACCOUNT_STATUS` | Get official business account status | read | `phone_number_id`; *fields* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_OWNED_BUSINESS_ACCOUNTS` | Get owned WhatsApp Business Accounts | read | `business_id`; *after, before, business_type, fields, find, first, last* |  | `whatsapp_business_messaging` |
| `WHATSAPP_GET_PHONE_NUMBER` | Get phone number | read | `phone_number_id`; *fields* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_PHONE_NUMBERS` | Get phone numbers | read | *limit, waba_id* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_PREVERIFIED_NUMBER` | Get preverified number | read | `preverified_phone_number_id`; *fields* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_PREVERIFIED_NUMBER_PARTNERS` | Get preverified number partners | read | `pre_verified_phone_number_id`; *after, before, fields, limit* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_QR_CODE` | Get qr code | read | `phone_number_id`, `qr_code_id`; *fields* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_REACH_ESTIMATE` | Get reach estimate | read | `date_interval`, `targeting_spec`, `waba_id` |  | `whatsapp_business_management` |
| `WHATSAPP_GET_SCHEDULES` | Get WhatsApp Schedules | read | *after, before, fields, filtering, limit, sort, waba_id* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_SETTINGS` | Get WhatsApp Cloud API Settings | read | `phone_number_id`; *include_sip_credentials* |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_GET_SOLUTION_ACCESS_TOKEN` | Get solution access token | read | `business_id`, `solution_id` |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_GET_SUBSCRIBED_APPS` | Get Subscribed Apps | read | *fields, waba_id* |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_GET_TEMPLATE_LIBRARY` | Get message template library | read | *after, category, country, language, limit, name_or_content, status* |  | `whatsapp_business_messaging`, `whatsapp_business_management` |
| `WHATSAPP_GET_TEMPLATE_STATUS` | Get template status | read | `template_id`; *fields* |  | `whatsapp_business_management` |
| `WHATSAPP_GET_UPLOAD_OFFSET` | Get upload offset | read | `upload_session_id` |  | `whatsapp_business_management` |
| `WHATSAPP_LIST_BLOCKED_USERS` | List blocked users | read | `phone_number_id`; *after, before, limit* |  | `whatsapp_business_messaging` |
| `WHATSAPP_LIST_BUSINESS_ACCOUNTS` | List business accounts | read | `business_id`; *after, before, business_type, fields, find, limit* |  | `business_management`, `whatsapp_business_management`, `whatsapp_business_messaging`, `public_profile` |
| `WHATSAPP_LIST_BUSINESS_SOLUTIONS` | List business solutions | read | `application_id`; *after, before, fields, limit, role* |  | `whatsapp_business_management` |
| `WHATSAPP_LIST_FLOWS` | List WhatsApp Flows | read | *after, limit, waba_id* |  | `whatsapp_business_management` |
| `WHATSAPP_LIST_GROUPS` | List WhatsApp Groups | read | `phone_number_id`; *after, before, limit* |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_LIST_ON_BEHALF_REQUESTS` | List on behalf requests | read | `waba_id`; *after, before, fields, limit* |  | `whatsapp_business_management` |
| `WHATSAPP_LIST_PREVERIFIED_NUMBERS` | List preverified numbers | read | `business_id`; *after, before, code_verification_status, fields, limit, phone_number* |  | `business_management` |
| `WHATSAPP_LIST_QR_CODES` | List qr codes | write | `phone_number_id`; *after, before, code, fields, limit* |  | `whatsapp_business_management` |
| `WHATSAPP_LIST_SOLUTIONS` | List WhatsApp Solutions | read | *after, before, fields, limit, waba_id* |  | `whatsapp_business_management` |
| `WHATSAPP_MANAGE_CALL` | Manage call | write | `action`, `phone_number_id`; *biz_opaque_callback_data, call_id, messaging_account_id, recipient, recording, session, to, transcription* |  | `whatsapp_business_messaging` |
| `WHATSAPP_MIGRATE_FLOWS` | Migrate flows | write | `destination_waba_id`, `source_waba_id`; *source_flow_names* |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_ONBOARD_MARKETING_PARTNERS` | Onboard marketing partners | write | `business_id`; *solution_id* |  | `business_management` |
| `WHATSAPP_PUBLISH_FLOW` | Publish WhatsApp Flow | write | `flow_id` |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_REGISTER_PHONE` | Register WhatsApp Business Phone Number | write | `phone_number_id`, `pin`; *backup, data_localization_region, messaging_product, meta_store_retention_minutes* |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_REJECT_BUSINESS_SOLUTION` | Reject business solution | write | `request_type`, `solution_id`; *partner_app_id, rejection_reason* |  | `whatsapp_business_management` |
| `WHATSAPP_REJECT_GROUP_JOIN_REQUESTS` | Reject group join requests | write | `group_id`, `join_requests`; *messaging_product* |  | `whatsapp_business_messaging` |
| `WHATSAPP_REJECT_SOLUTION_DEACTIVATION` | Reject solution deactivation | write | `reject_deactivation_request`, `solution_id`; *fields* |  | `whatsapp_business_management` |
| `WHATSAPP_REMOVE_ASSIGNED_USER` | Remove Assigned User from WhatsApp Business Account | write | `user`; *waba_id* |  | `whatsapp_business_management` |
| `WHATSAPP_REMOVE_GROUP_PARTICIPANTS` | Remove group participants | write | `group_id`, `participants`; *messaging_product* |  | `whatsapp_business_messaging` |
| `WHATSAPP_REQUEST_PREVERIFIED_NUMBER_CODE` | Request preverified number code | write | `code_method`, `language`, `pre_verified_phone_number_id`; *user_agent* |  | `whatsapp_business_management` |
| `WHATSAPP_REQUEST_SOLUTION_DEACTIVATION` | Request solution deactivation | write | `solution_id`; *reason* |  | `whatsapp_business_management` |
| `WHATSAPP_REQUEST_VERIFICATION_CODE` | Request verification code | write | `code_method`, `language`, `phone_number_id`; *user_agent* |  | `whatsapp_business_management` |
| `WHATSAPP_SEND_CONTACTS` | Send contacts | write | `contacts`, `phone_number_id`, `to_number`; *recipient_type, reply_to_message_id* |  | `whatsapp_business_messaging` |
| `WHATSAPP_SEND_ENCRYPTED_MESSAGE` | Send encrypted message | write | `encrypted_contents`, `phone_number_id`; *messaging_product* |  | `whatsapp_business_messaging` |
| `WHATSAPP_SEND_INTERACTIVE_BUTTONS` | Send interactive buttons | write | `body_text`, `buttons`, `phone_number_id`, `to_number`; *footer_text, header_text, reply_to_message_id* |  | `whatsapp_business_messaging` |
| `WHATSAPP_SEND_INTERACTIVE_LIST` | Send interactive list | write | `body_text`, `button_text`, `phone_number_id`, `sections`, `to_number`; *footer_text, header_text, reply_to_message_id* |  | `whatsapp_business_messaging` |
| `WHATSAPP_SEND_LOCATION` | Send location | write | `address`, `latitude`, `longitude`, `name`, `phone_number_id`, `to_number`; *reply_to_message_id* |  | `whatsapp_business_messaging` |
| `WHATSAPP_SEND_MARKETING_MESSAGE` | Send marketing message | write | `messaging_product`, `phone_number_id`, `recipient_type`, `template`, `type`; *bid_spec, message_activity_sharing, product_policy, recipient, to* |  | `whatsapp_business_messaging` |
| `WHATSAPP_SEND_MEDIA` | Send media | write | `link`, `media_type`, `phone_number_id`, `to_number`; *caption, reply_to_message_id* |  | `whatsapp_business_messaging` |
| `WHATSAPP_SEND_MEDIA_BY_ID` | Send media by | write | `media_id`, `media_type`, `phone_number_id`, `to_number`; *caption, filename, reply_to_message_id* |  | `whatsapp_business_messaging` |
| `WHATSAPP_SEND_MESSAGE` | Send message | write | `phone_number_id`, `text`, `to_number`; *message_id, preview_url* |  | `whatsapp_business_messaging` |
| `WHATSAPP_SEND_TEMPLATE_MESSAGE` | Send template message | write | `phone_number_id`, `template_name`, `to_number`; *components, language_code, reply_to_message_id* |  | `whatsapp_business_messaging` |
| `WHATSAPP_SET_BUSINESS_ENCRYPTION` | Set business encryption | write | `business_public_key`, `phone_number_id` |  | `whatsapp_business_management` |
| `WHATSAPP_SET_OBO_MOBILITY_INTENT` | Set obo mobility intent | write | `body`, `waba_id` |  | `whatsapp_business_management` |
| `WHATSAPP_SET_SOLUTION_MIGRATION_INTENT` | Set solution migration intent | write | `body`, `waba_id` |  | `whatsapp_business_management` |
| `WHATSAPP_SET_TWO_STEP_VERIFICATION` | Set WhatsApp Two-Step Verification PIN | write | `phone_number_id`, `pin` |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_SHARE_PREVERIFIED_NUMBERS` | Share preverified numbers | write | `business_id`, `partner_business_id`, `preverified_id` |  | `business_management` |
| `WHATSAPP_SUBSCRIBE_APP` | Subscribe to WABA Webhooks | write | *override_callback_uri, verify_token, waba_id* |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_UNBLOCK_USERS` | Unblock WhatsApp Users | write | `phone_number_id`; *block_users, messaging_product* |  | `whatsapp_business_messaging` |
| `WHATSAPP_UNSUBSCRIBE_APP` | Unsubscribe from WABA Webhooks | write | *waba_id* |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_UPDATE_BUSINESS_ACCOUNT` | Update WhatsApp Flow Metadata | write | `flow_id`; *categories, endpoint_uri, name* |  | `whatsapp_business_management` |
| `WHATSAPP_UPDATE_BUSINESS_ACCOUNT_SETTINGS` | Update business account settings | write | `waba_id`; *name, timezone_id* |  | `whatsapp_business_management` |
| `WHATSAPP_UPDATE_BUSINESS_COMPLIANCE_INFO` | Update business compliance info | write | `customer_care_details`, `entity_name`, `entity_type`, `grievance_officer_details`, `messaging_product`, `phone_number_id`; *entity_type_custom, is_registered* |  | `whatsapp_business_management` |
| `WHATSAPP_UPDATE_BUSINESS_PROFILE` | Update WhatsApp Business Profile | write | `phone_number_id`; *about, address, description, email, messaging_product, profile_picture_handle, vertical, websites* |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_UPDATE_BUSINESS_PROFILE_NODE` | Update business profile node | write | `messaging_product`, `whatsapp_business_profile_id`; *about, address, description, email, profile_picture_handle, vertical, websites* |  | `whatsapp_business_management` |
| `WHATSAPP_UPDATE_COMMERCE_SETTINGS` | Update WhatsApp commerce settings | write | `phone_number_id`; *is_cart_enabled, is_catalog_visible* |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_UPDATE_FLOW_JSON` | Update Flow JSON | write | `flow_id`; *asset_type, file_to_upload, rawContent* |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_UPDATE_GROUP` | Update group | write | `group_id`; *description, file, messaging_product, profile_picture_file, subject* |  | `whatsapp_business_messaging` |
| `WHATSAPP_UPDATE_OFFICIAL_BUSINESS_ACCOUNT_STATUS` | Update official business account status | write | `business_website_url`, `phone_number_id`, `primary_country_of_operation`; *additional_supporting_information, parent_business_or_brand, primary_language, supporting_links* |  | `whatsapp_business_management` |
| `WHATSAPP_UPDATE_SETTINGS` | Update WhatsApp Cloud API Settings | write | `phone_number_id`; *calling, payload_encryption, storage_configuration, user_identity_change* |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_UPLOAD_FILE_CHUNK` | Upload file chunk | write | `file`, `upload_session_id`; *file_offset* |  | `whatsapp_business_management` |
| `WHATSAPP_UPLOAD_MEDIA` | Upload media | write | `file_to_upload`, `media_type`, `phone_number_id` |  | `whatsapp_business_messaging` |
| `WHATSAPP_UPSERT_MESSAGE_TEMPLATE` | Upsert message template | write | `category`, `components`, `language`, `name`; *waba_id* |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_VERIFY_CODE` | Verify Phone Number OTP Code | write | `code`, `phone_number_id` |  | `whatsapp_business_management`, `whatsapp_business_messaging` |
| `WHATSAPP_VERIFY_PREVERIFIED_NUMBER_CODE` | Verify preverified number code | write | `code`, `pre_verified_phone_number_id`; *user_agent* |  | `whatsapp_business_management` |
| `WHATSAPP_WHO_AM_I` | Who Am I | read |  |  |  |
