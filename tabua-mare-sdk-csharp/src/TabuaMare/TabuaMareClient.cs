using System.Globalization;
using System.Net;
using System.Net.Http.Headers;
using System.Text.Json;

namespace TabuaMare;

/// <summary>
/// Cliente principal para a API Tábua de Marés v2.
/// Instancie uma vez e reutilize (thread-safe).
/// </summary>
public sealed class TabuaMareClient : IDisposable
{
    private readonly HttpClient _http;
    private readonly string _baseUrl;
    private readonly string? _apiKey;
    private bool _disposed;

    private static readonly JsonSerializerOptions JsonOpts = new()
    {
        PropertyNameCaseInsensitive = true,
    };

    /// <param name="options">Opções de configuração.</param>
    public TabuaMareClient(TabuaMareClientOptions? options = null)
    {
        options ??= new TabuaMareClientOptions();
        _baseUrl = options.BaseUrl.TrimEnd('/');
        _apiKey = options.ApiKey;

        _http = new HttpClient { Timeout = options.Timeout };
        _http.DefaultRequestHeaders.Accept.Add(
            new MediaTypeWithQualityHeaderValue("application/json"));

        if (!string.IsNullOrWhiteSpace(_apiKey))
        {
            _http.DefaultRequestHeaders.Authorization =
                new AuthenticationHeaderValue("Bearer", _apiKey);
            _http.DefaultRequestHeaders.TryAddWithoutValidation("X-Api-Key", _apiKey);
        }
    }

    internal TabuaMareClient(HttpClient httpClient, string baseUrl, string? apiKey = null)
    {
        _http = httpClient;
        _baseUrl = baseUrl.TrimEnd('/');
        _apiKey = apiKey;

        _http.DefaultRequestHeaders.Accept.Add(
            new MediaTypeWithQualityHeaderValue("application/json"));

        if (!string.IsNullOrWhiteSpace(_apiKey))
        {
            _http.DefaultRequestHeaders.Authorization =
                new AuthenticationHeaderValue("Bearer", _apiKey);
            _http.DefaultRequestHeaders.TryAddWithoutValidation("X-Api-Key", _apiKey);
        }
    }

    // -------------------------------------------------------------------------
    // States
    // -------------------------------------------------------------------------

    /// <summary>Lista todos os estados com portos disponíveis.</summary>
    public async Task<IReadOnlyList<string>> GetStatesAsync(
        CancellationToken cancellationToken = default)
    {
        var resp = await GetAsync<ApiResponse<List<string>>>(
            "/states", cancellationToken).ConfigureAwait(false);
        return resp.Data ?? [];
    }

    // -------------------------------------------------------------------------
    // Harbor Names
    // -------------------------------------------------------------------------

    /// <summary>Lista os portos de um estado (formato resumido).</summary>
    /// <param name="state">Sigla do estado (ex: <c>"pb"</c>).</param>
    public async Task<IReadOnlyList<HarborName>> GetHarborsByStateAsync(
        string state,
        CancellationToken cancellationToken = default)
    {
        ValidateNotEmpty(state, nameof(state));
        var resp = await GetAsync<ApiResponse<List<HarborName>>>(
            $"/harbor_names/{Uri.EscapeDataString(state)}", cancellationToken).ConfigureAwait(false);
        return resp.Data ?? [];
    }

    // -------------------------------------------------------------------------
    // Harbors
    // -------------------------------------------------------------------------

    /// <summary>Obtém detalhes de um ou mais portos por ID.</summary>
    /// <param name="ids">Um ou mais IDs (ex: <c>"pb01"</c>, <c>"pe01"</c>).</param>
    public async Task<IReadOnlyList<Harbor>> GetHarborsAsync(
        IEnumerable<string> ids,
        CancellationToken cancellationToken = default)
    {
        var list = ids.ToList();
        if (list.Count == 0)
            throw new TabuaMareValidationException(nameof(ids), "At least one harbor ID is required.");

        var path = list.Count == 1
            ? $"/harbors/{Uri.EscapeDataString(list[0])}"
            : $"/harbors/{Uri.EscapeDataString("[" + string.Join(",", list) + "]")}";

        var resp = await GetAsync<ApiResponse<List<Harbor>>>(
            path, cancellationToken).ConfigureAwait(false);
        return resp.Data ?? [];
    }

    /// <summary>Obtém detalhes de um único porto.</summary>
    public async Task<Harbor> GetHarborAsync(
        string id,
        CancellationToken cancellationToken = default)
    {
        ValidateNotEmpty(id, nameof(id));
        var harbors = await GetHarborsAsync([id], cancellationToken).ConfigureAwait(false);
        return harbors.FirstOrDefault()
            ?? throw new TabuaMareApiException(404, 0, $"Harbor '{id}' not found.");
    }

