using TabuaMare;
using Xunit;

namespace TabuaMare.Tests;

/// <summary>
/// Testes de integração contra a API real.
/// Execute com: TABUAMARE_INTEGRATION=true dotnet test
/// Opcionalmente: TABUAMARE_API_KEY=&lt;key&gt; dotnet test
/// </summary>
public sealed class IntegrationTests
{
    private static bool Enabled =>
        Environment.GetEnvironmentVariable("TABUAMARE_INTEGRATION") == "true";

    private static TabuaMareClient BuildClient()
    {
        var apiKey = Environment.GetEnvironmentVariable("TABUAMARE_API_KEY");
        return new TabuaMareClient(new TabuaMareClientOptions { ApiKey = apiKey });
    }

    [Fact]
    public async Task GetStates_ReturnsAtLeastOneState()
    {
        if (!Enabled) return;
        using var client = BuildClient();
        var states = await client.GetStatesAsync();
        Assert.NotEmpty(states);
    }

    [Fact]
    public async Task GetHarborsByState_ReturnsHarbors()
    {
        if (!Enabled) return;
        using var client = BuildClient();
        var harbors = await client.GetHarborsByStateAsync("pb");
        Assert.NotEmpty(harbors);
    }

    [Fact]
    public async Task GetHarbor_ReturnsPb01()
    {
        if (!Enabled) return;
        using var client = BuildClient();
        var harbor = await client.GetHarborAsync("pb01");
        Assert.Equal("pb01", harbor.Id);
    }

    [Fact]
    public async Task GetHarbors_MultipleIds_ReturnsBoth()
    {
        if (!Enabled) return;
        using var client = BuildClient();
        var harbors = await client.GetHarborsAsync(["pb01", "pe01"]);
        Assert.Equal(2, harbors.Count);
    }

    [Fact]
    public async Task GetTideTable_ReturnsData()
    {
        if (!Enabled) return;
        using var client = BuildClient();
        var tides = await client.GetTideTableAsync("pb01", 1, [1]);
        Assert.NotEmpty(tides);
    }

    [Fact]
    public async Task GetNearestHarborByState_ReturnsPbHarbor()
    {
        if (!Enabled) return;
        using var client = BuildClient();
        var harbor = await client.GetNearestHarborByStateAsync("pb", -7.11509, -34.864);
        Assert.Equal("PB", harbor.State, ignoreCase: true);
    }

    [Fact]
    public async Task GetGeoTideTable_ReturnsData()
    {
        if (!Enabled) return;
        using var client = BuildClient();
        var tides = await client.GetGeoTideTableAsync(-7.11509, -34.864, "pb", 1, [1]);
        Assert.NotEmpty(tides);
    }

    [Fact]
    public async Task GetUsage_WithApiKey_ReturnsUsage()
    {
        if (!Enabled) return;
        var apiKey = Environment.GetEnvironmentVariable("TABUAMARE_API_KEY");
        if (string.IsNullOrWhiteSpace(apiKey)) return;

        using var client = new TabuaMareClient(new TabuaMareClientOptions { ApiKey = apiKey });
        var usage = await client.GetUsageAsync();
        Assert.NotNull(usage.Plan);
    }
}
