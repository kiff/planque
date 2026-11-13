# Discord Bot: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

157 tools: 72 read, 85 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `DISCORDBOT_ADD_GUILD_MEMBER` | Add guild member | write | `access_token`, `guild_id`, `user_id`; *deaf, mute, nick, roles* |  |  |
| `DISCORDBOT_ADD_GUILD_MEMBER_ROLE` | Assign role to guild member | write | `guild_id`, `role_id`, `user_id` |  |  |
| `DISCORDBOT_ADD_MY_MESSAGE_REACTION` | Add Reaction To Message | write | `channel_id`, `emoji`, `message_id` |  |  |
| `DISCORDBOT_ADD_THREAD_MEMBER` | Add thread member | write | `channel_id`, `user_id` |  |  |
| `DISCORDBOT_BAN_USER_FROM_GUILD` | Ban user from guild | write | `guild_id`, `user_id`; *delete_message_seconds* |  |  |
| `DISCORDBOT_BULK_BAN_USERS_FROM_GUILD` | Bulk ban users from guild | write | `guild_id`, `user_ids`; *delete_message_seconds* |  |  |
| `DISCORDBOT_BULK_DELETE_MESSAGES` | Bulk Delete Messages | write | `channel_id`, `messages` |  |  |
| `DISCORDBOT_CREATE_AUTO_MODERATION_RULE` | Create auto moderation rule | write | `actions`, `event_type`, `guild_id`, `name`, `trigger_type`; *enabled, exempt_channels, exempt_roles, trigger_metadata* |  |  |
| `DISCORDBOT_CREATE_CHANNEL_INVITE` | Create channel invite | write | `channel_id`; *max_age, max_uses, target_application_id, target_type, target_user_id, temporary, unique* |  |  |
| `DISCORDBOT_CREATE_DM` | Initiate user channel with recipient | write | *access_tokens, nicks, recipient_id* |  |  |
| `DISCORDBOT_CREATE_GUILD_APPLICATION_COMMAND` | Create guild application command | write | `application_id`, `guild_id`, `name`; *default_member_permissions, description, description_localizations, dm_permission, name_localizations, options, type* |  |  |
| `DISCORDBOT_CREATE_GUILD_CHANNEL` | Create guild channel | write | `guild_id`, `name`; *available_tags, bitrate, default_auto_archive_duration, default_forum_layout, default_reaction_emoji, default_sort_order, default_thread_rate_limit_per_user, nsfw and 9 more* |  |  |
| `DISCORDBOT_CREATE_GUILD_EMOJI` | Create guild emoji | write | `guild_id`, `image`, `name`; *roles* |  |  |
| `DISCORDBOT_CREATE_GUILD_ROLE` | Create role with guild id | write | `guild_id`; *color, hoist, icon, mentionable, name, permissions, unicode_emoji* |  |  |
| `DISCORDBOT_CREATE_GUILD_SCHEDULED_EVENT` | Create guild scheduled event | write | `entity_type`, `guild_id`, `name`, `scheduled_start_time`; *channel_id, description, entity_metadata, image, privacy_level, scheduled_end_time* |  |  |
| `DISCORDBOT_CREATE_GUILD_STICKER` | Create guild sticker | write | `file`, `guild_id`, `name`, `tags`; *description* |  |  |
| `DISCORDBOT_CREATE_GUILD_TEMPLATE` | Create guild template | write | `guild_id`, `name`; *description* |  |  |
| `DISCORDBOT_CREATE_INTERACTION_RESPONSE` | Post interaction callback | write | `interaction_id`, `interaction_token`, `type`; *data* |  |  |
| `DISCORDBOT_CREATE_MESSAGE` | Send Message To Channel | write | `channel_id`; *allowed_mentions, components, content, embeds, flags, message_reference, sticker_ids, tts* |  |  |
| `DISCORDBOT_CREATE_STAGE_INSTANCE` | Create stage instance | write | `channel_id`, `topic`; *privacy_level, send_start_notification* |  |  |
| `DISCORDBOT_CREATE_THREAD` | Create new thread in channel | write | `channel_id`, `name`; *applied_tags, auto_archive_duration, invitable, message, rate_limit_per_user, type* |  |  |
| `DISCORDBOT_CREATE_THREAD_FROM_MESSAGE` | Create thread from message | write | `channel_id`, `message_id`, `name`; *auto_archive_duration, rate_limit_per_user* |  |  |
| `DISCORDBOT_CREATE_WEBHOOK` | Create channel webhook | write | `channel_id`, `name`; *avatar* |  |  |
| `DISCORDBOT_CROSSPOST_MESSAGE` | Crosspost Message | write | `channel_id`, `message_id` |  |  |
| `DISCORDBOT_DELETE_ALL_MESSAGE_REACTIONS` | Delete All Reactions | write | `channel_id`, `message_id` |  |  |
| `DISCORDBOT_DELETE_ALL_MESSAGE_REACTIONS_BY_EMOJI` | Delete All Reactions By Emoji | write | `channel_id`, `emoji`, `message_id` |  |  |
| `DISCORDBOT_DELETE_AUTO_MODERATION_RULE` | Delete auto moderation rule | write | `guild_id`, `rule_id` |  |  |
| `DISCORDBOT_DELETE_CHANNEL` | Delete channel by id | write | `channel_id` |  |  |
| `DISCORDBOT_DELETE_CHANNEL_PERMISSION_OVERWRITE` | Delete channel permission override | write | `channel_id`, `overwrite_id` |  |  |
| `DISCORDBOT_DELETE_GUILD` | Delete guild by id | write | `guild_id` |  |  |
| `DISCORDBOT_DELETE_GUILD_APPLICATION_COMMAND` | Delete guild application command | write | `application_id`, `command_id`, `guild_id` |  |  |
| `DISCORDBOT_DELETE_GUILD_EMOJI` | Delete guild emoji by id | write | `emoji_id`, `guild_id` |  |  |
| `DISCORDBOT_DELETE_GUILD_INTEGRATION` | Delete guild integration | write | `guild_id`, `integration_id` |  |  |
| `DISCORDBOT_DELETE_GUILD_MEMBER` | Delete guild member by id | write | `guild_id`, `user_id` |  |  |
| `DISCORDBOT_DELETE_GUILD_MEMBER_ROLE` | Delete guild member role | write | `guild_id`, `role_id`, `user_id` |  |  |
| `DISCORDBOT_DELETE_GUILD_ROLE` | Delete role from guild | write | `guild_id`, `role_id` |  |  |
| `DISCORDBOT_DELETE_GUILD_SCHEDULED_EVENT` | Delete guild scheduled event | write | `event_id`, `guild_id` |  |  |
| `DISCORDBOT_DELETE_GUILD_STICKER` | Delete guild sticker | write | `guild_id`, `sticker_id` |  |  |
| `DISCORDBOT_DELETE_GUILD_TEMPLATE` | Delete guild template by code | write | `guild_id`, `template_code` |  |  |
| `DISCORDBOT_DELETE_MESSAGE` | Delete Message | write | `channel_id`, `message_id` |  |  |
| `DISCORDBOT_DELETE_MY_MESSAGE_REACTION` | Remove My Reaction | write | `channel_id`, `emoji`, `message_id` |  |  |
| `DISCORDBOT_DELETE_ORIGINAL_WEBHOOK_MESSAGE` | Delete original webhook message | write | `webhook_id`, `webhook_token`; *thread_id* |  |  |
| `DISCORDBOT_DELETE_STAGE_INSTANCE` | Delete stage instance | write | `channel_id` |  |  |
| `DISCORDBOT_DELETE_THREAD_MEMBER` | Remove thread member | write | `channel_id`, `user_id` |  |  |
| `DISCORDBOT_DELETE_USER_MESSAGE_REACTION` | Remove User Reaction | write | `channel_id`, `emoji`, `message_id`, `user_id` |  |  |
| `DISCORDBOT_DELETE_WEBHOOK` | Delete webhook by id | write | `webhook_id` |  |  |
| `DISCORDBOT_DELETE_WEBHOOK_BY_TOKEN` | Delete webhook by token | write | `webhook_id`, `webhook_token` |  |  |
| `DISCORDBOT_DELETE_WEBHOOK_MESSAGE` | Delete webhook message | write | `message_id`, `webhook_id`, `webhook_token`; *thread_id* |  |  |
| `DISCORDBOT_EXECUTE_GITHUB_COMPATIBLE_WEBHOOK` | Execute GitHub-compatible webhook | write | `webhook_id`, `webhook_token`; *action, answer, check_run, check_suite, comment, commits, compare, discussion and 14 more* |  |  |
| `DISCORDBOT_EXECUTE_SLACK_COMPATIBLE_WEBHOOK` | Execute Slack-compatible webhook | write | `webhook_id`, `webhook_token`; *attachments, icon_url, text, thread_id, username, wait* |  |  |
| `DISCORDBOT_EXECUTE_WEBHOOK` | Execute webhook | write | `webhook_id`, `webhook_token`; *allowed_mentions, applied_tags, avatar_url, components, content, embeds, flags, poll and 5 more* |  |  |
| `DISCORDBOT_FOLLOW_CHANNEL` | Add follower to channel via webhook | write | `channel_id`, `webhook_channel_id` |  |  |
| `DISCORDBOT_GET_ACTIVE_GUILD_THREADS` | List active threads in guild | read | `guild_id` |  |  |
| `DISCORDBOT_GET_APPLICATION` | Get application details | read | `application_id` |  |  |
| `DISCORDBOT_GET_APPLICATION_COMMAND` | Get global application command | read | `application_id`, `command_id` |  |  |
| `DISCORDBOT_GET_APPLICATION_ROLE_CONNECTIONS_METADATA` | Get application role connections metadata | read | `application_id` |  |  |
| `DISCORDBOT_GET_APPLICATION_USER_ROLE_CONNECTION` | Get user role connection | read | `application_id` |  |  |
| `DISCORDBOT_GET_AUTO_MODERATION_RULE` | Get auto moderation rule | read | `guild_id`, `rule_id` |  |  |
| `DISCORDBOT_GET_BOT_GATEWAY` | Get bot gateway | read |  |  |  |
| `DISCORDBOT_GET_CHANNEL` | Retrieve channel details | read | `channel_id` |  |  |
| `DISCORDBOT_GET_CONNECTED_GUILD` | Get connected guild | read |  |  |  |
| `DISCORDBOT_GET_GATEWAY` | Get gateway URL | read |  |  |  |
| `DISCORDBOT_GET_GUILD` | Retrieve guild details | read | `guild_id`; *with_counts* |  |  |
| `DISCORDBOT_GET_GUILDS_ONBOARDING` | Get guild onboarding | read | `guild_id` |  |  |
| `DISCORDBOT_GET_GUILD_APPLICATION_COMMAND` | Get guild application command | read | `application_id`, `command_id`, `guild_id` |  |  |
| `DISCORDBOT_GET_GUILD_APPLICATION_COMMAND_PERMISSIONS` | Get guild command permissions | read | `application_id`, `command_id`, `guild_id` |  |  |
| `DISCORDBOT_GET_GUILD_BAN` | Get guild ban | read | `guild_id`, `user_id` |  |  |
| `DISCORDBOT_GET_GUILD_EMOJI` | Fetch emoji by guild and id | read | `emoji_id`, `guild_id` |  |  |
| `DISCORDBOT_GET_GUILD_MEMBER` | Retrieve guild member by user id | read | `guild_id`, `user_id` |  |  |
| `DISCORDBOT_GET_GUILD_PREVIEW` | Preview guild by id | read | `guild_id` |  |  |
| `DISCORDBOT_GET_GUILD_SCHEDULED_EVENT` | Get guild scheduled event | read | `event_id`, `guild_id`; *with_user_count* |  |  |
| `DISCORDBOT_GET_GUILD_STICKER` | Retrieve sticker from guild | read | `guild_id`, `sticker_id` |  |  |
| `DISCORDBOT_GET_GUILD_TEMPLATE` | Retrieve guild template with code | read | `template_code` |  |  |
| `DISCORDBOT_GET_GUILD_VANITY_URL` | Retrieve guild vanity url | read | `guild_id` |  |  |
| `DISCORDBOT_GET_GUILD_WEBHOOKS` | Retrieve guild webhooks | read | `guild_id` |  |  |
| `DISCORDBOT_GET_GUILD_WELCOME_SCREEN` | Retrieve guild welcome screen | read | `guild_id` |  |  |
| `DISCORDBOT_GET_GUILD_WIDGET` | Retrieve guild widget json | read | `guild_id` |  |  |
| `DISCORDBOT_GET_GUILD_WIDGET_PNG` | Get Guild Widget PNG | read | `guild_id`; *style* |  |  |
| `DISCORDBOT_GET_GUILD_WIDGET_SETTINGS` | Retrieve guild widget settings | read | `guild_id` |  |  |
| `DISCORDBOT_GET_MESSAGE` | Get Message | read | `channel_id`, `message_id` |  |  |
| `DISCORDBOT_GET_MY_APPLICATION` | Get my application | read |  |  |  |
| `DISCORDBOT_GET_MY_OAUTH2_APPLICATION` | Get my OAuth2 application | read |  |  |  |
| `DISCORDBOT_GET_ORIGINAL_WEBHOOK_MESSAGE` | Retrieve original webhook message | read | `webhook_id`, `webhook_token`; *thread_id* |  |  |
| `DISCORDBOT_GET_PUBLIC_KEYS` | Get OAuth2 public keys | read |  |  |  |
| `DISCORDBOT_GET_STAGE_INSTANCE` | Get stage instance | read | `channel_id` |  |  |
| `DISCORDBOT_GET_STICKER` | Get sticker | read | `sticker_id` |  |  |
| `DISCORDBOT_GET_THREAD_MEMBER` | Retrieve thread member by id | read | `channel_id`, `user_id`; *with_member* |  |  |
| `DISCORDBOT_GET_USER` | Get user | read | `user_id` |  |  |
| `DISCORDBOT_GET_WEBHOOK` | Retrieve webhook by id | read | `webhook_id` |  |  |
| `DISCORDBOT_GET_WEBHOOK_BY_TOKEN` | Retrieve webhook by token | read | `webhook_id`, `webhook_token` |  |  |
| `DISCORDBOT_GET_WEBHOOK_MESSAGE` | Retrieve webhook message | read | `message_id`, `webhook_id`, `webhook_token`; *thread_id* |  |  |
| `DISCORDBOT_INVITE_RESOLVE` | Resolve invite by code | read | `invite_code`; *with_counts, with_expiration* |  |  |
| `DISCORDBOT_INVITE_REVOKE` | Revoke invite by code | write | `invite_code` |  |  |
| `DISCORDBOT_JOIN_THREAD` | Join thread | write | `channel_id` |  |  |
| `DISCORDBOT_LEAVE_GUILD` | Leave Guild | write | `guild_id` |  |  |
| `DISCORDBOT_LEAVE_THREAD` | Leave thread | write | `channel_id` |  |  |
| `DISCORDBOT_LIST_APPLICATION_COMMANDS` | List global application commands | read | `application_id`; *with_localizations* |  |  |
| `DISCORDBOT_LIST_AUTO_MODERATION_RULES` | List auto moderation rules | read | `guild_id` |  |  |
| `DISCORDBOT_LIST_CHANNEL_INVITES` | List channel invites | read | `channel_id` |  |  |
| `DISCORDBOT_LIST_CHANNEL_WEBHOOKS` | List Channel Webhooks | read | `channel_id` |  |  |
| `DISCORDBOT_LIST_GUILD_APPLICATION_COMMANDS` | List guild application commands | read | `application_id`, `guild_id`; *with_localizations* |  |  |
| `DISCORDBOT_LIST_GUILD_APPLICATION_COMMAND_PERMISSIONS` | List guild command permissions | read | `application_id`, `guild_id` |  |  |
| `DISCORDBOT_LIST_GUILD_AUDIT_LOG_ENTRIES` | Get guild audit logs | read | `guild_id`; *action_type, after, before, limit, user_id* |  |  |
| `DISCORDBOT_LIST_GUILD_BANS` | List guild bans | read | `guild_id`; *after, before, limit* |  |  |
| `DISCORDBOT_LIST_GUILD_CHANNELS` | Retrieve guild channels | read | `guild_id` |  |  |
| `DISCORDBOT_LIST_GUILD_EMOJIS` | Retrieve guild emojis | read | `guild_id` |  |  |
| `DISCORDBOT_LIST_GUILD_INTEGRATIONS` | List guild integrations | read | `guild_id` |  |  |
| `DISCORDBOT_LIST_GUILD_INVITES` | List guild invites | read | `guild_id` |  |  |
| `DISCORDBOT_LIST_GUILD_MEMBERS` | Get guild members | read | `guild_id`; *after, limit* |  |  |
| `DISCORDBOT_LIST_GUILD_ROLES` | List guild roles | read | `guild_id` |  |  |
| `DISCORDBOT_LIST_GUILD_SCHEDULED_EVENTS` | List guild scheduled events | read | `guild_id`; *with_user_count* |  |  |
| `DISCORDBOT_LIST_GUILD_SCHEDULED_EVENT_USERS` | List scheduled event users | read | `event_id`, `guild_id`; *after, before, limit, with_member* |  |  |
| `DISCORDBOT_LIST_GUILD_STICKERS` | Retrieve guild stickers | read | `guild_id` |  |  |
| `DISCORDBOT_LIST_GUILD_TEMPLATES` | Get guild templates by guild id | read | `guild_id` |  |  |
| `DISCORDBOT_LIST_GUILD_VOICE_REGIONS` | List guild voice regions | read | `guild_id` |  |  |
| `DISCORDBOT_LIST_MESSAGES` | Fetch Messages From Channel | read | `channel_id`; *after, around, before, limit* |  |  |
| `DISCORDBOT_LIST_MESSAGE_REACTIONS_BY_EMOJI` | List Reactions By Emoji | read | `channel_id`, `emoji`, `message_id`; *after, limit* |  |  |
| `DISCORDBOT_LIST_MY_PRIVATE_ARCHIVED_THREADS` | List my private archived threads | read | `channel_id`; *before, limit* |  |  |
| `DISCORDBOT_LIST_PINNED_MESSAGES` | List Pinned Messages | read | `channel_id` |  |  |
| `DISCORDBOT_LIST_PRIVATE_ARCHIVED_THREADS` | List private archived threads | read | `channel_id`; *before, limit* |  |  |
| `DISCORDBOT_LIST_PUBLIC_ARCHIVED_THREADS` | List public archived threads | read | `channel_id`; *before, limit* |  |  |
| `DISCORDBOT_LIST_STICKER_PACKS` | List sticker packs | read |  |  |  |
| `DISCORDBOT_LIST_THREAD_MEMBERS` | List thread members | read | `channel_id`; *after, limit, with_member* |  |  |
| `DISCORDBOT_LIST_VOICE_REGIONS` | List voice regions | read |  |  |  |
| `DISCORDBOT_PIN_MESSAGE` | Pin Message | write | `channel_id`, `message_id` |  |  |
| `DISCORDBOT_PREVIEW_PRUNE_GUILD` | Preview guild prune | read | `guild_id`; *days, include_roles* |  |  |
| `DISCORDBOT_PRUNE_GUILD` | Prune inactive guild members | write | `guild_id`; *compute_prune_count, days, include_roles* |  |  |
| `DISCORDBOT_PUT_GUILDS_ONBOARDING` | Update guild onboarding configuration | write | `guild_id`; *default_channel_ids, enabled, mode, prompts* |  |  |
| `DISCORDBOT_REORDER_GUILD_ROLES` | Reorder guild roles | write | `guild_id`, `role_positions`; *audit_log_reason* |  |  |
| `DISCORDBOT_SEARCH_GUILD_MEMBERS` | Search guild members by username or nickname | read | `guild_id`, `query`; *limit* |  |  |
| `DISCORDBOT_SEARCH_GUILD_MESSAGES` | Search guild messages | read | `guild_id`; *author_id, channel_id, content, limit, max_id, min_id, offset, sort_by and 1 more* |  |  |
| `DISCORDBOT_SET_CHANNEL_PERMISSION_OVERWRITE` | Modify channel permissions | write | `channel_id`, `overwrite_id`, `type`; *allow, deny* |  |  |
| `DISCORDBOT_SYNC_GUILD_TEMPLATE` | Sync guild template | write | `guild_id`, `template_code` |  |  |
| `DISCORDBOT_TEST_AUTH` | Test bot token authentication | read |  |  |  |
| `DISCORDBOT_TRIGGER_TYPING_INDICATOR` | Trigger typing indicator | write | `channel_id` |  |  |
| `DISCORDBOT_UNBAN_USER_FROM_GUILD` | Unban user from guild | write | `guild_id`, `user_id` |  |  |
| `DISCORDBOT_UNPIN_MESSAGE` | Unpin Message | write | `channel_id`, `message_id` |  |  |
| `DISCORDBOT_UPDATE_AUTO_MODERATION_RULE` | Update auto moderation rule | write | `guild_id`, `rule_id`; *actions, enabled, event_type, exempt_channels, exempt_roles, name, trigger_metadata* |  |  |
| `DISCORDBOT_UPDATE_CHANNEL` | Update channel settings | write | `channel_id`; *available_tags, bitrate, default_auto_archive_duration, default_forum_layout, default_reaction_emoji, default_sort_order, default_thread_rate_limit_per_user, flags and 11 more* |  |  |
| `DISCORDBOT_UPDATE_GUILD` | Update guild settings | write | `guild_id`; *afk_channel_id, afk_timeout, banner, default_message_notifications, description, discovery_splash, explicit_content_filter, features and 14 more* |  |  |
| `DISCORDBOT_UPDATE_GUILD_APPLICATION_COMMAND` | Update guild application command | write | `application_id`, `command_id`, `guild_id`; *default_member_permissions, description, description_localizations, dm_permission, name, name_localizations, options* |  |  |
| `DISCORDBOT_UPDATE_GUILD_EMOJI` | Update guild emoji | write | `emoji_id`, `guild_id`; *name, roles* |  |  |
| `DISCORDBOT_UPDATE_GUILD_MEMBER` | Modify guild member details | write | `guild_id`, `user_id`; *channel_id, communication_disabled_until, deaf, flags, mute, nick, roles* |  |  |
| `DISCORDBOT_UPDATE_GUILD_ROLE` | Modify guild role | write | `guild_id`, `role_id`; *color, hoist, icon, mentionable, name, permissions, unicode_emoji* |  |  |
| `DISCORDBOT_UPDATE_GUILD_SCHEDULED_EVENT` | Update guild scheduled event | write | `event_id`, `guild_id`; *channel_id, description, entity_metadata, entity_type, image, name, privacy_level, scheduled_end_time and 2 more* |  |  |
| `DISCORDBOT_UPDATE_GUILD_STICKER` | Update guild sticker info | write | `guild_id`, `sticker_id`; *description, name, tags* |  |  |
| `DISCORDBOT_UPDATE_GUILD_TEMPLATE` | Update guild template | write | `guild_id`, `template_code`; *description, name* |  |  |
| `DISCORDBOT_UPDATE_GUILD_WELCOME_SCREEN` | Update guild welcome screen | write | `guild_id`; *description, enabled, welcome_channels* |  |  |
| `DISCORDBOT_UPDATE_GUILD_WIDGET_SETTINGS` | Update guild widget settings | write | `guild_id`; *channel_id, enabled* |  |  |
| `DISCORDBOT_UPDATE_MESSAGE` | Edit Message | write | `channel_id`, `message_id`; *allowed_mentions, attachments, components, content, embeds, flags* |  |  |
| `DISCORDBOT_UPDATE_MY_GUILD_MEMBER` | Update bot's nickname in guild | write | `guild_id`; *nick* |  |  |
| `DISCORDBOT_UPDATE_ORIGINAL_WEBHOOK_MESSAGE` | Update original webhook message | write | `webhook_id`, `webhook_token`; *allowed_mentions, attachments, components, content, embeds, flags, thread_id* |  |  |
| `DISCORDBOT_UPDATE_SELF_VOICE_STATE` | Update own voice state | write | `guild_id`; *channel_id, request_to_speak_timestamp, suppress* |  |  |
| `DISCORDBOT_UPDATE_VOICE_STATE` | Update user voice state | write | `guild_id`, `user_id`; *channel_id, suppress* |  |  |
| `DISCORDBOT_UPDATE_WEBHOOK` | Update webhook details | write | `webhook_id`; *avatar, channel_id, name* |  |  |
| `DISCORDBOT_UPDATE_WEBHOOK_BY_TOKEN` | Update webhook by token | write | `webhook_id`, `webhook_token`; *avatar, name* |  |  |
| `DISCORDBOT_UPDATE_WEBHOOK_MESSAGE` | Update webhook message | write | `message_id`, `webhook_id`, `webhook_token`; *allowed_mentions, attachments, components, content, embeds, flags, thread_id* |  |  |
