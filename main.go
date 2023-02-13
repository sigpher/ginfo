package main

import (
	"fmt"
	"ginfo/config"
	"ginfo/router"

	"github.com/spf13/viper"
)

func main() {
	router := router.SetupRouter()
	config.InitConfig()
	addr := viper.GetString("http.addr")
	err := router.Run(addr)
	if err != nil {
		fmt.Println("Running server faild, err:", err)
	}
}
