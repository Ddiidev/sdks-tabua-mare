# Guia de Testes - Tabua Mare SDK

## 🧪 Tipos de Testes

### 1. Testes Unitários

Testes que não fazem chamadas reais à API, usando mocks.

```bash
# Executar todos os testes unitários
go test -v
```

```bash
# Executar testes de integração
go test -v -tags=integration

# Com timeout maior
go test -v -tags=integration -timeout 30s
```

### 3. Teste Manual Completo

Um programa que testa todas as funcionalidades do SDK com a API real.

```bash
# Executar teste manual
go run cmd/test/main.go
```

## 📋 O que é testado

### Testes Unitários (`client_test.go`, `nearest_harbor_test.go`, `v2_endpoints_test.go`)
- ✅ Criação do cliente
- ✅ Configuração de URL base customizada
- ✅ Configuração de timeout
- ✅ Configuração de api_key e headers de autenticação
- ✅ Tratamento de respostas bem-sucedidas
- ✅ Tratamento de rate limiting (429) com Retry-After
- ✅ Tratamento de erros da API (envelope aninhado `error.message`)
- ✅ IDs de porto string com lista entre colchetes
- ✅ Endpoints v2: nearested-harbor, geo-tabua-mare, usage

### Testes de Integração (`integration_test.go`)
- ✅ Listagem de estados
- ✅ Listagem de portos por estado
- ✅ Obtenção de detalhes de porto
- ✅ Obtenção de tábua de marés
- ✅ Porto mais próximo por estado
- ✅ Tábua de marés por geolocalização

### Teste Manual (`cmd/test/main.go`)
- ✅ Listagem de estados
- ✅ Listagem de portos de Santa Catarina
- ✅ Detalhes de um porto específico
- ✅ Tábua de marés para dias específicos
- ✅ Consulta de múltiplos portos
- ✅ Porto mais próximo por estado
- ✅ Tábua de marés por geolocalização
- ✅ Uso da cota (opcional: defina `TABUAMARE_API_KEY`)
- ✅ Validação de erros

## 🚀 Exemplos de Uso

### Exemplo Básico

```bash
# Executar exemplo básico
go run examples/basic/main.go
```

### Exemplo Avançado

```bash
# Executar exemplo avançado
go run examples/advanced/main.go
```