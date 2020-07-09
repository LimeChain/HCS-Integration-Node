resource "google_compute_firewall" "default" {
  name    = "hcs-api-firewall"
  network = "default"
  project = var.project_id

  allow {
    protocol = "tcp"
    ports    = [var.api_port_peer1, var.api_port_peer2]
  }
}

resource "google_compute_address" "static_peer1" {
  name    = "peer1-static-ip-address"
  project = var.project_id
  region  = var.project_region
}

resource "google_compute_address" "static_peer2" {
  name      = "peer2-static-ip-address"
  project   = var.project_id
  region    = var.project_region
}

resource "google_compute_instance" "peer1" {
  name         = "peer1"
  machine_type = var.machine_type
  zone         = var.project_zone
  project      = var.project_id

  boot_disk {
    initialize_params {
      image = var.boot_disk_image
    }
  }

  scratch_disk {
    interface = "SCSI"
  }

  network_interface {
    network = "default"

    access_config {
      nat_ip = google_compute_address.static_peer1.address
    }
  }

  metadata = {
    sshKeys = "${var.ssh_operator_peer1}:${file(var.ssh_public_key_filepath_peer1)}"
  }
}

resource "google_compute_instance" "peer2" {
  name         = "peer2"
  machine_type = var.machine_type
  zone         = var.project_zone
  project      = var.project_id

  tags = ["hcs-node"]

  boot_disk {
    initialize_params {
      image = var.boot_disk_image
    }
  }

  // Local SSD disk
  scratch_disk {
    interface = "SCSI"
  }

  network_interface {
    network = "default"

    access_config {
      nat_ip = google_compute_address.static_peer2.address
    }
  }

  metadata = {
    sshKeys = "${var.ssh_operator_peer2}:${file(var.ssh_public_key_filepath_peer2)}"
  }
}