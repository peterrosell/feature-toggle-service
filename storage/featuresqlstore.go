package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	DB_USER = "featuretoggle"
	DB_PASSWORD = "ftftft"
	DB_NAME = "featuretoggle"
	DB_HOST = "europe"
)


type FeatureToggleStoreImpl struct {
	db *sql.DB
}

func NewFeatureToggleStoreImpl() *FeatureToggleStoreImpl {
	return new(FeatureToggleStoreImpl)
}

func (fs *FeatureToggleStoreImpl) Open() error {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err == nil {
		fs.db = db
	}
	return err
}

func (fs *FeatureToggleStoreImpl) Close() {
	fs.db.Close()
}



