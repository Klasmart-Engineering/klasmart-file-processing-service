package mp3

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/test"
	"testing"
)

var (
	mp3TestFileIn       = "../../test-data/piiFile.mp3"
	mp3TestFileOut      = "/tmp/piiFile.mp3-handled"
	mp3ExpectedChecksum = "bbd096ad51734bd9173ca05d537b3f263a3e71ee9173cac336fb80c88f734944"
)

func TestRemoveMetadataMP3(t *testing.T) {

	err := RemoveMetadata(context.Background(), mp3TestFileIn, mp3TestFileOut)
	if err != nil {
		t.Errorf("Failed to process file %s", mp3TestFileIn)
	}

	test.Checksum(t, mp3TestFileOut, mp3ExpectedChecksum)

	t.Cleanup(func() {
		test.RemoveFile(t, mp3TestFileOut)
	})
}
