# Supabase: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

230 tools: 111 read, 119 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `SUPABASE_ACCEPT_INVITE_EXTERNAL_JIT_ACCESS` | Accept Invite External Jit Access | write | `body`, `ref` |  |  |
| `SUPABASE_ALPHA_CREATE_API_KEY` | Create project api key | write | `name`, `ref`, `type`; *description* |  | `all` |
| `SUPABASE_ALPHA_DELETE_API_KEY` | Delete an API key from the project | write | `id`, `ref` |  |  |
| `SUPABASE_ALPHA_DELETE_THIRD_PARTY_AUTH_INTEGRATION` | Delete third party auth config | write | `ref`, `tpa_id` |  | `project.auth.write` |
| `SUPABASE_ALPHA_GET_THIRD_PARTY_INTEGRATION` | Get a third-party integration | read | `ref`, `tpa_id` |  | `project.auth.config.read` |
| `SUPABASE_ALPHA_LIST_THIRD_PARTY_AUTH_INTEGRATIONS` | List third-party auth integrations for project | read | `ref` |  |  |
| `SUPABASE_ALPHA_UPDATE_API_KEY` | Update an API key for the project | write | `id`, `ref`; *description* |  | `Secrets.Write` |
| `SUPABASE_APPEND_RESUMABLE_UPLOAD` | Append Resumable Upload | write | `file`, `object_path`, `project_ref`, `upload_offset`; *cache_control, content_type, x_upsert* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_APPEND_RESUMABLE_UPLOAD_SIGNED` | Append Resumable Upload Signed | write | `file`, `object_path`, `project_ref`, `upload_offset`, `x_signature`; *cache_control, content_type, x_upsert* |  |  |
| `SUPABASE_APPLY_A_MIGRATION` | Apply a database migration | write | `query`, `ref`; *idempotency_key, name, rollback* |  | `database:write`, `database_migrations_write` |
| `SUPABASE_APPLY_PROJECT_ADDON` | Apply Project Addon | write | `body`, `ref` |  | `infra_add_ons_write` |
| `SUPABASE_AUTHORIZE_JIT_ACCESS` | Authorize Jit Access | write | `body`, `ref` |  | `database:read`, `database_jit_read` |
| `SUPABASE_BETA_ACTIVATE_CUSTOM_HOSTNAME` | Beta activate custom hostname for project | write | `ref` |  | `Domains.Write` |
| `SUPABASE_BETA_ACTIVATE_VANITY_SUBDOMAIN` | Activate vanity subdomain for project | write | `ref`, `vanity_subdomain` |  | `projects.vanity-subdomain.write` |
| `SUPABASE_BETA_AUTHORIZE_USER_THROUGH_OAUTH` | Authorize user through OAuth | read | `client_id`, `redirect_uri`, `response_type`; *code_challenge, code_challenge_method, response_mode, scope, state* |  |  |
| `SUPABASE_BETA_CHECK_VANITY_SUBDOMAIN_AVAILABILITY` | Check vanity subdomain availability | read | `ref`, `vanity_subdomain` |  | `write.vanity-subdomain` |
| `SUPABASE_BETA_CREATE_READ_REPLICA` | Setup read replica for project | write | `read_replica_region`, `ref` |  | `Database:Write` |
| `SUPABASE_BETA_ENABLE_DATABASE_WEBHOOKS` | Enable project database webhooks | write | `ref` |  | `Database:Write` |
| `SUPABASE_BETA_GET_PROJECT_CUSTOM_HOSTNAME_CONFIG` | Beta get project's custom hostname config | write | `ref` |  | `Domains:Read`, `Domains:Write` |
| `SUPABASE_BETA_GET_PROJECT_NETWORK_BANS` | Retrieve network bans for project | read | `ref` |  |  |
| `SUPABASE_BETA_GET_PROJECT_NETWORK_RESTRICTIONS` | Get project network restrictions | read | `ref` |  | `Projects:Read` |
| `SUPABASE_BETA_GET_PROJECT_PGSODIUM_CONFIG` | Get project pgsodium config | read | `ref` |  |  |
| `SUPABASE_BETA_GET_PROJECT_SSL_ENFORCEMENT_CONFIG` | Get project SSL enforcement configuration | read | `ref` |  | `Database:Read`, `Database:Write` |
| `SUPABASE_BETA_GET_VANITY_SUBDOMAIN_CONFIG` | Get current vanity subdomain config | read | `ref` |  | `Domains:Read`, `Domains:Write` |
| `SUPABASE_BETA_REMOVE_NETWORK_BANS` | Remove project network bans | write | `ipv4_addresses`, `ref` |  | `Projects:Write` |
| `SUPABASE_BETA_REMOVE_READ_REPLICA` | Remove read replica | write | `database_identifier`, `ref` |  |  |
| `SUPABASE_BETA_RUN_SQL_QUERY` | Execute project database query | write | `query`, `ref`; *read_only* |  | `Database:Write` |
| `SUPABASE_BETA_UPDATE_PROJECT_NETWORK_RESTRICTIONS` | Beta update project network restrictions | write | `ref`; *dbAllowedCidrs, dbAllowedCidrsV6* |  | `Projects:Write` |
| `SUPABASE_BETA_UPGRADE_PROJECT_POSTGRES_VERSION` | Upgrade the project's PostgreSQL version | write | `ref`, `release_channel`, `target_version` |  | `projects.write` |
| `SUPABASE_CANCEL_PROJECT_RESTORATION` | Cancel Project Restoration | write | `ref` |  | `projects:write`, `project_admin_write` |
| `SUPABASE_COPY_STORAGE_OBJECT` | Copy Storage Object | write | `body`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_COUNT_ACTION_RUNS` | Count action runs | read | `ref` |  | `environment:read` |
| `SUPABASE_COUNT_FUNCTION_RESULTS` | Count Function Results | read | `function_name`, `project_ref`; *arguments, count, filters, limit, offset, order, schema, select* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_COUNT_TABLE_ROWS` | Count Table Rows | read | `project_ref`, `table`; *count, filters, limit, offset, order, schema, select* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_CREATE_A_PROJECT` | Create new project | write | `db_pass`, `name`, `organization_id`, `region`; *desired_instance_size, kps_enabled, plan, postgres_engine, release_channel, template_url* |  | `projects.write` |
| `SUPABASE_CREATE_BULK_SECRETS` | Bulk create secrets | write | `ref`, `secrets` |  | `secrets:write` |
| `SUPABASE_CREATE_DATABASE_BRANCH` | Create a database branch | write | `branch_name`, `ref`; *desired_instance_size, git_branch, persistent, postgres_engine, region, release_channel* |  | `Environment:Write` |
| `SUPABASE_CREATE_FUNCTION` | Create a function | write | `body`, `name`, `ref`, `slug`; *entrypoint_path, import_map, import_map_path, verify_jwt* |  | `Edge Functions:Write` |
| `SUPABASE_CREATE_LOGIN_ROLE` | Create CLI login role | write | `read_only`, `ref` |  | `database:write` |
| `SUPABASE_CREATE_ORGANIZATION` | Create an organization | write | `name` |  | `Organizations:Write` |
| `SUPABASE_CREATE_PROJECT_SIGNING_KEY` | Create project signing key | write | `algorithm`, `ref`; *status* |  | `secrets:write` |
| `SUPABASE_CREATE_RESUMABLE_UPLOAD` | Create Resumable Upload | write | `project_ref`, `upload_length`, `upload_metadata`; *cache_control, content_type, file, x_upsert* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_CREATE_RESUMABLE_UPLOAD_SIGNED` | Create Resumable Upload Signed | write | `project_ref`, `upload_length`, `upload_metadata`, `x_signature`; *cache_control, content_type, file, x_upsert* |  |  |
| `SUPABASE_CREATE_SIGNED_DOWNLOAD_URL` | Create Signed Download Url | write | `body`, `bucket_name`, `object_path`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_CREATE_SIGNED_DOWNLOAD_URLS` | Create Signed Download Urls | write | `body`, `bucket_name`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_CREATE_SIGNED_UPLOAD_URL` | Create Signed Upload URL | write | `bucket_name`, `object_path`, `project_ref`; *upsert* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_CREATE_SSO_PROVIDER` | Create SSO provider configuration | write | `ref`, `type`; *attribute__mapping__keys, domains, metadata_url, metadata_xml* |  |  |
| `SUPABASE_CREATE_STORAGE_BUCKET` | Create Storage Bucket | write | `body`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_CREATE_THIRD_PARTY_AUTH_INTEGRATION` | Create a new third-party auth integration | write | `ref`; *jwks_url, oidc_issuer_url* |  |  |
| `SUPABASE_DELETE_ALL_BRANCHES` | Delete all branches | write | `ref` |  | `Environment:Write` |
| `SUPABASE_DELETE_BUCKET_LIFECYCLE` | Delete Bucket Lifecycle | write | `bucket_id`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_DELETE_CUSTOM_HOSTNAME_CONFIG` | Delete custom hostname config | write | `ref` |  | `Domains:Write` |
| `SUPABASE_DELETE_DATABASE_BRANCH` | Delete branch by id | write | `branch_id` |  | `Environment:Write` |
| `SUPABASE_DELETE_EDGE_FUNCTION` | Delete edge function | write | `function_slug`, `ref` |  | `edge_functions:write` |
| `SUPABASE_DELETE_FUNCTION` | Delete an edge function by slug | write | `function_slug`, `ref` |  | `functions.write` |
| `SUPABASE_DELETE_INVITE_EXTERNAL_JIT_ACCESS` | Delete Invite External Jit Access | write | `invite_id`, `ref` |  | `database_jit_write` |
| `SUPABASE_DELETE_JIT_ACCESS` | Delete Jit Access | write | `ref`, `user_id` |  | `database_jit_write` |
| `SUPABASE_DELETE_LOGIN_ROLES` | Delete CLI login roles | write | `ref` |  | `database_write`, `database:write` |
| `SUPABASE_DELETE_PROJECT` | Delete a project | write | `ref` |  | `projects:write` |
| `SUPABASE_DELETE_PROJECT_VANITY_SUBDOMAIN` | Delete vanity subdomain for project | write | `ref` |  | `Domains:Write` |
| `SUPABASE_DELETE_SECRETS` | Bulk delete secrets | write | `ref`, `secret_names` |  | `secrets:write`, `edge_functions_secrets_write` |
| `SUPABASE_DELETE_SSO_PROVIDER` | Remove an SSO provider | write | `provider_id`, `ref` |  | `projects.auth.providers.saml.delete` |
| `SUPABASE_DELETE_STORAGE_BUCKET` | Delete Storage Bucket | write | `bucket_id`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_DELETE_STORAGE_OBJECT` | Delete Storage Object | write | `bucket_name`, `object_path`, `project_ref`; *version_id* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_DELETE_STORAGE_OBJECTS` | Delete Storage Objects | write | `body`, `bucket_name`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_DELETE_TABLE_ROWS` | Delete Table Rows | write | `project_ref`, `table`; *allow_all_rows, count, filters, returning, schema, select* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_DEPLOY_FUNCTION` | Deploy function | write | `ref`; *bundleOnly, file, file_content, file_url, slug* |  | `Edge Functions:Write` |
| `SUPABASE_DISABLE_PREVIEW_BRANCHING` | Disable preview branching | write | `ref` |  | `projects.read`, `projects.write` |
| `SUPABASE_DISABLE_PROJECT_READONLY` | Disable project readonly mode | write | `ref` |  |  |
| `SUPABASE_DOWNLOAD_PRIVATE_OBJECT` | Download private object | read | `bucket_name`, `object_path`, `project_ref`; *auth_mode, download_filename, version_id* |  | `secrets:read` |
| `SUPABASE_DOWNLOAD_PUBLIC_OBJECT` | Download Public Object | read | `bucket_name`, `object_path`, `project_ref`; *download, version_id* |  |  |
| `SUPABASE_DOWNLOAD_SIGNED_OBJECT` | Download Signed Object | read | `bucket_name`, `object_path`, `project_ref`, `token`; *download* |  |  |
| `SUPABASE_DOWNLOAD_STORAGE_OBJECT` | Download Storage Object | read | `bucket_name`, `object_path`, `project_ref`; *download, version_id* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_EMPTY_STORAGE_BUCKET` | Empty Storage Bucket | write | `bucket_id`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_EXCHANGE_OAUTH_TOKEN` | Exchange auth code for access and refresh token | write | `client_id`, `client_secret`, `grant_type`; *code, code_verifier, redirect_uri, refresh_token* |  | `No explicit required OAuth scopes documented for this endpoint` |
| `SUPABASE_GENERATE_TYPESCRIPT_TYPES` | Generate TypeScript types | read | `ref`; *included_schemas* |  | `Database:Read` |
| `SUPABASE_GETS_PROJECT_S_AUTH_CONFIG` | Get project's auth config | read | `ref` |  | `Auth:Read` |
| `SUPABASE_GETS_PROJECT_S_SERVICE_HEALTH_STATUS` | Get project's service health status | read | `ref`; *services* |  |  |
| `SUPABASE_GET_ACTION_RUN` | Get action run status | read | `ref`, `run_id` |  | `environment:read` |
| `SUPABASE_GET_ACTION_RUN_LOGS` | Get action run logs | read | `ref`, `run_id` |  | `environment:read` |
| `SUPABASE_GET_AUTHENTICATED_IMAGE_HEADERS` | Get Authenticated Image Headers | read | `bucket_name`, `object_path`, `project_ref`; *download, format, gravity, height, quality, resize, version_id, width and 2 more* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_GET_AUTHENTICATED_OBJECT_HEADERS` | Get Authenticated Object Headers | read | `bucket_name`, `object_path`, `project_ref`; *format, gravity, height, quality, resize, version_id, width, x_offset and 1 more* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_GET_AUTHENTICATED_OBJECT_INFO` | Get Authenticated Object Info | read | `bucket_name`, `object_path`, `project_ref`; *format, gravity, height, quality, resize, version_id, width, x_offset and 1 more* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_GET_AUTHENTICATED_OBJECT_INFO_HEADERS` | Get Authenticated Object Info Headers | read | `bucket_name`, `object_path`, `project_ref`; *format, gravity, height, quality, resize, version_id, width, x_offset and 1 more* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_GET_AUTH_CONFIG` | Get project auth config | read | `ref` |  | `auth_config_read` |
| `SUPABASE_GET_AVAILABLE_REGIONS` | Get Available Regions | read | `organization_slug`; *continent, desired_instance_size* |  | `organizations:read` |
| `SUPABASE_GET_BACKUP_SCHEDULE` | Get Backup Schedule | read | `ref` |  | `database:read`, `backups_read` |
| `SUPABASE_GET_BRANCH` | Get a database branch | read | `name`, `ref` |  | `environment:read` |
| `SUPABASE_GET_BUCKET_LIFECYCLE` | Get Bucket Lifecycle | read | `bucket_id`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_GET_BUCKET_LIFECYCLE_HEADERS` | Get Bucket Lifecycle Headers | read | `bucket_id`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_GET_DATABASE_BRANCH_CONFIG` | Get database branch config | read | `branch_id` |  |  |
| `SUPABASE_GET_DATABASE_DISK` | Get Database Disk | read | `ref` |  | `infra_disk_config_read` |
| `SUPABASE_GET_DATABASE_OPENAPI` | Get Database Openapi | read | `ref`; *schema_name* |  | `database:read`, `database_read` |
| `SUPABASE_GET_DISK_UTILIZATION` | Get Disk Utilization | read | `ref` |  | `infra_disk_config_read` |
| `SUPABASE_GET_FUNCTION` | Retrieve a function | read | `function_slug`, `ref` |  |  |
| `SUPABASE_GET_FUNCTION_BODY` | Retrieve a function body | read | `function_slug`, `ref` |  | `Edge Functions:Read` |
| `SUPABASE_GET_HEALTH` | Get API Health Status | read |  |  |  |
| `SUPABASE_GET_JIT_ACCESS` | Get Jit Access | read | `ref` |  | `database:read`, `database_jit_read` |
| `SUPABASE_GET_JIT_ACCESS_CONFIG` | Get JIT access config | read | `ref` |  | `database:read` |
| `SUPABASE_GET_MIGRATION` | Get a migration | read | `ref`, `version` |  | `database_migrations_read`, `database:read` |
| `SUPABASE_GET_ORGANIZATION` | Get information about an organization | read | `slug` |  | `organizations.read` |
| `SUPABASE_GET_ORGANIZATION_ENTITLEMENTS` | Get Organization Entitlements | read | `slug` |  | `organizations:read`, `organization_admin_read` |
| `SUPABASE_GET_PROFILE` | Get Profile | read |  |  |  |
| `SUPABASE_GET_PROJECT` | Get project | read | `ref` |  | `projects:read` |
| `SUPABASE_GET_PROJECT_API_KEY` | Get project API key | read | `id`, `ref`; *reveal* |  | `secrets:read` |
| `SUPABASE_GET_PROJECT_API_KEYS` | Get project API keys | read | `ref`; *reveal* |  | `Secrets:Read` |
| `SUPABASE_GET_PROJECT_CONFIG` | Get project config | read | `ref` |  | `database_config_read`, `database_read`, `database_ssl_config_read`, `database_network_restrictions_read`, `auth_config_read`, `data_api_config_read`, `realtime_config_read`, `storage_config_read` |
| `SUPABASE_GET_PROJECT_DISK_AUTOSCALE_CONFIG` | Get Project Disk Autoscale Config | read | `ref` |  | `infra_disk_config_read` |
| `SUPABASE_GET_PROJECT_FUNCTION_COMBINED_STATS` | Get Project Function Combined Stats | read | `function_id`, `interval`, `ref` |  | `analytics_usage_read` |
| `SUPABASE_GET_PROJECT_LOGS` | Get project logs | write | `ref`; *iso_timestamp_end, iso_timestamp_start, sql* |  | `analytics:read` |
| `SUPABASE_GET_PROJECT_METRICS` | Get project metrics | read | `ref` |  | `analytics:read` |
| `SUPABASE_GET_PROJECT_PGBOUNCER_CONFIG` | Get project PgBouncer config | read | `ref` |  |  |
| `SUPABASE_GET_PROJECT_POSTGREST_CONFIG` | Get project's PostgREST config | read | `ref` |  |  |
| `SUPABASE_GET_PROJECT_POSTGRES_CONFIG` | Get project postgres config | read | `ref` |  |  |
| `SUPABASE_GET_PROJECT_READONLY_MODE_STATUS` | Get project's read-only mode status | read | `ref` |  | `Database:Read` |
| `SUPABASE_GET_PROJECT_SIGNING_KEYS` | Get project signing keys | read | `ref` |  | `secrets:read`, `auth_signing_keys_read` |
| `SUPABASE_GET_PROJECT_SUPAVISOR_CONFIG` | Get project Supavisor configuration | read | `ref` |  | `Database:Read` |
| `SUPABASE_GET_PROJECT_UPGRADE_ELIGIBILITY` | Get Project Upgrade Eligibility | read | `ref` |  |  |
| `SUPABASE_GET_PROJECT_UPGRADE_STATUS` | Get project upgrade status | read | `ref` |  |  |
| `SUPABASE_GET_PROJECT_USAGE_API_COUNT` | Get Project Usage Api Count | read | `ref`; *interval* |  | `analytics_usage_read` |
| `SUPABASE_GET_PROJECT_USAGE_REQUEST_COUNT` | Get Project Usage Request Count | read | `ref` |  | `analytics_usage_read` |
| `SUPABASE_GET_PUBLIC_IMAGE_HEADERS` | Get Public Image Headers | read | `bucket_name`, `object_path`, `project_ref`; *download, format, gravity, height, quality, resize, version_id, width and 2 more* |  |  |
| `SUPABASE_GET_PUBLIC_OBJECT_HEADERS` | Get Public Object Headers | read | `bucket_name`, `object_path`, `project_ref`; *format, gravity, height, quality, resize, version_id, width, x_offset and 1 more* |  |  |
| `SUPABASE_GET_PUBLIC_OBJECT_INFO` | Get Public Object Info | read | `bucket_name`, `object_path`, `project_ref`; *format, gravity, height, quality, resize, version_id, width, x_offset and 1 more* |  |  |
| `SUPABASE_GET_REALTIME_CONFIG` | Get Realtime Config | read | `ref` |  | `realtime_config_read` |
| `SUPABASE_GET_RESUMABLE_UPLOAD_BASE_OPTIONS` | Get TUS resumable upload base options | read |  |  |  |
| `SUPABASE_GET_RESUMABLE_UPLOAD_OPTIONS` | Get TUS resumable upload options | read | `wildcard` |  | `all` |
| `SUPABASE_GET_RESUMABLE_UPLOAD_STATUS` | Get Resumable Upload Status | read | `object_path`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_GET_RESUMABLE_UPLOAD_STATUS_SIGNED` | Get Resumable Upload Status Signed | read | `object_path`, `project_ref`, `x_signature` |  |  |
| `SUPABASE_GET_SIGNED_IMAGE_HEADERS` | Get Signed Image Headers | read | `bucket_name`, `object_path`, `project_ref`, `token`; *download* |  |  |
| `SUPABASE_GET_SIGNED_OBJECT_HEADERS` | Get Signed Object Headers | read | `bucket_name`, `object_path`, `project_ref`, `token`; *download* |  |  |
| `SUPABASE_GET_SQL_SNIPPET` | Get a specific SQL snippet | read | `id` |  |  |
| `SUPABASE_GET_SSO_PROVIDER` | Get a SSO provider by its UUID | read | `provider_id`, `ref` |  | `Auth:Read`, `Auth:Write` |
| `SUPABASE_GET_STORAGE_BUCKET` | Get Storage Bucket | read | `bucket_id`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_GET_STORAGE_BUCKET_HEADERS` | Get Storage Bucket Headers | read | `bucket_id`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_GET_STORAGE_BUCKET_LIST_HEADERS` | Get Storage Bucket List Headers | read | `project_ref`; *limit, offset, search, sort_column, sort_order* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_GET_STORAGE_CONFIG` | Get Storage Config | read | `ref` |  | `storage_config_read` |
| `SUPABASE_GET_STORAGE_HEALTH` | Get Storage Health | read | `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_GET_STORAGE_HEALTH_HEADERS` | Get Storage Health Headers | read | `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_GET_STORAGE_OBJECT_HEADERS` | Get Storage Object Headers | read | `bucket_name`, `object_path`, `project_ref`; *format, gravity, height, quality, resize, version_id, width, x_offset and 1 more* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_GET_STORAGE_OBJECT_INFO` | Get Storage Object Info | read | `bucket_name`, `object_path`, `project_ref`; *format, gravity, height, quality, resize, version_id, width, x_offset and 1 more* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_GET_STORAGE_OBJECT_INFO_HEADERS` | Get Storage Object Info Headers | read | `bucket_name`, `object_path`, `project_ref`; *format, gravity, height, quality, resize, version_id, width, x_offset and 1 more* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_GET_TABLE_SCHEMAS` | Get Table Schemas | read | `project_ref`, `table_names`; *exclude_null_values, include_indexes, include_relationships* |  | `Database:Read` |
| `SUPABASE_HANDLE_RESUMABLE_UPLOAD_SIGN_OPTIONS` | OPTIONS for resumable upload sign | read |  |  |  |
| `SUPABASE_HANDLE_RESUMABLE_UPLOAD_SIGN_OPTIONS_WITH_ID` | OPTIONS for resumable upload sign | read | `wildcard` |  |  |
| `SUPABASE_INSERT_TABLE_ROWS` | Insert Table Rows | write | `project_ref`, `rows`, `table`; *columns, count, missing, on_conflict, resolution, return_preference, schema, select* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_INVITE_EXTERNAL_JIT_ACCESS` | Invite External Jit Access | write | `body`, `ref` |  | `database_jit_write` |
| `SUPABASE_INVOKE_DATABASE_FUNCTION` | Invoke Database Function | write | `function_name`, `project_ref`; *arguments, count, filters, limit, offset, order, schema, select* |  | `secrets:read` |
| `SUPABASE_INVOKE_EDGE_FUNCTION` | Invoke Edge Function | write | `function_slug`, `project_ref`; *auth_mode, body, headers, method, response_type* |  | `Edge Functions:Execute` |
| `SUPABASE_INVOKE_READ_ONLY_FUNCTION` | Invoke Read Only Function | read | `function_name`, `project_ref`; *arguments, count, filters, limit, offset, order, schema, select* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_LIST_ACTION_RUNS` | List Action Runs | read | `ref`; *limit, offset* |  | `environment:read`, `action_runs_read` |
| `SUPABASE_LIST_ALL_ORGANIZATIONS` | List all organizations | read |  |  | `organizations:read` |
| `SUPABASE_LIST_ALL_PROJECTS` | List all projects | read |  |  | `Projects.Read` |
| `SUPABASE_LIST_BACKUPS` | List project database backups | read | `ref` |  | `Database:Read` |
| `SUPABASE_LIST_BRANCHES` | List database branches | read | `ref`; *limit, offset* |  | `environment:read` |
| `SUPABASE_LIST_BUCKETS` | List all buckets | read | `ref` |  |  |
| `SUPABASE_LIST_DATABASE_BRANCHES` | List all database branches | read | `ref` |  | `Environment:Read` |
| `SUPABASE_LIST_FUNCTIONS` | List all functions | read | `ref` |  |  |
| `SUPABASE_LIST_JIT_ACCESS` | List Jit Access | read | `ref` |  | `database_jit_write` |
| `SUPABASE_LIST_MIGRATION_HISTORY` | List migration history | read | `ref` |  | `database_migrations_read` |
| `SUPABASE_LIST_ORGANIZATION_MEMBERS` | List members of an organization | read | `slug` |  | `organizations.members.read` |
| `SUPABASE_LIST_ORGANIZATION_PROJECTS` | List Organization Projects | read | `slug`; *limit, offset, search, sort, statuses* |  | `projects:read`, `organization_projects_read` |
| `SUPABASE_LIST_PROJECT_ADDONS` | List Project Addons | read | `ref` |  | `infra_add_ons_read` |
| `SUPABASE_LIST_PROJECT_RESTORE_VERSIONS` | List Project Restore Versions | read | `ref` |  | `projects:read`, `project_admin_read` |
| `SUPABASE_LIST_PROJECT_SECRETS` | List project secrets | read | `ref` |  | `secrets:read` |
| `SUPABASE_LIST_SECRETS` | List all secrets | read | `ref` |  | `Secrets Read` |
| `SUPABASE_LIST_SQL_SNIPPETS` | List SQL snippets for the logged in user | read | *project_ref* |  | `database:read` |
| `SUPABASE_LIST_SSO_PROVIDERS` | List all SSO providers | read | `ref` |  |  |
| `SUPABASE_LIST_STORAGE_BUCKETS` | List Storage Buckets | read | `project_ref`; *limit, offset, search, sort_column, sort_order* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_LIST_STORAGE_OBJECTS` | List Storage Objects | write | `body`, `bucket_name`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_LIST_STORAGE_OBJECTS_V2` | List Storage Objects V2 | write | `body`, `bucket_name`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_LIST_TABLES` | List Database Tables | read | `project_ref`; *include_metadata, include_system_schemas, include_views, schemas* |  | `Database:Read` |
| `SUPABASE_MERGE_BRANCH` | Merge a database branch | write | `branch_id_or_ref`; *migration_version* |  | `environment:write` |
| `SUPABASE_MODIFY_DATABASE_DISK` | Modify Database Disk | write | `body`, `ref` |  | `infra_disk_config_write` |
| `SUPABASE_MOVE_STORAGE_OBJECT` | Move Storage Object | write | `body`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_OAUTH_AUTHORIZE_PROJECT_CLAIM` | Oauth Authorize Project Claim | write | `client_id`, `project_ref`, `redirect_uri`, `response_type`; *code_challenge, code_challenge_method, response_mode, state* |  | `organization_admin_write`, `project_admin_write` |
| `SUPABASE_PATCH_MIGRATION` | Patch a migration | write | `ref`, `version`; *name, rollback* |  | `database:write` |
| `SUPABASE_PATCH_NETWORK_RESTRICTIONS` | Patch project network restrictions | write | `ref`; *add, remove* |  | `projects:write` |
| `SUPABASE_PAUSE_PROJECT` | Pause Project | write | `ref` |  | `projects:write`, `project_admin_write` |
| `SUPABASE_POST_RESUMABLE_UPLOAD` | Post Resumable Upload | write | `object_path`, `project_ref`, `upload_length`, `upload_metadata`; *cache_control, content_type, file, x_upsert* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_POST_RESUMABLE_UPLOAD_SIGNED` | Post Resumable Upload Signed | write | `object_path`, `project_ref`, `upload_length`, `upload_metadata`, `x_signature`; *cache_control, content_type, file, x_upsert* |  |  |
| `SUPABASE_PURGE_BUCKET_CACHE` | Purge Bucket Cache | write | `bucket_name`, `project_ref`; *transformations* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_PURGE_OBJECT_CACHE` | Purge Object Cache | write | `bucket_name`, `object_path`, `project_ref`; *transformations* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_PURGE_STORAGE_CACHE` | Purge Storage Cache | write | `project_ref`; *transformations* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_PUSH_BRANCH` | Push a database branch | write | `branch_id_or_ref`; *migration_version* |  | `environment:write` |
| `SUPABASE_REMOVE_PROJECT_ADDON` | Remove Project Addon | write | `addon_variant`, `ref` |  | `infra_add_ons_write` |
| `SUPABASE_RENDER_AUTHENTICATED_IMAGE` | Render Authenticated Image | read | `bucket_name`, `object_path`, `project_ref`; *download, format, gravity, height, quality, resize, version_id, width and 2 more* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_RENDER_PUBLIC_IMAGE` | Render Public Image | read | `bucket_name`, `object_path`, `project_ref`; *download, format, gravity, height, quality, resize, version_id, width and 2 more* |  |  |
| `SUPABASE_RENDER_SIGNED_IMAGE` | Render Signed Image | read | `bucket_name`, `object_path`, `project_ref`, `token`; *download* |  |  |
| `SUPABASE_REPLACE_BUCKET_LIFECYCLE` | Replace Bucket Lifecycle | write | `body`, `bucket_id`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_RESET_DATABASE_BRANCH` | Reset a database branch | write | `branch_id` |  | `Environment:Write` |
| `SUPABASE_RESTART_PROJECT` | Restart Project | write | `ref` |  | `projects:write`, `project_admin_write` |
| `SUPABASE_RESTORE_BRANCH` | Restore Branch | write | `branch_id_or_ref` |  | `environment:write`, `branching_development_write`, `branching_production_write` |
| `SUPABASE_RESTORE_PITR_BACKUP` | Restore database PITR backup | write | `recovery_time_target_unix`, `ref` |  |  |
| `SUPABASE_RESTORE_PROJECT` | Restore Project | write | `ref` |  | `projects:write`, `project_admin_write` |
| `SUPABASE_ROLLBACK_MIGRATIONS` | Rollback Migrations | write | `gte`, `ref` |  | `database:write`, `database_migrations_write` |
| `SUPABASE_RUN_READ_ONLY_QUERY` | Execute read-only database query | read | `query`, `ref`; *parameters* |  | `database_read` |
| `SUPABASE_SELECT_FROM_TABLE` | Select From Table | read | `project_ref`, `select`, `table`; *filters, limit, offset, order* |  | `Database:Read` |
| `SUPABASE_SHUTDOWN_REALTIME` | Shutdown Realtime | write | `ref` |  | `realtime_config_write` |
| `SUPABASE_TERMINATE_RESUMABLE_UPLOAD` | Terminate Resumable Upload | write | `object_path`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_TERMINATE_RESUMABLE_UPLOAD_SIGNED` | Terminate Resumable Upload Signed | write | `object_path`, `project_ref`, `x_signature` |  |  |
| `SUPABASE_UPDATE_ACTION_RUN_STATUS` | Update Action Run Status | write | `body`, `ref`, `run_id` |  | `environment:write`, `action_runs_write` |
| `SUPABASE_UPDATE_API_KEY` | Update API key | write | `id`, `ref`; *description, name, reveal, secret_jwt_template* |  | `secrets:write` |
| `SUPABASE_UPDATE_A_FUNCTION` | Update a function | write | `function_slug`, `ref`; *body, entrypoint_path, import_map, import_map_path, name, slug, verify_jwt* |  | `Edge Functions:Write` |
| `SUPABASE_UPDATE_BACKUP_SCHEDULE` | Update Backup Schedule | write | `body`, `ref` |  | `database:write`, `backups_write` |
| `SUPABASE_UPDATE_BRANCH` | Update branch configuration | write | `branch_id_or_ref`; *branch_name, git_branch, notify_url, persistent, request_review, reset_on_push, status* |  | `branching_production_write`, `branching_development_write` |
| `SUPABASE_UPDATE_DATABASE_BRANCH_CONFIG` | Update database branch config | write | `branch_id`; *branch_name, git_branch, persistent, reset_on_push, status* |  | `Environment:Write` |
| `SUPABASE_UPDATE_DATABASE_PASSWORD` | Update database password | write | `password`, `ref` |  | `database_write`, `database:write` |
| `SUPABASE_UPDATE_FUNCTIONS` | Bulk update functions | write | `functions`, `ref` |  | `Edge Functions:Write` |
| `SUPABASE_UPDATE_JIT_ACCESS` | Update Jit Access | write | `body`, `ref` |  | `database_jit_write` |
| `SUPABASE_UPDATE_JIT_ACCESS_CONFIG` | Update JIT access config | write | `ref`, `state` |  | `database:write` |
| `SUPABASE_UPDATE_PGSODIUM_CONFIG` | Update pgsodium root key | write | `ref`, `root_key` |  | `Secrets:Read`, `Secrets:Write` |
| `SUPABASE_UPDATE_PROJECT` | Update a project | write | `name`, `ref` |  | `projects:write` |
| `SUPABASE_UPDATE_PROJECT_AUTH_CONFIG` | Update project's auth config | write | `ref`; *disable_signup, external_anonymous_users_enabled, external_apple_client_id, external_apple_enabled, external_apple_secret, external_azure_client_id, external_azure_enabled, external_azure_secret and 51 more* |  | `Auth:Write` |
| `SUPABASE_UPDATE_PROJECT_CUSTOM_HOSTNAME` | Update project's custom hostname configuration | write | `custom_hostname`, `ref` |  | `project.custom_domain.write` |
| `SUPABASE_UPDATE_PROJECT_LEGACY_API_KEYS` | Update project legacy API keys | write | `enabled`, `ref` |  | `secrets:write` |
| `SUPABASE_UPDATE_PROJECT_POSTGREST_CONFIG` | Update project's PostgREST config | write | `ref`; *db_extra_search_path, db_pool, db_schema, max_rows* |  | `Rest:Read`, `Rest:Write` |
| `SUPABASE_UPDATE_PROJECT_POSTGRES_CONFIG` | Update project's postgres config | write | `ref`; *effective_cache_size, maintenance_work_mem, max_connections, max_locks_per_transaction, max_parallel_maintenance_workers, max_parallel_workers, max_parallel_workers_per_gather, max_slot_wal_keep_size and 9 more* |  | `Database:Write` |
| `SUPABASE_UPDATE_PROJECT_SUPAVISOR_CONFIG` | Update database pooler config | write | `ref`; *default_pool_size, pool_mode* |  | `Database:Read` |
| `SUPABASE_UPDATE_REALTIME_CONFIG` | Update Realtime Config | write | `body`, `ref` |  | `realtime_config_write` |
| `SUPABASE_UPDATE_SSL_ENFORCEMENT_CONFIG` | Update SSL enforcement config | write | `ref`; *requestedConfig__database* |  | `Database.Write` |
| `SUPABASE_UPDATE_SSO_PROVIDER` | Update an SSO provider by its UUID | write | `provider_id`, `ref`; *attribute__mapping__keys, domains, metadata_url, metadata_xml* |  | `Auth:Write` |
| `SUPABASE_UPDATE_STORAGE_BUCKET` | Update Storage Bucket | write | `body`, `bucket_id`, `project_ref` |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_UPDATE_STORAGE_CONFIG` | Update Storage Config | write | `body`, `ref` |  | `storage_config_write` |
| `SUPABASE_UPDATE_STORAGE_OBJECT` | Update Storage Object | write | `bucket_name`, `file`, `object_path`, `project_ref`; *cache_control, content_type, x_upsert* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_UPDATE_TABLE_ROWS` | Update Table Rows | write | `project_ref`, `table`, `updates`; *allow_all_rows, columns, count, filters, returning, schema, select* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_UPLOAD_SIGNED_OBJECT` | Upload Signed Object | write | `bucket_name`, `file`, `object_path`, `project_ref`, `token`; *cache_control, content_type, x_upsert* |  |  |
| `SUPABASE_UPLOAD_STORAGE_OBJECT` | Upload Storage Object | write | `bucket_name`, `file`, `object_path`, `project_ref`; *cache_control, content_type, x_upsert* |  | `secrets:read`, `api_gateway_keys_read`, `api_gateway_keys_secret_read` |
| `SUPABASE_UPSERT_MIGRATION` | Upsert migration | write | `query`, `ref`; *idempotency_key, name, rollback* |  | `database:write`, `database_migrations_write` |
| `SUPABASE_VERIFY_CUSTOM_HOSTNAME_DNS` | Reverify custom hostname | write | `ref` |  | `Domains:Write` |
