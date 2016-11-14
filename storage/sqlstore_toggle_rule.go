package storage

import (
	"database/sql"
	"github.com/peterrosell/feature-toggle-service/featuretree"
	"fmt"
	"time"
	"strings"
	"github.com/satori/go.uuid"
	"bytes"
	"errors"
)

const (
	INSERT_TOGGLE_RULE_SQL = "INSERT INTO toggle_rule(id, featureid, property, value, created, expires, enabled) values ($1,$2,$3,$4,$5,$6,$7)"
	DELETE_TOGGLE_RULE_SQL = "DELETE FROM toggle_rule WHERE toggle_rule.id = $1"
	READ_TOGGLE_RULE_WHERE_PART_SQL = "WHERE toggle_rule.id = $1"

	SEARCH_TOGGLE_RULE_SELECT_PART_SQL = "SELECT DISTINCT toggle_rule.* FROM toggle_rule "
	SEARCH_TOGGLE_RULE_FEATURE_JOIN_PART_SQL = "JOIN feature ON feature.id = toggle_rule.featureid "
	SEARCH_TOGGLE_RULE_NAME_PART_SQL = "feature.name = $1 "
	SEARCH_TOGGLE_RULE_FILTER_PART_SQL = "%s(p%d.property=$%d AND p%d.value = $%d) "
	SEARCH_TOGGLE_RULE_INNER_JOIN_SQL = "INNER JOIN toggle_rule AS p%d ON toggle_rule.featureid = p%d.featureid "
	SEARCH_ROGGLE_RULE_ENABLED_SQL = "SELECT DISTINCT feature.name, tr.property, tr.value FROM toggle_rule tr JOIN feature ON feature.id = tr.featureid WHERE feature.enabled = true and tr.enabled = true and tr.expires < $1"
)

func (fs *FeatureToggleStoreImpl) CreateToggleRule(toggleRule ToggleRule) (*string, error) {
	created := time.Now()

	tx, err := fs.db.Begin()
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to create trasaction, %v", err))
	}
	defer tx.Rollback()

	stmt, err := fs.db.Prepare(INSERT_TOGGLE_RULE_SQL)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to create prepared statement, %v", err))
	}
	defer stmt.Close()

	var id string
	if strings.Compare(toggleRule.Id, "") == 0 {
		id = uuid.NewV4().String()
	}
	for property, value := range toggleRule.Properties {
		_, err := stmt.Exec(id, toggleRule.FeatureId, property, value, created, toggleRule.Expires, toggleRule.Enabled)
		if ( err != nil) {
			return nil, errors.New(fmt.Sprintf("Failed to insert row with property '%s', %v", property, err))
		}
	}
	tx.Commit()

	return &id, nil
}

func (fs *FeatureToggleStoreImpl) ReadToggleRule(id string) (*ToggleRule, error) {
	var buffer bytes.Buffer

	buffer.WriteString(SEARCH_TOGGLE_RULE_SELECT_PART_SQL)
	buffer.WriteString(READ_TOGGLE_RULE_WHERE_PART_SQL)

	stmt, err := fs.db.Prepare(buffer.String())
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to create prepared statement, %v", err))
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to run query, %v", err))
	}
	defer rows.Close()

	toggleRules, err := rowsToToggleRule(rows)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get row data, %v", err))
	}
	if len(toggleRules) > 0 {
		return &toggleRules[0], nil
	}
	return nil, nil

}

func (fs *FeatureToggleStoreImpl) DeleteToggleRule(id string) (*bool, error) {

	stmt, err := fs.db.Prepare(DELETE_TOGGLE_RULE_SQL)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to create delete prepared statement, %v", err))
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("Failed to delete toggle rule, %v", err))
	}
	rowCount, err := res.RowsAffected()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to get rowsAffected, %v", err))
	}
	b := rowCount > 0
	return &b, nil
}

