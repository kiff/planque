# Maintaining planque

## Repository controls

Verify these in GitHub settings; documentation files alone do not enforce them:

- Private vulnerability reporting, secret scanning, push protection, Dependabot
  alerts, and Dependabot security updates.
- A main-branch ruleset requiring pull requests, at least one independent
  approval, review of the latest changes, resolved review conversations, and
  the relevant CI checks. Restrict force pushes and branch deletion.
- Read-only default workflow token permissions. Give a job additional access
  only where needed; never expose secrets to untrusted pull-request code.

[GitHub's repository security guide](https://docs.github.com/en/code-security/getting-started/quickstart-for-securing-your-repository)
explains the settings. CODEOWNERS routes reviews; it does not replace a ruleset.
Dependabot opens update PRs; updates still require review and CI. Do not add
blanket auto-merge or workflow approval for external contributions.

## Before merging

Check the final diff, review state, and required CI. Resolve requested changes;
a change after approval may require another review. Check dependency and
workflow changes for changes in permissions or secret access. Keep unrelated
work out of the PR.

For changes to grants, device verification, permits, token handling, relay
execution, or expiry, review failure and concurrent paths as well as success.
Use fake services for regression coverage before any separately authorized
live check.

## Before a release

1. Resolve release-blocking findings or document the evidence that rules them out.
2. Review Unreleased in CHANGELOG.md; include behavior, compatibility changes,
   migration or dependency order, and disclosed security fixes.
3. Tag and describe the exact reviewed commit when making a versioned release.
   Keep build, source version, and deployed revision traceable. Do not label
   unverified work as shipped.
4. Verify the resulting artifact and deployment using authorized test accounts.
   Do not close a rollout issue solely because a PR merged.

Use [SECURITY.md](../SECURITY.md) for private reports and coordinated disclosure.

## Running and deploying

```sh
deploy/aws.sh up                 # the hosted deployment
```

**Each merge to `main` deploys planque's code** (CI's `deploy` job, once the repository
variable `PLANQUE_DEPLOY` is `on` and the secret `AWS_DEPLOY_ROLE` names the deploy role):
the function is built as `deploy/aws.sh` builds it, its code replaced, and planque checked
live (`/healthz` answers 200, an unsigned read of a catalog 401). If it doesn't answer so,
the code it ran before is put back. Code only: the function's configuration, table, keys
and schedule stay as `deploy/aws.sh up` set them, so a change to those still needs it.

Configuration is in `cmd/planque/main.go`. The hosted deployment's table
keeps point-in-time recovery on: each `deploy/aws.sh up` asks for it, so a
table rebuilt or changed by hand gets it back.

**The live check** (`live/`) runs against a payment provider's real MCP server
in test mode and the hosted planque. Until KIFF ships the connector key role, KIFF is
the stand-in in `internal/fakes`, deployed with `FAKEKIFF=1 deploy/aws.sh up`.
It enrols a phone, asks for the provider, approves, waits for the sign-in at
the provider, then: a refund inside the Card goes through, one outside it is
never sent, permits for other arguments, another account, another audience,
expired or missing send nothing, a replayed permit gets the first answer,
parallel calls on several instances refresh the token once (the
rotated refresh token is kept), disconnecting revokes the token at the provider
(which then answers it with 401), and the export's signature verifies.

## The first end-to-end

What the first end-to-end had to do.

- From kiff/tafs' setup chat, an owner connects a payment provider through its
  MCP server:
  1. the approval on the phone;
  2. the provider's own sign-in;
  3. the token in planque;
  4. a refund inside the Card goes through;
  5. a refund outside it is never sent.
- A call without a valid permit, with a permit for another call or account,
  or with an expired permit sends nothing. A test covers each case.
- A grant request the phone didn't sign creates nothing.
- Disconnecting revokes the token at the vendor and deletes it here.
- The signed approvals can be exported as evidence of human oversight.
