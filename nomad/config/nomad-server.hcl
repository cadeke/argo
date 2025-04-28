datacenter = "dc1"
data_dir = "/opt/nomad/data"
log_level = "INFO"

server {
  enabled = true
  bootstrap_expect = 3
  encrypt = "jAJBQFP8uM0+VnspMh4ge1ICqFyKOBncgqlFzA1BHGA="
}

consul {
  address = "127.0.0.1:8500"
  grpc_address = "127.0.0.1:8502"
  server_auto_join = true
  client_auto_join = true
  server_service_name = "nomad"
}

tls {
  http = true
  rpc  = true

  ca_file   = "/etc/nomad.d/nomad-agent-ca.pem"
  cert_file = "/etc/nomad.d/vm01-dc1-server-nomad.pem"
  key_file  = "/etc/nomad.d/vm01-dc1-server-nomad-key.pem"

  verify_server_hostname = true
  verify_https_client    = false
}
