# Confluence

Checked: 2026-11-13 · Sources: [the generated part](generated/confluence.md)
(70 tools), and Atlassian's Confluence Cloud developer docs: the REST API v2
reference, the rate-limiting page and the OAuth 2.0 (3LO) apps page.

## How it connects

- **A sign-in to Atlassian**, held by the connection service; the taf never
  sees a password or token. At sign-in Atlassian can ask the owner which
  Confluence site to allow. planque's catalog lists the tools and, with a
  known sign-in, puts the ones its access doesn't cover under `unusable`.
- **The taf can do only what the owner can.** Even with write access, a page
  in a space the owner can't edit stays out of reach.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.

## Tools a job usually needs

Read means it changes nothing on the site. Returns are the ids a taf keeps;
they come from Atlassian's API reference, since the catalog lists none.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `CONFLUENCE_GET_SPACES` | Lists the spaces the owner can see, up to 250 a page | read | space ids, keys |
| `CONFLUENCE_SEARCH_PAGES` / `CONFLUENCE_GET_PAGES` | Finds pages by title or space, up to 250 a page | read | page ids, titles, a next-page link |
| `CONFLUENCE_CQL_SEARCH` | Searches all content with Confluence's query language | read | content ids, titles, excerpts |
| `CONFLUENCE_GET_PAGE_BY_ID` | Reads one page | read | the page, its version number, its space id |
| `CONFLUENCE_GET_CHILD_PAGES` | Lists the pages under a page | read | page ids |
| `CONFLUENCE_GET_PAGE_FOOTER_COMMENTS` | Reads a page's comments | read | comment ids |
| `CONFLUENCE_GET_ATTACHMENTS` / `CONFLUENCE_DOWNLOAD_ATTACHMENT` | Lists or downloads a page's files | read | attachment ids, the file |
| `CONFLUENCE_CREATE_PAGE` | Adds a page to a space, under a parent page if given | write, visible to others | the page id, its link |
| `CONFLUENCE_UPDATE_PAGE` | Replaces a page's title and whole body | write, visible to others | the page id, its new version number |
| `CONFLUENCE_CREATE_FOOTER_COMMENT` | Comments on a page or replies to a comment | write, visible to others | the comment id |
| `CONFLUENCE_ADD_CONTENT_LABEL` | Adds labels to a page | write | the labels |

**What to keep:** the space id and the id of every page the taf made, plus
the version number it last saw. Never search for a page the taf made: use
the id it kept.

**Updating replaces the whole page.** Read the page first, change what the
job asks, send the full body back, and set the version to the one read plus
one. A page others edit is shown to the owner before it is changed.

**Start read-only.** "Summarise our team's pages" needs only reads. A new
page is published as soon as it's made, so the taf shows the owner its text
first.

## Don't offer

- Deleting (`CONFLUENCE_DELETE_PAGE`, `CONFLUENCE_DELETE_BLOGPOST`,
  `CONFLUENCE_DELETE_SPACE`, and the `DELETE_..._PROPERTY` tools). A deleted
  blog post can be purged for good.
- Spaces: `CONFLUENCE_CREATE_SPACE`, `CONFLUENCE_CREATE_PRIVATE_SPACE` and
  the space property tools change how a space is set up.
- Admin and people: `CONFLUENCE_GET_AUDIT_LOGS`,
  `CONFLUENCE_GET_CONTENT_RESTRICTIONS` and `CONFLUENCE_SEARCH_USERS`.
- Content property tools (create, update, delete on pages, blog posts and
  whiteboards): hidden settings other apps rely on, not a job's content.

## Access and plans

- **Reading and writing need separate access.** A sign-in made for reading
  can't create or update pages; deleting needs its own access again.
- **Rate:** Atlassian counts points an hour per app, reset at the top of each
  UTC hour: a call costs at least one point, and more for each page, space
  or user it returns. By default the hourly allowance is shared by every
  site using the same app. Once it's spent, every call is refused until the
  hour resets. A job reading many pages should page, ask only for the fields
  it needs, and pause.
- **Big pages:** a new page over 5 MB is refused.
- **Content approvals:** where a space requires approval before publishing,
  an update from the taf is refused.
- **The sign-in ends** after 90 days without use, or when the owner changes
  their Atlassian password.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Not found (404) | The page or space doesn't exist, or the owner can't see or edit it: Confluence answers "not found" for both | "I can't reach that page. It may have been moved, deleted, or be restricted from you." |
| The sign-in no longer works (401) | The owner removed access, changed their password, or 90 days passed unused | "Confluence needs you to sign in again. Nothing was changed." |
| Refused for missing access | The sign-in was made for reading and the tool writes | "Confluence is only partly connected: this step needs access you didn't give. Want to sign in again and allow it?" |
| Update refused for the version number | Someone changed the page since the taf read it | Read the page again and retry once; else "Someone edited that page while I was working. I've left it as it is." |
| Refused (409) on update | The space needs approval before publishing | "This space needs an approval before changes go live, so I couldn't change the page. I can give you the new text to submit yourself." |
| Too large (413) | The page is over 5 MB | "That page is too big for Confluence. Can I split it into several pages?" |
| Too many calls (429) | The hourly allowance is spent; Confluence says how long to wait | "Confluence asked us to slow down. The taf will continue after the wait it gave." |
| Refused as invalid (400) | An input isn't valid | Fix it and try once; else "Confluence didn't accept that change. Nothing was written." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
