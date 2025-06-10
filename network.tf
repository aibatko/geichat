resource "docker_network" "app_network" {
  name   = "app-network"
  driver = "bridge"
}

resource "docker_volume" "db_data" {
  name = "db_data"
}