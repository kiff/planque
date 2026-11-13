# Google Maps: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

24 tools: 22 read, 2 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `GOOGLE_MAPS_AUTOCOMPLETE` | Autocomplete Place Predictions | read | `input`; *includePureServiceAreaBusinesses, includeQueryPredictions, includedPrimaryTypes, includedRegionCodes, inputOffset, languageCode, locationBias, locationRestriction and 3 more* |  | `cloud-platform` |
| `GOOGLE_MAPS_COMPUTE_ROUTE_MATRIX` | Compute Route Matrix | read | `destinations`, `origins`; *fieldMask, languageCode, routingPreference, travelMode, units* |  | `cloud-platform` |
| `GOOGLE_MAPS_FIND_STREET_VIEW_PANORAMAS` | Find Street View Panoramas | read | `locations`, `session`; *radius* |  | `cloud-platform` |
| `GOOGLE_MAPS_GEOCODE_ADDRESS_WITH_QUERY` | Geocode Address With Query | read | `address_query`; *language_code, region_code* |  | `cloud-platform`, `maps-platform.geocode`, `maps-platform.geocode.address` |
| `GOOGLE_MAPS_GEOCODE_DESTINATIONS` | Geocode Destinations | read | *addressQuery, languageCode, locationQuery, place, regionCode, travelModes* |  | `maps-platform.geocode`, `cloud-platform` |
| `GOOGLE_MAPS_GEOCODE_LOCATION` | Reverse Geocode Location | read | `latitude`, `longitude`; *languageCode, regionCode, types* |  | `cloud-platform`, `maps-platform.geocode`, `maps-platform.geocode.location` |
| `GOOGLE_MAPS_GEOCODE_PLACE` | Geocode Place by ID | read | `place_id`; *languageCode, regionCode* |  | `cloud-platform`, `maps-platform.geocode`, `maps-platform.geocode.place` |
| `GOOGLE_MAPS_GEOCODING_API` | Geocoding API | read | *address, bounds, components, extra_computations, key, language, latlng, location_type and 3 more* |  | `cloud-platform` |
| `GOOGLE_MAPS_GEOLOCATE` | Geolocate Device | read | *carrier, cellTowers, considerIp, homeMobileCountryCode, homeMobileNetworkCode, radioType, wifiAccessPoints* |  | `cloud-platform` |
| `GOOGLE_MAPS_GET2D_TILE` | Get 2D Map Tile | read | `session`, `x`, `y`, `z`; *key, orientation* |  | `cloud-platform` |
| `GOOGLE_MAPS_GET3D_TILES_ROOT` | Get 3D Tiles Root | read | *key* |  | `cloud-platform` |
| `GOOGLE_MAPS_GET_CURRENT_WEATHER` | Get Current Weather | read | `latitude`, `longitude`; *language_code, units_system* |  | `cloud-platform` |
| `GOOGLE_MAPS_GET_PLACE_DETAILS` | Get Place Details | read | `name`; *fieldMask, languageCode, regionCode, sessionToken* |  | `cloud-platform` |
| `GOOGLE_MAPS_GET_ROUTE` | Get Route | read | `destination_address`, `origin_address`; *avoid_ferries, avoid_highways, avoid_tolls, compute_alternative_routes, field_mask, language_code, routing_preference, travel_mode and 1 more* |  | `cloud-platform` |
| `GOOGLE_MAPS_GET_STREET_VIEW_METADATA` | Get Street View Metadata | read | `session`; *latitude, longitude, pano_id, radius* |  | `cloud-platform` |
| `GOOGLE_MAPS_GET_STREET_VIEW_TILE` | Get Street View Tile | read | `pano_id`, `session`, `x`, `y`, `zoom` |  | `cloud-platform` |
| `GOOGLE_MAPS_GET_TIME_ZONE` | Get Time Zone | read | `location`, `timestamp`; *language* |  |  |
| `GOOGLE_MAPS_LOOKUP_AERIAL_VIDEO` | Lookup Aerial Video | read | *address, videoId* |  | `cloud-platform` |
| `GOOGLE_MAPS_MAPS_EMBED_API` | Embed Google Map | read | `mode`; *api_key, directions_params, place_params, search_params, streetview_params, view_params* |  |  |
| `GOOGLE_MAPS_NEARBY_SEARCH` | Nearby search | read | `latitude`, `longitude`, `radius`; *excludedTypes, fieldMask, includedTypes, maxResultCount* |  | `cloud-platform` |
| `GOOGLE_MAPS_PLACE_PHOTO` | Get Place Photo | read | `photo_reference`; *maxheight, maxwidth* |  |  |
| `GOOGLE_MAPS_RENDER_AERIAL_VIDEO` | Render Aerial Video | write | `address` |  | `cloud-platform` |
| `GOOGLE_MAPS_TEXT_SEARCH` | Text Search | read | `textQuery`; *fieldMask, maxResultCount* |  | `cloud-platform` |
| `GOOGLE_MAPS_TILES_CREATE_SESSION` | Create Tiles Session | write | `language`, `mapType`, `region`; *highDpi, imageFormat, layerTypes, overlay, scale, styles* |  | `cloud-platform` |
