# Discord: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

23 tools: 18 read, 5 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `DISCORD_DELETE_USER_APPLICATION_ROLE_CONNECTION` | Delete User Application Role Connection | write | `application_id` |  | `role_connections.write` |
| `DISCORD_EDIT_APPLICATION_COMMAND_PERMISSIONS` | Edit Application Command Permissions | write | `application_id`, `command_id`, `guild_id`, `permissions` |  | `applications.commands.permissions.update` |
| `DISCORD_GET_APPLICATION_COMMAND_PERMISSIONS` | Get Application Command Permissions | read | `application_id`, `command_id`, `guild_id` |  | `applications.commands.permissions.update` |
| `DISCORD_GET_BATCH_APPLICATION_COMMAND_PERMISSIONS` | Get Batch Application Command Permissions | read | `application_id`, `guild_id` |  | `applications.commands.permissions.update` |
| `DISCORD_GET_CURRENT_USER_APPLICATION_ENTITLEMENTS` | Get Current User Application Entitlements | read | `application_id`; *after, before, exclude_deleted, exclude_ended, guild_id, limit, sku_ids* |  | `applications.entitlements` |
| `DISCORD_GET_GATEWAY` | Get Gateway | read |  |  |  |
| `DISCORD_GET_GUILD_TEMPLATE` | Get Guild Template | write | `code` |  |  |
| `DISCORD_GET_GUILD_WIDGET` | Get Guild Widget | read | `guild_id` |  |  |
| `DISCORD_GET_GUILD_WIDGET_PNG` | Get Guild Widget PNG | read | `guild_id`; *style* |  |  |
| `DISCORD_GET_MY_GUILD_MEMBER` | Get my guild member | read | `guild_id` |  | `guilds.members.read` |
| `DISCORD_GET_MY_OAUTH2_AUTHORIZATION` | Get my OAuth2 authorization | read |  |  | `identify` |
| `DISCORD_GET_MY_USER` | Get My User | read |  |  | `identify` |
| `DISCORD_GET_OPENID_CONNECT_USERINFO` | Get OpenID Connect userinfo | read |  |  | `openid`, `identify`, `email` |
| `DISCORD_GET_PUBLIC_KEYS` | Get Public Keys | read |  |  |  |
| `DISCORD_GET_SKU_SUBSCRIPTION` | Get SKU Subscription | read | `sku_id`, `subscription_id` |  |  |
| `DISCORD_GET_USER` | Get User | read | `user_id` |  |  |
| `DISCORD_GET_USER_APPLICATION_ROLE_CONNECTION` | Get User Application Role Connection | read | `application_id` |  | `role_connections.write` |
| `DISCORD_INVITE_RESOLVE` | Resolve Invite | write | `code`; *guild_scheduled_event_id, with_counts, with_expiration* |  |  |
| `DISCORD_LIST_MY_CONNECTIONS` | List My Connections | read |  |  | `connections` |
| `DISCORD_LIST_MY_GUILDS` | List My Guilds | read | *after, before, limit, with_counts* |  | `guilds` |
| `DISCORD_LIST_SKU_SUBSCRIPTIONS` | List SKU Subscriptions | read | `sku_id`; *after, before, limit, user_id* |  |  |
| `DISCORD_LIST_STICKER_PACKS` | List Sticker Packs | read |  |  |  |
| `DISCORD_UPDATE_USER_APPLICATION_ROLE_CONNECTION` | Update User Application Role Connection | write | `application_id`; *metadata, platform_name, platform_username* |  | `role_connections.write` |
