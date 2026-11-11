# go_backend_journey

Repo: https://github.com/michael-gt-12/go_backend_journey

My public learning log for a **65-day Go backend + AI backend track**, running from **Sun Oct 18, 2026 to Sun Jan 3, 2027**. I'm Michael, a 4th-year Software Engineering student at Adama Science & Technology University (ASTU), and I'm building this in public so every claim is backed by code, tests and commits.

**Goal:** go from Go basics to deployed, tested backend services, including four AI-powered backend projects, and finish with a portfolio I can point recruiters to.

## How I work

- About 3 focused hours a day: 30-45 min learning, about 2 hrs building, 15 min to commit, update this README and write the LinkedIn post.
- Every day has a "Done when" line that acts as my acceptance test. I push only when it is true.
- Sundays are buffer days for slipped work, the weekly check-in and the weekly recap post.
- Mid-exam week (Nov 23-27) is maintenance only. The plan resumes Sat Nov 28.
- Post format: `Day N/65 - what I built - what I learned - what is next - repo link`.

## The plan

| Block                                     | Days  | Dates           | Output                                                              |
| ----------------------------------------- | ----- | --------------- | ------------------------------------------------------------------- |
| Go foundations + Tasks 1-2                | 1-4   | Oct 18 - Oct 21 | Tour of Go, functions with tests, console library system            |
| Gin REST API + Postman (Task 3)           | 5-8   | Oct 22 - Oct 26 | Task Manager API, documented in Postman                             |
| MongoDB persistence (Task 4)              | 9-12  | Oct 27 - Oct 30 | API backed by MongoDB with the Go driver                            |
| JWT auth and roles (Task 5)               | 13-15 | Oct 31 - Nov 3  | Register/login, middleware, role-based access                       |
| Clean Architecture refactor (Task 6)      | 16-19 | Nov 4 - Nov 7   | Delivery / Domain / Infrastructure / Repositories / Usecases        |
| Testing with testify and mockery (Task 7) | 20-22 | Nov 9 - Nov 11  | Usecase and controller tests with mocks and suites                  |
| Concurrency + reminders (Task 8)          | 23-27 | Nov 12 - Nov 17 | Worker pool, scheduler, graceful shutdown, race-free tests          |
| PostgreSQL migration (Task 9)             | 28-33 | Nov 18 - Nov 30 | Migrations, pgx repositories, pagination, SKIP LOCKED, transactions |
| Docker + Compose (Task 10)                | 34-36 | Dec 1 - Dec 3   | Multi-stage image, one-command stack with healthchecks              |
| Redis (Task 11)                           | 37-39 | Dec 4 - Dec 7   | Cache-aside, Redis rate limiter, refresh-token rotation             |
| CI + Render deployment (Task 12)          | 40-41 | Dec 8 - Dec 9   | Green CI, live API on Render with smoke test                        |
| AI Project 1: SupportDesk AI              | 42-46 | Dec 10 - Dec 15 | Webhooks, queue, structured LLM triage, evals                       |
| AI Project 2: DocuAsk (RAG)               | 47-52 | Dec 16 - Dec 22 | Ingestion, pgvector, cited answers, streaming, eval                 |
| AI Project 3: LLM Gateway                 | 53-57 | Dec 23 - Dec 28 | Quotas, fallbacks, caching, metrics                                 |
| AI Project 4: OpsAgent                    | 58-62 | Dec 29 - Jan 2  | Tool-calling agent with guardrails, approvals, audit log            |
| Final buffer                              | -     | Sun Jan 3       | Finish anything unfinished, review, plan the next phase             |

## The four AI backend projects

1. **SupportDesk AI:** receives support tickets through signed webhooks, triages them with an LLM (category, priority, summary, suggested reply) and lets an agent review the result.
2. **DocuAsk:** upload documents and ask questions answered only from them, with citations (RAG over PDF/DOCX using Postgres + pgvector).
3. **LLM Gateway:** one internal API in front of LLM providers, with API keys, per-tenant quotas, retries, fallbacks, caching and Prometheus metrics.
4. **OpsAgent:** a tool-calling agent that looks things up and proposes actions safely, with prompt-injection defences, human approval and an audit log.

## Repository layout

