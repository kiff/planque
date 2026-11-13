# Google Drive: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

92 tools: 41 read, 51 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `GOOGLEDRIVE_ADD_PARENT` | Insert File Parent (v2) | write | `fileId`, `id`; *enforceSingleParent, supportsAllDrives, supportsTeamDrives* |  | `docs`, `drive`, `drive.appdata`, `drive.file` |
| `GOOGLEDRIVE_ADD_PROPERTY` | Insert Property (v2 API) | write | `fileId`, `property_key`, `property_value`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, quotaUser and 4 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.metadata` |
| `GOOGLEDRIVE_BATCH_GET_SPREADSHEET_VALUES` | Batch Get Spreadsheet Values | read | `ranges`, `spreadsheet_id`; *date_time_render_option, major_dimension, value_render_option* |  | `drive`, `drive.readonly` |
| `GOOGLEDRIVE_BATCH_UPDATE_SPREADSHEET_VALUES` | Batch Update Spreadsheet Values | write | `data`, `spreadsheet_id`, `value_input_option`; *include_values_in_response, response_date_time_render_option, response_value_render_option* |  | `drive` |
| `GOOGLEDRIVE_CLEAR_SPREADSHEET_VALUES` | Clear Spreadsheet Values | write | `range`, `spreadsheet_id` |  | `drive`, `spreadsheets` |
| `GOOGLEDRIVE_COPY_FILE_ADVANCED` | Copy file with advanced options | write | `fileId`; *access_token, alt, appProperties, callback, copyRequiresWriterPermission, createdTime, description, enforceSingleParent and 24 more* |  | `drive`, `drive.appdata`, `drive.file`, `drive.photos.readonly` |
| `GOOGLEDRIVE_CREATE_COMMENT` | Create Comment | write | `content`, `file_id`; *anchor, quoted_file_content_mime_type, quoted_file_content_value* |  | `drive`, `drive.file` |
| `GOOGLEDRIVE_CREATE_DRIVE` | Create Shared Drive | write | `name`; *backgroundImageFile, colorRgb, hidden, requestId, themeId* |  | `drive` |
| `GOOGLEDRIVE_CREATE_FILE` | Create File or Folder | write | *description, fields, file_to_upload, mimeType, name, parents, starred* |  | `drive`, `drive.file`, `drive.appdata` |
| `GOOGLEDRIVE_CREATE_FILE_FROM_TEXT` | Create a File from Text | write | `file_name`, `text_content`; *mime_type, parent_id* |  | `drive`, `drive.appdata`, `drive.file` |
| `GOOGLEDRIVE_CREATE_FOLDER` | Create a folder | write | `name`; *parent_id* |  | `drive`, `drive.file` |
| `GOOGLEDRIVE_CREATE_PERMISSION` | Create Permission | write | `file_id`, `role`, `type`; *allow_file_discovery, domain, email_address, email_message, expiration_time, move_to_new_owners_root, send_notification_email, supports_all_drives and 2 more* |  | `drive`, `drive.file` |
| `GOOGLEDRIVE_CREATE_PERMISSIONS_BATCH` | Create Permissions Batch | write | `permissions` |  | `drive`, `drive.file` |
| `GOOGLEDRIVE_CREATE_REPLY` | Create Reply | write | `comment_id`, `content`, `file_id`; *action, fields* |  | `drive`, `drive.file` |
| `GOOGLEDRIVE_CREATE_SHORTCUT_TO_FILE` | Create Shortcut to File/Folder | write | `name`, `target_id`; *ignoreDefaultVisibility, includeLabels, includePermissionsForView, keepRevisionForever, supportsAllDrives* |  | `drive`, `drive.file` |
| `GOOGLEDRIVE_DELETE_CHILD` | Delete Child (v2) | write | `childId`, `folderId`; *access_token, alt, callback, enforceSingleParent, fields, key, oauth_token, prettyPrint and 4 more* |  | `docs`, `drive`, `drive.file` |
| `GOOGLEDRIVE_DELETE_COMMENT` | Delete Comment | write | `comment_id`, `file_id` |  | `drive`, `drive.file` |
| `GOOGLEDRIVE_DELETE_DRIVE` | Delete Shared Drive | write | `driveId`; *allowItemDeletion, useDomainAdminAccess* |  | `drive` |
| `GOOGLEDRIVE_DELETE_PARENT` | Delete Parent (v2) | write | `fileId`, `parentId`; *access_token, alt, callback, enforceSingleParent, fields, key, oauth_token, prettyPrint and 4 more* |  | `docs`, `drive`, `drive.file` |
| `GOOGLEDRIVE_DELETE_PERMISSION` | Delete Permission | write | `file_id`, `permission_id`; *supportsAllDrives, useDomainAdminAccess* |  | `drive`, `drive.file` |
| `GOOGLEDRIVE_DELETE_PROPERTY` | Delete Property (v2 API) | write | `fileId`, `propertyKey`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, quotaUser and 4 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.metadata` |
| `GOOGLEDRIVE_DELETE_REPLY` | Delete Reply | write | `comment_id`, `file_id`, `reply_id` |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_DELETE_REVISION` | Delete Revision | write | `file_id`, `revision_id` |  | `drive`, `drive.appdata`, `drive.file` |
| `GOOGLEDRIVE_DOWNLOAD_FILE` | Download a file from Google Drive | read | `fileId`; *mime_type* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly`, `drive.appdata`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly` |
| `GOOGLEDRIVE_DOWNLOAD_FILE_OPERATION` | Download file via operation | read | `file_id`; *mime_type, revision_id* |  | `drive`, `drive.file`, `drive.readonly` |
| `GOOGLEDRIVE_DOWNLOAD_FILE_REVISION` | Download File Revision | read | `file_id`, `revision_id`; *acknowledge_abuse, file_name* |  | `drive`, `drive.readonly` |
| `GOOGLEDRIVE_EDIT_FILE` | Edit File | write | `content`, `file_id`; *mime_type* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_EMPTY_TRASH` | Empty Trash | write | *driveId, enforceSingleParent* |  | `drive` |
| `GOOGLEDRIVE_EXPORT_GOOGLE_WORKSPACE_FILE` | Export Google Workspace file | read | `fileId`, `mimeType` |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_FIND_FILE` | Find file | read | *bare_text_query_transformed, corpora, driveId, editors_field_removed, email_query_transformed, emailaddress_field_removed, fields, folder_id and 14 more* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_FIND_FOLDER` | Find folder | read | *full_text_contains, full_text_not_contains, modified_after, name_contains, name_exact, name_not_contains, page_size, page_token and 2 more* |  | `drive`, `drive.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.file`, `drive.appdata`, `drive.meet.readonly`, `drive.photos.readonly` |
| `GOOGLEDRIVE_GENERATE_IDS` | Generate File IDs | read | *count, space, type* |  | `drive`, `drive.appdata`, `drive.file` |
| `GOOGLEDRIVE_GET_ABOUT` | Get about | read | *fields* |  | `drive`, `drive.appdata`, `drive.file`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_GET_APP` | Get App | read | `appId` |  | `docs`, `drive`, `drive.appdata`, `drive.apps.readonly`, `drive.file`, `drive.metadata`, `drive.metadata.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_GET_CHANGES_START_PAGE_TOKEN` | Get Changes Start Page Token | read | *driveId, supportsAllDrives, supportsTeamDrives, teamDriveId* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_GET_CHILD` | Get Child Reference (v2) | read | `childId`, `folderId`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, quotaUser and 3 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_GET_COMMENT` | Get Comment | read | `commentId`, `fileId`; *includeDeleted* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_GET_DOCUMENT` | Get Document | read | `document_id`; *include_tabs_content, suggestions_view_mode* |  | `drive`, `drive.readonly` |
| `GOOGLEDRIVE_GET_DRIVE` | Get Shared Drive | read | `drive_id`; *use_domain_admin_access* |  | `drive`, `drive.readonly` |
| `GOOGLEDRIVE_GET_FILE_METADATA` | Get File Metadata | read | `fileId`; *fields, includeLabels, includePermissionsForView, supportsAllDrives* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_GET_FILE_PROPERTY` | Get Property (v2) | read | `fileId`, `propertyKey`; *visibility* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_GET_FORM` | Get Google Form | read | `form_id` |  | `drive`, `drive.readonly` |
| `GOOGLEDRIVE_GET_PARENT` | Get Parent Reference (v2) | read | `fileId`, `parentId`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, quotaUser and 3 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_GET_PERMISSION` | Get Permission | read | `file_id`, `permission_id`; *fields, supports_all_drives, use_domain_admin_access* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_GET_PERMISSION_ID_FOR_EMAIL` | Get Permission ID for Email | read | `email`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, quotaUser and 3 more* |  | `docs`, `drive`, `drive.appdata`, `drive.apps.readonly`, `drive.file`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_GET_PRESENTATION` | Get Google Slides Presentation | read | `presentation_id`; *fields* |  | `drive`, `drive.readonly` |
| `GOOGLEDRIVE_GET_PRESENTATION_PAGE_THUMBNAIL` | Get Presentation Page Thumbnail | read | `page_object_id`, `presentation_id`; *mime_type, thumbnail_size* |  | `drive.readonly`, `drive` |
| `GOOGLEDRIVE_GET_REPLY` | Get Reply | read | `commentId`, `fileId`, `replyId`; *includeDeleted* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_GET_REVISION` | Get Revision | read | `file_id`, `revision_id`; *acknowledge_abuse* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_GET_SPREADSHEET` | Get Spreadsheet | read | `spreadsheet_id`; *exclude_tables_in_banded_ranges, fields, include_grid_data, ranges* |  | `drive`, `drive.readonly` |
| `GOOGLEDRIVE_GET_SPREADSHEET_VALUES` | Get Spreadsheet Values | read | `range`, `spreadsheet_id`; *date_time_render_option, major_dimension, value_render_option* |  | `drive`, `drive.readonly` |
| `GOOGLEDRIVE_GOOGLE_DRIVE_DELETE_FOLDER_OR_FILE_ACTION` | Delete folder or file | write | `fileId`; *supportsAllDrives* |  | `drive`, `drive.appdata`, `drive.file` |
| `GOOGLEDRIVE_HIDE_DRIVE` | Hide Shared Drive | write | `drive_id` |  | `drive` |
| `GOOGLEDRIVE_INSERT_CHILD` | Insert Child Into Folder (v2) | write | `folderId`, `id`; *access_token, alt, callback, enforceSingleParent, fields, key, oauth_token, prettyPrint and 6 more* |  | `drive`, `drive.appdata`, `drive.file`, `docs` |
| `GOOGLEDRIVE_LIST_ACCESS_PROPOSALS` | List Access Proposals | read | `fileId`; *pageSize, pageToken* |  | `docs`, `drive`, `drive.file`, `drive.metadata`, `drive.metadata.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_LIST_APPROVALS` | List Approvals | read | `fileId`; *pageSize, pageToken* |  | `drive`, `drive.appdata`, `drive.file`, `drive.metadata`, `drive.metadata.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_LIST_CHANGES` | List Changes | read | *driveId, fields, includeCorpusRemovals, includeItemsFromAllDrives, includeLabels, includePermissionsForView, includeRemoved, pageSize and 4 more* |  | `drive.metadata.readonly`, `drive`, `drive.readonly`, `drive.metadata`, `drive.file`, `drive.appdata`, `drive.meet.readonly`, `drive.photos.readonly` |
| `GOOGLEDRIVE_LIST_CHILDREN_V2` | List Folder Children (v2) | read | `folderId`; *access_token, alt, callback, fields, key, maxResults, oauth_token, orderBy and 7 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_LIST_COMMENTS` | List Comments | read | `fileId`; *fields, includeDeleted, pageSize, pageToken, startModifiedTime* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_LIST_FILE_LABELS` | List File Labels | read | `file_id`; *max_results, page_token* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_LIST_FILE_PROPERTIES` | List Properties (v2 API) | read | `fileId`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, quotaUser and 3 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_LIST_FORM_RESPONSES` | List Google Form Responses | read | `form_id`; *filter, page_size, page_token* |  | `drive` |
| `GOOGLEDRIVE_LIST_LABEL_DEFINITIONS` | List Label Definitions | read | *customer, language_code, minimum_role, page_size, page_token, published_only* |  | `drive.labels.readonly`, `drive.labels`, `drive.admin.labels.readonly`, `drive.admin.labels` |
| `GOOGLEDRIVE_LIST_PERMISSIONS` | List Permissions | read | `fileId`; *fields, includePermissionsForView, pageSize, pageToken, supportsAllDrives, useDomainAdminAccess* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_LIST_REPLIES` | List Replies to Comment | read | `comment_id`, `file_id`; *fields, include_deleted, page_size, page_token* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_LIST_REVISIONS` | List File Revisions | read | `fileId`; *pageSize, pageToken, supportsAllDrives* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_LIST_SHARED_DRIVES` | List Shared Drives | read | *fields, pageSize, pageToken, q, useDomainAdminAccess* |  | `drive`, `drive.readonly` |
| `GOOGLEDRIVE_MODIFY_FILE_LABELS` | Modify File Labels | write | `file_id`, `label_modifications`; *kind* |  | `drive`, `drive.file`, `drive.metadata` |
| `GOOGLEDRIVE_MOVE_FILE` | Move File | write | `file_id`; *add_parents, include_labels, include_permissions_for_view, keep_revision_forever, ocr_language, remove_parents, supports_all_drives, use_content_as_indexable_text* |  | `drive`, `drive.file`, `drive.appdata`, `drive.metadata`, `drive.scripts` |
| `GOOGLEDRIVE_PATCH_PERMISSION` | Patch Permission | write | `file_id`, `permission_id`; *additional_roles, expiration_date, remove_expiration, role, supports_all_drives, transfer_ownership, use_domain_admin_access, with_link* |  | `docs`, `drive`, `drive.file` |
| `GOOGLEDRIVE_PATCH_PROPERTY` | Patch Property (v2 API) | write | `fileId`, `propertyKey`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, quotaUser and 5 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.metadata` |
| `GOOGLEDRIVE_READ_SPREADSHEET_VALUES` | Read Spreadsheet Values | read | `ranges`, `spreadsheet_id`; *date_time_render_option, major_dimension, value_render_option* |  | `drive`, `drive.readonly`, `spreadsheets`, `spreadsheets.readonly` |
| `GOOGLEDRIVE_RESUMABLE_UPLOAD` | Resumable Upload | write | `file_to_upload`; *chunkSize, file_id, folder_to_upload_to, metadata, queryParams* |  | `drive`, `drive.appdata`, `drive.file`, `drive.metadata`, `drive.scripts` |
| `GOOGLEDRIVE_STOP_WATCH_CHANNEL` | Stop Watch Channel | write | `id`, `resourceId`; *address, channelType, expiration, kind, params, payload, resourceUri, token* |  | `docs`, `drive`, `drive.appdata`, `drive.apps`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_TRASH_FILE` | Trash File | write | `file_id`; *fields, supportsAllDrives* |  | `drive`, `drive.appdata`, `drive.file`, `drive.metadata`, `drive.scripts` |
| `GOOGLEDRIVE_UNHIDE_DRIVE` | Unhide Shared Drive | write | `driveId` |  | `drive` |
| `GOOGLEDRIVE_UNTRASH_FILE` | Untrash File | write | `file_id`; *supportsAllDrives* |  | `drive`, `drive.file`, `drive.metadata`, `drive.appdata`, `drive.scripts` |
| `GOOGLEDRIVE_UPDATE_COMMENT` | Update Comment | write | `comment_id`, `file_id`; *content, fields, resolved* |  | `drive`, `drive.file` |
| `GOOGLEDRIVE_UPDATE_DRIVE` | Update Shared Drive | write | `driveId`; *backgroundImageFile, colorRgb, hidden, name, restrictions, themeId, useDomainAdminAccess* |  | `drive` |
| `GOOGLEDRIVE_UPDATE_FILE_METADATA_PATCH` | Update File Metadata (PATCH v2) | write | `fileId`; *addParents, description, includeLabels, includePermissionsForView, indexableText, labels, mimeType, modifiedDate and 14 more* |  | `docs`, `drive`, `drive.appdata`, `drive.apps.readonly`, `drive.file`, `drive.metadata`, `drive.scripts` |
| `GOOGLEDRIVE_UPDATE_FILE_PROPERTY` | Update Property (v2 API) | write | `fileId`, `propertyKey`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, property_value and 6 more* |  | `docs`, `drive`, `drive.appdata`, `drive.file`, `drive.metadata` |
| `GOOGLEDRIVE_UPDATE_FILE_PUT` | Update File (Metadata) | write | `fileId`; *add_parents, description, keep_revision_forever, mime_type, name, ocr_language, remove_parents, starred and 4 more* |  | `drive`, `drive.appdata`, `drive.file`, `drive.metadata`, `drive.scripts` |
| `GOOGLEDRIVE_UPDATE_FILE_REVISION_METADATA` | Update File Revision Metadata | write | `file_id`, `revision_id`; *keepForever, publishAuto, published, publishedOutsideDomain* |  | `drive`, `drive.appdata`, `drive.file` |
| `GOOGLEDRIVE_UPDATE_PERMISSION` | Update Permission | write | `fileId`, `permission`, `permissionId`; *enforceExpansiveAccess, removeExpiration, supportsAllDrives, transferOwnership, useDomainAdminAccess* |  | `drive`, `drive.file` |
| `GOOGLEDRIVE_UPDATE_REPLY` | Update Reply | write | `comment_id`, `content`, `file_id`, `reply_id`; *fields* |  | `drive`, `drive.file` |
| `GOOGLEDRIVE_UPDATE_SPREADSHEET_VALUES` | Update Spreadsheet Values | write | `range`, `spreadsheet_id`, `value_input_option`, `values`; *include_values_in_response, major_dimension, response_date_time_render_option, response_value_render_option* |  | `drive` |
| `GOOGLEDRIVE_UPLOAD_FILE` | Upload File | write | `file_to_upload`; *folder_to_upload_to* |  | `drive`, `drive.appdata`, `drive.file` |
| `GOOGLEDRIVE_UPLOAD_FROM_URL` | Upload File from URL to Drive | write | `name`, `source_url`; *mime_type, parent_folder_id, source_headers, supports_all_drives, verify_ssl* |  | `drive`, `drive.appdata`, `drive.file` |
| `GOOGLEDRIVE_UPLOAD_UPDATE_FILE` | Upload/Update File Content | write | `fileId`, `file_to_upload`; *addParents, keepRevisionForever, ocrLanguage, removeParents, supportsAllDrives, uploadType, useContentAsIndexableText* |  | `drive`, `drive.appdata`, `drive.file`, `drive.metadata`, `drive.scripts` |
| `GOOGLEDRIVE_WATCH_CHANGES` | Watch Drive Changes | write | `address`, `id`, `type`; *drive_id, expiration, include_corpus_removals, include_items_from_all_drives, include_labels, include_permissions_for_view, include_removed, page_size and 6 more* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_WATCH_FILE` | Watch File for Changes | write | `address`, `fileId`, `id`, `type`; *acknowledgeAbuse, expiration, includeLabels, includePermissionsForView, params, payload, supportsAllDrives, supportsTeamDrives and 1 more* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDRIVE_WRITE_SPREADSHEET_VALUES` | Write Spreadsheet Values | write | `range`, `spreadsheet_id`, `value_input_option`, `values`; *include_values_in_response, major_dimension, response_date_time_render_option, response_value_render_option* |  | `drive`, `spreadsheets` |
