package storage

import (
	"database/sql"
	"fmt"
	"time"
	"bytes"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

/**
CREATE TABLE feature_properties
(
    id text NOT NULL,
    name text NOT NULL,
    property text NOT NULL,
    value text NOT NULL,
    created datetime,
    expires date,
    CONSTRAINT feature_pkey PRIMARY KEY (id),
    CONSTRAINT feature_unique UNIQUE (name,property)
)
WITH (OIDS=FALSE);
 */

const (
	DB_USER = "featuretoggle"
	DB_PASSWORD = "ftftft"
	DB_NAME = "featuretoggle"
	DB_HOST = "europe"
)

const (
	INSERT_FEATURE_PROPERTY_SQL = "INSERT INTO feature_properties(name, property, value, created, expires, enabled) values ($1,$2,$3,$4,$5,$6)"
)

type FeatureStoreImpl struct {
	db *sql.DB
}

func NewFeatureStore() *FeatureStoreImpl {
	return new(FeatureStoreImpl)
}

func (fs *FeatureStoreImpl) Open() error {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err == nil {
		fs.db = db
	}
	return err
}

func (fs *FeatureStoreImpl) Close() {
	fs.db.Close()
}

func (fs *FeatureStoreImpl) CreateFeature(feature Feature) (*string, error) {
	created := time.Now()

	tx, err := fs.db.Begin()
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to create trasaction, %v", err))
	}
	defer tx.Rollback()

	stmt, err := fs.db.Prepare(INSERT_FEATURE_PROPERTY_SQL)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to create prepared statement, %v", err))
	}
	defer stmt.Close()

	for property, value := range feature.properties {
		_, err := stmt.Exec(feature.name, property, value, created, feature.expires, feature.enabled)
		if ( err != nil) {
			return nil, errors.New(fmt.Sprintf("Failed to insert row with property '%s', %v", property, err))
		}
		/*
		rowsAffected, err := res.RowsAffected()
		if ( err != nil) {
			return nil, errors.New(fmt.Sprintf("Failed to get rows affected for property '%s', %v", property, err))
		}
		fmt.Printf("%v\n", rowsAffected)
		*/
	}
	tx.Commit()

	return &(feature.name), nil
}

func (fs *FeatureStoreImpl) ReadFeature(name string) (*Feature, error) {
	var buffer bytes.Buffer

	SEARCH_SELECT_PART := "SELECT DISTINCT feature_properties.* FROM feature_properties "
	SEARCH_WHERE_PART := "WHERE feature_properties.name = $1"

	buffer.WriteString(SEARCH_SELECT_PART)
	buffer.WriteString(SEARCH_WHERE_PART)

	stmt, err := fs.db.Prepare(buffer.String())
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to create prepared statement, %v", err))
	}
	defer stmt.Close()

	rows, err := stmt.Query(name)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to run query, %v", err))
	}
	defer rows.Close()

	features, err := rowsToFeatures(rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get row data, %v", err))
	}
	if len(features) > 0 {
		return &features[0], nil
	}
	return nil, nil

}

func (fs *FeatureStoreImpl) DeleteFeature(name string) (*bool, error) {

	DELETE_SQL := "DELETE FROM feature_properties WHERE feature_properties.name = $1"

	stmt, err := fs.db.Prepare(DELETE_SQL)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to create delete prepared statement, %v", err))
	}
	defer stmt.Close()

	res, err := stmt.Exec(name)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to delete feature, %v", err))
	}
	rowCount, err := res.RowsAffected()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get rowsAffected, %v", err))
	}
	b := rowCount > 0
	return &b, nil
}

func (fs *FeatureStoreImpl) SearchFeature(name *string, filter Filter) ([]Feature, error) {
	var buffer bytes.Buffer

	SEARCH_SELECT_PART := "SELECT DISTINCT feature_properties.* FROM feature_properties "
	SEARCH_NAME_PART := "feature_properties.name = $1 "

	buffer.WriteString(SEARCH_SELECT_PART)
	for i := 0; i < len(filter); i++ {
		buffer.WriteString(getInnerJoinLine(i))
	}

	buffer.WriteString("WHERE ")
	if ( name != nil) {
		buffer.WriteString(SEARCH_NAME_PART)
	}

	var offset int
	if ( name == nil) {
		offset = 1
	} else {
		offset = 2
	}
	for i := 0; i < len(filter); i++ {
		skipStartingAnd := i == 0 && name == nil

		buffer.WriteString(getPropertyFilterLine(i, skipStartingAnd, offset))
	}

	searchQuery := buffer.String()
	fmt.Println(searchQuery)

	stmt, err := fs.db.Prepare(searchQuery)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to create prepared statement, %v", err))
	}
	defer stmt.Close()

	params := getParams(name, &filter)
	rows, err := stmt.Query(params...)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to run query, %v", err))
	}
	defer rows.Close()

	res, err := rowsToFeatures(rows)
	if( err != nil) {
		return nil, err
	}
	return res, nil
}

func getParams(name *string, filters *Filter) []interface{} {
	params := make([]interface{}, 0)
	if name != nil {
		params = append(params, name)
	}
	if ( filters != nil) {
		for propertyName, propertyValue := range *filters {
			params = append(params, propertyName, propertyValue)
		}
	}
	return params
}

func rowsToFeatures(rows *sql.Rows) ([]Feature, error) {
	featureMap := make(map[string]Feature)
	for rows.Next() {
		var name string
		var property string
		var value string
		var created string
		var expires string
		var enabled bool
		err := rows.Scan(&name, &property, &value, &created, &expires, &enabled)
		if ( err != nil) {
			return nil, errors.New(fmt.Sprintf("Failed to scan row, %v", err))
		}
		feature, ok := featureMap[name]
		if !ok {
			props := make(Properties)
			feature = Feature{name:name, properties:props}
			featureMap[name] = feature
		}
		feature.properties[property] = value
	}

	result := []Feature{}
	for _, feature := range featureMap {
		result = append(result, feature)
	}
	return result, nil
}

func getPropertyFilterLine(i int, skipStartingAnd bool, offset int) string {

	SEARCH_FILTER_PART := "%s(p%d.property=$%d AND p%d.value = $%d) "
	var prefix string
	if ( skipStartingAnd) {
		prefix = ""
	} else {
		prefix = "AND"
	}
	return fmt.Sprintf(SEARCH_FILTER_PART, prefix, i, offset + i * 2, i, offset + 1 + i * 2)
}

func getInnerJoinLine(i int) string {
	SEARCH_INNER_JOIN_SQL := "INNER JOIN feature_properties AS p%d ON feature_properties.name = p%d.name "
	return fmt.Sprintf(SEARCH_INNER_JOIN_SQL, i, i)
}


