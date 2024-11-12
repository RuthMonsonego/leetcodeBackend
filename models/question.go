package models

type Question struct {
    Code             int    `gorm:"primaryKey" json:"code"`
    Title            string `json:"title"`
    Description      string `json:"description"`
    TemplateForGo    string `json:"template_for_go"`
    TemplateForPython string `json:"template_for_python"`
}

type ExecuteRequest struct {
    Code         string `json:"code"`
    Language     string `json:"language"`
    QuestionCode int    `json:"questionCode"`
}
