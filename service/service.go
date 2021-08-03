package service

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/imq"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/log"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/runtime"
	"sync"
)

type FileProcessingService struct {
	mq                   imq.IMessageQueue
	handler              map[string]func(ctx context.Context, f *entity.HandleFileParams) error
	supportExtensionsMap map[string][]string
	mqChannels           []int
	quit                 chan struct{}
}

func (fp *FileProcessingService) Start() {
	//init MQ
	fp.initMQ()

	//init route
	fp.initProcessors()

	//subscribe topics
	fp.subscribeTopics()

	log.Info(context.Background(), "Service is starting...")
	//<-fp.quit
}

func (fp *FileProcessingService) Stop() {
	for i := range fp.mqChannels {
		fp.mq.Unsubscribe(fp.mqChannels[i])
	}
	//fp.quit <- struct{}{}
}

func (fp *FileProcessingService) MQ() imq.IMessageQueue {
	return fp.mq
}
func (fp *FileProcessingService) SupportExtensions() map[string][]string {
	return fp.supportExtensionsMap
}

func (fp *FileProcessingService) PendingMessages() (map[string][]string, error){
	res := make(map[string][]string)
	var err error
	for topic := range fp.supportExtensionsMap {
		res[topic], err = fp.mq.PendingMessage(context.Background(), topic)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (fp *FileProcessingService) initMQ() {
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
	fp.mq = mq
}

func (fp *FileProcessingService) subscribeTopics() {
	for topic, handler := range fp.handler {
		cid := fp.mq.SubscribeWithReconnect(topic, func(ctx context.Context, message string) error {
			//Update workers num
			runtime.GetWorkersInfo().Add()
			defer runtime.GetWorkersInfo().Done()
			log.Info(ctx,"receive topic: %v, message: %v \n", topic, message)
			return fp.handleMessage(ctx, topic, message, handler)
		})
		fp.mqChannels = append(fp.mqChannels, cid)
	}
}

func (fp *FileProcessingService) handleMessage(ctx context.Context,
	topic string,
	message string,
	handler func(ctx context.Context, f *entity.HandleFileParams) error) error {
	//parse file info

	fileInfo := entity.ParseFileInfo(topic, message)
	if fileInfo == nil {
		log.Warn(ctx, "ParseFileInfo failed, message: %v \n", message)
		return nil
	}
	//ignore unsupported extension
	supportExtension := fp.containsString(fileInfo.Extension, fp.supportExtensionsMap[topic])
	if !supportExtension {
		return nil
	}

	log.Debug(ctx, "downloading File, fileInfo: %#v \n", fileInfo)
	//download file
	fileParams, err := fp.downloadFile(ctx, fileInfo)
	if err != nil {
		log.Error(ctx, "downloadFile failed, err: %v \n", err)
		return err
	}
	log.Debug(ctx, "downloadFile success, fileParams: %#v \n", fileParams)
	defer fileParams.CleanLocalFile(ctx)
	defer fileParams.CleanOutputFile(ctx)

	err = fp.backupFile(ctx, fileInfo, fileParams.LocalFile)
	if err != nil {
		return err
	}

	//handle file
	err = handler(ctx, fileParams)
	if err != nil {
		log.Error(ctx, "Handle file failed, fileParams: %#v, err: %v", fileParams, err)
		log.Failed(ctx, "Handle file failed, fileParams: %#v, err: %v", fileParams, err)
		return nil
	}

	//upload file
	err = fp.uploadHandledFile(ctx, fileInfo, fileParams.DistPath)
	if err != nil {
		return err
	}

	log.Debug(ctx,"handle file success")
	return nil
}

var (
	_fileProcessingService *FileProcessingService
	_fileProcessingServiceOnce sync.Once
)

func GetFileProcessingService() *FileProcessingService {
	_fileProcessingServiceOnce.Do(func() {
		_fileProcessingService = &FileProcessingService{
			supportExtensionsMap: make(map[string][]string),
			handler:              make(map[string]func(ctx context.Context, f *entity.HandleFileParams) error),
			quit:                 make(chan struct{}),
		}
	})
	return _fileProcessingService
}
