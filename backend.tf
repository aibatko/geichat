resource "docker_image" "backend_image" {
  name = "local-backend:latest"
  build {
    context    = "./backend"
    dockerfile = "Dockerfile"
  }
  keep_locally = true
}

resource "docker_container" "backend" {
  name  = "backend"
  image = docker_image.backend_image.image_id

  ports {
    internal = 8080
    external = var.backend_port
  }

  env = local.backend_env

  networks_advanced {
    name = docker_network.app_network.name
  }

  depends_on = [docker_container.db]
  restart    = "unless-stopped"

  labels {
    label = "app.component"
    value = "backend"
  }

  labels {
    label = "app.environment"
    value = "local"
  }

  labels {
    label = "app.managed-by"
    value = "terraform"
  }

  provisioner "local-exec" {
    command = "sleep 10"
  }
}