package main

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"go-auth-service/config"
	"go-auth-service/internal/handlers"
	"go-auth-service/models/arg"
	"go-auth-service/resources"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
)

var DB *gorm.DB
var once sync.Once

func init() {
	once.Do(func() {
		var err error
		appConfig, err := config.LoadConfig()
		if err != nil {
			panic(err)
		}
		DB, err = resources.InitMysqlResource(appConfig)
		if err != nil {
			panic(err)
		}
	})

}
func main() {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, HTTP/2 World!"})
	})

	v1UserRouters := r.Group("/v1/users")
	{
		userHandel := handlers.NewUserHandel(DB)

		v1UserRouters.POST("/create", func(c *gin.Context) {
			var user arg.SaveOrUpdateUserArg
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			id, err := userHandel.CreateUser(&user)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, id)
		})

		v1UserRouters.GET("/get/:id", func(c *gin.Context) {
			id := c.Param("id")
			user, err := userHandel.GetUserById(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, user)
		})
	}

	// 启动HTTP/2服务器
	startServer(r)
}

func startServer(router *gin.Engine) {
	// 这里需要替换为你自己的证书和私钥路径
	certFile := "./cert.pem"
	keyFile := "./key.pem"

	server := &http.Server{
		Addr:    ":8083", // HTTPS默认端口是443
		Handler: router,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12, // 设置最低TLS版本
		},
	}

	log.Println("Starting HTTPS server on port 443...")
	err := server.ListenAndServeTLS(certFile, keyFile)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
