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
