package entry

import "go-auth-service/models"

type OAuthClient struct {
	models.BaseModel
	Name            string `gorm:"unique;size:100"`
	MaxHostCount    int
	Status          int    `gorm:"default:1;size:3"`
	InviteRebate    string `gorm:"default:0;size:6"` //邀请注册返利额度 为0则为禁用邀请返利
	PublicKey       string `gorm:"size:200"`
	PrivateKey      string `gorm:"size:200"`
	Password        string `gorm:"size:200"`
	CallBackAddress string
	Description     string `gorm:"size:200"`
}

type OAuthAccessRecord struct {
	models.BaseModel
	ClientID     string `gorm:"size:60"`
	AccessOrigin string `gorm:"size:500"`
	Description  string `gorm:"size:200"`
}

type OauthScopeGroup struct {
	models.BaseModel
	Name        string `gorm:"uniqueIndex:idx_client_scope;size:100"`
	Description string `gorm:"size:200"`
	ClientID    string `gorm:"uniqueIndex:idx_client_scope;size:60"`
}

type OauthScopes struct {
	models.BaseModel
	Scope       string `gorm:"uniqueIndex:idx_client_scope;size:100"`
	Description string `gorm:"size:200"`
	GroupID     string `gorm:"uniqueIndex:idx_client_scope;size:60"`
}
