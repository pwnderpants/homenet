package templates

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Templates holds all application templates
type Templates struct {
	templates map[string]*template.Template
}

// New creates a new template manager
func New() *Templates {
	return &Templates{
		templates: make(map[string]*template.Template),
	}
}

// LoadTemplates loads all templates from the web/templates directory
func (t *Templates) LoadTemplates() error {
	templatesDir := "web/templates"

	err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".html" {
			name := filepath.Base(path)
			tmpl, err := template.ParseFiles(path)
			if err != nil {
				return err
			}

			t.templates[name] = tmpl
			log.Printf("Loaded template: %s", name)
		}

		return nil
	})

	return err
}

// GetTemplate returns a template by name
func (t *Templates) GetTemplate(name string) (*template.Template, bool) {
	tmpl, exists := t.templates[name]
	return tmpl, exists
}

// ExecuteTemplate executes a template with the given data
func (t *Templates) ExecuteTemplate(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, exists := t.GetTemplate(name)
	if !exists {
		return fmt.Errorf("template %s not found", name)
	}

	return tmpl.Execute(w, data)
}
