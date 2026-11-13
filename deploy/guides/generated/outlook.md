# Outlook: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

317 tools: 137 read, 180 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `OUTLOOK_ACCEPT_EVENT` | Accept calendar event invite | write | `event_id`; *comment, response_type, send_response, user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_ADD_EVENT_ATTACHMENT` | Add event attachment | write | `event_id`, `name`, `odata_type`; *content_bytes, item, user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_ADD_MAIL_ATTACHMENT` | Add mail attachment | write | `message_id`; *attachment, contentBytes, contentId, contentLocation, contentType, isInline, item, name and 2 more* |  | `Mail.ReadWrite` |
| `OUTLOOK_BATCH_MOVE_MESSAGES` | Batch move messages | write | `destination_id`, `message_ids`; *response_detail, user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_BATCH_UPDATE_MESSAGES` | Batch update messages | write | `updates`; *response_detail, user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_CALENDAR_CREATE_EVENT` | Create Calendar Event | write | `end_datetime`, `start_datetime`, `subject`, `time_zone`; *attendees_info, body, calendar_id, categories, importance, is_all_day, is_html, is_online_meeting and 5 more* |  | `Calendars.ReadWrite` |
| `OUTLOOK_CANCEL_CALENDAR_EVENT` | Cancel user's calendar event | write | `calendar_id`, `event_id`, `user_id`; *Comment* |  | `Calendars.ReadWrite` |
| `OUTLOOK_CANCEL_CALENDAR_GROUP_CALENDAR_EVENT` | Cancel user calendar group event | write | `calendar_group_id`, `calendar_id`, `event_id`, `user_id`; *comment* |  | `Calendars.ReadWrite` |
| `OUTLOOK_CANCEL_EVENT` | Cancel user calendar event | write | `event_id`, `user_id`; *Comment* |  | `Calendars.ReadWrite` |
| `OUTLOOK_COPY_MAIL_FOLDER` | Copy user's mail folder | write | `destination_id`, `mail_folder_id`, `user_id` |  | `Mail.ReadWrite` |
| `OUTLOOK_COPY_MESSAGE` | Copy message to folder | write | `destination_id`, `message_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_COPY_MESSAGE_FROM_CHILD_FOLDER` | Copy message from child folder | write | `child_folder_id`, `destinationId`, `mail_folder_id`, `message_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_COPY_MESSAGE_FROM_MAIL_FOLDER` | Copy user message from folder | write | `destinationId`, `mail_folder_id`, `message_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_COPY_ME_MAIL_FOLDER` | Copy child mail folder | write | `child_folder_id`, `destination_id`, `mail_folder_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_COUNT_MESSAGES` | Count messages | read | *filter, search, user_id* |  | `Mail.ReadBasic`, `Mail.Read`, `Mail.ReadWrite` |
| `OUTLOOK_CREATE_ATTACHMENT_UPLOAD_SESSION` | Create attachment upload session | write | `attachment_item`, `message_id`; *user_id* |  | `Calendars.ReadWrite`, `Mail.ReadWrite` |
| `OUTLOOK_CREATE_ATTACHMENT_UPLOAD_SESSION_IN_CHILD_FOLDER` | Create attachment upload session in child folder | write | `attachment_item`, `child_folder_id`, `mail_folder_id`, `message_id`; *user_id* |  | `Calendars.ReadWrite`, `Mail.ReadWrite` |
| `OUTLOOK_CREATE_CALENDAR` | Create calendar | write | `name`; *color, hex_color, user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_CREATE_CALENDAR_EVENT_ATTACHMENT` | Create calendar event attachment | write | `calendar_id`, `event_id`, `name`, `odata_type`; *content_bytes, item* |  | `Calendars.ReadWrite` |
| `OUTLOOK_CREATE_CALENDAR_EVENT_ATTACHMENT_UPLOAD_SESSION` | Create calendar event attachment upload session | write | `attachment_item`, `calendar_id`, `event_id`; *user_id* |  | `Calendars.ReadWrite`, `Mail.ReadWrite` |
| `OUTLOOK_CREATE_CALENDAR_EVENT_IN_CALENDAR` | Create event in specific calendar | write | `calendar_id`, `end`, `start`, `subject`, `user_id`; *allow_new_time_proposals, attendees, body, is_online_meeting, location, online_meeting_provider, recurrence* |  | `Calendars.ReadWrite` |
| `OUTLOOK_CREATE_CALENDAR_GROUP` | Create calendar group | write | `name`, `user_id` |  | `Calendars.ReadWrite` |
| `OUTLOOK_CREATE_CALENDAR_GROUP_CALENDAR_EVENT_ATTACHMENT` | Create user calendar group event attachment | write | `calendar_group_id`, `calendar_id`, `event_id`, `name`, `odata_type`; *content_bytes, item, user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_CREATE_CALENDAR_GROUP_CALENDAR_EVENT_EXTENSION` | Create event extension | write | `calendar_group_id`, `calendar_id`, `event_id`, `extension_name`; *custom_properties, user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_CREATE_CALENDAR_PERMISSION` | Create Calendar Permission | write | `calendar_group_id`, `calendar_id`, `email_address`, `role`; *is_inside_organization, is_removable, user_id* |  | `Calendars.Read`, `Calendars.ReadWrite`, `Group.Read.All`, `Group.ReadWrite.All` |
| `OUTLOOK_CREATE_CAL_GROUP_EVENT_ATTACH_UPLOAD` | Create calendar group event attachment upload session | write | `attachment_item`, `calendar_group_id`, `calendar_id`, `event_id`, `user_id` |  | `Calendars.ReadWrite`, `Mail.ReadWrite` |
| `OUTLOOK_CREATE_CONTACT` | Create contact | write | *birthday, businessPhones, categories, companyName, department, displayName, emailAddresses, givenName and 7 more* |  | `Contacts.ReadWrite` |
| `OUTLOOK_CREATE_CONTACT_FOLDER` | Create contact folder | write | `display_name`; *parent_folder_id, user_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_CREATE_CONTACT_FOLDER_CHILD_FOLDER` | Create user contact folder child folder | write | `contact_folder_id`, `display_name`, `user_id` |  | `Contacts.ReadWrite` |
| `OUTLOOK_CREATE_DRAFT` | Create email draft | write | `body`, `subject`; *attachment, bcc_recipients, cc_recipients, is_html, to_recipients, user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_DRAFT_REPLY` | Create a draft reply | write | `message_id`; *bcc_emails, cc_emails, comment, user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_EMAIL_RULE` | Create Email Rule | write | `actions`, `conditions`, `display_name`; *is_enabled, sequence, user_id* |  | `Mail.ReadWrite`, `Mail.ReadWrite.Shared`, `MailboxSettings.ReadWrite` |
| `OUTLOOK_CREATE_EVENT_ATTACHMENT` | Create user calendar event attachment | write | `event_id`, `name`, `odata_type`, `user_id`; *content_bytes, item* |  | `Calendars.ReadWrite` |
| `OUTLOOK_CREATE_EVENT_ATTACHMENT_UPLOAD_SESSION` | Create event attachment upload session | write | `attachment_item`, `event_id`; *user_id* |  | `Calendars.ReadWrite`, `Mail.ReadWrite` |
| `OUTLOOK_CREATE_FORWARD_DRAFT` | Create user mail folder message forward draft | write | `mail_folder_id`, `message_id`; *comment, to_recipients, user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_GROUP_EVENT` | Create Microsoft 365 group event | write | `end`, `group_id`, `start`, `subject`; *allow_new_time_proposals, attendees, body, is_online_meeting, location, online_meeting_provider* |  | `Group.ReadWrite.All` |
| `OUTLOOK_CREATE_MAIL_FOLDER` | Create mail folder | write | `display_name`; *is_hidden, return_existing_if_exists, user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_MAIL_FOLDER_MESSAGE` | Create message in mail folder | write | `folder_id`; *bcc_recipients, body, cc_recipients, subject, to_recipients, user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_MAIL_FOLDER_MESSAGE_ATTACHMENT` | Create mail folder message attachment | write | `mail_folder_id`, `message_id`; *attachment, contentBytes, contentId, contentLocation, contentType, isInline, item, name and 2 more* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_MAIL_FOLDER_MESSAGE_ATTACHMENT_UPLOAD_SESSION` | Create mail folder message attachment upload session | write | `attachment_item`, `mail_folder_id`, `message_id`; *user_id* |  | `Calendars.ReadWrite`, `Mail.ReadWrite` |
| `OUTLOOK_CREATE_MAIL_FOLDER_MESSAGE_RULE` | Create User Mail Folder Message Rule | write | `actions`, `display_name`, `mail_folder_id`, `sequence`; *conditions, exceptions, is_enabled, user_id* |  | `MailboxSettings.ReadWrite`, `MailboxSettings.Read` |
| `OUTLOOK_CREATE_MESSAGE_ATTACHMENT` | Create message attachment | write | `message_id`, `name`, `odata_type`; *child_folder_id, content_bytes, content_id, content_location, content_type, is_inline, item, mail_folder_id and 1 more* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_ME_CALENDAR_EVENT_EXTENSION` | Create user calendar event extension | write | `calendar_id`, `event_id`, `extension_name`; *custom_properties, user_id* |  | `Calendars.ReadWrite`, `Contacts.ReadWrite`, `Directory.AccessAsUser.All`, `Group.ReadWrite.All`, `Mail.ReadWrite`, `Organization.ReadWrite.All`, `Tasks.ReadWrite`, `User.ReadWrite` |
| `OUTLOOK_CREATE_ME_CALENDAR_PERMISSION` | Create user calendar permission | write | `calendar_id`, `emailAddress`, `role`, `user_id`; *isInsideOrganization, isRemovable* |  | `Calendars.Read`, `Calendars.ReadWrite`, `Group.Read.All`, `Group.ReadWrite.All` |
| `OUTLOOK_CREATE_ME_CHILD_FOLDER_MESSAGE` | Create message in user's child folder | write | `child_folder_id`, `mail_folder_id`, `user_id`; *bcc_recipients, body, categories, cc_recipients, importance, subject, to_recipients* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_ME_CONTACT_FOLDERS_CONTACTS` | Create user contact in folder | write | `contactFolderId`, `userId`; *birthday, businessPhones, categories, companyName, department, displayName, emailAddresses, givenName and 6 more* |  | `Contacts.ReadWrite` |
| `OUTLOOK_CREATE_ME_CONTACT_FOLDERS_CONTACTS_EXTENSIONS` | Create user contact extension | write | `contact_folder_id`, `contact_id`, `extension_name`, `user_id`; *custom_properties* |  | `Contacts.ReadWrite` |
| `OUTLOOK_CREATE_ME_CONTACT_FOLDERS_EXTENSIONS` | Create contact extension | write | `child_folder_id`, `contact_folder_id`, `contact_id`, `extension_name`; *custom_properties, user_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_CREATE_ME_CONTACT_IN_CHILD_FOLDER` | Create contact in child folder | write | `child_folder_id`, `contact_folder_id`; *birthday, businessPhones, categories, companyName, department, displayName, emailAddresses, givenName and 7 more* |  | `Contacts.ReadWrite` |
| `OUTLOOK_CREATE_ME_EVENT` | Create calendar event for user | write | `end`, `start`, `subject`; *allow_new_time_proposals, attendees, body, is_online_meeting, location, online_meeting_provider, user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_CREATE_ME_EVENT_ATTACHMENT_UPLOAD_SESSION` | Create me event attachment upload session | write | `attachment_item`, `event_id`; *user_id* |  | `Calendars.ReadWrite`, `Mail.ReadWrite` |
| `OUTLOOK_CREATE_ME_FORWARD_DRAFT` | Create forward draft | write | `message_id`; *comment, user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_ME_INFERENCE_CLASSIFICATION_OVERRIDE` | Create user Focused Inbox override | write | `classify_as`, `sender_email_address`, `user_id` |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_ME_MESSAGE_REPLY_ALL_DRAFT` | Create reply-all draft for user message | write | `message_id`, `user_id`; *comment* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_ME_REPLY_ALL_DRAFT` | Create reply-all draft for child folder message | write | `child_folder_id`, `mail_folder_id`, `message_id`; *comment, user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_REPLY_ALL_DRAFT` | Create reply-all draft in folder | write | `mail_folder_id`, `message_id`; *comment, user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_TASK` | Create To Do task | write | `title`, `todo_task_list_id`; *body, categories, dueDateTime, importance, linkedResources, reminderDateTime, user_id* |  | `Tasks.ReadWrite` |
| `OUTLOOK_CREATE_USER_CALENDAR_EVENT_ATTACHMENT` | Create user calendar event attachment | write | `calendar_id`, `event_id`, `name`, `odata_type`, `user_id`; *content_bytes, item* |  | `Calendars.ReadWrite` |
| `OUTLOOK_CREATE_USER_CALENDAR_GROUP_CALENDAR` | Create calendar in user's calendar group | write | `calendar_group_id`, `name`; *color, user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_CREATE_USER_CALENDAR_GROUP_EVENT` | Create user calendar group event | write | `calendar_group_id`, `calendar_id`, `end`, `start`, `subject`, `user_id`; *attendees, body, categories, importance, is_online_meeting, location, online_meeting_provider, show_as* |  | `Calendars.ReadWrite` |
| `OUTLOOK_CREATE_USER_CONTACTS_EXTENSIONS` | Create user contact extension | write | `contact_id`, `extension_name`, `user_id`; *custom_properties* |  | `Contacts.ReadWrite` |
| `OUTLOOK_CREATE_USER_EVENT_CALENDAR_PERMISSION` | Create calendar permission via event | write | `email_address`, `event_id`, `role`; *is_inside_organization, is_removable, user_id* |  | `Calendars.Read`, `Calendars.ReadBasic`, `Calendars.ReadWrite`, `Group.Read.All`, `Group.ReadWrite.All` |
| `OUTLOOK_CREATE_USER_MAIL_CHILD_FOLDER_MSG_EXT` | Create user message extension | write | `child_folder_id`, `extension_name`, `mail_folder_id`, `message_id`; *custom_properties, user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_USER_MAIL_FOLDERS_CHILD_FOLDERS` | Create user mail folders child folders | write | `displayName`, `mailFolderId`, `userId`; *isHidden* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_USER_MAIL_FOLDER_MESSAGE_EXTENSION` | Create user mail folder message extension | write | `extension_name`, `mail_folder_id`, `message_id`, `user_id`; *custom_properties* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_USER_MAIL_FOLDER_MESSAGE_REPLY_DRAFT` | Create reply draft for user mail folder message | write | `mail_folder_id`, `message_id`, `user_id`; *comment* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_USER_MASTER_CATEGORY` | Create user master category | write | `display_name`; *color, user_id* |  | `MailboxSettings.ReadWrite` |
| `OUTLOOK_CREATE_USER_MESSAGE` | Create user message | write | *bcc_recipients, body, categories, cc_recipients, from, importance, internet_message_headers, sender and 3 more* |  | `Mail.ReadWrite` |
| `OUTLOOK_CREATE_USER_MESSAGE_ATTACHMENT` | Create user message attachment | write | `mail_folder_id`, `message_id`, `name`, `odata_type`, `user_id`; *content_bytes, content_id, content_location, content_type, is_inline, item* |  | `Mail.ReadWrite` |
| `OUTLOOK_DECLINE_EVENT` | Decline calendar event | write | `event_id`; *comment, proposed_new_time, send_response, user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_CALENDAR` | Delete calendar | write | `calendar_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_CALENDAR_EVENT` | Delete calendar event | write | `event_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_CALENDAR_EVENT_ATTACHMENT` | Delete user calendar event attachment | write | `attachment_id`, `calendar_id`, `event_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_CALENDAR_EVENT_FROM_SPECIFIC_CALENDAR` | Delete event from specific calendar | write | `calendar_id`, `event_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_CALENDAR_FROM_GROUP` | Delete user calendar group calendar | write | `calendar_group_id`, `calendar_id`, `user_id` |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_CALENDAR_GROUP` | Delete calendar group | write | `calendar_group_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_CALENDAR_GROUP_CALENDAR_EVENT` | Delete user calendar group calendar event | write | `calendar_group_id`, `calendar_id`, `event_id`, `user_id` |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_CALENDAR_GROUP_EVENT_ATTACHMENT` | Delete calendar group event attachment | write | `attachment_id`, `calendar_group_id`, `calendar_id`, `event_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_CALENDAR_GROUP_EVENT_PERMANENTLY` | Permanently Delete Calendar Group Event | write | `calendar_group_id`, `calendar_id`, `event_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_CALENDAR_PERMANENTLY` | Permanently Delete Calendar | write | `calendar_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_CALENDAR_PERMISSION` | Delete User Calendar Group Calendar Permission | write | `calendar_group_id`, `calendar_id`, `permission_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_CHILD_CONTACT_FOLDER_PERMANENTLY` | Permanently Delete Child Contact Folder | write | `child_folder_id`, `contact_folder_id`; *user_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_DELETE_CHILD_FOLDER_MESSAGE` | Delete Child Folder Message | write | `child_folder_id`, `mail_folder_id`, `message_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_DELETE_CONTACT` | Delete Contact | write | `contact_id`; *user_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_DELETE_CONTACTS_EXTENSIONS` | Delete contact extension | write | `child_folder_id`, `contact_folder_id`, `contact_id`, `extension_id`; *user_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_DELETE_CONTACT_FOLDER` | Delete contact folder | write | `contact_folder_id`; *user_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_DELETE_CONTACT_FOLDER_CHILD_FOLDER` | Delete user contact folder child folder | write | `child_folder_id`, `contact_folder_id`; *user_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_DELETE_CONTACT_FOLDER_CONTACT` | Delete Contact from Folder | write | `contact_folder_id`, `contact_id`; *user_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_DELETE_CONTACT_FOLDER_PERMANENTLY` | Permanently Delete Contact Folder | write | `contact_folder_id`; *user_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_DELETE_CONTACT_FROM_CHILD_FOLDER_PERMANENTLY` | Permanently Delete User Contact from Child Folder | write | `child_folder_id`, `contact_folder_id`, `contact_id`, `user_id` |  | `Contacts.ReadWrite` |
| `OUTLOOK_DELETE_CONTACT_FROM_FOLDER_PERMANENTLY` | Permanently Delete Contact from Folder | write | `contact_folder_id`, `contact_id`; *user_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_DELETE_CONTACT_PERMANENTLY` | Permanently Delete Contact | write | `contact_id`; *user_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_DELETE_EMAIL_RULE` | Delete Email Rule | write | `rule_id`; *user_id* |  | `MailboxSettings.ReadWrite` |
| `OUTLOOK_DELETE_EVENT_ATTACHMENT` | Delete event attachment | write | `attachment_id`, `event_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_EVENT_EXTENSION` | Delete event extension | write | `calendar_group_id`, `calendar_id`, `event_id`, `extension_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_EVENT_PERMANENTLY` | Permanently Delete Event | write | `event_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_MAIL_FOLDER` | Delete mail folder | write | `folder_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_DELETE_MAIL_FOLDER_MESSAGE` | Delete Mail Folder Message | write | `mail_folder_id`, `message_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_DELETE_MASTER_CATEGORY` | Delete master category | write | `category_id`; *user_id* |  | `MailboxSettings.ReadWrite` |
| `OUTLOOK_DELETE_MESSAGE` | Delete Message | write | `message_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_DELETE_MESSAGE_ATTACHMENT` | Delete Message Attachment | write | `attachment_id`, `mail_folder_id`, `message_id`; *child_folder_id, user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_DELETE_MESSAGE_EXTENSION` | Delete message extension | write | `extension_id`, `message_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_DELETE_ME_CALENDAR_PERMISSION` | Delete User Calendars Calendar Permission | write | `calendar_id`, `permission_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_ME_CONTACT_EXTENSION` | Delete contact extension | write | `contact_id`, `extension-id`; *user_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_DELETE_ME_CONTACT_FOLDER_CONTACT_EXTENSION` | Delete user contact folder contact extension | write | `contact_folder_id`, `contact_id`, `extension_id`; *user_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_DELETE_ME_EVENTS_ATTACHMENTS` | Delete user event attachment | write | `attachment_id`, `event_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_ME_EVENT_EXTENSION` | Delete user event extension | write | `event_id`, `extension_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_ME_INFERENCE_CLASSIFICATION_OVERRIDE` | Delete inference classification override | write | `override_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_DELETE_ME_MAIL_FOLDERS_CHILD_FOLDERS_MESSAGES_EXTENS` | Delete user message extension | write | `child_folder_id`, `extension_id`, `mail_folder_id`, `message_id`, `user_id` |  | `Mail.ReadWrite` |
| `OUTLOOK_DELETE_ME_MAIL_FOLDERS_MESSAGES_EXTENSIONS` | Delete message extension | write | `extension_id`, `mail_folder_id`, `message_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_DELETE_ME_MAIL_FOLDER_CHILD_FOLDER` | Delete mail folder child folder | write | `child_folder_id`, `mail_folder_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_DELETE_ME_MAIL_FOLDER_CHILD_FOLDER_PERMANENTLY` | Permanently Delete Child Mail Folder | write | `child_folder_id`, `mail_folder_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_DELETE_ME_MAIL_FOLDER_MESSAGE_RULE` | Delete mail folder message rule | write | `mail_folder_id`, `rule_id`; *user_id* |  | `MailboxSettings.ReadWrite` |
| `OUTLOOK_DELETE_ME_MESSAGES_ATTACHMENTS` | Delete Message Attachment | write | `attachment_id`, `message_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_DELETE_PRIMARY_CALENDAR_PERMISSION` | Delete User Calendar Permission | write | `permission_id`, `user_id` |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_TODO_TASK` | Delete To Do task | write | `task_id`, `todo_task_list_id`; *user_id* |  | `Tasks.ReadWrite` |
| `OUTLOOK_DELETE_USER_CALENDAR_EVENT_PERMANENTLY` | Permanently Delete User Calendar Event | write | `calendar_id`, `event_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_USER_CHILD_FOLDER_CONTACT` | Delete Contact from User's Child Folder | write | `child_folder_id`, `contact_folder_id`, `contact_id`, `user_id` |  | `Contacts.ReadWrite` |
| `OUTLOOK_DELETE_USER_CHILD_FOLDER_MESSAGE_PERMANENTLY` | Delete User Child Folder Message Permanently | write | `child_folder_id`, `mail_folder_id`, `message_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_DELETE_USER_EVENT_PERMANENTLY` | Permanently Delete User Event | write | `event_id`, `user_id` |  | `Calendars.ReadWrite` |
| `OUTLOOK_DELETE_USER_MAIL_FOLDER_PERMANENTLY` | Permanently Delete User Mail Folder | write | `mail_folder_id`, `user_id` |  | `Mail.ReadWrite` |
| `OUTLOOK_DISMISS_CALENDAR_EVENT_REMINDER` | Dismiss user calendar event reminder | write | `calendar_id`, `event_id`, `user_id` |  | `Calendars.ReadWrite` |
| `OUTLOOK_DISMISS_EVENT_REMINDER` | Dismiss event reminder | write | `event_id`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_DISMISS_EVENT_REMINDER_FROM_GROUP` | Dismiss user calendar group event reminder | write | `calendar_group_id`, `calendar_id`, `event_id`, `user_id` |  | `Calendars.ReadWrite` |
| `OUTLOOK_DISMISS_USER_EVENT_REMINDER` | Dismiss user event reminder | write | `event_id`, `user_id` |  | `Calendars.ReadWrite` |
| `OUTLOOK_DOWNLOAD_GROUP_POST_ATTACHMENT` | Download group post attachment | read | `attachment_id`, `file_name`, `group_id`, `post_id`, `thread_id` |  | `Group.Read.All` |
| `OUTLOOK_DOWNLOAD_OUTLOOK_ATTACHMENT` | Download Outlook attachment | read | `attachment_id`, `file_name`, `message_id`; *user_id* |  | `Mail.Read` |
| `OUTLOOK_DOWNLOAD_PROFILE_PHOTO` | Download profile photo | read | `size` |  | `ProfilePhoto.Read.All`, `ProfilePhoto.ReadWrite.All`, `User.Read`, `User.ReadBasic.All`, `User.Read.All`, `User.ReadWrite`, `User.ReadWrite.All` |
| `OUTLOOK_FIND_MEETING_TIMES` | Find Meeting Times | read | *attendees, is_organizer_optional, location_constraint, max_candidates, meeting_duration, minimum_attendee_percentage, prefer_timezone, return_suggestion_reasons and 2 more* |  | `Calendars.Read.Shared`, `Calendars.ReadWrite.Shared` |
| `OUTLOOK_FORWARD_MESSAGE` | Forward message | write | `message_id`, `to_recipients`; *comment, user_id* |  | `Mail.Send` |
| `OUTLOOK_FORWARD_USER_CALENDAR_EVENT` | Forward user calendar event | write | `calendar_id`, `event_id`, `to_recipients`; *comment, user_id* |  | `Calendars.Read` |
| `OUTLOOK_GET_CALENDAR_EVENT` | Get event from calendar | read | `calendar_id`, `event_id`; *include_sensitivity_label, user_id* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_GET_CALENDAR_EVENT_ATTACHMENT` | Get calendar event attachment | read | `attachment_id`, `calendar_id`, `event_id`; *user_id* |  | `Calendars.Read`, `Group.Read.All`, `Mail.Read` |
| `OUTLOOK_GET_CALENDAR_FROM_EVENT` | Get calendar from event | read | `calendar_id`, `event_id`; *user_id* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_GET_CALENDAR_FROM_GROUP` | Get calendar from calendar group | read | `calendar_group_id`, `calendar_id`, `user_id` |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_GET_CALENDAR_GROUP` | Get calendar group | read | `calendar_group_id`; *expand, select, user_id* |  | `Calendars.ReadBasic`, `Calendars.Read` |
| `OUTLOOK_GET_CALENDAR_GROUP_CALENDAR_EVENT_EXTENSION` | Get event extension | read | `calendar_group_id`, `calendar_id`, `event_id`, `extension_id`; *user_id* |  | `Calendars.Read` |
| `OUTLOOK_GET_CALENDAR_GROUP_SCHEDULE` | Get User Calendar Group Schedule | read | `calendar_group_id`, `calendar_id`, `end_time`, `schedules`, `start_time`, `user_id`; *availability_view_interval* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_GET_CALENDAR_PERMISSION` | Get user calendar permission | read | `permission_id`, `user_id` |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite`, `Group.Read.All`, `Group.ReadWrite.All` |
| `OUTLOOK_GET_CALENDAR_PERMISSION_FROM_CALENDAR` | Get user calendar permission | read | `calendar_id`, `permission_id`, `user_id` |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite`, `Group.Read.All`, `Group.ReadWrite.All` |
| `OUTLOOK_GET_CALENDAR_SCHEDULE` | Get Calendar Schedule | read | `end_time`, `schedules`, `start_time`; *availability_view_interval, user_id* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_GET_CALENDAR_VIEW` | Get Calendar View | read | `end_datetime`, `start_datetime`; *calendar_id, include_sensitivity_label, orderby, page_token, select, skip, timezone, top and 1 more* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_GET_CHILD_FOLDER_MESSAGE` | Get child folder message | read | `child_folder_id`, `mail_folder_id`, `message_id`; *include_sensitivity_label, select, user_id* |  | `Mail.ReadBasic`, `Mail.Read`, `Mail.ReadWrite` |
| `OUTLOOK_GET_CHILD_FOLDER_MESSAGE_CONTENT` | Get child folder message MIME content | read | `child_folder_id`, `mail_folder_id`, `message_id`; *user_id* |  | `Mail.Read`, `Mail.ReadBasic` |
| `OUTLOOK_GET_CHILD_MAIL_FOLDER` | Get child mail folder | read | `child_folder_id`, `mail_folder_id`; *select, user_id* |  | `Mail.Read`, `Mail.ReadBasic`, `Mail.ReadWrite` |
| `OUTLOOK_GET_CONTACT_EXTENSION` | Get contact extension | read | `contact_folder_id`, `contact_id`, `extension_id`; *user_id* |  | `Calendars.Read`, `Contacts.Read`, `Directory.Read.All`, `Group.Read.All`, `Mail.Read`, `Tasks.ReadWrite`, `User.Read` |
| `OUTLOOK_GET_CONTACT_FOLDER` | Get contact folder | read | `contact_folder_id`; *select, user_id* |  | `Contacts.Read`, `Contacts.ReadWrite` |
| `OUTLOOK_GET_CONTACT_FOLDERS` | List user contact folders | read | `user_id`; *expand, filter, orderby, page_token, select, skip, top* |  | `Contacts.Read`, `Contacts.ReadWrite` |
| `OUTLOOK_GET_CONTACT_FROM_FOLDER` | Get contact from folder | read | `contact_folder_id`, `contact_id`; *expand, select, user_id* |  | `Contacts.Read`, `Contacts.ReadWrite` |
| `OUTLOOK_GET_CONTACT_PHOTO_METADATA` | Get contact photo metadata | read | `contact_id`; *user_id* |  | `Contacts.Read`, `Contacts.ReadWrite` |
| `OUTLOOK_GET_DRAFTS_MAIL_FOLDER` | Get drafts mail folder | read | *select, user_id* |  | `Mail.Read`, `Mail.ReadBasic`, `Mail.ReadWrite` |
| `OUTLOOK_GET_DRIVE_ITEM` | Get drive item | read | `drive_id`, `item_id`; *expand, if_none_match, include_deleted_items, select* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `OUTLOOK_GET_EVENT` | Get calendar event | read | `event_id`; *include_sensitivity_label, user_id* |  | `Calendars.ReadBasic`, `Calendars.Read` |
| `OUTLOOK_GET_EVENT_ATTACHMENT` | Get event attachment | read | `attachment_id`, `event_id`; *user_id* |  | `Calendars.Read`, `Group.Read.All`, `Mail.Read` |
| `OUTLOOK_GET_EVENT_CALENDAR_FROM_GROUP` | Get event calendar from calendar group | read | `calendar_group_id`, `calendar_id`, `event_id`; *user_id* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_GET_GROUP` | Get Microsoft Graph group | read | `group_id`; *select* |  | `Group.Read.All` |
| `OUTLOOK_GET_GROUP_CALENDAR` | Get Microsoft 365 group calendar | read | `group_id`; *select* |  | `Group.Read.All` |
| `OUTLOOK_GET_GROUP_EVENT` | Get Microsoft 365 group event | read | `event_id`, `group_id`; *body_content_type, expand, select, timezone* |  | `Group.Read.All` |
| `OUTLOOK_GET_GROUP_POST_ATTACHMENT` | Get group post attachment | read | `attachment_id`, `group_id`, `post_id`, `thread_id`; *expand* |  | `Group.Read.All` |
| `OUTLOOK_GET_INFERENCE_CLASSIFICATION` | Get inference classification | read | *user_id* |  | `Mail.Read` |
| `OUTLOOK_GET_MAILBOX_SETTINGS` | Get mailbox settings | read | *select, user_id* |  | `MailboxSettings.Read`, `MailboxSettings.ReadWrite` |
| `OUTLOOK_GET_MAIL_DELTA` | Get mail delta | read | *delta_token, expand, folder_id, include_sensitivity_label, select, skip_token, top, user_id* |  | `Mail.ReadBasic`, `Mail.Read`, `Mail.ReadWrite` |
| `OUTLOOK_GET_MAIL_FOLDER` | Get mail folder | read | `mail_folder_id`; *select, user_id* |  | `Mail.Read`, `Mail.ReadBasic`, `Mail.ReadWrite` |
| `OUTLOOK_GET_MAIL_FOLDER_MESSAGE` | Get mail folder message | read | `mail_folder_id`, `message_id`; *include_sensitivity_label, select, user_id* |  | `Mail.ReadBasic`, `Mail.Read`, `Mail.ReadWrite` |
| `OUTLOOK_GET_MAIL_FOLDER_MESSAGE_RULE` | Get user mail folder message rule | read | `mail_folder_id`, `rule_id`, `user_id` |  | `MailboxSettings.Read` |
| `OUTLOOK_GET_MAIL_TIPS` | Get mail tips | read | `EmailAddresses`, `MailTipsOptions`; *user_id* |  | `Mail.Read`, `Mail.Read.Shared`, `MailTips.ReadBasic.Shared` |
| `OUTLOOK_GET_MASTER_CATEGORIES` | Get master categories | read | *filter, orderby, select, skip, top, user_id* |  | `MailboxSettings.Read` |
| `OUTLOOK_GET_MASTER_CATEGORY` | Get master category | read | `category_id`; *user_id* |  | `MailboxSettings.Read` |
| `OUTLOOK_GET_MESSAGE` | Get email message | read | `message_id`; *include_sensitivity_label, select, user_id* |  | `Mail.ReadBasic`, `Mail.Read`, `Mail.ReadWrite` |
| `OUTLOOK_GET_MESSAGE_EXTENSION` | Get user message extension | read | `extension_id`, `message_id`; *user_id* |  | `Mail.Read` |
| `OUTLOOK_GET_ME_CALENDAR` | Get user's default calendar | read | *user_id* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_GET_ME_CONTACTS` | Get contact | read | `contact_id`; *user_id* |  | `Contacts.Read`, `Contacts.ReadWrite` |
| `OUTLOOK_GET_ME_CONTACTS_EXTENSIONS` | Get user contact extension | read | `contact_id`, `extension_id`, `user_id` |  | `Calendars.Read`, `Contacts.Read`, `Directory.Read.All`, `Group.Read.All`, `Mail.Read`, `Tasks.ReadWrite`, `User.Read` |
| `OUTLOOK_GET_ME_CONTACT_FOLDERS_CHILD_FOLDER` | Get user child contact folder | read | `child_folder_id`, `contact_folder_id`, `user_id`; *select* |  | `Contacts.Read`, `Contacts.ReadWrite` |
| `OUTLOOK_GET_ME_CONTACT_FROM_CHILD_FOLDER` | Get contact from child folder | read | `child_folder_id`, `contact_folder_id`, `contact_id`, `user_id`; *expand, select* |  | `Contacts.Read`, `Contacts.ReadWrite` |
| `OUTLOOK_GET_ME_CONTACT_PHOTO` | Get contact photo | read | `contact_id`; *user_id* |  | `Contacts.Read`, `Contacts.ReadWrite` |
| `OUTLOOK_GET_ME_EVENT_ATTACHMENT` | Get user event attachment | read | `attachment_id`, `event_id`, `user_id` |  | `Calendars.Read`, `Group.Read.All`, `Mail.Read` |
| `OUTLOOK_GET_ME_EVENT_CALENDAR` | Get event calendar | read | `event_id`; *user_id* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_GET_ME_MAIL_FOLDERS_MESSAGES_EXTENSIONS` | Get message extension | read | `extension_id`, `mail_folder_id`, `message_id`; *user_id* |  | `Mail.Read` |
| `OUTLOOK_GET_ME_MESSAGE_MIME_CONTENT` | Get message MIME content | read | `message_id`; *user_id* |  | `Mail.Read`, `Mail.ReadBasic` |
| `OUTLOOK_GET_ME_OUTLOOK` | Get user outlook | read | *user_id* |  | `MailboxSettings.Read`, `MailboxSettings.ReadWrite` |
| `OUTLOOK_GET_NESTED_FOLDER_MESSAGE_ATTACHMENT` | Get attachment from nested folder message | read | `attachment_id`, `child_folder_ids`, `mail_folder_id`, `message_id`, `user_id` |  | `Calendars.Read`, `Group.Read.All`, `Mail.Read` |
| `OUTLOOK_GET_PLANNER_PLAN` | Get Planner plan | read | `plan_id` |  | `Tasks.Read` |
| `OUTLOOK_GET_PROFILE` | Get Outlook profile | read | *include_proxy_addresses, select, user_id* |  | `User.Read`, `User.ReadBasic.All`, `User.Read.All` |
| `OUTLOOK_GET_SCHEDULE` | Get schedule | read | `endTime`, `schedules`, `startTime`; *availabilityViewInterval, userId* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_GET_SUPPORTED_LANGUAGES` | Get supported languages | read | *user_id* |  | `User.Read`, `User.Read.All`, `User.ReadBasic.All` |
| `OUTLOOK_GET_SUPPORTED_TIME_ZONES` | Get supported time zones | read | *TimeZoneStandard, user_id* |  | `User.Read`, `User.Read.All`, `User.ReadBasic.All` |
| `OUTLOOK_GET_TODO_TASK` | Get To Do task | read | `task_id`, `todo_task_list_id`; *expand, select, user_id* |  | `Tasks.Read`, `Tasks.ReadWrite` |
| `OUTLOOK_GET_USER_BY_EMAIL` | Get user by email | read | `email`; *select* |  | `User.Read`, `User.ReadWrite`, `User.ReadBasic.All`, `User.Read.All`, `User.ReadWrite.All`, `Directory.Read.All`, `Directory.ReadWrite.All` |
| `OUTLOOK_GET_USER_CALENDAR` | Get user's calendar | read | `calendar_id`; *select, user_id* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_GET_USER_CALENDAR_ALLOWED_SHARING_ROLES` | Get allowed calendar sharing roles for user calendar | read | `User`, `calendar_id`, `user_id` |  | `Calendars.Read`, `Calendars.ReadBasic`, `Calendars.ReadWrite` |
| `OUTLOOK_GET_USER_CALENDAR_EVENT` | Get user calendar event | read | `event_id`, `user_id`; *include_sensitivity_label* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_GET_USER_CALENDAR_GROUP_CALENDAR_PERMISSION` | Get user calendar group calendar permission | read | `calendar_group_id`, `calendar_id`, `permission_id`, `user_id` |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite`, `Group.Read.All`, `Group.ReadWrite.All` |
| `OUTLOOK_GET_USER_CALENDAR_GROUP_EVENT` | Get user calendar group event | read | `calendar_group_id`, `calendar_id`, `event_id`; *include_sensitivity_label, user_id* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_GET_USER_CALENDAR_GROUP_EVENT_ATTACHMENT` | Get user calendar group event attachment | read | `attachment_id`, `calendar_group_id`, `calendar_id`, `event_id`; *user_id* |  | `Calendars.Read`, `Group.Read.All`, `Mail.Read` |
| `OUTLOOK_GET_USER_CHILD_FOLDER_MESSAGE` | Get message from child folder | read | `child_folder_id`, `mail_folder_id`, `message_id`, `user_id`; *include_sensitivity_label, select* |  | `Mail.ReadBasic`, `Mail.Read`, `Mail.ReadWrite` |
| `OUTLOOK_GET_USER_EVENT_EXTENSION` | Get event extension | read | `event_id`, `extension_id`; *user_id* |  | `Calendars.Read`, `Contacts.Read`, `Directory.Read.All`, `Group.Read.All`, `Mail.Read`, `Tasks.ReadWrite`, `User.Read` |
| `OUTLOOK_GET_USER_MESSAGES_ATTACHMENTS` | Get user message attachment | read | `attachment_id`, `mail_folder_id`, `message_id`, `user_id`; *child_folder_id* |  | `Calendars.Read`, `Group.Read.All`, `Mail.Read` |
| `OUTLOOK_LIST_ASSIGNED_PLANNER_TASKS` | List assigned Planner tasks | read | *page_token* |  | `Tasks.Read` |
| `OUTLOOK_LIST_CALENDARS` | List Outlook calendars | read | *filter, orderby, page_token, select, skip, top, user_id* |  | `Calendars.ReadBasic`, `Calendars.ReadWrite`, `Calendars.Read`, `Calendars.Read.Shared` |
| `OUTLOOK_LIST_CALENDAR_EVENT_ATTACHMENTS` | List user calendar event attachments | read | `calendar_id`, `event_id`, `user_id`; *filter, orderby, select, skip, top* |  | `Calendars.Read` |
| `OUTLOOK_LIST_CALENDAR_GROUPS` | List Outlook calendar groups | read | *filter, orderby, page_token, select, skip, top, user_id* |  | `Calendars.ReadBasic`, `Calendars.ReadWrite`, `Calendars.Read` |
| `OUTLOOK_LIST_CALENDAR_GROUP_CALENDARS` | List calendars in calendar group | read | `calendar_group_id`; *filter, orderby, page_token, select, skip, top, user_id* |  | `Calendars.ReadBasic`, `Calendars.ReadWrite`, `Calendars.Read`, `Calendars.Read.Shared` |
| `OUTLOOK_LIST_CALENDAR_GROUP_CALENDAR_EVENTS` | List user calendar group calendar events | read | `calendar_group_id`, `calendar_id`, `user_id`; *filter, include_sensitivity_label, orderby, page_token, select, skip, timezone, top* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_LIST_CALENDAR_GROUP_EVENT_ATTACHMENTS` | List user calendar group event attachments | read | `calendar_group_id`, `calendar_id`, `event_id`; *filter, orderby, select, skip, top, user_id* |  | `Calendars.Read` |
| `OUTLOOK_LIST_CALENDAR_PERMISSIONS` | List calendar permissions | read | `calendar_group_id`, `calendar_id`; *page_token, skip, top, user_id* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite`, `Group.Read.All`, `Group.ReadWrite.All` |
| `OUTLOOK_LIST_CALENDAR_VIEW_DELTA` | List calendar view delta | read | *deltatoken, end_datetime, include_sensitivity_label, skiptoken, start_datetime, user_id* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_LIST_CHATS` | List Teams chats | read | *expand, filter, orderby, skiptoken, top, user_id* |  | `Chat.ReadBasic`, `Chat.Read`, `Chat.ReadWrite` |
| `OUTLOOK_LIST_CHAT_MESSAGES` | List chat messages | read | `chat_id`; *filter, orderby, skiptoken, top* |  | `Chat.Read`, `Chat.ReadWrite`, `Chat.Read.All`, `Chat.ReadWrite.All` |
| `OUTLOOK_LIST_CHILD_FOLDER_CONTACTS` | List user child folder contacts | read | `child_folder_id`, `contact_folder_id`, `user_id`; *filter, orderby, page_token, select, skip, top* |  | `Contacts.Read`, `Contacts.ReadWrite` |
| `OUTLOOK_LIST_CHILD_FOLDER_MESSAGES` | List child folder messages | read | `child_folder_id`, `mail_folder_id`; *filter, include_sensitivity_label, orderby, page_token, select, skip, top, user_id* |  | `Mail.ReadBasic`, `Mail.Read`, `Mail.ReadWrite` |
| `OUTLOOK_LIST_CHILD_MAIL_FOLDERS` | List child mail folders | read | `parent_folder_id`; *filter, include_hidden_folders, page_token, select, skip, top, user_id* |  | `Mail.Read`, `Mail.ReadBasic`, `Mail.ReadWrite` |
| `OUTLOOK_LIST_CONTACTS_DELTA` | List contacts delta | read | `contact_folder_id`; *delta_token, select, skip_token, user_id* |  | `Contacts.Read`, `Contacts.ReadWrite` |
| `OUTLOOK_LIST_CONTACT_FOLDERS_DELTA` | List contact folders delta | read | *deltatoken, select, skiptoken, user_id* |  | `Contacts.Read`, `Contacts.ReadWrite` |
| `OUTLOOK_LIST_CONTACT_FOLDER_CHILD_FOLDERS` | List contact folder child folders | read | `folder_id`; *expand, filter, orderby, page_token, select, skip, top, user_id* |  | `Contacts.Read`, `Contacts.ReadWrite` |
| `OUTLOOK_LIST_DIRECT_GROUP_MEMBERSHIPS` | List direct group memberships | read | *filter, orderby, page_token, search, select, top* |  | `User.Read` |
| `OUTLOOK_LIST_EMAIL_RULES` | List Email Rules | read | *skip, top, user_id* |  | `MailboxSettings.Read` |
| `OUTLOOK_LIST_EVENTS` | List events | read | *calendar_id, expand_recurring_events, filter, include_sensitivity_label, orderby, page_token, select, skip and 3 more* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_LIST_EVENT_ATTACHMENTS` | List event attachments | read | `event_id`; *filter, orderby, select, skip, top, user_id* |  | `Calendars.Read` |
| `OUTLOOK_LIST_EVENT_CALENDAR_CALENDAR_PERMISSIONS` | List event calendar permissions | read | `event_id`; *filter, orderby, page_token, select, skip, top, user_id* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite`, `Group.Read.All`, `Group.ReadWrite.All` |
| `OUTLOOK_LIST_EVENT_INSTANCES` | List event instances | read | `endDateTime`, `event_id`, `startDateTime`; *count, expand, filter, include_sensitivity_label, orderby, page_token, search, select and 3 more* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_LIST_GROUP_CALENDAR_VIEW` | List Microsoft 365 group calendar view | read | *body_content_type, end_date_time, group_id, page_token, response_timezone, select, start_date_time, top* |  | `Group.Read.All`, `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_LIST_GROUP_CONVERSATIONS` | List group conversations | read | `group_id`; *filter, select* |  | `Group-Conversation.Read.All` |
| `OUTLOOK_LIST_GROUP_CONVERSATION_THREADS` | List group conversation threads | read | `conversation_id`, `group_id`; *select* |  | `Group-Conversation.Read.All` |
| `OUTLOOK_LIST_GROUP_POST_ATTACHMENTS` | List group post attachments | read | `group_id`, `post_id`, `thread_id`; *expand, page_token, select* |  | `Group-Conversation.Read.All`, `Group-Conversation.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `OUTLOOK_LIST_GROUP_THREADS` | List Outlook group threads | read | `group_id`; *filter, orderby, select, skip, top* |  | `Group-Conversation.Read.All` |
| `OUTLOOK_LIST_GROUP_THREAD_POSTS` | List group thread posts | read | `group_id`, `thread_id`; *expand, page_token, select, top* |  | `Group-Conversation.Read.All` |
| `OUTLOOK_LIST_INFERENCE_CLASSIFICATION_OVERRIDES` | List Inference Classification Overrides | read | *skip, top, user_id* |  | `Mail.Read` |
| `OUTLOOK_LIST_LICENSE_DETAILS` | List license details | read | *select* |  | `LicenseAssignment.Read.All`, `User.Read` |
| `OUTLOOK_LIST_MAIL_FOLDERS` | List mail folders | read | *count, filter, include_hidden_folders, orderby, page_token, select, skip, top and 1 more* |  | `Mail.Read`, `Mail.ReadBasic`, `Mail.ReadWrite` |
| `OUTLOOK_LIST_MAIL_FOLDERS_DELTA` | List mail folders delta | read | *deltatoken, select, skiptoken, user_id* |  | `Mail.Read`, `Mail.ReadBasic`, `Mail.ReadWrite` |
| `OUTLOOK_LIST_MAIL_FOLDER_MESSAGES` | List mail folder messages | read | `mail_folder_id`; *expand, filter, include_sensitivity_label, orderby, page_token, select, skip, top and 1 more* |  | `Mail.ReadBasic`, `Mail.Read`, `Mail.ReadWrite` |
| `OUTLOOK_LIST_MAIL_FOLDER_MESSAGE_ATTACHMENTS` | List mail folder message attachments | read | `mail_folder_id`, `message_id`; *response_detail, skip, top, user_id* |  | `Mail.Read` |
| `OUTLOOK_LIST_MAIL_FOLDER_MESSAGE_RULES` | List mail folder message rules | read | `mail_folder_id`; *skip, top, user_id* |  | `MailboxSettings.Read` |
| `OUTLOOK_LIST_MESSAGES` | List Messages | read | *categories, conversation_id, filter, folder, from_address, has_attachments, importance, include_sensitivity_label and 19 more* |  | `Mail.ReadBasic`, `Mail.Read`, `Mail.ReadWrite`, `Mail.Read.Shared` |
| `OUTLOOK_LIST_MESSAGE_ATTACHMENTS_FROM_CHILD_FOLDER` | List message attachments from child folder | read | `child_folder_id`, `mail_folder_id`, `message_id`; *filter, select, top, user_id* |  | `Mail.Read` |
| `OUTLOOK_LIST_ME_CALENDAR_PERMISSIONS` | List user calendar permissions | read | `user_id`; *page_token, skip, top* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite`, `Group.Read.All`, `Group.ReadWrite.All` |
| `OUTLOOK_LIST_NOTIFICATION_SUBSCRIPTIONS` | List notification subscriptions | read | *page_token* |  |  |
| `OUTLOOK_LIST_ORGANIZATIONS` | List organizations | read | *select* |  | `User.Read`, `Organization.Read.All` |
| `OUTLOOK_LIST_OUTLOOK_ATTACHMENTS` | List Outlook attachments | read | `message_id`; *response_detail, top, user_id* |  | `Mail.Read` |
| `OUTLOOK_LIST_PLACES` | List places | read | `placeType`; *count, filter, page_token, roomListId, select, skip, top* |  | `Place.Read.All`, `Place.ReadWrite.All` |
| `OUTLOOK_LIST_PLANNER_PLAN_BUCKETS` | List Planner plan buckets | read | `plan_id`; *page_token* |  | `Tasks.Read`, `Group.Read.All`, `Group.ReadWrite.All`, `Tasks.ReadWrite` |
| `OUTLOOK_LIST_PRIMARY_CALENDAR_PERMISSIONS` | List primary calendar permissions | read | *page_token, skip, top, user_id* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite`, `Group.Read.All`, `Group.ReadWrite.All` |
| `OUTLOOK_LIST_REMINDERS` | List event reminders | read | `endDateTime`, `startDateTime`; *userId* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_LIST_SENT_ITEMS_MESSAGES` | List sent items messages | read | *filter, include_sensitivity_label, orderby, page_token, select, skip, top, user_id* |  | `Mail.ReadBasic`, `Mail.Read`, `Mail.ReadWrite` |
| `OUTLOOK_LIST_TODO_TASKS` | List tasks from a To Do list | read | `todo_task_list_id`; *filter, orderby, page_token, skip, top, user_id* |  | `Tasks.Read`, `Tasks.ReadWrite` |
| `OUTLOOK_LIST_TO_DO_LISTS` | List To Do task lists | read | *filter, orderby, page_token, select, skip, top, user_id* |  | `Tasks.Read`, `Tasks.ReadWrite` |
| `OUTLOOK_LIST_USERS` | List users | read | *filter, page_token, search, select, skip, top* |  | `Directory.Read.All`, `Directory.ReadWrite.All`, `User.Read`, `User.Read.All`, `User.ReadBasic.All`, `User.ReadWrite`, `User.ReadWrite.All` |
| `OUTLOOK_LIST_USER_CALENDARS_CALENDAR_PERMISSIONS` | List user calendar permissions | read | `calendar_id`, `user_id`; *page_token, skip, top* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite`, `Group.Read.All`, `Group.ReadWrite.All` |
| `OUTLOOK_LIST_USER_CALENDARS_EVENTS` | List user calendars events | read | `calendar_id`, `user_id`; *filter, include_sensitivity_label, orderby, page_token, select, skip, timezone, top* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_LIST_USER_CALENDAR_EVENTS_ATTACHMENTS` | List user calendar event attachments | read | `event_id`; *filter, orderby, select, skip, top, user_id* |  | `Calendars.Read` |
| `OUTLOOK_LIST_USER_CALENDAR_EVENT_INSTANCES` | List calendar event instances | read | `calendar_id`, `endDateTime`, `event_id`, `startDateTime`, `user_id`; *include_sensitivity_label, page_token, skip, top* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_LIST_USER_CALENDAR_GROUPS_CALENDAR_VIEW` | Get Calendar View from User Calendar Group | read | `calendar_group_id`, `calendar_id`, `endDateTime`, `startDateTime`, `user_id`; *include_sensitivity_label, page_token, skip, top* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_LIST_USER_CALENDAR_GROUP_EVENT_INSTANCES` | List event instances | read | `calendar_group_id`, `calendar_id`, `endDateTime`, `event_id`, `startDateTime`, `user_id`; *include_sensitivity_label, page_token, skip, top* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_LIST_USER_CALENDAR_VIEW` | Get calendar view from user's calendar | read | `calendar_id`, `endDateTime`, `startDateTime`, `user_id`; *include_sensitivity_label, page_token, select, top* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_LIST_USER_CONTACTS` | List user contacts | read | `user_id`; *count, expand, filter, orderby, page_token, select, skip, top* |  | `Contacts.Read`, `Contacts.ReadWrite` |
| `OUTLOOK_MOVE_MAIL_FOLDER` | Move mail folder | write | `destination_id`, `folder_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_MOVE_MESSAGE` | Move message to folder | write | `destination_id`, `message_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_MOVE_MESSAGE_FROM_CHILD_FOLDER` | Move message from child folder | write | `child_folder_id`, `destination_id`, `mail_folder_id`, `message_id`, `user_id` |  | `Mail.ReadWrite` |
| `OUTLOOK_MOVE_MESSAGE_FROM_FOLDER` | Move message from folder | write | `destinationId`, `mail_folder_id`, `message_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_MOVE_ME_MAIL_FOLDER` | Move child mail folder | write | `child_folder_id`, `destination_id`, `mail_folder_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_PERMANENT_DELETE_MESSAGE` | Permanently Delete Message | write | `message_id`; *mail_folder_id, user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_PIN_MESSAGE` | Pin message | write | `chat_id`, `message_url`; *user_id* |  | `Chat.ReadWrite`, `Chat.ReadWrite.All`, `ChatMessage.Send` |
| `OUTLOOK_QUERY_EMAILS` | Query Emails | read | *filter, folder, include_sensitivity_label, orderby, page_token, select, skip, top and 1 more* |  | `Mail.ReadBasic`, `Mail.Read`, `Mail.ReadWrite`, `Mail.Read.Shared` |
| `OUTLOOK_REPLY_ALL` | Reply all | write | `comment`, `message_id`, `user_id` |  | `Mail.Send` |
| `OUTLOOK_REPLY_EMAIL` | Reply to Email | write | `comment`, `message_id`; *attachment, bcc_emails, cc_emails, is_html, user_id* |  | `Mail.ReadWrite`, `Mail.Send` |
| `OUTLOOK_REPORT_MESSAGE` | Report message | write | `is_message_move_requested`, `message_id`, `report_action` |  | `Mail.ReadWrite` |
| `OUTLOOK_SEARCH_EVENTS` | Search calendar events | read | *end_datetime, from_index, include_sensitivity_label, page_token, query, size, start_datetime* |  | `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_SEARCH_MESSAGES` | Search Outlook messages | read | *enable_top_results, fromEmail, from_index, hasAttachments, include_sensitivity_label, page_token, query, region and 3 more* |  | `Mail.Read`, `Acronym.Read.All`, `Bookmark.Read.All`, `Calendars.Read`, `Chat.Read`, `ExternalItem.Read.All`, `Files.Read.All`, `QnA.Read.All`, `Sites.Read.All`, `Mail.ReadBasic`, `Mail.ReadWrite` |
| `OUTLOOK_SEARCH_PEOPLE` | Search Outlook people | read | *filter, orderby, page_token, search, select, skip, top* |  | `People.Read` |
| `OUTLOOK_SEARCH_SITES` | Search SharePoint sites | read | `search` |  | `Sites.Read.All` |
| `OUTLOOK_SEND_DRAFT` | Send draft | write | `message_id`; *user_id* |  | `Mail.Send` |
| `OUTLOOK_SEND_EMAIL` | Send email | write | `body`, `subject`, `to`; *attachment, bcc_emails, cc_emails, from_address, is_html, save_to_sent_items, to_name, user_id* |  | `Mail.Send` |
| `OUTLOOK_SNOOZE_CALENDAR_GROUP_EVENT_REMINDER` | Snooze user calendar group event reminder | write | `calendar_group_id`, `calendar_id`, `event_id`, `new_reminder_time`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_SNOOZE_EVENT_REMINDER` | Snooze event reminder | write | `event_id`, `new_reminder_time`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_SNOOZE_USER_CALENDAR_EVENT_REMINDER` | Snooze user calendar event reminder | write | `calendar_id`, `event_id`, `newReminderTime`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_SNOOZE_USER_EVENT_REMINDER` | Snooze user event reminder | write | `event_id`, `new_reminder_time`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_TRANSLATE_EXCHANGE_IDS` | Translate Exchange IDs | read | `input_ids`, `source_id_type`, `target_id_type`; *user_id* |  | `User.Read`, `User.ReadBasic.All`, `User.Read.All` |
| `OUTLOOK_UPDATE_CALENDAR_EVENT` | Update calendar event | write | `event_id`; *attendees, body, categories, end_datetime, location, recurrence, show_as, start_datetime and 3 more* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_UPDATE_CALENDAR_EVENT_IN_CALENDAR` | Update event in specific calendar | write | `calendar_id`, `event_id`; *attendees, body, categories, end_datetime, importance, is_all_day, is_online_meeting, is_reminder_on and 9 more* |  | `Calendars.ReadBasic`, `Calendars.Read`, `Calendars.ReadWrite` |
| `OUTLOOK_UPDATE_CALENDAR_GROUP` | Update calendar group | write | `calendar_group_id`, `name`; *user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_UPDATE_CALENDAR_GROUPS_CALENDARS` | Update user calendar in calendar group | write | `calendar_group_id`, `calendar_id`; *color, hex_color, is_default_calendar, name, user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_UPDATE_CALENDAR_GROUPS_CALENDARS_EVENTS` | Update user calendar group event | write | `calendar_group_id`, `calendar_id`, `event_id`, `user_id`; *attendees, body, categories, end, importance, is_online_meeting, location, online_meeting_provider and 3 more* |  | `Calendars.ReadWrite` |
| `OUTLOOK_UPDATE_CALENDAR_GROUP_CALENDAR_PERMISSION` | Update Calendar Group Calendar Permission | write | `calendar_group_id`, `calendar_id`, `permission_id`, `role`; *user_id* |  | `Calendars.ReadWrite`, `Group.ReadWrite.All` |
| `OUTLOOK_UPDATE_CALENDAR_PERMISSION` | Update calendar permission | write | `calendar_id`, `permission_id`, `role`; *user_id* |  | `Calendars.ReadWrite`, `Group.ReadWrite.All` |
| `OUTLOOK_UPDATE_CHILD_FOLDER_CONTACT` | Update child folder contact | write | `child_folder_id`, `contact_folder_id`, `contact_id`; *assistantName, birthday, businessAddress, businessPhones, categories, children, companyName, department and 21 more* |  | `Contacts.ReadWrite` |
| `OUTLOOK_UPDATE_CONTACT` | Update Contact | write | `contact_id`; *birthday, business_phones, categories, company_name, department, display_name, email_addresses, given_name and 7 more* |  | `Contacts.ReadWrite` |
| `OUTLOOK_UPDATE_CONTACT_FOLDER` | Update user contact folder | write | `contact_folder_id`, `user_id`; *display_name, parent_folder_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_UPDATE_CONTACT_FOLDERS_CONTACTS` | Update contact in folder | write | `contact_folder_id`, `contact_id`; *birthday, business_address, business_phones, categories, company_name, department, display_name, email_addresses and 10 more* |  | `Contacts.ReadWrite` |
| `OUTLOOK_UPDATE_CONTACT_FOLDER_CHILD_FOLDER` | Update user contact folder child folder | write | `child_folder_id`, `contact_folder_id`, `user_id`; *display_name, parent_folder_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_UPDATE_EMAIL` | Update email message | write | `message_id`; *From, bcc_recipients, body, categories, cc_recipients, flag, importance, inference_classification and 9 more* |  | `Mail.ReadWrite`, `Mail.ReadWrite.Shared` |
| `OUTLOOK_UPDATE_EMAIL_RULE` | Update Email Rule | write | `ruleId`; *actions, conditions, displayName, isEnabled, sequence, user_id* |  | `MailboxSettings.ReadWrite` |
| `OUTLOOK_UPDATE_EVENT_EXTENSION` | Update event extension | write | `calendar_id`, `event_id`, `extension_id`, `extension_name`; *custom_properties, user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_UPDATE_EVENT_EXTENSION_IN_CALENDAR_GROUP` | Update event extension in calendar group | write | `calendar_group_id`, `calendar_id`, `event_id`, `extension_id`, `extension_name`; *custom_properties, user_id* |  | `Calendars.ReadWrite` |
| `OUTLOOK_UPDATE_INFERENCE_CLASSIFICATION` | Update Inference Classification | write | *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_UPDATE_MAILBOX_SETTINGS` | Update mailbox settings | write | *automatic_replies_setting, language, time_zone, user_id, working_hours* |  | `MailboxSettings.ReadWrite` |
| `OUTLOOK_UPDATE_MAIL_FOLDER` | Update mail folder | write | `display_name`, `mail_folder_id`; *user_id* |  | `Mail.ReadWrite` |
| `OUTLOOK_UPDATE_MASTER_CATEGORY` | Update master category | write | `category_id`, `color`; *user_id* |  | `MailboxSettings.ReadWrite` |
| `OUTLOOK_UPDATE_ME_CONTACTS_EXTENSIONS` | Update contact extension | write | `contact_folder_id`, `contact_id`, `extension_id`, `extension_name`; *custom_properties, user_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_UPDATE_TODO_CHECKLIST_ITEM` | Update To Do checklist item | write | `checklist_item_id`, `task_id`, `todo_task_list_id`; *checked_date_time, created_date_time, display_name, is_checked, user_id* |  | `Tasks.ReadWrite` |
| `OUTLOOK_UPDATE_TODO_TASK` | Update To Do task | write | `task_id`, `todo_task_list_id`; *body, categories, completedDateTime, dueDateTime, importance, reminderDateTime, status, title and 1 more* |  | `Tasks.ReadWrite` |
| `OUTLOOK_UPDATE_USER_CALENDAR` | Update user calendar | write | `user_id`; *color, name* |  | `Calendars.ReadWrite` |
| `OUTLOOK_UPDATE_USER_CALENDARS` | Update user calendar by ID | write | `calendar_id`, `user_id`; *color, is_default_calendar, name* |  | `Calendars.ReadWrite` |
| `OUTLOOK_UPDATE_USER_CALENDAR_EVENT` | Update user calendar event | write | `event_id`, `user_id`; *attendees, body, categories, end, importance, isAllDay, isOnlineMeeting, isReminderOn and 7 more* |  | `Calendars.ReadWrite` |
| `OUTLOOK_UPDATE_USER_CALENDAR_PERMISSION` | Update user calendar permission | write | `permission_id`, `role`, `user_id`; *calendar_id* |  | `Calendars.ReadWrite`, `Group.ReadWrite.All` |
| `OUTLOOK_UPDATE_USER_CHILD_FOLDER_MESSAGE` | Update user child folder message | write | `child_folder_id`, `mail_folder_id`, `message_id`, `user_id`; *bccRecipients, body, categories, ccRecipients, flag, from, importance, inferenceClassification and 8 more* |  | `Mail.ReadWrite` |
| `OUTLOOK_UPDATE_USER_CONTACTS_EXTENSIONS` | Update user contact extension | write | `child_folder_id`, `contact_folder_id`, `contact_id`, `extension_id`, `extension_name`, `user_id`; *custom_properties* |  | `Contacts.ReadWrite` |
| `OUTLOOK_UPDATE_USER_CONTACTS_EXTENSIONS_DIRECT` | Update user contact extension (v3) | write | `contact_id`, `extension-id`, `extension_name`; *custom_properties, user_id* |  | `Contacts.ReadWrite` |
| `OUTLOOK_UPDATE_USER_CONTACT_EXTENSION` | Update user contact extension | write | `contact_folder_id`, `contact_id`, `extension_id`, `extension_name`, `user_id`; *custom_properties* |  | `Contacts.ReadWrite` |
| `OUTLOOK_UPDATE_USER_EVENTS_EXTENSIONS` | Update user events extensions | write | `event_id`, `extension_id`, `extension_name`, `user_id`; *custom_properties* |  | `Calendars.ReadWrite` |
| `OUTLOOK_UPDATE_USER_INFERENCE_CLASSIFICATION_OVERRIDE` | Update User Inference Classification Override | write | `classify_as`, `override_id`, `user_id` |  | `Mail.ReadWrite` |
| `OUTLOOK_UPDATE_USER_MAIL_FOLDERS_CHILD_FOLDERS` | Update user mail folder child folder | write | `child_folder_id`, `mail_folder_id`, `user_id`; *display_name* |  | `Mail.ReadWrite` |
| `OUTLOOK_UPDATE_USER_MAIL_FOLDER_MESSAGE` | Update user mail folder message | write | `mail_folder_id`, `message_id`, `user_id`; *bccRecipients, body, categories, ccRecipients, flag, from, importance, inferenceClassification and 8 more* |  | `Mail.ReadWrite` |
| `OUTLOOK_UPDATE_USER_MAIL_FOLDER_MESSAGE_EXTENSION` | Update user message extension | write | `extension_id`, `mail_folder_id`, `message_id`, `user_id`; *custom_properties, extension_name* |  | `Mail.ReadWrite` |
| `OUTLOOK_UPDATE_USER_MAIL_FOLDER_MESSAGE_RULE` | Update user mail folder message rule | write | `mail_folder_id`, `rule_id`; *actions, conditions, display_name, exceptions, is_enabled, sequence, user_id* |  | `MailboxSettings.ReadWrite` |
| `OUTLOOK_UPDATE_USER_MESSAGE_EXTENSION` | Update user message extension | write | `extension_id`, `message_id`, `user_id`; *custom_properties, extension_name* |  | `Mail.ReadWrite` |
| `OUTLOOK_WHO_AM_I` | Who Am I | read |  |  |  |
