package util

import "github.com/spf13/viper"

type config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

//LoadConfig load configfile and returning config struct
func LoadConfig(path string) (cf config, err error) {
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
