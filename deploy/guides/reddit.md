# Reddit

Checked: 2026-11-13 · Sources: the connection service's Reddit tool
definitions as planque reads them (48 tools), Reddit's API reference, and
Reddit's Data API terms.

## How it connects

- **A sign-in to Reddit**, held by the connection service; the taf never
  sees a password or token. planque's catalog lists the tools and, with a
  known sign-in, puts the ones its access doesn't cover under `unusable`,
  with why.
- **The sign-in acts as one Reddit account.** Everything the taf posts
  appears under that account's name, in public.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.

## Tools a job usually needs

Read means it changes nothing on Reddit. The catalog's `returns` is empty
for Reddit today. Reddit names every post and comment by a "fullname": a
type prefix (`t3_` for a post, `t1_` for a comment) and an id. Editing,
deleting and replying all take that fullname. Lists come a page at a time,
at most 100 items, with a marker for the next page.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `REDDIT_GET_REDDIT_USER_ABOUT` | Reads an account's public profile | read | not documented |
| `REDDIT_GET_SUBREDDIT_RULES` | Reads a community's rules | read | the rules |
| `REDDIT_GET_SUBREDDIT_METADATA` | Reads a community's details | read | not documented |
| `REDDIT_SEARCH_ACROSS_SUBREDDITS` | Finds posts by words | read | post fullnames, a next-page marker |
| `REDDIT_RETRIEVE_REDDIT_POST` | Lists a community's posts | read | post fullnames, a next-page marker |
| `REDDIT_RETRIEVE_POST_COMMENTS` | Reads a post and its comments | read | the post, comment fullnames |
| `REDDIT_RETRIEVE_SPECIFIC_COMMENT` | Reads one comment or post | read | not documented |
| `REDDIT_LIST_USER_ACTIVITY` | Lists an account's own posts and comments | read | their fullnames |
| `REDDIT_LIST_INBOX_NOTIFICATIONS` | Lists mentions of the account | read | not documented |
| `REDDIT_CREATE_REDDIT_POST` | Posts a link or text to a community | write, public | not documented |
| `REDDIT_POST_REDDIT_COMMENT` | Replies to a post or a comment | write, public | not documented |
| `REDDIT_EDIT_REDDIT_COMMENT_OR_POST` | Changes the text of the account's own post or comment | write, public | not documented |
| `REDDIT_DELETE_REDDIT_POST` | Removes one of the account's own posts | write, public | not documented |

**What to keep:** the fullname of each post and comment the taf made, to
edit, delete or reply to it later. Reddit's reference doesn't show what a
new post's answer holds, so if it carries no fullname, the taf finds its
post in the account's own activity, not by searching the whole site.

**Everything written is public.** Each post, comment and edit is drafted,
shown to the owner, and sent only on their yes. Read the community's rules
first: many limit links, require text, or ban self-promotion. Offer deleting
only for posts the taf made, and each one waits for the owner's yes too.

## Don't offer

- Voting (`REDDIT_VOTE_ON_CONTENT`): it changes how other people see
  content, and it isn't a job for a taf.
- Private messages (`REDDIT_SEND_PRIVATE_MESSAGE`): writing to strangers
  is what Reddit's terms call spam.
- Moderation (`REDDIT_LIST_MODERATION_QUEUE`, `REDDIT_LIST_MODERATION_LOG`,
  `REDDIT_CONFIGURE_FLAIR_TEMPLATE`, the modmail tools): they act for a
  community, not for the person.
- Joining or leaving communities (`REDDIT_MANAGE_SUBREDDIT_SUBSCRIPTION`)
  and inbox housekeeping (`REDDIT_MARK_ALL_INBOX_MESSAGES_READ`,
  `REDDIT_TOGGLE_INBOX_REPLIES`): they change the account, and marking all
  read hides messages from the owner.
- Posting the same text to many communities: it can read as spam under
  Reddit's terms.

## Access and plans

- **Reading needs read access; posting and replying need posting access;
  editing and deleting need edit access.** A sign-in made for reading
  doesn't let the taf post.
- **Reddit's terms:**
  - commercial use of Reddit's API needs a separate agreement with Reddit;
  - spamming, incentivising or harassing people is not allowed;
  - data read from Reddit is kept only as long as the job needs it;
  - Reddit can block an app for good if it calls too much or abuses it.

- **Rate:** Reddit limits how often an app may call; the numbers weren't
  checked here. A taf posts at most a few times a day and pages through
  reads.
- **Communities set their own rules** for what a post must contain (a body
  required or not allowed, some link domains banned, some words banned).
  A post that breaks them is refused.

## When it fails

Reddit's reference lists few error codes, so these rows give none.

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access | The sign-in didn't give this tool's access: reading only, say, and the tool posts | "Reddit is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| The sign-in no longer works | The person removed the app, or the sign-in expired | "Reddit needs you to sign in again. Nothing was posted." |
| Post refused by the community | It breaks that community's rules (no links, text required, banned words) | "That community didn't accept the post: it doesn't fit their rules. Want to change it, or post somewhere else?" |
| Too many calls | Over Reddit's call limit | "Reddit asked us to slow down. The taf will try again later." |
| Not found | The post or comment was deleted | "That post isn't there any more. The taf will skip it." |
| Blocked by Reddit | Reddit has restricted this account or app | "Reddit has blocked this account or connection. The taf has stopped. Check the account on Reddit." |

After any refusal the taf stops that step and says so. A refused post is
never sent again with the text changed unless the owner changed it, and a
refused read never leads to a write.
