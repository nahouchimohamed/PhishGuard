output "queue_url" {
  value       = aws_sqs_queue.deliveries.url
  description = "SQS queue URL for phishguard-deliveries (LocalStack)."
}

output "lambda_name" {
  value       = aws_lambda_function.deliver.function_name
  description = "Name of the deliver Lambda function."
}

output "dlq_url" {
  value       = aws_sqs_queue.deliveries_dlq.url
  description = "Dead-letter queue URL for failed deliveries."
}