```
go_backend_journey/
├── go.mod              # one module for the daily exercises
├── README.md
├── Days/               # one folder per day: Day_1, Day_2, ...
├── Tasks/              # Task_1, Task_2, then task_manager (Tasks 3-12 grow in one codebase)
└── Projects/           # the four AI backend projects, one folder each
```

Folders under `Tasks/` and `Projects/` that have their own `go.mod` are separate Go modules. Later, the Task Manager API and each AI project may move into their own repos for CI and deployment.

## Tech

Go, Gin, MongoDB, PostgreSQL (pgx, pgvector), Redis, Docker and Compose, GitHub Actions, Render, Prometheus, testify and mockery.

## Progress

| Day    | Date         | Topic                                                                                              | Code                                                           |
| ------ | ------------ | -------------------------------------------------------------------------------------------------- | -------------------------------------------------------------- |
| 1      | Oct 18       | Environment setup, Go basics, hello-world, FizzBuzz CLI                                            | `Days/Day_1`                                                   |
| 2      | Oct 19       | Structs, methods, interfaces, custom errors, generics                                              | `Days/Day_2`                                                   |
| 3      | Oct 20, 2026 | Task 1 (word frequency, palindromes) + Task 2 skeleton                                             | `Tasks/Task_1`, `Tasks/Task_2/library_management`              |
| 4      | Oct 21, 2026 | Task 2: Console Library Management System (interfaces, maps, errors)                               | `Tasks/Task_2/library_management`                              |
| 5      | Oct 22, 2026 | Gin REST API basics: routing, JSON, handlers, request binding                                      | [Days/Day_5/gin-albums](Days/Day_5/gin-albums)                 |
| 6      | Oct 23       | Task 3: Gin project structure, in-memory store with mutex, task endpoints                          | [Tasks/Task_3/task_manager_api](Tasks/Task_3/task_manager_api) |
| 7      | Oct 24, 2026 | Task 3: POST/PUT/DELETE endpoints, validation, consistent errors                                   | Tasks/Task_3/task_manager                                      |
| 8      | Oct 26       | Postman API documentation, Task 3 finished                                                         | Tasks/Task_3/                                                  |
| 9      | Oct 27       | MongoDB fundamentals (documents, collections, ObjectID, CRUD, filters)                             | [Days/Day_9](Days/Day_9)                                       |
| 10     | Oct 28, 2026 | MongoDB Go driver: connect, context, CRUD                                                          | [Days/Day_10](Days/Day_10)                                     |
| 11     | Oct 29       | MongoDB in the Task Manager API: connect, GetAll, GetByID, Create                                  | Tasks/Task_4                                                   |
| 12     | Oct 30, 2026 | Task 4 finished: update/delete with MongoDB, error handling, config                                | Tasks/Task_4                                                   |
| 13     | Oct 31       | Authentication vs authorization, bcrypt password hashing, JWT structure                            | Days/Day_13                                                    |
| 14     | Nov 2        | Registration & login (Task 5, part 1): bcrypt + JWT                                                | Tasks/Task_5                                                   |
| 15     | Nov 3, 2026  | Task 5: JWT middleware and role-based access control                                               | Days/Day_15                                                    |
| 16     | Nov 4, 2026  | Clean Architecture: layers, dependency rule, Task 5 → Task 6 mapping                               | Tasks/Task_6                                                   |
| 17     | Nov 5        | Clean Architecture: Domain and Usecases                                                            | [Tasks/Task_6/task_manager](Tasks/Task_6/task_manager)         |
| Day 18 | Nov 6        | Task 6 part 2: Clean Architecture, repositories, infrastructure and delivery layers with manual DI | [Tasks/Task_6/task_manager](Tasks/Task_6/task_manager)         |
| 19     | Sat Nov 7    | Task 6 wrap-up: regression check, architecture docs, service tests                                 | Tasks/Task_6/task_manager                                      |
| 20     | Nov 9        | Testify (assert/require, table-driven tests, suites, mocks) + generating mocks with mockery        | [Days/Day_20](Days/Day_20)                                     |
| 21     | Nov 10       | Unit-testing usecases with Testify suites and mocks                                                | Tasks/Task_7/task_manager                                      |
| 22     | Nov 11       | Controller tests with mocked usecases + httptest (Task 7 complete)                                 | Tasks/Task_7/task_manager                                      |

## Day 1 (Oct 18)

