package views

import (
	"html/template"
	"net/http"
	"options-api/models"
)

type PageData struct {
	Fields  []models.Field
	Success bool
}

type View struct {
	template *template.Template
}

func NewView() (*View, error) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return nil, err
	}
	return &View{template: tmpl}, nil
}

func (v *View) RenderHome(w http.ResponseWriter, fields []models.Field, success bool) error {
	data := PageData{
		Fields:  fields,
		Success: success,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return v.template.ExecuteTemplate(w, "index.html", data)
}
