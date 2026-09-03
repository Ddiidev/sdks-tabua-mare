# Tabua Mare SDK - Go

[![Go Reference](https://pkg.go.dev/badge/github.com/Ddiidev/sdks-tabua-mare/go.svg)](https://pkg.go.dev/github.com/Ddiidev/sdks-tabua-mare/go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Ddiidev/sdks-tabua-mare/go)](https://goreportcard.com/report/github.com/Ddiidev/sdks-tabua-mare/go)
[![Go Tests](https://github.com/Ddiidev/sdks-tabua-mare/actions/workflows/go-test.yml/badge.svg)](https://github.com/Ddiidev/sdks-tabua-mare/actions/workflows/go-test.yml)
[![codecov](https://codecov.io/gh/Ddiidev/sdks-tabua-mare/branch/main/graph/badge.svg?flag=go)](https://codecov.io/gh/Ddiidev/sdks-tabua-mare)

SDK oficial em Go para a API Tábua de Marés (v2) - dados precisos sobre marés na costa brasileira.

> **⚠️ v2:** a API migrou para `https://tabuamare.api.br/api/v2`. A v1 foi descontinuada (HTTP 410). Os IDs de porto agora são strings (ex: `pb01`, `pe02`).

## 🚀 Instalação

### 1. Inicializar o módulo Go

Se ainda não tiver um módulo Go, inicialize um:

```bash
go mod init seu-projeto
```

### 2. Baixar o SDK

Execute o comando para baixar o SDK:

```bash
go get github.com/Ddiidev/sdks-tabua-mare/go
```

### 3. Atualizar dependências

Execute `go mod tidy` para garantir que todas as dependências estejam corretas:

```bash
go mod tidy
```

Isso irá criar/atualizar os arquivos `go.mod` e `go.sum` automaticamente.

## 📖 Uso Rápido

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    tabuamare "github.com/Ddiidev/sdks-tabua-mare/go"
)

func main() {
    // Criar cliente
    client := tabuamare.NewClient()
    ctx := context.Background()
    
    // Listar estados costeiros
    states, err := client.GetStates(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Estados:", states)
    
    // Listar portos de Santa Catarina
    harbors, err := client.GetHarborNames(ctx, "sc")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Portos em SC: %d\n", len(harbors))
    
    // Obter tábua de marés para dias específicos
    tides, err := client.GetTideTable(ctx, "pb01", 1, []int{1, 2, 3})
    if err != nil {
        log.Fatal(err)
    }
    
    // Exibir dados
    for _, tide := range tides {
        fmt.Printf("Porto: %s\n", tide.HarborName)
        for _, month := range tide.Months {
            for _, day := range month.Days {
                fmt.Printf("Dia %d: %d registros de maré\n", day.Day, len(day.Hours))
            }
        }
    }
    
    // Encontrar porto mais próximo de São Paulo (sem restrição de estado)
    nearestHarbor, err := client.GetNearestHarbor(ctx, -23.550520, -46.633308)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Porto mais próximo: %s (%s)\n", nearestHarbor.HarborName, nearestHarbor.ID)
}
```

## ✨ Funcionalidades

- ✅ Listagem de estados costeiros brasileiros
- ✅ Consulta de portos por estado
- ✅ Detalhes completos de portos (IDs string: `pb01`, `pe02`, ...)
- ✅ Tábua de marés por período
- ✅ **Porto mais próximo por coordenadas GPS** (geral e por estado)
- ✅ **Tábua de marés por geolocalização** (`GetGeoTideTable`)
- ✅ **Consulta de uso da cota** (`GetUsage` com api_key)
- ✅ **Autenticação opcional via api_key** (`WithAPIKey`)
- ✅ Suporte a múltiplos portos
- ✅ Validação de parâmetros
- ✅ Tratamento de erros robusto
- ✅ Suporte a context.Context
- ✅ Configuração flexível do cliente
- ✅ Zero dependências externas

## 🔧 Configuração Avançada

```go
import "time"

// Cliente com timeout customizado
client := tabuamare.NewClient(
    tabuamare.WithTimeout(60 * time.Second),
)

// Cliente com api_key (64 req/min grátis; sem chave: 16 req/min por IP)
client := tabuamare.NewClient(
    tabuamare.WithAPIKey("tm_live_..."),
)

// Cliente com URL base customizada
client := tabuamare.NewClient(
    tabuamare.WithBaseURL("https://custom-api.com"),
)

// Cliente com http.Client customizado
httpClient := &http.Client{
    Timeout: 30 * time.Second,
}
client := tabuamare.NewClient(
    tabuamare.WithHTTPClient(httpClient),
)
```

## 📚 Exemplos

Veja a pasta `examples/` para mais exemplos:

- `examples/basic/` - Uso básico do SDK
- `examples/advanced/` - Configurações avançadas e tratamento de erros

## 🧪 Testes

```bash
# Testes unitários
go test -v

# Testes de integração (requer internet)
go test -v -tags=integration

# Teste manual completo
go run cmd/test/main.go

# Usando Makefile
make test              # Testes unitários
make test-integration  # Testes de integração
make test-all         # Todos os testes
make check            # Verificações completas (fmt, vet, lint, test)
```

Veja [TESTING.md](TESTING.md) para mais detalhes.

## 🔧 Desenvolvimento

```bash
# Instalar ferramentas de desenvolvimento
make install-tools

# Executar verificações
make check

# Executar testes
make test-all

# Ver todos os comandos disponíveis
make help
```

## 📦 Publicação

Para publicar uma nova versão, veja [PUBLISHING.md](PUBLISHING.md).

## 📖 Documentação

- [Documentação da API](https://tabuamare.api.br/docs)
- [pkg.go.dev](https://pkg.go.dev/github.com/Ddiidev/sdks-tabua-mare/go)
- [Repositório da API](https://github.com/Ddiidev/tabua_mare_api)

## 🤝 Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou pull requests.

## 📄 Licença

MIT - veja [LICENSE](LICENSE) para detalhes.

## 🌊 Sobre a API

A API Tide Table fornece dados precisos sobre marés na costa brasileira, cobrindo 17 estados costeiros. Os dados são fornecidos pela DHN (Diretoria de Hidrografia e Navegação) e DNPVN.

**Desenvolvido com ❤️ no Brasil 🇧🇷**
