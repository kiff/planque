# Changelog

User-visible changes and compatibility notes are recorded here. Dates and
version headings refer to releases, not merely to merged pull requests.

## Unreleased

### Added

- **A managed app's tool is taken by its provider's name too.** Its catalog lists each tool's `slug` (`GMAIL_FETCH_EMAILS`) beside the name planque gives it, and setting tools takes the slug as well as the name. An app's guide and a search name the slug, so Coach no longer looks a tool up by a name the catalog doesn't hold.
- **A sign-in never asks for personal data the taf doesn't use.** When a managed app's usual sign-in asks for personal data (Google's date of birth, addresses, phone numbers, other email addresses, contacts, language) and the sign-in can't be narrowed to the taf's tools, or no tools are named, planque asks for the app's access less that data, as a sign-in of its own, recorded in the signed request's `scopes`; the card still says it's more than the taf uses. A personal scope stays only when a tool asked runs on nothing else (finding people, reading contacts), the narrowest read-only one; a tool added later that needs one is refused as not covered. When no such sign-in can be set up (the account's set-ups for the day used, the provider failing), nothing is asked: 429 or 503 with why, never the usual sign-in. The account's own address stays. Existing connections keep what they have.
- **A Stripe connection recorded by its managed app's id is called Stripe too.** A connection to Stripe's or Qonto's own server whose app was recorded as the managed app (`connector:<provider>:stripe`, from the reviewed server list) read as its server's own name ("stripe-mcp"); it now reads "Stripe" (or "Qonto"), on the next read.
- **A vendor's own server is called by its app's name.** A connection to a reviewed app's server (Stripe's, Qonto's, or any reviewed catalog entry's) is named by that app ("Stripe"), never by the name its server gives itself ("stripe-mcp"). Connections stored with the server's name read as Stripe and Qonto too. A server with no reviewed app keeps its own title, else its name.
- **A managed app's approval card is shorter and in plain words.** It names no internal service ("the connection service"): it says who keeps the sign-in in plain words ("Your Gmail sign-in is kept by our sign-in partner; the taf never sees it"). When the app will ask for more access than the taf uses, it says why (its sign-in can't be narrowed to this taf's actions), and the card ends saying the taf gets only the actions approved. It names no tools by the provider's names ("Fetch message by thread id"): the actions come next, by the titles the taf's request gives them. Details no longer show the taf's agent id (it stays in the signed request). In the 8 languages. A grant asked before keeps the words it was asked with.
- **A managed app's tools that can't be undone are never served.** Per app, in `connector.NeverServedTools`: for Stripe, the invoice void, delete, finalize, pay and mark-uncollectible tools, and voiding a credit note or credit grant. They never come back from the tools search or an app's catalog; asking for one, or setting it on a grant, is refused with `tool_not_served`; and a connection holding one from before no longer lists it, nor sends a call to it. Platform meta tools are refused as before.
- **The account's connections name their taf.** Each connection in `GET /v1/accounts/{tenant}/connections` carries `agent_name`, the taf's name its request gave, so tafs cloud can tell a removed taf's connection from another agent's whose id only looks like it.
- **Each connection is a row the apps draw as sent.** The account's list (and each grant's answer) carries `used_by` (the tafs on it, by name), `health` (`state`, `ok`, plain `words` in the language asked, its time in the owner's `X-Time-Zone`: "Working · last used 00:32", "Gmail's sign-in expired on 14 Nov") and `actions`: only those that can work for it now (`sign_in_again`, `remove` with what stops, `remove_unfinished`). An unfinished connection offers only to remove it, never Sign in again, so no button answers with a state like "this connection is listed". A connection's last use is kept beside the grant, from each call the vendor answered with no error. In the 8 languages.
- **A managed app's approval names its actions one way.** The approval names each action by the title the taf's request gave it (the owner's words, as the card's lines do), and a reused sign-in's text no longer lists the held sign-in's own tools as if they were the new taf's.
- **A connection's name for people, and an unfinished one said so.** A grant's view and the account's connections carry `title` ("Gmail · Demo Shop") and `since` (the day it was asked for), so every app names a connection the same way. A connection signed in with no tools chosen reads "Unfinished · signed in, no tools chosen yet", in the 8 languages, and still ends after its day.
- **Ending a taf names the connections that keep a shared sign-in.** `kept` in `…/agents/{agent}/end` now carries `grants` and `agents`: the live connections still holding the sign-in and their agents, so an app can offer to remove it too, by disconnecting those, when none is a taf the owner still has.
- **A managed app's connection is named by its account, never the provider's id.** Apps name a connection by `title`, which carries the account name the provider shows ("Demo Shop", an address), never the provider's id for the connection (`ca_…`). `account` is unchanged: the apps send it back for a connection's doors and events, and planque checks it there.
- **An app's guide, served.** `GET …/catalog/{entry}/guide` answers the app's adapted guide (`deploy/guides/<app>.md`, built into planque) for an agent of the account, before any connection, so Coach reads it before it offers or connects the app. An app with no adapted guide yet answers 404, with words.
- **A key grant declined before its key withdraws its declared tools at once.** `DELETE …/grants/{id}/tools` on a key grant whose tools are declared, before the owner gave the key, takes them out of KIFF and ends the request (`failed`, declined), instead of leaving them declared until its deadline.
- **A guide for every app, generated from the catalog.** `planque guides` reads every app that can be signed in to now, each one's tools, and the reviewed entries, and writes each app's generated part to `deploy/guides/generated/`: how it connects and, per tool, read or write, inputs, `returns` and the access it needs. `apps.md` lists every app in batches, the first ones first, for the adapted parts written by hand beside them (limits, plans, and what to tell the owner on an error). Gmail, Google Sheets, Stripe, GitHub, LinkedIn and Qonto are adapted first ([guides](deploy/guides/README.md)).
- **A managed app offers only the tools its sign-in lets run.** A tool whose scopes the sign-in wasn't granted (a search needing Drive access on a Sheets-only sign-in, say) is left out of the catalog's `tools` and listed under `unusable` with why, and choosing it is refused with the same words: no Card grants a tool the vendor would refuse with "insufficient scopes". A reused managed sign-in's text names its account, so an owner who holds two can tell them apart.
- **A catalog tool says what it gives back.** `returns` lists the ids and links in its answer (`data.spreadsheetId`, `data.url`), read from the vendor's description of it, for a managed app's tools and an MCP server's that describe their output: what a taf keeps to use next time.
- **A signed-in grant's first tools are declared in KIFF too,** as an addition's are, so its tools and the taf's Card take one approval.
- **A managed app the account already holds is reused.** Asked for by its entry, an app with a connected sign-in that covers the tools asked answers a connection on that sign-in, `listed`, with no browser and no new sign-in; it gets only the tools approved for the new taf. `new_sign_in: true` asks for one of its own.
- **Setup reads named with the sign-in.** A catalog entry connected first (Stripe's, Qonto's) takes `setup_reads` at its first request: approved with the sign-in, they're connected once signed in, each listed and read-only, with no second approval. A server's request can't name them: it connects first, then asks for its reads with its tools, as before.
- **Tools added to a connected grant are declared in KIFF while they wait.** Asked for, they're registered in KIFF at once, with a transport credential of their own, so the taf's Card can be drafted over them and the owner approves both with one signature. They stay unusable until approved: each call is refused, saying it waits for the owner's approval and nothing was sent. Approved, the same registration serves them. Declined (`DELETE …/grants/{id}/tools`, new), past their day, revised, or with the grant disconnected, they leave KIFF, and the grant keeps what it had; a grant's first tools, declined, leave it listed with its sign-in held.
- **A refused request's tools say what to fix, by code.** A 400 refusing a connection request's agent or tools now carries, beside the same `error` words, a stable `code` (such as `pin_shape`, `pin_unknown`, `required_argument`, `agent_required`), the `tool` and `field` at fault, and `want`, the form planque wants there, so Coach can say in plain words what to change. The connection-refused log names the code ([API](docs/api.md)).
- **A big API is described by words.** `api/operations` takes `q`, a few words from the job, and answers the matching operations, at most 50, best first, with `total` and `matched`. With no words, a description of up to 200 operations is listed whole. A bigger one (Stripe's 608, GitHub's 1,277) answers its count and how to ask, never every name. Listing reads past a snapshot's 500. A key request for a big API that names the job's operations takes them alone as its snapshot, so it's declared; a whole request is still at most 500. Such a grant holds only the operations its request named: one more is a new request naming it. Operations rank by how many of the words their name holds, then their title, so no summary stuffed with the job's words ranks above the operation a job names.
- **An API's operations, read before any grant.** `POST /v1/accounts/{tenant}/api/operations` with `{"api": {"url", "auth"}}` answers the vendor's public description as planque reads it: its operations' names, titles and read or write, plus what was left out and why. It reads the description document only: no key, no call to the vendor's API, nothing of the account's, and it keeps nothing. Coach names a key request's tools by these names, so the request is declared and the key and the taf's Card come in one approval. An account reads 100 descriptions a day.
- **Whether the managed sign-in works for an app, for the connect matrix.** `GET /v1/catalog/{entry}/sign-in` answers `{entry, state}`, with `state` one of `ready`, `on_the_owners_pick`, `needs_setup` or `not_at_provider`. It reads the provider's app metadata only: it makes no sign-in, starts no set-up and takes no allowance, and the provider's own ids aren't in the answer. Only the key whose SHA-256 is `PLANQUE_SIGNIN_READER_SHA256` reads it, and that key reads nothing else. Without the setting, the route isn't there.
- **A key by HTTP Basic.** A grant's key can go as HTTP Basic (`"in": "basic"`): the password, with the username given beside it (`"user": "me@example.com/token"`), or the username, with a password given (`"password": "X"`) or none. It matches a description's `http` `basic` scheme. The owner reads where the key goes, and the username, in the approval, in eight languages. A username with a colon or a newline, over 200 characters, or given with a password is refused, as is a key with a colon used as the username, and a version pin that would replace `Authorization`.
- **A description that names no security scheme takes the key as asked.** When a description names no scheme at all, the key goes where the request says, from the vendor's docs, and the owner reads where in the approval. An operation that names a scheme is still refused.
- **Form-encoded bodies.** An operation whose body is `application/x-www-form-urlencoded`, and only that, takes its fields as a form (`"body": "form"`): strings, numbers and lists of them, never nested. An operation that offers JSON too is still sent JSON.
- **An API at each account's own address.** A server whose URL has a variable with no default (`https://{subdomain}.vendor.com/api`) no longer refuses the description. It's kept as `accounts`, and a request with no `base` is told the API is at each account's own address, so Coach asks the owner for theirs. The owner reads the base in the approval, as before.
- **One fixed header beside the key.** `auth.fixed: {name, value}` is a header every call carries, as a vendor asks beside the key (`Harvest-Account-Id`). An operation that needs it as a parameter has it, and it's never an argument. It can't be `Authorization`, `Host`, `Cookie`, the request's own headers, a proxy's or a forwarded header, or the key's own header. Its value is one line of at most 200 characters, and a version pin can't replace it. The owner reads it in the approval, in eight languages.

- **API descriptions in YAML.** A grant by key reads a vendor's OpenAPI description in YAML as well as JSON, since many vendors publish YAML only. YAML is read as the JSON it stands for, then checked and limited as JSON is. Only plain data is read: one document, with no anchors, aliases, merge keys or tags beyond YAML's own, and an unquoted `200:` key reads as `"200"`. A YAML file within 16 MB that's over 16 MB as JSON is refused.
- **A refused connection is logged with its reason.** Each refusal of `POST …/grants` or `POST …/catalog/{entry}/grants` is logged as `connection refused`, with the route, the entry, the status and planque's words. A sign-in that fails to set up at the provider is logged as `app set-up failed`, with the provider's own reason, its secrets redacted. The request's body is never logged.

- **A reviewed entry is asked for with only the tools the taf uses.** `POST …/catalog/{entry}/grants` on a reviewed entry (such as Qonto by key) takes `tools`, the entry's tools the taf will use, as the tools search names them. The request, what the owner approves, its profile and what's declared in KIFF are those alone. Before, the body took only `agent` and `agent_name`, so Coach's `uses`, sent as `tools`, refused the whole request. A managed app's entry already took `tools`. A tool the entry doesn't have is refused, naming the ones it has. Left out, every tool of the entry is asked for, as before.
- **A key grant's tools are declared to KIFF when it's asked for.** A grant by key, asked with its tools, has them registered in KIFF at once, before the owner gives the key. A Card can then be drafted over them, and the owner signs the key and the Card together. Nothing is usable until the key is in: the relay refuses every call with `waits_for_key`, and nothing is sent to the vendor. Once the key is given, the same registration serves the grant, with its credential unchanged. A request cancelled, ended by its 7-day deadline, or ended with its taf takes its declared tools out of KIFF. A grant's view says `declared: true` while its tools wait for the key. If the declaration fails, the request connects as before, once the key is given.
- **A removed taf's connections end with it.** `POST /v1/accounts/{tenant}/agents/{agent}/end` ends every connection of a taf's agent, whatever its status: connected, waiting for the owner, or failed. A sign-in or key another taf's live connection uses stays, and the answer says which and for whom. `?preview=1` ends nothing and gives the same answer, so the app can say before the fingerprint what goes and what stays. Only a person who manages the account's tafs may call it, with their own session; no key ends anyone's connections. An agent with nothing to end is 200, never 404 ([README](README.md)).
- **The reviewed catalog kept as data: the verified tier.** Catalog entries can now live in planque's table, as revisions an operator drafts and reviews with an API credential, and another operator publishes or retires with an enrolled device's signature over the entry's JCS digest (5-minute window, single-use nonce). Who drafted or last changed a revision can't publish it. A key entry with a write that one of its arguments can make final isn't published until that argument is pinned to its undoable value. planque serves a revision only while its stored digest is the one signed, and keeps each publish and retire as a record. `PLANQUE_CATALOG_OPERATORS` names the operators. No deploy is needed to add an entry, and shipped entries are served as before ([connections](docs/connections.md#the-catalog-kept-as-data)).
- **Qonto by key, with invoice drafts from accepted quotes (`qonto-key@2`).** A new reviewed catalog entry: `@1`'s four reads, plus the quotes, one quote and the supplier invoices as reads, and one write, `qonto_draft_client_invoice`. That's Qonto's invoice create with `status: "draft"` fixed on the grant, so it makes a draft that goes to no one. Nothing on the entry finalises, sends, pays, cancels or deletes. Each draft names its quote (`quote_id`, the call's resource), its lines can't carry a discount, and its IBAN must be one of the business's own accounts, read with the same key before each call (`one_of`). The owner approves `@2` anew; `@1` and its grants are unchanged ([connections](docs/connections.md#qonto-key2-reads-and-invoice-drafts-from-quotes)).
- **Array arguments on a narrowed tool.** An argument may be `"type": "array"` with a reviewed `schema` (objects whose fields are all named), so a write that takes a list, such as an invoice's lines, can still have arguments fixed beside it. A call whose list has a field the schema doesn't name, the wrong kind of value, or a required field missing is refused before it reaches the vendor.
- **Stripe refunds for some accounts only.** `PLANQUE_STRIPE_REFUNDS=tenants:<id>[,<id>]` offers the reviewed refunds profile to those accounts alone; any other account is answered as if it were off. `on` still means every account, and anything else means off. Live mode stays refused in the code, whatever the list says ([managed connectors](docs/managed-connectors.md)).
- **Coach finds the reviewed catalog entries in its tool search.** `GET …/catalog/tools?q=…` now returns, on the first page, the tools of the entries in planque's reviewed catalog that a word of the search names (three letters or more), with the entry's id and state: first when a word names the entry's app ("qonto"), after the provider's apps when a word names only one of its tools ("invoices"), since a reviewed entry is one app among those a person may use. Before, only the managed provider was searched, and Coach had no other way to find an app, so an entry such as Qonto by key could be connected from Connect apps but not from a conversation. When `apps` names apps, only the provider is searched, as before. If the provider's search is down and a reviewed entry matched, the answer is 200 with `partial: true` and the entries found; with no match, it is 503 as before.
- **A sign-in to reuse is checked first.** Before a grant is offered with "no second sign-in" (\`reuse\` of a connected OAuth sign-in to an MCP server), planque checks the sign-in still works: a live token, refreshed when due, and one read of the server's tools, refreshed once if refused. One the vendor refused (signed out) isn't offered: the request signs in anew to the same server, on its first card, and its record says why; the refused grant is signed out as any refused refresh is. Before, the owner approved "no second sign-in" and the grant failed. A vendor that can't be reached isn't a refusal, and a key or a managed app's sign-in is checked where it's used, as before.
- **A grant's words, in the owner's language, and signed as shown.** A grant's `text` and `status_words` are written in the eight languages tafs ships, and a view is in the language its app asks in (`Accept-Language`), with `lang`. English stays the default, and a grant made before reads in English. `/approve` also takes `digest_v` 2, a digest over the request with the words shown and their language, and the approval keeps those words and their language as the receipt. Views carry `digest_v` 2 only with `PLANQUE_DIGEST=v2`, once both apps check it ([connections](docs/connections.md#the-words-in-the-owners-language)).
- **A vendor that refuses planque's client says so.** When a vendor's authorization server refuses planque's client registration (Qonto allow-lists MCP clients' redirect URIs), approving answers 502 with `code: "vendor_refused_client"` and words the apps can show, instead of a bare 502. One `client registration refused` log line gives the grant, host, status and OAuth error code. Nothing is connected, and the grant stays pending. A vendor that answers 408 or 429 hasn't refused anything: that stays the retryable 502.
- **Qonto, signed in to directly.** The catalog lists `qonto`, a direct sign-in to Qonto's MCP server (`https://mcp.qonto.com/mcp`), where the owner signs in on Qonto's page with Qonto's own approval step. The sign-in asks only for reviewed scopes: reading the business, its clients, its invoices and their e-invoicing state, and writing clients and client invoices. It never asks for transfers, cards, bank accounts, direct debits, team or members, though Qonto's server lists them, and it fails closed: the server being connected, or any resource on Qonto's host, gets the reviewed list. The managed connector's Qonto is no longer listed, found in a tool search, or connected (409, naming `qonto`): one way to connect a bank. `deploy/server-apps.json` lists Qonto's server under `qonto`.

### Changed

- **A merge to `main` deploys planque.** CI's `deploy` job builds the function as `deploy/aws.sh` does, replaces its code, checks it live (`/healthz` 200, an unsigned catalog read 401) and puts the previous code back if it doesn't answer so. Off until the repository variable `PLANQUE_DEPLOY` is `on` and the secret `AWS_DEPLOY_ROLE` is set.
- **One operation planque can't call is left out with its reason, never the whole description.** An operation that sets its own servers (or whose path does), has callbacks or links, is secured otherwise than the key, holds a reference that names nothing, or repeats another's name is in `skipped`, with why; the others still read. Only what spoils the whole (OpenAPI 2, an external reference, its size) refuses it. A description with nothing left to call names the first reason.
- **A description's digest is taken over its normal form.** `digest` is the SHA-256 of the description as JSON, with keys in order and numbers as written. So a description and its YAML twin, or the same JSON indented otherwise, have the same digest. It changes for a description read from now on. A grant already made keeps the digest it was approved with, and a reviewed entry is still checked against the bytes it was reviewed on.
- **A request nobody signs ends after 7 days.** A pending request (no sign-in held, no tools connected) now has a deadline, 7 days after it was asked. Past it, the sweep, or the next read, ends it as `failed`, "not approved in time", and it can't be approved anymore. Before, it waited for ever, outliving its taf. Reading it doesn't move the deadline. Tools added to a connected grant keep their own day.
- **A refused body names the field planque doesn't take.** "the body must be a JSON object with the documented fields" now ends with `unknown field "…"`, so a caller fixes its request on the first retry.
- **A server's metadata may name its origin.** The resource in a server's metadata must be the server planque calls, or the same scheme and host with a path made of its leading segments: Qonto names `https://mcp.qonto.com` for its server at `/mcp`. The sign-in is bound to the resource named. Another host, a port, a query or a sibling path (`/mcpx` for `/mcp`) is still refused.
- **Stripe's own catalog entry carries Stripe's logo**, as the managed apps do.

- **A Stripe connection's account comes from Stripe's own accounts tool.** Stripe's MCP sign-in names no account in its token answer (`stripe_user_id` is Connect's). Right after the owner signs in, planque calls `list_available_accounts_or_orgs` once, itself, with the new token, and keeps the account, with its mode, when it lists exactly one (`oauth.Metadata.AccountsTool`, set for Stripe's server only). With several or none it keeps none, and the read is on the connection's record. A refund connection pins that account, and is refused when the account is in live mode.
- **Stripe refunds are offered in test mode only.** A refund connection asked with `livemode: true` is refused (409) with the reason: in live mode, Stripe may ask the owner to confirm each refund in its Dashboard, which hasn't been checked yet. In test mode, a refund through the direct sign-in went through at once.

### Fixed

- **A tool's `returns` follows its answer's schema where it points.** The provider describes a call's answer as a schema generated from typed models: the answer's `data` is a reference into `$defs` (its docs show it with no type), and an optional field is `anyOf` with null. `returnsOf` read only inline `properties`, so every managed tool's `returns` came out empty. It now follows `$ref` within the schema (`$defs`, `definitions`), `anyOf`, `oneOf` and `allOf`; a model that contains itself is read once, and a reference out of the schema is never read.
- **A refused new sign-in says which account to sign in with.** A reconnect refused because its sign-in is to another account, or can't be shown to be this connection's, answered a bare 409, and both apps drew it as "It changed in the meantime", which looped the owner back to the same step. It now carries a `code` (`other_account`, `account_unknown`) and `words` in the owner's language, naming the account the taf uses ("You signed in to Gmail as …, which isn't this connection's account, so nothing changed. Sign in again as …, the account this taf uses."). A connection with no account name kept says at approve to remove access and connect again. Each refusal is logged with the grant and the code, never an account. In the 8 languages.
- **Sign in again on a signed-out managed connection, itself.** Approving a `needs_reconnect` managed connection again opens the provider's sign-in for that same grant; confirmed on the same vendor account (as the provider names it), the grant is connected again with its tools, its registration in KIFF and the taf's Card as they were, with no new Card and no Coach round. A sign-in to another account connects nothing: it's revoked, and the connection stays signed out (409). It fails closed: the provider shows no vendor id, only the account's name, so a new sign-in that comes back with no name is refused the same way, and a connection whose account never had a name isn't offered a new sign-in (409, saying to remove access and connect again). A call the vendor refuses with 401 answers in plain words ("…taps Connect again on this taf"), never the vendor's error body, which stays in the log.
- **A held sign-in the vendor no longer takes is never reused, nor shown as connected.** Before a managed app's held sign-in is reused for another taf, the provider is asked about it again; one it doesn't call active, or that can't be asked about now, isn't reused, and the owner is asked to sign in instead of reading "no second sign-in". A managed call the vendor refuses with 401 (a sign-in its refresh can't renew, behind a status the provider still calls active), or one the provider calls inactive, marks every connection on that sign-in `needs_reconnect` ("signed out", sign in again), and KIFF is told, as for a refused refresh.
- **Setup's read Card for a connection approved with its Card.** When the owner approves a connection together with the taf's Card under one fingerprint (`kiff.requests.v1`), the approval's message names the request by its hash. planque now keeps that request's own message (`approval.request`) and sends it with `POST /v1/me/setup-reads`, so KIFF can check what was approved and issue Setup's Card. Before, KIFF refused it, and Setup's first read was blocked although planque counted the read connected. The profile approval planque forwards (`POST /v1/me/profile-approvals`) carries the request too, so a profile approved with its Card applies. KIFF must take `approval.request` first.
- **Setup reads named at a key app's request connect with its key.** As after a sign-in, a key app's first request that names `setup_reads` connects them once the key is given, and Setup's Card is asked of KIFF. Before, only a sign-in did.
- **A narrowed argument keeps the values the vendor allows.** When planque narrowed a tool, the agent saw each argument as a bare type, so the vendor's `enum` (a desk's `open`, `pending`, `solved`) was dropped, and a guessed value (`new`) got an empty answer that read as "nothing to do". The argument now carries the vendor's `enum`, and a call with a value outside it is refused before anything is sent, naming the values it takes. Grants connected before keep their schema until their tools are connected again ([connections](docs/connections.md)).
- **A narrowed read called with no arguments still sends its arguments' object.** A grant that lifts only some of a tool's arguments into an object (Stripe's `stripe_api_read`: `limit` → `parameters.limit`) sent no `parameters` object at all when the agent set none of them, as a schedule's check does with `{}`, and the tool answered with an error, so the check never read. That object is now always sent, empty when nothing is set. Set arguments and fixed ones are unchanged.
- **A refund connection's account read takes the account at `parameters.id`.** `GetAccountsAccount` takes its account as `id`, as a test-mode read showed; it was pinned at `parameters.account`, which the first-call field check would have refused, so no refund could have been decided.

### Changed

- **A refund connection's case reads are checked against Stripe's description.** A narrowed read is now checked at its first call, like a narrowed write, where its connection holds the describing tool (`stripe_api_details`): a path its operation doesn't take is refused with Stripe's own list. Without that tool, a narrowed read isn't checked, as before.
- **A case needs a fixed currency.** A profile's `case` block is refused unless its operation has a fixed `amount.unit`: KIFF compares the charge's currency only with a fixed unit.

### Fixed

- **A field check on an ended sign-in says to sign in again.** When the vendor had ended the sign-in planque holds, the first-call check of a narrowed tool refused with `definition_unread` ("couldn't read the operation's fields"), and the taf told the owner to fix its tools. It now refuses with a new code, `signed_out`, in the words for an ended sign-in: the owner signs in again. The tool stays unchecked until a call after the new sign-in.
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
