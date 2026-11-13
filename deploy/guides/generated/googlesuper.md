# Google Super: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

541 tools: 246 read, 295 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `GOOGLESUPER_ACL_DELETE` | Delete ACL Rule | write | `calendar_id`, `rule_id` |  | `calendar`, `calendar.acls` |
| `GOOGLESUPER_ACL_GET` | Get ACL Rule | read | `calendar_id`, `rule_id` |  | `calendar`, `calendar.readonly`, `calendar.acls`, `calendar.acls.readonly` |
| `GOOGLESUPER_ACL_INSERT` | Create ACL Rule | write | `calendar_id`, `role`, `scope`; *send_notifications* |  | `calendar`, `calendar.acls` |
| `GOOGLESUPER_ACL_LIST` | List ACL Rules | read | `calendarId`; *maxResults, pageToken, showDeleted, syncToken* |  | `calendar.acls.readonly`, `calendar.acls`, `calendar` |
| `GOOGLESUPER_ACL_PATCH` | Patch ACL Rule | write | `calendar_id`, `rule_id`; *role, scope, send_notifications* |  | `calendar`, `calendar.acls` |
| `GOOGLESUPER_ACL_UPDATE` | Update ACL Rule | write | `calendar_id`, `role`, `rule_id`; *send_notifications* |  | `calendar`, `calendar.acls` |
| `GOOGLESUPER_ACL_WATCH` | Watch ACL Changes | write | `address`, `calendarId`, `id`; *params, token, type* |  | `calendar`, `calendar.acls`, `calendar.acls.readonly` |
| `GOOGLESUPER_ADD_ENRICHMENT` | Add Enrichment | write | `albumId`, `albumPosition`, `newEnrichmentItem` |  | `photoslibrary.appendonly`, `photoslibrary.edit.appcreateddata` |
| `GOOGLESUPER_ADD_LABEL_TO_EMAIL` | Modify email labels | write | `message_id`; *add_label_ids, remove_label_ids, user_id* |  | `https://mail.google.com/`, `gmail.modify` |
| `GOOGLESUPER_ADD_OR_REMOVE_TO_CUSTOMER_LIST` | Add or remove to customer list | write | `emails`, `resource_name`; *customer_id, operation* |  |  |
| `GOOGLESUPER_ADD_PARENT` | Insert File Parent (v2) | write | `fileId`, `id`; *enforceSingleParent, supportsAllDrives, supportsTeamDrives* |  | `docs`, `drive`, `drive.appdata`, `drive.file` |
| `GOOGLESUPER_ADD_PROPERTY` | Insert Property (v2 API) | write | `fileId`, `property_key`, `property_value`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, quotaUser and 4 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.metadata` |
| `GOOGLESUPER_ADD_SHEET` | Add Sheet to Existing Spreadsheet | write | `spreadsheet_id`; *force_unique, properties, title* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_AGGREGATE_COLUMN_DATA` | Aggregate Column Data | read | `operation`, `sheet_name`, `spreadsheet_id`, `target_column`; *additional_filters, case_sensitive, has_header_row, percentage_total, search_column, search_value* |  | `spreadsheets`, `spreadsheets.readonly`, `drive`, `drive.file`, `drive.readonly` |
| `GOOGLESUPER_APPEND_DIMENSION` | Append Dimension | write | `dimension`, `length`, `sheet_id`, `spreadsheet_id`; *include_spreadsheet_in_response, response_include_grid_data, response_ranges* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_ARCHIVE_CUSTOM_DIMENSION` | Archive Custom Dimension | write | `name` |  | `analytics.edit` |
| `GOOGLESUPER_AUTOCOMPLETE` | Autocomplete Place Predictions | read | `input`; *includePureServiceAreaBusinesses, includeQueryPredictions, includedPrimaryTypes, includedRegionCodes, inputOffset, languageCode, locationBias, locationRestriction and 3 more* |  | `cloud-platform` |
| `GOOGLESUPER_AUTO_RESIZE_DIMENSIONS` | Auto-Resize Rows or Columns | write | `dimension`, `end_index`, `spreadsheet_id`, `start_index`; *sheet_id, sheet_name* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_BATCH_ADD_MEDIA_ITEMS` | Batch Add Media Items | write | `albumId`, `mediaItemIds` |  | `photoslibrary.appendonly`, `photoslibrary.edit.appcreateddata` |
| `GOOGLESUPER_BATCH_CLEAR_VALUES_BY_DATA_FILTER` | Batch Clear Values By Data Filter | write | `dataFilters`, `spreadsheetId` |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_BATCH_CREATE_MEDIA_ITEMS` | Batch Create Media Items | write | *albumId, albumPosition, files, media_files, urls* |  | `photoslibrary` |
| `GOOGLESUPER_BATCH_DELETE_MESSAGES` | Batch delete Gmail messages | write | `messageIds`; *userId* |  | `https://mail.google.com/` |
| `GOOGLESUPER_BATCH_EVENTS` | Batch Events | write | `operations`; *fail_fast* |  | `calendar`, `calendar.app.created`, `calendar.events`, `calendar.events.owned` |
| `GOOGLESUPER_BATCH_EXECUTE` | Batch Execute Google Tasks Operations | write | `requests`; *max_requests* |  | `tasks` |
| `GOOGLESUPER_BATCH_GET` | Batch get spreadsheet | read | `spreadsheet_id`; *dateTimeRenderOption, empty_strings_filtered, majorDimension, ranges, valueRenderOption* |  | `spreadsheets.readonly`, `spreadsheets`, `drive.readonly`, `drive.file`, `drive` |
| `GOOGLESUPER_BATCH_GET_MEDIA_ITEMS` | Batch Get Media Items | read | `mediaItemIds` |  | `photoslibrary.readonly.appcreateddata` |
| `GOOGLESUPER_BATCH_GET_SPREADSHEET_VALUES` | Batch Get Spreadsheet Values | read | `ranges`, `spreadsheet_id`; *date_time_render_option, major_dimension, value_render_option* |  | `drive`, `drive.readonly` |
| `GOOGLESUPER_BATCH_MODIFY_MESSAGES` | Batch modify Gmail messages | write | `messageIds`; *addLabelIds, removeLabelIds, userId* |  | `https://mail.google.com/`, `gmail.modify` |
| `GOOGLESUPER_BATCH_MUTATE` | Batch Mutate Google Ads Resources | write | `mutate_operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_BATCH_RUN_PIVOT_REPORTS` | Batch Run Pivot Reports | read | `property`, `requests` |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_BATCH_RUN_REPORTS` | Batch Run Reports | read | `property`, `requests`; *unwrapped_filters* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_BATCH_UPDATE_MEMBERS` | Batch Update Google Meet Members | write | `parent`, `requests`; *updateMask* |  | `meetings.space.created` |
| `GOOGLESUPER_BATCH_UPDATE_SPREADSHEET_VALUES` | Batch Update Spreadsheet Values | write | `data`, `spreadsheet_id`, `value_input_option`; *include_values_in_response, response_date_time_render_option, response_value_render_option* |  | `drive` |
| `GOOGLESUPER_BATCH_UPDATE_VALUES_BY_DATA_FILTER` | Batch Update Values by Data Filter | write | `data`, `spreadsheetId`, `valueInputOption`; *includeValuesInResponse, responseDateTimeRenderOption, responseValueRenderOption* |  | `spreadsheets`, `drive`, `drive.file` |
| `GOOGLESUPER_CALENDARS_DELETE` | Delete Calendar | write | `calendar_id` |  | `calendar`, `calendar.app.created`, `calendar.calendars` |
| `GOOGLESUPER_CALENDARS_UPDATE` | Update Calendar | write | `calendarId`, `summary`; *description, location, timeZone* |  | `calendar`, `calendar.app.created`, `calendar.calendars` |
| `GOOGLESUPER_CALENDAR_LIST_DELETE` | Remove Calendar from List | write | `calendar_id` |  | `calendar`, `calendar.app.created`, `calendar.calendarlist` |
| `GOOGLESUPER_CALENDAR_LIST_GET` | Get Single Calendar by ID | read | `calendarId` |  | `calendar`, `calendar.readonly`, `calendar.app.created`, `calendar.calendarlist`, `calendar.calendarlist.readonly` |
| `GOOGLESUPER_CALENDAR_LIST_INSERT` | Insert Calendar into List | write | `id`; *backgroundColor, colorId, color_rgb_format, defaultReminders, foregroundColor, hidden, notificationSettings, selected and 1 more* |  | `calendar`, `calendar.calendarlist` |
| `GOOGLESUPER_CALENDAR_LIST_PATCH` | Patch Calendar List Entry | write | `calendar_id`; *backgroundColor, colorId, colorRgbFormat, defaultReminders, foregroundColor, hidden, notificationSettings, selected and 1 more* |  | `calendar`, `calendar.app.created`, `calendar.calendarlist` |
| `GOOGLESUPER_CALENDAR_LIST_UPDATE` | Update Calendar List Entry | write | `calendar_id`; *backgroundColor, colorId, colorRgbFormat, defaultReminders, foregroundColor, hidden, notificationSettings, selected and 1 more* |  | `calendar`, `calendar.app.created`, `calendar.calendarlist` |
| `GOOGLESUPER_CALENDAR_LIST_WATCH` | Watch Calendar List | write | `address`, `id`, `type`; *expiration, params, token* |  | `calendar`, `calendar.readonly`, `calendar.calendarlist`, `calendar.calendarlist.readonly` |
| `GOOGLESUPER_CHANNELS_STOP` | Stop Channel | write | `id`, `resourceId`; *token* |  | `calendar`, `calendar.acls`, `calendar.acls.readonly`, `calendar.app.created`, `calendar.calendarlist`, `calendar.calendarlist.readonly`, `calendar.events`, `calendar.events.freebusy`, `calendar.events.owned`, `calendar.events.owned.readonly`, `calendar.events.public.readonly`, `calendar.events.readonly`, `calendar.readonly`, `calendar.settings.readonly` |
| `GOOGLESUPER_CHECK_COMPATIBILITY` | Check Compatibility | read | `property`; *compatibilityFilter, dimensionFilter, dimensions, metricFilter, metrics* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_CLEAR_BASIC_FILTER` | Clear Basic Filter | write | `sheet_id`, `spreadsheet_id`; *include_spreadsheet_in_response, response_include_grid_data, response_ranges* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_CLEAR_CALENDAR` | Clear Calendar | write | `calendar_id` |  | `calendar`, `calendar.calendars` |
| `GOOGLESUPER_CLEAR_SPREADSHEET_VALUES` | Clear Spreadsheet Values | write | `range`, `spreadsheet_id` |  | `drive`, `spreadsheets` |
| `GOOGLESUPER_CLEAR_TASKS` | Clear tasks | write | `tasklist` |  | `tasks` |
| `GOOGLESUPER_CLEAR_VALUES` | Clear spreadsheet values | write | `range`, `spreadsheet_id` |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_COLORS_GET` | Get Color Definitions | read |  |  | `calendar`, `calendar.readonly`, `calendar.app.created`, `calendar.calendarlist`, `calendar.calendarlist.readonly`, `calendar.events.freebusy`, `calendar.events.owned`, `calendar.events.owned.readonly`, `calendar.events.public.readonly` |
| `GOOGLESUPER_COMPUTE_ROUTE_MATRIX` | Compute Route Matrix | read | `destinations`, `origins`; *fieldMask, languageCode, routingPreference, travelMode, units* |  | `cloud-platform` |
| `GOOGLESUPER_COPY_DOCUMENT` | Copy Google Document | write | `document_id`; *include_shared_drives, title* |  | `drive`, `drive.appdata`, `drive.file`, `drive.photos.readonly` |
| `GOOGLESUPER_COPY_FILE_ADVANCED` | Copy file with advanced options | write | `fileId`; *access_token, alt, appProperties, callback, copyRequiresWriterPermission, createdTime, description, enforceSingleParent and 24 more* |  | `drive`, `drive.appdata`, `drive.file`, `drive.photos.readonly` |
| `GOOGLESUPER_CREATE_ALBUM` | Create Album | write | `title` |  | `photoslibrary.appendonly` |
| `GOOGLESUPER_CREATE_AND_POPULATE_TABLE` | Create and Populate Table in Google Doc | write | `cellData`, `columns`, `documentId`, `rows`; *index, insertAtEndOfSegment, segmentId, tabId* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_AUDIENCE` | Create Audience | write | `description`, `displayName`, `filterClauses`, `membershipDurationDays`, `parent`; *eventTrigger, exclusionDurationMode* |  | `analytics.edit` |
| `GOOGLESUPER_CREATE_AUDIENCE_EXPORT` | Create Audience Export | write | `audience`, `parent`; *dimensions* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_CREATE_AUDIENCE_LIST` | Create Audience List | write | `audience`, `dimensions`, `parent`; *webhookNotification* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_CREATE_CALENDAR` | Create a new calendar | write | `summary`; *description, location, timezone* |  | `calendar`, `calendar.app.created`, `calendar.calendars` |
| `GOOGLESUPER_CREATE_CALLOUT_ASSET` | Create Callout Asset | write | `callout`; *customer_id, name, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_CREATE_CHANNEL_GROUP` | Create Channel Group | write | `displayName`, `groupingRule`, `parent`; *description, primary* |  | `analytics.edit` |
| `GOOGLESUPER_CREATE_CHART` | Create Chart in Google Sheets | write | `chart_type`, `data_range`, `sheet_id`, `spreadsheet_id`; *background_blue, background_green, background_red, chart_spec, legend_position, subtitle, title, x_axis_title and 1 more* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_CREATE_COMMENT` | Create Comment | write | `content`, `file_id`; *anchor, quoted_file_content_mime_type, quoted_file_content_value* |  | `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_CUSTOMER_CLIENT` | Create Customer Client | write | `customer_client`; *access_role, customer_id, email_address, validate_only* |  |  |
| `GOOGLESUPER_CREATE_CUSTOMER_LIST` | Create customer list | write | `name`; *customer_id, description* |  |  |
| `GOOGLESUPER_CREATE_CUSTOM_DIMENSION` | Create Custom Dimension | write | `displayName`, `parameterName`, `parent`, `scope`; *description, disallowAdsPersonalization* |  | `analytics.edit` |
| `GOOGLESUPER_CREATE_CUSTOM_METRIC` | Create Custom Metric | write | `displayName`, `measurementUnit`, `parameterName`, `parent`, `scope`; *description, restrictedMetricType* |  | `analytics.edit` |
| `GOOGLESUPER_CREATE_DATA_MANAGER_USER_LIST` | Create data manager user list | write | `account_types_id`, `accounts_id`, `display_name`; *account_access_status, description, ingested_user_list_info, integration_code, linked_account, login_account, membership_duration, membership_status and 3 more* |  |  |
| `GOOGLESUPER_CREATE_DATA_STREAM` | Create Data Stream | write | `parent`, `type`; *bundleId, defaultUri, displayName, packageName* |  | `analytics.edit` |
| `GOOGLESUPER_CREATE_DOCUMENT` | Create a document | write | `title`; *text* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_DOCUMENT_MARKDOWN` | Create Document Markdown | write | `title`; *image_assets, markdown_text* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_DRIVE` | Create Shared Drive | write | `name`; *backgroundImageFile, colorRgb, hidden, requestId, themeId* |  | `drive` |
| `GOOGLESUPER_CREATE_DRIVE_FILE` | Create Drive File | write | *app_properties, description, fields, mime_type, name, parent_folder_id, properties, shortcut_target_id and 2 more* |  | `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_EMAIL_DRAFT` | Create email draft | write | *attachment, bcc, body, cc, extra_recipients, include_signature, is_html, recipient_email and 3 more* |  | `https://mail.google.com/`, `gmail.compose`, `gmail.modify`, `gmail.readonly`, `contacts`, `contacts.readonly` |
| `GOOGLESUPER_CREATE_EVENT` | Create Event | write | *attendees, birthdayProperties, calendar_id, create_meeting_room, description, end_date, end_datetime, eventType and 19 more* |  | `calendar`, `calendar.app.created`, `calendar.events`, `calendar.events.owned`, `calendar.calendars`, `calendar.calendars.readonly`, `calendar.readonly` |
| `GOOGLESUPER_CREATE_EVENT_CREATE_RULE` | Create Event Create Rule | write | `destinationEvent`, `eventConditions`, `parent`; *parameterMutations, sourceCopyParameters* |  | `analytics.edit` |
| `GOOGLESUPER_CREATE_EXPANDED_DATA_SET` | Create Expanded Data Set | write | `expandedDataSet`, `parent` |  | `analytics.edit` |
| `GOOGLESUPER_CREATE_EXPERIMENT` | Create Experiment | write | `experiment`; *customer_id, validate_only* |  |  |
| `GOOGLESUPER_CREATE_EXPERIMENT_ARMS` | Create Experiment Arms | write | `experiment_arms`; *customer_id, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_CREATE_FILE` | Create File or Folder | write | *description, fields, file_to_upload, mimeType, name, parents, starred* |  | `drive`, `drive.file`, `drive.appdata` |
| `GOOGLESUPER_CREATE_FILE_FROM_TEXT` | Create a File from Text | write | `file_name`, `text_content`; *mime_type, parent_id* |  | `drive`, `drive.appdata`, `drive.file` |
| `GOOGLESUPER_CREATE_FILTER` | Create Gmail filter | write | `action`, `criteria`; *user_id* |  | `gmail.settings.basic` |
| `GOOGLESUPER_CREATE_FOLDER` | Create a folder | write | `name`; *parent_id* |  | `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_FOOTER` | Create Footer | write | `document_id`, `type`; *section_break_location* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_FOOTNOTE` | Create Footnote | write | `documentId`; *endOfSegmentLocation, location* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_GOOGLE_SHEET1` | Create a Google Sheet | write | *folder_id, folder_name, title* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_CREATE_HEADER` | Create Header | write | `documentId`; *sectionBreakLocation, text, type* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_KEY_EVENT` | Create Key Event | write | `countingMethod`, `eventName`, `parent`; *defaultValue* |  | `analytics.edit` |
| `GOOGLESUPER_CREATE_LABEL` | Create label | write | `label_name`; *background_color, label_list_visibility, message_list_visibility, text_color, user_id* |  | `https://mail.google.com/`, `gmail.labels`, `gmail.modify` |
| `GOOGLESUPER_CREATE_MEET` | Create Google Meet Space | write | *config* |  | `meetings.space.created` |
| `GOOGLESUPER_CREATE_MEMBER` | Create Member | write | `email`, `parent`; *fields, name, role* |  | `meetings.space.created` |
| `GOOGLESUPER_CREATE_NAMED_RANGE` | Create Named Range | write | `documentId`, `name`, `rangeEndIndex`, `rangeStartIndex`; *rangeSegmentId, tabId* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_PARAGRAPH_BULLETS` | Create Paragraph Bullets | write | `createParagraphBullets`, `document_id`; *tabId* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_PERMISSION` | Create Permission | write | `file_id`, `role`, `type`; *allow_file_discovery, domain, email_address, email_message, expiration_time, move_to_new_owners_root, send_notification_email, supports_all_drives and 2 more* |  | `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_PERMISSIONS_BATCH` | Create Permissions Batch | write | `permissions` |  | `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_PRESENTATION` | Create Google Slides Presentation | write | *locale, pageSize, presentationId, title* |  | `drive`, `drive.file`, `presentations` |
| `GOOGLESUPER_CREATE_PROPERTY` | Create Property | write | `displayName`, `parent`, `timeZone`; *currencyCode, industryCategory* |  | `analytics.edit` |
| `GOOGLESUPER_CREATE_RECURRING_AUDIENCE_LIST` | Create Recurring Audience List | write | `audience`, `dimensions`, `parent`; *activeDaysRemaining, webhookNotification* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_CREATE_REPLY` | Create Reply | write | `comment_id`, `content`, `file_id`; *action, fields* |  | `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_REPORT_TASK` | Create Report Task | write | `parent`; *reportDefinition* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_CREATE_ROLLUP_PROPERTY` | Create Rollup Property | write | `account`, `displayName`, `timeZone`; *sourceProperties* |  | `analytics.edit` |
| `GOOGLESUPER_CREATE_SHORTCUT_TO_FILE` | Create Shortcut to File/Folder | write | `name`, `target_id`; *ignoreDefaultVisibility, includeLabels, includePermissionsForView, keepRevisionForever, supportsAllDrives* |  | `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_SLIDES_MARKDOWN` | Create Slides from Markdown | write | `markdown_text`, `title` |  | `drive`, `drive.file`, `presentations`, `drive.readonly`, `spreadsheets`, `spreadsheets.readonly` |
| `GOOGLESUPER_CREATE_SPREADSHEET_COLUMN` | Create spreadsheet column | write | `spreadsheet_id`; *include_spreadsheet_in_response, inherit_from_before, insert_index, response_include_grid_data, response_ranges, sheet_id, sheet_name* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_CREATE_SPREADSHEET_ROW` | Create spreadsheet row | write | `spreadsheet_id`; *include_spreadsheet_in_response, inherit_from_before, insert_index, response_include_grid_data, response_ranges, sheet_id, sheet_name* |  | `spreadsheets`, `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_TAB` | Create Document Tab | write | `document_id`; *tab_properties* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_CREATE_TASK_LIST` | Create a task list | write | `tasklist_title` |  | `tasks` |
| `GOOGLESUPER_DELETE_CHART` | Delete Chart from Google Sheets | write | `chart_id`, `spreadsheet_id` |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_DELETE_CHILD` | Delete Child (v2) | write | `childId`, `folderId`; *access_token, alt, callback, enforceSingleParent, fields, key, oauth_token, prettyPrint and 4 more* |  | `docs`, `drive`, `drive.file` |
| `GOOGLESUPER_DELETE_COMMENT` | Delete Comment | write | `comment_id`, `file_id` |  | `drive`, `drive.file` |
| `GOOGLESUPER_DELETE_CONTENT_RANGE` | Delete Content Range in Document | write | `document_id`, `range` |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_DELETE_DATA_MANAGER_USER_LIST` | Delete data manager user list | write | `name`; *linked_account, login_account, validate_only* |  |  |
| `GOOGLESUPER_DELETE_DIMENSION` | Delete Dimension (Rows/Columns) | write | `spreadsheet_id`; *delete_dimension_request, dimension, end_index, include_spreadsheet_in_response, response_include_grid_data, response_ranges, sheet_id, sheet_name and 1 more* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_DELETE_DRAFT` | Delete Draft | write | `draft_id`; *user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.compose` |
| `GOOGLESUPER_DELETE_DRIVE` | Delete Shared Drive | write | `driveId`; *allowItemDeletion, useDomainAdminAccess* |  | `drive` |
| `GOOGLESUPER_DELETE_EVENT` | Delete event | write | `event_id`; *calendar_id, send_notifications, send_updates* |  | `calendar`, `calendar.app.created`, `calendar.events`, `calendar.events.owned` |
| `GOOGLESUPER_DELETE_FILTER` | Delete Gmail filter | write | `filter_id`; *user_id* |  | `gmail.settings.basic` |
| `GOOGLESUPER_DELETE_FOOTER` | Delete Footer | write | `document_id`, `footer_id`; *tab_id* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_DELETE_HEADER` | Delete Header | write | `document_id`, `header_id`; *tab_id* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_DELETE_KEY_EVENT` | Delete Key Event | write | `name` |  | `analytics.edit` |
| `GOOGLESUPER_DELETE_LABEL` | Delete label from account (permanent) | write | `label_id`; *user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.labels` |
| `GOOGLESUPER_DELETE_MEMBER` | Delete Google Meet Space Member | write | `name` |  | `meetings.space.created` |
| `GOOGLESUPER_DELETE_MESSAGE` | Delete message | write | `message_id`; *user_id* |  | `https://mail.google.com/` |
| `GOOGLESUPER_DELETE_NAMED_RANGE` | Delete Named Range | write | `deleteNamedRange`, `document_id` |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_DELETE_PARAGRAPH_BULLETS` | Delete Paragraph Bullets | write | `document_id`, `range`; *tab_id* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_DELETE_PARENT` | Delete Parent (v2) | write | `fileId`, `parentId`; *access_token, alt, callback, enforceSingleParent, fields, key, oauth_token, prettyPrint and 4 more* |  | `docs`, `drive`, `drive.file` |
| `GOOGLESUPER_DELETE_PERMISSION` | Delete Permission | write | `file_id`, `permission_id`; *supportsAllDrives, useDomainAdminAccess* |  | `drive`, `drive.file` |
| `GOOGLESUPER_DELETE_PROPERTY` | Delete Property (v2 API) | write | `fileId`, `propertyKey`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, quotaUser and 4 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.metadata` |
| `GOOGLESUPER_DELETE_REPLY` | Delete Reply | write | `comment_id`, `file_id`, `reply_id` |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly` |
| `GOOGLESUPER_DELETE_REVISION` | Delete Revision | write | `file_id`, `revision_id` |  | `drive`, `drive.appdata`, `drive.file` |
| `GOOGLESUPER_DELETE_SHEET` | Delete Sheet | write | `sheetId`, `spreadsheetId`; *includeSpreadsheetInResponse, responseIncludeGridData, responseRanges* |  | `spreadsheets`, `drive`, `drive.file` |
| `GOOGLESUPER_DELETE_TAB` | Delete Tab | write | `document_id`, `tab_id` |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_DELETE_TABLE_COLUMN` | Delete Table Column | write | `document_id`, `requests` |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_DELETE_TABLE_ROW` | Delete Table Row | write | `documentId`, `tableCellLocation` |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_DELETE_TASK` | Delete task | write | `task_id`, `tasklist_id` |  | `tasks` |
| `GOOGLESUPER_DELETE_TASK_LIST` | Delete task list | write | `tasklist_id` |  | `tasks` |
| `GOOGLESUPER_DELETE_THREAD` | Delete thread | write | `id`; *user_id* |  | `https://mail.google.com/` |
| `GOOGLESUPER_DOWNLOAD_FILE` | Download a file from Google Drive | read | `fileId`; *mime_type* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly`, `drive.appdata`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly` |
| `GOOGLESUPER_DOWNLOAD_FILE_OPERATION` | Download file via operation | read | `file_id`; *mime_type, revision_id* |  | `drive`, `drive.file`, `drive.readonly` |
| `GOOGLESUPER_DOWNLOAD_FILE_REVISION` | Download File Revision | read | `file_id`, `revision_id`; *acknowledge_abuse, file_name* |  | `drive`, `drive.readonly` |
| `GOOGLESUPER_DUPLICATE_CALENDAR` | Create a calendar | write | `summary` |  | `calendar`, `calendar.app.created`, `calendar.calendars` |
| `GOOGLESUPER_EDIT_FILE` | Edit File | write | `content`, `file_id`; *mime_type* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_EMPTY_TRASH` | Empty Trash | write | *driveId, enforceSingleParent* |  | `drive` |
| `GOOGLESUPER_END_ACTIVE_CONFERENCE` | End active conference | write | `space_name` |  | `meetings.space.created` |
| `GOOGLESUPER_END_EXPERIMENT` | End Experiment | write | `experiment`; *customer_id, validate_only* |  |  |
| `GOOGLESUPER_EVENTS_GET` | Get Event | read | `event_id`; *calendar_id, max_attendees, time_zone* |  | `calendar`, `calendar.readonly`, `calendar.events`, `calendar.events.readonly`, `calendar.app.created`, `calendar.events.freebusy`, `calendar.events.owned`, `calendar.events.owned.readonly`, `calendar.events.public.readonly` |
| `GOOGLESUPER_EVENTS_IMPORT` | Import Event | write | `end`, `iCalUID`, `start`; *attachments, attendees, calendar_id, colorId, conferenceDataVersion, description, extendedProperties, guestsCanInviteOthers and 12 more* |  | `calendar`, `calendar.events`, `calendar.app.created`, `calendar.events.owned` |
| `GOOGLESUPER_EVENTS_INSTANCES` | Get Event Instances | read | `calendarId`, `eventId`; *maxAttendees, maxResults, originalStart, pageToken, showDeleted, timeMax, timeMin, timeZone* |  | `calendar.events.public.readonly`, `calendar.events.owned.readonly`, `calendar.events.readonly`, `calendar.readonly`, `calendar.events.freebusy`, `calendar.events.owned`, `calendar.app.created`, `calendar.events`, `calendar` |
| `GOOGLESUPER_EVENTS_LIST` | List Events | read | *alwaysIncludeEmail, calendarId, eventTypes, fields, iCalUID, maxAttendees, maxResults, orderBy and 12 more* |  | `calendar`, `calendar.app.created`, `calendar.events`, `calendar.events.freebusy`, `calendar.events.owned`, `calendar.events.owned.readonly`, `calendar.events.public.readonly`, `calendar.events.readonly`, `calendar.readonly` |
| `GOOGLESUPER_EVENTS_LIST_ALL_CALENDARS` | List Events from All Calendars | read | `time_max`, `time_min`; *calendar_ids, event_types, max_results_per_calendar, q, response_detail, show_deleted, single_events* |  | `calendar`, `calendar.readonly`, `calendar.calendarlist`, `calendar.calendarlist.readonly`, `calendar.events`, `calendar.events.readonly`, `calendar.events.owned`, `calendar.events.owned.readonly`, `calendar.events.public.readonly`, `calendar.events.freebusy`, `calendar.app.created` |
| `GOOGLESUPER_EVENTS_MOVE` | Move Event | write | `calendar_id`, `destination`, `event_id`; *send_updates* |  | `calendar.events.owned`, `calendar.events`, `calendar` |
| `GOOGLESUPER_EVENTS_WATCH` | Watch Events | write | `address`, `calendarId`, `id`; *params, payload, token, type* |  | `calendar`, `calendar.app.created`, `calendar.events`, `calendar.events.freebusy`, `calendar.events.owned`, `calendar.events.owned.readonly`, `calendar.events.public.readonly`, `calendar.events.readonly`, `calendar.readonly` |
| `GOOGLESUPER_EXPORT_DOCUMENT_AS_PDF` | Export Google Doc as PDF | read | `file_id`; *filename* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly` |
| `GOOGLESUPER_EXPORT_DRIVE_FILE` | Export Drive Spreadsheet | read | `file_id`, `mime_type` |  | `drive`, `drive.readonly` |
| `GOOGLESUPER_EXPORT_GOOGLE_WORKSPACE_FILE` | Export Google Workspace file | read | `fileId`, `mimeType` |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly` |
| `GOOGLESUPER_EXPORT_SHEET_AS_PDF` | Export Sheet as PDF | read | `sheet_id`, `spreadsheet_id`; *bottom_margin, end_column_index, end_row_index, fit_to_width, left_margin, portrait, print_titles, repeat_frozen_rows and 6 more* |  | `drive`, `drive.readonly`, `spreadsheets`, `spreadsheets.readonly` |
| `GOOGLESUPER_FETCH_EMAILS` | Fetch emails | read | *ids_only, include_payload, include_spam_trash, label_ids, max_results, page_token, query, user_id and 1 more* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_FETCH_MESSAGE_BY_MESSAGE_ID` | Fetch message by message ID | read | `message_id`; *format, user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_FETCH_MESSAGE_BY_THREAD_ID` | Fetch Message by Thread ID | read | `thread_id`; *page_token, user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_FIND_EVENT` | Find event | read | *calendar_id, event_types, max_results, order_by, page_token, query, show_deleted, single_events and 3 more* |  | `calendar`, `calendar.readonly`, `calendar.events`, `calendar.events.readonly`, `calendar.events.owned`, `calendar.events.owned.readonly`, `calendar.events.public.readonly`, `calendar.events.freebusy`, `calendar.app.created` |
| `GOOGLESUPER_FIND_FILE` | Find file | read | *bare_text_query_transformed, corpora, driveId, editors_field_removed, email_query_transformed, emailaddress_field_removed, fields, folder_id and 14 more* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_FIND_FOLDER` | Find folder | read | *full_text_contains, full_text_not_contains, modified_after, name_contains, name_exact, name_not_contains, page_size, page_token and 2 more* |  | `drive`, `drive.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.file`, `drive.appdata`, `drive.meet.readonly`, `drive.photos.readonly` |
| `GOOGLESUPER_FIND_FREE_SLOTS` | Find free slots | read | *calendar_expansion_max, group_expansion_max, items, time_max, time_min, timezone* |  | `calendar`, `calendar.readonly`, `calendar.freebusy`, `calendar.events.freebusy` |
| `GOOGLESUPER_FIND_REPLACE` | Find and Replace in Spreadsheet | write | `find`, `replace`, `spreadsheetId`; *allSheets, endColumnIndex, endRowIndex, includeFormulas, matchCase, matchEntireCell, range, rangeSheetId and 5 more* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_FIND_STREET_VIEW_PANORAMAS` | Find Street View Panoramas | read | `locations`, `session`; *radius* |  | `cloud-platform` |
| `GOOGLESUPER_FORMAT_CELL` | Format cell | write | `spreadsheet_id`; *background_alpha, background_color, blue, bold, border_color, border_style, end_column_index, end_row_index and 27 more* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_FORWARD_MESSAGE` | Forward email message | write | `message_id`, `recipients`; *additional_text, bcc, cc, include_signature, user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly`, `gmail.compose`, `gmail.send` |
| `GOOGLESUPER_FREE_BUSY_QUERY` | Query Calendar Availability | read | `items`, `timeMax`, `timeMin`; *calendarExpansionMax, groupExpansionMax, timeZone* |  | `calendar`, `calendar.readonly`, `calendar.freebusy`, `calendar.events.freebusy` |
| `GOOGLESUPER_GENERATE_IDS` | Generate File IDs | read | *count, space, type* |  | `drive`, `drive.appdata`, `drive.file` |
| `GOOGLESUPER_GENERATE_KEYWORD_FORECAST_METRICS` | Generate Keyword Forecast Metrics | read | `campaign`; *currency_code, customer_id, forecast_period* |  |  |
| `GOOGLESUPER_GENERATE_KEYWORD_HISTORICAL_METRICS` | Generate Keyword Historical Metrics | read | `keywords`; *aggregate_metrics, customer_id, geo_target_constants, historical_metrics_options, include_adult_keywords, keyword_plan_network, language* |  |  |
| `GOOGLESUPER_GENERATE_KEYWORD_IDEAS` | Generate Keyword Ideas | read | *aggregate_metrics, customer_id, geo_target_constants, historical_metrics_options, include_adult_keywords, keyword_and_url_seed, keyword_annotation, keyword_plan_network and 6 more* |  |  |
| `GOOGLESUPER_GEOCODE_ADDRESS_WITH_QUERY` | Geocode Address With Query | read | `address_query`; *language_code, region_code* |  | `cloud-platform`, `maps-platform.geocode`, `maps-platform.geocode.address` |
| `GOOGLESUPER_GEOCODE_DESTINATIONS` | Geocode Destinations | read | *addressQuery, languageCode, locationQuery, place, regionCode, travelModes* |  | `maps-platform.geocode`, `cloud-platform` |
| `GOOGLESUPER_GEOCODE_LOCATION` | Reverse Geocode Location | read | `latitude`, `longitude`; *languageCode, regionCode, types* |  | `cloud-platform`, `maps-platform.geocode`, `maps-platform.geocode.location` |
| `GOOGLESUPER_GEOCODE_PLACE` | Geocode Place by ID | read | `place_id`; *languageCode, regionCode* |  | `cloud-platform`, `maps-platform.geocode`, `maps-platform.geocode.place` |
| `GOOGLESUPER_GEOCODING_API` | Geocoding API | read | *address, bounds, components, extra_computations, key, language, latlng, location_type and 3 more* |  | `cloud-platform` |
| `GOOGLESUPER_GEOLOCATE` | Geolocate Device | read | *carrier, cellTowers, considerIp, homeMobileCountryCode, homeMobileNetworkCode, radioType, wifiAccessPoints* |  | `cloud-platform` |
| `GOOGLESUPER_GET2D_TILE` | Get 2D Map Tile | read | `session`, `x`, `y`, `z`; *key, orientation* |  | `cloud-platform` |
| `GOOGLESUPER_GET3D_TILES_ROOT` | Get 3D Tiles Root | read | *key* |  | `cloud-platform` |
| `GOOGLESUPER_GET_ABOUT` | Get about | read | *fields* |  | `drive`, `drive.appdata`, `drive.file`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_GET_ACCOUNT` | Get Account | read | `name` |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_GET_ALBUM` | Get Album | read | `albumId` |  | `photoslibrary.readonly` |
| `GOOGLESUPER_GET_APP` | Get App | read | `appId` |  | `docs`, `drive`, `drive.appdata`, `drive.apps.readonly`, `drive.file`, `drive.metadata`, `drive.metadata.readonly`, `drive.readonly` |
| `GOOGLESUPER_GET_ATTACHMENT` | Get Gmail attachment | read | `attachment_id`, `file_name`, `message_id`; *user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_GET_ATTRIBUTION_SETTINGS` | Get Attribution Settings | read | `name` |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_GET_AUDIENCE` | Get Audience | read | `audienceId`, `propertyId` |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_GET_AUDIENCE_EXPORT` | Get Audience Export | read | `name` |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_GET_AUDIENCE_LIST` | Get Audience List | read | `name` |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_GET_AUTO_FORWARDING` | Get Auto-Forwarding Settings | read | *user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_GET_CALENDAR` | Get Google Calendar | read | *calendar_id* |  | `calendar`, `calendar.readonly`, `calendar.calendars`, `calendar.calendars.readonly`, `calendar.app.created` |
| `GOOGLESUPER_GET_CAMPAIGN_BY_ID` | Get Campaign By Id | read | `id`; *customer_id* |  |  |
| `GOOGLESUPER_GET_CAMPAIGN_BY_NAME` | Get campaign by name | read | `name`; *customer_id* |  |  |
| `GOOGLESUPER_GET_CHANGES_START_PAGE_TOKEN` | Get Changes Start Page Token | read | *driveId, supportsAllDrives, supportsTeamDrives, teamDriveId* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_GET_CHILD` | Get Child Reference (v2) | read | `childId`, `folderId`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, quotaUser and 3 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_GET_COMMENT` | Get Comment | read | `commentId`, `fileId`; *includeDeleted* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly` |
| `GOOGLESUPER_GET_CONDITIONAL_FORMAT_RULES` | Get conditional format rules | read | `spreadsheet_id`; *exclude_tables_in_banded_ranges, sheet_id, sheet_title* |  | `spreadsheets.readonly`, `spreadsheets`, `drive.readonly`, `drive.file`, `drive` |
| `GOOGLESUPER_GET_CONFERENCE_RECORD_BY_NAME` | Get conference record by name | read | `name` |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLESUPER_GET_CONTACTS` | Get contacts | read | *include_other_contacts, page_token, person_fields, resource_name* |  | `contacts`, `contacts.readonly`, `contacts.other.readonly` |
| `GOOGLESUPER_GET_CONVERSION_ACTION_TAG_SNIPPETS` | Get Conversion Action Tag Snippets | read | `conversion_action_id`; *customer_id* |  |  |
| `GOOGLESUPER_GET_CONVERSION_INGESTION_STATUS` | Get conversion ingestion status | read | `request_id` |  |  |
| `GOOGLESUPER_GET_CURRENT_DATE_TIME` | Get current date and time | read | *timezone* |  |  |
| `GOOGLESUPER_GET_CURRENT_USER` | Get Current Google Calendar User | read |  |  | `calendar`, `calendar.readonly`, `calendar.calendars`, `calendar.calendars.readonly` |
| `GOOGLESUPER_GET_CURRENT_WEATHER` | Get Current Weather | read | `latitude`, `longitude`; *language_code, units_system* |  | `cloud-platform` |
| `GOOGLESUPER_GET_CUSTOMER_LISTS` | Get customer lists | read | *customer_id, page_token* |  |  |
| `GOOGLESUPER_GET_CUSTOM_DIMENSION` | Get Custom Dimension | write | `name` |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_GET_DATA_MANAGER_USER_LIST` | Get data manager user list | read | `name`; *linked_account, login_account* |  |  |
| `GOOGLESUPER_GET_DATA_REDACTION_SETTINGS` | Get Data Redaction Settings | read | `name` |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_GET_DATA_RETENTION_SETTINGS` | Get Data Retention Settings | read | `name` |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_GET_DATA_SHARING_SETTINGS` | Get Data Sharing Settings | read | `name` |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_GET_DATA_STREAM` | Get Data Stream | read | `name` |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_GET_DATA_VALIDATION_RULES` | Get Data Validation Rules | read | `spreadsheetId`; *includeEmpty, ranges, sheetId, sheetTitle* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_GET_DOCS_DRIVE_FILE_METADATA` | Get Drive File Metadata | read | `fileId`; *fields, includeLabels, includePermissionsForView, supportsAllDrives* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_GET_DOCS_SPREADSHEET_VALUES` | Get Spreadsheet Values | read | `range`, `spreadsheet_id`; *date_time_render_option, major_dimension, value_render_option* |  | `drive`, `drive.readonly`, `drive.file`, `spreadsheets`, `spreadsheets.readonly` |
| `GOOGLESUPER_GET_DOCUMENT` | Get Document | read | `document_id`; *include_tabs_content, suggestions_view_mode* |  | `drive`, `drive.readonly` |
| `GOOGLESUPER_GET_DOCUMENT_BY_ID` | Get document by id | read | `id`; *includeTabsContent* |  | `documents`, `documents.readonly`, `drive`, `drive.readonly`, `drive.file` |
| `GOOGLESUPER_GET_DOCUMENT_END_INDEX` | Get Document End Index | read | `document_id`; *segment_id, tab_id* |  | `documents`, `documents.readonly`, `drive`, `drive.readonly`, `drive.file` |
| `GOOGLESUPER_GET_DOCUMENT_PLAINTEXT` | Get document plain text | read | `document_id`; *include_footers, include_footnotes, include_headers, include_tables, include_tabs_content, tab_ids, table_cell_delimiter, table_row_delimiter* |  | `documents`, `documents.readonly`, `drive`, `drive.readonly`, `drive.file` |
| `GOOGLESUPER_GET_DRAFT` | Get Draft | read | `draft_id`; *format, user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.compose`, `gmail.readonly` |
| `GOOGLESUPER_GET_DRIVE` | Get Shared Drive | read | `drive_id`; *use_domain_admin_access* |  | `drive`, `drive.readonly` |
| `GOOGLESUPER_GET_DRIVE_FILE_METADATA` | Get Drive File Metadata | read | `fileId`; *fields, includeLabels, includePermissionsForView, supportsAllDrives* |  | `drive`, `drive.metadata`, `drive.metadata.readonly`, `drive.readonly` |
| `GOOGLESUPER_GET_EXPERIMENT` | Get Experiment | read | `experiment`; *customer_id* |  |  |
| `GOOGLESUPER_GET_EXPERIMENT_ARM` | Get Experiment Arm | read | `experiment_arm`; *customer_id* |  |  |
| `GOOGLESUPER_GET_FILE_METADATA` | Get File Metadata | read | `fileId`; *fields, includeLabels, includePermissionsForView, supportsAllDrives* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_GET_FILE_PROPERTY` | Get Property (v2) | read | `fileId`, `propertyKey`; *visibility* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_GET_FILTER` | Get Gmail filter | read | `id`; *user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_GET_FORM` | Get Google Form | read | `form_id` |  | `drive`, `drive.readonly` |
| `GOOGLESUPER_GET_GLOBAL_SITE_TAG` | Get Global Site Tag | read | `name` |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_GET_GOOGLE_SIGNALS_SETTINGS` | Get Google Signals Settings | read | `name` |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_GET_KEY_EVENT` | Get Key Event | read | `name` |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_GET_LABEL` | Get label details | read | `id`; *user_id* |  | `https://mail.google.com/`, `gmail.labels`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_GET_LANGUAGE_SETTINGS` | Get Language Settings | read | *user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_GET_MEDIA_ITEM_DOWNLOAD` | Download Photos Media Item | read | `mediaItemId` |  | `photoslibrary.readonly` |
| `GOOGLESUPER_GET_MEET` | Get Meet details | read | `space_name` |  | `meetings.space.readonly`, `meetings.space.created`, `meetings.space.settings` |
| `GOOGLESUPER_GET_MEMBER` | Get Member | read | `name`; *fields* |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLESUPER_GET_METADATA` | Get Metadata | read | `name` |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_GET_PAGE_THUMBNAIL2` | Get Page Thumbnail v2 | read | `pageObjectId`, `presentationId`; *thumbnailProperties.mimeType, thumbnailProperties.thumbnailSize* |  | `drive`, `drive.file`, `drive.readonly`, `presentations`, `presentations.readonly` |
| `GOOGLESUPER_GET_PARENT` | Get Parent Reference (v2) | read | `fileId`, `parentId`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, quotaUser and 3 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_GET_PARTICIPANT` | Get Participant | read | `name` |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLESUPER_GET_PARTICIPANT_SESSION` | Get Participant Details | read | `name` |  | `meetings.space.readonly`, `meetings.space.created` |
| `GOOGLESUPER_GET_PEOPLE` | Get People | read | *other_contacts, page_size, page_token, person_fields, resource_name, sources, sync_token* |  | `contacts`, `contacts.other.readonly`, `contacts.readonly`, `userinfo.profile` |
| `GOOGLESUPER_GET_PERMISSION` | Get Permission | read | `file_id`, `permission_id`; *fields, supports_all_drives, use_domain_admin_access* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_GET_PERMISSION_ID_FOR_EMAIL` | Get Permission ID for Email | read | `email`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, quotaUser and 3 more* |  | `docs`, `drive`, `drive.appdata`, `drive.apps.readonly`, `drive.file`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_GET_PLACE_DETAILS` | Get Place Details | read | `name`; *fieldMask, languageCode, regionCode, sessionToken* |  | `cloud-platform` |
| `GOOGLESUPER_GET_PRESENTATION` | Get Google Slides Presentation | read | `presentation_id`; *fields* |  | `drive`, `drive.readonly` |
| `GOOGLESUPER_GET_PRESENTATION_PAGE_THUMBNAIL` | Get Presentation Page Thumbnail | read | `page_object_id`, `presentation_id`; *mime_type, thumbnail_size* |  | `drive.readonly`, `drive` |
| `GOOGLESUPER_GET_PROFILE` | Get Profile | read | *user_id* |  | `https://mail.google.com/`, `gmail.compose`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_GET_PROPERTY` | Get Property | read | `name` |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_GET_PROPERTY_QUOTAS_SNAPSHOT` | Get Property Quotas Snapshot | read | `property` |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_GET_RECORDING` | Get Recording | read | `name` |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLESUPER_GET_RECORDINGS_BY_CONFERENCE_RECORD_ID` | Get recordings by conference record ID | read | `conference_record_id` |  | `meetings.space.readonly`, `meetings.space.created` |
| `GOOGLESUPER_GET_RECURRING_AUDIENCE_LIST` | Get Recurring Audience List | read | `name` |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_GET_REPLY` | Get Reply | read | `commentId`, `fileId`, `replyId`; *includeDeleted* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly` |
| `GOOGLESUPER_GET_REPORT_TASK` | Get Report Task | read | `name` |  | `analytics`, `analytics.readonly` |
| `GOOGLESUPER_GET_REVISION` | Get Revision | read | `file_id`, `revision_id`; *acknowledge_abuse* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_GET_RMF_REPORT` | Get Required Google Ads Report | read | `report`; *customer_id, end_date, max_rows, optional_metrics, start_date* |  |  |
| `GOOGLESUPER_GET_ROUTE` | Get Route | read | `destination_address`, `origin_address`; *avoid_ferries, avoid_highways, avoid_tolls, compute_alternative_routes, field_mask, language_code, routing_preference, travel_mode and 1 more* |  | `cloud-platform` |
| `GOOGLESUPER_GET_SHEET_NAMES` | Get sheet names | read | `spreadsheet_id`; *exclude_hidden* |  | `spreadsheets.readonly`, `spreadsheets`, `drive.readonly`, `drive.file`, `drive` |
| `GOOGLESUPER_GET_SMART_NOTE` | Get Smart Note | read | `name` |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLESUPER_GET_SPREADSHEET` | Get Spreadsheet | read | `spreadsheet_id`; *exclude_tables_in_banded_ranges, fields, include_grid_data, ranges* |  | `drive`, `drive.readonly` |
| `GOOGLESUPER_GET_SPREADSHEET_BY_DATA_FILTER` | Get Spreadsheet by Data Filter | read | `spreadsheetId`; *dataFilters, excludeTablesInBandedRanges, includeGridData* |  | `spreadsheets`, `drive`, `drive.file` |
| `GOOGLESUPER_GET_SPREADSHEET_INFO` | Get spreadsheet info | read | *exclude_tables_in_banded_ranges, fields, include_grid_data, ranges, spreadsheet_id* |  | `spreadsheets.readonly`, `spreadsheets`, `drive.readonly`, `drive.file`, `drive` |
| `GOOGLESUPER_GET_SPREADSHEET_VALUES` | Get Spreadsheet Values | read | `range`, `spreadsheet_id`; *date_time_render_option, major_dimension, value_render_option* |  | `drive`, `drive.readonly` |
| `GOOGLESUPER_GET_STREET_VIEW_METADATA` | Get Street View Metadata | read | `session`; *latitude, longitude, pano_id, radius* |  | `cloud-platform` |
| `GOOGLESUPER_GET_STREET_VIEW_TILE` | Get Street View Tile | read | `pano_id`, `session`, `x`, `y`, `zoom` |  | `cloud-platform` |
| `GOOGLESUPER_GET_TASK` | Get Task | read | `task_id`, `tasklist_id` |  | `tasks`, `tasks.readonly` |
| `GOOGLESUPER_GET_TASK_LIST` | Get task list | read | `tasklist_id` |  | `tasks`, `tasks.readonly` |
| `GOOGLESUPER_GET_TIME_ZONE` | Get Time Zone | read | `location`, `timestamp`; *language* |  |  |
| `GOOGLESUPER_GET_TRANSCRIPT` | Get Transcript | read | `name` |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLESUPER_GET_TRANSCRIPTS_BY_CONFERENCE_RECORD_ID` | Get transcripts by conference record ID | read | `conference_record_id`; *page_size, page_token* |  | `meetings.space.readonly`, `meetings.space.created` |
| `GOOGLESUPER_GET_TRANSCRIPT_ENTRY` | Get Transcript Entry | read | `name` |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLESUPER_GET_VACATION_SETTINGS` | Get Vacation Settings | read | *user_id* |  | `gmail.settings.basic`, `gmail.modify`, `gmail.readonly`, `https://mail.google.com/` |
| `GOOGLESUPER_GOOGLE_DRIVE_DELETE_FOLDER_OR_FILE_ACTION` | Delete folder or file | write | `fileId`; *supportsAllDrives* |  | `drive`, `drive.appdata`, `drive.file` |
| `GOOGLESUPER_HIDE_DRIVE` | Hide Shared Drive | write | `drive_id` |  | `drive` |
| `GOOGLESUPER_IMPORT_MESSAGE` | Import message | write | `raw`; *deleted, internal_date_source, never_mark_spam, process_for_calendar, user_id* |  | `https://mail.google.com/`, `gmail.insert`, `gmail.modify` |
| `GOOGLESUPER_INGEST_AD_EVENTS` | Ingest ad events | write | `ad_events`, `login_account`; *encryption_info, validate_only* |  |  |
| `GOOGLESUPER_INGEST_AUDIENCE_MEMBERS` | Ingest audience members | write | `audience_members`, `destinations`; *consent, encoding, encryption_info, terms_of_service, validate_only* |  |  |
| `GOOGLESUPER_INGEST_CONVERSION_EVENTS` | Ingest conversion events | write | `destinations`, `events`; *consent, encoding, validate_only* |  |  |
| `GOOGLESUPER_INGEST_USERS` | Ingest users | write | `destinations`, `encoding`, `users`; *encryption_info, validate_only* |  |  |
| `GOOGLESUPER_INSERT_CHILD` | Insert Child Into Folder (v2) | write | `folderId`, `id`; *access_token, alt, callback, enforceSingleParent, fields, key, oauth_token, prettyPrint and 6 more* |  | `drive`, `drive.appdata`, `drive.file`, `docs` |
| `GOOGLESUPER_INSERT_DIMENSION` | Insert Dimension in Google Sheet | write | `insert_dimension`, `spreadsheet_id`; *include_spreadsheet_in_response, response_include_grid_data, response_ranges* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_INSERT_IMAGE_IN_TABLE_CELL` | Insert Image in Table Cell | write | `document_id`, `table_cell_location`, `uri`; *cell_content_index, object_size* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_INSERT_INLINE_IMAGE` | Insert Inline Image | write | `documentId`, `location`, `uri`; *objectSize* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_INSERT_MESSAGE` | Insert message into mailbox | write | `raw`; *deleted, internalDateSource, user_id* |  | `https://mail.google.com/`, `gmail.insert`, `gmail.modify` |
| `GOOGLESUPER_INSERT_PAGE_BREAK` | Insert Page Break | write | `documentId`, `insertPageBreak` |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_INSERT_TABLE_ACTION` | Insert Table in Google Doc | write | `columns`, `documentId`, `rows`; *index, insertAtEndOfSegment, segmentId, tabId* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_INSERT_TABLE_COLUMN` | Insert Table Column | write | `document_id`, `requests` |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_INSERT_TABLE_ROW` | Insert Table Row | write | `document_id`, `requests` |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_INSERT_TASK` | Insert Task | write | `tasklist_id`, `title`; *completed, deleted, due, etag, hidden, id, notes, status and 2 more* |  | `tasks` |
| `GOOGLESUPER_INSERT_TEXT_ACTION` | Insert Text into Document | write | `document_id`, `text_to_insert`; *append_to_end, insertion_index, segment_id, tab_id* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_INSERT_TEXT_IN_TABLE_CELL` | Insert Text in Table Cell | write | `column_index`, `document_id`, `row_index`, `table_start_index`, `text`; *tabId* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_LIST_ACCESSIBLE_CUSTOMERS` | List Accessible Customers | read | *customer_id* |  |  |
| `GOOGLESUPER_LIST_ACCESS_BINDINGS` | List Access Bindings | read | `parent`; *pageSize, pageToken* |  | `analytics.manage.users`, `analytics.manage.users.readonly` |
| `GOOGLESUPER_LIST_ACCESS_PROPOSALS` | List Access Proposals | read | `fileId`; *pageSize, pageToken* |  | `docs`, `drive`, `drive.file`, `drive.metadata`, `drive.metadata.readonly`, `drive.readonly` |
| `GOOGLESUPER_LIST_ACCOUNTS_V1_BETA` | List Accounts (v1beta) | read | *pageSize, pageToken, showDeleted* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_ACCOUNT_SUMMARIES` | List Account Summaries | read | *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_ADSENSE_LINKS` | List AdSense Links | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_ALBUMS` | List Albums | read | *excludeNonAppCreatedData, pageSize, pageToken* |  | `photoslibrary.readonly` |
| `GOOGLESUPER_LIST_ALL_TASKS` | List All Tasks Across All Lists | read | *completedMax, completedMin, dueMax, dueMin, max_tasks_total, showAssigned, showCompleted, showDeleted and 2 more* |  | `tasks` |
| `GOOGLESUPER_LIST_APPROVALS` | List Approvals | read | `fileId`; *pageSize, pageToken* |  | `drive`, `drive.appdata`, `drive.file`, `drive.metadata`, `drive.metadata.readonly`, `drive.readonly` |
| `GOOGLESUPER_LIST_AUDIENCES` | List Audiences | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_AUDIENCE_EXPORTS` | List Audience Exports | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_LIST_AUDIENCE_LISTS` | List Audience Lists | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_LIST_BIGQUERY_LINKS` | List BigQuery Links | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_BUILDINGS` | List Buildings | read | *customer, maxResults, pageToken* |  | `https://apps-apis.google.com/a/feeds/calendar/resource/`, `admin.directory.resource.calendar`, `admin.directory.resource.calendar.readonly` |
| `GOOGLESUPER_LIST_CALCULATED_METRICS` | List Calculated Metrics | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_CALENDARS` | List Google Calendars | read | *max_results, min_access_role, page_token, show_deleted, show_hidden, sync_token* |  | `calendar`, `calendar.readonly`, `calendar.calendarlist`, `calendar.calendarlist.readonly` |
| `GOOGLESUPER_LIST_CALENDAR_RESOURCES` | List Calendar Resources | read | *customer, maxResults, orderBy, pageToken, query* |  | `https://apps-apis.google.com/a/feeds/calendar/resource/`, `admin.directory.resource.calendar`, `admin.directory.resource.calendar.readonly` |
| `GOOGLESUPER_LIST_CHANGES` | List Changes | read | *driveId, fields, includeCorpusRemovals, includeItemsFromAllDrives, includeLabels, includePermissionsForView, includeRemoved, pageSize and 4 more* |  | `drive.metadata.readonly`, `drive`, `drive.readonly`, `drive.metadata`, `drive.file`, `drive.appdata`, `drive.meet.readonly`, `drive.photos.readonly` |
| `GOOGLESUPER_LIST_CHANNEL_GROUPS` | List Channel Groups | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_CHARTS` | List Charts in Google Sheets | read | `spreadsheet_id`; *chart_id, include_spec, sheet_id* |  | `spreadsheets.readonly`, `spreadsheets`, `drive.readonly`, `drive.file`, `drive` |
| `GOOGLESUPER_LIST_CHILDREN_V2` | List Folder Children (v2) | read | `folderId`; *access_token, alt, callback, fields, key, maxResults, oauth_token, orderBy and 7 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_LIST_COMMENTS` | List Comments | read | `fileId`; *fields, includeDeleted, pageSize, pageToken, startModifiedTime* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly` |
| `GOOGLESUPER_LIST_CONFERENCE_RECORDS` | List Conference Records | read | *filter, page_size, page_token* |  | `meetings.space.readonly`, `meetings.space.created` |
| `GOOGLESUPER_LIST_CONVERSION_EVENTS` | List Conversion Events | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_CSE_IDENTITIES` | List CSE identities | read | *page_size, page_token, user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly`, `gmail.settings.sharing` |
| `GOOGLESUPER_LIST_CSE_KEYPAIRS` | List CSE key pairs | read | *page_size, page_token, user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly`, `gmail.settings.sharing` |
| `GOOGLESUPER_LIST_CUSTOM_DIMENSIONS` | List Custom Dimensions | write | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_CUSTOM_METRICS` | List Custom Metrics | write | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_DATA_MANAGER_CONVERSION_ACTIONS` | List data manager conversion actions | read | *customer_id, name_contains, primary_for_goal, statuses, use_case* |  |  |
| `GOOGLESUPER_LIST_DATA_MANAGER_USER_LISTS` | List data manager user lists | read | `parent`; *filter, linked_account, login_account, page_size, page_token* |  |  |
| `GOOGLESUPER_LIST_DATA_STREAMS` | List DataStreams | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_DRAFTS` | List Drafts | read | *max_results, page_token, user_id, verbose* |  | `https://mail.google.com/`, `gmail.compose`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_LIST_DRIVE_CHANGES` | List Drive Changes | read | *drive_id, include_corpus_removals, include_items_from_all_drives, include_labels, include_permissions_for_view, include_removed, page_size, page_token and 3 more* |  | `drive`, `drive.metadata`, `drive.metadata.readonly`, `drive.readonly` |
| `GOOGLESUPER_LIST_DRIVE_ITEM_PERMISSIONS` | List Drive Item Permissions | read | `fileId`; *fields, includePermissionsForView, pageSize, pageToken, supportsAllDrives, useDomainAdminAccess* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_LIST_DV360_AD_LINKS` | List Display & Video 360 Advertiser Links | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_DV360_LINK_PROPOSALS` | List DisplayVideo360 Advertiser Link Proposals | read | `parent`; *pageSize, pageToken* |  | `analytics.edit`, `analytics.readonly` |
| `GOOGLESUPER_LIST_EVENT_CREATE_RULES` | List Event Create Rules | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_EVENT_EDIT_RULES` | List Event Edit Rules | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_EXPANDED_DATA_SETS` | List Expanded Data Sets | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_EXPERIMENTS` | List Experiments | read | *customer_id, experiment_type, page_token, status* |  |  |
| `GOOGLESUPER_LIST_EXPERIMENT_ARMS` | List Experiment Arms | read | *control, customer_id, experiment, page_token* |  |  |
| `GOOGLESUPER_LIST_FILE_LABELS` | List File Labels | read | `file_id`; *max_results, page_token* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.readonly` |
| `GOOGLESUPER_LIST_FILE_PROPERTIES` | List Properties (v2 API) | read | `fileId`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, quotaUser and 3 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_LIST_FILTERS` | List Gmail filters | read | *user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly`, `gmail.settings.basic` |
| `GOOGLESUPER_LIST_FIREBASE_LINKS` | List Firebase Links | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_FORM_RESPONSES` | List Google Form Responses | read | `form_id`; *filter, page_size, page_token* |  | `drive` |
| `GOOGLESUPER_LIST_FORWARDING_ADDRESSES` | List forwarding addresses | read | *user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly`, `gmail.settings.basic` |
| `GOOGLESUPER_LIST_GOOGLE_ADS_LINKS` | List Google Ads Links | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_HISTORY` | List Gmail history | read | `start_history_id`; *history_types, label_id, max_results, page_token, user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_LIST_KEY_EVENTS` | List Key Events | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_LABELS` | List Gmail labels | read | *include_details, user_id* |  | `https://mail.google.com/`, `gmail.labels`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_LIST_LABEL_DEFINITIONS` | List Label Definitions | read | *customer, language_code, minimum_role, page_size, page_token, published_only* |  | `drive.labels.readonly`, `drive.labels`, `drive.admin.labels.readonly`, `drive.admin.labels` |
| `GOOGLESUPER_LIST_MEASUREMENT_PROTOCOL_SECRETS` | List Measurement Protocol Secrets | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_MEDIA_ITEMS` | List Media Items (App-Created Only) | read | *pageSize, pageToken* |  | `photoslibrary.readonly.appcreateddata` |
| `GOOGLESUPER_LIST_MEMBERS` | List Members | read | `parent`; *fields, page_size, page_token* |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLESUPER_LIST_PARTICIPANTS` | List Participants | read | `parent`; *filter, page_size, page_token* |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLESUPER_LIST_PARTICIPANT_SESSIONS` | List Participant Sessions | read | `conference_record_id`, `participant_id`; *filter, page_size, page_token* |  | `meetings.space.readonly`, `meetings.space.created` |
| `GOOGLESUPER_LIST_PERMISSIONS` | List Permissions | read | `fileId`; *fields, includePermissionsForView, pageSize, pageToken, supportsAllDrives, useDomainAdminAccess* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_LIST_PROPERTIES_FILTERED` | List Property | read | `filter`; *pageSize, pageToken, showDeleted* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_RECORDINGS` | List Recordings | read | `parent`; *page_size, page_token* |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLESUPER_LIST_RECURRING_AUDIENCE_LISTS` | List Recurring Audience Lists | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_LIST_REPLIES` | List Replies to Comment | read | `comment_id`, `file_id`; *fields, include_deleted, page_size, page_token* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly` |
| `GOOGLESUPER_LIST_REPORTING_DATA_ANNOTATIONS` | List Reporting Data Annotations | read | `parent`; *filter, pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_REPORT_TASKS` | List Report Tasks | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_LIST_REVISIONS` | List File Revisions | read | `fileId`; *pageSize, pageToken, supportsAllDrives* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_LIST_SEARCH_ADS360_LINKS` | List Search Ads 360 Links | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_SEND_AS` | List send-as aliases | read | *user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_LIST_SHARED_DRIVES` | List Shared Drives | read | *fields, pageSize, pageToken, q, useDomainAdminAccess* |  | `drive`, `drive.readonly` |
| `GOOGLESUPER_LIST_SK_AD_NETWORK_CONVERSION_VALUE_SCHEMAS` | List SKAdNetwork Conversion Value Schemas | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_SMART_NOTES` | List Smart Notes | read | `parent`; *page_size, page_token* |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLESUPER_LIST_SMIME_INFO` | List S/MIME configs | read | `send_as_email`; *user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly`, `gmail.settings.basic`, `gmail.settings.sharing` |
| `GOOGLESUPER_LIST_SPREADSHEET_CHARTS` | Get Charts from Spreadsheet | read | `spreadsheet_id` |  | `drive`, `drive.readonly`, `drive.file`, `spreadsheets`, `spreadsheets.readonly` |
| `GOOGLESUPER_LIST_SUBPROPERTY_EVENT_FILTERS` | List Subproperty Event Filters | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_SUBPROPERTY_SYNC_CONFIGS` | List Subproperty Sync Configs | read | `parent`; *pageSize, pageToken* |  | `analytics.readonly`, `analytics.edit` |
| `GOOGLESUPER_LIST_SUB_ACCOUNTS` | List Sub Accounts | read | *customer_id, page_token* |  |  |
| `GOOGLESUPER_LIST_TASKS` | List Tasks | read | `tasklist_id`; *completedMax, completedMin, dueMax, dueMin, maxResults, pageToken, showCompleted, showDeleted and 2 more* |  | `tasks` |
| `GOOGLESUPER_LIST_TASK_LISTS` | List task lists | read | *maxResults, pageToken* |  | `tasks`, `tasks.readonly` |
| `GOOGLESUPER_LIST_THREADS` | List threads | read | *max_results, page_token, query, user_id, verbose* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_LIST_TRANSCRIPT_ENTRIES` | List Transcript Entries | read | `conference_record_id`, `transcript_id`; *page_size, page_token* |  | `meetings.space.created`, `meetings.space.readonly` |
| `GOOGLESUPER_LOOKUP_AERIAL_VIDEO` | Lookup Aerial Video | read | *address, videoId* |  | `cloud-platform` |
| `GOOGLESUPER_LOOKUP_SPREADSHEET_ROW` | Look up spreadsheet row | read | `query`, `spreadsheet_id`; *case_sensitive, date_time_render_option, normalize_whitespace, range, value_render_option* |  | `spreadsheets.readonly`, `spreadsheets`, `drive.readonly`, `drive.file`, `drive` |
| `GOOGLESUPER_MAPS_EMBED_API` | Embed Google Map | read | `mode`; *api_key, directions_params, place_params, search_params, streetview_params, view_params* |  |  |
| `GOOGLESUPER_MODIFY_FILE_LABELS` | Modify File Labels | write | `file_id`, `label_modifications`; *kind* |  | `drive`, `drive.file`, `drive.metadata` |
| `GOOGLESUPER_MODIFY_THREAD_LABELS` | Modify thread labels | write | `thread_id`; *add_label_ids, remove_label_ids, user_id* |  | `https://mail.google.com/`, `gmail.modify` |
| `GOOGLESUPER_MOVE_CHART` | Move or Resize Chart in Google Sheets | write | `chart_id`, `spreadsheet_id`; *anchor_column_index, anchor_row_index, anchor_sheet_id, height_pixels, new_sheet, offset_x_pixels, offset_y_pixels, target_sheet_id and 1 more* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_MOVE_FILE` | Move File | write | `file_id`; *add_parents, include_labels, include_permissions_for_view, keep_revision_forever, ocr_language, remove_parents, supports_all_drives, use_content_as_indexable_text* |  | `drive`, `drive.file`, `drive.appdata`, `drive.metadata`, `drive.scripts` |
| `GOOGLESUPER_MOVE_SPREADSHEET_TO_FOLDER` | Move Spreadsheet to Folder | write | `destination_folder_id`, `spreadsheet_id` |  | `drive`, `drive.file`, `drive.metadata` |
| `GOOGLESUPER_MOVE_TASK` | Move Task | write | `task`, `tasklist`; *destinationTasklist, parent, previous* |  | `tasks` |
| `GOOGLESUPER_MOVE_THREAD_TO_TRASH` | Trash thread | write | `thread_id`; *user_id* |  | `https://mail.google.com/`, `gmail.modify` |
| `GOOGLESUPER_MOVE_TO_TRASH` | Move to Trash | write | `message_id`; *user_id* |  | `https://mail.google.com/`, `gmail.modify` |
| `GOOGLESUPER_MUTATE_ADS` | Mutate Ads | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_AD_GROUPS` | Mutate Ad Groups | write | `operations`; *customer_id, partial_failure, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_AD_GROUP_ADS` | Mutate Ad Group Ads | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_AD_GROUP_ASSETS` | Mutate Ad Group Assets | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_AD_GROUP_BID_MODIFIERS` | Mutate Ad Group Bid Modifiers | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_AD_GROUP_CRITERIA` | Mutate Ad Group Criteria | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_ASSETS` | Mutate Assets | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_ASSET_GROUPS` | Mutate Asset Groups | write | `operations`; *customer_id, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_ASSET_GROUP_ASSETS` | Mutate Asset Group Assets | write | `operations`; *customer_id, partial_failure, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_ASSET_SETS` | Mutate Asset Sets | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_AUDIENCES` | Mutate Audiences | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_BIDDING_STRATEGIES` | Mutate Bidding Strategies | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_CAMPAIGNS` | Mutate Campaigns | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_CAMPAIGNS_V2` | Mutate Campaigns V2 | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_CAMPAIGN_ASSETS` | Mutate Campaign Assets | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_CAMPAIGN_BUDGETS` | Mutate Campaign Budgets | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_CAMPAIGN_CRITERIA` | Mutate Campaign Criteria | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_CAMPAIGN_LABELS` | Mutate Campaign Labels | write | `operations`; *customer_id, partial_failure, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_CONDITIONAL_FORMAT_RULES` | Mutate conditional format rules | write | `operation`, `sheet_id`, `spreadsheet_id`; *index, new_index, rule* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_MUTATE_CONVERSION_ACTIONS` | Mutate Conversion Actions | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_CUSTOMER` | Mutate Customer | write | `operation`; *customer_id, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_CUSTOMER_ASSETS` | Mutate Customer Assets | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_CUSTOM_AUDIENCES` | Mutate Custom Audiences | write | `operations`; *customer_id, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_KEYWORD_PLANS` | Mutate Keyword Plans | write | `operations`; *customer_id, partial_failure, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_KEYWORD_PLAN_AD_GROUPS` | Mutate Keyword Plan Ad Groups | write | `operations`; *customer_id, partial_failure, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_KEYWORD_PLAN_AD_GROUP_KEYWORDS` | Mutate Keyword Plan Ad Group Keywords | write | `operations`; *customer_id, partial_failure, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_KEYWORD_PLAN_CAMPAIGNS` | Mutate Keyword Plan Campaigns | write | `operations`; *customer_id, partial_failure, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_KEYWORD_PLAN_CAMPAIGN_KEYWORDS` | Mutate Keyword Plan Campaign Keywords | write | `operations`; *customer_id, partial_failure, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_LABELS` | Mutate Labels | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_PORTFOLIO_BIDDING_STRATEGIES` | Mutate Portfolio Bidding Strategies | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_SHARED_CRITERIA` | Mutate Shared Criteria | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_SHARED_SETS` | Mutate Shared Sets | write | `operations`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_MUTATE_USER_LISTS` | Mutate User Lists | write | `operations`; *customer_id, partial_failure, validate_only* |  |  |
| `GOOGLESUPER_NEARBY_SEARCH` | Nearby search | read | `latitude`, `longitude`, `radius`; *excludedTypes, fieldMask, includedTypes, maxResultCount* |  | `cloud-platform` |
| `GOOGLESUPER_PATCH_CALENDAR` | Patch Calendar | write | `calendar_id`; *description, location, summary, timezone* |  | `calendar`, `calendar.app.created`, `calendar.calendars` |
| `GOOGLESUPER_PATCH_EVENT` | Patch Event | write | `calendar_id`, `event_id`; *anyone_can_add_self, attachments, attendees, attendees_omitted, color_id, conference_data, conference_data_version, description and 23 more* |  | `calendar`, `calendar.app.created`, `calendar.events`, `calendar.events.owned` |
| `GOOGLESUPER_PATCH_LABEL` | Patch Label | write | `id`, `userId`; *color, labelListVisibility, messageListVisibility, name* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.labels` |
| `GOOGLESUPER_PATCH_PERMISSION` | Patch Permission | write | `file_id`, `permission_id`; *additional_roles, expiration_date, remove_expiration, role, supports_all_drives, transfer_ownership, use_domain_admin_access, with_link* |  | `docs`, `drive`, `drive.file` |
| `GOOGLESUPER_PATCH_PROPERTY` | Patch Property (v2 API) | write | `fileId`, `propertyKey`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, quotaUser and 5 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.metadata` |
| `GOOGLESUPER_PATCH_SEND_AS` | Patch send-as alias | write | `send_as_email`; *display_name, is_default, reply_to_address, signature, smtp_msa, treat_as_alias, user_id* |  | `gmail.settings.basic`, `gmail.settings.sharing` |
| `GOOGLESUPER_PATCH_TASK` | Patch Task | write | `task_id`, `tasklist_id`; *completed, deleted, due, etag, hidden, id, notes, status and 1 more* |  | `tasks` |
| `GOOGLESUPER_PATCH_TASK_LIST` | Patch task list | write | `tasklist_id`; *updated_title* |  | `tasks` |
| `GOOGLESUPER_PLACE_PHOTO` | Get Place Photo | read | `photo_reference`; *maxheight, maxwidth* |  |  |
| `GOOGLESUPER_PRESENTATIONS_BATCH_UPDATE` | Update Presentation (Batch/Markdown) | write | `presentationId`; *markdown_text, requests, writeControl* |  | `drive`, `drive.file`, `drive.readonly`, `presentations`, `spreadsheets`, `spreadsheets.readonly` |
| `GOOGLESUPER_PRESENTATIONS_COPY_FROM_TEMPLATE` | Copy Google Slides from Template | write | `template_presentation_id`; *new_title, parent_folder_id* |  | `drive`, `drive.appdata`, `drive.file`, `drive.photos.readonly` |
| `GOOGLESUPER_PRESENTATIONS_GET` | Get Presentation | read | *fields, presentationId, presentationName* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.photos.readonly`, `drive.readonly`, `presentations`, `presentations.readonly` |
| `GOOGLESUPER_PRESENTATIONS_PAGES_GET` | Get Presentation Page | read | `pageObjectId`, `presentationId` |  | `drive`, `drive.file`, `drive.readonly`, `presentations`, `presentations.readonly` |
| `GOOGLESUPER_PROVISION_ACCOUNT_TICKET` | Provision Account Ticket | write | `account`, `redirectUri` |  | `analytics.edit` |
| `GOOGLESUPER_QUERY_AUDIENCE_EXPORT` | Query Audience Export | read | `name`; *limit, offset* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_QUERY_AUDIENCE_LIST` | Query Audience List | read | `name`; *limit, offset* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_QUERY_REPORT_TASK` | Query Report Task | read | `name`; *limit, offset* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_QUICK_ADD` | Quick Add Event | write | `text`; *calendar_id, send_updates* |  | `calendar`, `calendar.events`, `calendar.events.owned`, `calendar.app.created` |
| `GOOGLESUPER_READ_SPREADSHEET_VALUES` | Read Spreadsheet Values | read | `ranges`, `spreadsheet_id`; *date_time_render_option, major_dimension, value_render_option* |  | `drive`, `drive.readonly`, `spreadsheets`, `spreadsheets.readonly` |
| `GOOGLESUPER_REMOVE_ALL_AUDIENCE_MEMBERS` | Remove all audience members | write | `destinations`; *remove_as_of_time, validate_only* |  |  |
| `GOOGLESUPER_REMOVE_ATTENDEE` | Remove attendee from event | write | `attendee_email`, `event_id`; *calendar_id* |  | `calendar`, `calendar.app.created`, `calendar.events`, `calendar.events.owned` |
| `GOOGLESUPER_REMOVE_AUDIENCE_MEMBERS` | Remove audience members | write | `audience_members`, `destinations`; *encoding, encryption_info, validate_only* |  |  |
| `GOOGLESUPER_REMOVE_USERS` | Remove users | write | `destinations`, `encoding`, `user_data`; *encryption_info, validate_only* |  |  |
| `GOOGLESUPER_RENDER_AERIAL_VIDEO` | Render Aerial Video | write | `address` |  | `cloud-platform` |
| `GOOGLESUPER_REPLACE_ALL_TEXT` | Replace All Text in Document | write | `document_id`, `find_text`, `replace_text`; *match_case, search_by_regex, tab_ids* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_REPLACE_IMAGE` | Replace Image in Document | write | `document_id`, `replace_image` |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_REPLY_TO_THREAD` | Reply to email thread | write | `thread_id`; *attachment, bcc, cc, extra_recipients, include_signature, is_html, message_body, recipient_email and 1 more* |  | `https://mail.google.com/`, `gmail.compose`, `gmail.modify`, `gmail.send`, `gmail.readonly`, `contacts`, `contacts.other.readonly`, `contacts.readonly` |
| `GOOGLESUPER_RESUMABLE_UPLOAD` | Resumable Upload | write | `file_to_upload`; *chunkSize, file_id, folder_to_upload_to, metadata, queryParams* |  | `drive`, `drive.appdata`, `drive.file`, `drive.metadata`, `drive.scripts` |
| `GOOGLESUPER_RUN_FUNNEL_REPORT` | Run Funnel Report | read | `funnel`, `property`; *dateRanges, dimensionFilter, funnelBreakdown, funnelNextAction, funnelVisualizationType, limit, returnPropertyQuota, segments* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_RUN_PIVOT_REPORT` | Run Pivot Report | read | `metrics`, `property`; *cohortSpec, comparisons, currencyCode, dateRanges, dimensionFilter, dimensions, keepEmptyRows, metricFilter and 2 more* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_RUN_REALTIME_REPORT` | Run Realtime Report | read | `property`; *dimensionFilter, dimensions, limit, metricAggregations, metricFilter, metrics, minuteRanges, orderBys and 1 more* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_RUN_REPORT` | Run Report | read | `property`; *cohortSpec, comparisons, currencyCode, dateRanges, dimensionFilter, dimensions, keepEmptyRows, limit and 7 more* |  | `analytics.readonly`, `analytics` |
| `GOOGLESUPER_SCHEDULE_EXPERIMENT` | Schedule Experiment | write | `resource_name`; *customer_id, validate_only* |  |  |
| `GOOGLESUPER_SEARCH_DEVELOPER_METADATA` | Search Developer Metadata | read | `dataFilters`, `spreadsheetId` |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_SEARCH_DOCUMENTS` | Search Documents | read | *created_after, include_shared_drives, include_trashed, max_results, modified_after, order_by, page_token, query and 3 more* |  | `drive`, `drive.appdata`, `drive.file`, `drive.readonly`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly` |
| `GOOGLESUPER_SEARCH_GAQL` | Search GAQL | read | `query`; *customer_id, page_token, search_settings, validate_only* |  |  |
| `GOOGLESUPER_SEARCH_GOOGLE_ADS_FIELDS` | Search Google Ads Fields | read | `query`; *page_size, page_token* |  |  |
| `GOOGLESUPER_SEARCH_MEDIA_ITEMS` | Search Media Items | read | *albumId, filters, orderBy, pageSize, pageToken* |  | `photoslibrary`, `photoslibrary.readonly` |
| `GOOGLESUPER_SEARCH_PEOPLE` | Search People | read | `query`; *other_contacts, page_size, person_fields* |  | `contacts`, `contacts.readonly`, `contacts.other.readonly` |
| `GOOGLESUPER_SEARCH_SPREADSHEETS` | Search Spreadsheets | read | *created_after, include_shared_drives, include_trashed, max_results, modified_after, order_by, page_token, query and 3 more* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_SEARCH_STREAM_GAQL` | Search Stream GAQL | read | `query`; *customer_id, summary_row_setting* |  |  |
| `GOOGLESUPER_SEND_DRAFT` | Send Draft | write | `draft_id`; *user_id* |  | `https://mail.google.com/`, `gmail.compose`, `gmail.modify` |
| `GOOGLESUPER_SEND_EMAIL` | Send Email | write | *attachment, bcc, body, cc, extra_recipients, from_email, include_signature, is_html and 3 more* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.compose`, `gmail.send` |
| `GOOGLESUPER_SEND_EVENTS` | Send Events | write | `api_secret`, `client_id`, `events`, `measurement_id`; *consent, timestamp_micros, user_id, user_properties* |  |  |
| `GOOGLESUPER_SETTINGS_GET` | Get Calendar Setting | read | `setting` |  | `calendar.settings.readonly`, `calendar.readonly`, `calendar` |
| `GOOGLESUPER_SETTINGS_GET_IMAP` | Get IMAP Settings | read | *user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_SETTINGS_GET_POP` | Get POP settings | read | *user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_SETTINGS_LIST` | List Settings | read | *maxResults, pageToken, syncToken* |  | `calendar.settings.readonly`, `calendar.readonly`, `calendar` |
| `GOOGLESUPER_SETTINGS_SEND_AS_GET` | Get send-as alias | read | `send_as_email`; *user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_SETTINGS_WATCH` | Watch Settings | write | `address`, `id`, `type`; *expiration, params, token* |  | `calendar.settings.readonly`, `calendar.readonly`, `calendar` |
| `GOOGLESUPER_SET_BASIC_FILTER` | Set Basic Filter | write | `filter`, `spreadsheetId`; *includeSpreadsheetInResponse, responseIncludeGridData, responseRanges* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_SET_DATA_VALIDATION_RULE` | Set Data Validation Rule | write | `end_column_index`, `end_row_index`, `mode`, `sheet_id`, `spreadsheet_id`, `start_column_index`, `start_row_index`; *condition_values, filtered_rows_included, formula, input_message, show_custom_ui, source_range_a1, strict, validation_type and 1 more* |  | `spreadsheets`, `drive`, `drive.file` |
| `GOOGLESUPER_SOFT_DELETE_PROPERTY` | Delete Property | write | `name` |  | `analytics.edit` |
| `GOOGLESUPER_SORT_RANGE` | Sort Range | write | `range`, `sortSpecs`, `spreadsheet_id`; *includeSpreadsheetInResponse, responseIncludeGridData, responseRanges* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_SPREADSHEETS_SHEETS_COPY_TO` | Copy Sheet to Another Spreadsheet | write | `destination_spreadsheet_id`, `sheet_id`, `spreadsheet_id` |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_SPREADSHEETS_VALUES_APPEND` | Append Values to Spreadsheet | write | `range`, `spreadsheetId`, `valueInputOption`, `values`; *includeValuesInResponse, insertDataOption, majorDimension, responseDateTimeRenderOption, responseValueRenderOption* |  | `spreadsheets`, `drive`, `drive.file` |
| `GOOGLESUPER_SPREADSHEETS_VALUES_BATCH_CLEAR` | Batch Clear Spreadsheet Values | write | `ranges`, `spreadsheet_id` |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_SPREADSHEETS_VALUES_BATCH_GET_BY_DATA_FILTER` | Batch Get Spreadsheet Values by Data Filter | read | `dataFilters`, `spreadsheetId`; *dateTimeRenderOption, majorDimension, valueRenderOption* |  | `drive`, `drive.file`, `spreadsheets` |
| `GOOGLESUPER_STOP_WATCH` | Stop watch notifications | write | *user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GOOGLESUPER_STOP_WATCH_CHANNEL` | Stop Watch Channel | write | `id`, `resourceId`; *address, channelType, expiration, kind, params, payload, resourceUri, token* |  | `docs`, `drive`, `drive.appdata`, `drive.apps`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_TEXT_SEARCH` | Text Search | read | `textQuery`; *fieldMask, maxResultCount* |  | `cloud-platform` |
| `GOOGLESUPER_TILES_CREATE_SESSION` | Create Tiles Session | write | `language`, `mapType`, `region`; *highDpi, imageFormat, layerTypes, overlay, scale, styles* |  | `cloud-platform` |
| `GOOGLESUPER_TRASH_FILE` | Trash File | write | `file_id`; *fields, supportsAllDrives* |  | `drive`, `drive.appdata`, `drive.file`, `drive.metadata`, `drive.scripts` |
| `GOOGLESUPER_UNHIDE_DRIVE` | Unhide Shared Drive | write | `driveId` |  | `drive` |
| `GOOGLESUPER_UNMERGE_TABLE_CELLS` | Unmerge Table Cells | write | `documentId`, `tableRange` |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_UNTRASH_FILE` | Untrash File | write | `file_id`; *supportsAllDrives* |  | `drive`, `drive.file`, `drive.metadata`, `drive.appdata`, `drive.scripts` |
| `GOOGLESUPER_UNTRASH_MESSAGE` | Untrash Message | write | `message_id`; *user_id* |  | `https://mail.google.com/`, `gmail.modify` |
| `GOOGLESUPER_UNTRASH_THREAD` | Untrash thread | write | `thread_id`; *user_id* |  | `https://mail.google.com/`, `gmail.modify` |
| `GOOGLESUPER_UPDATE_ALBUM` | Update Album | write | `albumId`; *coverPhotoMediaItemId, title* |  | `photoslibrary.edit.appcreateddata` |
| `GOOGLESUPER_UPDATE_CHART` | Update Chart in Google Sheets | write | `chart_id`, `chart_spec`, `spreadsheet_id` |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_UPDATE_COMMENT` | Update Comment | write | `comment_id`, `file_id`; *content, fields, resolved* |  | `drive`, `drive.file` |
| `GOOGLESUPER_UPDATE_CONTACT` | Update Google Contact | write | `person`, `resource_name`, `update_person_fields`; *person_fields, sources* |  | `contacts` |
| `GOOGLESUPER_UPDATE_DATA_MANAGER_USER_LIST` | Update data manager user list | write | `account_id`, `account_type`, `user_list_id`; *account_access_status, description, display_name, ingested_user_list_info, integration_code, linked_account, login_account, membership_duration and 4 more* |  |  |
| `GOOGLESUPER_UPDATE_DATA_RETENTION_SETTINGS` | Update Data Retention Settings | write | `name`, `updateMask`; *eventDataRetention, resetUserDataOnNewActivity, userDataRetention* |  | `analytics.edit` |
| `GOOGLESUPER_UPDATE_DATA_STREAM` | Update Data Stream | write | `name`; *defaultUri, displayName* |  | `analytics.edit` |
| `GOOGLESUPER_UPDATE_DIMENSION_PROPERTIES` | Update Dimension Properties (Hide/Unhide & Resize) | write | `dimension`, `end_index`, `spreadsheet_id`, `start_index`; *hidden_by_user, include_spreadsheet_in_response, pixel_size, response_include_grid_data, response_ranges, sheet_id, sheet_name* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_UPDATE_DOCUMENT_MARKDOWN` | Update Document Markdown | write | `id`, `markdown`; *tab_id* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_UPDATE_DOCUMENT_SECTION_MARKDOWN` | Update Document Section Markdown | write | `document_id`, `markdown_text`; *end_index, start_index, tab_id* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_UPDATE_DOCUMENT_STYLE` | Update Document Style | write | `document_id`, `document_style`; *fields, tab_id* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_UPDATE_DRAFT` | Update draft | write | `draft_id`; *attachment, bcc, body, cc, extra_recipients, include_signature, is_html, recipient_email and 3 more* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.compose` |
| `GOOGLESUPER_UPDATE_DRIVE` | Update Shared Drive | write | `driveId`; *backgroundImageFile, colorRgb, hidden, name, restrictions, themeId, useDomainAdminAccess* |  | `drive` |
| `GOOGLESUPER_UPDATE_EVENT` | Update Google event | write | `event_id`, `start_datetime`; *attendees, birthdayProperties, calendar_id, create_meeting_room, description, end_datetime, eventType, event_duration_hour and 15 more* |  | `calendar`, `calendar.events`, `calendar.events.owned`, `calendar.app.created` |
| `GOOGLESUPER_UPDATE_EXISTING_DOCUMENT` | Update existing document | write | `document_id`, `edit_docs` |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_UPDATE_EXPERIMENT` | Update Experiment | write | `experiment`; *customer_id, update_mask, validate_only* |  |  |
| `GOOGLESUPER_UPDATE_EXPERIMENT_ARMS` | Update Experiment Arms | write | `updates`; *customer_id, partial_failure, response_content_type, validate_only* |  |  |
| `GOOGLESUPER_UPDATE_FILE_METADATA_PATCH` | Update File Metadata (PATCH v2) | write | `fileId`; *addParents, description, includeLabels, includePermissionsForView, indexableText, labels, mimeType, modifiedDate and 14 more* |  | `docs`, `drive`, `drive.appdata`, `drive.apps.readonly`, `drive.file`, `drive.metadata`, `drive.scripts` |
| `GOOGLESUPER_UPDATE_FILE_PROPERTY` | Update Property (v2 API) | write | `fileId`, `propertyKey`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, property_value and 6 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.metadata` |
| `GOOGLESUPER_UPDATE_FILE_PUT` | Update File (Metadata) | write | `fileId`; *add_parents, description, keep_revision_forever, mime_type, name, ocr_language, remove_parents, starred and 4 more* |  | `drive`, `drive.appdata`, `drive.file`, `drive.metadata`, `drive.scripts` |
| `GOOGLESUPER_UPDATE_FILE_REVISION_METADATA` | Update File Revision Metadata | write | `file_id`, `revision_id`; *keepForever, publishAuto, published, publishedOutsideDomain* |  | `drive`, `drive.appdata`, `drive.file` |
| `GOOGLESUPER_UPDATE_IMAP_SETTINGS` | Update IMAP settings | write | *autoExpunge, enabled, expungeBehavior, maxFolderSize, user_id* |  | `gmail.settings.basic` |
| `GOOGLESUPER_UPDATE_LABEL` | Update Label | write | `id`; *color, labelListVisibility, messageListVisibility, name, userId* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.labels` |
| `GOOGLESUPER_UPDATE_LANGUAGE_SETTINGS` | Update Language Settings | write | `display_language`; *user_id* |  | `https://mail.google.com/`, `gmail.settings.basic` |
| `GOOGLESUPER_UPDATE_MEDIA_ITEM` | Update Media Item | write | `description`, `mediaItemId` |  | `photoslibrary.readonly` |
| `GOOGLESUPER_UPDATE_MEMBER` | Update Google Meet Member | write | `name`; *email, role, update_mask* |  | `meetings.space.created` |
| `GOOGLESUPER_UPDATE_PERMISSION` | Update Permission | write | `fileId`, `permission`, `permissionId`; *enforceExpansiveAccess, removeExpiration, supportsAllDrives, transferOwnership, useDomainAdminAccess* |  | `drive`, `drive.file` |
| `GOOGLESUPER_UPDATE_POP_SETTINGS` | Update POP settings | write | *access_window, disposition, user_id* |  | `https://mail.google.com/`, `gmail.settings.basic` |
| `GOOGLESUPER_UPDATE_PROPERTY` | Update Property | write | `name`, `updateMask`; *currencyCode, displayName, industryCategory, parent, propertyType, timeZone* |  | `analytics.edit` |
| `GOOGLESUPER_UPDATE_REPLY` | Update Reply | write | `comment_id`, `content`, `file_id`, `reply_id`; *fields* |  | `drive`, `drive.file` |
| `GOOGLESUPER_UPDATE_SEND_AS` | Update send-as alias | write | `send_as_email`; *display_name, is_default, reply_to_address, signature, smtp_msa, treat_as_alias, user_id* |  | `gmail.settings.basic`, `gmail.settings.sharing` |
| `GOOGLESUPER_UPDATE_SHEET_PROPERTIES` | Update Sheet Properties | write | `spreadsheetId`, `updateSheetProperties`; *includeSpreadsheetInResponse, responseIncludeGridData, responseRanges* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_UPDATE_SPACE` | Update Google Meet Space | write | `name`; *config, updateMask* |  | `meetings.space.created`, `meetings.space.settings` |
| `GOOGLESUPER_UPDATE_SPREADSHEET_PROPERTIES` | Update Spreadsheet Properties | write | `fields`, `properties`, `spreadsheetId`; *includeSpreadsheetInResponse, responseIncludeGridData, responseRanges* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESUPER_UPDATE_SPREADSHEET_VALUES` | Update Spreadsheet Values | write | `range`, `spreadsheet_id`, `value_input_option`, `values`; *include_values_in_response, major_dimension, response_date_time_render_option, response_value_render_option* |  | `drive` |
| `GOOGLESUPER_UPDATE_TABLE_ROW_STYLE` | Update Table Row Style | write | `documentId`, `updateTableRowStyle` |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_UPDATE_TAB_PROPERTIES` | Update Tab Properties | write | `document_id`, `tab_id`, `tab_properties`; *fields* |  | `documents`, `drive`, `drive.file` |
| `GOOGLESUPER_UPDATE_TASK_FULL` | Update Task (Full Replacement) | write | `id`, `task_id`, `tasklist_id`, `title`; *completed, deleted, due, etag, hidden, links, notes, parent and 2 more* |  | `tasks` |
| `GOOGLESUPER_UPDATE_TASK_LIST` | Update Task List | write | `tasklist_id`, `title` |  | `tasks` |
| `GOOGLESUPER_UPDATE_VACATION_SETTINGS` | Update Vacation Settings | write | *enableAutoReply, endTime, responseBodyHtml, responseBodyPlainText, responseSubject, restrictToContacts, restrictToDomain, startTime and 1 more* |  | `gmail.settings.basic` |
| `GOOGLESUPER_UPDATE_VALUES_BATCH` | Batch update spreadsheet values | write | `data`, `spreadsheet_id`, `valueInputOption`; *includeValuesInResponse, responseDateTimeRenderOption, responseValueRenderOption* |  | `drive`, `drive.file`, `spreadsheets` |
| `GOOGLESUPER_UPLOAD_CALL_CONVERSIONS` | Upload Call Conversions | write | `conversions`; *customer_id, partial_failure, validate_only* |  |  |
| `GOOGLESUPER_UPLOAD_CLICK_CONVERSIONS` | Upload Click Conversions | write | `conversions`; *customer_id, job_id, partial_failure, validate_only* |  |  |
| `GOOGLESUPER_UPLOAD_CONVERSION_ADJUSTMENTS` | Upload Conversion Adjustments | write | `conversion_adjustments`; *customer_id, job_id, partial_failure, validate_only* |  |  |
| `GOOGLESUPER_UPLOAD_DRIVE_FILE` | Upload Drive File | write | `file`; *description, fields, mime_type, name, parent_folder_id, supports_all_drives* |  | `drive`, `drive.file` |
| `GOOGLESUPER_UPLOAD_FILE` | Upload File | write | `file_to_upload`; *folder_to_upload_to* |  | `drive`, `drive.appdata`, `drive.file` |
| `GOOGLESUPER_UPLOAD_FROM_URL` | Upload File from URL to Drive | write | `name`, `source_url`; *mime_type, parent_folder_id, source_headers, supports_all_drives, verify_ssl* |  | `drive`, `drive.appdata`, `drive.file` |
| `GOOGLESUPER_UPLOAD_MEDIA` | Upload Media | write | *description, file_name, file_to_upload, url* |  | `photoslibrary.appendonly`, `photoslibrary` |
| `GOOGLESUPER_UPLOAD_UPDATE_FILE` | Upload/Update File Content | write | `fileId`, `file_to_upload`; *addParents, keepRevisionForever, ocrLanguage, removeParents, supportsAllDrives, uploadType, useContentAsIndexableText* |  | `drive`, `drive.appdata`, `drive.file`, `drive.metadata`, `drive.scripts` |
| `GOOGLESUPER_UPLOAD_USER_DATA` | Upload User Data | write | `operations`; *customer_id, customer_match_user_list_metadata* |  |  |
| `GOOGLESUPER_UPSERT_ROWS` | Upsert Rows (Smart Update/Insert) | write | `rows`, `sheetName`, `spreadsheetId`; *headers, keyColumn, key_column_index, normalization_message, strictMode, tableStart* |  | `spreadsheets`, `drive`, `drive.file` |
| `GOOGLESUPER_VALIDATE_EVENTS` | Validate Events | read | `api_secret`, `client_id`, `events`, `measurement_id`; *consent, timestamp_micros, user_id, user_properties* |  |  |
| `GOOGLESUPER_VALUES_GET` | Get spreadsheet values | read | `range`, `spreadsheet_id`; *date_time_render_option, end_row, major_dimension, start_row, value_render_option* |  | `drive`, `drive.readonly`, `drive.file`, `spreadsheets`, `spreadsheets.readonly` |
| `GOOGLESUPER_VALUES_UPDATE` | Update spreadsheet values | write | `range`, `spreadsheet_id`, `value_input_option`, `values`; *auto_expand_sheet, include_values_in_response, major_dimension, response_datetime_render_option, response_value_render_option* |  | `drive`, `drive.file`, `spreadsheets` |
| `GOOGLESUPER_WATCH_CHANGES` | Watch Drive Changes | write | `address`, `id`, `type`; *drive_id, expiration, include_corpus_removals, include_items_from_all_drives, include_labels, include_permissions_for_view, include_removed, page_size and 6 more* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_WATCH_FILE` | Watch File for Changes | write | `address`, `fileId`, `id`, `type`; *acknowledgeAbuse, expiration, includeLabels, includePermissionsForView, params, payload, supportsAllDrives, supportsTeamDrives and 1 more* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESUPER_WHO_AM_I` | Who Am I | read |  |  |  |
| `GOOGLESUPER_WRITE_SPREADSHEET_VALUES` | Write Spreadsheet Values | write | `range`, `spreadsheet_id`, `value_input_option`, `values`; *include_values_in_response, major_dimension, response_date_time_render_option, response_value_render_option* |  | `drive`, `spreadsheets` |
