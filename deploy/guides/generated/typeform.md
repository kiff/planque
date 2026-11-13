# Typeform: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

35 tools: 16 read, 19 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `TYPEFORM_CREATE_ACCOUNT_WORKSPACE` | Create Account Workspace | write | `account_id`, `name` |  | `workspaces:write` |
| `TYPEFORM_CREATE_FORM` | Create Form | write | `title`; *fields, logic, settings, thankyou_screens, theme, type, welcome_screens, workspace* |  | `forms:write` |
| `TYPEFORM_CREATE_IMAGE` | Create Image | write | `file_name`; *image, url* |  | `images:write` |
| `TYPEFORM_CREATE_OR_UPDATE_WEBHOOK` | Create or Update Webhook | write | `form_id`, `tag`, `url`; *enabled, event_types, secret, verify_ssl* |  | `webhooks:write` |
| `TYPEFORM_CREATE_THEME` | Create Theme | write | `colors`, `fields`, `font`; *background, has_transparent_button, name, rounded_corners* |  | `themes:write` |
| `TYPEFORM_CREATE_WORKSPACE` | Create Workspace | write | `name` |  | `workspaces:write` |
| `TYPEFORM_DELETE_FORM` | Delete Form | write | `form_id` |  | `forms:write` |
| `TYPEFORM_DELETE_IMAGE` | Delete Image | write | `image_id` |  | `images:write` |
| `TYPEFORM_DELETE_RESPONSES` | Delete Responses | write | `form_id`, `included_response_ids` |  | `responses:write` |
| `TYPEFORM_DELETE_THEME` | Delete Theme | write | `theme_id` |  | `themes:write` |
| `TYPEFORM_DELETE_WEBHOOK` | Delete Webhook | write | `form_id`, `tag` |  | `webhooks:write` |
| `TYPEFORM_DELETE_WORKSPACE` | Delete Workspace | write | `workspace_id` |  | `workspaces:write` |
| `TYPEFORM_GET_ABOUT_ME` | Get About Me | read |  |  | `accounts:read` |
| `TYPEFORM_GET_ALL_RESPONSE_FILES` | Get All Response Files | read | `form_id` |  | `responses:read` |
| `TYPEFORM_GET_BACKGROUND_BY_SIZE` | Get Background By Size | read | `image_id`, `size` |  | `images:read` |
| `TYPEFORM_GET_CHOICE_IMAGE_BY_SIZE` | Get Choice Image By Size | read | `image_id`, `size` |  | `images:read` |
| `TYPEFORM_GET_FORM` | Get Form | read | `form_id` |  | `forms:read` |
| `TYPEFORM_GET_FORM_MESSAGES` | Get Form Messages | read | `form_id` |  | `forms:read` |
| `TYPEFORM_GET_FORM_RESPONSES` | Get Form Responses | read | `form_id`; *after, answered_fields, before, excluded_response_ids, fields, included_response_ids, page_size, query and 4 more* |  | `responses:read` |
| `TYPEFORM_GET_IMAGE_BY_SIZE` | Get Image By Size | read | `image_id`, `size` |  | `images:read` |
| `TYPEFORM_GET_THEME` | Get Theme | read | `theme_id` |  | `themes:read` |
| `TYPEFORM_GET_WEBHOOK` | Get Webhook | read | `form_id`, `tag` |  | `webhooks:read` |
| `TYPEFORM_GET_WORKSPACE` | Get Workspace | read | `workspace_id` |  | `workspaces:read` |
| `TYPEFORM_LIST_FORMS` | List Forms | read | *order_by, page, page_size, search, sort_by, workspace_id* |  | `forms:read` |
| `TYPEFORM_LIST_IMAGES` | List Images | read |  |  | `images:read` |
| `TYPEFORM_LIST_THEMES` | List Typeform Themes | read | *page, page_size* |  | `images:write` |
| `TYPEFORM_LIST_WEBHOOKS` | List Form Webhooks | read | `form_id` |  | `webhooks:read` |
| `TYPEFORM_LIST_WORKSPACES` | List Workspaces | read | *page, page_size, search* |  | `workspaces:read` |
| `TYPEFORM_PATCH_FORM` | Patch Form | write | `form_id`, `operations` |  | `forms:write` |
| `TYPEFORM_PATCH_THEME` | Update Theme (Partial) | write | `theme_id`; *background, colors, fields, font, has_transparent_button, name, rounded_corners, screens* |  | `themes:write` |
| `TYPEFORM_UPDATE_FORM` | Update Form | write | `form_id`; *fields, hidden, logic, settings, thankyou_screens, theme, title, type and 3 more* |  | `forms:write` |
| `TYPEFORM_UPDATE_FORM_MESSAGES` | Update Form Messages | write | `form_id`; *block_short_text_placeholder, label_button_hint_default, label_button_no_answer_default, label_button_submit, label_error_required* |  | `forms:write` |
| `TYPEFORM_UPDATE_THEME` | Update Theme | write | `colors`, `font`, `name`, `theme_id`; *background, fields, has_transparent_button, rounded_corners, screens* |  | `themes:write` |
| `TYPEFORM_UPDATE_WORKSPACE` | Update Workspace | write | `operations`, `workspace_id` |  | `workspaces:write` |
| `TYPEFORM_UPLOAD_VIDEO` | Upload Video | write | `field_id`, `form_id`, `language` |  | `forms:write` |
