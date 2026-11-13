# Wise: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner pastes a key from the vendor's settings). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

9 tools: 8 read, 1 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `WISE_CREATE_QUOTE` | Create Quote | write | `profile_id`, `source_currency`, `target_currency`; *preferred_pay_in, source_amount, target_amount* |  |  |
| `WISE_GET_EXCHANGE_RATE` | Get Exchange Rate | read | `source_currency`, `target_currency` |  |  |
| `WISE_GET_QUOTE` | Get Quote | read | `profile_id`, `quote_id` |  |  |
| `WISE_LIST_ACTIVITIES` | List Activities | read | `profile_id`; *next_cursor, page_size, resource_type, since, status, until* |  |  |
| `WISE_LIST_BALANCES` | List Balances | read | `profile_id`; *balance_types* |  |  |
| `WISE_LIST_CURRENCIES` | List Supported Currencies | read | *locale* |  |  |
| `WISE_LIST_PROFILES` | List Profiles | read |  |  |  |
| `WISE_LIST_RECIPIENTS` | List Recipients | read | `profile_id`; *active, currencies, next_cursor, owned_by_customer, page_size, recipient_types* |  |  |
| `WISE_LIST_TRANSFERS` | List Transfers | read | `profile_id`; *created_at_or_after, created_at_or_before, next_cursor, page_size, source_currency, statuses, target_currency* |  |  |
