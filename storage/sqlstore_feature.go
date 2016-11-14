package storage

import (
	"fmt"
	"database/sql"
	"errors"
)

const (
	INSERT_FEATURE_SQL = "INSERT INTO feature(id, name, enabled, description) values ($1,$2,$3,$4)"
	READ_FEATURE_SQL = "SELECT * FROM feature WHERE id = $1"
	READ_FEATURE_BY_NAME_SQL = "SELECT * FROM feature WHERE name = $1"
	DELETE_FEATURE_SQL = "DELETE FROM feature WHERE id = $1"

)
func (fs *FeatureToggleStoreImpl) CreateFeature(feature Feature) (*string, error) {
	stmt, err := fs.db.Prepare(INSERT_FEATURE_SQL)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to create prepared statement, %v", err))
	}
	defer stmt.Close()

	_, err = stmt.Exec(feature.Id, feature.Name, feature.Enabled, feature.Description)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to insert property '%s', %v", feature.Name, err))
	}

	return &feature.Id, nil
}

func (fs *FeatureToggleStoreImpl) ReadFeature(id string) (*Feature, error) {
	stmt, err := fs.db.Prepare(READ_FEATURE_SQL)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("ReadFeature: Failed to create prepared statement, %v", err))
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("ReadFeature: Failed to select '%s', %v", id, err))
	}
	features, err := rowsToFeature(rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get row data, %v", err))
	}
	if len(features) > 0 {
		return &features[0], nil
	}
	return nil, nil
}

func (fs *FeatureToggleStoreImpl) ReadFeatureByName(name string) (*Feature, error) {
	stmt, err := fs.db.Prepare(READ_FEATURE_BY_NAME_SQL)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("ReadFeatureByName: Failed to create prepared statement, %v", err))
	}
	defer stmt.Close()

	rows, err := stmt.Query(name)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("ReadFeatureByName: Failed to select '%s', %v", name, err))
	}
	features, err := rowsToFeature(rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("ReadFeatureByName: Failed to get row data, %v", err))
	}
	if len(features) > 0 {
		return &features[0], nil
	}
	return nil, nil
}

func (fs *FeatureToggleStoreImpl) DeleteFeature(id string) (*bool, error) {
	stmt, err := fs.db.Prepare(DELETE_FEATURE_SQL)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("DeleteFeature: Failed to create prepared statement, %v", err))
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("DeleteFeature: Failed to delete '%s', %v", id, err))
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("DeleteFeature: Failed to get rowsAffected, %v", err))
	}
	b := rowCount > 0
	return &b, nil
}

func (fs *FeatureToggleStoreImpl) SearchFeature(name string) (*([]Feature), error) {
	return nil, errors.New("Not implemented")
}

func rowsToFeature(rows *sql.Rows) ([]Feature, error) {
	features := []Feature{}
	for rows.Next() {
		var id string
		var name string
		var description string
		var enabled bool
		err := rows.Scan(&id, &name, &enabled, &description)
		if ( err != nil) {
			return nil, errors.New(fmt.Sprintf("Feature: Failed to scan row, %v", err))
		}
		feature := Feature{id, name, enabled, description}
		features = append(features, feature)
	}
	return features, nil
}



