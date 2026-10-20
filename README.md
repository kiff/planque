# planque

**The safe place for your agents' keys.** An agent asks to act in a vendor
account (Stripe, Shopify, Gmail, any API). The owner approves on their own
phone, seeing exactly what will be granted. The vendor's keys stay in the
planque, in Europe or on the owner's own servers. Every call is made only
with a permit for that exact call.

*Planque* is French slang for a safe hiding place.

## For whom

People who let AI agents work in their business tools and need to know:
- who approved what;
- that the keys never sit in the agent;
- that each action stays within what was approved.

## How it works

planque is small on purpose. It doesn't keep a catalogue of vendors. It
follows the standard that modern vendors' MCP servers already speak:
OAuth 2.1 with discovery (RFC 9728, RFC 8414), dynamic client
registration (RFC 7591) and PKCE. Stripe's MCP server, for one, does all of
it.

1. **Ask.** An agent sends a grant request that says exactly what it wants:
   which vendor, what access, for which agent.
2. **Approve.** The owner approves in their own app, for example tafs, with
   Face ID. Their device signs that exact request with its key, and KIFF,
   which keeps the devices each person approves with, says whether it did.
   Nothing is granted without the signature.
3. **Connect.** planque discovers the vendor's sign-in server, registers
   itself, and sends the owner to the vendor's own sign-in page. The token
   lands in planque, encrypted with KMS, bound to that vendor account,
   refreshed safely, and revocable.
4. **Call.** For each call, planque checks a permit for that exact call (for
   example from KIFF, RFC 046), then calls the vendor itself. The agent never
   receives a token.

## Two ways to run it

- **Hosted:** we run it in the EU (AWS eu-west-1 to start). The owner runs
  nothing.
- **Self-hosted:** one container, on any cloud or server. For full
  sovereignty, the keys never leave the owner.

It's the same code either way.

## Done when (the first end-to-end)

- From kiff/tafs' setup chat, an owner connects **Stripe** through its MCP
  server:
  1. the approval on the phone;
  2. Stripe's own sign-in;
  3. the token in planque;
  4. a refund inside the Card goes through;
  5. a refund outside it is never sent.
- A call without a valid permit, with a permit for another call or account,
  or with an expired permit sends nothing. A test covers each case.
- A grant request the phone didn't sign creates nothing.
- Disconnecting revokes the token at the vendor and deletes it here.
- The signed approvals can be exported as evidence of human oversight
  (EU AI Act).

## The API

| Call | Who | What |
|---|---|---|
| `POST /v1/accounts/{tenant}/enrolment` | a device of the owner's that KIFF knows, signed (`d=`) | hands planque the account's KIFF **connector** key; any of the person's devices may, and the latest key is kept. A key that can do more than connect tools is refused. A signature without `d=` (a phone that enrolled its own key before KIFF kept devices) is no longer taken; a `device_key` in the body is accepted and ignored |
| `POST /v1/accounts/{tenant}/grants` | an agent, with its KIFF key | asks for a vendor's MCP server, some of its tools, for a taf; returns the plain words and the digest the owner signs; grants nothing. With `tools: []` it connects first: once signed in, planque only lists the vendor's tools |
| `GET /v1/accounts/{tenant}/grants/{id}` | an agent of the account | the grant's status, and `vendor`: the vendor's name as its MCP server gives it (`serverInfo` title, else name), else its host. Connected tools reach KIFF with the vendor's title for each (its MCP `title`, else `annotations.title`), for people to read; names a vendor gives are one line, at most 100 characters |
| `GET /v1/accounts/{tenant}/grants/{id}/catalog` | an agent of the account | for a grant asked for with no tools, once signed in (`listed`): each tool the vendor lists, with its description and arguments, and the account the sign-in acts in |
| `POST /v1/accounts/{tenant}/grants/{id}/describe` | the agent that asked, body `{"tool": "…"}` | while connecting first (`listed`), calls one vendor tool that the vendor marks read-only and that takes no arguments, and returns its answer. Such a tool describes the sign-in itself (which accounts, test or live), so the agent pins what the vendor reports instead of guessing (kiff/tafs#89). Only on a `listed` grant asked for with no tools after this change, whose approval text says the sign-in may be read; a connected grant, or one approved on the older text, is refused (409). Recorded on the grant before the call; at most 10 per grant; the vendor's answer is cut at 32 KiB. `readOnlyHint` is the vendor's word, not a guarantee: a vendor that mislabels a no-argument tool could have it change something, which is the owner's chosen vendor and their own sign-in, with no input to aim. A pinned argument's description in the tools approval is the vendor's text too, labelled "the vendor says" and capped at 200 characters |
| `POST /v1/accounts/{tenant}/grants/{id}/tools` | the agent that asked, body `{"tools": […]}` | sets a listed grant's tools, from its catalog; new text and digest, `pending` again. A tool name another of the account's connections holds (at any step but failed or disconnected) is refused (409), so the agent names it anew before the owner approves; a tool left unnamed gets planque's per-grant name |
| `POST /v1/accounts/{tenant}/grants/{id}/approve` | the phone, signed, body `{"digest": …}` | the owner's approval of exactly that request; returns the vendor's own sign-in link. For tools set on a grant whose sign-in is held, it connects at once, with no link |
| `GET /oauth/callback` | the vendor | `state` used once, `iss` checked (required when the server says it sends it), code exchanged with PKCE and the exact redirect; the token is sealed with KMS, unused, and the browser goes on to the owner's app with a confirmation code |
| `POST /v1/accounts/{tenant}/grants/{id}/confirm` | the phone that approved, signed, body `{"confirm": …}` | the tools are registered in KIFF's gateway in `relay` mode; for a grant with no tools, the vendor's tools are listed instead (`listed`) |
| `POST /mcp/{tenant}/{id}` | KIFF's gateway, with the relay's transport credential | MCP; each `tools/call` needs KIFF's permit (`KIFF-Permit`) for that exact call |
| `DELETE /v1/accounts/{tenant}/grants/{id}` | the phone, signed | removes the tools from KIFF, revokes the token at the vendor, deletes it here |
| `GET /v1/accounts/{tenant}/grants`, `…/export` | the phone, signed | every grant with what the owner read and their signature, checkable without planque |

A signed call carries `Planque-Device-Signature: t=<signed_at>,v1=<sig>,d=<device id>`:
the device's P-256 key (tafs' `ApprovalSigner`) signs, as base64 DER ECDSA
over SHA-256, the lines `planque.request.v1`, planque's base URL, tenant,
method, path, the body's SHA-256 and `signed_at`. A signed request is
accepted once, within 5 minutes, and only by the planque it names.

