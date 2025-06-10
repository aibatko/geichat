resource "docker_image" "db_image" {
  name = "local-db:latest"
  build {
    context    = "./db"
    dockerfile = "Dockerfile"
  }
  keep_locally = true
}

resource "docker_container" "db" {
  name  = "db"
  image = docker_image.db_image.image_id

  ports {
    internal = 5432
    external = var.database_port
  }

  env = local.database_env

  volumes {
    volume_name    = docker_volume.db_data.name
    container_path = "/var/lib/postgresql/data"
  }

  healthcheck {
    test     = ["CMD-SHELL", "pg_isready -U ${var.postgres_user} -d ${var.postgres_db}"]
    interval = "5s"
    timeout  = "5s"
    retries  = 5
  }

  networks_advanced {
    name = docker_network.app_network.name
  }

  restart = "unless-stopped"

  labels {
    label = "app.component"
    value = "database"
  }

  labels {
    label = "app.environment"
    value = "local"
  }

  labels {
    label = "app.managed-by"
    value = "terraform"
  }
}