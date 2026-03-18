1. FILES TO CREATE:
- `internal/application/usecase/upload_profile_picture.go`: Logic for file validation, UUID naming, and disk storage.
- `migrations/202310270001_add_profile_picture_to_users.sql`: SQL migration to add the column.

2. FILES TO MODIFY:
- `internal/domain/user.go`: Add `ProfilePictureURL` field to `User` struct.
- `internal/infrastructure/repository/user_repository.go`: Update `GetByID`, `Create`, and `Update` SQL queries.
- `internal/presentation/handler/user_handler.go`: Add `UploadProfilePicture` method to handle multipart/form-data.
- `main.go`: Register the new usecase, inject into handler, and configure static file serving for `/uploads`.

3. IMPLEMENTATION ORDER:
TASK 1: Domain & Database
- Files: `internal/domain/user.go`, `migrations/202310270001_add_profile_picture_to_users.sql`
- Validation: `go build ./...` passes; verify column `profile_picture_url` exists in local Postgres.
TASK 2: Repository Layer
- Files: `internal/infrastructure/repository/user_repository.go`
- Validation: Update repository tests or verify that `Update` method correctly persists a string in the new column.
TASK 3: Application Layer (Usecase)
- Files: `internal/application/usecase/upload_profile_picture.go`
- Validation: Unit test for file extension validation and UUID generation logic.
TASK 4: Presentation Layer (Handler & Routing)
- Files: `internal/presentation/handler/user_handler.go`, `main.go`
- Validation: `curl` request with `multipart/form-data` returns 200 and file exists in `./uploads/profile/`.

4. INTEGRATION CHECKLIST:
- MIGRATIONS: PK is likely `UUID` (based on `google/uuid` dependency and FR-5). New column is `TEXT`.
- ENV VARS: Existing: `DATABASE_URL`. New: `UPLOAD_PATH` (defaults to `./uploads`), `BASE_URL` (for returning full URLs if needed).
- LANGUAGE: Portuguese (e.g., "Usuário não encontrado", "Erro interno"). Follow this for all new error strings.
- FRONTEND ENTRY: New endpoint `POST /usuarios/:id/foto-perfil`. Triggered via profile edit page.
- AUTH/ROUTES: The route MUST be inside the `protected` group in `main.go` using `AuthMiddleware`.
- DEPLOY: Ensure the deployment environment has a persistent volume mounted at the `UPLOAD_PATH` directory.

5. RISKS:
- Directory Permissions: The application process must have write access to the `uploads/` folder.
- Path Traversal: Mitigated by using generated UUIDs instead of original filenames.
- Concurrent Uploads: `os.Create` with UUID prevents collisions, but disk space monitoring is required.