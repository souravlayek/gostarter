package main

import (
	"embed"
	"text/template"
)

//go:embed templates/**
var templateFiles embed.FS

//go:embed packages/**
var packageFiles embed.FS

// LoadTemplate loads a template by name
func LoadTemplate(name string) (*template.Template, error) {
	return template.ParseFS(templateFiles, "templates/"+name)
}