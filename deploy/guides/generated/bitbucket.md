# Bitbucket: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

116 tools: 86 read, 30 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `BITBUCKET_APPROVE_PULL_REQUEST` | Approve Pull Request | write | `pull_request_id`, `repo_slug`, `workspace` |  | `pullrequest` |
| `BITBUCKET_BROWSE_REPOSITORY_PATH` | Browse repository path | read | `repo_slug`, `workspace`; *commit, format, max_depth, page, pagelen, path, q, sort* |  | `repository` |
| `BITBUCKET_CREATE_BRANCH` | Create a branch | write | `name`, `repo_slug`, `target_hash`, `workspace` |  | `repository:write` |
| `BITBUCKET_CREATE_ISSUE` | Create an issue | write | `content`, `repo_slug`, `title`, `workspace`; *assignee, component_id, due_on, kind, milestone_id, priority, version_id* |  | `issue:write` |
| `BITBUCKET_CREATE_ISSUE_COMMENT` | Create an issue comment | write | `content`, `issue_id`, `repo_slug`, `workspace` |  | `issue:write` |
| `BITBUCKET_CREATE_PULL_REQUEST` | Create a pull request | write | `repo_slug`, `source_branch`, `title`, `workspace`; *close_source_branch, description, destination_branch, reviewers* |  | `pullrequest:write` |
| `BITBUCKET_CREATE_PULL_REQUEST_COMMENT` | Create pull request comment | write | `content_raw`, `pull_request_id`, `repo_slug`, `workspace`; *inline, parent_comment_id* |  | `pullrequest` |
| `BITBUCKET_CREATE_REPOSITORIES_COMMITS2` | List commits from revision (POST) | read | `repo_slug`, `revision`, `workspace`; *exclude, include, page, pagelen, path* |  | `repository` |
| `BITBUCKET_CREATE_REPOSITORIES_COMMIT_REPORTS_ANNOTATIONS` | Create commit report annotations | write | `annotations`, `commit`, `repo_slug`, `report_id`, `workspace` |  | `repository` |
| `BITBUCKET_CREATE_REPOSITORY` | Create repository | write | `repo_slug`, `workspace`; *description, fork_policy, has_issues, has_wiki, is_private, language, project_key* |  |  |
| `BITBUCKET_CREATE_SNIPPET_COMMENT` | Create snippet comment | write | `content`, `encoded_id`, `workspace`; *parent_id* |  | `snippet` |
| `BITBUCKET_CREATE_TEAMS_PIPELINES_CONFIG_VARIABLES` | Create team pipeline variable | write | `key`, `username`, `value`; *secured* |  | `pipeline:variable` |
| `BITBUCKET_CREATE_USERS_PIPELINES_CONFIG_VARIABLES` | Create user pipeline variable | write | `key`, `selected_user`, `value`; *secured* |  | `pipeline:variable` |
| `BITBUCKET_DECLINE_PULL_REQUEST` | Decline Pull Request | write | `pull_request_id`, `repo_slug`, `workspace` |  | `pullrequest:write`, `read:pullrequest:bitbucket`, `write:pullrequest:bitbucket` |
| `BITBUCKET_DELETE_COMMIT_COMMENT` | Delete commit comment | write | `comment_id`, `commit`, `repo_slug`, `workspace` |  | `repository` |
| `BITBUCKET_DELETE_ISSUE` | Delete issue | write | `issue_id`, `repo_slug`, `workspace` |  | `issue:write` |
| `BITBUCKET_DELETE_PULL_REQUEST_COMMENT` | Delete pull request comment | write | `comment_id`, `pull_request_id`, `repo_slug`, `workspace` |  | `pullrequest` |
| `BITBUCKET_DELETE_REPOSITORIES_COMMIT_REPORTS_ANNOTATIONS` | Delete commit report annotation | write | `annotation_id`, `commit`, `repo_slug`, `report_id`, `workspace` |  | `repository` |
| `BITBUCKET_DELETE_REPOSITORY` | Delete repository | write | `repo_slug`, `workspace` |  | `repository:delete` |
| `BITBUCKET_DELETE_SNIPPETS_WATCH` | Delete snippet watch | write | `encoded_id`, `workspace` |  | `snippet` |
| `BITBUCKET_DELETE_USER_PIPELINE_VARIABLE` | Delete user pipeline variable | write | `username`, `variable_uuid` |  | `pipeline:variable` |
| `BITBUCKET_DOWNLOAD_REPOSITORY_ARCHIVE` | Download repository archive | read | `repo_slug`, `revision`, `workspace`; *format* |  | `repository` |
| `BITBUCKET_GET_BRANCH` | Get branch | read | `name`, `repo_slug`, `workspace` |  | `repository` |
| `BITBUCKET_GET_COMMIT_BUILD_STATUS` | Get Commit Build Status | read | `commit`, `key`, `repo_slug`, `workspace` |  | `repository`, `read:repository:bitbucket` |
| `BITBUCKET_GET_COMMIT_CHANGES` | Get Commit Changes | read | `repo_slug`, `spec`, `workspace`; *page, pagelen* |  | `repository`, `read:repository:bitbucket` |
| `BITBUCKET_GET_COMMIT_COMMENT` | Get commit comment | read | `comment_id`, `commit`, `repo_slug`, `workspace` |  | `repository` |
| `BITBUCKET_GET_COMMIT_DIFF` | Get Commit Diff | read | `repo_slug`, `spec`, `workspace`; *binary, context, ignore_whitespace, merge, path, renames, topic* |  | `repository` |
| `BITBUCKET_GET_CURRENT_USER2` | Get current user (v2) | read |  |  | `account` |
| `BITBUCKET_GET_DEPLOYMENT_ENVIRONMENT_VARIABLES` | Get deployment environment variables | read | `environment_uuid`, `repo_slug`, `workspace` |  | `pipeline` |
| `BITBUCKET_GET_FILE_FROM_REPOSITORY` | Get file from repository | read | `commit`, `path`, `repo_slug`, `workspace` |  | `read` |
| `BITBUCKET_GET_HOOK_EVENTS` | Get hook events | read | `subject_type` |  |  |
| `BITBUCKET_GET_OPENID_CONFIGURATION` | Get OpenID configuration for OIDC in Pipelines | read | `workspace` |  |  |
| `BITBUCKET_GET_PIPELINE_STEP` | Get pipeline step | read | `pipeline_uuid`, `repo_slug`, `step_uuid`, `workspace` |  | `pipeline` |
| `BITBUCKET_GET_PIPELINE_STEP_LOG` | Get pipeline step log | read | `pipeline_uuid`, `repo_slug`, `step_uuid`, `workspace`; *range_end, range_start* |  | `pipeline` |
| `BITBUCKET_GET_PULL_REQUEST` | Get Pull Request | read | `pull_request_id`, `repo_slug`, `workspace` |  | `pullrequest` |
| `BITBUCKET_GET_PULL_REQUEST_ACTIVITY` | Get pull request activity | read | `pull_request_id`, `repo_slug`, `workspace`; *page, pagelen* |  | `pullrequest` |
| `BITBUCKET_GET_PULL_REQUEST_COMMENT` | Get pull request comment | read | `comment_id`, `pull_request_id`, `repo_slug`, `workspace` |  | `pullrequest` |
| `BITBUCKET_GET_PULL_REQUEST_COMMITS` | Get Pull Request Commits | read | `pull_request_id`, `repo_slug`, `workspace`; *pagelen, response_detail* |  | `pullrequest` |
| `BITBUCKET_GET_PULL_REQUEST_DIFF` | Get Pull Request Diff | read | `pull_request_id`, `repo_slug`, `workspace`; *max_chars* |  | `pullrequest` |
| `BITBUCKET_GET_PULL_REQUEST_DIFFSTAT` | Get Pull Request Diffstat | read | `pull_request_id`, `repo_slug`, `workspace`; *page, pagelen* |  | `pullrequest`, `read:pullrequest:bitbucket` |
| `BITBUCKET_GET_PULL_REQUEST_MERGE_BASE` | Get pull request merge base | read | `pull_request_id`, `repo_slug`, `workspace` |  | `pullrequest`, `repository` |
| `BITBUCKET_GET_REPOSITORIES_BRANCHING_MODEL` | Get Repository Branching Model | read | `repo_slug`, `workspace` |  | `repository` |
| `BITBUCKET_GET_REPOSITORIES_COMMIT` | Get Repository Commit | read | `commit`, `repo_slug`, `workspace` |  | `repository`, `read:repository:bitbucket` |
| `BITBUCKET_GET_REPOSITORIES_COMMIT_COMMENTS` | Get commit comments | read | `commit`, `repo_slug`, `workspace`; *page, pagelen* |  | `repository` |
| `BITBUCKET_GET_REPOSITORIES_COMMIT_REPORT` | Get commit report | read | `commit`, `repo_slug`, `report_id`, `workspace` |  | `repository` |
| `BITBUCKET_GET_REPOSITORIES_COMMIT_REPORTS` | Get Commit Reports | read | `commit`, `repo_slug`, `workspace` |  | `repository`, `read:repository:bitbucket` |
| `BITBUCKET_GET_REPOSITORIES_COMMIT_REPORTS_ANNOTATIONS` | Get commit report annotation | read | `annotation_id`, `commit`, `repo_slug`, `report_id`, `workspace` |  | `repository` |
| `BITBUCKET_GET_REPOSITORIES_COMMIT_STATUSES` | Get commit statuses | read | `commit`, `repo_slug`, `workspace` |  | `repository` |
| `BITBUCKET_GET_REPOSITORIES_EFFECTIVE_BRANCHING_MODEL` | Get effective branching model | read | `repo_slug`, `workspace` |  | `repository` |
| `BITBUCKET_GET_REPOSITORIES_ENVIRONMENTS2` | Get Repository Environment | read | `environment_uuid`, `repo_slug`, `workspace` |  | `pipeline`, `read:pipeline:bitbucket` |
| `BITBUCKET_GET_REPOSITORIES_FILEHISTORY` | Get file commit history | read | `commit`, `path`, `repo_slug`, `workspace`; *fields, page, pagelen, q, renames* |  | `repository`, `read:repository:bitbucket` |
| `BITBUCKET_GET_REPOSITORIES_ISSUES_VOTE` | Check if user voted for issue | read | `issue_id`, `repo_slug`, `workspace` |  | `issue` |
| `BITBUCKET_GET_REPOSITORIES_MERGE_BASE` | Get Repositories Merge Base | read | `repo_slug`, `revspec`, `workspace` |  | `repository` |
| `BITBUCKET_GET_REPOSITORIES_PIPELINES2` | Get specific pipeline | read | `pipeline_uuid`, `repo_slug`, `workspace` |  | `pipeline` |
| `BITBUCKET_GET_REPOSITORIES_PIPELINES_CONFIG_CACHES` | Get repository pipelines caches | read | `repo_slug`, `workspace` |  | `pipeline`, `read:pipeline:bitbucket` |
| `BITBUCKET_GET_REPOSITORIES_PIPELINES_CONFIG_RUNNERS` | Get repository pipeline runners | read | `repo_slug`, `workspace` |  | `runner` |
| `BITBUCKET_GET_REPOSITORIES_PIPELINES_CONFIG_SCHEDULES` | Get repository pipeline schedules | read | `repo_slug`, `workspace` |  | `pipeline`, `read:pipeline:bitbucket` |
| `BITBUCKET_GET_REPOSITORIES_PIPELINES_CONFIG_SSH_KNOWN_HOSTS` | Get repository SSH known hosts | read | `repo_slug`, `workspace` |  | `pipeline` |
| `BITBUCKET_GET_REPOSITORIES_PIPELINES_CONFIG_VARIABLES` | Get repository pipeline variables | read | `repo_slug`, `workspace`; *page, pagelen, q, sort* |  | `pipeline` |
| `BITBUCKET_GET_REPOSITORIES_PIPELINES_STEPS` | Get pipeline steps | read | `pipeline_uuid`, `repo_slug`, `workspace` |  | `pipeline`, `read:pipeline:bitbucket` |
| `BITBUCKET_GET_REPOSITORIES_PULLREQUESTS_ACTIVITY` | Get pull requests activity log | read | `repo_slug`, `workspace`; *ctx, page, pagelen* |  | `pullrequest` |
| `BITBUCKET_GET_REPOSITORIES_PULLREQUESTS_COMMENTS` | Get pull request comments | read | `pull_request_id`, `repo_slug`, `workspace`; *page, pagelen, sort* |  | `pullrequest` |
| `BITBUCKET_GET_REPOSITORIES_PULLREQUESTS_STATUSES` | Get pull request statuses | read | `pull_request_id`, `repo_slug`, `workspace`; *page, pagelen* |  | `pullrequest` |
| `BITBUCKET_GET_REPOSITORIES_REFS` | Get repository refs | read | `repo_slug`, `workspace`; *page, pagelen, q, sort* |  | `repository`, `read:repository:bitbucket` |
| `BITBUCKET_GET_REPOSITORIES_REFS_TAGS` | Get tag | read | `name`, `repo_slug`, `workspace` |  | `repository` |
| `BITBUCKET_GET_REPOSITORIES_SRC` | Get repositories src | read | `repo_slug`, `workspace`; *format* |  | `repository`, `read:repository:bitbucket` |
| `BITBUCKET_GET_REPOSITORIES_WATCHERS` | Get repository watchers | read | `repo_slug`, `workspace` |  | `repository` |
| `BITBUCKET_GET_REPOSITORY` | Get repository | read | `repo_slug`, `workspace` |  | `repository` |
| `BITBUCKET_GET_REPOSITORY_PATCH` | Get Repository Patch | read | `repo_slug`, `spec`, `workspace`; *max_chars* |  | `repository` |
| `BITBUCKET_GET_SNIPPET` | Get snippet | read | `encoded_id`, `workspace`; *accept_format* |  | `snippet` |
| `BITBUCKET_GET_SNIPPETS_WATCH` | Get snippet watch status | read | `encoded_id`, `workspace` |  | `snippet` |
| `BITBUCKET_GET_SSH_LATEST_KEYS` | Get SSH keys for user | read | `username`; *page, pagelen* |  | `account` |
| `BITBUCKET_GET_USER` | Get user | read | `selected_user` |  |  |
| `BITBUCKET_GET_USER_EMAILS` | Get user emails | read |  |  | `email`, `account` |
| `BITBUCKET_GET_USER_EMAILS2` | Get user email details | read | `email` |  | `email` |
| `BITBUCKET_GET_USER_PERMISSIONS_REPOSITORIES` | Get user permissions for repositories | read | *page, pagelen* |  | `account`, `repository` |
| `BITBUCKET_GET_USER_PERMISSIONS_WORKSPACES` | Get user permissions for workspaces | read | *page, pagelen, q* |  | `account` |
| `BITBUCKET_GET_USER_WORKSPACES` | Get user workspaces | read | *q* |  | `account`, `read:workspace:bitbucket`, `read:user:bitbucket` |
| `BITBUCKET_GET_WORKSPACE` | Get workspace | read | `workspace` |  |  |
| `BITBUCKET_GET_WORKSPACES_PULLREQUESTS` | Get Workspace Pull Requests by User | read | `selected_user`, `workspace`; *page, pagelen, state* |  | `pullrequest` |
| `BITBUCKET_GET_WORKSPACES_SEARCH_CODE` | Search code in workspace | read | `search_query`, `workspace`; *fields, page, pagelen* |  | `repository`, `read:repository:bitbucket` |
| `BITBUCKET_LIST_BRANCHES` | List branches | read | `repo_slug`, `workspace`; *page, pagelen, q, sort* |  | `repository` |
| `BITBUCKET_LIST_COMMITS` | List commits | read | `repo_slug`, `workspace`; *exclude, include, page, pagelen, path* |  | `repository` |
| `BITBUCKET_LIST_COMMITS_FROM_REVISION` | List commits from revision | read | `repo_slug`, `revision`, `workspace`; *exclude, include, page, pagelen, path* |  | `repository` |
| `BITBUCKET_LIST_COMMITS_ON_MASTER` | List commits on master | read | `repo_slug`, `workspace`; *page, pagelen* |  | `repository` |
| `BITBUCKET_LIST_DEPLOYMENTS` | List deployments | read | `repo_slug`, `workspace`; *page, pagelen, q, sort* |  | `pipeline` |
| `BITBUCKET_LIST_ISSUES` | List issues in a repository | read | `repo_slug`, `workspace`; *assignee, kind, page, pagelen, priority, state* |  | `issue`, `repository` |
| `BITBUCKET_LIST_PIPELINES` | List pipelines | read | `repo_slug`, `workspace`; *page, pagelen* |  | `pipeline` |
| `BITBUCKET_LIST_PIPELINE_STEP_TEST_CASES` | List pipeline step test cases | read | `pipeline_uuid`, `repo_slug`, `step_uuid`, `workspace`; *page, pagelen* |  | `pipeline` |
| `BITBUCKET_LIST_PULL_REQUESTS` | List pull requests | read | `repo_slug`, `workspace`; *page, pagelen, state* |  | `pullrequest` |
| `BITBUCKET_LIST_PULL_REQUESTS_FOR_COMMIT` | List pull requests for commit | read | `commit`, `repo_slug`, `workspace`; *page, pagelen* |  | `pullrequest` |
| `BITBUCKET_LIST_PULL_REQUEST_MERGEABILITY_CHECKS` | List pull request mergeability checks | read | `pull_request_id`, `repo_slug`, `workspace`; *exclude_git_mergeability, exclude_maximum_commits_behind* |  | `pullrequest` |
| `BITBUCKET_LIST_PULL_REQUEST_TASKS` | List pull request tasks | read | `pull_request_id`, `repo_slug`, `workspace`; *page, pagelen* |  | `pullrequest` |
| `BITBUCKET_LIST_REPOSITORIES` | List all public repositories | read | *after, page, pagelen, q, role, sort* |  | `repository`, `read:repository:bitbucket` |
| `BITBUCKET_LIST_REPOSITORIES_ENVIRONMENTS` | List repository environments | read | `repo_slug`, `workspace`; *page, pagelen, q, sort* |  | `pipeline` |
| `BITBUCKET_LIST_REPOSITORIES_IN_WORKSPACE` | List repositories in workspace | read | `workspace`; *page, pagelen, q, role, sort* |  | `repository` |
| `BITBUCKET_LIST_REPOSITORY_FORKS` | List repository forks | read | `repo_slug`, `workspace`; *page, pagelen* |  | `repository`, `read:repository:bitbucket` |
| `BITBUCKET_LIST_REPOSITORY_PATHS` | List repository paths | read | `commit`, `repo_slug`, `workspace`; *max_depth, page, pagelen, path, q, sort* |  | `repository` |
| `BITBUCKET_LIST_SNIPPETS` | List snippets | read | *page, pagelen, role* |  | `snippet` |
| `BITBUCKET_LIST_TAGS` | List tags | read | `repo_slug`, `workspace`; *page, pagelen, q, sort* |  | `repository` |
| `BITBUCKET_LIST_VERSIONS` | List versions | read | `repo_slug`, `workspace`; *page, pagelen, q, sort* |  | `issue` |
| `BITBUCKET_LIST_WORKSPACES` | List workspaces | read | *page, pagelen, q, sort* |  |  |
| `BITBUCKET_LIST_WORKSPACE_MEMBERS` | List workspace members | read | `workspace`; *page, pagelen* |  | `read:workspace` |
| `BITBUCKET_LIST_WORKSPACE_PROJECTS` | List workspace projects | read | `workspace`; *page, pagelen, q, sort* |  | `project` |
| `BITBUCKET_MERGE_PULL_REQUEST` | Merge Pull Request | write | `pull_request_id`, `repo_slug`, `workspace`; *close_source_branch, merge_strategy, message* |  | `pullrequest:write`, `read:pullrequest:bitbucket`, `write:pullrequest:bitbucket` |
| `BITBUCKET_REQUEST_PULL_REQUEST_CHANGES` | Request Pull Request Changes | write | `pull_request_id`, `repo_slug`, `workspace` |  | `pullrequest:write` |
| `BITBUCKET_RESOLVE_PULL_REQUEST_COMMENT` | Resolve pull request comment | write | `comment_id`, `pull_request_id`, `repo_slug`, `workspace`; *action* |  | `pullrequest` |
| `BITBUCKET_SEARCH_USER_REPOSITORIES_CODE` | Search code in user repositories | read | `search_query`, `selected_user`; *fields, page, pagelen* |  | `repository`, `project` |
| `BITBUCKET_STOP_PIPELINE` | Stop pipeline | write | `pipeline_uuid`, `repo_slug`, `workspace` |  | `pipeline:write` |
| `BITBUCKET_UPDATE_INSIGHTS_PROJECTS_REPOS_COMMITS_REPORTS` | Update commit insight report | write | `commit`, `details`, `repo_slug`, `report_id`, `report_type`, `title`, `workspace`; *data, link, logo_url, remote_link_enabled, reporter, result* |  | `repository` |
| `BITBUCKET_UPDATE_ISSUE` | Update an issue | write | `issue_id`, `repo_slug`, `workspace`; *assignee_account_id, component, content, kind, milestone, priority, state, title and 1 more* |  | `issue:write` |
| `BITBUCKET_UPDATE_PULL_REQUEST` | Update Pull Request | write | `pull_request_id`, `repo_slug`, `workspace`; *close_source_branch, description, draft, reviewers, title* |  | `pullrequest:write` |
| `BITBUCKET_UPDATE_REPOSITORIES_COMMIT_COMMENTS` | Update commit comment | write | `comment_id`, `commit`, `content_raw`, `repo_slug`, `workspace` |  | `repository` |
| `BITBUCKET_UPDATE_REPOSITORIES_COMMIT_REPORTS_ANNOTATIONS` | Update commit report annotation | write | `annotation_id`, `annotation_type`, `commit`, `repo_slug`, `report_id`, `severity`, `summary`, `workspace`; *details, external_id, line, link, path, result* |  | `repository` |
| `BITBUCKET_UPDATE_TEAMS_PIPELINES_CONFIG_VARIABLES` | Update team pipeline variable | write | `key`, `username`, `value`, `variable_uuid`; *secured* |  | `pipeline:variable` |
| `BITBUCKET_UPDATE_USERS_PIPELINES_CONFIG_VARIABLES` | Update user pipeline variable | write | `selected_user`, `variable_uuid`; *key, secured, value* |  | `pipeline:variable` |
