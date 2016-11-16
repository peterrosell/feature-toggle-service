package featuretree

import "testing"
import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"fmt"
)

var propertyNamesTest = []struct {
	names        []string // input
	expected string // expected result
}{
	{[]string{"prop 1"}, "propertyNames: [prop 1,]\n"},
	{[]string{"prop 1","prop 2"}, "propertyNames: [prop 1,prop 2,]\n"},
	{[]string{"prop 1","prop 2", "prop båt"}, "propertyNames: [prop 1,prop 2,prop båt,]\n"},
}

func TestNewFeatureTree(t *testing.T) {

	for _, tt := range propertyNamesTest {
		actual := NewFeatureTree(tt.names)
		fmt.Print(actual.String())
		assert.Equal(t, tt.expected, actual.String(), "should be equal")
	}
}

func TestGetOrCreateNode( t *testing.T) {
	node := NewNode("alice")

	assert.Equal(t, "alice", node.value, "value shall be set correctly")
	assert.Nil(t, node.features)
	assert.Nil(t, node.nodes)
}

func TestAddFeature_leaf(t *testing.T) {
	propertyName := "username"
	propertyValue := "adam"
	featureName := "feature 1"

	node := NewNode(propertyValue)

	props := Properties{}
	props[propertyName] = propertyValue

	feature := ToggleRule{featureName, props}

	node.addFeature([]string{}, feature)

	assert.Equal(t, propertyValue, node.value, "Node shall have value " + propertyValue)
	assert.True(t, contains(node.features, featureName), "Node shall have feature" + featureName)
}

func TestAddFeature_simple(t *testing.T) {
	propertyName := "username"
	propertyValue := "adam"

	tree := NewFeatureTree([]string{propertyName})

	props := Properties{}
	props[propertyName] = propertyValue

	feature := ToggleRule{"feature 1", props}

	tree.AddFeature(feature)

	assert.Equal(t, []string{propertyName}, tree.propertyNames, "should find property " + propertyName + " on tree")

	require.NotNil(t, tree.root.nodes[propertyValue], "should have node on value '" + propertyValue + "'")
	assert.Equal(t, "adam", tree.root.nodes[propertyValue].value, "should find value 'adam' on node")
	require.Equal(t, 1, len(tree.root.nodes[propertyValue].features), "number of features should be 1")
	assert.Equal(t, tree.root.nodes[propertyValue].features[0], "feature 1", "should find feature 1 on node")
}

func TestAddFeature_2_props__1_feature_prop(t *testing.T) {
	property1Name := "username"
	property2Name := "usertype"

	propertyValue := "adam"

	propertyNames := []string{property1Name, property2Name}
	tree := NewFeatureTree(propertyNames)

	//require.FailNow(t, "Stop here")
	props := Properties{}
	props[property1Name] = propertyValue

	feature := ToggleRule{"feature 1", props}

	tree.AddFeature(feature)

	assert.Equal(t, propertyNames, tree.propertyNames, "should find value 'username' and 'usertype' on tree")

	require.NotNil(t, tree.root.nodes[propertyValue], "should have node on value '" + propertyValue + "'")
	assert.Equal(t, "adam", tree.root.nodes[propertyValue].value, "should find value 'adam' on node")

	require.Equal(t, 1, len(tree.root.nodes[propertyValue].nodes), "number of node should be 1")
	require.NotNil(t, tree.root.nodes[propertyValue].nodes["*"], "should have a node on value '*'")
	assert.Equal(t, tree.root.nodes[propertyValue].nodes["*"].features[0], "feature 1", "should find feature 1 on node")
}

func TestAddFeature_3_props__1_feature_prop(t *testing.T) {
	property0Name := "userid"
	property1Name := "username"
	property2Name := "usertype"

	propertyValue := "adam"

	propertyNames := []string{property0Name, property1Name, property2Name}
	tree := NewFeatureTree(propertyNames)


	props := Properties{}
	props[property1Name] = propertyValue

	feature := ToggleRule{"feature 1", props}

	tree.AddFeature(feature)

	assert.Equal(t, propertyNames, tree.propertyNames, "should find value 'username' and 'usertype' on tree")

	require.NotNil(t, tree.root.nodes["*"], "should have node on value '*'")

	require.NotNil(t, tree.root.nodes["*"].nodes[propertyValue], "should have node on value '" + propertyValue + "'")
	assert.Equal(t, "adam", tree.root.nodes["*"].nodes[propertyValue].value, "should find value 'adam' on node")

	require.Equal(t, 1, len(tree.root.nodes["*"].nodes[propertyValue].nodes), "number of node should be 1")
	require.NotNil(t, tree.root.nodes["*"].nodes[propertyValue].nodes["*"], "should have a node on value '*'")
	assert.Equal(t, tree.root.nodes["*"].nodes[propertyValue].nodes["*"].features[0], "feature 1", "should find feature 1 on node")
}

