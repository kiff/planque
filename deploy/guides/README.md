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

## Adapted guides

Every app in [apps.md](apps.md) has its adapted guide, in the same order.

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
| Notion | [notion.md](notion.md) | [notion](generated/notion.md) |
| HubSpot | [hubspot.md](hubspot.md) | [hubspot](generated/hubspot.md) |
| Shopify | [shopify.md](shopify.md) | [shopify](generated/shopify.md) |
| Calendly | [calendly.md](calendly.md) | [calendly](generated/calendly.md) |
| Trello | [trello.md](trello.md) | [trello](generated/trello.md) |
| Asana | [asana.md](asana.md) | [asana](generated/asana.md) |
| Airtable | [airtable.md](airtable.md) | [airtable](generated/airtable.md) |
| Agent Mail | [agent_mail.md](agent_mail.md) | [agent_mail](generated/agent_mail.md) |
| Ashby | [ashby.md](ashby.md) | [ashby](generated/ashby.md) |
| Attio | [attio.md](attio.md) | [attio](generated/attio.md) |
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
| Excel | [excel.md](excel.md) | [excel](generated/excel.md) |
| Figma | [figma.md](figma.md) | [figma](generated/figma.md) |
| Freshdesk | [freshdesk.md](freshdesk.md) | [freshdesk](generated/freshdesk.md) |
| Google BigQuery | [googlebigquery.md](googlebigquery.md) | [googlebigquery](generated/googlebigquery.md) |
| Google Maps | [google_maps.md](google_maps.md) | [google_maps](generated/google_maps.md) |
| Google Meet | [googlemeet.md](googlemeet.md) | [googlemeet](generated/googlemeet.md) |
| Google Photos | [googlephotos.md](googlephotos.md) | [googlephotos](generated/googlephotos.md) |
| Google Super | [googlesuper.md](googlesuper.md) | [googlesuper](generated/googlesuper.md) |
| Google Tasks | [googletasks.md](googletasks.md) | [googletasks](generated/googletasks.md) |
| Gorgias | [gorgias.md](gorgias.md) | [gorgias](generated/gorgias.md) |
| Holded | [holded.md](holded.md) | [holded](generated/holded.md) |
| Intercom | [intercom.md](intercom.md) | [intercom](generated/intercom.md) |
| Jira | [jira.md](jira.md) | [jira](generated/jira.md) |
| Lexoffice | [lexoffice.md](lexoffice.md) | [lexoffice](generated/lexoffice.md) |
| Linear | [linear.md](linear.md) | [linear](generated/linear.md) |
| Mailchimp | [mailchimp.md](mailchimp.md) | [mailchimp](generated/mailchimp.md) |
| MailerLite | [mailerlite.md](mailerlite.md) | [mailerlite](generated/mailerlite.md) |
| Microsoft Teams | [microsoft_teams.md](microsoft_teams.md) | [microsoft_teams](generated/microsoft_teams.md) |
| Monday | [monday.md](monday.md) | [monday](generated/monday.md) |
| Moneybird | [moneybird.md](moneybird.md) | [moneybird](generated/moneybird.md) |
| Odoo | [odoo.md](odoo.md) | [odoo](generated/odoo.md) |
| OneDrive | [one_drive.md](one_drive.md) | [one_drive](generated/one_drive.md) |
| Pipedrive | [pipedrive.md](pipedrive.md) | [pipedrive](generated/pipedrive.md) |
| QuickBooks | [quickbooks.md](quickbooks.md) | [quickbooks](generated/quickbooks.md) |
| Reddit | [reddit.md](reddit.md) | [reddit](generated/reddit.md) |
| Salesforce | [salesforce.md](salesforce.md) | [salesforce](generated/salesforce.md) |
| Sentry | [sentry.md](sentry.md) | [sentry](generated/sentry.md) |
| Sevdesk | [sevdesk.md](sevdesk.md) | [sevdesk](generated/sevdesk.md) |
| Slackbot | [slackbot.md](slackbot.md) | [slackbot](generated/slackbot.md) |
| Supabase | [supabase.md](supabase.md) | [supabase](generated/supabase.md) |
| Telegram | [telegram.md](telegram.md) | [telegram](generated/telegram.md) |
| Typeform | [typeform.md](typeform.md) | [typeform](generated/typeform.md) |
| WhatsApp | [whatsapp.md](whatsapp.md) | [whatsapp](generated/whatsapp.md) |
| Wise | [wise.md](wise.md) | [wise](generated/wise.md) |
| Wrike | [wrike.md](wrike.md) | [wrike](generated/wrike.md) |
| YouTube | [youtube.md](youtube.md) | [youtube](generated/youtube.md) |
| Zendesk | [zendesk.md](zendesk.md) | [zendesk](generated/zendesk.md) |
| Zoho | [zoho.md](zoho.md) | [zoho](generated/zoho.md) |
| Zoho Books | [zoho_books.md](zoho_books.md) | [zoho_books](generated/zoho_books.md) |
| Zoho Invoice | [zoho_invoice.md](zoho_invoice.md) | [zoho_invoice](generated/zoho_invoice.md) |
| Qonto (key) | [qonto.md](qonto.md) | [qonto-key-1](generated/qonto-key-1.md), [qonto-key-2](generated/qonto-key-2.md) |
