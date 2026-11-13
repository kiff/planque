# Airtable: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

26 tools: 9 read, 17 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `AIRTABLE_CREATE_BASE` | Create base | write | `name`, `tables`, `workspaceId` |  | `schema.bases:write` |
| `AIRTABLE_CREATE_COMMENT` | Create Comment | write | `baseId`, `recordId`, `tableIdOrName`, `text`; *parentCommentId* |  | `data.recordComments:write` |
| `AIRTABLE_CREATE_FIELD` | Create Field | write | `baseId`, `name`, `tableIdOrName`; *description, options, type* |  | `schema.bases:read`, `schema.bases:write` |
| `AIRTABLE_CREATE_RECORDS` | Create records | write | `baseId`, `records`, `tableIdOrName`; *returnFieldsByFieldId, typecast* |  | `data.records:write` |
| `AIRTABLE_CREATE_RECORD_FROM_NATURAL_LANGUAGE` | Create Record From Natural Language | write | `baseId`, `nl_query`, `tableIdOrName`; *typecast* |  | `data.records:write`, `schema.bases:read` |
| `AIRTABLE_CREATE_TABLE` | Create table | write | `base_id`, `fields`, `name`; *description, if_exists, reorder_message* |  | `schema.bases:read`, `schema.bases:write` |
| `AIRTABLE_DELETE_COMMENT` | Delete Comment | write | `baseId`, `recordId`, `rowCommentId`, `tableIdOrName` |  | `data.recordComments:write` |
| `AIRTABLE_DELETE_MULTIPLE_RECORDS` | Delete multiple records | write | `baseId`, `recordIds`, `tableIdOrName` |  | `data.records:write` |
| `AIRTABLE_DELETE_RECORD` | Delete Record | write | `baseId`, `recordId`, `tableIdOrName` |  | `data.records:write` |
| `AIRTABLE_GET_ATTACHMENT_INFO` | Get Attachment Info | read | `attachmentFieldIdOrName`, `baseId`, `recordId`, `tableIdOrName` |  | `data.records:read` |
| `AIRTABLE_GET_BASE_METADATA` | Get Base Metadata | read | `baseId` |  | `schema.bases:read` |
| `AIRTABLE_GET_BASE_SCHEMA` | Get Base Schema | read | `baseId`; *include* |  | `schema.bases:read` |
| `AIRTABLE_GET_RECORD` | Get Record | read | `baseId`, `recordId`, `tableIdOrName`; *cellFormat, returnFieldsByFieldId, timeZone, userLocale* |  | `data.records:read` |
| `AIRTABLE_GET_USER_INFO` | Get user information | read |  |  |  |
| `AIRTABLE_LIST_BASES` | List bases | read | *offset* |  | `schema.bases:read` |
| `AIRTABLE_LIST_COMMENTS` | List Comments | read | `baseId`, `recordId`, `tableIdOrName`; *offset, pageSize* |  | `data.recordComments:read` |
| `AIRTABLE_LIST_RECORDS` | List records | read | `baseId`, `tableIdOrName`; *cellFormat, fields, filterByFormula, maxRecords, offset, pageSize, recordMetadata, returnFieldsByFieldId and 4 more* |  | `data.records:read` |
| `AIRTABLE_UPDATE_COMMENT` | Update Comment | write | `baseId`, `recordId`, `rowCommentId`, `tableIdOrName`, `text` |  | `data.recordComments:write` |
| `AIRTABLE_UPDATE_FIELD` | Update Field | write | `baseId`, `fieldId`, `tableId`; *description, name* |  | `schema.bases:write` |
| `AIRTABLE_UPDATE_MULTIPLE_RECORDS` | Update multiple records | write | `baseId`, `records`, `tableIdOrName`; *performUpsert, returnFieldsByFieldId, typecast* |  | `data.records:write` |
| `AIRTABLE_UPDATE_MULTIPLE_RECORDS_PUT` | Update multiple records (PUT) | write | `baseId`, `records`, `tableIdOrName`; *performUpsert, returnFieldsByFieldId, typecast* |  | `data.records:write` |
| `AIRTABLE_UPDATE_RECORD` | Update record | write | `baseId`, `fields`, `recordId`, `tableIdOrName`; *returnFieldsByFieldId, typecast* |  | `data.records:write` |
| `AIRTABLE_UPDATE_RECORD_PUT` | Update record (PUT) | write | `baseId`, `fields`, `recordId`, `tableIdOrName`; *returnFieldsByFieldId, typecast* |  | `data.records:write` |
| `AIRTABLE_UPDATE_TABLE` | Update Table | write | `baseId`, `tableIdOrName`; *dateDependencySettings, description, name* |  | `schema.bases:write` |
| `AIRTABLE_UPLOAD_ATTACHMENT` | Upload attachment | write | `attachmentFieldIdOrName`, `baseId`, `contentType`, `file`, `filename`, `recordId` |  | `data.records:write` |
| `AIRTABLE_WHO_AM_I` | Who Am I | read |  |  |  |
