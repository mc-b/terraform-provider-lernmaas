# all kvm hosts
data "maas_vm_hosts" "vm-hosts" {
  id  = "rack-01"
}

# 2 vms for vm host
resource "maas_vm_instance" "base" {
  count = length(data.maas_vm_hosts.vm-hosts.no) * 2
  kvm_no = data.maas_vm_hosts.vm-hosts.no[count.index % length(data.maas_vm_hosts.vm-hosts.no)]
  hostname = "base-${format("%02d", count.index + 10)}"
  description = "student ${format("%02d", count.index + 1)}"   
  zone =  "10-6-37-0"  
  pool = "webshop"
  user_data = data.template_file.userdata.rendered
}

# 1 vm on vm host with highest memory
resource "maas_vm_instance" "order" {
  kvm_no = data.maas_vm_hosts.vm-hosts.recommended
  cpu_count = 2
  memory = 2048
  storage = 12
  hostname = "base-61"
  zone =  "10-6-37-0"
  # pool = "webshop"
  # user_data = data.template_file.userdata.rendered
}