# Gmail: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

61 tools: 28 read, 33 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `GMAIL_ADD_LABEL_TO_EMAIL` | Modify email labels | write | `message_id`; *add_label_ids, remove_label_ids, user_id* |  | `https://mail.google.com/`, `gmail.modify` |
| `GMAIL_BATCH_DELETE_MESSAGES` | Batch delete Gmail messages | write | `messageIds`; *userId* |  | `https://mail.google.com/` |
| `GMAIL_BATCH_MODIFY_MESSAGES` | Batch modify Gmail messages | write | `messageIds`; *addLabelIds, removeLabelIds, userId* |  | `https://mail.google.com/`, `gmail.modify` |
| `GMAIL_CREATE_EMAIL_DRAFT` | Create email draft | write | *attachment, bcc, body, cc, extra_recipients, include_signature, is_html, recipient_email and 3 more* |  | `https://mail.google.com/`, `gmail.compose`, `gmail.modify`, `gmail.readonly`, `contacts`, `contacts.readonly` |
| `GMAIL_CREATE_FILTER` | Create Gmail filter | write | `action`, `criteria`; *user_id* |  | `gmail.settings.basic` |
| `GMAIL_CREATE_LABEL` | Create label | write | `label_name`; *background_color, label_list_visibility, message_list_visibility, text_color, user_id* |  | `https://mail.google.com/`, `gmail.labels`, `gmail.modify` |
| `GMAIL_DELETE_DRAFT` | Delete Draft | write | `draft_id`; *user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.compose` |
| `GMAIL_DELETE_FILTER` | Delete Gmail filter | write | `filter_id`; *user_id* |  | `gmail.settings.basic` |
| `GMAIL_DELETE_LABEL` | Delete label from account (permanent) | write | `label_id`; *user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.labels` |
| `GMAIL_DELETE_MESSAGE` | Delete message | write | `message_id`; *user_id* |  | `https://mail.google.com/` |
| `GMAIL_DELETE_THREAD` | Delete thread | write | `id`; *user_id* |  | `https://mail.google.com/` |
| `GMAIL_FETCH_EMAILS` | Fetch emails | read | *ids_only, include_payload, include_spam_trash, label_ids, max_results, page_token, query, user_id and 1 more* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_FETCH_MESSAGE_BY_MESSAGE_ID` | Fetch message by message ID | read | `message_id`; *format, user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_FETCH_MESSAGE_BY_THREAD_ID` | Fetch Message by Thread ID | read | `thread_id`; *page_token, user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_FORWARD_MESSAGE` | Forward email message | write | `message_id`, `recipients`; *additional_text, bcc, cc, include_signature, user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly`, `gmail.compose`, `gmail.send` |
| `GMAIL_GET_ATTACHMENT` | Get Gmail attachment | read | `attachment_id`, `file_name`, `message_id`; *user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_GET_AUTO_FORWARDING` | Get Auto-Forwarding Settings | read | *user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_GET_CONTACTS` | Get contacts | read | *include_other_contacts, page_token, person_fields, resource_name* |  | `contacts`, `contacts.readonly`, `contacts.other.readonly` |
| `GMAIL_GET_DRAFT` | Get Draft | read | `draft_id`; *format, user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.compose`, `gmail.readonly` |
| `GMAIL_GET_FILTER` | Get Gmail filter | read | `id`; *user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_GET_LABEL` | Get label details | read | `id`; *user_id* |  | `https://mail.google.com/`, `gmail.labels`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_GET_LANGUAGE_SETTINGS` | Get Language Settings | read | *user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_GET_PEOPLE` | Get People | read | *other_contacts, page_size, page_token, person_fields, resource_name, sources, sync_token* |  | `contacts`, `contacts.other.readonly`, `contacts.readonly`, `userinfo.profile` |
| `GMAIL_GET_PROFILE` | Get Profile | read | *user_id* |  | `https://mail.google.com/`, `gmail.compose`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_GET_VACATION_SETTINGS` | Get Vacation Settings | read | *user_id* |  | `gmail.settings.basic`, `gmail.modify`, `gmail.readonly`, `https://mail.google.com/` |
| `GMAIL_IMPORT_MESSAGE` | Import message | write | `raw`; *deleted, internal_date_source, never_mark_spam, process_for_calendar, user_id* |  | `https://mail.google.com/`, `gmail.insert`, `gmail.modify` |
| `GMAIL_INSERT_MESSAGE` | Insert message into mailbox | write | `raw`; *deleted, internalDateSource, user_id* |  | `https://mail.google.com/`, `gmail.insert`, `gmail.modify` |
| `GMAIL_LIST_CSE_IDENTITIES` | List CSE identities | read | *page_size, page_token, user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly`, `gmail.settings.sharing` |
| `GMAIL_LIST_CSE_KEYPAIRS` | List CSE key pairs | read | *page_size, page_token, user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly`, `gmail.settings.sharing` |
| `GMAIL_LIST_DRAFTS` | List Drafts | read | *max_results, page_token, user_id, verbose* |  | `https://mail.google.com/`, `gmail.compose`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_LIST_FILTERS` | List Gmail filters | read | *user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly`, `gmail.settings.basic` |
| `GMAIL_LIST_FORWARDING_ADDRESSES` | List forwarding addresses | read | *user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly`, `gmail.settings.basic` |
| `GMAIL_LIST_HISTORY` | List Gmail history | read | `start_history_id`; *history_types, label_id, max_results, page_token, user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_LIST_LABELS` | List Gmail labels | read | *include_details, user_id* |  | `https://mail.google.com/`, `gmail.labels`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_LIST_SEND_AS` | List send-as aliases | read | *user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_LIST_SMIME_INFO` | List S/MIME configs | read | `send_as_email`; *user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly`, `gmail.settings.basic`, `gmail.settings.sharing` |
| `GMAIL_LIST_THREADS` | List threads | read | *max_results, page_token, query, user_id, verbose* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_MODIFY_THREAD_LABELS` | Modify thread labels | write | `thread_id`; *add_label_ids, remove_label_ids, user_id* |  | `https://mail.google.com/`, `gmail.modify` |
| `GMAIL_MOVE_THREAD_TO_TRASH` | Trash thread | write | `thread_id`; *user_id* |  | `https://mail.google.com/`, `gmail.modify` |
| `GMAIL_MOVE_TO_TRASH` | Move to Trash | write | `message_id`; *user_id* |  | `https://mail.google.com/`, `gmail.modify` |
| `GMAIL_PATCH_LABEL` | Patch Label | write | `id`, `userId`; *color, labelListVisibility, messageListVisibility, name* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.labels` |
| `GMAIL_PATCH_SEND_AS` | Patch send-as alias | write | `send_as_email`; *display_name, is_default, reply_to_address, signature, smtp_msa, treat_as_alias, user_id* |  | `gmail.settings.basic`, `gmail.settings.sharing` |
| `GMAIL_REPLY_TO_THREAD` | Reply to email thread | write | `thread_id`; *attachment, bcc, cc, extra_recipients, include_signature, is_html, message_body, recipient_email and 1 more* |  | `https://mail.google.com/`, `gmail.compose`, `gmail.modify`, `gmail.send`, `gmail.readonly`, `contacts`, `contacts.other.readonly`, `contacts.readonly` |
| `GMAIL_SEARCH_PEOPLE` | Search People | read | `query`; *other_contacts, page_size, person_fields* |  | `contacts`, `contacts.readonly`, `contacts.other.readonly` |
| `GMAIL_SEND_DRAFT` | Send Draft | write | `draft_id`; *user_id* |  | `https://mail.google.com/`, `gmail.compose`, `gmail.modify` |
| `GMAIL_SEND_EMAIL` | Send Email | write | *attachment, bcc, body, cc, extra_recipients, from_email, include_signature, is_html and 3 more* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.compose`, `gmail.send` |
| `GMAIL_SETTINGS_GET_IMAP` | Get IMAP Settings | read | *user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_SETTINGS_GET_POP` | Get POP settings | read | *user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_SETTINGS_SEND_AS_GET` | Get send-as alias | read | `send_as_email`; *user_id* |  | `gmail.settings.basic`, `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_STOP_WATCH` | Stop watch notifications | write | *user_id* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.readonly` |
| `GMAIL_UNTRASH_MESSAGE` | Untrash Message | write | `message_id`; *user_id* |  | `https://mail.google.com/`, `gmail.modify` |
| `GMAIL_UNTRASH_THREAD` | Untrash thread | write | `thread_id`; *user_id* |  | `https://mail.google.com/`, `gmail.modify` |
| `GMAIL_UPDATE_CONTACT` | Update Google Contact | write | `person`, `resource_name`, `update_person_fields`; *person_fields, sources* |  | `contacts` |
| `GMAIL_UPDATE_DRAFT` | Update draft | write | `draft_id`; *attachment, bcc, body, cc, extra_recipients, include_signature, is_html, recipient_email and 3 more* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.compose` |
| `GMAIL_UPDATE_IMAP_SETTINGS` | Update IMAP settings | write | *autoExpunge, enabled, expungeBehavior, maxFolderSize, user_id* |  | `gmail.settings.basic` |
| `GMAIL_UPDATE_LABEL` | Update Label | write | `id`; *color, labelListVisibility, messageListVisibility, name, userId* |  | `https://mail.google.com/`, `gmail.modify`, `gmail.labels` |
| `GMAIL_UPDATE_LANGUAGE_SETTINGS` | Update Language Settings | write | `display_language`; *user_id* |  | `https://mail.google.com/`, `gmail.settings.basic` |
| `GMAIL_UPDATE_POP_SETTINGS` | Update POP settings | write | *access_window, disposition, user_id* |  | `https://mail.google.com/`, `gmail.settings.basic` |
| `GMAIL_UPDATE_SEND_AS` | Update send-as alias | write | `send_as_email`; *display_name, is_default, reply_to_address, signature, smtp_msa, treat_as_alias, user_id* |  | `gmail.settings.basic`, `gmail.settings.sharing` |
| `GMAIL_UPDATE_VACATION_SETTINGS` | Update Vacation Settings | write | *enableAutoReply, endTime, responseBodyHtml, responseBodyPlainText, responseSubject, restrictToContacts, restrictToDomain, startTime and 1 more* |  | `gmail.settings.basic` |
| `GMAIL_WHO_AM_I` | Who Am I | read |  |  |  |
