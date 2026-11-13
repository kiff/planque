# Qonto: generated

Checked: 2026-11-13 · Generated from planque's catalog: don't edit by hand, run `planque guides` again.

## How it connects

- **By the owner's API key**, reviewed in planque's catalog: the tools below are the ones installed, as reviewed. The key and the taf's Card come in one approval.

## Tools

8 tools: 7 read, 1 write (a tool not marked read-only counts as a write).

| Tool | What it does | Read or write | Inputs (required; *optional*) | Returns | Access, any one of |
|---|---|---|---|---|---|
| `qonto_organization` | Read the business and its bank accounts | read |  |  |  |
| `qonto_list_clients` | List the clients | read |  |  |  |
| `qonto_list_client_invoices` | List the client invoices | read |  |  |  |
| `qonto_get_client_invoice` | Read one client invoice | read |  |  |  |
| `qonto_list_quotes` | List the quotes | read |  |  |  |
| `qonto_get_quote` | Read one quote | read |  |  |  |
| `qonto_list_supplier_invoices` | List the supplier invoices | read |  |  |  |
| `qonto_draft_client_invoice` | Draft an invoice from an accepted quote | write | `client_id`, `currency`, `due_date`, `iban`, `issue_date`, `items`, `performance_end_date`, `performance_start_date`, `quote_id` |  |  |
