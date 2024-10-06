---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "maas_network_interface_vlan Resource - terraform-provider-maas"
subcategory: ""
description: |-
  Provides a resource to manage MAAS network Vlans.
---

# maas_network_interface_vlan (Resource)

Provides a resource to manage MAAS network Vlans.

## Example Usage

```terraform
resource "maas_network_interface_vlan" "example" {
  machine = maas_machine.example.id
  parent  = maas_network_interface_bond.example.name
  vlan    = data.maas_vlan.example.id
  fabric  = "fabric"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `fabric` (String) The identifier (name or ID) of the fabric for the new VLAN interface.
- `machine` (String) The identifier (system ID, hostname, or FQDN) of the machine with the VLAN interface.
- `parent` (String) Parent interface name for this VLAN interface.

### Optional

- `accept_ra` (Boolean) Accept router advertisements. (IPv6 only).
- `mtu` (Number) The MTU of the VLAN interface.
- `tags` (Set of String) A set of tag names to be assigned to the VLAN interface.
- `vlan` (Number) Database ID of the VLAN the VLAN interface is connected to.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
# A VLAN network interface can be imported using the machine identifier (system ID, hostname, or FQDN) and ID. e.g.
$ terraform import maas_network_interface_vlan.example vm1:id
```