func TestAddFeature_3_props__2_feature_prop(t *testing.T) {
	property0Name := "userid"
	property1Name := "username"
	property2Name := "usertype"

	property1Value := "adam"
	property1Value2 := "bert"

	featureName1 := "feature 1"
	featureName2 := "feature 2"

	propertyNames := []string{property0Name, property1Name, property2Name}
	tree := NewFeatureTree(propertyNames)


	props := Properties{}
	props[property1Name] = property1Value
	feature := ToggleRule{featureName1, props}
	tree.AddFeature(feature)

	props2 := Properties{}
	props2[property1Name] = property1Value2
	feature2 := ToggleRule{featureName2, props2}
	tree.AddFeature(feature2)

	assert.Equal(t, propertyNames, tree.propertyNames, "should find value 'username' and 'usertype' on tree")

	require.NotNil(t, tree.root.nodes["*"], "should have node on value '*'")

	require.NotNil(t, tree.root.nodes["*"].nodes[property1Value], "should have node on value '" + property1Value + "'")
	assert.Equal(t, property1Value, tree.root.nodes["*"].nodes[property1Value].value, "should find value 'adam' on node")

	require.Equal(t, 1, len(tree.root.nodes["*"].nodes[property1Value].nodes), "number of node should be 1")
	require.NotNil(t, tree.root.nodes["*"].nodes[property1Value].nodes["*"], "should have a node on value '*'")
	assert.Equal(t, tree.root.nodes["*"].nodes[property1Value].nodes["*"].features[0], featureName1, "should find feature 1 on node")

	//tree.root.nodes["*"].nodes.printMap()
	require.NotNil(t, tree.root.nodes["*"].nodes[property1Value2], "should have node on value '" + property1Value + "'")
	assert.Equal(t, property1Value2, tree.root.nodes["*"].nodes[property1Value2].value, "should find value '" + property1Value2 + "' on node")

	require.Equal(t, 1, len(tree.root.nodes["*"].nodes[property1Value2].nodes), "number of node should be 1")
	require.NotNil(t, tree.root.nodes["*"].nodes[property1Value2].nodes["*"], "should have a node on value '*'")
	assert.Equal(t, tree.root.nodes["*"].nodes[property1Value2].nodes["*"].features[0], featureName2, "should find feature 1 on node")
}

func TestAddFeature_3_props__3_feature_prop(t *testing.T) {
	property0Name := "userid"
	property1Name := "username"
	property2Name := "usertype"

	property1Value := "adam"
	property2Value2 := "beta"

	featureName1 := "feature 1"
	featureName2 := "feature 2"

	propertyNames := []string{property0Name, property1Name, property2Name}
	tree := NewFeatureTree(propertyNames)


	props := Properties{}
	props[property1Name] = property1Value
	feature := ToggleRule{featureName1, props}
	tree.AddFeature(feature)

	props2 := Properties{}
	props2[property1Name] = property1Value
	props2[property2Name] = property2Value2
	feature2 := ToggleRule{featureName2, props2}
	tree.AddFeature(feature2)

	print(tree.String())
	assert.Equal(t, propertyNames, tree.propertyNames, "should find value 'username' and 'usertype' on tree")

	require.NotNil(t, tree.root.nodes["*"], "should have node on value '*'")

	require.NotNil(t, tree.root.nodes["*"].nodes[property1Value], "should have node on value '" + property1Value + "'")
	assert.Equal(t, property1Value, tree.root.nodes["*"].nodes[property1Value].value, "should find value 'adam' on node")

	require.Equal(t, 2, len(tree.root.nodes["*"].nodes[property1Value].nodes), "number of node should be 2")
	require.NotNil(t, tree.root.nodes["*"].nodes[property1Value].nodes["*"], "should have a node on value '*'")
	assert.Equal(t, tree.root.nodes["*"].nodes[property1Value].nodes["*"].features[0], featureName1, "should find feature 1 on node")

	//tree.root.nodes["*"].nodes.printMap()
	require.NotNil(t, tree.root.nodes["*"].nodes[property1Value], "should have node on value '" + property1Value + "'")
	assert.Equal(t, property1Value, tree.root.nodes["*"].nodes[property1Value].value, "should find value '" + property1Value + "' on node")

	require.NotNil(t, tree.root.nodes["*"].nodes[property1Value].nodes[property2Value2], "should have a node on value '*'")
	assert.Equal(t, tree.root.nodes["*"].nodes[property1Value].nodes[property2Value2].features[0], featureName2, "should find '" + featureName2 + "' on node")
}


func TestFindFeatures(t *testing.T) {

	featureName1 := "feature 1"

	tree := createFeatureTree()

	properties := Properties{}
	properties["username"] = "adam"

	features := tree.FindFeatures(properties)

	assert.True(t, contains(features, featureName1), "Found features should contain '" + featureName1 + "'")

}

func TestFindFeatures_beta(t *testing.T) {

	tree := createFeatureTree()

	properties := Properties{}
	properties["usertype"] = "beta"

	features := tree.FindFeatures(properties)

	assert.Equal(t, 0, len(features), "No features should be found for only usertype=beta")

}

func TestFindFeatures_adam_beta(t *testing.T) {

	featureName1 := "feature 1"
	featureName2 := "feature 2"

	tree := createFeatureTree()

	properties := Properties{}
	properties["username"] = "adam"
	properties["usertype"] = "beta"

	features := tree.FindFeatures(properties)

	assert.True(t, contains(features, featureName1), "Found features should contain '" + featureName1 + "'")
	assert.True(t, contains(features, featureName2), "Found features should contain '" + featureName2 + "'")
}

func createFeatureTree() *ToggleRuleTree {
	property0Name := "userid"
	property1Name := "username"
	property2Name := "usertype"

	property1Value := "adam"
	property2Value2 := "beta"

	featureName1 := "feature 1"
	featureName2 := "feature 2"

	propertyNames := []string{property0Name, property1Name, property2Name}
	tree := NewFeatureTree(propertyNames)


	props := Properties{}
	props[property1Name] = property1Value
	feature := ToggleRule{featureName1, props}
	tree.AddFeature(feature)

	props2 := Properties{}
	props2[property1Name] = property1Value
	props2[property2Name] = property2Value2
	feature2 := ToggleRule{featureName2, props2}
	tree.AddFeature(feature2)
	return tree
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (m NodeMap) printMap() {
	for k, v := range m {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}
}