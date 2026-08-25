using TabuaMare;

var options = new TabuaMareClientOptions
{
    ApiKey = Environment.GetEnvironmentVariable("TABUAMARE_API_KEY"),
};

using var client = new TabuaMareClient(options);

Console.WriteLine("=== Estados disponíveis ===");
var states = await client.GetStatesAsync();
Console.WriteLine(string.Join(", ", states));

Console.WriteLine("\n=== Portos da Paraíba ===");
var harbors = await client.GetHarborsByStateAsync("pb");
foreach (var h in harbors)
    Console.WriteLine($"  {h.Id} - {h.Name}");

Console.WriteLine("\n=== Detalhes: pb01 ===");
var harbor = await client.GetHarborAsync("pb01");
Console.WriteLine($"  Nome: {harbor.Name} | Estado: {harbor.State} | TZ: {harbor.Timezone}");

Console.WriteLine("\n=== Tábua de Marés: pb01, janeiro, dias 1-3 ===");
var tides = await client.GetTideTableAsync("pb01", 1, [1, 2, 3]);
foreach (var t in tides)
{
    Console.WriteLine($"  Porto: {t.HarborName}");
    foreach (var month in t.Months)
        foreach (var day in month.Days)
            Console.WriteLine($"    Dia {day.Day} ({day.WeekdayName}): {day.Hours.Count} leituras");
}

Console.WriteLine("\n=== Porto mais próximo (João Pessoa/PB) ===");
var nearest = await client.GetNearestHarborByStateAsync("pb", -7.11509, -34.864);
Console.WriteLine($"  {nearest.Name} ({nearest.Id})");

Console.WriteLine("\n=== Tábua por Geolocalização ===");
var geoTides = await client.GetGeoTideTableAsync(-7.11509, -34.864, "pb", 1, [1]);
Console.WriteLine($"  Porto: {geoTides[0].HarborName}");

if (!string.IsNullOrWhiteSpace(options.ApiKey))
{
    Console.WriteLine("\n=== Uso da Cota ===");
    var usage = await client.GetUsageAsync();
    Console.WriteLine($"  Plano: {usage.Plan} | RPM: {usage.UsedRpm}/{usage.LimitRpm} | Mensal: {usage.UsedMonthly}/{usage.LimitMonthly}");
}
else
{
    Console.WriteLine("\n[Defina TABUAMARE_API_KEY para ver o uso da cota]");
}
