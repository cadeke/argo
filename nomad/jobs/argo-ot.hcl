job "ot" {
  datacenters = ["dc1"]
  type        = "service"

  group "ot-app" {
    count = 4

    update {
      max_parallel     = 1
      canary           = 0
      min_healthy_time = "10s"
      healthy_deadline = "3m"
      progress_deadline = "10m"
      auto_revert      = true
      stagger          = "30s"
    } 

    network {
      mode = "bridge"
      dns {
        servers = ["10.203.96.230", "10.203.96.231", "10.203.96.232"]
      }
    }

    service {
      name = "ot-app"

      connect {
        sidecar_service {}
      }
    }

    task "ot-app" {
      driver = "docker"
      config {
        image = "cadeke/argo-ot-app:v1.0"
      }

      env {
        QUERY_HOST = "lb.service.consul"
        QUERY_PORT = 9999
      }

template {
  data = <<EOH

API_KEY="1111-2222-3333-4444"
SOME_SECRET="some-value"
EOH

  destination = "secrets/argo.env"
  env         = true
}
    }
  }
}
