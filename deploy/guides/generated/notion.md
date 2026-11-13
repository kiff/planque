# Notion: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

56 tools: 28 read, 28 write (a tool not marked read-only counts as a write). 1 more couldn't be read and are left out.

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `NOTION_ADD_MULTIPLE_PAGE_CONTENT` | Add multiple content blocks (bulk, user-friendly) | write | `content_blocks`, `parent_block_id`; *after* |  | `Insert content` |
| `NOTION_APPEND_CODE_BLOCKS` | Append code blocks (code, quote, equation) | write | `block_id`, `children`; *after* |  |  |
| `NOTION_APPEND_LAYOUT_BLOCKS` | Append layout blocks (divider, TOC, columns) | write | `block_id`, `children`; *after* |  |  |
| `NOTION_APPEND_MEDIA_BLOCKS` | Append media blocks (image, video, audio, files) | write | `block_id`, `children`; *after* |  |  |
| `NOTION_APPEND_TABLE_BLOCKS` | Append table blocks | write | `block_id`, `tables`; *after* |  |  |
| `NOTION_APPEND_TASK_BLOCKS` | Append task blocks (to-do, toggle, callout) | write | `block_id`, `children`; *after* |  |  |
| `NOTION_APPEND_TEXT_BLOCKS` | Append text blocks (paragraphs, headings, lists) | write | `block_id`, `children`; *after* |  |  |
| `NOTION_ARCHIVE_NOTION_PAGE` | Archive Notion Page | write | `page_id`; *archive* |  | `read_content`, `update_content` |
| `NOTION_CREATE_COMMENT` | Create comment | write | `comment`; *discussion_id, parent_page_id* |  | `comment:write` |
| `NOTION_CREATE_DATABASE` | Create Notion Database | write | `parent_id`, `title`; *properties* |  |  |
| `NOTION_CREATE_FILE_UPLOAD` | Create Notion file upload | write | *content_type, external_url, filename, mode, number_of_parts* |  |  |
| `NOTION_CREATE_NOTION_PAGE` | Create Notion page | write | `parent_id`, `title`; *cover, icon, markdown* |  | `insert_content` |
| `NOTION_CREATE_VIEW` | Create Notion view | write | `data_source_id`, `name`, `type`; *configuration, create_database, database_id, filter, placement, position, quick_filters, sorts and 1 more* |  |  |
| `NOTION_CREATE_VIEW_QUERY` | Create Notion view query | read | `view_id`; *page_size* |  |  |
| `NOTION_DELETE_BLOCK` | Delete a block | write | `block_id` |  | `read_content`, `update_content` |
| `NOTION_DELETE_VIEW` | Delete Notion view | write | `view_id` |  |  |
| `NOTION_DELETE_VIEW_QUERY` | Delete Notion view query | write | `query_id`, `view_id` |  |  |
| `NOTION_DUPLICATE_PAGE` | Duplicate page | write | `page_id`, `parent_id`; *title* |  |  |
| `NOTION_FETCH_ALL_BLOCK_CONTENTS` | Fetch All Notion Block Contents | read | *block_id, max_blocks, max_depth, page_size, page_url, recursive* |  | `read:content` |
| `NOTION_FETCH_BLOCK_CONTENTS` | Fetch Notion Block Children | read | `block_id`; *page_size, start_cursor* |  | `read_content` |
| `NOTION_FETCH_BLOCK_METADATA` | Fetch Notion block metadata | read | `block_id` |  |  |
| `NOTION_FETCH_COMMENTS` | Fetch comments | read | *block_id, page_id, page_size, start_cursor* |  | `insert_content` |
| `NOTION_FETCH_DATA` | Fetch Notion Data | read | `fetch_type`; *original_page_size, page_size, page_size_was_capped, query, start_cursor* |  |  |
| `NOTION_FETCH_DATABASE` | Fetch Database | read | `database_id` |  |  |
| `NOTION_FETCH_ROW` | Fetch database row | read | `page_id` |  |  |
| `NOTION_GET_ABOUT_USER` | Get about user | read | `user_id` |  |  |
| `NOTION_GET_PAGE_MARKDOWN` | Get page markdown | read | `page_id`; *include_transcript* |  |  |
| `NOTION_GET_PAGE_PROPERTY_ACTION` | Get page property | read | `page_id`, `property_id`; *page_size, start_cursor* |  | `read_content` |
| `NOTION_GET_PLUGIN` | Get plugin | read | `plugin_id` |  | `read_content` |
| `NOTION_GET_VIEW_QUERY_RESULTS` | Get Notion view query results | read | `query_id`, `view_id`; *page_size, start_cursor* |  |  |
| `NOTION_INSERT_ROW_DATABASE` | Insert row database | write | `database_id`; *child_blocks, cover, icon, properties* |  |  |
| `NOTION_INSERT_ROW_FROM_NL` | Insert Row From Natural Language | write | `database_id`, `nl_query`; *cover, icon* |  |  |
| `NOTION_LIST_DATA_SOURCE_TEMPLATES` | List data source templates | read | `data_source_id`; *page_size, start_cursor* |  |  |
| `NOTION_LIST_FILE_UPLOADS` | List Notion file uploads | read | *page_size, start_cursor* |  |  |
| `NOTION_LIST_PLUGINS` | List plugins | read | *page_size, start_cursor* |  | `read_content` |
| `NOTION_LIST_USERS` | List users | read | *page_size, start_cursor* |  |  |
| `NOTION_LIST_VIEWS` | List Notion views | read | *data_source_id, database_id, page_size, start_cursor* |  |  |
| `NOTION_MOVE_PAGE` | Move Page | write | `page_id`, `parent` |  | `update_content` |
| `NOTION_QUERY_DATABASE` | Query database | read | `database_id`; *page_size, sorts, start_cursor* |  | `databases.read` |
| `NOTION_QUERY_DATABASE_WITH_FILTER` | Query database with filter | read | `database_id`; *composio_execution_message, filter, page_size, sorts, start_cursor* |  | `read_content` |
| `NOTION_QUERY_DATA_SOURCE` | Query data source | read | `data_source_id`; *filter, filter_properties, page_size, sorts, start_cursor* |  | `read_content` |
| `NOTION_RETRIEVE_COMMENT` | Retrieve Comment | read | `comment_id` |  |  |
| `NOTION_RETRIEVE_DATABASE_PROPERTY` | Retrieve Database Property | read | `database_id`, `property_id` |  |  |
| `NOTION_RETRIEVE_FILE_UPLOAD` | Retrieve Notion file upload | read | `file_upload_id` |  |  |
| `NOTION_RETRIEVE_PAGE` | Retrieve page | read | `page_id` |  |  |
| `NOTION_RETRIEVE_VIEW` | Retrieve Notion view | read | `view_id` |  |  |
| `NOTION_SEARCH_NOTION_PAGE` | Search Notion pages and databases | read | *direction, filter_properties, filter_property, filter_value, page_size, query, start_cursor, timestamp* |  |  |
| `NOTION_SEND_FILE_UPLOAD` | Send file upload | write | `file`, `file_upload_id`; *file_content_base64, part_number* |  |  |
| `NOTION_UPDATE_BLOCK` | Update block | write | `block_id`, `content`; *additional_properties, block_type, language* |  | `read_content`, `update_content` |
| `NOTION_UPDATE_PAGE` | Update Page | write | `page_id`; *archived, cover, icon, properties* |  |  |
| `NOTION_UPDATE_PAGE_MARKDOWN` | Update page markdown | write | `page_id`, `type`; *replace_content, update_content* |  | `update_content` |
| `NOTION_UPDATE_ROW_DATABASE` | Update Database Row (Page) | write | `row_id`; *cover, delete_row, icon, properties* |  | `read_content`, `update_content` |
| `NOTION_UPDATE_SCHEMA_DATABASE` | Update database schema | write | `database_id`; *description, properties, title* |  |  |
| `NOTION_UPDATE_VIEW` | Update Notion view | write | `view_id`; *configuration, filter, name, quick_filters, sorts* |  |  |
| `NOTION_UPSERT_ROW_DATABASE` | Upsert database rows | write | `items`; *data_source_id, database_id, options* |  | `insert_content`, `read_content`, `update_content` |
| `NOTION_WHO_AM_I` | Who Am I | read |  |  |  |
