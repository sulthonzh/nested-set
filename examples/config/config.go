package config

import (
	"os"

	"github.com/spf13/viper"
	"gitlab.com/sulthonzh/scraperpath-nested-set/examples/utils"
)

// Configuration ::
type Configuration struct {
	DB struct {
		Adapter  string `json:"adapter"`
		Name     string `json:"name"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Password string `json:"password"`
		User     string `json:"user"`
	} `json:"db"`
}

// Config ::
var Config = Configuration{}

func init() {
	switch true {
	case os.Getenv("ENV") == "PRODUCTION":
		viper.SetConfigName("config") // Production Config
	case os.Getenv("ENV") == "LOCAL.DEV":
		viper.SetConfigName("local_dev_config") // Local Development Config
	default:
		viper.SetConfigName("dev_config") // Dev Config
	}

	viper.AddConfigPath("/etc/nestedset/api")
	viper.AddConfigPath("$HOME/.nestedset/api")
	viper.AddConfigPath(".")

	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		utils.ExitOnFailure(err)
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		utils.ExitOnFailure(err)
	}
}
