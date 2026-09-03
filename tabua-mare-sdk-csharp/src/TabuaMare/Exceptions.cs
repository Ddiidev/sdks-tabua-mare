namespace TabuaMare;

/// <summary>Erro retornado pela API (HTTP 4xx/5xx).</summary>
public sealed class TabuaMareApiException : Exception
{
    /// <summary>Código HTTP da resposta.</summary>
    public int StatusCode { get; }

    /// <summary>Código de erro interno da API.</summary>
    public int ApiCode { get; }

    public TabuaMareApiException(int statusCode, int apiCode, string message)
        : base($"API error (HTTP {statusCode}, code {apiCode}): {message}")
    {
        StatusCode = statusCode;
        ApiCode = apiCode;
    }
}

/// <summary>
/// Limite de requisições excedido (HTTP 429).
/// Verifique <see cref="RetryAfter"/> para saber quando tentar novamente.
/// </summary>
public sealed class TabuaMareRateLimitException : Exception
{
    /// <summary>
    /// Tempo a aguardar antes de tentar novamente.
    /// <see langword="null"/> quando o header <c>Retry-After</c> não está presente.
    /// </summary>
    public TimeSpan? RetryAfter { get; }

    public TabuaMareRateLimitException(TimeSpan? retryAfter)
        : base(retryAfter.HasValue
            ? $"Rate limit exceeded: retry after {retryAfter.Value.TotalSeconds:0}s"
            : "Rate limit exceeded")
    {
        RetryAfter = retryAfter;
    }
}

/// <summary>Erro de validação de parâmetros no lado do cliente.</summary>
public sealed class TabuaMareValidationException : Exception
{
    /// <summary>Nome do campo inválido.</summary>
    public string Field { get; }

    public TabuaMareValidationException(string field, string message)
        : base($"Validation error on '{field}': {message}")
    {
        Field = field;
    }
}
