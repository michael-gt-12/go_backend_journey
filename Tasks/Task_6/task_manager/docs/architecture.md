# Task Manager API: Clean Architecture

Task 5 worked, but HTTP, business rules, MongoDB, bcrypt and JWT were mixed together. Task 6 splits them into layers.
**The API behaviour (routes, status codes, JSON) is unchanged.** The Task 5 Postman collection passes as before.

## 1. The dependency rule

Imports only point inward, toward Domain:

```
            +-------------------------------+
            |           Delivery            |  Gin controllers, router, main.go
            +---------------+---------------+
                            |
                            v
            +-------------------------------+
            |           Usecases            |  business rules
            +---------------+---------------+
                            |
                            v
            +-------------------------------+
            |            Domain             |  entities, errors, interfaces
            +-------------------------------+
                   ^                ^
                   |                |
        +----------+------+  +------+-----------+
        |  Repositories   |  |  Infrastructure  |
        |  (MongoDB)      |  |  (bcrypt, JWT,   |
        |                 |  |   middleware)    |
        +-----------------+  +------------------+
```

| Layer          | May import                                     | Must NOT import                                       |
| -------------- | ---------------------------------------------- | ----------------------------------------------------- |
| Domain         | standard library only                          | any other project package, Gin, Mongo, bcrypt, JWT    |
| Usecases       | Domain                                         | Gin, Mongo, bcrypt, JWT, Repositories, Infrastructure |
| Repositories   | Domain, Mongo driver                           | Usecases, Delivery, Gin                               |
| Infrastructure | Domain, bcrypt, JWT, Gin (middleware only)     | Usecases, Repositories                                |
| Delivery       | Usecases, Domain, Infrastructure (wiring only) | Mongo driver                                          |

`Delivery/main.go` is the only place that knows every concrete type. It builds the pieces and passes them in.

**Why:** business rules can be tested without a database (Task 7 mocks), and the database, hashing library or web framework can be swapped without touching them.

## 2. Folder structure

```
task_manager/
  Delivery/
    main.go                    # config, wiring (manual dependency injection), server start
    controllers/controller.go  # HTTP -> usecase calls, error -> status code mapping
    routers/router.go          # routes; receives controller + middleware as parameters
  Domain/
    domain.go                  # Task, User, Role, Credentials, domain errors, interfaces
  Infrastructure/
    auth_middleWare.go         # AuthRequired, AdminOnly (uses JWTService)
    jwt_service.go             # generates and validates tokens
    password_service.go        # bcrypt hash and compare
  Repositories/
    task_repository.go         # TaskRepository: MongoDB implementation
    user_repository.go         # UserRepository: MongoDB implementation
  Usecases/
    task_usecases.go           # TaskUsecase
    user_usecases.go           # UserUsecase
  docs/
    architecture.md            # this file
    api_documentation.md
```

## 3. How a request flows

Example: `POST /tasks` (admin creates a task)

1. **router** matches the route and runs `AuthRequired` then `AdminOnly` (Infrastructure, via `JWTService`).
2. **controller** binds and validates the JSON body, then calls `TaskUsecase.Create(ctx, task)`.
3. **usecase** applies the business rules and calls `TaskRepository.Create(ctx, task)` (an interface).
4. **repository** converts the task to a bson struct, runs the Mongo insert and returns a domain `Task` or a domain error.
5. **controller** maps the result or error to a status code (201, 400, 404, 500) and writes the JSON.

The usecase never sees HTTP, and the controller never sees MongoDB.

## 4. Where each Task 5 piece went

| Task 5                                       | Task 6                             | Note                                                         |
| -------------------------------------------- | ---------------------------------- | ------------------------------------------------------------ |
| `models/task.go`, `models/user.go`           | Domain/domain.go                   | ID is now `string`; no Mongo import                          |
| sentinel errors in `data/`                   | Domain/domain.go                   | shared by every layer                                        |
| `data/*` Mongo queries                       | Repositories/                      | only queries and bson; driver errors mapped to domain errors |
| ObjectID parsing -> `ErrInvalidID`           | Repositories/task_repository.go    | a Mongo detail                                               |
| first user becomes admin, empty-update check | Usecases/                          | business rules                                               |
| bcrypt hashing                               | Infrastructure/password_service.go | used through `PasswordService`                               |
| token creation and parsing                   | Infrastructure/jwt_service.go      | struct holding the secret, no globals                        |
| `AuthRequired`, `AdminOnly`                  | Infrastructure/auth_middleWare.go  | no business logic                                            |
| controllers, router                          | Delivery/                          | controller calls usecases only                               |
| startup and env loading                      | Delivery/main.go                   | repos -> usecases -> controller -> router                    |

## 5. Design decisions

1. **Interfaces are defined in Domain**, and implemented in Repositories and Infrastructure (dependency inversion). Usecases depend on `TaskRepository`, `UserRepository`, `PasswordService` and `JWTService`, never on concrete tools.
2. **ID is a `string` in Domain.** The Mongo repository converts to and from `ObjectID` internally, so the entity does not depend on the database.
3. **Domain errors live in Domain.** The repository maps driver errors to them, and the controller maps them to HTTP status codes. Nobody imports a storage package just to read an error.
4. **No global state.** The collection handles and the JWT secret now live in structs created in `main` and passed in. This replaces `data.Collection`, `data.UsersCollection` and `middleware.jwtSecret`.
5. **The usecase builds the token** (through `JWTService`), so the controller does not know how tokens are made.
6. **`context.Context` is passed down** from `c.Request.Context()`, so a cancelled request stops its database calls.
7. **Manual constructor injection** in `main.go`, with no DI framework. It is explicit and easy to read.
8. **Same error for unknown user and wrong password** at login, so the API does not reveal which usernames exist.
9. **`go.sum` is committed** so builds are reproducible.

## 6. Problems in Task 5 this refactor fixed

1. `models` imported the Mongo `primitive` package and carried bson tags.
2. `data/` mixed Mongo queries with business rules.
3. bcrypt was called directly inside the storage code.
4. Global mutable state (`Collection`, `UsersCollection`, `jwtSecret`).
5. The controller knew how tokens were generated.
6. The controller imported the storage package to read errors.
7. Functions used `context.Background()`, so requests could not cancel DB work.
8. `go.sum` was git-ignored.

## 7. How to add a new feature

Example: add a `Comment` on a task.

1. **Domain:** add the `Comment` struct, any new errors, and a `CommentRepository` interface.
2. **Usecases:** add `comment_usecases.go` with a `CommentUsecase` interface and implementation. It depends only on `CommentRepository` and the existing Domain types.
3. **Repositories:** add `comment_repository.go` that implements `CommentRepository` with Mongo.
4. **Delivery:** add handlers to the controller that call `CommentUsecase`, map errors to status codes, and register the routes in `router.go`.
5. **Wire it** in `Delivery/main.go`: repository -> usecase -> controller.
6. **Check the rule:** `grep -r "gin-gonic\|mongo-driver\|bcrypt\|golang-jwt" Domain Usecases` must print nothing.
7. **Document** the new endpoints in `docs/api_documentation.md`.

## 8. Testing

- `Infrastructure/password_service_test.go` and `jwt_service_test.go` cover the critical security components (valid, wrong, expired and tampered cases).
- Task 7 adds Testify suites for the usecases and controllers, with mocks of the four interfaces.
- Regression: the Task 5 Postman collection runs against this version with identical results.
