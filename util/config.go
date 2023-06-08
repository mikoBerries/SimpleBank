package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Enviroment          string        `mapstructure:"ENVIROMENT"`
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	DBMigratePath       string        `mapstructure:"DB_MIGRATE_PATH"`
	HttpServerAddress   string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress   string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefeshTokenDuration time.Duration `mapstructure:"REFESH_TOKEN_DURATION"`
}

// LoadConfig load configfile and returning config struct
func LoadConfig(path string) (cf Config, err error) {
	//set param to viper class
	viper.AddConfigPath(path)  //set Path to config file
	viper.SetConfigName("app") //set "config file" filename
	viper.SetConfigType("env") //set "config file" extension | support :JSON, TOML, YAML, HCL, envfile and Java properties config files

	viper.AutomaticEnv()

	//read config file
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	//unmarshal config file to stuct
	err = viper.Unmarshal(&cf)
	if err != nil {
		return
	}
	return
}
