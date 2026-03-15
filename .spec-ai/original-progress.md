STATUS: COMPLETE

## TASK 1: Structural Migration (infrastructure)
- [x] Move `domain`, `application`, `infrastructure`, `presentation` to `internal/`
- [x] Update all import paths to use `internal/` prefix
- [x] Run `go build ./...` to verify no import errors

## TASK 2: Global Import Update & Comment Removal (infrastructure)
- [x] Update all imports to `github.com/example/cadastro-de-usuarios/internal/...`
- [x] Remove all comments (// and /* */) from all `.go` files
- [x] Run `go build ./...` to verify build passes

## TASK 3: Domain Refinement (domain)
- [x] Rename identifiers for clarity (e.g., `h` to `handler`, `GetUser` to `FindUserByUuid`)
- [x] Define error constants in `internal/domain/errors.go`
- [x] Ensure entities have no infrastructure tags or only basic serialization tags
- [x] Rename interfaces by behavior (e.g., `UserStorer`, `UserFinder`)

## TASK 4: Use Case & Dependency Injection (application)
- [x] Break functions > 20 lines into smaller private functions
- [x] Ensure constructors take interfaces not concrete implementations
- [x] Verify `internal/application` has zero imports from `infrastructure` or `presentation`

## TASK 5: Presentation & Error Mapping (presentation)
- [x] Create `internal/presentation/error_mapper.go` for domain error to HTTP mapping
- [x] Remove business logic from handlers
- [x] Ensure handlers are < 30 lines and only handle I/O orchestration
- [x] Update handlers to use error mapper

## TASK 6: Composition Root & Cleanup (presentation)
- [x] Delete `CLEAN_ARCHITECTURE_PLAN.md` (if exists)
- [x] Move DB init logic to `internal/infrastructure/db.go`
- [x] Refactor `main.go` to be minimal Composition Root
- [x] Ensure `main.go` contains only wiring and `e.Start()`