resource "aws_sqs_queue" "dlq" {
  name                      = "file_processing_service_lambda_dlq"
  delay_seconds             = 0
  max_message_size          = 262144
  message_retention_seconds = 345600
  receive_wait_time_seconds = 0

  redrive_allow_policy = jsonencode({
    redrivePermission = "byQueue",
    sourceQueueArns   = ["arn:aws:sqs:${var.region}:${data.aws_caller_identity.current.account_id}:file_processing_service_lambda"]
  })
}

resource "aws_sqs_queue" "inbound" {
  name                       = "file_processing_service_lambda"
  delay_seconds              = 0
  max_message_size           = 262144
  message_retention_seconds  = 345600
  receive_wait_time_seconds  = 0
  visibility_timeout_seconds = 70
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dlq.arn
    maxReceiveCount     = var.retries
  })
}

resource "aws_sqs_queue_policy" "inbound" {
  queue_url = aws_sqs_queue.inbound.id

  policy = jsonencode(
    {
      Version = "2012-10-17"
      Statement = [
        {
          Action = "SQS:SendMessage"
          Effect = "Allow"
          Principal = {
            Service = "s3.amazonaws.com"
          }
          Resource = aws_sqs_queue.inbound.arn
          Condition = {
            StringEquals = {
              "aws:SourceAccount" = data.aws_caller_identity.current.account_id
            },
            ArnLike = {
              "aws:SourceArn" = "arn:aws:s3:::${var.inbound_bucket}"
            }
          }
        }
      ]
    }
  )
}
