package exiftool

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/test"
	"testing"
)

var (
	jpgTestFileIn        = "../../test-data/piiFile.jpg"
	jpgTestFileOut       = "/tmp/piiFile.jpg"
	jpgExpectedChecksum  = "df6bc1e2409f5046e497d84c85df8a1346519eb3201d63041457c250983831af"
	jpegTestFileIn       = "../../test-data/piiFile.jpeg"
	jpegTestFileOut      = "/tmp/piiFile.jpeg"
	jpegExpectedChecksum = "7f6020ef2577ed8bcbc30c6801cf983124f93f42b0b2a53c2e36e0a6666f2cb1"
)

func TestRemoveMetadataJPG(t *testing.T) {
	exifTool := GetExifTool()
	err := exifTool.Start()
	if err != nil {
		t.Fatal(err)
	}
	err = exifTool.RemoveMetadata(context.Background(),
		jpgTestFileIn,
		jpgTestFileOut,
		JpegTags,
	)
	if err != nil {
		t.Errorf("Failed to process file %s", jpgTestFileIn)
	}

	test.Checksum(t, jpgTestFileOut, jpgExpectedChecksum)

	t.Cleanup(func() {
		test.RemoveFile(t, jpgTestFileOut)
	})
}

func TestRemoveMetadataJPEG(t *testing.T) {
	exifTool := GetExifTool()
	err := exifTool.Start()
	if err != nil {
		t.Fatal(err)
	}
	err = exifTool.RemoveMetadata(context.Background(),
		jpegTestFileIn,
		jpegTestFileOut,
		JpegTags,
	)
	if err != nil {
		t.Errorf("Failed to process file %s", jpgTestFileIn)
	}

	test.Checksum(t, jpegTestFileOut, jpegExpectedChecksum)

	t.Cleanup(func() {
		test.RemoveFile(t, jpegTestFileOut)
	})
}
