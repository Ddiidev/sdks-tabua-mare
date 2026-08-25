package main

import (
	"context"
	"fmt"
	"log"
	"os"
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
		fmt.Printf("   - [%s] %s\n", h.ID, h.HarborName)
	}
	fmt.Println()

	// Teste 3: Detalhes de um Porto
	fmt.Println("🏖️  Teste 3: Obtendo Detalhes do Porto pb01")
	harbor, err := client.GetHarbor(ctx, "pb01")
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
	tides, err := client.GetTideTable(ctx, "pb01", 1, []int{1, 2, 3})
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
	fmt.Println("🔢 Teste 5: Obtendo Múltiplos Portos (IDs pb01, pe01)")
	multiHarbors, err := client.GetHarbors(ctx, "pb01", "pe01")
	if err != nil {
		log.Fatalf("❌ Erro: %v", err)
	}
	fmt.Printf("✅ Sucesso! Obtidos %d portos\n", len(multiHarbors))
	for _, h := range multiHarbors {
		fmt.Printf("   - %s (%s)\n", h.HarborName, h.State)
	}
	fmt.Println()

	// Teste 6: Porto mais próximo por estado
	fmt.Println("📍 Teste 6: Obtendo Porto Mais Próximo (PB, coordenadas de João Pessoa)")
	nearest, err := client.GetNearestHarborByState(ctx, "pb", -7.11509, -34.864)
	if err != nil {
		log.Fatalf("❌ Erro: %v", err)
	}
	fmt.Printf("✅ Sucesso! Porto mais próximo: %s (%s)\n\n", nearest.HarborName, nearest.ID)

	// Teste 7: Tábua de Marés por Geolocalização
	fmt.Println("🌍 Teste 7: Obtendo Tábua de Marés por Geolocalização")
	geoTides, err := client.GetGeoTideTable(ctx, -7.11509, -34.864, "pb", 1, []int{1})
	if err != nil {
		log.Fatalf("❌ Erro: %v", err)
	}
	fmt.Printf("✅ Sucesso! Porto: %s\n\n", geoTides[0].HarborName)

	// Teste 8: Uso da Cota (requer api_key via variável de ambiente TABUAMARE_API_KEY)
	if apiKey := os.Getenv("TABUAMARE_API_KEY"); apiKey != "" {
		fmt.Println("📊 Teste 8: Consultando Uso da Cota")
		authedClient := tabuamare.NewClient(tabuamare.WithAPIKey(apiKey))
		usage, err := authedClient.GetUsage(ctx)
		if err != nil {
			fmt.Printf("⚠️  Não foi possível consultar uso: %v\n\n", err)
		} else {
			fmt.Printf("✅ Plano: %s | RPM: %s/%s | Mensal: %s/%s\n\n",
				usage.Plan, usage.UsedRPM, usage.LimitRPM, usage.UsedMonthly, usage.LimitMonthly)
		}
	}

	// Teste 9: Validação de Erros
	fmt.Println("⚠️  Teste 9: Testando Validação de Erros")

	_, err = client.GetTideTable(ctx, "", 1, []int{1})
	if err != nil {
		fmt.Printf("✅ Erro esperado capturado: %v\n", err)
	}

	_, err = client.GetTideTable(ctx, "pb01", 13, []int{1})
	if err != nil {
		fmt.Printf("✅ Erro esperado capturado: %v\n", err)
	}

	_, err = client.GetHarborNames(ctx, "")
	if err != nil {
		fmt.Printf("✅ Erro esperado capturado: %v\n", err)
	}

	fmt.Println("\n🎉 Todos os testes concluídos com sucesso!")
}
