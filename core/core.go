package core

import "gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/core/exiftool"

func Init() error {
	//Start exiftool
	err := exiftool.GetExifTool().Start()
	if err != nil {
		return err
	}
	return nil
}
