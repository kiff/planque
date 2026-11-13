# Google Docs: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

46 tools: 9 read, 37 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `GOOGLEDOCS_COPY_DOCUMENT` | Copy Google Document | write | `document_id`; *include_shared_drives, title* |  | `drive`, `drive.appdata`, `drive.file`, `drive.photos.readonly` |
| `GOOGLEDOCS_CREATE_AND_POPULATE_TABLE` | Create and Populate Table in Google Doc | write | `cellData`, `columns`, `documentId`, `rows`; *index, insertAtEndOfSegment, segmentId, tabId* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_CREATE_DOCUMENT` | Create a document | write | `title`; *text* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_CREATE_DOCUMENT_MARKDOWN` | Create Document Markdown | write | `title`; *image_assets, markdown_text* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_CREATE_DRIVE_FILE` | Create Drive File | write | *app_properties, description, fields, mime_type, name, parent_folder_id, properties, shortcut_target_id and 2 more* |  | `drive`, `drive.file` |
| `GOOGLEDOCS_CREATE_FOOTER` | Create Footer | write | `document_id`, `type`; *section_break_location* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_CREATE_FOOTNOTE` | Create Footnote | write | `documentId`; *endOfSegmentLocation, location* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_CREATE_HEADER` | Create Header | write | `documentId`; *sectionBreakLocation, text, type* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_CREATE_NAMED_RANGE` | Create Named Range | write | `documentId`, `name`, `rangeEndIndex`, `rangeStartIndex`; *rangeSegmentId, tabId* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_CREATE_PARAGRAPH_BULLETS` | Create Paragraph Bullets | write | `createParagraphBullets`, `document_id`; *tabId* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_CREATE_TAB` | Create Document Tab | write | `document_id`; *tab_properties* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_DELETE_CONTENT_RANGE` | Delete Content Range in Document | write | `document_id`, `range` |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_DELETE_FOOTER` | Delete Footer | write | `document_id`, `footer_id`; *tab_id* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_DELETE_HEADER` | Delete Header | write | `document_id`, `header_id`; *tab_id* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_DELETE_NAMED_RANGE` | Delete Named Range | write | `deleteNamedRange`, `document_id` |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_DELETE_PARAGRAPH_BULLETS` | Delete Paragraph Bullets | write | `document_id`, `range`; *tab_id* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_DELETE_TAB` | Delete Tab | write | `document_id`, `tab_id` |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_DELETE_TABLE_COLUMN` | Delete Table Column | write | `document_id`, `requests` |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_DELETE_TABLE_ROW` | Delete Table Row | write | `documentId`, `tableCellLocation` |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_EXPORT_DOCUMENT_AS_PDF` | Export Google Doc as PDF | read | `file_id`; *filename* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.readonly` |
| `GOOGLEDOCS_GET_DOCS_DRIVE_FILE_METADATA` | Get Drive File Metadata | read | `fileId`; *fields, includeLabels, includePermissionsForView, supportsAllDrives* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDOCS_GET_DOCS_SPREADSHEET_VALUES` | Get Spreadsheet Values | read | `range`, `spreadsheet_id`; *date_time_render_option, major_dimension, value_render_option* |  | `drive`, `drive.readonly`, `drive.file`, `spreadsheets`, `spreadsheets.readonly` |
| `GOOGLEDOCS_GET_DOCUMENT_BY_ID` | Get document by id | read | `id`; *includeTabsContent* |  | `documents`, `documents.readonly`, `drive`, `drive.readonly`, `drive.file` |
| `GOOGLEDOCS_GET_DOCUMENT_END_INDEX` | Get Document End Index | read | `document_id`; *segment_id, tab_id* |  | `documents`, `documents.readonly`, `drive`, `drive.readonly`, `drive.file` |
| `GOOGLEDOCS_GET_DOCUMENT_PLAINTEXT` | Get document plain text | read | `document_id`; *include_footers, include_footnotes, include_headers, include_tables, include_tabs_content, tab_ids, table_cell_delimiter, table_row_delimiter* |  | `documents`, `documents.readonly`, `drive`, `drive.readonly`, `drive.file` |
| `GOOGLEDOCS_INSERT_IMAGE_IN_TABLE_CELL` | Insert Image in Table Cell | write | `document_id`, `table_cell_location`, `uri`; *cell_content_index, object_size* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_INSERT_INLINE_IMAGE` | Insert Inline Image | write | `documentId`, `location`, `uri`; *objectSize* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_INSERT_PAGE_BREAK` | Insert Page Break | write | `documentId`, `insertPageBreak` |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_INSERT_TABLE_ACTION` | Insert Table in Google Doc | write | `columns`, `documentId`, `rows`; *index, insertAtEndOfSegment, segmentId, tabId* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_INSERT_TABLE_COLUMN` | Insert Table Column | write | `document_id`, `requests` |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_INSERT_TABLE_ROW` | Insert Table Row | write | `document_id`, `requests` |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_INSERT_TEXT_ACTION` | Insert Text into Document | write | `document_id`, `text_to_insert`; *append_to_end, insertion_index, segment_id, tab_id* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_INSERT_TEXT_IN_TABLE_CELL` | Insert Text in Table Cell | write | `column_index`, `document_id`, `row_index`, `table_start_index`, `text`; *tabId* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_LIST_DRIVE_ITEM_PERMISSIONS` | List Drive Item Permissions | read | `fileId`; *fields, includePermissionsForView, pageSize, pageToken, supportsAllDrives, useDomainAdminAccess* |  | `drive`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLEDOCS_LIST_SPREADSHEET_CHARTS` | Get Charts from Spreadsheet | read | `spreadsheet_id` |  | `drive`, `drive.readonly`, `drive.file`, `spreadsheets`, `spreadsheets.readonly` |
| `GOOGLEDOCS_REPLACE_ALL_TEXT` | Replace All Text in Document | write | `document_id`, `find_text`, `replace_text`; *match_case, search_by_regex, tab_ids* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_REPLACE_IMAGE` | Replace Image in Document | write | `document_id`, `replace_image` |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_SEARCH_DOCUMENTS` | Search Documents | read | *created_after, include_shared_drives, include_trashed, max_results, modified_after, order_by, page_token, query and 3 more* |  | `drive`, `drive.appdata`, `drive.file`, `drive.readonly`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly` |
| `GOOGLEDOCS_UNMERGE_TABLE_CELLS` | Unmerge Table Cells | write | `documentId`, `tableRange` |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_UPDATE_DOCUMENT_MARKDOWN` | Update Document Markdown | write | `id`, `markdown`; *tab_id* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_UPDATE_DOCUMENT_SECTION_MARKDOWN` | Update Document Section Markdown | write | `document_id`, `markdown_text`; *end_index, start_index, tab_id* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_UPDATE_DOCUMENT_STYLE` | Update Document Style | write | `document_id`, `document_style`; *fields, tab_id* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_UPDATE_EXISTING_DOCUMENT` | Update existing document | write | `document_id`, `edit_docs` |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_UPDATE_TABLE_ROW_STYLE` | Update Table Row Style | write | `documentId`, `updateTableRowStyle` |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_UPDATE_TAB_PROPERTIES` | Update Tab Properties | write | `document_id`, `tab_id`, `tab_properties`; *fields* |  | `documents`, `drive`, `drive.file` |
| `GOOGLEDOCS_UPLOAD_DRIVE_FILE` | Upload Drive File | write | `file`; *description, fields, mime_type, name, parent_folder_id, supports_all_drives* |  | `drive`, `drive.file` |
