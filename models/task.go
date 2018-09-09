package models

import "github.com/jinzhu/gorm"

type Task struct {
	gorm.Model
	Title       string `gorm:"column:title"`
	Description string `gorm:"column:description"`
	Spec        string `gorm:"column:spec"`
	Timeout     int    `gorm:"column:timeout"`
	Concurrent  bool   `gorm:"column:concurrent"`
	Status      int    `gorm:"column:status"`
	LogId       int    `gorm:"column:log_id"`
	Command     string `gorm:"column:command"`
}

func (Task) TableName() string {
	return "tasks"
}
