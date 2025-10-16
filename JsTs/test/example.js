/**
 * Exemplo de uso do SDK Tábua de Marés
 */

const { TabuaMareClient } = require('../src/index.js');

async function runExamples() {
  const client = new TabuaMareClient();
  
  console.log('=== Tábua de Marés SDK - Exemplos ===\n');
  
  try {
    // 1. Listar estados
    console.log('1. Listando estados disponíveis...');
    const states = await client.getStates();
    console.log(`   Total de estados: ${states.total}`);
    console.log(`   Estados: ${states.data.join(', ')}\n`);
    
    // 2. Listar portos de um estado
    console.log('2. Listando portos da Paraíba (PB)...');
    const harbors = await client.getHarborsByState('pb');
    console.log(`   Total de portos: ${harbors.total}`);
    harbors.data.forEach(h => {
      console.log(`   - ${h.harbor_name} (ID: ${h.id})`);
    });
    console.log('');
    
    // 3. Obter detalhes de um porto
    console.log('3. Obtendo detalhes do porto ID 27...');
    const harborDetails = await client.getHarbors(27);
    const harbor = harborDetails.data[0];
    console.log(`   Nome: ${harbor.harbor_name}`);
    console.log(`   Estado: ${harbor.state.toUpperCase()}`);
    console.log(`   Nível médio: ${harbor.mean_level}m`);
    console.log(`   Timezone: ${harbor.timezone}\n`);
    
    // 4. Obter tábua de maré para dias específicos
    console.log('4. Obtendo tábua de maré para dias 1, 2 e 3 de janeiro...');
    const tides = await client.getTabuaMare(1, 1, [1, 2, 3]);
    const tideData = tides.data[0];
    console.log(`   Porto: ${tideData.harbor_name}`);
    console.log(`   Mês: ${tideData.months[0].month_name}`);
    console.log(`   Dias consultados: ${tideData.months[0].days.length}`);
    
    // Mostrar primeiro dia
    const firstDay = tideData.months[0].days[0];
    console.log(`\n   Exemplo - Dia ${firstDay.day} (${firstDay.weekday_name}):`);
    firstDay.hours.forEach(h => {
      console.log(`     ${h.hour} - Nível: ${h.level}m`);
    });
    console.log('');
    
    // 5. Obter tábua de maré para um período
    console.log('5. Obtendo tábua de maré para primeira semana de janeiro...');
    const weekTides = await client.getTabuaMareRange(1, 1, 1, 7);
    console.log(`   Dias retornados: ${weekTides.data[0].months[0].days.length}\n`);
    
    console.log('✅ Todos os exemplos executados com sucesso!');
    
  } catch (error) {
    console.error('❌ Erro:', error.message);
  }
}

// Executar exemplos
runExamples();
