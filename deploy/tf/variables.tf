variable "region" {
  type        = string
  description = "AWS region"
}

variable "lambda_function_name" {
  type = string
}

variable "lambda_function_image" {
  type = string
}

variable "lambda_function_timeout" {
  type    = number
  default = 60
}

variable "lambda_function_memory" {
  type    = number
  default = 512
}

variable "inbound_bucket" {
  type = string
}

variable "outbound_bucket" {
  type = string
}

variable "enable_trigger" {
  type    = bool
  default = true
}

variable "storage_driver" {
  type    = string
  default = "s3"
}

variable "storage_accelerate" {
  type    = bool
  default = false
}

variable "log_level" {
  type    = string
  default = "debug"
}

variable "log_std_out" {
  type    = bool
  default = true
}

variable "processors" {
  type    = list(string)
  default = ["exif"]
}

variable "retries" {
  type    = number
  default = 2
}
