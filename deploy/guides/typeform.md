# Typeform

Checked: 2026-11-13 · Sources: the connection service's Typeform tool
definitions as planque reads them (35 tools, [generated](generated/typeform.md)),
and Typeform's API reference.

## How it connects

- **A sign-in to Typeform**, held by the connection service. The catalog
  lists the tools and, with a known sign-in, puts the ones its access doesn't
  cover under `unusable`. Offer only what's under `tools`.
- **One approval for the sign-in, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `TYPEFORM_LIST_FORMS` | Lists the owner's forms | read | **the form's id** |
| `TYPEFORM_GET_FORM` | Reads a form's questions | read | each question's ref |
| `TYPEFORM_GET_FORM_RESPONSES` | Reads responses, newest first, since a time | read | **the last response's token or time** |
| `TYPEFORM_GET_ALL_RESPONSE_FILES` | Downloads files people uploaded | read | |
| `TYPEFORM_CREATE_FORM` | Creates a form | write | the form's id |
| `TYPEFORM_UPDATE_FORM` / `TYPEFORM_PATCH_FORM` | Changes a form | write | |

**What to keep:** the form's id (in its link); the time or token of the last
response handled, so the next run reads only newer ones and never handles a
response twice.

**Answers are matched by question ref, not by wording.** Read the form once,
keep each question's ref, and look answers up by it: an owner who rewords a
question doesn't break the taf.

## Don't offer

- `TYPEFORM_DELETE_RESPONSES`: it deletes what people sent, for good.
- `TYPEFORM_DELETE_FORM`, workspace and theme deletes.
- `TYPEFORM_UPDATE_FORM` on a live form, unless the job is about editing it:
  a full update replaces every question, and a question removed takes its
  answers out of the results.
- Webhooks (`TYPEFORM_CREATE_OR_UPDATE_WEBHOOK`): the taf reads on its
  schedule.

## Access and plans

- **Responses per month are capped by plan** (10 on the free plan, more on
  paid ones). Beyond the cap, new responses aren't collected, whatever the
  taf does.
- **Responses** answer up to 1,000 a page.
- **Rate:** about 2 calls a second per account.
- **Partial responses** (people who didn't finish) come back only if asked
  for and only on some plans.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403) | The sign-in didn't give this access, or the form belongs to another workspace | "Typeform is only partly connected, or that form isn't in your workspace." |
| Not found (404) | A deleted form, or a wrong id | "That form isn't there any more." |
| Too many calls (429) | Over Typeform's rate | "Typeform asked us to slow down. The taf will continue in a moment." |
| The sign-in no longer works (401) | Removed, or expired | "Typeform needs you to sign in again." |
| No new responses | Nobody answered, or the plan's monthly cap is reached | "No new responses since last time." Say the cap only if Typeform says so. |

After any refusal the taf stops that step and says so.
