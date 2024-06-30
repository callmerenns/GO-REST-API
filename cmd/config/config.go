package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Driver   string
}

type ApiConfig struct {
	ApiPort string
}

type TokenConfig struct {
	IssuerName       string `json:"IssuerName"`
	JwtSignatureKey  []byte `json:"JwtSignatureKey"`
	JwtSigningMethod *jwt.SigningMethodHMAC
	JwtExpiresTime   time.Duration
}

type Config struct {
	DbConfig
	ApiConfig
	TokenConfig
}

func (c *Config) readConfig() error {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}
	log.Println("Current directory:", cwd)

	log.Println("Looking for .env file in current directory")

	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("failed to load configuration: missing env file %v", err.Error())
	}

	log.Println(".env file loaded successfully")

	c.DbConfig = DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Driver:   os.Getenv("DB_DRIVER"),
	}

	c.ApiConfig = ApiConfig{ApiPort: os.Getenv("API_PORT")}

	tokenExpire, _ := strconv.Atoi(os.Getenv("TOKEN_EXPIRE"))
	c.TokenConfig = TokenConfig{
		IssuerName:       os.Getenv("TOKEN_ISSUE"),
		JwtSignatureKey:  []byte(os.Getenv("TOKEN_SECRET")),
		JwtSigningMethod: jwt.SigningMethodHS256,
		JwtExpiresTime:   time.Duration(tokenExpire) * time.Minute,
	}

	if c.Host == "" || c.Port == "" || c.User == "" || c.Name == "" || c.Driver == "" || c.ApiPort == "" ||
		c.IssuerName == "" || c.JwtExpiresTime < 0 || len(c.JwtSignatureKey) == 0 {
		return fmt.Errorf("missing required environment")
	}

	return nil

}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.readConfig(); err != nil {
		return nil, err
	}
	return cfg, nil
}
