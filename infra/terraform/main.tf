resource "aws_sqs_queue" "deliveries_dlq" {
  name = "phishguard-deliveries-dlq"
}

resource "aws_sqs_queue" "deliveries" {
  name                       = "phishguard-deliveries"
  visibility_timeout_seconds = 30
  receive_wait_time_seconds  = 10
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.deliveries_dlq.arn
    maxReceiveCount     = 3
  })
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "${path.module}/../../backend/bin/bootstrap"
  output_path = "${path.module}/../../backend/bin/function.zip"
}

resource "aws_iam_role" "lambda" {
  name = "phishguard-deliver-lambda"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
    }]
  })
}

resource "aws_lambda_function" "deliver" {
  function_name    = "phishguard-deliver"
  role             = aws_iam_role.lambda.arn
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  timeout          = 30

  environment {
    variables = {
      DATABASE_URL   = "postgres://phishguard:phishguard@host.docker.internal:5432/phishguard?sslmode=disable"
      REDIS_URL      = "redis://host.docker.internal:6379/0"
      SQS_QUEUE_URL  = aws_sqs_queue.deliveries.url
      AWS_REGION     = var.aws_region
      PUBLIC_APP_URL = "http://localhost:5173"
    }
  }
}

resource "aws_lambda_event_source_mapping" "deliveries" {
  event_source_arn = aws_sqs_queue.deliveries.arn
  function_name    = aws_lambda_function.deliver.arn
  batch_size       = 5
  enabled          = true
}

output "queue_url" {
  value = aws_sqs_queue.deliveries.url
}

output "lambda_name" {
  value = aws_lambda_function.deliver.function_name
}
