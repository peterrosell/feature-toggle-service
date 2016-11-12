package storage

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"time"
)

func TestFeatureStoreImpl_CreateFeature(t *testing.T) {
	fs := NewFeatureStore()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	feature := CreateFeature(randomSufix("Name-"), true, "prop 1", "val 1", "prop 2", "val 2", "prop 3", "val 3")
	id, err := fs.CreateFeature(feature)

	assert.NotNil(t, id, fmt.Sprintf("shall get an id in return, %v", err))
}

func TestFeatureStoreImpl_DeleteFeature(t *testing.T) {
	fs := NewFeatureStore()

	err := fs.Open()
	if ( err != nil) {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	name := randomSufix("name-")
	prop1 := "prop1"
	val1 := "val1"
	prop2 := "prop2"
	val2 := "val2"

	// setup database
	_, err = fs.CreateFeature(CreateFeature(name, true, prop1, val1, prop2, val2))
	require.Nil(t, err, "Failed to create feature, %v", err)
	res, err := fs.DeleteFeature(name)
	require.Nil(t, err, "Failed to delete feature %v", err)
	assert.True(t, *res, "Should get true as result")

}

func TestFeatureStoreImpl_ReadFeature(t *testing.T) {
	fs := NewFeatureStore()

	err := fs.Open()
	if ( err != nil) {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	name := randomSufix("name-")
	prop1 := "prop1"
	val1 := "val1"
	prop2 := "prop2"
	val2 := "val2"

	// setup database
	fs.CreateFeature(CreateFeature(name, true, prop1, val1, prop2, val2))

	feature, err := fs.ReadFeature(name)

	require.Nil(t, err, "Should not get an error, %v", err)
	require.NotNil(t, feature, "Result shall contain a feature")
	assert.Equal(t, val1, feature.properties[prop1], prop1 + " shall have value '" + val1 + "'")
	assert.Equal(t, val2, feature.properties[prop2], prop2 + " shall have value '" + val2 + "'")

}

func TestFeatureStore_SearchFeature(t *testing.T) {
	fs := NewFeatureStore()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	name := randomSufix("name-")
	prop1 := "prop1"
	val1 := "val1"
	prop2 := "prop2"
	val2 := "val2"

	// setup database
	fs.CreateFeature(CreateFeature(name, true, prop1, val1, prop2, val2))

	filter := make(Filter)
	filter[prop1] = val1
	filter[prop2] = val2
	features, err := fs.SearchFeature(&name, filter)

	PrintFeatures(features)
	require.Equal(t, 1, len(features), fmt.Sprintf("Result shall contain one feature, %v", features))
	assert.Equal(t, val1, features[0].properties[prop1], prop1 + " shall have value '" + val1 + "'")
	assert.Equal(t, val2, features[0].properties[prop2], prop2 + " shall have value '" + val2 + "'")

}

func TestFeatureStore_SearchFeature__no_name(t *testing.T) {
	fs := NewFeatureStore()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	name := randomSufix("name-")
	prop1 := "prop1"
	val1 := "val1"
	prop2 := "prop2"
	val2 := "val2"

	// setup database
	fs.CreateFeature(CreateFeature(name, true, prop1, val1, prop2, val2))

	filter := make(Filter)
	filter[prop1] = val1
	filter[prop2] = val2
	features, err := fs.SearchFeature(nil, filter)

	PrintFeatures(features)
	require.NotZero(t, len(features), "Result should contain one feature, %v", features)
	assert.NotZero(t, len(features[0].properties), "Feature '%s' should have properties.", features[0].name)

}

func randomSufix(text string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s%d", text, r.Int())
}

func CreateFeature(name string, enabled bool, propArgs... string) Feature {
	props := make(Properties)
	for i := 0; i < len(propArgs); i += 2 {
		props[propArgs[i]] = propArgs[i + 1]
	}

	feature := Feature{name:name, properties:props, enabled:enabled}
	return feature
}

func PrintFeatures(features []Feature) {
	fmt.Printf("Features: %d\n", len(features))
	for _, feature := range features {
		fmt.Printf("%v\n", feature)
	}
}