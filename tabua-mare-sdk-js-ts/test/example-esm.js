import { TabuaMareClient } from '../src/index.js';

console.log('📘 Exemplo de uso do Tabua Mare SDK\n');

const client = new TabuaMareClient();

async function exemploCompleto() {
  try {
    // 1. Listar todos os estados disponíveis
    console.log('1️⃣  Buscando estados...');
    const states = await client.getStates();
    console.log('   Estados:', states);

    // 2. Listar portos de Santa Catarina
    console.log('\n2️⃣  Buscando portos de SC...');
    const harbors = await client.getHarborsByState('sc');
    console.log('   Portos:', harbors);

    // 3. Obter detalhes de um porto
    if (harbors && harbors.length > 0) {
      const harborId = harbors[0].id;
      console.log(`\n3️⃣  Buscando detalhes do porto ${harborId}...`);
      const harbor = await client.getHarbors(harborId);
      console.log('   Detalhes:', harbor);
    }

    // 4. Obter tábua de maré do mês atual
    const now = new Date();
    const month = now.getMonth() + 1;
    console.log(`\n4️⃣  Buscando tábua de maré do mês ${month}...`);
    const tabuaMare = await client.getTabuaMareMonth('pb01', month);
    console.log('   Tábua de Maré:', tabuaMare);

  } catch (error) {
    console.error('\n❌ Erro na execução:', error.message);
    console.error('   Detalhes:', error);
    
    if (error.message.includes('502')) {
      console.log('\n⚠️  A API está temporariamente indisponível (erro 502).');
      console.log('   Isso pode acontecer se o servidor estiver em manutenção.');
    }
  }
}

exemploCompleto();
