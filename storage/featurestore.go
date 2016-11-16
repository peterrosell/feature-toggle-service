package storage

import (
	"time"
	"github.com/satori/go.uuid"
	"github.com/peterrosell/feature-toggle-service/featuretree"
)

type Filter map[string]string

type Properties map[string]string

type ToggleRule struct {
	Id         string
	FeatureId  string
	Enabled    bool
	Created    time.Time
	Expires    time.Time
	Properties Properties
}

type Feature struct {
	Id          string
	Name        string
	Enabled     bool
	Description string
}

type Property struct {
	Name        string
	Description string
}

type FeatureToggleStore interface {
	GetEnabledToggleRules() (*[]featuretree.ToggleRule, error)

	CreateFeature(feature Feature) (*string, error)
	ReadFeature(id string) (*Feature, error)
	ReadFeatureByName(name string) (*Feature, error)
	DeleteFeature(id string) (*bool, error)
	SearchFeature(name string) (*[]Feature, error)

	CreateProperty(property Property) (*string, error)
	ReadProperty(name string) (*Property, error)
	ReadAllPropertyNames() (*[]string, error)
	DeleteProperty(name string) (*bool, error)
	SearchProperty(name string) (*[]Property, error)


	CreateToggleRule(toggleRule ToggleRule) (*string, error)
	ReadToggleRule(id string) (*ToggleRule, error)
	DeleteToggleRule(id string) (*bool, error)
	SearchToggleRule(name *string, filter Filter) (*[]ToggleRule, error)

	Open() error
	Close()
}

func NewFeature(name string, enabled bool, description string) *Feature {
	return &Feature{uuid.NewV4().String(), name, enabled, description}
}

func NewProperty(name string, description string) *Property {
	return &Property{name, description}
}

func NewToggleRule(featureId string, enabled bool, propArgs... string) *ToggleRule {
	props := make(Properties)
	for i := 0; i < len(propArgs); i += 2 {
		props[propArgs[i]] = propArgs[i + 1]
	}

	toggleRule := ToggleRule{FeatureId:featureId, Properties:props, Enabled:enabled}
	return &toggleRule
}
