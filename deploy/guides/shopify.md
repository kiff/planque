# Shopify

Checked: 2026-11-14 · Sources: the connection service's Shopify tool
definitions as planque reads them (386 tools), and Shopify's Admin API
documentation (access scopes, rate limits, status codes).

## How it connects

- **A sign-in to the owner's Shopify store**, held by the connection
  service; the taf never sees a password or token. The owner approves the
  access on Shopify's own install page. planque's catalog lists the tools
  and, with a known sign-in, puts the ones its access doesn't cover under
  `unusable`, with why.
- **A held sign-in is reused** for a new taf when it covers the tools asked.
  Ask for every tool the job needs in the first request, so nothing is asked
  twice.

## Tools a job usually needs

Read means it changes nothing in the store. The catalog's `returns` lists
no paths for these tools yet; the ids below are what Shopify's answers
carry.

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

**What to keep:** the order id and customer id the taf worked on, to come
back to the same order next run; each draft order and refund id the taf
creates. Never search for something the taf made: use the id it kept.

**Check before a refund.** Run `SHOPIFY_CALCULATE_REFUND` first and show the
owner the amount. A refund moves real money, so it waits for their yes and
stays within the Card's limit.

## Don't offer

- Deleting or cancelling (`SHOPIFY_DELETE_*`, `SHOPIFY_CANCEL_ORDER`,
  `SHOPIFY_CLOSE_ORDER`, `SHOPIFY_CANCEL_FULFILLMENT`): they can't be undone, and a cancel can
  refund and restock in one go.
- Charging or money-like tools (`SHOPIFY_CREATE_ORDER_TRANSACTION`,
  `SHOPIFY_CREATE_GIFT_CARD`, `SHOPIFY_CREATE_DISCOUNT_CODE`,
  `SHOPIFY_CREATE_PRICE_RULE`, `SHOPIFY_CREATE_DISCOUNT_CODE_BATCH`) unless
  the job is exactly that.
- Turning a draft into a real order (`SHOPIFY_COMPLETE_DRAFT_ORDER`) and
  emailing customers (`SHOPIFY_SEND_INVOICE`, `SHOPIFY_SEND_CUSTOMER_INVITE`):
  the draft is there for the owner to review and send.
- App billing and access (`SHOPIFY_CREATE_APP_SUBSCRIPTION`,
  `SHOPIFY_CREATE_*_APPLICATION_CHARGE`, `SHOPIFY_CREATE_USAGE_CHARGE`,
  `SHOPIFY_CANCEL_APP_SUBSCRIPTION`, `SHOPIFY_REVOKE_APP_ACCESS_SCOPES`,
  `SHOPIFY_CREATE_DELEGATE_ACCESS_TOKEN`).
- Bulk and open-ended writes (`SHOPIFY_RUN_BULK_MUTATION_OPERATION`,
  `SHOPIFY_BULK_*`, `SHOPIFY_UPDATE_PRODUCT_VARIANTS_BULK`,
  `SHOPIFY_GRAPH_QL_*` write tools): one call changes many things at once.
- The storefront and its plumbing: themes and their assets, script tags,
  webhook subscriptions, storefront tokens, fulfilment services. They change
  what every visitor sees or where the store's data goes.
- Stock changes (`SHOPIFY_SET_INVENTORY_LEVEL`,
  `SHOPIFY_ADJUST_INVENTORY_LEVEL`) unless the owner asked for them: a wrong
  number oversells or hides a product.

## Access and plans

- **Orders older than 60 days** need a separate "all orders" access on top
  of order access. Without it, older orders simply aren't returned, so a
  job looking back further misses them.
- **Customer names, emails and addresses** are protected data on Shopify.
  An app must be approved by Shopify to read them on a live store. If
  the app isn't approved, Shopify doesn't return those fields.
- **Rate:** about 2 calls a second sustained, with short bursts up to 40,
  for each app on a store. On Shopify Plus stores it's ten times that. Over the limit,
  Shopify says how long to wait.
- **Paging:** large lists come a page at a time; 250 items is the most
  in one input list.
- **A frozen store** (an unpaid Shopify bill) refuses every call until the
  owner pays it.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403) | The sign-in didn't give this tool's access | "Shopify is only partly connected: this step needs access you didn't give. Want to sign in again and allow it, or do without this step?" |
| Older orders are missing | No "all orders" access, so only the last 60 days are visible | "Shopify only showed the last 60 days of orders. To look further back, sign in again and allow access to all orders." |
| Too many calls (429) | Over the store's call rate | "Shopify asked us to slow down. The taf will try again in a moment." |
| Payment required (402) | The store is frozen for an unpaid Shopify bill | "Your Shopify store is paused until its bill is paid. Nothing was changed." |
| Refused as invalid (422) | A required field is missing, or the action doesn't fit the order (refunding more than was paid, say) | "Shopify didn't accept that change to the order. The taf has stopped and will show you what it tried." |
| Not found (404) | The order, customer or product was deleted | "That isn't in your store any more. The taf will skip it." |
| The sign-in no longer works (401) | The app was removed from the store, or access expired | "Shopify needs you to connect the store again. Nothing was changed." |
| A server error (5xx) | Trouble on Shopify's side | "Shopify is having trouble right now. The taf will try again later." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write, and a refund is never retried without the owner's yes.
