output "application_urls" {
  description = "Application access URLs"
  value = {
    frontend = "http://localhost:${var.frontend_port}"
    backend  = "http://localhost:${var.backend_port}"
    database = "localhost:${var.database_port}"
  }
}

output "database_connection" {
  description = "Database connection details"
  value = {
    host     = "localhost"
    port     = var.database_port
    database = var.postgres_db
    username = var.postgres_user
  }
  sensitive = false
}

output "container_status" {
  description = "Container status information"
  value = {
    frontend_id = docker_container.frontend.id
    backend_id  = docker_container.backend.id
    database_id = docker_container.db.id
    network_id  = docker_network.app_network.id
    volume_id   = docker_volume.db_data.name
  }
}

output "environment_config" {
  description = "Environment configuration summary"
  value = {
    postgres_db = var.postgres_db
    postgres_user = var.postgres_user
    frontend_port = var.frontend_port
    backend_port = var.backend_port
    database_port = var.database_port
  }
  sensitive = false
}