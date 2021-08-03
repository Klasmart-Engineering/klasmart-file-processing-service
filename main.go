package main

import (
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/api"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/service"
)

func main(){
	//load config
	config.LoadYAML("./settings.yaml")

	//start Handler
	srv := service.GetFileProcessingService()
	srv.Start()

	//start API
	route := api.GetServer()
	route.Start()
}
