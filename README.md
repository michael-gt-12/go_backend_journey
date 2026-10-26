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

| Day | Date         | Topic                                                                     | Code                                                           |
| --- | ------------ | ------------------------------------------------------------------------- | -------------------------------------------------------------- |
| 1   | Oct 18       | Environment setup, Go basics, hello-world, FizzBuzz CLI                   | `Days/Day_1`                                                   |
| 2   | Oct 19       | Structs, methods, interfaces, custom errors, generics                     | `Days/Day_2`                                                   |
| 3   | Oct 20, 2026 | Task 1 (word frequency, palindromes) + Task 2 skeleton                    | `Tasks/Task_1`, `Tasks/Task_2/library_management`              |
| 4   | Oct 21, 2026 | Task 2: Console Library Management System (interfaces, maps, errors)      | `Tasks/Task_2/library_management`                              |
| 5   | Oct 22, 2026 | Gin REST API basics: routing, JSON, handlers, request binding             | [Days/Day_5/gin-albums](Days/Day_5/gin-albums)                 |
| 6   | Oct 23       | Task 3: Gin project structure, in-memory store with mutex, task endpoints | [Tasks/Task_3/task_manager_api](Tasks/Task_3/task_manager_api) |
| 7   | Oct 24, 2026 | Task 3: POST/PUT/DELETE endpoints, validation, consistent errors          | Tasks/Task_3/task_manager                                      |
| 8   | Oct 26       | Postman API documentation, Task 3 finished                                | Tasks/Task_3/                                                  |

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

## Follow along

I post a short update on LinkedIn every day under the format above. Feel free to open an issue if you spot a mistake or have a better resource.
