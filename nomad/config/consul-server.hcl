datacenter = "dc1"
data_dir = "/opt/consul"
server = true
bootstrap_expect = 3

bind_addr = "10.203.96.230"
client_addr = "0.0.0.0"

ui_config {
  enabled = true
}

encrypt = "Y0fMIuijSymLkCVXin064ZYF2FUKugBleIB8yiiWKqU="
retry_join = ["vm02", "vm03"]

connect {
  enabled = true
}

addresses {
  grpc = "0.0.0.0"
}

ports {
  grpc = 8502
}
