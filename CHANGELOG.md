# Changelog

User-visible changes and compatibility notes are recorded here. Dates and
version headings refer to releases, not merely to merged pull requests.

## Unreleased

### Fixed

- **A refund connection's account read takes the account at `parameters.id`.** `GetAccountsAccount` takes its account as `id`, as a test-mode read showed; it was pinned at `parameters.account`, which the first-call field check would have refused, so no refund could have been decided.

### Changed

- **A refund connection's case reads are checked against Stripe's description.** A narrowed read is now checked at its first call, like a narrowed write, where its connection holds the describing tool (`stripe_api_details`): a path its operation doesn't take is refused with Stripe's own list. Without that tool, a narrowed read isn't checked, as before.
- **A case needs a fixed currency.** A profile's `case` block is refused unless its operation has a fixed `amount.unit`: KIFF compares the charge's currency only with a fixed unit.

### Fixed

- **An existing managed Stripe connection stays in the catalog.** Since the app shows one Stripe, a managed Stripe connection made before wasn't listed anywhere the owner could disconnect it. The catalog now lists it, connected, while it exists; a new one still can't be asked for.

### Added

- **A direct Stripe connection keeps its Stripe account.** Stripe's sign-in answer names the account as `stripe_user_id`; planque now keeps it as the connection's `account`, for Stripe's server only (`oauth.Metadata.AccountField`). Before, a direct Stripe connection had no account. Connections made before this keep none until they sign in again.
- **One Stripe in the app: the direct sign-in, with test mode.** The catalog lists `stripe`, a direct sign-in to Stripe's MCP server (`https://mcp.stripe.com`), whose consent page offers live mode, sandboxes and test mode, so a taf can be tried on a test account first. The managed connector's Stripe is no longer listed, found in a tool search, or connected (409, naming `stripe`), and `deploy/server-apps.json` lists Stripe's server under `stripe`.
- **Stripe refunds get a reviewed request, on the direct Stripe connection (off by default).** With `PLANQUE_STRIPE_REFUNDS=on`, a direct Stripe connection that asks for `stripe_api_write` narrowed to `PostRefunds` takes planque's whole request: `stripe_api_details`, the three case reads (`GetAccountsAccount`, `GetChargesCharge`, `GetRefunds`) and the refund, each on the account Stripe's sign-in named and in the asked mode, with planque's profile `stripe-refunds@1`: the refund in euros, its marker inside the operation's parameters, and its `case` block. A case's marker may sit inside `parameters` only on a tool narrowed to one write, checked against Stripe's description of the operation. A drafted profile can no longer carry a `case` block (docs/managed-connectors.md).
- Setup's listing of a taf's grants (`GET /v1/accounts/{tenant}/agents/{agent}/grants`)
  gives every tool's `setup_reads`: true for one of Setup's own reads, false
  otherwise. Those reads are on Setup's reader Card, never the taf's, so a
  check that the taf's Card covers every tool it was given can leave them out.
  The field is always sent, so a reader can tell this listing from an older
  one that had no such field.
- A profile's refund operation can carry a `case` block: the reads of the
  same profile that give a refund's facts (which processor account the
  connection is on, the charge, its refunds), where the call names its
  charge, where a dispute shows, and optionally the app event that carries a
  refund's later state (`refund_event`). KIFF decides a refund on these. The
  reads may reference only the call's charge, and the refunds read finds
  KIFF's own refunds by the operation's marker. planque renames the reads to
  KIFF's tool names with the rest of the profile.
- The refund event a case block names is handed to KIFF as well as to the
  taf, once, under an id of its own (`PLANQUE_KIFF_EVENT_URL`, signed with
  `PLANQUE_KIFF_EVENT_SECRET`). An event kept before its profile stopped
  naming it is dropped, not sent. Without the setting, nothing changes.
- The hosted deployment carries KIFF's refund-event receiver:
  `deploy/aws.sh` passes `PLANQUE_KIFF_EVENT_URL` and
  `PLANQUE_KIFF_EVENT_SECRET` to the function, both or neither, and lets it
  read that one secret only (`secretsmanager:GetSecretValue` on its exact
  ARN). Half of the pair is refused before anything changes
  (docs/managed-connectors.md, "App events for a taf").

### Changed

