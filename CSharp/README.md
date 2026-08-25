# TabuaMare SDK — C#

SDK C# para a [API Tábua de Marés do Brasil](https://tabuamare.api.br) (v2).

## Instalação

```bash
dotnet add package TabuaMare.SDK
```

## Início rápido

```csharp
using TabuaMare;

using var client = new TabuaMareClient();

var states = await client.GetStatesAsync();
var harbor = await client.GetHarborAsync("pb01");
var tides  = await client.GetTideTableAsync("pb01", month: 1, days: [1, 2, 3]);
```

## Autenticação

Sem chave: **16 req/min** por IP.  
Com chave free: **64 req/min** + **32.000 req/mês**.

```csharp
using var client = new TabuaMareClient(new TabuaMareClientOptions
{
    ApiKey = "sua-chave-aqui",
});
```

Obtendo chave: [tabuamare.api.br](https://tabuamare.api.br)

## Referência da API

### Estados

```csharp
IReadOnlyList<string> states = await client.GetStatesAsync();
```

### Portos

```csharp
// Resumo dos portos de um estado
IReadOnlyList<HarborName> harbors = await client.GetHarborsByStateAsync("pb");

// Detalhes de um porto
Harbor harbor = await client.GetHarborAsync("pb01");

// Múltiplos portos
IReadOnlyList<Harbor> harbors = await client.GetHarborsAsync(["pb01", "pe01", "sc01"]);
```

### Tábua de Marés

```csharp
// Dias específicos de um mês
IReadOnlyList<TideTable> tides = await client.GetTideTableAsync("pb01", month: 3, days: [10, 11, 12]);

// Mês completo
IReadOnlyList<TideTable> tides = await client.GetTideTableForMonthAsync("pb01", month: 3);
```

### Porto mais próximo

```csharp
// Em qualquer estado
Harbor nearest = await client.GetNearestHarborAsync(lat: -7.115, lng: -34.864);

// Restrito a um estado
Harbor nearest = await client.GetNearestHarborByStateAsync("pb", lat: -7.115, lng: -34.864);
```

### Tábua por Geolocalização

```csharp
IReadOnlyList<TideTable> tides = await client.GetGeoTideTableAsync(
    lat: -7.115, lng: -34.864, state: "pb", month: 1, days: [1, 2, 3]);
```

### Uso da cota (requer ApiKey)

```csharp
UsageInfo usage = await client.GetUsageAsync();
Console.WriteLine($"Plano: {usage.Plan} | RPM: {usage.UsedRpm}/{usage.LimitRpm}");
```

## Tratamento de erros

| Exceção | Quando |
|---|---|
| `TabuaMareApiException` | API retorna 4xx/5xx com envelope de erro |
| `TabuaMareRateLimitException` | HTTP 429 — verifique `.RetryAfter` |
| `TabuaMareValidationException` | Parâmetro inválido (lado do cliente) |
| `InvalidOperationException` | `GetUsageAsync` sem `ApiKey` configurada |

```csharp
try
{
    var harbor = await client.GetHarborAsync("invalido");
}
catch (TabuaMareRateLimitException ex) when (ex.RetryAfter.HasValue)
{
    await Task.Delay(ex.RetryAfter.Value);
}
catch (TabuaMareApiException ex)
{
    Console.WriteLine($"Erro {ex.StatusCode}: {ex.Message}");
}
```

## Opções

```csharp
var options = new TabuaMareClientOptions
{
    ApiKey  = "chave",                          // opcional
    Timeout = TimeSpan.FromSeconds(10),         // padrão: 30s
    BaseUrl = "https://tabuamare.api.br/api/v2" // padrão
};
```

## Testes

```bash
# Unitários
dotnet test tests/TabuaMare.Tests/

# Integração (API real)
TABUAMARE_INTEGRATION=true dotnet test tests/TabuaMare.Tests/

# Com api_key
TABUAMARE_INTEGRATION=true TABUAMARE_API_KEY=<key> dotnet test tests/TabuaMare.Tests/
```

## Exemplo completo

```bash
dotnet run --project examples/Basic
```

## Licença

MIT
