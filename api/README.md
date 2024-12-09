
# GOLEDGER CHALLENGE REST API

Este projeto consiste em uma API escrita em Go que interage com um **smart contract** previamente "deployado" na **blockchain Besu**. A API permite interagir com o **smart contract**, armazenar o valor do contrato em um banco de dados SQL e fornecer endpoints para recuperar, definir e sincronizar valores entre a blockchain e o banco de dados.

## Database
A aplicação usa um banco de dados **PostgreSQL** para armazenar o valor da variável do smart contract. O schema do banco de dados inclui uma tabela para armazenar os valores do contrato.

## Endpoints

### 1. `POST /simple-storage/set/value`
**Descrição:**  
Define um novo valor para a variável do **smart contract**. Esse valor é enviado para o **smart contract**,

**Request Body:**
```json
{
  "value": "1997"
}
```

**Response:**
```json
{
  "message": "value set successfully"
}
```

**Error Response:**
```json
{
  "error": "invalid request"
}
```

### 2. `GET /simple-storage/get/value`
**Descrição:**  
Retorna o valor atual da variável do **smart contract**.

**Response:**
```json
{
  "value": 1997
}
```

**Error Response:**
```json
{
  "error": "Failed to retrieve contract value"
}
```

### 3. `GET /simple-storage/sync/value`
**Descrição:**  
Sincroniza o valor da variável do **smart contract** para o banco de dados SQL. Esta operação armazenará o valor atual do **smart contract** no banco de dados.

**Response:**
```json
{
  "message": "Value synchronized successfully",
  "value": 1997
}
```

**Error Response:**
```json
{
  "error": "Failed to synchronize value with database"
}
```

### 4. `GET /simple-storage/check/value`
**Descrição:**  
Compara o valor armazenado no banco de dados com o valor atual da variável do **smart contract**. Retorna `true` se forem iguais, caso contrário, retorna `false`.

**Response:**
```json
{
  "isEqual": true
}
```

**Error Response:**
```json
{
  "error": "Failed to compare values"
}
```

## Instruções

### 1. Clone o repositório
```
git clone https://github.com/JoaoVFerreira/goledger-challenge-besu.git
```

### 2. Crie o arquivo `.env` na raiz do projeto
Use o arquivo `.env.example` como base, que contém algumas sugestões:

```
CONTRACT_ADDRESS=your_contract_address
BLOCKCHAIN_NODE=http://your_besu_node_url
PRIVATE_KEY=your_private_key
DATABASE_URL=postgres://username:password@localhost:5432/your_db_name
PORT=8080
```

### 3. Instale os pacotes
```bash
go mod tidy
```

### 4. Inicie a API
```bash
go run ./cmd/main.go
```

A aplicação vai estar disponivel em `http://localhost:8080`.

## Tecnologias
- **Go**: Linguagem de programação.
- **PostgreSQL**: Banco de dados SQL.
- **Gin**: Framework para construção de APIs REST.
- **Blockchain**: Integração com a blockchain Besu.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
