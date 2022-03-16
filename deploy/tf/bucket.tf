resource "aws_s3_bucket_notification" "bucket_to_sqs" {
  bucket = var.inbound_bucket

  queue {
    events        = ["s3:ObjectCreated:*"]
    queue_arn     = aws_sqs_queue.inbound.arn
    filter_prefix = "assets/"
  }
}

resource "aws_s3_bucket" "outbound" {
  bucket = var.outbound_bucket
}

#resource "aws_s3_bucket_acl" "outbound" {
#  bucket = aws_s3_bucket.outbound.id
#  acl    = "private"
#}

#resource "aws_s3_bucket_versioning" "outbound" {
#  bucket = aws_s3_bucket.outbound.id
#  versioning_configuration {
#    status = "Disabled"
#  }
#}

data "aws_iam_policy_document" "outbound_allow_access_from_cloudfront_automation_s3" {
  statement {
    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::cloudfront:user/CloudFront Origin Access Identity ET8B6931D9SDO"]
    }
    actions = ["s3:GetObject"]
    resources = [
      "${aws_s3_bucket.outbound.arn}/*",
    ]
  }
  statement {
    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::942095822719:user/jenkins-terraform-user"]
    }
    actions = ["s3:*"]
    resources = [
      aws_s3_bucket.outbound.arn,
      "${aws_s3_bucket.outbound.arn}/*",
    ]
  }
  statement {
    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/kidsloop-global-loadtest-k8s-cms-bucket-access"]
    }
    actions = ["s3:*Object*"]
    resources = [
      aws_s3_bucket.outbound.arn,
      "${aws_s3_bucket.outbound.arn}/*",
    ]
  }
}

resource "aws_s3_bucket_policy" "outbound_allow_access_from_cloudfront" {
  bucket = aws_s3_bucket.outbound.id
  policy = data.aws_iam_policy_document.outbound_allow_access_from_cloudfront_automation_s3.json
}