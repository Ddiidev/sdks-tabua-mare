/**
 * Tábua de Marés SDK - TypeScript Definitions
 */

export interface TabuaMareOptions {
  baseUrl?: string;
}

export interface ApiResponse<T> {
  data: T;
  total: number;
  error?: {
    msg: string;
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
  id: number;
  harbor_name: string;
  state: string;
  timezone: string;
  card: string;
  geo_location: GeoLocation[];
  mean_level: number;
}

export interface HarborName {
  id: number;
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

export interface NearestHarbor extends Harbor {
  // O endpoint nearest-harbor retorna um Harbor regular, sem informação de distância
}

/**
 * Cliente principal da API Tábua de Marés
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
   * @param ids - ID ou IDs dos portos
   */
  getHarbors(ids: number | number[] | string): Promise<ApiResponse<Harbor[]>>;

  /**
   * Obtém a tábua de maré para um porto específico
   * @param harborId - ID do porto
   * @param month - Mês (1-12)
   * @param days - Dias no formato "[1,2,3]" ou array [1,2,3]
   */
  getTabuaMare(
    harborId: number,
    month: number,
    days: string | number[]
  ): Promise<ApiResponse<TabuaMare[]>>;

  /**
   * Obtém a tábua de maré para um período de dias
   * @param harborId - ID do porto
   * @param month - Mês (1-12)
   * @param startDay - Dia inicial
   * @param endDay - Dia final
   */
  getTabuaMareRange(
    harborId: number,
    month: number,
    startDay: number,
    endDay: number
  ): Promise<ApiResponse<TabuaMare[]>>;

  /**
   * Obtém a tábua de maré para o mês completo
   * @param harborId - ID do porto
   * @param month - Mês (1-12)
   */
  getTabuaMareMonth(
    harborId: number,
    month: number
  ): Promise<ApiResponse<TabuaMare[]>>;

  /**
   * Obtém o porto mais próximo de uma coordenada geográfica
   * @param lat - Latitude (-90 a 90)
   * @param lng - Longitude (-180 a 180)
   */
  getNearestHarbor(lat: number, lng: number): Promise<ApiResponse<NearestHarbor[]>>;
}
