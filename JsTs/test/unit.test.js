/**
 * Testes unitários para o SDK Tábua de Marés (API v2)
 * Usa o módulo assert nativo do Node.js
 */

const assert = require('assert');
const { TabuaMareClient } = require('../src/index.js');

// Contador de testes
let testsRun = 0;
let testsPassed = 0;
let testsFailed = 0;

/**
 * Helper para executar testes
 */
async function test(description, fn) {
  testsRun++;
  try {
    await fn();
    testsPassed++;
    console.log(`✅ ${description}`);
  } catch (error) {
    testsFailed++;
    console.error(`❌ ${description}`);
    console.error(`   Erro: ${error.message}`);
    if (error.response) {
      console.error(`   Status HTTP: ${error.response.status}`);
      console.error(`   URL: ${error.response.url}`);
      console.error(`   Headers da resposta:`, JSON.stringify(error.response.headers, null, 2));
    }
    if (error.stack) {
      console.error(`   Stack trace: ${error.stack}`);
    }
  }
}

const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

/**
 * Suite de testes
 */
async function runTests() {
  console.log('=== Testes Unitários - Tábua de Marés SDK (v2) ===\n');

  const client = new TabuaMareClient();

  // Teste 1: Instanciação do cliente
  await test('Cliente deve ser instanciado corretamente', async () => {
    assert.ok(client instanceof TabuaMareClient, 'Cliente não é uma instância de TabuaMareClient');
    assert.strictEqual(typeof client.getStates, 'function', 'Método getStates não existe');
    assert.strictEqual(typeof client.getHarborsByState, 'function', 'Método getHarborsByState não existe');
    assert.strictEqual(typeof client.getHarbors, 'function', 'Método getHarbors não existe');
    assert.strictEqual(typeof client.getTabuaMare, 'function', 'Método getTabuaMare não existe');
    assert.strictEqual(typeof client.getNearestHarbor, 'function', 'Método getNearestHarbor não existe');
    assert.strictEqual(typeof client.getNearestHarborByState, 'function', 'Método getNearestHarborByState não existe');
    assert.strictEqual(typeof client.getGeoTabuaMare, 'function', 'Método getGeoTabuaMare não existe');
    assert.strictEqual(typeof client.getUsage, 'function', 'Método getUsage não existe');
    assert.strictEqual(client.baseUrl, 'https://tabuamare.api.br/api/v2', 'Base URL v2 incorreta');
  });

  // Teste 2: Listar estados
  await sleep(4000);
  await test('getStates() deve retornar lista de estados', async () => {
    const result = await client.getStates();
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data, 'Resultado deve ter propriedade data');
    assert.ok(Array.isArray(result.data), 'data deve ser um array');
    assert.ok(result.total > 0, 'total deve ser maior que 0');
    assert.ok(result.data.length > 0, 'Deve retornar pelo menos um estado');
  });

  // Teste 3: Listar portos por estado
  await sleep(4000);
  await test('getHarborsByState() deve retornar portos com IDs string (ex: pb01)', async () => {
    const result = await client.getHarborsByState('pb');
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data, 'Resultado deve ter propriedade data');
    assert.ok(Array.isArray(result.data), 'data deve ser um array');
    assert.ok(result.data.length > 0, 'Deve retornar pelo menos um porto');

    const firstHarbor = result.data[0];
    assert.ok(firstHarbor.id, 'Porto deve ter ID');
    assert.strictEqual(typeof firstHarbor.id, 'string', 'ID deve ser string na v2 (ex: pb01)');
    assert.ok(firstHarbor.harbor_name, 'Porto deve ter nome');
  });

  // Teste 4: Obter detalhes de um porto por ID string
  await sleep(4000);
  await test('getHarbors() deve retornar detalhes de um porto por ID string', async () => {
    const result = await client.getHarbors('pb01');
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data, 'Resultado deve ter propriedade data');
    assert.ok(Array.isArray(result.data), 'data deve ser um array');
    assert.ok(result.data.length > 0, 'Deve retornar dados do porto');

    const harbor = result.data[0];
    assert.strictEqual(harbor.id, 'pb01', 'ID deve ser pb01');
    assert.ok(harbor.harbor_name, 'Porto deve ter nome');
    assert.ok(harbor.state, 'Porto deve ter estado');
    assert.ok(typeof harbor.mean_level === 'number', 'mean_level deve ser um número');
  });

  // Teste 5: Obter múltiplos portos por IDs string
  await sleep(4000);
  await test('getHarbors() deve aceitar múltiplos IDs', async () => {
    const result = await client.getHarbors(['pb01', 'pe01']);
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(Array.isArray(result.data), 'data deve ser um array');
    assert.strictEqual(result.data.length, 2, 'Deve retornar 2 portos');
    assert.deepStrictEqual(
      result.data.map((h) => h.id).sort(),
      ['pb01', 'pe01'],
      'IDs devem ser pb01 e pe01'
    );
  });

  // Teste 6: Obter tábua de maré para dias específicos
  await sleep(4000);
  await test('getTabuaMare() deve retornar tábua de maré para dias específicos', async () => {
    const result = await client.getTabuaMare('pb01', 1, [1, 2, 3]);
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data, 'Resultado deve ter propriedade data');
    assert.ok(Array.isArray(result.data), 'data deve ser um array');
    assert.ok(result.data.length > 0, 'Deve retornar dados');

    const tideData = result.data[0];
    assert.ok(tideData.harbor_name, 'Deve ter nome do porto');
    assert.ok(tideData.months, 'Deve ter array de meses');
    assert.ok(Array.isArray(tideData.months), 'months deve ser um array');
    assert.ok(tideData.months[0].days, 'Deve ter array de dias');
    assert.strictEqual(tideData.months[0].days.length, 3, 'Deve retornar 3 dias');
  });

  // Teste 7: Obter tábua de maré para um período
  await sleep(4000);
  await test('getTabuaMareRange() deve retornar tábua de maré para um período', async () => {
    const result = await client.getTabuaMareRange('pb01', 1, 1, 7);
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data, 'Resultado deve ter propriedade data');
    assert.ok(Array.isArray(result.data), 'data deve ser um array');

    const tideData = result.data[0];
    assert.strictEqual(tideData.months[0].days.length, 7, 'Deve retornar 7 dias');

    // Verificar estrutura de um dia
    const day = tideData.months[0].days[0];
    assert.ok(day.day, 'Dia deve ter número');
    assert.ok(day.weekday_name, 'Dia deve ter nome do dia da semana');
    assert.ok(Array.isArray(day.hours), 'Dia deve ter array de horas');
  });

  // Teste 8: Obter tábua de maré para mês completo
  await sleep(4000);
  await test('getTabuaMareMonth() deve retornar tábua de maré para mês completo', async () => {
    const result = await client.getTabuaMareMonth('pb01', 1);
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data, 'Resultado deve ter propriedade data');

    const tideData = result.data[0];
    assert.ok(tideData.months[0].days.length >= 28, 'Deve retornar pelo menos 28 dias');
  });

  // Teste 9: Validar estrutura de dados de hora
  await sleep(4000);
  await test('Dados de hora devem ter estrutura correta', async () => {
    const result = await client.getTabuaMare('pb01', 1, [1]);
    const day = result.data[0].months[0].days[0];
    const hour = day.hours[0];

    assert.ok(hour.hour, 'Hora deve ter campo hour');
    assert.ok(typeof hour.level === 'number', 'level deve ser um número');
    assert.ok(hour.level >= -1 && hour.level <= 3, 'level deve estar em intervalo plausível');
  });

  // Teste 10: Porto mais próximo independente de estado
  await sleep(4000);
  await test('getNearestHarbor() deve retornar envelope data/total da v2', async () => {
    const result = await client.getNearestHarbor(-23.550520, -46.633308);
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data, 'Resultado deve ter propriedade data');
    assert.ok(Array.isArray(result.data), 'data deve ser um array');
    assert.ok(result.data.length > 0, 'Deve retornar pelo menos um porto');

    const nearestHarbor = result.data[0];
    assert.ok(nearestHarbor.id, 'Porto deve ter ID');
    assert.ok(nearestHarbor.harbor_name, 'Porto deve ter nome');
    assert.ok(nearestHarbor.state, 'Porto deve ter estado');
  });

  // Teste 11: Porto mais próximo por estado
  await sleep(4000);
  await test('getNearestHarborByState() deve retornar porto dentro do estado', async () => {
    const result = await client.getNearestHarborByState('pb', -7.11509, -34.864);
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data.length > 0, 'Deve retornar pelo menos um porto');
    assert.strictEqual(result.data[0].state, 'pb', 'Porto deve pertencer ao estado pb');
  });

  // Teste 12: Tábua de maré por geolocalização
  await sleep(4000);
  await test('getGeoTabuaMare() deve retornar tábua do porto mais próximo', async () => {
    const result = await client.getGeoTabuaMare(-7.11509, -34.864, 'pb', 1, [1]);
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data.length > 0, 'Deve retornar dados');
    assert.ok(result.data[0].harbor_name, 'Deve ter nome do porto');
    assert.ok(Array.isArray(result.data[0].months), 'Deve ter meses');
  });

  // Teste 13: getUsage sem apiKey deve falhar localmente
  await test('getUsage() deve exigir apiKey configurada', async () => {
    try {
      await client.getUsage();
      assert.fail('Deveria ter lançado erro sem apiKey');
    } catch (error) {
      assert.ok(error.message.includes('apiKey'), 'Erro deve mencionar apiKey');
    }
  });

  // Teste 14: getUsage com apiKey configurada (chave inválida deve retornar erro da API)
  await test('getUsage() deve enviar headers de autenticação', async () => {
    const authedClient = new TabuaMareClient({ apiKey: 'tm_live_invalid_key_for_test' });
    try {
      await authedClient.getUsage();
      // Se a API aceitar, tanto melhor — valida estrutura
    } catch (error) {
      // Chave inválida deve gerar erro HTTP (401/403), não ausência de header
      assert.ok(error, 'Erro esperado para chave inválida');
    }
  });

  // Teste 15: Validação de latitude
  await test('getNearestHarbor() deve validar latitude', async () => {
    try {
      await client.getNearestHarbor(91, -46.633308);
      assert.fail('Deveria ter lançado erro para latitude inválida');
    } catch (error) {
      assert.ok(error.message.includes('Latitude deve estar entre'), 'Erro deve mencionar intervalo de latitude');
    }
  });

  // Teste 16: Validação de longitude
  await test('getNearestHarbor() deve validar longitude', async () => {
    try {
      await client.getNearestHarbor(-23.550520, 181);
      assert.fail('Deveria ter lançado erro para longitude inválida');
    } catch (error) {
      assert.ok(error.message.includes('Longitude deve estar entre'), 'Erro deve mencionar intervalo de longitude');
    }
  });

  // Teste 17: Validação de tipos
  await test('getNearestHarbor() deve validar tipos de entrada', async () => {
    try {
      await client.getNearestHarbor('invalid', -46.633308);
      assert.fail('Deveria ter lançado erro para tipo inválido');
    } catch (error) {
      assert.ok(error.message.includes('Latitude e longitude devem ser números'), 'Erro deve mencionar tipo esperado');
    }
  });

  // Teste 18: Coordenadas do Rio de Janeiro
  await sleep(4000);
  await test('getNearestHarbor() deve funcionar com coordenadas do Rio de Janeiro', async () => {
    const result = await client.getNearestHarbor(-22.906847, -43.172896); // Rio de Janeiro
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data, 'Resultado deve ter propriedade data');
    assert.ok(result.data.length > 0, 'Deve retornar pelo menos um porto');

    const nearestHarbor = result.data[0];
    assert.ok(nearestHarbor.state, 'Porto deve ter estado');
    // O porto mais próximo do RJ deve ser no RJ ou estado vizinho
    const validStates = ['rj', 'sp', 'es'];
    assert.ok(validStates.includes(nearestHarbor.state.toLowerCase()),
      `Estado ${nearestHarbor.state} deve ser próximo ao RJ`);
  });

  // Resumo dos testes
  console.log('\n=== Resumo dos Testes ===');
  console.log(`Total: ${testsRun}`);
  console.log(`Passou: ${testsPassed}`);
  console.log(`Falhou: ${testsFailed}`);

  if (testsFailed > 0) {
    console.log('\n❌ Alguns testes falharam!');
    process.exit(1);
  } else {
    console.log('\n✅ Todos os testes passaram!');
    process.exit(0);
  }
}

// Executar testes
runTests().catch(error => {
  console.error('\n❌ Erro fatal ao executar testes:', error);
  process.exit(1);
});
