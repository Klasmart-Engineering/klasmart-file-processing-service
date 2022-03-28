package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gitlab.badanamu.com.cn/calmisland/common-log/log"
	_ "gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/core"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/service"
)

type SqsBody struct {
	Records []struct {
		EventName string
		S3        struct {
			Bucket struct {
				Name string
			}
			Object struct {
				Key  string
				Size int
			}
		}
	}
}

func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) (string, error) {

	log.Debug(ctx, "SQS Record received",
		log.Int("size", len(sqsEvent.Records)))

	srv := service.GetFileProcessingService()

	for _, message := range sqsEvent.Records {

		data := &SqsBody{}
		err := json.Unmarshal([]byte(message.Body), &data)
		if err != nil {
			return "", err
		}

		// data.Records contains a single record always
		s3Bucket := data.Records[0].S3.Bucket.Name
		s3Key := data.Records[0].S3.Object.Key

		log.Debug(ctx, "S3 bucket details",
			log.String("Bucket", s3Bucket),
			log.String("key", s3Key))

		err = srv.Handle(ctx, s3Key)
		if err != nil {
			return "", err
		}

		log.Debug(ctx, "file successfully processed", log.String("key", s3Key))
	}
	return "Process finished successfully", nil
}

func main() {
	//Init core
	if err := core.Init(); err != nil {
		panic(err)
	}
	lambda.Start(HandleRequest)
}
