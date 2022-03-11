package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"gitlab.badanamu.com.cn/calmisland/kidsloop-file-processing-service/config"
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

	region, err := config.GetRegion()
	if err != nil {
		return "", err
	}
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return "", err
	}

	//load config
	config.MustLoad(awsSession)

	//Init core
	err = core.Init()
	if err != nil {
		fmt.Println("Init core failed, err:", err)
		return "", err
	}

	fmt.Printf("SQS Record size = %d\n", len(sqsEvent.Records))

	for _, message := range sqsEvent.Records {

		data := &SqsBody{}
		err := json.Unmarshal([]byte(message.Body), &data)
		if err != nil {
			return "", err
		}

		bucket := data.Records[0].S3.Bucket.Name
		key := data.Records[0].S3.Object.Key
		fmt.Printf("Bucket = %s, Key = %s \n", bucket, key)

		//start Handler
		srv := service.GetFileProcessingService()
		err = srv.Handle(ctx, key)
		if err != nil {
			return "", err
		}
		fmt.Printf("file %s successfully processed \n", key)
	}
	return "Process finished successfully", nil
}

func main() {
	lambda.Start(HandleRequest)
}
