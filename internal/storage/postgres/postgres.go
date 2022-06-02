package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mykhalskyio/image-api/internal/config"
	"github.com/mykhalskyio/image-api/internal/entity"
)

// postgres struct
type Postgres struct {
	db *sqlx.DB
}

// new postgres struct and connect
func NewPostgres(cfg *config.Config) (*Postgres, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=%s ",
		cfg.Postgres.Name,
		cfg.Postgres.User,
		cfg.Postgres.Pass,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Sslmode))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Postgres{db: db}, nil
}

// insert image to db
func (psql *Postgres) Insert(img *entity.Image) error {
	_, err := psql.db.Exec("INSERT INTO image (img_path_quality_original, img_path_quality_75, img_path_quality_50, img_path_quality_25) VALUES($1, $2, $3, $4);",
		img.ImagePathQualityOriginal,
		img.ImagePathQuality75,
		img.ImagePathQuality50,
		img.ImagePathQuality25)
	if err != nil {
		return err
	}
	return nil
}

// get image from db
func (psql *Postgres) Get(id int) (*entity.Image, error) {
	img := &entity.Image{}
	err := psql.db.QueryRowx("SELECT * FROM image WHERE id = $1", id).Scan(&img.Id,
		&img.ImagePathQualityOriginal,
		&img.ImagePathQuality75,
		&img.ImagePathQuality50,
		&img.ImagePathQuality25)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// delete from db
func (psql *Postgres) Delete(id int) error {
	_, err := psql.db.Exec("DELETE FROM image WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