- Installed Go and set up the editor
- Initialised this repo with `go mod init`
- Hello-world and a small FizzBuzz CLI
- Tour of Go: Basics

## Day 2 - interfaces, errors and generics

Built in `Days/Day_2`, each piece with table-driven tests:

- **`shape.go`**: a `Shape` interface (`Area`, `Perimeter`) implemented by `Rectangle` and `Circle`. Types satisfy interfaces implicitly, and a compile-time check (`var _ Shape = Rectangle{}`) guards it.
- **`errors.go`**: a custom `ValidationError` type and `ValidateAge`. Tests use `errors.As` to recover the concrete error type.
- **`generics.go`**: generic `Map[T, U]` and `Filter[T]`.
  What I learned: interfaces are satisfied implicitly, errors are just values with an `Error() string` method, and generics let one function work across types without losing type safety.

## Day 3: Task 1 and the Task 2 skeleton

**What I built**

- `WordFrequency`: case-insensitive word counts that ignore punctuation (`strings.FieldsFunc` plus `unicode` checks).
- `IsPalindrome`: ignores spaces, punctuation and case. It works on runes, so Amharic text is handled correctly.
- Table-driven tests (empty string, punctuation only, mixed case, digits, Amharic, accents).
- The Task 2 `library_management` skeleton: required folder structure, `Book` and `Member` models, and the 6-method `LibraryManager` interface.

**What I learned**

- A Go string is bytes, not characters. Ranging over it gives runes, which is what makes Unicode text safe.
- `strings.FieldsFunc` is more flexible than `strings.Fields` because I decide what a separator is.
- Designing the interface before the implementation shows what the service must do.

**How to run**

```bash
cd Tasks/Task_1 && go test -cover ./...
cd Tasks/Task_2/library_management && go run .
```

## Day 4 — Task 2: Library Management System

### What I built

A console-based library system. A `LibraryManager` interface defines six
operations, implemented by a `Library` service that stores books and members
in maps. A menu loop (built with `bufio`) lets the user add and remove books,
borrow and return them, list available books, and list the books a member has
borrowed.

### What I learned

- Go interfaces are implemented implicitly; `var _ LibraryManager = (*Library)(nil)`
  makes the compiler verify it.
- Maps hold copies of structs, so updated values must be written back.
- Errors are values: sentinel errors give clear, testable failure cases
  (book not found, member not found, already borrowed, not borrowed by member).
- Validate console input before using it.

### How to run

```bash
cd Tasks/Task_2/library_management
go run main.go
```

### How to test

```bash
go test ./...
```

See `docs/documentation.md` for the folder structure and error cases.

## Day 5 – Gin REST API basics (Oct 22, 2026)

### What I built

An albums REST API with Go and Gin, in `Days/Day_5/gin-albums`, following the official Go + Gin tutorial and extending it myself:

| Method | Endpoint      | Description                        | Status codes |
| ------ | ------------- | ---------------------------------- | ------------ |
| GET    | `/albums`     | List all albums                    | 200          |
| GET    | `/albums/:id` | Get one album                      | 200, 404     |
| POST   | `/albums`     | Add an album from a JSON body      | 201, 400     |
| DELETE | `/albums/:id` | Delete an album (my own extension) | 200, 404     |

Missing albums return a JSON body: `{"message": "album not found"}`. Data is stored in an in-memory slice, so it resets when the server restarts.

### What I learned

- How Gin routes requests and reads URL parameters with `c.Param("id")`.
- How `BindJSON` turns a request body into a struct, and how struct tags (`json:"title"`) control JSON field names.
- Which status code fits which result: 200 OK, 201 Created, 404 Not Found.
- Deleting from a slice with `append(albums[:i], albums[i+1:]...)`. The `...` is required.
- On Windows PowerShell, `curl` is an alias for `Invoke-WebRequest`. Use `curl.exe` to get the real curl (`-X`, `-i`, `-d`).

### How to run

```bash
cd Days/Day_5/gin-albums
go mod tidy
go run .          # server starts on localhost:8000
```

Test it from a second terminal:

```powershell
curl.exe http://localhost:8000/albums
curl.exe http://localhost:8000/albums/2
curl.exe -i http://localhost:8000/albums/99
curl.exe -X DELETE http://localhost:8000/albums/1
curl.exe -i -X POST http://localhost:8000/albums -H "Content-Type: application/json" -d "@album.json"
```

