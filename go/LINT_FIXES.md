# Correções do Linter

## ✅ Problemas Corrigidos

### 1. Erros não verificados (errcheck)

**Erro:**
```
Error return value of `w.Write` is not checked
```

**Correção em `client_test.go`:**
```go
// ❌ Antes
w.Write([]byte(`{"data": [], "total": 0}`))

// ✅ Depois  
_, _ = w.Write([]byte(`{"data": [], "total": 0}`))
```

Em testes, podemos ignorar explicitamente o erro usando `_, _`.

### 2. Parâmetros não utilizados (revive)

**Erro:**
```
unused-parameter: parameter 'r' seems to be unused
```

**Correção em `client_test.go`:**
```go
// ❌ Antes
func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}

// ✅ Depois
func(w http.ResponseWriter, _ *http.Request) {
    w.WriteHeader(http.StatusOK)
}
```

Use `_` para parâmetros que não serão usados.

### 3. fmt.Println com \n redundante (govet)

**Erro:**
```
fmt.Println arg list ends with redundant newline
```

**Correção em `cmd/test/main.go`:**
```go
// ❌ Antes
fmt.Println("Mensagem\n")

// ✅ Depois
fmt.Println("Mensagem")
```

`fmt.Println` já adiciona quebra de linha.

### 4. Configuração golangci-lint

**Atualizado `.github/workflows/go-test.yml`:**
- Action: `v4` → `v6`
- Version: `v1.61` → `v1.62`

## 🔧 Como Evitar

```bash
# Antes de fazer commit
make lint

# Ou
golangci-lint run
```

## 📚 Regras dos Linters

- **errcheck**: Sempre verifique erros (ou ignore explicitamente com `_`)
- **revive**: Renomeie parâmetros não usados para `_`
- **govet**: Não use `\n` com `fmt.Println`
- **gofmt**: Mantenha código formatado
- **staticcheck**: Siga boas práticas do Go
