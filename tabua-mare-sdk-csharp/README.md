# Tábua de Maré — cliente C#

Cliente .NET para a API REST de Tábua de Maré do Brasil.

## Instalação

```bash
dotnet add package TabuaMare.SDK
```

## Autenticação

Passe somente o valor da chave em `ApiKey`; o cliente envia `Authorization: Bearer <api_key>` e `X-Api-Key: <api_key>`. Não inclua `Bearer` no valor.

Consulte a explicação completa na [documentação da API](https://tabuamare.api.br/docs#api-key-header).

## Desenvolvimento

```bash
dotnet test tests/TabuaMare.Tests/TabuaMare.Tests.csproj
dotnet pack src/TabuaMare/TabuaMare.csproj --configuration Release
```
