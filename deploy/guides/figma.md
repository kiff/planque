# Figma

Checked: 2026-11-14 · Sources: the connection service's Figma tool
definitions as planque reads them (52 tools, [generated](generated/figma.md)),
and Figma's REST API reference.

## How it connects

- **A sign-in to Figma**, held by the connection service. The catalog lists
  the tools and, with a known sign-in, puts the ones its access doesn't cover
  under `unusable`. Offer only what's under `tools`.
- **Figma's API reads designs; it can't draw.** A taf can read files, export
  images, read and write comments. It can't change a design.

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `FIGMA_GET_PROJECTS_IN_A_TEAM` | Lists a team's projects | read | project ids |
| `FIGMA_GET_FILES_IN_A_PROJECT` | Lists a project's files | read | **file keys** |
| `FIGMA_GET_FILE_METADATA` | Reads a file's name, owner, last change | read | |
| `FIGMA_GET_FILE_NODES` | Reads chosen frames of a file | read | node ids |
| `FIGMA_RENDER_IMAGES_OF_FILE_NODES` | Exports frames as PNG, SVG or PDF | read | the image links, which expire |
| `FIGMA_GET_VERSIONS_OF_A_FILE` | Lists a file's saved versions | read | version ids |
| `FIGMA_GET_COMMENTS_IN_A_FILE` | Reads a file's comments | read | comment ids |
| `FIGMA_ADD_A_COMMENT_TO_A_FILE` | Adds a comment, or a reply | write, visible to others | the comment's id |
| `FIGMA_EXTRACT_DESIGN_TOKENS` | Reads a file's colours, type and spacing | read | |

**What to keep:** the file key (in the file's link, after `/design/`); the
ids of the frames the job exports. A team's id is in the team page's link:
the API can't list the owner's teams, so ask for that link once.

**Exported images are links that expire** after 30 days. Download what
needs keeping.

## Don't offer

- Webhooks (`FIGMA_CREATE_A_WEBHOOK` and the others): a taf reads on its
  schedule instead.
- Variables (`FIGMA_CREATE_MODIFY_DELETE_VARIABLES`): they change a design
  system everyone uses, and need an Enterprise plan.
- Activity logs, payments and user provisioning: they're for an
  organisation's admin.
- Deleting comments other than the taf's own.

## Access and plans

- **Reading files and comments works on any plan,** for files the owner can
  open, but reading files and exporting are capped: about 6 calls a month for
  a file on the free Starter plan, and about 20 a month on a view or collab
  seat. Offer file jobs only to an owner with a full or dev seat on a paid
  plan.
- **Library analytics, variables and activity logs** need Figma's Enterprise
  plan; activity logs also need an organisation admin.
- **Rate:** Figma limits calls per minute by plan, and image exports are
  heavier: exporting many frames at once should be split into small batches.
- **Big files** answer slowly when read whole; read the frames the job needs.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused (403) | The owner can't open that file, or the plan doesn't include the feature | "I can't open that Figma file. Check it's shared with you, or it may need a higher Figma plan." |
| Not found (404) | A wrong file key, or a deleted frame | "That Figma file or frame isn't there. Send me its link again." |
| Too many calls (429) | Over Figma's limit | "Figma asked us to slow down. The taf will continue in a minute." |
| An export comes back empty | The frame is hidden or empty | "That frame exported nothing. It may be hidden in the file." |
| The sign-in no longer works (401) | Removed, or expired | "Figma needs you to sign in again." |

After any refusal the taf stops that step and says so.
