# Changelog

User-visible changes and compatibility notes are recorded here. Dates and
version headings refer to releases, not merely to merged pull requests.

## Unreleased

### Added

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

>>>>>>> 72f2ff1 (One fingerprint, then taps: a routine call signed inside an approval window (kiff/tafs#494))
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
