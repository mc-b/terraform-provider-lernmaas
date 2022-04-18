terraform {
  required_providers {
    maas = {
      source = "mc-b/lernmaas"
      version = "0.0.3"
    }
  }
}

provider "maas" {
  # Configuration options
  api_version = "2.0"
  api_key = "nC83nVyLDWKF8zxvPq:VnYukvR9Yh3w9jRDUe:FdU7sFWJu8DHjAeRumRNmZfasNBDqgXa"
  api_url = "http://10.6.37.8:5240/MAAS"
}

data "maas_vm_host" "kvm-01" {
  name = "cloud-hf-39"
}

data "maas_vm_host" "kvm-02" {
  name = "cloud-hf-20"
}

data "maas_vm_hosts" "vm-hosts" {
  id  = "rack-01"
}

data "maas_machines" "machines" {
  id  = "rack-01"
}


resource "maas_vm_instance" "base" {
  count = length(data.maas_vm_hosts.vm-hosts.no) * 2
  kvm_no = data.maas_vm_hosts.vm-hosts.no[count.index % length(data.maas_vm_hosts.vm-hosts.no)]
  hostname = "base-${format("%02d", count.index + 10)}"
  description = "student ${format("%02d", count.index + 1)}"   
  zone =  "10-6-37-0"  
  pool = "webshop"
  user_data = data.template_file.userdata.rendered
}




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




