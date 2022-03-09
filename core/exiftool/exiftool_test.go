package exiftool

import (
	"context"
	"testing"
)

var (
	jpgTestFileIn  = "../test/piiFile.jpg"
	jpgTestFileOut = "/tmp/piiFile.jpg"
)

func TestStartExifToolJPG(t *testing.T) {
	exifTool := GetExifTool()
	err := exifTool.Start()
	if err != nil {
		panic(err)
	}
	err = exifTool.RemoveMetadata(context.Background(),
		jpgTestFileIn,
		jpgTestFileOut,
		JpegTags,
	)
	if err != nil {
		t.Errorf("Failed to process file %s", jpgTestFileIn)
	}
}

func TestStartExifTool(t *testing.T) {
	exifTool := GetExifTool()
	err := exifTool.Start()
	if err != nil {
		panic(err)
	}
	err = exifTool.RemoveMetadata(context.Background(),
		jpgTestFileIn,
		jpgTestFileOut,
		JpegTags,
	)
	if err != nil {
		t.Errorf("Failed to process file %s", jpgTestFileIn)
	}
}
