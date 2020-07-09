output "external_ip_peer1" {
  value = google_compute_instance.peer1.network_interface[0].access_config[0].nat_ip
}

output "internal_ip_peer1" {
  value = google_compute_instance.peer1.network_interface[0].network_ip
}

output "external_ip_peer2" {
  value = google_compute_instance.peer2.network_interface[0].access_config[0].nat_ip
}

output "internal_ip_peer2" {
  value = google_compute_instance.peer2.network_interface[0].network_ip
}