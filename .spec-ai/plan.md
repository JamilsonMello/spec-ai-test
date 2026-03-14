1. FILES TO CREATE
- `internal/domain/errors.go`: Named domain error constants/types.
- `internal/presentation/error_mapper.go`: Logic to map domain errors to Echo HTTP errors.

2. FILES TO MODIFY
- `CLEAN_ARCHITECTURE_PLAN.md`: DELETE.
- `main.go`: Refactor as Composition Root (DI and server start only).
- `internal/domain/*.go`: Move from `domain/`, remove tags, rename identifiers.
- `internal/application/*.go`: Move from `application/usecase/`, update to use interfaces.
- `internal/infrastructure/*.go`: Move from `infrastructure/` and `pkg/db/`.
- `internal/presentation/*.go`: Move from `presentation/`, implement error mapping.

3. IMPLEMENTATION ORDER
TASK 1: Structural Migration (Infrastructure)
- Files: Move all folders to `internal/`, update `go.mod` and all import paths.
- Validation: `go build ./...` passes with new import paths.
TASK 2: Domain Purity & Error Definition (Domain)
- Files: `internal/domain/*.go`
- Validation: Entities have no infrastructure tags; `errors.go` contains named error types.
TASK 3: Interface & Use Case Refinement (Application)
- Files: `internal/application/*.go`
- Validation: Use cases accept interfaces in constructors; functions > 20 lines are split.
TASK 4: Infrastructure Implementation (Infrastructure)
- Files: `internal/infrastructure/repository/*.go`, `internal/infrastructure/service/*.go`
- Validation: Implements domain interfaces; DB-specific logic isolated.
TASK 5: Presentation & Error Mapping (Presentation)
- Files: `internal/presentation/handler/*.go`, `internal/presentation/error_mapper.go`
- Validation: Handlers use `error_mapper` to return HTTP codes; no logic besides orchestration.
TASK 6: Composition Root & Global Cleanup (Root)
- Files: `main.go`, `**/*.go`
- Validation: `main.go` is minimal; `grep -r "//" .` returns zero results in `.go` files.

4. RISKS
- Import Path Cascades: Moving to `internal/` requires updating every file; high risk of compilation errors.
- Name Collisions: Renaming generic variables (e.g., `db`) to descriptive ones (e.g., `postgresDatabaseConnection`) across the whole project.
- Logic Loss: Removing all comments might remove "TODOs" or critical business context not yet captured in code.