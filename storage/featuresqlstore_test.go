package storage

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestFeatureStore_SearchFeature(t *testing.T) {
	fs := NewFeatureStore()

	err := fs.Open()
	if err != nil {
		fmt.Printf("%v\n", err)
		panic("Failed to open database")
	}
	defer fs.Close()

	filter := make(Filter)
	filter["prop1"] = "val1"
	filter["prop2"] = "val2"
	features := fs.SearchFeature("name1", filter)

	assert.Equal(t, 1, len(features), "Result shall contain one feature")
	assert.Equal(t, "val1", features[0].properties["prop1"], "prop1 shall have value 'val1'")
	assert.Equal(t, "val2", features[0].properties["prop2"], "prop2 shall have value 'val2'")

}

func PrintFeatures( features []Feature) {
	fmt.Printf( "Features: %d\n", len(features))
	for _, feature := range features {
		fmt.Printf("%v\n", feature)
	}
}