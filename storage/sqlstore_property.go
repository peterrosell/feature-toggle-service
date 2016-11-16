package storage

import (
	"fmt"
	"database/sql"
	"errors"
)

const (
	INSERT_PROPERTY_SQL = "INSERT INTO property(name, description) values ($1,$2)"
	READ_PROPERTY_SQL = "SELECT * FROM property WHERE name = $1"
	DELETE_PROPERTY_SQL = "DELETE FROM property WHERE name = $1"
	READ_ALL_PROPERTY_NAMES_SQL = "SELECT name FROM property"
)

func (fs *FeatureToggleStoreImpl) CreateProperty(property Property) (*string, error) {
	stmt, err := fs.db.Prepare(INSERT_PROPERTY_SQL)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("CreateProperty: Failed to create prepared statement, %v", err))
	}
	defer stmt.Close()

	_, err = stmt.Exec(property.Name, property.Description)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("CreateProperty: Failed to insert property '%s', %v", property.Name, err))
	}

	return &property.Name, nil
}

func (fs *FeatureToggleStoreImpl) ReadProperty(name string) (*Property, error) {
	stmt, err := fs.db.Prepare(READ_PROPERTY_SQL)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("ReadProperty: Failed to create prepared statement, %v", err))
	}
	defer stmt.Close()

	rows, err := stmt.Query(name)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("ReadProperty: Failed to select '%s', %v", name, err))
	}
	properties, err := rowsToProperty(rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("ReadProperty: Failed to get row data, %v", err))
	}
	if len(properties) > 0 {
		return &properties[0], nil
	}
	return nil, nil
}

func (fs *FeatureToggleStoreImpl) ReadAllPropertyNames() (*[]string, error) {
	rows, err := fs.db.Query(READ_ALL_PROPERTY_NAMES_SQL)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("ReadAllPropertyNames: Failed to create prepared statement, %v", err))
	}
	names, err := rowsToPropertyNames(rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("ReadAllPropertyNames: Failed to get row data, %v", err))
	}
	return &names, nil
}

func (fs *FeatureToggleStoreImpl) DeleteProperty(name string) (*bool, error) {
	stmt, err := fs.db.Prepare(DELETE_PROPERTY_SQL)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("DeleteProperty: Failed to create prepared statement, %v", err))
	}
	defer stmt.Close()

	res, err := stmt.Exec(name)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("DeleteProperty: Failed to delete '%s', %v", name, err))
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("DeleteProperty: Failed to get rowsAffected, %v", err))
	}
	b := rowCount > 0
	return &b, nil
}

func (fs *FeatureToggleStoreImpl) SearchProperty(name string) (*[]Property, error) {
	return nil, errors.New("Not implemented")
}

func rowsToProperty(rows *sql.Rows) ([]Property, error) {
	properties := []Property{}
	for rows.Next() {
		var name string
		var description string
		err := rows.Scan(&name, &description)
		if ( err != nil) {
			return nil, errors.New(fmt.Sprintf("Property: Failed to scan row, %v", err))
		}
		feature := Property{name, description}
		properties = append(properties, feature)
	}
	return properties, nil
}

func rowsToPropertyNames(rows *sql.Rows) ([]string, error) {
	names := []string{}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if ( err != nil) {
			return nil, errors.New(fmt.Sprintf("Property: Failed to scan row, %v", err))
		}
		names = append(names, name)
	}
	return names, nil
}

