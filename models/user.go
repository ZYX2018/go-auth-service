package models

type UserGroup struct {
	BaseModel
	Name         string `gorm:"unique;size:100"`
	Description  string `gorm:"size:200"`
	MaxHostCount int
}

type Role struct {
	BaseModel
	Name        string `gorm:"uniqueIndex:idx_group_role;size:100"`
	GroupId     string `gorm:"uniqueIndex:idx_group_role;size:60"`
	GroupName   string `gorm:"size:100"`
	Description string `gorm:"size:200"`
}

type User struct {
	BaseModel
	Username    string `gorm:"uniqueIndex:idx_group_user;size:100"`
	Password    string `gorm:"size:500"`
	Email       string `gorm:"size:200"`
	PhoneNumber string `gorm:"size:20"`
	GroupId     string `gorm:"uniqueIndex:idx_group_user;size:60"`
	GroupName   string `gorm:"size:100"`
}
