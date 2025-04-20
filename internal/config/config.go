package config

import (
	pvzv1 "github.com/AnikinSimon/avito-test-backend/internal/grpc/pvz/v1"
	jwttoken "github.com/AnikinSimon/avito-test-backend/internal/pkg/jwt"
	"github.com/AnikinSimon/avito-test-backend/internal/pkg/web"
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`

	HTTPServerCfg web.ServerConfig `mapstructure:"httpserver"`
	GRPCServerCfg pvzv1.Config     `mapstructure:"grpcserver"`

	TokenService jwttoken.TokenServiceConfig `mapstructure:"auth"`
}

func LoadConfig(cfgPath string) (config AppConfig, err error) {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	viper.SetConfigFile(cfgPath)
	viper.MergeInConfig()

	viper.Unmarshal(&config)

	log.Println("config:", config)

	return
}
