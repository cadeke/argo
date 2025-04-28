job "apis" {
  datacenters = ["dc1"]
  type        = "service"

  group "debug" {
    count = 1

    network {
      mode = "bridge"
      dns {
      	servers  = ["10.203.96.230", "10.203.96.231", "10.203.96.232"]
      }
    }

    task "debug" {
      driver = "docker"
      config {
        image = "nicolaka/netshoot:latest"
        args = ["sleep", "infinity"]
      }
    }

    service {
      name = "debug"
      tags = ["debug"]
      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "postgres"
	          local_bind_port = 5432
            }
            upstreams {
              destination_name = "memcached"
              local_bind_port = 11211
            }
	        upstreams {
              destination_name = "query-api"
              local_bind_port = 8080
            }
          }
        }
      }
    }
  }

  group "query-api" {
    count = 1

    network {
      mode = "bridge"
      port "http" {
	    static = 8080
        to = 8080
      }
      dns {
        servers = [
          "10.203.96.230",
          "10.203.96.231",
          "10.203.96.232",
        ]
      }
    }

    task "query-api" {
      driver = "docker"
      config {
        image = "cadeke/argo-q-api:nomad"
        ports = ["http"]
	    force_pull = true
      }

      env {
        POSTGRES_USER       = "admin"
        POSTGRES_PASSWORD   = "admin"
        POSTGRES_DB         = "argodb"
      }
    }

    service {
      name = "query-api"
      port = "http"
      tags = ["api"]

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "postgres"
	          local_bind_port = 5432
            }
            upstreams {
              destination_name = "memcached"
              local_bind_port = 11211
            }
          }
        }
      }
    }
  }
}
