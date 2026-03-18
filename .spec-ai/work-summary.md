# Work Summary — Current Implementation State

The previous session ran out of steps. This file has ALL the context you need.
Follow the progress.md checklist below — do NOT create new tasks or a new TODO.

## Spec (SOURCE OF TRUTH — every task must trace back to this)
# Implementar upload de foto de perfil do usuário

# FEATURE: Implementar upload de foto de perfil do usuário

## Context and Motivation
Atualmente, a entidade `User` não possui suporte para armazenar uma foto de perfil. Esta implementação visa adicionar o campo `ProfilePictureURL` à entidade de domínio em `domain/user.go`, persistir essa informação no PostgreSQL através do repositório existente e fornecer um caso de uso dedicado para o processamento de arquivos. O objetivo final é expor o endpoint `POST /usuarios/:id/foto-perfil` para que o frontend possa enviar imagens nos formatos suportados (JPG, JPEG, PNG).

## Out of Scope
- Redimensionamento automático de imagens.
- Upload para serviços de nuvem (S3, Cloudinary, etc) - o armazenamento será local.
- Deleção física de arquivos antigos quando um novo é enviado (apenas a URL no banco será sobrescrita).

## User Stories
- **Como usuário**, quero poder enviar uma imagem para o meu perfil para que eu possa ser identificado visualmente na plataforma.
- **Como sistema**, devo validar o tamanho e a extensão do arquivo enviado para garantir a segurança e integridade do armazenamento.

## Functional Requirements
- **FR-1**: Adicionar coluna `profile_picture_url` do tipo `TEXT` na tabela de usuários via migração.
- **FR-2**: O endpoint de upload deve aceitar apenas `multipart/form-data` com a chave `foto`.
- **FR-3**: Validar se o arquivo excede 5MB. Se exceder, retornar status 400.
- **FR-4**: Validar se a extensão é `.jpg`, `.jpeg` ou `.png`. Caso contrário, retornar status 400.
- **FR-5**: Gerar um nome de arquivo único usando UUID para evitar colisões no sistema de arquivos.
- **FR-6**: O caminho salvo no banco de dados deve ser o caminho relativo do arquivo no servidor.

## Database Changes
```sql
-- Migração: 202310270001_add_profile_picture_to_users.sql
ALTER TABLE users ADD COLUMN profile_picture_url TEXT;
```

## API Contracts
**POST /usuarios/:id/foto-perfil**
- **Request**: `multipart/form-data` | Body: `foto` (File)
- **Success 200**: `{ "message": "Foto atualizada com sucesso", "url": "/uploads/profile/uuid.png" }`
- **Error 400**: `{ "error": "Arquivo inválido ou muito grande", "code": "INVALID_FILE" }`
- **Error 404**: `{ "error": "Usuário não encontrado", "code": "USER_NOT_FOUND" }`
- **Error 500**: `{ "error": "Erro interno ao salvar arquivo", "code": "INTERNAL_ERROR" }`

## Key Entities
**User (em domain/user.go)**
- `ProfilePictureURL`: `string` (json: "profile_picture_url")

## Validation Rules
- **Field**: `foto`
- **Type**: `File`
- **Required**: `true`
- **Max Size**: `5242880 bytes (5MB)`
- **Allowed Extensions**: `image/jpeg`, `image/jpg`, `image/png`

## Error Handling
- Usuário não existe -> 404 Not Found -> `{ "error": "Usuário não encontrado" }`
- Arquivo > 5MB -> 400 Bad Request -> `{ "error": "Arquivo excede o limite de 5MB" }`
- Extensão proibida -> 400 Bad Request -> `{ "error": "Formato de arquivo não suportado" }`
- Erro de escrita em disco -> 500 Internal Server Error -> `{ "error": "Falha ao salvar imagem" }`

## Security Considerations
- Validar o Content-Type do arquivo no servidor, não apenas a extensão.
- Sanitizar o nome do arquivo original (usar UUID gerado internamente como nome final).
- Garantir que o diretório de uploads não tenha permissão de execução de scripts.

## Files to Modify
- **domain/user.go**: MODIFY - Adicionar `ProfilePictureURL string` à struct `User`.
- **infrastructure/repository/postgresql_user_repository.go**: MODIFY - Atualizar métodos `GetByID`, `Create` e `Update` para incluir a coluna `profile_picture_url`.
- **application/usecase/upload_profile_picture.go**: CREATE - Implementar lógica de validação, salvamento local e chamada ao repositório.
- **presentation/handler/user_handler.go**: MODIFY - Criar o método `UploadProfilePicture` que extrai o ID da URL e o arquivo do form-data.
- **main.go**: MODIFY - Registrar a rota `POST /usuarios/:id/foto-perfil` e configurar o diretório estático para servir as fotos.

