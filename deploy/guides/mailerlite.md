# MailerLite

Checked: 2026-11-14 · Sources: the connection service's MailerLite tool
definitions as planque reads them, and MailerLite's API documentation
(getting started, subscribers, groups, campaigns).

## How it connects

- **By the owner's API key**, pasted once and held by the connection
  service; the taf never sees it. The owner makes the key in MailerLite
  under Integrations, then the MailerLite API.
- **A key belongs to the person who made it.** If that person is removed
  from the account, or their user is deleted, the key stops working and the
  owner pastes a new one.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.

## Tools a job usually needs

Read means it changes nothing in the account. MailerLite's catalog entries
give no `returns` paths today; the returns below are what MailerLite's
documentation shows in each answer.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `MAILERLITE_GET_GROUPS` | Lists the account's groups | read | group ids, names, counts |
| `MAILERLITE_GET_GROUP_SUBSCRIBERS` | Lists the subscribers in one group | read | subscriber ids, emails, statuses |
| `MAILERLITE_GET_SUBSCRIBERS` | Lists subscribers, a page at a time | read | subscriber ids, emails, statuses, a page cursor |
| `MAILERLITE_GET_SUBSCRIBER` | Reads one subscriber, by id or email | read | the subscriber's id, status |
| `MAILERLITE_CREATE_SUBSCRIBER` | Adds a subscriber, or updates one with that email | write | the subscriber's id, status |
| `MAILERLITE_UPDATE_SUBSCRIBER` | Changes one subscriber's fields | write | the subscriber |
| `MAILERLITE_ADD_SUBSCRIBER_TO_GROUP` | Puts a subscriber in a group | write | the group |
| `MAILERLITE_REMOVE_SUBSCRIBER_FROM_GROUP` | Takes a subscriber out of a group | write | nothing (an empty answer) |
| `MAILERLITE_GET_CAMPAIGNS` | Lists campaigns by status or type | read | campaign ids, statuses, delivery stats |
| `MAILERLITE_GET_CAMPAIGN` | Reads one campaign | read | the campaign's id, status |
| `MAILERLITE_CREATE_CAMPAIGN` | Creates a campaign draft; sends nothing | write | the campaign's id, status |
| `MAILERLITE_UPDATE_CAMPAIGN` | Changes a draft | write | the campaign |
| `MAILERLITE_SEND_OR_SCHEDULE_CAMPAIGN` | Sends a campaign now, or at a set time | write, outside | the campaign's id, status |
| `MAILERLITE_CANCEL_CAMPAIGN` | Takes a scheduled campaign back to draft, before it starts sending | write | the campaign, as a draft |

**What to keep:** each group's id; a campaign's id from
`MAILERLITE_CREATE_CAMPAIGN`, to change, send or cancel that same draft; a
subscriber's id. Never search for a campaign the taf made: use the id it
kept.

**Drafts first, sending only on request.** Creating a campaign sends
nothing; `MAILERLITE_SEND_OR_SCHEDULE_CAMPAIGN` is what sends it, to
everyone it's addressed to. Offer it only when the owner asked for sending,
and each send waits for their yes. A campaign can be scheduled only once it
has recipients, a designed email and a working unsubscribe link.

**Adding is not overwriting.** `MAILERLITE_CREATE_SUBSCRIBER` on an email
already in the account updates only the fields sent, and leaves the rest and
the groups alone.

**People who opted out stay out.** Updating a subscriber who unsubscribed,
bounced or marked a campaign as spam doesn't make them active again. Adding
them again with the resubscribe option does, and MailerLite says to use it
only when the person asked to subscribe again. A taf never tries.

## Don't offer

- Deleting or forgetting subscribers (`MAILERLITE_DELETE_SUBSCRIBER`,
  `MAILERLITE_FORGET_SUBSCRIBER`): forgetting erases everything about the
  person within 30 days, and is for erasure requests only.
- Deleting groups, segments, fields or campaigns (`MAILERLITE_DELETE_GROUP`,
  `MAILERLITE_DELETE_SEGMENT`, `MAILERLITE_DELETE_FIELD`,
  `MAILERLITE_DELETE_CAMPAIGN`): a campaign that is sending is stopped first.
- Bulk imports and batches (`MAILERLITE_IMPORT_SUBSCRIBERS_TO_GROUP`,
  `MAILERLITE_EXECUTE_BATCH_REQUEST`): many changes in one step, hard to
  check.
- Automations, webhooks and double opt-in settings: they change what
  happens to every future subscriber.

## Access and plans

- **Rate:** 120 calls a minute. Imports, and batches made
  only of subscriber adds, have their own limit of 5 a minute. Over a
  limit, MailerLite says how long to wait.
- **Paging:** the subscriber list comes 25 at a time and a group's
  subscribers 50 at a time, unless asked for more; up to 1,000 a page.
- **Plans:** resend and multivariate campaigns need a paid plan or the free
  trial.
- **A subscriber's email** can change at most once a calendar month.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The key doesn't work (401) | The key is wrong, was deleted, or its maker left the account | "MailerLite needs a new key. Nothing was sent." |
| Refused, not allowed (403) | The account or this key isn't allowed to do this | "MailerLite doesn't let this connection do that step. Want to do without it?" |
| Not found (404) | The subscriber, group or campaign was deleted, or the id is wrong | "That isn't in MailerLite any more. The taf will skip it." |
| Rejected as invalid (422) | A field failed MailerLite's checks, such as a badly formed email | "MailerLite didn't accept that entry. The taf will skip it and tell you which." |
| Too many calls (429) | Over the per-minute limit | "MailerLite asked us to slow down. The taf will try again in a minute." |
| Conflict (409) | The change clashes with the item's current state | "MailerLite couldn't make that change right now. Nothing was changed." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write, and a failed send is never retried without the owner's
yes.
