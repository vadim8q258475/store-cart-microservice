package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env"
)

type Config struct {
	DBName     string `env:"DB_NAME,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBUser     string `env:"DB_USER,required"`
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`

	UserHost string `env:"USER_HOST,required"`
	UserPort string `env:"USER_PORT,required"`

	ProductHost string `env:"PRODUCT_HOST,required"`
	ProductPort string `env:"PRODUCT_PORT,required"`

	Port string `env:"PORT,required"`
}

func (c Config) String() string {
	var sb strings.Builder

	sb.WriteString("Product Service Settings:\n")
	sb.WriteString("DB:\n")
	sb.WriteString(fmt.Sprintf("DB_NAME: %s\n", c.DBName))
	sb.WriteString(fmt.Sprintf("DB_PASSWORD: %s\n", c.DBPassword))
	sb.WriteString(fmt.Sprintf("DB_NAME: %s\n", c.DBName))
	sb.WriteString(fmt.Sprintf("DB_HOST %s\n", c.DBHost))
	sb.WriteString(fmt.Sprintf("DB_PORT: %s\n", c.DBPort))
	sb.WriteString("USER:\n")
	sb.WriteString(fmt.Sprintf("USER_HOST %s\n", c.UserHost))
	sb.WriteString(fmt.Sprintf("USER_PORT: %s\n", c.UserPort))
	sb.WriteString("PRODUCT:\n")
	sb.WriteString(fmt.Sprintf("PRODUCT_HOST %s\n", c.ProductHost))
	sb.WriteString(fmt.Sprintf("PRODUCT_PORT: %s\n", c.ProductPort))
	sb.WriteString("CART:\n")
	sb.WriteString(fmt.Sprintf("PORT: %s\n", c.Port))

	return sb.String()
}

func MustLoadConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		panic(fmt.Errorf("parsing config error: %s", err.Error()))
	}
	return cfg
}
