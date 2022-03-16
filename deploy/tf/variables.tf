variable "region" {
  type        = string
  description = "AWS region"
}

variable "lambda_function_name" {
  type = string
  description = "lambda function name"
}

variable "lambda_function_image" {
  type = string
  description = "lambda function container image"
}

variable "lambda_function_timeout" {
  type    = number
  default = 60
  description = "lambda runtime timeout"
}

variable "lambda_function_memory" {
  type    = number
  default = 512
  description = "lambda max memory available"
}

variable "lambda_function_retries" {
  type    = number
  default = 2
  description = "Lambda number of retries for synchronous requests"
}

variable "inbound_bucket" {
  type = string
  description = "the file upload bucket name"
}

variable "outbound_bucket" {
  type = string
  description = "the processed file bucket name"
}

variable "enable_trigger" {
  type    = bool
  default = true
  description = "Enable or disable the lambda SQS trigger"
}

variable "env_storage_driver" {
  type    = string
  default = "s3"
  description = "env var for storage type"
}

variable "env_storage_accelerate" {
  type    = bool
  default = false
  description = "env var for storage acceleration"
}

variable "env_log_level" {
  type    = string
  default = "debug"
  description = "env var for logging level"
}

variable "env_log_std_out" {
  type    = bool
  default = true
  description = "env var for logging to standard out"
}

variable "env_processors" {
  type    = list(string)
  default = ["exif"]
  description = "env var for the process types to execute on files"
}
