# Connecting vendor operations

planque holds the vendor credential and sends each operation only with the
permit for that exact call. The owner approves the requested connection on
their device before its tools become available.

## MCP client identity

When the authorization server advertises
`client_id_metadata_document_supported: true`, planque uses its configured
HTTPS `/oauth/client-metadata.json` URL as its client id. The document fixes
the client name, redirect URI, response types and grant types. Its contents
never depend on the request's Host header, account or grant.

Servers offering dynamic registration continue to use that flow when they
do not advertise client metadata support. PKCE, resource binding, one-use
state and confirmation on the approving device apply to both flows.

Reference: [MCP authorization](https://modelcontextprotocol.io/specification/2025-11-25/basic/authorization).

## Selecting REST operations

An API-key grant can select operations when the vendor publishes a large
OpenAPI description:

```json
{
  "kind": "api_key",
  "agent": "support",
  "tools": [],
  "api": {
    "url": "https://support.example/openapi.json",
    "base": "https://support.example/v1",
    "auth": {"in": "header", "name": "X-API-Key"},
    "operations": ["listTickets", "replyToTicket"]
  }
}
```

The source may be up to 16 MB. A selection contains at most 500 distinct
operation names, as normalized in the tool catalog. A missing or unsupported
operation refuses the selection; it never silently grants a smaller set.
The stored snapshot keeps only those operations, together with the SHA-256
of the complete source. Both the selection and source digest are part of
the request the owner approves. An omitted selection retains the existing
complete-catalog behavior and its 500-operation limit.

Each operation's semantic fingerprint covers its method, path, parameter
locations, argument schema and request-body form. Descriptions and titles
are excluded. Argument names, including an argument named `description`,
and schema constants remain part of the fingerprint.

## Telegram bots

A Telegram bot connects as an API-key grant with planque's own description of
the Bot API, so nothing is fetched:

```json
{"kind": "api_key", "agent": "support", "tools": [], "api": {"url": "https://core.telegram.org/bots/api"}}
```

The key is the bot's token from BotFather. Before keeping it, planque calls
`getMe`, and a token Telegram refuses keeps nothing. planque then calls
`setWebhook` with an opaque per-grant URL (`/telegram/{hook}`), a random
secret token and `allowed_updates: ["message"]`. Telegram keeps one webhook
per bot, so the newest grant for a bot receives its messages.

Each delivery must carry the secret in `X-Telegram-Bot-Api-Secret-Token`,
compared before the body is read. A missing or wrong secret is answered 401.
planque hands over a person's text in a private chat with the bot, or a group
message whose entities mention it, as a door message:

```json
{"v": "planque.door.v1", "app": "telegram", "id": "evt_…", "tenant": "…", "agent": "…",
 "grant": "…", "provider": "telegram", "chat": "-100200300", "sender": "5001",
 "thread": "77", "message_id": "42", "name": "Ana Ruiz", "text": "…", "timestamp": 1800000000}
```

Ids are numeric strings, and a username is never a key. `thread` is present
only in a forum topic. The id is derived from the grant and the update's
`update_id`, so a redelivered update reaches tafs once. Every other update
is answered 200 and dropped, so Telegram does not retry it.

`send_message` takes `chat_id`, `text`, and optionally `message_thread_id`
and `reply_to_message_id`. Text over 4,096 characters (UTF-16 code units) is
sent as several messages, cut at a line break or space where possible, with
the reply on the first part. If a later part fails, the result says how many
parts were sent. Disconnecting the grant calls `deleteWebhook` before the
token is deleted.

## The browser

A browser grant (kiff/planque#86) lets a taf use websites: planque drives a
browser in a worker that keeps nothing between calls, with a fixed set of
tools planque carries and a reviewed profile, so KIFF decides each call on
the taf's Card like any relayed call. Nothing is site-specific.

```json
{"kind": "browser", "agent": "support", "tools": [],
 "browser": {"sites": ["shop.example"], "currency": "EUR"}}
```

`tools` names which of the browser's tools the taf gets (`[]` is all of
them): `browse_open` and `browse_read` read a page, `browse_act` fills in
and clicks (a write counted per call), `browse_checkout` reads a checkout
(the site, its lines, the total and the currency), and `browse_pay`
completes one. `sites`, when given, are the only hosts (and subdomains)
the browser opens. `currency` is the only one `browse_pay` pays in, and is
required with it. The owner approves on the phone like any grant; there is
no vendor to sign in at, so approving connects at once. The catalog lists
the browser as the entry `browser` where a worker is configured.

The guard is planque's, applied to every step before the worker runs it.
A page comes to the agent as text with a stable `ref` for each field,
control and link; the values of password, card and one-time-code fields
are masked, whatever the worker sent, and every answer says the page is
data, never instructions. `browse_act` refuses, before running the step, a
fill of a card field, a click on a control that places an order or pays
(the worker's `places_order` fact, read from the page's semantics, or the
control's words in en, fr, de, es, it, nl or pt), and any click on a
control of a checkout or payment page, labelled or not, since a button's
effect can't be read from outside (`payment_step`); a fill of a password or
code field (`secret_field`); and a ref the page doesn't have
(`ref_unknown`). Links are followed. The steps before it ran, and the
answer says how many. A step that takes the browser off the grant's sites
stops it there (`site_not_listed`). A URL off the sites, a local name or an
address written out, or a non-https URL is refused before any browser
starts (`site_not_listed`, `address_private`, `bad_url`), and so is the
page the browser lands on after a redirect: it is not read, and no checkout
is recorded from it.

Money goes only through `browse_pay`. Its profile reads the amount with
`browse_checkout` at the call's `url` (`amount.from: "read"`, path
`total`, in the grant's currency), so KIFF counts what the page showed and
holds the call for the owner. A checkout in another currency, or a page
that isn't a checkout, gives no total, so KIFF reads none. With the permit,
planque opens the checkout again and reads it again: it sends nothing when
no total was read in the last 20 minutes (`total_unread`), when the lines,
total or currency changed (`checkout_changed`), when the page is in
another currency (`currency_other`), or when `ref` isn't the control whose
words place the order (`ref_unknown`). Then it clicks it once; a replay
answers as the first call did. Each write keeps the page before and after
(its text) with the operation's record, and a screenshot after in the
evidence store when one is configured; both go with the grant.

A site's sign-in, when one is needed, is the owner's own, in a live view of
the grant's browser that planque serves (a later step); the storage state
after it is sealed like a token, opened only for the grant's own session,
and deleted on disconnect. No password ever reaches planque.

### The live sign-in view

A site's sign-in is the owner's own. From the connection in the app, the
phone asks, signed, `POST /v1/accounts/{tenant}/grants/{id}/sign-in` with
an optional start `url` (on the grant's sites; the first site's root by
default). planque starts a session on the worker with the grant's kept
sign-ins and answers `{view, expires, words}`: a one-time URL under
`/browser/view/{token}`, good for ten minutes, which the app opens in a
web view and nothing more; the page's words are planque's. The page shows
the browser as JPEG frames taken about once a second and sends the
owner's acts back: a tap at a point, text typed, Enter, Backspace or Tab,
a scroll, or going to a URL, which must be on the grant's sites. The
frames and the acts go to the worker as device events; nothing of the
page is read for an agent, and no frame is kept anywhere. A one-time
code works, since the owner types it. Tapping Done closes the session and
seals its cookies as the grant's sign-in, which every later session starts
with; the grant's record says a sign-in was made at the site, never what
was typed. A view not finished in ten minutes lapses and keeps nothing;
the worker ends its browser after two idle minutes. A sign-in that goes
through another host (a separate identity provider) needs that host on
the grant's sites too, or the browser's egress refuses it: Coach lists
it with the site when it asks for the grant. While a view is open,
the taf's calls on the grant are refused with `signing_in`. Asking for a
view again from the phone ends one still open (closed in the app before
Done, say) and keeps nothing from it. Disconnecting
ends an open view with the grant.

The page is served with `Cache-Control: no-store`, a content security
policy that allows nothing but its own frames and calls, and no referrer.
The token in the URL is the only credential, as a sign-in link's is: it
is random, hashed at rest, bound to one grant's one open view, and spent
by Done or the lapse.

### The worker

The worker (`cmd/browserworker`, `internal/browserworker`) is a headless
Chromium driven over its remote-debugging pipe with the standard library
alone: one browser process and one egress proxy per session, both gone when
the session closes, so nothing of a site persists but the cookies planque
seals. planque reaches it over its API with a shared key
(`PLANQUE_BROWSER_URL`, `PLANQUE_BROWSER_KEY`; the worker holds only the
key's SHA-256, `PLANQUE_WORKER_KEY_SHA256`), on a private address. An
installation with neither lists no browser and takes no browser grant.

```
POST   /sessions                 {grant, state, sites} → {session}
POST   /sessions/{id}/open       {url}                 → {page}
POST   /sessions/{id}/step       {step}                → {page}
POST   /sessions/{id}/screenshot                       → {png}
DELETE /sessions/{id}                                  → {state}
```

The browser goes out only through the session's proxy, which resolves each
name itself, connects only to public addresses (`egress.Public`, every
answer of the name, by address so a rebind reaches nothing), only to the
grant's sites when it has a list, over https on 443, and never to a name
written as an address or a local name. WebRTC, whose UDP would go around
an HTTP proxy, is taken away from every document the session loads before
the page's own script runs: the page's target, and every child target the
browser attaches with its debugger paused (a cross-site frame is its own
process under site isolation, which the worker forces on; frames in
frames; workers), each armed and only then released; a frame that can't
be armed is closed and never runs. The browser's policy
flag confining WebRTC to the proxy is set too, but measured alone it still
let STUN checks out. So a page, in any frame, opens no peer connection and
probes no address. Downloads are denied. A session unused for two minutes is ended
by the worker. Deployed, the worker runs in a network with no route to
planque's or anyone's private addresses: a second line behind the proxy.

The worker reports facts, in no language: each field's kind from its type
and autocomplete (`password`; `card` for a number or a security code;
`code` for a one-time code), its values masked at the source and again by
planque; each control, whether it submits a form, and `places_order` when
the page's semantics say so (a form with a card field, a form or button
posting to an order, checkout or payment path, schema.org buy actions, a
payment provider's button); links; a payment provider's frame; and a
checkout's total, read as one amount next to a "total" word in one
currency, in the languages above. Two different candidates, or none, is no
total. The lines are the rows (table rows, list items) with an amount in
the total's currency that aren't a total, a subtotal, shipping, tax or a
discount, in those languages: the item's words, the quantity from "× 2",
"x 2" or "qty 2" (else 1), and the price. A page that shows them
otherwise gives a total with no lines, which the owner's approval shows as
such.

The limit of reading from outside: a one-click button wired only through
script, with no words, no form and no telling attribute, on a page that
isn't a checkout, reads as an ordinary control, and `browse_act` clicks
it. An owner limits a grant to sites they know, and lists a site that can
take money with one click knowingly; the Card's limit on `browse_act`
calls bounds what such a click can do.

The tests run the real browser against the shop in HTML (`internal/fakes`,
served on loopback) when a Chromium is found (`PLANQUE_CHROME`, or the
usual names), and skip otherwise; the scripted worker in `internal/fakes`
fakes the same shop for planque's own tests. planque's end to end runs the
whole path through its guard and the real browser: the shop served over
TLS and reached as `https://shop.example` through the worker's proxy,
which the test maps there; look up an item, put it in the cart, read the
checkout, pay once with KIFF's permit, a replay that sends nothing, a
changed total that pays nothing, a card number that never reaches the
agent. The image is `deploy/browserworker/Dockerfile`: Chromium and the
binary, nothing else.

## Registered OAuth clients

The OAuth transport supports a configured public client, `client_secret_post`,
or `client_secret_basic`. A registered credential is bound to its client id
and exact token or revocation endpoint. Its secret is excluded from JSON
serialization and is never part of a grant or approval export. Token and
revocation requests do not follow redirects.

References: [OAuth client authentication](https://www.rfc-editor.org/rfc/rfc6749#section-2.3),
[token revocation](https://www.rfc-editor.org/rfc/rfc7009).

### Configuring an application

`PLANQUE_APPS_FILE` names a JSON file mapping stable application ids to
credential-free definitions. Each definition includes `name`, `kind`
(`mcp_oauth_app` or `rest_oauth_app`), `server`, a documentation `source`,
`oauth` metadata, `token_endpoint_auth_method` and `permission_words`.
REST definitions also include an `operations` snapshot of reviewed
operations, using the same request builder as an OpenAPI snapshot.

MCP applications require PKCE and exact resource binding. REST definitions
may set `no_resource_parameter` for providers that use scopes without
resource indicators. Only a confidential client may set `no_pkce` when
its provider does not support PKCE. Authorization parameters are restricted
to `access_type`, `prompt`, `include_granted_scopes` and `audience`; they
cannot override the client, redirect or one-use state.

`PLANQUE_APP_SECRETS` is a JSON object mapping configured application ids
to Secrets Manager references. The process's role needs `GetSecretValue`
only for those references. Each secret contains `client_id` and
`client_secret`; reads use `AWSCURRENT`. No reference or secret value is
accepted in a grant request. Without this setting, the operator may provision
credentials through the Go `appcredentials.Sealed` interface, backed by the
configured KV and sealer. Credential provisioning has no HTTP endpoint.

An agent requests a configured application using its id:

```json
{"application":"support-desk","agent":"support","tools":[]}
```

The owner first approves the connection on their device, signs in at the
vendor, and confirms that sign-in on the approving device. Only then may
the grant expose tools. A REST application uses its reviewed operations;
the caller cannot replace its source or base URL. MCP applications list
the server's tools through the existing MCP transport.

The signed request pins a copy of the public definition, the client id,
and the REST operation digest. A later configuration change does not
reinterpret an existing grant. Rotating a secret for the same client id is
supported; a different client id requires a new connection. Exchanges,
refreshes and revocations all use the approved endpoints. Vendor calls
retain exact-call permit verification and the existing outcome handling.

## Connection profiles and device approval

A grant request, or the request that chooses its tools, may include a
`profile` following KIFF's `kiff.profile.v1` contract. planque returns it as
`request.profile`; the existing grant digest covers it. The phone reads
KIFF's `kiff.words.v1` words using its own session, checks the profile and
words digests, and displays the ordered lines verbatim.

The signed `/approve` body includes both digests:

```json
{
  "digest": "<grant digest>",
  "kiff_profile": {
    "profile_sha256": "<canonical profile digest>",
    "words_sha256": "<canonical words digest>"
  }
}
```

The signed `/key` body carries the same object alongside `key_sha256`.
The key itself remains in its existing header. A missing or mismatched
profile digest is refused before sign-in or key storage. KIFF checks the
words digest independently when applying the profile.

planque derives read-only, amount, unit, scale and resource settings from
the profile when registering the tools. A decimal amount uses
`amount_format: "decimal_string"`; its argument or nested path must have
a string schema. It is forwarded as a string, without conversion to a
floating-point amount. KIFF measures it in minor units.

After registering the tools, planque sends the profile and the exact signed
request evidence to KIFF's gateway. The relay refuses calls until KIFF
confirms the matching profile and words digests. A failed application rolls
back the registered tools. A tool addition cannot remove the existing
profile or change its operations; the existing tools retain their profile
while the owner decides. Each operation record keeps its profile and
revision so a replay cannot adopt a later profile or send the write again.

Profiles can require arguments before any vendor call. A marker must name a
string field in the vendor schema; planque fills that field from the permit's
operation id and refuses a caller-supplied value there. Agent proposals
cannot claim reviewed catalog status or successful test calls. Profiles needing a nominal-currency lookup definition are refused until
that definition is supported. Version-header pins use the reviewed
configuration described below.

The synthetic profile/words vectors exercise both canonical digests. These
tests do not claim a live vendor check or a successful deployment.

## When an approved tool changes

The owner approves a tool as the vendor defines it then. If the vendor
changes it later, planque pauses the grant rather than send to a tool the
owner didn't approve:

- **Direct MCP connections:** before a call, planque reads the server's
  `tools/list` again, at most every 10 minutes per grant, and compares each
  granted tool's input schema with the one approved (its JCS SHA-256,
  `vendor_schema_sha256`). A tool gone, or with another schema, pauses the
  grant. If the list can't be read, nothing is paused and the call goes on
  to its usual handling (a sign-out, an unreachable server).
- **Managed connectors:** before each call, planque reads the provider's
  definition of the approved tool at its pinned version. One that reads
  and differs from the approved one pauses the grant. One that can't be
  read or checked (an unparseable answer, another tool or version) pauses
  nothing: the call isn't sent and says so, and the next call reads again.
- **REST grants** send only from the approved operation snapshot, so a
  vendor change can't alter what is sent.

A paused grant sends nothing. Every call is refused before it reaches the
vendor, with the code `operation_changed`; the grant's `status_words` say
"Paused · <what changed> · connect again to approve it", and its events
record the pause. A pause is never lifted on the same grant, not by a new
sign-in nor by adding tools: the owner connects again, approving the
tools' current definitions on a new grant.

## Reviewed vendor catalog

`PLANQUE_CATALOG_FILE` names a credential-free JSON object of reviewed
entries, keyed by a versioned id such as `support-desk@1`. Entries contain
`vendor`, `title`, optional HTTPS `icon`, `tools`, and a `profile`. Its
`source` must be `catalog`, `catalog_entry` must match the entry id, and the
profile must name exactly those tools. Test-proven status is not accepted
from catalog configuration. A changed definition receives a new revision.

For a registered OAuth connection, set `application` to an id in
`PLANQUE_APPS_FILE`. Its endpoints, scopes and REST operations come from
that application's definition. An MCP entry without a registered application
sets `server`. An API-key entry sets `kind: "api_key"`, `api` (the same shape
as a grant request), and `description_sha256`: the SHA-256 of the reviewed
OpenAPI source bytes. A changed source is refused before a grant is created.
The full snapshot is still frozen and signed with the grant.

Both catalog endpoints require a bearer KIFF key for the account:

- `GET /v1/accounts/{tenant}/catalog?q=...` searches vendor names and titles,
  without calling a vendor business action. Optional managed discovery contacts
  the configured connection provider. It returns `{entries, by_key: true}`
  and provider pagination when configured. Each
  entry has `entry`, `vendor`, `title`, optional `icon`, `how` (`oauth` or
  `key`), `state` (`available`, `connected` or `needs_reconnect`) and, for an
  active connection, `grant`. State comes only from that account's grants.
- `POST /v1/accounts/{tenant}/catalog/{entry}/grants` accepts only `{agent}`,
  identifying the taf's KIFF agent. It returns `{grant, needs_key}`, with the
  full grant view. The server supplies the reviewed tools and profile and
  allocates their final names before computing the approval digest. This
  request creates no sign-in link and grants no tools. The existing signed
  approval or key flow, including KIFF's profile and words digests, follows.

The catalog starts empty when neither a file nor a managed provider is configured.
Optional [managed connections](managed-connectors.md) add paginated upstream app
and tool discovery with their own credential custody and device-approved grants. Installing definitions
alone does not establish live vendor evidence; each supported entry needs
reviewed vendor documentation and a measured connection with test accounts.


### Version headers

A REST application's `version_header`, or an API-key catalog entry's
`version_header`, names the vendor's reviewed version header (for example
`Pay-Version`). A profile with `pin: {"kind":"header","value":"2030-01-01"}`
then binds that value to that header. A caller cannot choose the header name
in a grant request. The exact header and value are in
`request.api.version_header`, covered by the grant digest, and sent on each
REST operation alongside the stored credential. Later configuration changes
cannot change an existing grant's pin. Header names must identify a version;
control characters and collisions with the credential header are refused.
