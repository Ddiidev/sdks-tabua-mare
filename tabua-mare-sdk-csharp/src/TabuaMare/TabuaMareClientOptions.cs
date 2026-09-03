namespace TabuaMare;

/// <summary>Opções de configuração do <see cref="TabuaMareClient"/>.</summary>
public sealed class TabuaMareClientOptions
{
    /// <summary>URL base da API. Padrão: <c>https://tabuamare.api.br/api/v2</c>.</summary>
    public string BaseUrl { get; set; } = "https://tabuamare.api.br/api/v2";

    /// <summary>
    /// API key opcional. Quando definida, envia os headers
    /// <c>Authorization: Bearer &lt;key&gt;</c> e <c>X-Api-Key: &lt;key&gt;</c>.
    /// Sem chave: 16 req/min por IP. Com chave free: 64 req/min + 32.000/mês.
    /// </summary>
    public string? ApiKey { get; set; }

    /// <summary>Timeout das requisições. Padrão: 30 segundos.</summary>
    public TimeSpan Timeout { get; set; } = TimeSpan.FromSeconds(30);
}
