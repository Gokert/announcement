package configs

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type DbPsxConfig struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Dbname       string `yaml:"dbname"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Sslmode      string `yaml:"sslmode"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	Timer        int    `yaml:"timer"`
}

type DbRedisCfg struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	DbNumber int    `yaml:"db"`
	Timer    int    `yaml:"timer"`
}

func InitEnv() error {
	envMap := map[string]string{
		"REDIS_ADDR":     "127.0.0.1:6379",
		"REDIS_PASSWORD": "",
		"REDIS_DB":       "0",
		"REDIS_TIMER":    "15",
	}

	for key, defValue := range envMap {
		if err := setDefaultEnv(key, defValue); err != nil {
			return err
		}
	}

	return nil
}

func setDefaultEnv(key, value string) error {
	if _, exists := os.LookupEnv(key); !exists {
		err := os.Setenv(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetPsxConfig(cfgPath string) (*DbPsxConfig, error) {
	v := viper.GetViper()
	v.SetConfigFile(cfgPath)
	v.SetConfigType(strings.TrimPrefix(filepath.Ext(cfgPath), "."))

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &DbPsxConfig{
		User:         v.GetString("user"),
		Password:     v.GetString("password"),
		Dbname:       v.GetString("dbname"),
		Host:         v.GetString("host"),
		Port:         v.GetInt("port"),
		Sslmode:      v.GetString("sslmode"),
		MaxOpenConns: v.GetInt("max_open_conns"),
		Timer:        v.GetInt("timer"),
	}

	return cfg, nil
}

func GetRedisConfig() (*DbRedisCfg, error) {
	v := viper.GetViper()
	v.AutomaticEnv()

	cfg := &DbRedisCfg{
		Host:     v.GetString("REDIS_ADDR"),
		Password: v.GetString("REDIS_PASSWORD"),
		DbNumber: v.GetInt("REDIS_DB"),
		Timer:    v.GetInt("REDIS_TIMER"),
	}

	return cfg, nil
}
