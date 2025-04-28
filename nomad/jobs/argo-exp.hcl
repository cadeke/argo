job "exporters" {
  datacenters = ["dc1"]
  type        = "system"

  group "cadvisor" {
    constraint {
      operator = "distinct_hosts"
      value    = "true"
    }

    network {
      mode = "bridge"
      port "cadvisor" {
        static = 8080
        to     = 8080
      }
      dns {
        servers = ["10.203.96.230", "10.203.96.231", "10.203.96.232"]
      }
    }
  
    task "cadvisor" {
      driver = "docker"
      config {
        image = "gcr.io/cadvisor/cadvisor:latest"
        args  = [
          "-docker_only=true",
          "-housekeeping_interval=10s"
        ]
        volumes = [
          "/var/run/docker.sock:/var/run/docker.sock:ro",
          "/sys:/sys:ro",
          "/var/lib/docker/:/var/lib/docker:ro"
        ]
      }
    }

    service {
      name = "cadvisor"
      port = "cadvisor"
    }
  } 

  group "node-exporter" {
    constraint {
      operator = "distinct_hosts"
      value    = "true"
    }

    network {
      mode = "bridge"
      port "node-exporter" {
        static = 9100
        to     = 9100
      }
      dns {
        servers = ["10.203.96.230", "10.203.96.231", "10.203.96.232"]
      }
    }

    task "node-exporter" {
      driver = "docker"
      config {
        image = "prom/node-exporter:latest"
      }
    }

    service {
      name = "node-exporter"
      port = "node-exporter"
    }
  }
}
