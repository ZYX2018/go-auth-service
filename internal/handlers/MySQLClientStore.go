package handlers

import (
	"context"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"go-auth-service/models/entry"
	"gorm.io/gorm"
	"time"
)

type MySQLClientStore struct {
	db *gorm.DB
}

func NewMySQLClientStore(db *gorm.DB) *MySQLClientStore {
	return &MySQLClientStore{db: db}
}

func (s *MySQLClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	var client entry.OAuthClient
	err := s.db.First(&client, "id = ?", id).Error
	if err != nil {

		return nil, err
	}

	if client.ExpiresAt.Before(time.Now()) {
		return nil, errors.ErrInvalidClient // 客户端已过期
	}

	return &client, nil
}
