package models

type Question struct {
    ID               int        `gorm:"primaryKey" json:"id"`
    Title            string     `json:"title" binding:"required"`
    Description      string     `json:"description" binding:"required"`
    TemplateForGo    string     `json:"template_for_go" binding:"required"`
    TemplateForPython string    `json:"template_for_python" binding:"required"`
    Parameters       Parameters `gorm:"type:json" json:"parameters" binding:"required"`
}