## Implementation Guidance
1. `domain/user.go`: Adicionar campo à struct.
2. `infrastructure/repository/postgresql_user_repository.go`: Atualizar queries SQL (SELECT/INSERT/UPDATE).
3. `application/usecase/upload_profile_picture.go`: Criar struct e método `Execute`. Usar `os.Create` e `io.Copy` para salvar o arquivo.
4. `presentation/handler/user_handler.go`: Implementar handler usando o padrão dos handlers existentes.
5. `main.go`: Injetar dependências e registrar rota no roteador.

## Success Criteria
- **SC-1**: O campo `profile_picture_url` é persistido corretamente no banco de dados.
- **SC-2**: Arquivos maiores que 5MB são rejeitados com erro 400.
- **SC-3**: Apenas extensões JPG e PNG são aceitas.
- **SC-4**: O endpoint retorna a URL correta após o upload bem-sucedido.

## Changed Files
```
A  .spec-ai/claude-code-run1.log
A  .spec-ai/context.md
A  .spec-ai/original-progress.md
A  .spec-ai/plan.md
AM .spec-ai/progress.md
A  .spec-ai/reference.md
A  .spec-ai/rules.md
A  .spec-ai/session-log.md
A  .spec-ai/spec.md
```

## Diff Stats (unstaged)
```
 .spec-ai/progress.md | 7 ++++++-
 1 file changed, 6 insertions(+), 1 deletion(-)
```

## Diff Stats (staged)
```
 .spec-ai/claude-code-run1.log |  24 +++
 .spec-ai/context.md           |   5 +
 .spec-ai/original-progress.md |  12 ++
 .spec-ai/plan.md              |  36 ++++
 .spec-ai/progress.md          |  12 ++
 .spec-ai/reference.md         | 444 ++++++++++++++++++++++++++++++++++++++++++
 .spec-ai/rules.md             | 101 ++++++++++
 .spec-ai/session-log.md       |   7 +
 .spec-ai/spec.md              |  79 ++++++++
 9 files changed, 720 insertions(+)
```

## Implementation Progress (TASKS CHECKLIST)
Tasks marked [x] are DONE. Tasks marked [ ] are PENDING — continue these.

STATUS: IN_PROGRESS

## Tasks

- [x] TASK 1: Add ProfilePictureURL field to User struct in domain/user.go
- [x] TASK 2: Create SQL migration for profile_picture_url column
- [x] TASK 3: Update repository queries (SaveUser, GetUserByEmail, FindUserByUuid, UpdateUser, ListUsers) to include profile_picture_url
- [x] TASK 4: Create upload_profile_picture.go usecase
- [x] TASK 5: Add UploadProfilePicture handler method to user_handler.go
- [x] TASK 6: Add error mappings for upload errors in error_mapper.go
- [x] TASK 7: Wire usecase, update handler constructor, register route and static file serving in main.go
- [x] TASK 8: Build and verify


## Compliance — Missing Requirements (fix these)
- [ ] Requirement: File Validation - The implementation only validates the file extension; it fails to validate the Content-Type (image/jpeg, image/png) as required by Security Considerations and Validation Rules. Fix: In internal/application/usecase/upload_profile_picture.go, add a check for req.Header.Header.Get("Content-Type")
- [ ] Requirement: API Contract & Error Handling - Error responses (400, 404, 500) are missing the mandatory 'code' field (INVALID_FILE, USER_NOT_FOUND, INTERNAL_ERROR) in the JSON body. Fix: Update internal/presentation/handler/user_handler.go to include the 'code' field in the JSON response when an error occurs during upload


## Implementation Plan
Cross-reference with progress.md above to see what's DONE vs what's LEFT.

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

