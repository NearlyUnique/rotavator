package templates

import (
	"embed"
	"fmt"
	"log"

	"github.com/gobuffalo/plush"
)

//go:embed *
var content embed.FS

// ReadTemplate from embedded resources
func ReadTemplate(path string) string {
	body, err := content.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("template '%s', %v", path, err)
	}
	return string(body)
}

func RenderPage(path string, data map[string]any) string {
	ctx := plush.NewContextWith(data)
	content, err := plush.Render(ReadTemplate(path), ctx)
	if err != nil {
		log.Printf("ERROR: PAGE: %v", err)
		return err.Error()
	}

	ctx.Set("page", content)
	content, err = plush.Render(ReadTemplate("_layout.html"), ctx)
	if err != nil {
		log.Printf("ERROR: LAYOUT: %v", err)
		return err.Error()
	}

	return content
}
