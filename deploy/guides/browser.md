# Browser

Checked: 2026-11-21 · Sources: planque's browser tools and refusal codes,
and the plan for tafs in the owner's Chrome. Not yet run on an owner's
Chrome: check again after the first run.

## Which browser

One catalog entry, `browser`, covers both. Coach asks for the sites and the
tools; **the owner picks the browser on the phone's approval card**. Coach
never picks it.

- **The hosted browser:** planque runs it, keeping nothing between calls. It
  has **no sign-ins of the owner's**: it sees each site as a visitor does.
- **The owner's Chrome:** the tafs extension, in a separate Chrome profile
  named "tafs", paired from the phone. It uses the sign-ins the owner made
  in that profile, never their personal mail, bank or work sessions, and
  never sees or asks for a password. The owner's main profile isn't offered.
  - **Pairing:** on the phone, Account › Chrome › Pair a Chrome shows a
    short code; the owner types it into the extension and confirms on the
    phone with Face ID or the fingerprint. Several Chromes can be paired,
    each with its own name ("Work Mac").
  - **Coach never pairs** and never asks for the code.

A job that needs the owner's sign-ins needs their Chrome.

## Before offering

- **The job needs the owner's sign-ins and no Chrome is paired:** say how,
  and wait. Tell the owner: "This job needs your own sign-ins, so it runs in your
  Chrome. Pair it first: on your phone, open Account › Chrome ›
  Pair a Chrome. Then, in Chrome, add the tafs extension from the Chrome
  Web Store and type the code your phone shows. Use a separate Chrome
  profile for tafs, and sign in there only to the sites you want your tafs
  to use."
- **Sites:** only those the owner named. If the job needs more ("and
  the competitors too"), ask which; never add one.
- **Sensitive sites in Chrome** (a bank, a wallet, a password manager, a
  payment or health portal, a government account) stay blocked by the
  extension, even on the grant's list, until the owner allows that site by
  name on the phone. Say so before offering.
  Tell the owner: "{site} is a sensitive site, so a taf can use it in your
  Chrome only after you allow it by name on your phone."

## Tools a job usually needs

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `browse_open` | Opens a page: its text, links, fields and controls, each with a ref; secret fields masked | read | the refs, the page's read |
| `browse_read` | Reads a page's full text, to read rather than act | read | the text |
| `browse_act` | Fills fields, picks options and clicks, by the refs of one `browse_open` | write | the page after, its read |
| `browse_checkout` | Reads a checkout: site, lines, total and currency | read | the total held for the owner |
| `browse_pay` | Places the order on a checkout read with `browse_checkout`, once | write, pays | the outcome |

- Ask for `browse_open` and `browse_read` to read. Add `browse_act` only when
  the job fills in or submits something.
- Add `browse_checkout` and `browse_pay` only when the owner asks the taf to
  pay. Ask the currency (the grant pays in that one only), and say that
  **every payment waits for the owner's yes on the phone**, whatever the
  Card's budget. The yes covers the total `browse_checkout` read, for 15
  minutes.
- **Risky steps wait for the phone.** The card names the exact action and the
  site; a yes does one step, once.

## What the owner sees in Chrome

- The run's tabs are a background tab group named after the taf ("tafs ·
  Rae"); it never takes focus, and closes when the run ends.
- The side panel shows the run, with **Pause** and **Stop**. Touching the
  run's tab pauses it, with **Resume** in the panel.
- Unpairing ends it at once and wipes what the extension kept.

## What not to offer

- Anything on a site the owner didn't name.
- Typing a password, a card number or a code; solving a captcha; answering a
  page's "Are you sure…?" for the owner.
- Paying any way but `browse_pay`.
- The hosted browser in place of an absent Chrome: another grant, without
  the owner's sign-ins.

## When it fails

planque's browser refusals, by the answer's `code`:

| Code | What it means | Tell the owner |
|---|---|---|
| `site_not_listed` | Not one of the grant's sites; nothing was opened | "{site} isn't one of this taf's sites, so it stayed off it. Add it to the taf if you want it used." |
| `payment_step` | A click would pay outside `browse_pay`; it didn't run | Nothing: paying is `browse_pay`'s, if the Card allows it. |
| `secret_field` | A step would type into a password, card or code field; it didn't run | "{site} wants a password or a card number. The taf never types those: sign in yourself in your tafs Chrome profile if the job needs it." |
| `page_changed`, `ref_unknown` | The page isn't the read the steps came from; the rest didn't run | Nothing: read the page again, and build the steps from that read. |
| `checkout_changed` | The checkout changed after it was read; nothing was paid | "The checkout changed after you approved, so nothing was paid. The taf reads it again and asks you once more." |
| `total_unread` | No total and currency on the checkout; nothing was paid | "I couldn't read a total on {site}'s checkout, so nothing was paid." |
| `currency_other` | The checkout isn't in the grant's currency; nothing was paid | "{site}'s checkout isn't in {currency}, the currency you set, so nothing was paid." |
| `pay_not_approved` | A pay without the owner's yes | Ask for the yes on the phone; never pay another way. |
| `outcome_unknown` | An earlier pay on this checkout has no known result | Never pay again. "I'm not sure your order on {site} went through. Please check, and tell me in the app whether it happened." |
| `browser_unavailable` | The browser couldn't run; nothing was done | "The browser couldn't start. The taf tries again later. Nothing was done." |
| `bad_url`, `address_private` | Not an https address of a public site | Nothing: fix the address. |

**Coming with the Chrome routing** (not returned yet):

| Code | What it means | Tell the owner |
|---|---|---|
| `chrome_offline` | The owner's Chrome is closed or asleep; nothing was done, the run holds | "Your Chrome isn't reachable (closed, or the computer is asleep). This taf goes on when it is." Once per run. |
| `dialog_open` | A page dialog, with its text. An alert is dismissed; a confirm waits for the owner | "{site} is asking: \"{text}\". Answer it on your phone." |

**A sign-in page or a captcha:** "{site} needs you to sign in (or prove
you're a person) in your tafs Chrome profile. The taf waits."
