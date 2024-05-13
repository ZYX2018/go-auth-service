package config

import (
	"fmt"
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

	ServerConfig struct {
		cert string `yaml:"cert"`
		key  string `yaml:"key"`
	}
}

func LoadConfig() (*AppConfig, error) {
	f, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("config.yaml 文件关闭失败")
		}
	}(f)

	var cfg AppConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
