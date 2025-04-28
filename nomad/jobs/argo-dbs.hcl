job "dbs" {
  datacenters = ["dc1"]
  type        = "service"

  group "postgres" {
    count = 1

    constraint {
      attribute = "${node.unique.name}"
      value     = "vm04"
    }

    network {
      mode = "host"
      port "postgres" {
	    static = 5432
        to = 5432 
      }
    }

    task "postgres" {
      driver = "docker"
      config {
        image = "postgres:latest"
        ports = ["postgres"]
	    volumes = ["pgdata:/var/lib/postgresql/data"]
      }

      env {
        POSTGRES_USER     = "admin"
        POSTGRES_PASSWORD = "admin"
        POSTGRES_DB       = "argodb"
      }
    }

    service {
        name = "postgres"
        port = "postgres"
        tags = ["db"]
    }

    volume "pgdata" {
        type      = "host"
        source    = "postgres-storage"
        read_only = false
    }
  }

  group "memcached" {
    count = 1

    network {
      mode = "host"
      port "memcached" {
        static = 11211 
        to = 11211 
      }
    }

    task "memcached" {
      driver = "docker"
      config {
        image = "memcached:latest"
        command = "-m 64"
        ports = ["memcached"]
      }
    }

    service {
      name = "memcached"
      port = "memcached"
      tags = ["cache"]
    }
  }
}
