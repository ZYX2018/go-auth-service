package entry

import "go-auth-service/models"

type Role struct {
	models.BaseModel
	Name            string `gorm:"uniqueIndex:idx_client_role;size:50"`
	ClientId        string `gorm:"uniqueIndex:idx_client_role;size:60"`
	PermissionLevel int    `gorm:"default:0;size:3"` //权限等级 0 为普通注册默认角色 多角色时登陆默认使用权限等级最高角色
	ClientName      string `gorm:"size:100"`
	Description     string `gorm:"size:200"`
}

type User struct {
	models.BaseModel
	Username    string `gorm:"uniqueIndex:idx_client_user;size:100"`
	Password    string `gorm:"size:500"`
	Email       string `gorm:"size:200"`
	PhoneNumber string `gorm:"size:20"`
	Balance     string `gorm:"default:0"`
	InviteCode  string `gorm:"unique;size:100"`
	UserOrigin  int    `gorm:"default:0;size:3"`
	ClientId    string `gorm:"uniqueIndex:idx_client_user;size:60"`
}

type RelationRoleAndScope struct {
	models.BaseModel
	RoleId  string `gorm:"uniqueIndex:idx_client_role_scope;size:60"`
	ScopeId string `gorm:"uniqueIndex:idx_client_role_scope;size:60"`
}

type RelationUserAndRole struct {
	models.BaseModel
	UserId string `gorm:"uniqueIndex:idx_client_user_role;size:60"`
	RoleId string `gorm:"uniqueIndex:idx_client_user_role;size:60"`
}
