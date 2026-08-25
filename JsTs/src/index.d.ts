/**
 * Tábua de Marés SDK - TypeScript Definitions (API v2)
 */

export interface TabuaMareOptions {
  baseUrl?: string;
  /** API key opcional: aumenta o rate limit (header Authorization: Bearer + X-Api-Key) */
  apiKey?: string;
}

export interface ApiResponse<T> {
  data: T;
  total: number;
  error?: {
    message: string;
    code: number;
  };
}

export interface GeoLocation {
  lat: string;
  lng: string;
  decimal_lat: string;
  decimal_lng: string;
  lat_direction: string;
  lng_direction: string;
}

export interface Harbor {
  id: string;
  harbor_name: string;
  state: string;
  timezone: string;
  card: string;
  geo_location: GeoLocation[];
  mean_level: number;
}

export interface HarborName {
  id: string;
  year: number;
  harbor_name: string;
  data_collection_institution: string;
}

export interface TideHour {
  hour: string;
  level: number;
}

export interface TideDay {
  weekday_name: string;
  day: number;
  hours: TideHour[];
}

export interface TideMonth {
  month_name: string;
  month: number;
  days: TideDay[];
}

export interface TabuaMare {
  year: number;
  harbor_name: string;
  state: string;
  timezone: string;
  card: string;
  data_collection_institution: string;
  mean_level: number;
  months: TideMonth[];
}

export interface NearestHarbor extends Harbor {}

/** Campos numéricos chegam como string; "-1" indica limite ilimitado */
export interface UsageInfo {
  plan: string;
  limit_rpm: string;
  used_rpm: string;
  remaining_rpm: string;
  limit_monthly: string;
  used_monthly: string;
  remaining_monthly: string;
}

/**
 * Cliente principal da API Tábua de Marés (v2)
 */
export class TabuaMareClient {
  constructor(options?: TabuaMareOptions);

  /**
   * Lista todos os estados costeiros disponíveis
   */
  getStates(): Promise<ApiResponse<string[]>>;

  /**
   * Lista os portos de um estado específico
   * @param state - Sigla do estado (ex: 'sp', 'rj', 'sc')
   */
  getHarborsByState(state: string): Promise<ApiResponse<HarborName[]>>;

  /**
   * Obtém informações detalhadas de um ou mais portos
   * @param ids - ID ou IDs dos portos (ex: 'pb01' ou ['pb01','pe01'])
   */
  getHarbors(ids: string | string[]): Promise<ApiResponse<Harbor[]>>;

  /**
   * Obtém a tábua de maré para um porto específico
   * @param harborId - ID do porto (ex: 'pb01')
   * @param month - Mês (1-12)
   * @param days - Dias no formato "[1,2,3]" ou array [1,2,3]
   */
  getTabuaMare(
    harborId: string,
    month: number,
    days: string | number[]
  ): Promise<ApiResponse<TabuaMare[]>>;

  /**
   * Obtém a tábua de maré para um período de dias
   * @param harborId - ID do porto (ex: 'pb01')
   * @param month - Mês (1-12)
   * @param startDay - Dia inicial
   * @param endDay - Dia final
   */
  getTabuaMareRange(
    harborId: string,
    month: number,
    startDay: number,
    endDay: number
  ): Promise<ApiResponse<TabuaMare[]>>;

  /**
   * Obtém a tábua de maré para o mês completo
   * @param harborId - ID do porto (ex: 'pb01')
   * @param month - Mês (1-12)
   */
  getTabuaMareMonth(
    harborId: string,
    month: number
  ): Promise<ApiResponse<TabuaMare[]>>;

  /**
   * Obtém o porto mais próximo de uma coordenada dentro de um estado
   * @param state - Sigla do estado (ex: 'pb', 'rj', 'sp')
   * @param lat - Latitude (-90 a 90)
   * @param lng - Longitude (-180 a 180)
   */
  getNearestHarborByState(
    state: string,
    lat: number,
    lng: number
  ): Promise<ApiResponse<NearestHarbor[]>>;

  /**
   * Obtém o porto mais próximo de uma coordenada (sem restringir por estado)
   * @param lat - Latitude (-90 a 90)
   * @param lng - Longitude (-180 a 180)
   */
  getNearestHarbor(lat: number, lng: number): Promise<ApiResponse<NearestHarbor[]>>;

  /**
   * Obtém a tábua de maré do porto mais próximo das coordenadas, dentro de um estado
   * @param lat - Latitude (-90 a 90)
   * @param lng - Longitude (-180 a 180)
   * @param state - Sigla do estado (ex: 'pb', 'rj', 'sp')
   * @param month - Mês (1-12)
   * @param days - Dias no formato "[1,2,3]" ou array [1,2,3]
   */
  getGeoTabuaMare(
    lat: number,
    lng: number,
    state: string,
    month: number,
    days: string | number[]
  ): Promise<ApiResponse<TabuaMare[]>>;

  /**
   * Consulta o uso da cota da api_key configurada.
   * Requer apiKey no construtor. Não consome a cota mensal.
   */
  getUsage(): Promise<ApiResponse<UsageInfo[]>>;
}
