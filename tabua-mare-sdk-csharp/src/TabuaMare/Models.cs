using System.Text.Json.Serialization;

namespace TabuaMare;

/// <summary>Envelope padrão de resposta da API.</summary>
public sealed class ApiResponse<T>
{
    [JsonPropertyName("data")]
    public T? Data { get; init; }

    [JsonPropertyName("total")]
    public int Total { get; init; }

    [JsonPropertyName("error")]
    public ApiErrorDetail? Error { get; init; }
}

/// <summary>Detalhe de erro retornado pela API.</summary>
public sealed class ApiErrorDetail
{
    [JsonPropertyName("code")]
    public int Code { get; init; }

    [JsonPropertyName("message")]
    public string Message { get; init; } = string.Empty;
}

/// <summary>Informações básicas de um porto (endpoint harbor_names).</summary>
public sealed class HarborName
{
    [JsonPropertyName("id")]
    public string Id { get; init; } = string.Empty;

    [JsonPropertyName("year")]
    public int Year { get; init; }

    [JsonPropertyName("harbor_name")]
    public string Name { get; init; } = string.Empty;

    [JsonPropertyName("data_collection_institution")]
    public string DataCollectionInstitution { get; init; } = string.Empty;
}

/// <summary>Coordenadas geográficas de um porto.</summary>
public sealed class GeoLocation
{
    [JsonPropertyName("lat")]
    public string Lat { get; init; } = string.Empty;

    [JsonPropertyName("lng")]
    public string Lng { get; init; } = string.Empty;

    [JsonPropertyName("decimal_lat")]
    public string DecimalLat { get; init; } = string.Empty;

    [JsonPropertyName("decimal_lng")]
    public string DecimalLng { get; init; } = string.Empty;

    [JsonPropertyName("lat_direction")]
    public string LatDirection { get; init; } = string.Empty;

    [JsonPropertyName("lng_direction")]
    public string LngDirection { get; init; } = string.Empty;
}

/// <summary>Informações detalhadas de um porto.</summary>
public sealed class Harbor
{
    [JsonPropertyName("id")]
    public string Id { get; init; } = string.Empty;

    [JsonPropertyName("harbor_name")]
    public string Name { get; init; } = string.Empty;

    [JsonPropertyName("state")]
    public string State { get; init; } = string.Empty;

    [JsonPropertyName("timezone")]
    public string Timezone { get; init; } = string.Empty;

    [JsonPropertyName("card")]
    public string Card { get; init; } = string.Empty;

    [JsonPropertyName("geo_location")]
    public IReadOnlyList<GeoLocation> GeoLocation { get; init; } = [];

    [JsonPropertyName("mean_level")]
    public double MeanLevel { get; init; }
}

/// <summary>Nível de maré em um horário.</summary>
public sealed class TideHour
{
    [JsonPropertyName("hour")]
    public string Hour { get; init; } = string.Empty;

    [JsonPropertyName("level")]
    public double Level { get; init; }
}

/// <summary>Dados de maré de um dia.</summary>
public sealed class TideDay
{
    [JsonPropertyName("weekday_name")]
    public string WeekdayName { get; init; } = string.Empty;

    [JsonPropertyName("day")]
    public int Day { get; init; }

    [JsonPropertyName("hours")]
    public IReadOnlyList<TideHour> Hours { get; init; } = [];
}

/// <summary>Dados de maré de um mês.</summary>
public sealed class TideMonth
{
    [JsonPropertyName("month_name")]
    public string MonthName { get; init; } = string.Empty;

    [JsonPropertyName("month")]
    public int Month { get; init; }

    [JsonPropertyName("days")]
    public IReadOnlyList<TideDay> Days { get; init; } = [];
}

/// <summary>Tábua de marés completa de um porto.</summary>
public sealed class TideTable
{
    [JsonPropertyName("year")]
    public int Year { get; init; }

    [JsonPropertyName("harbor_name")]
    public string HarborName { get; init; } = string.Empty;

    [JsonPropertyName("state")]
    public string State { get; init; } = string.Empty;

    [JsonPropertyName("timezone")]
    public string Timezone { get; init; } = string.Empty;

    [JsonPropertyName("card")]
    public string Card { get; init; } = string.Empty;

    [JsonPropertyName("data_collection_institution")]
    public string DataCollectionInstitution { get; init; } = string.Empty;

    [JsonPropertyName("mean_level")]
    public double MeanLevel { get; init; }

    [JsonPropertyName("months")]
    public IReadOnlyList<TideMonth> Months { get; init; } = [];
}

/// <summary>
/// Consumo atual de rate-limit da api_key.
/// Campos numéricos chegam como string da API; "-1" indica limite ilimitado.
/// </summary>
public sealed class UsageInfo
{
    [JsonPropertyName("plan")]
    public string Plan { get; init; } = string.Empty;

    [JsonPropertyName("limit_rpm")]
    public string LimitRpm { get; init; } = string.Empty;

    [JsonPropertyName("used_rpm")]
    public string UsedRpm { get; init; } = string.Empty;

    [JsonPropertyName("remaining_rpm")]
    public string RemainingRpm { get; init; } = string.Empty;

    [JsonPropertyName("limit_monthly")]
    public string LimitMonthly { get; init; } = string.Empty;

    [JsonPropertyName("used_monthly")]
    public string UsedMonthly { get; init; } = string.Empty;

    [JsonPropertyName("remaining_monthly")]
    public string RemainingMonthly { get; init; } = string.Empty;
}
