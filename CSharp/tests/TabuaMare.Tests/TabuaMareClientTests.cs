using System.Net;
using System.Net.Http.Headers;
using System.Text;
using TabuaMare;
using Xunit;

namespace TabuaMare.Tests;

// ---------------------------------------------------------------------------
// Fake HttpMessageHandler para simular respostas da API sem rede
// ---------------------------------------------------------------------------

internal sealed class FakeHandler : HttpMessageHandler
{
    private readonly Func<HttpRequestMessage, HttpResponseMessage> _handler;

    public FakeHandler(Func<HttpRequestMessage, HttpResponseMessage> handler)
        => _handler = handler;

    public FakeHandler(string json, HttpStatusCode status = HttpStatusCode.OK)
        : this(_ => Json(json, status)) { }

    protected override Task<HttpResponseMessage> SendAsync(
        HttpRequestMessage request, CancellationToken cancellationToken)
        => Task.FromResult(_handler(request));

    internal static HttpResponseMessage Json(string json, HttpStatusCode status = HttpStatusCode.OK)
    {
        var resp = new HttpResponseMessage(status)
        {
            Content = new StringContent(json, Encoding.UTF8, "application/json")
        };
        return resp;
    }
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

internal static class ClientFactory
{
    public static TabuaMareClient Create(FakeHandler handler, string? apiKey = null)
    {
        var http = new HttpClient(handler);
        return new TabuaMareClient(http, "https://fake.local/api/v2", apiKey);
    }
}

// ---------------------------------------------------------------------------
// GetStatesAsync
// ---------------------------------------------------------------------------

public sealed class GetStatesTests
{
    [Fact]
    public async Task ReturnsStateList()
    {
        const string json = """{"data":["pb","pe","sc"],"total":3,"error":null}""";
        using var client = ClientFactory.Create(new FakeHandler(json));

        var states = await client.GetStatesAsync();

        Assert.Equal(3, states.Count);
        Assert.Contains("pb", states);
    }

    [Fact]
    public async Task ReturnsEmptyOnNullData()
    {
        const string json = """{"data":null,"total":0,"error":null}""";
        using var client = ClientFactory.Create(new FakeHandler(json));

        var states = await client.GetStatesAsync();

        Assert.Empty(states);
    }
}

// ---------------------------------------------------------------------------
// GetHarborsByStateAsync
// ---------------------------------------------------------------------------

public sealed class GetHarborsByStateTests
{
    private const string Json = """
        {
          "data": [
            {"id":"pb01","year":2024,"harbor_name":"Porto de Cabedelo","data_collection_institution":"CHM"}
          ],
          "total":1,"error":null
        }
        """;

    [Fact]
    public async Task ReturnsHarborNames()
    {
        using var client = ClientFactory.Create(new FakeHandler(Json));

        var harbors = await client.GetHarborsByStateAsync("pb");

        Assert.Single(harbors);
        Assert.Equal("pb01", harbors[0].Id);
        Assert.Equal("Porto de Cabedelo", harbors[0].Name);
    }

    [Fact]
    public async Task ThrowsOnEmptyState()
    {
        using var client = ClientFactory.Create(new FakeHandler(Json));

        await Assert.ThrowsAsync<TabuaMareValidationException>(
            () => client.GetHarborsByStateAsync(""));
    }

    [Fact]
    public async Task UsesStateInPath()
    {
        string? capturedPath = null;
        var handler = new FakeHandler(req =>
        {
            capturedPath = req.RequestUri!.AbsolutePath;
            return FakeHandler.Json(Json);
        });
        using var client = ClientFactory.Create(handler);

        await client.GetHarborsByStateAsync("sc");

        Assert.Contains("/harbor_names/sc", capturedPath);
    }
}

// ---------------------------------------------------------------------------
// GetHarborsAsync / GetHarborAsync
// ---------------------------------------------------------------------------

public sealed class GetHarborsTests
{
    private const string SingleJson = """
        {
          "data": [
            {"id":"pb01","harbor_name":"Porto de Cabedelo","state":"PB","timezone":"UTC-3","card":"","geo_location":[],"mean_level":0}
          ],
          "total":1,"error":null
        }
        """;

