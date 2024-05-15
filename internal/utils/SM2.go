package utils

import (
	"crypto/rand"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"go-auth-service/config"
	"os"
)

type SigningMethodSM2 struct {
}

func NewSigningMethodSM2() jwt.SigningMethod {
	return &SigningMethodSM2{}
}

func (m *SigningMethodSM2) Alg() string {
	return "SM2"
}

// Verify Returns nil if signature is valid
func (m *SigningMethodSM2) Verify(signingString string, sig []byte, key interface{}) error {
	sm2Key := key.(*sm2.PublicKey)
	isOk := sm2Key.Verify([]byte(signingString), sig)
	if isOk {
		return nil
	}
	return jwt.ErrSignatureInvalid
}

func (m *SigningMethodSM2) Sign(signingString string, key interface{}) ([]byte, error) {
	// Returns signature or error
	sm2Key := key.(*sm2.PrivateKey)
	return sm2Key.Sign(rand.Reader, []byte(signingString), nil)
}

// LoadSM2PrivateKey 加载私钥
func LoadSM2PrivateKey(config *config.AppConfig) *sm2.PrivateKey {
	keyPath := config.SM2.PrivateKey
	privateKeyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		panic("读取私钥文件失败")
	}
	privateKey, err := x509.ReadPrivateKeyFromPem(privateKeyPEM, nil)
	if err != nil {
		panic("解析私钥文件失败")
	}
	return privateKey
}

// LoadSM2PublicKey 加载公钥
func LoadSM2PublicKey(config *config.AppConfig) *sm2.PublicKey {
	keyPath := config.SM2.PublicKey
	publicKeyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		panic("读取公钥文件失败")
	}
	publicKey, err := x509.ReadPublicKeyFromPem(publicKeyPEM)
	if err != nil {
		panic("解析公钥文件失败")
	}
	return publicKey
}

// SM2SignString 签名
func SM2SignString(privateKey *sm2.PrivateKey, str string) string {
	signData, err := privateKey.Sign(rand.Reader, []byte(str), nil)
	if err != nil {
		panic("签名失败")
	}
	return string(signData)
}

// SM2DecryptString 解密
func SM2DecryptString(privateKey *sm2.PrivateKey, str string) string {
	decryptData, err := privateKey.DecryptAsn1([]byte(str))
	if err != nil {
		panic("解密失败")
	}
	return string(decryptData)
}

// SM2EncryptString 加密
func SM2EncryptString(publicKey *sm2.PublicKey, str string) string {
	encryptData, err := publicKey.EncryptAsn1([]byte(str), rand.Reader)
	if err != nil {
		panic("加密失败")
	}
	return string(encryptData)
}

// SM2VerifyString 验证签名
func SM2VerifyString(publicKey *sm2.PublicKey, dataString string, signString string) bool {
	return publicKey.Verify([]byte(dataString), []byte(signString))
}
