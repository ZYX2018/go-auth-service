package models

import (
	"github.com/go-oauth2/oauth2/utils/uuid"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        string `gorm:"size:60;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// 钩子：在文章被创建之前，通过UUID生成文章的ID
func (u *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	uid, _ := uuid.NewRandom()
	u.ID = uid.String()
	return
}
