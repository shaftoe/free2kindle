output "function_url" {
  description = "The HTTP URL of the deployed Cloud Function"
  value       = google_cloudfunctions2_function.api.service_config[0].uri
}

output "source_bucket_name" {
  description = "Name of the GCS bucket for source code"
  value       = google_storage_bucket.source.name
}
