package mp3

import (
	"context"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/test"
	"testing"
)

var (
	mp3TestFileIn       = "../../test-data/piiFile.mp3"
	mp3TestFileOut      = "/tmp/piiFile.mp3-handled"
	mp3ExpectedChecksum = "946fe9db5b82d1ea136d65c858f8a45ab30ed040274118d84891d2d7604b94aa"
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
