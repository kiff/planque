# GitHub: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

883 tools: 459 read, 424 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `GITHUB_ABORT_REPOSITORY_MIGRATION` | Abort Repository Migration | write | `migrationId`; *clientMutationId* |  | `repo`, `read:org`, `workflow`, `admin:org` |
| `GITHUB_ACCEPT_REPOSITORY_INVITATION` | Accept a repository invitation | write | `invitation_id` |  | `repo:invite`, `repo` |
| `GITHUB_ADD_APP_ACCESS_RESTRICTIONS` | Add app access restrictions | write | `apps`, `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_ADD_ASSIGNEES_TO_AN_ISSUE` | Add assignees to an issue | write | `issue_number`, `owner`, `repo`; *assignees* |  | `repo` |
| `GITHUB_ADD_A_REPOSITORY_COLLABORATOR` | Add a repository collaborator | write | `owner`, `repo`, `username`; *permission* |  | `repo` |
| `GITHUB_ADD_EMAIL_ADDRESS_FOR_AUTHENTICATED_USER` | Add email for auth user | write | `emails` |  | `user` |
| `GITHUB_ADD_FIELD_TO_USER_PROJECT` | Add field to user project | write | `data_type`, `name`, `project_number`, `username`; *iteration_configuration, single_select_options* |  | `project` |
| `GITHUB_ADD_ITEM_TO_USER_PROJECT` | Add item to user project | write | `project_number`, `type`, `username`; *id, number, owner, repo* |  | `project` |
| `GITHUB_ADD_LABELS_TO_AN_ISSUE` | Add labels to an issue | write | `issue_number`, `labels`, `owner`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_ADD_ORG_RUNNER_LABELS` | Add org runner labels | write | `labels`, `org`, `runner_id` |  | `admin:org` |
| `GITHUB_ADD_OR_UPDATE_TEAM_MEMBERSHIP_FOR_USER` | Add or update team membership for a user | write | `org`, `team_slug`, `username`; *role* |  | `admin:org`, `write:org` |
| `GITHUB_ADD_OR_UPDATE_TEAM_PROJECT_PERMISSIONS` | Add or update team project permissions | write | `org`, `project_id`, `team_slug`; *permission* |  | `admin:org` |
| `GITHUB_ADD_OR_UPDATE_TEAM_REPOSITORY_PERMISSIONS` | Add or update team repository permissions | write | `org`, `owner`, `repo`, `team_slug`; *permission* |  | `admin:org`, `repo` |
| `GITHUB_ADD_REPOSITORY_TO_APP_INSTALLATION` | Add a repository to an app installation | write | `installation_id`, `repository_id` |  | `repo` |
| `GITHUB_ADD_REPO_TO_ORG_SECRET_WITH_SELECTED_ACCESS` | Add repo to org secret with selected access | write | `org`, `repository_id`, `secret_name` |  | `admin:org`, `repo` |
| `GITHUB_ADD_REPO_TO_ORG_SECRET_WITH_SELECTED_VISIBILITY` | Add selected repo to org secret | write | `org`, `repository_id`, `secret_name` |  | `admin:org` |
| `GITHUB_ADD_RUNNER_LABELS` | Add runner labels | write | `labels`, `owner`, `repo`, `runner_id` |  | `repo` |
| `GITHUB_ADD_SELECTED_REPOSITORY_TO_ORGANIZATION_SECRET` | Add selected repository to an organization secret | write | `org`, `repository_id`, `secret_name` |  | `admin:org` |
| `GITHUB_ADD_SELECTED_REPOSITORY_TO_ORGANIZATION_VARIABLE` | Add selected repository to an organization variable | write | `name`, `org`, `repository_id` |  | `admin:org`, `repo` |
| `GITHUB_ADD_SELECTED_REPOSITORY_TO_USER_SECRET` | Add a selected repository to a user secret | write | `repository_id`, `secret_name` |  | `codespace`, `codespace:secrets` |
| `GITHUB_ADD_SOCIAL_ACCOUNTS_FOR_AUTHENTICATED_USER` | Add social accounts for the authenticated user | write | `account_urls` |  | `user` |
| `GITHUB_ADD_STATUS_CHECK_CONTEXTS` | Add status check contexts | write | `branch`, `contexts`, `owner`, `repo` |  | `repo` |
| `GITHUB_ADD_SUB_ISSUE` | Add sub-issue to an issue | write | `issue_id`; *client_mutation_id, replace_parent, sub_issue_id, sub_issue_url* |  | `repo`, `public_repo` |
| `GITHUB_ADD_TEAM_ACCESS_RESTRICTIONS` | Add team access restrictions | write | `branch`, `owner`, `repo`, `teams` |  | `repo` |
| `GITHUB_ADD_USERS_TO_CODESPACES_ACCESS_FOR_ORGANIZATION` | Add users to codespaces access for an organization | write | `org`, `selected_usernames` |  | `admin:org` |
| `GITHUB_ADD_USER_ACCESS_RESTRICTIONS` | Add user access restrictions | write | `branch`, `owner`, `repo`, `users` |  | `repo` |
| `GITHUB_API_ROOT` | Get GitHub API root | read |  |  |  |
| `GITHUB_APPROVE_WORKFLOW_RUN_FOR_FORK_PULL_REQUEST` | Approve a workflow run for a fork pull request | write | `owner`, `repo`, `run_id` |  | `repo` |
| `GITHUB_ASSIGN_ORGANIZATION_ROLE_TO_TEAM` | Assign an organization role to a team | write | `org`, `role_id`, `team_slug` |  | `admin:org` |
| `GITHUB_ASSIGN_ORGANIZATION_ROLE_TO_USER` | Assign an organization role to a user | write | `org`, `role_id`, `username` |  | `admin:org` |
| `GITHUB_BLOCK_USER` | Block a user | write | `username` |  | `user` |
| `GITHUB_BLOCK_USER_FROM_ORGANIZATION` | Block a user from an organization | write | `org`, `username` |  | `admin:org` |
| `GITHUB_CANCEL_GITHUB_PAGES_DEPLOYMENT` | Cancel a GitHub Pages deployment | write | `owner`, `pages_deployment_id`, `repo` |  |  |
| `GITHUB_CANCEL_SPONSORSHIP` | Cancel sponsorship | write | *client_mutation_id, sponsor_id, sponsor_login, sponsorable_id, sponsorable_login* |  | `user` |
| `GITHUB_CANCEL_WORKFLOW_RUN` | Cancel a workflow run | write | `owner`, `repo`, `run_id` |  | `repo` |
| `GITHUB_CHECK_FOLLOWING` | Check if following user | read | `username` |  | `user:follow`, `user` |
| `GITHUB_CHECK_IF_GIST_IS_STARRED` | Check if a gist is starred | read | `gist_id` |  | `gist` |
| `GITHUB_CHECK_IF_PULL_REQUEST_HAS_BEEN_MERGED` | Check if pull request merged | read | `owner`, `pull_number`, `repo` |  | `repo` |
| `GITHUB_CHECK_IF_PULL_REQUEST_MERGED` | Check if pull request merged | read | `owner`, `pull_number`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_CHECK_IF_USER_CAN_BE_ASSIGNED` | Check if a user can be assigned | read | `assignee`, `owner`, `repo` |  | `repo` |
| `GITHUB_CHECK_IF_USER_CAN_BE_ASSIGNED_TO_ISSUE` | Check if a user can be assigned to an issue | read | `assignee`, `issue_number`, `owner`, `repo` |  | `repo` |
| `GITHUB_CHECK_IF_USER_FOLLOWS_ANOTHER_USER` | Check if a user follows another user | read | `target_user`, `username` |  | `user:follow` |
| `GITHUB_CHECK_IF_USER_IS_BLOCKED_BY_AUTHENTICATED_USER` | Check if a user is blocked by the authenticated user | read | `username` |  | `read:user`, `user` |
| `GITHUB_CHECK_IF_USER_IS_BLOCKED_BY_ORGANIZATION` | Check if a user is blocked by an organization | read | `org`, `username` |  | `admin:org` |
| `GITHUB_CHECK_IF_USER_IS_REPOSITORY_COLLABORATOR` | Check if a user is a repository collaborator | read | `owner`, `repo`, `username` |  | `repo`, `read:org` |
| `GITHUB_CHECK_ORGANIZATION_MEMBERSHIP_FOR_USER` | Check organization membership for user | read | `org`, `username` |  | `read:org` |
| `GITHUB_CHECK_PERSON_FOLLOWED_BY_AUTH_USER` | Check if person is followed by authenticated user | read | `username` |  | `user:follow` |
| `GITHUB_CHECK_PRIVATE_VULNERABILITY_REPORTING_STATUS` | Check private vulnerability reporting status | read | `owner`, `repo` |  |  |
| `GITHUB_CHECK_PUBLIC_ORGANIZATION_MEMBERSHIP_FOR_USER` | Check public organization membership | read | `org`, `username` |  | `read:org` |
| `GITHUB_CHECK_REPOSITORY_STARRED` | Check if repository is starred | read | `owner`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_CHECK_TEAM_PERMISSIONS_FOR_A_PROJECT` | Check team permissions for a project | read | `org`, `project_id`, `team_slug` |  | `read:org` |
| `GITHUB_CHECK_TEAM_PERMISSIONS_FOR_A_REPOSITORY` | Check team permissions for a repository | read | `org`, `owner`, `repo`, `team_slug` |  | `admin:org`, `repo` |
| `GITHUB_CHECK_TOKEN` | Check a token | read | `access_token`, `client_id` |  |  |
| `GITHUB_CLEAR_PROJECT_V2_ITEM_FIELD_VALUE` | Clear Project V2 Item Field Value | write | `fieldId`, `itemId`, `projectId` |  | `project` |
| `GITHUB_CLEAR_REPOSITORY_CACHE_BY_KEY` | Clear repository cache by key | write | `key`, `owner`, `repo`; *ref* |  | `repo` |
| `GITHUB_CLEAR_SELF_HOSTED_RUNNER_ORG_LABELS` | Clear self-hosted runner org labels | write | `org`, `runner_id` |  | `admin:org` |
| `GITHUB_CLOSE_ISSUE` | Close an issue | write | `issue_number`, `owner`, `repo`; *state_reason* |  | `repo`, `public_repo` |
| `GITHUB_CLOSE_PULL_REQUEST` | Close a pull request | write | `owner`, `pull_number`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_COMMIT_MULTIPLE_FILES` | Commit Multiple Files | write | `branch`, `message`, `owner`, `repo`; *author, base_branch, committer, deletes, force, max_retries, upserts* |  | `repo` |
| `GITHUB_COMPARE_TWO_COMMITS` | Compare two commits | read | `basehead`, `owner`, `repo`; *page, per_page* |  |  |
| `GITHUB_CONFIGURE_JIT_RUNNER_FOR_ORG` | Configure JIT runner for an org | write | `labels`, `name`, `org`, `runner_group_id`; *work_folder* |  | `admin:org`, `repo` |
| `GITHUB_CONFIGURE_OIDC_SUBJECT_CLAIM_TEMPLATE` | Configure OIDC subject claim template | write | `include_claim_keys`, `org` |  | `read:org`, `write:org` |
| `GITHUB_CONVERT_ORG_MEMBER_TO_OUTSIDE_COLLABORATOR` | Convert an organization member to outside collaborator | write | `org`, `username`; *async* |  | `admin:org` |
| `GITHUB_CONVERT_PULL_REQUEST_TO_DRAFT` | Convert Pull Request to Draft | write | `pullRequestId`; *clientMutationId* |  | `repo`, `public_repo` |
| `GITHUB_CREATE_AN_AUTOLINK_REFERENCE_FOR_A_REPOSITORY` | Create an autolink reference for a repository | write | `key_prefix`, `owner`, `repo`, `url_template`; *is_alphanumeric* |  | `repo` |
| `GITHUB_CREATE_AN_ENVIRONMENT_VARIABLE` | Create an environment variable | write | `environment_name`, `name`, `owner`, `repo`, `value` |  | `repo` |
| `GITHUB_CREATE_AN_ISSUE` | Create an issue | write | `owner`, `repo`, `title`; *assignee, assignees, body, labels, milestone* |  | `public_repo`, `repo` |
| `GITHUB_CREATE_AN_ISSUE_COMMENT` | Create an issue comment | write | `body`, `issue_number`, `owner`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_CREATE_AN_ORGANIZATION_PROJECT` | Create an organization project | write | `name`, `org`; *body* |  | `project`, `write:org`, `admin:org`, `repo` |
| `GITHUB_CREATE_AN_ORGANIZATION_REPOSITORY` | Create an organization repository | write | `name`, `org`; *allow_auto_merge, allow_merge_commit, allow_rebase_merge, allow_squash_merge, auto_init, custom_properties, delete_branch_on_merge, description and 16 more* |  | `repo` |
| `GITHUB_CREATE_AN_ORGANIZATION_VARIABLE` | Create an organization variable | write | `name`, `org`, `value`, `visibility`; *selected_repository_ids* |  | `admin:org`, `repo` |
| `GITHUB_CREATE_AN_ORGANIZATION_WEBHOOK` | Create an organization webhook | write | `config__url`, `org`; *active, config__content__type, config__insecure__ssl, config__password, config__secret, config__username, events, name* |  | `admin:org_hook` |
| `GITHUB_CREATE_A_BLOB` | Create a blob | write | `content`, `owner`, `repo`; *encoding* |  | `repo` |
| `GITHUB_CREATE_A_CHECK_RUN` | Create a check run | write | `head_sha`, `name`, `owner`, `repo`; *actions, completed_at, conclusion, details_url, external_id, output__annotations, output__images, output__summary and 4 more* |  | `repo` |
| `GITHUB_CREATE_A_CHECK_SUITE` | Create a check suite | write | `head_sha`, `owner`, `repo` |  |  |
| `GITHUB_CREATE_A_CODESPACE_FOR_THE_AUTHENTICATED_USER` | Create a codespace for the authenticated user | write | *pull_request, ref, repository_id* |  | `codespace` |
| `GITHUB_CREATE_A_CODESPACE_FROM_A_PULL_REQUEST` | Create a codespace from a pull request | write | `owner`, `pull_number`, `repo`; *client_ip, devcontainer_path, display_name, geo, idle_timeout_minutes, location, machine, multi_repo_permissions_opt_out and 2 more* |  | `codespace` |
| `GITHUB_CREATE_A_CODESPACE_IN_A_REPOSITORY` | Create a codespace in a repository | write | `owner`, `repo`; *client_ip, devcontainer_path, display_name, geo, idle_timeout_minutes, location, machine, multi_repo_permissions_opt_out and 3 more* |  | `codespace` |
| `GITHUB_CREATE_A_COMMIT` | Create a commit | write | `message`, `owner`, `repo`, `tree`; *author__date, author__email, author__name, committer__date, committer__email, committer__name, parents, signature* |  | `repo` |
| `GITHUB_CREATE_A_COMMIT_COMMENT` | Create a commit comment | write | `body`, `commit_sha`, `owner`, `repo`; *line, path, position* |  | `repo`, `public_repo` |
| `GITHUB_CREATE_A_COMMIT_STATUS` | Create a commit status | write | `owner`, `repo`, `sha`, `state`; *context, description, target_url* |  | `repo`, `repo:status` |
| `GITHUB_CREATE_A_CUSTOM_ORGANIZATION_ROLE` | Create a custom organization role | write | `name`, `org`, `permissions`; *description* |  | `admin:org` |
| `GITHUB_CREATE_A_DEPLOYMENT` | Create a deployment | write | `owner`, `ref`, `repo`; *auto_merge, description, environment, payload, production_environment, required_contexts, task, transient_environment* |  | `repo`, `repo_deployment` |
| `GITHUB_CREATE_A_DEPLOYMENT_BRANCH_POLICY` | Create a deployment branch policy | write | `environment_name`, `name`, `owner`, `repo`; *type* |  | `repo` |
| `GITHUB_CREATE_A_DEPLOYMENT_STATUS` | Create a deployment status | write | `deployment_id`, `owner`, `repo`, `state`; *auto_inactive, description, environment, environment_url, log_url, target_url* |  | `repo_deployment` |
| `GITHUB_CREATE_A_DEPLOY_KEY` | Create a deploy key | write | `key`, `owner`, `repo`; *read_only, title* |  | `repo` |
| `GITHUB_CREATE_A_DISCUSSION` | Create a discussion | write | `body`, `org`, `team_slug`, `title`; *private* |  | `write:discussion` |
| `GITHUB_CREATE_A_DISCUSSION_COMMENT` | Create a discussion comment | write | `body`, `discussion_number`, `org`, `team_slug` |  | `write:discussion` |
| `GITHUB_CREATE_A_FORK` | Create a fork | write | `owner`, `repo`; *default_branch_only, name, organization* |  | `repo` |
| `GITHUB_CREATE_A_GIST` | Create a gist | write | `files`; *description, public* |  | `gist` |
| `GITHUB_CREATE_A_GIST_COMMENT` | Create a gist comment | write | `body`, `gist_id` |  | `gist` |
| `GITHUB_CREATE_A_GITHUB_APP_FROM_A_MANIFEST` | Create a GitHub App from a manifest | write | `code` |  |  |
| `GITHUB_CREATE_A_GITHUB_PAGES_DEPLOYMENT` | Create a github pages deployment | write | `oidc_token`, `owner`, `repo`; *artifact_id, artifact_url, environment, pages_build_version* |  |  |
| `GITHUB_CREATE_A_GITHUB_PAGES_SITE` | Create a github pages site | write | `owner`, `repo`; *build_type, source__branch, source__path* |  | `repo` |
| `GITHUB_CREATE_A_LABEL` | Create a label | write | `color`, `name`, `owner`, `repo`; *description* |  | `repo` |
| `GITHUB_CREATE_A_MILESTONE` | Create a milestone | write | `owner`, `repo`, `title`; *description, due_on, state* |  | `repo`, `public_repo` |
| `GITHUB_CREATE_A_PROJECT_CARD` | Create a project card | write | `column_id`; *content_id, content_type, note* |  | `project`, `public_repo`, `repo` |
| `GITHUB_CREATE_A_PULL_REQUEST` | Create a pull request | write | `base`, `head`, `owner`, `repo`; *body, draft, head_repo, issue, maintainer_can_modify, title* |  | `repo` |
| `GITHUB_CREATE_A_REFERENCE` | Create a reference | write | `owner`, `ref`, `repo`, `sha` |  | `repo` |
| `GITHUB_CREATE_A_REGISTRATION_TOKEN_FOR_AN_ORGANIZATION` | Create a registration token for an organization | write | `org` |  | `admin:org`, `repo` |
| `GITHUB_CREATE_A_REGISTRATION_TOKEN_FOR_A_REPOSITORY` | Create a registration token for a repository | write | `owner`, `repo` |  | `repo` |
| `GITHUB_CREATE_A_RELEASE` | Create a release | write | `owner`, `repo`, `tag_name`; *body, discussion_category_name, draft, generate_release_notes, make_latest, name, prerelease, target_commitish* |  | `repo`, `public_repo` |
| `GITHUB_CREATE_A_REMOVE_TOKEN_FOR_AN_ORGANIZATION` | Create a remove token for an organization | write | `org` |  | `admin:org`, `repo` |
| `GITHUB_CREATE_A_REMOVE_TOKEN_FOR_A_REPOSITORY` | Create a remove token for a repository | write | `owner`, `repo` |  | `repo` |
| `GITHUB_CREATE_A_REPLY_FOR_A_REVIEW_COMMENT` | Create a reply for a review comment | write | `body`, `comment_id`, `owner`, `pull_number`, `repo` |  | `repo` |
| `GITHUB_CREATE_A_REPOSITORY_DISPATCH_EVENT` | Create a repository dispatch event | write | `event_type`, `owner`, `repo`; *client_payload* |  | `repo` |
| `GITHUB_CREATE_A_REPOSITORY_FOR_THE_AUTHENTICATED_USER` | Create a repository for the authenticated user | write | `name`; *allow_auto_merge, allow_merge_commit, allow_rebase_merge, allow_squash_merge, auto_init, delete_branch_on_merge, description, gitignore_template and 14 more* |  | `repo` |
| `GITHUB_CREATE_A_REPOSITORY_FROM_AN_UNPUBLISHED_CODESPACE` | Create a repository from an unpublished codespace | write | `codespace_name`; *name, private* |  | `codespace` |
| `GITHUB_CREATE_A_REPOSITORY_PROJECT` | Create a repository project | write | `name`, `owner`, `repo`; *body* |  | `project`, `repo` |
| `GITHUB_CREATE_A_REPOSITORY_RULESET` | Create a repository ruleset | write | `enforcement`, `name`, `owner`, `repo`; *bypass_actors, conditions__ref__name__exclude, conditions__ref__name__include, rules, target* |  | `repo` |
| `GITHUB_CREATE_A_REPOSITORY_USING_A_TEMPLATE` | Create a repository using a template | write | `name`, `template_owner`, `template_repo`; *description, include_all_branches, owner, private* |  | `repo` |
| `GITHUB_CREATE_A_REPOSITORY_VARIABLE` | Create a repository variable | write | `name`, `owner`, `repo`, `value` |  | `repo` |
| `GITHUB_CREATE_A_REPOSITORY_WEBHOOK` | Create a repository webhook | write | `owner`, `repo`; *active, config__content__type, config__insecure__ssl, config__secret, config__url, events, name* |  | `admin:repo_hook`, `read:repo_hook`, `repo`, `write:repo_hook` |
| `GITHUB_CREATE_A_REVIEW_COMMENT_FOR_A_PULL_REQUEST` | Create a review comment for a pull request | write | `body`, `commit_id`, `owner`, `path`, `pull_number`, `repo`; *in_reply_to, line, position, side, start_line, start_side* |  | `public_repo`, `repo` |
| `GITHUB_CREATE_A_REVIEW_FOR_A_PULL_REQUEST` | Create a review for a pull request | write | `owner`, `pull_number`, `repo`; *body, comments, commit_id, event* |  | `repo` |
| `GITHUB_CREATE_A_SCOPED_ACCESS_TOKEN` | Create a scoped access token | write | `access_token`, `client_id`; *permissions__actions, permissions__administration, permissions__checks, permissions__codespaces, permissions__contents, permissions__dependabot__secrets, permissions__deployments, permissions__email__addresses and 44 more* |  | `repo`, `repo:status`, `repo_deployment`, `public_repo`, `repo:invite`, `security_events`, `admin:repo_hook`, `write:repo_hook`, `read:repo_hook`, `admin:org`, `write:org`, `read:org`, `admin:public_key`, `write:public_key`, `read:public_key`, `admin:org_hook`, `gist`, `notifications`, `user`, `read:user`, `user:email`, `user:follow`, `project`, `read:project`, `delete_repo`, `write:packages`, `read:packages`, `delete:packages`, `admin:gpg_key`, `write:gpg_key`, `read:gpg_key`, `codespace`, `workflow` |
| `GITHUB_CREATE_A_SNAPSHOT_OF_DEPENDENCIES_FOR_A_REPOSITORY` | Create a dependency snapshot | write | `owner`, `ref`, `repo`, `scanned`, `sha`, `version`; *detector__name, detector__url, detector__version, job__correlator, job__html__url, job__id, manifests, metadata* |  | `repo` |
| `GITHUB_CREATE_A_TAG_OBJECT` | Create a tag object | write | `message`, `object`, `owner`, `repo`, `tag`, `type`; *tagger__date, tagger__email, tagger__name* |  | `repo` |
| `GITHUB_CREATE_A_TEAM` | Create a team | write | `name`, `org`; *description, maintainers, notification_setting, parent_team_id, permission, privacy, repo_names* |  | `admin:org` |
| `GITHUB_CREATE_A_TEMPORARY_PRIVATE_FORK` | Create a temporary private fork | write | `ghsa_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_CREATE_A_TREE` | Create a tree | write | `owner`, `repo`, `tree`; *base_tree* |  | `repo` |
| `GITHUB_CREATE_A_USER_PROJECT` | Create a user project | write | `name`; *body* |  | `project`, `repo` |
| `GITHUB_CREATE_A_WORKFLOW_DISPATCH_EVENT` | Create a workflow dispatch event | write | `owner`, `ref`, `repo`, `workflow_id`; *inputs* |  | `repo` |
| `GITHUB_CREATE_BRANCH` | Create a branch | write | `branch_name`, `owner`, `repo`, `sha` |  | `repo` |
| `GITHUB_CREATE_COMMIT_SIGNATURE_PROTECTION` | Create commit signature protection | write | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_CREATE_DEPLOYMENT_PROTECTION_RULE` | Create a custom deployment protection rule on an environment | write | `environment_name`, `owner`, `repo`; *integration_id* |  | `repo` |
| `GITHUB_CREATE_DRAFT_ITEM_FOR_USER_PROJECT` | Create draft item for user project | write | `project_number`, `title`, `username`; *body* |  | `project` |
| `GITHUB_CREATE_INFERENCE_CHAT_COMPLETIONS` | Create inference chat completions | read | `messages`, `model`; *frequency_penalty, max_tokens, modalities, presence_penalty, response_format, seed, stop, stream and 5 more* |  |  |
| `GITHUB_CREATE_INFERENCE_EMBEDDINGS` | Create inference embeddings | read | `input`, `model`; *dimensions, encoding_format, user* |  |  |
| `GITHUB_CREATE_ISSUE_TYPE` | Create issue type | write | `is_enabled`, `name`, `org`; *color, description* |  | `admin:org` |
| `GITHUB_CREATE_JIT_RUNNER_CONFIG` | Create JIT runner config for repo | write | `labels`, `name`, `owner`, `repo`, `runner_group_id`; *work_folder* |  | `repo` |
| `GITHUB_CREATE_OR_UPDATE_AN_ENVIRONMENT` | Create or update an environment | write | `environment_name`, `owner`, `repo`; *deployment__branch__policy__custom__branch__policies, deployment__branch__policy__protected__branches, prevent_self_review, reviewers, wait_timer* |  | `repo` |
| `GITHUB_CREATE_OR_UPDATE_AN_ENVIRONMENT_SECRET` | Create or update an environment secret | write | `encrypted_value`, `environment_name`, `key_id`, `owner`, `repo`, `secret_name` |  | `repo` |
| `GITHUB_CREATE_OR_UPDATE_AN_ORGANIZATION_SECRET` | Create or update an organization secret | write | `encrypted_value`, `key_id`, `org`, `secret_name`, `visibility`; *selected_repository_ids* |  | `admin:org`, `repo` |
| `GITHUB_CREATE_OR_UPDATE_A_REPOSITORY_SECRET` | Create or update a repository secret | write | `owner`, `repo`, `secret_name`; *encrypted_value, key_id* |  | `repo` |
| `GITHUB_CREATE_OR_UPDATE_A_SECRET_FOR_THE_AUTHENTICATED_USER` | Create or update a secret for the authenticated user | write | `secret_name`; *encrypted_value, key_id, secret_value, selected_repository_ids* |  | `codespace`, `codespace:secrets` |
| `GITHUB_CREATE_OR_UPDATE_CUSTOM_PROPERTIES_FOR_AN_ORG` | Create or update custom properties for an organization | write | `org`, `properties` |  | `admin:org` |
| `GITHUB_CREATE_OR_UPDATE_CUSTOM_PROPERTY_ORG` | Create or update a custom property for an organization | write | `custom_property_name`, `org`, `value_type`; *allowed_values, default_value, description, required, values_editable_by* |  | `admin:org` |
| `GITHUB_CREATE_OR_UPDATE_FILE_CONTENTS` | Create or update file contents | write | `content`, `message`, `owner`, `path`, `repo`; *author__date, author__email, author__name, branch, committer__date, committer__email, committer__name, sha* |  | `public_repo`, `repo`, `workflow` |
| `GITHUB_CREATE_OR_UPDATE_GITHUB_PAGES_SITE` | Create or Update GitHub Pages Site | write | `owner`, `repo`; *build_type, cname, https_enforced, source_branch, source_path* |  | `repo` |
| `GITHUB_CREATE_OR_UPDATE_REPO_SECRET_WITH_ENCRYPTED_VALUE` | Create or update repo secret with encrypted value | write | `owner`, `repo`, `secret_name`; *encrypted_value, key_id* |  | `repo` |
| `GITHUB_CREATE_REACTION_FOR_AN_ISSUE` | Create reaction for an issue | write | `content`, `issue_number`, `owner`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_CREATE_REACTION_FOR_AN_ISSUE_COMMENT` | Create reaction for an issue comment | write | `comment_id`, `content`, `owner`, `repo` |  | `repo` |
| `GITHUB_CREATE_REACTION_FOR_A_COMMIT_COMMENT` | Create reaction for a commit comment | write | `comment_id`, `content`, `owner`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_CREATE_REACTION_FOR_A_PULL_REQUEST_REVIEW_COMMENT` | Create reaction for a pull request review comment | write | `comment_id`, `content`, `owner`, `repo` |  |  |
| `GITHUB_CREATE_REACTION_FOR_A_RELEASE` | Create reaction for a release | write | `content`, `owner`, `release_id`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_CREATE_REACTION_FOR_A_TEAM_DISCUSSION` | Create reaction for a team discussion | write | `content`, `discussion_number`, `org`, `team_slug` |  | `write:discussion` |
| `GITHUB_CREATE_REACTION_FOR_A_TEAM_DISCUSSION_COMMENT` | Create reaction for a team discussion comment | write | `comment_number`, `content`, `discussion_number`, `org`, `team_slug` |  | `write:discussion` |
| `GITHUB_CREATE_REPOSITORY` | Create repository | write | `name`; *allow_auto_merge, allow_merge_commit, allow_rebase_merge, allow_squash_merge, auto_init, delete_branch_on_merge, description, gitignore_template and 13 more* |  | `public_repo`, `repo` |
| `GITHUB_CREATE_SPONSORSHIP` | Create sponsorship | write | `amount`; *client_mutation_id, is_recurring, privacy_level, sponsor_id, sponsor_login, sponsorable_id, sponsorable_login* |  | `user` |
| `GITHUB_CREATE_SPONSORS_TIER` | Create Sponsors Tier | write | `amount`, `description`, `sponsorable_id`; *client_mutation_id* |  | `user` |
| `GITHUB_CREATE_UPDATE_ORG_SECRET_WITH_LIB_SODIUM` | Create or Update Organization Secret with LibSodium | write | `org`, `secret_name`, `visibility`; *encrypted_value, key_id, selected_repository_ids* |  | `admin:org` |
| `GITHUB_CREATE_USER_LIST` | Create user list | write | `name`; *clientMutationId, description* |  | `user` |
| `GITHUB_CREATE_VIEW_FOR_USER_PROJECT` | Create view for user project | write | `layout`, `name`, `project_number`, `user_id`; *filter, visible_fields* |  | `project`, `repo` |
| `GITHUB_DECLINE_A_REPOSITORY_INVITATION` | Decline a repository invitation | write | `invitation_id` |  | `repo`, `repo:invite` |
| `GITHUB_DELETE_ACCESS_RESTRICTIONS` | Delete access restrictions | write | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_ADMIN_BRANCH_PROTECTION` | Delete admin branch protection | write | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_APP_AUTHORIZATION` | Delete an app authorization | write | `access_token`, `client_id` |  |  |
| `GITHUB_DELETE_APP_TOKEN` | Delete an app token | write | `access_token`, `client_id` |  |  |
| `GITHUB_DELETE_ARTIFACT` | Delete an artifact | write | `artifact_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_AUTOLINK_REFERENCE` | Delete an autolink reference from a repository | write | `autolink_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_A_FILE` | Delete a file | write | `message`, `owner`, `path`, `repo`; *author__email, author__name, branch, committer__email, committer__name, sha* |  | `repo` |
| `GITHUB_DELETE_A_PACKAGE_VERSION_FOR_THE_AUTHENTICATED_USER` | Delete a package version for the authenticated user | write | `package_name`, `package_type`, `package_version_id` |  | `read:packages`, `delete:packages` |
| `GITHUB_DELETE_A_REFERENCE` | Delete a reference | write | `owner`, `ref`, `repo` |  | `repo` |
| `GITHUB_DELETE_A_RELEASE` | Delete a release | write | `owner`, `release_id`, `repo` |  | `delete_repo`, `repo` |
| `GITHUB_DELETE_A_RELEASE_ASSET` | Delete a release asset | write | `asset_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_A_RELEASE_REACTION` | Delete a release reaction | write | `owner`, `reaction_id`, `release_id`, `repo` |  | `repo` |
| `GITHUB_DELETE_A_REPOSITORY` | Delete a repository | write | `owner`, `repo` |  | `delete_repo`, `repo` |
| `GITHUB_DELETE_A_REPOSITORY_INVITATION` | Delete a repository invitation | write | `invitation_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_A_REPOSITORY_RULESET` | Delete a repository ruleset | write | `owner`, `repo`, `ruleset_id` |  | `repo` |
| `GITHUB_DELETE_A_REPOSITORY_SECRET` | Delete a repository secret | write | `owner`, `repo`, `secret_name` |  | `repo` |
| `GITHUB_DELETE_A_REPOSITORY_SUBSCRIPTION` | Delete a repository subscription | write | `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_A_REPOSITORY_VARIABLE` | Delete a repository variable | write | `name`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_A_REPOSITORY_WEBHOOK` | Delete a repository webhook | write | `hook_id`, `owner`, `repo` |  | `admin:repo_hook`, `repo` |
| `GITHUB_DELETE_A_REVIEW_COMMENT_FOR_A_PULL_REQUEST` | Delete a review comment for a pull request | write | `comment_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_A_SECRET_FOR_THE_AUTHENTICATED_USER` | Delete a secret for the authenticated user | write | `secret_name` |  | `codespace`, `codespace:secrets` |
| `GITHUB_DELETE_A_SELF_HOSTED_RUNNER_FROM_AN_ORGANIZATION` | Delete a self hosted runner from an organization | write | `org`, `runner_id` |  | `admin:org`, `repo` |
| `GITHUB_DELETE_A_SELF_HOSTED_RUNNER_FROM_A_REPOSITORY` | Delete a self hosted runner from a repository | write | `owner`, `repo`, `runner_id` |  | `repo` |
| `GITHUB_DELETE_A_TEAM` | Delete a team | write | `org`, `team_slug` |  | `admin:org` |
| `GITHUB_DELETE_A_THREAD_SUBSCRIPTION` | Delete a thread subscription | write | `thread_id` |  | `notifications`, `repo` |
| `GITHUB_DELETE_A_WORKFLOW_RUN` | Delete a workflow run | write | `owner`, `repo`, `run_id` |  | `repo` |
| `GITHUB_DELETE_BRANCH_PROTECTION` | Delete branch protection | write | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_CODESPACE` | Delete a codespace for the authenticated user | write | `codespace_name` |  | `codespace` |
| `GITHUB_DELETE_CODE_SCANNING_ANALYSIS` | Delete a code scanning analysis from a repository | write | `analysis_id`, `owner`, `repo`; *confirm_delete* |  | `repo`, `public_repo` |
| `GITHUB_DELETE_COMMIT_COMMENT` | Delete a commit comment | write | `comment_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_COMMIT_COMMENT_REACTION` | Delete a commit comment reaction | write | `comment_id`, `owner`, `reaction_id`, `repo` |  | `repo` |
| `GITHUB_DELETE_COMMIT_SIGNATURE_PROTECTION` | Delete commit signature protection | write | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_CUSTOM_ORGANIZATION_ROLE` | Delete a custom organization role | write | `org`, `role_id` |  | `admin:org` |
| `GITHUB_DELETE_DEPENDEBOT_SECRET_BY_NAME` | Delete dependabot secret by name | write | `owner`, `repo`, `secret_name` |  | `repo` |
| `GITHUB_DELETE_DEPLOYMENT` | Delete a deployment | write | `deployment_id`, `owner`, `repo` |  | `repo`, `repo_deployment` |
| `GITHUB_DELETE_DEPLOYMENT_BRANCH_POLICY` | Delete a deployment branch policy | write | `branch_policy_id`, `environment_name`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_DEPLOY_KEY` | Delete a deploy key | write | `key_id`, `owner`, `repo` |  | `admin:public_key`, `repo` |
| `GITHUB_DELETE_DISCUSSION` | Delete a discussion | write | `discussion_number`, `org`, `team_slug` |  | `write:discussion` |
| `GITHUB_DELETE_DISCUSSION_COMMENT` | Delete a discussion comment | write | `comment_number`, `discussion_number`, `org`, `team_slug` |  | `write:discussion` |
| `GITHUB_DELETE_EMAIL_ADDRESS` | Delete authenticated user email address | write | `emails` |  | `user`, `user:email` |
| `GITHUB_DELETE_ENVIRONMENT` | Delete an environment | write | `environment_name`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_ENVIRONMENT_SECRET` | Delete an environment secret | write | `environment_name`, `owner`, `repo`, `secret_name` |  | `repo` |
| `GITHUB_DELETE_ENVIRONMENT_VARIABLE` | Delete an environment variable | write | `environment_name`, `name`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_GIST` | Delete a gist | write | `gist_id` |  | `gist` |
| `GITHUB_DELETE_GIST_COMMENT` | Delete a gist comment | write | `comment_id`, `gist_id` |  | `gist` |
| `GITHUB_DELETE_GITHUB_ACTIONS_CACHE_BY_ID` | Delete GitHub Actions cache by ID | write | `cache_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_GITHUB_PAGES_SITE` | Delete a GitHub Pages site | write | `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_ISSUE_COMMENT` | Delete an issue comment | write | `comment_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_ISSUE_COMMENT_REACTION` | Delete an issue comment reaction | write | `comment_id`, `owner`, `reaction_id`, `repo` |  | `public_repo`, `repo` |
| `GITHUB_DELETE_ISSUE_REACTION` | Delete an issue reaction | write | `issue_number`, `owner`, `reaction_id`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_DELETE_LABEL` | Delete a label | write | `name`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_MILESTONE` | Delete a milestone | write | `milestone_number`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_ORGANIZATION` | Delete an organization | write | `org` |  | `admin:org` |
| `GITHUB_DELETE_ORGANIZATION_SECRET` | Delete an organization secret | write | `org`, `secret_name` |  | `admin:org`, `repo` |
| `GITHUB_DELETE_ORGANIZATION_VARIABLE` | Delete an organization variable | write | `name`, `org` |  | `admin:org`, `repo` |
| `GITHUB_DELETE_ORGANIZATION_WEBHOOK` | Delete an organization webhook | write | `hook_id`, `org` |  | `admin:org_hook` |
| `GITHUB_DELETE_ORG_CODESPACE` | Delete a codespace from the organization | write | `codespace_name`, `org`, `username` |  | `admin:org` |
| `GITHUB_DELETE_ORG_PACKAGE` | Delete a package for an organization | write | `org`, `package_name`, `package_type` |  | `read:packages`, `delete:packages` |
| `GITHUB_DELETE_PACKAGE` | Delete a package for the authenticated user | write | `package_name`, `package_type` |  | `read:packages`, `delete:packages` |
| `GITHUB_DELETE_PACKAGE_VERSION` | Delete package version | write | `package_version_id`; *client_mutation_id* |  | `read:packages`, `delete:packages`, `repo` |
| `GITHUB_DELETE_PACKAGE_VERSION_FOR_AN_ORGANIZATION` | Delete package version for an organization | write | `org`, `package_name`, `package_type`, `package_version_id` |  | `read:packages`, `delete:packages` |
| `GITHUB_DELETE_PACKAGE_VERSION_FOR_A_USER` | Delete package version for a user | write | `package_name`, `package_type`, `package_version_id`, `username` |  | `read:packages`, `delete:packages` |
| `GITHUB_DELETE_PENDING_REVIEW` | Delete a pending review for a pull request | write | `owner`, `pull_number`, `repo`, `review_id` |  | `repo` |
| `GITHUB_DELETE_PROJECT` | Delete a project | write | `project_id` |  | `project`, `repo` |
| `GITHUB_DELETE_PROJECT_COLUMN` | Delete a project field (column) | write | `field_id` |  | `project`, `repo` |
| `GITHUB_DELETE_PROJECT_ITEM_FOR_USER` | Delete project item for user | write | `item_id`, `project_number`, `username` |  | `project`, `repo` |
| `GITHUB_DELETE_PULL_REQUEST_COMMENT_REACTION` | Delete a pull request comment reaction | write | `comment_id`, `owner`, `reaction_id`, `repo` |  | `repo` |
| `GITHUB_DELETE_PULL_REQUEST_REVIEW_PROTECTION` | Delete pull request review protection | write | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_DELETE_REPO_CODESPACE_SECRET_BY_NAME` | Delete repo codespace secret by name | write | `owner`, `repo`, `secret_name` |  | `repo` |
| `GITHUB_DELETE_SOCIAL_ACCOUNTS_FOR_THE_AUTHENTICATED_USER` | Delete social accounts for the authenticated user | write | `account_urls` |  | `user` |
| `GITHUB_DELETE_TEAM_DISCUSSION_COMMENT_REACTION` | Delete team discussion comment reaction | write | `comment_number`, `discussion_number`, `org`, `reaction_id`, `team_slug` |  | `write:discussion` |
| `GITHUB_DELETE_TEAM_DISCUSSION_REACTION` | Delete team discussion reaction | write | `discussion_number`, `org`, `reaction_id`, `team_slug` |  | `write:discussion` |
| `GITHUB_DELETE_USER_LIST` | Delete User List | write | `listId`; *clientMutationId* |  | `user` |
| `GITHUB_DELETE_USER_PACKAGE` | Delete a package for a user | write | `package_name`, `package_type`, `username` |  | `read:packages`, `delete:packages` |
| `GITHUB_DELETE_WORKFLOW_RUN_LOGS` | Delete workflow run logs | write | `owner`, `repo`, `run_id` |  | `repo` |
| `GITHUB_DISABLE_A_CUSTOM_PROTECTION_RULE_FOR_AN_ENVIRONMENT` | Disable a custom protection rule for an environment | write | `environment_name`, `owner`, `protection_rule_id`, `repo` |  | `repo` |
| `GITHUB_DISABLE_A_WORKFLOW` | Disable a workflow | write | `owner`, `repo`, `workflow_id` |  | `repo` |
| `GITHUB_DISABLE_PRIVATE_VULN_REPORTING_FOR_REPO` | Disable private vulnerability reporting for a repository | write | `owner`, `repo` |  | `repo` |
| `GITHUB_DISABLE_REPOSITORY_ACTIONS_IN_ORG` | Disable Repository Actions In Org | write | `org`, `repository_id` |  | `admin:org` |
| `GITHUB_DISMISS_A_REVIEW_FOR_A_PULL_REQUEST` | Dismiss a review for a pull request | write | `message`, `owner`, `pull_number`, `repo`, `review_id`; *event* |  | `repo` |
| `GITHUB_DOWNLOAD_AN_ARTIFACT` | Download an artifact | read | `artifact_id`, `owner`, `repo`; *archive_format* |  | `repo` |
| `GITHUB_DOWNLOAD_A_REPOSITORY_ARCHIVE_TAR` | Download a repository tarball | read | `owner`, `ref`, `repo` |  | `public_repo`, `repo` |
| `GITHUB_DOWNLOAD_A_REPOSITORY_ARCHIVE_ZIP` | Download a repository archive ZIP | read | `owner`, `ref`, `repo` |  | `repo` |
| `GITHUB_DOWNLOAD_JOB_LOGS_FOR_A_WORKFLOW_RUN` | Download job logs for a workflow run | read | `job_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_DOWNLOAD_WORKFLOW_RUN_ATTEMPT_LOGS` | Download workflow run attempt logs | read | `attempt_number`, `owner`, `repo`, `run_id` |  | `repo` |
| `GITHUB_DOWNLOAD_WORKFLOW_RUN_LOGS` | Download workflow run logs | read | `owner`, `repo`, `run_id` |  | `repo` |
| `GITHUB_ENABLE_A_WORKFLOW` | Enable a workflow | write | `owner`, `repo`, `workflow_id` |  | `repo`, `workflow` |
| `GITHUB_ENABLE_GITHUB_ACTIONS_IN_SELECTED_REPOSITORIES` | Enable GitHub actions in selected repositories | write | `org`, `selected_repository_ids` |  | `admin:org` |
| `GITHUB_ENABLE_PRIVATE_VULN_REPORTING_FOR_REPO` | Enable private vulnerability reporting for a repository | write | `owner`, `repo` |  | `repo` |
| `GITHUB_ENABLE_REPO_FOR_GITHUB_ACTIONS` | Enable repo for Github Actions | write | `org`, `repository_id` |  | `admin:org` |
| `GITHUB_ENCRYPT_AND_UPDATE_DEV_SECRET` | Encrypt and update dev secret | write | `owner`, `repo`, `secret_name`; *encrypted_value, key_id, secret_value* |  | `repo` |
| `GITHUB_ENCRYPT_ORG_DEV_ENV_SECRET` | Encrypt org dev env secret | write | `org`, `secret_name`, `visibility`; *encrypted_value, key_id, secret_value, selected_repository_ids* |  | `admin:org` |
| `GITHUB_EXPORT_A_CODESPACE_FOR_THE_AUTHENTICATED_USER` | Export a codespace for the authenticated user | write | `codespace_name` |  | `codespace` |
| `GITHUB_EXPORT_SBOM_FOR_REPO` | Export an SBOM for a repository | read | `owner`, `repo` |  |  |
| `GITHUB_FIND_PULL_REQUESTS` | Find Pull Requests | read | *assignee, author, base_branch, created_since, for_authenticated_user, head_branch, is_merged, label and 12 more* |  | `public_repo`, `repo` |
| `GITHUB_FIND_REPOSITORIES` | Find Repositories | read | `query`; *archived, for_authenticated_user, fork_filter, language, max_stars, min_forks, min_stars, order and 6 more* |  | `repo`, `public_repo` |
| `GITHUB_FOLLOW_ORGANIZATION_GRAPHQL` | Follow Organization (GraphQL) | write | `organizationId`; *clientMutationId* |  | `user`, `user:follow` |
| `GITHUB_FOLLOW_USER` | Follow a user | write | `username` |  | `user:follow` |
| `GITHUB_FOLLOW_USER_GRAPHQL` | Follow User (GraphQL) | write | `userId`; *clientMutationId* |  | `user:follow` |
| `GITHUB_FORCE_CANCEL_WORKFLOW_RUN` | Force cancel a workflow run | write | `owner`, `repo`, `run_id` |  | `repo` |
| `GITHUB_FORK_GIST` | Fork a gist | write | `gist_id` |  | `gist` |
| `GITHUB_GENERATE_RELEASE_NOTES` | Generate release notes content for a release | read | `owner`, `repo`, `tag_name`; *configuration_file_path, previous_tag_name, target_commitish* |  | `repo` |
| `GITHUB_GET_ACCESS_RESTRICTIONS` | Get access restrictions | read | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_ADMIN_BRANCH_PROTECTION` | Get admin branch protection | read | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_ALL_API_VERSIONS` | Get all API versions | read |  |  |  |
| `GITHUB_GET_ALL_CONTRIBUTOR_COMMIT_ACTIVITY` | Get all contributor commit activity | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_ALL_DEPLOYMENT_PROTECTION_RULES_FOR_ENV` | Get all deployment protection rules for an environment | read | `environment_name`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_AN_APP` | Get an app | read | `app_slug` |  |  |
| `GITHUB_GET_AN_ARTIFACT` | Get an artifact | read | `artifact_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_AN_ASSIGNMENT` | Get an assignment | read | `assignment_id` |  |  |
| `GITHUB_GET_AN_AUTOLINK_REFERENCE_OF_A_REPOSITORY` | Get an autolink reference of a repository | read | `autolink_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_AN_ENVIRONMENT` | Get an environment | read | `environment_name`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_AN_ENVIRONMENT_PUBLIC_KEY` | Get an environment public key | read | `environment_name`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_AN_ENVIRONMENT_SECRET` | Get an environment secret | read | `environment_name`, `owner`, `repo`, `secret_name` |  | `repo` |
| `GITHUB_GET_AN_ENVIRONMENT_VARIABLE` | Get an environment variable | read | `environment_name`, `name`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_AN_ISSUE` | Get an issue | read | `issue_number`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_AN_ISSUE_COMMENT` | Get an issue comment | read | `comment_id`, `owner`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_GET_AN_ISSUE_EVENT` | Get an issue event | read | `event_id`, `owner`, `repo` |  | `public_repo`, `repo` |
| `GITHUB_GET_AN_ORGANIZATION` | Get an organization | read | `org` |  | `admin:org`, `read:org` |
| `GITHUB_GET_AN_ORGANIZATION_PUBLIC_KEY` | Get an organization public key | read | `org` |  | `admin:org`, `repo` |
| `GITHUB_GET_AN_ORGANIZATION_ROLE` | Get an organization role | read | `org`, `role_id` |  | `admin:org` |
| `GITHUB_GET_AN_ORGANIZATION_SECRET` | Get an organization secret | read | `org`, `secret_name` |  | `admin:org`, `repo` |
| `GITHUB_GET_AN_ORGANIZATION_VARIABLE` | Get an organization variable | read | `name`, `org` |  | `admin:org`, `repo` |
| `GITHUB_GET_AN_ORGANIZATION_WEBHOOK` | Get an organization webhook | read | `hook_id`, `org` |  | `admin:org_hook`, `read:org` |
| `GITHUB_GET_ASSIGNMENT_GRADES` | Get assignment grades | read | `assignment_id` |  |  |
| `GITHUB_GET_AUTHENTICATED_USER_ORGANIZATION_MEMBERSHIP` | Get authenticated user organization membership | read | `org` |  | `read:org`, `user` |
| `GITHUB_GET_A_BLOB` | Get a blob | read | `file_sha`, `owner`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_GET_A_BRANCH` | Get a branch | read | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_A_COMMIT` | Get a commit | read | `owner`, `ref`, `repo`; *page, per_page* |  | `repo`, `public_repo` |
| `GITHUB_GET_A_MILESTONE` | Get a milestone | read | `milestone_number`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_A_PACKAGE_FOR_AN_ORGANIZATION` | Get a package for an organization | read | `org`, `package_name`, `package_type` |  | `read:packages` |
| `GITHUB_GET_A_PACKAGE_FOR_A_USER` | Get a package for a user | read | `package_name`, `package_type`, `username` |  | `read:packages`, `repo` |
| `GITHUB_GET_A_PACKAGE_FOR_THE_AUTHENTICATED_USER` | Get a package for the authenticated user | read | `package_name`, `package_type` |  | `read:packages`, `repo` |
| `GITHUB_GET_A_PACKAGE_VERSION_FOR_AN_ORGANIZATION` | Get a package version for an organization | read | `org`, `package_name`, `package_type`, `package_version_id` |  | `read:packages` |
| `GITHUB_GET_A_PACKAGE_VERSION_FOR_A_USER` | Get a package version for a user | read | `package_name`, `package_type`, `package_version_id`, `username` |  | `read:packages` |
| `GITHUB_GET_A_PACKAGE_VERSION_FOR_THE_AUTHENTICATED_USER` | Get authenticated user package version | read | `package_name`, `package_type`, `package_version_id` |  | `read:packages`, `repo` |
| `GITHUB_GET_A_PROJECT` | Get a project | read | `project_number`; *org, username* |  | `repo`, `public_repo`, `project`, `read:project` |
| `GITHUB_GET_A_PROJECT_CARD` | Get a project card | read | `card_id` |  | `public_repo`, `read:project`, `repo` |
| `GITHUB_GET_A_PULL_REQUEST` | Get a pull request | read | `owner`, `pull_number`, `repo` |  | `repo` |
| `GITHUB_GET_A_REFERENCE` | Get a reference | read | `owner`, `ref`, `repo` |  | `public_repo`, `repo` |
| `GITHUB_GET_A_RELEASE` | Get a release | read | `owner`, `release_id`, `repo` |  | `public_repo`, `repo` |
| `GITHUB_GET_A_RELEASE_ASSET` | Get a release asset | read | `asset_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_A_RELEASE_BY_TAG_NAME` | Get a release by tag name | read | `owner`, `repo`, `tag` |  | `repo` |
| `GITHUB_GET_A_REPOSITORY` | Get a repository | read | `owner`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_GET_A_REPOSITORY_PUBLIC_KEY` | Get a repository public key | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_A_REPOSITORY_README` | Get a repository readme | read | `owner`, `repo`; *ref* |  | `repo` |
| `GITHUB_GET_A_REPOSITORY_README_FOR_A_DIRECTORY` | Get a repository readme for a directory | read | `dir`, `owner`, `repo`; *ref* |  | `public_repo`, `repo` |
| `GITHUB_GET_A_REPOSITORY_RULESET` | Get a repository ruleset | read | `owner`, `repo`, `ruleset_id`; *includes_parents* |  | `repo` |
| `GITHUB_GET_A_REPOSITORY_RULE_SUITE` | Get a repository rule suite | read | `owner`, `repo`, `rule_suite_id` |  | `admin:repo_hook`, `repo` |
| `GITHUB_GET_A_REPOSITORY_SECRET` | Get a repository secret | read | `owner`, `repo`, `secret_name` |  | `repo` |
| `GITHUB_GET_A_REPOSITORY_SUBSCRIPTION` | Get a repository subscription | read | `owner`, `repo` |  |  |
| `GITHUB_GET_A_REPOSITORY_VARIABLE` | Get a repository variable | read | `name`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_A_REPOSITORY_WEBHOOK` | Get a repository webhook | read | `hook_id`, `owner`, `repo` |  | `admin:repo_hook`, `read:repo_hook`, `repo` |
| `GITHUB_GET_A_REVIEW_COMMENT_FOR_A_PULL_REQUEST` | Get a review comment for a pull request | read | `comment_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_A_REVIEW_FOR_A_PULL_REQUEST` | Get a review for a pull request | read | `owner`, `pull_number`, `repo`, `review_id` |  |  |
| `GITHUB_GET_A_SECRET_FOR_THE_AUTHENTICATED_USER` | Get a secret for the authenticated user | read | `secret_name` |  | `codespace`, `codespace:secrets` |
| `GITHUB_GET_A_SELF_HOSTED_RUNNER_FOR_AN_ORGANIZATION` | Get a self hosted runner for an organization | read | `org`, `runner_id` |  | `admin:org` |
| `GITHUB_GET_A_SELF_HOSTED_RUNNER_FOR_A_REPOSITORY` | Get a self hosted runner for a repository | read | `owner`, `repo`, `runner_id` |  | `repo` |
| `GITHUB_GET_A_TAG` | Get a tag | read | `owner`, `repo`, `tag_sha` |  | `repo` |
| `GITHUB_GET_A_TEAM_BY_NAME` | Get a team by name | read | `org`, `team_slug` |  | `read:org` |
| `GITHUB_GET_A_THREAD` | Get a thread | read | `thread_id` |  | `notifications` |
| `GITHUB_GET_A_THREAD_SUBSCRIPTION_FOR_THE_AUTHENTICATED_USER` | Get a thread subscription for the authenticated user | read | `thread_id` |  | `notifications`, `repo` |
| `GITHUB_GET_A_TREE` | List repository files and directories | read | `owner`, `repo`, `tree_sha`; *recursive* |  | `repo` |
| `GITHUB_GET_A_USER` | Get a user | read | `username` |  | `read:user` |
| `GITHUB_GET_A_WEBHOOK_CONFIGURATION_FOR_AN_ORGANIZATION` | Get a webhook configuration for an organization | read | `hook_id`, `org` |  | `admin:org_hook` |
| `GITHUB_GET_A_WEBHOOK_CONFIGURATION_FOR_A_REPOSITORY` | Get a webhook configuration for a repository | read | `hook_id`, `owner`, `repo` |  | `admin:repo_hook`, `read:repo_hook`, `repo` |
| `GITHUB_GET_A_WEBHOOK_DELIVERY_FOR_AN_ORGANIZATION_WEBHOOK` | Get a webhook delivery for an organization webhook | read | `delivery_id`, `hook_id`, `org` |  | `admin:org_hook` |
| `GITHUB_GET_A_WORKFLOW` | Get a workflow | read | `owner`, `repo`; *workflow_id, workflow_name* |  | `repo` |
| `GITHUB_GET_A_WORKFLOW_RUN` | Get a workflow run | read | `owner`, `repo`, `run_id`; *exclude_pull_requests* |  | `repo` |
| `GITHUB_GET_A_WORKFLOW_RUN_ATTEMPT` | Get a workflow run attempt | read | `attempt_number`, `owner`, `repo`, `run_id`; *exclude_pull_requests* |  | `repo` |
| `GITHUB_GET_BILLING_USAGE_REPORT_USER` | Get billing usage report for user | read | `username`; *day, month, year* |  | `user` |
| `GITHUB_GET_BRANCH_PROTECTION` | Get branch protection | read | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_CATALOG_MODELS` | Get catalog models | read |  |  |  |
| `GITHUB_GET_CHECK_RUN` | Get a check run | read | `check_run_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_CHECK_SUITE` | Get a check suite | read | `check_suite_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_CLASSROOM` | Get a classroom | read | `classroom_id` |  |  |
| `GITHUB_GET_CODEQL_DATABASE` | Get a CodeQL database for a repository | read | `language`, `owner`, `repo` |  | `public_repo`, `repo`, `security_events` |
| `GITHUB_GET_CODESPACE` | Get a codespace for the authenticated user | read | `codespace_name` |  | `codespace` |
| `GITHUB_GET_CODE_OF_CONDUCT` | Get a code of conduct | read | `key` |  |  |
| `GITHUB_GET_CODE_OF_CONDUCT_GRAPH_QL` | Get Code of Conduct (GraphQL) | write | `key` |  |  |
| `GITHUB_GET_CODE_SCANNING_ALERT` | Get a code scanning alert | read | `alert_number`, `owner`, `repo` |  | `security_events`, `public_repo` |
| `GITHUB_GET_CODE_SCANNING_ANALYSIS` | Get a code scanning analysis for a repository | read | `analysis_id`, `owner`, `repo` |  | `security_events`, `public_repo` |
| `GITHUB_GET_CODE_SCANNING_DEFAULT_SETUP` | Get a code scanning default setup configuration | read | `owner`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_GET_COMMIT_COMMENT` | Get a commit comment | read | `comment_id`, `owner`, `repo` |  | `public_repo`, `repo` |
| `GITHUB_GET_COMMIT_OBJECT` | Get a commit object | read | `commit_sha`, `owner`, `repo` |  | `public_repo`, `repo` |
| `GITHUB_GET_COMMIT_SIGNATURE_PROTECTION` | Get commit signature protection | read | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_COMMIT_STATUSES` | Get commit statuses | read | `owner`, `ref`, `repo`; *page, per_page* |  | `public_repo`, `repo`, `repo:status` |
| `GITHUB_GET_COMMUNITY_PROFILE_METRICS` | Get community profile metrics | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_CONTEXTUAL_INFORMATION_FOR_A_USER` | Get contextual information for a user | read | `username`; *subject_id, subject_type* |  | `repo` |
| `GITHUB_GET_CUSTOM_DEPLOYMENT_PROTECTION_RULE` | Get a custom deployment protection rule | write | `environment_name`, `owner`, `protection_rule_id`, `repo` |  | `repo` |
| `GITHUB_GET_DEFAULT_ATTRIBUTES_FOR_A_CODESPACE` | Get default attributes for a codespace | read | `owner`, `repo`; *client_ip, ref* |  | `codespace` |
| `GITHUB_GET_DEFAULT_WORKFLOW_PERMISSIONS_FOR_AN_ORGANIZATION` | Get default workflow permissions for an organization | read | `org` |  | `admin:org` |
| `GITHUB_GET_DEFAULT_WORKFLOW_PERMISSIONS_FOR_A_REPOSITORY` | Get default workflow permissions for a repository | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_DEPENDABOT_ALERT` | Get Dependabot alert | read | `alert_number`, `owner`, `repo` |  | `repo`, `public_repo`, `security_events` |
| `GITHUB_GET_DEPENDENCY_DIFF` | Get a diff of the dependencies between commits | read | `basehead`, `owner`, `repo`; *name* |  | `repo` |
| `GITHUB_GET_DEPLOYMENT` | Get a deployment | read | `deployment_id`, `owner`, `repo` |  | `repo`, `repo_deployment` |
| `GITHUB_GET_DEPLOYMENT_BRANCH_POLICY` | Get a deployment branch policy | read | `branch_policy_id`, `environment_name`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_DEPLOYMENT_STATUS` | Get a deployment status | read | `deployment_id`, `owner`, `repo`, `status_id` |  | `repo_deployment` |
| `GITHUB_GET_DEPLOY_KEY` | Get a deploy key | read | `key_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_DETAILS_ABOUT_A_CODESPACE_EXPORT` | Get details about a codespace export | read | `codespace_name`, `export_id` |  | `codespace` |
| `GITHUB_GET_DISCUSSION` | Get a discussion | read | `discussion_number`, `org`, `team_slug` |  | `read:discussion` |
| `GITHUB_GET_DISCUSSION_COMMENT` | Get a discussion comment | read | `comment_number`, `discussion_number`, `org`, `team_slug` |  | `read:discussion` |
| `GITHUB_GET_EMOJIS` | Get emojis | read |  |  |  |
| `GITHUB_GET_ENTERPRISE_MEMBER_INVITATION` | Get Enterprise Member Invitation | read | `enterprise_slug`, `user_login` |  | `admin:enterprise`, `read:enterprise` |
| `GITHUB_GET_FEEDS` | Get feeds | read |  |  |  |
| `GITHUB_GET_GIST` | Get a gist | read | `gist_id` |  | `gist` |
| `GITHUB_GET_GIST_COMMENT` | Get a gist comment | read | `comment_id`, `gist_id` |  | `gist` |
| `GITHUB_GET_GIST_REVISION` | Get a gist revision | read | `gist_id`, `sha` |  | `gist` |
| `GITHUB_GET_GITHUB_ACTIONS_CACHE_USAGE_FOR_AN_ORGANIZATION` | Get GitHub Actions cache usage for an organization | read | `org` |  | `read:org` |
| `GITHUB_GET_GITHUB_ACTIONS_CACHE_USAGE_FOR_A_REPOSITORY` | Get github actions cache usage for a repository | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_GITHUB_ACTIONS_PERMISSIONS_FOR_AN_ORGANIZATION` | Get github actions permissions for an organization | read | `org` |  | `admin:org` |
| `GITHUB_GET_GITHUB_ACTIONS_PERMISSIONS_FOR_A_REPOSITORY` | Get GitHub Actions permissions for a repository | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_GITHUB_BILLING_PREMIUM_REQUEST_USAGE_REPORT_USER` | Get billing premium request usage | read | `username`; *day, model, month, product, year* |  | `read:user`, `user` |
| `GITHUB_GET_GITHUB_BILLING_USAGE_SUMMARY_FOR_USER` | Get billing usage summary for user | read | `username`; *day, month, product, repository, sku, year* |  | `read:user`, `user` |
| `GITHUB_GET_GITHUB_META_INFORMATION` | Get github meta information | read |  |  |  |
| `GITHUB_GET_GITHUB_PAGES_BUILD` | Get github pages build | read | `build_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_GITIGNORE_TEMPLATE` | Get a gitignore template | read | `name` |  |  |
| `GITHUB_GET_GLOBAL_SECURITY_ADVISORY` | Get a global security advisory | read | `ghsa_id` |  |  |
| `GITHUB_GET_GRAPHQL_NODE` | Fetch GraphQL Node | write | `id`; *fragment* |  | `public_repo`, `repo` |
| `GITHUB_GET_GRAPHQL_RATE_LIMIT` | Get GraphQL rate limit | write | *dry_run* |  |  |
| `GITHUB_GET_INFORMATION_ABOUT_A_SARIF_UPLOAD` | Get information about a sarif upload | read | `owner`, `repo`, `sarif_id` |  | `security_events`, `public_repo` |
| `GITHUB_GET_INTERACTION_RESTRICTIONS_FOR_AN_ORGANIZATION` | Get interaction restrictions for an organization | read | `org` |  | `admin:org`, `read:org` |
| `GITHUB_GET_INTERACTION_RESTRICTIONS_FOR_A_REPOSITORY` | Get interaction restrictions for a repository | read | `owner`, `repo` |  | `admin:repo_hook` |
| `GITHUB_GET_INTERACTION_RESTRICTIONS_FOR_PUBLIC_REPOS` | Get interaction restrictions for public repositories | read |  |  |  |
| `GITHUB_GET_ISSUE_COMMENTS` | Get issue comments | read | `issue_number`, `owner`, `repo`; *page, per_page, since* |  | `repo`, `public_repo` |
| `GITHUB_GET_LABEL` | Get a label | read | `name`, `owner`, `repo` |  | `public_repo`, `repo` |
| `GITHUB_GET_LATEST_PAGES_BUILD` | Get latest pages build | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_LICENSE` | Get a license | read | `license` |  |  |
| `GITHUB_GET_LICENSE_GRAPH_QL` | Get License (GraphQL) | write | `key` |  |  |
| `GITHUB_GET_OCTOCAT` | Get octocat | read | *s* |  |  |
| `GITHUB_GET_OIDC_SUBJECT_CLAIM_TEMPLATE` | Customize OIDC subject claim template | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_ORGANIZATION_AUDIT_LOG` | Get organization audit log | read | `org`; *after, before, include, order, per_page, phrase* |  | `read:audit_log` |
| `GITHUB_GET_ORGANIZATION_CUSTOM_PROPERTY` | Get a custom property for an organization | write | `custom_property_name`, `org` |  | `admin:org`, `read:org` |
| `GITHUB_GET_ORGANIZATION_MEMBERSHIP_FOR_USER` | Get organization membership for user | read | `org`, `username` |  | `read:org` |
| `GITHUB_GET_ORGANIZATION_PROJECT_ITEM` | Get organization project item | read | `item_id`, `org`, `project_number`; *fields* |  | `project`, `read:project` |
| `GITHUB_GET_ORGANIZATION_PUBLIC_KEY` | Get organization public key for Codespaces secrets | read | `org` |  | `admin:org` |
| `GITHUB_GET_ORGANIZATION_PUBLIC_KEY_FOR_DEPENDABOT` | Fetch Public Key For Secret Encryption | read | `org` |  | `admin:org` |
| `GITHUB_GET_ORGANIZATION_TEAMS` | Get organization teams | read | `org`; *page, per_page, team_type* |  | `read:org` |
| `GITHUB_GET_ORG_ALLOWED_ACTIONS` | Get allowed actions and workflows for an org | read | `org` |  | `admin:org`, `read:org` |
| `GITHUB_GET_ORG_DEV_ENVIRONMENT_SECRET_SAFELY` | Get org dev environment secret safely | read | `org`, `secret_name` |  | `admin:org` |
| `GITHUB_GET_ORG_OIDC_SUBJECT_CLAIM_TEMPLATE` | Get custom OIDC subject claim template | read | `org` |  | `read:org` |
| `GITHUB_GET_PAGES_DNS_HEALTH_CHECK` | Get a dns health check for github pages | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_PAGES_SITE` | Get a GitHub Pages site | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_PAGE_VIEWS` | Get page views | read | `owner`, `repo`; *per* |  | `repo` |
| `GITHUB_GET_PARENT_ISSUE` | Get parent issue | read | `issue_number`, `owner`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_GET_PENDING_DEPLOYMENTS_FOR_A_WORKFLOW_RUN` | Get pending deployments for a workflow run | read | `owner`, `repo`, `run_id` |  | `repo` |
| `GITHUB_GET_PROJECT_PERMISSION_FOR_A_USER` | Get project permission for a user | read | `project_id`, `username` |  | `repo`, `project`, `read:project` |
| `GITHUB_GET_PUBLIC_KEY_FOR_SECRET_ENCRYPTION` | Get public key for secret encryption | read | `owner`, `repo` |  | `codespace`, `repo` |
| `GITHUB_GET_PUBLIC_KEY_FOR_THE_AUTHENTICATED_USER` | Get public key for the authenticated user | read |  |  | `codespace`, `codespace:secrets` |
| `GITHUB_GET_PULL_REQUESTS` | Get Pull Requests | read | `owner`, `repo`; *base, direction, head, page, per_page, sort, state* |  | `repo`, `public_repo` |
| `GITHUB_GET_PULL_REQUEST_REVIEW_PROTECTION` | Get pull request review protection | read | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_RATE_LIMIT_STATUS_FOR_THE_AUTHENTICATED_USER` | Get rate limit status for the authenticated user | read |  |  |  |
| `GITHUB_GET_RAW_REPOSITORY_CONTENT` | Get Raw Repository Content | write | `owner`, `path`, `repo`; *ref* |  |  |
| `GITHUB_GET_REPOSITORY_CLONES` | Get repository clones | read | `owner`, `repo`; *per* |  | `repo` |
| `GITHUB_GET_REPOSITORY_CONTENT` | Get repository content | read | `owner`, `path`, `repo`; *ref* |  | `public_repo`, `repo`, `workflow` |
| `GITHUB_GET_REPOSITORY_DISCUSSION_THREAD` | Get repository Discussion thread | read | `discussion_number`, `owner`, `repo`; *comments_after, comments_first, replies_first* |  | `repo` |
| `GITHUB_GET_REPOSITORY_GRAPH_QL` | Get Repository via GraphQL | write | `name`, `owner`; *followRenames* |  | `public_repo`, `repo` |
| `GITHUB_GET_REPOSITORY_OWNER` | Get Repository Owner | read | `login` |  | `read:user`, `read:org` |
| `GITHUB_GET_REPOSITORY_PERMISSIONS_FOR_A_USER` | Get repository permissions for a user | read | `owner`, `repo`, `username` |  | `repo` |
| `GITHUB_GET_REPOSITORY_PUBLIC_KEY_FOR_ENCRYPTION` | Retrieve repo public key for encryption | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_REPOSITORY_SECRET_SECURELY` | Get repository secret securely | read | `owner`, `repo`, `secret_name` |  | `repo` |
| `GITHUB_GET_REPOSITORY_SECURITY_ADVISORY` | Get repository security advisory | read | `ghsa_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_REPOSITORY_STAR_HISTORY` | Get repository star history | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_GET_REPOSITORY_WEBHOOK_DELIVERY` | Get a delivery for a repository webhook | read | `delivery_id`, `hook_id`, `owner`, `repo` |  |  |
| `GITHUB_GET_REPOS_ATTESTATIONS` | Get repository attestations | read | `owner`, `repo`, `subject_digest`; *after, before, per_page, predicate_type* |  | `repo` |
| `GITHUB_GET_REPO_ALLOWED_ACTIONS` | Get allowed actions and workflows for a repository | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_REPO_DEV_ENV_SECRET` | Get repository development environment secret | read | `owner`, `repo`, `secret_name` |  | `repo` |
| `GITHUB_GET_RULES_FOR_A_BRANCH` | Get rules for a branch | read | `branch`, `owner`, `repo`; *page, per_page* |  |  |
| `GITHUB_GET_SINGLE_ORG_SECRET_WITHOUT_DECRYPTION` | Get Organization Dependabot Secret | read | `org`, `secret_name` |  | `admin:org` |
| `GITHUB_GET_STATUS_CHECKS_PROTECTION` | Get status checks protection | read | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_TEAMS_WITH_ACCESS_TO_THE_PROTECTED_BRANCH` | Get teams with access to the protected branch | read | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_TEAM_MEMBERS` | Get team members | read | `org`, `team_slug`; *page, per_page, role* |  | `read:org` |
| `GITHUB_GET_TEAM_MEMBERSHIP_FOR_A_USER` | Get team membership for a user | read | `org`, `team_slug`, `username` |  | `read:org` |
| `GITHUB_GET_THE_AUTHENTICATED_USER` | Get the authenticated user | read |  |  |  |
| `GITHUB_GET_THE_COMBINED_STATUS_FOR_A_SPECIFIC_REFERENCE` | Get the combined status for a specific reference | read | `owner`, `ref`, `repo`; *page, per_page* |  | `repo:status` |
| `GITHUB_GET_THE_HOURLY_COMMIT_COUNT_FOR_EACH_DAY` | Get the hourly commit count for each day | read | `owner`, `repo` |  |  |
| `GITHUB_GET_THE_LAST_YEAR_OF_COMMIT_ACTIVITY` | Get the last year of commit activity | read | `owner`, `repo` |  | `public_repo`, `repo` |
| `GITHUB_GET_THE_LATEST_RELEASE` | Get the latest release | read | `owner`, `repo` |  | `public_repo`, `repo` |
| `GITHUB_GET_THE_LICENSE_FOR_A_REPOSITORY` | Get the license for a repository | read | `owner`, `repo`; *ref* |  | `public_repo` |
| `GITHUB_GET_THE_REVIEW_HISTORY_FOR_A_WORKFLOW_RUN` | Get the review history for a workflow run | read | `owner`, `repo`, `run_id` |  | `repo` |
| `GITHUB_GET_THE_STATUS_OF_A_GITHUB_PAGES_DEPLOYMENT` | Get the status of a GitHub Pages deployment | read | `owner`, `pages_deployment_id`, `repo` |  | `repo` |
| `GITHUB_GET_THE_WEEKLY_COMMIT_ACTIVITY` | Get the weekly commit activity | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_THE_WEEKLY_COMMIT_COUNT` | Get the weekly commit count | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_THE_ZEN_OF_GITHUB` | Get the Zen of GitHub | read |  |  |  |
| `GITHUB_GET_TOP_REFERRAL_PATHS` | Get top referral paths | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_TOP_REFERRAL_SOURCES` | Get top referral sources | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_USERS_ATTESTATIONS` | Get user attestations by subject digest | read | `subject_digest`, `username`; *after, before, per_page, predicate_type* |  | `user` |
| `GITHUB_GET_USERS_PROJECTS_V2` | Get user's projects v2 | read | `username`; *after, before, per_page, q* |  | `project`, `read:project`, `repo` |
| `GITHUB_GET_USERS_WITH_ACCESS_TO_THE_PROTECTED_BRANCH` | Get users with access to the protected branch | read | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_USER_BY_ID` | Get user by ID | read | `account_id` |  | `read:user` |
| `GITHUB_GET_USER_PROJECT_ITEM` | Get user project item | read | `item_id`, `project_number`, `username`; *fields* |  | `project`, `read:project`, `repo` |
| `GITHUB_GET_VIEWER_GRAPHQL` | Get Viewer via GraphQL | write |  |  |  |
| `GITHUB_GET_WORKFLOW_EXTERNAL_ACCESS` | Get workflow external access level | read | `owner`, `repo` |  | `repo` |
| `GITHUB_GET_WORKFLOW_RUN_JOB` | Get a job for a workflow run | read | `job_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_GET_WORKFLOW_RUN_USAGE` | Get workflow run usage | read | `owner`, `repo`, `run_id` |  | `repo` |
| `GITHUB_GET_WORKFLOW_USAGE` | Get workflow usage | read | `owner`, `repo`, `workflow_id` |  | `repo` |
| `GITHUB_IS_REPOSITORY_STARRED_BY_THE_USER` | Check if repo starred by auth user | read | `owner`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_LIST_ACCEPTED_ASSIGNMENTS_FOR_AN_ASSIGNMENT` | List accepted assignments for an assignment | read | `assignment_id`; *page, per_page* |  |  |
| `GITHUB_LIST_ACCESSIBLE_REPOSITORIES` | List repositories accessible via GitHub App installation | read | `installation_id`; *page, per_page* |  | `repo`, `public_repo` |
| `GITHUB_LIST_APPS_WITH_ACCESS_TO_THE_PROTECTED_BRANCH` | Get apps with access to the protected branch | read | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_LIST_APP_INSTALLATIONS` | List app installations for user token | read | *page, per_page* |  | `user` |
| `GITHUB_LIST_ARTIFACTS_FOR_A_REPOSITORY` | List artifacts for a repository | read | `owner`, `repo`; *name, page, per_page* |  | `repo` |
| `GITHUB_LIST_ASSIGNEES` | List assignees | read | `owner`, `repo`; *page, per_page* |  | `repo`, `public_repo` |
| `GITHUB_LIST_ASSIGNMENTS_FOR_A_CLASSROOM` | List assignments for a classroom | read | `classroom_id`; *page, per_page* |  |  |
| `GITHUB_LIST_ATTESTATION_REPOSITORIES` | List attestation repositories | read | `org`; *after, before, per_page* |  | `read:org` |
| `GITHUB_LIST_AUTHENTICATED_USER_SSH_KEYS` | List authenticated user SSH keys | read | *page, per_page* |  | `read:public_key` |
| `GITHUB_LIST_AVAILABLE_MACHINE_TYPES_FOR_A_REPOSITORY` | List available machine types for a repository | read | `owner`, `repo`; *client_ip, location, ref* |  | `codespace` |
| `GITHUB_LIST_BRANCHES` | List branches | read | `owner`, `repo`; *page, per_page, protected* |  | `repo`, `public_repo` |
| `GITHUB_LIST_BRANCHES_FOR_HEAD_COMMIT` | List branches for head commit | read | `commit_sha`, `owner`, `repo` |  |  |
| `GITHUB_LIST_CHECK_RUNS_FOR_A_REF` | List check runs for a git ref | read | `owner`, `ref`, `repo`; *app_id, check_name, filter, page, per_page, status* |  | `public_repo`, `repo` |
| `GITHUB_LIST_CHECK_RUNS_IN_A_CHECK_SUITE` | List check runs in a check suite | read | `check_suite_id`, `owner`, `repo`; *check_name, filter, page, per_page, status* |  | `repo` |
| `GITHUB_LIST_CHECK_RUN_ANNOTATIONS` | List check run annotations | read | `check_run_id`, `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_CHECK_SUITES_FOR_A_GIT_REFERENCE` | List check suites for a git reference | read | `owner`, `ref`, `repo`; *app_id, check_name, page, per_page* |  | `repo` |
| `GITHUB_LIST_CHILD_TEAMS` | List child teams | read | `org`, `team_slug`; *page, per_page* |  | `read:org` |
| `GITHUB_LIST_CLASSROOMS` | List classrooms | read | *page, per_page* |  |  |
| `GITHUB_LIST_CODEOWNERS_ERRORS` | List codeowners errors | read | `owner`, `repo`; *ref* |  | `repo` |
| `GITHUB_LIST_CODEQL_DATABASES_FOR_A_REPOSITORY` | List CodeQL databases for a repository | read | `owner`, `repo` |  | `security_events`, `repo`, `public_repo` |
| `GITHUB_LIST_CODESPACES_FOR_A_USER_IN_ORGANIZATION` | List codespaces for a user in organization | read | `org`, `username`; *page, per_page* |  | `admin:org` |
| `GITHUB_LIST_CODESPACES_FOR_THE_AUTHENTICATED_USER` | List codespaces for the authenticated user | read | *page, per_page, repository_id* |  | `codespace` |
| `GITHUB_LIST_CODESPACES_FOR_THE_ORGANIZATION` | List codespaces for the organization | read | `org`; *page, per_page* |  | `admin:org` |
| `GITHUB_LIST_CODES_OF_CONDUCT` | Get all codes of conduct | read |  |  |  |
| `GITHUB_LIST_CODES_OF_CONDUCT_GRAPHQL` | List Codes of Conduct (GraphQL) | write |  |  |  |
| `GITHUB_LIST_CODE_SCANNING_ALERTS_FOR_AN_ORGANIZATION` | List code scanning alerts for an organization | read | `org`; *after, before, direction, page, per_page, severity, sort, state and 2 more* |  | `security_events`, `repo`, `public_repo` |
| `GITHUB_LIST_CODE_SCANNING_ALERTS_FOR_A_REPOSITORY` | List code scanning alerts for a repository | read | `owner`, `repo`; *direction, page, per_page, ref, severity, sort, state, tool_guid and 1 more* |  | `security_events`, `public_repo` |
| `GITHUB_LIST_CODE_SCANNING_ANALYSES_FOR_A_REPOSITORY` | List code scanning analyses for a repository | read | `owner`, `repo`; *direction, page, per_page, ref, sarif_id, sort, tool_guid, tool_name* |  | `security_events`, `public_repo` |
| `GITHUB_LIST_COMMENTS_FOR_A_PULL_REQUEST_REVIEW` | List comments for a pull request review | read | `owner`, `pull_number`, `repo`, `review_id`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_COMMENT_CHANGES` | List comment changes across organization | read | `org`; *page, per_page* |  | `public_repo`, `repo` |
| `GITHUB_LIST_COMMITS` | List commits | read | `owner`, `repo`; *author, committer, page, path, per_page, sha, since, until* |  | `repo` |
| `GITHUB_LIST_COMMITS_ON_A_PULL_REQUEST` | List commits on a pull request | read | `owner`, `pull_number`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_COMMIT_COMMENTS` | List commit comments | read | `commit_sha`, `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_COMMIT_COMMENTS_FOR_A_REPOSITORY` | List commit comments for a repository | read | `owner`, `repo`; *page, per_page* |  |  |
| `GITHUB_LIST_COMMONLY_USED_LICENSES` | Get all commonly used licenses | read | *featured, page, per_page* |  | `public_repo` |
| `GITHUB_LIST_CUSTOM_PROPERTIES_FOR_AN_ORGANIZATION` | Get all custom properties for an organization | write | `org` |  | `admin:org`, `read:org`, `write:org` |
| `GITHUB_LIST_CUSTOM_PROPERTY_VALUES_FOR_A_REPOSITORY` | Get all custom property values for a repository | write | `owner`, `repo` |  |  |
| `GITHUB_LIST_CUSTOM_PROPERTY_VALUES_FOR_ORG_REPOS` | List custom property values for organization repositories | write | `org`; *page, per_page, repository_query* |  | `admin:org`, `read:org` |
| `GITHUB_LIST_DELIVERIES_FOR_AN_ORGANIZATION_WEBHOOK` | List deliveries for an organization webhook | read | `hook_id`, `org`; *cursor, per_page, redelivery* |  | `admin:org_hook` |
| `GITHUB_LIST_DELIVERIES_FOR_A_REPOSITORY_WEBHOOK` | List deliveries for a repository webhook | read | `hook_id`, `owner`, `repo`; *cursor, per_page, redelivery* |  | `admin:repo_hook`, `read:repo_hook`, `repo` |
| `GITHUB_LIST_DEPLOYMENTS` | List deployments | read | `owner`, `repo`; *environment, page, per_page, ref, sha, task* |  | `repo_deployment` |
| `GITHUB_LIST_DEPLOYMENT_BRANCH_POLICIES` | List deployment branch policies | read | `environment_name`, `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_DEPLOYMENT_STATUSES` | List deployment statuses | read | `deployment_id`, `owner`, `repo`; *page, per_page* |  | `repo`, `repo_deployment` |
| `GITHUB_LIST_DEPLOY_KEYS` | List deploy keys | read | `owner`, `repo`; *page, per_page* |  | `admin:public_key`, `repo` |
| `GITHUB_LIST_EMAIL_ADDRESSES_FOR_THE_AUTHENTICATED_USER` | List email addresses for the authenticated user | read | *page, per_page* |  | `user:email` |
| `GITHUB_LIST_ENVIRONMENTS` | List environments | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_ENVIRONMENT_CUSTOM_DEPLOYMENT_RULES` | List environment custom deployment rules | write | `environment_name`, `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_ENVIRONMENT_SECRETS` | List environment secrets | read | `environment_name`, `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_ENVIRONMENT_VARIABLES` | List environment variables | read | `environment_name`, `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_EVENTS_FOR_THE_AUTHENTICATED_USER` | List events for the authenticated user | read | `username`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_EVENTS_FOR_USER` | List events for user | read | `username`; *page, per_page* |  |  |
| `GITHUB_LIST_EVENTS_RECEIVED_BY_THE_AUTHENTICATED_USER` | List events received by the authenticated user | read | `username`; *page, per_page* |  |  |
| `GITHUB_LIST_FOLLOWERS_OF_A_USER` | List followers of a user | read | `username`; *page, per_page* |  |  |
| `GITHUB_LIST_FOLLOWERS_OF_THE_AUTHENTICATED_USER` | List followers of the authenticated user | read | *page, per_page* |  | `read:user`, `user:follow` |
| `GITHUB_LIST_FOLLOWING_FOR_USER` | List following for user | read | `username`; *page, per_page* |  |  |
| `GITHUB_LIST_FORKS` | List forks | read | `owner`, `repo`; *page, per_page, sort* |  | `public_repo` |
| `GITHUB_LIST_GISTS_FOR_A_USER` | List gists for a user | read | `username`; *page, per_page, since* |  | `gist` |
| `GITHUB_LIST_GISTS_FOR_THE_AUTHENTICATED_USER` | List gists for the authenticated user | read | *page, per_page, since* |  | `gist` |
| `GITHUB_LIST_GIST_COMMENTS` | List gist comments | read | `gist_id`; *page, per_page* |  | `gist` |
| `GITHUB_LIST_GIST_COMMITS` | List gist commits | read | `gist_id`; *page, per_page* |  | `gist` |
| `GITHUB_LIST_GIST_FORKS` | List gist forks | read | `gist_id`; *page, per_page* |  | `gist` |
| `GITHUB_LIST_GITHUB_ACTIONS_CACHES_FOR_A_REPOSITORY` | List github actions caches for a repository | read | `owner`, `repo`; *direction, key, page, per_page, ref, sort* |  | `repo` |
| `GITHUB_LIST_GITHUB_PAGES_BUILDS` | List GitHub Pages builds | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_GITIGNORE_TEMPLATES` | Get all gitignore templates | read |  |  | `public_repo` |
| `GITHUB_LIST_GLOBAL_SECURITY_ADVISORIES` | List global security advisories | read | *affects, after, before, cve_id, cwes, direction, ecosystem, epss_percentage and 10 more* |  | `repo` |
| `GITHUB_LIST_GPG_KEYS_FOR_A_USER` | List GPG keys for a user | read | `username`; *page, per_page* |  |  |
| `GITHUB_LIST_INSTALLATION_REQUESTS_FOR_THE_AUTHENTICATED_APP` | List installation requests for the authenticated app | read | *page, per_page* |  |  |
| `GITHUB_LIST_INSTANCES_OF_A_CODE_SCANNING_ALERT` | List instances of a code scanning alert | read | `alert_number`, `owner`, `repo`; *page, per_page, ref* |  | `public_repo`, `security_events` |
| `GITHUB_LIST_ISSUES_ASSIGNED_TO_THE_AUTHENTICATED_USER` | List issues for user | read | *direction, filter, labels, page, per_page, since, sort, state* |  | `repo` |
| `GITHUB_LIST_ISSUE_COMMENTS` | List issue comments | read | `issue_number`, `owner`, `repo`; *page, per_page, since* |  | `repo`, `public_repo` |
| `GITHUB_LIST_ISSUE_COMMENTS_FOR_A_REPOSITORY` | List issue comments for a repository | read | `owner`, `repo`; *direction, page, per_page, since, sort* |  | `repo` |
| `GITHUB_LIST_ISSUE_DEPENDENCIES_BLOCKED_BY` | List issue dependencies blocked by | read | `issue_number`, `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_ISSUE_DEPENDENCIES_BLOCKING` | List issue dependencies blocking | read | `issue_number`, `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_ISSUE_EVENTS` | List issue events | read | `issue_number`, `owner`, `repo`; *page, per_page* |  | `public_repo`, `repo` |
| `GITHUB_LIST_ISSUE_EVENTS_FOR_A_REPOSITORY` | List issue events for a repository | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_ISSUE_TYPES` | List issue types | read | `org` |  | `read:org` |
| `GITHUB_LIST_JOBS_FOR_A_WORKFLOW_RUN` | List jobs for a workflow run | read | `owner`, `repo`, `run_id`; *filter, page, per_page* |  | `repo` |
| `GITHUB_LIST_JOBS_FOR_A_WORKFLOW_RUN_ATTEMPT` | List jobs for a workflow run attempt | read | `attempt_number`, `owner`, `repo`, `run_id`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_LABELS_FOR_AN_ISSUE` | List labels for an issue | read | `issue_number`, `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_LABELS_FOR_A_REPOSITORY` | List labels for a repository | read | `owner`, `repo`; *page, per_page* |  | `public_repo`, `repo` |
| `GITHUB_LIST_LABELS_FOR_A_SELF_HOSTED_RUNNER_FOR_A_REPOSITORY` | List labels for a self-hosted runner for a repository | read | `owner`, `repo`, `runner_id` |  | `repo` |
| `GITHUB_LIST_LABELS_FOR_ISSUES_IN_A_MILESTONE` | List labels for issues in a milestone | read | `milestone_number`, `owner`, `repo`; *page, per_page* |  | `public_repo`, `repo` |
| `GITHUB_LIST_LABELS_FOR_SELF_HOSTED_RUNNER_ORG` | List labels for a self-hosted runner for an organization | read | `org`, `runner_id` |  | `admin:org` |
| `GITHUB_LIST_LICENSES_GRAPH_QL` | List Licenses (GraphQL) | write |  |  |  |
| `GITHUB_LIST_MACHINE_TYPES_FOR_A_CODESPACE` | List machine types for a codespace | read | `codespace_name` |  | `codespace` |
| `GITHUB_LIST_MATCHING_REFERENCES` | List matching references | read | `owner`, `ref`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_MILESTONES` | List milestones | read | `owner`, `repo`; *direction, page, per_page, sort, state* |  | `repo` |
| `GITHUB_LIST_NOTIFICATIONS` | List notifications for authenticated user | read | *all, before, page, participating, per_page, since* |  | `notifications`, `repo` |
| `GITHUB_LIST_ORGANIZATIONS` | List organizations | read | *per_page, since* |  | `read:org`, `user` |
| `GITHUB_LIST_ORGANIZATIONS_FOR_A_USER` | List organizations for a user | read | `username`; *page, per_page* |  | `read:org`, `user` |
| `GITHUB_LIST_ORGANIZATIONS_FOR_THE_AUTHENTICATED_USER` | List organizations for the authenticated user | read | *page, per_page* |  | `user`, `read:org` |
| `GITHUB_LIST_ORGANIZATIONS_FOR_USER` | List organizations for a user | read | `username`; *page, per_page* |  |  |
| `GITHUB_LIST_ORGANIZATION_APP_INSTALLATIONS` | List organization app installations | read | `org`; *page, per_page* |  | `read:org` |
| `GITHUB_LIST_ORGANIZATION_CODESPACES_SECRETS` | List org level codespaces secrets | read | `org`; *page, per_page* |  | `admin:org` |
| `GITHUB_LIST_ORGANIZATION_EVENTS_FOR_THE_AUTHENTICATED_USER` | List organization events for the authenticated user | read | `org`, `username`; *page, per_page* |  | `read:org` |
| `GITHUB_LIST_ORGANIZATION_MEMBERS` | List organization members | read | `org`; *filter, page, per_page, role* |  | `read:org`, `admin:org` |
| `GITHUB_LIST_ORGANIZATION_MEMBERSHIPS_FOR_THE_AUTHENTICATED_U` | List organization memberships for authenticated user | read | *page, per_page, state* |  | `read:org`, `user` |
| `GITHUB_LIST_ORGANIZATION_PROJECTS` | List organization projects | read | `org`; *after, before, per_page, q* |  | `read:org`, `repo`, `project` |
| `GITHUB_LIST_ORGANIZATION_REPOSITORIES` | List organization repositories | read | `org`; *direction, page, per_page, sort, type* |  | `repo`, `read:org` |
| `GITHUB_LIST_ORGANIZATION_ROLES_FOR_AN_ORGANIZATION` | Get all organization roles for an organization | read | `org` |  | `admin:org` |
| `GITHUB_LIST_ORGANIZATION_SECRETS` | List organization secrets | read | `org`; *page, per_page* |  | `admin:org` |
| `GITHUB_LIST_ORGANIZATION_VARIABLES` | List organization variables | read | `org`; *page, per_page* |  | `admin:org` |
| `GITHUB_LIST_ORGANIZATION_WEBHOOKS` | List organization webhooks | read | `org`; *page, per_page* |  | `admin:org_hook`, `read:org` |
| `GITHUB_LIST_ORG_FINE_GRAINED_PERMISSIONS` | List organization fine-grained permissions | read | `org` |  | `admin:org` |
| `GITHUB_LIST_ORG_ISSUES_FOR_THE_AUTHENTICATED_USER` | List organization issues assigned to the authenticated user | read | `org`; *direction, filter, labels, page, per_page, since, sort, state* |  | `repo`, `public_repo`, `read:org` |
| `GITHUB_LIST_ORG_PACKAGE_VERSIONS` | List package versions for an organization package | read | `org`, `package_name`, `package_type`; *page, per_page, state* |  | `read:packages` |
| `GITHUB_LIST_ORG_PROJECT_VIEW_ITEMS` | List organization project view items | read | `org`, `project_number`, `view_number`; *after, before, fields, per_page* |  | `project`, `read:project`, `repo` |
| `GITHUB_LIST_ORG_REPOS_WITH_GITHUB_ACTIONS_ENABLED` | List organization repositories with Github Actions enabled | read | `org`; *page, per_page* |  | `admin:org` |
| `GITHUB_LIST_ORG_RESOURCES_WITH_PERSONAL_TOKENS` | List org resources with personal tokens | read | `org`; *direction, last_used_after, last_used_before, owner, page, per_page, permission, repository and 2 more* |  | `admin:org`, `read:org` |
| `GITHUB_LIST_ORG_RESOURCE_ACCESS_TOKENS` | List org resource access tokens | read | `org`; *direction, last_used_after, last_used_before, owner, page, per_page, permission, repository and 2 more* |  | `admin:org`, `read:org` |
| `GITHUB_LIST_ORG_SECRETS_WITHOUT_VALUES` | List organization Dependabot secrets without values | read | `org`; *page, per_page* |  | `admin:org` |
| `GITHUB_LIST_OUTSIDE_COLLABORATORS_FOR_AN_ORGANIZATION` | List outside collaborators for an organization | read | `org`; *filter, page, per_page* |  | `read:org` |
| `GITHUB_LIST_OWNED_PACKAGE_VERSIONS` | List owned package versions | read | `package_name`, `package_type`; *page, per_page, state* |  | `delete:packages`, `read:packages`, `write:packages` |
| `GITHUB_LIST_PACKAGES_FOR_AN_ORGANIZATION` | List packages for an organization | read | `org`, `package_type`; *page, per_page, visibility* |  | `read:packages` |
| `GITHUB_LIST_PACKAGES_FOR_A_USER` | List packages for a user | read | `package_type`, `username`; *page, per_page, visibility* |  | `read:packages` |
| `GITHUB_LIST_PACKAGES_FOR_THE_AUTHENTICATED_USER` | List packages for the authenticated user's namespace | read | `package_type`; *page, per_page, visibility* |  | `delete:packages`, `read:packages`, `write:packages` |
| `GITHUB_LIST_PACKAGE_VERSIONS_FOR_A_PACKAGE_OWNED_BY_A_USER` | List package versions for a package owned by a user | read | `package_name`, `package_type`, `username`; *page, per_page* |  | `read:packages` |
| `GITHUB_LIST_PENDING_TEAM_INVITATIONS` | List pending team invitations | read | `org`, `team_slug`; *page, per_page* |  | `read:org` |
| `GITHUB_LIST_PROJECT_COLLABORATORS` | List project collaborators (classic) | read | `project_id`; *affiliation, page, per_page* |  | `repo`, `public_repo` |
| `GITHUB_LIST_PROJECT_COLUMNS` | List project columns | read | `project_id`; *page, per_page* |  | `repo`, `public_repo` |
| `GITHUB_LIST_PROJECT_FIELDS_FOR_USER` | List project fields for user | read | `project_number`, `username`; *after, before, per_page* |  | `project`, `read:project` |
| `GITHUB_LIST_PROJECT_ITEMS_FOR_USER` | List project items for user | read | `project_number`, `username`; *after, before, fields, per_page, q* |  | `project`, `read:project` |
| `GITHUB_LIST_PROJECT_VIEW_ITEMS_FOR_USER` | List items for user project view | read | `project_number`, `username`, `view_number`; *after, before, fields, per_page* |  | `read:project` |
| `GITHUB_LIST_PUBLIC_EMAILS_FOR_AUTH_USER` | List public email addresses for the authenticated user | read | *page, per_page* |  | `user:email` |
| `GITHUB_LIST_PUBLIC_EVENTS` | List public events | read | *page, per_page* |  | `repo` |
| `GITHUB_LIST_PUBLIC_EVENTS_FOR_A_NETWORK_OF_REPOSITORIES` | List public events for a network of repositories | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_PUBLIC_EVENTS_FOR_A_USER` | List public events for a user | read | `username`; *page, per_page* |  |  |
| `GITHUB_LIST_PUBLIC_EVENTS_RECEIVED_BY_A_USER` | List public events received by a user | read | `username`; *page, per_page* |  |  |
| `GITHUB_LIST_PUBLIC_GISTS` | List public gists | read | *page, per_page, since* |  | `gist` |
| `GITHUB_LIST_PUBLIC_KEYS_FOR_A_USER` | List public keys for a user | read | `username`; *page, per_page* |  | `read:public_key` |
| `GITHUB_LIST_PUBLIC_ORGANIZATION_EVENTS` | List public organization events | read | `org`; *page, per_page* |  |  |
| `GITHUB_LIST_PUBLIC_ORGANIZATION_MEMBERS` | List public organization members | read | `org`; *page, per_page* |  | `read:org` |
| `GITHUB_LIST_PUBLIC_REPOSITORIES` | List public repositories | read | *since* |  | `public_repo`, `repo` |
| `GITHUB_LIST_PULL_REQUESTS` | List pull requests | read | `owner`, `repo`; *base, direction, head, page, per_page, sort, state* |  | `repo` |
| `GITHUB_LIST_PULL_REQUESTS_ASSOCIATED_WITH_A_COMMIT` | List pull requests associated with a commit | read | `commit_sha`, `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_PULL_REQUESTS_FILES` | List pull requests files | read | `owner`, `pull_number`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_REACTIONS_FOR_AN_ISSUE` | List reactions for an issue | read | `issue_number`, `owner`, `repo`; *content, page, per_page* |  | `repo`, `public_repo` |
| `GITHUB_LIST_REACTIONS_FOR_AN_ISSUE_COMMENT` | List reactions for an issue comment | read | `comment_id`, `owner`, `repo`; *content, page, per_page* |  | `public_repo`, `repo` |
| `GITHUB_LIST_REACTIONS_FOR_A_COMMIT_COMMENT` | List reactions for a commit comment | read | `comment_id`, `owner`, `repo`; *content, page, per_page* |  | `public_repo`, `repo` |
| `GITHUB_LIST_REACTIONS_FOR_A_PULL_REQUEST_REVIEW_COMMENT` | List reactions for a pull request review comment | read | `comment_id`, `owner`, `repo`; *content, page, per_page* |  | `repo` |
| `GITHUB_LIST_REACTIONS_FOR_A_RELEASE` | List reactions for a release | read | `owner`, `release_id`, `repo`; *content, page, per_page* |  | `repo`, `public_repo` |
| `GITHUB_LIST_REACTIONS_FOR_A_TEAM_DISCUSSION` | List reactions for a team discussion | read | `discussion_number`, `org`, `team_slug`; *content, page, per_page* |  | `read:discussion` |
| `GITHUB_LIST_REACTIONS_FOR_A_TEAM_DISCUSSION_COMMENT` | List reactions for a team discussion comment | read | `comment_number`, `discussion_number`, `org`, `team_slug`; *content, page, per_page* |  | `read:discussion` |
| `GITHUB_LIST_RELEASES` | List releases | read | `owner`, `repo`; *page, per_page* |  | `repo`, `public_repo` |
| `GITHUB_LIST_RELEASE_ASSETS` | List release assets | read | `owner`, `release_id`, `repo`; *page, per_page* |  | `public_repo`, `repo` |
| `GITHUB_LIST_REPOSITORIES` | List repositories | read | *affiliation, before, direction, page, per_page, since, sort, type and 1 more* |  | `repo`, `public_repo` |
| `GITHUB_LIST_REPOSITORIES_ACCESSIBLE_TO_THE_APP_INSTALLATION` | List repos accessible to app installation | read | *page, per_page* |  |  |
| `GITHUB_LIST_REPOSITORIES_FOR_A_USER` | List repositories for a user | read | `username`; *direction, page, per_page, sort, type* |  | `repo` |
| `GITHUB_LIST_REPOSITORIES_FOR_THE_AUTHENTICATED_USER` | List repositories for the authenticated user | read | *before, direction, page, per_page, since, sort, type* |  | `public_repo`, `repo` |
| `GITHUB_LIST_REPOSITORIES_STARRED_BY_A_USER` | List repositories starred by a user | read | `username`; *direction, page, per_page, sort* |  | `read:user`, `repo` |
| `GITHUB_LIST_REPOSITORIES_STARRED_BY_THE_AUTHENTICATED_USER` | List repositories starred by the authenticated user | read | *direction, page, per_page, sort* |  | `read:user`, `repo` |
| `GITHUB_LIST_REPOSITORIES_WATCHED_BY_A_USER` | List repositories watched by a user | read | `username`; *page, per_page* |  |  |
| `GITHUB_LIST_REPOSITORIES_WATCHED_BY_THE_AUTHENTICATED_USER` | List repositories watched by the authenticated user | read | *page, per_page* |  | `public_repo`, `read:user`, `repo` |
| `GITHUB_LIST_REPOSITORIES_WITH_GITHUB_ACTIONS_CACHE_USAGE` | List repos with ghactions cache usage | read | `org`; *page, per_page* |  | `read:org` |
| `GITHUB_LIST_REPOSITORY_ACTIVITIES` | List repository activities | read | `owner`, `repo`; *activity_type, actor, after, before, direction, per_page, ref, time_period* |  | `repo` |
| `GITHUB_LIST_REPOSITORY_AUTOLINKS` | Get all autolinks of a repository | read | `owner`, `repo` |  | `repo` |
| `GITHUB_LIST_REPOSITORY_COLLABORATORS` | List repository collaborators | read | `owner`, `repo`; *affiliation, page, per_page, permission* |  | `read:org`, `repo` |
| `GITHUB_LIST_REPOSITORY_CONTRIBUTORS` | List repository contributors | read | `owner`, `repo`; *anon, page, per_page* |  | `repo` |
| `GITHUB_LIST_REPOSITORY_DEPENDABOT_ALERTS` | List repository Dependabot alerts | read | `owner`, `repo`; *after, before, direction, ecosystem, manifest, package, per_page, relationship and 4 more* |  | `repo`, `public_repo`, `security_events` |
| `GITHUB_LIST_REPOSITORY_EVENTS` | List repository events | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_REPOSITORY_INVITATIONS` | List repository invitations | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_REPOSITORY_ISSUES` | List repository issues | read | `owner`, `repo`; *assignee, creator, direction, labels, mentioned, milestone, page, per_page and 4 more* |  | `repo` |
| `GITHUB_LIST_REPOSITORY_LANGUAGES` | List repository languages | read | `owner`, `repo` |  | `public_repo`, `repo` |
| `GITHUB_LIST_REPOSITORY_ORGANIZATION_SECRETS` | List repository organization secrets | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_REPOSITORY_ORGANIZATION_VARIABLES` | List repository organization variables | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_REPOSITORY_PROJECTS` | List repository projects | read | `owner`, `repo`; *page, per_page, state* |  | `repo` |
| `GITHUB_LIST_REPOSITORY_RULESETS` | Get all repository rulesets | read | `owner`, `repo`; *includes_parents, page, per_page* |  | `repo` |
| `GITHUB_LIST_REPOSITORY_RULE_SUITES` | List repository rule suites | read | `owner`, `repo`; *actor_name, page, per_page, ref, rule_suite_result, time_period* |  |  |
| `GITHUB_LIST_REPOSITORY_SECRETS` | List repository secrets | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_REPOSITORY_SECRETS_WITHOUT_DECRYPTING` | List repository secrets without decrypting | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_REPOSITORY_SECURITY_ADVISORIES` | List repository security advisories | read | `owner`, `repo`; *after, before, direction, per_page, sort, state* |  | `repo` |
| `GITHUB_LIST_REPOSITORY_TAGS` | List repository tags | read | `owner`, `repo`; *page, per_page* |  | `public_repo`, `repo` |
| `GITHUB_LIST_REPOSITORY_TEAMS` | List repository teams | read | `owner`, `repo`; *page, per_page* |  | `public_repo`, `repo` |
| `GITHUB_LIST_REPOSITORY_TOPICS` | Get all repository topics | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_REPOSITORY_VARIABLES` | List repository variables | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_REPOSITORY_WEBHOOKS` | List repository webhooks | read | `owner`, `repo`; *page, per_page* |  | `read:repo_hook`, `repo` |
| `GITHUB_LIST_REPOSITORY_WORKFLOWS` | List repository workflows | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_REPO_ACCESS_BY_TOKEN` | List repo access by token | read | `org`, `pat_request_id`; *page, per_page* |  |  |
| `GITHUB_LIST_REPO_CODESPACES` | List repo codespaces for authenticated user | read | `owner`, `repo`; *page, per_page* |  | `codespace` |
| `GITHUB_LIST_REPO_DEV_CONTAINER_CONFIGS_FOR_USER` | List repo dev container configs for user | read | `owner`, `repo`; *page, per_page* |  | `codespace` |
| `GITHUB_LIST_REPO_INVITATIONS_FOR_AUTH_USER` | List authenticated user's repository invitations | read | *page, per_page* |  | `repo`, `repo:invite` |
| `GITHUB_LIST_REPO_NOTIFICATIONS_FOR_AUTH_USER` | List repository notifications for the authenticated user | read | `owner`, `repo`; *all, before, page, participating, per_page, since* |  | `notifications`, `repo` |
| `GITHUB_LIST_REPO_SECRETS_WITHOUT_VALUES` | List repo secrets without values | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_REQUESTED_REVIEWERS_FOR_A_PULL_REQUEST` | Get all requested reviewers for a pull request | read | `owner`, `pull_number`, `repo` |  |  |
| `GITHUB_LIST_REVIEWS_FOR_A_PULL_REQUEST` | List reviews for a pull request | read | `owner`, `pull_number`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_REVIEW_COMMENTS_IN_A_REPOSITORY` | List review comments in a repository | read | `owner`, `repo`; *direction, page, per_page, since, sort* |  | `repo` |
| `GITHUB_LIST_REVIEW_COMMENTS_ON_A_PULL_REQUEST` | List review comments on a pull request | read | `owner`, `pull_number`, `repo`; *direction, page, per_page, since, sort* |  | `repo` |
| `GITHUB_LIST_RUNNER_APPLICATIONS_FOR_AN_ORGANIZATION` | List runner applications for an organization | read | `org` |  | `admin:org`, `repo` |
| `GITHUB_LIST_RUNNER_APPLICATIONS_FOR_A_REPOSITORY` | List runner applications for a repository | read | `owner`, `repo` |  | `repo` |
| `GITHUB_LIST_SECRETS_FOR_THE_AUTHENTICATED_USER` | List secrets for the authenticated user | read | *page, per_page* |  | `codespace`, `codespace:secrets` |
| `GITHUB_LIST_SECRET_SCANNING_ALERTS_FOR_A_REPOSITORY` | List secret scanning alerts for a repository | read | `owner`, `repo`; *after, assignee, before, direction, hide_secret, is_multi_repo, is_publicly_leaked, page and 6 more* |  | `public_repo`, `repo`, `security_events` |
| `GITHUB_LIST_SECURITY_ADVISORIES_GRAPHQL` | List security advisories (GraphQL) | write | *after, before, classifications, epssPercentage, epssPercentile, first, identifier, last and 3 more* |  |  |
| `GITHUB_LIST_SECURITY_VULNERABILITIES` | List security vulnerabilities | read | *after, before, classifications, ecosystem, first, last, order_by_direction, order_by_field and 2 more* |  |  |
| `GITHUB_LIST_SELECTED_REPOSITORIES_FOR_AN_ORGANIZATION_SECRET` | List selected repositories for a Codespaces organization secret | read | `org`, `secret_name`; *page, per_page* |  | `admin:org` |
| `GITHUB_LIST_SELECTED_REPOSITORIES_FOR_A_USER_SECRET` | List selected repositories for a user secret | read | `secret_name`; *page, per_page* |  | `codespace`, `codespace:secrets` |
| `GITHUB_LIST_SELECTED_REPOSITORIES_FOR_SECRET_ACCESS` | List selected repos for secret access | read | `org`, `secret_name`; *page, per_page* |  | `admin:org` |
| `GITHUB_LIST_SELECTED_REPOS_ORG_VARIABLE` | List selected repositories for an organization variable | read | `name`, `org`; *page, per_page* |  | `admin:org` |
| `GITHUB_LIST_SELF_HOSTED_RUNNERS_FOR_AN_ORGANIZATION` | List self hosted runners for an organization | read | `org`; *name, page, per_page* |  | `admin:org` |
| `GITHUB_LIST_SELF_HOSTED_RUNNERS_FOR_A_REPOSITORY` | List self hosted runners for a repository | read | `owner`, `repo`; *name, page, per_page* |  | `repo` |
| `GITHUB_LIST_SOCIAL_ACCOUNTS_FOR_A_USER` | List social accounts for a user | read | `username`; *page, per_page* |  | `read:user` |
| `GITHUB_LIST_SOCIAL_ACCOUNTS_FOR_THE_AUTHENTICATED_USER` | List social accounts for the authenticated user | read | *page, per_page* |  | `user` |
| `GITHUB_LIST_SPONSORABLES` | List Sponsorable Users and Organizations | read | *after, before, ecosystem, first, last, only_dependencies, order_by_direction, order_by_field and 1 more* |  |  |
| `GITHUB_LIST_SSH_SIGNING_KEYS_FOR_A_USER` | List SSH signing keys for user | read | `username`; *page, per_page* |  |  |
| `GITHUB_LIST_STARGAZERS` | List stargazers | read | `owner`, `repo`; *page, per_page* |  | `repo` |
| `GITHUB_LIST_STARRED_GISTS` | List starred gists | read | *page, per_page, since* |  | `gist` |
| `GITHUB_LIST_STARRED_REPOSITORIES_FOR_USER` | List starred repositories for user | read | `username`; *direction, page, per_page, sort* |  |  |
| `GITHUB_LIST_STATUS_CHECK_CONTEXTS` | Get all status check contexts | read | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_LIST_SUBSCRIPTIONS_FOR_THE_AUTHENTICATED_USER` | List marketplace subscriptions for authenticated user | read | *page, per_page* |  | `read:user`, `user:email` |
| `GITHUB_LIST_SUBSCRIPTIONS_FOR_THE_AUTHENTICATED_USER_STUBBED` | List stubbed subscriptions for the authenticated user | read | *page, per_page* |  |  |
| `GITHUB_LIST_SUB_ISSUES` | List sub-issues for an issue | read | `issue_number`, `owner`, `repo`; *page, per_page* |  | `repo`, `public_repo` |
| `GITHUB_LIST_TEAMS` | List teams | read | `org`; *page, per_page* |  | `read:org` |
| `GITHUB_LIST_TEAMS_FOR_THE_AUTHENTICATED_USER` | List teams for the authenticated user | read | *page, per_page* |  | `user`, `repo`, `read:org` |
| `GITHUB_LIST_TEAMS_THAT_ARE_ASSIGNED_TO_AN_ORGANIZATION_ROLE` | List teams assigned to an organization role | read | `org`, `role_id`; *page, per_page* |  | `admin:org`, `read:org` |
| `GITHUB_LIST_TEAM_MEMBERS` | List team members | read | `org`, `team_slug`; *page, per_page, role* |  | `read:org` |
| `GITHUB_LIST_TEAM_PROJECTS` | List team projects | read | `org`, `team_slug`; *page, per_page* |  | `read:org`, `repo` |
| `GITHUB_LIST_TEAM_REPOSITORIES` | List team repositories | read | `org`, `team_slug`; *page, per_page* |  | `read:org` |
| `GITHUB_LIST_THE_PEOPLE_A_USER_FOLLOWS` | List the people a user follows | read | `username`; *page, per_page* |  | `user:follow` |
| `GITHUB_LIST_THE_PEOPLE_THE_AUTHENTICATED_USER_FOLLOWS` | List people the authenticated user follows | read | *page, per_page* |  | `read:user`, `user:follow` |
| `GITHUB_LIST_TIMELINE_EVENTS_FOR_AN_ISSUE` | List timeline events for an issue | read | `issue_number`, `owner`, `repo`; *page, per_page* |  | `public_repo`, `repo` |
| `GITHUB_LIST_TOKEN_ACCESS_REPOSITORIES` | List token access repositories | read | `org`, `pat_id`; *page, per_page* |  |  |
| `GITHUB_LIST_USERS` | List users | read | *per_page, since* |  |  |
| `GITHUB_LIST_USERS_BLOCKED_BY_AN_ORGANIZATION` | List users blocked by an organization | read | `org`; *page, per_page* |  | `admin:org`, `read:org` |
| `GITHUB_LIST_USERS_BLOCKED_BY_THE_AUTHENTICATED_USER` | List users blocked by the authenticated user | read | *page, per_page* |  |  |
| `GITHUB_LIST_USERS_THAT_ARE_ASSIGNED_TO_AN_ORGANIZATION_ROLE` | List users assigned org role | read | `org`, `role_id`; *page, per_page* |  | `admin:org`, `read:org` |
| `GITHUB_LIST_USER_PROJECTS` | List user projects | read | `username`; *page, per_page, state* |  | `project`, `public_repo`, `read:project`, `repo` |
| `GITHUB_LIST_WATCHERS` | List watchers | read | `owner`, `repo`; *page, per_page* |  | `public_repo`, `repo` |
| `GITHUB_LIST_WORKFLOW_RUNS_FOR_A_REPOSITORY` | List workflow runs for a repository | read | `owner`, `repo`; *actor, branch, check_suite_id, created, event, exclude_pull_requests, head_sha, page and 2 more* |  | `repo` |
| `GITHUB_LIST_WORKFLOW_RUNS_FOR_A_WORKFLOW` | List workflow runs for a workflow | read | `owner`, `repo`, `workflow_id`; *actor, branch, check_suite_id, created, event, exclude_pull_requests, head_sha, page and 2 more* |  | `repo` |
| `GITHUB_LIST_WORKFLOW_RUN_ARTIFACTS` | List workflow run artifacts | read | `owner`, `repo`, `run_id`; *name, page, per_page* |  | `repo` |
| `GITHUB_LOCK_AN_ISSUE` | Lock an issue | write | `issue_number`, `owner`, `repo`; *lock_reason* |  | `repo` |
| `GITHUB_LOOKUP_ENTERPRISE_ADMINISTRATOR_INVITATION` | Lookup Enterprise Administrator Invitation | read | `enterprise_slug`, `role`, `user_login` |  | `admin:enterprise`, `read:enterprise`, `repo` |
| `GITHUB_LOOKUP_GITHUB_RESOURCE` | Lookup GitHub Resource | read | `url` |  | `repo`, `public_repo` |
| `GITHUB_LOOKUP_NODES_BY_IDS` | Lookup nodes by IDs | read | `ids` |  | `admin:org`, `public_repo`, `read:org`, `read:project`, `read:user`, `repo`, `user` |
| `GITHUB_LOOKUP_ORGANIZATION_GRAPH_QL` | Lookup Organization (GraphQL) | write | `login` |  | `read:org`, `repo` |
| `GITHUB_LOOKUP_TOPIC_BY_NAME` | Look Up Topic by Name | read | `name`; *include_related_topics, related_topics_count* |  |  |
| `GITHUB_MANAGE_ACCESS_CONTROL_FOR_ORGANIZATION_CODESPACES` | Manage access control for organization codespaces | write | `org`, `visibility`; *selected_usernames* |  | `admin:org` |
| `GITHUB_MANAGE_CUSTOM_PROPERTIES_FOR_ORG_REPOS` | Manage custom properties for org repos | write | `org`, `properties`, `repository_names` |  | `admin:org`, `read:org`, `write:org` |
| `GITHUB_MANAGE_ISSUE_BLOCKING_DEPENDENCY` | Manage issue blocking dependency | write | `blocking_issue_id`, `issue_number`, `operation`, `owner`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_MANAGE_SECRETS_IN_SELECTED_REPOSITORIES` | List selected repositories for a GitHub Actions organization secret | read | `org`, `secret_name`; *page, per_page* |  | `admin:org` |
| `GITHUB_MARK_A_THREAD_AS_DONE` | Mark a thread as done | write | `thread_id` |  | `notifications` |
| `GITHUB_MARK_A_THREAD_AS_READ` | Mark a thread as read | write | `thread_id` |  | `notifications` |
| `GITHUB_MARK_NOTIFICATIONS_AS_READ` | Mark notifications as read | write | *last_read_at, read* |  | `notifications` |
| `GITHUB_MARK_PULL_REQUEST_READY_FOR_REVIEW` | Mark pull request ready for review | write | `pull_request_id`; *client_mutation_id* |  | `repo` |
| `GITHUB_MARK_REPOSITORY_NOTIFICATIONS_AS_READ` | Mark repository notifications as read | write | `owner`, `repo`; *last_read_at* |  | `repo` |
| `GITHUB_MERGE_A_BRANCH` | Merge a branch | write | `base`, `head`, `owner`, `repo`; *commit_message* |  | `repo` |
| `GITHUB_MERGE_A_PULL_REQUEST` | Merge a pull request | write | `owner`, `pull_number`, `repo`; *commit_message, commit_title, merge_method, sha* |  | `repo` |
| `GITHUB_PING_AN_ORGANIZATION_WEBHOOK` | Ping an organization webhook | write | `hook_id`, `org` |  | `admin:org_hook` |
| `GITHUB_PING_A_REPOSITORY_WEBHOOK` | Ping a repository webhook | write | `hook_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_PRIVATELY_REPORT_A_SECURITY_VULNERABILITY` | Privately report a security vulnerability | write | `description`, `owner`, `repo`, `summary`; *cvss_vector_string, cwe_ids, severity, start_private_fork, vulnerabilities* |  | `repo` |
| `GITHUB_PUBLISH_SPONSORS_TIER_GRAPH_QL` | Publish Sponsors Tier (GraphQL) | write | `tierId`; *clientMutationId* |  | `user`, `admin:org` |
| `GITHUB_REDELIVER_A_DELIVERY_FOR_AN_ORGANIZATION_WEBHOOK` | Redeliver a delivery for an organization webhook | write | `delivery_id`, `hook_id`, `org` |  | `admin:org_hook` |
| `GITHUB_REDELIVER_A_DELIVERY_FOR_A_REPOSITORY_WEBHOOK` | Redeliver a delivery for a repository webhook | write | `delivery_id`, `hook_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_REMOVE_ALL_LABELS_FROM_AN_ISSUE` | Remove all labels from an issue | write | `issue_number`, `owner`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_REMOVE_ALL_ORGANIZATION_ROLES_FOR_A_TEAM` | Remove all organization roles for a team | write | `org`, `team_slug` |  | `admin:org` |
| `GITHUB_REMOVE_ALL_ORGANIZATION_ROLES_FOR_A_USER` | Remove all organization roles for a user | write | `org`, `username` |  | `admin:org` |
| `GITHUB_REMOVE_AN_ORGANIZATION_MEMBER` | Remove an organization member | write | `org`, `username` |  | `admin:org` |
| `GITHUB_REMOVE_AN_ORGANIZATION_ROLE_FROM_A_TEAM` | Remove an organization role from a team | write | `org`, `role_id`, `team_slug` |  | `admin:org` |
| `GITHUB_REMOVE_AN_ORGANIZATION_ROLE_FROM_A_USER` | Remove an organization role from a user | write | `org`, `role_id`, `username` |  | `admin:org` |
| `GITHUB_REMOVE_APP_ACCESS_RESTRICTIONS` | Remove app access restrictions | write | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_REMOVE_ASSIGNEES_FROM_AN_ISSUE` | Remove assignees from an issue | write | `issue_number`, `owner`, `repo`; *assignees* |  | `repo` |
| `GITHUB_REMOVE_A_CUSTOM_PROPERTY_FOR_AN_ORGANIZATION` | Remove a custom property for an organization | write | `custom_property_name`, `org` |  | `admin:org` |
| `GITHUB_REMOVE_A_LABEL_FROM_AN_ISSUE` | Remove a label from an issue | write | `issue_number`, `name`, `owner`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_REMOVE_A_PROJECT_FROM_A_TEAM` | Remove a project from a team | write | `org`, `project_id`, `team_slug` |  | `repo` |
| `GITHUB_REMOVE_A_REPOSITORY_COLLABORATOR` | Remove a repository collaborator | write | `owner`, `repo`, `username` |  | `repo` |
| `GITHUB_REMOVE_A_REPOSITORY_FROM_AN_APP_INSTALLATION` | Remove a repository from an app installation | write | `installation_id`, `repository_id` |  | `repo`, `repo:status`, `repo_deployment`, `public_repo`, `delete_repo`, `notifications`, `gist` |
| `GITHUB_REMOVE_A_REPOSITORY_FROM_A_TEAM` | Remove a repository from a team | write | `org`, `owner`, `repo`, `team_slug` |  | `write:org` |
| `GITHUB_REMOVE_A_SELECTED_REPOSITORY_FROM_A_USER_SECRET` | Remove a selected repository from a user secret | write | `repository_id`, `secret_name` |  | `codespace`, `codespace:secrets` |
| `GITHUB_REMOVE_CUSTOM_LABELS_FROM_SELF_HOSTED_REPO_RUNNER` | Remove custom labels from self hosted repo runner | write | `owner`, `repo`, `runner_id` |  | `repo` |
| `GITHUB_REMOVE_CUSTOM_LABEL_FROM_REPO_RUNNER` | Remove custom label from repo runner | write | `name`, `owner`, `repo`, `runner_id` |  | `repo` |
| `GITHUB_REMOVE_CUSTOM_LABEL_FROM_SELF_HOSTED_RUNNER` | Remove custom label from self hosted runner | write | `name`, `org`, `runner_id` |  | `admin:org` |
| `GITHUB_REMOVE_INTERACTION_RESTRICTIONS_FOR_AN_ORGANIZATION` | Remove interaction restrictions for an organization | write | `org` |  | `admin:org`, `write:org` |
| `GITHUB_REMOVE_INTERACTION_RESTRICTIONS_FOR_A_REPOSITORY` | Remove interaction restrictions for a repository | write | `owner`, `repo` |  |  |
| `GITHUB_REMOVE_ORG_DEV_ENV_SECRET_BY_NAME` | Remove org dev env secret by name | write | `org`, `secret_name` |  | `admin:org` |
| `GITHUB_REMOVE_ORG_SECRET_BY_NAME` | Remove org secret by name | write | `org`, `secret_name` |  | `admin:org` |
| `GITHUB_REMOVE_OUTSIDE_COLLABORATOR_FROM_AN_ORGANIZATION` | Remove outside collaborator from an organization | write | `org`, `username` |  | `admin:org` |
| `GITHUB_REMOVE_PUBLIC_ORG_MEMBERSHIP` | Remove public org membership | write | `org`, `username` |  | `admin:org`, `user` |
| `GITHUB_REMOVE_PUBLIC_REPO_INTERACTION_RESTRICTIONS` | Remove user public repo interaction restrictions | write |  |  | `user` |
| `GITHUB_REMOVE_REPO_FROM_ORG_DEV_ENV_SECRET` | Remove repository from organization dev environment secret | write | `org`, `repository_id`, `secret_name` |  | `admin:org` |
| `GITHUB_REMOVE_REPO_FROM_ORG_SECRET_WITH_SELECTED_VISIBILITY` | Remove selected repo from org secret | write | `org`, `repository_id`, `secret_name` |  | `admin:org` |
| `GITHUB_REMOVE_REQUESTED_REVIEWERS_FROM_A_PULL_REQUEST` | Remove requested reviewers from a pull request | write | `owner`, `pull_number`, `repo`, `reviewers`; *team_reviewers* |  |  |
| `GITHUB_REMOVE_REVIEW_REQUEST` | Remove review request | write | `owner`, `pull_number`, `repo`, `reviewers`; *team_reviewers* |  | `repo`, `public_repo` |
| `GITHUB_REMOVE_SELECTED_REPO_FROM_ORG_SECRET` | Remove selected repository from an organization secret | write | `org`, `repository_id`, `secret_name` |  | `admin:org`, `repo` |
| `GITHUB_REMOVE_SELECTED_REPO_FROM_ORG_VARIABLE` | Remove selected repository from an organization variable | write | `name`, `org`, `repository_id` |  | `admin:org`, `repo` |
| `GITHUB_REMOVE_STATUS_CHECK_CONTEXTS` | Remove status check contexts | write | `branch`, `contexts`, `owner`, `repo` |  | `repo` |
| `GITHUB_REMOVE_STATUS_CHECK_PROTECTION` | Remove status check protection | write | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_REMOVE_TEAM_ACCESS_RESTRICTIONS` | Remove team access restrictions | write | `branch`, `owner`, `repo`, `teams` |  | `repo` |
| `GITHUB_REMOVE_TEAM_MEMBERSHIP` | Remove team membership | write | `org`, `team_slug`, `username` |  | `admin:org`, `write:org` |
| `GITHUB_REMOVE_TEAM_MEMBERSHIP_FOR_A_USER` | Remove team membership for a user | write | `org`, `team_slug`, `username` |  | `admin:org` |
| `GITHUB_REMOVE_USERS_CODESPACES_ACCESS_ORG` | Remove users from codespaces access for an organization | write | `org`, `selected_usernames` |  | `admin:org` |
| `GITHUB_REMOVE_USER_ACCESS_RESTRICTIONS` | Remove user access restrictions | write | `branch`, `owner`, `repo`, `users` |  | `repo` |
| `GITHUB_RENAME_A_BRANCH` | Rename a branch | write | `branch`, `new_name`, `owner`, `repo` |  | `repo` |
| `GITHUB_RENDER_A_MARKDOWN_DOCUMENT` | Render a markdown document | read | `text`; *context, mode* |  |  |
| `GITHUB_RENDER_MARKDOWN_RAW` | Render Markdown (Raw Mode) | write | `markdown_content` |  |  |
| `GITHUB_REOPEN_ISSUE` | Reopen an issue | write | `issue_number`, `owner`, `repo` |  | `repo`, `public_repo` |
| `GITHUB_REOPEN_PULL_REQUEST` | Reopen Pull Request | write | `owner`, `pull_number`, `repo` |  | `public_repo`, `repo` |
| `GITHUB_REPLACE_ALL_REPOSITORY_TOPICS` | Replace all repository topics | write | `names`, `owner`, `repo` |  | `repo` |
| `GITHUB_REPLACE_ORG_SECRET_VISIBILITY_TO_SELECTED` | Replace org secret visibility to selected | write | `org`, `secret_name`, `selected_repository_ids` |  | `admin:org` |
| `GITHUB_REPLACE_REPO_ACCESS_ON_ORG_DEV_ENV_SECRET_SET` | Replace repository access for an org Codespaces secret | write | `org`, `secret_name`, `selected_repository_ids` |  | `admin:org` |
| `GITHUB_REQUEST_A_GITHUB_PAGES_BUILD` | Request a github pages build | write | `owner`, `repo` |  | `repo` |
| `GITHUB_REQUEST_REVIEWERS` | Request reviewers for a pull request | write | `owner`, `pull_number`, `repo`; *reviewers, team_reviewers* |  | `repo` |
| `GITHUB_REQUEST_REVIEWERS_FOR_A_PULL_REQUEST` | Add reviewers to a specific pull request | write | `owner`, `pull_number`, `repo`; *reviewers, team_reviewers* |  | `repo` |
| `GITHUB_REREQUEST_A_CHECK_RUN` | Rerequest a check run | write | `check_run_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_REREQUEST_A_CHECK_SUITE` | Rerequest a check suite | write | `check_suite_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_RESET_A_TOKEN` | Reset a token | write | `access_token`, `client_id` |  |  |
| `GITHUB_RESTORE_A_PACKAGE_FOR_AN_ORGANIZATION` | Restore a package for an organization | write | `org`, `package_name`, `package_type`; *token* |  | `write:packages`, `read:packages` |
| `GITHUB_RESTORE_A_PACKAGE_FOR_A_USER` | Restore a package for a user | write | `package_name`, `package_type`, `username`; *token* |  | `write:packages`, `read:packages` |
| `GITHUB_RESTORE_A_PACKAGE_FOR_THE_AUTHENTICATED_USER` | Restore a package for the authenticated user | write | `package_name`, `package_type`; *token* |  | `write:packages`, `read:packages` |
| `GITHUB_RESTORE_A_PACKAGE_VERSION_FOR_THE_AUTHENTICATED_USER` | Restore a package version for the authenticated user | write | `package_name`, `package_type`, `package_version_id` |  | `read:packages`, `write:packages` |
| `GITHUB_RESTORE_PACKAGE_VERSION_FOR_AN_ORGANIZATION` | Restore package version for an organization | write | `org`, `package_name`, `package_type`, `package_version_id` |  | `read:packages`, `write:packages` |
| `GITHUB_RESTORE_PACKAGE_VERSION_FOR_A_USER` | Restore package version for a user | write | `package_name`, `package_type`, `package_version_id`, `username` |  | `read:packages`, `write:packages` |
| `GITHUB_RETIRE_SPONSORS_TIER` | Retire Sponsors Tier | write | `tier_id`; *client_mutation_id* |  | `user` |
| `GITHUB_REVIEW_ACCESS_WITH_PERSONAL_TOKEN` | Review access with personal token | write | `action`, `org`, `pat_request_id`; *reason* |  | `admin:org` |
| `GITHUB_REVIEW_DEPLOYMENT_PROTECTION_RULES` | Review Custom Deployment Rules For Workflow Run | write | `environment_name`, `owner`, `repo`, `run_id`, `state`; *comment* |  | `repo` |
| `GITHUB_REVIEW_PENDING_DEPLOYMENTS_FOR_A_WORKFLOW_RUN` | Review pending deployments for a workflow run | write | `comment`, `environment_ids`, `owner`, `repo`, `run_id`, `state` |  | `repo` |
| `GITHUB_REVIEW_RESOURCE_REQUESTS_WITH_FINE_GRAINED_TOKENS` | Review PAT Requests in Bulk | write | `action`, `org`, `pat_request_ids`; *reason* |  | `admin:org` |
| `GITHUB_REVOKE_AN_INSTALLATION_ACCESS_TOKEN` | Revoke an installation access token | write |  |  |  |
| `GITHUB_RE_RUN_A_JOB_FROM_A_WORKFLOW_RUN` | Re-run a job from a workflow run | write | `job_id`, `owner`, `repo`; *enable_debug_logging* |  | `repo` |
| `GITHUB_RE_RUN_A_WORKFLOW` | Rerun a workflow | write | `owner`, `repo`, `run_id`; *enable_debug_logging* |  | `repo` |
| `GITHUB_RE_RUN_FAILED_JOBS_FROM_A_WORKFLOW_RUN` | Rerun failed jobs from a workflow run | write | `owner`, `repo`, `run_id`; *enable_debug_logging* |  | `repo` |
| `GITHUB_RUN_GRAPH_QL_QUERY` | Run GraphQL Query | write | `query`; *variables* |  |  |
| `GITHUB_RUN_GRAPH_QL_RELAY_QUERY` | Run GraphQL Relay Query | write | `query`; *variables* |  | `repo`, `user`, `read:user`, `read:org`, `public_repo`, `admin:org`, `write:org`, `gist`, `notifications`, `project`, `read:project`, `delete_repo`, `write:packages`, `read:packages`, `delete:packages`, `admin:repo_hook`, `write:repo_hook`, `read:repo_hook`, `admin:org_hook`, `admin:public_key`, `write:public_key`, `read:public_key`, `admin:gpg_key`, `write:gpg_key`, `read:gpg_key`, `codespace`, `workflow`, `user:email`, `user:follow`, `repo:status`, `repo_deployment`, `repo:invite`, `security_events`, `read:audit_log` |
| `GITHUB_SEARCH_CODE` | Search code | read | `q`; *order, page, per_page, sort* |  | `public_repo` |
| `GITHUB_SEARCH_CODE_ALL_PAGES` | Search code (all pages) | read | `q`; *accept, limit, order, page, per_page, sort* |  |  |
| `GITHUB_SEARCH_COMMITS` | Search commits | read | `q`; *order, page, per_page, sort* |  | `repo` |
| `GITHUB_SEARCH_COMMITS_BY_AUTHOR` | Search commits by author across repos | read | `q`; *composio_execution_message, order, page, per_page, sort* |  |  |
| `GITHUB_SEARCH_GITHUB_GRAPHQL` | GitHub GraphQL Search | write | `query`, `type`; *after, first* |  | `repo`, `public_repo` |
| `GITHUB_SEARCH_ISSUES` | Search issues and pull requests | read | `q`; *order, page, per_page, sort* |  |  |
| `GITHUB_SEARCH_ISSUES_AND_PULL_REQUESTS` | Search issues and pull requests | read | `q`; *order, page, per_page, response_detail, sort* |  | `repo`, `public_repo` |
| `GITHUB_SEARCH_LABELS` | Search labels | read | `q`, `repository_id`; *order, page, per_page, sort* |  |  |
| `GITHUB_SEARCH_REPOSITORIES` | Search repositories | read | `q`; *order, page, per_page, sort* |  | `public_repo` |
| `GITHUB_SEARCH_TOPICS` | Search topics | read | `q`; *page, per_page* |  |  |
| `GITHUB_SEARCH_USERS` | Search users | read | `q`; *order, page, per_page, sort* |  |  |
| `GITHUB_SET_ADMIN_BRANCH_PROTECTION` | Set admin branch protection | write | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_SET_APP_ACCESS_RESTRICTIONS` | Set app access restrictions | write | `apps`, `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_SET_A_REPOSITORY_SUBSCRIPTION` | Set a repository subscription | write | `owner`, `repo`; *ignored, subscribed* |  | `repo`, `public_repo` |
| `GITHUB_SET_A_THREAD_SUBSCRIPTION` | Set a thread subscription | write | `thread_id`; *ignored* |  | `notifications` |
| `GITHUB_SET_DEFAULT_WORKFLOW_PERMISSIONS_FOR_AN_ORGANIZATION` | Set default workflow permissions for an organization | write | `org`; *can_approve_pull_request_reviews, default_workflow_permissions* |  | `admin:org` |
| `GITHUB_SET_DEFAULT_WORKFLOW_PERMISSIONS_FOR_A_REPOSITORY` | Set default workflow permissions for a repository | write | `owner`, `repo`; *can_approve_pull_request_reviews, default_workflow_permissions* |  | `repo` |
| `GITHUB_SET_EMAIL_VISIBILITY` | Set primary email visibility | write | `visibility` |  | `user:email` |
| `GITHUB_SET_GITHUB_ACTIONS_PERMISSIONS_FOR_AN_ORGANIZATION` | Set GitHub Actions permissions for an organization | write | `enabled_repositories`, `org`; *allowed_actions, sha_pinning_required* |  | `admin:org` |
| `GITHUB_SET_GITHUB_ACTIONS_PERMISSIONS_FOR_A_REPOSITORY` | Set github actions permissions for a repository | write | `enabled`, `owner`, `repo`; *allowed_actions* |  | `repo` |
| `GITHUB_SET_INTERACTION_RESTRICTIONS_FOR_AN_ORGANIZATION` | Set interaction restrictions for an organization | write | `limit`, `org`; *expiry* |  | `admin:org` |
| `GITHUB_SET_INTERACTION_RESTRICTIONS_FOR_A_REPOSITORY` | Set interaction restrictions for a repository | write | `limit`, `owner`, `repo`; *expiry* |  |  |
| `GITHUB_SET_INTERACTION_RESTRICTIONS_FOR_YOUR_PUBLIC_REPOS` | Set interaction restrictions for your public repositories | write | `limit`; *expiry* |  | `read:user` |
| `GITHUB_SET_LABELS_FOR_AN_ISSUE` | Set labels for an issue | write | `issue_number`, `labels`, `owner`, `repo` |  | `public_repo`, `repo` |
| `GITHUB_SET_OIDC_SUBJECT_CLAIM_TEMPLATE` | Set custom OIDC subject claim template | write | `owner`, `repo`, `use_default`; *include_claim_keys* |  | `repo` |
| `GITHUB_SET_ORG_ALLOWED_ACTIONS` | Set allowed actions and workflows for an organization | write | `org`; *github_owned_allowed, patterns_allowed, verified_allowed* |  | `admin:org` |
| `GITHUB_SET_ORG_RUNNER_LABELS` | Set custom labels for a self-hosted runner for an organization | write | `labels`, `org`, `runner_id` |  | `admin:org` |
| `GITHUB_SET_REPO_ALLOWED_ACTIONS` | Set repo allowed actions | write | `owner`, `repo`; *github_owned_allowed, patterns_allowed, verified_allowed* |  | `repo` |
| `GITHUB_SET_REPO_RUNNER_LABELS` | Set custom labels for repo runner | write | `labels`, `owner`, `repo`, `runner_id` |  | `repo` |
| `GITHUB_SET_SELECTED_REPOSITORIES_FOR_AN_ORGANIZATION_SECRET` | Set selected repositories for an organization secret | write | `org`, `secret_name`, `selected_repository_ids` |  | `admin:org`, `repo` |
| `GITHUB_SET_SELECTED_REPOSITORIES_FOR_A_USER_SECRET` | Set selected repositories for a user secret | write | `secret_name`, `selected_repository_ids` |  | `codespace`, `codespace:secrets` |
| `GITHUB_SET_SELECTED_REPOS_FOR_ORG_VARIABLE` | Set selected repositories for an organization variable | write | `name`, `org`, `selected_repository_ids` |  | `admin:org` |
| `GITHUB_SET_STATUS_CHECK_CONTEXTS` | Set status check contexts | write | `branch`, `contexts`, `owner`, `repo` |  | `repo` |
| `GITHUB_SET_TEAM_ACCESS_RESTRICTIONS` | Set team access restrictions | write | `branch`, `owner`, `repo`, `teams` |  | `repo` |
| `GITHUB_SET_USER_ACCESS_RESTRICTIONS` | Set user access restrictions | write | `branch`, `owner`, `repo`, `users` |  | `repo` |
| `GITHUB_SET_USER_INTERACTION_LIMIT_GRAPH_QL` | Set User Interaction Limit (GraphQL) | write | `limit`, `userId`; *clientMutationId, expiry* |  | `user` |
| `GITHUB_SET_WORKFLOW_EXTERNAL_ACCESS` | Set workflow access level outside repository | write | `access_level`, `owner`, `repo` |  | `repo` |
| `GITHUB_START_A_CODESPACE_FOR_THE_AUTHENTICATED_USER` | Start a codespace for the authenticated user | write | `codespace_name` |  | `codespace` |
| `GITHUB_STAR_A_GIST` | Star a gist | write | `gist_id` |  | `gist` |
| `GITHUB_STAR_A_REPOSITORY_FOR_THE_AUTHENTICATED_USER` | Star a specified repository (requires owner and repo) | write | `owner`, `repo` |  | `repo` |
| `GITHUB_STOP_A_CODESPACE_FOR_AN_ORGANIZATION_USER` | Stop a codespace for an organization user | write | `codespace_name`, `org`, `username` |  | `admin:org` |
| `GITHUB_STOP_A_CODESPACE_FOR_THE_AUTHENTICATED_USER` | Stop a codespace for the authenticated user | write | `codespace_name` |  | `codespace` |
| `GITHUB_SUBMIT_A_REVIEW_FOR_A_PULL_REQUEST` | Submit a review for a pull request | write | `event`, `owner`, `pull_number`, `repo`, `review_id`; *body* |  | `repo`, `repo:status`, `public_repo` |
| `GITHUB_SYNC_FORK_BRANCH_WITH_UPSTREAM` | Sync a fork branch with the upstream repository | write | `branch`, `owner`, `repo` |  | `repo` |
| `GITHUB_TEST_REPOSITORY_WEBHOOK` | Test the push repository webhook | write | `hook_id`, `owner`, `repo` |  |  |
| `GITHUB_TRANSFER_A_REPOSITORY` | Transfer a repository | write | `new_owner`, `owner`, `repo`; *new_name, team_ids* |  | `repo` |
| `GITHUB_UNBLOCK_A_USER` | Unblock a user | write | `username` |  |  |
| `GITHUB_UNBLOCK_USER_FROM_ORGANIZATION` | Unblock a user from an organization | write | `org`, `username` |  | `admin:org` |
| `GITHUB_UNFOLLOW_A_USER` | Unfollow a user | write | `username` |  | `user`, `user:follow` |
| `GITHUB_UNFOLLOW_ORGANIZATION` | Unfollow organization | write | `organizationId`; *clientMutationId* |  | `user` |
| `GITHUB_UNFOLLOW_USER` | Unfollow user | write | `user_id`; *client_mutation_id* |  | `user:follow`, `user` |
| `GITHUB_UNLOCK_AN_ISSUE` | Unlock an issue | write | `issue_number`, `owner`, `repo` |  | `public_repo`, `repo` |
| `GITHUB_UNLOCK_ORGANIZATION_REPOSITORY` | Unlock an organization repository | write | `migration_id`, `org`, `repo_name` |  |  |
| `GITHUB_UNLOCK_USER_REPOSITORY` | Unlock a user repository | write | `migration_id`, `repo_name` |  | `repo` |
| `GITHUB_UNSTAR_A_GIST` | Unstar a gist | write | `gist_id` |  |  |
| `GITHUB_UNSTAR_REPOSITORY` | Unstar a repository | write | `owner`, `repo` |  | `public_repo`, `repo` |
| `GITHUB_UPDATE_AN_ENVIRONMENT_VARIABLE` | Update an environment variable | write | `environment_name`, `name`, `owner`, `repo`; *new_name, value* |  | `repo` |
| `GITHUB_UPDATE_AN_ISSUE` | Update an issue | write | `issue_number`, `owner`, `repo`; *assignee, assignees, body, labels, milestone, state, state_reason, title* |  | `repo`, `public_repo` |
| `GITHUB_UPDATE_AN_ISSUE_COMMENT` | Update an issue comment | write | `body`, `comment_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_UPDATE_AN_ORGANIZATION` | Update an organization | write | `org`; *advanced_security_enabled_for_new_repositories, billing_email, blog, company, default_repository_permission, dependabot_alerts_enabled_for_new_repositories, dependabot_security_updates_enabled_for_new_repositories, dependency_graph_enabled_for_new_repositories and 21 more* |  | `admin:org`, `repo` |
| `GITHUB_UPDATE_AN_ORGANIZATION_VARIABLE` | Update an organization variable | write | `name`, `org`; *new_name, selected_repository_ids, value, visibility* |  | `admin:org` |
| `GITHUB_UPDATE_AN_ORGANIZATION_WEBHOOK` | Update an organization webhook | write | `hook_id`, `org`; *active, config__content__type, config__insecure__ssl, config__secret, config__url, events, name* |  | `admin:org_hook` |
| `GITHUB_UPDATE_A_CHECK_RUN` | Update a check run | write | `check_run_id`, `owner`, `repo`; *actions, completed_at, conclusion, details_url, external_id, name, output__annotations, output__images and 5 more* |  | `repo` |
| `GITHUB_UPDATE_A_CODE_SCANNING_ALERT` | Update a code scanning alert | write | `alert_number`, `owner`, `repo`, `state`; *dismissed_comment, dismissed_reason* |  | `security_events`, `public_repo` |
| `GITHUB_UPDATE_A_COMMIT_COMMENT` | Update a commit comment | write | `body`, `comment_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_UPDATE_A_CUSTOM_ORGANIZATION_ROLE` | Update a custom organization role | write | `org`, `role_id`; *base_role, description, name, permissions* |  | `admin:org` |
| `GITHUB_UPDATE_A_DEPLOYMENT_BRANCH_POLICY` | Update a deployment branch policy | write | `branch_policy_id`, `environment_name`, `name`, `owner`, `repo` |  | `repo` |
| `GITHUB_UPDATE_A_DISCUSSION` | Update a discussion | write | `discussion_number`, `org`, `team_slug`; *body, title* |  | `write:discussion` |
| `GITHUB_UPDATE_A_DISCUSSION_COMMENT` | Update a discussion comment | write | `body`, `comment_number`, `discussion_number`, `org`, `team_slug` |  | `write:discussion` |
| `GITHUB_UPDATE_A_GIST` | Update a gist | write | `gist_id`; *description, files* |  | `gist` |
| `GITHUB_UPDATE_A_GIST_COMMENT` | Update a gist comment | write | `body`, `comment_id`, `gist_id` |  | `gist` |
| `GITHUB_UPDATE_A_LABEL` | Update a label | write | `name`, `owner`, `repo`; *color, description, new_name* |  | `repo`, `public_repo` |
| `GITHUB_UPDATE_A_MILESTONE` | Update a milestone | write | `milestone_number`, `owner`, `repo`; *description, due_on, state, title* |  | `repo` |
| `GITHUB_UPDATE_A_PROJECT` | Update a project | write | `project_number`; *closed, org, public, readme, short_description, title, username* |  | `repo`, `project` |
| `GITHUB_UPDATE_A_PULL_REQUEST` | Update a pull request | write | `owner`, `pull_number`, `repo`; *base, body, maintainer_can_modify, state, title* |  | `repo`, `public_repo` |
| `GITHUB_UPDATE_A_PULL_REQUEST_BRANCH` | Update a pull request branch | write | `owner`, `pull_number`, `repo`; *expected_head_sha* |  | `repo` |
| `GITHUB_UPDATE_A_REFERENCE` | Update a reference | write | `owner`, `ref`, `repo`, `sha`; *force* |  | `repo` |
| `GITHUB_UPDATE_A_RELEASE` | Update a release | write | `owner`, `release_id`, `repo`; *body, discussion_category_name, draft, make_latest, name, prerelease, tag_name, target_commitish* |  | `repo` |
| `GITHUB_UPDATE_A_RELEASE_ASSET` | Update a release asset | write | `asset_id`, `owner`, `repo`; *label, name, state* |  | `repo` |
| `GITHUB_UPDATE_A_REPOSITORY` | Update a repository | write | `owner`, `repo`; *allow_auto_merge, allow_forking, allow_merge_commit, allow_rebase_merge, allow_squash_merge, allow_update_branch, archived, default_branch and 19 more* |  | `repo` |
| `GITHUB_UPDATE_A_REPOSITORY_INVITATION` | Update a repository invitation | write | `invitation_id`, `owner`, `repo`; *permissions* |  | `repo:invite` |
| `GITHUB_UPDATE_A_REPOSITORY_RULESET` | Update a repository ruleset | write | `owner`, `repo`, `ruleset_id`; *bypass_actors, conditions__ref__name__exclude, conditions__ref__name__include, enforcement, name, rules, target* |  | `admin:repo_hook`, `repo` |
| `GITHUB_UPDATE_A_REPOSITORY_VARIABLE` | Update a repository variable | write | `name`, `owner`, `repo`; *new_name, value* |  | `repo` |
| `GITHUB_UPDATE_A_REPOSITORY_WEBHOOK` | Update a repository webhook | write | `hook_id`, `owner`, `repo`; *active, add_events, config__content__type, config__insecure__ssl, config__secret, config__url, events, remove_events* |  | `admin:repo_hook` |
| `GITHUB_UPDATE_A_TEAM` | Update a team | write | `org`, `team_slug`; *description, name, notification_setting, parent_team_id, permission, privacy* |  | `admin:org`, `read:org` |
| `GITHUB_UPDATE_BRANCH_PROTECTION` | Update branch protection | write | `branch`, `owner`, `repo`; *allow_deletions, allow_force_pushes, allow_fork_syncing, block_creations, enforce_admins, lock_branch, required__pull__request__reviews__dismiss__stale__reviews, required__pull_request__reviews__approving__review__count and 16 more* |  | `repo` |
| `GITHUB_UPDATE_CODESPACE` | Update a codespace for the authenticated user | write | `codespace_name`; *display_name, machine, recent_folders* |  | `codespace` |
| `GITHUB_UPDATE_CODE_SCANNING_DEFAULT_SETUP` | Update a code scanning default setup configuration | write | `owner`, `repo`; *languages, query_suite, runner_label, runner_type, state* |  | `repo`, `public_repo` |
| `GITHUB_UPDATE_ISSUE_TYPE` | Update issue type | write | `is_enabled`, `issue_type_id`, `name`, `org`; *color, description* |  | `admin:org` |
| `GITHUB_UPDATE_ORGANIZATION_WEBHOOK_CONFIGURATION` | Update a webhook configuration for an organization | write | `hook_id`, `org`; *content_type, insecure_ssl, secret, url* |  | `admin:org_hook` |
| `GITHUB_UPDATE_PULL_REQUEST_REVIEW` | Update a review for a pull request | write | `body`, `owner`, `pull_number`, `repo`, `review_id` |  |  |
| `GITHUB_UPDATE_PULL_REQUEST_REVIEW_COMMENT` | Update a review comment for a pull request | write | `body`, `comment_id`, `owner`, `repo` |  | `repo` |
| `GITHUB_UPDATE_PULL_REQUEST_REVIEW_PROTECTION` | Update pull request review protection | write | `branch`, `owner`, `repo`; *bypass__pull__request__allowances__apps, bypass__pull__request__allowances__teams, bypass__pull__request__allowances__users, dismiss_stale_reviews, dismissal__restrictions__apps, dismissal__restrictions__teams, dismissal__restrictions__users, require_code_owner_reviews and 2 more* |  | `repo`, `admin:repo_hook` |
| `GITHUB_UPDATE_REPOSITORY_PREFERENCES_FOR_CHECK_SUITES` | Update repository preferences for check suites | write | `owner`, `repo`; *auto_trigger_checks* |  | `repo` |
| `GITHUB_UPDATE_REPOSITORY_WEBHOOK_CONFIGURATION` | Update a webhook configuration for a repository | write | `hook_id`, `owner`, `repo`; *content_type, insecure_ssl, secret, url* |  | `write:repo_hook`, `repo` |
| `GITHUB_UPDATE_RESOURCE_ACCESS_WITH_TOKENS` | Update resource access with tokens | write | `action`, `org`, `pat_ids` |  |  |
| `GITHUB_UPDATE_STATUS_CHECK_PROTECTION` | Update status check protection | write | `branch`, `owner`, `repo`; *checks, contexts, strict* |  | `repo` |
| `GITHUB_UPDATE_SUBSCRIPTION` | Update Subscription | write | `state`, `subscribableId`; *clientMutationId* |  | `notifications` |
| `GITHUB_UPDATE_THE_AUTHENTICATED_USER` | Update the authenticated user | write | *bio, blog, company, email, hireable, location, name, twitter_username* |  | `user` |
| `GITHUB_UPDATE_TOKEN_ORG_ACCESS` | Revoke Token Organization Access | write | `org`, `pat_id`; *action* |  | `admin:org` |
| `GITHUB_UPDATE_USER_LIST` | Update User List | write | `list_id`; *client_mutation_id, description, is_private, name* |  |  |
| `GITHUB_UPDATE_USER_LISTS_FOR_ITEM` | Update user lists for item | write | `item_id`, `list_ids`; *client_mutation_id, suggested_list_ids* |  | `user` |
| `GITHUB_UPDATE_USER_PROJECT_ITEM` | Update Project V2 Item for User | write | `fieldId`, `itemId`, `projectId`, `value` |  | `project` |
| `GITHUB_UPDATE_USER_STATUS` | Change User Status | write | *clientMutationId, emoji, expiresAt, limitedAvailability, message, organizationId* |  | `user` |
| `GITHUB_UPLOAD_RELEASE_ASSET` | Upload a release asset | write | `asset`, `content_type`, `owner`, `release_id`, `repo`; *label* |  |  |
| `GITHUB_UPLOAD_SARIF_ANALYSIS` | Upload an analysis as sarif data | write | `commit_sha`, `owner`, `ref`, `repo`, `sarif`; *checkout_uri, started_at, tool_name, validate* |  | `public_repo`, `security_events` |
| `GITHUB_UPSERT_REPO_CUSTOM_PROPERTY_VALUES` | Create or update repository custom property values | write | `owner`, `properties`, `repo` |  | `repo` |
| `GITHUB_VALIDATE_CREDENTIAL` | Validate credential | read |  |  | `user`, `read:user` |
| `GITHUB_VERIFY_DEV_CONTAINER_PERMISSIONS_ACCEPTED` | Verify dev container permissions accepted | read | `devcontainer_path`, `owner`, `ref`, `repo` |  | `codespace` |
| `GITHUB_WHO_AM_I` | Who Am I | read |  |  |  |
