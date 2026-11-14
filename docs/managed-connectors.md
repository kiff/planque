# Managed connections

`internal/connector` defines the boundary to a connection platform. The
Composio adapter uses its v3.1 HTTP API and the Go standard library. It
reuses the platform's app catalog, managed authentication, tool definitions
and concrete execution API. It does not run another agent framework.

The adapter supports paginated app and tool discovery, a toolkit-scoped
sign-in session, account inspection and one versioned action per execution.
It requires explicit account selection. Before execution it checks the
account's owner, app and active status, then checks the selected tool's
pinned version and schema. Platform meta-executors are refused. A lost or
malformed execution answer remains unknown and is never automatically
sent again. A refusal provides no proof that the action had no effect.

Some of an app's tools are never served, however a taf asks: those that
move money or close a bill for good, listed per app in `connector.NeverServedTools` (for Stripe,
voiding, deleting, finalizing or paying an invoice, voiding a credit note or
credit grant, and marking an invoice uncollectible). Like the meta tools,
they never come back from the tools search, an app's catalog or the
generated guides. Asking for one in a catalog request's `tools`, or setting
it on a grant, is a 400 with `code` `tool_not_served` and `tool` its slug
(`field` `tools` or `tool`). A connection that holds one from before no
longer lists it, and a call to it is not sent (`tool_not_served`); adding
tools to that connection is refused the same way until it's removed.

The same goes for a vendor's own server, where one broad tool writes to any
operation (Stripe's `stripe_api_write`): the operations it's never narrowed
to are in each describer's `Never` (for Stripe, voiding, paying, finalizing,
deleting or marking an invoice uncollectible, attaching a payment to one,
and voiding a credit note or credit grant). Narrowing it to one is a 400
with `code` `tool_not_served`, `tool` the tool's name and `field` the
operation's argument, before anything is read. A call through it that names
one, pinned on a connection from before or set by the taf on a broad tool,
in any spelling, is not sent (`tool_not_served`). It fails closed: a call
whose arguments aren't an object, or whose operation isn't one string, is
not sent either.

Project credentials are injected server-side. `SecretFromStore` reads one
operator-configured secret containing `api_key`, with a one-minute cache.
No request can choose a secret reference. Credential requests follow no
redirects, and backing-service error bodies are not returned.

With provider-managed authentication, the connection service holds the vendor
credentials: OAuth tokens and their refresh, or a key or bearer token entered
on its hosted connection page. Neither the taf nor planque receives that key
or token. The adapter requests upstream revocation
when disconnecting the exact selected account. A successful deletion
immediately disables that connection, while upstream revocation remains
pending; the response keeps the background job id when supplied.

## Configure the provider

Set `PLANQUE_COMPOSIO_SECRET` to a Secrets Manager reference containing
`{"api_key":"..."}`. The runtime needs permission to read that exact secret.
The reference enables the Composio provider; without it, existing direct
MCP, registered OAuth and API-key connections continue to work independently.
No agent request can enable a provider or supply its project key.

For the AWS deployment, pass the same reference as `PLANQUE_COMPOSIO_SECRET`
to `deploy/aws.sh up`. The script resolves the existing secret's exact ARN,
grants the function only `secretsmanager:GetSecretValue` on that resource,
and sets the runtime variable. It does not create or read the secret value.
The deployment identity needs permission to describe that secret and update
the function's existing role/configuration. Fake-service deployments receive
no connector-secret permission or runtime setting.


## Connect an app

The account-authenticated catalog endpoints include upstream discovery:

