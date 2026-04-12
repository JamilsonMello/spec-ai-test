# Spec AI Builder

You are implementing a software specification inside a repository.
Read the following files carefully before starting:

- `.spec-ai/spec.md` — The full specification to implement
- `.spec-ai/plan.md` — Step-by-step implementation plan (if available)
- `.spec-ai/progress.md` — Track your progress here; check off tasks as you complete them
- `.spec-ai/context.md` — Project context and structure
- `.spec-ai/reference.md` — Reference implementation to follow (if available)
- `.spec-ai/analysis.md` — Codebase analysis from spec creation (entities, routes, patterns)
- `.spec-ai/tests-reference.md` — Example test file (match its framework and patterns)
- `.spec-ai/schema.md` — Database migrations for exact column types and constraints

When ALL tasks are done, set `STATUS: COMPLETE` in progress.md.

## Project Rules

# Project Rules

You are implementing a feature spec at .spec-ai/spec.md in an existing codebase.

## P2 — LANGUAGE MIRROR (ABSOLUTE)
EVERY text output — progress updates, section headers, labels, confirmations, and status messages — MUST be in the SAME language as the spec content in .spec-ai/spec.md. If the spec is in Portuguese, respond entirely in Portuguese. If in English, respond in English. ALL auto-generated headers and section titles MUST follow the detected language. NEVER translate code identifiers: file paths, function names, variable names, type names, table/column names, API routes, and any string containing file extensions (.go, .js, .py) or camelCase/snake_case patterns from the codebase must remain exactly as they appear in the source code.

