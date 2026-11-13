# LinkedIn

Checked: 2026-11-13 · Sources: the connection service's LinkedIn tool
definitions as planque reads them (29 tools), and LinkedIn's API
documentation.

## How it connects

- **A sign-in to LinkedIn**, held by the connection service. The catalog
  lists the tools and, with a known sign-in, puts the ones its access doesn't
  cover under `unusable`. A held sign-in is reused for a new taf when it
  covers the tools asked.
- **The managed sign-in acts as the person.** It can post as them and read
  their own profile. **Acting for a company page** (posting as the page,
  reading its statistics) needs page access that LinkedIn gives only to apps
  it has approved for that. Our managed sign-in can't ask for it today, so
  page tools come under `unusable`. Don't offer a page job yet.

## Tools a job usually needs

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| Get my info | Reads the signed-in person's name and id | read | the person's id |
| Create a post | Posts text, a link or an image as the person | write, public | the post's id |
| Get a post's content | Reads one of the person's posts | read | the post |
| Delete a post | Removes one of the person's posts | write | nothing |
| Get company info | Reads a company page's public details | read | the page's id |

**What to keep:** the person's id from "Get my info", needed as the author of
every post; each post's id, to read or delete it later.

**Posts are public.** Every post is drafted, shown to the owner, and sent
only on their yes. Offer "Delete a post" only for posts the taf made.

## Don't offer

- Messaging, connection requests, or reading other people's feeds: LinkedIn's
  API doesn't allow them to outside apps, and the tools that seem to promise
  them fail.
- Page posting and page statistics, until page access is set up.
- Scheduling posts in bulk: LinkedIn treats bursts as spam.

## Access and plans

- **No paid plan is needed** to post as oneself. Paid LinkedIn plans add
  nothing to what the API allows.
- **Rate:** LinkedIn's daily limits are per person and per app, and it
  doesn't publish all of them. A taf posts at most a few times a day.
- **Images** are uploaded first, then posted with the upload's id.
- **A sign-in lasts about 60 days.** After that the owner signs in again;
  the taf says so rather than failing silently.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused (403), not allowed | The tool needs page access, or a permission LinkedIn didn't give | "LinkedIn doesn't let this connection do that yet. Posting as yourself works; posting as your page doesn't." |
| Refused (422), duplicate | The same text was just posted | "LinkedIn refused the post because it's the same as a recent one. Want to change it?" |
| Too many calls (429) | Over LinkedIn's daily limit | "LinkedIn's daily limit is reached. The taf will try again tomorrow." |
| The sign-in no longer works (401) | Expired after about 60 days, or removed | "LinkedIn needs you to sign in again. Nothing was posted." |
| Not found (404) | The post was deleted | "That post isn't there any more. The taf will skip it." |

After any refusal the taf stops that step and says so. A refused post is
never retried with the text changed unless the owner changed it.
