# planque

planque keeps the keys and sign-ins AI agents use for other apps (payments, a
shop, email), so the agent never holds them. The owner approves each
connection on their phone, and each call goes out only with a permit for that
exact call.

*Planque* is French slang for a safe hiding place.

## For whom

People who let AI agents work in their business tools and need to know:
- who approved what;
- that the keys never sit in the agent;
- that each action stays within what was approved.

## How it works

1. **Ask.** An agent sends a grant request that says exactly what it wants:
   which vendor, what access, for which agent.
2. **Approve.** The owner approves in their own app, for example tafs, with
   their fingerprint or face. Their device signs that exact request with its key, and KIFF,
   which keeps the devices each person approves with, says whether it did.
   Nothing is granted without the signature.
3. **Connect.** planque sends the owner to the vendor's own sign-in page. The
   token lands in planque, encrypted, bound to that vendor account, refreshed
   safely, and revocable.
4. **Call.** For each call, planque checks a permit for that exact call (for
   example KIFF's execution permit), then calls the vendor itself. The agent
   never receives a token.

## Two ways to run it

- **Hosted:** we run it in the EU. The owner runs nothing.
- **Self-hosted:** one container, on any cloud or server. Direct connections
  keep credentials on your own servers; [managed connections](docs/managed-connectors.md)
  keep them with the configured provider.

It's the same code either way.

## Try it

```sh
go test ./...                    # every step against a fake vendor and a fake KIFF
PLANQUE_BASE_URL=http://localhost:8080 PLANQUE_LOCAL_KEY=$(openssl rand -base64 32) \
  go run ./cmd/planque           # self-hosted, records in memory
```

## Docs

- [The API](docs/api.md): every route, who may call it, and what it takes.
- [Connections](docs/connections.md): the sign-in standards, keys, the browser,
  connection profiles and the reviewed catalog.
- [Managed connections](docs/managed-connectors.md).
- [Maintaining](docs/maintaining.md): deploying, the live check and releases.

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
