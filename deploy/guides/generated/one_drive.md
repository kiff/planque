# OneDrive: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

101 tools: 58 read, 43 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `ONE_DRIVE_ADD_WORKBOOK_WORKSHEET` | Add Workbook Worksheet | write | `drive_id`, `item_id`; *name, session_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_CALCULATE_WORKBOOK` | Calculate Workbook | write | `calculation_type`, `drive_id`, `item_id`; *session_id* |  | `Files.ReadWrite` |
| `ONE_DRIVE_CHECKIN_ITEM` | Check In Drive Item | write | `drive_id`, `drive_item_id`; *check_in_as, comment* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_CHECKOUT_ITEM` | Checkout Drive Item | write | `drive_id`, `drive_item_id` |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_CLEAR_WORKBOOK_RANGE` | Clear Workbook Range | write | `address`, `drive_id`, `item_id`, `worksheet_id_or_name`; *apply_to, session_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_CLOSE_WORKBOOK_SESSION` | Close Workbook Session | write | `item_id`, `session_id`; *drive_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_COPY_ITEM` | Copy Item | write | `item_id`; *children_only, conflict_behavior, drive_id, group_id, include_all_version_history, name, parent_reference, site_id and 1 more* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_CREATE_ITEM_PERMISSION` | Create Drive Item Permission | write | `drive_id`, `drive_item_id`, `granted_to_v2`, `roles` |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Files.SelectedOperations.Selected`, `Sites.ReadWrite.All`, `Sites.FullControl.All`, `Sites.Manage.All`, `Sites.Selected`, `Lists.SelectedOperations.Selected`, `ListItems.SelectedOperations.Selected` |
| `ONE_DRIVE_CREATE_LINK` | Create Sharing Link | write | `item_id`, `type`; *drive_id, expiration_date_time, group_id, password, recipients, retain_inherited_permissions, scope, site_id and 1 more* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_CREATE_SHAREPOINT_LIST_ITEM` | Create SharePoint List Item | write | `fields`, `list_id`, `site_id` |  | `Sites.ReadWrite.All` |
| `ONE_DRIVE_CREATE_UPLOAD_SESSION` | Create Upload Session | write | `file_name`, `parent_item_id`; *conflict_behavior, defer_commit, description, drive_id, file_size, if_match_etag, if_none_match_etag* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_CREATE_UPLOAD_SESSION_BY_PATH` | Create Upload Session By Path | write | `destination_path`; *conflict_behavior, defer_commit, description, drive_id, drive_item_source, file_size, file_system_info, if_match_etag and 2 more* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_CREATE_WORKBOOK_SESSION` | Create Workbook Session | write | `item_id`; *drive_id, persist_changes* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_DELETE_ITEM` | Delete Item | write | `item_id`; *drive_id, if_match, user_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_DELETE_ITEM_PERMANENTLY` | Permanently Delete Drive Item | write | `drive_id`, `item_id` |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_DELETE_ITEM_PERMISSION` | Delete Drive Item Permission | write | `item_id`, `perm_id`; *drive_id, group_id, site_id, user_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_DELETE_SHAREPOINT_ITEM_PERMISSION` | Delete SharePoint Item Permission | write | `list_id`, `permission_id`, `site_id`; *item_id* |  | `Sites.ReadWrite.All`, `Sites.Manage.All`, `Sites.FullControl.All` |
| `ONE_DRIVE_DELETE_SHARE_PERMISSION` | Delete Shares Permission | write | `shared_drive_item_id` |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_DELETE_WORKBOOK_RANGE` | Delete Workbook Range | write | `drive_id`, `item_id`, `range`, `shift`, `worksheet` |  | `Files.ReadWrite` |
| `ONE_DRIVE_DELETE_WORKBOOK_WORKSHEET` | Delete Workbook Worksheet | write | `drive_id`, `item_id`, `worksheet_id_or_name`; *session_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_DISCARD_CHECKOUT` | Discard Checkout | write | `drive_id`, `drive_item_id` |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_DOWNLOAD_FILE` | Download a file | read | `file_name`, `item_id`; *drive_id, format, if_none_match, user_id* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_DOWNLOAD_FILE_BY_PATH` | Download file by path | read | `file_name`, `item_path`; *user_id* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_DOWNLOAD_ITEM_AS_FORMAT` | Download item as format | read | `file_name`, `format`; *drive_id, item_id, path_and_filename, user_id* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_DOWNLOAD_ITEM_THUMBNAIL` | Download Item Thumbnail | read | `item_id`; *drive_id, file_name, size, thumbnail_id* |  | `Files.Read`, `Files.Read.All`, `Sites.Read.All` |
| `ONE_DRIVE_DOWNLOAD_ITEM_VERSION` | Download Drive Item Version Content | read | `file_name`, `item_id`, `version_id`; *drive_id* |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_DOWNLOAD_THUMBNAIL_BY_PATH` | Download Thumbnail by Path | read | `item_path`; *file_name, size, thumbnail_id* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_FOLLOW_ITEM` | Follow Drive Item | write | `drive_id`, `drive_item_id` |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_FORMAT_WORKBOOK_RANGE` | Format Workbook Range | write | `address`, `item_id`, `worksheet_id_or_name`; *base_format, borders, drive_id, fill, font, session_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_DOCUMENT_COMMENTS` | Get document comments | read | `item_id`; *drive_id, user_id* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_DRIVE` | Get Drive | read | `drive_id`; *expand_fields, select_fields* |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All`, `User.Read` |
| `ONE_DRIVE_GET_DRIVE_ITEM_BY_SHARING_URL` | Get DriveItem by Sharing URL | read | *expand_children, prefer_redeem, select_fields, share_id_or_encoded_url, sharing_url* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_DRIVE_ITEM_LIST_ITEM` | Get Drive Item List Item | read | `drive_id`, `item_id`; *expand_fields, select_fields* |  | `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_FOLLOWED_ITEM` | Get Drives Following | read | `drive_id`, `drive_item_id` |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_GROUP_DRIVE` | Get Group Drive | read | `group_id`; *select_fields* |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All`, `User.Read` |
| `ONE_DRIVE_GET_ITEM` | Get Item Metadata | read | `item_id`; *drive_id, expand_relations, select_fields* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_ITEM_BY_PATH` | Get Item by Path | read | `item_path`; *drive_id, expand_relations, select_fields* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_ITEM_PERMISSIONS` | Get Item Permissions | read | *drive_id, group_id, if_none_match, item_id, item_path, page_token, select, site_id and 1 more* |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_ITEM_THUMBNAILS` | Get Item Thumbnails | read | `item_id`; *drive_id, group_id, original_orientation, select, site_id, user_id* |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_ITEM_VERSIONS` | Get Item Versions | read | `item_id`; *drive_id, group_id, site_id, user_id* |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_ROOT` | Get Drive Root Folder | read | *select_fields, user_id* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_SHARE` | Get Shared Item by ShareId | read | `share_id_or_encoded_sharing_url`; *expand_children, prefer_redeem* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_SHAREPOINT_LIST_ITEMS` | Get SharePoint List Items | read | `list_id`, `site_id`; *count, expand, filter, orderby, page_token, select, skip, top* |  | `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_SITE` | Get Site Details | read | `site_id`; *expand, select* |  | `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_SITE_PAGE_CONTENT` | Get SharePoint Site Page Content | read | `page_id`, `site_id`; *expand, select* |  | `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_SPECIAL_FOLDER` | Get Drive Special Folder | read | `special_folder_name`; *expand_relations, select_fields, user_id* |  | `Files.ReadWrite.AppFolder`, `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_USER` | Get User Profile | read | `user_id`; *select* |  | `User.Read`, `User.ReadWrite`, `User.ReadBasic.All`, `User.Read.All`, `User.ReadWrite.All`, `Directory.Read.All`, `Directory.ReadWrite.All` |
| `ONE_DRIVE_GET_WORKBOOK_RANGE_FILL` | Get Workbook Range Fill | read | `address`, `item_id`, `worksheet_id_or_name`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `ONE_DRIVE_GET_WORKBOOK_USED_RANGE` | Get Workbook Used Range | read | `drive_id`, `item_id`, `worksheet_id_or_name`; *session_id, values_only* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GET_WORKBOOK_WORKSHEET` | Get Workbook Worksheet | read | `drive_id`, `item_id`, `worksheet_id_or_name`; *select, session_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_GRANT_SHAREPOINT_ITEM_PERMISSION` | Grant SharePoint Item Permission | write | `granted_to_v2`, `list_id`, `roles`, `site_id`; *item_id* |  | `Sites.ReadWrite.All`, `Sites.Manage.All`, `Sites.FullControl.All` |
| `ONE_DRIVE_GRANT_SHARE_PERMISSION` | Grant Shares Permission | write | `encoded_sharing_url`, `recipients`, `roles` |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_INSPECT_WORKBOOK_TABLE` | Inspect Workbook Table | read | `item_id`, `table_id_or_name`; *column_page_size, drive_id, page_token, session_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_INVITE_USER_TO_ITEM` | Invite User to Drive Item | write | `item_id`, `recipients`, `roles`; *drive_id, expiration_date_time, group_id, message, password, require_sign_in, retain_inherited_permissions, send_invitation and 2 more* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_LIST_ACTIVITIES` | List Drive Activities | read | *drive_id, page_token, top* |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_LIST_ALL_DRIVE_ITEMS` | List All Drive Items | read | `drive_id`; *expand, page_token, select, top* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_LIST_BUNDLES` | List Drive Bundles | read | `drive_id`; *expand, filter, orderby, page_token, select, top* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All` |
| `ONE_DRIVE_LIST_DRIVES` | List Drives | read | *expand, group_id, orderby, page_token, select, site_id, top, user_id* |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_LIST_FOLDER_CHILDREN` | List Folder Children | read | *drive_id, expand, folder_item_id, folder_path, orderby, page_token, select, site_id and 2 more* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_LIST_FOLLOWED_SITES` | List Followed Sites | read | *page_token, select_fields* |  | `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_LIST_ITEM_ACTIVITIES` | List Drive Item Activities | read | `drive_id`, `item_id` |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_LIST_ONEDRIVE_SHARED_ITEMS` | List OneDrive Shared Items | read | *fields, from, size* |  | `Files.Read.All`, `Sites.Read.All` |
| `ONE_DRIVE_LIST_ROOT_DRIVE_CHANGES` | List Root Drive Changes | read | *drive_id, expand, select, token, top* |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_LIST_SHARED_DRIVE_ITEM_CHILDREN` | List Shared Drive Item Children | read | *expand_fields, if_none_match_etag, orderby, page_token, select_fields, share_id_or_encoded_url, sharing_url, top* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_LIST_SHAREPOINT_ITEM_PERMISSIONS` | List SharePoint Item Permissions | read | `list_id`, `site_id`; *item_id, page_token, select, top* |  | `Sites.Read.All`, `Sites.ReadWrite.All`, `Sites.Manage.All`, `Sites.FullControl.All` |
| `ONE_DRIVE_LIST_SHAREPOINT_LIST_ITEMS_DELTA` | List SharePoint List Items Delta | read | `list_id`, `site_id`; *expand, select, token, top* |  | `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_LIST_SHARE_PERMISSIONS` | List Shares Permission | read | `shared_drive_item_id` |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_LIST_SITE_COLUMNS` | List Site Columns | read | `site_id`; *count, expand, filter, orderby, select, top* |  | `Sites.Read.All`, `Sites.Manage.All`, `Sites.FullControl.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_LIST_SITE_ITEMS_DELTA` | List Site Drive Items Delta | read | `site_id`; *expand, select, token, top* |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_LIST_SITE_LISTS` | List Site Lists | read | `site_id`; *expand, filter, orderby, page_token, select, top* |  | `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_LIST_SITE_RECYCLE_BIN_ITEMS` | List Site Recycle Bin Items | read | `site_id`; *page_token, select_fields, top* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_LIST_SITE_SUBSITES` | List Site Subsites | read | `site_id`; *expand, filter, orderby, page_token, select, top* |  | `Sites.Read.All` |
| `ONE_DRIVE_LIST_SUBSCRIPTIONS` | List Subscriptions | read | *page_token* |  | `Subscription.Read.All` |
| `ONE_DRIVE_LIST_WORKBOOK_NAMES` | List Workbook Names | read | `drive_id`, `item_id`; *page_token, select_fields, session_id* |  | `Files.ReadWrite` |
| `ONE_DRIVE_LIST_WORKBOOK_WORKSHEETS` | List Workbook Worksheets | read | `drive_id`, `item_id`; *select, session_id, skip, top* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_MOVE_ITEM` | Move Item | write | `itemId`, `parentReference`; *description, driveId, groupId, name, siteId, userId* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_ONEDRIVE_CREATE_FOLDER` | Create folder | write | `name`; *conflict_behavior, description, parent_folder, user_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_ONEDRIVE_CREATE_TEXT_FILE` | Create a new text file | write | `content`, `name`; *conflict_behavior, folder, user_id* |  | `Files.ReadWrite` |
| `ONE_DRIVE_ONEDRIVE_FIND_FILE` | Find Item | read | `name`; *cursor, folder, limit, response_detail, user_id* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_ONEDRIVE_FIND_FOLDER` | Find Folder | read | *expand, folder, name, orderby, page_token, select, top, user_id* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_ONEDRIVE_LIST_ITEMS` | List OneDrive items | read | *page_token, select, top, user_id* |  | `Files.Read`, `Files.Read.All` |
| `ONE_DRIVE_ONEDRIVE_UPLOAD_FILE` | Upload file | write | *conflict_behavior, defer_commit, description, drive_id, file, file_name, file_system_info, folder and 6 more* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_PREVIEW_DRIVE_ITEM` | Preview Drive Item | read | `item_id`; *drive_id, group_id, page, share_id, site_id, user_id, zoom* |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_READ_WORKBOOK_RANGE` | Read Workbook Range | read | `address`, `drive_id`, `item_id`, `worksheet_id_or_name`; *select_fields, session_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_RESTORE_DRIVE_ITEM` | Restore Deleted Item | write | `item_id`; *drive_id, name, parent_reference_id* |  | `Files.ReadWrite.All` |
| `ONE_DRIVE_RESTORE_ITEM_VERSION` | Restore Item Version | write | `item_id`, `version_id`; *drive_id, user_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_REVOKE_SHARING_LINK_GRANTS` | Revoke Sharing Link Grants | write | `drive_id`, `grantees`, `item_id`, `permission_id` |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_SEARCH_DRIVE_ITEMS` | Search Drive Items | read | `query`; *fields, from, size* |  | `Files.Read.All`, `Sites.Read.All` |
| `ONE_DRIVE_SEARCH_ITEMS` | Search Items | read | `q`; *drive_id, expand, orderby, page_token, search_scope, select, stripped_annotations, top and 4 more* |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_SEARCH_SHAREPOINT_SITES` | Search SharePoint Sites | read | *order_by, page_size, page_token, query* |  | `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_TRACK_FOLDER_CHANGES` | Track Folder Changes | read | `folder_id`; *drive_id, expand_relations, page_token, select_fields, top* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_UNFOLLOW_ITEM` | Delete Drive Following | write | `id` |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_UPDATE_DRIVE_ITEMS_PERMISSIONS` | Update Drive Item Permissions | write | `item_id`, `permission_id`, `roles`; *drive_id, group_id, site_id, user_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_UPDATE_DRIVE_ITEM_METADATA` | Update Drive Item Metadata | write | `item_id`; *description, drive_id, file_system_info, group_id, if_match, name, parent_reference_drive_id, parent_reference_id and 2 more* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_UPDATE_FILE_CONTENT` | Update File Content | write | `item_id`; *conflict_behavior, defer_commit, drive_id, drive_item_source, file, file_system_info, group_id, if_match_etag and 4 more* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_UPDATE_SHAREPOINT_ITEM_PERMISSION` | Update SharePoint Item Permission | write | `list_id`, `permission_id`, `roles`, `site_id`; *item_id* |  | `Sites.ReadWrite.All`, `Sites.Manage.All`, `Sites.FullControl.All` |
| `ONE_DRIVE_UPDATE_SHAREPOINT_LIST_ITEM_FIELDS` | Update SharePoint List Item Fields | write | `fields`, `item_id`, `list_id`, `site_id`; *if_match* |  | `Sites.ReadWrite.All` |
| `ONE_DRIVE_UPDATE_TABLE_COLUMN_RANGE` | Update Table Column Range | write | `column_index`, `drive_id`, `item_id`, `table_id_or_name`; *column_hidden, formulas, formulas_local, formulas_r1c1, number_format, range_kind, row_hidden, session_id and 1 more* |  | `Files.ReadWrite` |
| `ONE_DRIVE_UPDATE_WORKBOOK_RANGE` | Update Workbook Range | write | `address`, `drive_id`, `item_id`, `worksheet_id_or_name`; *formulas, number_format, session_id, values* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_UPDATE_WORKBOOK_WORKSHEET` | Update Workbook Worksheet | write | `drive_id`, `item_id`, `worksheet_id_or_name`; *name, position, session_id, visibility* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `ONE_DRIVE_WRITE_WORKBOOK_CELL` | Write Workbook Cell | write | `column`, `drive_id`, `item_id`, `row`, `worksheet_id_or_name`; *formulas, number_format, session_id, values* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
