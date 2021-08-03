package main

import (
	"context"
	"fmt"
	"gitlab.badanamu.com.cn/calmisland/imq"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/storage"
	"strings"
	"sync"
	"time"
)

const(
	publishTimes = 10

	publishInterval = time.Millisecond * 500
)
var currentMQ imq.IMessageQueue

func extensionIsImage(extension string) bool{
	if extension == "jpg" || extension == "jpeg" || extension == "gif" || extension == "png" {
		return true
	}
	return false
}

func fetchAllImageObjects() []string{
	s3 := storage.DefaultStorage(context.Background())
	objs, err := s3.ListAll()
	if err != nil {
		panic(err)
	}
	ret := make([]string, 0)
	for i := range objs {
		fileParts := strings.Split(objs[i], ".")
		if extensionIsImage(fileParts[len(fileParts) - 1]) {
			ret = append(ret, objs[i])
		}
	}
	return ret
}
func initMQ(){
	mq, err := imq.CreateMessageQueue(imq.Config{
		Drive:                  "redis-list",
		RedisHost:              config.Get().MQ.RedisHost,
		RedisPort:              config.Get().MQ.RedisPort,
		RedisPassword:          config.Get().MQ.RedisPassword,
		RedisFailedPersistence: config.Get().MQ.RedisFailedPersistence,
		RedisHandlerThread: 		config.Get().MQ.MaxWorker,
	})
	if err != nil {
		panic(err)
	}
	currentMQ = mq
}

func quickPublish(files []string){
	wg := new(sync.WaitGroup)
	for i := 0; i < publishTimes; i ++ {
		wg.Add(1)
		go func() {
			for j := range files {
				fmt.Println("Publish file: " + files[j])
				currentMQ.Publish(context.Background(), entity.MQPrefix + "attachment", files[j])
				time.Sleep(publishInterval)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func main(){
	//load config
	config.LoadYAML("../settings.yaml")

	//initial mq
	initMQ()

	//fetch all images in s3
	files := fetchAllImageObjects()

	//publish
	quickPublish(files)
}
