# Refatoração Global para Clean Code e Clean Architecture

# REFACTORING: Global Refactoring for Clean Code and Clean Architecture

## Context and Motivation
O projeto atual contém arquivos de documentação interna (como `CLEAN_ARCHITECTURE_PLAN.md`) e diversos comentários nos arquivos `.go` que poluem a leitura do código. Além disso, a estrutura de pastas e a nomenclatura de variáveis, funções e tipos precisam ser aprimoradas para refletir a semântica do negócio sem a necessidade de explicações externas. Esta refatoração visa alinhar o repositório aos padrões de Clean Code e garantir que a regra de dependência da Clean Architecture seja estritamente respeitada (as camadas internas não devem conhecer as externas).

## Out of Scope
- Criação de novas funcionalidades ou endpoints.
- Alteração de esquemas de banco de dados (exceto renomeação de campos em structs se necessário para clareza).
- Adição de novas bibliotecas de terceiros.

## Current vs Expected Behavior
- **Atual:** Arquivos `.go` possuem comentários de linha e bloco; nomes de funções são genéricos; o arquivo `main.go` contém lógica de inicialização misturada; existe um arquivo de plano de arquitetura na raiz.
- **Esperado:** Zero comentários em arquivos `.go`; nomes altamente descritivos; `main.go` limpo e apenas orquestrando a injeção de dependência; estrutura de pastas dividida em `domain`, `application`, `infrastructure` e `presentation`; arquivo `CLEAN_ARCHITECTURE_PLAN.md` removido.

## Functional Requirements
- **FR-1:** Remover o arquivo `CLEAN_ARCHITECTURE_PLAN.md` da raiz do projeto.
- **FR-2:** Remover todos os comentários (// e /* */) de todos os arquivos com extensão `.go`.
- **FR-3:** Renomear identificadores (variáveis, structs, interfaces, funções) para nomes que descrevam exatamente sua intenção (ex: de `h` para `handler`, de `GetUser` para `FindUserByUuid`).
- **FR-4:** Quebrar funções com mais de 20 linhas ou múltiplas responsabilidades em funções menores e privadas.
- **FR-5:** Mover arquivos para as camadas corretas: 
  - `domain`: Entidades e interfaces de repositório.
  - `application`: Casos de uso e DTOs.
  - `infrastructure`: Implementações de banco de dados, clientes HTTP e adaptadores externos.
  - `presentation`: Handlers HTTP, controladores e middlewares.

## Key Entities
As entidades existentes no diretório `domain` (ou equivalente [verify path in codebase]) devem ser puras, contendo apenas dados e lógica de domínio, sem tags de frameworks de infraestrutura se possível, ou usando apenas o necessário para serialização básica.

## Error Handling
- Substituir comentários de erro por tipos de erro nomeados e constantes claras no pacote de domínio.
- Centralizar o mapeamento de erros de domínio para códigos HTTP na camada de `presentation`.

## Files to Modify
- **Raiz do Projeto:**
  - `CLEAN_ARCHITECTURE_PLAN.md`: DELETE.
  - `main.go`: MODIFY - Refatorar para atuar apenas como o Composition Root, instanciando dependências e iniciando o servidor.
- **Diretórios .go (Recursivo):**
  - `**/*.go`: MODIFY - Remover todos os comentários e aplicar renomeação descritiva.
- **Novos Diretórios (Estrutura Alvo):**
  - `internal/domain/`: CREATE/MOVE - Entidades e interfaces.
  - `internal/application/`: CREATE/MOVE - Services/UseCases.
  - `internal/infrastructure/`: CREATE/MOVE - Repositories/Clients.
  - `internal/presentation/`: CREATE/MOVE - Handlers/Controllers.

## Implementation Guidance
1. **Limpeza Inicial:** Excluir `CLEAN_ARCHITECTURE_PLAN.md` e rodar um script/comando para remover comentários em todos os arquivos `.go`.
2. **Refatoração de Nomes:** Começar pela camada de `domain` e subir para `application`. Renomear interfaces para nomes que reflitam o comportamento (ex: `Reader`, `Writer`, `Processor`).
3. **Injeção de Dependência:** Garantir que todos os constructors (ex: `NewUseCase`, `NewHandler`) recebam interfaces como argumentos, não implementações concretas.
4. **Organização de Pastas:** Mover os arquivos existentes para a nova estrutura de pastas definida nos Requisitos Funcionais, ajustando os `import` paths conforme necessário.
5. **Main Clean-up:** Mover lógicas de configuração de banco e roteamento para funções auxiliares ou pacotes de infraestrutura, mantendo o `main.go` minimalista.

## Success Criteria
- **SC-1:** Nenhum arquivo `.go` contém comentários após a execução.
- **SC-2:** O arquivo `CLEAN_ARCHITECTURE_PLAN.md` não existe mais.
- **SC-3:** O projeto compila e executa com sucesso após as mudanças de pacotes e nomes.
- **SC-4:** A estrutura de diretórios reflete exatamente as quatro camadas da Clean Architecture.
- **SC-5:** Não há dependências circulares entre os pacotes.