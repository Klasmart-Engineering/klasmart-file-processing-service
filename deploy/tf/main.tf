resource "aws_lambda_function" "fps" {
  function_name = var.lambda_function_name
  description   = "The file processing service lambda"
  role          = aws_iam_role.fps.arn
  architectures = ["x86_64"]
  image_uri     = var.lambda_function_image
  package_type  = "Image"
  memory_size   = var.lambda_function_memory
  timeout       = var.lambda_function_timeout

  environment {
    variables = {
      log_level          = var.log_level
      log_std_out        = var.log_std_out
      processors         = join(",", var.processors)
      storage_accelerate = var.storage_accelerate
      storage_bucket     = var.inbound_bucket
      storage_bucket_out = var.outbound_bucket
      storage_driver     = var.storage_driver
      storage_region     = var.region
    }
  }

  dead_letter_config {
    target_arn = aws_sqs_queue.dlq.arn
  }

  depends_on = [
    aws_iam_role.fps,
    aws_cloudwatch_log_group.fps,
  ]
}

resource "aws_lambda_event_source_mapping" "trigger" {
  batch_size       = 6
  enabled          = var.enable_trigger
  event_source_arn = aws_sqs_queue.inbound.arn
  function_name    = var.lambda_function_name

  depends_on = [
    aws_lambda_function.fps,
  ]
}

resource "aws_cloudwatch_log_group" "fps" {
  name              = "/aws/lambda/${var.lambda_function_name}"
  retention_in_days = 14
}

resource "aws_iam_role" "fps" {
  name = var.lambda_function_name
  path = "/service-role/"

  assume_role_policy = jsonencode(
    {
      Version = "2012-10-17"
      Statement = {
        Action = ["sts:AssumeRole"]
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    }
  )
}

resource "aws_iam_policy" "fps" {
  name        = "file_processing_lambda_sqs"
  path        = "/"
  description = "IAM policy for file processing service lambda"

  policy = jsonencode(
    {
      Version = "2012-10-17"
      Statement = [
        {
          Action = [
            "sqs:DeleteMessage",
            "sqs:ReceiveMessage",
            "sqs:GetQueueAttributes",
          ]
          Effect   = "Allow"
          Resource = aws_sqs_queue.inbound.arn
        },
        {
          Action   = "sqs:SendMessage"
          Effect   = "Allow"
          Resource = aws_sqs_queue.dlq.arn
        },
        {
          Action   = "logs:CreateLogGroup"
          Effect   = "Allow"
          Resource = "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:*"
        },
        {
          Action = [
            "logs:CreateLogStream",
            "logs:PutLogEvents"
          ]
          Effect   = "Allow"
          Resource = "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${var.lambda_function_name}:*"
        },
        {
          Action = [
            "s3:PutObject",
            "s3:GetObject"
          ]
          Effect = "Allow"
          Resource = [
            "arn:aws:s3:::${var.inbound_bucket}/*",
            "arn:aws:s3:::${var.outbound_bucket}/*"
          ]
        }
      ]
    }
  )
}

resource "aws_iam_role_policy_attachment" "fps" {
  role       = aws_iam_role.fps.name
  policy_arn = aws_iam_policy.fps.arn
}