## Day 6 - Oct 23: Task Manager API (Task 3, structure)

### What I built

- Set up the Task 3 project layout: `models/`, `data/`, `controllers/`, `router/`, `docs/` and a `main.go` that only starts the app
- A `Task` model with ID, Title, Description, DueDate and Status
- An in-memory store (a map protected by `sync.RWMutex`) with seeded sample tasks and server-generated IDs
- A data layer that holds the business logic and returns values or errors, and a controller layer that turns them into HTTP responses
- Endpoints: `GET /tasks`, `GET /tasks/:id`, plus `POST`, `PUT` and `DELETE` (started early), with a 404 and `{"error": "..."}` for missing tasks
- Tested every endpoint in Postman

### What I learned

1. Keep layers separate: the data layer should not know about Gin, it returns `(Task, error)` and the controller picks the status code.
2. Go maps store copies of structs, so a changed task must be written back into the map. Map order is also random, so I sort before returning a list.
3. A mutex protects shared data when requests arrive at the same time: `RLock` for reads, `Lock` for writes.
4. `go mod verify` can reveal a corrupted module cache, and `go clean -modcache` followed by `go mod tidy` fixes it.

### How to run

```bash
cd Tasks/Task_3/task_manager_api
go run .
```

The server starts on `http://localhost:8080`. Try `GET /tasks`, `GET /tasks/1` and `GET /tasks/999` (404).

## Day 7 - Task Manager API: Create, Update, Delete + Validation

### What I built

- `POST /tasks` creates a task and returns `201 Created`
- `PUT /tasks/:id` updates a task and returns `200 OK` with the updated task
- `DELETE /tasks/:id` deletes a task (returns 204 No Content)
- Request validation with Gin binding tags:
  - `title` is required
  - `due_date` must be a valid RFC 3339 date
  - `status` must be one of `pending`, `in_progress` or `done`
- One consistent error format everywhere: `{"error": "message"}`
- `400` for invalid input and `404` when the task does not exist
- Thread-safe in-memory storage (`map` + `sync.RWMutex`) with a sentinel `ErrTaskNotFound` error

### What I learned

1. `ShouldBindJSON` lets me control the error response, while `BindJSON` writes its own plain-text 400.
2. Validation rules live in struct tags (`binding:"required,oneof=..."`), so handlers stay short.
3. A shared `respondError` helper keeps every error response identical, which makes the API easier to document and test.
4. Write operations need `mu.Lock()`, while reads can use `RLock()`.

### How to run

cd Tasks/Task_3/task_manager
go run main.go

Create a task (on Windows PowerShell use `curl.exe`):

curl.exe -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d "{\"title\":\"Study Go\",\"description\":\"Gin\",\"due_date\":\"2026-11-01T10:00:00Z\",\"status\":\"pending\"}"

Update: curl.exe -X PUT http://localhost:8080/tasks/1 -H "Content-Type: application/json" -d "@body.json"
Delete: curl.exe -X DELETE http://localhost:8080/tasks/1

### Day 8 – Oct 26: Postman documentation (Task 3 complete)

**What I built**

- A Postman collection for all 5 Task Manager endpoints (GET all, GET by ID, POST, PUT, DELETE)
- Saved example responses for the success, 400 and 404 cases
- A description for every request and an overview for the collection
- Published public API docs: `https://documenter.getpostman.com/view/57227945/2sBYHUULJ5`
- `docs/api_documentation.md` with the payloads, responses, status codes and run instructions

**What I learned**

1. Good API docs let someone use the API without reading the code.
2. Saved examples in Postman double as documentation and as a regression test.
3. Keeping one error format (`{"error": "message"}`) makes the docs simpler and the API predictable.

**How to run**

```bash
cd Tasks/Task_3
go run main.go
```

The server runs on `http://localhost:8080`. See `docs/api_documentation.md` or the Postman link for all endpoints.

## Day 9 — MongoDB Fundamentals

### What I built

- Installed MongoDB Community Server and mongosh locally and connected through Compass
- Created the `taskmanager` database with a `tasks` collection holding 10 sample documents
- Ran 5 queries by hand: filter by status, sort by due date, count by status (aggregation), update many, delete one
- Wrote `mongo-notes.md` documenting each query, its result, and the ObjectID vs string ID difference

