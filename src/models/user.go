package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"Name"`
	Phone    string `json:"Phone"`
	Role     string `json:"Role"`
	Password string
}
