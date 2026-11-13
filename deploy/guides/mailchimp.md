# Mailchimp

Checked: 2026-11-13 · Sources: the connection service's Mailchimp tool
definitions as planque reads them, and Mailchimp's Marketing API reference
(rate limits and timeouts, errors, OAuth, campaigns, lists, members).

## How it connects

- **A sign-in to Mailchimp** on Mailchimp's own page, held by the connection
  service; the taf never sees a password or token. planque's catalog lists
  the tools and, with a known sign-in, puts the ones its access doesn't
  cover under `unusable`, with why.
- **The sign-in doesn't expire.** It keeps working until the owner removes
  the app's access, or the person who signed in is removed from the
  Mailchimp account. Trying again won't bring it back; only a new sign-in
  does.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.

## Tools a job usually needs

Read means it changes nothing in the account. Mailchimp's catalog entries
give no `returns` paths today; the returns below are what Mailchimp's API
reference shows in each answer.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `MAILCHIMP_GET_LISTS_INFO` | Lists the account's audiences | read | audience ids, names |
| `MAILCHIMP_LIST_MEMBERS_INFO` | Lists the contacts in one audience, a page at a time | read | contact ids, emails, statuses |
| `MAILCHIMP_GET_MEMBER_INFO` | Reads one contact of an audience | read | the contact's id, status, tags |
| `MAILCHIMP_SEARCH_MEMBERS` | Finds contacts by email or name, in one audience or all | read | exact and partial matches |
| `MAILCHIMP_ADD_OR_UPDATE_LIST_MEMBER` | Adds a contact, or updates them if already there | write | the contact's id, status |
| `MAILCHIMP_ADD_OR_REMOVE_MEMBER_TAGS` | Adds or removes tags on one contact | write | nothing (an empty answer) |
| `MAILCHIMP_ADD_MEMBER_NOTE` | Adds a note to a contact | write | the note's id |
| `MAILCHIMP_LIST_CAMPAIGNS` | Lists campaigns and their status | read | campaign ids, statuses, archive links |
| `MAILCHIMP_ADD_CAMPAIGN` | Creates a campaign draft; sends nothing | write | the campaign id, its archive link |
| `MAILCHIMP_SET_CAMPAIGN_CONTENT` | Sets a draft's content | write | the content as saved |
| `MAILCHIMP_GET_CAMPAIGN_SEND_CHECKLIST` | Checks whether a draft is ready to send | read | ready or not, items to fix |
| `MAILCHIMP_SEND_TEST_EMAIL` | Sends a test of a draft to named addresses | write, outside | nothing (an empty answer) |
| `MAILCHIMP_SCHEDULE_CAMPAIGN` | Schedules a campaign for a set time | write, outside | nothing (an empty answer) |
| `MAILCHIMP_SEND_CAMPAIGN` | Sends a campaign to its audience now | write, outside | nothing (an empty answer) |
| `MAILCHIMP_GET_CAMPAIGN_REPORT` | Reads a sent campaign's results | read | the campaign id, emails sent, opens, clicks, unsubscribes |

**What to keep:** each audience's id; a campaign's id from
`MAILCHIMP_ADD_CAMPAIGN`, to set its content, check it, test it and send
that same draft; a contact's id. A contact's id is worked out from their
email in lower case, so the same email always leads to the same contact.
Never search for a campaign the taf made: use the id it kept.

**Drafts first, sending only on request.** "Prepare the newsletter and send
me a test" is the draft tools plus `MAILCHIMP_SEND_TEST_EMAIL` to the
owner's own address. Sending or scheduling a campaign reaches every contact
in its audience, so offer `MAILCHIMP_SEND_CAMPAIGN` and
`MAILCHIMP_SCHEDULE_CAMPAIGN` only when the owner asked for sending, and
each send waits for their yes. Run the send checklist before asking.

**Adding contacts:** `MAILCHIMP_ADD_OR_UPDATE_LIST_MEMBER` takes the status
to give a new contact. A job adds people as subscribed only when they
agreed to it; otherwise it keeps them out of the audience.

## Don't offer

- Deleting an audience (`MAILCHIMP_DELETE_LIST`): it loses the audience's
  whole history, its unsubscribes, complaints and bounces, and the contacts'
  addresses unless they were exported.
- Deleting a contact for good
  (`MAILCHIMP_DELETE_LIST_MEMBER`): it erases their details and they can
  never be added again. To take someone off, archive them
  (`MAILCHIMP_ARCHIVE_LIST_MEMBER`), and only when asked.
- Resending or deleting campaigns (`MAILCHIMP_RESEND_CAMPAIGN`,
  `MAILCHIMP_DELETE_CAMPAIGN`): a second send reaches people again, and a
  deletion loses the campaign.
- Changing an audience's settings (`MAILCHIMP_UPDATE_LISTS`): it changes
  what every future signup sees.

## Access and plans

- **Rate:** at most 10 calls open at once for each Mailchimp user, across
  all their apps. Over that, calls are refused. A job runs its calls one
  after another.
- **Calls stop after 120 seconds.** A big read should ask for fewer fields
  and page through.
- **Paging:** lists come 10 at a time unless asked for more, up to 1,000 a
  page.
- **Scheduled times** fall on the quarter hour only (:00, :15, :30, :45).
- **Some features follow the account's plan.** An endpoint the plan doesn't
  include is refused. Which plan includes which isn't checked here.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused, not allowed (403) | The sign-in's user role or the account's plan doesn't allow this step | "Mailchimp doesn't let this connection do that step. It may need a higher role on the account or a different plan. Want to do without it?" |
| Refused, account disabled (403, user disabled) | The Mailchimp account has been deactivated | "Your Mailchimp account is switched off, so the taf can't use it. Check the account in Mailchimp." |
| The sign-in no longer works (401) | Access was removed, or the person who signed in left the account | "Mailchimp needs you to sign in again. Nothing was sent." |
| Too many calls (429) | More than 10 calls at once for this user | "Mailchimp asked us to slow down. The taf will try again in a moment." |
| Not found (404) | The audience, contact or campaign was deleted, or the id is wrong | "That isn't in Mailchimp any more. The taf will skip it." |
| Refused, not allowed in this state (405) | The campaign's state blocks it, such as deleting one that is sending | "Mailchimp won't do that while the campaign is in its current state. Nothing was changed." |
| Rejected as invalid (400) | A field failed Mailchimp's checks | "Mailchimp didn't accept that entry. The taf will skip it and tell you which." |
| Mailchimp error (500) | Trouble on Mailchimp's side | "Mailchimp is having trouble right now. The taf will try again later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write, and a failed send is never retried without the owner's
yes: a call that timed out may still have gone through.
