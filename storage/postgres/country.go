package postgres

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"ret/api/models"

	"github.com/google/uuid"
)

type CountryRepo struct {
	db *sql.DB
}

func NewCountryRepo(db *sql.DB) *CountryRepo {
	return &CountryRepo{
		db: db,
	}
}

func (p *CountryRepo) Create(req models.CreateCountry) (*models.Country, error) {
	var id string

	err := p.db.QueryRow(`INSERT INTO countries(guid, title, code, continent, updated_at) VALUES ($1, $2, $3, $4, now()) RETURNING guid`, uuid.New().String(), req.Title, req.Code, req.Continent).
		Scan(&id)
	if err != nil {
		return nil, err
	}

	return p.GetById(models.CountryPrimaryKey{Id: id})
}

func (c *CountryRepo) GetById(req models.CountryPrimaryKey) (*models.Country, error) {
	var (
		Guid      sql.NullString
		Title     sql.NullString
		Code      sql.NullString
		Continent sql.NullString
		CreatedAt sql.NullString
		UpdatedAt sql.NullString
	)

	err := c.db.QueryRow(`SELECT guid, title, code, continent, created_at, updated_at FROM countries WHERE guid = $1`, req.Id).
		Scan(
			&Guid,
			&Title,
			&Code,
			&Continent,
			&CreatedAt,
			&UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	return &models.Country{
		Guid:      Guid.String,
		Title:     Title.String,
		Code:      Code.String,
		Continent: Continent.String,
		CreatedAt: CreatedAt.String,
		UpdatedAt: UpdatedAt.String,
	}, nil
}

func (c *CountryRepo) GetList(req models.GetListCountryRequest) (*models.GetListCountryResponse, error) {
	var countries = models.GetListCountryResponse{}
	offset := req.Offset
	limit := req.Limit

	if offset < 0 {
		offset = 0
	}

	if limit <= 0 {
		limit = 10
	}

	rows, err := c.db.Query(`SELECT COUNT(*) OVER(), guid, title, code, continent, created_at, updated_at FROM countries`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			Guid      sql.NullString
			Title     sql.NullString
			Code      sql.NullString
			Continent sql.NullString
			CreatedAt sql.NullString
			UpdatedAt sql.NullString
		)

		err = rows.Scan(
			&countries.Count,
			&Guid,
			&Title,
			&Code,
			&Continent,
			&CreatedAt,
			&UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		countries.Countries = append(countries.Countries, models.Country{
			Guid:      Guid.String,
			Title:     Title.String,
			Code:      Code.String,
			Continent: Continent.String,
			CreatedAt: CreatedAt.String,
			UpdatedAt: UpdatedAt.String,
		})
	}

	return &countries, nil
}

func (c *CountryRepo) Update(req models.UpdateCountry) (*models.Country, error) {
	_, err := c.db.Exec(`UPDATE countries SET guid=$1, title=$2, code=$3, continent=$4, updated_at=now() WHERE guid = $5`, req.Guid, req.Title, req.Code, req.Continent, req.Guid)
	if err != nil {
		return nil, err
	}

	return c.GetById(models.CountryPrimaryKey{Id: req.Guid})
}

func (c *CountryRepo) Delete(req models.CountryPrimaryKey) error {
	_, err := c.db.Exec(`DELETE FROM countries WHERE guid = $1`, req.Id)
	if err != nil {
		return err
	}

	return nil
}



func (s *CountryRepo) ImportFromFileCountry(filePath string) error {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var countries []models.Country
	if err := json.Unmarshal(fileContent, &countries); err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, country := range countries {
		_, err := tx.Exec(
			`INSERT INTO countries (guid, title, code, continent) VALUES ($1, $2, $3, $4)`,
			country.Guid, country.Title, country.Code, country.Continent,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