1. `GET /v1/accounts/{tenant}/catalog?provider=composio&q=...` returns
   `entries`, `provider` and `next_cursor`. Pass the same search/provider and
   `cursor=...` for the next page. `skipped` counts entries that fail adapter
   validation; they do not discard valid entries or the next-page cursor.
   Configured reviewed entries remain available.
   With no search, only an app the provider project can sign someone in to
   is listed: one with an enabled auth config, read with the scoped key and
   cached for five minutes. This includes API-key and bearer-token configs,
   whose hosted page collects the owner's credential. If the auth configs
   can't be read, the provider's apps aren't listed. A search lists every app
   the provider has (kiff/planque#87), each as it stands: one with an enabled
   config; one the provider can sign people in to with its own credentials
   (`composio_managed_auth_schemes`), takes only a key or token, or needs no
   sign-in, which is set up on the owner's pick (step 2); or one that needs
   our own developer app first, listed as `needs_setup` and not started until
   an operator enables it in the provider's dashboard, which the next read
   picks up. One read walks the provider's catalog (1000 apps a page, at most
   ten pages); stopped with pages left, it says `partial: true`. The first
   page with no search is kept with the sign-ins it was filtered by, so
   opening Connect apps again reads nothing upstream.
   **Finding tools by words.** `GET /v1/accounts/{tenant}/catalog/tools?q=…`
   is the provider's text search over tool names and descriptions
   (`GET /tools?query=`), at most 50 a page, narrowed upstream to one app when
   `apps` names one and here when it names several; each tool comes with its
   app as it stands, read once per app, so Coach can propose apps before any
   connection exists. The reviewed catalog's entries are in it, on the
   first page and when `apps` names none, with the entry's id and state:
   every tool of an entry whose app a word of `q` names (three letters or
   more) first, and a tool a word names only by its own words after the
   provider's apps. Coach finds apps only this way. With the provider
   down, the entries found still come back, with `partial: true`. Platform
   meta tools never come back. An app's events
   are listed the same way (`…/catalog/{entry}/events`, and
   `…/grants/{id}/events` for a connection), from `GET /triggers_types`, each
   pinned to a version; the approval of one stays as below.
2. `POST /v1/accounts/{tenant}/catalog/connector:composio:{app}/grants`
   accepts `{"agent":"...", "agent_name":"..."}` (`agent_name` optional, the
   taf's name for people; without it the words say "this taf", and the agent
   id is under Details). An app the provider doesn't have is 404; one that
   needs our own developer credentials is 409 until an operator enables it.
   An app the provider signs people in to by itself, with no sign-in set up
   yet, is **set up on this pick**: planque creates one provider-managed auth
   config (`POST /auth_configs` with `use_composio_managed_auth` and nothing
   of ours in it), takes the answer only when the provider says it manages
   it, and keeps its id under a per-app record that is also the lock, so two
   owners picking the app at once make one. The grant pins that config as a
   pre-set app's grant does, and is marked `on_demand: true`, as is every
   later grant on a sign-in set up this way, so Coach drafts its writes and
   money to ask the owner by default. An account sets up at most ten apps
   this way a day (429 after). Every grant view says who holds the sign-in
   (`custody`: `provider` here, `planque` for a sign-in planque holds). It returns
   a pending grant and no sign-in link (`needs_key: false`). Its text names
   the app and the taf, and says who keeps the sign-in in plain words, never a
   service's name: "Your <app> sign-in is kept by our sign-in partner; the
   taf never sees it" (the owner signs these words, so a third party keeping
   it is said before the approval). For `API_KEY` it says the owner will enter
   their app key on our sign-in partner's page, which keeps it, and the taf
   never sees it; `BEARER_TOKEN` says the same about a token. The config ID and scheme are part of the signed request,
   so later catalog changes cannot rewrite an existing grant's words. If an
   app has several enabled configs, the adapter selects the lowest ID in
   lexical order, independently of page ordering, and pins that config.

   **The narrowest sign-in (gap 5 of kiff/tafs#476).** A request may name the
   taf's tools (`tools`, the provider's slugs, as the tools search names
   them). For an app on the provider's own sign-in, planque reads that
   sign-in's scope ceiling and the scopes it always asks (the app's
   `composio_managed_auth` and `auth_config_details`), asks the provider for
   the tools' least-privilege scopes within the ceiling
   (`toolkits/{slug}/scopes/recommended`), and pins a managed sign-in that
   asks only for those (`credentials.scopes`). It is made once per app and
   scope set, named `planque-scoped-<app>-<hash>`, found again by that name,
   and never chosen as an app's usual sign-in. A new one counts against the
   account's on-demand set-ups for the day, as setting up an app does. Slugs
   are compared in upper case, as the provider names them. The scopes and the tools are
   in the signed request (`asked_tools`, `scopes`), the words say "<app> is
   asked only for what this taf's actions need.", and the grant binds only those
   tools later (409 for another: a new connection asks for it).

   When the provider can't narrow it (no ceiling, a tool outside it, its
   recommendation failing, an operator's own sign-in, or the account's
   set-ups for the day used), the sign-in asks
   for the app's usual access, recorded as `usual_access` so a later
   narrowing can replace it, and the words say it plainly, with why: "<app>
   will ask for more access than this taf uses, because its sign-in can't be
   narrowed to this taf's actions.", and the card ends "Sign in, then approve
   its actions: the taf gets only those, nothing before." When the app's
   usual sign-in asks for personal data (Google's date of birth, addresses,
   phone numbers, other email addresses, contacts, language), it's never
   used: a sign-in that can't be narrowed, or one asked with no tools, asks
   for the managed ceiling less those, pinned as its own (wide), recorded in
   `scopes`, and still said to be more than the taf uses. A personal scope
   stays only for a tool asked that runs on nothing else (the narrowest
   read-only one); a tool added later that needs a dropped one isn't
   covered. When no such sign-in can be set up (the day's set-ups used, the
   provider failing), nothing is asked: 429 or 503, with why. The account's
   own address stays. Neither names the tools: the asked ones are the provider's names
   ("Fetch message by thread id"), and the actions to approve come next, each
   by the title the taf's request gives it. KIFF's Card stays the boundary
   either way: the taf calls only the tools its Card grants. A reused sign-in
   carries its narrowing to the new grant.
