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
  without contacting a vendor. It returns `{entries, by_key: true}`. Each
  entry has `entry`, `vendor`, `title`, optional `icon`, `how` (`oauth` or
  `key`), `state` (`available`, `connected` or `needs_reconnect`) and, for an
  active connection, `grant`. State comes only from that account's grants.
- `POST /v1/accounts/{tenant}/catalog/{entry}/grants` accepts only `{agent}`,
  identifying the taf's KIFF agent. It returns `{grant, needs_key}`, with the
  full grant view. The server supplies the reviewed tools and profile and
  allocates their final names before computing the approval digest. This
  request creates no sign-in link and grants no tools. The existing signed
  approval or key flow, including KIFF's profile and words digests, follows.

The catalog starts empty when no file is configured. Installing definitions
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
