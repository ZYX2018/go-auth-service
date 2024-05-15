package handlers

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	redisStore "github.com/go-oauth2/redis/v4"
	"github.com/go-redis/redis/v8"
	"go-auth-service/config"
	"go-auth-service/internal/utils"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
)

func initOauth(config *config.AppConfig, db *gorm.DB, rdb *redis.Client) {

	// 配置Redis
	redisTokenStore := redisStore.NewRedisStoreWithCli(rdb, config.Server.NameSpace)

	// 配置OAuth2管理器
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.MustTokenStorage(redisTokenStore, nil)

	// 使用MySQL存储客户端信息
	clientStore := NewMySQLClientStore(db)
	manager.MapClientStorage(clientStore)

	// 使用SM2 JWT生成器
	signMethod := utils.NewSigningMethodSM2()
	jwtGen := utils.NewSM2JWTAccessGenerate(config, signMethod)
	manager.MapAccessGenerate(jwtGen)

	// 使用MySQL存储用户信息
	userStore := NewUserHandel(db)

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		user, err := userStore.GetUserByUsername(username, clientID)
		if err != nil || user.Password == password {
			return "", err
		}
		return user.ID, nil
	})

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// 加载服务器证书和密钥
	cert, err := tls.LoadX509KeyPair("path/to/your_server_cert.pem", "path/to/your_server_key.pem")
	if err != nil {
		log.Fatal(err)
	}

	// 创建证书池并添加客户端CA证书
	clientCACert, err := ioutil.ReadFile("path/to/your_client_ca_cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	// 启用HTTP/2和双向TLS
	srvTLS := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientCAs:    clientCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert, // 要求并验证客户端证书
			MinVersion:   tls.VersionTLS12,
			// 其他TLS配置
		},
	}

	log.Fatal(srvTLS.ListenAndServeTLS("path/to/your_server_cert.pem", "path/to/your_server_key.pem"))
}
