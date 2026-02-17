package web

import (
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/berkaycubuk/subtrack/internal/utils"
)

//go:embed templates/*.html
var templateFS embed.FS

var funcMap = template.FuncMap{
	"formatDate": utils.FormatDate,
	"formatPrice": func(price float64, currency string) string {
		return fmt.Sprintf("%.2f %s", price, currency)
	},
	"formatInputDate": func(t time.Time) string {
		return utils.FormatDate(t)
	},
}

func parseTemplate(name string) *template.Template {
	return template.Must(template.New("layout.html").Funcs(funcMap).ParseFS(
		templateFS, "templates/layout.html", "templates/"+name,
	))
}
