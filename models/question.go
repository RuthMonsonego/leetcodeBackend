package models

// Question represents a coding challenge with templates for different languages.
type Question struct {
    Code             int    `gorm:"primaryKey" json:"code"`
    Title            string `json:"title"`
    Description      string `json:"description"`
    TemplateForGo    string `json:"template_for_go"`
    TemplateForPython string `json:"template_for_python"`
}

// ExecuteRequest holds the code execution request details.
type ExecuteRequest struct {
    Code         string `json:"code"`
    Language     string `json:"language"`
    QuestionCode int    `json:"questionCode"`
}
