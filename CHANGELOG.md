# Changelog

User-visible changes and compatibility notes are recorded here. Dates and
version headings refer to releases, not merely to merged pull requests.

## Unreleased

### Changed

- An enrolment keeps the enrolled connector key while KIFF still takes it for
  the account, and answers `200 {"kept": true}`: a second device setting up no
  longer swaps the key, which left the tools the first key connected unowned
  in KIFF. A new key replaces the enrolled one only once KIFF no longer takes
  it.

- A broad tool narrowed to one operation (Stripe's `stripe_api_write`) takes
  only the operation's fields: planque reads them with the vendor's
  describing tool, once the owner approved it on the grant, and refuses a
  field the operation doesn't take, with the list it does. A narrowed broad
  write is added to a grant after its describing tool is approved; asked for
  at once, it is refused. Narrowed reads are unchanged.

### Fixed

- A read of planque's store is sent again when the connection drops before an
  answer, so one dropped connection no longer fails the owner's request. Writes
  are never sent again.
- An operation's outcome is always recorded, even when the vendor's answer is
  too large to keep; the answer is then kept as a short stand-in, and a replay
  still never sends again.

### Added

- An operator-installed, versioned vendor catalog lists account connection
  state and requests reviewed operations for a chosen taf. Catalog grants
  freeze their profile and final tool names before device approval. API-key
  entries refuse a source description changed since review.

- **Device-approved connection profiles:** a grant's digest covers its
  profile, and `/approve` and `/key` accept the phone-signed profile/words
  digests. Tools use the profile's measurements, including decimal-string
  amounts. The relay waits for KIFF to confirm the profile; a failed
  application rolls back registration. Required arguments and permit-derived
  markers are checked before sending, and operation records keep the profile
  revision used for their replay.

- **Registered vendor applications:** operator-configured MCP and REST
  applications use credentials from a sealed store or configured Secrets
  Manager references. The owner's approval pins the client identity,
  endpoints, access descriptions and reviewed operations. Secret rotation
  for that client works without changing the approval; replacing the
  client identity needs a new approval.

- **MCP client metadata documents:** a vendor that supports CIMD can use
  planque's public HTTPS client document instead of a registration request.
  The document contains only the configured client identity and redirect.
- **Larger API descriptions and selected operations:** descriptions may be
  up to 16 MB. `api.operations` selects up to 500 operations from a larger
  description; the selection and the complete source digest are included
  in the request the owner approves. An operation fingerprint distinguishes
  request changes from wording changes.

- **A write's amount read from its resource:** a tool may name `amount_lookup` (`{tool, argument,
  path}`), a read-only tool of the same request, the read's string argument that takes the call's
  `resource_argument`, and where the amount is in its answer. planque registers it in KIFF like
  `unit_lookup`, by the read's name there; KIFF reads the amount before deciding, and the call
  carries none, so a Card counts the amount the vendor says, never one the agent sends. It needs
  `amount_argument` and `resource_argument`, and the amount can't also be one of a narrowed tool's
  arguments.
- **API-key grants** (kiff/planque#64): a vendor with a REST API and an OpenAPI description, and no
  MCP server, is connected with a key the owner gives on their phone. The description is read
  once into a snapshot the owner approves by its digest; each operation is a tool, called through
  the relay with KIFF's permit like any other. A 4xx is the vendor's refusal; a 5xx, a redirect or a
  lost answer is unknown and never sent again. The key is sealed, never shown again, and deleted on
  disconnect.
- **A vendor's refusal is never evidence that nothing happened:** an API-key write that the vendor
  refuses answers the agent with the vendor's status and words, and KIFF keeps it counted. No tool
  is registered with an outcome contract: a status alone can't say nothing ran (a batch endpoint
  can answer 400 after part of it did). This replaces the brief release on a fixed list of 4xx
  statuses.
- **Reusing a sign-in** (kiff/planque#36): `reuse` names a connected grant whose sign-in or key a
  second grant uses, approved with no second sign-in.

- A vendor's sign-out is told to KIFF once, when the grant turns
  `needs_reconnect`, so the people who see the taf hear of it at once. A
  notice KIFF didn't take is sent again by the hourly sweep, for up to a week;
  one whose grant was signed in again is dropped.
- Every grant carries `status_words`: its status in the owner's words, so the
  apps show the same line for a connection and keep no words of their own.
- A connected grant gains tools with one approval and no second sign-in,
  while its tools keep working (kiff/tafs#138). The asking agent lists the
  grants it asked for a taf's agent, to add to them.
- Deleting the account, signed from the owner's phone, erases everything
  planque keeps for it: tools out of KIFF, vendor tokens revoked, records and
  the enrolment deleted (kiff/tafs#185).
- Security reporting, contribution and agent guidance, issue and pull-request
  templates, code owners, and dependency-update configuration.

Earlier history is available in [merged pull requests](https://github.com/kiff/planque/pulls?q=is%3Apr+is%3Amerged)
and [releases](https://github.com/kiff/planque/releases). This file does not
retroactively assign release dates or claim that merged work was deployed.
