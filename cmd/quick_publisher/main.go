package main

import (
	"context"
	"fmt"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
	"math/rand"
	"strings"
	"sync"
	"time"

	"gitlab.badanamu.com.cn/calmisland/imq"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/storage"
)

const (
	publishTimes = 4

	publishInterval = time.Millisecond * 500
)

var currentMQ imq.IMessageQueue

func extensionIsImage(extension string) bool {
	if extension == "jpg" || extension == "jpeg" || extension == "gif" || extension == "png" {
		return true
	}
	return false
}

func fetchAllImageObjects() []string {
	s3 := storage.DefaultStorage(context.Background())
	objs, err := s3.ListAll()
	if err != nil {
		panic(err)
	}
	ret := make([]string, 0)
	for i := range objs {
		fileParts := strings.Split(objs[i], ".")
		if extensionIsImage(fileParts[len(fileParts)-1]) {
			ret = append(ret, objs[i])
		}
	}
	return ret
}
func initMQ() {
	mq, err := imq.CreateMessageQueue(imq.Config{
		Drive:                  "redis-list",
		RedisHost:              config.Get().MQ.RedisHost,
		RedisPort:              config.Get().MQ.RedisPort,
		RedisPassword:          config.Get().MQ.RedisPassword,
		RedisFailedPersistence: config.Get().MQ.RedisFailedPersistence,
		RedisHandlerThread:     config.Get().MQ.MaxWorker,
	})
	if err != nil {
		panic(err)
	}
	currentMQ = mq
}

func quickPublish(files []string) {
	wg := new(sync.WaitGroup)
	for i := 0; i < publishTimes; i++ {
		sliceOutOfOrder(files)
		newFiles := make([]string, len(files))
		copy(newFiles, files)
		wg.Add(1)
		go func(f []string) {
			for j := range f {
				fmt.Println("Publish file: " + f[j])
				currentMQ.Publish(context.Background(), entity.MQPrefix+"attachment", f[j])
				time.Sleep(publishInterval)
			}
			wg.Done()
		}(newFiles)
	}
	wg.Wait()
}
func sliceOutOfOrder(in []string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(in), func(i, j int) {
		in[i], in[j] = in[j], in[i]
	})
}

func main() {
	//load config
	config.LoadYAML("../../settings.yaml")

	//initial mq
	initMQ()

	//fetch all images in s3
	files := fetchAllImageObjects()
	//publish
	quickPublish(files)
}