### What I learned

1. A document is one record (like one `Task`), and a collection is a group of documents
2. `_id` is an auto-generated `ObjectId`, not a string, so lookups by ID need a conversion
3. Filters choose which documents are affected, and `$set` changes only the fields I name

### How to run

1. Make sure the MongoDB server is running (the MongoDB Server service on Windows)
2. Open the shell in Compass (or run `mongosh`) and connect to `mongodb://localhost:27017`
3. Switch to the database with `use taskmanager`
4. Run the commands from `Days/Day_9/mongo-notes.md` one at a time

## Day 10 - MongoDB Go Driver (Oct 28, 2026)

**Folder:** `Days/Day_10/mongo-scratch`

### What I built

- A scratch Go program that connects to a local MongoDB and runs the full CRUD cycle on the `taskmanager.tasks` collection.
- `Connect` + `Ping` to confirm the server is reachable.
- `InsertOne`, `FindOne`, `Find` (with a filter), `UpdateOne` (with `$set`) and `DeleteOne`.
- A `Task` struct with `bson` tags, and a helper that gives every database call its own 5-second `context.WithTimeout`.
- The connection string is read from the `MONGO_URI` env var. `.env` is git-ignored and `.env.example` is committed.

### What I learned

1. `Connect` alone doesn't prove the database is reachable, so `Ping` is needed to fail fast.
2. `Find` returns a cursor, not a slice, so it must be read with `cursor.All` and closed with `defer cursor.Close(ctx)`.
3. `MatchedCount` and `DeletedCount` tell me whether anything actually changed, which is how the API will return 404 in Task 4.
4. `_id` is an `ObjectID`, not a string, so IDs from URLs need `primitive.ObjectIDFromHex`.

### How to run

```bash
cd Days/Day_10/mongo-scratch
cp .env.example .env        # PowerShell: Copy-Item .env.example .env
go mod tidy
go run .
```

## Day 11 (Oct 29): Task 4 part 1, MongoDB in the Task Manager API

### What I built

- Copied the Task 3 project into `Tasks/Task_4` and replaced the in-memory storage with MongoDB, keeping the same folder structure.
- `data.Connect`: connects with a timeout and runs `Ping`, so the app exits immediately if the database is unreachable (fail fast).
- Configuration from environment variables (`MONGO_URI`, `DB_NAME`, `PORT`) with `.env` git-ignored and a `.env.example` committed.
- Task model with `bson` tags and `primitive.ObjectID` as the `_id`.
- `GET /tasks` (always returns `[]`, never `null`), `GET /tasks/:id` (400 for an invalid ObjectID, 404 when not found) and `POST /tasks` (the database generates the ID).
- PUT and DELETE return 501 for now; they come on Day 12.

### What I learned

1. `mongo.Connect` doesn't prove the server is reachable; `Ping` does.
2. An invalid ObjectID (400) and a missing document (404) are different errors and need different responses.
3. Putting the ObjectID in the model keeps the API response shape unchanged, so Task 3 clients still work.

### How to run

```bash
cd Tasks/Task_4
cp .env.example .env      # set MONGO_URI, DB_NAME, PORT
go mod tidy
go run main.go
```

Then try:

```bash
curl.exe http://localhost:8000/tasks
curl.exe -X POST http://localhost:8000/tasks -H "Content-Type: application/json" -d "{\"title\":\"Test\",\"description\":\"d\",\"due_date\":\"2026-11-10T10:00:00Z\",\"status\":\"pending\"}"
```

## Day 12 - Oct 30, 2026: Finishing Task 4 (MongoDB persistence)

### What I built

- `UpdateTask` (PUT /tasks/:id) using `UpdateOne` with `$set`; returns the updated task.
- `DeleteTask` (DELETE /tasks/:id) using `DeleteOne`; returns 200 with `{"message": "task deleted"}`.
- Separate error handling for each failure type through sentinel errors (`ErrInvalidID`, `ErrNotFound`, `ErrEmptyUpdate`, `ErrTimeout`, `ErrDBUnavailable`) and one `respondError` helper in the controller:
  - invalid ObjectID or bad body: 400
  - task not found: 404
  - database timeout: 504
  - database unreachable: 503
  - anything else: 500 with a generic message (details only in server logs)
