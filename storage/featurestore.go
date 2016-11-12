package storage

import "time"

type Properties map[string]string

type Feature struct {
	name       string
	enabled    bool
	created    time.Time
	expires    time.Time
	properties Properties
}

type Filter map[string]string

type FeatureStore interface {
	CreateFeature( feature Feature) string
	ReadFeature( id string) Feature
	DeleteFeature(id string) bool
	SearchFeature(filter Filter) []Feature
}
