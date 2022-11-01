package models

import "github.com/DeLuci/coog-music/internal/forms"

type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	Form            *forms.Form
	CSRFToken       string
	IsAuthenticated int
}