- Configuration from environment variables (`MONGO_URI`, `PORT`) with a `.env.example`.
- Updated `docs/api_documentation.md` with the MongoDB setup, connection settings, error table and a Task 3 to Task 4 change list.
- Extended the Postman collection with update/delete cases (success, 400, 404, DB down).

### What I learned

1. `MatchedCount` and `ModifiedCount` mean different things. Updating a task with identical values modifies nothing but is still a success, so the 404 check must use `MatchedCount == 0`. For deletes, use `DeletedCount == 0`.
2. Network errors, database errors and validation errors should be different responses. A DB timeout (504) is not the same as a bad request (400), and driver error details should never reach the client.
3. Configuration belongs in environment variables (12-factor). Real secrets are never committed; a `.env.example` documents what is needed.

### How to run

1. Start MongoDB (local `mongod` or an Atlas cluster).
2. Copy `.env.example` to `.env` and set `MONGO_URI` and `PORT`.
3. From `Tasks/Task_4`:
   ```bash
   go mod tidy
   go run main.go
   ```
4. Import the Postman collection, set `baseUrl`, and run the requests in order.

## Day 13 – Authentication, Password Hashing and JWT (Oct 31)

### What I built

- `auth-scratch/main.go`: a scratch program that
  - hashes and verifies passwords with bcrypt (and shows the same password gives different hashes),
  - creates an HS256 JWT with custom claims (user ID, username, role) and a 15-minute expiry,
  - parses and validates the token,
  - shows three rejected cases: an expired token, a tampered token (role changed to admin) and a token checked with the wrong secret.
- `docs/auth-notes.md`: notes on authentication vs authorization, password hashing, JWT structure, the secret and token expiry.

### What I learned

1. Authentication (who you are, 401) and authorization (what you can do, 403) are different checks.
2. bcrypt salts every hash automatically, so the same password produces a different hash each time, and it is slow on purpose.
3. A JWT is `header.payload.signature`. The payload is readable by anyone, but the signature stops anyone from changing it without the secret.

### How to run

```bash
cd Days/Day_13/auth-scratch
go mod tidy
# PowerShell
$env:JWT_SECRET="some-long-random-secret-for-testing"
go run .
```

## Day 14 - Registration & Login (Task 5, part 1)

### What I built

- User model: ID, Username (unique index in MongoDB), Password (bcrypt hash only), Role (admin/user)
- `POST /register`: validates username and password, hashes the password with bcrypt, stores the user
- `POST /login`: verifies credentials and returns a signed JWT
- Duplicate username returns 409; wrong password returns 401 with a generic message
- The password hash is never returned in any response

### What I learned

1. Passwords are hashed with bcrypt and never stored in plaintext. Hashing lives in the data layer.
2. A generic 401 message on login stops attackers from discovering which usernames exist.
3. A unique index in MongoDB enforces unique usernames, so the check is safe even under concurrent requests.

### How to run

```bash
cd Tasks/Task_5
cp .env.example .env   # set MONGO_URI, JWT_SECRET, PORT
go run main.go

# register
curl.exe -X POST http://localhost:<PORT>/register -H "Content-Type: application/json" -d "{\"username\":\"michael\",\"password\":\"secret123\"}"

# login
curl.exe -X POST http://localhost:<PORT>/login -H "Content-Type: application/json" -d "{\"username\":\"michael\",\"password\":\"secret123\"}"
```

## Day 15: JWT Middleware and Role-Based Access (Nov 3, 2026)

**What I built**

- `middleware/auth_middleware.go`: reads the `Authorization: Bearer <token>` header and validates the signature and expiry. It puts the user (ID, username, role) in the Gin context.
- An admin-only check in the middleware.
- Access rules:
  - Any logged-in user can `GET` tasks.
  - Only admins can `POST`, `PUT` and `DELETE` tasks.
  - Only admins can call `POST /promote/:id`.
- Postman tests with and without a token, with the token stored in a variable.
- Updated `docs/api_documentation.md` with registration, login, protected endpoints and request/response formats.

**What I learned**

1. How Gin middleware works: `c.Next()`, `c.Abort()` and passing data through `c.Set()` / `c.Get()`.
2. 401 vs 403: 401 means "I don't know who you are" (missing or invalid token), and 403 means "I know who you are, but you can't do this" (wrong role).
3. Always validate the token's signing method and expiry, and never trust the role from anywhere except the verified claims.

