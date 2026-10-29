# Changelog

User-visible changes and compatibility notes are recorded here. Dates and
version headings refer to releases, not merely to merged pull requests.

## Unreleased

### Added

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