    // -------------------------------------------------------------------------
    // Tide tables
    // -------------------------------------------------------------------------

    /// <summary>Obtém a tábua de marés de um porto para dias específicos de um mês.</summary>
    /// <param name="harborId">ID do porto (ex: <c>"pb01"</c>).</param>
    /// <param name="month">Mês (1–12).</param>
    /// <param name="days">Dias do mês desejados.</param>
    public async Task<IReadOnlyList<TideTable>> GetTideTableAsync(
        string harborId,
        int month,
        IEnumerable<int> days,
        CancellationToken cancellationToken = default)
    {
        ValidateNotEmpty(harborId, nameof(harborId));
        ValidateMonth(month);
        var dayList = days.ToList();
        if (dayList.Count == 0)
            throw new TabuaMareValidationException(nameof(days), "At least one day is required.");

        var daysParam = string.Join(",", dayList);
        var path = $"/tabua-mare/{Uri.EscapeDataString(harborId)}/{month}/{daysParam}";

        var resp = await GetAsync<ApiResponse<List<TideTable>>>(
            path, cancellationToken).ConfigureAwait(false);
        return resp.Data ?? [];
    }

    /// <summary>Obtém a tábua de marés de um porto para um mês completo.</summary>
    public async Task<IReadOnlyList<TideTable>> GetTideTableForMonthAsync(
        string harborId,
        int month,
        CancellationToken cancellationToken = default)
    {
        ValidateNotEmpty(harborId, nameof(harborId));
        ValidateMonth(month);

        var path = $"/tabua-mare/{Uri.EscapeDataString(harborId)}/{month}";
        var resp = await GetAsync<ApiResponse<List<TideTable>>>(
            path, cancellationToken).ConfigureAwait(false);
        return resp.Data ?? [];
    }

    // -------------------------------------------------------------------------
    // Nearest harbor
    // -------------------------------------------------------------------------

    /// <summary>Obtém o porto mais próximo de uma coordenada geográfica.</summary>
    /// <param name="lat">Latitude (-90 a 90).</param>
    /// <param name="lng">Longitude (-180 a 180).</param>
    public async Task<Harbor> GetNearestHarborAsync(
        double lat,
        double lng,
        CancellationToken cancellationToken = default)
    {
        ValidateCoordinates(lat, lng);
        var coord = FormatCoord(lat, lng);
        var path = $"/nearest-harbor/{Uri.EscapeDataString(coord)}";
        var resp = await GetAsync<ApiResponse<List<Harbor>>>(
            path, cancellationToken).ConfigureAwait(false);
        return resp.Data?.FirstOrDefault()
            ?? throw new TabuaMareApiException(404, 0, "No nearest harbor found.");
    }

    /// <summary>Obtém o porto mais próximo de uma coordenada dentro de um estado.</summary>
    /// <param name="state">Sigla do estado (ex: <c>"pb"</c>).</param>
    /// <param name="lat">Latitude (-90 a 90).</param>
    /// <param name="lng">Longitude (-180 a 180).</param>
    public async Task<Harbor> GetNearestHarborByStateAsync(
        string state,
        double lat,
        double lng,
        CancellationToken cancellationToken = default)
    {
        ValidateNotEmpty(state, nameof(state));
        ValidateCoordinates(lat, lng);
        var coord = FormatCoord(lat, lng);
        var path = $"/nearest-harbor/{Uri.EscapeDataString(state)}/{Uri.EscapeDataString(coord)}";
        var resp = await GetAsync<ApiResponse<List<Harbor>>>(
            path, cancellationToken).ConfigureAwait(false);
        return resp.Data?.FirstOrDefault()
            ?? throw new TabuaMareApiException(404, 0, "No nearest harbor found.");
    }

    // -------------------------------------------------------------------------
    // Geo tide table
    // -------------------------------------------------------------------------

