package service

import (
	"errors"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/processor"
	"strings"
)

func (fp *FileProcessingService) initProcessors() error {
	processors := config.Get().Core.Processors
	if processors == "" {
		return errors.New("processors environment variable must be set")
	}

	for _, p := range strings.Split(processors, ",") {
		fp.addProcessor(p, processor.GetProcessor(p))
	}
	return nil
}

func (fp *FileProcessingService) addProcessor(key string,
	processor processor.IFileProcessor) {
	fp.handler[key] = processor.HandleFile
	fp.supportExtensionsMap[key] = processor.SupportExtensions()
}
