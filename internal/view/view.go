package view

import (
	"fmt"
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
	options := &ViewOptions{
		Path:      env.GetString("VIEW_PATH", "./views"),
		Extension: env.GetString("VIEW_EXTENSION", ".html"),
	}

	rootTemplate := template.New("root")
	err := loadTemplatesConcurrently(options, rootTemplate)
	if err != nil {
		return err
	}

	templateCache.Store("root", rootTemplate)
	return nil
}

func loadTemplatesConcurrently(options *ViewOptions, rootTemplate *template.Template) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

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

		wg.Add(1)
		go func(path, templateName string) {
			defer wg.Done()
			if _, err := rootTemplate.New(templateName).ParseFiles(path); err != nil {
				log.Printf("Failed to parse template %s: %v", path, err)
				select {
				case errCh <- err:
				default:
				}
				return
			}
			log.Printf("Template loaded: %s", templateName)
		}(path, templateName)

		return nil
	})

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		return err
	}

	return err
}

func logAndRespondError(w http.ResponseWriter, message string, statusCode int) {
	log.Println(message)
	http.Error(w, message, statusCode)
}

func RenderTemplate(w http.ResponseWriter, templateName string, data ...interface{}) {
	rootTemplate, err := getRootTemplate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var templateData interface{}
	if len(data) > 0 {
		templateData = data[0]
	}

	err = executeTemplate(w, rootTemplate, templateName, templateData)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}

func getRootTemplate() (*template.Template, error) {
	value, ok := templateCache.Load("root")
	if !ok {
		return nil, fmt.Errorf("root template not found")
	}

	rootTemplate, ok := value.(*template.Template)
	if !ok {
		return nil, fmt.Errorf("template cache corrupted")
	}

	return rootTemplate, nil
}

func executeTemplate(w http.ResponseWriter, rootTemplate *template.Template, templateName string, templateData interface{}) error {
	var builder strings.Builder
	err := rootTemplate.ExecuteTemplate(&builder, templateName, templateData)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(builder.String()))
	return err
}

func Render404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	RenderTemplate(w, "page/404")
}

func Render500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	RenderTemplate(w, "page/500")
}
