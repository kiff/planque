# AGENTS.md

For AI assistants working on this repository. Read this first, then the README.

## Rules

- **This repository is public, and so are its Actions logs.** Never write
  people's emails or phones, account, key or team IDs, ARNs, or any secret
  into code, commits, PRs, issues or workflow inputs. Keys live in KMS or the
  environment, never in code.
- **Lean MVP, by value.** Each piece of work is a feature someone can use
  end to end, with a test that proves it. Build what's needed for it now,
  not what a platform might need later.
- **The agent never holds a vendor key.** planque calls the vendor itself, and
  only with a valid permit for that exact call. Nothing is granted without the
  owner's device-signed approval.
- **Small and portable.** Go and the standard library first. One binary that
  runs hosted (AWS eu-west-1) or self-hosted. Pin any dependency exactly.
- **Docs describe what this is and does.**

## Public contribution and review rules

Read [CONTRIBUTING.md](CONTRIBUTING.md) for scope and checks,
[SECURITY.md](SECURITY.md) for private vulnerability reporting, and
[docs/maintaining.md](docs/maintaining.md) for release checks.

- Follow the user's authorized scope. Audit work does not authorize fixes,
  deployment, merging, or changes to customer accounts. An explicit user
  instruction takes precedence over general workflow guidance here.
- Before posting an issue, comment, PR, or review, verify the authorized
  GitHub identity. Do not impersonate another maintainer. An author cannot
  approve their own PR; request an independent review of the final revision.
- Search existing work before opening an issue or changing shared code. Use
  an isolated branch or worktree and preserve other contributors' edits.
- Public commits, PRs, issues, fixtures, screenshots, and Actions logs must
  contain no live credentials, customer data, private infrastructure details,
  or private assistant-session links. Use synthetic evidence. Report new
  vulnerabilities through the private channel, not a public exploit report.
- State what checks actually ran, their results, and any limitations. Do not
  substitute prompt wording tests for measured agent behavior or a happy-path
  test for approval, failure, retry, concurrency, and account-isolation checks.
- Update Unreleased in [CHANGELOG.md](CHANGELOG.md) for user-visible or
  compatibility changes. Do not claim a deployment, release, or security fix
  is complete solely because code merged.
