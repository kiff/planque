# Excel: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

125 tools: 63 read, 62 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `EXCEL_ADD_CHART` | Add Chart | write | `item_id`, `seriesby`, `sourcedata`, `type`, `worksheet`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_ADD_NAMED_ITEM` | Add Named Item | write | `item_id`, `name`; *comment, drive_id, formula, formula_mode, reference, session_id, worksheet* |  | `Files.ReadWrite` |
| `EXCEL_ADD_SHAREPOINT_WORKSHEET` | Add SharePoint Worksheet | write | `item_id`, `name`, `site_id`; *session_id* |  | `Files.ReadWrite` |
| `EXCEL_ADD_TABLE` | Add Table | write | `address`, `hasHeaders`, `item_id`; *drive_id, session_id, worksheet* |  | `Files.ReadWrite` |
| `EXCEL_ADD_TABLE_COLUMN` | Add Table Column | write | `index`, `item_id`, `table_id`, `values`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_ADD_TABLE_ROW` | Add Table Row | write | `item_id`, `table_id`, `values`; *drive_id, index, session_id* |  | `Files.ReadWrite` |
| `EXCEL_ADD_WORKBOOK_PERMISSION` | Add Workbook Permission | write | `item_id`, `recipients`, `roles`; *drive_id, expirationDateTime, message, password, requireSignIn, retainInheritedPermissions, sendInvitation, site_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `EXCEL_ADD_WORKSHEET` | Add Worksheet | write | `item_id`, `name`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_APPLY_TABLE_FILTER` | Apply Table Filter | write | `column_id`, `criteria`, `item_id`, `table_id`, `worksheet`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_APPLY_TABLE_SORT` | Apply Table Sort | write | `fields`, `item_id`, `table_id`, `worksheet`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_AUTOFIT_RANGE` | Autofit Range | write | `dimension`, `item_id`, `target`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_CALCULATE_FUNCTION` | Calculate Workbook Function | read | `function_name`, `item_id`; *arguments, drive_id, session_id* |  | `Files.Read` |
| `EXCEL_CALCULATE_WORKBOOK` | Calculate Excel Workbook | write | `calculation_type`, `item_id`; *drive_id* |  | `Files.ReadWrite` |
| `EXCEL_CLEAR_RANGE` | Clear Range | write | `address`, `item_id`, `worksheet_id`; *applyTo, drive_id, session_id, site_id* |  | `Files.ReadWrite` |
| `EXCEL_CLEAR_TABLE_FILTER` | Clear Table Filter | write | `column_id`, `item_id`, `table_id`, `worksheet`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_CLEAR_TABLE_FILTERS` | Clear All Table Filters | write | `item_id`, `table_id`; *drive_id, session_id, worksheet* |  | `Files.ReadWrite` |
| `EXCEL_CLEAR_TABLE_SORT` | Clear Table Sort | write | `item_id`, `table_id`; *drive_id, session_id, worksheet* |  | `Files.ReadWrite` |
| `EXCEL_CLOSE_SESSION` | Close Excel Session | write | `item_id`, `session_id`; *drive_id* |  | `Files.ReadWrite` |
| `EXCEL_CONVERT_TABLE_TO_RANGE` | Convert Table To Range | write | `item_id`, `table_id`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_COPY_DRIVE_ITEM` | Copy Drive Item | write | `drive_id`, `item_id`; *children_only, conflict_behavior, destination_drive_id, destination_folder_id, include_all_version_history, name* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `EXCEL_CREATE_COMMENT` | Create Comment | write | `cell_address`, `content`, `item_id`; *drive_id* |  | `Files.ReadWrite` |
| `EXCEL_CREATE_FOLDER` | Create Folder | write | `drive_id`, `name`, `parent_item_id`; *conflict_behavior* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `EXCEL_CREATE_UPLOAD_SESSION` | Create Upload Session | write | `drive_id`, `file_path`; *conflict_behavior, defer_commit* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `EXCEL_CREATE_WORKBOOK` | Create Workbook | write | `path`; *drive_id, worksheet_data, worksheet_names* |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `EXCEL_CREATE_WORKBOOK_SHARING_LINK` | Create Workbook Sharing Link | write | `item_id`, `type`; *drive_id, expirationDateTime, retainInheritedPermissions, scope* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `EXCEL_DELETE_DRIVE_ITEM` | Delete Drive Item | write | `item_id`; *drive_id, if_match, prefer* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `EXCEL_DELETE_NAMED_ITEM` | Delete Named Item | write | `item_id`, `name`; *drive_id, session_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All` |
| `EXCEL_DELETE_RANGE` | Delete Range | write | `address`, `item_id`, `shift`, `worksheet_id`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_DELETE_TABLE` | Delete Table | write | `item_id`, `table_id`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_DELETE_TABLE_COLUMN` | Delete Table Column | write | `column_id`, `item_id`, `table_id`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_DELETE_TABLE_ROW` | Delete Table Row | write | `item_id`, `row_index`, `table_id`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_DELETE_WORKSHEET` | Delete Worksheet | write | `item_id`, `worksheet`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_DOWNLOAD_WORKBOOK` | Download Workbook | read | `item_id`; *drive_id, if_none_match* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `EXCEL_EXPORT_WORKBOOK_TO_PDF` | Export Workbook to PDF | read | `item_id`; *drive_id, drive_scope, return_download_url_only, save_pdf_to_path* |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `EXCEL_GET_APPLICATION` | Get Excel Application | read | `item_id`; *drive_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_CHART` | Get chart | read | `item_id`, `worksheet`; *chart_index, chart_name, drive_id, expand, select_fields, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_CHART_AXIS` | Get Chart Axis | read | `axis_type`, `chart_name`, `item_id`, `worksheet`; *drive_id, expand, select, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_CHART_COMPONENT` | Get chart component | read | `chart_name`, `item_id`, `selector`, `worksheet`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_CHART_DATA_LABELS` | Get Chart Data Labels | read | `chart_id`, `item_id`, `worksheet_id`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_CHART_IMAGE` | Get chart image | read | `chart_name`, `item_id`, `worksheet`; *drive_id, fitting_mode, height, session_id, width* |  | `Files.ReadWrite` |
| `EXCEL_GET_CHART_LEGEND` | Get Chart Legend | read | `chart_name`, `item_id`, `worksheet_id`; *drive_id, expand, select, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_CHART_POINT` | Get chart point | read | `chart_name`, `item_id`, `worksheet`; *drive_id, expand, point_id, point_index, select_fields, series_id, series_index, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_CHART_SERIES` | Get chart series | read | `chart_name`, `item_id`, `worksheet`; *drive_id, expand, select_fields, series_id, series_index, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_COMMENT` | Get Comment | read | `comment_id`, `item_id`; *drive_id, expand, select, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_COMMENT_REPLY` | Get Comment Reply | read | `comment_id`, `item_id`, `reply_id`; *drive_id, expand, select, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_DERIVED_RANGE` | Get Derived Range | read | `item_id`, `operation`, `target`; *another_range, column_index, column_offset, count, delta_columns, delta_rows, drive_id, row_index and 3 more* |  | `Files.ReadWrite` |
| `EXCEL_GET_DRIVE_ITEM` | Get Drive Item | read | *drive_id, expand, if_none_match, include_deleted_items, item_id, path, select* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `EXCEL_GET_NAMED_ITEM` | Get Named Item | read | `item_id`, `name`; *drive_id, select_fields, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_NAMED_ITEM_RANGE` | Get Named Item Range | read | `item_id`, `name`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_OPERATION_STATUS` | Get Excel Operation Status | read | `item_id`, `operation_id`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_PIVOT_TABLE` | Get Pivot Table | read | `item_id`, `pivot_table_id`, `worksheet`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_PROJECTED_WORKSHEET_SUBRESOURCE` | Get Projected Worksheet Subresource | read | `item_id`, `subresource`, `worksheet_id`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_RANGE` | Get Range | read | `address`, `item_id`, `worksheet_id`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_RANGE_FORMAT` | Get Range Format | read | `component`, `item_id`, `target`; *drive_id, session_id, side_index* |  | `Files.ReadWrite` |
| `EXCEL_GET_SESSION` | Create Excel Session | write | `item_id`; *drive_id, persist_changes* |  | `Files.ReadWrite` |
| `EXCEL_GET_SHAREPOINT_RANGE` | Get SharePoint Range | read | `address`, `item_id`, `site_id`, `worksheet_id`; *session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_SHAREPOINT_WORKSHEET` | Get SharePoint Worksheet | read | `item_id`, `site_id`, `worksheet`; *session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_TABLE` | Get Table | read | `item_id`, `table_id`; *drive_id, expand, select_fields, session_id, worksheet* |  | `Files.ReadWrite` |
| `EXCEL_GET_TABLE_COLUMN` | Get table column | read | `column_id`, `item_id`, `table_id`; *drive_id, session_id, site_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_TABLE_COLUMN_CELL` | Get Table Column Cell | read | `column_id`, `item_id`, `row`, `table_id`; *column, drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_TABLE_COLUMN_RANGE` | Get Table Column Range | read | `column_id`, `item_id`, `table_id`; *drive_id, range_kind, session_id, worksheet* |  | `Files.ReadWrite` |
| `EXCEL_GET_TABLE_RANGE` | Get Table Range | read | `item_id`, `table_id`; *drive_id, range_kind, session_id, worksheet* |  | `Files.ReadWrite` |
| `EXCEL_GET_TABLE_ROW` | Get Table Row | read | `item_id`, `row_index`, `table_id`; *drive_id, session_id, worksheet* |  | `Files.ReadWrite` |
| `EXCEL_GET_TABLE_ROW_RANGE` | Get Table Row Range | read | `item_id`, `row_index`, `table_id`; *drive_id, session_id, worksheet* |  | `Files.ReadWrite` |
| `EXCEL_GET_TABLE_SORT` | Get Table Sort | read | `item_id`, `table_id`; *drive_id, session_id, worksheet* |  | `Files.ReadWrite` |
| `EXCEL_GET_VISIBLE_RANGE` | Get Visible Range | read | `address`, `item_id`, `result`, `worksheet_id`; *drive_id, row_index, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_WORKBOOK` | Get workbook | read | `drive_id`, `item_id`; *expand, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_WORKSHEET` | Get Worksheet | read | `item_id`, `worksheet`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_WORKSHEET_PROTECTION` | Get Worksheet Protection | read | `item_id`, `worksheet`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_WORKSHEET_ROW_RANGE` | Get Worksheet Row Range | read | `item_id`, `row_index`, `worksheet`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_GET_WORKSHEET_USED_RANGE` | Get Worksheet Used Range | read | `drive_id`, `item_id`, `worksheet_id`; *session_id, values_only* |  | `Files.ReadWrite` |
| `EXCEL_INSERT_RANGE` | Insert Range | write | `address`, `item_id`, `shift`, `worksheet_id`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_LIST_CHARTS` | List Charts | read | `item_id`, `worksheet`; *drive_id, session_id, skip, top* |  | `Files.ReadWrite` |
| `EXCEL_LIST_CHART_POINTS` | List chart points | read | `chart_name`, `item_id`, `worksheet`; *drive_id, expand, select_fields, series_id, series_index, session_id, skip, top* |  | `Files.ReadWrite` |
| `EXCEL_LIST_CHART_SERIES` | List Chart Series | read | `chart_name`, `item_id`, `worksheet`; *drive_id, session_id, skip, top* |  | `Files.ReadWrite` |
| `EXCEL_LIST_COMMENTS` | List Comments | read | `item_id`; *drive_id, session_id, skip, top* |  | `Files.ReadWrite` |
| `EXCEL_LIST_COMMENT_REPLIES` | List Comment Replies | read | `comment_id`, `item_id`; *count, drive_id, expand, filter, orderby, search, select, session_id and 2 more* |  | `Files.ReadWrite` |
| `EXCEL_LIST_DRIVE_ITEM_CHILDREN` | List Drive Item Children | read | `drive_id`, `item_id`; *skip_token, top* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `EXCEL_LIST_DRIVE_ITEM_VERSIONS` | List Drive Item Versions | read | `item_id`; *drive_id, top* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `EXCEL_LIST_FILES` | List Drive Files | read | *drive_id, path, skip_token, top* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All` |
| `EXCEL_LIST_NAMED_ITEMS` | List Named Items | read | `item_id`; *drive_id, session_id, skip, top* |  | `Files.ReadWrite` |
| `EXCEL_LIST_PIVOT_TABLES` | List Pivot Tables | read | `item_id`, `worksheet`; *drive_id, session_id, skip, top* |  | `Files.ReadWrite` |
| `EXCEL_LIST_SHAREPOINT_TABLES` | List SharePoint Tables | read | `item_id`, `site_id`, `worksheet`; *session_id, skip, top* |  | `Files.ReadWrite` |
| `EXCEL_LIST_SHAREPOINT_WORKSHEETS` | List SharePoint Worksheets | read | `item_id`, `site_id`; *drive_id, session_id, skip, top* |  | `Files.ReadWrite` |
| `EXCEL_LIST_SITE_LISTS` | List Site Lists | read | `site_id`; *select_fields, skip_token, top* |  | `Sites.Read.All` |
| `EXCEL_LIST_TABLES` | List Tables | read | `item_id`, `worksheet`; *drive_id, session_id, skip, top* |  | `Files.ReadWrite` |
| `EXCEL_LIST_TABLE_COLUMNS` | List Table Columns | read | `item_id`, `table_id`; *drive_id, session_id, skip, top* |  | `Files.ReadWrite` |
| `EXCEL_LIST_TABLE_ROWS` | List Table Rows | read | `item_id`, `table_id`; *drive_id, session_id, skip, top* |  | `Files.ReadWrite` |
| `EXCEL_LIST_WORKBOOK_PERMISSIONS` | List Workbook Permissions | read | `drive_id`, `item_id`; *expand, select, skip, skip_token, top* |  | `Files.Read`, `Files.Read.All`, `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `EXCEL_LIST_WORKBOOK_TABLES` | List Workbook Tables | read | `item_id`; *count, drive_id, expand, filter, order_by, search, select_fields, session_id and 2 more* |  | `Files.Read` |
| `EXCEL_LIST_WORKSHEETS` | List Worksheets | read | `item_id`; *drive_id, session_id, skip, top* |  | `Files.ReadWrite` |
| `EXCEL_MERGE_CELLS` | Merge Cells | write | `item_id`, `range`, `worksheet`; *across, drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_PROTECT_WORKSHEET` | Protect Worksheet | write | `item_id`, `worksheet`; *drive_id, options, session_id* |  | `Files.ReadWrite` |
| `EXCEL_REAPPLY_TABLE_FILTERS` | Reapply All Table Filters | write | `item_id`, `table_id`; *drive_id, session_id, worksheet* |  | `Files.ReadWrite` |
| `EXCEL_REAPPLY_TABLE_SORT` | Reapply Table Sort | write | `item_id`, `table_id`; *drive_id, session_id, worksheet* |  | `Files.ReadWrite` |
| `EXCEL_REFRESH_PIVOT_TABLES` | Refresh Pivot Tables | write | `item_id`, `worksheet`; *drive_id, pivot_table_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_REFRESH_SESSION` | Refresh Excel Session | write | `item_id`, `session_id`; *drive_id* |  | `Files.ReadWrite` |
| `EXCEL_RESOLVE_SHARED_WORKBOOK` | Resolve Shared Workbook | read | `share`; *select* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `EXCEL_REVOKE_WORKBOOK_PERMISSION` | Revoke Workbook Permission | write | `item_id`, `permission_id`; *drive_id, if_match* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `EXCEL_RUN_GRAPH_MUTATION_BATCH` | Run Graph Mutation Batch | write | `requests` |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `EXCEL_RUN_GRAPH_READ_BATCH` | Run Graph Read Batch | write | `requests` |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `EXCEL_SEARCH_FILES` | Search Drive Files | read | `query`; *drive_id, file_extensions, scope, skip_token, top* |  | `Files.Read`, `Files.ReadWrite`, `Files.Read.All`, `Files.ReadWrite.All`, `Sites.Read.All`, `Sites.ReadWrite.All` |
| `EXCEL_SEARCH_MICROSOFT_365_CONTENT` | Search Microsoft 365 Content | read | `query`; *entity_types, fields, from_index, region, size* |  | `Sites.Read.All`, `Sites.ReadWrite.All` |
| `EXCEL_SEARCH_SITES` | Search Sites | read | `query`; *select_fields, skip_token, top* |  | `Sites.Read.All` |
| `EXCEL_SET_CHART_DATA` | Set chart data | write | `chart_name`, `item_id`, `source_data`, `worksheet`; *drive_id, series_by, session_id* |  | `Files.ReadWrite` |
| `EXCEL_SET_CHART_POSITION` | Set chart position | write | `chart_name`, `item_id`, `start_cell`, `worksheet`; *drive_id, end_cell, session_id* |  | `Files.ReadWrite` |
| `EXCEL_SORT_RANGE` | Sort Range | write | `fields`, `item_id`, `range_type`, `worksheet`; *drive_id, session_id, site_id* |  | `Files.ReadWrite` |
| `EXCEL_UNMERGE_RANGE` | Unmerge Range | write | `item_id`, `range`, `worksheet`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_UNPROTECT_WORKSHEET` | Unprotect Worksheet | write | `item_id`, `worksheet`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_UPDATE_CHART` | Update Chart | write | `chart_id`, `item_id`, `worksheet`; *drive_id, height, left, name, session_id, top, width* |  | `Files.ReadWrite` |
| `EXCEL_UPDATE_CHART_COMPONENT` | Update chart component | write | `chart_name`, `item_id`, `settings`, `worksheet`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_UPDATE_CHART_LEGEND` | Update Chart Legend | write | `chart_id`, `item_id`, `worksheet`; *drive_id, overlay, position, session_id, visible* |  | `Files.ReadWrite` |
| `EXCEL_UPDATE_CHART_SERIES` | Update chart series | write | `chart_name`, `item_id`, `name`, `worksheet`; *drive_id, series_id, series_index, session_id* |  | `Files.ReadWrite` |
| `EXCEL_UPDATE_DRIVE_ITEM` | Update Drive Item | write | `item_id`; *drive_id, file_system_info, if_match, name, parent_item_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `EXCEL_UPDATE_NAMED_ITEM` | Update Named Item | write | `item_id`, `name`; *comment, drive_id, session_id, visible* |  | `Files.ReadWrite` |
| `EXCEL_UPDATE_RANGE` | Update Range | write | `address`, `item_id`, `values`, `worksheet_id`; *drive_id, session_id* |  | `Files.ReadWrite` |
| `EXCEL_UPDATE_RANGE_FORMAT` | Update Range Format | write | `component`, `item_id`, `target`; *border, drive_id, fill, font, format, protection, session_id* |  | `Files.ReadWrite` |
| `EXCEL_UPDATE_SHAREPOINT_RANGE` | Update SharePoint Range | write | `address`, `item_id`, `site_id`, `values`, `worksheet_id`; *session_id* |  | `Files.ReadWrite` |
| `EXCEL_UPDATE_TABLE` | Update Table | write | `item_id`, `table_id`; *drive_id, name, session_id, showHeaders, showTotals, style* |  | `Files.ReadWrite` |
| `EXCEL_UPDATE_TABLE_COLUMN` | Update Table Column | write | `column_id`, `item_id`, `table_id`; *drive_id, name, session_id, values, worksheet* |  | `Files.ReadWrite` |
| `EXCEL_UPDATE_TABLE_COLUMN_RANGE` | Update Table Column Range | write | `column_index`, `item_id`, `table_id`, `values`; *drive_id, session_id, worksheet* |  | `Files.ReadWrite` |
| `EXCEL_UPDATE_TABLE_ROW` | Update Table Row | write | `item_id`, `row_index`, `table_id`, `values`; *drive_id, session_id, worksheet* |  | `Files.ReadWrite` |
| `EXCEL_UPDATE_WORKSHEET` | Update Worksheet | write | `item_id`, `worksheet`; *drive_id, name, position, session_id, site_id* |  | `Files.ReadWrite` |
| `EXCEL_UPLOAD_WORKBOOK` | Upload Workbook from URL | write | `source_url`; *conflict_behavior, destination_path, drive_id* |  | `Files.ReadWrite`, `Files.ReadWrite.All`, `Sites.ReadWrite.All` |
| `EXCEL_WHO_AM_I` | Who Am I | read |  |  |  |
