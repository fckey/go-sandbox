provider "google" {
  credentials = "${file("terraform-key.json")}"
  project     = "$GCP_PROJECT_ID"
  region      = "us-central1"
}

resource "google_composer_environment" "composer_terraform" {
  provider = "google"
  project = "$GCP_PROJECT_ID"
  region  = "us-central1"
  name    = "terra-composer"
  config {
    node_count = "3"
    node_config {
      zone            = "us-central1-c"
      machine_type    = "n1-standard-2"
      disk_size_gb    = "32"
    }
    software_config {
      image_version            = "composer-1.9.0-airflow-1.10.6"
      python_version           = "3"
    }
  }
  lifecycle {
    prevent_destroy = true
  }
}
