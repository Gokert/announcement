package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
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

type GrpcConfig struct {
	Addr           string `yaml:"addr"`
	Port           string `yaml:"port"`
	ConnectionType string `yaml:"connection_type"`
}

func InitEnv() error {
	envMap := map[string]string{
		"REDIS_ADDR":      "127.0.0.1:6379",
		"AUTH_PSX_HOST":   "127.0.0.1",
		"MARKET_PSX_HOST": "127.0.0.1",
		"AUTH_GRPC_ADDR":  "127.0.0.1",
	}

	for key, defValue := range envMap {
		if err := setDefaultEnv(key, defValue); err != nil {
			return err
		}
	}

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("load .env error: %s", err.Error())
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

func GetAuthPsxConfig() (*DbPsxConfig, error) {
	v := viper.GetViper()
	v.AutomaticEnv()

	cfg := &DbPsxConfig{
		User:         v.GetString("AUTH_PSX_USER"),
		Password:     v.GetString("AUTH_PSX_PASSWORD"),
		Dbname:       v.GetString("AUTH_PSX_DBNAME"),
		Host:         v.GetString("AUTH_PSX_HOST"),
		Port:         v.GetInt("AUTH_PSX_PORT"),
		Sslmode:      v.GetString("AUTH_PSX_SSLMODE"),
		MaxOpenConns: v.GetInt("AUTH_PSX_MAXCONNS"),
		Timer:        v.GetInt("AUTH_PSX_TIMER"),
	}

	return cfg, nil
}

func GetMarketPsxConfig() (*DbPsxConfig, error) {
	v := viper.GetViper()
	v.AutomaticEnv()

	cfg := &DbPsxConfig{
		User:         v.GetString("MARKET_PSX_USER"),
		Password:     v.GetString("MARKET_PSX_PASSWORD"),
		Dbname:       v.GetString("MARKET_PSX_DBNAME"),
		Host:         v.GetString("MARKET_PSX_HOST"),
		Port:         v.GetInt("MARKET_PSX_PORT"),
		Sslmode:      v.GetString("MARKET_PSX_SSLMODE"),
		MaxOpenConns: v.GetInt("MARKET_PSX_MAXCONNS"),
		Timer:        v.GetInt("MARKET_PSX_TIMER"),
	}

	return cfg, nil
}

func GetGrpcConfig() (*GrpcConfig, error) {
	v := viper.GetViper()
	v.AutomaticEnv()

	cfg := &GrpcConfig{
		Addr:           v.GetString("AUTH_GRPC_ADDR"),
		Port:           v.GetString("AUTH_GRPC_PORT"),
		ConnectionType: v.GetString("AUTH_CONN_TYPE"),
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
