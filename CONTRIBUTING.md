# Contributing to planque

Read [README.md](README.md), [AGENTS.md](AGENTS.md), and
[LICENSE.md](LICENSE.md) before changing the project. This is a public,
source-available repository under FSL-1.1-ALv2, with the future Apache-2.0
grant described in the license. Public visibility does not remove its current
use restrictions.

## Discuss and scope the change

Search existing issues and pull requests first. Report ordinary bugs with the
issue template; report vulnerabilities privately through [SECURITY.md](SECURITY.md).
For a new feature, protocol change, or changed authority boundary, explain the
user need and proposed behavior on an issue before building it. Reuse an
existing issue for related work. A focused PR can fix several closely related
problems without creating an issue for each edit.

Keep changes in this project's responsibility.
planque holds vendor credentials and verifies exact-call permits; KIFF decides
authority and tafs presents the owner's approval.
Document cross-repository prerequisites using public interfaces and public
references that an outside contributor can read.

## Prepare a pull request

1. Work on a feature branch. Do not push directly to main.
2. Use synthetic fixtures. Remove credentials, customer data, personal
   identifiers, internal infrastructure details, and private assistant-session
   URLs from code, screenshots, logs, commits, and issue/PR text.
3. Explain the previous behavior, the resulting behavior, and why it matters.
   Link the relevant issue; use `Closes #N` only when the whole issue is done.
4. Add regression coverage for changed behavior, including relevant failure,
   retry, concurrency, and cross-account cases. For documentation-only edits,
   checking links and configuration syntax is sufficient.
5. Add a short entry under Unreleased in [CHANGELOG.md](CHANGELOG.md) for
   user-visible behavior, security fixes safe to disclose, or compatibility
   changes. Pure refactors and routine documentation corrections need no entry.
6. State exactly what was tested and what could not be tested. Describe any
   rollout, migration, or cross-repository dependency.

Run the checks used by CI:

```sh
gofmt -l .
go vet ./...
go vet -tags live ./live
go test -race ./...
```

`gofmt -l .` must print nothing. The normal suite uses fake services. Tests
under `live/` are separate, opt-in checks against real services; do not run
them with production credentials or against other people's accounts.

## Review and release

A maintainer reviews the final revision, with required CI passing. Authors do
not approve their own changes. Resolve requested changes with evidence and
obtain a new review when the implementation changes. Avoid unrelated generated
files, dependency updates, and formatting churn.

A merged PR is not evidence of a deployment, completed rollout, or successful
business operation. Record those separately when relevant. A docs-only PR does
not authorize a release, deployment, or change to customer accounts.

AI-assisted contributions follow the same rules. The person or agent submitting
the PR remains responsible for its accuracy, tests, licensing, and review.
