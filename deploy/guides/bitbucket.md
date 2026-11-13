# Bitbucket

Checked: 2026-11-13 · Sources: the Bitbucket tool definitions as planque
reads them (116 tools, generated/bitbucket.md), Atlassian's Bitbucket Cloud
REST API reference (authentication, scopes, pagination, error layout) and
its API request limits page.

## How it connects

- **A sign-in to Bitbucket** with the owner's Atlassian account, held by the
  connection service; the taf never sees a password or token. The catalog
  lists the tools and, with a known sign-in, puts the ones its access doesn't
  cover under `unusable`, with why.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request, so nothing is asked
  twice.
- **A sign-in left unused for three months lapses.** Bitbucket's sign-ins
  renew themselves while used; after three months without use the owner
  signs in again.

## Tools a job usually needs

Read means it changes nothing in Bitbucket. The generated part gives no
returns; the ids below are the ones Bitbucket's API reference shows.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `BITBUCKET_LIST_WORKSPACES` | Lists the workspaces the owner belongs to | read | workspace slugs |
| `BITBUCKET_LIST_REPOSITORIES_IN_WORKSPACE` | Lists a workspace's repositories | read | repository slugs, links |
| `BITBUCKET_GET_FILE_FROM_REPOSITORY` / `BITBUCKET_GET_REPOSITORIES_SRC` | Reads one file or folder at a branch or commit | read | the file's content |
| `BITBUCKET_LIST_BRANCHES` / `BITBUCKET_LIST_COMMITS` | Reads branches and recent commits | read | branch names, commit hashes |
| `BITBUCKET_LIST_PULL_REQUESTS` / `BITBUCKET_GET_PULL_REQUEST` | Reads pull requests and their state | read | pull request ids, links, source and destination branches |
| `BITBUCKET_GET_PULL_REQUEST_DIFF` | Reads a pull request's changes | read | the diff |
| `BITBUCKET_GET_REPOSITORIES_PULLREQUESTS_COMMENTS` | Reads a pull request's comments | read | comment ids |
| `BITBUCKET_LIST_ISSUES` | Reads a repository's issues | read | issue ids, links |
| `BITBUCKET_LIST_PIPELINES` / `BITBUCKET_GET_PIPELINE_STEP_LOG` | Reads build runs and a step's log | read | pipeline uuids, build numbers, the log |
| `BITBUCKET_CREATE_PULL_REQUEST_COMMENT` | Comments on a pull request | write, visible to others | the comment id |
| `BITBUCKET_CREATE_ISSUE` | Opens an issue | write, visible to others | the issue id and link |
| `BITBUCKET_CREATE_ISSUE_COMMENT` | Comments on an issue | write, visible to others | the comment id |
| `BITBUCKET_CREATE_BRANCH` | Makes a branch from a commit | write | the branch name and its commit |
| `BITBUCKET_CREATE_PULL_REQUEST` | Opens a pull request between two branches | write, visible to others | the pull request id and link |

**What to keep:** the workspace and repository slugs, and a pull request's
or issue's id, to come back to the same one. Never search for something the
taf made: use the id it kept.

**Start read-only.** "Summarise open pull requests" or "tell me when a build
fails" needs only read tools. Comments, issues and pull requests are seen by
the team and can notify them: offer a writing tool only when the owner asked
for it, and show each one to the owner first.

## Don't offer

- Deleting: `BITBUCKET_DELETE_REPOSITORY`, `BITBUCKET_DELETE_ISSUE`,
  `BITBUCKET_DELETE_PULL_REQUEST_COMMENT`, `BITBUCKET_DELETE_COMMIT_COMMENT`
  and the other delete tools.
- Deciding on code: `BITBUCKET_MERGE_PULL_REQUEST`,
  `BITBUCKET_DECLINE_PULL_REQUEST`, `BITBUCKET_APPROVE_PULL_REQUEST`,
  `BITBUCKET_REQUEST_PULL_REQUEST_CHANGES`. They act as a reviewer and
  change what others build on.
- Admin and settings: `BITBUCKET_CREATE_REPOSITORY`, and every pipeline
  variable tool (`BITBUCKET_CREATE_TEAMS_PIPELINES_CONFIG_VARIABLES`,
  `BITBUCKET_UPDATE_USERS_PIPELINES_CONFIG_VARIABLES` and the like).
  Variables can hold secrets; don't read them either
  (`BITBUCKET_GET_REPOSITORIES_PIPELINES_CONFIG_VARIABLES`,
  `BITBUCKET_GET_DEPLOYMENT_ENVIRONMENT_VARIABLES`).
- Stopping builds: `BITBUCKET_STOP_PIPELINE`.
- People and permissions: `BITBUCKET_LIST_WORKSPACE_MEMBERS`,
  `BITBUCKET_GET_USER_EMAILS`, `BITBUCKET_GET_USER_PERMISSIONS_WORKSPACES`,
  `BITBUCKET_GET_USER_PERMISSIONS_REPOSITORIES`, `BITBUCKET_GET_SSH_LATEST_KEYS`.
- Bulk downloads: `BITBUCKET_DOWNLOAD_REPOSITORY_ARCHIVE`.

## Access and plans

- **Pull request access is separate from repository access.** Reading
  repositories doesn't let a taf see pull requests. Reading pull requests
  also lets it comment on them; opening a pull request needs the right to
  write pull requests, which also lets it write to the repository.
- **Issues have their own access,** apart from repositories.
- **Rate:** Bitbucket counts calls per user over a rolling hour. Reading
  repository data allows about 1,000 calls an hour for a signed-in user,
  file downloads 5,000, archives 1,500. Without a sign-in, only 60 an hour
  for everything. Larger allowances on paid plans apply to workspace,
  project or repository tokens, not to a person's sign-in.
- **Big answers:** a list answers a page of 10 to 100 items, with a link to
  the next page. Follow that link; don't build page numbers.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The sign-in no longer works (401) | The owner removed access, or the sign-in lapsed after three months unused | "Bitbucket needs you to sign in again. Nothing was written." |
| Refused for missing access | The sign-in didn't give this tool's access, or the owner's account can't use that repository | "Bitbucket is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Not found | The workspace, repository, pull request or file doesn't exist, or this account can't see it | "I can't find that in Bitbucket. Check the name, and that your account can open it." |
| Refused, too many calls | Over the hourly limit for this user | "Bitbucket's hourly limit is reached. The taf will continue when it frees up." |
| Refused, with an error per field | The input isn't valid: a branch that doesn't exist, a title missing | Fix it and try once; else "Bitbucket didn't accept that change. Nothing was written." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
