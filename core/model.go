package core

import "time"

//Model is the base struct for all models containing the common fields for them (ID, etc.).
//It differs from gorm.Model by having json tags with lowercase names.
type Model struct {
	ID        int        `gorm:"primary_key",json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index",json:"deletedAt"`
}
