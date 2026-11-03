# Managed connections

`internal/connector` defines the boundary to a connection platform. The
Composio adapter uses its v3.1 HTTP API and the Go standard library. It
reuses the platform's app catalog, managed authentication, tool definitions
and concrete execution API. It does not run another agent framework.

The adapter supports paginated app and tool discovery, a toolkit-scoped
sign-in session, account inspection and one versioned action per execution.
It requires explicit account selection. Before execution it checks the
account's owner, app and active status, then checks the selected tool's
pinned version and schema. Platform meta-executors are refused. A lost or
malformed execution answer remains unknown and is never automatically
sent again. A refusal provides no proof that the action had no effect.

Project credentials are injected server-side. `SecretFromStore` reads one
operator-configured secret containing `api_key`, with a one-minute cache.
No request can choose a secret reference. Credential requests follow no
redirects, and backing-service error bodies are not returned.

With Composio-managed authentication, Composio holds the vendor OAuth
credentials and handles refresh. The adapter requests upstream revocation
when disconnecting the exact selected account. A successful deletion
immediately disables that connection, while upstream revocation remains
pending; the response keeps the background job id when supplied.

## Configure the provider

Set `PLANQUE_COMPOSIO_SECRET` to a Secrets Manager reference containing
`{"api_key":"..."}`. The runtime needs permission to read that exact secret.
The reference enables the Composio provider; without it, existing direct
MCP, registered OAuth and API-key connections continue to work independently.
No agent request can enable a provider or supply its project key.

## Connect an app

The account-authenticated catalog endpoints include upstream discovery:

1. `GET /v1/accounts/{tenant}/catalog?provider=composio&q=...` returns
   `entries`, `provider` and `next_cursor`. Pass the same search/provider and
   `cursor=...` for the next page. Configured reviewed entries remain available.
2. `POST /v1/accounts/{tenant}/catalog/connector:composio:{app}/grants`
   accepts only `{"agent":"..."}`. It returns a pending grant and no sign-in
   link. Its text explains that the provider holds the vendor credentials.
3. The owner's device signs `/approve` with that grant's digest. Only then
   is a scoped hosted sign-in created. The callback independently checks the
   server-allocated user, app and exact connection; callback account/status
   parameters establish no authority. The approving device confirms the
   browser's one-use code through the existing `/confirm` route.
4. After confirmation, the asking agent reads
   `GET /v1/accounts/{tenant}/grants/{id}/catalog?cursor=...`. Tool discovery
   is paginated and cached separately, rather than copied into the grant.
   Each tool has a stable gateway-compatible alias. Read the relevant page
   before selecting actions; discovery snapshots expire after 24 hours.
5. `/tools` selects the concrete aliases and supplies a drafted profile
   (`source: "drafted"`, `transport: "rest"`, `pin: {"kind":"none"}`).
   The grant digest covers the exact connected account, upstream tool
   versions and schemas as well as the profile. Upstream catalog membership
   never establishes reviewed profile status or test evidence. At most 20
   actions and 64 KB of selected upstream definitions fit one grant.
6. The device signs the updated digest and checked profile/words digests.
   Registration and KIFF profile application finish before any action can
   execute. Each action still requires an exact-call permit and the relay's
   once-only operation record. Permit expiry is checked again immediately
   before the provider execution request, after upstream preflight.

Disconnect disables the specific provider connection and reports upstream
credential revocation as pending. Unfinished sign-ins and unapproved held
connections are disabled by the existing sweep. Replacing a connected
account requires disconnecting and requesting a new grant; an existing
approved profile never silently switches accounts.

Synthetic tests exercise real signed HTTP handlers, the Composio HTTP
adapter and the KIFF test service, including account isolation, profile
application, exact-call permits, concurrent replay, lost answers and expiry.
They do not establish live app compatibility. A configured project and
owner-controlled test account are needed to validate an individual app.

Protocol references:

- [App catalog](https://docs.composio.dev/reference/api-reference/toolkits/getToolkits)
- [Session configuration](https://docs.composio.dev/docs/configuring-sessions)
- [Authentication and token custody](https://docs.composio.dev/docs/authentication)
- [Versioned tool definitions](https://docs.composio.dev/reference/api-reference/tools/getToolsByToolSlug)
- [Concrete execution](https://docs.composio.dev/reference/api-reference/tools/postToolsExecuteByToolSlug)
- [Account deletion and upstream revocation](https://docs.composio.dev/reference/api-reference/connected-accounts/deleteConnectedAccountsByNanoid)
