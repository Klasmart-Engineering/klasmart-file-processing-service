package core

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/entity"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/test"
	"testing"
)

var (
	mp3TestFileIn       = "../test-data/piiFile.docx"
	mp3TestFileOut      = "/tmp/piiFile.docx-handled"
	mp3ExpectedChecksum = "77f0b99c9f6602f28fb807c2a401c3527683dd274c61996b905462b3f570b85e"
)

func TestRemoveMetadataOOXML(t *testing.T) {

	fileInfo := entity.ParseFileInfo(mp3TestFileIn)
	fileParams := &entity.HandleFileParams{
		Extension: fileInfo.Extension,
		Name:      fileInfo.Name,
		LocalPath: mp3TestFileIn,
	}

	err := GetRemoveOOXMLMetaDataHandler().Do(context.Background(), fileParams)
	if err != nil {
		t.Errorf("Failed to process file %s", mp3TestFileIn)
	}

	test.Checksum(t, mp3TestFileOut, mp3ExpectedChecksum)

	t.Cleanup(func() {
		test.RemoveFile(t, mp3TestFileOut)
	})
}
