package main

import (
	"context"
	"fmt"
	"log"
	"time"

	tabuamare "github.com/Ddiidev/sdks-tabua-mare/go"
)

func main() {
	// Criar cliente com configurações customizadas
	client := tabuamare.NewClient(
		tabuamare.WithTimeout(60 * time.Second),
	)

	ctx := context.Background()

	// Exemplo 1: Obter múltiplos portos
	fmt.Println("=== Múltiplos Portos ===")
	harbors, err := client.GetHarbors(ctx, 1, 2, 3)
	if err != nil {
		log.Fatalf("Erro ao obter portos: %v", err)
	}
	for _, harbor := range harbors {
		fmt.Printf("%s (%s)\n", harbor.HarborName, harbor.State)
	}
	fmt.Println()

	// Exemplo 2: Obter tábua de marés para o mês inteiro
	fmt.Println("=== Tábua de Marés do Mês ===")
	tides, err := client.GetTideTableForMonth(ctx, 1, 1)
	if err != nil {
		log.Fatalf("Erro ao obter tábua de marés: %v", err)
	}
	fmt.Printf("Total de registros: %d\n", len(tides))

	// Exemplo 3: Tratamento de erros
	fmt.Println("\n=== Tratamento de Erros ===")
	_, err = client.GetHarbor(ctx, -1)
	if err != nil {
		fmt.Printf("Erro esperado: %v\n", err)
	}

	_, err = client.GetTideTable(ctx, 1, 13, []int{1})
	if err != nil {
		fmt.Printf("Erro esperado: %v\n", err)
	}

	// Exemplo 4: Usando context com timeout
	fmt.Println("\n=== Context com Timeout ===")
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	states, err := client.GetStates(ctxTimeout)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}
	fmt.Printf("Estados obtidos com sucesso: %d\n", len(states))
}
