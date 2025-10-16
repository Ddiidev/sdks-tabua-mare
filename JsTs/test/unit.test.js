/**
 * Testes unitários para o SDK Tábua de Marés
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

/**
 * Suite de testes
 */
async function runTests() {
  console.log('=== Testes Unitários - Tábua de Marés SDK ===\n');

  const client = new TabuaMareClient();

  // Teste 1: Instanciação do cliente
  await test('Cliente deve ser instanciado corretamente', async () => {
    assert.ok(client instanceof TabuaMareClient, 'Cliente não é uma instância de TabuaMareClient');
    assert.strictEqual(typeof client.getStates, 'function', 'Método getStates não existe');
    assert.strictEqual(typeof client.getHarborsByState, 'function', 'Método getHarborsByState não existe');
    assert.strictEqual(typeof client.getHarbors, 'function', 'Método getHarbors não existe');
    assert.strictEqual(typeof client.getTabuaMare, 'function', 'Método getTabuaMare não existe');
  });

  // Teste 2: Listar estados
  await test('getStates() deve retornar lista de estados', async () => {
    const result = await client.getStates();
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data, 'Resultado deve ter propriedade data');
    assert.ok(Array.isArray(result.data), 'data deve ser um array');
    assert.ok(result.total > 0, 'total deve ser maior que 0');
    assert.ok(result.data.length > 0, 'Deve retornar pelo menos um estado');
  });

  // Teste 3: Listar portos por estado
  await test('getHarborsByState() deve retornar portos de um estado', async () => {
    const result = await client.getHarborsByState('pb');
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data, 'Resultado deve ter propriedade data');
    assert.ok(Array.isArray(result.data), 'data deve ser um array');
    assert.ok(result.data.length > 0, 'Deve retornar pelo menos um porto');

    const firstHarbor = result.data[0];
    assert.ok(firstHarbor.id, 'Porto deve ter ID');
    assert.ok(firstHarbor.harbor_name, 'Porto deve ter nome');
  });

  // Teste 4: Obter detalhes de um porto
  await test('getHarbors() deve retornar detalhes de um porto', async () => {
    const result = await client.getHarbors(27);
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data, 'Resultado deve ter propriedade data');
    assert.ok(Array.isArray(result.data), 'data deve ser um array');
    assert.ok(result.data.length > 0, 'Deve retornar dados do porto');

    const harbor = result.data[0];
    assert.ok(harbor.harbor_name, 'Porto deve ter nome');
    assert.ok(harbor.state, 'Porto deve ter estado');
    assert.ok(typeof harbor.mean_level === 'number', 'mean_level deve ser um número');
  });

  // Teste 5: Obter tábua de maré para dias específicos
  await test('getTabuaMare() deve retornar tábua de maré para dias específicos', async () => {
    const result = await client.getTabuaMare(1, 1, [1, 2, 3]);
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data, 'Resultado deve ter propriedade data');
    assert.ok(Array.isArray(result.data), 'data deve ser um array');
    assert.ok(result.data.length > 0, 'Deve retornar dados');

    const tideData = result.data[0];
    assert.ok(tideData.harbor_name, 'Deve ter nome do porto');
    assert.ok(tideData.months, 'Deve ter array de meses');
    assert.ok(Array.isArray(tideData.months), 'months deve ser um array');
    assert.ok(tideData.months[0].days, 'Deve ter array de dias');
    assert.ok(tideData.months[0].days.length === 3, 'Deve retornar 3 dias');
  });

  // Teste 6: Obter tábua de maré para um período
  await test('getTabuaMareRange() deve retornar tábua de maré para um período', async () => {
    const result = await client.getTabuaMareRange(1, 1, 1, 7);
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data, 'Resultado deve ter propriedade data');
    assert.ok(Array.isArray(result.data), 'data deve ser um array');

    const tideData = result.data[0];
    assert.ok(tideData.months[0].days.length === 7, 'Deve retornar 7 dias');

    // Verificar estrutura de um dia
    const day = tideData.months[0].days[0];
    assert.ok(day.day, 'Dia deve ter número');
    assert.ok(day.weekday_name, 'Dia deve ter nome do dia da semana');
    assert.ok(Array.isArray(day.hours), 'Dia deve ter array de horas');
  });

  // Teste 7: Obter tábua de maré para mês completo
  await test('getTabuaMareMonth() deve retornar tábua de maré para mês completo', async () => {
    const result = await client.getTabuaMareMonth(1, 1);
    assert.ok(result, 'Resultado não deve ser nulo');
    assert.ok(result.data, 'Resultado deve ter propriedade data');

    const tideData = result.data[0];
    assert.ok(tideData.months[0].days.length >= 28, 'Deve retornar pelo menos 28 dias');
  });

  // Teste 8: Validar estrutura de dados de hora
  await test('Dados de hora devem ter estrutura correta', async () => {
    const result = await client.getTabuaMare(1, 1, [1]);
    const day = result.data[0].months[0].days[0];
    const hour = day.hours[0];

    assert.ok(hour.hour, 'Hora deve ter campo hour');
    assert.ok(typeof hour.level === 'number', 'level deve ser um número');
    assert.ok(hour.level >= 0, 'level deve ser maior ou igual a 0');
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
