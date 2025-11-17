package tabuamare

// StatesResponse representa a resposta da listagem de estados
type StatesResponse struct {
	Data  []string  `json:"data"`
	Total int       `json:"total"`
	Error *APIError `json:"error,omitempty"`
}

// HarborName representa informações básicas de um porto
type HarborName struct {
	ID                        int    `json:"id"`
	Year                      int    `json:"year"`
	HarborName                string `json:"harbor_name"`
	DataCollectionInstitution string `json:"data_collection_institution"`
}

// HarborNamesResponse representa a resposta da listagem de portos
type HarborNamesResponse struct {
	Data  []HarborName `json:"data"`
	Total int          `json:"total"`
	Error *APIError    `json:"error,omitempty"`
}

// GeoLocation representa as coordenadas geográficas de um porto
type GeoLocation struct {
	Lat          string `json:"lat"`
	Lng          string `json:"lng"`
	DecimalLat   string `json:"decimal_lat"`
	DecimalLng   string `json:"decimal_lng"`
	LatDirection string `json:"lat_direction"`
	LngDirection string `json:"lng_direction"`
}

// Harbor representa informações detalhadas de um porto
type Harbor struct {
	ID          int           `json:"id"`
	HarborName  string        `json:"harbor_name"`
	State       string        `json:"state"`
	Timezone    string        `json:"timezone"`
	Card        string        `json:"card"`
	GeoLocation []GeoLocation `json:"geo_location"`
	MeanLevel   float64       `json:"mean_level"`
}

// HarborsResponse representa a resposta da consulta de portos
type HarborsResponse struct {
	Data  []Harbor  `json:"data"`
	Total int       `json:"total"`
	Error *APIError `json:"error,omitempty"`
}

// TideHour representa um horário e nível de maré
type TideHour struct {
	Hour  string  `json:"hour"`
	Level float64 `json:"level"`
}

// TideDay representa os dados de maré de um dia
type TideDay struct {
	WeekdayName string     `json:"weekday_name"`
	Day         int        `json:"day"`
	Hours       []TideHour `json:"hours"`
}

// TideMonth representa os dados de maré de um mês
type TideMonth struct {
	MonthName string    `json:"month_name"`
	Month     int       `json:"month"`
	Days      []TideDay `json:"days"`
}

// TideTable representa a tábua de marés completa
type TideTable struct {
	Year                      int         `json:"year"`
	HarborName                string      `json:"harbor_name"`
	State                     string      `json:"state"`
	Timezone                  string      `json:"timezone"`
	Card                      string      `json:"card"`
	DataCollectionInstitution string      `json:"data_collection_institution"`
	MeanLevel                 float64     `json:"mean_level"`
	Months                    []TideMonth `json:"months"`
}

// TideTableResponse representa a resposta da consulta de tábua de marés
type TideTableResponse struct {
	Data  []TideTable `json:"data"`
	Total int         `json:"total"`
	Error *APIError   `json:"error,omitempty"`
}

// NearestHarbor representa o porto mais próximo de uma coordenada
type NearestHarbor struct {
	Harbor
	Distance float64 `json:"distance_km"`
}

// NearestHarborResponse representa a resposta da consulta de porto mais próximo
type NearestHarborResponse struct {
	Data  []NearestHarbor `json:"data"`
	Total int             `json:"total"`
	Error *APIError       `json:"error,omitempty"`
}
