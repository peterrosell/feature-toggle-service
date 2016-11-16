package storage

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	"strings"
)

func TestFeatureToggleStoreImpl_CreateProperty(t *testing.T) {
	var fs FeatureToggleStore = NewFeatureToggleStoreImpl()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	property := NewProperty(randomSufix("Prop-"), "p description")
	propertyName, err := fs.CreateProperty(*property);

	require.NotNil(t, propertyName, "Should get property name, %v", err)
}

func TestFeatureToggleStoreImpl_ReadProperty(t *testing.T) {
	var fs FeatureToggleStore = NewFeatureToggleStoreImpl()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	property := NewProperty(randomSufix("Prop-"), "p description")
	propertyName, err := fs.CreateProperty(*property);

	require.NotNil(t, propertyName, "Should get property name, %v", err)

	p, err := fs.ReadProperty(*propertyName)
	require.NotNil(t, *p, "Shoudl get property from property name '%s', %v", propertyName, err)
	assert.Equal(t, property.Name, p.Name, "Should get property name, %v", err)
	assert.Equal(t, property.Description, p.Description, "Should get property description, %v", err)
}

func TestFeatureToggleStoreImpl_ReadAllPropertyNames(t *testing.T) {
	var fs FeatureToggleStore = NewFeatureToggleStoreImpl()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	propName1 := randomSufix("Prop-")
	property := NewProperty(propName1, "p description")
	propertyName, err := fs.CreateProperty(*property);

	require.NotNil(t, propertyName, "Should get property 1 name, %v", err)

	propName2 := randomSufix("Prop-")
	property = NewProperty(propName2, "p description")
	propertyName, err = fs.CreateProperty(*property);

	require.NotNil(t, propertyName, "Should get property 2 name, %v", err)

	p, err := fs.ReadAllPropertyNames()
	require.NotNil(t, *p, "Shoudl get property names, %v", err)
	assert.True(t, contains(p, propName1), "Should find property name, %v", err)
	assert.True(t, contains(p, propName2), "Should find property name, %v", err)
}

func contains(ss *[]string, str string) bool {
	for _, s := range *ss {
		if strings.Compare(s, str) == 0 {
			return true
		}
	}
	return false
}

func TestFeatureToggleStoreImpl_DeleteProperty(t *testing.T) {
	var fs FeatureToggleStore = NewFeatureToggleStoreImpl()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	property := NewProperty(randomSufix("Name-"), "p description")
	propertyName, err := fs.CreateProperty(*property);

	require.NotNil(t, propertyName, "Should get property name, %v", err)

	res, err := fs.DeleteProperty(*propertyName)

	require.True(t, *res, "Should get true from delete operation for property '%s', %v", propertyName, err)
}


