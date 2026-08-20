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
  type             = "zip"
  source_file      = "${path.module}/../../backend/bin/bootstrap"
  output_path      = "${path.module}/../../backend/bin/function.zip"
  output_file_mode = "0755"
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

resource "aws_iam_role_policy_attachment" "lambda_sqs" {
  role       = aws_iam_role.lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole"
}

resource "aws_lambda_function" "deliver" {
  function_name    = "phishguard-deliver"
  role             = aws_iam_role.lambda.arn
  handler          = "bootstrap"
  runtime          = "provided.al2023"
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256
  timeout          = 5

  environment {
    variables = {
      DATABASE_URL   = var.database_url
      REDIS_URL      = var.redis_url
      SQS_QUEUE_URL  = aws_sqs_queue.deliveries.url
      PUBLIC_APP_URL = var.public_app_url
    }
  }
}

resource "aws_lambda_event_source_mapping" "deliveries" {
  event_source_arn = aws_sqs_queue.deliveries.arn
  function_name    = aws_lambda_function.deliver.arn
  batch_size       = 5
  enabled          = true
}
