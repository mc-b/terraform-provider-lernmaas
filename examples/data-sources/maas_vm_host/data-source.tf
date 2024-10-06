# only kvm host with name cloud-hf-39
data "maas_vm_host" "kvm-01" {
  name = "cloud-hf-39"
}

data "maas_vm_host" "kvm-02" {
  name = "cloud-hf-20"
}