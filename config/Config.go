package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type AppConfig struct {
	Mysql struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		UserName string `yaml:"userName"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"mysql"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`

	Server struct {
		Port         int    `yaml:"port"`
		cert         string `yaml:"cert"`
		key          string `yaml:"key"`
		ClientCACert string `yaml:"client_ca_cert"`
		NameSpace    string `yaml:"name_space"`
	} `yaml:"server"`

	SM2 struct {
		PrivateKey string `yaml:"private_key"`
		PublicKey  string `yaml:"public_key"`
	} `yaml:"sm2"`
}

func LoadConfig() *AppConfig {
	config := &AppConfig{}
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic("读取配置文件失败")
	}
	if err := yaml.Unmarshal(data, config); err != nil {
		panic("解析配置文件失败")
	}
	return config
}
