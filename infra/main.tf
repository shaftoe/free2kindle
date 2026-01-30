terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_cloudfunctions2_function" "api" {
  name        = "free2kindle-api"
  location    = var.region
  description = "Free2Kindle API - Article conversion and delivery service"

  build_config {
    runtime     = "go122"
    entry_point = "HTTP"
    source {
      storage_source {
        bucket = var.source_bucket
        object = google_storage_bucket_object.function_source.name
      }
    }
  }

  service_config {
    available_memory               = "256M"
    timeout_seconds                = 540
    max_instance_count             = 10
    min_instance_count             = 0
    environment_variables          = var.environment_variables
    ingress_settings               = "ALLOW_ALL"
    all_traffic_on_latest_revision = true
  }
}

resource "google_cloudfunctions2_function_iam_member" "public_invoker" {
  project        = google_cloudfunctions2_function.api.project
  location       = google_cloudfunctions2_function.api.location
  cloud_function = google_cloudfunctions2_function.api.name
  role           = "roles/cloudfunctions.invoker"
  member         = "allUsers"
}

resource "google_storage_bucket" "source" {
  name                        = "${var.project_id}-free2kindle-source"
  location                    = var.region
  force_destroy               = true
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "function_source" {
  name   = "free2kindle-function-source.zip"
  bucket = google_storage_bucket.source.name
  source = var.source_zip_path
}
