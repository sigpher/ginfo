package config

import "github.com/spf13/viper"

func InitConfig() {
	// workDir, _ := os.Getwd()
	viper.SetConfigName("settings")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
