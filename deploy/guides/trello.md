# Trello

Checked: 2026-11-13 · Sources: the connection service's Trello tool
definitions as planque reads them (329 tools), and Trello's REST API and help
documentation.

## How it connects

- **A sign-in to Trello**, held by the connection service; the taf never sees
  a password or token. planque's catalog lists the tools and, with a known
  sign-in, puts the ones its access doesn't cover under `unusable`, with why.
- **Read and write are separate access.** A sign-in made for reading boards
  doesn't let the taf add or move cards.
- **An account that already holds Trello reuses that sign-in** for a new taf
  when it covers the tools asked. Ask for every tool the job needs in the
  first request, so nothing is asked twice.
- **The taf sees what the signed-in person sees:** their boards, and the
  boards of workspaces they belong to. Nothing else.

## Tools a job usually needs

Read means it changes nothing on any board. Returns are what a taf keeps to
use next time; the catalog's `returns` gives exact paths where it has them.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `TRELLO_GET_MEMBERS_BOARDS_BY_ID_MEMBER` | Lists the boards the signed-in person can open | read | board ids, names, links |
| `TRELLO_GET_BOARDS_BY_ID_BOARD` | Reads one board's details | read | the board id, its link |
| `TRELLO_GET_BOARDS_LISTS_BY_ID_BOARD` | Lists a board's columns (lists) | read | list ids and names |
| `TRELLO_GET_LISTS_CARDS_BY_ID_LIST` | Lists the cards in one column | read | card ids, names, links |
| `TRELLO_GET_BOARDS_CARDS_BY_ID_BOARD` | Lists every card on a board | read | card ids, their list ids |
| `TRELLO_GET_CARDS_BY_ID_CARD` | Reads one card: title, description, due date, labels | read | the card, its list id, its link |
| `TRELLO_GET_CARDS_ACTIONS_BY_ID_CARD` | Reads a card's history and comments | read | action ids, comment text |
| `TRELLO_GET_SEARCH` | Finds boards and cards by words | read | matching board and card ids |
| `TRELLO_ADD_CARDS` | Adds a card to a column | write | the new card's id and link |
| `TRELLO_UPDATE_CARDS_BY_ID_CARD` | Changes a card's title, description, due date or labels | write | the card |
| `TRELLO_UPDATE_CARDS_ID_LIST_BY_ID_CARD` | Moves a card to another column | write | the card |
| `TRELLO_ADD_CARDS_ACTIONS_COMMENTS_BY_ID_CARD` | Adds a comment to a card | write | the comment's id |
| `TRELLO_ADD_CARDS_CHECKLISTS_BY_ID_CARD` | Adds a checklist to a card | write | not documented by Trello; read the card again for its checklist ids |
| `TRELLO_ADD_CHECKLISTS_CHECK_ITEMS_BY_ID_CHECKLIST` | Adds an item to a checklist | write | not documented by Trello |

**What to keep:** the board id and the column (list) ids the job works in,
so it never searches for them again; each card's id once it adds one, to
update, move or comment on that same card later. Never search by title for
a card the taf made: use the id it kept.

**Start read-only.** "Tell me what's overdue on the board" needs only read
tools. Offer adding, moving or commenting only when the owner asked for it,
and name the board and columns the taf may touch.

## Don't offer

- **Deleting** (`TRELLO_DELETE_CARDS_BY_ID_CARD` and the other delete tools):
  a deleted card can't be brought back. Moving a card to a "Done" column, or
  archiving it with `TRELLO_UPDATE_CARDS_CLOSED_BY_ID_CARD`, keeps it.
- **Board and workspace admin:** creating or closing boards, changing members
  or permissions, and workspace settings. They change things for everyone on
  the board.
- **Webhooks and power-ups:** they set up lasting links that outlive the job.

## Access and plans

- **Free plan:** a free workspace has at most 10 collaborators. Anyone
  invited to one of its boards counts. Adding a person can fail once that's
  reached.
- **Attachments:** up to 10 MB a file on the free plan, 250 MB on paid plans.
  There's no total storage cap.
- **Rate:** Trello allows about 100 calls every 10 seconds for one sign-in,
  and fewer for search and for reading people's profiles. A job that reads a
  large board should read whole lists or the board in one call, not card by
  card.
- **Big reads:** asking for many cards with their full history in one call
  can be refused. Read the cards first, then each card's history.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (401 or 403) | The sign-in gave read access only, or the person isn't on that board | "Trello is only partly connected: this step needs access you didn't give, or that board isn't shared with you. Want to sign in again and allow it, or do without this step?" |
| The sign-in no longer works (401, "invalid token") | The owner revoked it in Trello's settings, or it expired | "Trello needs you to sign in again. Nothing was changed." |
| Too many calls (429) | Over Trello's short-term limit | "Trello asked us to slow down. The taf will try again in a moment." |
| Too much asked at once (429, too many cards) | One call asked for too many cards with their history | "That board is too big to read in one go. The taf will read it in parts." |
| A board, list or card isn't found (404) | It was deleted, or moved to a board the person can't open | "That card or board isn't there any more. The taf will skip it." |
| An upload is refused | The file is over the plan's size limit | "That file is too big for your Trello plan. Want a link to it instead?" |
| A collaborator can't be added | The free workspace already has 10 | "Your free Trello workspace is full. Remove someone or upgrade, then try again." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
