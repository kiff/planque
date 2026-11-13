# Reddit

Checked: 2026-11-13 · Sources: the connection service's Reddit tool
definitions as planque reads them (48 tools), Reddit's live API reference
(reddit.com/dev/api/oauth), its Data API Wiki (updated 3 Nov 2026) and its
Data API Terms.

## How it connects

- **A sign-in to the owner's Reddit account**, held by the connection
  service; the taf never sees a password or token. Each tool needs its own
  kind of access (reading, posting, editing, private messages, moderation).
  With a known sign-in, the catalog puts the tools it doesn't cover under
  `unusable`, with why.
- **Everything the taf writes is posted as the owner,** under their
  username, in public.
- **Ask for every tool in the first request.** A held sign-in is reused for
  a new taf only when it covers the tools asked.

## Tools a job usually needs

Read changes nothing on Reddit. Returns are the ids Reddit's reference
documents: posts and comments carry a "fullname" (`t3_…` for a post, `t1_…`
for a comment), and lists carry an `after` cursor for the next page.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `REDDIT_RETRIEVE_REDDIT_POST` | Lists a subreddit's posts, new, hot or top | read | post fullnames, links, an `after` cursor |
| `REDDIT_SEARCH_ACROSS_SUBREDDITS` | Searches posts, in one subreddit or all | read | post fullnames, an `after` cursor |
| `REDDIT_RETRIEVE_POST_COMMENTS` | Reads a post and its comment tree | read | comment fullnames, stubs for more |
| `REDDIT_EXPAND_MORE_COMMENTS` | Fetches comments left out of a long thread, up to 100 at a time | read | comment fullnames |
| `REDDIT_GET_POSTS_BY_FULLNAME` | Reads posts the taf already knows | read | the posts |
| `REDDIT_GET_SUBREDDIT_RULES` | Reads a subreddit's rules | read | the rules |
| `REDDIT_GET_SUBREDDIT_METADATA` | Reads a subreddit's description and settings | read | its fullname, subscriber count |
| `REDDIT_LIST_SUBREDDIT_POST_FLAIRS` | Lists the flairs a post may carry | read | flair template ids |
| `REDDIT_LIST_INBOX_NOTIFICATIONS` | Lists mentions of the owner | read | message fullnames |
| `REDDIT_CREATE_REDDIT_POST` | Posts a link or text to a subreddit | write, public | the post's id, fullname and link |
| `REDDIT_POST_REDDIT_COMMENT` | Replies to a post or comment | write, public | the new comment's fullname |
| `REDDIT_EDIT_REDDIT_COMMENT_OR_POST` | Changes the text of the owner's own comment or text post | write, public | the edited item |

**What to keep:** the fullname of every post or comment the taf wrote, to
edit or follow it later. Keep ids, not copies of other people's posts:
Reddit's terms require removing content that has been deleted on Reddit.

**Start read-only.** "Watch these subreddits and summarise" or "tell me when
I'm mentioned" is reads only. Offer posting or commenting only when the
owner asked for it; each post or comment is shown to them before it goes
out. Read `REDDIT_GET_SUBREDDIT_RULES` before offering to post anywhere.

## Don't offer

- Deleting: `REDDIT_DELETE_REDDIT_POST`, `REDDIT_DELETE_REDDIT_COMMENT`.
- Moderation and subreddit settings: `REDDIT_CONFIGURE_FLAIR_TEMPLATE`,
  `REDDIT_GET_SUBREDDIT_TRAFFIC`, `REDDIT_LIST_MODERATION_QUEUE`,
  `REDDIT_LIST_MODERATION_LOG`, `REDDIT_LIST_MODMAIL_CONVERSATIONS`,
  `REDDIT_GET_MODMAIL_CONVERSATION`, `REDDIT_GET_USER_FLAIR`. They only
  work in subreddits the owner moderates, and aren't a taf's job.
- Voting (`REDDIT_VOTE_ON_CONTENT`): an automated vote counts as the
  owner's own.
- Private messages to people (`REDDIT_SEND_PRIVATE_MESSAGE`): unsolicited
  messages from an automated account read as spam.
- Account changes and bulk actions: `REDDIT_MANAGE_SUBREDDIT_SUBSCRIPTION`,
  `REDDIT_MARK_ALL_INBOX_MESSAGES_READ`, `REDDIT_TOGGLE_INBOX_REPLIES`.
- `REDDIT_PREPARE_MEDIA_UPLOAD`: it isn't in Reddit's public reference.

## Access and plans

- **Free access, with rules.** Reddit's terms allow use within its free
  limits; commercial use, or anything beyond them, needs a separate
  agreement with Reddit. Reddit's content may not be used to train AI
  models without the rightsholders' permission.
- **Rate:** 100 calls a minute, averaged over 10 minutes, per registered
  app rather than per user, so other people's use of the same app counts.
  A job that reads many subreddits each run should page and pause.
- **Pages:** a list returns 25 items by default, at most 100.
- **Private messages** are a separate access; a sign-in made for reading
  doesn't include them.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The sign-in no longer works | The owner removed access, or the sign-in expired | "Reddit needs you to sign in again. Nothing was posted." |
| Refused for missing access | The sign-in didn't give this tool's access, posting say | "Reddit is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without it?" |
| Refused on a subreddit | It's private, or the owner is banned there | "I can't reach that subreddit from your account. It may be private, or closed to you." |
| Too many calls | Over the per-minute limit | "Reddit asked us to slow down. The taf will try again in a few minutes." |
| A post is refused by the subreddit | The subreddit's rules or settings don't allow it | "That subreddit didn't accept the post. Its rules may not allow it from your account. Nothing was posted." |
| Not found | The post, comment or subreddit was deleted | "That isn't on Reddit any more. The taf will skip it." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
