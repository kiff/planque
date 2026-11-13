# Wise

Checked: 2026-11-13 · Sources: the connection service's Wise tool
definitions as planque reads them (9 tools, [generated](generated/wise.md)),
and Wise's API reference.

## How it connects

- **By a key**, made by the owner in Wise's settings (Integrations and tools,
  API tokens), pasted once. Wise lets a token be **read only**: ask for that.
- **The tools only read and quote.** None moves money: there's no transfer
  in this app, by design.
- **One approval for the key, one for the tools and the Card.**

## Tools

| Tool | What it does | Read or write | Keep |
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

**A quote isn't a transfer.** It expires after a while and costs nothing.
The owner makes the transfer in Wise.

## Don't offer

- Any job that sends money: this app can't, and a key that can should never
  be pasted for a taf.
- A full-access token when a read-only one does.

## Access and plans

- **No paid plan;** Wise charges per transfer, not for the API.
- **Personal tokens can't do everything in the EU and UK:** some reads of
  balance statements need an extra check the owner does in Wise. When a read
  is refused for this, the owner sees why in Wise.
- **Rate:** Wise limits calls per second and doesn't publish every limit; a
  daily read stays far below it.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| The key is refused (401) | Wrong, revoked or expired | "Wise didn't accept the key. Make a new read-only token in Wise's settings and paste it again." |
| Refused (403) | The token can't read this, or Wise wants an extra check | "Wise wants an extra check for this read. Open Wise to confirm, then I'll try again." |
| A quote is refused (422) | A currency pair or amount Wise doesn't take | "Wise can't price that: the currency or amount isn't supported." |
| Too many calls (429) | Over Wise's rate | "Wise asked us to slow down. The taf will try again in a minute." |

After any refusal the taf stops that step and says so.
