package bootstrap

import (
	"log"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/gin-gonic/gin"
)

func SetupTemplatesAndStatic(r *gin.Engine) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Unable to get current file path")
	}
	currentDir := filepath.Dir(filename)
	projectRoot := filepath.Join(currentDir, "..", "..", "..")

	// r.LoadHTMLGlob(filepath.Join(projectRoot, "templates", "*.html"))
	// r.Static("/static", filepath.Join(projectRoot, "static"))

	templatesPath := filepath.Join(projectRoot, "templates", "*.html")
	staticPath := filepath.Join(projectRoot, "static")

	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"subtract": func(a, b int) int {
			if a < b {
				return 0
			}
			return a - b
		},
		"max": func(a, b int) int {
			if a > b {
				return a
			}
			return b
		},
		"min": func(a, b int) int {
			if a < b {
				return a
			}
			return a
		},
	})

	r.LoadHTMLGlob(templatesPath)
	r.Static("/static", staticPath)
}