func (fs *FeatureToggleStoreImpl) SearchToggleRule(name *string, filter Filter) (*[]ToggleRule, error) {
	var buffer bytes.Buffer

	buffer.WriteString(SEARCH_TOGGLE_RULE_SELECT_PART_SQL)
	if ( name != nil) {
		buffer.WriteString(SEARCH_TOGGLE_RULE_FEATURE_JOIN_PART_SQL)
	}
	for i := 0; i < len(filter); i++ {
		buffer.WriteString(getInnerJoinLine(i))
	}

	buffer.WriteString("WHERE ")
	if ( name != nil) {
		buffer.WriteString(SEARCH_TOGGLE_RULE_NAME_PART_SQL)
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

	res, err := rowsToToggleRule(rows)
	if ( err != nil) {
		return nil, err
	}
	return &res, nil
}

func (fs *FeatureToggleStoreImpl) GetEnabledToggleRules() (*[]featuretree.ToggleRule, error) {
	var buffer bytes.Buffer

	searchQuery := buffer.String()
	fmt.Println(searchQuery)

	stmt, err := fs.db.Prepare(SEARCH_ROGGLE_RULE_ENABLED_SQL)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("GetEnabledToggleRules: Failed to create prepared statement, %v", err))
	}
	defer stmt.Close()

	now := time.Now()
	rows, err := stmt.Query(now)
	if ( err != nil) {
		return nil, errors.New(fmt.Sprintf("GetEnabledToggleRules: Failed to run query, %v", err))
	}
	defer rows.Close()

	res, err := rowsToTreeToggleRule(rows)
	if ( err != nil) {
		return nil, err
	}
	return &res, nil
}

func getParams(featurename *string, filters *Filter) []interface{} {
	params := make([]interface{}, 0)
	if featurename != nil {
		params = append(params, featurename)
	}
	if ( filters != nil) {
		for propertyName, propertyValue := range *filters {
			params = append(params, propertyName, propertyValue)
		}
	}
	return params
}

func rowsToToggleRule(rows *sql.Rows) ([]ToggleRule, error) {
	ruleMap := make(map[string]ToggleRule)
	for rows.Next() {
		var id string
		var featureid string
		var property string
		var value string
		var created string
		var expires string
		var enabled bool
		err := rows.Scan(&id, &featureid, &property, &value, &created, &expires, &enabled)
		if ( err != nil) {
			return nil, errors.New(fmt.Sprintf("Failed to scan row, %v", err))
		}
		rule, ok := ruleMap[featureid]
		if !ok {
			props := make(Properties)
			rule = ToggleRule{Id:id, FeatureId:featureid, Properties:props}
			ruleMap[featureid] = rule
		}
		rule.Properties[property] = value
	}

	result := []ToggleRule{}
	for _, rule := range ruleMap {
		result = append(result, rule)
	}
	return result, nil
}

func rowsToTreeToggleRule(rows *sql.Rows) ([]featuretree.ToggleRule, error) {
	ruleMap := make(map[string]featuretree.ToggleRule)
	for rows.Next() {
		var featurename string
		var property string
		var value string
		err := rows.Scan(&featurename, &property, &value)
		if ( err != nil) {
			return nil, errors.New(fmt.Sprintf("Failed to scan row, %v", err))
		}
		rule, ok := ruleMap[featurename]
		if !ok {
			props := make(featuretree.Properties)
			rule = featuretree.ToggleRule{Name:featurename, Properties:props}

			ruleMap[featurename] = rule
		}
		rule.Properties[property] = value
	}

	result := []featuretree.ToggleRule{}
	for _, rule := range ruleMap {
		result = append(result, rule)
	}
	return result, nil
}

func getPropertyFilterLine(i int, skipStartingAnd bool, offset int) string {

	var prefix string
	if ( skipStartingAnd) {
		prefix = ""
	} else {
		prefix = "AND"
	}
	return fmt.Sprintf(SEARCH_TOGGLE_RULE_FILTER_PART_SQL, prefix, i, offset + i * 2, i, offset + 1 + i * 2)
}

func getInnerJoinLine(i int) string {
	return fmt.Sprintf(SEARCH_TOGGLE_RULE_INNER_JOIN_SQL, i, i)
}

