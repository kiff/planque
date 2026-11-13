# Microsoft Teams: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

237 tools: 126 read, 111 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `MICROSOFT_TEAMS_ADD_CHAT_MEMBER` | Add chat member | write | `chat_id`, `user_id`; *roles, visible_history_start_date_time* |  | `ChatMember.ReadWrite`, `Chat.ReadWrite`, `Chat.ReadWrite.All` |
| `MICROSOFT_TEAMS_ADD_GROUP_MEMBER` | Add User to Group | write | `group_id`, `user_id` |  |  |
| `MICROSOFT_TEAMS_ADD_MEMBER_TO_CHANNEL` | Add member to channel | write | `channel_id`, `team_id`, `user_id`; *roles* |  | `ChannelMember.ReadWrite.All`, `ChannelMember.ReadWrite.Group` |
| `MICROSOFT_TEAMS_ADD_TAB` | Add tab to channel | write | `channel_id`, `display_name`, `team_id`, `teams_app_id`; *configuration* |  | `TeamsTab.Create`, `Directory.ReadWrite.All`, `Group.ReadWrite.All`, `TeamsTab.ReadWrite.All`, `TeamsTab.ReadWriteForTeam`, `TeamsTab.ReadWriteSelfForTeam` |
| `MICROSOFT_TEAMS_ADD_TEAM_MEMBER` | Add member to team | write | `team_id`, `user_id`; *roles* |  | `TeamMember.ReadWriteNonOwnerRole.All`, `TeamMember.ReadWrite.All` |
| `MICROSOFT_TEAMS_ADD_TEAM_MEMBERS_BULK` | Add team members (bulk) | write | `team_id`, `values` |  | `TeamMember.ReadWriteNonOwnerRole.All`, `TeamMember.ReadWrite.All` |
| `MICROSOFT_TEAMS_ADD_WORKBOOK_TABLE_COLUMN` | Add workbook table column | write | `drive_id`, `index`, `item_id`, `table_id`; *name, session_id, values* |  |  |
| `MICROSOFT_TEAMS_ARCHIVE_CHANNEL` | Archive channel | write | `channel_id`, `team_id`; *should_set_spo_site_read_only_for_members* |  | `ChannelSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_ARCHIVE_GROUP_TEAM_CHANNEL` | Archive team channel | write | `channel_id`, `group_id`; *should_set_spo_site_read_only_for_members* |  | `ChannelSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_ARCHIVE_TEAM` | Archive Teams team | write | `team_id`; *should_set_spo_site_read_only_for_members* |  | `TeamSettings.ReadWrite.All`, `Directory.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_CHATS_GET_ALL_CHATS` | Get all chats | read | *filter, limit, page_token, select, top, user_id* |  | `Chat.ReadBasic`, `Chat.Read`, `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_CLEAR_AUTOMATIC_LOCATION` | Clear Automatic Location | write | `user_id` |  | `Presence.ReadWrite` |
| `MICROSOFT_TEAMS_CLEAR_ME_PRESENCE_USER_PREFERRED` | Clear User Preferred Presence | write | *user_id* |  | `Presence.ReadWrite` |
| `MICROSOFT_TEAMS_CLEAR_PRESENCE` | Clear Presence | write | `session_id`, `user_id` |  | `Presence.ReadWrite` |
| `MICROSOFT_TEAMS_CLEAR_PRESENCE_LOCATION` | Clear Presence Location | write | *user_id* |  | `Presence.ReadWrite` |
| `MICROSOFT_TEAMS_CLEAR_WORKBOOK_RANGE` | Clear Workbook Range | write | `address`, `drive_id`, `item_id`, `worksheet`; *apply_to, workbook_session_id* |  |  |
| `MICROSOFT_TEAMS_CLONE_TEAM` | Clone Team | write | `displayName`, `mailNickname`, `team_id`; *classification, description, partsToClone, visibility* |  | `Directory.ReadWrite.All`, `Group.ReadWrite.All`, `Team.Create`, `TeamSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_CREATE_CALL_OPERATION` | Create Call Operation | write | `call_id`; *client_context, odata_type, status* |  |  |
| `MICROSOFT_TEAMS_CREATE_CHANNEL` | Create a channel | write | `name`, `team_id`; *created_date_time, description, is_favorite_by_default, members, membership_type, shared_with_teams* |  | `Channel.Create`, `Group.ReadWrite.All`, `Directory.ReadWrite.All` |
| `MICROSOFT_TEAMS_CREATE_CONTENT_SHARING_SESSION` | Create Content Sharing Session | write | `call_id` |  |  |
| `MICROSOFT_TEAMS_CREATE_DRIVE_FOLDER` | Create drive folder | write | `drive_id`, `name`, `parent_item_id`; *conflict_behavior* |  |  |
| `MICROSOFT_TEAMS_CREATE_DRIVE_UPLOAD_SESSION` | Create Drive Upload Session | write | `drive_id`, `file_name`, `parent_item_id`; *conflict_behavior, defer_commit, if_match, if_none_match* |  |  |
| `MICROSOFT_TEAMS_CREATE_GROUP_CALENDAR_EVENT` | Create group calendar event | write | `end`, `group_id`, `start`, `subject`; *attendees, body, importance, is_all_day, location, sensitivity, show_as, transaction_id* |  |  |
| `MICROSOFT_TEAMS_CREATE_GROUP_TEAM_CHANNEL` | Create channel in group team | write | `display_name`, `group_id`; *description, is_favorite_by_default, membership_type* |  | `Channel.Create`, `Group.ReadWrite.All`, `Directory.ReadWrite.All` |
| `MICROSOFT_TEAMS_CREATE_MEETING` | Create online meeting | write | `end_date_time`, `start_date_time`, `subject`; *allow_attendee_to_enable_camera, allow_attendee_to_enable_mic, allow_breakout_rooms, allow_copying_and_sharing_meeting_content, allow_live_share, allow_meeting_chat, allow_participants_to_change_name, allow_power_point_sharing and 18 more* |  | `OnlineMeetings.ReadWrite` |
| `MICROSOFT_TEAMS_CREATE_OFFER_SHIFT_REQUEST` | Create user offer shift request | write | `recipient_user_id`, `sender_shift_id`, `team_id`; *odata_type, recipient_action_date_time, recipient_action_message, sender_message, user_id* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_CREATE_OPEN_SHIFT` | Create open shifts | write | `shared_open_shift`, `team_id`; *draft_open_shift, is_staged_for_deletion, scheduling_group_id* |  | `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_CREATE_OPEN_SHIFT_CHANGE_REQUEST` | Create open shift change request | write | `openShiftId`, `senderUserId`, `team_id`; *message* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_CREATE_OR_GET_ONLINE_MEETING` | Create or get online meeting | write | `externalId`; *chatInfo, endDateTime, participants, startDateTime, subject, userId* |  | `OnlineMeetings.ReadWrite` |
| `MICROSOFT_TEAMS_CREATE_OR_UPDATE_SCHEDULE` | Create or update schedule | write | `team_id`; *enabled, isActivitiesIncludedWhenCopyingShiftsEnabled, offerShiftRequestsEnabled, openShiftsEnabled, swapShiftsRequestsEnabled, timeClockEnabled, timeOffRequestsEnabled, timeZone* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_CREATE_PLANNER_BUCKET` | Create Planner bucket | write | `name`, `plan_id`; *order_hint* |  |  |
| `MICROSOFT_TEAMS_CREATE_PLANNER_TASK` | Create Planner task | write | `plan_id`, `title`; *assignee_ids, bucket_id* |  |  |
| `MICROSOFT_TEAMS_CREATE_SCHEDULE_DAY_NOTE` | Create schedule day note | write | `dayNoteDate`, `teamId`; *draftDayNote, sharedDayNote* |  | `Group.ReadWrite.All`, `Schedule.ReadWrite.All` |
| `MICROSOFT_TEAMS_CREATE_SCHEDULING_GROUP` | Create scheduling group | write | `display_name`, `team_id`, `user_ids`; *code, is_active* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_CREATE_SHIFT` | Create team schedule shift | write | `team_id`, `userId`; *draftShift, isStagedForDeletion, schedulingGroupId, sharedShift* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_CREATE_SUBSCRIPTION` | Create Graph subscription | write | `change_type`, `expiration_date_time`, `notification_url`, `resource`; *client_state, encryption_certificate, encryption_certificate_id, include_resource_data, latest_supported_tls_version, lifecycle_notification_url* |  |  |
| `MICROSOFT_TEAMS_CREATE_TEAM` | Create Team | write | `displayName`; *channels, classification, description, firstChannelName, funSettings, guestSettings, installedApps, memberSettings and 5 more* |  | `Team.Create`, `Group.ReadWrite.All`, `Directory.ReadWrite.All` |
| `MICROSOFT_TEAMS_CREATE_TEAM_FROM_GROUP` | Create Team From Group | write | `group_id`; *funSettings, guestSettings, memberSettings, messagingSettings* |  | `Directory.ReadWrite.All`, `Group.ReadWrite.All`, `Team.Create` |
| `MICROSOFT_TEAMS_CREATE_TIME_OFF` | Create time off | write | `team_id`, `userId`; *draftTimeOff, isStagedForDeletion, sharedTimeOff* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_CREATE_TIME_OFF_REASON` | Create time off reason | write | `displayName`, `iconType`, `isActive`, `team_id`; *code* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_CREATE_TIME_OFF_REQUEST` | Create time off request | write | `endDateTime`, `senderUserId`, `startDateTime`, `team_id`, `timeOffReasonId`; *odata_type* |  | `Group.ReadWrite.All`, `Schedule.ReadWrite.All` |
| `MICROSOFT_TEAMS_CREATE_USER_ONLINE_MEETING` | Create user online meeting | write | `endDateTime`, `startDateTime`, `subject`, `user_id`; *externalId, isBroadcast, meetingTemplateId* |  | `OnlineMeetings.ReadWrite` |
| `MICROSOFT_TEAMS_CREATE_WORKBOOK_SESSION` | Create workbook session | write | `drive_id`, `item_id`; *persist_changes* |  |  |
| `MICROSOFT_TEAMS_DELETE_CALL_OPERATION` | Delete Call Operation | write | `call_id`, `comms_operation_id`; *if_match* |  |  |
| `MICROSOFT_TEAMS_DELETE_CHANNEL` | Delete team channel | write | `channel_id`, `team_id`; *if_match* |  | `Channel.Delete.All`, `Directory.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_DELETE_DAY_NOTE` | Delete team schedule day note | write | `day_note_id`, `if_match`, `team_id` |  | `Schedule.ReadWrite.All` |
| `MICROSOFT_TEAMS_DELETE_DRIVE_ITEM` | Delete Drive Item | write | `drive_id`, `item_id`; *if_match, prefer* |  |  |
| `MICROSOFT_TEAMS_DELETE_OPEN_SHIFT` | Delete open shift | write | `open_shift_id`, `team_id`; *if_match* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_DELETE_PLANNER_TASK` | Delete Planner task | write | `etag`, `task_id` |  |  |
| `MICROSOFT_TEAMS_DELETE_SCHEDULING_GROUP` | Delete team scheduling group | write | `scheduling_group_id`, `team_id`; *if_match* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_DELETE_SHIFT` | Delete team schedule shift | write | `shift_id`, `team_id`; *if_match* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_DELETE_SOFT_MESSAGE` | Soft-delete Teams channel message | write | `channel_id`, `message_id`, `team_id` |  | `ChannelMessage.ReadWrite`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_DELETE_SUBSCRIPTION` | Delete Graph subscription | write | `subscription_id` |  |  |
| `MICROSOFT_TEAMS_DELETE_TAB` | Delete channel tab | write | `channel_id`, `tab_id`, `team_id`; *if_match* |  | `TeamsTab.ReadWriteSelfForTeam`, `TeamsTab.ReadWriteForTeam`, `TeamsTab.ReadWrite.All`, `Group.ReadWrite.All`, `Directory.ReadWrite.All` |
| `MICROSOFT_TEAMS_DELETE_TEAM` | Delete Teams team | write | `team_id` |  | `Directory.AccessAsUser.All`, `Directory.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_DELETE_TIME_OFF` | Delete time off | write | `team_id`, `timeoff_id`; *if_match* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_DELETE_TIME_OFF_REASON` | Delete time off reason | write | `team_id`, `time_off_reason_id`; *if_match* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_DELETE_TIME_OFF_REQUEST` | Delete time off request | write | `team_id`, `time_off_request_id`; *if_match* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_DELETE_USER_ONLINE_MEETING` | Delete User Online Meeting | write | `online_meeting_id`, `user_id`; *if_match* |  | `OnlineMeetings.ReadWrite` |
| `MICROSOFT_TEAMS_DOWNLOAD_CHANNEL_MESSAGE_HOSTED_CONTENT` | Download channel message hosted content | read | `channel_id`, `hosted_content_id`, `message_id`, `team_id`; *consistency_level* |  |  |
| `MICROSOFT_TEAMS_DOWNLOAD_CHAT_MESSAGE_HOSTED_CONTENT` | Download chat message hosted content | read | `chat_id`, `hosted_content_id`, `message_id`; *consistency_level* |  |  |
| `MICROSOFT_TEAMS_DOWNLOAD_DRIVE_ITEM` | Download drive item | read | `drive_id`; *format, if_none_match, item_id, item_path* |  |  |
| `MICROSOFT_TEAMS_DOWNLOAD_SHARED_DRIVE_ITEM` | Download shared drive item | read | `share_id`; *filename* |  |  |
| `MICROSOFT_TEAMS_GET_CALL_OPERATION` | Get call operation | read | `call_id`, `comms_operation_id`; *expand, select* |  |  |
| `MICROSOFT_TEAMS_GET_CHANNEL` | Get team channel | read | `channel_id`, `team_id`; *expand, select* |  | `Channel.ReadBasic.All`, `ChannelSettings.Read.All`, `ChannelSettings.ReadWrite.All`, `Directory.Read.All`, `Directory.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_CHANNEL_MESSAGE` | Get Teams channel message | read | `channel_id`, `message_id`, `team_id` |  | `ChannelMessage.Read.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_CHANNEL_MESSAGE_REPLY` | Get channel message reply | read | `channel_id`, `message_id`, `reply_id`, `team_id`; *expand, select* |  | `ChannelMessage.Read.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_CHAT` | Get chat | read | `chat_id`; *expand* |  | `Chat.ReadBasic`, `Chat.Read`, `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_GET_CHAT_MEMBER` | Get chat member | read | `chat_id`, `conversation_member_id` |  | `Chat.ReadBasic`, `ChatMember.Read`, `ChatMember.ReadWrite`, `Chat.Read`, `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_GET_CHAT_MESSAGE` | Get chat message | read | `chat_id`, `message_id` |  | `Chat.Read`, `Chat.ReadWrite`, `Chat.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_DAY_NOTE` | Get day note | read | `daynote_id`, `team_id`; *expand, select* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_DRIVE_ITEM` | Get drive item | read | `drive_id`; *if_none_match, include_children, include_deleted_items, item_id, item_path, select* |  |  |
| `MICROSOFT_TEAMS_GET_DRIVE_ITEM_BY_PATH` | Get drive item by path | read | `drive_id`, `item_path`; *expand, select* |  |  |
| `MICROSOFT_TEAMS_GET_FILES_FOLDER` | Get channel files folder | read | `channel_id`, `group_id`; *expand, select* |  | `Files.Read.All`, `Files.Read`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_GROUP_CALENDAR_EVENT` | Get group calendar event | read | `event_id`, `group_id`; *body_content_type, expand, select, timezone* |  |  |
| `MICROSOFT_TEAMS_GET_GROUP_DETAILS` | Get Microsoft 365 group details | read | `group_id`; *select* |  |  |
| `MICROSOFT_TEAMS_GET_GROUP_DRIVE` | Get group drive | read | `group_id`; *select* |  |  |
| `MICROSOFT_TEAMS_GET_MEETING_RECORDING_CONTENT` | Get meeting recording content | read | `online_meeting_id`, `recording_id`; *user_id* |  | `OnlineMeetingRecording.Read.All` |
| `MICROSOFT_TEAMS_GET_MEETING_TRANSCRIPT_CONTENT` | Get meeting transcript content | read | `online_meeting_id`, `transcript_id`; *format, user_id* |  | `OnlineMeetingTranscript.Read.All` |
| `MICROSOFT_TEAMS_GET_MY_PROFILE` | Get User Profile | read | *select, user_id* |  | `User.Read`, `User.ReadWrite`, `User.ReadBasic.All`, `User.Read.All`, `User.ReadWrite.All`, `Directory.Read.All`, `Directory.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_MY_PROFILE_PHOTO` | Get my profile photo | read |  |  | `ProfilePhoto.Read.All`, `ProfilePhoto.ReadWrite.All`, `User.Read`, `User.ReadWrite`, `User.ReadBasic.All`, `User.Read.All`, `User.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_OFFER_SHIFT_REQUEST` | Get offer shift request | read | `offer_shift_request_id`, `team_id`; *expand, select* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_ONLINE_MEETING` | Get online meeting | read | `meeting_id`; *expand, select, user_id* |  | `OnlineMeetings.Read`, `OnlineMeetings.ReadWrite` |
| `MICROSOFT_TEAMS_GET_OPEN_SHIFT` | Get open shift | read | `openshift_id`, `team_id`; *expand, select* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_OPEN_SHIFT_CHANGE_REQUEST` | Get open shift change request | read | `open_shift_change_request_id`, `team_id`; *expand, select* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_PLANNER_BUCKET` | Get Planner bucket | read | `bucket_id` |  |  |
| `MICROSOFT_TEAMS_GET_PLANNER_PLAN` | Get Planner plan | read | `plan_id`; *expand, select* |  |  |
| `MICROSOFT_TEAMS_GET_PLANNER_PLAN_DETAILS` | Get Planner plan details | read | `plan_id` |  |  |
| `MICROSOFT_TEAMS_GET_PLANNER_TASK` | Get Planner task | read | `task_id`; *expand, select* |  |  |
| `MICROSOFT_TEAMS_GET_PLANNER_TASK_DETAILS` | Get Planner task details | read | `task_id` |  |  |
| `MICROSOFT_TEAMS_GET_PRESENCE` | Get User Presence | read | `user_id`; *expand, select* |  | `Presence.Read`, `Presence.Read.All` |
| `MICROSOFT_TEAMS_GET_PRIMARY_CHANNEL` | Get primary channel | read | `team_id`; *expand, select* |  | `Channel.ReadBasic.All`, `ChannelSettings.Read.All`, `ChannelSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_SCHEDULE` | Get schedule | read | `team_id`; *expand, select* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_SCHEDULING_GROUP` | Get scheduling group | read | `schedulinggroup_id`, `team_id`; *expand, select* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_SHARED_DRIVE_ITEM` | Get shared drive item | read | `sharing_url`; *include_children, prefer* |  |  |
| `MICROSOFT_TEAMS_GET_SHIFT` | Get shift | read | `shift_id`, `team_id`; *expand, select* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_SITE_DRIVE_ITEM` | Get site drive item | read | `item_id`, `site_id`; *expand, select* |  |  |
| `MICROSOFT_TEAMS_GET_SWAP_SHIFTS_CHANGE_REQUEST` | Get swap shifts change request | read | `swap_shifts_change_request_id`, `team_id`; *expand, select* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_TAB` | Get channel tab | read | `channel_id`, `tab_id`, `team_id`; *expand, select* |  | `TeamsTab.Read.All`, `TeamsTab.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_TEAMS_APP_DEFINITION` | Get installed app in team | read | `app_installation_id`, `team_id`; *expand, select* |  | `TeamsAppInstallation.ReadWriteSelfForTeam`, `Directory.Read.All`, `Directory.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All`, `TeamsAppInstallation.ManageSelectedForTeam`, `TeamsAppInstallation.ReadForTeam`, `TeamsAppInstallation.ReadForUser`, `TeamsAppInstallation.ReadWriteAndConsentForTeam`, `TeamsAppInstallation.ReadWriteAndConsentSelfForTeam`, `TeamsAppInstallation.ReadWriteForTeam` |
| `MICROSOFT_TEAMS_GET_TEAM_BY_ID` | Get team by ID | read | `team_id`; *expand, select* |  |  |
| `MICROSOFT_TEAMS_GET_TEAM_FROM_GROUP` | Get Team | read | `group_id`; *expand, select* |  | `Directory.Read.All`, `Directory.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All`, `Team.ReadBasic.All`, `TeamSettings.Read.All`, `TeamSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_TEAM_MEMBER` | Get team member | read | `membership_id`, `team_id` |  | `TeamMember.Read.All`, `TeamMember.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_TEAM_OPERATION` | Get Team Operation Status | read | `operation_id`, `team_id` |  | `Team.ReadBasic.All`, `TeamSettings.Read.All`, `TeamSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_TIME_OFF` | Get time off by ID | read | `team_id`, `timeoff_id`; *expand, select* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_TIME_OFF_REASON` | Get time off reason | read | `team_id`, `time_off_reason_id` |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_TIME_OFF_REQUEST` | Get time off request | read | `team_id`, `time_off_request_id` |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_GET_USER_CHAT` | Get user chat | read | `chat_id`, `user_id`; *expand, select* |  | `Chat.ReadBasic`, `Chat.Read`, `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_GET_USER_TEAMWORK` | Get User Teamwork | read | `user_id`; *expand, select* |  | `UserTeamwork.Read` |
| `MICROSOFT_TEAMS_HIDE_CHAT_FOR_USER` | Hide chat for user | write | `chat_id`, `user` |  | `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_INSPECT_WORKBOOK` | Inspect workbook | read | `drive_id`, `item_id`; *values_only, workbook_session_id, worksheet, worksheet_select* |  |  |
| `MICROSOFT_TEAMS_LIST_APPLICATION_OWNERS` | List application owners | read | `application_id`; *count, cursor, expand, filter, orderby, search, select, top* |  |  |
| `MICROSOFT_TEAMS_LIST_ASSOCIATED_TEAMS` | List user's associated teams | read | `user_id`; *page_token* |  | `Team.ReadBasic.All`, `TeamSettings.Read.All`, `TeamSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_CHANNEL_MEMBERS` | List channel members | read | `channel_id`, `team_id`; *cursor, filter, select, top* |  |  |
| `MICROSOFT_TEAMS_LIST_CHANNEL_MESSAGE_CHANGES` | List channel message changes | read | `channel_id`, `team_id`; *cursor, delta_token, filter, top* |  |  |
| `MICROSOFT_TEAMS_LIST_CHANNEL_TABS` | List channel tabs | read | `channel_id`, `team_id`; *count, expand, filter, orderby, page_token, search, select, skip and 1 more* |  | `TeamsTab.Read.All`, `TeamsTab.ReadWrite.All`, `TeamsTab.ReadWriteForTeam`, `TeamsTab.ReadWriteSelfForTeam`, `Group.Read.All`, `Group.ReadWrite.All`, `Directory.Read.All`, `Directory.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_CHAT_MESSAGE_HOSTED_CONTENTS` | List chat message hosted contents | read | `chat_id`, `message_id` |  |  |
| `MICROSOFT_TEAMS_LIST_CHAT_TABS` | List chat tabs | read | `chat_id`; *cursor, expand, filter, select* |  |  |
| `MICROSOFT_TEAMS_LIST_COMMUNICATIONS_CALLS_OPERATIONS` | List Communications Call Operations | read | `call_id`; *count, expand, filter, orderby, page_token, search, select, skip and 1 more* |  |  |
| `MICROSOFT_TEAMS_LIST_DELETED_TEAMS` | List deleted teams | read | *count, expand, filter, orderby, page_token, search, select, skip and 1 more* |  | `Team.ReadBasic.All` |
| `MICROSOFT_TEAMS_LIST_DIRECTORY_ROLE_MEMBERS` | List directory role members | read | `role_id`; *select* |  |  |
| `MICROSOFT_TEAMS_LIST_DRIVE_FOLDER_CHILDREN` | List drive folder children | read | `drive_id`, `item_path`; *cursor, expand, orderby, select, top* |  |  |
| `MICROSOFT_TEAMS_LIST_DRIVE_ITEM_CHILDREN` | List Drive Item Children | read | `drive_id`, `item_id`; *expand, orderby, select, skip_token, top* |  |  |
| `MICROSOFT_TEAMS_LIST_DRIVE_ITEM_PERMISSIONS` | List drive item permissions | read | `drive_id`, `item_id`; *select* |  |  |
| `MICROSOFT_TEAMS_LIST_DRIVE_ROOT_CHANGES` | List Drive Root Changes | read | `drive_id`; *delta_exclude_parent, expand, page_token, select, token, top* |  |  |
| `MICROSOFT_TEAMS_LIST_GROUPS` | List directory groups | read | *count, cursor, expand, filter, orderby, search, select, top* |  |  |
| `MICROSOFT_TEAMS_LIST_GROUP_APP_ROLE_ASSIGNMENTS` | List group application role assignments | read | `group_id`; *cursor, select, top* |  |  |
| `MICROSOFT_TEAMS_LIST_GROUP_CALENDAR_EVENTS` | List group calendar events | read | `group_id`; *cursor, filter, orderby, prefer, select, top* |  |  |
| `MICROSOFT_TEAMS_LIST_GROUP_CALENDAR_VIEW` | List group calendar view | read | `end_date_time`, `group_id`, `start_date_time`; *cursor, filter, orderby, prefer, select, top* |  |  |
| `MICROSOFT_TEAMS_LIST_GROUP_DRIVE_ITEM_CHILDREN` | List group drive folder contents | read | `group_id`, `item_id`; *cursor, expand, orderby, select, top* |  |  |
| `MICROSOFT_TEAMS_LIST_GROUP_MEMBERS` | List group members | read | `group_id`; *count, cursor, expand, filter, search, select, top* |  |  |
| `MICROSOFT_TEAMS_LIST_GROUP_OWNERS` | List group owners | read | `group_id`; *count, cursor, expand, filter, search, select, top* |  |  |
| `MICROSOFT_TEAMS_LIST_GROUP_PLANNER_PLANS` | List group Planner plans | read | `group_id`; *cursor, select, top* |  |  |
| `MICROSOFT_TEAMS_LIST_GROUP_TEAM_OPERATIONS` | List Group Team Operations | read | `group_id`; *page_token* |  | `Team.ReadBasic.All`, `TeamSettings.Read.All`, `TeamSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_GROUP_THREAD_POSTS` | List group thread posts | read | `group_id`, `thread_id`; *cursor, expand, select* |  |  |
| `MICROSOFT_TEAMS_LIST_INCOMING_CHANNELS` | List team incoming channels | read | `team_id`; *filter, page_token, select* |  | `Channel.ReadBasic.All`, `ChannelSettings.Read.All`, `ChannelSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_INSTALLED_APPS` | List installed apps in team | read | `group_id`; *count, expand, filter, orderby, search, select, skip, top* |  | `TeamsAppInstallation.ReadForTeam`, `TeamsAppInstallation.ReadWriteSelfForTeam`, `TeamsAppInstallation.ManageSelectedForTeam`, `TeamsAppInstallation.ReadWriteAndConsentForTeam`, `TeamsAppInstallation.ReadWriteAndConsentSelfForTeam`, `TeamsAppInstallation.ReadWriteForTeam`, `Directory.Read.All`, `Directory.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_MEETING_ATTENDANCE_REPORTS` | List meeting attendance reports | read | `meeting_id`; *count, cursor, expand, filter, order_by, search, select, skip and 1 more* |  |  |
| `MICROSOFT_TEAMS_LIST_MEETING_RECORDINGS` | List meeting recordings | read | `online_meeting_id`; *filter, page_token, select, top, user_id* |  | `OnlineMeetingRecording.Read.All` |
| `MICROSOFT_TEAMS_LIST_MEETING_TRANSCRIPTS` | List meeting transcripts | read | `online_meeting_id`, `user_id`; *filter, page_token, select, top* |  | `OnlineMeetingTranscript.Read.All` |
| `MICROSOFT_TEAMS_LIST_MESSAGES_TO_SELF` | List messages to yourself | read | *cursor, filter, order_by, top* |  | `Chat.Read`, `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_LIST_MESSAGE_REPLIES` | List channel message replies | read | `channel_id`, `message_id`, `team_id`; *count, expand, filter, orderby, page_token, search, select, skip and 1 more* |  | `ChannelMessage.Read.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_MY_PLANNER_PLANS` | List my Planner plans | read | *cursor* |  |  |
| `MICROSOFT_TEAMS_LIST_MY_PLANNER_TASKS` | List my Planner tasks | read | *cursor* |  |  |
| `MICROSOFT_TEAMS_LIST_OFFER_SHIFT_REQUESTS` | List offer shift requests | read | `team_id`; *count, expand, filter, orderby, page_token, search, select, skip and 1 more* |  | `Group.ReadWrite.All`, `Schedule.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_ONLINE_MEETINGS` | List user online meetings | read | `filter`; *count, expand, orderby, page_token, search, select, skip, top and 1 more* |  | `OnlineMeetings.Read`, `OnlineMeetings.ReadWrite` |
| `MICROSOFT_TEAMS_LIST_OPEN_SHIFTS` | List open shifts | read | `team_id`; *count, expand, filter, orderby, page_token, search, select, skip and 1 more* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_OPEN_SHIFT_CHANGE_REQUESTS` | List open shift change requests | read | `team_id`; *count, expand, filter, orderby, page_token, search, select, skip and 1 more* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_ORGANIZATIONS` | List organizations | read | *select* |  |  |
| `MICROSOFT_TEAMS_LIST_PEOPLE` | List People | read | *filter, orderby, page_token, search, select, skip, top, user_id* |  | `People.Read`, `People.Read.All` |
| `MICROSOFT_TEAMS_LIST_PINNED_MESSAGES` | List pinned messages | read | `chat_id`; *expand, page_token* |  | `Chat.Read`, `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_LIST_PLANNER_BUCKET_TASKS` | List Planner bucket tasks | read | `bucket_id`; *cursor, select* |  |  |
| `MICROSOFT_TEAMS_LIST_PLANNER_PLAN_BUCKETS` | List Planner plan buckets | read | `plan_id`; *cursor, select* |  |  |
| `MICROSOFT_TEAMS_LIST_PLANNER_PLAN_TASKS` | List Planner plan tasks | read | `plan_id`; *cursor* |  |  |
| `MICROSOFT_TEAMS_LIST_PLANNER_TASK_MESSAGES` | List Planner Task Messages (Beta) | read | `task_id`; *page_token* |  |  |
| `MICROSOFT_TEAMS_LIST_SCHEDULE_DAY_NOTES` | List team schedule day notes | read | `team_id`; *page_token* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_SCHEDULING_GROUPS` | List scheduling groups | read | `team_id`; *count, expand, filter, orderby, search, select, skip, top* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_SERVICE_PRINCIPALS` | List service principals | read | *count, cursor, expand, filter, orderby, search, select, top* |  |  |
| `MICROSOFT_TEAMS_LIST_SHIFTS` | List shifts | read | `team_id`; *count, expand, filter, orderby, page_token, search, select, skip and 1 more* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_SUBSCRIPTIONS` | List Graph subscriptions | read | *cursor* |  |  |
| `MICROSOFT_TEAMS_LIST_SWAP_SHIFTS_CHANGE_REQUESTS` | List swap shifts change requests | read | `team_id`; *count, expand, filter, orderby, page_token, search, select, skip and 1 more* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_TEAM_MEMBERS` | List team members | read | `team_id`; *filter, limit, page_token, select, top* |  | `TeamMember.Read.All`, `TeamMember.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_TEAM_OPERATIONS` | List Team Operations | read | `team_id`; *page_token* |  | `Team.ReadBasic.All`, `TeamSettings.Read.All`, `TeamSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_TEAM_PERMISSION_GRANTS` | List teams permission grants | read | `team_id`; *page_token* |  | `TeamsAppInstallation.ReadForTeam`, `TeamsAppInstallation.ReadWriteSelfForTeam`, `TeamsAppInstallation.ReadWriteForTeam` |
| `MICROSOFT_TEAMS_LIST_TEAM_TAGS` | List team tags | read | `team_id`; *cursor, filter, top* |  |  |
| `MICROSOFT_TEAMS_LIST_TEAM_TAG_MEMBERS` | List team tag members | read | `tag_id`, `team_id`; *cursor, filter, top* |  |  |
| `MICROSOFT_TEAMS_LIST_TIME_OFF` | List team time off entries | read | `team_id`; *count, expand, filter, orderby, page_token, search, select, skip and 1 more* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_TIME_OFF_REASONS` | List team schedule time off reasons | read | `team_id`; *count, expand, filter, orderby, search, select, skip, top* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_TIME_OFF_REQUESTS` | List time off requests | read | `team_id`; *count, expand, filter, orderby, page_token, search, select, skip and 1 more* |  | `Schedule.Read.All`, `Schedule.ReadWrite.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_USERS` | List users | read | *consistencyLevel, count, expand, filter, orderby, page_token, search, select and 1 more* |  | `Directory.Read.All`, `Directory.ReadWrite.All`, `User.Read.All`, `User.ReadBasic.All`, `User.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_USER_CHAT_MEMBERS` | List user's chat members | read | `chat_id`, `user_id` |  | `Chat.ReadBasic`, `ChatMember.ReadWrite`, `Chat.Read`, `Chat.ReadWrite`, `ChatMember.Read`, `Chat.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_USER_CHAT_MESSAGES` | List user chat messages | read | `chat_id`, `user_id`; *filter, limit, orderby, page_token, top* |  | `Chat.Read`, `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_LIST_USER_DIRECT_REPORTS` | List user direct reports | read | `user_id`; *count, cursor, expand, filter, select* |  |  |
| `MICROSOFT_TEAMS_LIST_USER_JOINED_TEAMS` | List user joined teams | read | `user_id`; *page_token* |  | `Team.ReadBasic.All`, `TeamSettings.Read.All`, `TeamSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_LIST_USER_LICENSE_DETAILS` | List user license details | read | `user_id`; *cursor, select* |  |  |
| `MICROSOFT_TEAMS_LIST_USER_MEMBERSHIPS` | List user directory memberships | read | `user_id`; *count, cursor, filter, search, select, top* |  |  |
| `MICROSOFT_TEAMS_LIST_WORKBOOK_TABLES` | List workbook tables | read | `drive_id`, `item_id`; *select, session_id, skip, top* |  |  |
| `MICROSOFT_TEAMS_LIST_WORKBOOK_WORKSHEETS` | List workbook worksheets | read | `drive_id`, `item_id`; *cursor, select, workbook_session_id* |  |  |
| `MICROSOFT_TEAMS_MARK_CHAT_READ_FOR_USER` | Mark chat as read for user | write | `chat_id`, `user` |  | `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_MARK_CHAT_UNREAD_FOR_USER` | Mark chat as unread for user | write | `chat_id`, `user`; *lastMessageReadDateTime* |  | `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_PIN_MESSAGE` | Pin message in chat | write | `chat_id`, `message_id` |  | `ChatMessage.Send`, `Chat.ReadWrite`, `Chat.ReadWrite.All` |
| `MICROSOFT_TEAMS_POST_MESSAGE_REPLY` | Reply to Teams channel message | write | `channel_id`, `content`, `message_id`, `team_id`; *attachments, content_type, hosted_contents, importance, locale, mentions, subject, summary* |  | `ChannelMessage.Send`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_PROVISION_CHANNEL_EMAIL` | Provision channel email address | write | `channel_id`, `team_id` |  | `ChannelSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_READ_WORKSHEET_RANGE` | Read worksheet range | read | `drive_id`, `item_id`, `worksheet`; *address, select, used_range, values_only, workbook_session_id* |  |  |
| `MICROSOFT_TEAMS_REMOVE_CHANNEL_EMAIL` | Remove channel email address | write | `channel_id`, `team_id` |  | `ChannelSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_REMOVE_CHANNEL_MEMBER` | Remove channel member | write | `channel_id`, `membership_id`, `team_id` |  |  |
| `MICROSOFT_TEAMS_REMOVE_CHAT_MEMBER` | Remove chat member | write | `chat_id`, `membership_id`; *if_match* |  | `ChatMember.ReadWrite` |
| `MICROSOFT_TEAMS_REMOVE_GROUP_MEMBER` | Remove Microsoft 365 group member | write | `group_id`, `member_id` |  |  |
| `MICROSOFT_TEAMS_REMOVE_TEAM_MEMBER` | Remove team member | write | `membership_id`, `team_id`; *if_match* |  | `TeamMember.ReadWriteNonOwnerRole.All`, `TeamMember.ReadWrite.All` |
| `MICROSOFT_TEAMS_REMOVE_TEAM_MEMBERS` | Remove team members | write | `members`, `team_id` |  | `TeamMember.ReadWrite.All` |
| `MICROSOFT_TEAMS_SEARCH_FILES` | Search Teams files | read | `query`; *from_index, region, size* |  | `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `MICROSOFT_TEAMS_SEARCH_MESSAGES` | Search Teams messages | read | *enable_top_results, from_index, query, region, size* |  | `Chat.Read`, `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_SEND_ACTIVITY_NOTIFICATION` | Send Activity Notification to Recipients | write | `activityType`, `previewText`, `recipients`, `topic`; *chainId, teamsAppId, templateParameters* |  | `TeamsActivity.Send` |
| `MICROSOFT_TEAMS_SEND_MESSAGE_TO_SELF` | Send a message to yourself | write | `content`; *content_type* |  | `ChatMessage.Send`, `Chat.ReadWrite`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_SET_CHANNEL_MESSAGE_REACTION` | Set channel message reaction | write | `channel_id`, `message_id`, `reaction_type`, `team_id` |  |  |
| `MICROSOFT_TEAMS_SET_CHAT_MESSAGE_REACTION` | Set chat message reaction | write | `chat_id`, `message_id`, `reaction_type` |  |  |
| `MICROSOFT_TEAMS_SET_PRESENCE` | Set Presence | write | `activity`, `availability`, `session_id`, `user_id`; *expiration_duration* |  | `Presence.ReadWrite` |
| `MICROSOFT_TEAMS_SET_PRESENCE_AUTOMATIC_LOCATION` | Set Presence Automatic Location | write | `user_id`, `work_location_type`; *place_id* |  | `Presence.ReadWrite` |
| `MICROSOFT_TEAMS_SET_USER_PREFERRED_PRESENCE` | Set User Preferred Presence | write | `activity`, `availability`, `user_id`; *expiration_duration* |  | `Presence.ReadWrite` |
| `MICROSOFT_TEAMS_SET_USER_PRESENCE_MANUAL_LOCATION` | Set User Presence Manual Location | write | `user_id`, `workLocationType`; *placeId* |  | `Presence.ReadWrite` |
| `MICROSOFT_TEAMS_SHARE_TEAM_SCHEDULE` | Share Team Schedule | write | `endDateTime`, `notifyTeam`, `startDateTime`, `team_id` |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_TEAMS_CREATE_CHAT` | Create Chat | write | `chatType`, `members`; *installedApps, topic* |  | `Chat.Create`, `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_TEAMS_LIST` | List Teams | read | *count, filter, page_token, select, top* |  | `Team.ReadBasic.All`, `TeamSettings.Read.All`, `TeamSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_TEAMS_LIST_CHANNELS` | List team channels | read | `team_id`; *filter, include_shared_channels, page_token, select* |  | `Channel.ReadBasic.All` |
| `MICROSOFT_TEAMS_TEAMS_LIST_CHANNEL_MESSAGES` | List channel messages | read | `channel_id`, `team_id`; *expand, page_token, top* |  | `ChannelMessage.Read.All`, `Group.Read.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_TEAMS_POST_CHANNEL_MESSAGE` | Send Channel Message | write | `channel_id`, `content`, `team_id`; *attachments, content_type, hosted_contents, importance, locale, mentions, subject, summary* |  | `ChannelMessage.Send`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_TEAMS_POST_CHAT_MESSAGE` | Send message to Teams chat | write | `chat_id`, `content`; *attachments, content_type, hostedContents, importance, locale, mentions, subject* |  | `ChatMessage.Send`, `Chat.ReadWrite`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_UNARCHIVE_CHANNEL` | Unarchive channel | write | `channel_id`, `team_id` |  | `ChannelSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_UNARCHIVE_GROUP_TEAM_CHANNEL` | Unarchive group team channel | write | `channel_id`, `group_id` |  | `ChannelSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_UNARCHIVE_TEAM` | Unarchive Teams team | write | `team_id` |  | `Directory.ReadWrite.All`, `Group.ReadWrite.All`, `TeamSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_UNHIDE_CHAT_FOR_USER` | Unhide Chat For User | write | `chat_id`, `user` |  | `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_UNPIN_MESSAGE` | Unpin chat message | write | `chat_id`, `pinned_message_id`; *if_match* |  | `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_UPDATE_CALL` | Update Call | write | `call_id`; *callChainId, callbackUri, myParticipantId, odata_type, subject, tenantId* |  |  |
| `MICROSOFT_TEAMS_UPDATE_CALL_OPERATION` | Update Call Operation | write | `call_id`, `comms_operation_id`; *client_context, result_info, status* |  |  |
| `MICROSOFT_TEAMS_UPDATE_CALL_PARTICIPANT` | Update Call Participant | write | `callId`, `participantId`; *isInLobby* |  |  |
| `MICROSOFT_TEAMS_UPDATE_CHANNEL` | Update Teams channel | write | `channel_id`, `group_id`; *description, display_name, is_favorite_by_default* |  | `ChannelSettings.ReadWrite.All`, `Directory.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_UPDATE_CHANNEL_MESSAGE` | Update Teams channel message | write | `channel_id`, `message_id`, `team_id`; *content, content_type, importance, model, policy_violation, subject* |  | `ChannelMessage.ReadWrite`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_UPDATE_CHAT` | Update Chat | write | `chat_id`; *topic* |  | `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_UPDATE_CHAT_MESSAGE` | Update Teams chat message | write | `chat_id`, `content`, `message_id`; *content_type, model* |  | `Chat.ReadWrite` |
| `MICROSOFT_TEAMS_UPDATE_CONTENT_SHARING_SESSION` | Update Content Sharing Session | write | `call_id`, `content_sharing_session_id`; *id, odata_type* |  |  |
| `MICROSOFT_TEAMS_UPDATE_DAY_NOTE` | Update schedule day note | write | `dayNote_id`, `team_id`; *dayNoteDate, draftDayNote, sharedDayNote* |  | `Schedule.ReadWrite.All` |
| `MICROSOFT_TEAMS_UPDATE_DRIVE_ITEM` | Update Drive Item | write | `drive_id`, `item_id`; *created_date_time, description, if_match, last_modified_date_time, name, parent_reference_id* |  |  |
| `MICROSOFT_TEAMS_UPDATE_ONLINE_MEETING` | Update online meeting | write | `online_meeting_id`, `user_id`; *end_date_time, is_entry_exit_announced, lobby_bypass_settings, participants, record_automatically, start_date_time, subject, watermark_protection* |  | `OnlineMeetings.ReadWrite` |
| `MICROSOFT_TEAMS_UPDATE_OPEN_SHIFT` | Update open shift | write | `openShiftId`, `teamId`; *draftOpenShift, isStagedForDeletion, schedulingGroupId, sharedOpenShift* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_UPDATE_PLANNER_BUCKET` | Update Planner bucket | write | `bucket_id`, `etag`; *name, order_hint, return_representation* |  |  |
| `MICROSOFT_TEAMS_UPDATE_PLANNER_TASK` | Update Planner task | write | `etag`, `task_id`; *applied_categories, assignee_priority, assignments, bucket_id, conversation_thread_id, due_date_time, order_hint, percent_complete and 3 more* |  |  |
| `MICROSOFT_TEAMS_UPDATE_PLANNER_TASK_DETAILS` | Update Planner task details | write | `etag`, `task_id`; *checklist, description, preview_type, references, return_representation* |  |  |
| `MICROSOFT_TEAMS_UPDATE_SCHEDULING_GROUP` | Update scheduling group | write | `scheduling_group_id`, `team_id`; *code, display_name, is_active, user_ids* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_UPDATE_SUBSCRIPTION` | Renew Graph subscription | write | `subscription_id`; *expiration_date_time, notification_url* |  |  |
| `MICROSOFT_TEAMS_UPDATE_TAB` | Update Teams channel tab | write | `channel_id`, `tab_id`, `team_id`; *configuration, display_name, web_url* |  | `TeamsTab.ReadWrite.All`, `TeamsTab.ReadWriteForTeam`, `Group.ReadWrite.All`, `Directory.ReadWrite.All` |
| `MICROSOFT_TEAMS_UPDATE_TEAM` | Update Team | write | `team_id`; *classSettings, classification, description, displayName, funSettings, guestSettings, memberSettings, messagingSettings and 2 more* |  | `Directory.ReadWrite.All`, `Group.ReadWrite.All`, `TeamSettings.ReadWrite.All` |
| `MICROSOFT_TEAMS_UPDATE_TEAM_MEMBER` | Update team member | write | `conversation_member_id`, `team_id`; *display_name, odata_type, roles, visible_history_start_date_time* |  | `TeamMember.ReadWriteNonOwnerRole.All`, `TeamMember.ReadWrite.All` |
| `MICROSOFT_TEAMS_UPDATE_TEAM_SCHEDULE_SHIFT` | Update team schedule shift | write | `shiftId`, `teamId`; *draftShift, isStagedForDeletion, schedulingGroupId, sharedShift, userId* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_UPDATE_TIME_OFF` | Update time off | write | `team_id`, `timeoff_id`, `user_id`; *draft_time_off, is_staged_for_deletion, shared_time_off* |  | `Group.Read.All`, `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_UPDATE_TIME_OFF_REASON` | Update time off reason | write | `team_id`, `time_off_reason_id`; *code, display_name, icon_type, is_active* |  | `Schedule.ReadWrite.All`, `Group.ReadWrite.All` |
| `MICROSOFT_TEAMS_UPLOAD_DRIVE_FILE` | Upload drive file | write | `content`, `drive_id`, `filename`, `parent_item_id` |  |  |
| `MICROSOFT_TEAMS_WHO_AM_I` | Who Am I | read |  |  |  |
