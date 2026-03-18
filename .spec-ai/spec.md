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