# Trello: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

329 tools: 140 read, 189 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `TRELLO_ADD_BOARDS` | Add board | write | *closed, desc, idBoardSource, idOrganization, keepFromSource, labelNames__blue, labelNames__green, labelNames__orange and 23 more* |  |  |
| `TRELLO_ADD_CARDS` | Add card | write | `idList`; *closed, desc, due, fileSource, idAttachmentCover, idBoard, idCardSource, idLabels and 6 more* |  |  |
| `TRELLO_ADD_CARDS_ACTIONS_COMMENTS_BY_ID_CARD` | Add comment to card | write | `idCard`, `text` |  | `write` |
| `TRELLO_ADD_CARDS_ATTACHMENTS_BY_ID_CARD` | Add attachment to card | write | `idCard`; *file, mimeType, name, url* |  |  |
| `TRELLO_ADD_CARDS_CHECKLISTS_BY_ID_CARD` | Add checklist to card via id | write | `idCard`; *idChecklistSource, name, value* |  |  |
| `TRELLO_ADD_CARDS_ID_LABELS_BY_ID_CARD` | Add label to card | write | `idCard`, `value` |  |  |
| `TRELLO_ADD_CARD_VOTE` | Add member vote to card | write | `idCard`, `value` |  |  |
| `TRELLO_ADD_CHECKLISTS` | Add checklist to card | write | *idBoard, idCard, idChecklistSource, name, pos* |  |  |
| `TRELLO_ADD_CHECKLISTS_CHECK_ITEMS_BY_ID_CHECKLIST` | Add check item to checklist | write | `idChecklist`; *checked, name, pos* |  |  |
| `TRELLO_ADD_LABELS` | Create label on board | write | *color, idBoard, name* |  |  |
| `TRELLO_ADD_LISTS` | Add new list to board | write | *closed, idBoard, idListSource, name, pos, subscribed* |  |  |
| `TRELLO_ADD_MEMBER_TO_CARD` | Add card member by id | write | `idCard`; *value* |  |  |
| `TRELLO_ARCHIVE_ALL_LIST_CARDS` | Archive all cards in list | write | `idList` |  |  |
| `TRELLO_CONVERT_CHECKLIST_ITEM_TO_CARD` | Convert checklist item to card | write | `idCard`, `idCheckItem`, `idChecklist` |  |  |
| `TRELLO_CREATE_ACTION_REACTION` | Add reaction to action | write | `idAction`; *native, shortName, skinVariation, unified* |  | `read`, `write` |
| `TRELLO_CREATE_BOARD_LABEL` | Add a label to a board | write | `idBoard`; *color, name* |  |  |
| `TRELLO_CREATE_CARD_LABEL` | Add labels to card | write | `idCard`; *color, name, value* |  |  |
| `TRELLO_CREATE_CARD_STICKER` | Add sticker to card | write | `idCard`; *image, left, rotate, top, zIndex* |  |  |
| `TRELLO_CREATE_CUSTOM_FIELD` | Create custom field | write | `idModel`, `name`, `pos`, `type`; *displayCardFront, options* |  | `write` |
| `TRELLO_CREATE_MEMBER_BOARD_BACKGROUND` | Add member board background | write | `idMember`; *brightness, file, tile* |  |  |
| `TRELLO_CREATE_MEMBER_BOARD_STAR` | Add board star to member | write | `idBoard`, `idMember`, `pos` |  |  |
| `TRELLO_CREATE_MEMBER_CUSTOM_BOARD_BACKGROUND` | Add custom board background for a member | write | `idMember`; *brightness, file, tile, url* |  |  |
| `TRELLO_CREATE_MEMBER_CUSTOM_EMOJI` | Add custom emoji for a member | write | `idMember`; *file, name* |  |  |
| `TRELLO_CREATE_MEMBER_CUSTOM_STICKER` | Add custom sticker for member | write | `file`, `idMember` |  |  |
| `TRELLO_CREATE_MEMBER_SAVED_SEARCH` | Add saved search for member | write | `idMember`, `name`, `pos`, `query` |  |  |
| `TRELLO_CREATE_ORGANIZATION` | Create organization | write | *desc, displayName, name, prefs__associatedDomain, prefs__boardVisibilityRestrict__org, prefs__boardVisibilityRestrict__private, prefs__boardVisibilityRestrict__public, prefs__externalMembersDisabled and 4 more* |  |  |
| `TRELLO_CREATE_SESSION` | Create new session | write | *idBoard, status* |  |  |
| `TRELLO_CREATE_WEBHOOK` | Create webhook | write | `callbackURL`, `idModel`; *active, description* |  | `read`, `write` |
| `TRELLO_DELETE_ACTION` | Delete action by id | write | `idAction` |  |  |
| `TRELLO_DELETE_ACTION_REACTION` | Delete action reaction | write | `id`, `idAction` |  | `read`, `write` |
| `TRELLO_DELETE_BOARD` | Delete board by id | write | `id` |  | `read`, `write` |
| `TRELLO_DELETE_CARDS_BY_ID_CARD` | Delete card by ID | write | `idCard` |  |  |
| `TRELLO_DELETE_CARDS_MEMBERS_VOTED_BY_ID_CARD_BY_ID_MEMBER` | Delete member vote from card | write | `idCard`, `idMember` |  |  |
| `TRELLO_DELETE_CARDS_STICKERS_BY_ID_CARD_BY_ID_STICKER` | Delete card sticker | write | `idCard`, `idSticker` |  |  |
| `TRELLO_DELETE_CARD_ATTACHMENT` | Delete card attachment | write | `idAttachment`, `idCard` |  |  |
| `TRELLO_DELETE_CARD_CHECKLIST` | Delete card checklist | write | `idCard`, `idChecklist` |  |  |
| `TRELLO_DELETE_CARD_CHECK_ITEM` | Delete card check item | write | `idCard`, `idCheckItem` |  | `write` |
| `TRELLO_DELETE_CARD_COMMENT` | Delete card action comment | write | `idAction`, `idCard` |  |  |
| `TRELLO_DELETE_CHECKLISTS_BY_ID_CHECKLIST` | Delete checklist by id | write | `idChecklist` |  |  |
| `TRELLO_DELETE_CHECKLIST_ITEM` | Delete checklist item | write | `idCheckItem`, `idChecklist` |  |  |
| `TRELLO_DELETE_LABELS_BY_ID_LABEL` | Delete label by id | write | `idLabel` |  |  |
| `TRELLO_DELETE_MEMBER_BOARD_BACKGROUND` | Delete member board background | write | `idBoardBackground`, `idMember` |  |  |
| `TRELLO_DELETE_MEMBER_CUSTOM_STICKER` | Delete member custom sticker | write | `idCustomSticker`, `idMember` |  |  |
| `TRELLO_DELETE_MEMBER_SAVED_SEARCH` | Delete member saved search | write | `idMember`, `idSavedSearch` |  |  |
| `TRELLO_DELETE_ORGANIZATIONS_BY_ID_ORG` | Delete organization by id | write | `idOrg` |  |  |
| `TRELLO_DELETE_ORGANIZATIONS_LOGO_BY_ID_ORG` | Delete organization logo by id | write | `idOrg` |  |  |
| `TRELLO_DELETE_ORGANIZATIONS_MEMBERS_BY_ID_ORG_BY_ID_MEMBER` | Delete organization member by IDs | write | `idMember`, `idOrg` |  |  |
| `TRELLO_DELETE_ORG_ASSOCIATED_DOMAIN` | Delete org associated domain | write | `idOrg` |  |  |
| `TRELLO_DELETE_ORG_INVITE_RESTRICT` | Delete org invite restriction | write | `idOrg`, `value` |  |  |
| `TRELLO_DELETE_TOKENS_BY_TOKEN` | Delete token | write | *token_to_delete* |  |  |
| `TRELLO_DELETE_TOKENS_WEBHOOKS_BY_TOKEN_BY_ID_WEBHOOK` | Delete token webhook by id | write | `idWebhook` |  |  |
| `TRELLO_DELETE_WEBHOOKS_BY_ID_WEBHOOK` | Delete webhook by id | write | `idWebhook` |  |  |
| `TRELLO_DISABLE_BOARD_POWER_UP` | Delete Power-Up from board | write | `idBoard`, `powerUp` |  |  |
| `TRELLO_DISMISS_MEMBER_MESSAGE` | Dismiss member message | write | `idMember`, `value` |  |  |
| `TRELLO_ENABLE_BOARD_POWER_UP` | Add board power-up | write | `idBoard`; *value* |  |  |
| `TRELLO_GENERATE_BOARD_CALENDAR_KEY` | Generate board calendar key | write | `idBoard` |  |  |
| `TRELLO_GENERATE_BOARD_EMAIL_KEY` | Generate email key for board | write | `idBoard` |  |  |
| `TRELLO_GET_ACTIONS_BOARD_BY_ID_ACTION` | Get board by action id | read | `idAction`; *fields* |  |  |
| `TRELLO_GET_ACTIONS_BOARD_BY_ID_ACTION_BY_FIELD` | Get action's board field | read | `field`, `idAction` |  |  |
| `TRELLO_GET_ACTIONS_BY_ID_ACTION` | Get action by ID | read | `idAction`; *display, entities, fields, member, memberCreator, memberCreator_fields, member_fields* |  |  |
| `TRELLO_GET_ACTIONS_BY_ID_ACTION_BY_FIELD` | Get action field by id | read | `field`, `idAction` |  |  |
| `TRELLO_GET_ACTIONS_CARD_BY_ID_ACTION` | Get card by action ID | read | `idAction`; *fields* |  |  |
| `TRELLO_GET_ACTIONS_CARD_BY_ID_ACTION_BY_FIELD` | Get action's card field | read | `field`, `idAction` |  |  |
| `TRELLO_GET_ACTIONS_DISPLAY_BY_ID_ACTION` | Get action display by ID | read | `idAction` |  |  |
| `TRELLO_GET_ACTIONS_ENTITIES_BY_ID_ACTION` | Get action entities by id | read | `idAction` |  |  |
| `TRELLO_GET_ACTIONS_LIST_BY_ID_ACTION` | Get an action's list | read | `idAction`; *fields* |  |  |
| `TRELLO_GET_ACTIONS_LIST_BY_ID_ACTION_BY_FIELD` | Get field of action's list | read | `field`, `idAction` |  |  |
| `TRELLO_GET_ACTIONS_MEMBER_BY_ID_ACTION` | Get action's member by id | read | `idAction`; *fields* |  |  |
| `TRELLO_GET_ACTIONS_MEMBER_BY_ID_ACTION_BY_FIELD` | Get Action Subject Member Field | read | `field`, `idAction` |  |  |
| `TRELLO_GET_ACTIONS_MEMBER_CREATOR_BY_ID_ACTION` | Fetch action member creator | read | `idAction`; *fields* |  |  |
| `TRELLO_GET_ACTIONS_MEMBER_CREATOR_BY_ID_ACTION_BY_FIELD` | Get action member creator field | read | `field`, `idAction` |  |  |
| `TRELLO_GET_ACTIONS_ORGANIZATION_BY_ID_ACTION` | Fetch organization action by id | read | `idAction`; *fields* |  |  |
| `TRELLO_GET_ACTIONS_ORGANIZATION_BY_ID_ACTION_BY_FIELD` | Get action's organization field | read | `field`, `idAction` |  |  |
| `TRELLO_GET_ACTIONS_REACTIONS_BY_ID` | Get action's reaction by ID | read | `id`, `idAction`; *emoji, member* |  | `read` |
| `TRELLO_GET_ACTIONS_REACTIONS_BY_ID_ACTION` | List action reactions | read | `idAction`; *emoji, member* |  | `read` |
| `TRELLO_GET_BATCH` | Get batch | read | `urls` |  |  |
| `TRELLO_GET_BOARDS_ACTIONS_BY_ID_BOARD` | Get board actions by id | read | `idBoard`; *before, display, entities, fields, filter, format, idModels, limit and 6 more* |  |  |
| `TRELLO_GET_BOARDS_BOARD_PLUGINS_BY_ID` | Get board Power-Ups | read | `id` |  | `read` |
| `TRELLO_GET_BOARDS_BOARD_STARS_BY_ID_BOARD` | Get board stars by board ID | read | `idBoard`; *filter* |  |  |
| `TRELLO_GET_BOARDS_BY_ID_BOARD` | Get board by id | read | `idBoard`; *action_fields, action_member, action_memberCreator, action_memberCreator_fields, action_member_fields, actions, actions_display, actions_entities and 29 more* |  |  |
| `TRELLO_GET_BOARDS_BY_ID_BOARD_BY_FIELD` | Get board field | read | `field`, `idBoard` |  |  |
| `TRELLO_GET_BOARDS_CARDS_BY_ID_BOARD` | Get cards by board ID | read | `idBoard`; *actions, attachment_fields, attachments, before, checkItemStates, checklists, fields, filter and 5 more* |  |  |
| `TRELLO_GET_BOARDS_CARDS_BY_ID_BOARD_BY_FILTER` | Get cards by filter from board | read | `filter`, `idBoard` |  |  |
| `TRELLO_GET_BOARDS_CARDS_BY_ID_BOARD_BY_ID_CARD` | Get card from board | read | `idBoard`, `idCard`; *action_fields, action_memberCreator_fields, actions, actions_display, actions_entities, actions_limit, attachment_fields, attachments and 8 more* |  |  |
| `TRELLO_GET_BOARDS_CHECKLISTS_BY_ID_BOARD` | Get board checklists | read | `idBoard`; *card_fields, cards, checkItem_fields, checkItems, fields, filter* |  |  |
| `TRELLO_GET_BOARDS_CUSTOM_FIELDS` | Get board custom fields | write | `id` |  | `read` |
| `TRELLO_GET_BOARDS_DELTAS_BY_ID_BOARD` | Retrieve board deltas | read | `idBoard`, `ixLastUpdate`, `tags` |  |  |
| `TRELLO_GET_BOARDS_LABELS_BY_ID_BOARD` | Get board labels by ID | read | `idBoard`; *fields, limit* |  |  |
| `TRELLO_GET_BOARDS_LABELS_BY_ID_BOARD_BY_ID_LABEL` | Get a board label | read | `idBoard`, `idLabel`; *fields* |  |  |
| `TRELLO_GET_BOARDS_LISTS_BY_ID_BOARD` | Get board's lists | read | `idBoard`; *card_fields, cards, fields, filter* |  |  |
| `TRELLO_GET_BOARDS_MEMBERSHIPS_BY_ID_BOARD` | List board memberships | read | `idBoard`; *filter, member, member_fields* |  |  |
| `TRELLO_GET_BOARDS_MEMBERSHIPS_BY_ID_BOARD_BY_ID_MEMBERSHIP` | Get board membership | read | `idBoard`, `idMembership`; *member, member_fields* |  |  |
| `TRELLO_GET_BOARDS_MEMBERS_BY_ID_BOARD` | Get board members | read | `idBoard`; *activity, fields, filter* |  |  |
| `TRELLO_GET_BOARDS_MEMBERS_CARDS_BY_ID_BOARD_BY_ID_MEMBER` | Get member cards from board | read | `idBoard`, `idMember`; *actions, attachment_fields, attachments, board, board_fields, checkItemStates, checklists, fields and 5 more* |  |  |
| `TRELLO_GET_BOARDS_MEMBERS_INVITED_BY_ID_BOARD` | Get invited board members | read | `idBoard`; *fields* |  |  |
| `TRELLO_GET_BOARDS_MEMBERS_INVITED_BY_ID_BOARD_BY_FIELD` | Retrieve invited board member field | read | `field`, `idBoard` |  |  |
| `TRELLO_GET_BOARDS_MY_PREFS_BY_ID_BOARD` | Get my board preferences | read | `idBoard` |  |  |
| `TRELLO_GET_BOARDS_ORGANIZATION_BY_ID_BOARD` | Get board organization | read | `idBoard`; *fields* |  |  |
| `TRELLO_GET_BOARDS_ORGANIZATION_BY_ID_BOARD_BY_FIELD` | Get board organization field | read | `field`, `idBoard` |  |  |
| `TRELLO_GET_BOARDS_PLUGINS_BY_ID_BOARD` | Get board's Power-Ups | read | `id`; *filter* |  | `read` |
| `TRELLO_GET_CARDS_ACTIONS_BY_ID_CARD` | Get card actions by id | read | `idCard`; *before, display, entities, fields, filter, format, idModels, limit and 6 more* |  |  |
| `TRELLO_GET_CARDS_ATTACHMENTS_BY_ID_CARD` | Get card attachments | read | `idCard`; *fields, filter* |  |  |
| `TRELLO_GET_CARDS_ATTACHMENTS_BY_ID_CARD_BY_ID_ATTACHMENT` | Get specific card attachment | read | `idAttachment`, `idCard`; *fields* |  |  |
| `TRELLO_GET_CARDS_BOARD_BY_ID_CARD` | Get board by card id | read | `idCard`; *fields* |  |  |
| `TRELLO_GET_CARDS_BOARD_BY_ID_CARD_BY_FIELD` | Get board field by card ID | read | `field`, `idCard` |  |  |
| `TRELLO_GET_CARDS_BY_ID_CARD` | Get card by ID | read | `idCard`; *action_fields, action_memberCreator_fields, actions, actions_display, actions_entities, actions_limit, attachment_fields, attachments and 15 more* |  |  |
| `TRELLO_GET_CARDS_BY_ID_CARD_BY_FIELD` | Get card field by id | read | `field`, `idCard` |  |  |
| `TRELLO_GET_CARDS_CHECKLISTS_BY_ID_CARD` | Get card checklists by ID | read | `idCard`; *card_fields, cards, checkItem_fields, checkItems, fields, filter* |  |  |
| `TRELLO_GET_CARDS_CHECK_ITEM` | Get card check item | read | `id`, `idCheckItem`; *fields* |  | `read` |
| `TRELLO_GET_CARDS_CHECK_ITEM_STATES_BY_ID_CARD` | Get card check item states | read | `idCard`; *fields* |  |  |
| `TRELLO_GET_CARDS_CUSTOM_FIELD_ITEMS_BY_ID_CARD` | Get card custom field items | write | `idCard` |  |  |
| `TRELLO_GET_CARDS_LIST_BY_ID_CARD` | Get list by card ID | read | `idCard`; *fields* |  |  |
| `TRELLO_GET_CARDS_LIST_BY_ID_CARD_BY_FIELD` | Get card list field | read | `field`, `idCard` |  |  |
| `TRELLO_GET_CARDS_MEMBERS_BY_ID_CARD` | Get card members | read | `idCard`; *fields* |  |  |
| `TRELLO_GET_CARDS_MEMBERS_VOTED_BY_ID_CARD` | Get card members voted | read | `idCard`; *fields* |  |  |
| `TRELLO_GET_CARDS_PLUGIN_DATA_BY_ID` | Get card plugin data | read | `id` |  | `read` |
| `TRELLO_GET_CARDS_STICKERS_BY_ID_CARD` | Get card stickers by ID card | read | `idCard`; *fields* |  |  |
| `TRELLO_GET_CARDS_STICKERS_BY_ID_CARD_BY_ID_STICKER` | Get card sticker by id | read | `idCard`, `idSticker`; *fields* |  |  |
| `TRELLO_GET_CHECKLISTS_BOARD_BY_ID_CHECKLIST` | Get board for a checklist | read | `idChecklist`; *fields* |  |  |
| `TRELLO_GET_CHECKLISTS_BOARD_BY_ID_CHECKLIST_BY_FIELD` | Get checklist's board field | read | `field`, `idChecklist` |  |  |
| `TRELLO_GET_CHECKLISTS_BY_ID_CHECKLIST` | Get checklist by ID | read | `idChecklist`; *card_fields, cards, checkItem_fields, checkItems, fields* |  |  |
| `TRELLO_GET_CHECKLISTS_BY_ID_CHECKLIST_BY_FIELD` | Get checklist field | read | `field`, `idChecklist` |  |  |
| `TRELLO_GET_CHECKLISTS_CARDS_BY_ID_CHECKLIST` | Get checklist cards | read | `idChecklist`; *actions, attachment_fields, attachments, before, checkItemStates, checklists, fields, filter and 5 more* |  |  |
| `TRELLO_GET_CHECKLISTS_CHECK_ITEMS_BY_ID_CHECKLIST` | Get checklist items by ID | read | `idChecklist`; *fields, filter* |  |  |
| `TRELLO_GET_CHECK_ITEM_BY_ID` | Get check item by id | read | `idCheckItem`, `idChecklist`; *fields* |  |  |
| `TRELLO_GET_CUSTOM_FIELD` | Get custom field | write | `id` |  | `read` |
| `TRELLO_GET_LABELS_BOARD_BY_ID_LABEL` | Get board by label ID | read | `idLabel`; *fields* |  |  |
| `TRELLO_GET_LABELS_BOARD_BY_ID_LABEL_BY_FIELD` | Retrieve board field by label id | read | `field`, `idLabel` |  |  |
| `TRELLO_GET_LABELS_BY_ID_LABEL` | Get label by id | read | `idLabel`; *fields* |  |  |
| `TRELLO_GET_LISTS_ACTIONS_BY_ID_LIST` | Get list actions by ID | read | `idList`; *before, display, entities, fields, filter, format, idModels, limit and 6 more* |  |  |
| `TRELLO_GET_LISTS_BOARD_BY_ID_LIST` | Get board by list ID | read | `idList`; *fields* |  |  |
| `TRELLO_GET_LISTS_BOARD_BY_ID_LIST_BY_FIELD` | Get board field by list ID | read | `field`, `idList` |  |  |
| `TRELLO_GET_LISTS_BY_ID_LIST` | Get list by ID | read | `idList`; *board, board_fields, card_fields, cards, fields* |  |  |
| `TRELLO_GET_LISTS_BY_ID_LIST_BY_FIELD` | Get list field value | read | `field`, `idList` |  |  |
| `TRELLO_GET_LISTS_CARDS_BY_ID_LIST` | Get cards in a list | read | `idList`; *actions, attachment_fields, attachments, before, checkItemStates, checklists, fields, filter and 5 more* |  |  |
| `TRELLO_GET_MEMBERS_ACTIONS_BY_ID_MEMBER` | Get member actions by ID | read | `idMember`; *before, display, entities, fields, filter, format, idModels, limit and 6 more* |  |  |
| `TRELLO_GET_MEMBERS_BOARDS_BY_ID_MEMBER` | Get member boards by id | read | *action_fields, actions, actions_entities, actions_format, actions_limit, actions_since, fields, filter and 5 more* |  |  |
| `TRELLO_GET_MEMBERS_BOARDS_BY_ID_MEMBER_BY_FILTER` | Get member boards with filter | read | `idMember`; *fields, filter* |  |  |
| `TRELLO_GET_MEMBERS_BOARDS_INVITED_BY_ID_MEMBER` | Get member's invited boards | read | *fields, idMember* |  |  |
| `TRELLO_GET_MEMBERS_BOARDS_INVITED_BY_ID_MEMBER_BY_FIELD` | Get member's invited board field | read | `field`, `idMember` |  |  |
| `TRELLO_GET_MEMBERS_BOARD_BACKGROUNDS_BY_ID_MEMBER` | Get member board backgrounds | read | `idMember`; *filter* |  |  |
| `TRELLO_GET_MEMBERS_BOARD_STARS_BY_ID_MEMBER` | Get member board stars | read | `idMember` |  |  |
| `TRELLO_GET_MEMBERS_BOARD_STARS_BY_ID_MEMBER_BY_ID_BOARD_STAR` | Get member board star | read | `idBoardStar`, `idMember` |  |  |
| `TRELLO_GET_MEMBERS_BY_ID_MEMBER` | Get member by ID | read | `idMember`; *action_before, action_fields, action_since, actions, actions_display, actions_entities, actions_limit, boardBackgrounds and 44 more* |  |  |
| `TRELLO_GET_MEMBERS_BY_ID_MEMBER_BY_FIELD` | Get member field by ID | read | `field`, `idMember` |  |  |
| `TRELLO_GET_MEMBERS_CARDS_BY_ID_MEMBER` | Get cards by member ID | read | `idMember`; *actions, attachment_fields, attachments, before, checkItemStates, checklists, fields, filter and 5 more* |  |  |
| `TRELLO_GET_MEMBERS_CUSTOM_BOARD_BACKGROUNDS_BY_ID_MEMBER` | Retrieve member custom board backgrounds | write | `idMember`; *filter* |  |  |
| `TRELLO_GET_MEMBERS_CUSTOM_EMOJI_BY_ID_MEMBER` | Get member custom emoji | write | `idMember`; *filter* |  |  |
| `TRELLO_GET_MEMBERS_CUSTOM_STICKERS_BY_ID_MEMBER` | Get member custom stickers | write | `idMember`; *filter* |  |  |
| `TRELLO_GET_MEMBERS_CUSTOM_STICKERS_BY_ID_MEMBER_BY_ID_STICKE` | Get member custom sticker | write | `idCustomSticker`, `idMember`; *fields* |  |  |
| `TRELLO_GET_MEMBERS_DELTAS_BY_ID_MEMBER` | Get member deltas by ID | read | `idMember`, `ixLastUpdate`, `tags` |  |  |
| `TRELLO_GET_MEMBERS_ME` | Get authenticated member | read | *actions, boardBackgrounds, boardStars, boards, boardsInvited, cards, customBoardBackgrounds, customEmoji and 2 more* |  | `read` |
| `TRELLO_GET_MEMBERS_NOTIFICATIONS_BY_ID_MEMBER` | Get member notifications by id | read | `idMember`; *before, display, entities, fields, filter, limit, memberCreator, memberCreator_fields and 3 more* |  |  |
| `TRELLO_GET_MEMBERS_ORGANIZATIONS_BY_ID_MEMBER` | Get a specified member's organizations | read | `idMember`; *fields, filter, paid_account* |  |  |
| `TRELLO_GET_MEMBERS_ORGANIZATIONS_INVITED_BY_ID_MEMBER` | Retrieve member s invited organizations | read | `idMember`; *fields* |  |  |
| `TRELLO_GET_MEMBERS_SAVED_SEARCHES_BY_ID_MEMBER` | Get member saved searches | read | `idMember` |  |  |
| `TRELLO_GET_MEMBERS_SAVED_SEARCHES_BY_ID_MEMBER_BY_ID_SAVED_S` | Get Member Saved Search | read | `idMember`, `idSavedSearch` |  |  |
| `TRELLO_GET_MEMBERS_TOKENS_BY_ID_MEMBER` | Retrieve member tokens | read | `idMember`; *filter* |  |  |
| `TRELLO_GET_MEMBER_BOARD_BACKGROUND` | Get Member Board Background | read | `idBoardBackground`, `idMember`; *fields* |  |  |
| `TRELLO_GET_MEMBER_CUSTOM_BG` | Get member custom board background | write | `idBoardBackground`, `idMember`; *fields* |  |  |
| `TRELLO_GET_MEMBER_CUSTOM_EMOJI` | Get member custom emoji | write | `idCustomEmoji`, `idMember`; *fields* |  |  |
| `TRELLO_GET_MEMBER_INVITED_ORG` | Get field of member's invited organization | read | `field`, `idMember` |  |  |
| `TRELLO_GET_MEMBER_ORG_CARDS` | Get organization member cards | read | `idMember`, `idOrg`; *actions, attachment_fields, attachments, board, board_fields, checkItemStates, checklists, fields and 5 more* |  |  |
| `TRELLO_GET_NOTIFICATIONS_BOARD_BY_ID_NOTIFICATION` | Get notification board by ID | read | `idNotification`; *fields* |  |  |
| `TRELLO_GET_NOTIFICATIONS_BOARD_BY_ID_NOTIFICATION_BY_FIELD` | Get notification's board field | read | `field`, `idNotification` |  |  |
| `TRELLO_GET_NOTIFICATIONS_BY_ID_NOTIFICATION` | Get notification by ID | read | `idNotification`; *board, board_fields, card, card_fields, display, entities, fields, list and 6 more* |  |  |
| `TRELLO_GET_NOTIFICATIONS_BY_ID_NOTIFICATION_BY_FIELD` | Get a notification field | read | `field`, `idNotification` |  |  |
| `TRELLO_GET_NOTIFICATIONS_CARD_BY_ID_NOTIFICATION` | Get notification's card | read | `idNotification`; *fields* |  |  |
| `TRELLO_GET_NOTIFICATIONS_CARD_BY_ID_NOTIFICATION_BY_FIELD` | Get notification card field | read | `field`, `idNotification` |  |  |
| `TRELLO_GET_NOTIFICATIONS_DISPLAY_BY_ID_NOTIFICATION` | Get notification display by id | read | `idNotification` |  |  |
| `TRELLO_GET_NOTIFICATIONS_ENTITIES_BY_ID_NOTIFICATION` | Get Notification Entities | read | `idNotification` |  |  |
| `TRELLO_GET_NOTIFICATIONS_LIST_BY_ID_NOTIFICATION` | Retrieve notification list by id | read | `idNotification`; *fields* |  |  |
| `TRELLO_GET_NOTIFICATIONS_LIST_BY_ID_NOTIFICATION_BY_FIELD` | Get notification list field | read | `field`, `idNotification` |  |  |
| `TRELLO_GET_NOTIFICATIONS_MEMBER_BY_ID_NOTIFICATION` | Get notification member by id | read | `idNotification`; *fields* |  |  |
| `TRELLO_GET_NOTIFICATIONS_MEMBER_BY_ID_NOTIFICATION_BY_FIELD` | Retrieve notification member field | read | `field`, `idNotification` |  |  |
| `TRELLO_GET_NOTIFICATIONS_MEMBER_CREATOR_BY_ID_NOTIFICATION` | Get member creator of a notification | read | `idNotification`; *fields* |  |  |
| `TRELLO_GET_NOTIFICATIONS_MEMBER_CREATOR_BY_ID_NOTIFICATION_B` | Get notification creator field | read | `field`, `idNotification` |  |  |
| `TRELLO_GET_NOTIFICATIONS_ORGANIZATION_BY_ID_NOTIFICATION` | Get notification organization | read | `idNotification`; *fields* |  |  |
| `TRELLO_GET_NOTIFICATIONS_ORGANIZATION_BY_ID_NOTIFICATION_BY_` | Get notification organization field | read | `field`, `idNotification` |  |  |
| `TRELLO_GET_ORGANIZATIONS_ACTIONS_BY_ID_ORG` | Get organization actions by ID | read | `idOrg`; *before, display, entities, fields, filter, format, idModels, limit and 6 more* |  |  |
| `TRELLO_GET_ORGANIZATIONS_BOARDS_BY_ID_ORG` | Get organization boards | read | `idOrg`; *action_fields, actions, actions_entities, actions_format, actions_limit, actions_since, fields, filter and 4 more* |  |  |
| `TRELLO_GET_ORGANIZATIONS_BY_ID_ORG` | Get organization by ID | read | `idOrg`; *action_fields, actions, actions_display, actions_entities, actions_limit, board_action_fields, board_actions, board_actions_display and 17 more* |  |  |
| `TRELLO_GET_ORGANIZATIONS_BY_ID_ORG_BY_FIELD` | Get organization field by id | read | `field`, `idOrg` |  |  |
| `TRELLO_GET_ORGANIZATIONS_DELTAS_BY_ID_ORG` | Get organization deltas by id | read | `idOrg`, `ixLastUpdate`, `tags` |  |  |
| `TRELLO_GET_ORGANIZATIONS_MEMBERSHIPS_BY_ID_ORG` | Get organization memberships | read | `idOrg`; *filter, member, member_fields* |  |  |
| `TRELLO_GET_ORGANIZATIONS_MEMBERS_BY_ID_ORG` | Retrieve organization members by id | read | `idOrg`; *activity, fields, filter* |  |  |
| `TRELLO_GET_ORGANIZATIONS_MEMBERS_INVITED_BY_ID_ORG` | Get invited organization members | read | `idOrg`; *fields* |  |  |
| `TRELLO_GET_ORGANIZATIONS_MEMBERS_INVITED_BY_ID_ORG_BY_FIELD` | Get organization invited member field | read | `field`, `idOrg` |  |  |
| `TRELLO_GET_ORGANIZATIONS_PLUGIN_DATA` | Get organization plugin data | read | `id` |  | `read` |
| `TRELLO_GET_ORGANIZATIONS_TAGS_BY_ID` | Get organization tags | read | `id` |  | `read` |
| `TRELLO_GET_ORG_MEMBERSHIP` | Get organization membership | read | `idMembership`, `idOrg`; *member, member_fields* |  |  |
| `TRELLO_GET_ORG_NEW_BILLABLE_GUESTS_BY_BOARD` | Get new billable guests | read | `id`, `idBoard` |  | `read` |
| `TRELLO_GET_SEARCH` | Search Trello | read | `query`; *board_fields, boards_limit, card_attachments, card_board, card_fields, card_list, card_members, card_stickers and 11 more* |  |  |
| `TRELLO_GET_SEARCH_MEMBERS` | Search for members | read | `query`; *idBoard, idOrganization, limit, onlyOrgMembers* |  |  |
| `TRELLO_GET_SESSIONS_SOCKET` | Get sessions socket | read |  |  |  |
| `TRELLO_GET_TOKENS_BY_TOKEN` | Get token by token | read | `token`; *fields, webhooks* |  |  |
| `TRELLO_GET_TOKENS_BY_TOKEN_BY_FIELD` | Get token field | read | `field` |  |  |
| `TRELLO_GET_TOKENS_MEMBER_BY_TOKEN` | Get token member | read | *fields* |  |  |
| `TRELLO_GET_TOKENS_MEMBER_BY_TOKEN_BY_FIELD` | Retrieve token member field | read | `field` |  |  |
| `TRELLO_GET_TOKENS_WEBHOOKS_BY_TOKEN` | Get webhooks for token | read |  |  |  |
| `TRELLO_GET_TOKENS_WEBHOOKS_BY_TOKEN_BY_ID_WEBHOOK` | Get token webhook by ID | read | `idWebhook` |  |  |
| `TRELLO_GET_TYPES_BY_ID` | Get type by id | read | `id` |  |  |
| `TRELLO_GET_WEBHOOKS_BY_ID_WEBHOOK` | Get webhook by ID | read | `idWebhook` |  |  |
| `TRELLO_GET_WEBHOOKS_BY_ID_WEBHOOK_BY_FIELD` | Get webhook field by id | read | `field`, `idWebhook` |  |  |
| `TRELLO_LIST_EMOJI` | List available emoji | read | *locale, spritesheets* |  | `read` |
| `TRELLO_MARK_ALL_NOTIFICATIONS_READ` | Mark all notifications as read | write |  |  |  |
| `TRELLO_MARK_BOARD_AS_VIEWED` | Mark board as viewed | write | `idBoard` |  |  |
| `TRELLO_MARK_CARD_NOTIFICATIONS_READ` | Mark card notifications read | write | `idCard` |  |  |
| `TRELLO_MOVE_ALL_LIST_CARDS` | Move all cards in list to another list | write | `idBoard`, `idList`, `idListDestination` |  |  |
| `TRELLO_REMOVE_BOARD_MEMBER` | Delete member from board | write | `idBoard`, `idMember` |  |  |
| `TRELLO_REMOVE_BOARD_STAR` | Delete member board star | write | `idBoardStar`, `idMember` |  |  |
| `TRELLO_REMOVE_CARD_LABEL_BY_COLOR` | Delete card label by color | write | `color`, `idCard` |  |  |
| `TRELLO_REMOVE_LABEL_FROM_CARD` | Remove label from card | write | `idCard`, `idLabel` |  |  |
| `TRELLO_REMOVE_MEMBER_FROM_CARD` | Remove member from card | write | `idCard`, `idMember` |  |  |
| `TRELLO_REMOVE_MEMBER_FROM_ORG` | Delete member from organization | write | `idMember`, `idOrg` |  |  |
| `TRELLO_SET_CARD_CUSTOM_FIELD_VALUE` | Set card custom field value | write | `idCard`, `idCustomField`; *checked, date, idValue, number, text* |  | `write` |
| `TRELLO_UPDATE_ACTIONS_BY_ID_ACTION` | Update action text | write | `idAction`; *text* |  |  |
| `TRELLO_UPDATE_ACTIONS_TEXT_BY_ID_ACTION` | Update action text | write | `idAction`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_BY_ID_BOARD` | Update board by ID | write | `idBoard`; *closed, desc, idBoardSource, idOrganization, keepFromSource, labelNames__blue, labelNames__green, labelNames__orange and 23 more* |  |  |
| `TRELLO_UPDATE_BOARDS_CLOSED_BY_ID_BOARD` | Update board closed status | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_DESC_BY_ID_BOARD` | Update board description | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_ID_ORGANIZATION_BY_ID_BOARD` | Update board organization | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_LABEL_NAMES_BLUE_BY_ID_BOARD` | Update a board's blue label name | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_LABEL_NAMES_GREEN_BY_ID_BOARD` | Update board's green label name | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_LABEL_NAMES_ORANGE_BY_ID_BOARD` | Update board label orange name | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_LABEL_NAMES_PURPLE_BY_ID_BOARD` | Update purple label name | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_LABEL_NAMES_RED_BY_ID_BOARD` | Update red label name | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_LABEL_NAMES_YELLOW_BY_ID_BOARD` | Update yellow label name on board | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_MEMBERS_BY_ID_BOARD` | Update board members | write | `idBoard`; *email, fullName, type* |  |  |
| `TRELLO_UPDATE_BOARDS_MEMBERS_BY_ID_BOARD_BY_ID_MEMBER` | Update board member attributes | write | `idBoard`, `idMember`; *email, fullName, type* |  |  |
| `TRELLO_UPDATE_BOARDS_MY_PREFS_EMAIL_POSITION_BY_ID_BOARD` | Modify board email position preference | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_MY_PREFS_ID_EMAIL_LIST_BY_ID_BOARD` | Update board email list preference | write | `idBoard`, `value` |  |  |
| `TRELLO_UPDATE_BOARDS_MY_PREFS_SHOW_SIDEBAR_BY_ID_BOARD` | Update board sidebar visibility | write | `idBoard`, `value` |  |  |
| `TRELLO_UPDATE_BOARDS_NAME_BY_ID_BOARD` | Update board name | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_PREFS_BACKGROUND_BY_ID_BOARD` | Update board background | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_PREFS_CALENDAR_FEED_ENABLED_BY_ID_BOARD` | Update board calendar feed enabled status | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_PREFS_CARD_AGING_BY_ID_BOARD` | Update board card aging preference | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_PREFS_CARD_COVERS_BY_ID_BOARD` | Update board card cover preference | write | `idBoard`, `value` |  |  |
| `TRELLO_UPDATE_BOARDS_PREFS_COMMENTS_BY_ID_BOARD` | Update board comment preferences | write | `idBoard`, `value` |  |  |
| `TRELLO_UPDATE_BOARDS_PREFS_INVITATIONS_BY_ID_BOARD` | Update board invitation preferences | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_PREFS_PERMISSION_LEVEL_BY_ID_BOARD` | Update board prefs permission level | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_PREFS_SELF_JOIN_BY_ID_BOARD` | Update board self-join preference | write | `idBoard`, `value` |  |  |
| `TRELLO_UPDATE_BOARDS_PREFS_VOTING_BY_ID_BOARD` | Update board voting preferences | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARDS_SUBSCRIBED_BY_ID_BOARD` | Update board subscription status | write | `idBoard`, `value` |  |  |
| `TRELLO_UPDATE_BOARD_MEMBERSHIP` | Update board membership | write | `idBoard`, `idMembership`; *member_fields, type* |  |  |
| `TRELLO_UPDATE_BOARD_SIDEBAR_ACTIONS_PREFS` | Update board sidebar actions preferences | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARD_SIDEBAR_ACTIVITY_PREFS` | Update my board sidebar activity preference | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARD_SIDEBAR_MEMBERS_PREF` | Update board sidebar members preference | write | `idBoard`; *value* |  |  |
| `TRELLO_UPDATE_BOARD_STARS_POSITION_BY_ID_MEMBER_BY_ID_BOARD_` | Update board star position | write | `idBoardStar`, `idMember`; *value* |  |  |
| `TRELLO_UPDATE_CARDS_ACTIONS_COMMENTS_BY_ID_CARD_BY_ID_ACTION` | Update card comment | write | `idAction`, `idCard`; *text* |  |  |
| `TRELLO_UPDATE_CARDS_BY_ID_CARD` | Update card attributes | write | `idCard`; *closed, desc, due, fileSource, idAttachmentCover, idBoard, idCardSource, idLabels and 8 more* |  |  |
| `TRELLO_UPDATE_CARDS_CHECKLISTS_CHECK_ITEMS_NAME_BY_ID_CARD_B` | Update checklist item name in card | write | `idCard`, `idCheckItem`, `idChecklist`; *value* |  |  |
| `TRELLO_UPDATE_CARDS_CHECKLIST_ITEMS_BY_ID_CARD_BY_ID_CHECKLI` | Update check item on card | write | `id`, `idCheckItem`; *due, dueReminder, idChecklist, idMember, name, pos, state* |  | `read`, `write` |
| `TRELLO_UPDATE_CARDS_CHECKLIST_ITEMS_POSITION_BY_ID_CARD_BY_I` | Update checklist item position | write | `idCard`, `idCheckItem`, `idChecklist`, `value` |  |  |
| `TRELLO_UPDATE_CARDS_CHECKLIST_ITEMS_STATE_BY_ID_CARD_BY_ID_C` | Update card checklist item state | write | `idCard`, `idCheckItem`, `idChecklist`; *value* |  |  |
| `TRELLO_UPDATE_CARDS_CLOSED_BY_ID_CARD` | Update card closed status by id | write | `idCard`, `value` |  |  |
| `TRELLO_UPDATE_CARDS_DESC_BY_ID_CARD` | Update card description | write | `idCard`; *value* |  | `write` |
| `TRELLO_UPDATE_CARDS_DUE_BY_ID_CARD` | Update card due date | write | `idCard`; *value* |  |  |
| `TRELLO_UPDATE_CARDS_ID_ATTACHMENT_COVER_BY_ID_CARD` | Update card attachment cover | write | `idCard`; *value* |  |  |
| `TRELLO_UPDATE_CARDS_ID_BOARD_BY_ID_CARD` | Move card to board | write | `idCard`; *idList, value* |  |  |
| `TRELLO_UPDATE_CARDS_ID_LIST_BY_ID_CARD` | Update card list ID | write | `idCard`, `value` |  |  |
| `TRELLO_UPDATE_CARDS_ID_MEMBERS_BY_ID_CARD` | Set card members | write | `idCard`; *value* |  |  |
| `TRELLO_UPDATE_CARDS_LABELS_BY_ID_CARD` | Update card labels | write | `idCard`, `value` |  |  |
| `TRELLO_UPDATE_CARDS_NAME_BY_ID_CARD` | Update card name | write | `idCard`, `value` |  |  |
| `TRELLO_UPDATE_CARDS_POS_BY_ID_CARD` | Update card position | write | `idCard`; *value* |  |  |
| `TRELLO_UPDATE_CARDS_STICKERS_BY_ID_CARD_BY_ID_STICKER` | Update card sticker attributes | write | `idCard`, `idSticker`; *image, left, rotate, top, zIndex* |  |  |
| `TRELLO_UPDATE_CARDS_SUBSCRIBED_BY_ID_CARD` | Update card subscribed status | write | `idCard`, `value` |  |  |
| `TRELLO_UPDATE_CHECKLISTS_BY_ID_CHECKLIST` | Update checklist | write | `idChecklist`; *idBoard, idCard, idChecklistSource, name, pos* |  |  |
| `TRELLO_UPDATE_CHECKLISTS_BY_ID_CHECKLIST_BY_FIELD` | Update checklist field | write | `field`, `id`, `value` |  | `read`, `write` |
| `TRELLO_UPDATE_CHECKLISTS_ID_CARD_BY_ID_CHECKLIST` | Move checklist to card | write | `idChecklist`; *value* |  |  |
| `TRELLO_UPDATE_CHECKLISTS_NAME_BY_ID_CHECKLIST` | Update checklist name | write | `idChecklist`; *value* |  |  |
| `TRELLO_UPDATE_CHECKLISTS_POS_BY_ID_CHECKLIST` | Update checklist position by id | write | `idChecklist`; *value* |  |  |
| `TRELLO_UPDATE_CHECKLIST_ITEM_BY_IDS` | Update a checklist item | write | `idCard`, `idCheckItem`, `idChecklistCurrent`; *idChecklist, name, pos, state* |  |  |
| `TRELLO_UPDATE_CUSTOM_FIELD` | Update custom field | write | `idCustomField`; *displayCardFront, name, pos* |  | `read`, `write` |
| `TRELLO_UPDATE_CUSTOM_FIELD_OPTION` | Update custom field option | write | `idCustomField`, `idCustomFieldOption`; *color, text* |  | `write` |
| `TRELLO_UPDATE_LABELS_BY_ID_LABEL` | Update label attributes | write | `idLabel`; *color, idBoard, name* |  |  |
| `TRELLO_UPDATE_LABELS_BY_ID_LABEL_BY_FIELD` | Update label field | write | `field`, `id`, `value` |  | `read`, `write` |
| `TRELLO_UPDATE_LABELS_COLOR_BY_ID_LABEL` | Update label color | write | `idLabel`; *value* |  |  |
| `TRELLO_UPDATE_LABELS_NAME_BY_ID_LABEL` | Update label name | write | `idLabel`; *value* |  |  |
| `TRELLO_UPDATE_LISTS_BY_ID_LIST` | Update list attributes | write | `idList`; *closed, idBoard, idListSource, name, pos, subscribed* |  |  |
| `TRELLO_UPDATE_LISTS_CLOSED_BY_ID_LIST` | Update lists closed status | write | `idList`; *value* |  |  |
| `TRELLO_UPDATE_LISTS_ID_BOARD_BY_ID_LIST` | Move list to board | write | `idList`, `value`; *pos* |  |  |
| `TRELLO_UPDATE_LISTS_NAME_BY_ID_LIST` | Update list name | write | `idList`, `value` |  |  |
| `TRELLO_UPDATE_LISTS_POS_BY_ID_LIST` | Update list position | write | `idList`; *value* |  |  |
| `TRELLO_UPDATE_LISTS_SUBSCRIBED_BY_ID_LIST` | Update list subscription | write | `idList`, `value` |  |  |
| `TRELLO_UPDATE_MEMBERS_AVATAR_SOURCE_BY_ID_MEMBER` | Update member avatar source | write | `idMember`, `value` |  |  |
| `TRELLO_UPDATE_MEMBERS_BIO_BY_ID_MEMBER` | Update member bio | write | `idMember`; *value* |  |  |
| `TRELLO_UPDATE_MEMBERS_BOARD_BACKGROUNDS_BY_ID_MEMBER_BY_ID_B` | Update Member Board Background | write | `idBoardBackground`, `idMember`; *brightness, file, tile* |  |  |
| `TRELLO_UPDATE_MEMBERS_BOARD_STARS_BY_ID_MEMBER_BY_ID_BOARD_S` | Update Member Board Star | write | `idBoardStar`, `idMember`; *idBoard, pos* |  |  |
| `TRELLO_UPDATE_MEMBERS_BY_ID_MEMBER` | Update member by ID | write | `idMember`; *avatarSource, bio, fullName, initials, prefs__colorBlind, prefs__locale, prefs__minutesBetweenSummaries, username* |  |  |
| `TRELLO_UPDATE_MEMBERS_CUSTOM_BOARD_BACKGROUNDS_BY_ID_MEMBER_` | Update member's custom board background | write | `idBoardBackground`, `idMember`; *brightness, file, tile* |  |  |
| `TRELLO_UPDATE_MEMBERS_FULL_NAME_BY_ID_MEMBER` | Update member full name by id | write | `idMember`; *value* |  |  |
| `TRELLO_UPDATE_MEMBERS_INITIALS_BY_ID_MEMBER` | Update member initials | write | `idMember`; *value* |  |  |
| `TRELLO_UPDATE_MEMBERS_PREFS_COLOR_BLIND_BY_ID_MEMBER` | Update member color blind preference | write | `value`; *idMember* |  |  |
| `TRELLO_UPDATE_MEMBERS_PREFS_LOCALE_BY_ID_MEMBER` | Update member locale preference | write | `idMember`; *value* |  |  |
| `TRELLO_UPDATE_MEMBERS_PREFS_SUMMARY_INTERVAL_BY_ID_MEMBER` | Update Member Summary Interval | write | `value`; *idMember* |  |  |
| `TRELLO_UPDATE_MEMBERS_SAVED_SEARCHES_BY_ID_MEMBER_BY_ID_SAVE` | Update member saved search | write | `idMember`, `idSavedSearch`; *name, pos, query* |  |  |
| `TRELLO_UPDATE_MEMBERS_SAVED_SEARCHES_NAME_BY_ID_MEMBER_BY_ID` | Update member saved search name | write | `idMember`, `idSavedSearch`, `value` |  |  |
| `TRELLO_UPDATE_MEMBERS_SAVED_SEARCHES_POS_BY_ID_MEMBER_BY_ID_` | Update member saved search position | write | `idMember`, `idSavedSearch`, `value` |  |  |
| `TRELLO_UPDATE_MEMBERS_SAVED_SEARCHES_QUERY_BY_ID_MEMBER_BY_I` | Update member saved search query | write | `idMember`, `idSavedSearch`, `value` |  |  |
| `TRELLO_UPDATE_MEMBERS_USERNAME_BY_ID_MEMBER` | Update member username | write | `idMember`, `value` |  |  |
| `TRELLO_UPDATE_MEMBER_AVATAR` | Upload member avatar | write | `idMember`; *file* |  |  |
| `TRELLO_UPDATE_NOTIFICATIONS_BY_ID_NOTIFICATION` | Update notification status by id | write | `idNotification`; *unread* |  |  |
| `TRELLO_UPDATE_NOTIFICATIONS_UNREAD_BY_ID_NOTIFICATION` | Update notification unread status | write | `idNotification`; *value* |  |  |
| `TRELLO_UPDATE_ORGANIZATIONS_BY_ID_ORG` | Update organization attributes | write | `idOrg`; *desc, displayName, name, prefs__associatedDomain, prefs__boardVisibilityRestrict__org, prefs__boardVisibilityRestrict__private, prefs__boardVisibilityRestrict__public, prefs__externalMembersDisabled and 4 more* |  |  |
| `TRELLO_UPDATE_ORGANIZATIONS_DESC_BY_ID_ORG` | Update organization description | write | `idOrg`; *value* |  |  |
| `TRELLO_UPDATE_ORGANIZATIONS_DISPLAY_NAME_BY_ID_ORG` | Update organization display name | write | `idOrg`, `value` |  |  |
| `TRELLO_UPDATE_ORGANIZATIONS_MEMBERSHIPS_BY_ID_ORG_BY_ID_MEMB` | Update org membership | write | `idMembership`, `idOrg`; *member_fields, type* |  |  |
| `TRELLO_UPDATE_ORGANIZATIONS_MEMBERS_BY_ID_ORG` | Update an organization member | write | `idOrg`; *email, fullName, type* |  |  |
| `TRELLO_UPDATE_ORGANIZATIONS_MEMBERS_BY_ID_ORG_BY_ID_MEMBER` | Update organization member | write | `idMember`, `idOrg`, `type` |  |  |
| `TRELLO_UPDATE_ORGANIZATIONS_MEMBERS_DEACTIVATION_BY_ID_ORG_B` | Update org member deactivation status | write | `idMember`, `idOrg`, `value` |  |  |
| `TRELLO_UPDATE_ORGANIZATIONS_NAME_BY_ID_ORG` | Update organization name by id | write | `idOrg`, `value` |  |  |
| `TRELLO_UPDATE_ORGANIZATIONS_PREFS_ASSOCIATED_DOMAIN_BY_ID_OR` | Update an organization's associated domain preferences | write | `idOrg`; *value* |  |  |
| `TRELLO_UPDATE_ORGANIZATIONS_PREFS_BOARD_VISIBILITY_BY_ID_ORG` | Update organization board visibility preference | write | `idOrg`; *value* |  |  |
| `TRELLO_UPDATE_ORGANIZATIONS_PREFS_EXTERNAL_MEMBERS_ACCESS_BY` | Update an organization's external members access | write | `idOrg`; *value* |  |  |
| `TRELLO_UPDATE_ORGANIZATIONS_PREFS_PERMISSION_LEVEL_BY_ID_ORG` | Update organization permission level | write | `idOrg`; *value* |  |  |
| `TRELLO_UPDATE_ORGANIZATIONS_PREFS_PRIVATE_BOARD_VISIBILITY_B` | Update organization private board visibility restriction | write | `idOrg`, `value` |  |  |
| `TRELLO_UPDATE_ORGANIZATIONS_PREFS_PUBLIC_BOARD_VISIBILITY_BY` | Update organization public board visibility | write | `idOrg`; *value* |  |  |
| `TRELLO_UPDATE_ORGANIZATIONS_WEBSITE_BY_ID_ORG` | Update organization website | write | `idOrg`; *value* |  |  |
| `TRELLO_UPDATE_ORGANIZATION_LOGO` | Add organization logo | write | `file`, `idOrg` |  |  |
| `TRELLO_UPDATE_ORG_GOOGLE_APPS_VERSION` | Update organization Google Apps version | write | `idOrg`; *value* |  |  |
| `TRELLO_UPDATE_ORG_INVITE_RESTRICTIONS` | Update organization invite restriction | write | `idOrg`, `value` |  |  |
| `TRELLO_UPDATE_SESSIONS_BY_ID_SESSION` | Update session by id | write | `idSession`; *idBoard, status* |  |  |
| `TRELLO_UPDATE_WEBHOOKS_ACTIVE_BY_ID_WEBHOOK` | Update webhook active status | write | `idWebhook`, `value` |  |  |
| `TRELLO_UPDATE_WEBHOOKS_BY_ID_WEBHOOK` | Update webhooks | write | `idWebhook`; *active, callbackURL, description, idModel* |  |  |
| `TRELLO_UPDATE_WEBHOOKS_CALLBACK_URL_BY_ID_WEBHOOK` | Update webhook callback url | write | `idWebhook`, `value` |  |  |
| `TRELLO_UPDATE_WEBHOOKS_DESCRIPTION_BY_ID_WEBHOOK` | Update webhook description | write | `idWebhook`; *value* |  |  |
| `TRELLO_UPDATE_WEBHOOKS_ID_MODEL_BY_ID_WEBHOOK` | Update webhook idModel | write | `idWebhook`, `value` |  |  |
| `TRELLO_WHO_AM_I` | Who Am I | read |  |  |  |
