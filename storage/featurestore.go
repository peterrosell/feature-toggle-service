package storage

type Properties map[string]string

type Feature struct {
	name       string
	properties Properties
}

type Filter map[string]string

type FeatureStore interface {
	CreateFeature( feature Feature) string
	ReadFeature( id string) Feature
	DeleteFeature(id string) bool
	SearchFeature(filter Filter) []Feature
}