    /// <summary>
    /// Obtém a tábua de marés pelo porto geograficamente mais próximo das coordenadas.
    /// </summary>
    /// <param name="lat">Latitude.</param>
    /// <param name="lng">Longitude.</param>
    /// <param name="state">Sigla do estado.</param>
    /// <param name="month">Mês (1–12).</param>
    /// <param name="days">Dias desejados.</param>
    public async Task<IReadOnlyList<TideTable>> GetGeoTideTableAsync(
        double lat,
        double lng,
        string state,
        int month,
        IEnumerable<int> days,
        CancellationToken cancellationToken = default)
    {
        ValidateCoordinates(lat, lng);
        ValidateNotEmpty(state, nameof(state));
        ValidateMonth(month);
        var dayList = days.ToList();
        if (dayList.Count == 0)
            throw new TabuaMareValidationException(nameof(days), "At least one day is required.");

        var coord = FormatCoord(lat, lng);
        var daysParam = string.Join(",", dayList);
        var path = $"/geo-tabua-mare/{Uri.EscapeDataString(coord)}/{Uri.EscapeDataString(state)}/{month}/{daysParam}";

        var resp = await GetAsync<ApiResponse<List<TideTable>>>(
            path, cancellationToken).ConfigureAwait(false);
        return resp.Data ?? [];
    }

    // -------------------------------------------------------------------------
    // Usage
    // -------------------------------------------------------------------------

    /// <summary>
    /// Consulta o consumo atual da api_key. Requer que <see cref="TabuaMareClientOptions.ApiKey"/>
    /// esteja configurada.
    /// </summary>
    public async Task<UsageInfo> GetUsageAsync(CancellationToken cancellationToken = default)
    {
        if (string.IsNullOrWhiteSpace(_apiKey))
            throw new InvalidOperationException(
                "GetUsageAsync requires an ApiKey. Set TabuaMareClientOptions.ApiKey.");

        var resp = await GetAsync<ApiResponse<UsageInfo>>(
            "/usage", cancellationToken).ConfigureAwait(false);
        return resp.Data
            ?? throw new TabuaMareApiException(0, 0, "Empty usage response.");
    }

    // -------------------------------------------------------------------------
    // Core HTTP
    // -------------------------------------------------------------------------

    private async Task<T> GetAsync<T>(string path, CancellationToken ct)
    {
        var url = _baseUrl + path;
        using var response = await _http.GetAsync(url, ct).ConfigureAwait(false);

        if (response.StatusCode == HttpStatusCode.TooManyRequests)
        {
            TimeSpan? retryAfter = null;
            if (response.Headers.RetryAfter?.Delta is { } delta)
                retryAfter = delta;
            throw new TabuaMareRateLimitException(retryAfter);
        }

        var body = await response.Content.ReadAsStringAsync(ct).ConfigureAwait(false);

        T result;
        try
        {
            result = JsonSerializer.Deserialize<T>(body, JsonOpts)
                     ?? throw new InvalidOperationException("Empty response body.");
        }
        catch (JsonException ex)
        {
            throw new TabuaMareApiException(
                (int)response.StatusCode, 0,
                $"Failed to deserialize response: {ex.Message}. Body: {body}");
        }

        if (result is ApiResponse<object> { Error: { Code: > 0 } err })
            throw new TabuaMareApiException((int)response.StatusCode, err.Code, err.Message);

        if (!response.IsSuccessStatusCode)
        {
            var errDetail = TryExtractError(body);
            throw new TabuaMareApiException(
                (int)response.StatusCode,
                errDetail?.Code ?? 0,
                errDetail?.Message ?? response.ReasonPhrase ?? "Unknown error");
        }

        return result;
    }

    private static ApiErrorDetail? TryExtractError(string body)
    {
        try
        {
            using var doc = JsonDocument.Parse(body);
            if (doc.RootElement.TryGetProperty("error", out var errorEl))
            {
                return JsonSerializer.Deserialize<ApiErrorDetail>(errorEl.GetRawText(), JsonOpts);
            }
        }
        catch { }
        return null;
    }

    // -------------------------------------------------------------------------
    // Helpers / validation
    // -------------------------------------------------------------------------

    private static void ValidateNotEmpty(string value, string paramName)
    {
        if (string.IsNullOrWhiteSpace(value))
            throw new TabuaMareValidationException(paramName, "Value cannot be empty.");
    }

    private static void ValidateMonth(int month)
    {
        if (month < 1 || month > 12)
            throw new TabuaMareValidationException(nameof(month), "Month must be between 1 and 12.");
    }

    private static void ValidateCoordinates(double lat, double lng)
    {
        if (lat < -90 || lat > 90)
            throw new TabuaMareValidationException(nameof(lat), "Latitude must be between -90 and 90.");
        if (lng < -180 || lng > 180)
            throw new TabuaMareValidationException(nameof(lng), "Longitude must be between -180 and 180.");
    }

    private static string FormatCoord(double lat, double lng)
        => string.Create(CultureInfo.InvariantCulture, $"[{lat:F6},{lng:F6}]");

    // -------------------------------------------------------------------------
    // IDisposable
    // -------------------------------------------------------------------------

    public void Dispose()
    {
        if (_disposed) return;
        _disposed = true;
        _http.Dispose();
    }
}
