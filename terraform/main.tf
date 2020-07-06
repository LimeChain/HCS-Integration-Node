resource "google_compute_address" "static" {
  name = "peer2-static-ip-address"
  project = "hcs-integration-node"
  region    = "us-central1"
}

resource "google_compute_instance" "peer2" {
  name         = "peer2"
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"
  project      = "hcs-integration-node"

  tags = ["hcs-node"]

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1604-lts"
    }
  }

  // Local SSD disk
  scratch_disk {
    interface = "SCSI"
  }

  network_interface {
    network = "default"

    access_config {
      nat_ip = google_compute_address.static.address
    }
  }

  metadata = {
    sshKeys = "ubuntu:${file(var.ssh_public_key_filepath)}"
  }
}