    [Fact]
    public async Task SingleIdUsesSimplePath()
    {
        string? capturedPath = null;
        var handler = new FakeHandler(req =>
        {
            capturedPath = req.RequestUri!.AbsolutePath;
            return FakeHandler.Json(SingleJson);
        });
        using var client = ClientFactory.Create(handler);

        await client.GetHarborsAsync(["pb01"]);

        Assert.Contains("/harbors/pb01", capturedPath);
    }

    [Fact]
    public async Task MultipleIdsUseBracketPath()
    {
        string? capturedUrl = null;
        var handler = new FakeHandler(req =>
        {
            capturedUrl = req.RequestUri!.ToString();
            return FakeHandler.Json(SingleJson);
        });
        using var client = ClientFactory.Create(handler);

        await client.GetHarborsAsync(["pb01", "pe01"]);

        Assert.NotNull(capturedUrl);
        Assert.True(capturedUrl!.Contains("[pb01,pe01]") || capturedUrl.Contains("%5Bpb01%2Cpe01%5D"),
            $"URL deveria conter os IDs entre colchetes. Actual: {capturedUrl}");
    }

    [Fact]
    public async Task ThrowsOnEmptyIds()
    {
        using var client = ClientFactory.Create(new FakeHandler(SingleJson));

        await Assert.ThrowsAsync<TabuaMareValidationException>(
            () => client.GetHarborsAsync([]));
    }

    [Fact]
    public async Task GetHarborReturnsFirstResult()
    {
        using var client = ClientFactory.Create(new FakeHandler(SingleJson));

        var harbor = await client.GetHarborAsync("pb01");

        Assert.Equal("pb01", harbor.Id);
        Assert.Equal("PB", harbor.State);
    }
}

// ---------------------------------------------------------------------------
// GetTideTableAsync
// ---------------------------------------------------------------------------

public sealed class GetTideTableTests
{
    private const string Json = """
        {
          "data": [
            {
              "year":2024,"harbor_name":"Porto de Cabedelo","state":"PB",
              "timezone":"UTC-3","card":"","data_collection_institution":"CHM",
              "mean_level":0,
              "months":[{"month_name":"Janeiro","month":1,"days":[]}]
            }
          ],
          "total":1,"error":null
        }
        """;

    [Fact]
    public async Task ReturnsTideData()
    {
        using var client = ClientFactory.Create(new FakeHandler(Json));

        var result = await client.GetTideTableAsync("pb01", 1, [1, 2]);

        Assert.Single(result);
        Assert.Equal("Porto de Cabedelo", result[0].HarborName);
    }

    [Fact]
    public async Task ThrowsOnInvalidMonth()
    {
        using var client = ClientFactory.Create(new FakeHandler(Json));

        await Assert.ThrowsAsync<TabuaMareValidationException>(
            () => client.GetTideTableAsync("pb01", 13, [1]));
    }

    [Fact]
    public async Task ThrowsOnEmptyDays()
    {
        using var client = ClientFactory.Create(new FakeHandler(Json));

        await Assert.ThrowsAsync<TabuaMareValidationException>(
            () => client.GetTideTableAsync("pb01", 1, []));
    }

    [Fact]
    public async Task ThrowsOnEmptyHarborId()
    {
        using var client = ClientFactory.Create(new FakeHandler(Json));

        await Assert.ThrowsAsync<TabuaMareValidationException>(
            () => client.GetTideTableAsync("", 1, [1]));
    }
}

// ---------------------------------------------------------------------------
// Error handling
// ---------------------------------------------------------------------------

public sealed class ErrorHandlingTests
{
    [Fact]
    public async Task Throws429AsRateLimitException()
    {
        var handler = new FakeHandler(req =>
        {
            var resp = new HttpResponseMessage(HttpStatusCode.TooManyRequests);
            resp.Headers.RetryAfter = new RetryConditionHeaderValue(TimeSpan.FromSeconds(60));
            return resp;
        });
        using var client = ClientFactory.Create(handler);

        var ex = await Assert.ThrowsAsync<TabuaMareRateLimitException>(
            () => client.GetStatesAsync());

        Assert.NotNull(ex.RetryAfter);
        Assert.Equal(60, ex.RetryAfter!.Value.TotalSeconds);
    }