- planque's pages for a person in a browser speak the app's language: a
  vendor sign-in's result (not completed, another account, a used link) and
  the browser sign-in view, with its buttons, its note and its answers, in
  the eight languages the tafs apps ship. The language is the one the app
  asked in when it approved or opened the view, kept on the grant, else the
  browser's (the language it rates highest), else English (README, "The
  pages' language"). Before, every page was in English. A browser that
  didn't hand its sign-in back now logs why instead of showing the worker's
  error on the page.

### Fixed

- A managed connector's sign-in that comes back to planque answers the
  person with a page in their language, as a vendor sign-in does, not
  planque's JSON. A sign-in cancelled at the provider reads "The sign-in at
  … was not completed. Nothing was granted." and the connection keeps
  waiting. A missing, unknown or used link, and a failure on planque's side,
  read their own pages. The cause is still logged.

- The hosted deployment's table keeps point-in-time recovery on. Each
  `deploy/aws.sh up` asks for it, whether the table is new or already there,
  so a rebuilt table gets it back. Before, the script never asked for it, and
  a bad write couldn't be undone unless it had been turned on by hand.

- A tool is read-only only when its vendor marks it so and nothing about it
  says it can change things: a GraphQL tool tagged read-only that "may perform
  mutations", a free-form document, program or raw request among its
  arguments at any depth, a query that doesn't say what it takes, a method
  beside a url or path, or a name with `GRAPH_QL`, `EXECUTE`, `RAW` or
  `CUSTOM` makes it a write. When unsure, it's a write. Before, the vendor's mark
  alone made a tool read-only, and the owner was told it "reads only" while it
  could change things (docs/managed-connectors.md, "Which tools are read-only").

