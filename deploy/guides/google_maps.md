# Google Maps

Checked: 2026-11-13 · Sources: the connection service's Google Maps tool
definitions as planque reads them (24 tools, [generated](generated/google_maps.md)),
and Google's Maps Platform reference.

## How it connects

- **A sign-in to Google**, held by the connection service. The catalog lists
  the tools; all but two only read public map data. Nothing of the owner's is read: no saved
  places, no location history.
- **One approval for the sign-in, one for the tools and the Card.**

## Tools a job usually needs

| Tool | What it does | Read or write | Keep |
|---|---|---|---|
| `GOOGLE_MAPS_TEXT_SEARCH` | Finds places from words ("bakeries in Lyon") | read | **each place's id** |
| `GOOGLE_MAPS_NEARBY_SEARCH` | Finds places of a type near a point | read | place ids |
| `GOOGLE_MAPS_GET_PLACE_DETAILS` | Reads a place's address, hours, phone, site | read | |
| `GOOGLE_MAPS_GEOCODE_ADDRESS_WITH_QUERY` | Turns an address into coordinates | read | the coordinates |
| `GOOGLE_MAPS_GEOCODE_LOCATION` | Turns coordinates into an address | read | |
| `GOOGLE_MAPS_GET_ROUTE` | Gives a route's distance and time, by car, foot, bike or transit | read | |
| `GOOGLE_MAPS_COMPUTE_ROUTE_MATRIX` | Gives times between several starts and ends | read | |
| `GOOGLE_MAPS_GET_TIME_ZONE` | Gives a place's time zone | read | |
| `GOOGLE_MAPS_AUTOCOMPLETE` | Suggests places as a name is typed | read | place ids |

**What to keep:** a place's id. A place's details change (opening hours, a
new phone number); its id doesn't. Read the details again when they matter.

**Ask for fewer fields.** A search can say which fields it wants. Asking only
for the name, address and id costs less than asking for everything, and
answers faster.

## Don't offer

- Map tiles, Street View and aerial video (`GOOGLE_MAPS_GET2D_TILE`,
  `GOOGLE_MAPS_GET_STREET_VIEW_TILE`, `GOOGLE_MAPS_RENDER_AERIAL_VIDEO`):
  they're for drawing maps in an app, not for a job.
- `GOOGLE_MAPS_GEOLOCATE`: it guesses where a device is from its network,
  which a taf has no reason to know.
- Lists of places scraped for contact details in bulk: Google's terms forbid
  storing place data beyond a short time, other than the place's id.

## Access and plans

- **Maps calls are paid for per call,** with a free allowance each month.
  A job that searches hundreds of places a day can exhaust it, and calls
  are then refused until the next month.
- **Results are limited:** a search answers at most 20 places, with a few
  more pages for some searches.
- **Keeping place data:** only a place's id may be kept long term; addresses
  and hours are read again when needed.

## When it fails

| What happens | Why | Tell the owner |
|---|---|---|
| No results (`ZERO_RESULTS`) | Nothing matched | "I didn't find any place matching that. Want to try other words or a wider area?" |
| Refused, over the allowance (`OVER_QUERY_LIMIT`, 429) | The month's or minute's allowance is used | "Google Maps' allowance is used up for now. The taf will try again later." |
| Refused (`REQUEST_DENIED`, 403) | The connection's access to Maps isn't set up | "Google Maps isn't available on this connection right now." |
| A bad request (`INVALID_REQUEST`, 400) | An address or field the API didn't accept | Fix the input and try once; tell the owner only if it fails again. |
| Not found (404) | A place id that no longer exists (the place closed) | "That place isn't on Google Maps any more. It may have closed." |

After any refusal the taf stops that step and says so.
