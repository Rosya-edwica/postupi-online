package config

import "github.com/spf13/viper"

type Config struct {
	DBUrl string `mapstructure:"DB_URL"`
}

func LoadConfig(type_parsing string) (c Config, err error) {
	configPath := getConfigPath(type_parsing)

	viper.AddConfigPath(configPath)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)
	return
}

func getConfigPath(typeParsing string) (path string) {
	if typeParsing == "college" {
		return "./pkg/common/config/envs/college"
	}
	return "./pkg/common/config/envs/vuz"
}
