# Ashby

Checked: 2026-11-14 · Sources: the connection service's Ashby tool
definitions as planque reads them (191 tools), and Ashby's API documentation
(authentication, responses, the candidate and application references).

## How it connects

- **By the owner's API key**, pasted once and held by the connection
  service; the taf never sees it. In Ashby, an admin creates the key and
  picks, module by module, what it may read or write (candidates, jobs,
  interviews, offers and so on).
- **A key without a module leaves its tools out.** With a known sign-in the
  catalog puts the tools the key can't run under `unusable`. Ask for every
  tool the job needs in the first request, so nothing is asked twice.
- **Confidential jobs and private fields are hidden by default.** A key sees
  them only if the admin ticked the matching box when setting it up.

## Tools a job usually needs

Read means it changes nothing in Ashby. Returns are what a taf keeps to use
next time.

| Tool | What it does | Read or write | Returns |
|---|---|---|---|
| `ASHBY_LIST_JOBS` / `ASHBY_SEARCH_JOBS` | Lists the jobs, or finds one by title | read | job ids |
| `ASHBY_GET_JOB_INFO` | Reads one job | read | the job id |
| `ASHBY_LIST_INTERVIEW_STAGES` | Lists a job plan's stages | read | stage ids |
| `ASHBY_LIST_APPLICATIONS` | Lists applications, a page at a time | read | application ids, candidate ids, a page cursor |
| `ASHBY_GET_APPLICATION_INFO` | Reads one application: job, stage, status | read | the application id, its stage |
| `ASHBY_SEARCH_CANDIDATES` | Finds candidates by email or name | read | candidate ids |
| `ASHBY_GET_CANDIDATE_INFO` | Reads one candidate's profile | read | the candidate id |
| `ASHBY_LIST_CANDIDATE_NOTES` | Reads the notes on a candidate | read | note ids |
| `ASHBY_LIST_APPLICATION_FEEDBACK` | Reads interviewers' feedback on an application | read | feedback ids |
| `ASHBY_LIST_INTERVIEW_SCHEDULES` | Lists scheduled interviews | read | schedule ids |
| `ASHBY_CREATE_CANDIDATE_NOTE` | Adds a note to a candidate | write | the note id |
| `ASHBY_CREATE_CANDIDATE` | Adds a candidate | write | the candidate id |
| `ASHBY_CREATE_APPLICATION` | Puts a candidate forward for a job | write | the application id |
| `ASHBY_CHANGE_APPLICATION_STAGE` | Moves an application to another stage | write | the application, with its new stage |

**What to keep:** each candidate's id and application id, to read or update
the same person next time; job and stage ids, which every move needs.
Before adding a candidate, search by email first, so the same person isn't
added twice.

**Start read-only.** "Summarise new applicants each morning" is the read
tools alone. A note is the lightest write. Offer a stage move only when the
owner asked for it, and say each move waits for their yes: moving a
candidate is what a candidate and the hiring team see.

## Don't offer

- Deleting and anonymising tools (applications, referrals, candidates):
  they can't be undone, and anonymising erases a person's data.
- Offer tools (create, approve, reject): an offer is a decision for people.
- Rejection through a stage move with an archive reason, unless the owner
  asked for exactly that, candidate by candidate.
- Hiring-process setup (custom fields, interviewer pools, job and opening
  setup, hiring-team members): they change how the whole team works.
- Email and sourcing-sequence tools: they write to candidates from the
  team's addresses.
- Notetaker transcripts and private fields, unless the job needs them and
  the owner says so.

## Access and plans

- **Each module is read or write on the key.** A key with candidate read
  only lets a taf read candidates and applications, not add notes or move
  stages.
- **Acting as a named person** in Ashby needs a separate permission on the
  key.
- **Rate:** 1,000 calls a minute for each key. A taf pages through lists
  rather than reading everything each run.

## When it fails

Ashby answers most refusals as a normal answer marked unsuccessful, with an
error code inside, rather than an error status.

| What happens | Why | Tell the owner |
|---|---|---|
| Refused for missing access (403, missing endpoint permission) | The key's module for this step is off, or read only | "Ashby is only partly connected: this step needs access your key doesn't give. Ask your Ashby admin to allow it on the key, or do without this step." |
| The key doesn't work (401 or 403) | No key was sent, or the key is wrong or switched off | "Ashby needs a new key from your admin. Nothing was changed." |
| An application or candidate isn't found | It was removed, or it belongs to a confidential job the key can't see | "The taf can't see that candidate. It may be deleted, or on a confidential job. Your admin can allow confidential jobs on the key." |
| Answer marked unsuccessful, with a message | A field is wrong, such as a stage from another job's plan | "Ashby didn't accept this change. The taf will stop and show you what it tried." |

After any refusal the taf stops that step and says so. A refused read never
leads to a write.
