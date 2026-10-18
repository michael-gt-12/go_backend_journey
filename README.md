# go_backend_journey

Repo: https://github.com/michael-gt-12/go_backend_journey

My public learning log for a **62-day Go backend + AI backend track**, running from **Sun Oct 18, 2026 to Sun Jan 3, 2027**. I'm Michael, a 4th-year Software Engineering student at Adama Science & Technology University (ASTU), and I'm building this in public so every claim is backed by code, tests and commits.

**Goal:** go from Go basics to deployed, tested backend services, including four AI-powered backend projects, and finish with a portfolio I can point recruiters to.

## How I work

- About 3 focused hours a day: 30-45 min learning, about 2 hrs building, 15 min to commit, update this README and write the LinkedIn post.
- Every day has a "Done when" line that acts as my acceptance test. I push only when it is true.
- Sundays are buffer days for slipped work, the weekly check-in and the weekly recap post.
- Mid-exam week (Nov 23-27) is maintenance only. The plan resumes Sat Nov 28.
- Post format: `Day N/90 - what I built - what I learned - what is next - repo link`.

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

- [ ] Day 1: environment setup, Go basics, repo and README
- [ ] Days 2-4: Tasks 1-2
- [ ] Days 5-8: Task 3 (Gin + Postman)
- [ ] Days 9-12: Task 4 (MongoDB)
- [ ] Days 13-15: Task 5 (JWT)
- [ ] Days 16-19: Task 6 (Clean Architecture)
- [ ] Days 20-22: Task 7 (testing)
- [ ] Days 23-27: Task 8 (concurrency)
- [ ] Days 28-33: Task 9 (PostgreSQL)
- [ ] Days 34-36: Task 10 (Docker)
- [ ] Days 37-39: Task 11 (Redis)
- [ ] Days 40-41: Task 12 (CI + Render)
- [ ] Days 42-46: SupportDesk AI
- [ ] Days 47-52: DocuAsk
- [ ] Days 53-57: LLM Gateway
- [ ] Days 58-62: OpsAgent

## Day 1 (Oct 18)

- Installed Go and set up the editor
- Initialised this repo with `go mod init`
- Hello-world and a small FizzBuzz CLI
- Tour of Go: Basics

## Follow along

I post a short update on LinkedIn every day under the format above. Feel free to open an issue if you spot a mistake or have a better resource.
