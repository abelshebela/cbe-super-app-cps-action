# Test Setup

## Folder layout

- `test/testinit/`
  - Shared helpers for tests (logger, request context helpers)
- `test/unit/`
  - `test/unit/handler/` - HTTP handler unit tests
  - `test/unit/service/` - service unit tests
  - `test/unit/storage/` - storage/repository unit tests (no real DB)
- `test/integration/`
  - Integration tests (real DB/redis/etc). Keep these isolated from unit tests.
- `test/e2e/`
  - End-to-end tests (start router/server and call HTTP endpoints)

## Running tests

- Unit tests only:
  - `go test ./test/unit/...`
- All tests under `./test`:
  - `go test ./test/...`

## Conventions

- Put domain-specific tests under a domain folder, for example:
  - `test/unit/handler/bps_action/...`
  - `test/unit/service/roles/...`
  - `test/unit/storage/bps_action_role/...`
- Use `test/testinit` helpers to build request context values (`user_id`, `role_code`, etc.) instead of duplicating that logic.