## Changes Made (diff preview)
```diff
diff --git a/.spec-ai/claude-code-run1.log b/.spec-ai/claude-code-run1.log
new file mode 100644
index 0000000..076eb16
--- /dev/null
+++ b/.spec-ai/claude-code-run1.log
@@ -0,0 +1,24 @@
+Implementation complete. Here's a summary of all changes:
+
+**Files modified:**
+- `internal/domain/user.go` — Added `ProfilePictureURL string` field to `User` struct
+- `internal/domain/repository_interfaces.go` — Added `UpdateProfilePictureURL` method to `UserRepository` interface
+- `internal/infrastructure/repository/user_repository.go` — Updated all SQL queries (SaveUser, GetUserByEmail, FindUserByUuid, UpdateUser, ListUsers) to include `profile_picture_url` column; added `UpdateProfilePictureURL` method
+- `internal/presentation/handler/user_handler.go` — Added `UploadProfilePictureUseCase` to handler struct and `UploadProfilePicture` handler method (multipart/form-data with key `foto`)
+- `internal/presentation/handler/error_mapper.go` — Added error mappings for `ErrFileTooLarge`, `ErrUnsupportedFileFormat`, `ErrSaveFileFailed`, `ErrUserNotFoundUpload`
+- `main.go` — Wired `UploadProfilePictureUseCase`, updated handler constructor, registered `POST /usuarios/:id/foto-perfil` in protected group, added static file serving for `/uploads`
+
+**Files created:**
+- `migrations/202310270001_add_profile_picture_to_users.sql` — ALTER TABLE with UP/DOWN migration
+- `internal/application/usecase/upload_profile_picture.go` — Validates file size (5MB), extension (.jpg/.jpeg/.png), generates UUID filename, saves to disk, updates DB
+
+**Spec verification:**
+- FR-1: Migration adds `profile_picture_url TEXT` column
+- FR-2: Endpoint accepts `multipart/form-data` with key `foto`
+- FR-3: Files >5MB rejected with 400
+- FR-4: Only .jpg/.jpeg/.png accepted
+- FR-5: UUID-based filenames prevent collisions
+- FR-6: Relative path (`/uploads/profile/uuid.ext`) stored in DB
+- Route is inside the `protected` group with `AuthMiddleware`
+- All error messages in Portuguese matching codebase language
+- Build passes cleanly
diff --git a/.spec-ai/context.md b/.spec-ai/context.md
new file mode 100644
index 0000000..e584fbc
--- /dev/null
+++ b/.spec-ai/context.md
@@ -0,0 +1,5 @@
+# Project Reference
+
+## Entry Point (main.go)
+Read this file to understand the DI wiring pattern before adding new dependencies.
+
diff --git a/.spec-ai/original-progress.md b/.spec-ai/original-progress.md
new file mode 100644
index 0000000..a8fd87d
--- /dev/null
+++ b/.spec-ai/original-progress.md
@@ -0,0 +1,12 @@
+STATUS: COMPLETE
+
+## Tasks
+
+- [x] TASK 1: Add ProfilePictureURL field to User struct in domain/user.go
+- [x] TASK 2: Create SQL migration for profile_picture_url column
+- [x] TASK 3: Update repository queries (SaveUser, GetUserByEmail, FindUserByUuid, UpdateUser, ListUsers) to include profile_picture_url
+- [x] TASK 4: Create upload_profile_picture.go usecase
+- [x] TASK 5: Add UploadProfilePicture handler method to user_handler.go
+- [x] TASK 6: Add error mappings for upload errors in error_mapper.go
+- [x] TASK 7: Wire usecase, update handler constructor, register route and static file serving in main.go
+- [x] TASK 8: Build and verify
diff --git a/.spec-ai/plan.md b/.spec-ai/plan.md
new file mode 100644
index 0000000..a9b2f08
--- /dev/null
+++ b/.spec-ai/plan.md
@@ -0,0 +1,36 @@
+1. FILES TO CREATE:
+- `internal/application/usecase/upload_profile_picture.go`: Logic for file validation, UUID naming, and disk storage.
+- `migrations/202310270001_add_profile_picture_to_users.sql`: SQL migration to add the column.
+
+2. FILES TO MODIFY:
+- `internal/domain/user.go`: Add `ProfilePictureURL` field to `User` struct.
+- `internal/infrastructure/repository/user_repository.go`: Update `GetByID`, `Create`, and `Update` SQL queries.
+- `internal/presentation/handler/user_handler.go`: Add `UploadProfilePicture` method to handle multipart/form-data.
+- `main.go`: Register the new usecase, inject into handler, and configure static file serving for `/uploads`.
+
+3. IMPLEMENTATION ORDER:
+TASK 1: Domain & Database
+- Files: `internal/domain/user.go`, `migrations/202310270001_add_profile_picture_to_users.sql`
+- Validation: `go build ./...` passes; verify column `profile_picture_url` exists in local Postgres.
+TASK 2: Repository Layer
+- Files: `internal/infrastructure/repository/user_repository.go`
+- Validation: Update repository tests or verify that `Update` method correctly persists a string in the new column.
+TASK 3: Application Layer (Usecase)
+- Files: `internal/application/usecase/upload_profile_picture.go`
+- Validation: Unit test for file extension validation and UUID generation logic.
+TASK 4: Presentation Layer (Handler & Routing)
+- Files: `internal/presentation/handler/user_handler.go`, `main.go`
+- Validation: `curl` request with `multipart/form-data` returns 200 and file exists in `./uploads/profile/`.
+
+4. INTEGRATION CHECKLIST:
+- MIGRATIONS: PK is likely `UUID` (based on `google/uuid` dependency and FR-5). New column is `TEX
... (diff truncated — read individual files only if needed)
```

## Session History
Each entry shows what a previous session accomplished and where it stopped.

# Session Log
History of what each session accomplished.

### Session 1 — completed
- Progress: 8/8 tasks complete
- Last completed: TASK 8: Build and verify



## What To Do Next
1. Continue the UNCHECKED [ ] tasks in progress.md above. Do NOT create new tasks or a new TODO.
2. Re-read .spec-ai/spec.md before each task — implement EXACTLY what the spec says.
3. For each pending task: search for types/functions → implement → build → mark [x].
4. When ALL tasks are done: verify each spec requirement, build, set STATUS: COMPLETE.
IMPORTANT: Ignore any TODO list that OpenCode generates. Only .spec-ai/progress.md matters.
