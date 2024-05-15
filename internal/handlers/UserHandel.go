package handlers

import (
	"github.com/shopspring/decimal"
	"go-auth-service/models"
	"go-auth-service/models/arg"
	"go-auth-service/models/entry"
	"gorm.io/gorm"
)

type UserHandel interface {
	CreateUser(arg *arg.SaveOrUpdateUserArg) (*models.BaseModel, error)
	GetUserByUsername(username string, clientId string) (*entry.User, error)
	GetUserById(id string) (*entry.User, error)
	DeleteUserById(id string) (*models.BaseModel, error)
}

type userHandel struct {
	db *gorm.DB
}

func NewUserHandel(db *gorm.DB) UserHandel {
	return &userHandel{db: db}
}

func (handel *userHandel) CreateUser(arg *arg.SaveOrUpdateUserArg) (*models.BaseModel, error) {
	user := &entry.User{
		Username:   arg.Username,
		Password:   arg.Password,
		ClientId:   arg.ClientId,
		UserOrigin: arg.UserOrigin,
	}
	err := handel.db.Transaction(func(tx *gorm.DB) error {
		//解密密码
		var client = entry.OAuthClient{}
		result := handel.db.First(&client, "id = ?", arg.ClientId)
		if result.Error != nil {
			panic("客户端不存在")
		}
		oldUser := &entry.User{}
		result = handel.db.First(&oldUser, "username = ?", arg.Username)
		if result.Error == nil {
			panic("用户名已存在")
		}

		if arg.RegisterMethod == 0 {
			user.Email = arg.Username
		}
		if arg.RegisterMethod == 1 {
			user.PhoneNumber = arg.Username
		}
		//生成hash的邀请码

		//判断是否三邀请注册 ， 若是注册则 需给被邀请人返利
		if arg.UserOrigin == 1 {
			//判断邀请码是否正确
			if arg.InviteCode == "" {
				panic("邀请码不能为空")
			}
			//获取邀请人信息
			var inviteUser = entry.User{}
			result = handel.db.First(&inviteUser, "invite_code = ?", arg.InviteCode)
			if result.Error != nil {
				panic("邀请码不存在")
			}
			//给邀请人返利
			if arg.ClientId == "" {
				panic("客户端不能为空")
			}

			oldRebate, _ := decimal.NewFromString(inviteUser.Balance)
			addRebate, _ := decimal.NewFromString(client.InviteRebate)
			newRebate := oldRebate.Add(addRebate)
			result = handel.db.Model(&inviteUser).Update("balance", newRebate.String())
			if result.Error != nil {
				panic("邀请人返利失败")
			}
			user.CreatedBy = inviteUser.Username
			user.UpdateBy = inviteUser.Username
		}

		result = handel.db.Create(&user)
		if result.Error != nil {
			panic("创建用户失败")
		}
		//提交事务
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &models.BaseModel{ID: user.ID}, nil
}

func (handel *userHandel) GetUserByUsername(userName string, clientId string) (*entry.User, error) {
	if userName == "" {
		panic("用户名不能为空")
	}
	var userModel entry.User
	result := handel.db.First(&userModel, "userName = ? and clientId", userName, clientId)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userModel, nil
}

func (handel *userHandel) GetUserById(id string) (*entry.User, error) {
	if id == "" {
		panic("id不能为空")
	}
	var userModel entry.User
	result := handel.db.First(&userModel, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userModel, nil
}
func (handel *userHandel) DeleteUserById(id string) (*models.BaseModel, error) {
	if id == "" {
		panic("id不能为空")
	}
	result := handel.db.Delete(&entry.User{}, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &models.BaseModel{
		ID: id,
	}, nil
}
