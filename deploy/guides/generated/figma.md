# Figma: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

52 tools: 40 read, 12 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `FIGMA_ADD_A_COMMENT_TO_A_FILE` | Add a comment to a file | write | `file_key`, `message`; *client_meta, comment_id* |  | `file_comments:write` |
| `FIGMA_ADD_A_REACTION_TO_A_COMMENT` | Add a reaction to a comment | write | `comment_id`, `emoji`, `file_key` |  | `file_comments:write` |
| `FIGMA_CREATE_A_WEBHOOK` | Create a webhook | write | `endpoint`, `event_type`, `passcode`; *context, context_id, description, status, team_id* |  | `webhooks:write` |
| `FIGMA_CREATE_DEV_RESOURCES` | Create dev resources | write | `dev_resources` |  | `file_dev_resources:write` |
| `FIGMA_CREATE_MODIFY_DELETE_VARIABLES` | Create, modify, or delete variables | write | `file_key`; *variableCollections, variableModeValues, variableModes, variables* |  | `file_variables:write` |
| `FIGMA_DELETE_A_COMMENT` | Delete a comment | write | `comment_id`, `file_key` |  | `file_comments:write` |
| `FIGMA_DELETE_A_REACTION` | Delete a reaction | write | `comment_id`, `emoji`, `file_key` |  | `file_comments:write` |
| `FIGMA_DELETE_A_WEBHOOK` | Delete a webhook | write | `webhook_id` |  | `webhooks:write` |
| `FIGMA_DELETE_DEV_RESOURCE` | Delete dev resource | write | `dev_resource_id`, `file_key` |  | `file_dev_resources:write` |
| `FIGMA_DESIGN_TOKENS_TO_TAILWIND` | Design tokens to tailwind | write | `tokens`; *config_format, include_font_imports, prefix* |  |  |
| `FIGMA_DETECT_BACKGROUND` | Detect Background Layers | read | `file_key`, `target_node_ids`; *search_depth* |  | `file_content:read` |
| `FIGMA_DISCOVER_FIGMA_RESOURCES` | Discover Figma Resources | read | *figma_url, file_key, max_depth, project_id, team_id* |  | `projects:read` |
| `FIGMA_DOWNLOAD_FIGMA_IMAGES` | Download Figma Images | read | `file_key`, `images`; *scale, svg_include_id, svg_outline_text, svg_simplify_stroke* |  | `file_content:read` |
| `FIGMA_EXTRACT_DESIGN_TOKENS` | Extract design tokens | read | `file_key`; *extract_from_nodes, include_local_styles, include_variables* |  | `file_read`, `file_variables:read` |
| `FIGMA_EXTRACT_PROTOTYPE_INTERACTIONS` | Extract Prototype Interactions | read | `file_key`; *analyze_components, include_animations* |  | `files:read` |
| `FIGMA_GET_ACTIVITY_LOGS` | Get activity logs | read | *end_time, events, limit, order, start_time* |  | `org:activity_log_read` |
| `FIGMA_GET_A_WEBHOOK` | Get a webhook | read | `webhook_id` |  | `webhooks:write` |
| `FIGMA_GET_COMMENTS_IN_A_FILE` | Get comments in a file | read | `file_key`; *as_md* |  | `files:read` |
| `FIGMA_GET_COMPONENT2` | Get component | read | `key` |  | `library_assets:read` |
| `FIGMA_GET_COMPONENT_SET` | Get component set | read | `key` |  | `files:read` |
| `FIGMA_GET_CURRENT_USER` | Get current user | read |  |  |  |
| `FIGMA_GET_DEV_RESOURCES` | Get dev resources | read | `file_key`; *node_ids* |  | `file_dev_resources:read` |
| `FIGMA_GET_FILES_IN_A_PROJECT` | Get files in a project | read | `project_id`; *branch_data* |  | `files:read` |
| `FIGMA_GET_FILE_COMPONENTS` | Get file components | read | `file_key` |  | `library_content:read` |
| `FIGMA_GET_FILE_COMPONENT_SETS` | Get file component sets | read | `file_key` |  | `library_content:read` |
| `FIGMA_GET_FILE_JSON` | Get file json | read | `file_key`; *branch_data, depth, geometry, ids, plugin_data, response_detail, version* |  | `file_content:read` |
| `FIGMA_GET_FILE_METADATA` | Get file metadata | read | `file_key` |  | `file_metadata:read` |
| `FIGMA_GET_FILE_NODES` | Get file nodes | read | `file_key`, `ids`; *depth, geometry, plugin_data, version* |  | `file_content:read` |
| `FIGMA_GET_FILE_STYLES` | Get file styles | read | `file_key` |  | `files:read` |
| `FIGMA_GET_IMAGE_FILLS` | Get image fills | read | `file_key` |  | `files:read` |
| `FIGMA_GET_LIBRARY_ANALYTICS_COMPONENT_ACTION_DATA` | Get library analytics component action data | read | `file_key`, `group_by`; *cursor, end_date, start_date* |  | `library_analytics:read` |
| `FIGMA_GET_LIBRARY_ANALYTICS_COMPONENT_USAGE_DATA` | Get library analytics component usage data | read | `file_key`, `group_by`; *cursor* |  | `library_analytics:read` |
| `FIGMA_GET_LIBRARY_ANALYTICS_STYLE_ACTION_DATA` | Get library analytics style action data | read | `file_key`, `group_by`; *cursor, end_date, start_date* |  | `library_analytics:read` |
| `FIGMA_GET_LIBRARY_ANALYTICS_STYLE_USAGE_DATA` | Get library analytics style usage data | read | `file_key`, `group_by`; *cursor* |  | `library_analytics:read` |
| `FIGMA_GET_LIBRARY_ANALYTICS_VARIABLE_ACTION_DATA` | Get library analytics variable action data | read | `file_key`, `group_by`; *cursor, end_date, start_date* |  | `library_analytics:read` |
| `FIGMA_GET_LIBRARY_ANALYTICS_VARIABLE_USAGE_DATA` | Get library analytics variable usage data | read | `file_key`, `group_by`; *cursor* |  | `library_analytics:read` |
| `FIGMA_GET_LOCAL_VARIABLES` | Get local variables | read | `file_key` |  | `file_variables:read` |
| `FIGMA_GET_PAYMENTS` | Get payments | read | *community_file_id, plugin_id, plugin_payment_token, user_id, widget_id* |  |  |
| `FIGMA_GET_PROJECTS_IN_A_TEAM` | Get projects in a team | read | `team_id` |  | `projects:read` |
| `FIGMA_GET_PUBLISHED_VARIABLES` | Get published variables | read | `file_key` |  | `file_variables:read` |
| `FIGMA_GET_REACTIONS_FOR_A_COMMENT` | Get reactions for a comment | read | `comment_id`, `file_key`; *cursor* |  | `file_comments:write` |
| `FIGMA_GET_SCIM_SERVICE_PROVIDER_CONFIG` | Get SCIM service provider config | read |  |  |  |
| `FIGMA_GET_STYLE` | Get style | read | `key` |  | `files:read` |
| `FIGMA_GET_TEAM_COMPONENTS` | Get team components | read | `team_id`; *after, before, page_size* |  | `files:read`, `team_library_content:read` |
| `FIGMA_GET_TEAM_COMPONENT_SETS` | Get team component sets | read | `team_id`; *after, before, page_size* |  | `team_library_content:read` |
| `FIGMA_GET_TEAM_STYLES` | Get team styles | read | `team_id`; *after, before, page_size* |  | `files:read` |
| `FIGMA_GET_TEAM_WEBHOOKS` | Get webhooks | read | `context_id`; *context* |  | `webhooks:read` |
| `FIGMA_GET_VERSIONS_OF_A_FILE` | Get versions of a file | read | `file_key`; *after, before, page_size* |  | `file_versions:read` |
| `FIGMA_GET_WEBHOOK_REQUESTS` | Get webhook requests | read | `webhook_id` |  | `webhooks:read` |
| `FIGMA_RENDER_IMAGES_OF_FILE_NODES` | Render images of file nodes | read | `file_key`, `ids`; *contents_only, format, scale, svg_include_id, svg_include_node_id, svg_outline_text, svg_simplify_stroke, use_absolute_bounds and 1 more* |  | `file_content:read` |
| `FIGMA_UPDATE_A_WEBHOOK` | Update a webhook | write | `webhook_id`; *description, endpoint, event_type, passcode, status* |  | `webhooks:write` |
| `FIGMA_UPDATE_DEV_RESOURCES` | Update dev resources | write | `dev_resources` |  | `file_dev_resources:write` |
