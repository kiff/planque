# Ashby: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner pastes a key from the vendor's settings). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

191 tools: 81 read, 110 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `ASHBY_ADD_CANDIDATE_PROJECT` | Add Candidate to Project | write | `candidateId`, `projectId` |  |  |
| `ASHBY_ADD_CANDIDATE_TAG` | Add Candidate Tag | write | `candidateId`, `tagId` |  |  |
| `ASHBY_ADD_COMPLETED_TO_CANDIDATE_ASSESSMENT` | Add Completed Assessment to Candidate | write | `assessment`, `candidateId`, `partnerId`, `timestamp` |  | `candidatesWrite` |
| `ASHBY_ADD_EMAIL_MESSAGE_CANDIDATE` | Add Existing Email Message to Candidate | write | `body`, `candidateId`, `emailProviderEmailId`, `from`, `subject`, `to`; *cc, isPrivate, messageIdHeader, messageUrl, sentAt, threadId, userId* |  | `candidatesWrite` |
| `ASHBY_ADD_HIRING_TEAM_MEMBER` | Add Hiring Team Member | write | `roleId`, `teamMemberId`; *applicationId, jobId, openingId* |  |  |
| `ASHBY_ADD_HIRING_TEAM_MEMBER_APPLICATION` | Add Hiring Team Member to Application | write | `applicationId`, `roleId`, `teamMemberId` |  | `candidatesWrite` |
| `ASHBY_ADD_INTERVIEWER_POOL_USER` | Add User to Interviewer Pool | write | `interviewerPoolId`, `userId`; *interviewerPoolTrainingPathStageId* |  |  |
| `ASHBY_ADD_OPENING_JOB` | Add Opening Job | write | `jobId`, `openingId` |  | `jobsWrite` |
| `ASHBY_ADD_OPENING_LOCATION` | Add Opening Location | write | `locationId`, `openingId` |  | `jobsWrite` |
| `ASHBY_ADD_SEQUENCE` | Add Sequence | write | `candidateId`, `sequenceTemplateId`, `start`; *applicationId, from* |  | `sourcingWrite` |
| `ASHBY_ANONYMIZE_CANDIDATE` | Anonymize Candidate | write | `candidateId` |  | `candidatesWrite` |
| `ASHBY_APPROVE_OFFER` | Approve Offer | write | `offerId`; *approvalStepId, userId* |  |  |
| `ASHBY_ARCHIVE_CANDIDATE_TAG` | Archive Candidate Tag | write | `tagId` |  | `hiringProcessMetadataWrite` |
| `ASHBY_ARCHIVE_DEPARTMENT` | Archive Department | write | `departmentId` |  |  |
| `ASHBY_ARCHIVE_INTERVIEWER_POOL` | Archive Interviewer Pool | write | `interviewerPoolId` |  | `hiringProcessMetadataWrite` |
| `ASHBY_ARCHIVE_LOCATION` | Archive Location | write | `locationId` |  |  |
| `ASHBY_ARCHIVE_PROJECT` | Archive Project | write | `projectId` |  | `projectsWrite` |
| `ASHBY_CANCEL_INTERVIEW_SCHEDULE` | Cancel Interview Schedule | write | `id`; *allowReschedule* |  | `interviewsWrite` |
| `ASHBY_CANCEL_SEQUENCE` | Cancel Sequence | write | `sequenceId`; *reason* |  | `candidatesWrite`, `sourcingWrite` |
| `ASHBY_CHANGE_APPLICATION_SOURCE` | Change Application Source | write | `applicationId`, `sourceId`; *creditedToUserId* |  |  |
| `ASHBY_CHANGE_APPLICATION_STAGE` | Change Application Stage | write | `applicationId`, `interviewStageId`; *archiveReasonId* |  |  |
| `ASHBY_CREATE_APPLICATION` | Create Application | write | `candidateId`, `jobId`; *creditedToUserId, interviewStageId, sourceId* |  |  |
| `ASHBY_CREATE_APPLICATION_FEEDBACK_REQUEST` | Create Application Feedback Request | write | `applicationId`, `interviewId`, `interviewerUserId` |  | `candidatesWrite` |
| `ASHBY_CREATE_CANDIDATE` | Create Candidate | write | `name`; *email, githubUrl, linkedInUrl, phoneNumber, websiteUrl* |  |  |
| `ASHBY_CREATE_CANDIDATE_NOTE` | Create Candidate Note | write | `candidateId`, `note`; *sendNotifications* |  | `candidatesWrite` |
| `ASHBY_CREATE_CANDIDATE_TAG` | Create Candidate Tag | write | `title` |  |  |
| `ASHBY_CREATE_CUSTOM_FIELD` | Create Custom Field | write | `fieldType`, `objectType`, `title`; *description, isDateOnlyField, isExposableToCandidate, selectableValues* |  | `hiringProcessMetadataWrite` |
| `ASHBY_CREATE_DEPARTMENT` | Create Department | write | `name`; *parentId* |  |  |
| `ASHBY_CREATE_EMAIL_SYNC_BLOCKED_ADDRESS` | Create Email Sync Blocked Address | write | `addresses` |  | `organizationWrite` |
| `ASHBY_CREATE_INTERVIEWER_PAUSE_USER` | Create Interviewer Pause User | write | `userId`; *comment, endsAt, startsAt* |  | `organizationWrite` |
| `ASHBY_CREATE_INTERVIEWER_POOL` | Create Interviewer Pool | write | `title`; *requiresTraining* |  |  |
| `ASHBY_CREATE_INTERVIEW_SCHEDULE` | Create Interview Schedule | write | `applicationId`, `interviewEvents` |  | `interviewsWrite` |
| `ASHBY_CREATE_JOB` | Create Job | write | `brandId`, `locationId`, `teamId`, `title`; *defaultInterviewPlanId, jobTemplateId* |  |  |
| `ASHBY_CREATE_JOB_POSTING` | Create Job Posting | write | `jobId`; *applicationFormDefinitionId, description, isListed, locationId, secondaryLocationIds, status, title, workplaceType* |  | `jobsWrite` |
| `ASHBY_CREATE_LOCATION` | Create Location | write | `name`, `type`; *address, isRemote, parentLocationId, workplaceType* |  |  |
| `ASHBY_CREATE_OFFER` | Create Offer | write | `offerForm`, `offerFormId`, `offerProcessId` |  |  |
| `ASHBY_CREATE_OPENING` | Create Opening | write | *description, employmentType, identifier, isBackfill, jobIds, locationIds, openingState, targetHireDate and 3 more* |  |  |
| `ASHBY_CREATE_PROJECT` | Create Project | write | `title`; *confidential, customFields, description* |  | `projectsWrite` |
| `ASHBY_CREATE_REFERRAL` | Create Referral | write | `creditedToUserId`, `fieldSubmissions`, `id`; *createdAt* |  | `candidatesWrite` |
| `ASHBY_CREATE_SURVEY_REQUEST` | Create Survey Request | write | `applicationId`, `candidateId`, `surveyFormDefinitionId` |  |  |
| `ASHBY_CREATE_SURVEY_SUBMISSION` | Create Survey Submission | write | `applicationId`, `candidateId`, `submittedValues`, `surveyFormDefinitionId` |  |  |
| `ASHBY_DELETE_APPLICATION` | Delete Application | write | `applicationId` |  | `candidatesDelete` |
| `ASHBY_DELETE_EMAIL_SYNC_BLOCKED_ADDRESS` | Delete Email Sync Blocked Address | write | `addresses` |  | `organizationWrite` |
| `ASHBY_DELETE_INTERVIEWER_PAUSE_USER` | Delete User Interviewer Pause | write | `interviewerPauseId` |  | `organizationWrite` |
| `ASHBY_DELETE_REFERRAL` | Delete Referral | write | `referralId`; *clearAttribution* |  | `candidatesDelete` |
| `ASHBY_DISCARD_SEQUENCE` | Discard Sequence Draft | write | `sequenceId` |  | `sourcingWrite` |
| `ASHBY_GENERATE_REPORT` | Generate Report | read | `reportId`; *includeHeadersInData, requestId, resultStyle* |  | `reportsRead` |
| `ASHBY_GET_API_KEY_INFO` | Get API Key Info | read |  |  |  |
| `ASHBY_GET_APPLICATION_INFO` | Get Application Info | read | `applicationId` |  |  |
| `ASHBY_GET_CANDIDATE_INFO` | Get Candidate Info | read | `candidateId` |  |  |
| `ASHBY_GET_CUSTOM_FIELD_INFO` | Get Custom Field Info | write | `customFieldId` |  | `hiringProcessMetadataRead` |
| `ASHBY_GET_DEPARTMENT_INFO` | Get Department Info | read | `departmentId` |  |  |
| `ASHBY_GET_FEEDBACK_FORM_DEFINITION` | Get Feedback Form Definition | read | `feedbackFormDefinitionId` |  | `hiringProcessMetadataRead` |
| `ASHBY_GET_FILE_INFO` | Get File Info | read | `fileHandle` |  | `candidatesRead` |
| `ASHBY_GET_INTERVIEWER_POOL_INFO` | Get Interviewer Pool Info | read | `interviewerPoolId` |  |  |
| `ASHBY_GET_INTERVIEWER_USER_SETTINGS` | Get Interviewer User Settings | read | `userId` |  |  |
| `ASHBY_GET_INTERVIEW_INFO` | Get Interview Info | read | `interviewId` |  |  |
| `ASHBY_GET_INTERVIEW_STAGE_INFO` | Get Interview Stage Info | read | `interviewStageId` |  |  |
| `ASHBY_GET_JOB_INFO` | Get Job Info | read | `jobId` |  |  |
| `ASHBY_GET_JOB_INTERVIEW_PLAN_INFO` | Get Job Interview Plan Info | read | `jobId` |  |  |
| `ASHBY_GET_JOB_POSTING_INFO` | Get Job Posting Info | read | `jobPostingId` |  |  |
| `ASHBY_GET_LOCATION_INFO` | Get Location Info | read | `locationId` |  |  |
| `ASHBY_GET_OFFER_INFO` | Get Offer Info | read | `offerId` |  | `offersRead` |
| `ASHBY_GET_OPENING_INFO` | Get Opening Info | read | `openingId` |  |  |
| `ASHBY_GET_REFERRAL_FORM` | Get Referral Form | read |  |  |  |
| `ASHBY_GET_SURVEY_FORM_DEFINITION` | Get Survey Form Definition | read | `surveyFormDefinitionId` |  | `hiringProcessMetadataRead` |
| `ASHBY_GET_USER_INFO` | Get User Info | read | `userId` |  |  |
| `ASHBY_INFO_INTERVIEW_BRIEFING` | Info Interview Briefing | read | `interviewEventId`; *expand* |  | `interviewsRead` |
| `ASHBY_INFO_NOTETAKER_TRANSCRIPT` | Info Notetaker Transcript | read | `notetakerTranscriptId` |  | `notetakerRead` |
| `ASHBY_INFO_PROJECT` | Info Project | read | `projectId` |  | `projectsRead`, `candidatesRead` |
| `ASHBY_INFO_SEQUENCE` | Info Sequence | read | `sequenceId` |  | `sourcingRead` |
| `ASHBY_INFO_SEQUENCE_TEMPLATE` | Info Sequence Template | read | `sequenceTemplateId` |  | `sourcingRead` |
| `ASHBY_INFO_TAKE_HOME_ASSIGNMENT` | Get Take-Home Assignment Info | read | `takeHomeAssignmentId`; *expand* |  | `interviewsRead` |
| `ASHBY_LIST_APPLICATIONS` | List Applications | read | *createdAfter, cursor, expand, jobId, limit, status, syncToken* |  |  |
| `ASHBY_LIST_APPLICATION_CRITERIA_EVALUATIONS` | List Application Criteria Evaluations | read | `applicationId`; *cursor, limit* |  |  |
| `ASHBY_LIST_APPLICATION_FEEDBACK` | List Application Feedback | read | `applicationId`; *cursor, syncToken* |  |  |
| `ASHBY_LIST_APPLICATION_HIRING_TEAM_ROLE` | List Application Hiring Team Roles | read |  |  |  |
| `ASHBY_LIST_APPLICATION_HISTORY` | List Application History | read | `applicationId`; *cursor, limit* |  |  |
| `ASHBY_LIST_APPROVAL` | List Approvals | read | *cursor, perPage, syncToken* |  |  |
| `ASHBY_LIST_ARCHIVE_REASONS` | List Archive Reasons | read |  |  |  |
| `ASHBY_LIST_BRAND` | List Brands | read | *cursor, limit* |  |  |
| `ASHBY_LIST_CANDIDATES` | List Candidates | read | *cursor, perPage, syncToken* |  |  |
| `ASHBY_LIST_CANDIDATE_CLIENT_INFO` | List Candidate Client Info | read | `candidateId`; *createdAfter, cursor, limit, syncToken* |  | `candidatesRead` |
| `ASHBY_LIST_CANDIDATE_NOTES` | List Candidate Notes | read | `candidateId`; *cursor, limit, syncToken* |  |  |
| `ASHBY_LIST_CANDIDATE_PROJECTS` | List Candidate Projects | read | `candidateId`; *cursor, limit* |  |  |
| `ASHBY_LIST_CANDIDATE_TAGS` | List Candidate Tags | read | *cursor, perPage, syncToken* |  |  |
| `ASHBY_LIST_CLOSE_REASONS` | List Close Reasons | read |  |  |  |
| `ASHBY_LIST_COMMUNICATION_TEMPLATES` | List Communication Templates | read |  |  |  |
| `ASHBY_LIST_CUSTOM_FIELDS` | List Custom Fields | write |  |  |  |
| `ASHBY_LIST_DEPARTMENTS` | List Departments | read | *cursor, limit, syncToken* |  |  |
| `ASHBY_LIST_EMAIL_SENDER` | List Email Sender | read |  |  | `sourcingRead` |
| `ASHBY_LIST_EMAIL_SYNC_BLOCKED_ADDRESS` | List Email Sync Blocked Address | read | *cursor, limit* |  | `organizationRead` |
| `ASHBY_LIST_FEEDBACK_FORM_DEFINITIONS` | List Feedback Form Definitions | read | *cursor, limit* |  |  |
| `ASHBY_LIST_FRAUD_CHECKS_CANDIDATE` | List Candidate Fraud Checks | read | `candidateId`; *cursor, limit* |  | `candidatesRead` |
| `ASHBY_LIST_HIRING_TEAM_ROLE` | List Hiring Team Roles | read | *namesOnly* |  |  |
| `ASHBY_LIST_INTERVIEWER_PAUSES_USER` | List User Interviewer Pauses | read | `userId` |  | `organizationRead` |
| `ASHBY_LIST_INTERVIEWER_POOLS` | List Interviewer Pools | read | *cursor, perPage, syncToken* |  |  |
| `ASHBY_LIST_INTERVIEWS` | List Interview Types | read | *cursor, perPage, syncToken* |  |  |
| `ASHBY_LIST_INTERVIEW_EVENTS` | List Interview Events | read | `interviewScheduleId`; *cursor, limit, syncToken* |  | `interviewsRead` |
| `ASHBY_LIST_INTERVIEW_PLANS` | List Interview Plans | read | *cursor, includeArchived, perPage, syncToken* |  |  |
| `ASHBY_LIST_INTERVIEW_SCHEDULES` | List Interview Schedules | read | *cursor, perPage, syncToken* |  |  |
| `ASHBY_LIST_INTERVIEW_STAGES` | List Interview Stages | read | `interviewPlanId` |  |  |
| `ASHBY_LIST_INTERVIEW_STAGE_GROUPS` | List Interview Stage Groups | read |  |  |  |
| `ASHBY_LIST_JOBS` | List Jobs | read | *cursor, perPage, syncToken* |  |  |
| `ASHBY_LIST_JOB_BOARDS` | List Job Boards | read |  |  |  |
| `ASHBY_LIST_JOB_POSTINGS` | List Job Postings | read | *cursor, perPage, syncToken* |  |  |
| `ASHBY_LIST_JOB_TEMPLATES` | List Job Templates | read | *createdAfter, cursor, limit, syncToken* |  |  |
| `ASHBY_LIST_LOCATIONS` | List Locations | read | *cursor, limit* |  |  |
| `ASHBY_LIST_OFFERS` | List Offers | read | *createdAfter, cursor, perPage, syncToken* |  |  |
| `ASHBY_LIST_OPENINGS` | List Openings | read | *cursor, perPage, syncToken* |  |  |
| `ASHBY_LIST_PROJECTS` | List Projects | read | *cursor, perPage, syncToken* |  |  |
| `ASHBY_LIST_REFERRAL` | List Referrals | read | *applicationId, createdAfter, cursor, includeVisibilityConditions, limit, syncToken* |  | `candidatesRead` |
| `ASHBY_LIST_SEQUENCE` | List Sequence | read | *candidateId, cursor, limit, syncToken* |  | `sourcingRead` |
| `ASHBY_LIST_SEQUENCE_TEMPLATE` | List Sequence Template | read | *cursor, includeArchived, limit, syncToken* |  | `sourcingRead` |
| `ASHBY_LIST_SOURCES` | List Sources | read | *cursor, syncToken* |  |  |
| `ASHBY_LIST_SOURCE_TRACKING_LINKS` | List Source Tracking Links | read | *cursor, syncToken* |  |  |
| `ASHBY_LIST_SURVEY_FORM_DEFINITIONS` | List Survey Form Definitions | read | *cursor, limit* |  |  |
| `ASHBY_LIST_SURVEY_REQUEST` | List Survey Requests | read | `surveyType`; *applicationId, candidateId, createdAfter, cursor, limit, syncToken* |  | `candidatesRead` |
| `ASHBY_LIST_SURVEY_SUBMISSIONS` | List Survey Submissions | read | `surveyType`; *createdAfter, cursor, limit, syncToken* |  |  |
| `ASHBY_LIST_TAKE_HOME_ASSIGNMENT` | List Take Home Assignments | read | *applicationId, candidateId, cursor, expand, limit, syncToken* |  | `interviewsRead` |
| `ASHBY_LIST_USERS` | List Users | read | *cursor, perPage, syncToken* |  |  |
| `ASHBY_LIST_VERSION` | List Version | read |  |  |  |
| `ASHBY_MOVE_DEPARTMENT` | Move Department | write | `departmentId`; *parentId* |  | `organizationWrite` |
| `ASHBY_MOVE_LOCATION` | Move Location | write | `locationId`, `parentLocationHierarchyId` |  | `organizationWrite` |
| `ASHBY_REJECT_OFFER` | Reject Offer Approval | write | `approvalStepId`, `offerVersionId`, `userId`; *comment, excludeFormDefinition* |  | `offersWrite` |
| `ASHBY_REMOVE_HIRING_TEAM_MEMBER` | Remove Hiring Team Member | write | `roleId`, `teamMemberId`; *applicationId, jobId, openingId* |  | `jobsWrite` |
| `ASHBY_REMOVE_HIRING_TEAM_MEMBER_APPLICATION` | Remove Application Hiring Team Member | write | `applicationId`, `roleId`, `teamMemberId` |  | `candidatesWrite` |
| `ASHBY_REMOVE_INTERVIEWER_POOL_USER` | Remove User from Interviewer Pool | write | `interviewerPoolId`, `userId` |  | `hiringProcessMetadataWrite` |
| `ASHBY_REMOVE_OPENING_JOB` | Remove Opening Job | write | `jobId`, `openingId` |  | `jobsWrite` |
| `ASHBY_REMOVE_OPENING_LOCATION` | Remove Opening Location | write | `locationId`, `openingId` |  | `jobsWrite` |
| `ASHBY_REMOVE_PROJECT_CANDIDATE` | Remove Candidate from Project | write | `candidateId`, `projectId` |  | `candidatesWrite` |
| `ASHBY_REMOVE_TAG_CANDIDATE` | Remove Tag Candidate | write | `candidateId`, `tagId` |  | `candidatesWrite` |
| `ASHBY_RESTORE_DEPARTMENT` | Restore Department | write | `departmentId` |  | `organizationWrite` |
| `ASHBY_RESTORE_INTERVIEWER_POOL` | Restore Interviewer Pool | write | `interviewerPoolId` |  | `hiringProcessMetadataWrite` |
| `ASHBY_RESTORE_LOCATION` | Restore Location | write | `locationId` |  | `organizationWrite` |
| `ASHBY_RESTORE_PROJECT` | Restore Project | write | `projectId` |  | `projectsWrite` |
| `ASHBY_SEARCH_CANDIDATES` | Search Candidates | read | *email, name* |  |  |
| `ASHBY_SEARCH_JOBS` | Search Jobs | read | `title` |  |  |
| `ASHBY_SEARCH_OPENING` | Search Opening | read | `identifier` |  |  |
| `ASHBY_SEARCH_PROJECTS` | Search Projects | read | `title` |  |  |
| `ASHBY_SEARCH_SUGGESTIONS_CANDIDATE` | Search Candidate Suggestions | read | *email, limit, name* |  | `candidatesRead` |
| `ASHBY_SEARCH_USERS` | Search Users | read | `email` |  |  |
| `ASHBY_SET_CUSTOM_FIELD_VALUE` | Set Custom Field Value | write | `fieldId`, `fieldValue`, `objectId`, `objectType` |  | `candidatesWrite` |
| `ASHBY_SET_CUSTOM_FIELD_VALUES` | Set Custom Field Values | write | `objectId`, `objectType`, `values` |  | `candidatesWrite` |
| `ASHBY_SET_CUSTOM_FIELD_VALUES_PROJECT` | Set Custom Field Values Project | write | `projectId`, `values` |  | `projectsWrite` |
| `ASHBY_SET_CUSTOM_FIELD_VALUES_USER` | Set Custom Field Values User | write | `userId`, `values` |  | `organizationWrite` |
| `ASHBY_SET_CUSTOM_FIELD_VALUE_PROJECT` | Set Project Custom Field Value | write | `fieldValue`, `projectId`; *fieldId, referenceIdentifier* |  | `projectsWrite` |
| `ASHBY_SET_CUSTOM_FIELD_VALUE_USER` | Set User Custom Field Value | write | `fieldValue`, `userId`; *fieldId, referenceIdentifier* |  | `organizationWrite` |
| `ASHBY_SET_DECIDED_AT_OFFER` | Set Offer Decision Timestamp | write | `decidedAt`, `offerId`; *excludeFormDefinition* |  | `offersWrite` |
| `ASHBY_SET_FRAUD_STATUS_CANDIDATE` | Set Candidate Fraud Status | write | `candidateId`, `fraudStatus` |  | `candidatesWrite` |
| `ASHBY_SET_JOB_STATUS` | Set Job Status | write | `jobId`, `status` |  |  |
| `ASHBY_SET_OPENING_ARCHIVED` | Set Opening Archived | write | `archive`, `openingId` |  | `jobsWrite` |
| `ASHBY_SET_OPENING_OPENING_STATE` | Set Opening State | write | `openingId`, `openingState`; *closeReasonId* |  |  |
| `ASHBY_SET_STATUS_JOB_POSTING` | Set Status Job Posting | write | `jobPostingId`, `status` |  | `jobsWrite` |
| `ASHBY_SET_STATUS_OFFER` | Set Offer Acceptance Status | write | `acceptanceStatus`, `offerId`; *excludeFormDefinition* |  | `offersWrite` |
| `ASHBY_START_APPROVAL_PROCESS_JOB` | Start Job Approval Process | write | `jobId`; *note* |  | `approvalsWrite` |
| `ASHBY_START_APPROVAL_PROCESS_OFFER` | Start Offer Approval Process | write | `offerVersionId`; *excludeFormDefinition, note* |  | `offersWrite` |
| `ASHBY_START_OFFER` | Start Offer | write | `offerProcessId` |  |  |
| `ASHBY_START_OFFER_PROCESS` | Start Offer Process | write | `applicationId` |  |  |
| `ASHBY_START_SEQUENCE` | Start Sequence Draft | write | `sequenceId` |  | `sourcingWrite` |
| `ASHBY_SUBMIT_APPLICATION_FEEDBACK` | Submit Application Feedback | write | `applicationId`, `fieldSubmissions`, `formDefinitionId`, `interviewEventId`, `userId` |  | `candidatesWrite` |
| `ASHBY_SUBMIT_APPLICATION_FORM` | Submit Application Form | write | `applicationForm`, `jobPostingId`; *allowSubmissionForUnpublishedJobPosting, tagIds, utmData* |  | `candidatesWrite` |
| `ASHBY_SYNCHRONOUS_REPORT` | Synchronous Report | read | `reportId`; *includeHeadersInData, resultStyle* |  | `reportsRead` |
| `ASHBY_TRANSFER_APPLICATION` | Transfer Application | write | `applicationId`, `interviewPlanId`, `interviewStageId`, `jobId`; *startAutomaticActivities* |  |  |
| `ASHBY_UPDATE_APPLICATION` | Update Application | write | `applicationId`; *createdAt, creditedToUserId, sendNotifications, sourceId* |  |  |
| `ASHBY_UPDATE_APPLICATION_HISTORY` | Update Application History | write | `applicationHistory`, `applicationId` |  | `candidatesWrite` |
| `ASHBY_UPDATE_APPROVAL_DEFINITION` | Update Approval Definition | write | `approvalStepDefinitions`, `entityId`, `entityType`; *submitApprovalRequest* |  | `approvalsWrite` |
| `ASHBY_UPDATE_ASSESSMENT` | Update Assessment | write | `assessment_id`, `timestamp`; *assessment_profile_url, assessment_result, assessment_status, cancelled_reason, metadata* |  | `candidatesWrite` |
| `ASHBY_UPDATE_CANDIDATE` | Update Candidate | write | `candidateId`; *alternateEmail, company, creditedToUserId, email, githubUrl, linkedInUrl, name, phoneNumber and 5 more* |  |  |
| `ASHBY_UPDATE_COMPENSATION_JOB` | Update Job Compensation | write | `compensationTiers`, `jobId` |  |  |
| `ASHBY_UPDATE_COMPENSATION_JOB_POSTING` | Update Job Posting Compensation | write | `compensationTiers`, `jobPostingId`; *includeUnpublishedJobPostings* |  | `jobsWrite` |
| `ASHBY_UPDATE_COMPENSATION_SETTINGS_JOB` | Update Job Compensation Settings | write | `jobId`; *compensationPhilosophy, compensationTierGuideUrl, shouldShowOnJobBoard, shouldShowOnJobPostings* |  | `jobsWrite` |
| `ASHBY_UPDATE_COMPENSATION_SETTINGS_JOB_POSTING` | Update Job Posting Compensation Settings | write | `jobPostingId`; *compensationPhilosophy, compensationTierGuideUrl, includeUnpublishedJobPostings, shouldShowOnJobBoard, shouldShowOnJobPostings* |  | `jobsWrite` |
| `ASHBY_UPDATE_DEPARTMENT` | Update Department | write | `departmentId`; *name* |  |  |
| `ASHBY_UPDATE_INTERVIEWER_POOL` | Update Interviewer Pool | write | `interviewerPoolId`; *requiresTraining, title* |  | `hiringProcessMetadataWrite` |
| `ASHBY_UPDATE_INTERVIEW_SCHEDULE` | Update Interview Schedule | write | `interviewScheduleId`; *allowFeedbackDeletion, interviewEvent, interviewEventIdToCancel* |  | `interviewsWrite` |
| `ASHBY_UPDATE_JOB` | Update Job | write | `jobId`; *confidential, customRequisitionId, defaultInterviewPlanId, locationId, teamId, title* |  |  |
| `ASHBY_UPDATE_JOB_POSTING` | Update Job Posting | write | `jobPostingId`; *description, isListed, suppressDescriptionClosing, suppressDescriptionOpening, title* |  |  |
| `ASHBY_UPDATE_LOCATION_ADDRESS` | Update Location Address | write | `address`, `locationId` |  |  |
| `ASHBY_UPDATE_LOCATION_EXTERNAL_NAME` | Update Location External Name | write | `locationId`; *externalName* |  | `organizationWrite` |
| `ASHBY_UPDATE_LOCATION_NAME` | Update Location Name | write | `locationId`, `name` |  |  |
| `ASHBY_UPDATE_LOCATION_REMOTE_STATUS` | Update Location Remote Status | write | `isRemote`, `locationId` |  |  |
| `ASHBY_UPDATE_LOCATION_WORKPLACE_TYPE` | Update Location Workplace Type | write | `locationId`, `workplaceType` |  |  |
| `ASHBY_UPDATE_OFFER` | Update Offer | write | `offerForm`, `offerId`; *excludeFormDefinition* |  | `offersWrite` |
| `ASHBY_UPDATE_OPENING` | Update Opening | write | `openingId`; *description, employmentType, identifier, isBackfill, targetHireDate, targetStartDate, teamId* |  | `jobsWrite` |
| `ASHBY_UPDATE_PROJECT` | Update Project | write | `projectId`; *description, title* |  | `projectsWrite` |
| `ASHBY_UPDATE_SELECTABLE_VALUES_CUSTOM_FIELD` | Update Selectable Values Custom Field | write | `customFieldId`, `selectableValues` |  | `hiringProcessMetadataWrite` |
| `ASHBY_UPDATE_STAGE_SEQUENCE` | Update Sequence Stage | write | `sequenceId`, `stageId`; *body, subject* |  | `sourcingWrite` |
| `ASHBY_UPDATE_USER_INTERVIEWER_SETTINGS` | Update User Interviewer Settings | write | `userId`; *dailyLimit, weeklyLimit* |  |  |
| `ASHBY_UPLOAD_FILE_CANDIDATE` | Upload File Candidate | write | `candidateId`, `fileHandle` |  | `candidatesWrite` |
| `ASHBY_UPLOAD_RESUME_CANDIDATE` | Upload Resume Candidate | write | `candidateId`, `resumeHandle` |  | `candidatesWrite` |
