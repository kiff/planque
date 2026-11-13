# WhatsApp

Checked: 2026-11-13 · Sources: the connection service's WhatsApp tool
definitions as planque reads them (127 tools, [generated](generated/whatsapp.md)),
and Meta's WhatsApp Business Platform reference.

## What it is

**WhatsApp's business platform, not the owner's personal WhatsApp.** A taf
sends from a **business phone number** registered with Meta, never from the
owner's own phone or chats. It can't read the owner's personal messages.

## How it connects

- **A sign-in to the owner's Meta business account,** held by the connection
  service, which must already have a WhatsApp business account and a
  registered number. Setting those up is done in Meta's business tools, not
  by the taf.
- **One approval for the sign-in, one for the tools and the Card.**

## The 24-hour rule

- **Within 24 hours of a customer's last message,** the business may reply
  with any text (`WHATSAPP_SEND_MESSAGE`).
- **Outside that window,** it may only send a **template** Meta approved
  beforehand (`WHATSAPP_SEND_TEMPLATE_MESSAGE`). A plain message is refused.

A job that writes first (a reminder, a confirmation) needs an approved
template, named on the Card.

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `WHATSAPP_GET_PHONE_NUMBERS` | Lists the business's numbers | read | **the number's id** |
| `WHATSAPP_GET_MESSAGE_TEMPLATES` | Lists approved templates | read | template names and languages |
| `WHATSAPP_SEND_TEMPLATE_MESSAGE` | Sends an approved template | write, outside | the message's id |
| `WHATSAPP_SEND_MESSAGE` | Sends text, within 24 hours of the customer's message | write, outside | the message's id |
| `WHATSAPP_SEND_MEDIA` | Sends an image or document | write, outside | |
| `WHATSAPP_SEND_INTERACTIVE_BUTTONS` | Sends a message with reply buttons | write, outside | |

**What to keep:** the sending number's id; each message's id, to match
replies and delivery.

**Every message reaches a real customer's phone.** The Card bounds who and
how many, and the first message is shown to the owner.

## Don't offer

- Marketing messages to lists (`WHATSAPP_SEND_MARKETING_MESSAGE`): they're
  charged, need opt-in, and are the fastest way to get a number flagged.
- Account, number and template administration: registering or removing
  numbers, creating or deleting templates, business solutions, groups,
  blocking users, flows.

## Access and plans

- **Meta charges per conversation started by a template**, by country and
  kind (utility, marketing, authentication). Replies within 24 hours to a
  customer's message are free.
- **New numbers start with a daily cap** on how many customers they can
  message first (from 250 a day), rising with good use.
- **Quality:** if many people block or report the number, Meta lowers its
  cap or suspends it.
- **Receiving messages** needs a webhook that this app doesn't set up; a taf
  here sends, it doesn't read customers' replies.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused: outside 24 hours (131047) | A plain message outside the window | "That customer last wrote over a day ago, so only an approved template can be sent." |
| Template refused (132000 and others) | Not approved, wrong language, or wrong fields | "WhatsApp refused the template. Check it's approved in that language." |
| Daily cap reached | The number's first-contact cap | "Your WhatsApp number's daily limit is reached. The rest go tomorrow." |
| Not a WhatsApp user | The number isn't on WhatsApp | "That number isn't on WhatsApp, so the message wasn't sent." |
| Payment refused | No valid payment method in Meta's business account | "Meta needs a payment method on your WhatsApp account before it sends this." |
| The sign-in no longer works (401) | Removed, or expired | "WhatsApp needs you to sign in again. Nothing was sent." |

After any refusal the taf stops that step and says so.
