package postgres

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"ret/api/models"

	"github.com/google/uuid"
)

type AirportRepo struct {
	db *sql.DB
}

func NewAirportRepo(db *sql.DB) *AirportRepo {
	return &AirportRepo{
		db: db,
	}
}

func (p *AirportRepo) Create(req models.CreateAirport) (*models.Airport, error) {
	var id string

	err := p.db.QueryRow(`
		INSERT INTO buildings(
			guid,
			title,
			country_id,
			city_id,
			longitude,
			radius,
			image,
			address,
			timezone_id,
			country,
			city,
			search_text,
			code,
			product_count,
			gmt,
			updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,NOW()) RETURNING guid`,
		uuid.New().String(),
		req.Title,
		req.CountryId,
		req.CityId,
		req.Longitude,
		req.Radius,
		req.Image,
		req.Adress,
		req.TimezoneId,
		req.Country,
		req.City,
		req.SearchText,
		req.Code,
		req.ProductCount,
		req.Gmt,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	return p.GetById(models.AirportPrimaryKey{Id: id})
}

func (c *AirportRepo) GetById(req models.AirportPrimaryKey) (*models.Airport, error) {

	var (
		Id           sql.NullString
		Title        sql.NullString
		CountryId    sql.NullString
		CityId       sql.NullString
		Latitude     sql.NullFloat64
		Longitude    sql.NullFloat64
		Radius       sql.NullString
		Image        sql.NullString
		Adress       sql.NullString
		TimezoneId   sql.NullString
		Country      sql.NullString
		City         sql.NullString
		SearchText   sql.NullString
		Code         sql.NullString
		ProductCount sql.NullInt16
		Gmt          sql.NullString
		CreatedAt    sql.NullString
		UpdatedAt    sql.NullString
	)

	err := c.db.QueryRow(`
		SELECT
			guid,
			title,
			country_id,
			city_id,
			latitude,
			longitude,
			radius,
			image,
			address,
			timezone_id,
			country,
			city,
			search_text,
			code,
			product_count,
			gmt,
			created_at,
			updated_at
		FROM buildings
		WHERE guid = $1
	`, req.Id).Scan(
		&Id,
		&Title,
		&CountryId,
		&CityId,
		&Latitude,
		&Longitude,
		&Radius,
		&Image,
		&Adress,
		&TimezoneId,
		&Country,
		&City,
		&SearchText,
		&Code,
		&ProductCount,
		&Gmt,
		&CreatedAt,
		&UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Airport{
		Guid:           Id.String,
		Title:        Title.String,
		CountryId:    CountryId.String,
		CityId:       CityId.String,
		Latitude:     Latitude.Float64,
		Longitude:    Longitude.Float64,
		Radius:       Radius.String,
		Image:        Image.String,
		Adress:       Adress.String,
		TimezoneId:   TimezoneId.String,
		Country:      Country.String,
		City:         City.String,
		SearchText:   SearchText.String,
		Code:         Code.String,
		ProductCount: int(ProductCount.Int16),
		Gmt:          Gmt.String,
		CreatedAt:    CreatedAt.String,
		UpdatedAt:    UpdatedAt.String,
	}, nil
}

func (c *AirportRepo) GetList(req models.GetListAirportRequest) (*models.GetListAirportResponse, error) {
	var airports = models.GetListAirportResponse{}
	offset := req.Offset
	limit := req.Limit

	if offset <= 0 {
		offset = 0
	}

	if limit <= 0 {
		limit = 10
	}

	rows, err := c.db.Query(`
		SELECT
			COUNT(*) OVER(),
			guid,
			title,
			country_id,
			city_id,
			longitude,
			latitude, 
			radius,
			image,
			address,
			timezone_id,
			country,
			city,
			search_text,
			code,
			product_count,
			gmt,
			created_at,
			updated_at
		FROM buildings
		LIMIT $1 OFFSET $2  
	`, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			Id           sql.NullString
			Title        sql.NullString
			CountryId    sql.NullString
			CityId       sql.NullString
			Latitude     sql.NullFloat64
			Longitude    sql.NullFloat64
			Radius       sql.NullString
			Image        sql.NullString
			Adress       sql.NullString
			TimezoneId   sql.NullString
			Country      sql.NullString
			City         sql.NullString
			SearchText   sql.NullString
			Code         sql.NullString
			ProductCount sql.NullInt16
			Gmt          sql.NullString
			CreatedAt    sql.NullString
			UpdatedAt    sql.NullString
		)

		err = rows.Scan(
			&airports.Count,
			&Id,
			&Title,
			&CountryId,
			&CityId,
			&Longitude,
			&Latitude,
			&Radius,
			&Image,
			&Adress,
			&TimezoneId,
			&Country,
			&City,
			&SearchText,
			&Code,
			&ProductCount,
			&Gmt,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		airports.Airports = append(airports.Airports, models.Airport{
			Guid:           Id.String,
			Title:        Title.String,
			CountryId:    CountryId.String,
			CityId:       CityId.String,
			Longitude:    Longitude.Float64,
			Latitude:     Latitude.Float64,
			Radius:       Radius.String,
			Image:        Image.String,
			Adress:       Adress.String,
			TimezoneId:   TimezoneId.String,
			Country:      Country.String,
			City:         City.String,
			SearchText:   SearchText.String,
			Code:         Code.String,
			ProductCount: int(ProductCount.Int16),
			Gmt:          Gmt.String,
			CreatedAt:    CreatedAt.String,
			UpdatedAt:    UpdatedAt.String,
		})
	}
	return &airports, nil
}

func (c *AirportRepo) Update(req models.UpdateAirport) (*models.Airport, error) {
	_, err := c.db.Exec(`
		UPDATE buildings
		SET
			guid=$1,
			title=$2,
			country_id=$3,
			city_id=$4,
			longitude=$5,
			radius=$6,
			image=$7,
			address=$8,
			timezone_id=$9,
			country=$10,
			city=$11,
			search_text=$12,
			code=$13,
			product_count=$14,
			gmt=$15,
			updated_at=NOW()
		WHERE guid = $16
	`, req.Id, req.Title, req.CountryId, req.CityId, req.Longitude, req.Radius, req.Image, req.Adress, req.TimezoneId, req.Country, req.City, req.SearchText, req.Code, req.ProductCount, req.Gmt)

	if err != nil {
		return nil, err
	}

	return c.GetById(models.AirportPrimaryKey{Id: req.Id})
}

func (c *AirportRepo) Delete(req models.AirportPrimaryKey) error {
	_, err := c.db.Exec(`DELETE FROM buildings WHERE guid = $1`, req.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *AirportRepo) ImportFromFileAirport(filePath string) error {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var airports []models.Airport
	if err := json.Unmarshal(fileContent, &airports); err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, airport := range airports {
		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM countries WHERE guid = $1", airport.CountryId).Scan(&count)
		if err != nil {
			return err
		}

		_, err = tx.Exec(`
			INSERT INTO buildings (
				guid, title, country_id, city_id, latitude, longitude, radius, image, address, timezone_id, country, city, search_text, code, product_count, gmt
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
			airport.Guid, airport.Title, airport.CountryId, airport.CityId, airport.Latitude, airport.Longitude, airport.Radius, airport.Image, airport.Adress, airport.TimezoneId, airport.Country, airport.City, airport.SearchText, airport.Code, airport.ProductCount, airport.Gmt)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
