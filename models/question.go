package models

type Question struct {
    Code             int    `gorm:"primaryKey" json:"code"`
    Title            string `json:"title"`
    Description      string `json:"description"`
    TemplateForGo    string `json:"template_for_go"`
    TemplateForPython string `json:"template_for_python"`
}
