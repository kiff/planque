# Odoo: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner pastes a key from the vendor's settings). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

12 tools: 5 read, 7 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `ODOO_CALL_ODOO_JSONRPC` | Odoo JSON-RPC Call | write | `method`, `service`; *args, database, request_id* |  |  |
| `ODOO_CREATE_CONTACT` | Create Contact | write | `name`; *city, comment, company_type, country_id, customer_rank, email, function, is_company and 14 more* |  |  |
| `ODOO_GET_APPLICANT_RESULT` | Get Applicant Parsing Result | read | `account_token`, `document_token`; *version* |  |  |
| `ODOO_GET_BANK_STATEMENT_RESULT` | Get Bank Statement Result | read | `account_token`, `document_token`, `version` |  |  |
| `ODOO_GET_EXPENSE_RESULT` | Get Expense Extraction Result | read | `account_token`, `document_token`, `version` |  |  |
| `ODOO_GET_INVOICE_RESULT` | Get Invoice Parsing Result | read | `account_token`, `document_token`, `version` |  |  |
| `ODOO_LIST_DATABASES` | List Databases | read |  |  |  |
| `ODOO_PARSE_APPLICANT` | Parse Applicant Resume | write | `account_token`, `document`; *request_id, version* |  |  |
| `ODOO_PARSE_BANK_STATEMENT` | Parse Bank Statement Document | write | `account_token`, `documents`, `version`; *request_id* |  |  |
| `ODOO_PARSE_EXPENSE` | Parse Expense Document | write | `account_token`, `documents`, `version` |  |  |
| `ODOO_PARSE_INVOICE` | Parse Invoice Document | write | `account_token`, `documents`, `version`; *dbuuid, request_id, user_infos, webhook_url* |  |  |
| `ODOO_UPDATE_RECORD` | Update Odoo Record | write | `model`, `record_ids`, `values` |  |  |
