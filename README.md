# planque

**The safe place for your agents' keys.** An agent asks to act in a vendor
account: payments, a shop, email, any API. The owner approves on their own
phone, seeing exactly what will be granted. Direct connections keep vendor keys in planque, in Europe or on the owner's
own servers. Optional managed connections keep vendor credentials at the
configured connection provider. Every call is made only
with a permit for that exact call.

*Planque* is French slang for a safe hiding place.

## For whom

People who let AI agents work in their business tools and need to know:
- who approved what;
- that the keys never sit in the agent;
- that each action stays within what was approved.

## How it works

planque follows the standard that modern vendors' MCP servers speak:
OAuth 2.1 with discovery (RFC 9728, RFC 8414), dynamic client
registration (RFC 7591) and PKCE.

An MCP authorization server can also use a **Client ID Metadata Document**
instead of dynamic registration. When it advertises that support, planque
uses its configured HTTPS `/oauth/client-metadata.json` URL as the client id.
That public document names this installation and its exact callback; it
contains no account, grant or credential data.

1. **Ask.** An agent sends a grant request that says exactly what it wants:
   which vendor, what access, for which agent.
2. **Approve.** The owner approves in their own app, for example tafs, with
   Face ID. Their device signs that exact request with its key, and KIFF,
   which keeps the devices each person approves with, says whether it did.
   Nothing is granted without the signature.
3. **Connect.** planque discovers the vendor's sign-in server, registers
   itself, and sends the owner to the vendor's own sign-in page. The token
   lands in planque, encrypted with a managed key, bound to that vendor account,
   refreshed safely, and revocable.
