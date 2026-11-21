# Typeform

Checked: 2026-11-14 · Sources: the connection service's Typeform tool
definitions as planque reads them (35 tools, [generated](generated/typeform.md)),
and Typeform's API reference.
Failure table checked against planque's answers on 2026-11-16.

## How it connects

- **A sign-in to Typeform**, held by the connection service. The grant's catalog
  lists the tools and, with a known sign-in, puts the ones its access doesn't
  cover under `unusable`. Offer only what's under `tools`.
- **Accounts that keep responses in Typeform's EU data centre** use a
  separate API address. A connection made to the usual address doesn't
  reach their forms.
- **One approval for the sign-in, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Returns |
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

**Creating a form works on the free plan,** and a form a taf makes is public by its link. A create past the Card's daily limit is refused as "not sent", and Coach then offers the owner a card to raise the limit: the taf never asks for more itself.

**Which account signed in:** the connection doesn't show it, and Google's chooser lists the owner's other accounts. A job that must be sure lists the forms first and shows them to the owner.

**Answers are matched by question ref, not by wording.** Read the form once,
keep each question's ref, and look answers up by it: an owner who rewords a
question doesn't break the taf.

## Don't offer

- `TYPEFORM_DELETE_RESPONSES`: it deletes what people sent, for good.
- `TYPEFORM_DELETE_FORM`, workspace and theme deletes.
- `TYPEFORM_UPDATE_FORM` on a live form, unless the job is about editing it:
  a full update replaces the whole form, and a question left out is deleted
  with its answers.
- Webhooks (`TYPEFORM_CREATE_OR_UPDATE_WEBHOOK`): the taf reads on its
  schedule.

## Access and plans

- **Responses per month are capped by plan;** each plan has its own limit.
- **A form with more than 1,000 responses** is read in parts, by date or
  from the last response handled.
- **Rate:** 2 calls a second per account.
- **Partial responses** (people who didn't finish) come back only when
  asked for. Collecting them on forms with logic or scoring needs one of
  Typeform's Growth plans.

## When it fails

The rows are Typeform's own answers: a refusal reaches the taf as "Typeform
refused it (it answered …): …", with Typeform's reason. planque's own
refusals, and an expired sign-in, are in the part below.

| What the taf gets | Why | Tell the owner |
|---|---|---|
| Refused for missing access | The sign-in didn't give this access, or the form belongs to another workspace | "Typeform is only partly connected, or that form isn't in your workspace." |
| Refused, not found | A deleted form, or a wrong id | "That form isn't there any more." |
| Refused, too many calls | Over the 2 calls a second | "Typeform asked us to slow down. The taf will continue in a moment." |
| Typeform's sign-in has expired (401) | The sign-in isn't accepted any more | "Typeform needs you to sign in again." |
| A result with no new responses | Nobody answered, or the plan's monthly cap is reached | "No new responses since last time." Say the cap only if Typeform says so. |
