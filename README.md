# Desafio Técnico Rate Limiter

**Objetivo**: Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

**Descrição**: O objetivo deste desafio é criar um rate limiter em Go que possa ser utilizado para controlar o tráfego de requisições para um serviço web. O rate limiter deve ser capaz de limitar o número de requisições com base em dois critérios:

1. **Endereço IP**: O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.

2. **Token de Acesso**: O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:
   API_KEY: TOKEN
3. As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.

**Requisitos**:

- O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web
- O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.
- O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
- As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env” na pasta raiz.
- Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.
- O sistema deve responder adequadamente quando o limite é excedido:
  - Código HTTP: **429**
  - Mensagem: **you have reached the maximum number of requests or actions allowed within a certain time frame**
- Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.
- Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.
- A lógica do limiter deve estar separada do middleware.

**Entrega**:

- O código-fonte completo da implementação.
- Documentação explicando como o rate limiter funciona e como ele pode ser configurado.
- Testes automatizados demonstrando a eficácia e a robustez do rate limiter.
- Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
- O servidor web deve responder na porta 8080.

## Rate Limiter

O Rate Limiter e um middleware que possibilita limitar o número de requisições por segundo com base em um endereço IP ou um token de acesso.
Nessa implementação, deve ser passar para a função que cria a estrutura o número de requisições permitidas por segundo sem o header API_KEY definido, o número de requisições permitidas por segundo com o header API_KEY definido, o tempo em segundos quando bloqueado e um estrutura do cache.

A estrutura de Cache deve ser a interface a seguir:

```go
type Cache interface {
  Get(key string) (int, bool)
  Set(key string, value int)
}
```

As configurações são feitas através de variáveis de ambiente ou em um arquivo “.env” na pasta raiz.

```env
PORT=8080 # porta do servidor web

MAX_IP_REQUESTS_PER_SECOND=10 # maximo de requisições por segundo sem o header API_KEY
MAX_TOKEN_REQUESTS_PER_SECOND=100 # maximo de requisições por segundo com o header API_KEY
TIME_BLOCK_IN_SECOND=5 # tempo de bloqueio em segundos

REDIS_HOST=localhost # host do redis (cache escolhido para a implementação)
REDIS_PORT=6379 # porta do redis
```
