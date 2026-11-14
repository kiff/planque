# GitHub

Checked: 2026-11-14 · Sources: the connection service's GitHub tool
definitions as planque reads them, GitHub's REST API description as
`api/operations` reads it (1,277 operations), and test runs on a phone on
13 Nov 2026.

## How it connects

- **Public repositories need no sign-in.** Reading a public repository's
  files, releases or issues can go by the web read tool. Offer a sign-in only
  when the job needs the owner's private repositories or writes anything.
- **A sign-in to GitHub**, held by the connection service, is the usual way.
  The catalog lists its tools and, with a known sign-in, puts the ones its
  access doesn't cover under `unusable`. A held sign-in is reused for a new
  taf when it covers the tools asked.
- **Or a key**, a fine-grained token the owner makes in GitHub's settings,
  through `api/operations`. GitHub's API is big: ask for operations by a few
  words from the job, and name only those in the request.

## Tools a job usually needs

There are hundreds; a job needs a handful. These are the ones seen working.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| Get repository content | Reads one file or folder of a repository | read | the file's content, its sha |
| List releases / Get latest release | Reads a repository's releases | read | release ids, tags, links |
| List issues / Get an issue | Reads issues and their comments | read | issue numbers, links |
| List pull requests / Get a pull request | Reads pull requests and their state | read | numbers, links, head sha |
| List commits | Reads a branch's recent commits | read | commit shas, links |
| Create an issue | Opens an issue | write, visible to others | the issue number and link |
| Comment on an issue or pull request | Adds a comment | write, visible to others | the comment id and link |
| Create or update a file | Commits one file to a branch | write | the commit sha |

**What to keep:** an issue's or pull request's number and the repository's
`owner/name`; a file's sha, to update that same version next time.

**Writes are public in a public repository,** and notify people. Every
comment or issue the taf writes is shown to the owner first.

## Don't offer

- Deleting a repository, branch or release; changing a repository's
  settings, collaborators or visibility; managing secrets, keys or
  webhooks. None is a job a taf should hold.
- Merging pull requests, or pushing to the default branch: they change what
  others build on.
- Anything in an organisation's admin area.

## Access and plans

- **Reading and writing issues, files and releases** works without a paid
  plan in the cases checked here. Which other features need one isn't
  checked here.
- **An organisation can restrict outside apps.** Then the sign-in works for
  the owner's own repositories but not the organisation's until an admin
  approves the app.
- **Rate:** about 5,000 calls an hour for a signed-in user, and only 60 an
  hour without a sign-in. Searches have a lower limit, about 30 a minute.
- **Big answers:** a list answers a page, at most 100 items. Files over 1 MB
  come back without their content and need a second, raw read; files over
  100 MB can't be read this way at all.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Not found (404) on a private repository | The sign-in can't see it: a private repository answers "not found", not "forbidden" | "I can't see that repository. It may be private and not shared with this sign-in." |
| Refused (403), organisation | The organisation hasn't approved the app | "Your GitHub organisation needs to approve this connection. An admin can do it in the organisation's settings." |
| Refused (403 or 429), rate limit | Over the hourly limit, or too many writes in a short time | "GitHub's hourly limit is reached. The taf will continue when it resets." |
| Refused (409), conflict | The file or branch changed since it was read | Read it again and retry once; tell the owner only if it fails again. |
| Refused (422) | The input isn't valid: a branch that exists, a title missing | Fix it and try once; else "GitHub didn't accept that change. Nothing was written." |
| The sign-in no longer works (401) | Revoked, or the token expired | "GitHub needs you to sign in again. Nothing was written." |

After any refusal the taf stops that step and says so.
