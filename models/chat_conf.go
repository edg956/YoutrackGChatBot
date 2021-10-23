package models

import (
	"gorm.io/gorm"
)

type SpaceConfiguration struct {
	gorm.Model
	token     string
	space     string
	projectId string
	boardId   string
	sprintId  string
}
