package feature_toggle_api

type Properties map[string]string

type Feature struct {
	name       string
	properties Properties
}

