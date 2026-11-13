# Sentry: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

210 tools: 110 read, 100 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `SENTRY_ACCESS_PROJECT_INFORMATION` | Access project information | read | `organization_id_or_slug`, `project_id_or_slug` |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_ADD_A_SYMBOL_SOURCE_TO_A_PROJECT` | Add a symbol source to a project | write | `layout__casing`, `layout__type`, `name`, `organization_id_or_slug`, `project_id_or_slug`, `type`; *access_key, bucket, client_email, filters__filetypes, filters__path_patterns, filters__requires_checksum, id, password and 6 more* |  | `project:admin`, `project:write` |
| `SENTRY_ADD_ORGANIZATION_MEMBER_VIA_EMAIL` | Add organization member via email | write | `email`, `organization_id_or_slug`; *orgRole, reinvite, sendInvite, teamRoles* |  | `member:admin`, `member:invite`, `member:write` |
| `SENTRY_ADD_REMOVE_USER_EMAIL_BY_ID` | Add or remove user email by id | write | `email`, `user_id` |  |  |
| `SENTRY_ADD_TEAM_MEMBER_IN_ORGANIZATION` | Add team member in organization | write | `member_id`, `organization_id_or_slug`, `team_id_or_slug` |  | `org:read`, `org:write`, `team:write` |
| `SENTRY_ADD_TEAM_TO_PROJECT` | Add team to project | write | `organization_id_or_slug`, `project_id_or_slug`, `team_id_or_slug` |  | `project:admin`, `project:write` |
| `SENTRY_BULK_MUTATE_AN_ORGANIZATIONS_ISSUES` | Bulk mutate an organization's issues | write | `organization_id_or_slug`; *assignedTo, discard, environment, hasSeen, id, inbox, isBookmarked, isPublic and 11 more* |  | `event:admin`, `event:write` |
| `SENTRY_BULK_REMOVE_AN_ORGANIZATIONS_ISSUES` | Bulk remove an organization's issues | write | `organization_id_or_slug`; *environment, id, limit, project, query, sort, viewId* |  | `event:admin` |
| `SENTRY_CREATE_AN_ALERT_FOR_AN_ORGANIZATION` | Create an alert for an organization | write | `actionFilters`, `config`, `name`, `organization_id_or_slug`, `triggers`; *enabled, environment* |  | `alerts:write`, `org:admin`, `org:write` |
| `SENTRY_CREATE_DASHBOARD_WITH_WIDGETS` | Create dashboard with widgets | write | `organization_id_or_slug`, `title`; *end, environment, filters, id, period, projects, start, utc and 1 more* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_CREATE_EXTERNAL_USER_FOR_ORGANIZATION` | Create external user for organization | write | `external_name`, `id`, `integration_id`, `organization_id_or_slug`, `provider`, `user_id`; *external_id* |  | `org:admin`, `org:write` |
| `SENTRY_CREATE_ORGANIZATION_ALERT_RULE` | Create organization alert rule | write | `aggregate`, `name`, `organization_id_or_slug`, `projects`, `query`, `thresholdType`, `timeWindow`, `triggers`; *activationCondition, comparisonDelta, dataset, environment, eventTypes, monitorType, owner, queryType and 1 more* |  | `alerts:write`, `org:admin`, `org:read`, `org:write` |
| `SENTRY_CREATE_ORGANIZATION_MONITOR` | Create organization monitor | write | `config`, `name`, `organization_id_or_slug`, `project`; *is_muted, owner, slug, status* |  | `alerts:write`, `org:admin`, `org:read`, `org:write` |
| `SENTRY_CREATE_ORGANIZATION_TEAM` | Create organization team | write | `organization_id_or_slug`; *name, slug* |  | `org:admin`, `org:write`, `team:write` |
| `SENTRY_CREATE_PROJECT_KEY_WITH_OPTIONAL_RATE_LIMITING` | Create project key with optional rate limiting | write | `organization_id_or_slug`, `project_id_or_slug`; *name, rateLimit__count, rateLimit__window* |  | `project:admin`, `project:write` |
| `SENTRY_CREATE_PROJECT_RULE_FOR_ALERTS` | Create project rule for alerts | write | `actionMatch`, `actions`, `conditions`, `frequency`, `name`, `organization_id_or_slug`, `project_id_or_slug`; *environment, filterMatch, filters, owner* |  | `alerts:write`, `project:admin`, `project:write` |
| `SENTRY_CREATE_PROJECT_WEBHOOK_SUBSCRIPTION` | Create project webhook subscription | write | `events`, `organization_id_or_slug`, `project_id_or_slug`, `url` |  | `project:write` |
| `SENTRY_CREATE_RELEASE_DEPLOY_FOR_ORG` | Create release deploy for org | write | `environment`, `organization_id_or_slug`, `version`; *dateFinished, dateStarted, name, projects, url* |  | `org:ci`, `project:admin`, `project:releases`, `project:write` |
| `SENTRY_CREATE_RELEASE_FOR_ORGANIZATION` | Create release for organization | write | `organization_id_or_slug`, `projects`, `version`; *commits, dateReleased, ref, refs, url* |  | `project:releases` |
| `SENTRY_CREATE_SCIM_GROUP_FOR_ORGANIZATION` | Create scim group for organization | write | `displayName`, `organization_id_or_slug` |  | `team:admin`, `team:write` |
| `SENTRY_CREATE_SENTRY_EXTERNAL_ISSUE_LINK` | Create sentry external issue link | write | `identifier`, `issueId`, `project`, `uuid`, `webUrl` |  | `event:write` |
| `SENTRY_CREATE_TEAM_PROJECT_FOR_ORGANIZATION` | Create team project for organization | write | `name`, `organization_id_or_slug`, `team_id_or_slug`; *default_rules, platform, slug* |  | `project:admin`, `project:write` |
| `SENTRY_CREATE_USER_FOR_SAML_INTEGRATION` | Create user for SAML integration | write | `organization_id_or_slug`, `userName`; *sentryOrgRole* |  | `member:admin`, `member:write` |
| `SENTRY_DELETE_AN_ALERT` | Delete an alert | write | `organization_id_or_slug`, `workflow_id` |  | `alerts:write`, `org:admin`, `org:write` |
| `SENTRY_DELETE_A_CUSTOM_INTEGRATION` | Delete a custom integration | write | `sentry_app_id_or_slug` |  | `org:admin` |
| `SENTRY_DELETE_BULK_ALERTS` | Delete bulk alerts | write | `organization_id_or_slug`; *id, project, query* |  | `alerts:write`, `org:admin`, `org:write` |
| `SENTRY_DELETE_DSYMS_FOR_PROJECT` | Delete dsyms for project | write | `id`, `organization_id_or_slug`, `project_id_or_slug` |  | `project:admin`, `project:releases` |
| `SENTRY_DELETE_EXTERNAL_ISSUE_BY_UUID` | Delete external issue by uuid | write | `external_issue_id`, `uuid` |  | `event:admin` |
| `SENTRY_DELETE_EXTERNAL_TEAM_BY_ID` | Delete an external team by id | write | `external_team_id`, `organization_id_or_slug`, `team_id_or_slug` |  | `team:admin` |
| `SENTRY_DELETE_EXTERNAL_USER_FROM_ORGANIZATION` | Delete external user from organization | write | `external_user_id`, `organization_id_or_slug` |  | `org:admin`, `org:write` |
| `SENTRY_DELETE_MEMBER_FROM_TEAM` | Delete member from team | write | `member_id`, `organization_id_or_slug`, `team_id_or_slug` |  | `org:admin`, `org:read`, `org:write`, `team:admin` |
| `SENTRY_DELETE_ORGANIZATION_ALERT_RULE` | Delete organization alert rule | write | `alert_rule_id`, `organization_id_or_slug` |  | `alerts:write`, `org:admin`, `org:write` |
| `SENTRY_DELETE_ORGANIZATION_DASHBOARD` | Delete organization dashboard | write | `dashboard_id`, `organization_id_or_slug` |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_DELETE_ORGANIZATION_DISCOVER_QUERY` | Delete organization discover query | write | `organization_id_or_slug`, `query_id` |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_DELETE_ORGANIZATION_INTEGRATION` | Delete organization integration | write | `integration_id`, `organization_id_or_slug` |  | `org:admin`, `org:integrations` |
| `SENTRY_DELETE_ORGANIZATION_ISSUE` | Delete organization issue | write | `issue_id`, `organization_id_or_slug` |  | `event:admin` |
| `SENTRY_DELETE_ORGANIZATION_MEMBER` | Delete organization member | write | `member_id`, `organization_id_or_slug` |  | `member:admin`, `member:read`, `member:write` |
| `SENTRY_DELETE_ORGANIZATION_MONITOR` | Delete organization monitor | write | `monitor_id_or_slug`, `organization_id_or_slug`; *environment* |  | `alerts:write`, `project:admin`, `project:write` |
| `SENTRY_DELETE_ORGANIZATION_RELEASE` | Delete organization release | write | `organization_id_or_slug`, `version` |  | `project:admin`, `project:releases` |
| `SENTRY_DELETE_ORGANIZATION_RELEASE_FILE` | Delete organization release file | write | `file_id`, `organization_id_or_slug`, `version` |  | `project:releases` |
| `SENTRY_DELETE_ORG_NOTIFICATION_ACTION` | Delete org notification action | write | `action_id`, `organization_id_or_slug` |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_DELETE_PROJECT_BY_ID` | Delete project by id | write | `organization_id_or_slug`, `project_id_or_slug` |  | `project:admin` |
| `SENTRY_DELETE_PROJECT_HOOK` | Delete project hook | write | `hook_id`, `organization_id_or_slug`, `project_id_or_slug` |  | `project:admin` |
| `SENTRY_DELETE_PROJECT_ISSUES` | Delete project issues | write | `organization_id_or_slug`, `project_id_or_slug`; *id* |  | `event:admin` |
| `SENTRY_DELETE_PROJECT_KEY` | Delete project key | write | `key_id`, `organization_id_or_slug`, `project_id_or_slug` |  | `project:admin` |
| `SENTRY_DELETE_PROJECT_MONITOR` | Delete project monitor | write | `monitor_id_or_slug`, `organization_id_or_slug`, `project_id_or_slug`; *environment* |  | `alerts:write`, `project:admin`, `project:write` |
| `SENTRY_DELETE_PROJECT_REPLAY` | Delete project replay | write | `organization_id_or_slug`, `project_id_or_slug`, `replay_id` |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_DELETE_PROJECT_RULE` | Delete project rule | write | `organization_id_or_slug`, `project_id_or_slug`, `rule_id` |  | `alerts:write`, `project:admin`, `project:write` |
| `SENTRY_DELETE_PROJECT_TEAM_ASSOCIATION` | Delete project team association | write | `organization_id_or_slug`, `project_id_or_slug`, `team_id_or_slug` |  | `project:admin`, `project:write` |
| `SENTRY_DELETE_RELEASE_FILE_BY_ID` | Delete release file by id | write | `file_id`, `organization_id_or_slug`, `project_id_or_slug`, `version` |  | `project:releases` |
| `SENTRY_DELETE_SYMBOL_SOURCE_FROM_PROJECT` | Delete a symbol source from a project | write | `id`, `organization_id_or_slug`, `project_id_or_slug` |  | `project:admin` |
| `SENTRY_DELETE_TEAM_BY_ORGANIZATION_OR_TEAM_SLUG` | Delete team by organization and team slug | write | `organization_id_or_slug`, `team_id_or_slug` |  | `team:admin` |
| `SENTRY_DELETE_TEAM_FROM_ORG_SCIM_V2` | Delete a team from an organization (SCIM v2) | write | `organization_id_or_slug`, `team_id` |  |  |
| `SENTRY_DELETE_USER_EMAILS_BY_ID` | Delete user emails by id | write | `email`, `user_id` |  |  |
| `SENTRY_DELETE_USER_FROM_ORG` | Delete user from org | write | `member_id`, `organization_id_or_slug` |  | `member:admin` |
| `SENTRY_DISABLE_SPIKE_PROTECTION` | Disable spike protection | write | `organization_id_or_slug`, `projects` |  | `project:read`, `project:write`, `project:admin` |
| `SENTRY_FETCH_ALERTS` | Fetch alerts | read | `organization_id_or_slug`; *id, project, query, sortBy* |  | `alerts:read`, `org:admin`, `org:read`, `org:write` |
| `SENTRY_FETCH_AN_ALERT` | Fetch an alert | read | `organization_id_or_slug`, `workflow_id` |  | `org:read`, `org:write`, `org:admin`, `alerts:read` |
| `SENTRY_FETCH_AN_ORGANIZATIONS_MONITORS` | Fetch an organization's monitors | read | `organization_id_or_slug`; *id, project, query, sortBy* |  | `alerts:read`, `org:read`, `org:write`, `org:admin` |
| `SENTRY_FETCH_A_MONITOR` | Fetch a monitor | read | `detector_id`, `organization_id_or_slug` |  | `alerts:read`, `alerts:write`, `project:admin`, `project:read`, `project:write`, `org:read`, `org:write`, `org:admin` |
| `SENTRY_FETCH_ISSUE_EVENT_BY_ID` | Fetch issue event by id | read | `event_id`, `issue_id`, `organization_id_or_slug`; *environment* |  | `event:admin`, `event:read`, `event:write` |
| `SENTRY_FETCH_ORGANIZATION_ALERT_RULES` | Fetch organization alert rules | read | `organization_id_or_slug` |  | `alerts:read`, `org:admin`, `org:read`, `org:write` |
| `SENTRY_FETCH_ORGANIZATION_RELEASE_THRESHOLD_STATUSES` | Fetch organization release threshold statuses | read | `end`, `organization_id_or_slug`, `start`; *environment, projectSlug, release* |  | `org:ci`, `project:admin`, `project:read`, `project:releases`, `project:write` |
| `SENTRY_FETCH_ORGANIZATION_REPLAY_COUNT` | Fetch organization replay count | read | `organization_id_or_slug`; *end, environment, project, query, start, statsPeriod* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_FETCH_ORGANIZATION_REPLAY_DETAILS` | Fetch organization replay details | read | `organization_id_or_slug`, `replay_id`; *cursor, end, environment, field, per_page, project, query, sort and 2 more* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_FETCH_ORGANIZATION_REPLAY_SELECTORS` | Fetch organization replay selectors | read | `organization_id_or_slug`; *cursor, end, environment, per_page, project, query, sort, start and 1 more* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_FETCH_PROJECT_ENVIRONMENT_DETAILS` | Fetch project environment details | read | `environment`, `organization_id_or_slug`, `project_id_or_slug` |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_FETCH_PROJECT_OWNERSHIP_DETAILS` | Fetch project ownership details | read | `organization_id_or_slug`, `project_id_or_slug` |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_FETCH_PROJECT_RELEASE_FILES` | Fetch project release files | read | `organization_id_or_slug`, `project_id_or_slug`, `version` |  | `project:releases` |
| `SENTRY_FETCH_PROJECT_REPLAY_CLICKS` | Fetch project replay clicks | read | `organization_id_or_slug`, `project_id_or_slug`, `replay_id`; *cursor, environment, per_page, query* |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_FETCH_REPLAY_RECORDING_SEGMENT` | Fetch replay recording segment | read | `organization_id_or_slug`, `project_id_or_slug`, `replay_id`, `segment_id` |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_FETCH_TAG_VALUES_FOR_ISSUE` | Fetch Tag Values for Issue | read | `issue_id`, `key`, `organization_id_or_slug`; *environment, sort* |  | `event:admin`, `event:read`, `event:write` |
| `SENTRY_GET_ALERTS` | Get alerts and monitors | read | `organization`; *expand, sort, team* |  | `project:admin`, `project:write` |
| `SENTRY_GET_INTEGRATION_DETAILS_BY_ORG` | Get integration details by org | read | `integration_id`, `organization_id_or_slug` |  | `org:admin`, `org:integrations`, `org:read`, `org:write` |
| `SENTRY_GET_OLDEST_EVENT` | Get oldest event | read | `issue_id`, `organization_id_or_slug`; *environment* |  | `event:admin`, `event:read`, `event:write` |
| `SENTRY_GET_ORGANIZATION_BY_ID_OR_SLUG` | Get organization by id or slug | read | `organization_id_or_slug`; *detailed* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_GET_ORGANIZATION_DETAILS` | Get organization details | read | *cursor, owner, query, sortBy* |  | `org:read`, `org:write`, `org:admin` |
| `SENTRY_GET_ORGANIZATION_ENVIRONMENTS` | Get organization environments | read | `organization_id_or_slug`; *visibility* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_GET_ORGANIZATION_ISSUE_DETAILS` | Get organization issue details | read | `issue_id`, `organization_id_or_slug` |  | `event:read` |
| `SENTRY_GET_ORGANIZATION_RELEASE_VERSION` | Get organization release version | read | `organization_id_or_slug`, `version`; *adoptionStages, health, healthStatsPeriod, project_id, query, sort, status, summaryStatsPeriod* |  | `org:ci`, `project:admin`, `project:read`, `project:releases`, `project:write` |
| `SENTRY_GET_ORGANIZATION_SESSIONS` | Get organization sessions | read | `field`, `organization_id_or_slug`; *end, environment, groupBy, includeSeries, includeTotals, interval, orderBy, per_page and 4 more* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_GET_ORGANIZATION_STATS_SUMMARY` | Get organization stats summary | read | `field`, `organization_id_or_slug`; *category, download, end, interval, outcome, project, reason, start and 1 more* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_GET_PROFILING_FLAMEGRAPH` | Get profiling flamegraph | read | `organization_id_or_slug`; *data_source, end, environment, expand, fingerprint, project, query, start and 1 more* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_GET_PROJECT_EVENTS` | Get project events | read | `organization_id_or_slug`, `project_id_or_slug`; *cursor, end, full, sample, start, statsPeriod* |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_GET_PROJECT_EVENT_STATS` | Get project event stats | read | `organization_id_or_slug`, `project_id_or_slug`; *resolution, since, stat, until* |  | `project:read` |
| `SENTRY_GET_PROJECT_HOOK` | Get project hook | read | `hook_id`, `organization_id_or_slug`, `project_id_or_slug` |  | `project:read` |
| `SENTRY_GET_PROJECT_LIST` | Get project list | read | *cursor* |  |  |
| `SENTRY_GET_PROJECT_MONITOR_BY_ID` | Get project monitor by id | read | `monitor_id_or_slug`, `organization_id_or_slug`, `project_id_or_slug` |  | `alerts:read`, `alerts:write`, `project:admin`, `project:read`, `project:write` |
| `SENTRY_GET_PROJECT_RULE_DETAILS` | Get project rule details | read | `organization_id_or_slug`, `project_id_or_slug`, `rule_id` |  | `alerts:read`, `alerts:write`, `project:admin`, `project:read`, `project:write` |
| `SENTRY_GET_SCIM_GROUP_BY_TEAM_ID` | Get SCIM group by team ID | read | `organization_id_or_slug`, `team_id` |  |  |
| `SENTRY_GET_SENTRY_APP_INSTALLATIONS_FOR_ORGANIZATION` | Get sentry app installations for organization | read | `organization_id_or_slug` |  | `org:integrations`, `org:read` |
| `SENTRY_GET_TEAM_MEMBERS_BY_ID_OR_SLUG` | Get team members by id or slug | read | `organization_id_or_slug`, `team_id_or_slug`; *cursor* |  | `team:admin`, `team:read`, `team:write` |
| `SENTRY_GET_WHO_VIEWED_REPLAY_BY_PROJECT` | Get who viewed replay by project | read | `organization_id_or_slug`, `project_id_or_slug`, `replay_id` |  | `event:admin`, `event:read`, `event:write` |
| `SENTRY_INGEST_EVENT_VIA_DSN` | Ingest event via DSN | write | `dsn`; *environment, event_id, exception_type, exception_value, extra, fingerprint, level, logger and 8 more* |  |  |
| `SENTRY_LIST_AN_ORGANIZATIONS_ISSUES` | List an organization's issues | read | `organization_id_or_slug`; *collapse, cursor, end, environment, expand, groupStatsPeriod, limit, project and 6 more* |  | `event:read`, `event:write`, `event:admin` |
| `SENTRY_LIST_AN_ORGANIZATION_S_CLIENT_KEYS` | List an organization's client keys | read | `organization_id_or_slug`; *cursor, status, team* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_LIST_A_PROJECTS_USER_FEEDBACK` | List a project's user feedback | read | `organization_id_or_slug`, `project_id_or_slug`; *cursor, per_page* |  | `project:read` |
| `SENTRY_LIST_A_TAGS_VALUES_FOR_AN_ISSUE` | List a tag's values for an issue | read | `issue_id`, `key`, `organization_id_or_slug`; *environment* |  | `event:admin`, `event:read`, `event:write` |
| `SENTRY_LIST_ORGANIZATIONS` | List organizations | read | *cursor, owner, query, sortBy* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_LIST_ORGANIZATION_DASHBOARDS` | List organization dashboards | read | `organization_id_or_slug`; *cursor, per_page* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_LIST_ORGANIZATION_MEMBERS` | List organization members | read | `organization_id_or_slug` |  | `member:admin`, `member:read`, `member:write` |
| `SENTRY_LIST_ORGANIZATION_RELEASES` | List organization releases | read | `organization_id_or_slug`; *query* |  | `project:releases` |
| `SENTRY_LIST_ORGANIZATION_REPOSITORIES` | List organization repositories | read | `organization_id_or_slug` |  | `org:read` |
| `SENTRY_LIST_ORGANIZATION_USER_TEAMS` | List User's Teams in Organization | read | `organization_id_or_slug` |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_LIST_PROJECT_TAGS` | List project tags | read | `organization_id_or_slug`, `project_id_or_slug` |  | `project:read` |
| `SENTRY_LIST_PROJECT_USERS` | List project users | read | `organization_id_or_slug`, `project_id_or_slug`; *query* |  | `project:read` |
| `SENTRY_LIST_SCIM_V2_ORGANIZATION_USERS` | List SCIM v2 organization users | read | `organization_id_or_slug`; *count, excludedAttributes, filter, startIndex* |  | `member:admin`, `member:read`, `member:write` |
| `SENTRY_LIST_SEER_AI_MODELS` | List Seer AI models | read |  |  |  |
| `SENTRY_LIST_TEAMS_IN_ORGANIZATION` | List teams in organization | read | `organization_id_or_slug`; *cursor, detailed* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_MANAGE_TEAM_EXTERNAL_INTEGRATIONS` | Create an external team | write | `external_name`, `integration_id`, `organization_id_or_slug`, `provider`, `team_id_or_slug`; *external_id* |  | `team:admin`, `team:write` |
| `SENTRY_MODIFY_ORGANIZATION_NOTIFICATION_ACTION` | Modify organization notification action | write | `action_id`, `organization_id_or_slug`, `projects`, `service_type`, `trigger_type`; *integration_id, target_display, target_identifier* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_MODIFY_ORGANIZATION_SETTINGS` | Modify organization settings | write | `organization_id_or_slug`; *alertsMemberWrite, allowJoinRequests, allowSharedIssues, attachmentsRole, avatar, avatarType, cancelDeletion, codecovAccess and 24 more* |  | `org:admin`, `org:write` |
| `SENTRY_MODIFY_RELEASE_FILE_ATTRIBUTES` | Modify release file attributes | write | `file_id`, `organization_id_or_slug`, `version`; *dist, name* |  | `project:releases` |
| `SENTRY_MUTATE_AN_ORGANIZATIONS_MONITORS` | Mutate an organization's monitors | write | `enabled`, `organization_id_or_slug`; *id, project, query* |  | `alerts:write`, `org:admin`, `org:read`, `org:write` |
| `SENTRY_MUTATE_AN_ORGANIZATION_S_ALERTS` | Mutate an organization's alerts | write | `enabled`, `organization_id_or_slug`; *id, project, query* |  | `alerts:write`, `org:admin`, `org:write` |
| `SENTRY_PATCH_SCIM_GROUP_OPERATIONS` | Patch SCIM group operations | write | `Operations`, `organization_id_or_slug`, `team_id` |  |  |
| `SENTRY_PATCH_USER_ACTIVE_STATUS_IN_ORGANIZATION` | Deactivate an organization member | write | `Operations`, `member_id`, `organization_id_or_slug` |  | `member:admin`, `member:write` |
| `SENTRY_POST_SPIKE_PROTECTION_FOR_ORGANIZATION` | Enable spike protection for an organization | write | `organization_id_or_slug`, `projects` |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_QUERY_EXPLORE_EVENTS_IN_TIMESERIES_FORMAT` | Query explore events in timeseries format | read | `dataset`, `organization_id_or_slug`; *comparisonDelta, disableAggregateExtrapolation, end, environment, excludeOther, groupBy, interval, preventMetricAggregates and 7 more* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_AGGREGATED_TEST_RESULT_METRICS` | Retrieve aggregated test result metrics | read | `organization_id_or_slug`, `owner`, `repository`; *branch, interval* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_ALERT_RULE_DETAILS` | Retrieve alert rule details | read | `alert_rule_id`, `organization_id_or_slug` |  | `alerts:read`, `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_A_CUSTOM_INTEGRATION_BY_ID_OR_SLUG` | Retrieve a custom integration by id or slug | write | `sentry_app_id_or_slug` |  | `event:admin`, `event:read`, `event:write`, `member:read`, `org:read`, `project:read`, `project:releases`, `team:read` |
| `SENTRY_RETRIEVE_A_MONITOR` | Retrieve a monitor | read | `monitor_id_or_slug`, `organization_id_or_slug`; *environment* |  | `alerts:read`, `alerts:write`, `project:admin`, `project:read`, `project:write` |
| `SENTRY_RETRIEVE_COMMITS_FOR_ORGANIZATION_REPO` | Retrieve an organization's repository commits | read | `organization_id_or_slug`, `repo_id` |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_COMMIT_FILES_FOR_RELEASE` | Retrieve commit files for release | read | `organization_id_or_slug`, `version` |  | `project:releases` |
| `SENTRY_RETRIEVE_CUSTOM_INTEG_ISSUE_LINKS_GIVEN_SENTRY_ISSUE` | Retrieve custom integration issue links | write | `issue_id`, `organization_id_or_slug` |  | `event:admin`, `event:read`, `event:write` |
| `SENTRY_RETRIEVE_DATA_FORWARDERS_FOR_AN_ORGANIZATION` | Retrieve data forwarders for an organization | read | `organization_id_or_slug` |  | `org:read` |
| `SENTRY_RETRIEVE_DSYM_FILES_FOR_PROJECT` | Retrieve dsym files for project | read | `organization_id_or_slug`, `project_id_or_slug` |  | `project:read` |
| `SENTRY_RETRIEVE_EVENT_ID_FOR_ORGANIZATION` | Retrieve event id for organization | read | `event_id`, `organization_id_or_slug` |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_FILES_FOR_RELEASE` | Retrieve files for release | read | `organization_id_or_slug`, `version` |  | `project:releases` |
| `SENTRY_RETRIEVE_ISSUE_EVENTS_BY_ID` | Retrieve issue events by id | read | `issue_id`, `organization_id_or_slug`; *end, environment, full, query, sample, start, statsPeriod* |  | `event:admin`, `event:read`, `event:write` |
| `SENTRY_RETRIEVE_ISSUE_HASHES_FOR_ORGANIZATION` | Retrieve issue hashes for organization | read | `issue_id`, `organization_id_or_slug`; *cursor, full* |  | `event:read` |
| `SENTRY_RETRIEVE_ISSUE_TAG_DETAILS` | Retrieve issue tag details | read | `issue_id`, `key`; *environment* |  |  |
| `SENTRY_RETRIEVE_MONITOR_CHECKINS` | Retrieve monitor checkins | read | `monitor_id_or_slug`, `organization_id_or_slug`, `project_id_or_slug` |  | `alerts:read`, `alerts:write`, `project:admin`, `project:read`, `project:write` |
| `SENTRY_RETRIEVE_MONITOR_CHECKINS_BY_ORG` | Retrieve monitor checkins by org | read | `monitor_id_or_slug`, `organization_id_or_slug` |  | `alerts:read`, `alerts:write`, `project:admin`, `project:read`, `project:write` |
| `SENTRY_RETRIEVE_NOTIFICATION_ACTION_BY_ORG_ID` | Retrieve notification action by org id | read | `action_id`, `organization_id_or_slug` |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_ORGANIZATION_DASHBOARD` | Retrieve organization dashboard | read | `dashboard_id`, `organization_id_or_slug` |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_ORGANIZATION_EVENTS` | Retrieve organization events | read | `field`, `organization_id_or_slug`; *end, environment, per_page, project, query, sort, start, statsPeriod* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_ORGANIZATION_INTEGRATIONS_LIST` | Retrieve organization integrations list | read | `organization_id_or_slug`; *features, includeConfig, providerKey* |  | `org:admin`, `org:integrations`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_ORGANIZATION_MEMBER` | Retrieve organization member | read | `member_id`, `organization_id_or_slug` |  | `member:admin`, `member:read`, `member:write` |
| `SENTRY_RETRIEVE_ORGANIZATION_MONITORS` | Retrieve organization monitors | read | `organization_id_or_slug`; *environment, owner, project* |  | `alerts:read`, `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_ORGANIZATION_PROJECTS` | Retrieve organization projects | read | `organization_id_or_slug`; *cursor* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_ORGANIZATION_RELAY_USAGE` | Retrieve organization relay usage | read | `organization_id_or_slug` |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_ORGANIZATION_RELEASE_COMMITS` | Retrieve organization release commits | read | `organization_id_or_slug`, `version` |  | `project:releases` |
| `SENTRY_RETRIEVE_ORGANIZATION_REPLAYS` | Retrieve organization replays | read | `organization_id_or_slug`; *cursor, end, environment, field, per_page, project, query, sort and 2 more* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_ORGANIZATION_SCIM_GROUPS` | Retrieve organization SCIM groups | read | `organization_id_or_slug`; *count, excludedAttributes, filter, startIndex* |  | `team:admin`, `team:read`, `team:write` |
| `SENTRY_RETRIEVE_ORGANIZATION_STATS_V2` | Retrieve organization stats v2 | read | `field`, `groupBy`, `organization_id_or_slug`; *category, end, interval, outcome, project, reason, start, statsPeriod* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_ORG_INTEGRATION_CONFIG` | Retrieve org integration config | read | `organization_id_or_slug`; *providerKey* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_PAGINATED_REPOSITORY_TOKENS` | Retrieve paginated repository tokens | read | `organization_id_or_slug`, `owner`; *cursor, limit, navigation, sortBy* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_PROJECT_ENVIRONMENTS` | Retrieve project environments | read | `organization_id_or_slug`, `project_id_or_slug`; *visibility* |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_RETRIEVE_PROJECT_EVENT_BY_ID` | Retrieve project event by id | read | `event_id`, `organization_id_or_slug`, `project_id_or_slug` |  | `project:read` |
| `SENTRY_RETRIEVE_PROJECT_FILTER_DATA` | Retrieve project filter data | read | `organization_id_or_slug`, `project_id_or_slug` |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_RETRIEVE_PROJECT_HOOKS` | Retrieve project hooks | read | `organization_id_or_slug`, `project_id_or_slug`; *cursor* |  | `project:read` |
| `SENTRY_RETRIEVE_PROJECT_ISSUES_LIST` | Retrieve project issues list | read | `organization_id_or_slug`, `project_id_or_slug`; *cursor, hashes, query, shortIdLookup, statsPeriod* |  | `event:read` |
| `SENTRY_RETRIEVE_PROJECT_KEYS_BY_ORG_AND_PROJECT` | List a project's client keys | read | `organization_id_or_slug`, `project_id_or_slug`; *cursor, status* |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_RETRIEVE_PROJECT_KEY_DETAILS` | Retrieve project key details | read | `key_id`, `organization_id_or_slug`, `project_id_or_slug` |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_RETRIEVE_PROJECT_MEMBERS_LIST` | Retrieve project members list | read | `organization_id_or_slug`, `project_id_or_slug` |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_RETRIEVE_PROJECT_RELEASE_FILE_DETAILS` | Retrieve project release file details | read | `file_id`, `organization_id_or_slug`, `project_id_or_slug`, `version`; *download* |  | `project:releases` |
| `SENTRY_RETRIEVE_PROJECT_RULES_BY_ORG_AND_PROJECT_ID` | Retrieve project rules by org and project id | read | `organization_id_or_slug`, `project_id_or_slug` |  | `alerts:read`, `alerts:write`, `project:admin`, `project:read`, `project:write` |
| `SENTRY_RETRIEVE_PROJECT_SYMBOL_SOURCES` | Retrieve project symbol sources | read | `organization_id_or_slug`, `project_id_or_slug`; *id* |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_RETRIEVE_PROJECT_TAG_VALUES` | Retrieve project tag values | read | `key`, `organization_id_or_slug`, `project_id_or_slug` |  | `project:read` |
| `SENTRY_RETRIEVE_PROJECT_TEAMS` | Retrieve project teams | read | `organization_id_or_slug`, `project_id_or_slug` |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_RETRIEVE_RELEASE_COMMITS` | Retrieve release commits | read | `organization_id_or_slug`, `project_id_or_slug`, `version` |  | `project:releases` |
| `SENTRY_RETRIEVE_RELEASE_DEPLOYMENTS` | Retrieve release deployments | read | `organization_id_or_slug`, `version` |  | `org:ci`, `project:admin`, `project:read`, `project:releases`, `project:write` |
| `SENTRY_RETRIEVE_RELEASE_FILE_BY_ID` | Retrieve release file by id | read | `file_id`, `organization_id_or_slug`, `version`; *download* |  | `project:releases` |
| `SENTRY_RETRIEVE_REPLAY_RECORDING_SEGMENTS` | Retrieve replay recording segments | read | `organization_id_or_slug`, `project_id_or_slug`, `replay_id`; *cursor, per_page* |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_RETRIEVE_SAVED_DISCOVER_QUERIES` | Retrieve saved discover queries | read | `organization_id_or_slug`; *cursor, per_page, query, sortBy* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_SAVED_DISCOVER_QUERY_FOR_ORGANIZATION` | Retrieve saved discover query for organization | read | `organization_id_or_slug`, `query_id` |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_SEER_ISSUE_FIX_STATE` | Retrieve Seer issue fix state | read | `issue_id`, `organization_id_or_slug` |  | `event:read`, `event:write`, `event:admin` |
| `SENTRY_RETRIEVE_SHORT_ID_FOR_ORGANIZATION` | Retrieve short id for organization | read | `organization_id_or_slug`, `short_id` |  |  |
| `SENTRY_RETRIEVE_SOURCE_MAP_DEBUG_EVENT` | Retrieve source map debug event | read | `event_id`, `exception_idx`, `frame_idx`, `organization_id_or_slug`, `project_id_or_slug` |  |  |
| `SENTRY_RETRIEVE_TEAM_INFO_VIA_ORGANIZATION_ID_OR_SLUG` | Retrieve team info | read | `organization_id_or_slug`, `team_id_or_slug`; *collapse, expand* |  | `team:admin`, `team:read`, `team:write` |
| `SENTRY_RETRIEVE_TEAM_PROJECTS` | Retrieve team projects | read | `organization_id_or_slug`, `team_id_or_slug`; *cursor* |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_RETRIEVE_THE_CUSTOM_INTEGS_CREATED_BY_AN_ORG` | Retrieve the custom integrations created by an organization | write | `organization_id_or_slug` |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_RETRIEVE_USER_EMAIL_INFORMATION` | Retrieve user email information | read | `user_id` |  |  |
| `SENTRY_RETRIEVE_USER_VIA_SCIM_API` | Retrieve user via SCIM API | read | `member_id`, `organization_id_or_slug` |  | `member:admin`, `member:read`, `member:write` |
| `SENTRY_SAVE_ORGANIZATION_DISCOVER_QUERY` | Save organization discover query | write | `name`, `organization_id_or_slug`; *display, end, environment, fields, interval, orderby, projects, query and 5 more* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_START_SEER_ISSUE_FIX` | Start Seer issue fix | write | `issue_id`, `organization_id_or_slug`; *event_id, instruction, pr_to_comment_on_url, stopping_point* |  | `event:write`, `event:admin` |
| `SENTRY_SUBMIT_NOTIFICATION_ACTION_API_DATA` | Submit notification action API data | write | `organization_id_or_slug`, `service_type`, `trigger_type`; *integration_id, projects, target_display, target_identifier* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_SWITCH_TEAM_ROLE_FOR_MEMBER` | Switch team role for member | write | `member_id`, `organization_id_or_slug`, `team_id_or_slug`; *teamRole* |  | `member:admin`, `member:read`, `member:write`, `org:admin`, `org:read`, `org:write`, `team:admin`, `team:write` |
| `SENTRY_TOGGLE_PROJECT_FILTER_STATUS` | Toggle project filter status | write | `filter_id`, `organization_id_or_slug`, `project_id_or_slug`; *active, subfilters* |  | `project:admin`, `project:write` |
| `SENTRY_UPDATE_AN_ALERT_BY_ID` | Update an alert by ID | write | `name`, `organization_id_or_slug`, `workflow_id`; *action_filters, config, enabled, environment, id, triggers* |  | `alerts:write`, `org:admin`, `org:write` |
| `SENTRY_UPDATE_AN_EXISTING_CUSTOM_INTEGRATION` | Update an existing custom integration | write | `name`, `scopes`, `sentry_app_id_or_slug`; *allowedOrigins, author, events, isAlertable, isInternal, overview, redirectUrl, schema and 2 more* |  | `org:admin`, `org:write` |
| `SENTRY_UPDATE_A_MONITOR` | Update a monitor | write | `config`, `monitor_id_or_slug`, `name`, `organization_id_or_slug`, `project`; *is_muted, owner, slug, status* |  | `alerts:write`, `project:admin`, `project:write` |
| `SENTRY_UPDATE_A_MONITOR_BY_ID` | Update a monitor by ID | write | `detectorId`, `name`, `organizationIdOrSlug`, `type`; *conditionGroup, config, dataSources, description, enabled, owner* |  | `alerts:write`, `org:admin`, `org:read`, `org:write` |
| `SENTRY_UPDATE_A_PROJECTS_SYMBOL_SOURCE` | Update a symbol source in a project | write | `id`, `name`, `organization_id_or_slug`, `project_id_or_slug`, `type`; *access_key, bucket, client_email, filters__filetypes, filters__path_patterns, filters__requires_checksum, layout__casing, layout__type and 8 more* |  |  |
| `SENTRY_UPDATE_ENVIRONMENT_VISIBILITY` | Update environment visibility | write | `environment`, `isHidden`, `organization_id_or_slug`, `project_id_or_slug` |  | `project:admin`, `project:write` |
| `SENTRY_UPDATE_EXTERNAL_TEAM_INTEGRATION` | Update external team integration | write | `external_name`, `external_team_id`, `integration_id`, `organization_id_or_slug`, `provider`, `team_id_or_slug`; *external_id* |  | `team:admin`, `team:write` |
| `SENTRY_UPDATE_EXTERNAL_USER_FOR_ORGANIZATION` | Update external user for organization | write | `external_name`, `external_user_id`, `id`, `integration_id`, `organization_id_or_slug`, `provider`, `user_id`; *external_id* |  | `org:admin`, `org:write` |
| `SENTRY_UPDATE_ISSUE_ATTRIBUTES_IN_ORGANIZATION` | Update issue attributes in organization | write | `issue_id`, `organization_id_or_slug`; *assignedTo, hasSeen, isBookmarked, isPublic, isSubscribed, status, statusDetails__inCommit, statusDetails__inNextRelease and 1 more* |  | `event:write` |
| `SENTRY_UPDATE_MONITOR_FOR_PROJECT` | Update monitor for project | write | `config`, `monitor_id_or_slug`, `name`, `organization_id_or_slug`, `project`, `project_id_or_slug`; *is_muted, owner, slug, status* |  | `alerts:write`, `project:admin`, `project:write` |
| `SENTRY_UPDATE_ORGANIZATION_ALERT_RULES` | Update organization alert rules | write | `aggregate`, `alert_rule_id`, `name`, `organization_id_or_slug`, `projects`, `query`, `thresholdType`, `timeWindow`, `triggers`; *activationCondition, comparisonDelta, dataset, environment, eventTypes, monitorType, owner, queryType and 1 more* |  | `alerts:write`, `org:admin`, `org:write` |
| `SENTRY_UPDATE_ORGANIZATION_DASHBOARD` | Update organization dashboard | write | `dashboard_id`, `organization_id_or_slug`; *end, environment, filters, id, period, projects, start, title and 2 more* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_UPDATE_ORGANIZATION_MEMBER_ROLE` | Update organization member role | write | `member_id`, `organization_id_or_slug`; *orgRole, teamRoles* |  | `member:admin`, `member:invite`, `member:write` |
| `SENTRY_UPDATE_PROJECT_DETAILS` | Update project details | write | `organization_id_or_slug`, `project_id_or_slug`; *highlightContext, highlightTags, isBookmarked, name, platform, resolveAge, slug, subjectPrefix and 1 more* |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_UPDATE_PROJECT_ISSUE_STATUS_AND_DETAILS` | Update project issue status and details | write | `organization_id_or_slug`, `project_id_or_slug`; *assignedTo, hasSeen, id, ignoreDuration, isBookmarked, isPublic, merge, status and 8 more* |  | `event:write` |
| `SENTRY_UPDATE_PROJECT_KEY_CONFIGURATION` | Update project key configuration | write | `key_id`, `organization_id_or_slug`, `project_id_or_slug`; *browserSdkVersion, dynamicSdkLoaderOptions__hasDebug, dynamicSdkLoaderOptions__hasPerformance, dynamicSdkLoaderOptions__hasReplay, isActive, name, rateLimit__count, rateLimit__window* |  | `project:admin`, `project:write` |
| `SENTRY_UPDATE_PROJECT_OWNERSHIP_SETTINGS` | Update project ownership settings | write | `organization_id_or_slug`, `project_id_or_slug`; *autoAssignment, codeownersAutoSync, fallthrough, raw* |  | `project:admin`, `project:read`, `project:write` |
| `SENTRY_UPDATE_PROJECT_RELEASE_FILE_DETAILS` | Update project release file details | write | `file_id`, `organization_id_or_slug`, `project_id_or_slug`, `version`; *dist, name* |  | `project:releases` |
| `SENTRY_UPDATE_PROJECT_RULE_BY_ID` | Update Project Rule By Id | write | `actionMatch`, `actions`, `conditions`, `frequency`, `name`, `organization_id_or_slug`, `project_id_or_slug`, `rule_id`; *environment, filterMatch, filters, owner* |  | `alerts:write`, `project:admin`, `project:write` |
| `SENTRY_UPDATE_RELEASE_DETAILS_FOR_ORGANIZATION` | Update release details for organization | write | `organization_id_or_slug`, `version`; *commits, dateReleased, ref, refs, url* |  | `org:ci`, `project:admin`, `project:releases`, `project:write` |
| `SENTRY_UPDATE_SAVED_QUERY_FOR_ORGANIZATION` | Update a saved query for an organization | write | `name`, `organization_id_or_slug`, `query_id`; *display, end, environment, fields, interval, orderby, projects, query and 5 more* |  | `org:admin`, `org:read`, `org:write` |
| `SENTRY_UPDATE_TEAM_INFORMATION_BY_ORGANIZATION_ID` | Update team information by organization id | write | `organization_id_or_slug`, `slug`, `team_id_or_slug` |  | `team:admin`, `team:write` |
| `SENTRY_UPDATE_USER_EMAIL` | Update user email | write | `email`, `user_id` |  |  |
| `SENTRY_UPDATE_WEBHOOK_CONFIGURATION_XP` | Update webhook configuration | write | `events`, `hook_id`, `organization_id_or_slug`, `project_id_or_slug`, `url` |  | `project:write` |
| `SENTRY_UPLOAD_DSYMS_FILE_TO_PROJECT` | Upload dsyms file to project | write | `organization_id_or_slug`, `project_id_or_slug`; *file* |  | `project:write` |
| `SENTRY_UPLOAD_FILE_TO_PROJECT_RELEASE` | Upload file to project release | write | `organization_id_or_slug`, `project_id_or_slug`, `version`; *dist, file, header, name* |  | `project:releases` |
| `SENTRY_UPLOAD_RELEASE_FILE_TO_ORGANIZATION` | Upload release file to organization | write | `organization_id_or_slug`, `version`; *dist, file, header, name* |  | `project:releases` |
| `SENTRY_VALIDATE_CREDENTIAL` | Validate credential | read |  |  |  |
| `SENTRY_VIEW_ORGANIZATION_NOTIFICATION_ACTIONS` | View organization notification actions | read | `organization_id_or_slug`; *project, project_id_or_slug, triggerType* |  | `org:admin`, `org:read`, `org:write` |
