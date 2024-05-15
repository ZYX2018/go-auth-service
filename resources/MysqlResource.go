package resources

import (
	"fmt"
	"go-auth-service/config"
	"go-auth-service/models/entry"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func InitMysqlResource(appConfig *config.AppConfig) *gorm.DB {
	fmt.Println("config", appConfig)
	//encodedPassword := url.QueryEscape(appConfig.Mysql.Password)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		appConfig.Mysql.UserName, appConfig.Mysql.Password,
		appConfig.Mysql.Host, appConfig.Mysql.Port,
		appConfig.Mysql.DBName)
	fmt.Println("sql", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("mysql 链接失败")
	}
	AutoMigrateAllModels(db)
	return db
}
func AutoMigrateAllModels(db *gorm.DB) {
	modelsToMigrate := []interface{}{
		&entry.User{},
		&entry.Role{},
		&entry.RelationUserAndRole{},
		&entry.RelationRoleAndScope{},
		&entry.OauthScopes{},
		&entry.OauthScopeGroup{},
		&entry.OAuthClient{},
		&entry.OAuthAccessRecord{},

		// 添加其他模型...
	}

	for _, model := range modelsToMigrate {
		err := db.AutoMigrate(model)
		if err != nil {
			log.Fatal("Failed to auto migrate model:", err)
		}
	}
}
