job "ot" {
  datacenters = ["dc1"]
  type        = "service"

  group "ot-app" {
    count = 1

    network {
      mode = "bridge"
    }

    service {
      name = "ot-app"

      connect {
        sidecar_service {
	      proxy {
	        upstreams {
	          destination_name = "query-api"
              local_bind_port = 8080
	        }
	      }
	    }
      }
    }

    task "ot-app" {
      driver = "docker"
      config {
        image = "cadeke/argo-ot-app:nomad"
      }
    }
  }
}
