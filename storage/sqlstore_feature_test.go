package storage

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestFeatureToggleStoreImpl_CreateFeature(t *testing.T) {
	var fs FeatureToggleStore = NewFeatureToggleStoreImpl()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	feature := NewFeature(randomSufix("Feature-"), true, "f description")
	featureId, err := fs.CreateFeature(*feature);

	require.NotNil(t, featureId, "Should get featureId, %v", err)
}

func TestFeatureToggleStoreImpl_ReadFeature(t *testing.T) {
	var fs FeatureToggleStore = NewFeatureToggleStoreImpl()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	feature := NewFeature(randomSufix("Feature-"), true, "f description")
	featureId, err := fs.CreateFeature(*feature);

	require.NotNil(t, featureId, "Should get featureId, %v", err)

	f, err := fs.ReadFeature(*featureId)
	require.NotNil(t, *f, "Should get feature from featureId %s, err", featureId, err)
	assert.Equal(t, feature.Name, f.Name, "Should get feature name, %v", err)
	assert.Equal(t, feature.Description, f.Description, "Should get feature description, %v", err)
	assert.Equal(t, feature.Enabled, f.Enabled, "Should get feature enabled, %v", err)
}

func TestFeatureToggleStoreImpl_ReadFeatureByName(t *testing.T) {
	var fs FeatureToggleStore = NewFeatureToggleStoreImpl()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	featureName := randomSufix("Feature-")
	feature := NewFeature(featureName, true, "f description")
	featureId, err := fs.CreateFeature(*feature);

	require.NotNil(t, featureId, "Should get featureId, %v", err)

	f, err := fs.ReadFeatureByName(featureName)
	require.NotNil(t, *f, "Should get feature from featureId %s, err", featureId, err)
	assert.Equal(t, feature.Name, f.Name, "Should get feature name, %v", err)
	assert.Equal(t, feature.Description, f.Description, "Should get feature description, %v", err)
	assert.Equal(t, feature.Enabled, f.Enabled, "Should get feature enabled, %v", err)
}

func TestFeatureToggleStoreImpl_ReadFeatureByName__unknown(t *testing.T) {
	var fs FeatureToggleStore = NewFeatureToggleStoreImpl()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	f, err := fs.ReadFeatureByName("unknown feature")
	assert.Nil(t, f, "Should not get a feature, %v", f)
	assert.Nil(t, err, "Should not get an error, %v", err)
}

func TestFeatureToggleStoreImpl_DeleteFeature(t *testing.T) {
	var fs FeatureToggleStore = NewFeatureToggleStoreImpl()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	feature := NewFeature(randomSufix("Feature-"), true, "f description")
	featureId, err := fs.CreateFeature(*feature);

	require.NotNil(t, featureId, "Should get featureId, %v", err)

	res,err := fs.DeleteFeature(*featureId)

	require.True(t, *res, "Should get true from delete operation for featureId %s, %v", featureId, err)
}


