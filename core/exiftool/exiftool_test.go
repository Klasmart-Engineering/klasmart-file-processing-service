package exiftool

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/test"
	"testing"
)

var (
	jpgTestFileIn        = "../../test-data/piiFile.jpg"
	jpgTestFileOut       = "/tmp/piiFile.jpg-handled"
	jpgExpectedChecksum  = "df6bc1e2409f5046e497d84c85df8a1346519eb3201d63041457c250983831af"
	jpegTestFileIn       = "../../test-data/piiFile.jpeg"
	jpegTestFileOut      = "/tmp/piiFile.jpeg-handled"
	jpegExpectedChecksum = "7f6020ef2577ed8bcbc30c6801cf983124f93f42b0b2a53c2e36e0a6666f2cb1"
	movTestFileIn        = "../../test-data/piiFile.mov"
	movTestFileOut       = "/tmp/piiFile.mov-handled"
	movExpectedChecksum  = "1a540fd3d7519e73ad0b22510f7b5c97e7f747d04c00e469b8978fce2466b516"
	mp4TestFileIn        = "../../test-data/piiFile.mp4"
	mp4TestFileOut       = "/tmp/piiFile.mp4-handled"
	mp4ExpectedChecksum  = "53a5d36e734ac8e2825a02d877bc2c8ac323c98a585a1324cee2cd8149474027"
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
		t.Errorf("Failed to process file %s", jpegTestFileIn)
	}

	test.Checksum(t, jpegTestFileOut, jpegExpectedChecksum)

	t.Cleanup(func() {
		test.RemoveFile(t, jpegTestFileOut)
	})
}

func TestRemoveMetadataMOV(t *testing.T) {
	exifTool := GetExifTool()
	err := exifTool.Start()
	if err != nil {
		t.Fatal(err)
	}
	err = exifTool.RemoveMetadata(context.Background(),
		movTestFileIn,
		movTestFileOut,
		MovTags,
	)
	if err != nil {
		t.Errorf("Failed to process file %s", movTestFileIn)
	}

	test.Checksum(t, movTestFileOut, movExpectedChecksum)

	t.Cleanup(func() {
		test.RemoveFile(t, movTestFileOut)
	})
}

func TestRemoveMetadataMP4(t *testing.T) {
	exifTool := GetExifTool()
	err := exifTool.Start()
	if err != nil {
		t.Fatal(err)
	}
	err = exifTool.RemoveMetadata(context.Background(),
		mp4TestFileIn,
		mp4TestFileOut,
		Mp4Tags,
	)
	if err != nil {
		t.Errorf("Failed to process file %s", mp4TestFileIn)
	}

	test.Checksum(t, mp4TestFileOut, mp4ExpectedChecksum)

	t.Cleanup(func() {
		test.RemoveFile(t, mp4TestFileOut)
	})
}
