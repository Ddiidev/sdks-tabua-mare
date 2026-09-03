/**
 * Tábua de Marés SDK
 * SDK JavaScript para integração com a API Tábua de Marés (v2)
 */

const BASE_URL = 'https://tabuamare.api.br/api/v2';

/**
 * Cliente HTTP universal (funciona em Node.js e Browser)
 */
class HttpClient {
  async get(url, headers = {}) {
    // Detecta ambiente
    if (typeof window !== 'undefined' && window.fetch) {
      // Browser
      const response = await fetch(url, { headers });

      if (!response.ok) {
        throw new Error(`HTTP Error: ${response.status}`);
      }
      return response.json();
    } else if (typeof require !== 'undefined') {
      // Node.js
      const https = require('https');
      const urlParsed = new URL(url);

      const requestHeaders = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
        'Accept': 'application/json, text/plain, */*',
        'Accept-Language': 'pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7',
        'Connection': 'keep-alive',
        ...headers
      };

      const options = {
        hostname: urlParsed.hostname,
        port: urlParsed.port || 443,
        path: urlParsed.pathname + urlParsed.search,
        method: 'GET',
        headers: requestHeaders
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
 * Valida coordenadas geográficas
 */
function validateLatLng(lat, lng) {
  if (typeof lat !== 'number' || typeof lng !== 'number') {
    throw new Error('Latitude e longitude devem ser números');
  }

  if (isNaN(lat) || isNaN(lng) || !isFinite(lat) || !isFinite(lng)) {
    throw new Error('Latitude e longitude devem ser números válidos');
  }

  if (lat < -90 || lat > 90) {
    throw new Error('Latitude deve estar entre -90 e 90 graus');
  }

  if (lng < -180 || lng > 180) {
    throw new Error('Longitude deve estar entre -180 e 180 graus');
  }

  return `[${lat.toFixed(6)},${lng.toFixed(6)}]`;
}

/**
 * Converte array/string para o formato [dias] esperado pela API
 */
function formatDays(days) {
  if (Array.isArray(days)) {
    return `[${days.join(',').replace(/\s/g, '')}]`;
  } else if (typeof days === 'string' && !days.startsWith('[')) {
    return `[${days}]`;
  }
  return days;
}

/**
 * Cliente principal da API Tábua de Marés (v2)
 */
class TabuaMareClient {
  constructor(options = {}) {
    this.baseUrl = options.baseUrl || BASE_URL;
    this.apiKey = options.apiKey || '';
    this.http = new HttpClient();
  }

  /**
   * Monta os headers de autenticação quando há api_key configurada
   */
  authHeaders() {
    if (!this.apiKey) {
      return {};
    }
    return {
      'Authorization': `Bearer ${this.apiKey}`,
      'X-Api-Key': this.apiKey
    };
  }

  /**
   * Executa uma requisição GET autenticada (quando aplicável)
   */
  async request(path) {
    return this.http.get(`${this.baseUrl}${path}`, this.authHeaders());
  }

  /**
   * Lista todos os estados costeiros disponíveis
   * @returns {Promise<{data: string[], total: number}>}
   */
  async getStates() {
    return this.request('/states');
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
    return this.request(`/harbor_names/${state.toLowerCase()}`);
  }

  /**
   * Obtém informações detalhadas de um ou mais portos
   * @param {string|string[]} ids - ID ou IDs dos portos (ex: 'pb01' ou ['pb01','pe01'])
   * @returns {Promise<{data: Array, total: number}>}
   */
  async getHarbors(ids) {
    if (!ids || (Array.isArray(ids) && ids.length === 0)) {
      throw new Error('Informe pelo menos um ID de porto');
    }

    const idsStr = Array.isArray(ids) ? ids.join(',') : String(ids);
    return this.request(`/harbors/[${idsStr}]`);
  }

  /**
   * Obtém a tábua de maré para um porto específico
   * @param {string} harborId - ID do porto (ex: 'pb01')
   * @param {number} month - Mês (1-12)
   * @param {string|number[]} days - Dias no formato "[1,2,3]" ou array [1,2,3]
   * @returns {Promise<{data: Array, total: number}>}
   */
  async getTabuaMare(harborId, month, days) {
    if (!harborId || !month) {
      throw new Error('harborId e month são obrigatórios');
    }

    const encodedDays = encodeURIComponent(formatDays(days));
    return this.request(`/tabua-mare/${harborId}/${month}/${encodedDays}`);
  }

  /**
   * Obtém a tábua de maré para um período de dias
   * @param {string} harborId - ID do porto (ex: 'pb01')
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
   * @param {string} harborId - ID do porto (ex: 'pb01')
   * @param {number} month - Mês (1-12)
   * @returns {Promise<{data: Array, total: number}>}
   */
  async getTabuaMareMonth(harborId, month) {
    return this.getTabuaMare(harborId, month, '[1-31]');
  }

  /**
   * Obtém o porto mais próximo de uma coordenada dentro de um estado
   * @param {string} state - Sigla do estado (ex: 'pb', 'rj', 'sp')
   * @param {number} lat - Latitude (-90 a 90)
   * @param {number} lng - Longitude (-180 a 180)
   * @returns {Promise<{data: Array, total: number}>}
   */
  async getNearestHarborByState(state, lat, lng) {
    if (!state || typeof state !== 'string') {
      throw new Error('Estado inválido');
    }

    const latLng = validateLatLng(lat, lng);
    return this.request(`/nearested-harbor/${state.toLowerCase()}/${latLng}`);
  }

  /**
   * Obtém o porto mais próximo de uma coordenada (sem restringir por estado)
   * @param {number} lat - Latitude (-90 a 90)
   * @param {number} lng - Longitude (-180 a 180)
   * @returns {Promise<{data: Array, total: number}>}
   */
  async getNearestHarbor(lat, lng) {
    const latLng = validateLatLng(lat, lng);
    return this.request(`/nearest-harbor-independent-state/${latLng}`);
  }

  /**
   * Obtém a tábua de maré do porto mais próximo das coordenadas, dentro de um estado
   * @param {number} lat - Latitude (-90 a 90)
   * @param {number} lng - Longitude (-180 a 180)
   * @param {string} state - Sigla do estado (ex: 'pb', 'rj', 'sp')
   * @param {number} month - Mês (1-12)
   * @param {string|number[]} days - Dias no formato "[1,2,3]" ou array [1,2,3]
   * @returns {Promise<{data: Array, total: number}>}
   */
  async getGeoTabuaMare(lat, lng, state, month, days) {
    if (!state || typeof state !== 'string') {
      throw new Error('Estado inválido');
    }

    if (!month) {
      throw new Error('month é obrigatório');
    }

    const latLng = validateLatLng(lat, lng);
    const encodedDays = encodeURIComponent(formatDays(days));
    return this.request(`/geo-tabua-mare/${latLng}/${state.toLowerCase()}/${month}/${encodedDays}`);
  }

  /**
   * Consulta o uso da cota da api_key configurada.
   * Requer apiKey no construtor. Não consome a cota mensal.
   * @returns {Promise<{data: Array<{plan, limit_rpm, used_rpm, remaining_rpm, limit_monthly, used_monthly, remaining_monthly}>, total: number}>}
   */
  async getUsage() {
    if (!this.apiKey) {
      throw new Error('getUsage() requer apiKey: informe options.apiKey no construtor');
    }
    return this.request('/usage');
  }
}

// Exporta para diferentes ambientes
if (typeof module !== 'undefined' && module.exports) {
  module.exports = { TabuaMareClient };
}

if (typeof window !== 'undefined') {
  window.TabuaMareClient = TabuaMareClient;
}
