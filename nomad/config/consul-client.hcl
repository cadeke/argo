datacenter = "dc1"
data_dir = "/opt/consul"

bind_addr = "10.203.96.233"
encrypt = "Y0fMIuijSymLkCVXin064ZYF2FUKugBleIB8yiiWKqU="
retry_join = ["vm01"]

connect {
  enabled = true
}

addresses {
  grpc = "0.0.0.0"
}

ports {
  grpc = 8502
}
