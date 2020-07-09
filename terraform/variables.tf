variable "project_id" {
  default = "hcs-integration-node"
}

variable "project_zone" {
  default = "us-central1-a"
}

variable "project_region" {
  default = "us-central1"
}

variable "machine_type" {
  default = "n1-standard-1"
}
variable "boot_disk_image" {
  default = "ubuntu-os-cloud/ubuntu-1604-lts"
}

variable "ssh_operator_peer1" {
  default = "ubuntu_peer1"
}

variable "ssh_operator_peer2" {
  default = "ubuntu_peer2"
}

variable "ssh_public_key_filepath_peer1" {
  default = "ubuntu_peer1.pub"
}

variable "ssh_public_key_filepath_peer2" {
  default = "ubuntu_peer2.pub"
}

variable "api_port_peer1" {
  default = 8181
}

variable "api_port_peer2" {
  default = 8182
}