# LinkedIn: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner signs in on the vendor's own page). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

29 tools: 16 read, 13 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `LINKEDIN_CREATE_ARTICLE_OR_URL_SHARE` | Create article or URL share | write | `author`, `specificContent`, `visibility`; *lifecycleState* |  | `w_member_social`, `w_organization_social` |
| `LINKEDIN_CREATE_COMMENT_ON_POST` | Create comment on LinkedIn post | write | `actor`, `message`, `object`, `target_urn`; *content, parent_comment* |  | `w_member_social`, `w_organization_social` |
| `LINKEDIN_CREATE_LINKED_IN_POST` | Create a LinkedIn post | write | `author`, `commentary`; *container, contentCallToActionLabel, contentLandingPage, distribution, images, isReshareDisabledByAuthor, lifecycleState, reshareContext and 1 more* |  | `w_member_social`, `w_organization_social` |
| `LINKEDIN_CREATE_VIDEO_POST` | Create video post | write | `commentary`, `video_urn`; *title, visibility* |  | `openid`, `profile`, `w_member_social` |
| `LINKEDIN_DELETE_COMMENT` | Delete Comment | write | `comment_id`, `entity_urn`; *actor* |  | `w_member_social`, `w_organization_social`, `w_member_social_feed`, `w_organization_social_feed` |
| `LINKEDIN_DELETE_LINKED_IN_POST` | Delete LinkedIn Post | write | `share_id` |  | `w_member_social`, `w_organization_social` |
| `LINKEDIN_DELETE_POST` | Delete Post | write | `post_urn` |  | `w_member_social`, `w_organization_social` |
| `LINKEDIN_DELETE_UGC_POST` | Delete UGC Post (Legacy) | write | `ugc_post_urn` |  | `w_member_social`, `w_organization_social`, `w_compliance` |
| `LINKEDIN_GET_AD_TARGETING_FACETS` | Get ad targeting facets | read |  |  |  |
| `LINKEDIN_GET_AUDIENCE_COUNTS` | Get audience counts | read | `targetingCriteria`; *q* |  |  |
| `LINKEDIN_GET_COMPANY_INFO` | Get company info | read | *count, role, start, state* |  | `r_organization_admin` |
| `LINKEDIN_GET_GROUP_DEFINITION` | Get group definition | read | `group_id` |  |  |
| `LINKEDIN_GET_IMAGE` | Get image details | read | `image_urn` |  | `r_organization_social`, `rw_ads`, `w_member_social` |
| `LINKEDIN_GET_IMAGES` | Get images | read | `ids` |  | `r_organization_social`, `rw_ads`, `w_member_social` |
| `LINKEDIN_GET_MY_INFO` | Get my info | read |  |  | `profile`, `r_liteprofile`, `r_basicprofile` |
| `LINKEDIN_GET_NETWORK_SIZE` | Get network size | read | `organization_id`; *edgeType* |  | `rw_organization_admin` |
| `LINKEDIN_GET_ORG_PAGE_STATS` | Get organization page statistics | read | `organization`; *timeGranularityType, timeRangeEnd, timeRangeStart* |  | `rw_organization_admin` |
| `LINKEDIN_GET_PERSON` | Get person profile | read | `person_id` |  | `profile`, `r_liteprofile`, `r_basicprofile` |
| `LINKEDIN_GET_POST_CONTENT` | Get post content | read | `post_id` |  | `r_organization_social`, `r_member_social` |
| `LINKEDIN_GET_SHARE_STATS` | Get share statistics | read | `organizational_entity`; *time_intervals* |  | `rw_organization_admin` |
| `LINKEDIN_GET_VIDEOS` | Get videos | read | *video_ids, video_urn* |  | `w_organization_social`, `w_member_social` |
| `LINKEDIN_INITIALIZE_DOCUMENT_UPLOAD` | Initialize document upload | write | `owner` |  | `w_member_social`, `w_organization_social` |
| `LINKEDIN_INITIALIZE_IMAGE_UPLOAD` | Initialize image upload | write | `owner` |  | `rw_ads`, `w_member_social`, `w_organization_social`, `w_power_creators` |
| `LINKEDIN_LIST_REACTIONS` | List reactions on entity | read | `entity`; *count, sort, start* |  | `r_organization_social_feed`, `r_member_social_feed` |
| `LINKEDIN_REGISTER_IMAGE_UPLOAD` | Register image upload | write | `owner_urn`; *recipe, supported_upload_mechanism* |  | `w_member_social`, `w_organization_social`, `rw_ads`, `w_compliance` |
| `LINKEDIN_SEARCH_AD_TARGETING_ENTITIES` | Search ad targeting entities | read | `facet`, `query`; *count, queryVersion, start* |  | `r_ads`, `rw_ads` |
| `LINKEDIN_UPLOAD_DOCUMENT` | Upload document | write | *document_url, file* |  | `openid`, `profile`, `w_member_social` |
| `LINKEDIN_UPLOAD_VIDEO` | Upload video | write | *file, video_url* |  | `openid`, `profile`, `w_member_social` |
| `LINKEDIN_WHO_AM_I` | Who Am I | read |  |  |  |
