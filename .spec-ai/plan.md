1. FILES TO CREATE
- `internal/presentation/error_mapper.go`: Logic to map domain errors to HTTP codes.
- `internal/domain/errors.go`: Named error constants and types (refactored from `domain/errors.go`).

2. FILES TO MODIFY
- `main.go`: Refactor as Composition Root, update all import paths.
- `internal/domain/*.go`: Remove comments, rename identifiers, ensure purity.
- `internal/application/usecase/*.go`: Remove comments, rename identifiers, use interfaces.
- `internal/infrastructure/repository/*.go`: Remove comments, rename identifiers.
- `internal/presentation/handler/*.go`: Remove comments, rename identifiers, use error mapper.

3. IMPLEMENTATION ORDER
TASK 1: Structural Migration (infrastructure)
- Files: Move `domain`, `application`, `infrastructure`, `presentation` to `internal/`.
- Validation: `go build ./...` fails only on import errors; directory structure matches spec.
TASK 2: Global Import Update & Comment Removal (infrastructure)
- Files: All `.go` files.
- Action: Update imports to `github.com/example/cadastro-de-usuarios/internal/...`. Run `grep -r "//" .` to ensure no comments remain.
- Validation: `go build ./...` passes.
TASK 3: Domain Refinement (domain)
- Files: `internal/domain/*.go`
- Action: Rename `h` to `handler`, `GetUser` to `FindUserByUuid`, etc. Define error constants.
- Validation: No infrastructure tags in entities; all interfaces named by behavior (e.g., `UserStorer`).
TASK 4: Use Case & Dependency Injection (application)
- Files: `internal/application/usecase/*.go`
- Action: Break functions > 20 lines. Ensure constructors take interfaces.
- Validation: `internal/application` has zero imports from `infrastructure` or `presentation`.
TASK 5: Presentation & Error Mapping (presentation)
- Files: `internal/presentation/handler/*.go`, `internal/presentation/error_mapper.go`
- Action: Remove logic from handlers; call Use Cases; map domain errors to HTTP.
- Validation: Handlers are < 30 lines and only handle I/O and orchestration.
TASK 6: Composition Root & Cleanup (presentation)
- Files: `main.go`, `CLEAN_ARCHITECTURE_PLAN.md`
- Action: Delete `CLEAN_ARCHITECTURE_PLAN.md`. Move DB init logic to `internal/infrastructure/db.go`.
- Validation: `main.go` contains only wiring and `e.Start()`.

4. RISKS
- Circular Dependencies: Moving files into a nested `internal` structure might expose tight coupling.
- Import Hell: Every file requires an import change; manual errors are likely without `goimports`.
- Logic Loss: Removing all comments might obscure "why" certain complex regex or hacks exist.