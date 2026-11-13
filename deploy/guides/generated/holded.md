# Holded: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **A sign-in, held by the connection service** (the owner pastes a key from the vendor's settings). The taf never sees a password or token.
- With a known sign-in, the grant's catalog lists under `unusable`, with why, each tool below its access doesn't cover. Offer only its `tools`.
- An app's held sign-in is reused for a new taf when it covers the tools asked: ask for every tool the job needs in the first request.

## Tools

10 tools: 4 read, 6 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `HOLDED_CREATE_CONTACT` | Create Contact | write | `name`; *billing_address, contact_type, email, is_person, mobile, phone, tax_id, trade_name and 2 more* |  |  |
| `HOLDED_CREATE_INVOICE` | Create Invoice | write | `contact_id`, `items`; *body, currency, currency_change, description, discount, due_date, issue_date, language and 4 more* |  |  |
| `HOLDED_DELETE_CONTACT` | Delete Contact | write | `contact_id` |  |  |
| `HOLDED_DELETE_INVOICE` | Delete Invoice | write | `invoice_id` |  |  |
| `HOLDED_FIND_CONTACTS` | Find Contacts | read | *cursor, email, external_id, limit, mobile, name, phone, tax_id* |  |  |
| `HOLDED_FIND_INVOICES` | Find Invoices | read | *contact_id, cursor, document_number, due_date_end, due_date_start, issue_date_end, issue_date_start, limit and 2 more* |  |  |
| `HOLDED_GET_CONTACT` | Get Contact | read | `contact_id` |  |  |
| `HOLDED_GET_INVOICE` | Get Invoice | read | `invoice_id` |  |  |
| `HOLDED_UPDATE_CONTACT` | Update Contact | write | `changes`, `contact_id` |  |  |
| `HOLDED_UPDATE_INVOICE` | Update Invoice | write | `changes`, `invoice_id` |  |  |
