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
		tabuamare.WithTimeout(60*time.Second),
		tabuamare.WithAPIKey("SUA_API_KEY"), // opcional: aumenta o rate limit
	)

	ctx := context.Background()

	// Exemplo 1: Obter múltiplos portos
	fmt.Println("=== Múltiplos Portos ===")
	harbors, err := client.GetHarbors(ctx, "pb01", "pe01")
	if err != nil {
		log.Fatalf("Erro ao obter portos: %v", err)
	}
	for _, harbor := range harbors {
		fmt.Printf("%s (%s)\n", harbor.HarborName, harbor.State)
	}
	fmt.Println()

	// Exemplo 2: Obter tábua de marés para o mês inteiro
	fmt.Println("=== Tábua de Marés do Mês ===")
	tides, err := client.GetTideTableForMonth(ctx, "pb01", 1)
	if err != nil {
		log.Fatalf("Erro ao obter tábua de marés: %v", err)
	}
	fmt.Printf("Total de registros: %d\n", len(tides))

	// Exemplo 2.1: Tábua de marés por geolocalização
	fmt.Println("=== Tábua de Marés por Geolocalização ===")
	geoTides, err := client.GetGeoTideTable(ctx, -7.11509, -34.864, "pb", 1, []int{1})
	if err != nil {
		log.Fatalf("Erro ao obter tábua de marés por geolocalização: %v", err)
	}
	fmt.Printf("Porto mais próximo: %s\n", geoTides[0].HarborName)

	// Exemplo 3: Tratamento de erros
	fmt.Println("\n=== Tratamento de Erros ===")
	_, err = client.GetHarbor(ctx, "")
	if err != nil {
		fmt.Printf("Erro esperado: %v\n", err)
	}

	_, err = client.GetTideTable(ctx, "pb01", 13, []int{1})
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

	// Exemplo 5: Consulta de uso da cota (requer api_key válida)
	fmt.Println("\n=== Uso da Cota ===")
	if usage, err := client.GetUsage(ctx); err != nil {
		fmt.Printf("Configure uma api_key válida para consultar o uso: %v\n", err)
	} else {
		fmt.Printf("Plano: %s | RPM usado: %s/%s | Mensal usado: %s/%s\n",
			usage.Plan, usage.UsedRPM, usage.LimitRPM, usage.UsedMonthly, usage.LimitMonthly)
	}
}
