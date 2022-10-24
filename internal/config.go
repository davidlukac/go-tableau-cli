package internal

import "github.com/spf13/viper"

type Config struct {
	TableauURL string `mapstructure:"TABLEAU_URL"`
	Username   string `mapstructure:"TABLEAU_USERNAME"`
	Password   string `mapstructure:"TABLEAU_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".local")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.MergeInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