**How to run**

```bash
cd Tasks/Task_5
export JWT_SECRET=your_secret   # PowerShell: $env:JWT_SECRET="your_secret"
go run main.go
```

1. `POST /register`, then `POST /login` to get a token.
2. Send the token as `Authorization: Bearer <token>`.
3. Without a token, protected routes return 401. With a non-admin token, write routes return 403.

### Day 16: Clean Architecture planning (Nov 4, 2026)

**What I built**

- Reviewed my Task 5 Task Manager API file by file and mapped every piece to its Clean Architecture layer
- Wrote `docs/architecture.md` with the target folder structure, the dependency rule, a Task 5 → Task 6 mapping table, and a list of problems in the current design

**What I learned**

- The dependency rule: imports only point inward, so Domain imports nothing and Usecases import only Domain
- Dependency inversion: usecases define interfaces (like a task repository) and the Mongo code implements them
- My Task 5 code mixed layers: Mongo queries sat next to business rules, bcrypt ran inside storage code, and the controller generated JWTs

**Where to look**
`Tasks/Task_6/task_manager_api/docs/architecture.md`

This is a planning day, so there is nothing to run yet. The code refactor starts on Day 17.

### Day 17: Clean Architecture, Domain and Usecases (Task 6)

**What I built**

- `Domain/domain.go`: the `Task` and `User` entities, role constants and shared domain errors (`ErrNotFound`, `ErrInvalidInput`, `ErrUsernameTaken`, `ErrInvalidCredentials`)
- `Repositories/task_repository.go` and `user_repository.go`: repository interfaces only (the Mongo implementations come on Day 18)
- `Usecases/task_usecases.go`: `TaskUsecase` with validation (title, due date, status) and not-found handling
- `Usecases/user_usecases.go`: `UserUsecase` with register (first user becomes admin), login (same error for wrong password and unknown user) and promote
- `PasswordService` and `TokenService` interfaces, so hashing and JWT stay out of the business layer

**What I learned**

1. Usecases depend only on interfaces, so they know nothing about Gin or MongoDB.
2. Domain errors let the controller map a business failure to an HTTP status without the usecase knowing about HTTP.
3. A missing `go.mod` makes every `module/package` import fail, so `go mod init` comes first.

**How to run**

```
cd Tasks/Task_6/task_manager
go build ./Domain/... ./Repositories/... ./Usecases/...
```

The layers have no `main` yet; Day 18 wires everything together.

## Day 18 - Clean Architecture: repositories, infrastructure and delivery (Task 6, part 2)

**What I built**

- Mongo implementations of `TaskRepository` and `UserRepository`, each in the same file as its interface (`Repositories/`), with a `taskDoc`/`userDoc` mapping so bson details stay out of `Domain`.
- A unique index on `username` and mapping of duplicate-key errors to `Domain.ErrUsernameTaken`.
- Infrastructure services: `password_service.go` (bcrypt), `jwt_service.go` (HS256 token generation and validation with user ID, username and role claims) and `auth_middleWare.go` (Bearer token check and an admin-only guard).
- Delivery layer: `controllers/controller.go` calls only the usecase interfaces and maps domain errors to HTTP status codes, `routers/router.go` sets up public, logged-in and admin-only routes, and `main.go` wires everything by hand.

**What I learned**

- Manual dependency injection: `main.go` is the only place that knows the concrete types. It builds repositories, then usecases, then the controller, then the router.
- Defining `PasswordService` and `TokenService` interfaces in the usecase layer lets Infrastructure implement them without Usecases importing Infrastructure.
- Keeping `bson` and `ObjectID` inside the repository layer keeps `Domain` and `Usecases` free of framework imports.

**How to run**

```bash
cd Tasks/Task_6/task_manager
cp .env.example .env        # set MONGO_URI, MONGO_DB, PORT, JWT_SECRET
go mod tidy
go run ./Delivery
```

Routes: `POST /register`, `POST /login`, `GET /tasks`, `GET /tasks/:id`, and admin-only `POST /tasks`, `PUT /tasks/:id`, `DELETE /tasks/:id`, `POST /promote/:id`.

### Day 19: Task 6 wrap-up (Sat Nov 7)

