/**
 * Tábua de Marés SDK
 * SDK JavaScript para integração com a API Tábua de Marés
 */

const BASE_URL = 'https://tabuamare.devtu.qzz.io/api/v1';

/**
 * Cliente HTTP universal (funciona em Node.js e Browser)
 */
class HttpClient {
  async get(url) {
    // Detecta ambiente
    if (typeof window !== 'undefined' && window.fetch) {
      // Browser
      const response = await fetch(url);

      if (!response.ok) {
        throw new Error(`HTTP Error: ${response.status}`);
      }
      return response.json();
    } else if (typeof require !== 'undefined') {
      // Node.js
      const https = require('https');
      const urlParsed = new URL(url);

      const headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
        'Accept': 'application/json, text/plain, */*',
        'Accept-Language': 'pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7',
        'Connection': 'keep-alive'
      };

      const options = {
        hostname: urlParsed.hostname,
        port: urlParsed.port || 443,
        path: urlParsed.pathname + urlParsed.search,
        method: 'GET',
        headers: headers
      };

      return new Promise((resolve, reject) => {
        const req = https.request(options, (res) => {
          let data = '';
          res.on('data', (chunk) => data += chunk);
          res.on('end', () => {
            if (res.statusCode >= 200 && res.statusCode < 300) {
              resolve(JSON.parse(data));
            } else {
              const error = new Error(`HTTP Error: ${res.statusCode}`);
              error.response = { status: res.statusCode, url: url };
              reject(error);
            }
          });
        });

        req.on('error', reject);
        req.end();
      });
    } else {
      throw new Error('Ambiente não suportado');
    }
  }
}

/**
 * Cliente principal da API Tábua de Marés
 */
class TabuaMareClient {
  constructor(options = {}) {
    this.baseUrl = options.baseUrl || BASE_URL;
    this.http = new HttpClient();
  }

  /**
   * Lista todos os estados costeiros disponíveis
   * @returns {Promise<{data: string[], total: number}>}
   */
  async getStates() {
    return this.http.get(`${this.baseUrl}/states`);
  }

  /**
   * Lista os portos de um estado específico
   * @param {string} state - Sigla do estado (ex: 'sp', 'rj', 'sc')
   * @returns {Promise<{data: Array, total: number}>}
   */
  async getHarborsByState(state) {
    if (!state || typeof state !== 'string') {
      throw new Error('Estado inválido');
    }
    return this.http.get(`${this.baseUrl}/harbor_names/${state.toLowerCase()}`);
  }

  /**
   * Obtém informações detalhadas de um ou mais portos
   * @param {number|number[]|string} ids - ID ou IDs dos portos
   * @returns {Promise<{data: Array, total: number}>}
   */
  async getHarbors(ids) {
    const idsStr = Array.isArray(ids) ? ids.join(',') : String(ids);
    return this.http.get(`${this.baseUrl}/harbors/${idsStr}`);
  }

  /**
   * Obtém a tábua de maré para um porto específico
   * @param {number} harborId - ID do porto
   * @param {number} month - Mês (1-12)
   * @param {string|number[]} days - Dias no formato "[1,2,3]" ou array [1,2,3]
   * @returns {Promise<{data: Array, total: number}>}
   */
  async getTabuaMare(harborId, month, days) {
    if (!harborId || !month) {
      throw new Error('harborId e month são obrigatórios');
    }

    // Converte array para string no formato esperado
    let daysStr = days;
    if (Array.isArray(days)) {
      daysStr = `[${days.join(',').replace(/\s/g, '')}]`;
    } else if (typeof days === 'string' && !days.startsWith('[')) {
      daysStr = `[${days}]`;
    }

    // Encode URL para lidar com caracteres especiais
    const encodedDays = encodeURIComponent(daysStr);
    return this.http.get(`${this.baseUrl}/tabua-mare/${harborId}/${month}/${encodedDays}`);
  }

  /**
   * Obtém a tábua de maré para um período de dias
   * @param {number} harborId - ID do porto
   * @param {number} month - Mês (1-12)
   * @param {number} startDay - Dia inicial
   * @param {number} endDay - Dia final
   * @returns {Promise<{data: Array, total: number}>}
   */
  async getTabuaMareRange(harborId, month, startDay, endDay) {
    return this.getTabuaMare(harborId, month, `[${startDay}-${endDay}]`);
  }

  /**
   * Obtém a tábua de maré para o mês completo
   * @param {number} harborId - ID do porto
   * @param {number} month - Mês (1-12)
   * @returns {Promise<{data: Array, total: number}>}
   */
  async getTabuaMareMonth(harborId, month) {
    return this.getTabuaMare(harborId, month, '[1-31]');
  }
}

// Exporta para diferentes ambientes
if (typeof module !== 'undefined' && module.exports) {
  module.exports = { TabuaMareClient };
}

if (typeof window !== 'undefined') {
  window.TabuaMareClient = TabuaMareClient;
}
