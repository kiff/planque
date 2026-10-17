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
   Face ID. The phone signs that exact request with its device key.
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
| `POST /v1/accounts/{tenant}/enrolment` | the owner's phone, signed with the key it enrols | binds the phone's key and planque's KIFF **connector** key to the account; a key that can do more than connect tools is refused |
| `POST /v1/accounts/{tenant}/grants` | an agent, with its KIFF key | asks for a vendor's MCP server, some of its tools, for a taf; returns the plain words and the digest the owner signs; grants nothing |
| `GET /v1/accounts/{tenant}/grants/{id}` | an agent of the account | the grant's status |
| `POST /v1/accounts/{tenant}/grants/{id}/approve` | the phone, signed, body `{"digest": …}` | the owner's approval of exactly that request; returns the vendor's own sign-in link |
| `GET /oauth/callback` | the vendor | `state` used once, `iss` checked (required when the server says it sends it), code exchanged with PKCE and the exact redirect; the token is sealed with KMS, unused, and the browser goes on to the owner's app with a confirmation code |
| `POST /v1/accounts/{tenant}/grants/{id}/confirm` | the phone that approved, signed, body `{"confirm": …}` | the tools are registered in KIFF's gateway in `relay` mode |
| `POST /mcp/{tenant}/{id}` | KIFF's gateway, with the relay's transport credential | MCP; each `tools/call` needs KIFF's permit (`KIFF-Permit`) for that exact call |
| `DELETE /v1/accounts/{tenant}/grants/{id}` | the phone, signed | removes the tools from KIFF, revokes the token at the vendor, deletes it here |
| `GET /v1/accounts/{tenant}/grants`, `…/export` | the phone, signed | every grant with what the owner read and their signature, checkable without planque |

A signed call carries `Planque-Device-Signature: t=<signed_at>,v1=<sig>`:
the phone's P-256 key (tafs' `ApprovalSigner`) signs, as base64 DER ECDSA
over SHA-256, the lines `planque.request.v1`, planque's base URL, tenant,
method, path, the body's SHA-256 and `signed_at`. A signed request is
accepted once, within 5 minutes, and only by the planque it names.

**The sign-in is bound to the phone that approved.** The vendor sends the
browser back to planque, which keeps the token unused and sends that
browser on to the owner's app (`PLANQUE_APP_RETURN`, set by the operator)
with a one-time code. The phone that approved signs
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
