package view

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/abhilash26/tigerfly/internal/env"
)

type ViewOptions struct {
	Path      string
	Extension string
}

var templateCache = sync.Map{}

func PreloadAllTemplates() error {
	rootTemplate := template.New("root")

	options := &ViewOptions{
		Path:      env.GetString("VIEW_PATH", "./views/"),
		Extension: env.GetString("VIEW_EXTENSION", ".html"),
	}

	err := filepath.Walk(options.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), options.Extension) {
			return nil
		}

		relativePath, err := filepath.Rel(options.Path, path)
		if err != nil {
			return err
		}

		templateName := strings.TrimSuffix(relativePath, options.Extension)
		if _, err := rootTemplate.New(templateName).ParseFiles(path); err != nil {
			log.Printf("Failed to parse template %s: %v", path, err)
			return err
		}

		log.Printf("Template loaded: %s", templateName)
		return nil
	})
	if err != nil {
		return err
	}

	templateCache.Store("root", rootTemplate)
	return nil
}

func logAndRespondError(w http.ResponseWriter, message string, statusCode int) {
	log.Println(message)
	http.Error(w, message, statusCode)
}

func RenderTemplate(w http.ResponseWriter, templateName string, data ...interface{}) {
	// Load the root template from the cache
	value, ok := templateCache.Load("root")
	if !ok {
		http.Error(w, "Template Not Found", http.StatusInternalServerError)
		log.Println("Root template not found")
		return
	}

	rootTemplate, ok := value.(*template.Template)
	if !ok {
		http.Error(w, "Template Cache Corrupted", http.StatusInternalServerError)
		log.Println("Failed to assert template type from cache")
		return
	}

	var templateData interface{}
	if len(data) > 0 {
		templateData = data[0]
	}

	var builder strings.Builder
	err := rootTemplate.ExecuteTemplate(&builder, templateName, templateData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
		return
	}
	w.Write([]byte(builder.String()))
}

func Render404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	RenderTemplate(w, "page/404")
}

func Render500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	RenderTemplate(w, "page/500")
}