3. The owner's device signs `/approve` with that grant's digest. Only then
   is a scoped hosted sign-in created: `/tool_router/session` pins the approved
   config through `auth_configs`, and `/tool_router/session/{id}/link` creates
   the hosted link for OAuth, keys and bearer tokens alike. A session that
   reports a different or missing pinned config is refused before `/link`.
   The answer is `{status: "approved", sign_in_url, app_name, expires_at}`.
   `app_name` is the app the owner signs in to, as people name it ("Google
   Sheets"): the app's card says "Sign in to <app_name>", never the
   provider's host that `sign_in_url` goes through.
   The callback independently checks the
   server-allocated user, app and exact connection; callback account/status
   parameters establish no authority. The approving device confirms the
   browser's one-use code through the existing `/confirm` route.

   A sign-in that never comes back to the app (the provider's page, or a
   browser closed before its redirect, so no callback; or a callback whose
   code never reached the app) is never stuck until its deadline. The device
   that approved signs `/approve` again, and planque asks the provider:
   - signed in there: the grant is finished as the callback finishes it, and
     the answer is `{status: "signed_in", confirm, confirm_by}`. The device
     confirms with that code through `/confirm`, as with the browser's. A
     code the callback minted before no longer confirms;
   - not signed in there: a fresh sign-in link for the same grant,
     `{status: "approved", sign_in_url, app_name, expires_at}`, and the unfinished
     provider account is revoked.

   Another device gets 403. `DELETE …/grants/{id}` cancels a connection
   still waiting on its sign-in, as it disconnects any other.
4. After confirmation, the asking agent reads
   `GET /v1/accounts/{tenant}/grants/{id}/catalog?cursor=...`. Tool discovery
   is paginated and cached separately, rather than copied into the grant.
   Each tool has a stable gateway-compatible alias. Read the relevant page
   before selecting actions; discovery snapshots expire after 24 hours.
5. `/tools` selects the concrete aliases and supplies a drafted profile
   (`source: "drafted"`, `transport: "rest"`, `pin: {"kind":"none"}`).
   The grant digest covers the exact connected account, upstream tool
   versions and schemas as well as the profile. Upstream catalog membership
   never establishes reviewed profile status or test evidence. At most 20
   actions and 64 KB of selected upstream definitions fit one grant.
6. The device signs the updated digest and checked profile/words digests.
   Registration and KIFF profile application finish before any action can
   execute. Each action still requires an exact-call permit and the relay's
   once-only operation record. Permit expiry is checked again immediately
   before the provider execution request, after upstream preflight.
7. **What a call says when it doesn't go through.** A refusal tells the taf
   who refused it and why: the vendor's own answer the provider passed on
   ("Stripe refused it: No such invoice: 'in_1'."), or the provider's own
   refusal with its status ("The connection provider refused it (it answered
   403): …"). An answer that never came says the provider's status and words
   too. Each still makes no claim about whether an effect occurred. The
   reason is one line of at most 300 characters, with keys, bearer tokens,
   webhook secrets and JWTs redacted (`connector.Reason`). planque's own
   refusals before the provider read "Not sent: …".
8. **Which account, and which mode.** The grant view of a managed connection
   carries `account_name`, the provider's name for the vendor account, and
   `mode`: `live` or `test` once a permitted call returned an object that
   says (`livemode`), `unknown` until then. Unknown is treated as live.
   `mode_words` says it for every app: "Live", "Test", or "Mode not known
   yet: treated as live".
   planque never calls the vendor outside a permit to find out, and an
   empty list says nothing. Live wins: any live object read makes it
   `live` for good, and `test` only stands while no live one was seen.

Provider-managed key apps use the same hosted path as OAuth. The catalog's
existing `how: "oauth"` value selects that path; it does not ask the phone for
a key. Operator-reviewed direct API-key entries retain `how: "key"` and their
separate `/key` flow. Grants made before the config binding keep their signed
fields and existing wording.

Disconnect disables the specific provider connection and reports upstream
credential revocation as pending. Unfinished sign-ins and unapproved held
connections are disabled by the existing sweep. Replacing a connected
account requires disconnecting and requesting a new grant; an existing
approved profile never silently switches accounts.

Synthetic tests exercise real signed HTTP handlers, the Composio HTTP
adapter and the KIFF test service, including account isolation, profile
application, exact-call permits, concurrent replay, lost answers and expiry.
They do not establish live app compatibility. A configured project and
owner-controlled test account are needed to validate an individual app.
Key and bearer fixtures prove that the hosted `/link` flow reaches `connected`
only after an active account is checked and both device approvals complete;
they do not submit credentials to a live provider.

Protocol references:

- [App catalog](https://docs.composio.dev/reference/api-reference/toolkits/getToolkits)
- [Session configuration](https://docs.composio.dev/docs/configuring-sessions)
- [Authentication config IDs and schemes](https://docs.composio.dev/reference/api-reference/auth-configs/getAuthConfigs)
- [Authentication and token custody](https://docs.composio.dev/docs/authentication)
- [Versioned tool definitions](https://docs.composio.dev/reference/api-reference/tools/getToolsByToolSlug)
- [Concrete execution](https://docs.composio.dev/reference/api-reference/tools/postToolsExecuteByToolSlug)
- [Account deletion and upstream revocation](https://docs.composio.dev/reference/api-reference/connected-accounts/deleteConnectedAccountsByNanoid)

## Which tools are read-only

A tool is read-only only when its vendor marks it so (the provider's
`readOnlyHint` tag, or MCP's `readOnlyHint` on a vendor's own server) **and**
nothing about it says it can change things. Unknown means a write.
`connector.MayChange` counts a tool as a write when:

- its name has `GRAPH_QL`, `GQL`, `EXECUTE`, `PROXY`, `RAW` or `CUSTOM` as a
  word, or says it's an API call or HTTP request;
- its description says it can change things ("including mutations", "may
  perform mutations", any or arbitrary request, a custom document);
- an argument, at any depth (inside objects, array items and `anyOf`,
  `oneOf` or `allOf`), may take a free-form document, program, body or raw
  request: `document`, `mutation`, `body`, `payload`, `raw`, `sql`,
  `statement`, `script`, `command`, `code` and the like;
- a `query` argument says it's GraphQL, or doesn't say what it takes;
- a `method` argument sits beside a `url`, `path` or `endpoint`;
- an argument, or the schema, is given by reference (`$ref`): it isn't
  followed;
- its schema can't be read.

When unsure, it's a write. Some reads come out writes too ("custom
properties", raw file content, a discount `code` to look up): that only costs
Setup a read it can't make. A search's query whose description says so
("Gmail advanced search query") is a read's argument and stays one.
This matters past planque: KIFF's gateway, Setup's read Card, the meter and
the runtime all treat a read-only tool as one that only reads, and the owner
is told "reads only". A GraphQL tool a vendor tags read-only while it "may
perform mutations" is granted as a write.

## Read-only live discovery check

With an operator-controlled project already configured, run:

```sh
LIVE_COMPOSIO_SECRET=your-secret-reference go test -tags live -run '^TestLiveComposioDiscovery$' -v ./live
```

This resolves the project credential from Secrets Manager, follows all public
app catalog pages, checks one toolkit's paginated concrete tool definitions,
and compares an exact version's schema with discovery. `LIVE_COMPOSIO_TOOLKIT`
selects the public toolkit (default `github`). AWS credentials come from the
environment or the operator's AWS CLI profile; neither credentials nor raw
provider responses are printed. The check creates no sign-in sessions or
connected accounts and executes no vendor action. Passing it establishes
live discovery/definition conformance only; owner-confirmed authentication,
a permitted test action and revocation still need separate evidence.

## Incoming Slack door events

Managed Slackbot connections can deliver channel messages and direct
messages through Composio triggers. This uses the `slackbot` toolkit and
its managed OAuth app. A message arriving does not authorize a reply;
replies use the taf's approved send tool through KIFF.

Configure `PLANQUE_DOOR_URL` as the fixed HTTPS receiver and
`PLANQUE_DOOR_SECRET` as a Secrets Manager reference containing
`{"secret":"..."}`. The receiver uses the same literal UTF-8 secret.
Deployment resolves this secret's exact ARN and grants read access alongside
the configured project secret. Requests cannot supply another receiver or
credential reference.

After the approving device has confirmed the connection:

- Device-signed `GET /v1/accounts/{tenant}/grants/{id}/channels` lists the
  channels a door may be set on, through the connection's own pinned
  `SLACKBOT_LIST_ALL_CHANNELS` action (refused unless it is read-only):
  `{channels: [{id, name, private}]}`, public and private, not archived, at
  most 200, no content.
- `GET /v1/accounts/{tenant}/grants/{id}/triggers/types?slug=...` reads the
  current definition and exact account, authenticated with the KIFF key of
  the connection's own agent or of the agent that asked for it. Supported
  slugs are `SLACKBOT_CHANNEL_MESSAGE_RECEIVED` and
  `SLACKBOT_DIRECT_MESSAGE_RECEIVED`.
- Device-signed `POST /v1/accounts/{tenant}/grants/{id}/triggers` takes
  `{account, slug, channels}`: a supported slug, and 1–100 distinct exact
  Slack conversation IDs. planque resolves the definition and pins its
  version, and builds each upstream config itself, so no subscription
  covers more than the approved channels. A channel-message door makes one
  subscription per channel, `{channel_id, is_bot_message: false}` (never a
  null `channel_id`, which would mean every channel). A direct-message door
  takes DM IDs (`D…`) and makes one subscription with no upstream filter;
  delivery is routed to its DM IDs only. Any other slug is 400. The answer is
  `{id, trigger, triggers, account, channels, definition}`, `id` and
  `trigger` naming the first subscription. If any step fails after an
  upstream subscription was made, every subscription of the request is
  deleted upstream and its record disabled, even if the request is gone; one
  the provider does not delete at once is noted and deleted by the sweep.
  A subscription already on before the request, and its record, stay as
  they were. A note carries its record's version: it is deleted upstream
  only after a write conditioned on that version marks it as being removed
  with a two-minute lease, exclusive through the provider call (taken over
  only once lapsed, settled only by a write conditioned on the claim),
  a newer record retires the note, and an approval that reuses an id being
  removed is 409 until the removal is done. The approval is limited to
  64 KB. The adapter independently rechecks account ownership and the
  pinned trigger definition before registering an instance. No listener is
  created without device approval.
- `GET .../triggers` lists this connection's records, for the same agents as
  `.../triggers/types`. Device-signed `DELETE .../triggers/{trigger}` takes an
  optional `{account}` (the connection's own account when left out) and
  disables local delivery before deleting the upstream listener.
  Disconnecting stops listeners and revokes the connection. Active door
  subscriptions keep a confirmed sign-in in use, even when no business
  action has been selected.
- A grant's view names its catalog `entry` (`connector:{provider}:{toolkit}`)
  and `provider`.

The first approved subscription registers this installation's fixed
`/connectors/composio/webhook` as the project webhook. An existing subscription
must already use that URL, V3 payloads and trigger-message events; planque
refuses to replace another destination. The returned signing secret stays
sealed in the existing KV/KMS store. A project key without subscription or
trigger write permissions is refused explicitly. Read-only discovery alone
does not establish those permissions.

Incoming raw bytes are authenticated with `webhook-id`, `webhook-timestamp`
and `webhook-signature` before routing. Signatures use HMAC-SHA256 over
`id.timestamp.rawBody`, with a five-minute tolerance. The handler matches the
stored trigger, connected account, user and allowed channel. Bot messages
and unrelated events are filtered. Tenant and KIFF agent come from the
approved connection record.

The receiver gets `{v:"planque.door.v1", id, tenant, agent, grant, provider,
account, trigger, trigger_slug, channel, sender, text, ts, thread?, timestamp}`.
`ts` is the message's own Slack ts, verbatim; `thread` its parent when it is already in a thread.
`sender` is the Slack user ID; the receiver's sender controls decide whether
work starts. Headers `Planque-Event-Id`, `Planque-Event-Timestamp` and
`Planque-Event-Signature` use the same signing format and the shared handoff
secret. Message IDs are stable for the connected account, channel and Slack
message timestamp, including across upstream retries.

Delivery is durable and payloads remain sealed while pending. The receiver
must deduplicate the stable ID and return 202 when taken or already taken.
Only 5xx or a lost answer is retried, with a lease and the scheduled sweep as
backup. Other statuses, including 401, 404 and 409, are terminal. Completed
or cancelled deliveries drop their message content. Disabling a route cancels
its queued deliveries; deleting the account removes its global event and
routing records. The oversight export includes the signed trigger approvals,
without project credentials or signing secrets.

## App events for a taf

Any other event a connection's app lists can wake a taf: a new deal, a paid
invoice, a reaction. Configure `PLANQUE_EVENT_URL` as its fixed HTTPS
receiver; events are signed with the same `PLANQUE_DOOR_SECRET`. Without it,
no app event can be approved (503).

The refund event a profile's case block names also goes to KIFF, which
follows a refund it decided until the payment processor says how it ended.
`PLANQUE_KIFF_EVENT_URL` is KIFF's fixed HTTPS receiver and
`PLANQUE_KIFF_EVENT_SECRET` names the Secrets Manager secret those events are
signed with: both or neither, at start and in `deploy/aws.sh`, which gives the
function a read of that one secret and nothing wider. Without them, no event
goes to KIFF. Either way, an event reaches KIFF only on a subscription the
owner approved for the taf.

- `GET .../triggers/types?slug=...` reads any event the connection's app
  lists, with its config schema and `door` saying whether it is a Slack door's
  event instead.
- Device-signed `POST .../triggers` takes `{account, slug, config}`. `config`
  is the owner's filters exactly as signed, at most 8 KB: an object whose
  every field the event's config schema declares, each required field
  present, each value of its declared type and, where the schema lists
  values, one of them. Anything else is refused, so a subscription never
  listens wider than what the owner read. It makes exactly one subscription
  with that config and answers `{id, subscription, account, event, config,
  definition}`. A door's events still take `channels`, and an app event never
  does.
- Taking back on failure, the cleanup sweep, disabling and disconnecting work
  as for doors.

### Stripe refunds through KIFF's governed path

A refund's money rules are reviewed code, never a model's draft. Refunds go
through the direct Stripe connection (catalog entry `stripe`, a sign-in to
`https://mcp.stripe.com` whose consent page offers test mode). When
`PLANQUE_STRIPE_REFUNDS` offers refunds to the account (`on`: every account,
for self-hosting and tests; `tenants:<id>[,<id>]`: those accounts only, the
others answered as if it were off; anything else: off) and such a connection
asks for the refund, planque
supplies the whole request itself (`stripe-refunds@1`,
`internal/planque/stripe_refunds.go`), the way it supplies the browser's.

- **What's asked:** one tool and no profile: `stripe_api_write` with
  `fixed_arguments` `{"stripe_api_operation_id": "PostRefunds", "livemode":
  false}`. Only test mode is offered for now: in live mode, Stripe may ask the
  owner to confirm each refund in its Dashboard, and that hasn't been checked,
  so `livemode: true` is refused (409). Anything else is refused, the caller's own
  account, arguments or profile included. Refunds go on a Stripe connection
  of their own, asked before it has tools; the taf's other Stripe tools go on
  another connection, with the profile Coach drafts.
- **What's granted:** five tools, each on the connection's account and in the asked mode:
  - `stripe_api_details`, Setup's read of an operation's fields;
  - `stripe_api_read` narrowed to `GetAccountsAccount` (the account read),
    `GetChargesCharge` (the charge, by `charge`) and `GetRefunds` (its
    refunds, by `charge`, at most 100);
  - `stripe_api_write` narrowed to `PostRefunds`, with `charge` and `amount`.

  **Which account:** Stripe's sign-in answer names none. Right after the
  owner signs in, planque calls Stripe's own `list_available_accounts_or_orgs`
  once, itself, with the new token, never on a caller's word. When it lists
  exactly one account, planque keeps it, with its mode, as the connection's
  `account`; with several or none, it keeps none. A refund connection without
  an account is asked to sign in again choosing one account in test mode
  (409), and one whose account is in live mode is refused (409).
- **The profile:** the refund moves money, in `EUR` minor units. planque
  writes its marker at `parameters.metadata.kiff_op`, inside the operation's
  fields, and the refund Stripe makes carries it at `metadata.kiff_op`, where
  the case's refunds read finds it. The case reads `id` from the account,
  `currency`, `amount_captured` and `disputed` from the charge, and `data[]`
  with `id`, `amount`, `status` and `metadata.kiff_op` from the refunds list.
  Stripe writes the currency in lower case (`eur`), and KIFF compares it
  without case. KIFF refuses a refund on a charge in another currency.
- **Checked against Stripe's own description:** the narrowed tools' paths,
  the marker included, are checked against the fields `stripe_api_details`
  lists for each operation, once the owner approved it, before the refund is
  first sent. A marker inside `parameters` is allowed only on a tool narrowed
  to one write, and only under exactly that prefix.
- **Off by default.** Stripe's shapes were read on a test-mode connection;
  turn it on only after the account read is confirmed there too.

A drafted profile with a `case` block is refused on every connection.

The receiver gets `{v:"planque.event.v1", id, tenant, agent, grant, provider,
account, subscription, event, payload, timestamp}`: `event` is the event's
slug, `subscription` the approved subscription and `payload` the event's data
as the app sent it. The ID is stable for the connection, subscription and the
provider's event, so a redelivery is the same event. An event whose data is
over 64 KB is dropped and noted on the connection, never cut. Delivery,
sealing, retries and terminal answers are as for door messages.