- A managed app's sign-in that never came back to the app no longer waits
  for its deadline behind a 409 ("this connection already has a sign-in").
  The device that approved it signs `/approve` again, and planque asks the
  provider. If the account is signed in there, the grant finishes as the
  callback would, and the device confirms with the code in the answer. If
  not, the answer is a fresh sign-in link for the same grant, and the
  unfinished account is revoked (kiff/planque#87, on the phone).

### Added

- A reviewed list of which app a vendor's own server is
  (`PLANQUE_SERVER_APPS_FILE`, shipped as `deploy/server-apps.json`). A grant on
  a listed server gets the app the list names, the managed app the owner
  already uses: `https://mcp.stripe.com` is `connector:composio:stripe`, and
  `https://mcp.agentmail.to` is `connector:composio:agent_mail`. An older grant
  gets it beside its request, and a new one in it, as for the catalog. The list
  comes before the catalog. A server not on it is matched as before, and one
  with two rows naming different apps gets no entry.
- Every grant names its app's catalog `entry`, so the apps list a connection
  under its app by `entry` alone (the follow-up to kiff/tafs#581). A grant
  asked by server, key or registered application names the one reviewed
  entry for the same server, reached the same way, in `request.catalog_entry`,
  signed. A grant made before gets the entry recorded beside its request,
  never in it, the first time the owner's grants list shows it: the signed
  request, its digest and its approval stay as they were. Where no entry
  matches, or more than one does, there's none. The catalog counts a
  connection Coach started as connected to its app.
- A connection the vendor refuses is fixed in place, from the app (gap 13 of
  kiff/tafs#476). `POST …/key` on a `connected` API-key grant replaces its
  key once one read with the new key works: a read-only operation that needs
  no argument, the grant's own first, or `getMe` for a Telegram bot, whose
  new token must be the same bot's. A key the vendor refuses, or one planque
  can't check, changes nothing, and the old key stays. `POST …/approve` on a
  `connected` sign-in grant signs in again while the connection goes on
  working: the new sign-in is held apart, in the same vendor account or not
  at all, until the approving device confirms it, then replaces the one in
  use, which is revoked at the vendor. One not confirmed in time is revoked
  and goes. Either way the grant id, its tools and its Card stay as they
  were. Once the vendor named a grant's account, a new sign-in that names
  none is refused like one in another account, for a reconnect too. A managed app can't sign in again in place (409, saying to disconnect
  it and connect it again), since planque can't tell which account a new
  sign-in is.
- A managed app asked for with the taf's tools (`tools`, the provider's
  slugs) signs in for only what they need: the provider's least-privilege
  scopes within its own sign-in's ceiling, on a scoped managed sign-in made
  once per app and scope set. The scopes are in the signed request and its
  words, and the grant binds only those tools. When the provider can't
  narrow it, the sign-in asks for the app's usual access, recorded on the
  grant, and the words say so: "<app> will ask for more access than this taf
  uses. The taf can still only use: …" (gap 5 of kiff/tafs#476). A new
  scoped sign-in counts against the account's on-demand set-ups for the day;
  once they're used, the sign-in asks for the app's usual access.
- A managed app's sign-in can be reused for another taf (`"reuse"` naming a
  connected managed grant, `tools: []`): the new connection is signed in
  and listed at once, on the same provider account, with no second sign-in
  or key, and its words say so. The owner approves its tools on their
  device as for any connection, and nothing is granted before that. The
  account stays while either connection uses it (kiff/tafs#476, gap 2).
- A vendor's refusal through a managed app carries the vendor's HTTP status
  when the provider passes it on ("Agent Mail refused it (it answered 403):
  Forbidden."), in the form KIFF's gateway reads, so a refused key can be
  told from a missing item (kiff/tafs#476, gap 10).

- The hourly sweep revokes a sign-in kept for another connection once no
  live connection uses it. A disconnect or expiry that keeps a shared
  sign-in (a reused token, or a provider's connected account) indexes it,
  so two connections on one sign-in disconnected at the same moment, each
  seeing the other as live, no longer leave it live at the vendor. A
  sign-in a live connection uses is never touched, and a provider account
  is revoked once.
- Setup's listings of the connections it asked for say how each reaches
  its vendor (`kind`): the vendor's own sign-in (`vendor_sign_in`), a
  provider's app (`provider`, with its `entry`, `mode` and `mode_words`),
  a REST key (`key`, with its API's URL, base and where the key goes,
  never the key) or a browser. So Coach can offer a new taf the kind of
  connection the account already uses for that vendor.
- `GET /v1/accounts/{tenant}/connections`: how the account's tafs reach
  their vendors (agent, vendor, kind, what asks for the same again, mode),
  only the connections that hold a sign-in,
  for any of the account's keys or sessions, as its catalog is. No grant id,
  approval, text, tool or key. tafs cloud gives it to Coach when a
  conversation starts.

### Changed

- What the owner reads names things as people do (kiff/tafs#476). The taf
  is named by `agent_name`, which every grant request now takes (an MCP
  server, an API key, a browser, a reviewed catalog entry, as a managed
  app's already did), else "this taf": never by its agent id. A managed
  connection's actions are named by their titles, and its account as
  "your connected <app> account". The ids go under Details (`terms`), after
  the tools' own lines: each action's slug and exact version, the connected
  account, and the taf's agent.
- The taf's tools and Setup's own reads are two lists: "<the taf> may:"
  lists only the taf's, and Setup's reads follow under "Setup reads, to set
  this up". A connection with only Setup's reads has no list for the taf.
- Compatibility: a request without `agent_name` keeps its encoding and
  digest, and a grant's stored words never change. A managed grant made
  before titles keeps its words, ids and all. One with a title and no name
  says "this taf" from its next approval on.

### Fixed

- The browser proxy's count of refused connections is read under the lock
  its writer holds (`Proxy.Refused()`): a test read it unlocked, a data race
  under `-race`.

- The last connection on a shared sign-in says what its disconnect did: "the
  sign-in it shared was revoked" (or, for a key, "the key it shared was
  deleted"), not that the sign-in stays.

- Disconnecting one taf's connection never signs another taf out. A sign-in
  another live connection still uses (a second taf that reused it, or a
  provider's connected account two connections pin) stays: only this
  connection ends, its tools out of KIFF, and its record says the sign-in
  stayed. The sign-in is revoked when the last connection on it is
  disconnected, so it never outlives them. An expiring held sign-in another
  connection uses is kept the same way.

- A refused call through a managed connection says who refused it and why:
  the vendor's answer, or the provider's with its status, redacted and on
  one line. Before, every refusal was one fixed line, the same for any call,
  so nobody could tell a wrong account or mode from a broken tool.
- A managed connection's grant shows its account's name and its mode
  (`live`, `test`, or `unknown`, treated as live), from the objects
  permitted calls returned, live winning, with no call of planque's own.
- Adding tools to a connection whose profile the owner approved: a profile
  of the new operations alone is merged into the approved one, which stays
  as it is. Before, every such addition was refused, since the caller has
  no copy of the approved profile to repeat, and an app had to be connected
  again, its key given again.

### Added

- Two optional fields on a profile's operation, for KIFF's "Did this happen?"
  ask: `record_url`, where the owner looks at the vendor's record of the
  resource (an https address with the resource's id path-escaped into
  `{resource}`, reviewed catalog only), and `label_read`, a read of the same
  profile that gives the resource's own name (for an operation with a
  `resource`). A profile that names neither is unchanged.

- A call whose outcome comes later (the slow-tools design): a write's
  profile may say `completes: {event, key, event_key, timeout_s, words}`.
  The relay then refuses the call while the taf doesn't listen for that
  event (`not_listening`), answers it as `pending` with the handle and the
  deadline once sent, and delivers the completing event with the same
  handle as its `key`, so the run that made the call can wait for it and
  be resumed by it.
- One fingerprint, then taps (kiff/tafs#494): a routine signed call (approving
  a sign-in or a key, confirming, a trigger, enrolling) is signed by the
  device key inside an approval window the phone's approval key opened at
  KIFF (`Planque-Device-Signature: t=…,v2=…,w=<window id>`, over
  `planque.request.v2` with the window's id last), or by the approval key
  itself; planque asks KIFF whether the window is open (`window` in
  `POST /v1/devices/verify`) and refuses a signature naming one that isn't,
  403 `window_required`. Safe calls (disconnect, lists, the export, deleting
  a trigger) keep the device key. `PLANQUE_WINDOW_REQUIRED=true` refuses a
  device-key signature with no window on a routine call; off, it is taken as
  before. Each approval records `key` and `window`.
- A browser as a connection (`"kind": "browser"`, kiff/planque#86): a taf
  uses websites in a browser planque runs for it alone, with planque's own
  tools (`browse_open`, `browse_read`, `browse_act`, `browse_checkout`,
  `browse_pay`) and reviewed profile, optionally limited to listed sites.
  Money moves only through `browse_pay`, whose amount KIFF reads from the
  checkout page and holds for the owner; a step that would pay another way,
  or type a password, card number or code, is refused before it runs, and
  secret fields are masked in every page. Offered only where a browser
  worker is configured.
- The live sign-in view of a browser grant: the owner signs in to a site
  themselves, on the phone, in a view of the grant's own browser that
  planque serves (`POST …/grants/{id}/sign-in`, signed, answers a one-time
  URL good for ten minutes). The page shows the browser as frames and
  sends the owner's taps, typing and keys to it; a URL must be on the
  grant's sites; nothing typed is seen by any agent and no frame is kept.
  Done seals the cookies as the grant's sign-in; a view that lapses keeps
  nothing; the taf's calls wait while a view is open (`signing_in`). The
  worker takes input and frames for it, holds at most three sessions at
  once (`PLANQUE_WORKER_MAX_SESSIONS`), and answers "busy" beyond that.
- The browser worker (`cmd/browserworker`, its image in
  `deploy/browserworker`): a headless Chromium planque drives over its
  debugging pipe, one browser and one egress proxy per session, gone when
  the session ends. The proxy resolves each name itself and connects only
  to public addresses and the grant's sites, over https; WebRTC is taken
  away from every page, so no UDP leaves the browser; downloads are
  denied. The worker reads a page's facts in no language (field kinds,
  controls that place an order, a checkout's total) with secrets masked at
  the source, and a checkout's lines (item, quantity, price). planque
  reaches it with `PLANQUE_BROWSER_URL` and `PLANQUE_BROWSER_KEY`. A card number shown in a page's text, or echoed
  in a plain field, reaches the agent as its last four digits only.
- Any app on demand (kiff/planque#87): a catalog search lists every app the
  managed provider has, each as it stands (`state: "needs_setup"` for one
  that needs our own developer credentials first); the unsearched list keeps
  to the apps with a sign-in set up. Picking an app the provider signs
  people in to by itself sets its sign-in up at that moment, once, with
  nothing of ours in it; the grant pins it and is `on_demand: true`, as is
  every later grant on it. At most ten such set-ups an account a day. Every
  grant view says who holds the sign-in (`custody`). A catalog read that
  stops with pages left says `partial: true`.
- Tools found by words across the provider's apps
  (`GET /v1/accounts/{tenant}/catalog/tools?q=…`), each with its app as it
  stands, and an app's events listed before and after a connection
  (`…/catalog/{entry}/events`, `…/grants/{id}/events`), so Coach proposes
  apps and events before connecting any (kiff/planque#87).
- A Telegram bot as a connection: an API-key grant whose `api.url` is
  `https://core.telegram.org/bots/api`. The bot's token is checked with
  Telegram before it is kept, and planque sets the bot's webhook with a
  secret of its own. People's text messages to the bot (private chats, and
  group messages that mention it) go to tafs as `planque.door.v1` with
  `"app": "telegram"`, once per update. The taf answers with `send_message`,
  a governed call that splits text over 4,096 characters. Disconnecting
  deletes the webhook.
- A tool the owner approved that changes at the vendor pauses its grant.
  For a direct MCP connection, planque reads the server's tools again
  before a call, at most every 10 minutes, and compares each approved
  tool's input schema with the one approved (a canonical hash kept with
  the tool from now on; a tool approved earlier is given one at its first
  read). A tool gone or changed pauses the grant. For a managed connector,
  a definition the provider answers that differs from the approved one
  does; one that can't be read or checked pauses nothing. A paused grant
  sends nothing: every call is refused with the code `operation_changed`,
  its `status_words` say "Paused · …", and an event records it. The
  grant's `status` stays `connected`: an app shows the pause through
  `status_words`, and telling the owner in KIFF comes later. A pause stays
  on that grant, through a new sign-in too; the owner connects again,
  approving the current definitions. A tool list planque can't read
  pauses nothing: the call goes on to its usual handling. REST grants are
  unchanged: they send only from the approved snapshot.

### Fixed

- The browser connection refuses every form a browser reads as an address,
  not only a dotted one: a short or hexadecimal one (`127.1`, `0x7f.1`),
  one with a leading zero, or one written with another script's full stop.
  A site's name in another script is refused; its `xn--` form opens.
- A connection an account key asked for (an app's own Connect apps screen)
  can be finished: the taf it is for and the account's Setup agent (KIFF
  role `setup_reader`) may list it, read its catalog and propose its tools.
  No other agent may, and the owner still approves the tools on the phone.
- The connection service's sign-in page moved to a new host; planque accepts
  it. Each sign-in attempt names its connection afresh, so an approval that
  failed after the service made its link can be approved again.
- Connecting a managed app works again. The hosted session no longer sends
  a workbench setting the connection service now refuses.

### Changed

- Managed API-key and bearer-token apps explain that the owner enters their
  credential on the connection service's page and that the service keeps it.
  New grants bind the enabled authentication config and its scheme to the
  device approval and hosted session. Existing grants and direct key entry
  keep their previous behavior.

- An enrolment keeps the enrolled connector key while KIFF still takes it for
  the account, and answers `200 {"kept": true}`: a second device setting up no
  longer swaps the key, which left the tools the first key connected unowned
  in KIFF. A new key replaces the enrolled one only once KIFF no longer takes
  it.

- A broad tool narrowed to one operation (Stripe's `stripe_api_write`) takes
  only the operation's fields: planque reads them with the vendor's
  describing tool, once the owner approved it on the grant, and refuses a
  field the operation doesn't take, with the list it does. A narrowed broad
  write is added to a grant after its describing tool is approved; asked for
  at once, it is refused. Narrowed reads are unchanged.

### Fixed

- Managed discovery accepts app and tool slugs beginning with an underscore.
  Unsupported catalog entries are counted and skipped without losing valid
  entries or the next-page cursor; direct lookup and execution stay strict.

- A read of planque's store is sent again when the connection drops before an
  answer, so one dropped connection no longer fails the owner's request. Writes
  are never sent again.
- An operation's outcome is always recorded, even when the vendor's answer is
  too large to keep; the answer is then kept as a short stand-in, and a replay
  still never sends again.

### Added

- Device-approved Slackbot event subscriptions bind exact accounts, versioned
  definitions and allowed conversations. Signed project webhooks feed a fixed,
  signed door receiver, with durable replay suppression, transient delivery
  retries and cancellation on disconnect. Oversight exports include the
  approvals; account deletion removes the associated event records.

- An optional Composio connection provider reuses upstream app discovery and
  hosted authentication. Device-approved grants pin each selected account,
  concrete tool version and schema; exact-call permits and once-only operation
  records govern execution. Discovery is paginated, project credentials are
  resolved server-side, and disconnect distinguishes disabled access from
  pending upstream credential revocation.

- Reviewed REST definitions can name a version header. Profile approval pins
  its value in the grant and every request; callers cannot replace the header
  definition or overwrite authentication with it.

- An operator-installed, versioned vendor catalog lists account connection
  state and requests reviewed operations for a chosen taf. Catalog grants
  freeze their profile and final tool names before device approval. API-key
  entries refuse a source description changed since review.

- **Device-approved connection profiles:** a grant's digest covers its
  profile, and `/approve` and `/key` accept the phone-signed profile/words
  digests. Tools use the profile's measurements, including decimal-string
  amounts. The relay waits for KIFF to confirm the profile; a failed
  application rolls back registration. Required arguments and permit-derived
  markers are checked before sending, and operation records keep the profile
  revision used for their replay.

- **Registered vendor applications:** operator-configured MCP and REST
  applications use credentials from a sealed store or configured Secrets
  Manager references. The owner's approval pins the client identity,
  endpoints, access descriptions and reviewed operations. Secret rotation
  for that client works without changing the approval; replacing the
  client identity needs a new approval.

- **MCP client metadata documents:** a vendor that supports CIMD can use
  planque's public HTTPS client document instead of a registration request.
  The document contains only the configured client identity and redirect.
- **Larger API descriptions and selected operations:** descriptions may be
  up to 16 MB. `api.operations` selects up to 500 operations from a larger
  description; the selection and the complete source digest are included
  in the request the owner approves. An operation fingerprint distinguishes
  request changes from wording changes.

- **A write's amount read from its resource:** a tool may name `amount_lookup` (`{tool, argument,
  path}`), a read-only tool of the same request, the read's string argument that takes the call's
  `resource_argument`, and where the amount is in its answer. planque registers it in KIFF like
  `unit_lookup`, by the read's name there; KIFF reads the amount before deciding, and the call
  carries none, so a Card counts the amount the vendor says, never one the agent sends. It needs
  `amount_argument` and `resource_argument`, and the amount can't also be one of a narrowed tool's
  arguments.
- **API-key grants** (kiff/planque#64): a vendor with a REST API and an OpenAPI description, and no
  MCP server, is connected with a key the owner gives on their phone. The description is read
  once into a snapshot the owner approves by its digest; each operation is a tool, called through
  the relay with KIFF's permit like any other. A 4xx is the vendor's refusal; a 5xx, a redirect or a
  lost answer is unknown and never sent again. The key is sealed, never shown again, and deleted on
  disconnect.
- **A vendor's refusal is never evidence that nothing happened:** an API-key write that the vendor
  refuses answers the agent with the vendor's status and words, and KIFF keeps it counted. No tool
  is registered with an outcome contract: a status alone can't say nothing ran (a batch endpoint
  can answer 400 after part of it did). This replaces the brief release on a fixed list of 4xx
  statuses.
- **Reusing a sign-in** (kiff/planque#36): `reuse` names a connected grant whose sign-in or key a
  second grant uses, approved with no second sign-in.

- A vendor's sign-out is told to KIFF once, when the grant turns
  `needs_reconnect`, so the people who see the taf hear of it at once. A
  notice KIFF didn't take is sent again by the hourly sweep, for up to a week;
  one whose grant was signed in again is dropped.
- Every grant carries `status_words`: its status in the owner's words, so the
  apps show the same line for a connection and keep no words of their own.
- A connected grant gains tools with one approval and no second sign-in,
  while its tools keep working (kiff/tafs#138). The asking agent lists the
  grants it asked for a taf's agent, to add to them.
- Deleting the account, signed from the owner's phone, erases everything
  planque keeps for it: tools out of KIFF, vendor tokens revoked, records and
  the enrolment deleted (kiff/tafs#185).
- Security reporting, contribution and agent guidance, issue and pull-request
  templates, code owners, and dependency-update configuration.

Earlier history is available in [merged pull requests](https://github.com/kiff/planque/pulls?q=is%3Apr+is%3Amerged)
and [releases](https://github.com/kiff/planque/releases). This file does not
retroactively assign release dates or claim that merged work was deployed.
