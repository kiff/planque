# Wrike: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

144 tools: 75 read, 69 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `WRIKE_BULK_MODIFY_GROUPS` | Bulk modify group members | write | `members` |  | `amReadWriteGroup` |
| `WRIKE_COPY_FOLDER` | Copy folder | write | `folderId`, `parent`, `title`; *addResponsibles, copyCustomFields, copyCustomStatuses, copyDescriptions, copyParents, copyResponsibles, copyStatuses, entryLimit and 3 more* |  | `Default`, `wsReadWrite` |
| `WRIKE_COPY_FOLDER_ASYNC` | Copy folder async | write | `folder_id`, `parent`, `title`; *add_responsibles, copy_custom_fields, copy_custom_statuses, copy_descriptions, copy_parents, copy_responsibles, copy_statuses, entry_limit and 3 more* |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_ACCOUNT_WEBHOOKS` | Create account webhooks | write | `hookUrl`; *customItemTypes, events, parameterisedEvents, secret* |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_ASSETS` | Create equipment asset | write | `name` |  | `amReadWriteUser` |
| `WRIKE_CREATE_CUSTOMFIELD` | Create custom field | write | `title`, `type`; *aggregation, allowOtherValues, contacts, currency, decimalPlaces, inheritanceType, readOnly, shareds and 3 more* |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_FOLDER` | Create a folder | write | `folderId`, `title`; *customColumns, customFields, customItemTypeId, description, metadata, plainTextCustomFields, shareds, userAccessRoles and 1 more* |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_FOLDERS_APPROVALS` | Create approval in folder | write | `folderId` |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_FOLDERS_COMMENTS` | Create folder comment | write | `folderId`, `text` |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_FOLDERS_TIMELOG_LOCK_PERIOD` | Create timelog lock period for folder | write | `folderId`, `period` |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_FOLDER_BOOKING` | Create folder booking | write | `bookingDates`, `folderId`; *effortAllocation, responsibleId, responsiblePlaceholderId* |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_FOLDER_WEBHOOK` | Create a webhook on folder | write | `folder_id`, `hookUrl`; *events, recursive, secret* |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_GROUP` | Create a group | write | `title`; *avatar__color, avatar__letters, members, metadata, parent* |  | `amReadOnlyGroup`, `amReadWriteGroup` |
| `WRIKE_CREATE_INVITATION` | Create invitation | write | `email`; *external, firstName, lastName, message, role, subject, userTypeId* |  | `amReadWriteInvitation` |
| `WRIKE_CREATE_SPACES_TIMELOG_LOCK_PERIOD` | Create space timelog lock period | write | `period`, `spaceId` |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_SPACE_WEBHOOK` | Create space webhook | write | `hook_url`, `space_id`; *events, recursive, secret* |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_TASK` | Create task in folder | write | `folderId`, `title`; *billingType, customFields, customItemTypeId, customStatus, description, follow, followers, importance and 11 more* |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_TASKS_APPROVALS` | Create task approval | write | `taskId` |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_TASKS_DEPENDENCIES` | Create task dependency | write | `predecessorId`, `relationType`, `taskId` |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_TASKS_TIMELOGS` | Create task timelog | write | `hours`, `taskId`, `trackedDate`; *categoryId, comment* |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_TASK_ATTACHMENT` | Create task attachment | write | `file`, `taskId` |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_TASK_COMMENT` | Create task comment | write | `taskId`, `text` |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_TIMESHEETS` | Create timesheet | write | `periodStartDate` |  | `Default`, `wsReadWrite` |
| `WRIKE_CREATE_WORKFLOWS` | Create workflows | write | `name` |  | `amReadWriteWorkflow` |
| `WRIKE_DELETE_APPROVALS` | Delete approval | write | `approvalId` |  | `Default`, `wsReadWrite` |
| `WRIKE_DELETE_ASSET` | Delete asset | write | `assetId` |  | `amReadWriteUser` |
| `WRIKE_DELETE_ATTACHMENTS` | Delete attachment | write | `attachment_id` |  | `Default`, `wsReadWrite` |
| `WRIKE_DELETE_BOOKING` | Delete booking | write | `bookingId` |  | `Default`, `wsReadWrite` |
| `WRIKE_DELETE_COMMENT` | Delete comment | write | `commentId` |  | `Default`, `wsReadWrite` |
| `WRIKE_DELETE_CUSTOM_FIELD` | Delete custom field by id | write | `customFieldId` |  | `Default` |
| `WRIKE_DELETE_DEPENDENCIES` | Delete dependencies | write | `dependencyId` |  | `Default`, `wsReadWrite` |
| `WRIKE_DELETE_FOLDER` | Delete folder | write | `folderId` |  | `wsReadWrite` |
| `WRIKE_DELETE_FOLDERS_TIMELOG_LOCK_PERIOD` | Delete timelog lock period for folder | write | `folderId`, `period` |  | `Default`, `wsReadWrite` |
| `WRIKE_DELETE_GROUP` | Delete group by id | write | `groupId` |  | `amReadWriteGroup` |
| `WRIKE_DELETE_INIVTATION` | Delete invitation | write | `invitationId` |  | `amReadWriteInvitation` |
| `WRIKE_DELETE_JOB_ROLE` | Delete job role | write | `jobroleId` |  | `amReadWriteUser` |
| `WRIKE_DELETE_SPACE` | Delete space | write | `spaceId` |  | `Default`, `wsReadWrite` |
| `WRIKE_DELETE_SPACES_TIMELOG` | Delete spaces timelog lock period | write | `period`, `spaceId` |  | `Default`, `wsReadWrite` |
| `WRIKE_DELETE_TASK` | Delete task | write | `taskId` |  | `Default`, `wsReadWrite` |
| `WRIKE_DELETE_TIMELOG` | Delete timelog | write | `timelogId` |  | `Default`, `wsReadWrite` |
| `WRIKE_DELETE_WEBHOOK` | Delete webhook | write | `webhook_id` |  | `Default`, `wsReadWrite` |
| `WRIKE_FETCH_ALL_TASKS` | Fetch all tasks | read | *descendants, due_date, fields, metadata, page_size, resolve_user_names, status* |  |  |
| `WRIKE_GET_ACCOUNT_INFORMATION` | Get account information | read | *fields* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_ALL_CUSTOM_FIELDS` | Get all custom fields | write |  |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_ALL_WEBHOOKS` | Get all webhooks | read |  |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_APPROVALS` | Get approvals | read |  |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_APPROVALS_BY_IDS` | Get approvals by IDs | read | `approval_ids` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_ASYNC_JOB` | Get async job status | read | `async_job_id` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_ATTACHMENTS` | Get attachments by IDs | read | `attachment_ids` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_ATTACHMENTS_DOWNLOAD` | Download attachment | read | `attachmentId` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_ATTACHMENTS_PREVIEW` | Get attachment preview | read | `attachmentId` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_ATTACHMENTS_URL` | Get attachment URL | read | `attachmentId` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_BOOKINGS` | Get bookings | read | `bookingId` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_COLORS` | Get colors | read |  |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_COMMENTS` | Get comments | read | *plainText, updatedDate* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_COMMENTS_BY_IDS` | Get comments by IDs | read | `comment_ids`; *plainText* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_CONTACTS` | Get contacts | read | *deleted, fields, metadata* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_CONTACTS_HISTORY` | Get contacts history | read | `contactIds` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_CONTACT_HOURLY_RATES` | Get contact hourly rates | read | `contactIds` |  | `Default`, `wsReadWrite` |
| `WRIKE_GET_CONTACT_TIMELOGS` | Get contact timelogs | read | `contactId`; *createdDate, descendants, me, plainText, subTasks, trackedDate* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_CUSTOM_ITEM_TYPES` | Get custom item types | write |  |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_DEPENDENCIES` | Get dependencies by IDs | read | `dependency_ids` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_FOLDERS` | Get folders | read | *contractTypes, customFields, customItemTypes, deleted, descendants, fields, metadata, nextPageToken and 6 more* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_FOLDERS_APPROVALS` | Get folder approvals | read | `folder_id` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_FOLDERS_ATTACHMENTS` | Get folder attachments | read | `folderId`; *withUrls* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_FOLDERS_BY_IDS` | Get folders by IDs | read | `folder_ids` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_FOLDERS_COMMENTS` | Get folder comments | read | `folderId`; *plainText* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_FOLDERS_HISTORY` | Get folders history | read | `folderIds` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_FOLDERS_HOURLY` | Get folder hourly rates | read | `folderId` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_FOLDERS_ROLLUPS` | Get folder rollups | read | `folderId`; *descendants, nextPageToken, pageSize* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_FOLDERS_TASKS` | Get folder tasks | read | `folder_id`; *descendants, due_date, fields, metadata, page_size, status* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_FOLDERS_TIMELOGS` | Get folders timelogs | read | `folderId`; *createdDate, descendants, me, plainText, subTasks, trackedDate* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_FOLDERS_TIMELOG_LOCK_PERIODS` | Get folder timelog lock periods | read | `folderId` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_FOLDER_BOOKINGS` | Get folder bookings | read | `folderId` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_IDS` | Convert legacy v2 IDs to v4 | read | `ids`, `type` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_INFORMATION_ABOUT_SPECIFIC_CONTACT` | Get specific contact information | read | `contactId`; *fields* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_JOB_ROLES` | Get job roles by IDs | read | `jobrole_ids` |  | `amReadOnlyUser`, `amReadWriteUser` |
| `WRIKE_GET_PLACEHOLDERS` | Get placeholders by IDs | read | `placeholderIds` |  | `amReadOnlyUser`, `amReadWriteUser` |
| `WRIKE_GET_PLACEHOLDER_HOURLY_RATES` | Get placeholder hourly rates | read | `placeholderIds` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_SPACE` | Get space by ID | read | `spaceId`; *fields* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_SPACES_CUSTOMFIELDS` | Get space custom fields | read | `spaceId` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_SPACES_TASKS` | Get spaces tasks | read | `spaceId`; *descendants, fields, importance, sortField, sortOrder, status, subTasks, title* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_SPACES_WORKFLOWS` | Get space workflows | read | `spaceId` |  | `Default`, `amReadOnlyWorkflow`, `amReadWriteWorkflow`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_SPACE_CUSTOM_ITEM_TYPES` | Get space custom item types | write | `spaceId` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_SPACE_FOLDERS` | Get space folders | read | `spaceId`; *authors, customStatuses, fields, owners, project, title* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_SPACE_TIMELOG_LOCK_PERIODS` | Get space timelog lock periods | read | `spaceId` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_SPECIFIC_USER` | Get specific user | read | `userId` |  | `amReadOnlyUser`, `amReadWriteUser` |
| `WRIKE_GET_TASKS_APPROVALS` | Get task approvals | read | `taskId` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_TASKS_ATTACHMENTS` | Get task attachments | read | `taskId`; *withUrls* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_TASKS_COMMENTS` | Get task comments | read | `taskId`; *plainText* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_TASKS_HISTORY` | Get tasks history | read | `taskIds` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_TASKS_ROLLUPS` | Get task rollups | read | `taskId` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_TASKS_TIMELOG_LOCK_PERIODS` | Get task timelog lock periods | read | `taskId` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_TASK_BY_ID` | Get task by id | read | `taskId`; *fields* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_TASK_DEPENDENCIES` | Get task dependencies | read | `taskId` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_TASK_TIMELOGS` | Get task timelogs | read | `taskId`; *createdDate, descendants, me, plainText, subTasks, trackedDate* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_TIMELOGS` | Get timelogs | read | *createdDate, descendants, me, plainText, subTasks, trackedDate* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_TIMELOGS_BY_IDS` | Get timelogs by IDs | read | `timelogId`; *fields* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_TIMELOG_CATEGORIES` | Get timelog categories | read |  |  | `Default`, `amReadOnlyTimelogCategory`, `amReadWriteTimelogCategory`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_TIMESHEETS` | Get timesheets | read |  |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_TIMESHEET_SUBMISSION_RULES` | Get timesheet submission rules | read |  |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_VERSION` | Get API version | read |  |  |  |
| `WRIKE_GET_WEBHOOK_BY_ID` | Get webhook by ID | read | `webhook_id` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_GET_WORKSCHEDULES_TIMESHEET` | Get work schedule timesheet rules | read | `workschedule_id` |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_LAUNCH_FOLDER_BLUEPRINT_ASYNC` | Launch folder blueprint async | write | `folder_blueprint_id`, `parent`, `title`; *copy_attachments, copy_custom_fields, copy_descriptions, copy_responsibles, entry_limit, notify_responsibles, reschedule_date, reschedule_mode and 1 more* |  |  |
| `WRIKE_LAUNCH_TASK_BLUEPRINT_ASYNC` | Launch Task Blueprint Async | write | `task_blueprint_id`, `title`; *copy_attachments, copy_custom_fields, copy_descriptions, copy_responsibles, entry_limit, notify_responsibles, parent_id, reschedule_date and 3 more* |  |  |
| `WRIKE_LIST_ATTACHMENTS` | List all attachments | read | *createdDate, versions, withUrls* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_LIST_FOLDER_BLUEPRINTS` | List Folder Blueprints | read |  |  |  |
| `WRIKE_LIST_PLACEHOLDERS` | List all placeholders | read |  |  | `amReadOnlyUser`, `amReadWriteUser` |
| `WRIKE_LIST_SPACES` | List spaces | read | *withInvitations* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_LIST_SPACE_FOLDER_BLUEPRINTS` | List space folder blueprints | read | `space_id` |  |  |
| `WRIKE_LIST_SPACE_TASK_BLUEPRINTS` | List space task blueprints | read | `space_id`; *limit, next_page_token, page_size* |  |  |
| `WRIKE_LIST_SUBFOLDERS_BY_FOLDER_ID` | List subfolders by folder id | read | `folderId`; *contractTypes, customFields, customItemTypes, descendants, fields, metadata, nextPageToken, pageSize and 5 more* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_LIST_TASK_BLUEPRINTS` | List Task Blueprints | read | *limit, next_page_token, page_size* |  |  |
| `WRIKE_MODIFY_ACCOUNT` | Update account metadata | write | *metadata* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_MODIFY_FOLDER` | Modify folder attributes | write | `folderId`; *addAccessRoles, addParents, addShareds, clearCustomColumns, convertToCustomItemType, customColumns, customFields, description and 8 more* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_MODIFY_GROUP` | Modify group | write | `groupId`; *addInvitations, addMembers, avatar__color, avatar__letters, metadata, parent, removeInvitations, removeMembers and 1 more* |  | `amReadWriteGroup` |
| `WRIKE_MODIFY_TASK` | Modify task | write | `taskId`; *addFollowers, addParents, addResponsiblePlaceholders, addResponsibles, addShareds, addSuperTasks, billingType, convertToCustomItemType and 20 more* |  | `Default`, `wsReadWrite` |
| `WRIKE_QUERY_CUSTOMFIELDS` | Retrieve custom field by id | read | `customFieldId` |  |  |
| `WRIKE_QUERY_INVITATIONS` | Query invitations | read |  |  | `amReadOnlyInvitation`, `amReadWriteInvitation` |
| `WRIKE_QUERY_JOB_ROLES` | Query job roles | read |  |  | `amReadOnlyUser`, `amReadWriteUser` |
| `WRIKE_QUERY_SPECIFIC_GROUP` | Get group by id | read | `groupId`; *fields* |  | `amReadOnlyGroup`, `amReadWriteGroup` |
| `WRIKE_QUERY_WORKFLOWS` | Query workflows | read |  |  | `Default`, `wsReadOnly` |
| `WRIKE_RETRIEVE_LIST_OF_GROUPS` | Retrieve list of groups | read | *fields, metadata, pageSize, pageToken* |  | `amReadOnlyGroup`, `amReadWriteGroup` |
| `WRIKE_SEARCH_EDISCOVERY` | Search eDiscovery | read | *scopes, terms* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_UPDATE_APPROVALS` | Update approval | write | `approvalId`; *description, dueDate* |  | `Default`, `wsReadWrite` |
| `WRIKE_UPDATE_ASSET` | Update asset | write | `assetId`, `name` |  | `amReadWriteUser` |
| `WRIKE_UPDATE_ATTACHMENT` | Update attachment | write | `attachmentId`, `file_to_upload`; *contentType* |  | `Default`, `wsReadWrite` |
| `WRIKE_UPDATE_BOOKINGS` | Update booking | write | `bookingDates`, `bookingId` |  | `Default`, `wsReadWrite` |
| `WRIKE_UPDATE_COMMENT` | Update comment | write | `commentId`, `text` |  | `Default`, `wsReadWrite` |
| `WRIKE_UPDATE_CUSTOM_FIELD_BY_ID` | Update custom field by id | write | `customFieldId`; *addShareds, removeShareds, settings, sharing, spaceId, title, type* |  | `Default`, `wsReadWrite` |
| `WRIKE_UPDATE_DEPENDENCY` | Update dependency | write | `dependencyId`; *predecessorId, relationType, successorId* |  | `Default`, `wsReadWrite` |
| `WRIKE_UPDATE_FOLDERS_ROLLUPS` | Update folder rollup settings | write | `folderId`, `rollups` |  | `Default`, `wsReadWrite` |
| `WRIKE_UPDATE_INVITATION` | Update invitation | write | `invitationId`; *external, resend, role, userTypeId* |  | `amReadWriteInvitation` |
| `WRIKE_UPDATE_JOB_ROLE` | Update job role | write | `jobroleId`; *shortTitle, title* |  | `amReadWriteUser` |
| `WRIKE_UPDATE_METADATA_ON_SPECIFIC_CONTACT` | Update metadata on a specific contact | write | `contactId`; *customFields, fields, jobRoleId, metadata* |  | `Default`, `wsReadOnly`, `wsReadWrite` |
| `WRIKE_UPDATE_SPECIFIC_USER` | Update a specific user | write | `userId`; *active, profile__accountId, profile__external, profile__role, userTypeId* |  | `amReadWriteUser` |
| `WRIKE_UPDATE_TASKS_ROLLUPS` | Update task rollup settings | write | `rollups`, `taskId` |  | `Default`, `wsReadWrite` |
| `WRIKE_UPDATE_TIMELOGS` | Update timelog entry | write | `timelogId`; *categoryId, comment, hours, plainText, trackedDate* |  | `Default`, `wsReadWrite` |
| `WRIKE_UPDATE_TIMESHEETS` | Update timesheet | write | `approvalStatus`, `timesheetId` |  | `Default`, `wsReadWrite` |
| `WRIKE_UPDATE_TIMESHEET_ROW` | Update timesheet row | write | `timesheet_rowId`; *categoryId* |  | `Default`, `wsReadWrite` |
| `WRIKE_UPDATE_WEBHOOK_STATE` | Update webhook state | write | `status`, `webhookId` |  | `Default`, `wsReadWrite` |
| `WRIKE_UPDATE_WORKFLOW` | Update workflow | write | `workflowId`; *description, name* |  | `amReadWriteWorkflow` |
| `WRIKE_UPDATE_WORKSCHEDULE_TIMESHEET_RULES` | Update work schedule timesheet rules | write | `workschedule_id`; *enabled, frequency, rule_type, track_exceptions_mode* |  | `Default`, `wsReadWrite` |
