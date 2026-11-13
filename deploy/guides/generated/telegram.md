# Telegram: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner pastes a key from the vendor's settings). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

18 tools: 7 read, 11 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `TELEGRAM_ANSWER_CALLBACK_QUERY` | Answer Callback Query | write | `callback_query_id`; *cache_time, show_alert, text, url* |  |  |
| `TELEGRAM_CREATE_CHAT_INVITE_LINK` | Export Chat Invite Link | write | `chat_id` |  |  |
| `TELEGRAM_DELETE_MESSAGE` | Delete Message | write | `chat_id`, `message_id` |  |  |
| `TELEGRAM_EDIT_MESSAGE` | Edit Message | write | `text`; *chat_id, disable_web_page_preview, inline_message_id, message_id, parse_mode, reply_markup* |  |  |
| `TELEGRAM_FORWARD_MESSAGE` | Forward Message | write | `chat_id`, `from_chat_id`, `message_id`; *disable_notification* |  |  |
| `TELEGRAM_GET_CHAT` | Get Chat Info | read | `chat_id` |  |  |
| `TELEGRAM_GET_CHAT_ADMINISTRATORS` | Get Chat Administrators | read | `chat_id` |  |  |
| `TELEGRAM_GET_CHAT_HISTORY` | Get Chat History | read | `chat_id`; *limit, message_id, offset* |  |  |
| `TELEGRAM_GET_CHAT_MEMBER` | Get Chat Member | read | `chat_id`, `user_id` |  |  |
| `TELEGRAM_GET_CHAT_MEMBERS_COUNT` | Get Chat Members Count | read | `chat_id` |  |  |
| `TELEGRAM_GET_ME` | Get Bot Info | read |  |  |  |
| `TELEGRAM_GET_UPDATES` | Get Updates | read | *allowed_updates, limit, offset, timeout* |  |  |
| `TELEGRAM_SEND_DOCUMENT` | Send Document | write | `chat_id`, `document`; *caption, disable_content_type_detection, disable_notification, parse_mode, reply_markup, reply_to_message_id, thumbnail* |  |  |
| `TELEGRAM_SEND_LOCATION` | Send Location | write | `chat_id`, `latitude`, `longitude`; *disable_notification, heading, horizontal_accuracy, live_period, proximity_alert_radius, reply_markup, reply_to_message_id* |  |  |
| `TELEGRAM_SEND_MESSAGE` | Send Message | write | `chat_id`, `text`; *disable_notification, disable_web_page_preview, parse_mode, reply_markup, reply_to_message_id* |  |  |
| `TELEGRAM_SEND_PHOTO` | Send Photo | write | `chat_id`, `photo`; *caption, disable_notification, parse_mode, reply_markup, reply_to_message_id* |  |  |
| `TELEGRAM_SEND_POLL` | Send Poll | write | `chat_id`, `options`, `question`; *allows_multiple_answers, close_date, correct_option_id, disable_notification, explanation, explanation_parse_mode, is_anonymous, is_closed and 4 more* |  |  |
| `TELEGRAM_SET_MY_COMMANDS` | Set Bot Commands | write | `commands`; *language_code, scope* |  |  |
