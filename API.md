# Documentacao da API

## Endpoints

### POST /usuarios

Cadastra um novo usuario no sistema.

**URL**: `/usuarios`
**Metodo**: `POST`
**Autenticacao**: Nenhuma (endpoint publico)

---

#### Request Body

| Campo       | Tipo   | Obrigatorio | Descricao                                                                 |
|-------------|--------|-------------|---------------------------------------------------------------------------|
| `name`      | string | Sim         | Nome do usuario. Minimo 2, maximo 50 caracteres. Apenas letras e espacos. |
| `surname`   | string | Sim         | Sobrenome do usuario. Minimo 2, maximo 50 caracteres. Apenas letras e espacos. |
| `email`     | string | Sim         | E-mail valido e unico no sistema.                                         |
| `birthDate` | string | Sim         | Data de nascimento no formato `YYYY-MM-DD` (ISO 8601).                    |

**Exemplo de requisicao:**

```json
{
  "name": "Joao",
  "surname": "Silva",
  "email": "joao.silva@email.com",
  "birthDate": "1990-05-15"
}
```

---

#### Regras de Negocio

- **Idade minima**: O usuario deve ter no minimo **18 anos** com base na data atual. Caso contrario, a requisicao sera rejeitada.
- **Formato de data**: O campo `birthDate` deve seguir o formato `YYYY-MM-DD`. Datas futuras nao sao aceitas.
- **Unicidade de e-mail**: O campo `email` deve ser unico. Nao e permitido cadastrar um e-mail ja existente na base de dados.
- **Validacao de nome/sobrenome**: Devem conter entre 2 e 50 caracteres, utilizando apenas letras e espacos.

---

#### Respostas

##### Sucesso — 200 OK

Retorna os dados do usuario cadastrado.

```json
{
  "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "name": "Joao",
  "surname": "Silva",
  "email": "joao.silva@email.com",
  "birthDate": "1990-05-15"
}
```

| Campo       | Tipo   | Descricao                          |
|-------------|--------|------------------------------------|
| `id`        | string | UUID gerado automaticamente.       |
| `name`      | string | Nome do usuario.                   |
| `surname`   | string | Sobrenome do usuario.              |
| `email`     | string | E-mail do usuario.                 |
| `birthDate` | string | Data de nascimento (`YYYY-MM-DD`). |

---

##### Erro — 400 Bad Request

Retornado quando o payload nao pode ser interpretado ou quando o e-mail ja esta em uso.

**Payload invalido:**

```json
{
  "error": "Invalid request payload"
}
```

**E-mail ja cadastrado:**

```json
{
  "error": "email já está em uso"
}
```

---

##### Erro — 422 Unprocessable Entity

Retornado quando os campos possuem formato valido de JSON, mas nao atendem as regras de validacao.

**Nome invalido:**

```json
{
  "error": "nome deve ter entre 2 e 50 caracteres e conter apenas letras e espaços"
}
```

**Sobrenome invalido:**

```json
{
  "error": "sobrenome deve ter entre 2 e 50 caracteres e conter apenas letras e espaços"
}
```

**E-mail invalido:**

```json
{
  "error": "email inválido"
}
```

**Data de nascimento invalida:**

```json
{
  "error": "data de nascimento inválida"
}
```

**Data de nascimento no futuro:**

```json
{
  "error": "data de nascimento não pode ser no futuro"
}
```

**Usuario menor de 18 anos:**

```json
{
  "error": "usuário deve ter no mínimo 18 anos"
}
```

---

##### Erro — 500 Internal Server Error

Retornado em caso de falha inesperada no servidor ou no banco de dados.

```json
{
  "error": "Internal server error"
}
```

---

## Modelos de Dados

### User

| Campo               | Tipo      | Descricao                          |
|---------------------|-----------|-------------------------------------|
| `id`                | UUID      | Identificador unico do usuario.     |
| `name`              | string    | Nome do usuario.                    |
| `surname`           | string    | Sobrenome do usuario.               |
| `email`             | string    | E-mail do usuario (unico).          |
| `birthDate`         | date      | Data de nascimento (`YYYY-MM-DD`).  |
| `role`              | string    | Papel do usuario (padrao: `user`).  |
| `profile_picture_url` | string | URL da foto de perfil.              |
| `createdAt`         | timestamp | Data de criacao do registro.        |

---

## Erros Comuns

| Codigo HTTP | Causa                                      | Mensagem de erro                                                              |
|-------------|--------------------------------------------|-------------------------------------------------------------------------------|
| 400         | Payload JSON mal formatado                 | `Invalid request payload`                                                     |
| 400         | E-mail ja cadastrado                       | `email já está em uso`                                                        |
| 422         | Nome fora das regras                       | `nome deve ter entre 2 e 50 caracteres e conter apenas letras e espaços`      |
| 422         | Sobrenome fora das regras                  | `sobrenome deve ter entre 2 e 50 caracteres e conter apenas letras e espaços` |
| 422         | Formato de e-mail invalido                 | `email inválido`                                                              |
| 422         | Data de nascimento invalida                | `data de nascimento inválida`                                                 |
| 422         | Data de nascimento no futuro               | `data de nascimento não pode ser no futuro`                                   |
| 422         | Usuario menor de 18 anos                   | `usuário deve ter no mínimo 18 anos`                                          |
| 500         | Falha interna no servidor ou banco de dados | `Internal server error`                                                       |
