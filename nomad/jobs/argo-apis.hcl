job "apis" {
  datacenters = ["dc1"]
  type        = "service"

  group "query-api" {
    count = 1

    network {
      mode = "bridge"
      port "http" { to = 8080 }
      dns {
      	servers  = ["10.203.96.230", "10.203.96.231", "10.203.96.232"]
      }
    }

    service {
      name = "query-api"
      port = "http"

      check {
        type     = "http"
        path     = "/health"
        interval = "10s"
        timeout  = "2s"
	    port     = "http"
	    address_mode = "alloc"
      }

	  tags = [
            "urlprefix-/api/ip2domain",
            "urlprefix-/api/domain2ip"
          ]
	    }
      }
    }

    task "query-api" {
      driver = "docker"
      config {
        image = "cadeke/argo-q-api:latest"
        ports = ["http"]
      }

      env {
        POSTGRES_HOST     = "postgres.service.consul"
        POSTGRES_PORT     = 5432
        POSTGRES_USER     = "admin"
        POSTGRES_PASSWORD = "admin"
        POSTGRES_DB       = "argodb"
        MEMCACHED_HOST    = "memcache.service.consul"
        MEMCACHED_PORT    = 11211
      }
    }
  }

  group "debug-query-api" {
    count = 1

    network {
      mode = "bridge"
      dns {
      	servers  = ["10.203.96.230", "10.203.96.231", "10.203.96.232"]
      }
    }

    task "query-api-debug" {
      driver = "docker"
      config {
        image = "nicolaka/netshoot:latest"
        args = ["sleep", "infinity"]
      }
    }

    service {
      name = "debug-query-api"
    }
  }

  group "admin-api" {
    count = 1

    network {
      mode = "bridge"
      port "http" { to = 8080 }
      dns {
        servers  = ["10.203.96.230", "10.203.96.231", "10.203.96.232"]
      }
    }

    service {
      name = "admin-api"
      port = "http"
      tags = [
        "urlprefix-/api/get",
        "urlprefix-/api/list",
        "urlprefix-/api/add",
        "urlprefix-/api/delete",
        "urlprefix-/api/update"
      ] 

      check {
        type     = "http"
        path     = "/health"
        interval = "10s"
        timeout  = "2s"
      }
    }

    task "admin-api" {
      driver = "docker"
      config {
        image = "cadeke/argo-a-api:latest"
        ports = ["http"]
      }

      env {
        POSTGRES_HOST     = "postgres.service.consul"
        POSTGRES_PORT     = 5432
        POSTGRES_USER     = "admin"
        POSTGRES_PASSWORD = "admin"
        POSTGRES_DB       = "argodb"
      }
    }
  }
}
