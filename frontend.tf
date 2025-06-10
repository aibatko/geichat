resource "docker_image" "frontend_image" {
  name = "local-frontend:latest"
  build {
    context    = "./frontend"
    dockerfile = "Dockerfile"
  }
  keep_locally = true
}

resource "docker_container" "frontend" {
  name  = "frontend"
  image = docker_image.frontend_image.image_id

  ports {
    internal = 80
    external = var.frontend_port
  }

  networks_advanced {
    name = docker_network.app_network.name
  }

  depends_on = [docker_container.backend]
  restart    = "unless-stopped"

  labels {
    label = "app.component"
    value = "frontend"
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