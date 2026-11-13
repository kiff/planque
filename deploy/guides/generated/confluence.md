# Confluence: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

70 tools: 44 read, 26 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `CONFLUENCE_ADD_CONTENT_LABEL` | Add Content Label | write | `id`, `labels` |  | `write:confluence-content` |
| `CONFLUENCE_CQL_SEARCH` | CQL Search | read | `cql`; *excerpt, expand, includeArchivedSpaces, limit, start* |  | `search:confluence`, `read:content-details:confluence` |
| `CONFLUENCE_CREATE_BLOGPOST` | Create Blogpost | write | `body`, `spaceId`, `status`, `title` |  | `write:page:confluence` |
| `CONFLUENCE_CREATE_BLOGPOST_PROPERTY` | Create Blogpost Property | write | `id`, `key`, `value` |  | `read:page:confluence`, `write:page:confluence` |
| `CONFLUENCE_CREATE_CONTENT_PROPERTY_FOR_WHITEBOARD` | Create Whiteboard Property | write | `key`, `value`, `whiteboard_id` |  | `read:whiteboard:confluence`, `write:whiteboard:confluence` |
| `CONFLUENCE_CREATE_FOOTER_COMMENT` | Create Footer Comment | write | `body`; *attachmentId, blogPostId, customContentId, pageId, parentCommentId* |  | `write:comment:confluence` |
| `CONFLUENCE_CREATE_INLINE_COMMENT` | Create Inline Comment | write | `body`; *blogPostId, inlineCommentProperties, pageId, parentCommentId* |  | `write:comment:confluence` |
| `CONFLUENCE_CREATE_PAGE` | Create Page | write | `spaceId`, `title`; *body, parentId* |  | `read:space:confluence`, `write:page:confluence` |
| `CONFLUENCE_CREATE_PAGE_PROPERTY` | Create Page Property | write | `key`, `page_id`, `value` |  | `read:page:confluence`, `write:page:confluence` |
| `CONFLUENCE_CREATE_PRIVATE_SPACE` | Create Private Space | write | `key`, `name`; *description* |  | `write:confluence-space`, `read:space-details:confluence`, `write:space:confluence`, `write:space.permission:confluence` |
| `CONFLUENCE_CREATE_SPACE` | Create Space | write | `key`, `name`; *description, metadata, type* |  | `write:confluence-space` |
| `CONFLUENCE_CREATE_SPACE_PROPERTY` | Create Space Property | write | `key`, `spaceId`, `value`; *version* |  | `read:space:confluence`, `write:space:confluence` |
| `CONFLUENCE_CREATE_WHITEBOARD` | Create Whiteboard | write | `spaceId`; *parentId, title* |  | `write:whiteboard:confluence` |
| `CONFLUENCE_DELETE_BLOGPOST` | Delete Blog Post | write | `id`; *purge* |  | `delete:page:confluence` |
| `CONFLUENCE_DELETE_BLOGPOST_PROPERTY` | Delete Blogpost Property | write | `id`, `propertyId` |  | `read:page:confluence`, `write:page:confluence` |
| `CONFLUENCE_DELETE_CONTENT_PROPERTY_FOR_PAGE_BY_ID` | Delete Page Content Property | write | `pageId`, `propertyId` |  | `read:page:confluence`, `write:page:confluence` |
| `CONFLUENCE_DELETE_CONTENT_PROPERTY_FOR_WHITEBOARD_BY_ID` | Delete Whiteboard Content Property | write | `propertyId`, `whiteboardId` |  | `read:whiteboard:confluence`, `write:whiteboard:confluence` |
| `CONFLUENCE_DELETE_PAGE` | Delete Page | write | `id` |  | `delete:page:confluence` |
| `CONFLUENCE_DELETE_SPACE` | Delete Space | write | `space_key` |  | `write:confluence-space`, `read:content.metadata:confluence`, `delete:space:confluence` |
| `CONFLUENCE_DELETE_SPACE_PROPERTY` | Delete Space Property | write | `propertyId`, `spaceIdOrKey` |  | `read:space:confluence`, `write:space:confluence` |
| `CONFLUENCE_DOWNLOAD_ATTACHMENT` | Download Attachment | read | `attachmentId`, `pageId` |  | `read:attachment:confluence` |
| `CONFLUENCE_FIND_TEXT_SELECTIONS` | Find Text Selections | read | `page_id`, `search_text`; *case_insensitive, context_chars, max_results* |  | `read:page:confluence` |
| `CONFLUENCE_GET_ATTACHMENTS` | Get Attachments | read | `pageId`; *cursor, limit, mediaType* |  | `read:attachment:confluence` |
| `CONFLUENCE_GET_ATTACHMENT_LABELS` | Get Attachment Labels | read | `id`; *cursor, limit, prefix, sort* |  | `read:attachment:confluence` |
| `CONFLUENCE_GET_AUDIT_LOGS` | Get Audit Logs | read | *endDate, limit, searchString, start, startDate* |  | `read:audit-log:confluence` |
| `CONFLUENCE_GET_BLOGPOST_BY_ID` | Get Blogpost by ID | read | `id` |  | `read:page:confluence` |
| `CONFLUENCE_GET_BLOGPOST_LABELS` | Get Blogpost Labels | read | `id`; *cursor, limit* |  | `read:page:confluence` |
| `CONFLUENCE_GET_BLOGPOST_LIKE_COUNT` | Get Blogpost Like Count | read | `id` |  | `read:page:confluence` |
| `CONFLUENCE_GET_BLOGPOST_OPERATIONS` | Get Blogpost Operations | read | `id` |  | `read:page:confluence` |
| `CONFLUENCE_GET_BLOGPOST_VERSIONS` | Get Blogpost Versions | read | `blogpostId`; *cursor, limit* |  | `read:page:confluence` |
| `CONFLUENCE_GET_BLOGPOST_VERSION_DETAILS` | Get Blogpost Version Details | read | `blogpostId`, `versionNumber` |  | `read:page:confluence` |
| `CONFLUENCE_GET_BLOG_POSTS` | Get Blog Posts | read | *cursor, limit* |  | `read:page:confluence` |
| `CONFLUENCE_GET_BLOG_POSTS_FOR_LABEL` | Get Blog Posts For Label | read | `id`; *cursor, limit* |  | `read:page:confluence` |
| `CONFLUENCE_GET_CHILD_PAGES` | Get Child Pages | read | `id`; *cursor, limit, sort* |  | `read:page:confluence` |
| `CONFLUENCE_GET_CONTENT_PROPERTIES_FOR_BLOG_POST` | Get Blog Post Content Properties | read | `blogpost_id`; *cursor, limit* |  | `read:page:confluence` |
| `CONFLUENCE_GET_CONTENT_PROPERTIES_FOR_PAGE` | Get Page Content Properties | read | `page_id`; *cursor, limit* |  | `read:page:confluence` |
| `CONFLUENCE_GET_CONTENT_RESTRICTIONS` | Get Content Restrictions | read | `id`; *expand* |  | `read:confluence-content.all`, `read:content-details:confluence` |
| `CONFLUENCE_GET_CURRENT_USER` | Get Current User | read |  |  | `read:confluence-user` |
| `CONFLUENCE_GET_INLINE_COMMENTS_FOR_BLOG_POST` | Get Inline Comments for Blog Post | read | `id`; *body_format, cursor, limit, sort* |  | `read:comment:confluence` |
| `CONFLUENCE_GET_LABELS` | Get Labels | read | *cursor, limit* |  | `read:label:confluence` |
| `CONFLUENCE_GET_LABELS_FOR_PAGE` | Get Page Labels | read | `id`; *cursor, limit* |  | `read:page:confluence` |
| `CONFLUENCE_GET_LABELS_FOR_SPACE` | Get Labels for Space | read | `id`; *cursor, limit, prefix, sort* |  | `read:space:confluence` |
| `CONFLUENCE_GET_LABELS_FOR_SPACE_CONTENT` | Get Labels for Space Content | read | `id`; *cursor, limit, prefix* |  | `read:space:confluence` |
| `CONFLUENCE_GET_PAGES` | Get Pages | read | *body_format, cursor, limit, sort, space_id, status, subtype, title* |  | `read:page:confluence` |
| `CONFLUENCE_GET_PAGE_ANCESTORS` | Get Page Ancestors | read | `id`; *limit* |  | `read:content.metadata:confluence` |
| `CONFLUENCE_GET_PAGE_BY_ID` | Get Page by ID | read | `id`; *draft, version* |  | `read:page:confluence` |
| `CONFLUENCE_GET_PAGE_FOOTER_COMMENTS` | Get Page Footer Comments | read | `id`; *body_format, cursor, limit, sort, status* |  | `read:comment:confluence` |
| `CONFLUENCE_GET_PAGE_INLINE_COMMENTS` | Get Page Inline Comments | read | `id`; *body_format, cursor, limit, resolution_status, sort, status* |  | `read:comment:confluence` |
| `CONFLUENCE_GET_PAGE_LIKE_COUNT` | Get Page Like Count | read | `id` |  | `read:page:confluence` |
| `CONFLUENCE_GET_PAGE_VERSIONS` | Get Page Versions | read | `id`; *cursor, limit* |  | `read:page:confluence` |
| `CONFLUENCE_GET_SPACES` | Get Spaces | read | *cursor, expand, label, limit, spaceKey, status, type* |  | `read:space:confluence` |
| `CONFLUENCE_GET_SPACE_BY_ID` | Get Space by ID | read | `id` |  | `read:space:confluence` |
| `CONFLUENCE_GET_SPACE_CONTENTS` | Get Space Contents | read | `spaceKey`; *cursor, limit, status, type* |  | `read:page:confluence`, `read:space:confluence` |
| `CONFLUENCE_GET_SPACE_PROPERTIES` | Get Space Properties | read | `id`; *cursor, key, limit* |  | `read:space:confluence` |
| `CONFLUENCE_GET_TASKS` | Get Tasks | read | *assigned_to, blogpost_id, body_format, completed_by, created_by, cursor, include_blank_tasks, limit and 4 more* |  | `read:task:confluence` |
| `CONFLUENCE_GET_USER_ANONYMOUS` | Get Anonymous User | read |  |  | `read:confluence-user`, `read:user:confluence` |
| `CONFLUENCE_LIST_SPACES` | List Spaces | read | *cursor, limit* |  | `read:space:confluence` |
| `CONFLUENCE_SEARCH_CONTENT` | Search Content | read | `query`; *cursor, expand, limit, spaceKey, start* |  | `read:page:confluence` |
| `CONFLUENCE_SEARCH_PAGES` | Search Pages | read | *body_format, cursor, limit, space_id, status, title* |  | `read:page:confluence` |
| `CONFLUENCE_SEARCH_SPACES` | Search Spaces | read | *cursor, expand, label, limit, query, spaceKey, status, type* |  | `read:space:confluence` |
| `CONFLUENCE_SEARCH_USERS` | Search Users | read | *accountId, limit, start* |  | `read:content-details:confluence` |
| `CONFLUENCE_UPDATE_BLOGPOST` | Update Blogpost | write | `body`, `id`, `spaceId`, `title`, `version`; *status* |  | `write:page:confluence` |
| `CONFLUENCE_UPDATE_BLOGPOST_PROPERTY` | Update Blogpost Property | write | `id`, `propertyId`, `value` |  | `read:page:confluence`, `write:page:confluence` |
| `CONFLUENCE_UPDATE_CONTENT_PROPERTY_FOR_PAGE_BY_ID` | Update Page Content Property | write | `pageId`, `propertyId`, `value` |  | `read:page:confluence`, `write:page:confluence` |
| `CONFLUENCE_UPDATE_CONTENT_PROPERTY_FOR_WHITEBOARD_BY_ID` | Update Whiteboard Content Property | write | `propertyId`, `value`, `whiteboardId` |  | `read:whiteboard:confluence`, `write:whiteboard:confluence` |
| `CONFLUENCE_UPDATE_PAGE` | Update Page | write | `body`, `id`, `title`, `version`; *spaceId* |  | `read:page:confluence`, `write:page:confluence` |
| `CONFLUENCE_UPDATE_SPACE_PROPERTY` | Update Space Property | write | `propertyId`, `spaceId`, `value` |  | `read:space:confluence`, `write:space:confluence` |
| `CONFLUENCE_UPDATE_TASK` | Update Task | write | `id`, `status`; *body_format* |  | `write:task:confluence` |
| `CONFLUENCE_VALIDATE_CREDENTIAL` | Validate Credential | read |  |  | `read:confluence-user`, `read:content-details:confluence` |
| `CONFLUENCE_WHO_AM_I` | Who Am I | read |  |  |  |
