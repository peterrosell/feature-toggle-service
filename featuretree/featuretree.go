package featuretree

import (
	"bytes"
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
		node.features = append(node.features, rule.Name)
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

func (tree *ToggleRuleTree) AddFeature(rule ToggleRule) {
	tree.root.addFeature(tree.propertyNames, rule)
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
	buffer.WriteString("propertyNames: [")
	for _, name := range tree.propertyNames {
		buffer.WriteString(name)
		buffer.WriteString(",")
	}
	buffer.WriteString("]")
	return buffer.String()
}

func NewFeatureTree(propertyNames []string) *ToggleRuleTree {
	tree := ToggleRuleTree{Node{}, propertyNames}
	return &tree
}