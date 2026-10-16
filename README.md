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

## Licence

FSL-1.1-ALv2, like kiff/agents. You can read, run and self-host it; you
can't offer it as a competing hosted service for two years, after which
each version becomes Apache 2.0.
