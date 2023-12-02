package helpers

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"regexp"
)

// IsValidPhone ...
func IsValidPhone(phone string) bool {
	r := regexp.MustCompile(`^\+998[0-9]{2}[0-9]{7}$`)
	return r.MatchString(phone)
}

// IsValidEmail ...
func IsValidEmail(email string) bool {
	r := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	return r.MatchString(email)
}

// IsValidLogin ...
func IsValidLogin(login string) bool {
	r := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{5,29}$`)
	return r.MatchString(login)
}

// IsValidUUID ...
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func NewNullString(s string) sql.NullString {

	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}



type City struct {
	Guid        string  `json:"guid"`
	Title       string  `json:"title"`
	CountryID   string  `json:"country_id"`
	CityCode    string  `json:"city_code"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Offset      string  `json:"offset"`
	TimezoneID  string  `json:"timezone_id"`
	CountryName string  `json:"country_name"`
}


func ReadJSONFile(filePath string) ([]City, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []City
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
