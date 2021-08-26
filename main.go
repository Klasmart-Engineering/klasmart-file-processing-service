package main

import (
	"fmt"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/api"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/core"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/service"
)

func main() {
	//load config
	config.MustLoad("./settings.yaml")

	//start Handler
	srv := service.GetFileProcessingService()
	srv.Start()

	//Init core
	err := core.Init()
	if err != nil {
		fmt.Println("Init core failed, err:", err)
		panic(err)
	}

	//start API
	route := api.GetServer()
	route.Start()
}
