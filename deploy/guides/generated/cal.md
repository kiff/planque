# Cal: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

170 tools: 78 read, 92 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `CAL_ADD_ATTENDEE` | Add attendee | write | `bookingId`, `email`, `name`, `timeZone` |  |  |
| `CAL_ADD_MEMBER_TO_TEAM_USING_ORG_AND_TEAM_ID` | Add member to team | write | `orgId`, `teamId`, `userId`; *accepted, disableImpersonation, role* |  |  |
| `CAL_ADD_ORGANIZATION_ATTRIBUTE_OPTION` | Add organization attribute option | write | `attributeId`, `orgId`, `slug`, `value` |  |  |
| `CAL_CANCEL_BOOKING_VIA_UID` | Cancel booking via uid | write | `bookingUid`; *cancelSubsequentBookings, cancellationReason, seatUid* |  |  |
| `CAL_CHECK_CALENDAR_VERSION2` | Check calendar availability | read | `calendar` |  |  |
| `CAL_CHECK_GCAL_SYNCHRONIZATION_STATUS` | Check gcal synchronization status | read |  |  |  |
| `CAL_CHECK_ICS_FEED_CALENDAR_ENDPOINT` | Check ics feed calendar endpoint | read |  |  |  |
| `CAL_CHECK_STRIPE_STATUS` | Check Stripe status | read |  |  |  |
| `CAL_CHECK_TEAM_STRIPE_INTEGRATION_STATUS` | Check team stripe integration status | read | `teamId` |  |  |
| `CAL_CONFIRM_BOOKING_BY_UID` | Confirm booking by uid | write | `bookingUid` |  |  |
| `CAL_CONNECT_TO_CALENDAR` | Connect to calendar | read | `calendar` |  |  |
| `CAL_CREATE_MEMBERSHIP_FOR_ORGANIZATION` | Create membership for organization | write | `orgId`, `userId`; *accepted, disableImpersonation, role* |  |  |
| `CAL_CREATE_ORGANIZATION_ATTRIBUTES` | Create organization attributes | write | `name`, `options`, `orgId`, `slug`, `type`; *enabled* |  |  |
| `CAL_CREATE_ORGANIZATION_TEAM_FORM_WORKFLOW` | Create organization team form workflow | write | `activation`, `name`, `orgId`, `steps`, `teamId`, `trigger` |  | `*` |
| `CAL_CREATE_ORGANIZATION_WEBHOOK_BY_ORG_ID` | Create organization webhook by org ID | write | `active`, `orgId`, `subscriberUrl`, `triggers`; *payloadTemplate, secret* |  |  |
| `CAL_CREATE_OR_UPDATE_TEAM_PROFILE` | Create or update team profile | write | `name`; *appIconLogo, appLogo, autoAcceptCreator, bannerUrl, bio, brandColor, calVideoLogo, darkBrandColor and 10 more* |  |  |
| `CAL_CREATE_PERSONAL_EVENT_TYPE` | Create personal event type | write | `lengthInMinutes`, `slug`, `title`; *afterEventBuffer, allowReschedulingCancelledBookings, allowReschedulingPastBookings, beforeEventBuffer, bookerActiveBookingsLimit, bookerLayouts, bookingFields, bookingLimitsCount and 39 more* |  | `EVENT_TYPE_WRITE` |
| `CAL_CREATE_PHONE_CALL_EVENT` | Create phone call event | write | `calApiKey`, `eventTypeId`, `numberToCall`, `orgId`, `teamId`, `yourPhoneNumber`; *beginMessage, generalPrompt, guestCompany, guestEmail, guestName, schedulerName, templateType* |  |  |
| `CAL_CREATE_PHONE_CALL_FOR_EVENT_TYPE` | Create phone call for event type | write | `calApiKey`, `eventTypeId`, `numberToCall`, `teamId`, `yourPhoneNumber`; *beginMessage, enabled, generalPrompt, guestCompany, guestEmail, guestName, schedulerName, templateType* |  |  |
| `CAL_CREATE_TEAM_EVENT_TYPE` | Create team event type | write | `lengthInMinutes`, `schedulingType`, `slug`, `teamId`, `title`; *afterEventBuffer, assignAllTeamMembers, beforeEventBuffer, bookingFields, bookingLimitsCount, bookingLimitsDuration, bookingWindow, color__darkThemeHex and 23 more* |  |  |
| `CAL_CREATE_TEAM_EVENT_TYPES_WITH_CUSTOM_OPTIONS` | Create team event types with custom options | write | `hosts`, `lengthInMinutes`, `lengthInMinutesOptions`, `orgId`, `schedulingType`, `slug`, `teamId`, `title`; *afterEventBuffer, assignAllTeamMembers, beforeEventBuffer, bookingFields, bookingLimitsCount, bookingLimitsDuration, bookingWindow, color__darkThemeHex and 21 more* |  |  |
| `CAL_CREATE_TEAM_INVITE_LINK` | Create team invite link | write | `teamId` |  |  |
| `CAL_CREATE_TEAM_IN_ORGANIZATION` | Create team in organization | write | `name`, `orgId`; *appIconLogo, appLogo, autoAcceptCreator, bannerUrl, bio, brandColor, calVideoLogo, darkBrandColor and 10 more* |  |  |
| `CAL_CREATE_TEAM_MEMBERSHIP_WITH_ROLE` | Create team membership with role | write | `teamId`, `userId`; *accepted, disableImpersonation, role* |  |  |
| `CAL_CREATE_USER_AVAILABILITY_SCHEDULE` | Create user availability schedule | write | `isDefault`, `name`, `timeZone`; *availability, overrides* |  |  |
| `CAL_CREATE_USER_SCHEDULE_IN_ORGANIZATION` | Create user schedule in organization | write | `isDefault`, `name`, `orgId`, `timeZone`, `userId`; *availability, overrides* |  |  |
| `CAL_CREATE_WEBHOOK_FOR_EVENT_TYPE` | Create webhook for event type | write | `active`, `eventTypeId`, `subscriberUrl`, `triggers`; *payloadTemplate, secret* |  |  |
| `CAL_DECLINE_BOOKING_WITH_REASON` | Decline booking with reason | write | `bookingUid`; *reason* |  |  |
| `CAL_DELETE_ALL_TEAM_EVENT_TYPE_WEBHOOKS` | Delete all team event type webhooks | write | `eventTypeId`, `teamId` |  |  |
| `CAL_DELETE_CONFERENCE_APP_CONNECTION` | Delete conference app connection | write | `app` |  |  |
| `CAL_DELETE_DESTINATION_CALENDAR_BY_ID` | Delete destination calendar by id | write | `id` |  |  |
| `CAL_DELETE_EVENT_TYPE_BY_ID` | Delete event type by id | write | `eventTypeId` |  |  |
| `CAL_DELETE_EVENT_TYPE_IN_TEAM` | Delete event type in team | write | `eventTypeId`, `teamId` |  |  |
| `CAL_DELETE_MEMBERSHIP_IN_TEAM` | Delete membership in team | write | `membershipId`, `orgId`, `teamId` |  |  |
| `CAL_DELETE_ORGANIZATION_ATTRIBUTE` | Delete organization attribute | write | `attributeId`, `orgId` |  |  |
| `CAL_DELETE_ORGANIZATION_ATTRIBUTE_OPTION` | Delete organization attribute option | write | `attributeId`, `optionId`, `orgId` |  |  |
| `CAL_DELETE_ORGANIZATION_MEMBERSHIP` | Delete organization membership | write | `membershipId`, `orgId` |  |  |
| `CAL_DELETE_ORG_WEBHOOK` | Delete org webhook | write | `orgId`, `webhookId` |  |  |
| `CAL_DELETE_SCHEDULE_BY_ID` | Delete schedule by id | write | `scheduleId` |  |  |
| `CAL_DELETE_SELECTED_CALENDARS` | Delete selected calendars | write | `credentialId`, `externalId`, `integration` |  |  |
| `CAL_DELETE_SELECTED_SLOT` | Delete selected slot | write | `uid` |  |  |
| `CAL_DELETE_TEAM_BY_ID` | Delete team by id | write | `teamId` |  |  |
| `CAL_DELETE_TEAM_EVENT_TYPE_IN_ORGANIZATION` | Delete team event type in organization | write | `eventTypeId`, `orgId`, `teamId` |  | `EVENT_TYPE_WRITE` |
| `CAL_DELETE_TEAM_FROM_ORGANIZATION` | Delete team from organization | write | `orgId`, `teamId` |  |  |
| `CAL_DELETE_TEAM_MEMBERSHIPS_BY_ID` | Delete team memberships by id | write | `membershipId`, `teamId` |  |  |
| `CAL_DELETE_USER_ATTRIBUTE_OPTION` | Delete user attribute option | write | `attributeOptionId`, `orgId`, `userId` |  |  |
| `CAL_DELETE_USER_FROM_ORGANIZATION` | Delete user from organization | write | `orgId`, `userId` |  |  |
| `CAL_DELETE_USER_SCHEDULE` | Delete user schedule | write | `orgId`, `scheduleId`, `userId` |  |  |
| `CAL_DELETE_WEBHOOKS_FOR_EVENT_TYPE` | Delete webhooks for event type | write | `eventTypeId` |  |  |
| `CAL_DELETE_WEBHOOK_BY_ID` | Delete webhook by id | write | `webhookId` |  |  |
| `CAL_DELETE_WEBHOOK_FOR_EVENT_TYPE` | Delete webhook for event type | write | `eventTypeId`, `webhookId` |  |  |
| `CAL_DISCONNECT_CALENDAR_USING_CREDENTIAL_ID` | Disconnect calendar using credential id | write | `calendar`, `id` |  |  |
| `CAL_EDIT_ATTENDEE_BY_ID` | Edit attendee by ID | write | `bookingId`, `id`; *email, name, timeZone* |  |  |
| `CAL_EDIT_AVAILABILITY_BY_ID` | Edit availability by ID | write | `id`; *days, endTime, scheduleId, startTime* |  |  |
| `CAL_EDIT_BOOKING_BY_ID` | Edit booking by ID | write | `id`; *description, end, start, status, title* |  |  |
| `CAL_EDIT_EVENT_TYPE_BY_ID` | Edit event type by ID | write | `id`; *afterEventBuffer, beforeEventBuffer, currency, description, disableGuests, hidden, hideCalendarNotes, hosts and 23 more* |  | `EVENT_TYPE_WRITE`, `EVENT_TYPE_READ` |
| `CAL_EDIT_SELECTED_CALENDAR_BY_ID` | Edit selected calendar by ID | write | `external_id`, `integration`, `user_id`; *new_credential_id, new_external_id, new_integration* |  |  |
| `CAL_FETCH_ALL_BOOKINGS` | Fetch all bookings | read | *afterStart, attendeeEmail, attendeeName, beforeEnd, eventTypeId, eventTypeIds, skip, sortCreated and 6 more* |  |  |
| `CAL_FETCH_EVENT_TYPE_DETAILS` | Fetch event type details | read | `eventTypeId`, `orgId`, `teamId` |  |  |
| `CAL_FETCH_ORGANIZATION_ATTRIBUTE_BY_ID` | Fetch organization attribute by id | read | `attributeId`, `orgId` |  |  |
| `CAL_FETCH_PROVIDER_ACCESS_TOKEN` | Fetch provider access token | read | `clientId` |  |  |
| `CAL_FETCH_SCHEDULE_BY_ID` | Fetch schedule by id | read | `scheduleId` |  |  |
| `CAL_FETCH_USER_SCHEDULE_BY_ORG_ID` | Fetch user schedule by org id | read | `orgId`, `scheduleId`, `userId` |  |  |
| `CAL_FETCH_WEBHOOK_BY_EVENT_TYPE_ID` | Fetch webhook by event type id | read | `eventTypeId`, `webhookId` |  |  |
| `CAL_GET_ALL_TIMEZONES` | Get all timezones | read |  |  |  |
| `CAL_GET_AVAILABLE_SLOTS_INFO` | Get available slots info | read | `endTime`, `eventTypeId`, `startTime`; *duration, eventTypeSlug, orgSlug, rescheduleUid, slotFormat, timeZone, usernameList* |  |  |
| `CAL_GET_BOOKING_REFERENCES` | Get booking references | read | `bookingUid`, `orgId`, `teamId`; *type* |  |  |
| `CAL_GET_BOOKING_REFERENCE_BY_ID` | Get booking reference by id | read | `id` |  |  |
| `CAL_GET_CONFERENCE_OAUTH_AUTHORIZATION_URL` | Get conference OAuth authorization url | write | `app`, `onErrorReturnTo`, `returnTo` |  |  |
| `CAL_GET_DEFAULT_SCHEDULE_DETAILS` | Get default schedule details | read |  |  |  |
| `CAL_GET_DESTINATION_CALENDARS` | Get destination calendars | read |  |  |  |
| `CAL_GET_EVENT_TYPE_BY_TEAM_ID` | Get event type by team id | read | `eventTypeId`, `teamId` |  |  |
| `CAL_GET_EVENT_TYPE_PRIVATE_LINKS` | Get event type private links | read | `eventTypeId` |  | `EVENT_TYPE_READ` |
| `CAL_GET_GOOGLE_CALENDAR_OAUTH_AUTHENTICATION_URL` | Get google calendar oauth authentication url | write |  |  |  |
| `CAL_GET_OAUTH_CLIENTS_USER` | Get oauth clients user | read | `clientId`; *limit, x_cal_secret_key* |  |  |
| `CAL_GET_ORGANIZATION_ATTRIBUTE_ASSIGNED_OPTIONS` | Get organization attribute assigned options | read | `attributeId`, `orgId`; *assignedOptionIds, skip, take, teamIds* |  |  |
| `CAL_GET_ORGANIZATION_ATTRIBUTE_ASSIGNED_OPTIONS_BY_SLUG` | Get organization attribute assigned options by slug | read | `attributeSlug`, `orgId`; *assignedOptionIds, skip, take, teamIds* |  | `organization.attributes.read` |
| `CAL_GET_ORGANIZATION_ID` | Get organization ID | read |  |  |  |
| `CAL_GET_ORGANIZATION_SCHEDULES` | Get organization schedules | read | `orgId`; *skip, take* |  |  |
| `CAL_GET_ORGANIZATION_TEAMS_EVENT_TYPES` | Get organization teams event types | read | `orgId`; *skip, sortCreatedAt, take* |  |  |
| `CAL_GET_ORGANIZATION_TEAM_WORKFLOWS` | Get organization team workflows | read | `orgId`, `teamId`; *skip, take* |  | `workflow.read` |
| `CAL_GET_ORGANIZATION_USER_SCHEDULES` | Get organization user schedules | read | `orgId`, `userId` |  |  |
| `CAL_GET_PRIVATE_LINKS` | Get private links for team event type | read | `event_type_id`, `org_id`, `team_id` |  | `EVENT_TYPE_READ` |
| `CAL_GET_SCHEDULE_FOR_USER_IN_TEAM` | Get schedule for user in team | read | `orgId`, `teamId`, `userId` |  |  |
| `CAL_GET_SELECTED_CALENDAR_BY_ID` | Get selected calendar by ID | read | `external_id`, `integration`, `user_id` |  |  |
| `CAL_GET_STRIPE_CONNECT_INFO` | Get stripe connect info | read |  |  |  |
| `CAL_GET_STRIPE_CONNECT_URL` | Get Stripe Connect URL for team | read | `onErrorReturnTo`, `orgId`, `returnTo`, `teamId` |  |  |
| `CAL_GET_TEAMS_LIST` | Get teams list | read |  |  |  |
| `CAL_GET_TEAM_BOOKINGS` | Get team bookings | read | `teamId`; *afterStart, attendeeEmail, attendeeName, beforeEnd, bookingUid, eventTypeId, eventTypeIds, skip and 5 more* |  |  |
| `CAL_GET_TEAM_DEFAULT_CONFERENCING_APP` | Get team default conferencing app | read | `orgId`, `teamId` |  |  |
| `CAL_GET_TEAM_DETAILS_BY_ORG_ID_AND_TEAM_ID` | Get team details by organization ID and team ID | read | `orgId`, `teamId` |  |  |
| `CAL_GET_TEAM_EVENT_TYPE_WEBHOOK` | Get team event type webhook | read | `eventTypeId`, `teamId`, `webhookId` |  | `EVENT_TYPE_READ` |
| `CAL_GET_TEAM_EVENT_TYPE_WEBHOOKS` | Get team event type webhooks | read | `eventTypeId`, `teamId`; *skip, take* |  | `EVENT_TYPE_READ` |
| `CAL_GET_TEAM_INFORMATION_BY_TEAM_ID` | Get team information by team ID | read | `teamId` |  |  |
| `CAL_GET_TEAM_ROUTING_FORMS` | Get team routing forms | read | `orgId`, `teamId`; *afterCreatedAt, afterUpdatedAt, beforeCreatedAt, beforeUpdatedAt, routedToBookingUid, skip, sortCreatedAt, sortUpdatedAt and 1 more* |  |  |
| `CAL_GET_TEAM_SCHEDULES` | Get team schedules | read | `orgId`, `teamId`; *eventTypeId, skip, take* |  | `SCHEDULE_READ` |
| `CAL_GET_VERIFIED_PHONE_NUMBERS` | Get verified phone numbers | read | `orgId`, `teamId`; *skip, take* |  |  |
| `CAL_GET_WEBHOOK_BY_ID` | Get webhook by id | read | `webhookId` |  |  |
| `CAL_HANDLE_CONFERENCING_OAUTH_CALLBACK_FOR_APP` | Handle conferencing oauth callback for app | write | `app`, `code`, `state` |  |  |
| `CAL_LIST_ACCOUNT_WORKFLOWS` | List account workflows | read | *skip, take* |  |  |
| `CAL_LIST_ATTENDEES` | List all attendees | read |  |  |  |
| `CAL_LIST_BOOKING_REFERENCES` | List booking references | read | *cursor, take* |  | `READ_BOOKING` |
| `CAL_LIST_EVENT_TYPES` | List event types | read | *eventSlug, orgId, orgSlug, username, usernames* |  |  |
| `CAL_LIST_ORGANIZATION_MEMBERSHIPS` | List organization memberships | read | `orgId`; *skip, take* |  |  |
| `CAL_LIST_TEAM_EVENT_TYPES_BY_ORG_AND_TEAM_ID` | List team event types by org and team id | read | `orgId`, `teamId`; *eventSlug* |  |  |
| `CAL_MARK_BOOKING_ABSENT_FOR_UID` | Mark booking absent for UID | write | `bookingUid`; *attendees, host* |  |  |
| `CAL_MODIFY_ORGANIZATION_MEMBERSHIP_BY_ID` | Modify organization membership by id | write | `membershipId`, `orgId`; *accepted, disableImpersonation, role* |  |  |
| `CAL_MODIFY_ORG_ATTRIBUTE_BY_ID` | Modify org attribute by id | write | `attributeId`, `orgId`; *enabled, name, slug, type* |  |  |
| `CAL_PATCH_ORGANIZATION_ATTRIBUTE_OPTION` | Patch organization attribute option | write | `attributeId`, `optionId`, `orgId`; *slug, value* |  |  |
| `CAL_PATCH_ORGANIZATION_USER_DETAILS` | Patch organization user details | write | `orgId`, `userId`; *allowDynamicBooking, avatarUrl, bio, brandColor, darkBrandColor, defaultScheduleId, email, hideBranding and 6 more* |  |  |
| `CAL_PATCH_TEAM_DETAILS_BY_ID` | Patch team details by ID | write | `teamId`; *appIconLogo, appLogo, bannerUrl, bio, bookingLimits, brandColor, calVideoLogo, darkBrandColor and 12 more* |  |  |
| `CAL_PATCH_WEBHOOK_EVENT_TYPE` | Patch webhook event type | write | `eventTypeId`, `webhookId`; *active, payloadTemplate, secret, subscriberUrl, triggers* |  |  |
| `CAL_POST_CALENDAR_CREDENTIALS` | Post calendar credentials | write | `calendar`, `password`, `username` |  |  |
| `CAL_POST_CONFERENCING_APP_CONNECT` | Connect conferencing app | write | `app` |  |  |
| `CAL_POST_NEW_BOOKING_REQUEST` | Create a new booking | write | `language`, `responses`, `start`, `timeZone`; *bookingFieldsResponses, emailVerificationCode, eventTypeId, eventTypeSlug, guests, instant, lengthInMinutes, location and 5 more* |  |  |
| `CAL_POST_ORG_ATTRIBUTE_OPTION_BY_USER` | Assign or create attribute option for user | write | `attributeId`, `orgId`, `userId`; *attributeOptionId, value* |  |  |
| `CAL_POST_SELECTED_CALENDARS` | Add selected calendar | write | `credentialId`, `externalId`, `integration` |  |  |
| `CAL_POST_USER_TO_ORGANIZATION` | Post user to organization | write | `email`, `orgId`; *appTheme, autoAccept, avatarUrl, brandColor, darkBrandColor, defaultScheduleId, hideBranding, locale and 6 more* |  |  |
| `CAL_POST_WEBHOOK_EVENT` | Create webhook subscription | write | `active`, `subscriberUrl`, `triggers`; *payloadTemplate, secret* |  |  |
| `CAL_REASSIGN_BOOKING_TO_ANOTHER_USER` | Reassign booking to another user | write | `bookingUid`, `userId`; *reason* |  |  |
| `CAL_REASSIGN_BOOKING_WITH_UID` | Reassign booking with uid | write | `bookingUid` |  |  |
| `CAL_REQUEST_EMAIL_VERIFICATION_CODE` | Request email verification code | write | `email`, `teamId` |  |  |
| `CAL_RESCHEDULE_BOOKING_BY_UID` | Reschedule booking by uid | write | `bookingUid`, `start`; *emailVerificationCode, rescheduledBy, reschedulingReason* |  |  |
| `CAL_RESERVE_SLOT_FOR_EVENT` | Reserve slot for event | write | `eventTypeId`, `slotUtcEndDate`, `slotUtcStartDate`; *bookingUid* |  |  |
| `CAL_RETRIEVE_ATTRIBUTE_OPTIONS_FOR_ORG` | Retrieve attribute options for org | read | `attributeId`, `orgId` |  |  |
| `CAL_RETRIEVE_BOOKING_DETAILS_BY_UID` | Retrieve booking details by uid | read | `bookingUid` |  |  |
| `CAL_RETRIEVE_CALENDAR_BUSY_TIMES` | Retrieve calendar busy times | read | `credentialId`, `externalId`, `loggedInUsersTz`; *dateFrom, dateTo* |  |  |
| `CAL_RETRIEVE_CALENDAR_LIST` | Retrieve calendar list | read |  |  |  |
| `CAL_RETRIEVE_CURRENT_TEAM_FOR_ORGANIZATION` | Retrieve current team for organization | read | `orgId`; *skip, take* |  |  |
| `CAL_RETRIEVE_DEFAULT_CONFERENCING_SETTINGS` | Retrieve default conferencing settings | read |  |  |  |
| `CAL_RETRIEVE_EVENT_TYPE_BY_ID` | Retrieve event type by id | read | `eventTypeId` |  |  |
| `CAL_RETRIEVE_MEMBERSHIP_FROM_ORGANIZATION` | Retrieve membership from organization | read | `membershipId`, `orgId` |  |  |
| `CAL_RETRIEVE_MY_INFORMATION` | Retrieve my information | read |  |  |  |
| `CAL_RETRIEVE_OAUTH_CLIENT_USER_BY_ID` | Retrieve OAuth client user by ID | read | `clientId`, `userId`, `x_cal_secret_key` |  |  |
| `CAL_RETRIEVE_ORGANIZATION_ATTRIBUTES` | Retrieve organization attributes | read | `orgId`; *skip, take* |  |  |
| `CAL_RETRIEVE_ORGANIZATION_ATTRIBUTES_OPTIONS` | Retrieve organization attributes options | read | `orgId`, `userId` |  |  |
| `CAL_RETRIEVE_ORGANIZATION_WEBHOOKS_BY_ORG_ID` | Retrieve organization webhooks by org ID | read | `orgId`; *skip, take* |  |  |
| `CAL_RETRIEVE_ORGANIZATION_WEBHOOK_BY_ID` | Retrieve organization webhook by id | read | `orgId`, `webhookId` |  |  |
| `CAL_RETRIEVE_PROVIDER_DETAILS` | Retrieve provider details | read | `clientId`, `x_cal_secret_key` |  |  |
| `CAL_RETRIEVE_SCHEDULES_LIST` | Retrieve schedules list | read |  |  |  |
| `CAL_RETRIEVE_TEAM_DETAILS_IN_ORGANIZATION` | Retrieve team details in organization | read | `orgId`; *skip, take* |  |  |
| `CAL_RETRIEVE_TEAM_EVENT_TYPES` | Retrieve team event types | read | `teamId`; *eventSlug, hostsLimit* |  |  |
| `CAL_RETRIEVE_TEAM_MEMBERSHIPS` | Retrieve team memberships | read | `teamId`; *skip, take* |  |  |
| `CAL_RETRIEVE_TEAM_MEMBERSHIPS_FOR_ORGANIZATION` | Retrieve team memberships for organization | read | `orgId`, `teamId`; *skip, take* |  |  |
| `CAL_RETRIEVE_TEAM_MEMBERSHIP_BY_ID` | Retrieve team membership by id | read | `membershipId`, `orgId`, `teamId` |  |  |
| `CAL_RETRIEVE_TEAM_MEMBERSHIP_DETAILS` | Retrieve team membership details | read | `membershipId`, `teamId` |  |  |
| `CAL_RETRIEVE_USERS_IN_ORGANIZATION` | Retrieve users in organization | read | `orgId`; *emails, skip, take* |  |  |
| `CAL_RETRIEVE_V2_CONFERENCING_INFO` | Retrieve v2 conferencing info | read |  |  |  |
| `CAL_RETRIEVE_WEBHOOKS_FOR_EVENT_TYPE` | Retrieve webhooks for event type | read | `eventTypeId`; *skip, take* |  |  |
| `CAL_RETRIEVE_WEBHOOKS_LIST` | Retrieve webhooks list | read | *skip, take* |  |  |
| `CAL_RETRIEVE_WEBHOOK_DETAILS_FOR_OAUTH_CLIENT` | Retrieve webhook details for OAuth client | read | `clientId`; *skip, take, x_cal_secret_key* |  |  |
| `CAL_SAVE_CALENDAR_ENTRY` | Save calendar entry | write | `calendar`, `code`, `state` |  |  |
| `CAL_SAVE_CALENDAR_ICS_FEEDS` | Save calendar ics feeds | write | `urls`; *readOnly* |  |  |
| `CAL_SAVE_OAUTH_CREDENTIALS_VIA_GCAL_API` | Save OAuth credentials via GCal API | write | `code`, `state` |  |  |
| `CAL_SAVE_STRIPE_DETAILS` | Save stripe details | write | `code`, `state` |  |  |
| `CAL_SET_DEFAULT_CONFERENCING_APP` | Set default conferencing app | write | `app` |  |  |
| `CAL_UPDATE_DESTINATION_CALENDAR_INTEGRATION` | Update destination calendar integration | write | `externalId`, `integration` |  |  |
| `CAL_UPDATE_OAUTH_CLIENT_USER_SETTINGS` | Update OAuth client user settings | write | `clientId`, `userId`; *avatarUrl, bio, defaultScheduleId, email, locale, metadata, name, timeFormat and 3 more* |  |  |
| `CAL_UPDATE_OAUTH_CLIENT_WEBHOOK` | Update oauth client webhook | write | `clientId`, `webhookId`; *active, payloadTemplate, secret, subscriberUrl, triggers, x_cal_secret_key* |  |  |
| `CAL_UPDATE_PRIVATE_LINK` | Update private link | write | `eventTypeId`, `linkId`, `orgId`, `teamId`; *expiresAt, maxUsageCount* |  | `eventType.update` |
| `CAL_UPDATE_SCHEDULE_BY_ID` | Update schedule by ID | write | `scheduleId`; *availability, isDefault, name, overrides, timeZone* |  |  |
| `CAL_UPDATE_TEAM_EVENT_TYPE` | Update team event type | write | `eventTypeId`, `teamId`; *afterEventBuffer, allowReschedulingCancelledBookings, allowReschedulingPastBookings, assignAllTeamMembers, beforeEventBuffer, bookerActiveBookingsLimit, bookerLayouts, bookingFields and 42 more* |  | `EVENT_TYPE_WRITE` |
| `CAL_UPDATE_TEAM_EVENT_TYPE_WEBHOOK` | Update team event type webhook | write | `eventTypeId`, `teamId`, `webhookId`; *active, payloadTemplate, secret, subscriberUrl, triggers, version* |  | `EVENT_TYPE_WRITE` |
| `CAL_UPDATE_TEAM_INFORMATION_BY_ID` | Update team information by id | write | `orgId`, `teamId`; *appIconLogo, appLogo, bannerUrl, bio, bookingLimits, brandColor, calVideoLogo, darkBrandColor and 12 more* |  |  |
| `CAL_UPDATE_TEAM_MEMBERSHIP_BY_ID` | Update team membership by id | write | `membershipId`, `teamId`; *accepted, disableImpersonation, role* |  |  |
| `CAL_UPDATE_TEAM_MEMBERSHIP_PROPERTIES` | Update team membership properties | write | `membershipId`, `orgId`, `teamId`; *accepted, disableImpersonation, role* |  |  |
| `CAL_UPDATE_USER_PROFILE_DETAILS` | Update user profile details | write | *avatarUrl, defaultScheduleId, email, locale, name, timeFormat, timeZone, weekStart* |  |  |
| `CAL_UPDATE_USER_SCHEDULE_IN_ORGANIZATION` | Update user schedule in organization | write | `orgId`, `scheduleId`, `userId`; *availability, isDefault, name, overrides, timeZone* |  |  |
| `CAL_UPDATE_WEBHOOK_BY_ID` | Update webhook by id | write | `webhookId`; *active, payloadTemplate, secret, subscriberUrl, triggers* |  |  |
| `CAL_UPDATE_WEBHOOK_FOR_ORGANIZATION` | Update webhook for organization | write | `orgId`, `webhookId`; *active, payloadTemplate, secret, subscriberUrl, triggers* |  |  |
