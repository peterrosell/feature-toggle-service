package storage

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"time"
)

func TestFeatureToggleStoreImpl_CreateToggleRule(t *testing.T) {
	var fs FeatureToggleStore = NewFeatureToggleStoreImpl()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	feature := NewFeature(randomSufix("Name-"), true, "f description")
	featureId, err := fs.CreateFeature(*feature); require.NotNil(t, featureId, "Should get featureId, %v", err)

	prop1 := Property{randomSufix("prop-"), "p description 1"}
	prop2 := Property{randomSufix("prop-"), "p description 2"}
	prop3 := Property{randomSufix("prop-"), "p description 3"}
	p, err := fs.CreateProperty(prop1); require.NotNil(t, p, "Should get propertyName, %v", err)
	p, err = fs.CreateProperty(prop2); require.NotNil(t, p, "Should get propertyName, %v", err)
	p, err = fs.CreateProperty(prop3); require.NotNil(t, p, "Should get propertyName, %v", err)

	id, err := fs.CreateToggleRule(*NewToggleRule(*featureId, true, prop1.Name, "val 1", prop2.Name, "val 2", prop3.Name, "val 3"))

	assert.NotNil(t, id, fmt.Sprintf("shall get an id in return, %v", err))
}

func TestFeatureToggleStoreImpl_DeleteToggleRule(t *testing.T) {
	var fs FeatureToggleStore = NewFeatureToggleStoreImpl()

	err := fs.Open()
	if ( err != nil) {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	feature := NewFeature(randomSufix("Name-"), true, "f description")
	featureId, err := fs.CreateFeature(*feature); require.NotNil(t, featureId, "Should get featureId, %v", err)

	prop1 := Property{randomSufix("prop-"), "p description 1"}
	prop2 := Property{randomSufix("prop-"), "p description 2"}
	prop3 := Property{randomSufix("prop-"), "p description 3"}
	p, err := fs.CreateProperty(prop1); require.NotNil(t, p, "Should get propertyName, %v", err)
	p, err = fs.CreateProperty(prop2); require.NotNil(t, p, "Should get propertyName, %v", err)
	p, err = fs.CreateProperty(prop3); require.NotNil(t, p, "Should get propertyName, %v", err)

	// setup database
	ruleId, err := fs.CreateToggleRule(*NewToggleRule(*featureId, true, prop1.Name, "val1", prop2.Name, "val2"))
	fmt.Printf("ruleID %s\n", *ruleId)
	require.Nil(t, err, "Failed to create toggle rule, %v", err)
	res, err := fs.DeleteToggleRule(*ruleId)
	require.Nil(t, err, "Failed to delete toggle rule %v", err)
	assert.True(t, *res, "Should get true as result")

}

func TestFeatureToggleStoreImpl_ReadToggleRule(t *testing.T) {
	var fs FeatureToggleStore = NewFeatureToggleStoreImpl()

	err := fs.Open()
	if ( err != nil) {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	val1 := "val1"
	val2 := "val2"

	feature := NewFeature(randomSufix("Name-"), true, "f description")
	featureId, err := fs.CreateFeature(*feature); require.NotNil(t, featureId, "Should get featureId, %v", err)

	prop1 := Property{randomSufix("prop-"), "p description 1"}
	prop2 := Property{randomSufix("prop-"), "p description 2"}
	p, err := fs.CreateProperty(prop1); require.NotNil(t, p, "Should get propertyName, %v", err)
	p, err = fs.CreateProperty(prop2); require.NotNil(t, p, "Should get propertyName, %v", err)

	// setup database
	ruleId, err := fs.CreateToggleRule(*NewToggleRule(*featureId, true, prop1.Name, val1, prop2.Name, val2))

	toggleRule, err := fs.ReadToggleRule(*ruleId)

	require.Nil(t, err, "Should not get an error, %v", err)
	require.NotNil(t, toggleRule, "Result shall contain a toggle rule")
	assert.Equal(t, val1, toggleRule.Properties[prop1.Name], prop1.Name + " shall have value '" + val1 + "'")
	assert.Equal(t, val2, toggleRule.Properties[prop2.Name], prop2.Name + " shall have value '" + val2 + "'")
}

func TestFeatureStore_SearchToggleRule(t *testing.T) {
	var fs FeatureToggleStore = NewFeatureToggleStoreImpl()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	val1 := "val1"
	val2 := "val2"

	feature := NewFeature(randomSufix("Name-"), true, "f description")
	featureId, err := fs.CreateFeature(*feature); require.NotNil(t, featureId, "Should get featureId, %v", err)

	prop1 := Property{randomSufix("prop-"), "p description 1"}
	prop2 := Property{randomSufix("prop-"), "p description 2"}
	p, err := fs.CreateProperty(prop1); require.NotNil(t, p, "Should get propertyName, %v", err)
	p, err = fs.CreateProperty(prop2); require.NotNil(t, p, "Should get propertyName, %v", err)

	// setup database
	fs.CreateToggleRule(*NewToggleRule(*featureId, true, prop1.Name, val1, prop2.Name, val2))

	filter := make(Filter)
	filter[prop1.Name] = val1
	filter[prop2.Name] = val2
	toggleRules, err := fs.SearchToggleRule(&feature.Name, filter)

	PrintFeatures(*toggleRules)
	require.Equal(t, 1, len(*toggleRules), fmt.Sprintf("Result shall contain one toggle rule, %v", *toggleRules))
	assert.Equal(t, val1, (*toggleRules)[0].Properties[prop1.Name], prop1.Name + " shall have value '" + val1 + "'")
	assert.Equal(t, val2, (*toggleRules)[0].Properties[prop2.Name], prop2.Name + " shall have value '" + val2 + "'")

}

func TestFeatureStore_SearchToggleRule__no_name(t *testing.T) {
	var fs FeatureToggleStore = NewFeatureToggleStoreImpl()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	val1 := "val1"
	val2 := "val2"

	feature := NewFeature(randomSufix("Name-"), true, "f description")
	featureId, err := fs.CreateFeature(*feature); require.NotNil(t, featureId, "Should get featureId, %v", err)

	prop1 := Property{randomSufix("prop-"), "p description 1"}
	prop2 := Property{randomSufix("prop-"), "p description 2"}
	p, err := fs.CreateProperty(prop1); require.NotNil(t, p, "Should get propertyName, %v", err)
	p, err = fs.CreateProperty(prop2); require.NotNil(t, p, "Should get propertyName, %v", err)

	// setup database
	fs.CreateToggleRule(*NewToggleRule(*featureId, true, prop1.Name, val1, prop2.Name, val2))

	filter := make(Filter)
	filter[prop1.Name] = val1
	filter[prop2.Name] = val2
	features, err := fs.SearchToggleRule(nil, filter)

	PrintFeatures(*features)
	require.NotZero(t, len(*features), "Result should contain one toggle rule, %v", *features)
	assert.NotZero(t, len((*features)[0].Properties), "Feature '%s' should have properties.", (*features)[0].FeatureId)
}

func TestFeatureToggleStoreImpl_GetEnabledToggleRules(t *testing.T) {
	var fs FeatureToggleStore = NewFeatureToggleStoreImpl()

	err := fs.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open database, %v", err))
	}
	defer fs.Close()

	val1 := "val1"
	val2 := "val2"

	feature := NewFeature(randomSufix("Name-"), true, "f description")
	featureId, err := fs.CreateFeature(*feature); require.NotNil(t, featureId, "Should get featureId, %v", err)

	prop1 := Property{randomSufix("prop-"), "p description 1"}
	prop2 := Property{randomSufix("prop-"), "p description 2"}
	p, err := fs.CreateProperty(prop1); require.NotNil(t, p, "Should get propertyName, %v", err)
	p, err = fs.CreateProperty(prop2); require.NotNil(t, p, "Should get propertyName, %v", err)

	// setup database
	fs.CreateToggleRule(*NewToggleRule(*featureId, true, prop1.Name, val1, prop2.Name, val2))

	rules, err := fs.GetEnabledToggleRules()

	require.NotNil(t, rules, "Should get rules, %v\n", err)
	require.True(t, len(*rules) > 0, "Should get one or more rules")
	assert.NotEmpty(t, (*rules)[0].Name, "Rule should have a name")
	assert.True(t, len((*rules)[0].Properties) > 0, "Rule should have one or more properties")

}

func randomSufix(text string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s%d", text, r.Int())
}

func PrintFeatures(toggleRules []ToggleRule) {
	fmt.Printf("Toggle rules: %d\n", len(toggleRules))
	for _, toggleRule := range toggleRules {
		fmt.Printf("%v\n", toggleRule)
	}
}