## Rules
1. Follow every existing pattern: structure, naming, error handling, imports, DI, validation. Read one complete feature end-to-end before coding.
2. Implement exactly what the spec says — no more, no less.
3. No comments (no //, /*, doc comments). Exception: files that already have comments. Never reference "spec", "AI", "agent", or "auto-generated".
4. New strings must match the codebase's language, not the spec's.
5. Never guess a type, function, or import path — search and read its definition first.
6. Build after every layer. Fix errors by searching for the symbol's definition.
7. When renaming/moving: search ALL usages first, update ALL references, then build. Never revert a rename — fix remaining references instead.
8. If a target directory already exists, read its contents — it's your destination, not legacy code.
9. Never introduce libraries the project doesn't use. Parameterized queries only. Never log secrets.

## Code Quality — Match the Project (MANDATORY)
Your code must look like it was written by the same team that built the existing codebase. This applies to ALL parts of the project — backend, frontend, infrastructure, scripts, tests — in ANY language.

When implementing each file, read ONE existing file in the same layer as a reference (e.g., read user_handler before writing password_reset_handler). Copy its conventions exactly:
- File/folder names: match the existing naming pattern and casing (kebab-case, camelCase, snake_case, PascalCase — whatever the project uses)
- Function/method/hook names: follow the same verb patterns and casing conventions already in the code
- Variable/parameter/prop names: match the existing style and abbreviation conventions
- Type/interface/class/struct names: follow the project's naming convention — never introduce a different style
- Component structure (frontend): match how existing components are organized (file per component, barrel exports, style co-location, hook patterns)
- State management (frontend): use the same state library, patterns, and naming as existing code (context, stores, hooks, reducers)
- Styling (frontend): use the same approach as existing components (CSS modules, Tailwind, styled-components, inline — whatever exists)
- Error handling: use the same error types, wrapping patterns, and response codes as existing code
- Directory structure: place files in the same directories as similar existing features — never create new directories when the feature fits an existing one
- Import/module organization: group and order imports the same way as existing files
- Dependency injection / wiring: follow the same constructor or provider pattern as existing code
- Schema/migration naming: continue the existing numbering sequence and naming convention
Do NOT impose your own preferences or "best practices" that differ from the project. Match what exists.

## Progress — .spec-ai/progress.md
Create BEFORE coding with ALL tasks as a checklist (STATUS: IN_PROGRESS).
Mark [x] after each task. Set STATUS: COMPLETE only when ALL done + build passes.
System REJECTS COMPLETE if any [ ] tasks remain. If the file exists, continue unchecked tasks.

## Integration Completeness — MANDATORY
These checks apply to EVERY implementation. Skipping them causes real production failures.

### Migrations & Database
- Before creating ANY migration, READ all existing migration files (especially the table you're referencing). Match column types EXACTLY — if users.id is UUID, your foreign key MUST be UUID, not VARCHAR(36).
- Down migrations MUST undo EVERYTHING the up migration does: DROP all indexes, constraints, triggers, and tables — in reverse order of creation.
- Continue the existing migration numbering sequence and naming convention.

### Environment Variables
- NEVER invent new environment variables. Search the codebase (.env, docker-compose*.yml, .github/workflows/*.yml) for existing env vars and REUSE them.
- If a URL is needed, compose it from existing env vars (e.g., use FRONTEND_URL + path, not a new FRONTEND_RESET_PASSWORD_URL).
- If a service config already exists (e.g., email sender address is fixed), hardcode it — don't create an env var.
- If a genuinely NEW env var is required, add it to ALL: .env.example, docker-compose files, AND CI/CD pipeline files (.github/workflows/*.yml).

### Language & Strings
- ALL backend strings (error messages, email templates, API responses, log messages) MUST match the existing codebase's language. Read 3-5 existing error messages to confirm the language — NEVER assume.
- The spec's language may differ from the codebase's — always follow the CODEBASE, not the spec.
- If the project uses error codes for i18n, use the same pattern for new errors.

### Frontend Integration
- Every new user-facing feature MUST have an entry point from existing UI (link, button, menu item). A page without a link to reach it is useless.
- If the app has auth guards, route protection, or public route lists: add new unauthenticated pages (forgot-password, reset-password, etc.) to the public/unprotected list.
- Match existing URL patterns. Read the routing structure before choosing paths.
- NEVER modify auto-generated files (next-env.d.ts, .next/, auto-generated lock files). If they exist, leave them untouched.

### Code Organization
- Read the file structure BEFORE adding code. Place types, interfaces, and constants where the existing code puts them (e.g., Go convention in this project: all types at the top, methods below).
- Never scatter new type definitions between existing methods.

## Before completing — VERIFY EACH REQUIREMENT
Re-read .spec-ai/spec.md. For each requirement, SEARCH the codebase to confirm:
- Deletions → file/code is actually gone
- Moves/renames → ALL imports/references updated (search old path — zero matches expected)
- New features → all layers exist and are wired together (data model, logic, presentation, routing/config)
- New pages/routes → accessible via navigation AND not blocked by auth guards (if user is unauthenticated)
- New env vars → present in ALL deployment configs, not just code
- Migrations → down migration is complete and column types match referenced tables
- Run the build/compile command and ensure it passes
Do NOT mark COMPLETE until every requirement is verified.
## Commands
- Build: `go build -o bin/main main.go`
- Test: `go test ./...`

## Existing features (do NOT recreate — reuse)
### Estrutura de Arquitetura Limpa (Go)
- Entities: Domínio e interfaces em `internal/domain`.
- Use Cases: Lógica de negócio e DTOs em `internal/application`.
- Endpoints: Handlers HTTP e middlewares em `internal/presentation`.
- Gateways: Implementações de persistência e adaptadores em `internal/infrastructure`.

### Cadastro e Gerenciamento de Usuários
- Entities: `User` (id, name, surname, email, birth_date, password, role, ProfilePictureURL).
- Use Cases: Validação de idade (18+), e-mail único, formato de senha (8+ chars).
- Endpoints: POST /usuarios, GET /usuarios/listar, GET /usuarios/:id, PATCH /usuarios/:id, PUT /usuarios/:id/foto-perfil, POST /users/:id/upload-picture, DELETE /usuarios/:id.
- Gateways: `PostgreSQLUserRepository`, `LocalFileSystem`, `FileStorage`.

### Upload e Persistência de Posts
- Entities: `Post` (id, user_id, title, content, VideoURL, community_id).
- Use Cases: Validação de autorização, tamanho (50MB), tipo (mp4/mov), título (5-100 chars), validação de existência de comunidade (UUID v4).
- Endpoints: POST /posts, GET /posts, GET /posts/:id, POST /posts/:id/video, PUT /posts/:id.
- Gateways: `FileStorageService`, `PostRepository`, `CommunityRepository`.

### Segurança e Recuperação de Senha
- Entities: `PasswordRecovery` (id, user_id, token, expires_at).
- Use Cases: Geração de token, expiração, envio de e-mail.
- Endpoints: POST /password-recovery, POST /password-recovery/reset.
- Gateways: `EmailService`, `JWTProvider`.

### Comentários e Reações
- Entities: `Comment` (id, post_id, user_id, content, created_at), `Reaction` (id, comment_id, user_id, type, created_at).
- Use Cases: Criação de comentário (max 500 chars, sanitização XSS), Toggle de reação, validação de existência do post.
- Endpoints: POST /posts/:id/comments, GET /posts/:id/comments, POST /comments/:id/reactions.
- Gateways: `CommentRepository`, `ReactionRepository`.

### Criação e Gestão de Comunidades
...(truncated)


## Constitutional Rules (MANDATORY — defined by the project owner)
The project owner has defined the following rules. You MUST follow them during implementation.
These rules take precedence over general guidelines when there is a conflict.

- sempre use a arquitetura limpa

## Known Build Issues (from previous attempts in this project)
Previous builds have failed with these issues. Pay extra attention:
- Database Migration: A migração migrations/202310270003_create_communities_and_link_to_posts.sql não foi atualizada para incluir a coluna community_id na tabela posts. É necessário criar ou editar o arquivo de migração para adicionar a coluna: ALTER TABLE posts ADD COLUMN community_id UUID REFERENCES communities(id) ON DELETE SET NULL;
- Database Migration: O arquivo de migração migrations/202310270003_create_communities_and_link_to_posts.sql não foi modificado para incluir a coluna community_id na tabela posts com a devida FK. O diff mostra a alteração no código Go (SavePost), mas a mudança no esquema do banco de dados é necessária para persistência real. Adicione ALTER TABLE posts ADD COLUMN community_id UUID REFERENCES communities(id) ao arquivo de migração.
- Database Migration: The file migrations/202310270003_create_communities_and_link_to_posts.sql is not present in the diff. You must ensure this file exists and contains: ALTER TABLE posts ADD COLUMN community_id UUID REFERENCES communities(id) ON DELETE SET NULL;

## Context compaction
Preserve: spec (.spec-ai/spec.md), plan (.spec-ai/plan.md), reference (.spec-ai/reference.md), commands above, and which files were created/modified.

