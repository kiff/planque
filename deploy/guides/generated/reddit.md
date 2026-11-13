# Reddit: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

48 tools: 35 read, 13 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `REDDIT_CONFIGURE_FLAIR_TEMPLATE` | Configure flair template | write | `subreddit`; *allowable_content, background_color, css_class, flair_template_id, flair_type, max_emojis, mod_only, text and 2 more* |  | `modflair` |
| `REDDIT_CREATE_REDDIT_POST` | Create a Reddit post | write | `subreddit`, `title`; *flair_id, kind, text, url* |  | `submit` |
| `REDDIT_DELETE_REDDIT_COMMENT` | Delete Reddit comment | write | `id` |  | `edit` |
| `REDDIT_DELETE_REDDIT_POST` | Delete a Reddit post | write | `id` |  | `edit` |
| `REDDIT_EDIT_REDDIT_COMMENT_OR_POST` | Edit comment or post | write | `text`, `thing_id` |  | `edit` |
| `REDDIT_EXPAND_MORE_COMMENTS` | Expand More Comments | read | `children`, `link_id`; *more_id, sort* |  | `read` |
| `REDDIT_GET` | Get Reddit listing by sort | read | `sort`; *after, before, count, limit, show, time_filter* |  | `read` |
| `REDDIT_GET_AUTHENTICATED_KARMA_BREAKDOWN` | Get authenticated karma breakdown | read |  |  | `mysubreddits` |
| `REDDIT_GET_CONTROVERSIAL_POSTS` | Get controversial posts from all subreddits | read | *after, before, count, limit, show, t* |  | `read` |
| `REDDIT_GET_INBOX_MESSAGE` | Get inbox message | read | `message_id` |  | `privatemessages` |
| `REDDIT_GET_ME_PREFS` | Get user preferences | read | *fields* |  | `identity` |
| `REDDIT_GET_MODMAIL_CONVERSATION` | Get modmail conversation | read | `conversation_id` |  | `modmail` |
| `REDDIT_GET_POSTS_BY_FULLNAME` | Get Reddit posts by fullname | read | `fullnames` |  | `read` |
| `REDDIT_GET_RANDOM` | Get random Reddit post | read | *subreddit* |  | `read` |
| `REDDIT_GET_REDDIT_USER_ABOUT` | Get user information | read | `username` |  | `identity`, `read` |
| `REDDIT_GET_SCOPES` | Get OAuth scopes | read |  |  |  |
| `REDDIT_GET_SUBREDDITS_SEARCH` | Search subreddits | read | `q`; *after, before, count, limit, show, show_users, sort, sr_detail* |  | `read` |
| `REDDIT_GET_SUBREDDIT_METADATA` | Get subreddit metadata | read | `subreddit` |  | `read` |
| `REDDIT_GET_SUBREDDIT_RULES` | Get subreddit rules | read | `subreddit`; *raw_json* |  | `read` |
| `REDDIT_GET_SUBREDDIT_TRAFFIC` | Get subreddit traffic | read | `subreddit` |  | `modconfig` |
| `REDDIT_GET_USERNAME_AVAILABLE` | Check username availability | read | `user` |  | `identity` |
| `REDDIT_GET_USER_COMMENTS_RSS` | Get Reddit user comments RSS | read | `username`; *after, before, context, count, item_type, limit, show, sort and 2 more* |  | `history` |
| `REDDIT_GET_USER_FLAIR` | Get user flair | read | `subreddit`; *after, before, limit* |  | `modflair` |
| `REDDIT_GET_WIKI_PAGE` | Get subreddit wiki page | read | `page`, `subreddit`; *comparison_revision, raw_json, revision* |  | `wikiread` |
| `REDDIT_LIST_CUSTOM_FEED` | List custom Reddit feed | write | `custom_feed`, `username`; *after, before, count, limit, show, sort, sr_detail, time_filter* |  | `read` |
| `REDDIT_LIST_INBOX_MESSAGES` | List inbox messages | read | `inbox_view`; *after, before, count, limit, max_replies, mid, show, sr_detail* |  | `privatemessages` |
| `REDDIT_LIST_INBOX_NOTIFICATIONS` | List inbox mentions | read | *after, before, count, limit, max_replies, mid, show, sr_detail* |  | `privatemessages` |
| `REDDIT_LIST_MODERATION_LOG` | List moderation log | read | `subreddit`; *action, after, before, count, limit, moderator, show, sr_detail* |  | `modlog` |
| `REDDIT_LIST_MODERATION_QUEUE` | List moderation queue | read | `queue`, `subreddit`; *after, before, content_type, count, limit, show, sr_detail* |  | `read` |
| `REDDIT_LIST_MODMAIL_CONVERSATIONS` | List modmail conversations | read | *after, limit, sort, state, subreddits* |  | `modmail` |
| `REDDIT_LIST_POPULAR_SUBREDDITS` | List popular subreddits | read | *after, before, count, limit, show, sr_detail* |  | `read` |
| `REDDIT_LIST_SUBREDDIT_COMMENTS` | List subreddit comments | read | `subreddit`; *after, before, count, limit, show, sr_detail* |  | `read` |
| `REDDIT_LIST_SUBREDDIT_MODERATORS` | List subreddit moderators | read | `subreddit`; *after, before, count, limit, show, sr_detail, user* |  | `read` |
| `REDDIT_LIST_SUBREDDIT_POST_FLAIRS` | List subreddit post flairs | read | `subreddit` |  | `flair` |
| `REDDIT_LIST_SUBREDDIT_WIKI_PAGES` | List subreddit wiki pages | read | `subreddit` |  | `wikiread` |
| `REDDIT_LIST_SUBSCRIBED_SUBREDDITS` | List subscribed subreddits | read | *after, before, count, limit, show, sr_detail* |  | `mysubreddits` |
| `REDDIT_LIST_USER_ACTIVITY` | List user activity | read | `activity_type`, `username`; *after, before, content_type, context, count, limit, show, sort and 2 more* |  | `history` |
| `REDDIT_MANAGE_SUBREDDIT_SUBSCRIPTION` | Manage subreddit subscription | write | `action`, `subreddit`; *skip_initial_defaults* |  | `subscribe` |
| `REDDIT_MARK_ALL_INBOX_MESSAGES_READ` | Mark all inbox messages read | write | *filter_types* |  | `privatemessages` |
| `REDDIT_POST_REDDIT_COMMENT` | Post a comment | write | `text`, `thing_id` |  | `submit`, `privatemessages` |
| `REDDIT_PREPARE_MEDIA_UPLOAD` | Prepare Reddit media upload | write | `filename`, `mime_type` |  | `submit` |
| `REDDIT_RETRIEVE_POST_COMMENTS` | Retrieve Comments for a Post | read | `article`; *depth, limit, sort* |  | `read` |
| `REDDIT_RETRIEVE_REDDIT_POST` | Retrieve posts from subreddit | read | `subreddit`; *after, before, count, max_results, show, sort, sr_detail, time_filter* |  | `read` |
| `REDDIT_RETRIEVE_SPECIFIC_COMMENT` | Retrieve specific comment or post | read | `id` |  | `read` |
| `REDDIT_SEARCH_ACROSS_SUBREDDITS` | Search across subreddits | read | `search_query`; *after, before, limit, restrict_sr, result_type, sort, time_filter* |  | `read` |
| `REDDIT_SEND_PRIVATE_MESSAGE` | Send private message | write | `recipient`, `subject`, `text`; *from_subreddit* |  | `privatemessages` |
| `REDDIT_TOGGLE_INBOX_REPLIES` | Enable or disable inbox replies | write | `id`, `state` |  | `edit` |
| `REDDIT_VOTE_ON_CONTENT` | Vote on Reddit content | write | `direction`, `thing_id`; *rank* |  | `vote` |
