package storage

import "time"

type Properties map[string]string

type Feature struct {
	name       string
	properties Properties
	created    time.Time
	expires    time.Time
	enabled    bool
}

type Filter map[string]string

type FeatureStore interface {
	CreateFeature( feature Feature) string
	ReadFeature( id string) Feature
	DeleteFeature(id string) bool
	SearchFeature(filter Filter) []Feature
}
