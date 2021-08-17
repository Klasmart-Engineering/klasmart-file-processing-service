package exiftool

import (
	"context"
	"testing"
)

func TestStartExifTool(t *testing.T) {
	exifTool := GetExifTool()
	err := exifTool.Start()
	if err != nil {
		panic(err)
	}
	err = exifTool.RemoveMetadata(context.Background(),
		"D:\\Work\\Temp\\cms_test.jpg",
		"D:\\Work\\Temp\\cms_test_out.jpg",
		jpegTags,
	)
	if err != nil {
		panic(err)
	}
	err = exifTool.RemoveMetadata(context.Background(),
		"D:\\Work\\Temp\\cms_test2.jpg",
		"D:\\Work\\Temp\\cms_test_out2.jpg",
		jpegTags,
	)
	if err != nil {
		panic(err)
	}
	t.Log("Done")
}
