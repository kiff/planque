# Shopify

Checked: 2026-11-14 · Sources: the connection service's Shopify tool
definitions as planque reads them (386 tools), and Shopify's Admin API
documentation (access scopes, rate limits, status codes).
Failure table checked against planque's answers on 2026-11-16.
Sign-in checked against planque's integrations record on 2026-11-16.
Order access checked on a run through the owner's phone, 2026-11-20.

## How it connects

- **A sign-in to the owner's Shopify store**, held by the connection
  service; the taf never sees a password or token. The owner approves the
  access on Shopify's own install page. The grant's catalog lists the tools
  and, with a known sign-in, puts the ones its access doesn't cover under
  `unusable`, with why.
- **Type only the store's subdomain** (`my-shop` for `my-shop.myshopify.com`)
  where the connection page asks for it; its auto-fill can be wrong.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request.

## Tools a job usually needs

A tool marked read fetches from the store; one marked write creates, changes
or sends. The ids below are what Shopify's answers carry.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `SHOPIFY_GET_SHOP_INFO` | Reads the store's name, currency and time zone | read | the store's id |
| `SHOPIFY_GET_ORDERS_WITH_FILTERS` | Lists orders by date, status, payment or fulfilment state | read | order ids, a page link |
| `SHOPIFY_GET_ORDER` | Reads one order: items, customer, totals, status | read | the order id, its customer id |
| `SHOPIFY_GET_CUSTOMERS_SEARCH` | Finds customers by email, name or other fields | read | customer ids |
| `SHOPIFY_GET_CUSTOMER_ORDERS` | Lists one customer's orders | read | order ids |
| `SHOPIFY_GET_ORDER_FULFILLMENTS` | Reads an order's shipments and tracking | read | fulfilment ids, tracking numbers |
| `SHOPIFY_GET_ORDER_REFUNDS` | Lists the refunds already made on an order | read | refund ids |
| `SHOPIFY_CALCULATE_REFUND` | Works out what a refund would be, changing nothing | read | the amounts and lines it would refund |
| `SHOPIFY_SEARCH_PRODUCTS` | Finds products by title, type, tag or status | read | product ids, a page cursor |
| `SHOPIFY_GET_INVENTORY_LEVELS` | Reads stock per item and location | read | inventory item ids, location ids |
| `SHOPIFY_UPDATE_ORDER` | Changes an order's note, tags or contact details | write | the order id |
| `SHOPIFY_CREATE_DRAFT_ORDER` | Drafts an order for the owner to review; charges nothing | write | the draft order id |
| `SHOPIFY_UPDATE_FULFILLMENT_TRACKING` | Adds or corrects a shipment's tracking | write | the fulfilment id |
| `SHOPIFY_CREATE_REFUND` | Refunds money on an order | write, money | the refund id |

**What the sign-in asks for:** every scope the managed sign-in offers (34).
The Card decides what a taf may do: offer the tools the job needs, and
write tools when the owner asks for them.

**What to keep:** the order and customer ids the taf worked on, and each
draft order and refund id it creates. Never search for something the taf
made: use the id it kept.

**Check before a refund.** Run `SHOPIFY_CALCULATE_REFUND` first and show the
owner the amount. A refund moves real money: it waits for their yes.

## Don't offer

- Deleting or cancelling (`SHOPIFY_DELETE_*`, `SHOPIFY_CANCEL_ORDER`,
  `SHOPIFY_CLOSE_ORDER`, `SHOPIFY_CANCEL_FULFILLMENT`): they can't be undone.
- Charging or money-like tools (`SHOPIFY_CREATE_ORDER_TRANSACTION`,
  `SHOPIFY_CREATE_GIFT_CARD`, discount codes and price rules) unless the job
  is exactly that.
- Turning a draft into a real order (`SHOPIFY_COMPLETE_DRAFT_ORDER`) and
  emailing customers: the draft is for the owner to review and send.
- App billing and access (app subscriptions and charges,
  `SHOPIFY_REVOKE_APP_ACCESS_SCOPES`, `SHOPIFY_CREATE_DELEGATE_ACCESS_TOKEN`).
- Bulk and open-ended writes (`SHOPIFY_RUN_BULK_MUTATION_OPERATION`,
  `SHOPIFY_BULK_*`, `SHOPIFY_UPDATE_PRODUCT_VARIANTS_BULK`,
  `SHOPIFY_GRAPH_QL_*` write tools): one call changes many things at once.
- The storefront and its plumbing (themes, script tags, webhooks, storefront
  tokens, fulfilment services): they change what every visitor sees or where
  the data goes.
- Stock changes (`SHOPIFY_SET_INVENTORY_LEVEL`,
  `SHOPIFY_ADJUST_INVENTORY_LEVEL`) unless the owner asked: a wrong number
  oversells or hides a product.

## Access and plans

- **Orders older than 60 days** need a separate "all orders" access.
- **Orders and customers' names, emails and addresses** are protected data.
  A public app needs Shopify's approval on a live store. The owner's own app
  on a development store needs no review, but its "protected customer data"
  form (reasons and data-protection answers) must be saved **and
  released as a new app version** (the "Release" button; "Publier" in French). Saved alone it stays a draft:
  orders are refused while draft orders work. Guide the owner tap by tap;
  don't do it for them.
- **Shopify reports write access when it grants write:** it covers reading
  orders too.
- **Rate:** about 2 calls a second, bursts up to 40; ten times that on
  Shopify Plus. Over it, Shopify says how long to wait.
- **Paging:** lists come a page at a time; 250 items at most in one input.

## When it fails

The rows are Shopify's own answers: a refusal reaches the taf as "Shopify
refused it (it answered …): …", with Shopify's reason. planque's own
refusals, and an expired sign-in, are in the part below.

| What the taf gets | Why | Tell the owner |
|---|---|---|
| Refused (403), not approved to access orders | The protected-data form isn't released yet | "Shopify hasn't approved this app to read orders yet: release its settings as a new app version, then ask again. Draft orders still work." |
| Refused for missing access (403) | The sign-in didn't give this tool's access | "Shopify is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| A result missing older orders | No "all orders" access, so only the last 60 days are visible | "Shopify only showed the last 60 days of orders. To look further back, sign in again and allow access to all orders." |
| Refused (429), too many calls | Over the store's call rate | "Shopify asked us to slow down. The taf will try again in a moment." |
| Refused (402), payment required | The store is frozen for an unpaid Shopify bill: every call fails; not a sign-in problem | "Your Shopify store is paused until its bill is paid. Nothing was changed." |
| Refused as invalid (422) | A required field is missing, or the action doesn't fit the order | "Shopify didn't accept that change to the order. The taf has stopped and will show you what it tried." |
| Refused (404), not found | The order, customer or product was deleted | "That isn't in your store any more. The taf will skip it." |
| Shopify's sign-in has expired (401) | The app was removed from the store, or access expired | "Shopify needs you to connect the store again. Nothing was changed." |
| Refused (500 or above) | Trouble on Shopify's side | "Shopify is having trouble right now. The taf will try again later." |

A refund is never retried without the owner's yes.
