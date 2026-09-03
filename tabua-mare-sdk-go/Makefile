.PHONY: help test test-integration test-all lint fmt vet build clean examples

help: ## Mostra esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

test: ## Executa testes unitários
	go test -v -race -timeout 30s ./...

test-integration: ## Executa testes de integração
	go test -v -race -timeout 60s -tags=integration ./...

test-all: ## Executa todos os testes
	go test -v -race -timeout 60s -tags=integration ./...

test-coverage: ## Executa testes com cobertura
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Relatório de cobertura gerado em coverage.html"

lint: ## Executa linter
	golangci-lint run --timeout=5m

fmt: ## Formata o código
	go fmt ./...
	gofmt -s -w .

vet: ## Executa go vet
	go vet ./...

build: ## Compila o pacote
	go build -v ./...

build-examples: ## Compila os exemplos
	go build -v ./examples/basic
	go build -v ./examples/advanced
	go build -v ./cmd/test

run-basic: ## Executa exemplo básico
	go run examples/basic/main.go

run-advanced: ## Executa exemplo avançado
	go run examples/advanced/main.go

run-test: ## Executa teste manual
	go run cmd/test/main.go

clean: ## Remove arquivos gerados
	rm -f coverage.out coverage.html
	rm -f examples/basic/basic
	rm -f examples/advanced/advanced
	rm -f cmd/test/test
	go clean -cache -testcache

deps: ## Baixa dependências
	go mod download
	go mod verify

tidy: ## Organiza dependências
	go mod tidy

check: fmt vet lint test ## Executa todas as verificações

ci: check test-integration ## Executa pipeline de CI

install-tools: ## Instala ferramentas de desenvolvimento
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.DEFAULT_GOAL := help
