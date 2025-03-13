package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config 结构体存储所有配置信息
type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`

	JWT struct {
		Secret string `yaml:"secret"`
		Expire int    `yaml:"expire"`
	} `yaml:"jwt"`
}

// GlobalConfig 作为全局变量存储配置信息
var GlobalConfig Config

// LoadConfig 读取并解析 YAML 配置文件
func LoadConfig(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}

	err = yaml.Unmarshal(data, &GlobalConfig)
	if err != nil {
		log.Fatalf("解析 YAML 配置失败: %v", err)
	}

	log.Println("配置加载成功")
}
