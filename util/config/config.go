package config

import (
	"screen_stream/util"

	"github.com/spf13/viper"
)

type Config struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func Load(path string) (config Config,err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()


	if err = viper.ReadInConfig(); err != nil{
		return Config{}, err
	} else if err = viper.Unmarshal(&config); err != nil{
		return 
	} 
	
	config.Password = util.Hash([]byte(config.Password))
	config.Username = util.Hash([]byte(config.Username))

	return
}