    [Fact]
    public async Task Throws429WithoutRetryAfterHeader()
    {
        var handler = new FakeHandler(_ => new HttpResponseMessage(HttpStatusCode.TooManyRequests));
        using var client = ClientFactory.Create(handler);

        var ex = await Assert.ThrowsAsync<TabuaMareRateLimitException>(
            () => client.GetStatesAsync());

        Assert.Null(ex.RetryAfter);
    }

    [Fact]
    public async Task ThrowsApiExceptionOnErrorEnvelope()
    {
        const string json = """{"data":null,"total":0,"error":{"code":404,"message":"harbor not found"}}""";
        var handler = new FakeHandler(json, HttpStatusCode.NotFound);
        using var client = ClientFactory.Create(handler);

        var ex = await Assert.ThrowsAsync<TabuaMareApiException>(
            () => client.GetStatesAsync());

        Assert.Equal(404, ex.StatusCode);
        Assert.Contains("harbor not found", ex.Message);
    }
}

// ---------------------------------------------------------------------------
// Auth headers
// ---------------------------------------------------------------------------

public sealed class AuthHeaderTests
{
    [Fact]
    public async Task SendsAuthHeadersWhenApiKeySet()
    {
        string? authHeader = null;
        string? xApiKeyHeader = null;

        var handler = new FakeHandler(req =>
        {
            authHeader = req.Headers.Authorization?.ToString();
            req.Headers.TryGetValues("X-Api-Key", out var vals);
            xApiKeyHeader = vals?.FirstOrDefault();
            return FakeHandler.Json("""{"data":[],"total":0,"error":null}""");
        });

        using var client = ClientFactory.Create(handler, apiKey: "test-key-123");
        await client.GetStatesAsync();

        Assert.Equal("Bearer test-key-123", authHeader);
        Assert.Equal("test-key-123", xApiKeyHeader);
    }

    [Fact]
    public async Task GetUsageThrowsWithoutApiKey()
    {
        using var client = ClientFactory.Create(
            new FakeHandler("""{"data":null,"total":0,"error":null}"""));

        await Assert.ThrowsAsync<InvalidOperationException>(
            () => client.GetUsageAsync());
    }
}

// ---------------------------------------------------------------------------
// Coordinate validation
// ---------------------------------------------------------------------------

public sealed class CoordinateValidationTests
{
    [Theory]
    [InlineData(-91, 0)]
    [InlineData(91, 0)]
    public async Task ThrowsOnInvalidLatitude(double lat, double lng)
    {
        using var client = ClientFactory.Create(
            new FakeHandler("""{"data":null,"total":0,"error":null}"""));

        await Assert.ThrowsAsync<TabuaMareValidationException>(
            () => client.GetNearestHarborAsync(lat, lng));
    }

    [Theory]
    [InlineData(0, -181)]
    [InlineData(0, 181)]
    public async Task ThrowsOnInvalidLongitude(double lat, double lng)
    {
        using var client = ClientFactory.Create(
            new FakeHandler("""{"data":null,"total":0,"error":null}"""));

        await Assert.ThrowsAsync<TabuaMareValidationException>(
            () => client.GetNearestHarborAsync(lat, lng));
    }

    [Fact]
    public async Task FormatsCoordinatesWithSixDecimalPlaces()
    {
        string? capturedPath = null;
        var json = """
            {"data":[{"id":"pb01","harbor_name":"Porto de Cabedelo","state":"PB","timezone":"UTC-3","card":"","geo_location":[],"mean_level":0}],"total":1,"error":null}
            """;
        var handler = new FakeHandler(req =>
        {
            capturedPath = Uri.UnescapeDataString(req.RequestUri!.AbsolutePath);
            return FakeHandler.Json(json);
        });
        using var client = ClientFactory.Create(handler);

        await client.GetNearestHarborAsync(-7.115090, -34.864000);

        Assert.NotNull(capturedPath);
        Assert.Contains("[-7.115090,-34.864000]", capturedPath);
    }
}
