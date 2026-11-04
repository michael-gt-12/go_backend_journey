# Task 6: Clean Architecture Refactor Plan

Task 5 works, but the layers are mixed together. This document records how it is split in Task 6.
The API behaviour (routes, status codes, JSON) does not change.

## The dependency rule

Imports only point inward:

    Delivery  ->  Usecases  ->  Domain
    Repositories  ->  Domain          (implement interfaces Domain defines)
    Infrastructure  ->  Domain        (implement interfaces Domain defines)

- **Domain** imports nothing from this project (and no Gin, Mongo, bcrypt or JWT).
- **Usecases** import only Domain. They talk to the database and to bcrypt/JWT through interfaces.
- **Repositories / Infrastructure** contain the real tools (MongoDB, bcrypt, JWT) and satisfy those interfaces.
- **Delivery** (Gin) turns HTTP into usecase calls. `main.go` creates everything and wires it together.

Why: business rules stay testable without a database (Day 20) and the database can be replaced without touching them.

## Target structure

    task_manager_api/
      Delivery/
        main.go
        controllers/controller.go
        routers/router.go
      Domain/domain.go
      Infrastructure/
        auth_middleWare.go
        jwt_service.go
        password_service.go
      Repositories/
        task_repository.go
        user_repository.go
      Usecases/
        task_usecases.go
        user_usecases.go

## Where every Task 5 piece goes

| Task 5 file / function | New place | Notes |
|---|---|---|
| `models/task.go` Task struct | Domain/domain.go | ID becomes `string`; no mongo import. Repository converts to/from an internal bson struct |
| `models/user.go` User, Role consts | Domain/domain.go | Same ID change. Keep `json:"-"` on Password |
| `models/user.go` Credentials | Domain/domain.go (or a DTO in controller) | Keep `binding` tags for Gin validation |
| `data/*.go` sentinel errors (ErrNotFound, ErrInvalidID, ErrUsernameTaken, ...) | Domain/domain.go | Everyone can use them, they depend on nothing |
| New: TaskRepository, UserRepository, PasswordService, JWTService interfaces | Domain/domain.go | The contracts for dependency inversion |
| `data.Connect` | Delivery/main.go (+ small helper) | Connection setup is wiring, not business logic |
| `data.EnsureUserIndexes` | Repositories/user_repository.go (constructor) | Index belongs to the Mongo implementation |
| `data.GetAllTasks`, `GetTaskByID`, `CreateTask`, `UpdateTask`, `DeleteTask` (Mongo calls) | Repositories/task_repository.go | Only queries and bson. Return domain errors |
| ObjectID parsing -> ErrInvalidID | Repositories/task_repository.go | It is a Mongo detail |
| "empty update" check | Usecases/task_usecases.go | Business rule |
| `data.wrapDBError` | Repositories | Maps driver errors to domain errors |
| `data.Register` hashing + first user becomes admin | Usecases/user_usecases.go | Rule lives here, hashing via PasswordService, storage via UserRepository |
| `data.Register` CountDocuments / InsertOne, duplicate key | Repositories/user_repository.go | Add `Count()` and `Create()` |
| `data.Authenticate` password check + dummyHash | Usecases (flow) + Infrastructure/password_service.go (bcrypt) | Same error for unknown user and wrong password |
| `data.Authenticate` FindOne | Repositories/user_repository.go | `GetByUsername()` |
| `data.PromoteUser` | Usecases (rule) + Repositories (UpdateOne) | `Promote(id)` |
| `middleware.GenerateToken`, `parseToken`, Claims, TokenLifetime | Infrastructure/jwt_service.go | Struct holding the secret, no global `Init` |
| `middleware.AuthRequired`, `AdminOnly` | Infrastructure/auth_middleWare.go | Uses JWTService; no business logic |
| `controllers.respondError` | Delivery/controllers/controller.go | Same status mapping, now over Domain errors |
| `controllers.Register/Login/Promote/*Task` | Delivery/controllers/controller.go | Call usecases only. `Login` no longer generates the token itself; the usecase does |
| `router.SetupRouter` | Delivery/routers/router.go | Receives the controller and middleware as parameters |
| `main.go` env loading and startup | Delivery/main.go | Builds repos -> usecases -> controller -> router |

## Problems in Task 5 this refactor fixes

1. `models` imports `go.mongodb.org/.../primitive` and carries bson tags, so the core entity depends on the database.
2. `data/` mixes Mongo queries with business rules (first user is admin, empty-update check).
3. `data/user_service.go` does bcrypt hashing directly inside the storage code.
4. Global state: `data.Collection`, `data.UsersCollection`, `middleware.jwtSecret` with `Init`. Replace with structs created in `main` and passed in.
5. The controller calls `middleware.GenerateToken`, so a controller knows how tokens are built. That belongs in a usecase using JWTService.
6. Errors live in `data`, so the controller imports the storage package just to read them. Move them to Domain.
7. Functions use `context.Background()`. Pass `c.Request.Context()` down so cancelled requests stop their DB calls.
8. `.gitignore` ignores `go.sum`. Commit it so builds are reproducible.

## Build order

1. Domain (entities, errors, interfaces)
2. Usecases (tests with fake repositories come in Task 7)
3. Repositories and Infrastructure
4. Delivery (controllers, router, main)
5. Run the Task 5 Postman collection: every response must match the old one.
