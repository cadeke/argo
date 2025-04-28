job "monitoring" {
  datacenters = ["dc1"]
  type        = "service"

  group "prometheus" {
    count = 1

    network {
      mode = "bridge"
      port "prometheus" {
	    static = 9090
        to = 9090
      }
      dns {
        servers  = ["10.203.96.230", "10.203.96.231", "10.203.96.232"]
      }
    }

    task "prometheus" {
      driver = "docker"
      config {
        image = "prom/prometheus:latest"
        args = [
          "--config.file=/etc/prometheus/prometheus.yml",
        ]
        ports = ["prometheus"]
        volumes = [
          "local/prometheus.yml:/etc/prometheus/prometheus.yml"
        ]
      }

      template {
        data = <<EOF
global:
  scrape_interval: 15s

scrape_configs:

  - job_name: "cadvisor"
    static_configs:
      - targets: ["cadvisor.service.consul:8080"]

  - job_name: "node-exporter"
    static_configs:
      - targets: ["node-exporter.service.consul:9100"]
EOF
        destination = "local/prometheus.yml"
      }

      service {
        name = "prometheus"
        port = "prometheus"
      }
    }
  }

  group "grafana" {
    count = 1

    network {
      mode = "bridge"
      port "grafana" {
	    static = 3000
        to = 3000
      }
      dns {
        servers  = ["10.203.96.230", "10.203.96.231", "10.203.96.232"]
      }
    }

    volume "grafana" {
      type      = "host"
      source    = "grafana-storage"
      read_only = false
    }

    task "grafana" {
      driver = "docker"
      config {
        image = "grafana/grafana:latest"
        ports = ["grafana"]
      }

    volume_mount {
      volume = "grafana"
      destination = "/var/lib/grafana"
      read_only = false
    }

      env {
        GF_SECURITY_ADMIN_USER     = "admin"
        GF_SECURITY_ADMIN_PASSWORD = "admin"
      }

      service {
        name = "grafana"
        port = "grafana"
        #tags = ["urlprefix-/grafana"]
      }
    }
  }
}
