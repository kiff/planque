# Dropbox: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

175 tools: 93 read, 82 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `DROPBOX_ACTIVATE_TEAM_FOLDER` | Activate team folder | write | `team_folder_id` |  | `team_data.content.write` |
| `DROPBOX_ADD_FILE_MEMBER` | Add file member | write | `file`, `members`; *access_level, add_message_as_comment, custom_message, quiet* |  | `sharing.write` |
| `DROPBOX_ADD_FILE_PROPERTIES` | Add file properties | write | `path`, `property_groups`; *Dropbox-API-Select-User* |  | `files.metadata.write` |
| `DROPBOX_ADD_FILE_TAGS` | Add tag to file or folder | write | `path`, `tag_text`; *Dropbox-API-Select-User* |  | `files.metadata.write` |
| `DROPBOX_ADD_FOLDER_MEMBER_ACTION` | Add folder member | write | `members`, `shared_folder_id`; *custom_message, quiet* |  | `sharing.write` |
| `DROPBOX_ADD_MEMBER_SPACE_LIMITS_EXCLUDED_USERS` | Add users to space limits exclusion list | write | `users` |  | `members.write` |
| `DROPBOX_ADD_SHARING_ALLOWLIST` | Add sharing allowlist | write | *domains, emails* |  | `team_info.write` |
| `DROPBOX_ADD_TEAM_GROUP_MEMBERS` | Add members to team group | write | `group`, `members`; *return_members* |  | `groups.write` |
| `DROPBOX_ADD_TEAM_MEMBERS` | Add team members | write | `new_members`; *force_async* |  | `members.write` |
| `DROPBOX_ADD_TEAM_MEMBERS_SECONDARY_EMAILS` | Add secondary emails to team members | write | `new_secondary_emails` |  | `members.write` |
| `DROPBOX_ADD_TEAM_PROPERTIES_TEMPLATE` | Add property template for team | write | `description`, `fields`, `name` |  | `files.team_metadata.write` |
| `DROPBOX_ALPHA_UPLOAD_FILE` | Alpha Upload File | write | `content`, `path`; *autorename, client_modified, dropbox_api_select_user, mode, mute, property_groups* |  | `files.content.write` |
| `DROPBOX_APPEND_UPLOAD_SESSION` | Append data to upload session | write | `content`, `cursor`; *close* |  | `files.content.write` |
| `DROPBOX_APPEND_UPLOAD_SESSION_BATCH` | Append to multiple upload sessions | write | `data`, `entries` |  | `files.content.write` |
| `DROPBOX_ARCHIVE_TEAM_FOLDER` | Archive team folder | write | `team_folder_id`; *force_async_off* |  | `team_data.content.write` |
| `DROPBOX_CHECK_COPY_BATCH` | Check copy batch job status | read | `async_job_id` |  | `files.content.write` |
| `DROPBOX_CHECK_DELETE_BATCH` | Check delete batch status | read | `async_job_id` |  | `files.content.write` |
| `DROPBOX_CHECK_FOLDER_BATCH` | Check folder batch status | read | `async_job_id` |  | `files.content.write` |
| `DROPBOX_CHECK_JOB_STATUS` | Check sharing job status | read | `async_job_id` |  | `sharing.write` |
| `DROPBOX_CHECK_MOVE_BATCH` | Check move batch job status | read | `async_job_id` |  | `files.content.write` |
| `DROPBOX_CHECK_MOVE_FORMER_MEMBER_FILES_JOB_STATUS` | Check move former member files job status | read | `async_job_id` |  | `members.write` |
| `DROPBOX_CHECK_REMOVE_MEMBER` | Check remove member job status | read | `async_job_id` |  | `sharing.write` |
| `DROPBOX_CHECK_SAVE_URL_STATUS` | Check save URL job status | read | `async_job_id` |  | `files.content.write` |
| `DROPBOX_CHECK_SHARE_JOB_STATUS` | Check share job status | read | `async_job_id` |  | `sharing.write` |
| `DROPBOX_CHECK_TEAM_FOLDER_ARCHIVE` | Check team folder archive status | read | `async_job_id` |  | `team_data.team_space` |
| `DROPBOX_CHECK_UPLOAD_BATCH` | Check upload batch status | read | `async_job_id` |  | `files.content.write` |
| `DROPBOX_CHECK_USER` | Check user | read | *query* |  | `account_info.read` |
| `DROPBOX_COPY_BATCH` | Copy multiple files or folders | write | `entries`; *autorename, dropbox_api_select_user* |  | `files.content.write` |
| `DROPBOX_COPY_FILE_OR_FOLDER` | Copy file or folder | write | `from_path`, `to_path`; *allow_ownership_transfer, allow_shared_folder, autorename, dropbox_api_select_user* |  | `files.content.write` |
| `DROPBOX_COUNT_FILE_REQUESTS` | Count file requests | read |  |  | `file_requests.read` |
| `DROPBOX_CREATE_FILE_REQUEST` | Create file request | write | `destination`, `title`; *deadline, open* |  | `file_requests.write` |
| `DROPBOX_CREATE_FOLDER` | Create folder | write | `path`; *autorename* |  | `files.content.write` |
| `DROPBOX_CREATE_FOLDER_BATCH` | Create multiple folders | write | `paths`; *autorename, force_async* |  | `files.content.write` |
| `DROPBOX_CREATE_PAPER_DOCUMENT` | Create paper document | write | `content`, `path`; *import_format* |  | `files.content.write` |
| `DROPBOX_CREATE_PAPER_FOLDER` | Create Paper folder | write | `name`; *is_team_folder, parent_folder_id* |  | `files.content.write` |
| `DROPBOX_CREATE_SHARED_LINK` | Create shared link | write | `path`; *settings* |  | `sharing.write` |
| `DROPBOX_CREATE_TEAM_FOLDER` | Create team folder | write | `name`; *sync_setting* |  | `team_data.member` |
| `DROPBOX_CREATE_TEAM_GROUP` | Create team group | write | `group_name`; *group_external_id, group_management_type* |  | `groups.write` |
| `DROPBOX_DELETE_ALL_CLOSED_FILE_REQUESTS` | Delete all closed file requests | write |  |  | `file_requests.write` |
| `DROPBOX_DELETE_BATCH` | Delete multiple files/folders | write | `entries` |  | `files.content.write` |
| `DROPBOX_DELETE_FILE` | Delete file or folder | write | `path`; *parent_rev* |  | `files.content.write` |
| `DROPBOX_DELETE_FILE_REQUESTS` | Delete file requests | write | `ids` |  | `file_requests.write` |
| `DROPBOX_DELETE_MANUAL_CONTACTS_BATCH` | Delete manual contacts batch | write | `email_addresses`; *dropbox_api_select_user* |  | `contacts.write` |
| `DROPBOX_DELETE_TEAM_FOLDER_PERMANENTLY` | Permanently delete archived team folder | write | `team_folder_id` |  | `team_data.content.write` |
| `DROPBOX_DELETE_TEAM_GROUP` | Delete team group | write | `selector_type`; *group_external_id, group_id* |  | `groups.write` |
| `DROPBOX_DELETE_TEAM_MEMBERS_SECONDARY_EMAILS` | Delete team members secondary emails | write | `emails_to_delete` |  | `members.write` |
| `DROPBOX_DELETE_TEAM_MEMBER_PROFILE_PHOTO` | Delete team member profile photo | write | `user` |  | `members.write` |
| `DROPBOX_DOWNLOAD_ZIP` | Download Folder as Zip | read | `path` |  | `files.content.read` |
| `DROPBOX_EXPORT_FILE` | Export File | read | `path`; *export_format, mode* |  | `files.content.read` |
| `DROPBOX_FILES_SEARCH` | Search files and folders | read | `query`; *max_results, mode, path, start* |  | `files.metadata.read` |
| `DROPBOX_FINISH_UPLOAD_SESSION` | Finish upload session | write | `commit`, `cursor`; *content* |  | `files.content.write` |
| `DROPBOX_FINISH_UPLOAD_SESSION_BATCH` | Finish batch upload sessions | write | `entries` |  | `files.content.write` |
| `DROPBOX_GET_ABOUT_ME` | Get about me | read |  |  | `account_info.read` |
| `DROPBOX_GET_ACCOUNT` | Get account | read | `account_id` |  | `sharing.read` |
| `DROPBOX_GET_ACCOUNT_BATCH` | Get account batch | read | `account_ids` |  | `sharing.read` |
| `DROPBOX_GET_AVAILABLE_TEAM_MEMBER_ROLES` | Get available team member roles | read |  |  | `members.read` |
| `DROPBOX_GET_COPY_REFERENCE` | Get copy reference | read | `path` |  | `files.content.write` |
| `DROPBOX_GET_FILE_LOCK_BATCH` | Get file lock batch | read | `entries` |  | `files.content.read` |
| `DROPBOX_GET_FILE_METADATA_BATCH` | Get shared file metadata batch | read | `files`; *actions* |  | `sharing.read` |
| `DROPBOX_GET_FILE_PREVIEW` | Get file preview | read | `path` |  | `files.content.read` |
| `DROPBOX_GET_FILE_REQUEST` | Get file request details | read | `id` |  | `file_requests.read` |
| `DROPBOX_GET_FILE_TAGS` | Get file or folder tags | read | `paths` |  | `files.metadata.read` |
| `DROPBOX_GET_FOLDER_CURSOR` | Get folder state cursor | read | *include_deleted, include_has_explicit_shared_members, include_media_info, include_mounted_folders, include_property_groups, limit, path, recursive and 1 more* |  | `files.metadata.read` |
| `DROPBOX_GET_JWKS` | Get JWKS | read |  |  |  |
| `DROPBOX_GET_METADATA` | Get file or folder metadata | read | `path`; *include_deleted, include_has_explicit_shared_members, include_media_info* |  | `files.metadata.read` |
| `DROPBOX_GET_METADATA_ALPHA` | Get file or folder metadata (alpha) | read | `path`; *include_deleted, include_has_explicit_shared_members, include_media_info, include_property_groups* |  | `files.metadata.read` |
| `DROPBOX_GET_OPENID_CONFIG` | Get OpenID config | read |  |  |  |
| `DROPBOX_GET_SHARED_FILE_METADATA` | Get shared file metadata | read | `file`; *actions* |  | `sharing.read` |
| `DROPBOX_GET_SHARED_FOLDER_METADATA` | Get shared folder metadata | read | `shared_folder_id`; *actions* |  | `sharing.read` |
| `DROPBOX_GET_SHARED_LINK_FILE` | Download file from shared link | read | `url`; *link_password, path* |  | `sharing.read` |
| `DROPBOX_GET_SHARED_LINK_METADATA` | Get shared link metadata | read | `url`; *link_password, path* |  | `sharing.read` |
| `DROPBOX_GET_SPACE_USAGE` | Get space usage | read |  |  | `account_info.read` |
| `DROPBOX_GET_TEAM_FEATURE_VALUES` | Get team feature values | read | `features` |  | `team_info.read` |
| `DROPBOX_GET_TEAM_FOLDER_INFO` | Get team folder info | read | `team_folder_ids` |  | `team_data.content.read` |
| `DROPBOX_GET_TEAM_GROUPS_INFO` | Get team groups info | read | `selector_type`; *group_external_ids, group_ids* |  | `groups.read` |
| `DROPBOX_GET_TEAM_GROUPS_JOB_STATUS` | Get team groups job status | read | `async_job_id` |  | `groups.write` |
| `DROPBOX_GET_TEAM_INFO` | Get team info | read |  |  | `team_info.read` |
| `DROPBOX_GET_TEAM_LOG_EVENTS` | Get team audit log events | read | *category, cursor, limit, time* |  | `events.read` |
| `DROPBOX_GET_TEAM_LOG_EVENTS_CONTINUE` | Get team log events (continue) | read | `cursor` |  | `events.read` |
| `DROPBOX_GET_TEAM_MEMBERS_ADD_JOB_STATUS` | Get team members add job status v2 | read | `async_job_id` |  | `members.write` |
| `DROPBOX_GET_TEAM_MEMBERS_INFO` | Get team members info | read | `members` |  | `members.read` |
| `DROPBOX_GET_TEAM_MEMBER_CUSTOM_QUOTA` | Get team member custom quota | write | *users* |  | `members.read` |
| `DROPBOX_GET_TEAM_PROPERTIES_TEMPLATE` | Get property template for team | read | `template_id` |  | `files.team_metadata.read`, `files.team_metadata.write` |
| `DROPBOX_GET_TEMPORARY_LINK` | Get temporary link | read | `path` |  | `files.content.read` |
| `DROPBOX_GET_TEMPORARY_UPLOAD_LINK` | Get temporary upload link | write | `commit_info`; *duration* |  | `files.content.write` |
| `DROPBOX_GET_THUMBNAIL` | Get thumbnail | read | `path`; *format, mode, size* |  | `files.content.read` |
| `DROPBOX_GET_THUMBNAIL_BATCH` | Get thumbnail batch | read | `entries` |  | `files.content.read` |
| `DROPBOX_GET_THUMBNAIL_V2` | Get thumbnail v2 | read | `resource_type`; *format, link_path, mode, path, size, url* |  | `files.content.read` |
| `DROPBOX_GET_USER_FEATURES` | Get user features | read | `features` |  | `account_info.read` |
| `DROPBOX_LIST_EXCLUDED_USERS_CONTINUE` | List member space limits excluded users (continue) | read | `cursor` |  | `members.read` |
| `DROPBOX_LIST_FILES_IN_FOLDER` | List files in folder | read | *include_deleted, include_media_info, include_non_downloadable_files, limit, path, recursive* |  | `files.metadata.read` |
| `DROPBOX_LIST_FILE_MEMBERS` | List file members | read | `file`; *actions, include_inherited, limit* |  | `sharing.read` |
| `DROPBOX_LIST_FILE_MEMBERS_BATCH` | List file members batch | read | `files`; *limit* |  | `sharing.read` |
| `DROPBOX_LIST_FILE_REQUESTS` | List file requests | read | *limit* |  | `file_requests.read` |
| `DROPBOX_LIST_FILE_REQUESTS_CONTINUE` | List file requests (continue) | read | `cursor` |  | `file_requests.read` |
| `DROPBOX_LIST_FILE_REVISIONS` | List file revisions | read | `path`; *limit, mode* |  | `files.metadata.read` |
| `DROPBOX_LIST_FOLDERS_CONTINUE` | List shared folders continue | read | `cursor` |  | `sharing.read` |
| `DROPBOX_LIST_FOLDER_CONTINUE` | List folder continue | read | `cursor` |  | `files.metadata.read` |
| `DROPBOX_LIST_FOLDER_MEMBERS` | List folder members | read | `shared_folder_id`; *actions, limit* |  | `sharing.read` |
| `DROPBOX_LIST_FOLDER_MEMBERS_CONTINUE` | List folder members continue | read | `cursor` |  | `sharing.read` |
| `DROPBOX_LIST_MEMBER_LINKED_APPS` | List member linked apps | read | `team_member_id` |  | `sessions.list` |
| `DROPBOX_LIST_MEMBER_SPACE_LIMITS_EXCLUDED_USERS` | List excluded users from space limits | read | *limit* |  | `members.read` |
| `DROPBOX_LIST_MOUNTABLE_FOLDERS` | List mountable folders | read | *actions, limit* |  | `sharing.read` |
| `DROPBOX_LIST_MOUNTABLE_FOLDERS_CONTINUE` | List mountable folders continue | read | `cursor` |  | `sharing.read` |
| `DROPBOX_LIST_PAPER_DOCS` | List Paper docs | read | *filter_by, limit, sort_by, sort_order* |  | `files.metadata.read` |
| `DROPBOX_LIST_PAPER_DOCS_CONTINUE` | List Paper documents (continue) | read | `cursor` |  | `files.metadata.read` |
| `DROPBOX_LIST_RECEIVED_FILES` | List received files | read | *actions, limit* |  | `sharing.read` |
| `DROPBOX_LIST_SHARED_FOLDERS` | List shared folders | read | *limit* |  | `sharing.read` |
| `DROPBOX_LIST_SHARED_LINKS` | List shared links | read | *cursor, direct_only, path* |  | `sharing.read` |
| `DROPBOX_LIST_SHARING_ALLOWLIST` | List team sharing allowlist | read | *limit* |  | `team_info.read` |
| `DROPBOX_LIST_TEAM_DEVICES` | List team devices | read | *cursor, include_desktop_clients, include_mobile_clients, include_web_sessions* |  | `sessions.list` |
| `DROPBOX_LIST_TEAM_FOLDERS` | List team folders | read | *limit* |  | `team_data.content.read` |
| `DROPBOX_LIST_TEAM_FOLDERS_CONTINUE` | List team folders continue | read | `cursor` |  | `team_data.content.read` |
| `DROPBOX_LIST_TEAM_GROUPS` | List team groups | read | *limit* |  | `groups.read` |
| `DROPBOX_LIST_TEAM_GROUPS_CONTINUE` | List team groups continue | read | `cursor` |  | `groups.read` |
| `DROPBOX_LIST_TEAM_GROUPS_MEMBERS_CONTINUE` | List team group members continue | read | `cursor` |  | `groups.read` |
| `DROPBOX_LIST_TEAM_GROUP_MEMBERS` | List team group members | read | `selector_type`; *group_external_id, group_id, limit* |  | `groups.read` |
| `DROPBOX_LIST_TEAM_LINKED_APPS` | List members linked apps | read | *cursor* |  | `sessions.list` |
| `DROPBOX_LIST_TEAM_MEMBERS` | List team members | read | *include_removed, limit* |  | `members.read` |
| `DROPBOX_LIST_TEAM_MEMBERS_CONTINUE` | List team members continue (v2) | read | `cursor` |  | `members.read` |
| `DROPBOX_LIST_TEAM_MEMBER_DEVICES` | List team member devices | read | `team_member_id`; *include_desktop_clients, include_mobile_clients, include_web_sessions* |  | `sessions.list` |
| `DROPBOX_LIST_TEAM_NAMESPACES` | List team namespaces | read | *limit* |  | `team_data.member` |
| `DROPBOX_LIST_TEAM_NAMESPACES_CONTINUE` | List team namespaces continue | read | `cursor` |  | `team_data.member` |
| `DROPBOX_LIST_TEAM_SHARING_ALLOWLIST_CONTINUE` | List team sharing allowlist continue | read | `cursor` |  | `team_info.read` |
| `DROPBOX_LIST_TEMPLATES_FOR_TEAM` | List property templates for team | read |  |  | `files.team_metadata.write` |
| `DROPBOX_LIST_USER_TEMPLATES` | List user templates | read |  |  | `files.metadata.read` |
| `DROPBOX_MODIFY_SHARED_LINK_SETTINGS` | Modify shared link settings | write | `settings`, `url`; *remove_expiration* |  | `sharing.write` |
| `DROPBOX_MOUNT_FOLDER` | Mount shared folder | write | `shared_folder_id` |  | `sharing.write` |
| `DROPBOX_MOVE_BATCH` | Move multiple files or folders in batch | write | `entries`; *allow_ownership_transfer, autorename* |  | `files.content.write` |
| `DROPBOX_MOVE_FILE_OR_FOLDER` | Move file or folder | write | `from_path`, `to_path`; *allow_ownership_transfer, allow_shared_folder, autorename* |  | `files.content.write` |
| `DROPBOX_OVERWRITE_FILE_PROPERTIES` | Overwrite file properties | write | `path`, `property_groups`; *Dropbox-API-Select-User* |  | `files.metadata.write` |
| `DROPBOX_READ_FILE` | Read a file | read | `path` |  | `files.content.read` |
| `DROPBOX_REMOVE_FILE_MEMBER` | Remove file member | write | `file`, `member` |  | `sharing.write` |
| `DROPBOX_REMOVE_FILE_PROPERTIES` | Remove File Properties | write | `path`, `property_template_ids`; *Dropbox-API-Select-User* |  | `files.metadata.write` |
| `DROPBOX_REMOVE_FILE_PROPERTIES_TEMPLATE_FOR_TEAM` | Remove file properties template for team | write | `template_id` |  | `files.team_metadata.write` |
| `DROPBOX_REMOVE_FILE_TAG` | Remove tag from file or folder | write | `path`, `tag_text`; *Dropbox-API-Select-User* |  | `files.metadata.write` |
| `DROPBOX_REMOVE_FOLDER_MEMBER` | Remove folder member | write | `member`, `shared_folder_id`; *leave_a_copy* |  | `sharing.write` |
| `DROPBOX_REMOVE_GROUP_MEMBERS` | Remove group members | write | `group`, `users`; *return_members* |  | `groups.write` |
| `DROPBOX_REMOVE_SHARING_ALLOWLIST` | Remove sharing allowlist | write | *domains, emails* |  | `team_info.write` |
| `DROPBOX_REMOVE_TEAM_MEMBER_CUSTOM_QUOTA` | Remove team member custom quota | write | `users` |  | `members.write` |
| `DROPBOX_REMOVE_TEAM_MEMBER_SPACE_LIMITS_EXCLUDED_USERS` | Remove member space limits excluded users | write | `users` |  | `members.write` |
| `DROPBOX_RENAME_TEAM_FOLDER` | Rename team folder | write | `name`, `team_folder_id` |  | `team_data.content.write` |
| `DROPBOX_RESEND_SECONDARY_EMAIL_VERIFICATION` | Resend secondary email verification | write | `emails_to_resend` |  | `members.write` |
| `DROPBOX_RESTORE_FILE` | Restore file to specific revision | write | `path`, `rev` |  | `files.content.write` |
| `DROPBOX_REVOKE_SHARED_LINK` | Revoke shared link | write | `url` |  | `sharing.write` |
| `DROPBOX_SAVE_COPY_REFERENCE` | Save copy reference | write | `copy_reference`, `path` |  | `files.content.write` |
| `DROPBOX_SAVE_URL` | Save file from URL | write | `path`, `url`; *Dropbox-API-Select-User, timeout_seconds, wait* |  | `files.content.write` |
| `DROPBOX_SEARCH_CONTINUE` | Continue Dropbox Search | read | `cursor` |  | `files.metadata.read` |
| `DROPBOX_SEARCH_FILE_OR_FOLDER` | Search File or Folder | read | *Dropbox-API-Select-User, cursor, match_field_options, options, query* |  | `files.metadata.read` |
| `DROPBOX_SEARCH_FILE_PROPERTIES` | Search File Properties | read | *cursor, queries, template_filter* |  | `files.metadata.read` |
| `DROPBOX_SEND_TEAM_MEMBER_WELCOME_EMAIL` | Send team member welcome email | write | `.tag`; *email, external_id, team_member_id* |  | `members.write` |
| `DROPBOX_SET_PROFILE_PHOTO` | Set profile photo | write | `photo`; *Dropbox-API-Select-User* |  | `account_info.write` |
| `DROPBOX_SET_TEAM_MEMBER_CUSTOM_QUOTA` | Set team member custom quota | write | `users_and_quotas` |  | `members.write` |
| `DROPBOX_SHARE_FOLDER` | Share folder | write | `path`; *acl_update_policy, actions, force_async, link_settings, member_policy, shared_link_policy, viewer_info_policy* |  | `sharing.write` |
| `DROPBOX_START_UPLOAD_SESSION` | Start upload session | write | `content`; *close* |  | `files.content.write` |
| `DROPBOX_START_UPLOAD_SESSION_BATCH` | Start batch of upload sessions | write | `num_sessions`; *session_type* |  | `files.content.write` |
| `DROPBOX_UNLOCK_FILE_BATCH` | Unlock multiple files | write | `entries` |  | `files.content.write` |
| `DROPBOX_UNMOUNT_FOLDER` | Unmount shared folder | write | `shared_folder_id` |  | `sharing.write` |
| `DROPBOX_UNSHARE_FILE` | Unshare file | write | `file` |  | `sharing.write` |
| `DROPBOX_UNSHARE_FOLDER` | Unshare folder | write | `shared_folder_id`; *leave_a_copy* |  | `sharing.write` |
| `DROPBOX_UPDATE_FILE_MEMBER` | Update file member | write | `access_level`, `file`, `member` |  | `sharing.write` |
| `DROPBOX_UPDATE_FILE_PROPERTIES` | Update file properties | write | `path`, `update_property_groups`; *Dropbox-API-Select-User* |  | `files.metadata.write` |
| `DROPBOX_UPDATE_FILE_REQUEST` | Update file request | write | `id`; *deadline, description, destination, open, title* |  | `file_requests.write` |
| `DROPBOX_UPDATE_FOLDER_MEMBER` | Update folder member | write | `access_level`, `member`, `shared_folder_id` |  | `sharing.write` |
| `DROPBOX_UPDATE_FOLDER_POLICY` | Update folder policy | write | `shared_folder_id`; *acl_update_policy, actions, link_settings, member_policy, shared_link_policy, viewer_info_policy* |  | `sharing.write` |
| `DROPBOX_UPDATE_GROUP_MEMBER_ACCESS_TYPE` | Set group member access type | write | `access_type`, `group`, `user`; *return_members* |  | `groups.write` |
| `DROPBOX_UPDATE_PAPER_DOCUMENT` | Update Paper Document | write | `content`, `doc_update_policy`, `import_format`, `path`; *paper_revision* |  | `files.content.write` |
| `DROPBOX_UPDATE_PROPERTY_TEMPLATE_FOR_TEAM` | Update property template for team | write | `template_id`; *add_fields, description, name* |  | `files.team_metadata.write` |
| `DROPBOX_UPDATE_TEAM_FOLDER_SYNC_SETTINGS` | Update team folder sync settings | write | `sync_setting`, `team_folder_id`; *content_sync_settings* |  | `team_data.team_space` |
| `DROPBOX_UPDATE_TEAM_GROUP` | Update team group | write | `group`; *new_group_external_id, new_group_management_type, new_group_name, return_members* |  | `groups.write` |
| `DROPBOX_UPDATE_TEAM_MEMBER_ADMIN_PERMISSIONS` | Set team member admin permissions | write | `new_role`, `user` |  | `members.write` |
| `DROPBOX_UPDATE_TEAM_MEMBER_PROFILE` | Set team member profile | write | `user`; *new_email, new_external_id, new_given_name, new_is_directory_restricted, new_persistent_id, new_surname* |  | `members.write` |
| `DROPBOX_UPDATE_TEAM_MEMBER_PROFILE_PHOTO` | Set team member profile photo | write | `photo`, `user` |  | `members.write` |
| `DROPBOX_UPLOAD_FILE` | Upload File | write | `content`, `path`; *autorename, mode, mute, strict_conflict* |  | `files.content.write` |
| `DROPBOX_WHO_AM_I` | Who Am I | read |  |  |  |
