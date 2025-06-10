locals {
  common_labels = {
    app_name        = "docker-app"
    app_environment = "local"
    app_managed_by  = "terraform"
    app_created     = formatdate("YYYY-MM-DD_hh-mm-ss", timestamp())
  }

  port_mappings = {
    frontend = {
      internal = 80
      external = var.frontend_port
    }
    backend = {
      internal = 8080
      external = var.backend_port
    }
    database = {
      internal = 5432
      external = var.database_port
    }
  }

  backend_env = [
    "DB_HOST=db",
    "DB_USER=${var.postgres_user}",
    "DB_PASSWORD=${var.postgres_password}",
    "DB_NAME=${var.postgres_db}",
    "JWT_SECRET=${var.jwt_secret}"
  ]

  database_env = [
    "POSTGRES_DB=${var.postgres_db}",
    "POSTGRES_USER=${var.postgres_user}",
    "POSTGRES_PASSWORD=${var.postgres_password}"
  ]
}