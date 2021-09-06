package eyed3

import (
	"context"
	"testing"
)

func TestRemoveMetadata(t *testing.T) {
	err := GetEyeD3Tool().RemoveMP3MetaData(context.Background(), "D:\\1.mp3")
	if err != nil {
		panic(err)
	}
	t.Log("Done")
}
