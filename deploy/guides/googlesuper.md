# Google Super

Checked: 2026-11-14 · Sources: the connection service's Google Super tool
definitions as planque reads them (541 tools, [generated](generated/googlesuper.md)),
and Google's API references for the apps it spans.

## What it is

**One sign-in for many Google apps at once:** Gmail, Calendar, Drive, Docs,
Sheets, Slides, Tasks, Meet, Maps, Photos, Contacts, Analytics and more, in
a single app. Its tools are the same kinds as each app's own, under one
name prefix.

## When to offer it, and when not

- **Prefer each Google app on its own** (Gmail, Google Sheets, Google
  Calendar…). Each asks Google only for what its tools need, its guide
  is precise, and the owner's Card reads plainly.
- **Offer Google Super only** when a job truly spans four or more Google apps
  and the owner would otherwise sign in four times. Even then, ask for the
  few tools the job needs: the sign-in's access follows the tools asked,
  never the whole app.
- **Its sign-in is not the same as the single apps'.** An owner who already
  connected Gmail signs in again for Google Super, and the other way round.

## How it connects

- **A sign-in to Google**, held by the connection service. The catalog lists
  the tools and, with a known sign-in, puts the ones its access doesn't cover
  under `unusable`, with why. Offer only what's under `tools`.
- **One approval for the sign-in, one for the tools and the Card.** With
  hundreds of tools, the Card must name each tool the job uses, in words.

## Tools a job usually needs

The same as each app's guide, under this app's names. For example:

| Tool | Same as | Read or write |
|---|---|---|
| `GOOGLESUPER_FETCH_EMAILS` | Gmail's `GMAIL_FETCH_EMAILS` | read |
| `GOOGLESUPER_CREATE_EMAIL_DRAFT` | Gmail's draft | write |
| `GOOGLESUPER_SEND_EMAIL` | Gmail's send | write, outside |
| `GOOGLESUPER_EVENTS_LIST` / `GOOGLESUPER_FIND_FREE_SLOTS` | Calendar's reads | read |
| `GOOGLESUPER_CREATE_EVENT` | Calendar's create | write |
| `GOOGLESUPER_VALUES_GET` | Sheets' read | read |
| `GOOGLESUPER_FIND_FILE` / `GOOGLESUPER_UPLOAD_FILE` | Drive's find and upload | read / write |
| `GOOGLESUPER_CREATE_DOCUMENT` | Docs' create | write |
| `GOOGLESUPER_INSERT_TASK` | Tasks' add | write |

**Read each app's own guide** for what to keep, what not to offer, its
limits and its errors: [Gmail](gmail.md), [Google Sheets](googlesheets.md),
[Google Calendar](googlecalendar.md), [Google Drive](googledrive.md),
[Google Docs](googledocs.md), [Google Tasks](googletasks.md),
[Google Meet](googlemeet.md), [Google Maps](google_maps.md) and
[Google Photos](googlephotos.md). They all apply here.

## Don't offer

Everything each app's guide says not to offer, and in particular:
- deleting tools, of any app;
- sharing and permissions (Drive, Calendar);
- Gmail settings and forwarding;
- Google Ads changes of any kind (campaigns, budgets, bids, ads, accounts):
  they spend the owner's money;
- sending people's emails or other customer data to Ads or Analytics
  audiences;
- Analytics account and property admin, and shared drive admin (create,
  hide, change, delete).

## Access and plans

- **The limits are each app's own,** per user: Gmail's sending caps, Drive's
  storage, Sheets' per-minute writes.
- **A business admin** can block outside apps, and is more likely to for a
  sign-in asking for many apps at once.

## When it fails

As in each app's guide. One more:

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403) on one app's tool | The sign-in was made for other apps' tools | "Your Google connection doesn't cover this step yet. Sign in again and allow it, or I'll do without." |

After any refusal the taf stops that step and says so.
