# all kvm hosts
data "maas_vm_hosts" "vm-hosts" {
  id = "rack-01"
}