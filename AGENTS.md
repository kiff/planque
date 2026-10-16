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
