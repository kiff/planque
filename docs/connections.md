# Connecting vendor operations

planque holds the vendor credential and sends each operation only with the
permit for that exact call. The owner approves the requested connection on
their device before its tools become available.

It follows the sign-in standard vendors' MCP servers speak: OAuth 2.1 with
discovery (RFC 9728, RFC 8414), dynamic client registration (RFC 7591) and
PKCE.

## MCP client identity

When the authorization server advertises
`client_id_metadata_document_supported: true`, planque uses its configured
HTTPS `/oauth/client-metadata.json` URL as its client id. The document fixes
the client name, redirect URI, response types and grant types. Its contents
never depend on the request's Host header, account or grant.

Servers offering dynamic registration continue to use that flow when they
do not advertise client metadata support. PKCE, resource binding, one-use
state and confirmation on the approving device apply to both flows.

A sign-in link that came back once stays one-use: landed on again (a
vendor's "click here" after its redirect, a reopened tab), it runs nothing
and moves no grant. Until the link's own deadline planque keeps a mark of
it, by the state's hash, never the state, so the page can say what
happened: the sign-in already went through (its grant is signed in or
connected), so the owner goes back to tafs, where it finishes; or the link
was used and the sign-in starts again from tafs. A link planque never gave
is unknown, as before. The page names the app, never a grant, an account or
a code.

**The resource.** A server's metadata must name the server planque calls:
the same URL, or the same scheme and host with a path made of the server's
leading segments, as Qonto names `https://mcp.qonto.com` for its server at
`/mcp`. The sign-in is bound to the resource the metadata names. Another
host, a port, a query or a sibling path (`/mcpx` for `/mcp`) is refused.

**A vendor that refuses planque's client.** Some authorization servers let
only the clients they allow register (Qonto allow-lists the redirect URIs of
MCP clients). When one refuses planque's registration with a 4xx, approving
the grant answers 502 with `code: "vendor_refused_client"` and words the apps
can show: the vendor doesn't allow tafs to sign in yet, and nothing was
connected. One log line says it too, `client registration refused`, with the
grant, the vendor's host, its status and its OAuth error code, never its
description. The grant stays pending.

**The scopes.** planque asks for the scopes a server's metadata lists
(`scopes_supported`), except at a server it reviewed. Qonto's sign-in asks
only for reading the business, its clients, its invoices and their
e-invoicing state, and writing clients and client invoices. It never asks
for transfers, cards, bank accounts, direct debits, team or members.

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

## Qonto by key

Qonto's MCP server can't take planque's callback yet, so Qonto also connects
by key through the reviewed catalog entry `qonto-key@1` (`deploy/catalog.json`):

```json
{"agent": "invoicing"}
```

posted to `POST /v1/accounts/{tenant}/catalog/qonto-key@1/grants`.

- **The key** is the owner's own, from Qonto's settings: `<organisation login>:<secret key>`,
  sent as is in `Authorization`, with no prefix. The owner gives it on the phone
  (the signed `/key`); no agent sees it.
- **The description** is planque's reviewed copy of Qonto's Business API for
  these operations (`deploy/specs/qonto-invoicing-openapi.json`), taken from the
  OpenAPI on Qonto's reference pages, with the prose left out and `/v2` moved
  into the server. Qonto publishes no single file to fetch. The entry pins its
  SHA-256, so a changed copy is refused, never used.
- **Reads only:** the business and its bank accounts, clients, client invoices
  and one client invoice. Creating an invoice isn't offered: Qonto's create makes
  a real (`unpaid`) invoice unless `status: draft` is sent, and a draft-only tool
  needs `status` fixed on the grant, with the invoice's lines still the agent's
  to give. A later entry (`qonto-key@2`, below) adds it, approved anew.

### `qonto-key@2`: reads, and invoice drafts from quotes

`qonto-key@2` is `@1` plus the quotes and supplier invoices, as reads, and
one write: an invoice drafted from an accepted quote. It's posted the same way,
to `POST /v1/accounts/{tenant}/catalog/qonto-key@2/grants`, and the owner
approves it anew: a grant from `@1` stays as it was.

- **Its description** is `deploy/specs/qonto-invoicing-openapi-2.json`: `@1`'s
  reviewed copy, with `GET /quotes`, `GET /quotes/{id}` and
  `GET /supplier_invoices` added from Qonto's reference pages, their prose left
  out. The entry pins its SHA-256. `@1`'s file is unchanged, since `@1`'s grants
  pin it.
- **The reads:** `qonto_list_quotes` (filter `approved`), `qonto_get_quote`
  (its client, lines, prices, and `invoice_ids`, empty while no invoice comes
  from it) and `qonto_list_supplier_invoices`, beside `@1`'s four.
- **The one write, `qonto_draft_client_invoice`:** Qonto's
  `POST /client_invoices` with `status: "draft"` fixed on the grant. A draft
  goes to no one, and it's never finalised or sent: no operation that does
  either is offered, nor any that pays, cancels or deletes.
  - **Its arguments:** `quote_id` (required, the call's resource, so the Card
    and the record name the quote), the client, the dates, the currency, the
    IBAN the customer pays to, and `items`.
  - **The IBAN** must be one of the business's own accounts: before each
    call, planque reads `GET /organization` with the same key and refuses an
    IBAN it doesn't list (spaces and case aside), so one planted in a quote or
    a client's note never becomes the payee. A read that fails refuses the
    call.
  - **`items`** is an array argument with a reviewed schema: Qonto's own line
    fields, without `discount`. A line with a field the schema doesn't name is
    refused before anything reaches Qonto, and so is the invoice's own
    `discount`, which isn't one of the arguments.
  - **The write's effect** is `none_until_approved`, and `claim: invoice_draft`
    means one taf at a time drafts from a given quote.

**Array arguments** (`narrow.go`): a narrowed tool's argument may be
`"type": "array"` with a `schema`, a JSON Schema of `type`, `properties`,
`required`, `items`, `enum`, `maxLength`, `maxItems` and `description` only.
Its items are objects, each naming its fields, plain values (strings,
integers, numbers or booleans: a draft's cc addresses), or rows of plain
values, one level deep (a sheet's `values`: `[["Ana", 12]]`), each cell
checked as a plain value; a row of objects or of rows is refused. Annotations a
provider's generated schema carries (`title`, `default`, `examples`,
`example`, `$comment`, `deprecated`, `readOnly`, `writeOnly`, `nullable`,
`format`, which planque never checked) are dropped at every level, and an
optional value written as `anyOf` (or `oneOf`) of it and `{"type": "null"}`
is read as the value alone, before the schema is checked: the request is
stored, declared and signed without them. Any other keyword (`pattern`,
`minimum`, `minItems`, `const`, `$ref`, another `anyOf`) restricts or can't
be checked, so it is refused (`argument_schema`), never dropped. The agent
sees that schema. A call whose value has a field the schema doesn't name, the wrong kind
of value, a value outside its list, or a missing required field is refused.

**Object arguments** (`narrow.go`, `checkObject`): an argument may also be
`"type": "object"` with a `schema` naming each field the taf may set (a sheet
row's update, whose `values` is a text per column):
`{"path": "values", "type": "object", "schema": {"type": "object", "properties": {"Introduction": {"type": "string"}}}}`.
Each field is a plain value or an array, as an array argument's schema is.
It is never another object (name a nested object's fields as arguments of
their own, at their paths), never `null`, a type list or alternatives. A
field's name is 1 to 64 characters with no dot (spaces are fine: "First
name"), since a field's path is the argument's and its name split on dots,
which is how the disclosure line finds the text. An object names at most 50
fields, and takes `type`, `properties`, `required` and `description` only
(no `enum`, `items` or lengths on the object itself). The
same keywords are taken, and the same annotations dropped. A field the
schema doesn't name is refused, whatever `additionalProperties` says. An
object that names its fields is stored closed (`additionalProperties: false`),
even when it came open (`true`, or a schema such as the vendor's "a text per
column"): closing only narrows it, and what is stored and signed says so. An
open object that names no field is refused, and so is an object with no schema
or no `properties`: each is told to name the fields the taf writes, by the names
the job or the owner gives (the column it marks), with an example. The approval names the fields
(`values (Introduction)`). An object argument is never the amount or the
resource, and it can't share its path with a pin or another argument inside
it. Lifting each field as its own argument (`values.Introduction`, a string)
still works. An `argument_type` refusal shows both forms, so a retry works
whichever the taf picks. KIFF never gets an argument's schema: planque checks
each call against it. A refusal of an argument's type or schema carries its
`shape`, the declaration by keys and types alone (field names, an enum's
count, never a value or a description), and planque logs it.

**The vendor's allowed values** (`narrow.go`, `schemaOver`): a narrowed
string, number or boolean argument keeps the `enum` the vendor's description
gives at its path (at most 50 values, each of the argument's type), so the
agent is shown the values the vendor takes instead of a bare string. A call
with another value is refused before anything is sent, naming the values it
takes: a desk asked for `status=new` when its statuses are open, pending and
solved answers an empty list, not an error, and the taf would read that as
"no tickets". A grant connected before keeps its schema until its tools are
connected again.

**Arguments checked against the account** (`one_of`): a string argument may
name a read-only tool of the same grant and the dotted path of the values it
allows (`{"tool": "get_organization", "path": "organization.bank_accounts.iban"}`).
On an API-key grant, planque makes that read with the grant's key before each
call and refuses a value it doesn't list.

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
completes one. `sites` are the only hosts (and subdomains) the browser
opens. They're required, at least one, and none may be a shared domain
where anyone has a site (`github.io`, `vercel.app`, `myshopify.com`) or a
registry's suffix (`co.uk`): the Public Suffix List
(`internal/browser/public_suffix_list.dat`, refreshed with
`tools/publicsuffix.sh`) decides. `currency` is the only one `browse_pay` pays in, and is
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
(`ref_unknown`). Each step is exactly one of its kinds in the tool's schema (fill, select,
click, check, press, scroll, back), with only its own keys, so the taf is told
before the call is held; `option` goes only with select. Each step names its ref by the page's own words: `label`,
the words the read gave it, with `submits: true` for a control that submits
a form and `to`, the host a link goes to. A held call shows them to the
owner, so planque refuses a step whose words aren't the page's
(`label_mismatch`): a missing or different label, a submit not said to
submit (a click or Enter on a submitting control, Enter in a field, or a
field typed key by key with a line break), a link to another host, or a
control or link with no words, which the owner couldn't read. `submits` is
the read's fact about the page: a button a script makes submit reads as not
submitting, so `payment_step` stays the net for paying. Links are followed. The steps before it ran, and the
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
no total was read in the last 15 minutes (`total_unread`), when the lines,
total or currency changed (`checkout_changed`), when the page is in
another currency (`currency_other`), or when `ref` isn't the control whose
words place the order (`ref_unknown`). Then it marks the checkout and
clicks once; a replay answers as the first call did. The mark is written
before the click, as a conditional write, and a success removes it: a click
whose outcome is not known (it errors, or never returns) leaves it for 24
hours, so a pay on the same lines, total and currency, from any operation or
at the same moment, is refused (`outcome_unknown`) and a new read of it gives
no total; a different cart is a new checkout. Each write keeps the page before and after
(its text) with the operation's record, and a screenshot after in the
evidence store when one is configured; both go with the grant.

A site's sign-in, when one is needed, is the owner's own, in a live view of
the grant's browser that planque serves (a later step); the storage state
after it is sealed like a token, opened only for the grant's own session,
and deleted on disconnect. planque never stores or logs a password: in the
live view, what the owner types goes to the browser through planque's
memory as an act, and is kept nowhere.

### The live sign-in view

A site's sign-in is the owner's own. From the connection in the app, the
phone asks, signed, `POST /v1/accounts/{tenant}/grants/{id}/sign-in` with
an optional start `url` (on the grant's sites; the first site's root by
default). planque starts a session on the worker with the grant's kept
sign-ins and answers `{view, expires, words}`: a one-time URL under
`/browser/view/{token}`, good for ten minutes, which the app opens in a
web view and nothing more; the page's words are planque's, in the language
the app asked in when it opened the view (see "The pages' language" in the
README). The page shows
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

A configured worker is offered to no account by itself. `PLANQUE_BROWSER_OFFER`
names which: `on` (every account), `test` (only those in
`PLANQUE_BROWSER_TENANTS`, `tenants:<id>[,<id>]`) or `off`; unset, none. An
account off the list sees no "Use websites" in the catalog, its ask for a
browser grant is refused (403 `browser_off`), and a browser grant it already
holds opens nothing: each call is refused before a browser starts
(`browser_off`), and so is the sign-in view. A malformed value stops planque
starting.

Paying is a setting of its own: the hosted browser pays (`browse_pay`) only
with `PLANQUE_BROWSER_PAY=on`. Unset, it opens, reads, fills and reads a
checkout, but `browse_pay` isn't offered, a grant that asks for it is
refused, and a call to it is refused before any browser starts (`pay_off`),
whatever the owner approved. A click that would pay is refused by
`browse_act` either way (`payment_step`), and downloads are denied.

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
written as an address or a local name. A redirect is checked hop by hop, as
the browser follows it, so a listed site that redirects off its sites
reaches nothing; a page left at a plain-http address is the proxy's
refusal, and reads as refused, never as a page. WebRTC, whose UDP would go around
an HTTP proxy, is taken away from every document the session loads before
the page's own script runs: the page's target, and every child target the
browser attaches with its debugger paused (a cross-site frame is its own
process under site isolation, which the worker forces on; frames in
frames; workers), each armed and only then released; a frame that can't
be armed is closed and never runs. The browser's policy
flag confining WebRTC to the proxy is set too, but measured alone it still
let STUN checks out. So a page, in any frame, opens no peer connection and
probes no address. Downloads are denied. A session unused for two minutes is ended
by the worker. Deployed (`deploy/browserworker/aws.sh`), a network ACL on
the worker's subnets is the second line behind the proxy: inside the VPC it
lets through only planque's calls to the worker's port and their replies,
and out of it only TCP 443 to public addresses; private, link-local and
shared ranges are denied. Each task runs one session, so a page that broke
out of Chromium would hold no other account's session, and the image is
rebuilt each week for the browser's security updates.

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

## A call whose outcome comes later

A write can answer later, as one of the app's events: a call placed ends
minutes later, a transfer settles, a document renders. Its profile
operation says so:

```json
{"tool": "place_call", "kind": "write", "effect": "sends",
 "completes": {"event": "CALL_ENDED", "key": "id", "event_key": "call.id", "timeout_s": 900,
               "words": "The call runs on the vendor's side; its outcome comes back when it ends."}}
```

`event` is the app's event slug; `key` is a dotted path to the handle in the
call's answer (its structured content, else its text read as JSON);
`event_key` the path to the same handle in the event's payload; `timeout_s`
60 seconds to 3 days; `words` one line the owner reads for the wait. A read
never completes later.

Such a call is refused before sending while the taf has no active
subscription to that event on the connection (`not_listening`): nothing
would ever resume the run. Sent and answered, the relay's result is the
vendor's answer with a line added, and in its structured content
`{"pending": true, "sent": true, "wait_for": {"grant", "event", "key", "words"}, "timeout_s"}`;
an answer with no handle at its path is left as it was, saying so. A
replay answers the same. When the event comes, its delivery
(`planque.event.v1`) carries `key`: the value at `event_key` in the payload,
for an event one of the connection's operations completes on; absent
otherwise. The run that made the call is matched to it by `(tenant, grant,
event, key)` on tafs' side, without reading the vendor's payload.

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
to Secrets Manager references. With `PLANQUE_APP_SECRET_PREFIX`
(`planque/apps/` as deployed), an application with no reference of its own
uses the secret named that prefix plus its id (`planque/apps/googlesheets`),
so a new application needs only its definition and its secret. The id must
still be an application id by its pattern, so no grant can name another
secret. The process's role needs `GetSecretValue` only for those references,
or under that prefix. As deployed, the definitions are `deploy/apps.json`,
reviewed like the catalog, with a security read of every scope, since a
scope's breadth is the definition's to keep narrow: no code can check it. Each secret contains `client_id` and
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
operation id and refuses a caller-supplied value there. An operation may say
where the owner looks at the vendor's own record of its resource
(`record_url`, an https address with the resource's id path-escaped into
`{resource}`; reviewed catalog only, as `no_effect` is) and which read of the
profile gives the resource's own name (`label_read`, for an operation with a
`resource`): KIFF shows both when it asks the owner "Did this happen?". A write
with a `resource` may name the kind of action it is (`claim`, such as
`dispute_answer` or `reminder`): KIFF then lets one agent at a time act on one
object for it, and refuses another agent's call of that kind on that object
while the first agent's is being decided, held or sent. A write may also name
where its amount lookup's read gives a related charge (`claim_related`,
`{kind: "refund", path}`): KIFF shows the agent the refunds already made on it,
and blocks nothing. Agent proposals
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

**When planque changes one of its own tools** (the browser's `browse_act`, say), nothing pauses, and an
existing grant keeps the schema it was approved with. KIFF's gateway registers each of a connection's
tools with the input schema planque listed when the connection was made, keeps it, and shows the taf that
copy; it doesn't read planque's list again. planque checks every call against its current rules. So a
change that adds a required argument reaches only grants connected after it. On an older grant the taf
never sees the argument, and its calls are refused at run time: when #476 made a step carry `label`, a
grant connected before it could only send steps that are refused as `label_mismatch`. That fails closed,
but the owner can be asked to approve a held call that could never run. Today the remedy is to connect
again, which registers the current schema. A refresh would take three things:
1. KIFF re-reading planque's `tools/list` for each connection after a planque deploy.
2. Comparing each tool's schema with the stored one.
3. Applying the change by what it does: one that only narrows a tool (a new required argument, a removed
   option) without asking, and one that widens what the tool can do only on the owner's approval, as a new
   grant would. KIFF already replaces a schema this way for its own control server's tools, when a batch
   asks it to; a grant's tools have no such path.

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

An entry's `titles` holds its `title` in other languages, by language tag
(`"de"`, `"pt-br"`), each starting with the vendor as the English one does.
`GET /v1/accounts/{tenant}/catalog` shows each title in the asker's language:
its `lang` parameter, else the first language its `Accept-Language` names; a
tag's own title, then its language's, else the English. planque's own
browser entry, "Use websites", is titled in the same languages. The titles
are words only: they never change what an entry offers or what a grant signs.

A key entry's address is one of three things:

- **Fixed:** `"base": "https://api.example.com/v1"`.
- **An account's address:** `"base": "https://{account}.example.com/api/v2"`,
  for a vendor that gives each customer a subdomain. The template is `https://`,
  then `{account}` once and alone as the host's first label, then the vendor's
  own domain written out (two labels or more), then an optional path. No
  other variable, no port, no user info; `{domain}.{tld}` and a variable after
  the domain are refused when the entry is loaded.
- **A closed list of regions:** `"bases": ["https://api.example.com",
  "https://api.eu.example.com"]` with no `base`, two to eight fixed addresses,
  each on the pinned path the profile reviewed.

A grant request for such an entry names the address it wants: `account` (the
owner's own name at the vendor, `"acme"` for `acme.example.com`) or `base` (one
of the regions, exactly as listed). Taken exactly as given and lowercased, an
account is one label of letters, digits and dashes, 63 characters at most, the
rule planque applies to an account's address on any key grant: no dot, `@`,
colon, slash, space or non-ASCII character. `acme.evil.com`,
`acme.zendesk.com.evil.com`, `acme@evil`, `acme:444`, an empty or 64-character
label and an IDN look-alike are refused with a 400. The catalog read lists the
template as `account` and the list as `regions`, so the asking agent knows
which to give.

**A login in the Basic username.** Some vendors take the customer's own login
beside the key: an API token is the password and the username is an email
address (Zendesk's `login/token`). An entry says so with `"auth": {"in":
"basic", "user": "{login}/token"}`, or `"user": "{login}"`: those two forms
only, and any other brace in a username is refused when the catalog loads. The
grant request then names `login`: one email address of letters, digits and
`. _ % + -` before the `@`, a domain of at least two labels, at most 254
characters. A space, colon, slash, comma, quote, bracket or non-ASCII character,
a second address, or a login on an entry that has no `{login}` is a 400. The
grant keeps the one username it was approved with, the owner reads it in the
approval as for any Basic username, the key stays the password and is never
shown, and the entry itself stays a template. The catalog read lists such an
entry with `"login": true`.

The address is fixed in the grant: the signed request holds the one address it
was approved with (`api.base`), the approval shows its host ("acme.example.com"),
and no tool has an argument that can change it. Changing it is a new
connection. An ad-hoc grant's `api.base` is still one address; a template or a
list there is refused.

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

### Managed marks

Whether a vendor's apps may be connected through the connection platform is a
**managed mark**, kept in planque's table beside the catalog, never inside an
entry (a published entry can't change; a mark must be able to). A mark names
the vendor, the exact hosts its apps may use (a host, or `*.domain` for a
vendor that serves one per customer, as `vendorHosts` in the review checks
writes them) and the date of the published subprocessors-list update it
follows.

- **Only an enrolled operator's device sets one.** `POST /v1/marks/managed
  {name, hosts, follows, operator, nonce, signed_at, signature}`: the device
  signs the SHA-256 of `{name, hosts (sorted), follows}` with the catalog's
  signing message (`mark-managed`), planque computes it from what was sent, the
  signature holds for five minutes and counts once. No credential, session or
  generator can set a mark, and a row written into the table any other way, or
  changed after it was signed, isn't read as one.
- **No wildcard on a domain anyone can have a name under.** A mark may take
  `*.domain` only for a vendor's own per-customer domain (`*.zendesk.com`, or
  `*.acme.co.uk`, a vendor's own domain under a registry suffix), never for a
  registry suffix of two labels (`*.co.uk`) or for shared hosting
  (`*.github.io`, `*.herokuapp.com`) **or anything beneath it**
  (`*.s3.amazonaws.com`, `*.acme.github.io`): that would let any user of the
  platform answer for the vendor. `public_suffixes.json` lists them in two
  parts, `registry` and `hosting`, kept short on purpose (the owner's device
  shows the hosts it signs, too); a suffix is added there in a PR that names it.
- **The gate.** A managed connection needs the vendor's mark, or a test tenant,
  or, for a vendor with no mark at all, what `managed_apps.json` (the seed)
  allows: a place on it in its list mode, or not being one of its exceptions
  in its all mode (docs/managed-connectors.md). A **withdrawn** mark wins over
  the seed.
- **The catalog read** says `managed: true` for an entry only when its vendor
  has an active mark and **every endpoint the entry has** (its API addresses,
  each region, an account template read on the vendor's own domain, or its app's
  endpoints) is on the mark's hosts. An entry that calls itself by a managed
  vendor's name but points elsewhere, or that adds a host the mark doesn't list,
  is not managed until the owner signs the mark again with that host.
- **Rotating an operator.** A mark counts only while the operator whose device
  signed it is still in `PLANQUE_CATALOG_OPERATORS`. Removing an operator, or
  replacing their device, stops every mark they signed from counting (and
  every managed connection it allowed from starting) until an enrolled
  operator signs it again. Plan a rotation as: enrol the new device, have it
  sign each mark again, then remove the old one. Connections already made stay.
  **A row that exists but no longer counts reads off, never the seed:** an
  operator rotated out also drops the withdrawals they signed, and a seeded
  vendor must not become managed again for it. The same holds for a row altered
  after signing, and while the table can't be read and nothing was read before:
  the gate refuses managed connections (a test tenant is still let through).
- **Withdrawing** (`POST /v1/marks/managed/{vendor}/withdraw`, signed over the
  mark's digest as it is, `unmark-managed`) stops new managed connections within
  seconds (an instance keeps what it read for five). It revokes nothing: a
  connection already made stays, and what happens to it is the owner's decision.
- **History.** Every change is kept as one more line (`LOG#`), in order, with
  who signed it and when, and is only ever added to. `GET /v1/marks` lists the
  current marks, with no signature or nonce in it.
- A generator never sets a mark, and a regenerated entry never inherits one.

### The catalog kept as data

The reviewed catalog is also kept in planque's table, so an entry is added
without a deploy. These entries are the **verified** tier: a person vouched for
each one. Connecting never requires one. Each revision (`name@n`) moves
`draft` → `reviewed` → `published` → `retired`, and never changes once
published: a change is a new revision. A grant keeps the revision it was
approved with, so retiring one stops only new grants.

`PLANQUE_CATALOG_OPERATORS` lists who may work on it, as
`[{id, token_sha256, device_key}]`. Each operator has two identities that
never stand in for each other: an API credential (`Authorization: Bearer`,
kept only as its SHA-256) drafts, changes and reviews, and an enrolled device
(a P-256 key, base64 SPKI DER, catalog-only and apart from any account's
phone) publishes and retires. An agent drafting through the API has a
credential and no device, so it has no path to `published`.

- `POST /v1/catalog-entries` `{entry}` drafts a new revision, with its `id`,
  checked as a shipped entry is. A revision that exists is never replaced (409,
  naming the next one).
- `PUT /v1/catalog-entries/{id}` `{entry}` changes a draft or a reviewed
  revision, which is a draft again, changed by that operator.
- `POST /v1/catalog-entries/{id}/review` `{digest}` reviews the draft as read.
- `POST /v1/catalog-entries/{id}/publish` and `…/retire`
  `{operator, digest, nonce, signed_at, signature}`: the operator's device
  signs, over the SHA-256 of
  `planque.catalog.v1\n<planque base URL>\n<publish|retire>\n<id>\n<digest>\n<nonce>\n<signed_at>`,
  where `digest` is the JCS SHA-256 of the entry with its id. The signature
  holds for 5 minutes, and each nonce counts once. Who drafted or last changed
  a revision can't publish it. A revision is published only reviewed, and
  only when no write it holds can be made final by an argument it leaves
  unpinned: for a key entry, planque reads its pinned description, and walks
  each write's arguments, nested objects included, by their dotted paths (as
  `fixed_arguments` name them, up to four parts). An enum that holds an
  undoable value (`draft`, `pending`, …) and any other, a boolean named for a
  final step (`send`, `issue`, …, pinned `false`) or for the undoable state
  (`draft`, `is_draft`, pinned `true`) must be fixed to the undoable value, or
  the tool left out. OpenAPI 3.1's list types (`["string", "null"]`) are read.
  A write whose arguments can't be read isn't published: the check fails
  closed.
- `GET /v1/catalog-entries` lists every revision for an operator.
  `GET /v1/catalog-entries/{id}` reads one, with the digest a device signs,
  and needs no credential: an entry holds no secret. It also carries
  `described`, outside the digest: for a key entry, each tool as planque reads
  the pinned description (`method`, `path`, `read_only`, and the
  `final_arguments` with the value each would be pinned to), or `why` it
  can't be read. A device draws read or write from this, never from the
  entry's own labels. An MCP entry's tools are listed by its server once
  connected, so it says so.

planque serves a published revision only while its stored digest is the one
its publish signature covers, and that signature verifies, with the key of an
operator still enrolled, over that revision and digest. A row changed after
that, or marked published by anything but an enrolled device, isn't served,
and removing an operator stops serving what they published. A write's read or
write status at publish is what the pinned description says, never the
entry's own label, and a tool the description doesn't have is refused. The
row's signed `published` and `retired` blocks are the record: each is set
once and never edited. A log of them, by time, is kept beside it. A search
or a list offers only the newest revision of each name. The table is read at most once a
minute, alongside `PLANQUE_CATALOG_FILE`'s entries, which are served as
before until they're published here. A revision retired here isn't served,
whether or not it was shipped.

### Which app a grant is for

Every grant's view names its app's catalog `entry` when it's certain, so an
app lists each connection under its app by `entry` alone:

- **A grant asked from the catalog** names its entry in `request.catalog_entry`,
  as before. A managed app's is `connector:{provider}:{toolkit}`, the browser's
  `browser`, and a Telegram bot's `api:telegram`.
- **A vendor's own server on the reviewed server list** gets the app the list
  names: the managed app the owner already uses for that vendor, such as
  `https://mcp.stripe.com` → `connector:composio:stripe`. The list comes first,
  whatever the catalog holds for that server.
- **A grant asked by server, key or registered application** otherwise names the
  one reviewed entry for the same server reached the same way: an MCP sign-in, or
  a REST API by key or by application.
- Either way, a new grant has it in `request.catalog_entry`, so the owner signs it
  with the rest.
- **A grant made before requests named their app** gets its entry recorded
  beside the request, never in it: the request is what the owner signed, and
  stays byte for byte as signed, so its digest still matches the approval. The
  owner's grants list records it the first time it lists that grant, once.
- **When no entry matches, or more than one does,** there's no `entry`. It's
  never guessed. Two rows of the server list that match one server with
  different apps leave it without one too.

**The server list.** `PLANQUE_SERVER_APPS_FILE` names a JSON object from a
vendor's own HTTPS server to a catalog entry: `connector:{provider}:{toolkit}`
or a reviewed entry's versioned id, at most 500 rows.
- A row that's an origin (`https://mcp.stripe.com`) names every path on it,
  since a vendor's MCP host is its own. A row with a path names that path alone.
- Anything else in the file refuses the whole list at start.
- It's configuration, reviewed like the catalog: `deploy/server-apps.json`
  ships with the function, and `deploy/aws.sh` points the variable at it.

The catalog's `state` counts any connection to the app, one Coach started
included.

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

## The words, in the owner's language

A grant's words (`text`) and its status line (`status_words`) are written
in each language tafs ships: en, de, fr, es, it, pt, nl and pl
(`internal/words`, held to tafs' glossary: taf and Card stay names). A view
of a grant, wherever planque shows one, is in the language its app asks in
(`Accept-Language`, by q), and says which (`lang`). A language planque
doesn't speak, or a grant made before it wrote its words in each, reads in
English (`lang: "en"`). The exact terms under Details (`terms`), the
vendor's, the app's and the taf's own names, and a managed connection made
before titles stay as signed.

What the phone signs is the grant's `digest`:

- **v1, today, with no `digest_v`:** the SHA-256 of the request's canonical
  JSON (RFC 8785). It covers the request, not the words.
- **v2, `digest_v: 2`:** the SHA-256 of the canonical JSON of
  `{"lang": lang, "request": request, "text": text}`, where `request` is the
  very object the view sent. It covers the words the phone showed, and
  their language. An app that doesn't know `digest_v`, or reads v2 without
  `lang` or `text`, signs nothing.

`/approve` takes either digest. With v2, the approval keeps the words
signed and their language (`approval.lang`, `approval.text`): the receipt of
what the owner read. A digest over words planque didn't write, or over a
grant's words relabelled with another language, approves nothing (409).
Views carry v2 only with `PLANQUE_DIGEST=v2`, once both apps check it.

The shared vector, which both apps pin: the request
`{"agent":"books","asked_by":"setup-agent","created_at":"2026-11-11T16:00:00Z","grant":"g_02","server":"https://mcp.qonto.com/mcp","tenant":"t_01","tools":[{"read_only":true,"tool":"list_invoices"}],"v":"planque.grant.v1"}`,
`lang` `fr` and `text` `books pourra lire vos factures sur mcp.qonto.com`
give `f1ca3ab371ecf97a82d134204e4931ef39ea474b7cd6a3e233b36f97499fbb32`.
