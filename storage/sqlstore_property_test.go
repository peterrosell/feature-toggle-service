package storage

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
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
	require.NotNil(t, *p, "Shoudl get property from property name '%s', err", propertyName, err)
	assert.Equal(t, property.Name, p.Name, "Should get property name, %v", err)
	assert.Equal(t, property.Description, p.Description, "Should get property description, %v", err)
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

	res,err := fs.DeleteProperty(*propertyName)

	require.True(t, *res, "Should get true from delete operation for property '%s', %v", propertyName, err)
}


