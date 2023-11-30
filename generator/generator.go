package generator

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/a-h/templ"
)

func SaveTemplComponent(path string, component templ.Component) (*os.File, error) {
	file, err := generateStaticFile(path)
	if err != nil {
		return nil, err
	}
	err = component.Render(context.Background(), file)
	if err != nil {

		return nil, err
	}
	return file, nil
}

func generateStaticFile(filename string) (*os.File, error) {
	filepath := fmt.Sprintf("static/%s.html", filename)
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("failed to create output file. path: %s, error: %v", filepath, err)
		return file, err
	}
	return file, nil
}
