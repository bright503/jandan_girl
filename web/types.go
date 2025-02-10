package web

import (
	"html/template"
	"jandan_girl/models"
)

type PostWarp struct {
	*models.Post
	HtmlContent template.HTML
}
