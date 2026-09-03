# Tábua de Maré — cliente C#

Cliente .NET para consultar a API REST de Tábua de Maré do Brasil (v2). A biblioteca encapsula as requisições HTTP, desserializa as respostas JSON e oferece métodos tipados para estados, portos, tábuas mensais e consultas por geolocalização.

## Requisitos

- .NET 9.0 ou superior
- Uma API key da [Tábua de Maré](https://tabuamare.api.br/dashboard) para endpoints autenticados e limites ampliados

## Instalação

```bash
dotnet add package TabuaMare.SDK
```

Ou adicione diretamente ao `.csproj`:

```xml
<PackageReference Include="TabuaMare.SDK" Version="2.0.0" />
```

## Autenticação

A chave é informada na configuração do cliente. Passe somente o valor da chave; não inclua o prefixo `Bearer`.

```csharp
using TabuaMare;

var client = new TabuaMareClient(new TabuaMareClientOptions
{
    ApiKey = Environment.GetEnvironmentVariable("TABUAMARE_API_KEY")
});
```

O cliente envia a chave nos headers `Authorization: Bearer <api_key>` e `X-Api-Key: <api_key>`. A explicação completa está na [documentação sobre API key](https://tabuamare.api.br/docs#api-key-header).

## Exemplos

### Listar estados e portos

```csharp
using TabuaMare;

var client = new TabuaMareClient();
var states = await client.GetStatesAsync();
var harbors = await client.GetHarborsByStateAsync("pb");

Console.WriteLine(string.Join(", ", states));
foreach (var harbor in harbors)
    Console.WriteLine($"{harbor.Id}: {harbor.Name}");
```

### Consultar a tábua de um porto

```csharp
var tables = await client.GetTideTableAsync(
    harborId: "pb01",
    month: 1,
    days: new[] { 1, 2, 3 });

foreach (var month in tables[0].Months)
foreach (var day in month.Days)
    Console.WriteLine($"Dia {day.Day}: {day.Hours.Count} eventos");
```

### Encontrar o porto mais próximo

```csharp
var nearest = await client.GetNearestHarborAsync(-7.11509, -34.86414);
Console.WriteLine($"{nearest.Name} ({nearest.Id})");
```

As respostas seguem os tipos retornados pela API v2. Os IDs de porto são strings, como `pb01` e `pe02`.

## Desenvolvimento

```bash
dotnet test tests/TabuaMare.Tests/TabuaMare.Tests.csproj
dotnet pack src/TabuaMare/TabuaMare.csproj --configuration Release
```

## Links

- [Documentação da API](https://tabuamare.api.br/docs)
- [Dashboard e API keys](https://tabuamare.api.br/dashboard)
- [Repositório](https://github.com/Ddiidev/sdks-tabua-mare)

## Licença

MIT.
