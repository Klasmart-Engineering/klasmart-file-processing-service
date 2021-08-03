package service

import (
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/processor"
)


func (fp *FileProcessingService) initProcessors(){
	fp.addRouteProcessor("attachment", processor.GetAttachmentProcessor())
}

func (fp *FileProcessingService) addRouteProcessor(classify string,
	processor processor.IFileProcessor){
	topic := entity.MQPrefix + classify
	fp.handler[topic] = processor.HandleFile
	fp.supportExtensionsMap[topic] = processor.SupportExtensions()
}

