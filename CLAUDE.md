# Spec AI Builder

You are implementing a software specification inside a repository.
Read the following files carefully before starting:

- `.spec-ai/spec.md` — The full specification to implement
- `.spec-ai/plan.md` — Step-by-step implementation plan (if available)
- `.spec-ai/progress.md` — Track your progress here; check off tasks as you complete them
- `.spec-ai/context.md` — Project context and structure
- `.spec-ai/reference.md` — Reference implementation to follow (if available)

When ALL tasks are done, set `STATUS: COMPLETE` in progress.md.

## Project Rules

# Project Rules

You are implementing a feature spec at .spec-ai/spec.md in an existing codebase.

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
- Build: `go build ./...`

## Existing features (do NOT recreate — reuse)
### Estrutura de Arquitetura Limpa (Go)
- Entities: Entidades de domínio puras e interfaces de repositório localizadas em `internal/domain`.
- Use Cases: Lógica de negócio e DTOs organizados em `internal/application` seguindo a regra de dependência.
- Endpoints: Handlers HTTP, controladores e middlewares estruturados em `internal/presentation`.
- Gateways: Implementações de banco de dados, clientes HTTP e adaptadores em `internal/infrastructure`.

### Remoção de Documentação Obsoleta
- Entities: Deleção definitiva dos arquivos `CLEAN_ARCHITECTURE_PLAN.md` e `README.md` legados.
- Use Cases: Exclusão de guias de refatoração antigos; limpeza de referências cruzadas e links órfãos.
- Endpoints: N/A.
- Gateways: File System (OS) para remoção física de arquivos de documentação interna da raiz do projeto.

### Refatoração Global para Clean Code
- Entities: Identificadores descritivos (ex: `FindUserByUuid`); tipos de erro nomeados e constantes de domínio.
- Use Cases: `main.go` como Composition Root; injeção de dependência via interfaces; funções com responsabilidade única (< 20 linhas).
- Endpoints: Mapeamento centralizado de erros de domínio para códigos de status HTTP na camada de apresentação.
- Gateways: Remoção total de comentários (// e /* */) em arquivos `.go`; abstração de repositórios via interfaces `Reader` e `Writer`.

## Context compaction
Preserve: spec (.spec-ai/spec.md), plan (.spec-ai/plan.md), reference (.spec-ai/reference.md), commands above, and which files were created/modified.

