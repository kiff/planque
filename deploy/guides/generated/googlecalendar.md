# Google Calendar: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

47 tools: 19 read, 28 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `GOOGLECALENDAR_ACL_DELETE` | Delete ACL Rule | write | `calendar_id`, `rule_id` |  | `calendar`, `calendar.acls` |
| `GOOGLECALENDAR_ACL_GET` | Get ACL Rule | read | `calendar_id`, `rule_id` |  | `calendar`, `calendar.readonly`, `calendar.acls`, `calendar.acls.readonly` |
| `GOOGLECALENDAR_ACL_INSERT` | Create ACL Rule | write | `calendar_id`, `role`, `scope`; *send_notifications* |  | `calendar`, `calendar.acls` |
| `GOOGLECALENDAR_ACL_LIST` | List ACL Rules | read | `calendarId`; *maxResults, pageToken, showDeleted, syncToken* |  | `calendar.acls.readonly`, `calendar.acls`, `calendar` |
| `GOOGLECALENDAR_ACL_PATCH` | Patch ACL Rule | write | `calendar_id`, `rule_id`; *role, scope, send_notifications* |  | `calendar`, `calendar.acls` |
| `GOOGLECALENDAR_ACL_UPDATE` | Update ACL Rule | write | `calendar_id`, `role`, `rule_id`; *send_notifications* |  | `calendar`, `calendar.acls` |
| `GOOGLECALENDAR_ACL_WATCH` | Watch ACL Changes | write | `address`, `calendarId`, `id`; *params, token, type* |  | `calendar`, `calendar.acls`, `calendar.acls.readonly` |
| `GOOGLECALENDAR_BATCH_EVENTS` | Batch Events | write | `operations`; *fail_fast* |  | `calendar`, `calendar.app.created`, `calendar.events`, `calendar.events.owned` |
| `GOOGLECALENDAR_CALENDARS_DELETE` | Delete Calendar | write | `calendar_id` |  | `calendar`, `calendar.app.created`, `calendar.calendars` |
| `GOOGLECALENDAR_CALENDARS_UPDATE` | Update Calendar | write | `calendarId`, `summary`; *description, location, timeZone* |  | `calendar`, `calendar.app.created`, `calendar.calendars` |
| `GOOGLECALENDAR_CALENDAR_LIST_DELETE` | Remove Calendar from List | write | `calendar_id` |  | `calendar`, `calendar.app.created`, `calendar.calendarlist` |
| `GOOGLECALENDAR_CALENDAR_LIST_GET` | Get Single Calendar by ID | read | `calendarId` |  | `calendar`, `calendar.readonly`, `calendar.app.created`, `calendar.calendarlist`, `calendar.calendarlist.readonly` |
| `GOOGLECALENDAR_CALENDAR_LIST_INSERT` | Insert Calendar into List | write | `id`; *backgroundColor, colorId, color_rgb_format, defaultReminders, foregroundColor, hidden, notificationSettings, selected and 1 more* |  | `calendar`, `calendar.calendarlist` |
| `GOOGLECALENDAR_CALENDAR_LIST_PATCH` | Patch Calendar List Entry | write | `calendar_id`; *backgroundColor, colorId, colorRgbFormat, defaultReminders, foregroundColor, hidden, notificationSettings, selected and 1 more* |  | `calendar`, `calendar.app.created`, `calendar.calendarlist` |
| `GOOGLECALENDAR_CALENDAR_LIST_UPDATE` | Update Calendar List Entry | write | `calendar_id`; *backgroundColor, colorId, colorRgbFormat, defaultReminders, foregroundColor, hidden, notificationSettings, selected and 1 more* |  | `calendar`, `calendar.app.created`, `calendar.calendarlist` |
| `GOOGLECALENDAR_CALENDAR_LIST_WATCH` | Watch Calendar List | write | `address`, `id`, `type`; *expiration, params, token* |  | `calendar`, `calendar.readonly`, `calendar.calendarlist`, `calendar.calendarlist.readonly` |
| `GOOGLECALENDAR_CHANNELS_STOP` | Stop Channel | write | `id`, `resourceId`; *token* |  | `calendar`, `calendar.acls`, `calendar.acls.readonly`, `calendar.app.created`, `calendar.calendarlist`, `calendar.calendarlist.readonly`, `calendar.events`, `calendar.events.freebusy`, `calendar.events.owned`, `calendar.events.owned.readonly`, `calendar.events.public.readonly`, `calendar.events.readonly`, `calendar.readonly`, `calendar.settings.readonly` |
| `GOOGLECALENDAR_CLEAR_CALENDAR` | Clear Calendar | write | `calendar_id` |  | `calendar`, `calendar.calendars` |
| `GOOGLECALENDAR_COLORS_GET` | Get Color Definitions | read |  |  | `calendar`, `calendar.readonly`, `calendar.app.created`, `calendar.calendarlist`, `calendar.calendarlist.readonly`, `calendar.events.freebusy`, `calendar.events.owned`, `calendar.events.owned.readonly`, `calendar.events.public.readonly` |
| `GOOGLECALENDAR_CREATE_CALENDAR` | Create a new calendar | write | `summary`; *description, location, timezone* |  | `calendar`, `calendar.app.created`, `calendar.calendars` |
| `GOOGLECALENDAR_CREATE_EVENT` | Create Event | write | *attendees, birthdayProperties, calendar_id, create_meeting_room, description, end_date, end_datetime, eventType and 19 more* |  | `calendar`, `calendar.app.created`, `calendar.events`, `calendar.events.owned`, `calendar.calendars`, `calendar.calendars.readonly`, `calendar.readonly` |
| `GOOGLECALENDAR_DELETE_EVENT` | Delete event | write | `event_id`; *calendar_id, send_notifications, send_updates* |  | `calendar`, `calendar.app.created`, `calendar.events`, `calendar.events.owned` |
| `GOOGLECALENDAR_DUPLICATE_CALENDAR` | Create a calendar | write | `summary` |  | `calendar`, `calendar.app.created`, `calendar.calendars` |
| `GOOGLECALENDAR_EVENTS_GET` | Get Event | read | `event_id`; *calendar_id, max_attendees, time_zone* |  | `calendar`, `calendar.readonly`, `calendar.events`, `calendar.events.readonly`, `calendar.app.created`, `calendar.events.freebusy`, `calendar.events.owned`, `calendar.events.owned.readonly`, `calendar.events.public.readonly` |
| `GOOGLECALENDAR_EVENTS_IMPORT` | Import Event | write | `end`, `iCalUID`, `start`; *attachments, attendees, calendar_id, colorId, conferenceDataVersion, description, extendedProperties, guestsCanInviteOthers and 12 more* |  | `calendar`, `calendar.events`, `calendar.app.created`, `calendar.events.owned` |
| `GOOGLECALENDAR_EVENTS_INSTANCES` | Get Event Instances | read | `calendarId`, `eventId`; *maxAttendees, maxResults, originalStart, pageToken, showDeleted, timeMax, timeMin, timeZone* |  | `calendar.events.public.readonly`, `calendar.events.owned.readonly`, `calendar.events.readonly`, `calendar.readonly`, `calendar.events.freebusy`, `calendar.events.owned`, `calendar.app.created`, `calendar.events`, `calendar` |
| `GOOGLECALENDAR_EVENTS_LIST` | List Events | read | *alwaysIncludeEmail, calendarId, eventTypes, fields, iCalUID, maxAttendees, maxResults, orderBy and 12 more* |  | `calendar`, `calendar.app.created`, `calendar.events`, `calendar.events.freebusy`, `calendar.events.owned`, `calendar.events.owned.readonly`, `calendar.events.public.readonly`, `calendar.events.readonly`, `calendar.readonly` |
| `GOOGLECALENDAR_EVENTS_LIST_ALL_CALENDARS` | List Events from All Calendars | read | `time_max`, `time_min`; *calendar_ids, event_types, max_results_per_calendar, q, response_detail, show_deleted, single_events* |  | `calendar`, `calendar.readonly`, `calendar.calendarlist`, `calendar.calendarlist.readonly`, `calendar.events`, `calendar.events.readonly`, `calendar.events.owned`, `calendar.events.owned.readonly`, `calendar.events.public.readonly`, `calendar.events.freebusy`, `calendar.app.created` |
| `GOOGLECALENDAR_EVENTS_MOVE` | Move Event | write | `calendar_id`, `destination`, `event_id`; *send_updates* |  | `calendar.events.owned`, `calendar.events`, `calendar` |
| `GOOGLECALENDAR_EVENTS_WATCH` | Watch Events | write | `address`, `calendarId`, `id`; *params, payload, token, type* |  | `calendar`, `calendar.app.created`, `calendar.events`, `calendar.events.freebusy`, `calendar.events.owned`, `calendar.events.owned.readonly`, `calendar.events.public.readonly`, `calendar.events.readonly`, `calendar.readonly` |
| `GOOGLECALENDAR_FIND_EVENT` | Find event | read | *calendar_id, event_types, max_results, order_by, page_token, query, show_deleted, single_events and 3 more* |  | `calendar`, `calendar.readonly`, `calendar.events`, `calendar.events.readonly`, `calendar.events.owned`, `calendar.events.owned.readonly`, `calendar.events.public.readonly`, `calendar.events.freebusy`, `calendar.app.created` |
| `GOOGLECALENDAR_FIND_FREE_SLOTS` | Find free slots | read | *calendar_expansion_max, group_expansion_max, items, time_max, time_min, timezone* |  | `calendar`, `calendar.readonly`, `calendar.freebusy`, `calendar.events.freebusy` |
| `GOOGLECALENDAR_FREE_BUSY_QUERY` | Query Calendar Availability | read | `items`, `timeMax`, `timeMin`; *calendarExpansionMax, groupExpansionMax, timeZone* |  | `calendar`, `calendar.readonly`, `calendar.freebusy`, `calendar.events.freebusy` |
| `GOOGLECALENDAR_GET_CALENDAR` | Get Google Calendar | read | *calendar_id* |  | `calendar`, `calendar.readonly`, `calendar.calendars`, `calendar.calendars.readonly`, `calendar.app.created` |
| `GOOGLECALENDAR_GET_CURRENT_DATE_TIME` | Get current date and time | read | *timezone* |  |  |
| `GOOGLECALENDAR_GET_CURRENT_USER` | Get Current Google Calendar User | read |  |  | `calendar`, `calendar.readonly`, `calendar.calendars`, `calendar.calendars.readonly` |
| `GOOGLECALENDAR_LIST_BUILDINGS` | List Buildings | read | *customer, maxResults, pageToken* |  | `https://apps-apis.google.com/a/feeds/calendar/resource/`, `admin.directory.resource.calendar`, `admin.directory.resource.calendar.readonly` |
| `GOOGLECALENDAR_LIST_CALENDARS` | List Google Calendars | read | *max_results, min_access_role, page_token, show_deleted, show_hidden, sync_token* |  | `calendar`, `calendar.readonly`, `calendar.calendarlist`, `calendar.calendarlist.readonly` |
| `GOOGLECALENDAR_LIST_CALENDAR_RESOURCES` | List Calendar Resources | read | *customer, maxResults, orderBy, pageToken, query* |  | `https://apps-apis.google.com/a/feeds/calendar/resource/`, `admin.directory.resource.calendar`, `admin.directory.resource.calendar.readonly` |
| `GOOGLECALENDAR_PATCH_CALENDAR` | Patch Calendar | write | `calendar_id`; *description, location, summary, timezone* |  | `calendar`, `calendar.app.created`, `calendar.calendars` |
| `GOOGLECALENDAR_PATCH_EVENT` | Patch Event | write | `calendar_id`, `event_id`; *anyone_can_add_self, attachments, attendees, attendees_omitted, color_id, conference_data, conference_data_version, description and 23 more* |  | `calendar`, `calendar.app.created`, `calendar.events`, `calendar.events.owned` |
| `GOOGLECALENDAR_QUICK_ADD` | Quick Add Event | write | `text`; *calendar_id, send_updates* |  | `calendar`, `calendar.events`, `calendar.events.owned`, `calendar.app.created` |
| `GOOGLECALENDAR_REMOVE_ATTENDEE` | Remove attendee from event | write | `attendee_email`, `event_id`; *calendar_id* |  | `calendar`, `calendar.app.created`, `calendar.events`, `calendar.events.owned` |
| `GOOGLECALENDAR_SETTINGS_GET` | Get Calendar Setting | read | `setting` |  | `calendar.settings.readonly`, `calendar.readonly`, `calendar` |
| `GOOGLECALENDAR_SETTINGS_LIST` | List Settings | read | *maxResults, pageToken, syncToken* |  | `calendar.settings.readonly`, `calendar.readonly`, `calendar` |
| `GOOGLECALENDAR_SETTINGS_WATCH` | Watch Settings | write | `address`, `id`, `type`; *expiration, params, token* |  | `calendar.settings.readonly`, `calendar.readonly`, `calendar` |
| `GOOGLECALENDAR_UPDATE_EVENT` | Update Google event | write | `event_id`, `start_datetime`; *attendees, birthdayProperties, calendar_id, create_meeting_room, description, end_datetime, eventType, event_duration_hour and 15 more* |  | `calendar`, `calendar.events`, `calendar.events.owned`, `calendar.app.created` |
