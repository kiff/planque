# App guides

One short guide per app a taf can connect, for every app planque offers.
Coach reads an app's guide before it offers a job that uses the app. tafs
cloud preselects the guides of the apps a job connects or offers, so the
choice is deterministic, never the model's.

planque serves each adapted guide as it's built
(`GET /v1/accounts/{tenant}/catalog/{entry}/guide`, `guides.go` beside this
file): a new or changed guide reaches Coach with the next deploy.

## Two parts per app

- **Generated:** `generated/<slug>.md`, written by `planque guides` from the
  catalog as planque reads it, for every app at once. It says how the app
  connects and, for each tool: what it does, read or write, its inputs, what
  it returns (the ids and links a taf keeps, as the catalog's `returns` gives
  them) and the access it needs (as the catalog's `unusable` checks it).
  Never edited by hand: run the command again.
- **Adapted:** `<slug>.md`, beside this file, written by a person in our own
  words from the vendor's published documentation and from real runs: the
  tools a job usually needs and what to keep, what not to offer, plan and
  rate limits, and the app's error answers with what to tell the owner.
  Nothing is copied from elsewhere, and no outside product is named.

`apps.md` lists every app, the first ones first (what Explore shows and what
people connect first), in batches of about ten: each batch is one person's
share of the adapted parts.

## Generating

Where planque's connection-service key is reachable (the same
`PLANQUE_COMPOSIO_SECRET` and AWS credentials as the server):

```
go run ./cmd/planque guides -out deploy/guides
```

It reads every app that can be signed in to now and each one's tools, plus
the reviewed entries in `deploy/catalog.json`, and reads nothing of any
owner's. Without the key it writes the reviewed entries only.

## Rules for the adapted part

- **Dated.** The first line under the title is `Checked: <date>`. A guide older
  than 90 days is checked again before Coach relies on it.
- **Short.** Under 1,200 words each.
- **Plain words for the owner.** The "Tell the owner" lines are what Coach
  says. They never name the connection service, a status code or a scope.
- **Not a permission.** A guide never widens what a taf may do. The Card and
  the grant decide; a guide only helps Coach offer what will work.

## Adapted so far

| App | Adapted | Generated |
|---|---|---|
| Gmail | [gmail.md](gmail.md) | [gmail](generated/gmail.md) |
| Google Sheets | [googlesheets.md](googlesheets.md) | [googlesheets](generated/googlesheets.md) |
| Stripe (Stripe's own server) | [stripe.md](stripe.md) | [stripe](generated/stripe.md) (the managed app); its own server lists its tools at sign-in |
| GitHub | [github.md](github.md) | [github](generated/github.md) |
| LinkedIn | [linkedin.md](linkedin.md) | [linkedin](generated/linkedin.md) |
| Google Calendar | [googlecalendar.md](googlecalendar.md) | [googlecalendar](generated/googlecalendar.md) |
| Google Drive | [googledrive.md](googledrive.md) | [googledrive](generated/googledrive.md) |
| Google Docs | [googledocs.md](googledocs.md) | [googledocs](generated/googledocs.md) |
| Outlook | [outlook.md](outlook.md) | [outlook](generated/outlook.md) |
| Slack | [slack.md](slack.md) | [slack](generated/slack.md) |
| Qonto (key) | [qonto.md](qonto.md) | [qonto-key-1](generated/qonto-key-1.md), [qonto-key-2](generated/qonto-key-2.md) |
| bexio | [bexio.md](bexio.md) | [bexio](generated/bexio.md) |
| Bitbucket | [bitbucket.md](bitbucket.md) | [bitbucket](generated/bitbucket.md) |
| Cal | [cal.md](cal.md) | [cal](generated/cal.md) |
| Clay | [clay.md](clay.md) | [clay](generated/clay.md) |
| ClickUp | [clickup.md](clickup.md) | [clickup](generated/clickup.md) |
| Confluence | [confluence.md](confluence.md) | [confluence](generated/confluence.md) |
| Discord | [discord.md](discord.md) | [discord](generated/discord.md) |
| Discord Bot | [discordbot.md](discordbot.md) | [discordbot](generated/discordbot.md) |
| Dropbox | [dropbox.md](dropbox.md) | [dropbox](generated/dropbox.md) |
| Dynamics 365 | [dynamics365.md](dynamics365.md) | [dynamics365](generated/dynamics365.md) |
