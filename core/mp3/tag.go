package mp3

import (
	"context"
	"github.com/frolovo22/tag"
	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	"os"
)

func RemoveMetadata(ctx context.Context, file, outFile string) error {
	in, err := os.Open(file)
	if err != nil {
		log.Error(ctx, "Can't create output file",
			log.Err(err),
			log.Any("file", file))
		return err
	}
	defer in.Close()

	tags, err := tag.Read(in)
	if err != nil {
		log.Error(ctx, "Can't read output file",
			log.Err(err),
			log.Any("file", file))
		return err
	}

	err = tags.DeleteAll()
	if err != nil {
		log.Error(ctx, "Failed to delete tags",
			log.Err(err),
			log.Any("file", file))
		return err
	}
	err = tags.SaveFile(outFile)
	if err != nil {
		log.Error(ctx, "Failed to save processed tags",
			log.Err(err),
			log.Any("file", outFile))
		return err
	}
	return nil
}
