# NEXT.md: restart notes for Sat Nov 28 (Day 32)

Last session: Sat Nov 21 (Day 31). Mid-exam week Nov 23-27 is maintenance only.

## Where I stopped

- Task 9 lives in `Tasks/Task_9/task_manager`, using local PostgreSQL 18.7 (no Docker).
- Done: migrations, Postgres user and task repositories, `DB_DRIVER` switch for both, pagination and status filter with `X-` headers, ownership in the usecase layer (404 for other users' tasks), updated mocks, tests and API docs, `docs/performance.md` (index and OFFSET vs keyset). All tests pass and the manual curl check passed for a normal user, a second user and an admin.
- Not done: `docs/schema.md`, `POST /tasks/reassign`, the repository integration tests, running the Postman collection with `DB_DRIVER=mongo`.

## Next three steps

1. **Re-sync (30 min).**
   - Run `go build ./...` and `go test ./... -cover`.
   - Run the Postman collection with `DB_DRIVER=postgres` (and `mongo` if time allows).
   - Fix anything that rotted during exam week.
2. **Day 32: SKIP LOCKED test.**
   - Run two scheduler processes (two terminals, different `PORT`) against the same Postgres database.
   - Create 2-3 tasks due within `REMINDER_WINDOW`. Each reminder must be sent exactly once across both processes.
3. **Day 33 (Mon Nov 30): reassign and integration tests.**
   - Add `POST /tasks/reassign` (admin only), `{from_user_id, to_user_id}`, in ONE transaction: check both users exist, move all tasks, return the count, roll back on any failure.
   - Write repository integration tests against `TEST_DATABASE_URL` (migrate first, truncate between tests): reassign success plus a rollback case, ClaimDueTasks, the unique-username 409.
   - Write `docs/schema.md`, update the API doc and Postman, tag `v0.9`.

## Run it

```powershell
cd Tasks/Task_9/task_manager
psql -U postgres -h localhost -d taskmanager   # database: taskmanager, user: postgres
migrate -path migrations -database $env:DATABASE_URL up
$env:DB_DRIVER="postgres"; go run ./Delivery
```

- Env vars: `DB_DRIVER`, `DATABASE_URL`, `JWT_SECRET`, `PORT`, `REMINDER_*`. `MONGO_URI` is only needed when `DB_DRIVER=mongo`.
- `DATABASE_URL` lives in `.env`, so PowerShell can't see it. Load it first (see below) or pass `-U postgres -h localhost -d taskmanager` to `psql`.
- `go test -race` doesn't work locally on Windows (needs cgo). Use the Docker command from Day 23 or rely on GitHub Actions.

## Gotchas

- The auth middleware keys in `controller.go` (`ctxUserIDKey`, `ctxRoleKey`) must match the `c.Set(...)` calls in `Infrastructure/auth_middleWare.go`.
- Tasks need a valid Postgres user as owner, so test with users registered on Postgres.
- Small leftovers, if not done yet: in `scanTask` add `t.DueDate = t.DueDate.UTC()` (responses currently show `+03:00`), and change the 404 text in `api_documentation.md` to `not found`, which is what the API returns.
- Seed users (`seed_user_%`) from the benchmarks must stay deleted, otherwise lists and reminder tests get cluttered.
