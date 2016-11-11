package api

type Properties map[string]string

type Feature struct {
	name       string
	properties Properties
}

