# YouTube: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

67 tools: 32 read, 35 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `YOUTUBE_ADD_VIDEO_TO_PLAYLIST` | Add Video to Playlist | write | `playlistId`, `videoId`; *position* |  | `youtube`, `youtube.force-ssl`, `youtube.readonly`, `youtubepartner` |
| `YOUTUBE_BIND_LIVE_BROADCAST` | Bind Live Broadcast | write | `id`, `part`; *onBehalfOfContentOwner, onBehalfOfContentOwnerChannel, streamId* |  | `youtube`, `youtube.force-ssl` |
| `YOUTUBE_CREATE_CHANNEL_SECTION` | Insert Channel Section | write | `snippet`; *contentDetails* |  | `youtubepartner`, `youtube`, `youtube.force-ssl` |
| `YOUTUBE_CREATE_COMMENT_REPLY` | Insert Comment Reply | write | `parentId`, `textOriginal` |  | `youtube.force-ssl` |
| `YOUTUBE_CREATE_LIVE_BROADCAST` | Create Live Broadcast | write | `privacyStatus`, `scheduledStartTime`, `title`; *availabilityConfig, broadcastStreamDelayMs, categoryId, description, enableAutoStart, enableAutoStop, enableClosedCaptions, enableDvr and 7 more* |  | `youtube`, `youtube.force-ssl` |
| `YOUTUBE_CREATE_LIVE_STREAM` | Create YouTube Live Stream | write | `frame_rate`, `ingestion_type`, `resolution`, `title`; *description, is_reusable, on_behalf_of_content_owner, on_behalf_of_content_owner_channel* |  | `youtube`, `youtube.force-ssl` |
| `YOUTUBE_CREATE_PLAYLIST` | Create Playlist | write | `title`; *description, privacyStatus* |  | `youtube`, `youtube.force-ssl`, `youtube.readonly`, `youtubepartner` |
| `YOUTUBE_DELETE_CAPTION_TRACK` | Delete caption track | write | `confirmDelete`, `id`; *onBehalfOfContentOwner* |  | `youtube.force-ssl`, `youtubepartner` |
| `YOUTUBE_DELETE_CHANNEL_SECTION` | Delete Channel Section | write | `id`; *onBehalfOfContentOwner* |  | `youtubepartner`, `youtube`, `youtube.force-ssl` |
| `YOUTUBE_DELETE_COMMENT` | Delete Comment | write | `id` |  | `youtube.force-ssl`, `youtube` |
| `YOUTUBE_DELETE_LIVE_BROADCAST` | Delete Live Broadcast | write | `confirmDelete`, `id`; *onBehalfOfContentOwner, onBehalfOfContentOwnerChannel* |  | `youtube`, `youtube.force-ssl` |
| `YOUTUBE_DELETE_PLAYLIST` | Delete Playlist | write | `confirmDelete`, `id`; *onBehalfOfContentOwner* |  | `youtubepartner`, `youtube`, `youtube.force-ssl` |
| `YOUTUBE_DELETE_PLAYLIST_ITEM` | Delete Playlist Item | write | `id` |  | `youtubepartner`, `youtube`, `youtube.force-ssl` |
| `YOUTUBE_DELETE_VIDEO` | Delete Video | write | `confirmDelete`, `videoId`; *onBehalfOfContentOwner* |  | `youtube`, `youtube.force-ssl`, `youtubepartner` |
| `YOUTUBE_GET_CHANNEL_ACTIVITIES` | Get Channel Activities | read | `channelId`; *maxResults, pageToken, part, publishedAfter, publishedBefore* |  | `youtube.readonly` |
| `YOUTUBE_GET_CHANNEL_ID_BY_HANDLE` | Get channel ID by handle | read | `channel_handle` |  | `youtube.readonly` |
| `YOUTUBE_GET_CHANNEL_STATISTICS` | Get Channel Statistics | read | *forHandle, forUsername, id, mine, part* |  | `youtube.readonly` |
| `YOUTUBE_GET_PLAYLISTS_BY_ID` | Get Playlists by ID | read | `id`, `part`; *hl, maxResults, onBehalfOfContentOwner, onBehalfOfContentOwnerChannel, pageToken* |  | `youtube.readonly`, `youtube`, `youtube.force-ssl` |
| `YOUTUBE_GET_VIDEO_DETAILS_BATCH` | Video Details Batch | read | `id`; *hl, parts* |  | `youtube.readonly` |
| `YOUTUBE_GET_VIDEO_RATING` | Get Video Rating | read | `id`; *onBehalfOfContentOwner* |  | `youtubepartner`, `youtube`, `youtube.force-ssl` |
| `YOUTUBE_LIST_CAPTION_TRACK` | List captions | read | `video_id`; *part* |  | `youtube`, `youtube.force-ssl`, `youtube.readonly` |
| `YOUTUBE_LIST_CHANNELS` | List Channels | read | *forHandle, forUsername, hl, id, managedByMe, maxResults, mine, onBehalfOfContentOwner and 2 more* |  | `youtube.readonly`, `youtubepartner-channel-audit` |
| `YOUTUBE_LIST_CHANNEL_SECTIONS` | List Channel Sections | read | `part`; *channelId, hl, id, mine, onBehalfOfContentOwner* |  | `youtube`, `youtube.force-ssl`, `youtube.readonly`, `youtubepartner` |
| `YOUTUBE_LIST_CHANNEL_VIDEOS` | List channel videos | read | *channelId, maxResults, mine, pageToken, part* |  | `youtube.readonly` |
| `YOUTUBE_LIST_COMMENTS` | List Comments | read | *id, maxResults, pageToken, parentId, part, textFormat* |  | `youtube.force-ssl`, `youtube.readonly` |
| `YOUTUBE_LIST_COMMENT_THREADS2` | List Comment Threads | read | `part`; *allThreadsRelatedToChannelId, channelId, id, maxResults, moderationStatus, order, pageToken, searchTerms and 2 more* |  | `youtube.force-ssl` |
| `YOUTUBE_LIST_I18N_LANGUAGES` | List I18n Languages | read | *hl, part* |  | `youtube.readonly` |
| `YOUTUBE_LIST_I18N_REGIONS` | List I18n Regions | read | *hl, part* |  | `youtube`, `youtube.force-ssl`, `youtube.readonly`, `youtubepartner` |
| `YOUTUBE_LIST_JOB_REPORTS` | List YouTube Job Reports | read | `jobId`; *createdAfter, onBehalfOfContentOwner, pageSize, pageToken, startTimeAtOrAfter, startTimeBefore* |  | `yt-analytics.readonly`, `yt-analytics-monetary.readonly` |
| `YOUTUBE_LIST_LIVE_BROADCASTS` | List Live Broadcasts | read | `part`; *broadcastStatus, broadcastType, id, maxResults, mine, onBehalfOfContentOwner, onBehalfOfContentOwnerChannel, pageToken* |  | `youtube.readonly`, `youtube`, `youtube.force-ssl` |
| `YOUTUBE_LIST_LIVE_CHAT_MESSAGES` | List Live Chat Messages | read | `liveChatId`; *hl, maxResults, pageToken, part, profileImageSize* |  | `youtube`, `youtube.force-ssl`, `youtube.readonly` |
| `YOUTUBE_LIST_LIVE_STREAMS` | List Live Streams | read | `part`; *id, maxResults, mine, onBehalfOfContentOwner, onBehalfOfContentOwnerChannel, pageToken* |  | `youtube.readonly`, `youtube`, `youtube.force-ssl` |
| `YOUTUBE_LIST_MOST_POPULAR_VIDEOS` | List Most Popular Videos | read | *chart, maxResults, pageToken, part, regionCode, videoCategoryId* |  | `youtube`, `youtube.force-ssl`, `youtube.readonly`, `youtubepartner` |
| `YOUTUBE_LIST_PLAYLIST_IMAGES` | List Playlist Images | read | *id, maxResults, onBehalfOfContentOwner, onBehalfOfContentOwnerChannel, pageToken, parent, part* |  | `youtubepartner`, `youtube`, `youtube.force-ssl` |
| `YOUTUBE_LIST_PLAYLIST_ITEMS` | List Playlist Items | read | `playlistId`; *fields, maxResults, onBehalfOfContentOwner, pageToken, part, videoId* |  | `youtube.readonly` |
| `YOUTUBE_LIST_RATED_VIDEOS` | List Rated Videos | read | `part`, `rating`; *content_owner_id, localization_language, max_height, max_results, max_width, page_token* |  | `youtube.readonly`, `youtube`, `youtube.force-ssl`, `youtubepartner` |
| `YOUTUBE_LIST_REPORTING_JOBS` | List YouTube Reporting Jobs | read | *includeSystemManaged, onBehalfOfContentOwner, pageSize, pageToken* |  | `yt-analytics.readonly`, `yt-analytics-monetary.readonly` |
| `YOUTUBE_LIST_REPORT_TYPES` | List YouTube Report Types | read | *includeSystemManaged, onBehalfOfContentOwner, pageSize, pageToken* |  | `yt-analytics.readonly`, `yt-analytics-monetary.readonly` |
| `YOUTUBE_LIST_SUPER_CHAT_EVENTS` | List Super Chat Events | read | *hl, maxResults, pageToken, part* |  | `youtube`, `youtube.force-ssl`, `youtube.readonly`, `youtubepartner` |
| `YOUTUBE_LIST_USER_PLAYLISTS` | List user playlists | read | *maxResults, pageToken, part* |  | `youtube.readonly` |
| `YOUTUBE_LIST_USER_SUBSCRIPTIONS` | List user subscriptions | read | *maxResults, pageToken, part* |  | `youtube`, `youtube.force-ssl`, `youtube.readonly`, `youtubepartner` |
| `YOUTUBE_LIST_VIDEO_ABUSE_REPORT_REASONS` | List Video Abuse Report Reasons | read | *hl, part* |  | `youtube`, `youtube.force-ssl`, `youtube.upload`, `youtubepartner` |
| `YOUTUBE_LIST_VIDEO_CATEGORIES` | List Video Categories | read | *hl, id, part, regionCode* |  | `youtube.readonly` |
| `YOUTUBE_LOAD_CAPTIONS` | Download YouTube caption track | read | `id`; *tfmt* |  | `youtube.force-ssl` |
| `YOUTUBE_MULTIPART_UPLOAD_VIDEO` | Multipart upload video | write | `categoryId`, `description`, `privacyStatus`, `title`, `videoFile`; *tags* |  | `youtube.upload`, `youtube`, `youtubepartner`, `youtube.force-ssl` |
| `YOUTUBE_POST_COMMENT` | Post Comment on Video | write | `channelId`, `textOriginal`, `videoId` |  | `youtube.force-ssl` |
| `YOUTUBE_QUERY_ANALYTICS` | Query YouTube Analytics | read | `endDate`, `ids`, `metrics`, `startDate`; *dimensions, filters, includeHistoricalChannelData, maxResults, sort, startIndex* |  | `youtube.readonly`, `youtube`, `youtubepartner`, `yt-analytics-monetary.readonly` |
| `YOUTUBE_RATE_VIDEO` | Rate Video | write | `id`, `rating` |  | `youtubepartner`, `youtube`, `youtube.force-ssl` |
| `YOUTUBE_REPORT_VIDEO_ABUSE` | Report Video for Abuse | write | `reasonId`, `videoId`; *comments, language, secondaryReasonId* |  | `youtube.force-ssl` |
| `YOUTUBE_SEARCH_YOU_TUBE` | Search YouTube | read | `q`; *maxResults, order, pageToken, part, publishedAfter, publishedBefore, regionCode, relevanceLanguage and 1 more* |  | `youtube.readonly` |
| `YOUTUBE_SET_COMMENT_MODERATION_STATUS` | Set Comment Moderation Status | write | `id`, `moderationStatus`; *banAuthor* |  | `youtube.force-ssl` |
| `YOUTUBE_SUBSCRIBE_CHANNEL` | Subscribe to channel | write | `channelId` |  | `youtube`, `youtube.force-ssl`, `youtube.readonly`, `youtubepartner` |
| `YOUTUBE_TRANSITION_LIVE_BROADCAST` | Transition Live Broadcast | write | `broadcastStatus`, `id`, `part`; *onBehalfOfContentOwner, onBehalfOfContentOwnerChannel* |  | `youtube`, `youtube.force-ssl` |
| `YOUTUBE_UNSUBSCRIBE_CHANNEL` | Unsubscribe from channel | write | `subscriptionId` |  | `youtube`, `youtube.force-ssl`, `youtubepartner` |
| `YOUTUBE_UPDATE_CAPTION` | Update caption track | write | `id`, `snippet` |  | `youtube.force-ssl`, `youtubepartner` |
| `YOUTUBE_UPDATE_CHANNEL` | Update channel | write | `id`; *brandingSettings, invideoPromotion, localizations, onBehalfOfContentOwner, part* |  | `youtube`, `youtube.force-ssl`, `youtube.readonly`, `youtubepartner`, `youtubepartner-channel-audit` |
| `YOUTUBE_UPDATE_CHANNEL_SECTION` | Update Channel Section | write | `id`; *contentDetails, snippet* |  | `youtubepartner`, `youtube`, `youtube.force-ssl` |
| `YOUTUBE_UPDATE_COMMENT` | Update Comment | write | `id`, `textOriginal` |  | `youtube.force-ssl` |
| `YOUTUBE_UPDATE_LIVE_BROADCAST` | Update Live Broadcast | write | `broadcastId`; *availabilityConfig, broadcastStreamDelayMs, categoryId, description, enableAutoStart, enableAutoStop, enableClosedCaptions, enableDvr and 9 more* |  | `youtube`, `youtube.force-ssl` |
| `YOUTUBE_UPDATE_PLAYLIST` | Update Playlist | write | `id`, `snippet`; *onBehalfOfContentOwner, part, status* |  | `youtubepartner`, `youtube`, `youtube.force-ssl` |
| `YOUTUBE_UPDATE_PLAYLIST_ITEM` | Update Playlist Item | write | `id`, `snippet`; *contentDetails, onBehalfOfContentOwner, part, status* |  | `youtubepartner`, `youtube`, `youtube.force-ssl` |
| `YOUTUBE_UPDATE_THUMBNAIL` | Update thumbnail | write | `thumbnailUrl`, `videoId` |  | `youtube`, `youtube.upload` |
| `YOUTUBE_UPDATE_VIDEO` | Update video | write | `video_id`; *category_id, composio_execution_message, description, privacy_status, tags, title* |  | `youtube.force-ssl`, `youtube.upload` |
| `YOUTUBE_UPLOAD_CAPTION_MEDIA` | Upload caption media | write | `captionFile`, `id`; *onBehalfOfContentOwner* |  | `youtube.force-ssl`, `youtubepartner` |
| `YOUTUBE_UPLOAD_CAPTION_TRACK` | Upload caption track | write | `file`, `language`, `name`, `videoId`; *isDraft, onBehalfOfContentOwner* |  | `youtube.force-ssl`, `youtubepartner` |
| `YOUTUBE_UPLOAD_CHANNEL_BANNER` | Upload channel banner | write | `bannerFile`; *onBehalfOfContentOwner* |  | `youtube.upload`, `youtube`, `youtube.force-ssl` |
| `YOUTUBE_UPLOAD_VIDEO` | Upload video | write | `categoryId`, `description`, `privacyStatus`, `tags`, `title`, `videoFilePath` |  | `youtube`, `youtube.force-ssl`, `youtube.readonly`, `youtubepartner` |
