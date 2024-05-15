package entry

import (
	"go-auth-service/models"
	"gorm.io/gorm"
	"time"
)

type OAuthClient struct {
	models.BaseModel
	Name           string `gorm:"unique;size:100"`
	MaxHostCount   int
	Status         int `gorm:"default:1;size:3"`
	IsPublicClient bool
	InviteRebate   string `gorm:"default:0;size:6"` //邀请注册返利额度 为0则为禁用邀请返利
	PublicKey      string `gorm:"size:200"`
	PrivateKey     string `gorm:"size:200"`
	Secret         string `gorm:"size:200"`
	Domain         string
	UserID         string `gorm:"size:60"`
	Description    string `gorm:"size:200"`
	ExpiresAt      time.Time
}

func (c *OAuthClient) IsPublic() bool {
	return c.IsPublicClient
}

// GetID 实现 oauth2.ClientInfo 接口
func (c *OAuthClient) GetID() string {
	return c.ID
}

func (c *OAuthClient) GetSecret() string {
	return c.Secret
}

func (c *OAuthClient) GetDomain() string {
	return c.Domain
}

func (c *OAuthClient) GetUserID() string {
	return c.UserID
}

func (c *OAuthClient) GetExpiresAt() time.Time {
	return c.ExpiresAt
}

func (c *OAuthClient) Save(db *gorm.DB) error {
	return db.Create(c).Error
}

func (c *OAuthClient) Update(db *gorm.DB) error {
	return db.Save(c).Error
}

func FindClientByID(db *gorm.DB, id string) (*OAuthClient, error) {
	var client OAuthClient
	if err := db.First(&client, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &client, nil
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
