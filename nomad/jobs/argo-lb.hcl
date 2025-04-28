job "lb" {
  datacenters = ["dc1"]
  type = "service"

  group "lb" {

    count = 1

    constraint {
      distinct_hosts = true
    }

    network {
      mode = "bridge" 
      port "lb" {
        static = 9999
	    to = 9999
      }
      port "ui" {
        static = 9998
	    to = 9998
      }
      dns {
        servers = ["10.203.96.230", "10.203.96.231", "10.203.96.232"]  # Consul DNS IPs
      }
    }

    task "fabio" {
      driver = "docker"
      config {
        image = "fabiolb/fabio"
	    ports = ["lb", "ui"]
        args = [
          "-proxy.strategy=rr",
          "-registry.consul.addr=consul.service.consul:8500"
        ]
      }

      env {
        registry_consul_addr = "127.0.0.1:8500"
      }
    }

    service {
      name = "lb"
      port = "lb"

      tags = ["lb", "fabio"]

      check {
        type     = "http"
        path     = "/health"
        interval = "10s"
        timeout  = "2s"
        port     = "ui"
      }
    }
  }
}
