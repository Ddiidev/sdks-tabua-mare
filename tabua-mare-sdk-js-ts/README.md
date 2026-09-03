# Tábua de Maré — cliente JavaScript/TypeScript

Cliente JavaScript/TypeScript para integração com a API de Tábua de Maré do Brasil.

## Sobre o Projeto

Este SDK faz parte do ecossistema Tábua da Maré:

- 🌊 **API Principal**: [tabua_mare_api](https://github.com/Ddiidev/tabua_mare_api) - API REST que serve os dados de tábua da maré
- 📄 **Processador de Dados**: [tabua_mare_convert_pdf2db](https://github.com/Ddiidev/tabua_mare_convert_pdf2db) - Converte PDFs de tábua da maré para o banco de dados
- 📦 **Clientes**: [sdks-tabua-mare](https://github.com/Ddiidev/sdks-tabua-mare) - Monorepo com clientes para diversas linguagens

Para mais informações, acesse: [tabuamare.api.br](https://tabuamare.api.br)

> **⚠️ v2:** a API migrou para `https://tabuamare.api.br/api/v2`. A v1 foi descontinuada (HTTP 410). Os IDs de porto agora são strings (ex: `pb01`, `pe02`). Opcionalmente configure uma `apiKey` para mais requisições por minuto.

## Características

- ✅ Funciona em Node.js e Browser
- ✅ Suporte completo a TypeScript
- ✅ Zero dependências
- ✅ API simples e intuitiva
- ✅ Totalmente tipado
- ✅ Autenticação opcional via api_key
- ✅ Porto mais próximo (geral e por estado) e tábua por geolocalização

## Autenticação

A chave é opcional para endpoints públicos. Quando configurada em `apiKey`, o cliente envia:

```http
Authorization: Bearer <api_key>
X-Api-Key: <api_key>
```

Informe somente o valor da chave em `apiKey`; não inclua o prefixo `Bearer`.

```javascript
const client = new TabuaMareClient({
  apiKey: process.env.TABUAMARE_API_KEY,
});
```

Obtenha ou gerencie a chave no [dashboard](https://tabuamare.api.br/dashboard). O método `getUsage()` exige uma chave configurada. Consulte a seção [Como enviar a API Key](https://tabuamare.api.br/docs#api-key-header) para limites e códigos de resposta.

## Instalação

### Via NPM/Yarn/PNPM/Bun

```bash
# NPM
npm install tabua-mare-sdk

# Yarn
yarn add tabua-mare-sdk

# PNPM
pnpm add tabua-mare-sdk

# Bun
bun add tabua-mare-sdk
```

### Via CDN

```html
<!-- unpkg -->
<script src="https://unpkg.com/tabua-mare-sdk/src/index.js"></script>

<!-- unpkg (minificado) -->
<script src="https://unpkg.com/tabua-mare-sdk/src/index.min.js"></script>

<!-- jsDelivr -->
<script src="https://cdn.jsdelivr.net/npm/tabua-mare-sdk/src/index.js"></script>

<!-- jsDelivr (minificado) -->
<script src="https://cdn.jsdelivr.net/npm/tabua-mare-sdk/src/index.min.js"></script>

<!-- Uso -->
<script>
  const client = new TabuaMareClient();
  client.getStates().then(data => console.log(data));
</script>
```

### Cópia Direta

Ou copie o arquivo `src/index.js` diretamente para seu projeto.

## Uso

### Node.js (ES Modules - Recomendado)

```javascript
import { TabuaMareClient } from 'tabua-mare-sdk';

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

    // 5. Obter porto mais próximo de coordenadas
    console.log('\n5️⃣  Buscando porto mais próximo de Florianópolis...');
    const nearestHarbor = await client.getNearestHarbor(-27.5954, -48.5480);
    console.log('   Porto mais próximo:', nearestHarbor);

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
```

### Node.js (CommonJS)

```javascript
const { TabuaMareClient } = require('tabua-mare-sdk');

const client = new TabuaMareClient();

// Usando async/await
async function exemplo() {
  const states = await client.getStates();
  console.log('Estados:', states);
}

exemplo();

// Ou usando .then()
client.getStates().then(data => {
  console.log('Estados:', data);
});
```

### CDN (Browser)

```html
<!DOCTYPE html>
<html>
<head>
  <title>Tábua de Marés</title>
</head>
<body>
  <script src="https://unpkg.com/tabua-mare-sdk/src/index.js"></script>
  <script>
    const client = new TabuaMareClient();
    
    client.getStates().then(data => {
      console.log('Estados:', data);
    }).catch(error => {
      console.error('Erro:', error);
    });
  </script>
</body>
</html>
```

### TypeScript

```typescript
import { TabuaMareClient } from 'tabua-mare-sdk';

(async () => {
  const client = new TabuaMareClient();

  try {
    console.log('1️⃣  getStates()');
    const states = await client.getStates();
    console.log(JSON.stringify(states, null, 2));
    console.log('\n───────────────────────────────────────\n');

    console.log('2️⃣  getHarborsByState("sc")');
    const harbors = await client.getHarborsByState('sc');
    console.log(JSON.stringify(harbors, null, 2));
    console.log('\n───────────────────────────────────────\n');

    console.log('3️⃣  getHarbors("pb01")');
    const harbor = await client.getHarbors('pb01');
    console.log(JSON.stringify(harbor, null, 2));
    console.log('\n───────────────────────────────────────\n');

    console.log('4️⃣  getTabuaMare("pb01", 1, [1, 2, 3])');
    const tabuaMare = await client.getTabuaMare('pb01', 1, [1, 2, 3]);
    console.log(JSON.stringify(tabuaMare, null, 2));
    console.log('\n───────────────────────────────────────\n');

    console.log('5️⃣  getTabuaMareMonth("pb01", 1)');
    const tabuaMareMonth = await client.getTabuaMareMonth('pb01', 1);
    console.log(JSON.stringify(tabuaMareMonth, null, 2));
    console.log('\n───────────────────────────────────────\n');

    console.log('6️⃣  getNearestHarbor(-27.5954, -48.5480)');
    const nearestHarbor = await client.getNearestHarbor(-27.5954, -48.5480);
    console.log(JSON.stringify(nearestHarbor, null, 2));
    console.log('\n───────────────────────────────────────\n');

    console.log('7️⃣  getGeoTabuaMare(-27.5954, -48.5480, "sc", 1, [1, 2])');
    const geoTide = await client.getGeoTabuaMare(-27.5954, -48.5480, 'sc', 1, [1, 2]);
    console.log(JSON.stringify(geoTide, null, 2));

  } catch (error) {
    console.error('❌ Erro:', error);
  }
})()
```
