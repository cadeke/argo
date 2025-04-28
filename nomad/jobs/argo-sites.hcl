job "sites" {
  datacenters = ["dc1"]
  type        = "service"

  group "query-site" {
    count = 1

    network {
      mode = "bridge"
      port "http" { to = 80 }
    }

    service {
      name = "query-site"
      port = "http"
      tags = ["urlprefix-/query/ strip=/query/"]
      check {
        type     = "http"
        path     = "/"
        interval = "10s"
        timeout  = "2s"
      }
    }

    task "query-site" {
      driver = "docker"
      config {
        image = "cadeke/argo-q-site:nomad"
        ports = ["http"]
      }
    }
  }

  group "admin-site" {
    count = 1

    network {
      mode = "bridge"
      port "http" { to = 80 }
    }

    service {
      name = "admin-site"
      port = "http"
      tags = ["urlprefix-/admin/ strip=/admin/"]
      check {
        type     = "http"
        path     = "/"
        interval = "10s"
        timeout  = "2s"
      }
    }

    task "admin-site" {
      driver = "docker"
      config {
        image = "cadeke/argo-a-site:nomad"
        ports = ["http"]
      }
    }
  }
}
