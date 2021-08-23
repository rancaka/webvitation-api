package models

type Template struct {
	TemplateID   string `db:"template_id"`
	TemplateName string `db:"template_name"`
	Timestamp
}
