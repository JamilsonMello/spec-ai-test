# CLAUDE.md

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

## Progress — .spec-ai/progress.md
Create BEFORE coding with ALL tasks as a checklist (STATUS: IN_PROGRESS).
Mark [x] after each task. Set STATUS: COMPLETE only when ALL done + build passes.
System REJECTS COMPLETE if any [ ] tasks remain. If the file exists, continue unchecked tasks.

## Before completing — VERIFY EACH REQUIREMENT
Re-read .spec-ai/spec.md. For each requirement, SEARCH the codebase to confirm:
- Deletions → file/code is actually gone
- Moves → ALL imports updated (search old path — zero matches expected)
- Comment removal → search for // and /* in target files
- New features → route + handler + use case + DI wiring all exist
Do NOT mark COMPLETE until every requirement is verified.
## Commands
- Build: `go build ./...`

## Existing features (do NOT recreate — reuse)
### Estrutura de Arquitetura Limpa (Go)
- Entities: Entidades de domínio puras e objetos de valor localizados em `internal/domain`.
- Use Cases: Lógica de negócio e DTOs organizados em `internal/application` seguindo a regra de dependência.
- Endpoints: Handlers HTTP, controladores e middlewares estruturados em `internal/presentation`.
- Gateways: Repositórios, clientes HTTP e adaptadores externos implementados em `internal/infrastructure`.

### Remoção de Documentação Obsoleta
- Entities: Arquivos `CLEAN_ARCHITECTURE_PLAN.md` e `README.md`.
- Use Cases: Exclusão definitiva de guias de refatoração legados; limpeza de referências cruzadas e links órfãos.
- Endpoints: N/A.
- Gateways: File System (OS) para deleção e atualização de arquivos de documentação.

### Refatoração Global para Clean Code
- Entities: Tipos de erro nomeados; constantes de domínio; identificadores descritivos (ex: `FindUserByUuid`).
- Use Cases: `main.go` como Composition Root; injeção de dependência via interfaces; funções com responsabilidade única (< 20 linhas).
- Endpoints: Mapeamento centralizado de erros de domínio para códigos de status HTTP na camada de apresentação.
- Gateways: Remoção total de comentários (// e /* */) em arquivos `.go`; abstração de repositórios via interfaces `Reader`/`Writer`.

## Context compaction
Preserve: spec (.spec-ai/spec.md), plan (.spec-ai/plan.md), reference (.spec-ai/reference.md), commands above, and which files were created/modified.
