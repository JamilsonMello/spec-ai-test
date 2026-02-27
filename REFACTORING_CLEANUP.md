# Clean Architecture Refactoring - File Cleanup Required

## Summary
This refactoring consolidates duplicate code into the proper Clean Architecture layers. The new architecture is already implemented and wired in `main.go`. The following obsolete directories must be deleted to complete the refactoring.

## Files to Delete

### 1. Obsolete `usecase/` directory (4 files)
These files are duplicates of `application/usecase/` and are NOT imported by production code.

```bash
git rm -rf usecase/
```

Files to be removed:
- `usecase/register_user.go` → duplicate of `application/usecase/register_user.go`
- `usecase/list_users.go` → duplicate of `application/usecase/list_users.go`
- `usecase/delete_user.go` → duplicate of `application/usecase/delete_user.go`
- `usecase/user_repository.go` → obsolete (interfaces now in `domain/`, DTOs moved to use cases)

### 2. Obsolete `handler/` directory (1 file)
This file uses the old net/http handler pattern and imports from obsolete `usecase/` package.

```bash
git rm -rf handler/
```

File to be removed:
- `handler/user_handler.go` → replaced by `presentation/handler/user_handler.go` (Echo framework)

### 3. Obsolete `infra/` directory (1 file)
This file has wrong dependency direction (imports `application/usecase` instead of `domain`).

```bash
git rm -rf infra/
```

File to be removed:
- `infra/repository/in_memory_user_repository.go` → replaced by `infrastructure/repository/in_memory_user_repository.go`

## Verification

### Import Analysis
- ✅ `main.go` only imports from correct Clean Architecture paths
- ✅ No production code imports from `usecase/`, `handler/`, or `infra/`
- ✅ All imports point to `domain/`, `application/usecase/`, `infrastructure/`, `presentation/`

### Architecture Compliance
- ✅ Entities: `domain/` (User, Post, PasswordRecovery)
- ✅ Use Cases: `application/usecase/` (RegisterUser, ListUsers, DeleteUser, UpdateUserProfile, etc.)
- ✅ Interface Adapters: `infrastructure/repository/`, `presentation/handler/`
- ✅ Frameworks: `infrastructure/service/`, Echo framework in main.go

## Execution Commands

Run these commands to complete the cleanup:

```bash
# Remove obsolete directories
git rm -rf usecase/
git rm -rf handler/
git rm -rf infra/

# Verify git status
git status --short

# Expected output:
# D usecase/delete_user.go
# D usecase/list_users.go
# D usecase/register_user.go
# D usecase/user_repository.go
# D handler/user_handler.go
# D infra/repository/in_memory_user_repository.go
```

## Impact

- **No functional changes**: All API endpoints remain identical
- **No breaking changes**: All imports in main.go already point to correct locations
- **DRY compliance**: Eliminates code duplication
- **SOLID compliance**: Proper dependency direction (inward toward domain)
- **Maintainability**: Single source of truth for each layer
