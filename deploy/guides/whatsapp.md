# WhatsApp

Checked: 2026-11-14 · Sources: the connection service's WhatsApp tool
definitions as planque reads them (127 tools, [generated](generated/whatsapp.md)),
and WhatsApp's Business Platform documentation.

## What it is

**WhatsApp's business platform, not the owner's personal WhatsApp.** A taf
sends from a **business phone number** on the platform, never from the
owner's own phone or chats.

## How it connects

- **A sign-in to the owner's business account,** held by the connection
  service, which must already have a WhatsApp business account and a
  registered number. Setting those up is done in WhatsApp's own business
  tools, not by the taf.
- **One approval for the sign-in, one for the tools and the Card.**

## The 24-hour rule

- **Within 24 hours of a customer's last message or call,** the business may
  reply freely (`WHATSAPP_SEND_MESSAGE`). Each new message from them starts
  the 24 hours again.
- **Outside that window,** it may only send a **template** approved
  beforehand (`WHATSAPP_SEND_TEMPLATE_MESSAGE`). A plain message is refused.

A job that writes first (a reminder, a confirmation) needs an approved
template, named on the Card.

## Tools a job usually needs

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `WHATSAPP_GET_PHONE_NUMBERS` | Lists the business's numbers | read | not checked here; keep **the number's id** |
| `WHATSAPP_GET_MESSAGE_TEMPLATES` | Lists approved templates | read | not checked here |
| `WHATSAPP_SEND_TEMPLATE_MESSAGE` | Sends an approved template | write, outside | the message's id |
| `WHATSAPP_SEND_MESSAGE` | Sends text, within 24 hours of the customer's message | write, outside | the message's id |
| `WHATSAPP_SEND_MEDIA` | Sends an image or document | write, outside | the message's id |
| `WHATSAPP_SEND_INTERACTIVE_BUTTONS` | Sends a message with reply buttons | write, outside | the message's id |

**What to keep:** the sending number's id; each message's id, which delivery
and read reports name.

**Every message reaches a real customer's phone.** The Card bounds who and
how many, and the first message is shown to the owner.

## Don't offer

- Marketing messages to lists (`WHATSAPP_SEND_MARKETING_MESSAGE`): each one
  is charged, and WhatsApp's policy requires the customer's opt-in before a
  business writes to them first.
- Messaging anyone who hasn't opted in to hear from the business.
- Account, number and template administration: registering or removing
  numbers, creating or deleting templates, business solutions, groups,
  blocking users, flows.

## Access and plans

- **Charged per delivered message,** by country and kind (marketing,
  utility, authentication, service). Customers' messages to the business
  cost nothing, and each number's first 1,000 service messages a month are
  free.
- **A cap on first contact:** a business can message at most so many
  customers outside the 24 hours in any rolling day, shared by all its
  numbers. A new business starts at 250; the cap can rise to 2,000, 10,000,
  100,000 or no cap.
- **Replies and delivery reports** arrive by webhook, which this app doesn't
  set up; a taf here sends, it doesn't read customers' replies.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused: outside 24 hours (131047) | A plain message outside the window | "That customer last wrote over a day ago, so only an approved template can be sent." |
| Template refused (132000 and others) | The values sent don't match the template's fields, or the template can't be used | "WhatsApp refused the template message. The taf has stopped and will show you what it tried." |
| First-contact cap reached | The business's cap on customers messaged first in a day | "Your WhatsApp number's daily limit is reached. The rest go tomorrow." |
| Can't be delivered (131026) | The number isn't on WhatsApp, or its app is out of date | "That number can't get WhatsApp messages right now, so nothing was delivered." |
| Payment refused (131042) | A problem with the business account's payment method | "WhatsApp needs a working payment method on your business account before it sends this." |
| Sending too fast (130429) | Over the platform's sending rate | "WhatsApp asked us to slow down. The rest goes in a moment." |
| The sign-in no longer works (code 0) | The access expired or was withdrawn | "WhatsApp needs you to sign in again. Nothing was sent." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
