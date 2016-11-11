package storage

import (
	"database/sql"
	"fmt"
	"bytes"
	_ "github.com/lib/pq"
)

/**
CREATE TABLE feature_properties
(
    id NOT NULL,
    name NOT NULL,
    property text NOT NULL,
    value text NOT NULL,
    created date,
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
	INSERT_FEATURE_PROPERTY_SQL = "INSERT INTO feature_properties(id, name, property, value, created) values($1,$2,$3,$4,$5) returning created;"
	DELETE_FEATURE_SQL = "DELETE FROM feature_properties WHERE name=$1"
	SEARCH_FEATURE_SQL = "SELECT FROM feature_properties WHERE property=$1"
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (*FeatureStoreImpl) CreateFeature(feature Feature) string {
	/*
	db, err := sql.Open("postgres", dbinfo)
	db.Begin()
	db.Prepare("update userinfo set username=$1 where uid=$2")

	feature.name
	feature.properties

	err = db.QueryRow(INSERT_FEATURE_PROPERTY_SQL, "astaxie", "研发部门", "2012-12-09").Scan(&lastInsertId)
	checkErr(err)
	*/
	return ""
}

func (*FeatureStoreImpl) ReadFeature(id string) Feature {
	return Feature{}
}

func (*FeatureStoreImpl) DeleteFeature(id string) bool {
	return false
}

func (fs *FeatureStoreImpl) SearchFeature(name string, filter Filter) []Feature {
	var buffer bytes.Buffer

	SEARCH_SELECT_PART := "SELECT DISTINCT feature_properties.* FROM feature_properties "
	SEARCH_WHERE_PART := "WHERE feature_properties.name = $1 "

	buffer.WriteString(SEARCH_SELECT_PART)
	for i := 0; i < len(filter); i++ {
		buffer.WriteString(getInnerJoinLine(i))
	}
	buffer.WriteString(SEARCH_WHERE_PART)
	i := 0
	for k, v := range filter {
		buffer.WriteString(getPropertyFilterLine(i, k, v))
		i++
	}

	searchQuery := buffer.String()
	//fmt.Println(searchQuery)

	stmt, err := fs.db.Prepare(searchQuery)
	checkErr(err)
	defer stmt.Close()

	params := getParams(name, filter)
	//params := make([]interface{}, len(filter)*2+1)
	rows, err := stmt.Query(params...)
	checkErr(err)
	defer rows.Close()

	return rowsToFeatures(rows)

	//INNER JOIN feature_properties AS p1 ON feature_properties.name = p1.name
	//INNER JOIN feature_properties AS p2 ON feature_properties.name = p2.name
	//and (p1.property='prop1' and p1.value = 'val1')
	//and (p2.property='prop2' and p2.value = 'val2')
	//;"
}
func getParams(name string, filters Filter) []interface{} {
	params := make([]interface{}, 1)
	params[0] = name
	for propertyName, propertyValue := range filters {
		params = append(params, propertyName, propertyValue)
	}
	return params
}

func rowsToFeatures(rows *sql.Rows) []Feature {
	featureMap := make(map[string]Feature)
	for rows.Next() {
		var id string
		var name string
		var property string
		var value string
		var created string
		var expires string
		err := rows.Scan(&id, &name, &property, &value, &created, &expires)
		checkErr(err)
		feature, ok := featureMap[name]
		if !ok {
			props := make(Properties)
			//props[property] = value
			feature = Feature{name:name, properties:props}
			featureMap[name] = feature
		}
		feature.properties[property] = value
	}

	result := []Feature{}
	for _, feature := range featureMap {
		result = append(result, feature)
	}
	return result
}

func getPropertyFilterLine(i int, property string, value string) string {
	//and (p1.property='prop1' and p1.value = 'val1')
	SEARCH_FILTER_PART := "AND (p%d.property=$%d AND p%d.value = $%d)"
	return fmt.Sprintf(SEARCH_FILTER_PART, i, 2 + i * 2, i, 3 + i * 2)
}

func getInnerJoinLine(i int) string {
	SEARCH_INNER_JOIN_SQL := "INNER JOIN feature_properties AS p%d ON feature_properties.name = p%d.name "
	return fmt.Sprintf(SEARCH_INNER_JOIN_SQL, i, i)
}


