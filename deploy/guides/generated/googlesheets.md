# Google Sheets: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

50 tools: 17 read, 33 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `GOOGLESHEETS_ADD_SHEET` | Add Sheet to Existing Spreadsheet | write | `spreadsheet_id`; *force_unique, properties, title* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_AGGREGATE_COLUMN_DATA` | Aggregate Column Data | read | `operation`, `sheet_name`, `spreadsheet_id`, `target_column`; *additional_filters, case_sensitive, has_header_row, percentage_total, search_column, search_value* |  | `spreadsheets`, `spreadsheets.readonly`, `drive`, `drive.file`, `drive.readonly` |
| `GOOGLESHEETS_APPEND_DIMENSION` | Append Dimension | write | `dimension`, `length`, `sheet_id`, `spreadsheet_id`; *include_spreadsheet_in_response, response_include_grid_data, response_ranges* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_AUTO_RESIZE_DIMENSIONS` | Auto-Resize Rows or Columns | write | `dimension`, `end_index`, `spreadsheet_id`, `start_index`; *sheet_id, sheet_name* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_BATCH_CLEAR_VALUES_BY_DATA_FILTER` | Batch Clear Values By Data Filter | write | `dataFilters`, `spreadsheetId` |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_BATCH_GET` | Batch get spreadsheet | read | `spreadsheet_id`; *dateTimeRenderOption, empty_strings_filtered, majorDimension, ranges, valueRenderOption* |  | `spreadsheets.readonly`, `spreadsheets`, `drive.readonly`, `drive.file`, `drive` |
| `GOOGLESHEETS_BATCH_UPDATE_VALUES_BY_DATA_FILTER` | Batch Update Values by Data Filter | write | `data`, `spreadsheetId`, `valueInputOption`; *includeValuesInResponse, responseDateTimeRenderOption, responseValueRenderOption* |  | `spreadsheets`, `drive`, `drive.file` |
| `GOOGLESHEETS_CLEAR_BASIC_FILTER` | Clear Basic Filter | write | `sheet_id`, `spreadsheet_id`; *include_spreadsheet_in_response, response_include_grid_data, response_ranges* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_CLEAR_VALUES` | Clear spreadsheet values | write | `range`, `spreadsheet_id` |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_CREATE_CHART` | Create Chart in Google Sheets | write | `chart_type`, `data_range`, `sheet_id`, `spreadsheet_id`; *background_blue, background_green, background_red, chart_spec, legend_position, subtitle, title, x_axis_title and 1 more* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_CREATE_GOOGLE_SHEET1` | Create a Google Sheet | write | *folder_id, folder_name, title* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_CREATE_SPREADSHEET_COLUMN` | Create spreadsheet column | write | `spreadsheet_id`; *include_spreadsheet_in_response, inherit_from_before, insert_index, response_include_grid_data, response_ranges, sheet_id, sheet_name* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_CREATE_SPREADSHEET_ROW` | Create spreadsheet row | write | `spreadsheet_id`; *include_spreadsheet_in_response, inherit_from_before, insert_index, response_include_grid_data, response_ranges, sheet_id, sheet_name* |  | `spreadsheets`, `drive`, `drive.file` |
| `GOOGLESHEETS_DELETE_CHART` | Delete Chart from Google Sheets | write | `chart_id`, `spreadsheet_id` |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_DELETE_DIMENSION` | Delete Dimension (Rows/Columns) | write | `spreadsheet_id`; *delete_dimension_request, dimension, end_index, include_spreadsheet_in_response, response_include_grid_data, response_ranges, sheet_id, sheet_name and 1 more* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_DELETE_SHEET` | Delete Sheet | write | `sheetId`, `spreadsheetId`; *includeSpreadsheetInResponse, responseIncludeGridData, responseRanges* |  | `spreadsheets`, `drive`, `drive.file` |
| `GOOGLESHEETS_EXPORT_DRIVE_FILE` | Export Drive Spreadsheet | read | `file_id`, `mime_type` |  | `drive`, `drive.readonly` |
| `GOOGLESHEETS_EXPORT_SHEET_AS_PDF` | Export Sheet as PDF | read | `sheet_id`, `spreadsheet_id`; *bottom_margin, end_column_index, end_row_index, fit_to_width, left_margin, portrait, print_titles, repeat_frozen_rows and 6 more* |  | `drive`, `drive.readonly`, `spreadsheets`, `spreadsheets.readonly` |
| `GOOGLESHEETS_FIND_REPLACE` | Find and Replace in Spreadsheet | write | `find`, `replace`, `spreadsheetId`; *allSheets, endColumnIndex, endRowIndex, includeFormulas, matchCase, matchEntireCell, range, rangeSheetId and 5 more* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_FORMAT_CELL` | Format cell | write | `spreadsheet_id`; *background_alpha, background_color, blue, bold, border_color, border_style, end_column_index, end_row_index and 27 more* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_GET_CONDITIONAL_FORMAT_RULES` | Get conditional format rules | read | `spreadsheet_id`; *exclude_tables_in_banded_ranges, sheet_id, sheet_title* |  | `spreadsheets.readonly`, `spreadsheets`, `drive.readonly`, `drive.file`, `drive` |
| `GOOGLESHEETS_GET_DATA_VALIDATION_RULES` | Get Data Validation Rules | read | `spreadsheetId`; *includeEmpty, ranges, sheetId, sheetTitle* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_GET_DRIVE_FILE_METADATA` | Get Drive File Metadata | read | `fileId`; *fields, includeLabels, includePermissionsForView, supportsAllDrives* |  | `drive`, `drive.metadata`, `drive.metadata.readonly`, `drive.readonly` |
| `GOOGLESHEETS_GET_SHEET_NAMES` | Get sheet names | read | `spreadsheet_id`; *exclude_hidden* |  | `spreadsheets.readonly`, `spreadsheets`, `drive.readonly`, `drive.file`, `drive` |
| `GOOGLESHEETS_GET_SPREADSHEET_BY_DATA_FILTER` | Get Spreadsheet by Data Filter | read | `spreadsheetId`; *dataFilters, excludeTablesInBandedRanges, includeGridData* |  | `spreadsheets`, `drive`, `drive.file` |
| `GOOGLESHEETS_GET_SPREADSHEET_INFO` | Get spreadsheet info | read | *exclude_tables_in_banded_ranges, fields, include_grid_data, ranges, spreadsheet_id* |  | `spreadsheets.readonly`, `spreadsheets`, `drive.readonly`, `drive.file`, `drive` |
| `GOOGLESHEETS_INSERT_DIMENSION` | Insert Dimension in Google Sheet | write | `insert_dimension`, `spreadsheet_id`; *include_spreadsheet_in_response, response_include_grid_data, response_ranges* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_LIST_CHARTS` | List Charts in Google Sheets | read | `spreadsheet_id`; *chart_id, include_spec, sheet_id* |  | `spreadsheets.readonly`, `spreadsheets`, `drive.readonly`, `drive.file`, `drive` |
| `GOOGLESHEETS_LIST_DRIVE_CHANGES` | List Drive Changes | read | *drive_id, include_corpus_removals, include_items_from_all_drives, include_labels, include_permissions_for_view, include_removed, page_size, page_token and 3 more* |  | `drive`, `drive.metadata`, `drive.metadata.readonly`, `drive.readonly` |
| `GOOGLESHEETS_LOOKUP_SPREADSHEET_ROW` | Look up spreadsheet row | read | `query`, `spreadsheet_id`; *case_sensitive, date_time_render_option, normalize_whitespace, range, value_render_option* |  | `spreadsheets.readonly`, `spreadsheets`, `drive.readonly`, `drive.file`, `drive` |
| `GOOGLESHEETS_MOVE_CHART` | Move or Resize Chart in Google Sheets | write | `chart_id`, `spreadsheet_id`; *anchor_column_index, anchor_row_index, anchor_sheet_id, height_pixels, new_sheet, offset_x_pixels, offset_y_pixels, target_sheet_id and 1 more* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_MOVE_SPREADSHEET_TO_FOLDER` | Move Spreadsheet to Folder | write | `destination_folder_id`, `spreadsheet_id` |  | `drive`, `drive.file`, `drive.metadata` |
| `GOOGLESHEETS_MUTATE_CONDITIONAL_FORMAT_RULES` | Mutate conditional format rules | write | `operation`, `sheet_id`, `spreadsheet_id`; *index, new_index, rule* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_SEARCH_DEVELOPER_METADATA` | Search Developer Metadata | read | `dataFilters`, `spreadsheetId` |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_SEARCH_SPREADSHEETS` | Search Spreadsheets | read | *created_after, include_shared_drives, include_trashed, max_results, modified_after, order_by, page_token, query and 3 more* |  | `drive`, `drive.appdata`, `drive.file`, `drive.meet.readonly`, `drive.metadata`, `drive.metadata.readonly`, `drive.photos.readonly`, `drive.readonly` |
| `GOOGLESHEETS_SET_BASIC_FILTER` | Set Basic Filter | write | `filter`, `spreadsheetId`; *includeSpreadsheetInResponse, responseIncludeGridData, responseRanges* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_SET_DATA_VALIDATION_RULE` | Set Data Validation Rule | write | `end_column_index`, `end_row_index`, `mode`, `sheet_id`, `spreadsheet_id`, `start_column_index`, `start_row_index`; *condition_values, filtered_rows_included, formula, input_message, show_custom_ui, source_range_a1, strict, validation_type and 1 more* |  | `spreadsheets`, `drive`, `drive.file` |
| `GOOGLESHEETS_SORT_RANGE` | Sort Range | write | `range`, `sortSpecs`, `spreadsheet_id`; *includeSpreadsheetInResponse, responseIncludeGridData, responseRanges* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_SPREADSHEETS_SHEETS_COPY_TO` | Copy Sheet to Another Spreadsheet | write | `destination_spreadsheet_id`, `sheet_id`, `spreadsheet_id` |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_SPREADSHEETS_VALUES_APPEND` | Append Values to Spreadsheet | write | `range`, `spreadsheetId`, `valueInputOption`, `values`; *includeValuesInResponse, insertDataOption, majorDimension, responseDateTimeRenderOption, responseValueRenderOption* |  | `spreadsheets`, `drive`, `drive.file` |
| `GOOGLESHEETS_SPREADSHEETS_VALUES_BATCH_CLEAR` | Batch Clear Spreadsheet Values | write | `ranges`, `spreadsheet_id` |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_SPREADSHEETS_VALUES_BATCH_GET_BY_DATA_FILTER` | Batch Get Spreadsheet Values by Data Filter | read | `dataFilters`, `spreadsheetId`; *dateTimeRenderOption, majorDimension, valueRenderOption* |  | `drive`, `drive.file`, `spreadsheets` |
| `GOOGLESHEETS_UPDATE_CHART` | Update Chart in Google Sheets | write | `chart_id`, `chart_spec`, `spreadsheet_id` |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_UPDATE_DIMENSION_PROPERTIES` | Update Dimension Properties (Hide/Unhide & Resize) | write | `dimension`, `end_index`, `spreadsheet_id`, `start_index`; *hidden_by_user, include_spreadsheet_in_response, pixel_size, response_include_grid_data, response_ranges, sheet_id, sheet_name* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_UPDATE_SHEET_PROPERTIES` | Update Sheet Properties | write | `spreadsheetId`, `updateSheetProperties`; *includeSpreadsheetInResponse, responseIncludeGridData, responseRanges* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_UPDATE_SPREADSHEET_PROPERTIES` | Update Spreadsheet Properties | write | `fields`, `properties`, `spreadsheetId`; *includeSpreadsheetInResponse, responseIncludeGridData, responseRanges* |  | `spreadsheets`, `drive.file`, `drive` |
| `GOOGLESHEETS_UPDATE_VALUES_BATCH` | Batch update spreadsheet values | write | `data`, `spreadsheet_id`, `valueInputOption`; *includeValuesInResponse, responseDateTimeRenderOption, responseValueRenderOption* |  | `drive`, `drive.file`, `spreadsheets` |
| `GOOGLESHEETS_UPSERT_ROWS` | Upsert Rows (Smart Update/Insert) | write | `rows`, `sheetName`, `spreadsheetId`; *headers, keyColumn, key_column_index, normalization_message, strictMode, tableStart* |  | `spreadsheets`, `drive`, `drive.file` |
| `GOOGLESHEETS_VALUES_GET` | Get spreadsheet values | read | `range`, `spreadsheet_id`; *date_time_render_option, end_row, major_dimension, start_row, value_render_option* |  | `drive`, `drive.readonly`, `drive.file`, `spreadsheets`, `spreadsheets.readonly` |
| `GOOGLESHEETS_VALUES_UPDATE` | Update spreadsheet values | write | `range`, `spreadsheet_id`, `value_input_option`, `values`; *auto_expand_sheet, include_values_in_response, major_dimension, response_datetime_render_option, response_value_render_option* |  | `drive`, `drive.file`, `spreadsheets` |
