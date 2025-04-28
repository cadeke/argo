datacenter = "dc1"
data_dir = "/opt/nomad/data"
plugin_dir = "/opt/nomad/plugins"
log_level = "INFO"

client {
  enabled = true
  servers = ["10.203.96.230", "10.203.96.231", "10.203.96.232"]
  
  host_volume "postgres-storage" {
    path = "/mnt/argo-storage/postgres"
    read_only = false
  }

  host_volume "grafana-storage" {
    path = "/mnt/argo-storage/grafana"
    read_only = false
  }

  cni_path = "/opt/cni/bin"
  bridge_network_name = "nomad"
  bridge_network_subnet = "172.26.64.0/20"
}

plugin "docker" {
  config {
    allow_privileged = true
    allow_caps = ["NET_ADMIN", "NET_BROADCAST", "NET_RAW", "CHOWN", "SETGID", "SETUID"]
    volumes { enabled = true }
  }
}

consul {
  address = "127.0.0.1:8500"
  grpc_address = "127.0.0.1:8502"
  client_auto_join = true
  auto_advertise = true
  client_service_name = "nomad-client"
}

tls {
  http = true
  rpc  = true

  ca_file   = "/etc/nomad.d/nomad-agent-ca.pem"
  cert_file = "/etc/nomad.d/vm04-dc1-client-nomad.pem"
  key_file  = "/etc/nomad.d/vm04-dc1-client-nomad-key.pem"

  verify_server_hostname = true
  verify_https_client    = true
}
