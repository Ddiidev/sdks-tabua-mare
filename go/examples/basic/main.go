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

	// Exemplo 1: Listar estados
	fmt.Println("=== Listando Estados ===")
	states, err := client.GetStates(ctx)
	if err != nil {
		log.Fatalf("Erro ao listar estados: %v", err)
	}
	fmt.Printf("Total de estados: %d\n", len(states))
	fmt.Printf("Estados: %v\n\n", states)

	// Exemplo 2: Listar portos de Santa Catarina
	fmt.Println("=== Listando Portos de SC ===")
	harbors, err := client.GetHarborNames(ctx, "sc")
	if err != nil {
		log.Fatalf("Erro ao listar portos: %v", err)
	}
	for _, harbor := range harbors {
		fmt.Printf("ID: %d - %s\n", harbor.ID, harbor.HarborName)
	}
	fmt.Println()

	// Exemplo 3: Obter detalhes de um porto
	fmt.Println("=== Detalhes do Porto ===")
	harbor, err := client.GetHarbor(ctx, 1)
	if err != nil {
		log.Fatalf("Erro ao obter porto: %v", err)
	}
	fmt.Printf("Nome: %s\n", harbor.HarborName)
	fmt.Printf("Estado: %s\n", harbor.State)
	fmt.Printf("Nível médio: %.2f m\n\n", harbor.MeanLevel)

	// Exemplo 4: Obter tábua de marés para dias específicos
	fmt.Println("=== Tábua de Marés ===")
	tides, err := client.GetTideTable(ctx, 1, 1, []int{1, 2, 3})
	if err != nil {
		log.Fatalf("Erro ao obter tábua de marés: %v", err)
	}

	for _, tide := range tides {
		fmt.Printf("Porto: %s\n", tide.HarborName)
		for _, month := range tide.Months {
			fmt.Printf("Mês: %s\n", month.MonthName)
			for _, day := range month.Days {
				fmt.Printf("  Dia %d (%s):\n", day.Day, day.WeekdayName)
				for _, hour := range day.Hours {
					fmt.Printf("    %s - %.2f m\n", hour.Hour, hour.Level)
				}
			}
		}
	}
}