**KIFF keeps the devices** (kiff-cloud #1267). planque keeps none: it asks
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

**Permits** follow kiff-cloud RFC 046: `typ` `kiff-permit+jwt`, `alg`
`Ed25519`, a key from KIFF's JWKS (refetched every 5 minutes, an unknown
`kid` at most once a minute, trusted at most an hour without a fresh
fetch), the relay's audience, the account, the tool, the RFC 8785 hash of
the arguments, and a lifetime of at most 60 seconds. Each operation runs
at most once: a replay gets the first answer, and a call whose answer was
lost is never sent again.

**Narrowing a vendor tool.** Some vendors offer one broad tool: Stripe's
MCP server writes to any endpoint through `stripe_api_write`. A grant
narrows it with `fixed_arguments` (sent on every call, never set by the
agent) and `arguments` (the only ones the agent sets, each placed at a path
in the vendor's arguments, at the top level where KIFF's Card can govern
them). The owner reads the narrowing in the text they sign:

```json
{"tool": "stripe_api_write", "name": "stripe_refund", "amount_argument": "amount",
 "fixed_arguments": {"stripe_api_operation_id": "PostRefunds", "stripe_context": "acct_…", "livemode": false},
 "arguments": {"payment_intent": {"path": "parameters.payment_intent", "type": "string", "required": true},
               "amount": {"path": "parameters.amount", "type": "integer", "required": true}}}
```

## Run it

```sh
go test ./...                    # every step against a fake vendor and a fake KIFF
deploy/aws.sh up                 # Lambda + function URL, DynamoDB, KMS in eu-west-1
PLANQUE_BASE_URL=http://localhost:8080 PLANQUE_LOCAL_KEY=$(openssl rand -base64 32) \
  go run ./cmd/planque           # self-hosted, records in memory
```

Configuration is in `cmd/planque/main.go`.

**The live check** (`live/`) runs against Stripe's real MCP server in test
mode and planque on AWS. Until KIFF ships the connector key role, KIFF is
the stand-in in `internal/fakes`, deployed with `FAKEKIFF=1 deploy/aws.sh up`.
It enrols a phone, asks for Stripe, approves, waits for the sign-in at
Stripe, then: a refund inside the Card goes through, one outside it is
never sent, permits for other arguments, another account, another audience,
expired or missing send nothing, a replayed permit gets the first answer,
parallel calls on several Lambda instances refresh the token once (the
rotated refresh token is kept), disconnecting revokes the token at Stripe
(Stripe then answers it with 401), and the export's signature verifies.

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
