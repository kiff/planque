## What the taf gets back, for every app

planque adds this part to every guide. It says what a failed call returns,
whatever the app, so Coach can tell the owner what actually went wrong. The
app's own table above says what {app}'s answers mean.

**Five kinds of answer, told apart by their first words:**

- **"Not sent: …" or "Refused by planque: …":** planque stopped the call, and
  nothing reached {app}. The answer's `code`, when it has one, says why (the
  table below). Nothing to undo.
- **"{app} refused it (it answered 403): …"** (or {app}'s address in place of
  its name, "refused it (403)" or "returned an error"): {app} itself refused
  the call, with its status and its own reason. The table above says what
  that means for {app}. Read before trying again: a refused write is taken
  as not done.
- **"The connection provider refused it":** the service between planque and
  {app} refused, not {app}. Try once more later. Then: "The connection to
  {app} had a problem, and the step didn't go through. The taf will try
  again later."
- **"{app}'s sign-in has expired (it answered 401)":** {app} no longer takes
  the sign-in. Every taf using it waits until the owner signs in again.
  Tell the owner: "{app} needs you to sign in again: tap Connect on this
  taf. Nothing was done."
- **"did not return a confirmed result", "may or may not have done it", or
  "its answer was lost":** the call may have reached {app}. planque never
  sends it again. Read {app} to see whether it happened before any new call,
  and never repeat a write blindly.

A call past the taf's Card never reaches planque: KIFF holds it for the
owner. That's an approval to ask for, not a failure.

**When planque refuses (the answer's `code`):**

| Code | Why | Tell the owner |
|---|---|---|
| `not_granted` | The tool isn't one the owner approved for this taf | "This step needs something you haven't allowed this taf to do in {app}. Want to add it?" |
| `waits_for_approval` | A tool the taf asked to add, not approved yet | "This step waits for your approval in the app." |
| `waits_for_key` | The tool needs a key the owner hasn't given yet | "This step needs your {app} key first." |
| `operation_changed` | {app} changed a tool since the owner approved it, so the connection is paused | "{app} changed how this step works since you approved it, so it's paused. Please look at it again and approve it." |
| `connection_unavailable` | The connection, or its {app} account, can't be used now | "I couldn't use your {app} connection just now. Nothing was done. The taf will try again later; if it keeps failing, connect {app} again." |
| `tool_not_served` | planque never runs this tool: it moves money, or closes a bill for good | Don't offer it. If a job needs it, the owner does it in {app}. |
| `halted`, `halted_before` | The owner halted this taf, or the whole account | "Your taf is halted, so nothing was sent to {app}." |
| `halt_unknown`, `record_unavailable`, `permit_expired`, `definition_unread` | A check before sending couldn't finish | Try once more. Then: "That step didn't go through, and nothing was sent. The taf will try again later." |
| `not_listening` | The call finishes on an event from {app} that the taf doesn't listen for | Subscribe to the event first; tell the owner only if that's refused. |
| Any other code | A check before sending didn't pass | "That step didn't go through, and nothing was sent to {app}." |

After any of these the taf stops that step and says so. A refused read
never leads to a write.
