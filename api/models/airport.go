package models

type Airport struct {
	Guid           string  `json:"guid"`
	Title        string  `json:"title"`
	CountryId    string  `json:"country_id"`
	CityId       string  `json:"city_id"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Radius       string  `json:"radius"`
	Image        string  `json:"image"`
	Adress       string  `json:"adress"`
	TimezoneId   string  `json:"timezone_id"`
	Country      string  `json:"country"`
	City         string  `json:"city"`
	SearchText   string  `json:"search_text"`
	Code         string  `json:"code"`
	ProductCount int     `json:"product_count"`
	Gmt          string  `json:"gmt"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}


type CreateAirport struct {
	Title        string  `json:"title"`
	CountryId    string  `json:"country_id"`
	CityId       string  `json:"city_id"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Radius       string  `json:"radius"`
	Image        string  `json:"image"`
	Adress       string  `json:"adress"`
	TimezoneId   string  `json:"timezone_id"`
	Country      string  `json:"country"`
	City         string  `json:"city"`
	SearchText   string  `json:"search_text"`
	Code         string  `json:"code"`
	ProductCount int     `json:"product_count"`
	Gmt          string  `json:"gmt"`
}

type UpdateAirport struct {
	Id           string  `json:"id"`
	Title        string  `json:"title"`
	CountryId    string  `json:"country_id"`
	CityId       string  `json:"city_id"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Radius       string  `json:"radius"`
	Image        string  `json:"image"`
	Adress       string  `json:"adress"`
	TimezoneId   string  `json:"timezone_id"`
	Country      string  `json:"country"`
	City         string  `json:"city"`
	SearchText   string  `json:"search_text"`
	Code         string  `json:"code"`
	ProductCount int     `json:"product_count"`
	Gmt          string  `json:"gmt"`
}

type AirportPrimaryKey struct {
	Id string `json:"id"`
}

type GetListAirportRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type GetListAirportResponse struct {
	Count    int       `json:"count"`
	Airports []Airport `json:"airports"`
}


type Airports struct {
	ID           struct {
		Oid string `json:"$oid"`
	} `json:"_id"`
	Guid         string  `json:"guid"`
	Title        string  `json:"title"`
	CountryID    string  `json:"country_id"`
	CityID       string  `json:"city_id"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Radius       float64 `json:"radius"`
	Image        string  `json:"image"`
	Address      string  `json:"adress"`
	TimezoneID   string  `json:"timezone_id"`
	Country      string  `json:"country"`
	City         string  `json:"city"`
	SearchText   string  `json:"search_text"`
	Code         string  `json:"code"`
	ProductCount int     `json:"product_count"`
	GMT          string  `json:"gmt"`
	V            int     `json:"__v"`
	CreatedAt    struct {
		Date string `json:"$date"`
	} `json:"createdAt"`
	UpdatedAt struct {
		Date string `json:"$date"`
	} `json:"updatedAt"`
}
