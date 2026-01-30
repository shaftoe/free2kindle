variable "project_id" {
  description = "Google Cloud project ID"
  type        = string
  default     = "free2kindle"
}

variable "region" {
  description = "Google Cloud region"
  type        = string
  default     = "us-central1"
}

variable "source_bucket" {
  description = "GCS bucket name for function source code"
  type        = string
  default     = "free2kindle-source"
}

variable "source_zip_path" {
  description = "Local path to the function source code zip file"
  type        = string
  default     = "../bin/free2kindle-source.zip"
}

variable "environment_variables" {
  description = "Environment variables for the Cloud Function (Mailjet API keys, API secret, etc.)"
  type        = map(string)
  sensitive   = true
}
