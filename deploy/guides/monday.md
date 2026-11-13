# Monday

Checked: 2026-11-13 · Sources: the connection service's Monday tool
definitions as planque reads them (122 tools), and monday.com's API
documentation (rate limits, errors, items, updates, sign-in for outside
apps).

## How it connects

- **A sign-in to Monday**, held by the connection service; the taf never
  sees a password or token. The person picks what to allow on Monday's own
  page. planque's catalog lists the tools and, with a known sign-in, puts the
  ones its access doesn't cover under `unusable`, with why.
- **The taf acts for the person who signed in**, within what they allowed,
  and sees only the boards they can open.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.

## Tools a job usually needs

Read means it changes nothing in Monday. The catalog gives no `returns`
paths for Monday today; the ids below are those Monday's API reference
shows where it documents the answer.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `MONDAY_GET_ME` | Reads who the sign-in acts as | read | not documented |
| `MONDAY_GET_WORKSPACES` | Lists the workspaces the person can open | read | not documented |
| `MONDAY_BOARDS` | Lists boards, or reads chosen ones | read | not documented |
| `MONDAY_COLUMNS` | Reads a board's columns, to know where each value goes | read | not documented |
| `MONDAY_LIST_GROUPS` | Lists a board's groups (its sections) | read | not documented |
| `MONDAY_ITEMS_PAGE` | Reads a board's items, a page at a time | read | item ids and names, a page cursor |
| `MONDAY_LIST_ITEMS_BY_COLUMN_VALUES` | Finds items whose column holds a given value | read | not documented |
| `MONDAY_GET_ITEMS` | Reads chosen items by id | read | not documented |
| `MONDAY_GET_UPDATES` | Reads the comments (updates) on items | read | update ids, text, author, date |
| `MONDAY_CREATE_ITEM` | Adds an item to a board, optionally in a group with values | write | the new item's id |
| `MONDAY_CHANGE_SIMPLE_COLUMN_VALUE` | Sets one column of one item | write | not documented |
| `MONDAY_UPDATE_ITEM` | Changes an item's values | write | not documented |
| `MONDAY_MOVE_ITEM_TO_GROUP` | Moves an item to another group on its board | write | the item's id |
| `MONDAY_CREATE_UPDATE` | Posts a comment on an item, or a reply | write | the new update's id |

**What to keep:** each board's id and the ids of its columns and groups;
each item's id, to change or comment on that same item later. Never look an
item up by its name: use the id the taf kept.

**Read the columns first.** A value has to fit its column's type; a status
or date sent in the wrong form is refused.

**Start read-only.** "Tell me what's stuck on this board" needs only read
tools. Offer writes when the owner asked for changes. Comments can mention
and notify people.

## Don't offer

- Deleting and archiving tools (items, boards, groups, columns, docs,
  workspaces, folders, updates): they remove work other people rely on.
- People and permission tools (inviting, activating or deactivating users,
  changing roles, teams, workspace members, board permissions, the email
  domain): they change who can do what in the account.
- Board structure tools (creating or changing columns and boards,
  duplicating boards): they change the board for everyone using it.
- `MONDAY_GET_MUTATION_COMPLEXITY`: it runs any change it's given, so it
  counts as a write. `MONDAY_CUSTOM_ACTIVITY` is also marked as a write.
- Webhook tools: the taf reads on its schedule, it doesn't listen.

## Access and plans

- **Daily calls depend on the account's plan,** counted from midnight UTC:
  200 on Free, 1,000 on Basic and Standard, 10,000 on Pro and trials,
  25,000 on Enterprise. A job reading many boards a day won't fit Free.
- **A minute limit too:** 1,000 calls a minute (2,500 on Pro, 5,000 on
  Enterprise), and at most 40 at once (100 on Pro, 250 on Enterprise).
- **Complexity:** each call is scored by how much data it asks for, against
  a budget per minute. Smaller pages use less of it.
- **A board holds at most 10,000 items.** Adding more is refused.
- **Item names** are 1 to 255 characters.
- **Pages:** at most 500 items a page when read through a board.
- **Monday answers many refusals as an ordinary answer** with an error
  inside, not as an error status. The taf reads the answer for errors even
  when the call looks successful.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Missing access (`missingRequiredPermissions`) | The sign-in didn't allow this kind of action | "Monday is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Not allowed (403) | The person can't do this on that board, or can't use the API (view-only, inactive, or email not confirmed) | "Monday didn't let the taf do this. Check your role on that board, or do without this step." |
| Board or item not found | The id is wrong, it was deleted, or the person can't open it | "The taf can't see that board or item in Monday. It may be deleted, or not shared with you." |
| A value doesn't fit its column | Wrong form for the column's type | "Monday didn't accept a value for that column. The taf has stopped and will show you what it tried." |
| Board full | The board has 10,000 items | "That Monday board is full. Archive some items, or pick another board." |
| Board busy (423) | Another change is being saved on the board | "Monday was busy saving another change. The taf will try again in a moment." |
| Too many calls (429) | Over the minute, at-once or complexity limit; Monday says when to retry | "Monday asked us to slow down. The taf will try again shortly." |
| Daily limit reached | The plan's daily calls are used up | "Monday's daily limit for your plan is reached. The taf will carry on after midnight UTC." |
| The sign-in no longer works (401) | Access was removed, or the app was uninstalled | "Monday needs you to sign in again. Nothing was changed." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write: a taf that can't find its item asks, rather than adding
another one.
