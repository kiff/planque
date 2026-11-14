# Wise

Checked: 2026-11-14 · Sources: the connection service's Wise tool
definitions as planque reads them (9 tools, [generated](generated/wise.md)),
and Wise's API reference.

## How it connects

- **By a personal API token** from a Wise business account, pasted once.
  The owner makes it in Wise's account settings (Accounting integrations,
  Developer tools, API tokens), on the website, with 2-step verification
  set up. It stays active until the owner deletes it there.
- **Ask for a read-only token.** Wise lets the owner set a token's access
  to read only or full when making it, and change it later. Ask for read
  only when the job only reads, which is every tool here. Treat a
  full-access token as able to move money, keep any token to this one
  connection, and the owner can delete it in Wise when the job ends.
- **The tools only read and quote.** None moves money: there's no transfer
  in this app, by design.
- **One approval for the key, one for the tools and the Card.**

## Tools

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `WISE_LIST_PROFILES` | Lists the owner's personal and business profiles | read | **the profile's id** |
| `WISE_LIST_BALANCES` | Reads each currency balance | read | |
| `WISE_LIST_ACTIVITIES` | Lists recent activity: payments in and out | read | |
| `WISE_LIST_TRANSFERS` | Lists transfers and their status | read | transfer ids |
| `WISE_LIST_RECIPIENTS` | Lists saved recipients | read | |
| `WISE_GET_EXCHANGE_RATE` | Reads the current rate between two currencies | read | |
| `WISE_LIST_CURRENCIES` | Lists the currencies Wise handles | read | |
| `WISE_CREATE_QUOTE` | Prices a transfer: rate, fee, arrival | write, but moves nothing | the quote's id |
| `WISE_GET_QUOTE` | Reads a quote again | read | |

**What to keep:** the profile's id. A business owner has two profiles; the
taf names which one it reads, in words, on the Card.

**A quote isn't a transfer.** It holds the rate for 30 minutes and moves
nothing. The owner makes any transfer in Wise.

## Don't offer

- Any job that sends money. None of these tools makes or funds a transfer,
  and a taf never asks for one.
- A full-access token for a job that only reads: the token itself could
  create transfers, even though the taf's tools can't.

## Access and plans

- **A personal token reaches only part of Wise's API:** quotes,
  recipients, transfers and their status. Balance statements aren't
  available to it, except for accounts based in the US, Canada, Australia,
  New Zealand, Singapore and Malaysia.
- **In the UK and Europe, some actions need a stronger check** (reading a
  statement, say). Wise refuses the first call and asks for it.
- **Rate:** not checked here; a daily read makes few calls.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The token is refused (401) | Wrong, or deleted in Wise | "Wise didn't accept the token. If you still want this job, make a new one (read only) in Wise's settings and paste it again; nothing was moved." |
| Refused (403) | The token can't reach this, or Wise wants a stronger check | "Wise doesn't let this connection do that step, or wants a stronger check that only you can give in Wise. The taf will do without it." |
| A quote is refused | A currency pair or amount Wise doesn't take | "Wise can't price that: the currency or amount isn't supported." |
| Too many calls | Over Wise's rate | "Wise asked us to slow down. The taf will try again in a minute." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