4. **Call.** For each call, planque checks a permit for that exact call (for
   example KIFF's execution permit), then calls the vendor itself. The agent never
   receives a token.

An optional [managed connection provider](docs/managed-connectors.md) adds
paginated app discovery and hosted authentication without an app-specific
adapter for every vendor. The same device
approval, exact-call permits and once-only relay govern each selected action.
Incoming chat messages for a taf's door use separately approved subscriptions
and a signed door handoff; receiving a message grants no reply authority.

**Connection profiles.** A grant request can carry a profile: what each of the
connection's operations does, in terms KIFF's Card can count (where its amount
and resource are, and how KIFF can tell it happened), with KIFF's words for it.
The owner's phone signs the profile and the words together with the grant.
planque hands that signed evidence to KIFF, which applies the profile only if
the signature covers it; if KIFF doesn't apply it, nothing is connected. Each
operation keeps the profile revision it was decided under, so a replay never
sends it again under another.

## Two ways to run it

- **Hosted:** we run it in the EU. The owner runs nothing.
- **Self-hosted:** one container, on any cloud or server. Direct connections
  keep credentials on the owner's infrastructure; managed connections use
  the configured provider's credential custody.

It's the same code either way.

## Done when (the first end-to-end)

- From kiff/tafs' setup chat, an owner connects a payment provider through its
  MCP server:
  1. the approval on the phone;
  2. the provider's own sign-in;
  3. the token in planque;
  4. a refund inside the Card goes through;
  5. a refund outside it is never sent.
- A call without a valid permit, with a permit for another call or account,
  or with an expired permit sends nothing. A test covers each case.
- A grant request the phone didn't sign creates nothing.
- Disconnecting revokes the token at the vendor and deletes it here.
- The signed approvals can be exported as evidence of human oversight.

## The API

| Call | Who | What |
|---|---|---|
| `POST /v1/accounts/{tenant}/enrolment` | a device of the owner's that KIFF knows, signed (`d=`) | hands planque the account's KIFF **connector** key; any of the person's devices may. The enrolled key is kept while KIFF still takes it for the account (**200** `{"kept": true}`): the tools it connected are that key's in KIFF. A new key replaces it (**201**) only once KIFF no longer takes the enrolled one. A key that can do more than connect tools is refused. A signature without `d=` (a phone that enrolled its own key before KIFF kept devices) is no longer taken; a `device_key` in the body is accepted and ignored |
| `POST /v1/accounts/{tenant}/grants` | an agent, with its KIFF key | asks for a vendor's MCP server, some of its tools, for a taf; returns the plain words and the digest the owner signs; grants nothing. With `tools: []` it connects first: once signed in, planque only lists the vendor's tools. `agent_name`, optional, is the taf's name for people, on one line: the words name the taf by it, else "this taf", and its agent id is under Details (`terms`) |
| `GET /v1/accounts/{tenant}/grants/{id}` | an agent of the account | the grant's status, `status_words`: that status in the owner's words, one line every app shows as it is ("Connected", "Signed in · the tools are being chosen", "Connected · new tools waiting for your approval"…), and `vendor`: the vendor's name as its MCP server gives it (`serverInfo` title, else name), else its host. Connected tools reach KIFF with the vendor's title for each (its MCP `title`, else `annotations.title`), for people to read; names a vendor gives are one line, at most 100 characters |
| `GET /v1/accounts/{tenant}/grants/{id}/catalog` | an agent of the account | for a grant asked for with no tools, once signed in (`listed`): each tool the vendor lists, with its description and arguments, and the account the sign-in acts in |
| `POST /v1/accounts/{tenant}/grants/{id}/read` | | answers **410**. Setup's reads of a grant's `setup_reads` tools go through **KIFF's gateway** like any call (kiff/tafs#138): when the grant connects, planque asks KIFF for Setup's own read Card (`POST /v1/me/setup-reads`, with the owner's signed approval, which KIFF checks itself), at most 50 reads a day. Each read is decided on that Card, has KIFF's permit, and is kept in KIFF's record, and the relay sends it only with that permit. `setup_reads` is only for a tool asked for as `read_only`, which the vendor must mark read-only, and the owner's text names the read on its own line. (`…/describe`, which read a sign-in before its tools were approved, answers 410 too: kiff/planque#26.) |
| `POST /v1/accounts/{tenant}/grants/{id}/tools` | the agent that asked, body `{"tools": […]}` | sets a listed grant's tools, from its catalog; new text and digest, `pending` again. A tool name another of the account's connections holds (at any step but failed or disconnected) is refused (409), so the agent names it anew before the owner approves; a tool left unnamed gets planque's per-grant name |
| `POST /v1/accounts/{tenant}/grants/{id}/tools` on a **connected** grant | the agent that asked, body `{"tools": […]}`: only the tools to add | adds them to what it has (kiff/tafs#138): `pending` with the sign-in held, and the owner's text says what's added and that its tools keep working meanwhile. They do: the relay serves the grant's connected tools until the owner approves, which needs **no second sign-in** and registers only the added tools in KIFF, with a transport credential of their own. Not approved within a day, or refused by KIFF, the addition goes and the grant stays connected with what it had. A tool it already has is refused (409) |
| `GET /v1/accounts/{tenant}/agents/{agent}/grants` | the agent that asked | the grants it asked for that agent, not failed or disconnected: id, vendor, server, status, account, `kind` and tools (each with `read_only` and `setup_reads`, always sent: true for one of Setup's own reads, which the taf's Card never covers), so Setup, improving a taf, adds to the connection the taf has instead of asking for another. `kind` is how the grant reaches its vendor: `vendor_sign_in` (the vendor's own MCP server), `provider` (a provider's app, with its `entry`, and `mode` and `mode_words` as the grant's own view says them), `key` (a REST key, with `api`: its description's URL, base and where the key goes, never the key) or `browser`, so a new taf is offered the kind the account already uses for that vendor |
| `GET /v1/accounts/{tenant}/catalog?q=…` | a bearer KIFF key for the account | searches configured entries and optional provider catalogs, with pagination and this account’s connection state; see [catalog configuration](docs/connections.md#reviewed-vendor-catalog). With no search, a provider's apps are those with a sign-in set up; a search finds every app the provider has, an app that needs an operator's own developer credentials first as `state: "needs_setup"`. `partial: true` says the read stopped with pages left (kiff/planque#87) |
| `GET /v1/accounts/{tenant}/catalog/tools?q=…&apps=…&cursor=…` | a bearer KIFF key for the account | finds tools across the provider's apps by words (kiff/planque#87), narrowed to the apps named in `apps` (slugs, comma-separated, at most 20): each with its `entry`, `app`, `app_title`, `app_state` (`available` or `needs_setup`), `slug`, `name`, `description` and `read_only`, so Coach proposes apps before connecting any. A platform meta tool never comes back. `next_cursor` pages |
| `GET /v1/accounts/{tenant}/catalog/{entry}/events?cursor=…` | a bearer KIFF key for the account | a provider app's events, before any connection: `slug`, `name`, `description`, `version`, the filters' schema (`config`) and whether planque serves it as a door; Coach picks one as it picks an action, and the owner approves its exact filters on the connection as before |
| `GET /v1/accounts/{tenant}/grants/{id}/events?cursor=…` | the connection's own agent, or the one that asked | the same list for a connection's app |
| `POST /v1/accounts/{tenant}/catalog/{entry}/grants` | a bearer KIFF key, body `{"agent": …}` | requests reviewed entry tools or a managed app sign-in for the taf, with names allocated before approval; returns `{grant, needs_key}`, with no sign-in link before device approval. A managed app the provider signs people in to by itself is **set up on this pick** when it isn't yet: one provider-managed sign-in, made once under a lock and pinned by every grant on it, which the grant marks `on_demand: true`; at most 10 such set-ups an account a day (429). An app that needs our own developer credentials answers 409 until an operator enables it. Every grant says who holds the sign-in (`custody`: `provider` or `planque`). For a managed app, `tools`, optional, names the provider's tools the taf will use, as the tools search names them (at most 20): the sign-in then asks only for what they need, and serves only those tools; when the provider can't narrow it, it asks for the app's usual access, and the words say so ([managed connectors](docs/managed-connectors.md)) |
| `GET /v1/accounts/{tenant}/asked-grants` | the agent that asked | every grant it asked for, for any taf's agent, each with its agent, when it was asked and its `kind` (as above): so Setup that lost a taf's agent id finds its connection again instead of asking for another |
| `GET /v1/accounts/{tenant}/connections` | any of the account's keys or sessions, as for its catalog: a taf's own key reads every taf's connections, which suits one owner; with tafs shared between people it should narrow to the caller's own agent, or to sessions | how the account's tafs reach their vendors: each connection that holds a sign-in (`connected`, `listed`, `needs_reconnect`, or `pending` with its sign-in held, never one asked for and never signed in to), with its `agent`, `vendor`, `kind` and what asks for the same again (`server`, `entry` or `api`), its `status`, and a provider's `mode` and `mode_words`. No grant id, approval, text, tool or key. tafs cloud gives it to Coach when a conversation starts, for the tafs the person may see, so a new taf is offered the kind of connection the account already uses for a vendor |
| `POST /v1/accounts/{tenant}/grants/{id}/approve` | the phone, signed, body `{"digest": …}` | the owner's approval of exactly that request; returns the vendor's own sign-in link. For tools set on a grant whose sign-in is held, it connects at once, with no link. On a `needs_reconnect` grant, the same digest signed again returns a new sign-in link: the reconnect (below). On a `connected` sign-in grant, it returns a new sign-in link while the connection goes on working with the one it has: signing in again (below). On a managed connection still waiting on its sign-in, the approving device signing it again finishes it when the provider has it signed in (`{status: "signed_in", confirm, confirm_by}`, then `/confirm`), else returns a fresh link ([managed connectors](docs/managed-connectors.md)) |
| `POST /v1/accounts/{tenant}/grants` with `"kind": "api_key"` | an agent, with its KIFF key | a vendor with a REST API and no MCP server (kiff/planque#64). `api`: `{"url": its OpenAPI description, "auth": {"in": "header"\|"query", "name", "prefix"?}, "base"?}`. planque reads the description once (JSON, OpenAPI 3.0 or 3.1, at most 16 MB and 500 operations per snapshot; `api.operations` may select a subset from a larger catalog) and keeps it as the grant's snapshot: the owner signs its SHA-256 with the server and where the key goes, and every call uses the snapshot, whatever the URL serves later. Refused, with the reason: YAML, external `$ref`s, servers set per path or operation, an operation secured otherwise than the key, callbacks and links, and any way of sending the key but a header or the query. Each operation is a tool (its `operationId`, else its method and path), with its path, query and JSON body as arguments; `GET`/`HEAD` are reads. A vendor's refusal is never taken as evidence that nothing happened: the call stays counted in KIFF |
| `POST /v1/accounts/{tenant}/grants` with `"kind": "browser"` (or `POST …/catalog/browser/grants`) | an agent, with its KIFF key | a browser as a connection (kiff/planque#86): the taf uses websites in a browser planque runs for it alone, with planque's own tools and reviewed profile (`browse_open`, `browse_read`, `browse_act`, `browse_checkout`, `browse_pay`; `tools` names a subset, `[]` all of them). `browser`: `{"sites": [hosts the browser may open; none, any public site], "currency": the only currency `browse_pay` pays in, required with it}`. Approving on the phone connects at once: there is no vendor to sign in at. Money moves only through `browse_pay`, whose amount KIFF reads from the checkout page with `browse_checkout`; a `browse_act` step that would pay or place an order (the worker's fact or the control's words, in seven languages), any click on a control of a checkout or payment page, or typing into a password, card or code field, is refused and nothing more runs (`payment_step`, `secret_field`); a URL off the sites, a local address or a non-https URL, and the page a redirect lands on, are refused before anything is read (`site_not_listed`, `address_private`, `bad_url`); `browse_pay` sends nothing when no total was read (`total_unread`), the checkout changed (`checkout_changed`) or it is in another currency (`currency_other`). Secret fields are masked in every page. Offered only where a browser worker is configured; without one, 503 |
| `POST /v1/accounts/{tenant}/grants/{id}/sign-in` | the phone, signed, body `{"url"?}` | the live sign-in view of a browser grant (kiff/planque#86): planque starts a session on the worker with the grant's kept sign-ins and answers `{view, expires, words}`, a one-time URL good for ten minutes, which the app opens and nothing more. On the page the owner sees the browser as frames and taps, types and presses keys in it; a URL must be on the grant's sites. Tapping Done seals the session's cookies as the grant's sign-in; a view not finished in ten minutes lapses and keeps nothing. While a view is open, the taf's calls on the grant are refused (`signing_in`). No frame is kept, and no page is read for an agent |
| the browser worker (`cmd/browserworker`) | planque, with `PLANQUE_BROWSER_KEY` | the headless browser a browser grant uses: one Chromium and one egress proxy per session, gone at its end; reports a page's facts with secrets masked at the source; reached at `PLANQUE_BROWSER_URL` on a private address; its image is `deploy/browserworker/Dockerfile`. See [docs/connections.md](docs/connections.md#the-browser) |
| `POST /v1/accounts/{tenant}/grants/{id}/key` | the phone, signed, body `{"digest", "key_sha256"}`, the key in `Planque-Vendor-Key` | the owner's approval of an API-key grant, with the key: the signed body holds only the key's SHA-256, so the approval kept as evidence never holds the key. The key is sealed like a sign-in's token and never shown again. With no tools the grant is `listed` (its catalog is the snapshot's operations, each with its method and path); with tools it connects at once. Tools set later are approved with `…/approve`, with no second key. A `needs_reconnect` grant takes a new key the same way. A `connected` grant takes a replacement (gap 13 of kiff/tafs#476): planque reads once with the new key and keeps it only if that works, else nothing changes and the old key stays (400 when the vendor refuses it, 502 when it can't tell). The read is a read-only operation that needs no argument, the grant's own tools first, and the grant's event names it; an API with none can't take a replacement (409). A Telegram token is checked with `getMe` and must be the same bot's. The grant, its tools and its Card stay as they were, and a grant that reuses the key gets the new one. Disconnecting deletes the key and the snapshot |
| `POST /v1/accounts/{tenant}/grants` with `"reuse": "<grant id>"` | an agent, with its KIFF key | a second grant to the same vendor and account uses a connected grant's sign-in or key (kiff/planque#36): the owner's text says so, and approving needs no second sign-in or key. The token isn't copied: both grants use one, refreshed once. A grant that isn't connected, of another account or to another server can't be reused (409). When the first is disconnected, the second is `needs_reconnect` and signs in, or takes a key, on its own. A managed app's connected grant is reused with `tools: []`: the new grant is signed in and listed on the same provider account, its tools are picked from the catalog, and the owner approves them; the account stays while either uses it |
| `GET /oauth/callback` | the vendor | `state` used once, `iss` checked (required when the server says it sends it), code exchanged with PKCE and the exact redirect; the token is sealed, unused, and the browser goes on to the owner's app with a confirmation code |
| `POST /v1/accounts/{tenant}/grants/{id}/confirm` | the phone that approved, signed, body `{"confirm": …}` | the tools are registered in KIFF's gateway in `relay` mode; for a grant with no tools, the vendor's tools are listed instead (`listed`); for a reconnect, the grant is `connected` again with the new token, its tools in KIFF as they were; for a sign-in again on a connected grant, its new token replaces the one in use, and the old one is revoked at the vendor |
| `POST /mcp/{tenant}/{id}` | KIFF's gateway, with the relay's transport credential | MCP; each `tools/call` needs KIFF's permit (`KIFF-Permit`) for that exact call |
| *(a vendor signs planque out)* | | when the vendor refuses the refresh (`invalid_grant`), or a token expires with no refresh token, the grant is **`needs_reconnect`** ("Signed out by the vendor · sign in again"). Its tools stay in KIFF; each call through it is answered "not sent", saying the owner signs in again, and reaches no vendor. The owner's phone approves the same digest again, signs in (once the vendor named the grant's account, a sign-in to another account, or one that names none, is refused), and confirms: the grant is `connected` again. A new sign-in nobody confirms goes, and the grant waits again. planque tells KIFF of each sign-out once (`POST /v1/me/connections/signed-out`, with the connector key: the grant, the taf's agent, the vendor's name, why and when), so the people who see the taf hear of it; one KIFF didn't take is sent again by the sweep |
| *(signing in again on a connected grant)* | | the owner fixes a sign-in the vendor refuses without disconnecting (gap 13 of kiff/tafs#476). The phone approves the same digest, and planque returns a sign-in link: the grant stays `connected`, and its calls go on with the sign-in it has. The vendor's sign-in comes back held apart, in the same vendor account or not at all: once the vendor named the grant's account, one in another account, or one that names none, is revoked, and nothing changes. The approving phone confirms it, and it replaces the sign-in in use for every grant that uses it; the old one is revoked at the vendor, and the grant's tools and KIFF stay as they were. One not confirmed in 15 minutes is revoked and goes, and the connection stays as it was. A managed app can't sign in again in place, since planque can't tell which account a new sign-in is (409, saying to disconnect it and connect it again), and a grant on another's sign-in signs in again on that one |
| `DELETE /v1/accounts/{tenant}/grants/{id}` | the phone, signed | removes the tools from KIFF, revokes the token at the vendor, deletes it here |
| `GET /v1/accounts/{tenant}/grants`, `…/export` | the phone, signed (the export, by someone KIFF says may export the account: `account.export`, admins) | the grants in use with what the owner read and their signature, checkable without planque. A disconnected grant's record stays until the account ends (its vendor tokens are deleted at disconnect); the list leaves it out unless `?include=disconnected`, and the export has every grant. Each grant names its app's catalog `entry` when it's certain, so an app lists a connection under its app by `entry` alone ([which app a grant is for](docs/connections.md#which-app-a-grant-is-for)) |
| `DELETE /v1/accounts/{tenant}` | the phone, signed by someone KIFF says may delete the account (`account.delete`, admins), as KIFF's own deletion requires; tafs cloud forwards it, before KIFF's own deletion | deletes the account (kiff/tafs#185): each grant's tools leave KIFF's gateway, its vendor tokens are revoked where the vendor allows and deleted, and every record goes, the enrolment with the connector key last. Answers what KIFF or a vendor didn't confirm in `problems`. Asked again once nothing is kept, **204** |
| `POST /v1/account-deletions` | `{intent, signature}`: the owner's phone's one signed intent to delete the account (`tafs.delete-account.v2`: tenant, device, issued, expires at most 10 minutes later, nonce), which tafs cloud carries | **accepts** it, deleting nothing, once KIFF has: KIFF accepts first, from the phone's session, checking its approval key, and planque asks KIFF (`POST /v1/account-deletions/status`) that it accepted exactly this intent. Its nonce is spent once. A later intent takes an earlier one's place. **202** |
| `POST /v1/account-deletions/finish` | `{intent, signature}`, only the one accepted, which KIFF still vouches for | erases the account as `DELETE /v1/accounts/{tenant}` does, with no deadline; the same when asked again. **403** `not_accepted` for any other |
| `POST /v1/account-deletions/status` | `{intent, signature}` | whether that intent was accepted (`accepted`), or the account is `erased`, deleting nothing. **403** `not_accepted` otherwise |

A signed call carries `Planque-Device-Signature: t=<signed_at>,v1=<sig>,d=<device id>`:
the device's P-256 key (tafs' `ApprovalSigner`) signs, as base64 DER ECDSA
over SHA-256, the lines `planque.request.v1`, planque's base URL, tenant,
method, path, the body's SHA-256 and `signed_at`. A signed request is
accepted once, within 5 minutes, and only by the planque it names.

**One fingerprint, then taps (kiff/tafs#494).** Each signed route has a
kind. *Safe* routes take authority away or read (disconnect, the grants
list, a door's channels, deleting a trigger, the export): the device key
alone signs them. *Routine* routes give a little (approving a sign-in or a
key, confirming, creating a trigger, enrolling): the phone's approval key
signs them behind its fingerprint, as KIFF knows it, or the device key signs
them inside an approval window the approval key opened at KIFF: then the
header is `t=<signed_at>,v2=<sig>,w=<window id>` over the same lines under
`planque.request.v2` with the window's id last, and planque asks KIFF
whether that window is open for the device. A signature naming a window
that isn't open is refused, 403 with `"code": "window_required"`: the app
opens one and signs again. With `PLANQUE_WINDOW_REQUIRED=true`, a device-key
signature naming no window is refused the same way on a routine route;
until then it is taken as before, for apps from before the window. Each
grant's approval records which key signed (`key`: `device` or `critical`)
and the window, so an audit tells a tap from a fingerprint.

**KIFF keeps the devices.** planque keeps none: it asks
KIFF's `POST /v1/devices/verify`, with the account's connector key, whether
device `d` signed those bytes and whether its person may connect tools
(`tools.connect`). KIFF answers for its own account's devices only, and a
removed device, a person who left or one who can't connect tools is
refused. **So any member who may connect tools on the account (an admin or
an editor in KIFF) approves from any of their devices,** not only the
first phone set up. Each grant's record keeps who approved and on which
device (KIFF's person and device ids), and its events say it in the owner's
words, with the device's name from KIFF ("approved on tafs on iPad"), never an
id. KIFF says nothing about when, so planque keeps the 5-minute window
and the replay check. A phone enrolled with its own key before this signs
without `d=` and is checked against that key, until it signs as a KIFF
device.

**The sign-in is bound to the device that approved.** The vendor sends the
browser back to planque, which keeps the token unused and sends that
browser on to the owner's app (`PLANQUE_APP_RETURN`, set by the operator)
with a one-time code. The device that approved signs
`POST …/grants/{id}/confirm` with that code within 15 minutes; only then
are the tools registered. A link someone else opens hands its code to
their browser, never to the approver. Vendors are reached only at public
addresses, checked when each connection is made.

**Permits** are KIFF's execution permits: `typ` `kiff-permit+jwt`, `alg`
`Ed25519`, a key from KIFF's JWKS (refetched every 5 minutes, an unknown
`kid` at most once a minute, trusted at most an hour without a fresh
fetch), the relay's audience, the account, the tool, the RFC 8785 hash of
the arguments, and a lifetime of at most 60 seconds. Each operation runs
at most once: a replay gets the first answer, and a call whose answer was
lost is never sent again.

**Narrowing a vendor tool.** Some vendors offer one broad tool that writes to
any endpoint of their API. A grant
narrows it with `fixed_arguments` (sent on every call, never set by the
agent) and `arguments` (the only ones the agent sets, each placed at a path
in the vendor's arguments, at the top level where KIFF's Card can govern
them). A `title` says what the tool does in the owner's words: the approval
text, KIFF's Card and the record name the tool by it. The exact terms, the
narrowing included, are in the grant's `terms`, which the app shows under
Details. The phone signs the whole request either way:

```json
{"tool": "api_write", "name": "refund", "title": "Refund a payment", "amount_argument": "amount",
 "fixed_arguments": {"operation_id": "PostRefunds", "account": "acct_…", "livemode": false},
 "arguments": {"payment_intent": {"path": "parameters.payment_intent", "type": "string", "required": true},
               "amount": {"path": "parameters.amount", "type": "integer", "required": true}}}
```

The fields a narrowed tool may reach are the operation's own: planque reads
them with the vendor's describing tool and
refuses a path into `parameters` the operation doesn't take, with the list
it does, before the owner is asked. It reads nothing through a vendor's
tools before the owner approves them, so that tool is connected first (with
`setup_reads`), and a broad write is added once it's approved; a narrowed
broad write asked for any other way is refused. The read is on the grant's
record.

## Run it

```sh
go test ./...                    # every step against a fake vendor and a fake KIFF
deploy/aws.sh up                 # the hosted deployment
PLANQUE_BASE_URL=http://localhost:8080 PLANQUE_LOCAL_KEY=$(openssl rand -base64 32) \
  go run ./cmd/planque           # self-hosted, records in memory
```

Configuration is in `cmd/planque/main.go`. The hosted deployment's table
keeps point-in-time recovery on: each `deploy/aws.sh up` asks for it, so a
table rebuilt or changed by hand gets it back.

**The live check** (`live/`) runs against a payment provider's real MCP server
in test mode and the hosted planque. Until KIFF ships the connector key role, KIFF is
the stand-in in `internal/fakes`, deployed with `FAKEKIFF=1 deploy/aws.sh up`.
It enrols a phone, asks for the provider, approves, waits for the sign-in at
the provider, then: a refund inside the Card goes through, one outside it is
never sent, permits for other arguments, another account, another audience,
expired or missing send nothing, a replayed permit gets the first answer,
parallel calls on several instances refresh the token once (the
rotated refresh token is kept), disconnecting revokes the token at the provider
(which then answers it with 401), and the export's signature verifies.

## Licence

FSL-1.1-ALv2, like kiff/agents. You can read, run and self-host it; you
can't offer it as a competing hosted service for two years, after which
each version becomes Apache 2.0.

## Contributing and project policies

This public repository is source-available under [FSL-1.1-ALv2](LICENSE.md).
The license contains its future Apache-2.0 grant.

- [Contributing and checks](CONTRIBUTING.md)
- [Report a security issue privately](SECURITY.md)
- [Changelog](CHANGELOG.md)
- [Community conduct](CODE_OF_CONDUCT.md)
- [Instructions for coding agents](AGENTS.md)
- [Maintainer and release checklist](docs/maintaining.md)
