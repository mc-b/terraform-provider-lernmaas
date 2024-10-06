package maas

import (
	"context"

	"github.com/canonical/gomaasclient/client"
	"github.com/canonical/gomaasclient/entity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMaasMachines() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMachinesRead,
		Description: "Provides details about all existing MAAS Machines.",

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Internal ID.",
			},
			"system_id": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The VM host IDs.",
			},
			"hostname": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The VM host names.",
			},
			"zone": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The deployed MAAS machine zone names.",
			},
			"pool": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The deployed MAAS machine pool names.",
			},
			"power_type": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A power management type (e.g. `ipmi`).",
			},
			"architecture": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The architecture type of the machine. Defaults to `amd64/generic`.",
			},
			"description": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The descriptions.",
			},
		},
	}
}

func dataSourceMachinesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*client.Client)
	machines, err := client.Machines.Get(&entity.MachinesParams{})
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	system_id := make([]string, len(machines))
	names := make([]string, len(machines))
	zones := make([]string, len(machines))
	pools := make([]string, len(machines))
	power_types := make([]string, len(machines))
	architectures := make([]string, len(machines))
	descriptions := make([]string, len(machines))

	for i, machine := range machines {
		system_id[i] = machine.SystemID
		names[i] = machine.Hostname
		zones[i] = machine.Zone.Name
		pools[i] = machine.Pool.Name
		power_types[i] = machine.PowerType
		architectures[i] = machine.Architecture
		descriptions[i] = machine.Description
	}

	tfState := map[string]interface{}{
		"id":           id,
		"system_id":    system_id,
		"hostname":     names,
		"zone":         zones,
		"pool":         pools,
		"power_type":   power_types,
		"architecture": architectures,
		"description":  descriptions,
	}
	if err := setTerraformState(d, tfState); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
