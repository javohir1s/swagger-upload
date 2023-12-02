package postgres

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"ret/api/models"
	"ret/pkg/helpers"

	"github.com/google/uuid"
)

type CityRepo struct {
	db *sql.DB
}

func NewCityRepo(db *sql.DB) *CityRepo {
	return &CityRepo{
		db: db,
	}
}

func (p *CityRepo) Create(req models.CreateCity) (*models.City, error) {
	id := uuid.New().String()
	query := `
		INSERT INTO cities(
			"guid",
			"title",
			"country_id",
			"city_code",
			"latitude",
			"longitude",
			"offset",
			"timezone_id",
			"country_name",
			"updated_at"
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW()) RETURNING guid`

	var createdID string
	err := p.db.QueryRow(query,
		id,
		req.Title,
		helpers.NewNullString(req.CountryId),
		req.CityCode,
		req.Latitude,
		req.Longitude,
		req.Offset,
		helpers.NewNullString(req.TimezoneId),
		req.CountryName,
	).Scan(&createdID)

	if err != nil {
		return nil, err
	}

	return p.GetById(models.CityPrimaryKey{Id: createdID})
}

func (c *CityRepo) GetById(req models.CityPrimaryKey) (*models.City, error) {
	query := `
		SELECT
			"guid",
			"title",
			"country_id",
			"city_code",
			"latitude",
			"longitude",
			"offset",
			"timezone_id",
			"country_name",
			"created_at",
			"updated_at"
		FROM cities
		WHERE guid = $1
	`

	var (
		Guid        sql.NullString
		Title       sql.NullString
		CountryId   sql.NullString
		CityCode    sql.NullString
		Latitude    sql.NullString
		Longitude   sql.NullString
		Offset      sql.NullString
		TimezoneId  sql.NullString
		CountryName sql.NullString
		CreatedAt   sql.NullString
		UpdatedAt   sql.NullString
	)

	err := c.db.QueryRow(query, req.Id).Scan(
		&Guid,
		&Title,
		&CountryId,
		&CityCode,
		&Latitude,
		&Longitude,
		&Offset,
		&TimezoneId,
		&CountryName,
		&CreatedAt,
		&UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &models.City{
		Guid:        Guid.String,
		Title:       Title.String,
		CountryId:   CountryId.String,
		CityCode:    CityCode.String,
		Latitude:    Latitude.String,
		Longitude:   Longitude.String,
		Offset:      Offset.String,
		TimezoneId:  TimezoneId.String,
		CountryName: CountryName.String,
		CreatedAt:   CreatedAt.String,
		UpdatedAt:   UpdatedAt.String,
	}, nil
}

func (c *CityRepo) GetList(req models.GetListCityRequest) (*models.GetListCityResponse, error) {
	var (
		resp = models.GetListCityResponse{}
	)
	offset := req.Offset
	limit := req.Limit

	if offset < 0 {
		offset = 0
	}

	if limit <= 0 {
		limit = 10
	}

	rows, err := c.db.Query(`
		SELECT
			COUNT(*) OVER(),
			"guid",
			"title",
			"country_id",
			"city_code",
			"latitude",
			"longitude",
			"offset",
			"timezone_id",
			"country_name",
			"created_at",
			"updated_at"
		FROM cities
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			Guid        sql.NullString
			Title       sql.NullString
			CountryId   sql.NullString
			CityCode    sql.NullString
			Latitude    sql.NullString
			Longitude   sql.NullString
			Offset      sql.NullString
			TimezoneId  sql.NullString
			CountryName sql.NullString
			CreatedAt   sql.NullString
			UpdatedAt   sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&Guid,
			&Title,
			&CountryId,
			&CityCode,
			&Latitude,
			&Longitude,
			&Offset,
			&TimezoneId,
			&CountryName,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Cities = append(resp.Cities, models.City{
			Guid:        Guid.String,
			Title:       Title.String,
			CountryId:   CountryId.String,
			CityCode:    CityCode.String,
			Latitude:    Latitude.String,
			Longitude:   Longitude.String,
			Offset:      Offset.String,
			TimezoneId:  TimezoneId.String,
			CountryName: CountryName.String,
			CreatedAt:   CreatedAt.String,
			UpdatedAt:   UpdatedAt.String,
		})
	}

	return &resp, nil
}

func (c *CityRepo) Update(req models.UpdateCity) (*models.City, error) {
	_, err := c.db.Exec(`UPDATE cities SET guid=$1, title=$2, country_id=$3, city_code=$4, latitude=$5, longitude=$6, offset=$7, timezone_id=$8, country_name=$9, updated_at=NOW() WHERE guid = $10`, req.Guid, req.Title, req.CountryId, req.CityCode, req.Latitude, req.Longitude, req.Offset, req.TimezoneId, req.CountryName, req.Id)
	if err != nil {
		return nil, err
	}

	return c.GetById(models.CityPrimaryKey{Id: req.Id})
}

func (c *CityRepo) Delete(req models.CityPrimaryKey) error {

	_, err := c.db.Exec(`DELETE FROM cities WHERE guid = $1`, req.Id)

	if err != nil {
		return err
	}

	return nil

}

func (s *CityRepo) ImportFromFile(filePath string) error {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var cities []models.City
	if err := json.Unmarshal(fileContent, &cities); err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, city := range cities {
		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM countries WHERE guid = $1", city.CountryId).Scan(&count)
		if err != nil {
			return err
		}

		var countryID interface{}
		if count == 0 {
			countryID = nil
		} else {
			countryID = city.CountryId
		}

		_, err = tx.Exec(`
			INSERT INTO cities (
				guid, title, country_id, city_code, latitude, longitude, "offset", timezone_id, country_name
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			city.Guid, city.Title, countryID, city.CityCode, city.Latitude, city.Longitude, city.Offset, city.TimezoneId, city.CountryName)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
