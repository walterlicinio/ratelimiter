# Rate Limiter em Go com Redis

Rate limiter implementado em Go que limita o número de requisições por segundo com base no endereço IP ou token de acesso, utilizando Redis para persistência.

## Funcionalidades
- Limitação por endereço IP
- Limitação por token de acesso
- Configuração de limites de requisições por segundo
- Configuração de tempo de bloqueio
- Suporte a persistência de dados com Redis

## Pré-requisitos

- Docker
- Docker Compose

## Instalação

### Clone o Repositório

Clone o repositório para a sua máquina local:

```sh
git clone https://github.com/walterlicinio/ratelimiter.git
cd ratelimiter
```

### Crie um Arquivo `.env` e um `tokens.json`

Na raiz do projeto, crie um arquivo `.env` com o seguinte conteúdo:

```plaintext
REDIS_ADDR=redis:6379
REDIS_PASSWORD=senharedis
RATE_LIMIT_IP=5
RATE_LIMIT_TOKEN=10
BLOCK_TIME_SECONDS=300
```

Na mesma raiz, crie um arquivo `tokens.json` com o mapeamento string:int parra os tokens e seu tempo de expiração:
Por exemplo:
```
{
    "abc123": 2,
    "fifo123": 5
}
```

## Execução com Docker Compose

### 1. Construir e Executar os Containers

Use Docker Compose para construir e executar os containers para a aplicação e Redis:

```sh
docker-compose up --build
```

### 2. Acessar a Aplicação

A aplicação estará disponível em `http://localhost:8080`.

## Endpoints

### GET /

Endpoint principal que retorna "Hello, World!".

```plaintext
GET /
```

## Testes

Para rodar os testes automatizados:

```sh
go test -v -count=1 ./tests
```

Os testes cobrem as seguintes verificações:

1. **Limitação por Endereço IP**:
   - Permite até 5 requisições de um IP por segundo.
   - A 6ª requisição do mesmo IP é bloqueada.
   - O identificador (IP) é marcado como bloqueado após exceder o limite.
   - Após esperar o tempo de expiracão do bloqueio, uma nova requisição deve ser permitida.
   - Verifica se o bloqueio é removido após o período de expiração.

2. **Limitação por Token**:
   - Usa um mapa de limites específicos para cada token (`allowedTokens`).
   - Permite até 3 requisições para "token1" por segundo.
   - A 4ª requisição com "token1" é bloqueada.
   - O identificador (token) é marcado como bloqueado após exceder o limite específico.
   - Após esperar o tempo de expiracão do bloqueio, uma nova requisição deve ser permitida.
   - Verifica se o bloqueio é removido após o período de expiração.

   
Os testes utilizam [`miniredis`](https://github.com/alicebob/miniredis), um servidor Redis em memória, para testar a lógica sem precisar de um servidor Redis em execução real.

