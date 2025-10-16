package main

import (
	"context"
	"fmt"
	"log"
	"time"

	tabuamare "github.com/Ddiidev/sdks-tabua-mare/go"
)

func main() {
	fmt.Println("🌊 Testando Tabua Mare SDK")

	client := tabuamare.NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Teste 1: Listar Estados
	fmt.Println("📍 Teste 1: Listando Estados")
	states, err := client.GetStates(ctx)
	if err != nil {
		log.Fatalf("❌ Erro: %v", err)
	}
	fmt.Printf("✅ Sucesso! Encontrados %d estados\n", len(states))
	fmt.Printf("   Estados: %v\n", states)
	fmt.Println()

	// Teste 2: Listar Portos de SC
	fmt.Println("⚓ Teste 2: Listando Portos de Santa Catarina")
	harbors, err := client.GetHarborNames(ctx, "sc")
	if err != nil {
		log.Fatalf("❌ Erro: %v", err)
	}
	fmt.Printf("✅ Sucesso! Encontrados %d portos\n", len(harbors))
	for _, h := range harbors {
		fmt.Printf("   - [%d] %s\n", h.ID, h.HarborName)
	}
	fmt.Println()

	// Teste 3: Detalhes de um Porto
	fmt.Println("🏖️  Teste 3: Obtendo Detalhes do Porto ID 1")
	harbor, err := client.GetHarbor(ctx, 1)
	if err != nil {
		log.Fatalf("❌ Erro: %v", err)
	}
	fmt.Printf("✅ Sucesso!\n")
	fmt.Printf("   Nome: %s\n", harbor.HarborName)
	fmt.Printf("   Estado: %s\n", harbor.State)
	fmt.Printf("   Timezone: %s\n", harbor.Timezone)
	fmt.Printf("   Nível Médio: %.2f m\n", harbor.MeanLevel)
	if len(harbor.GeoLocation) > 0 {
		fmt.Printf("   Localização: %s, %s\n", harbor.GeoLocation[0].DecimalLat, harbor.GeoLocation[0].DecimalLng)
	}
	fmt.Println()

	// Teste 4: Tábua de Marés
	fmt.Println("📊 Teste 4: Obtendo Tábua de Marés (Janeiro, dias 1-3)")
	tides, err := client.GetTideTable(ctx, 1, 1, []int{1, 2, 3})
	if err != nil {
		log.Fatalf("❌ Erro: %v", err)
	}
	fmt.Printf("✅ Sucesso!\n")
	for _, tide := range tides {
		fmt.Printf("   Porto: %s\n", tide.HarborName)
		for _, month := range tide.Months {
			fmt.Printf("   Mês: %s\n", month.MonthName)
			for _, day := range month.Days {
				fmt.Printf("     📅 Dia %d (%s):\n", day.Day, day.WeekdayName)
				for _, hour := range day.Hours {
					fmt.Printf("        🕐 %s - %.2f m\n", hour.Hour, hour.Level)
				}
			}
		}
	}
	fmt.Println()

	// Teste 5: Múltiplos Portos
	fmt.Println("🔢 Teste 5: Obtendo Múltiplos Portos (IDs 1, 2, 3)")
	multiHarbors, err := client.GetHarbors(ctx, 1, 2, 3)
	if err != nil {
		log.Fatalf("❌ Erro: %v", err)
	}
	fmt.Printf("✅ Sucesso! Obtidos %d portos\n", len(multiHarbors))
	for _, h := range multiHarbors {
		fmt.Printf("   - %s (%s)\n", h.HarborName, h.State)
	}
	fmt.Println()

	// Teste 6: Validação de Erros
	fmt.Println("⚠️  Teste 6: Testando Validação de Erros")

	_, err = client.GetTideTable(ctx, -1, 1, []int{1})
	if err != nil {
		fmt.Printf("✅ Erro esperado capturado: %v\n", err)
	}

	_, err = client.GetTideTable(ctx, 1, 13, []int{1})
	if err != nil {
		fmt.Printf("✅ Erro esperado capturado: %v\n", err)
	}

	_, err = client.GetHarborNames(ctx, "")
	if err != nil {
		fmt.Printf("✅ Erro esperado capturado: %v\n", err)
	}

	fmt.Println("\n🎉 Todos os testes concluídos com sucesso!")
}
