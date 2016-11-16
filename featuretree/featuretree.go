package featuretree

import (
	"bytes"
	"fmt"
	"strings"
	"github.com/pkg/errors"
)

const unspecifiedProperty = "*"

type NodeMap map[string]*Node

type Node struct {
	value    string
	features []string
	nodes    NodeMap
}

type ToggleRuleTree struct {
	root          Node
	propertyNames []string
}

type Properties map[string]string

type ToggleRule struct {
	Name       string
	Properties Properties
}

func NewNode(key string) *Node {
	return &Node{value:key}
}

func (node *Node) getNodes() NodeMap {
	if ( node.nodes == nil) {
		node.nodes = make(NodeMap)
	}
	return node.nodes
}

func (node *Node) getOrCreateNode(value string) *Node {
	nextNode, ok := node.nodes[value]
	if !ok {
		newNode := NewNode(value)
		node.getNodes()[value] = newNode
		return newNode
	}
	return nextNode
}

func (node *Node) addFeature(propertyNames []string, rule ToggleRule) {
	if len(propertyNames) == 0 {
		node.features = addToFeatureList(node.features, rule.Name)
	} else {
		val, ok := rule.Properties[propertyNames[0]]
		nextNode := &Node{}
		if ok {
			nextNode = node.getOrCreateNode(val)
		} else {
			// feature does not have a property on this level, add to wildcard
			nextNode = node.getOrCreateNode("*")
		}
		(nextNode).addFeature(propertyNames[1:], rule)
	}
}
func addToFeatureList(featureList []string, feature string) []string {
	for _, f := range featureList {
		if strings.Compare(f, feature) == 0 {
			return featureList
		}
	}
	return append(featureList, feature)
}

func (tree *ToggleRuleTree) AddFeature(rule ToggleRule) error {
	err := tree.validateToggleRule(rule)
	if err != nil {
		return errors.New("Ignoring feature. " + err.Error())
	}
	tree.root.addFeature(tree.propertyNames, rule)
	return nil
}

func (tree *ToggleRuleTree) validateToggleRule(rule ToggleRule) error {
	for propName, _ := range rule.Properties {
		var found bool = false
		for _, name := range tree.propertyNames {
			fmt.Printf("propName %v == name %v\n", propName, name)
			if strings.Compare(name, propName) == 0 {
				found = true
				break
			}
		}
		if !found {
			return errors.New(fmt.Sprintf("Property '%s' is unknown.", propName))
		}
	}
	return nil
}

func (node *Node) findFeature(propertyNames []string, properties Properties) []string {
	if node.features != nil {
		return node.features
	} else {
		features := []string{}
		if len(propertyNames) == 0 {
			return []string{}
		} else {
			nextPropertyName := propertyNames[0]
			fmt.Printf("props: %v, nextName: %v\n", properties, nextPropertyName)
			if val, ok := properties[nextPropertyName]; ok {
				if nextNode, ok := node.nodes[val]; ok {
					features = nextNode.findFeature(propertyNames[1:], properties)
				}
			}
			if nextNode, ok := node.nodes[unspecifiedProperty]; ok {
				unspecFeatures := nextNode.findFeature(propertyNames[1:], properties)
				features = append(features, unspecFeatures...)
			}

		}
		return features
	}
}

func (tree *ToggleRuleTree) FindFeatures(properties Properties) []string {

	return tree.root.findFeature(tree.propertyNames, properties)
}

func (tree *ToggleRuleTree) String() string {
	var buffer bytes.Buffer
	tree.writeProperties(&buffer)
	//buffer.WriteString("\n")
	tree.root.writeToBuf(&buffer, 0)

	return buffer.String()
}

func (tree *ToggleRuleTree) writeProperties(buffer *bytes.Buffer) {
	buffer.WriteString("propertyNames: [")
	for _, name := range tree.propertyNames {
		buffer.WriteString(name)
		buffer.WriteString(",")
	}
	buffer.WriteString("]")
}

func (node *Node) writeToBuf(buf *bytes.Buffer, indent int) {
	if ( strings.Compare(node.value, "") != 0) {
		buf.WriteString(node.value)
		//buf.WriteString("\n")
	}
	if ( node.nodes != nil) {
		for _, n := range node.nodes {
			addIndent(buf, indent)
			buf.WriteString("->")
			n.writeToBuf(buf, indent + 2)
		}
	} else {
		// write features
		if indent > 0 {
			addIndent(buf, indent)
			buf.WriteString("=")
		}
		if ( len(node.features) > 0) {
			buf.WriteString(node.features[0])
		}
		for i, f := range node.features {
			if ( i > 0) {
				buf.WriteString(",")
				buf.WriteString(f)
			}
		}
		buf.WriteString("\n")

	}

}

func addIndent(buf *bytes.Buffer, indent int) {
	for i := 0; i < indent; i++ {
		buf.WriteString(" ")
	}
}

func NewFeatureTree(propertyNames []string) *ToggleRuleTree {
	tree := ToggleRuleTree{Node{}, propertyNames}
	return &tree
}