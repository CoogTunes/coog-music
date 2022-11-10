package models

import "github.com/CoogTunes/coog-music/internal/forms"

type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	Form            *forms.Form
	UserData        Users
	CSRFToken       string
	IsAuthenticated int
}