**What I built**

- Final `docs/architecture.md`: layer diagram, dependency rule, folder tree, request flow, design decisions, and a guide for adding a feature
- Unit tests for `password_service` (hash, verify, wrong password, salting) and `jwt_service` (valid, expired, tampered, wrong secret)
- Re-ran the Task 5 Postman collection against the refactored API to confirm identical behaviour
- Tagged `v0.3`

**What I learned**

- The dependency rule: Domain imports nothing, and Usecases depend only on interfaces
- Moving the interfaces into Domain is what lets the database, bcrypt and JWT be swapped without touching business rules
- A refactor is only proven when the old tests pass unchanged

**How to run**
cd Tasks/Task_6/task_manager
go test ./Infrastructure/...
go run ./Delivery

## Day 20 - Testify and Mockery (Mon, Nov 9)

### What I built

- Scratch exercises in `Days/Day_20` covering:
  - `assert` vs `require` on a small `Divide` function
  - a table-driven test using `t.Run`
  - a suite that embeds `suite.Suite` and uses `SetupTest`
  - a mocked `TaskRepository` with `mock.On(...).Return(...)` and `AssertExpectations`
- The four mocks for Task 7, generated with mockery into `Tasks/Task_7/task_manager/Mocks/`:
  - `TaskRepositoryMock.go`
  - `UserRepositoryMock.go`
  - `TaskUsecaseMock.go`
  - `UserUsecaseMock.go`
- Confirmed the layers talk only through the four interfaces, and that `go build ./...` passes with the mocks.

### What I learned

1. `require` stops the test on failure and `assert` keeps going, so I use `require` for preconditions like `err` checks and `assert` for independent checks.
2. A suite with `SetupTest` gives every test fresh mocks, and `AssertExpectations` catches the case where the code never calls the dependency it should.
3. Mockery generates mocks from interfaces, so after an interface changes I regenerate the mock instead of editing it by hand.

### How to run

\`\`\`bash

# scratch exercises

cd Days/Day_20
go test -v ./...

# regenerate a mock (example)

cd Tasks/Task_7/task_manager
mockery --name=TaskRepository --dir=Repositories --output=Mocks --outpkg=Mocks --filename=TaskRepositoryMock.go --structname=TaskRepositoryMock
\`\`\`

### Day 21 (Nov 10): Unit-testing usecases with Testify suites

**What I built**

- Usecases/task_usecases_test.go: a suite covering GetAllTasks, GetTaskByID, CreateTask, UpdateTask and DeleteTask
- Usecases/user_usecases_test.go: a suite covering Register, Login and Promote
- 32 table-driven subtests with mocked repositories (no database), at 98.3% coverage

**What I learned**

1. A suite's SetupTest gives every test fresh mocks, so no state leaks between cases.
2. A mock panics on any unexpected call, so a passing "rejected" case proves the repository was never touched.
3. Because usecases depend only on interfaces, I can test every failure path (not found, duplicate username, wrong password, hashing and token errors) in milliseconds.

**How to run**
cd Tasks/Task_7/task_manager
go test ./Usecases -v
go test ./Usecases -cover

### Day 22 – Nov 11: Testing controllers with mocked usecases

**What I built**

- Delivery/controllers/controller_test.go: a Testify suite that injects TaskUsecaseMock and UserUsecaseMock into the controllers and sends requests through the real Gin router (with the real JWT middleware) using httptest.
- 36 table-driven subtests covering success, invalid JSON/input (400), unauthorized (401), forbidden (403), conflict (409) and usecase errors mapped to 404/500.
- Every test checks the parsed request, the usecase call and its arguments, and the status code, and finishes with AssertExpectations.
- GitHub Actions workflow that runs go vet and go test -race on every push.

**What I learned**

1. httptest lets me test HTTP handlers without a real server or database.
2. Because controllers depend only on usecase interfaces, mocks make the tests fast and isolated: Clean Architecture paying off.
3. Testing through the real router also tests the middleware, so 401 and 403 are verified end to end.

**How to run**
cd Tasks/Task_7/task_manager
go test ./...
go test -cover ./... # controllers 100%, usecases 98.3%

## Follow along

I post a short update on LinkedIn every day under the format above. Feel free to open an issue if you spot a mistake or have a better resource.
