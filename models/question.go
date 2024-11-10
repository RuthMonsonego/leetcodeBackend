package models

import "gorm.io/gorm"

type Question struct {
    gorm.Model
    Code          int      `json:"code"`
    Title         string   `json:"title"`
    Description   string   `json:"description"`
    TemplateForGo string   `json:"templateforgo"`
    TemplateForPy string   `json:"templateforpy"`
}
