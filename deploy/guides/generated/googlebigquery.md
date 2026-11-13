# Google BigQuery: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

63 tools: 32 read, 31 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `GOOGLEBIGQUERY_CANCEL_JOB` | Cancel BigQuery Job | write | `job_id`, `project_id`; *location* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_CREATE_CAPACITY_COMMITMENT` | Create Capacity Commitment | write | `parent`, `plan`, `slotCount`; *capacityCommitmentId, edition, enforceSingleAdminProjectPerOrg, multiRegionAuxiliary, renewalPlan* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_CREATE_CONNECTION` | Create BigQuery Connection | write | `parent`; *aws, azure, cloudResource, cloudSpanner, cloudSql, connectionId, description, friendlyName and 3 more* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_CREATE_DATAEXCHANGES_LISTINGS` | Create Analytics Hub Listing | write | `bigqueryDataset`, `displayName`, `listingId`, `parent`; *categories, dataProvider, description, documentation, primaryContact, publisher, requestAccess* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_CREATE_DATASET` | Create BigQuery Dataset | write | `dataset_id`, `location`, `project_id`; *access, defaultPartitionExpirationMs, defaultTableExpirationMs, description, friendly_name, labels* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_CREATE_DATA_EXCHANGE` | Create Analytics Hub Data Exchange | write | `dataExchangeId`, `displayName`, `parent`; *description, documentation, icon, primaryContact, sharingEnvironmentConfig* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_CREATE_LISTING` | Create Analytics Hub Listing | write | `bigqueryDataset`, `displayName`, `parent`; *access_token, alt, callback, categories, dataProvider, description, documentation, fields and 13 more* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_CREATE_LOCATIONS_DATAPOLICIES` | Create BigQuery Data Policy (v2beta1) | write | `dataPolicy`, `dataPolicyId`, `parent` |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_CREATE_QUERY_TEMPLATE` | Create Analytics Hub Query Template | write | `displayName`, `parent`, `queryTemplateId`, `routine`; *description* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_CREATE_RESERVATION` | Create BigQuery Reservation | write | `parent`, `slotCapacity`; *autoscale, concurrency, edition, ignoreIdleSlots, multiRegionAuxiliary, reservationId* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_CREATE_RESERVATION_ASSIGNMENT` | Create BigQuery Reservation Assignment | write | `assignee`, `jobType`, `parent`; *assignmentId* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_CREATE_ROUTINE` | Create BigQuery Routine | write | `dataset_id`, `definitionBody`, `project_id`, `routineReference`, `routineType`; *arguments, dataGovernanceType, description, determinismLevel, importedLibraries, language, remoteFunctionOptions, returnTableType and 4 more* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_CREATE_TABLE` | Create BigQuery Table | write | `dataset_id`, `project_id`, `tableReference`; *clustering, defaultCollation, description, encryptionConfiguration, expirationTime, externalDataConfiguration, friendlyName, labels and 6 more* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_DELETE_DATASET` | Delete BigQuery Dataset | write | `dataset_id`, `project_id`; *delete_contents* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_DELETE_JOB_METADATA` | Delete BigQuery Job Metadata | write | `job_id`, `project_id`; *location* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_DELETE_MODEL` | Delete BigQuery ML Model | write | `dataset_id`, `model_id`, `project_id` |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_DELETE_ROUTINE` | Delete BigQuery Routine | write | `dataset_id`, `project_id`, `routine_id` |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_DELETE_TABLE` | Delete BigQuery Table | write | `dataset_id`, `project_id`, `table_id` |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_GET_BIGQUERY_MODEL` | Get BigQuery ML Model | read | `dataset_id`, `model_id`, `project_id` |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_GET_CONNECTION_IAM_POLICY` | Get BigQuery Connection IAM Policy | read | `resource`; *options* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_GET_DATASET` | Get BigQuery Dataset Metadata | read | `dataset_id`, `project_id`; *access_policy_version, dataset_view* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_GET_JOB` | Get BigQuery Job | read | `job_id`, `project_id`; *location* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_GET_QUERY_RESULTS` | Get BigQuery Query Results | read | `jobId`, `projectId`; *access_token, alt, callback, fields, formatOptions.useInt64Timestamp, key, location, maxResults and 9 more* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_GET_ROUTINE` | Get BigQuery Routine | read | `dataset_id`, `project_id`, `routine_id`; *read_mask* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_GET_ROUTINE_IAM_POLICY` | Get BigQuery Routine IAM Policy | read | `dataset_id`, `project_id`, `routine_id`; *options* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_GET_SERVICE_ACCOUNT` | Get BigQuery Service Account | read | `project_id`; *access_token, alt, callback, fields, key, oauth_token, prettyPrint, quota_user and 3 more* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_GET_TABLE_IAM_POLICY` | Get BigQuery Table IAM Policy | read | `resource`; *requested_policy_version* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_GET_TABLE_SCHEMA` | Get BigQuery Table Schema | read | `dataset_id`, `project_id`, `table_id`; *selected_fields* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_INSERT_ALL` | Insert Data into BigQuery Table | write | `dataset_id`, `project_id`, `rows`, `table_id`; *ignoreUnknownValues, skipInvalidRows, templateSuffix, traceId* |  | `bigquery`, `cloud-platform`, `bigquery.insertdata` |
| `GOOGLEBIGQUERY_INSERT_JOB` | Insert BigQuery Job | write | `configuration`, `project_id`; *jobReference* |  | `bigquery`, `cloud-platform`, `devstorage.full_control`, `devstorage.read_only`, `devstorage.read_write` |
| `GOOGLEBIGQUERY_INSERT_JOB_WITH_UPLOAD` | Insert BigQuery Job with Upload | write | `configuration`, `file_to_upload`, `project_id`; *jobReference* |  | `bigquery`, `cloud-platform`, `devstorage.full_control`, `devstorage.read_only`, `devstorage.read_write` |
| `GOOGLEBIGQUERY_LIST_ANALYTICS_HUB_LISTINGS` | List Analytics Hub Listings | read | `parent`; *access_token, alt, callback, fields, key, oauth_token, pageSize, pageToken and 5 more* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_LIST_BIG_QUERY_CONNECTIONS` | List BigQuery Connections | read | `parent`; *pageSize, pageToken* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_LIST_CAPACITY_COMMITMENTS` | List BigQuery Capacity Commitments | read | `parent`; *page_size, page_token* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_LIST_DATAEXCHANGES_LISTINGS` | List Data Exchange Listings | read | `parent`; *access_token, alt, callback, fields, key, oauth_token, pageSize, pageToken and 5 more* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_LIST_DATASETS` | List BigQuery Datasets | read | `project_id`; *all, filter, max_results, page_token* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_LIST_JOBS` | List BigQuery Jobs | read | `project_id`; *all_users, max_creation_time, max_results, min_creation_time, page_token, parent_job_id, projection, state_filter* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_LIST_LOCATIONS` | List BigQuery Data Transfer Locations | read | `name`; *filter, pageSize, pageToken* |  | `cloud-platform` |
| `GOOGLEBIGQUERY_LIST_LOCATIONS_CONNECTIONS` | List Connections in Location | read | `maxResults`, `parent`; *pageToken* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_LIST_LOCATIONS_DATAPOLICIES` | List BigQuery Location Data Policies | read | `parent`; *filter, pageSize, pageToken* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_LIST_MODELS` | List BigQuery Models | read | `dataset_id`, `project_id`; *max_results, page_token* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_LIST_ORGANIZATION_DATA_EXCHANGES` | List Organization Data Exchanges | read | `organization`; *access_token, alt, callback, fields, key, oauth_token, pageSize, pageToken and 5 more* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_LIST_PROJECTS` | List BigQuery Projects | read | *max_results, page_token* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_LIST_QUERY_TEMPLATES` | List Analytics Hub Query Templates | read | `parent`; *pageSize, pageToken* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_LIST_RESERVATIONS` | List BigQuery Reservations | read | `parent`; *access_token, alt, callback, fields, key, oauth_token, page_size, page_token and 5 more* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_LIST_RESERVATION_ASSIGNMENTS` | List BigQuery Reservation Assignments | read | `parent`; *page_size, page_token* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_LIST_RESERVATION_GROUPS` | List BigQuery Reservation Groups | read | `parent`; *page_size, page_token* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_LIST_ROUTINES` | List BigQuery Routines | read | `dataset_id`, `project_id`; *filter, maxResults, pageToken, readMask* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_LIST_ROW_ACCESS_POLICIES` | List BigQuery Row Access Policies | read | `dataset_id`, `project_id`, `table_id`; *page_size, page_token* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_LIST_TABLES` | List BigQuery Tables | read | `dataset_id`, `project_id`; *max_results, page_token* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_LIST_TABLE_DATA` | List BigQuery Table Data | read | `dataset_id`, `project_id`, `table_id`; *format_options_use_int64_timestamp, max_results, page_token, selected_fields, start_index* |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_PATCH_DATASET` | Patch BigQuery Dataset | write | `dataset_id`, `project_id`; *access, datasetReference, defaultCollation, defaultEncryptionConfiguration, defaultPartitionExpirationMs, defaultRoundingMode, defaultTableExpirationMs, description and 8 more* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_PATCH_MODEL` | Patch BigQuery ML Model | write | `dataset_id`, `model_id`, `project_id`; *description, encryptionConfiguration, expirationTime, friendlyName, labels* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_PATCH_TABLE` | Patch BigQuery Table | write | `dataset_id`, `project_id`, `table_id`; *autodetect_schema, description, expirationTime, friendlyName, labels, requirePartitionFilter* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_QUERY` | Query | write | `project_id`, `query`; *location, max_results, timeout_ms, use_legacy_sql* |  |  |
| `GOOGLEBIGQUERY_SEARCH_ALL_ASSIGNMENTS` | Search All BigQuery Reservation Assignments | read | `parent`; *access_token, alt, callback, fields, key, oauth_token, page_size, page_token and 6 more* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_SET_ROUTINE_IAM_POLICY` | Set BigQuery Routine IAM Policy | write | `dataset_id`, `policy`, `project_id`, `routine_id`; *update_mask* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_TEST_ROUTINE_IAM_PERMISSIONS` | Test BigQuery Routine IAM Permissions | read | `dataset_id`, `permissions`, `project_id`, `routine_id` |  | `bigquery`, `cloud-platform`, `bigquery.readonly`, `cloud-platform.read-only` |
| `GOOGLEBIGQUERY_UNDELETE_DATASET` | Undelete BigQuery Dataset | write | `dataset_id`, `project_id`; *deletion_time* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_UPDATE_CONNECTION` | Update BigQuery Connection | write | `name`, `updateMask`; *aws, azure, cloudResource, cloudSpanner, cloudSql, description, friendlyName, kmsKeyName and 2 more* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_UPDATE_DATASET` | Update BigQuery Dataset | write | `datasetReference`, `dataset_id`, `location`, `project_id`; *access, defaultCollation, defaultEncryptionConfiguration, defaultPartitionExpirationMs, defaultRoundingMode, defaultTableExpirationMs, description, externalDatasetReference and 6 more* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_UPDATE_ROUTINE` | Update BigQuery Routine | write | `dataset_id`, `definitionBody`, `project_id`, `routineReference`, `routineType`, `routine_id`; *arguments, dataGovernanceType, description, determinismLevel, importedLibraries, language, remoteFunctionOptions, returnTableType and 4 more* |  | `bigquery`, `cloud-platform` |
| `GOOGLEBIGQUERY_UPDATE_TABLE` | Update BigQuery Table | write | `dataset_id`, `project_id`, `tableReference`, `table_id`; *autodetect_schema, clustering, defaultCollation, description, encryptionConfiguration, expirationTime, externalDataConfiguration, friendlyName and 5 more* |  | `bigquery`, `cloud-platform` |
