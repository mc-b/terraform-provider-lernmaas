package maas

import (
	"context"

	"github.com/canonical/gomaasclient/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMaasVMHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVMHostsRead,
		Description: "Provides details about all existing MAAS VM hosts.",

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Internal ID.",
			},
			"name": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The VM host names.",
			},
			"no": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "The VM host internal IDs (for create VM Instances).",
			},
			"recommended": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The VM host internal ID with the most free memory (for create VM Instances).",
			},
			"system_id": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The VM host IDs.",
			},
		},
	}
}

func dataSourceVMHostsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*client.Client)
	vmHosts, err := client.VMHosts.Get()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	numbers := make([]int, len(vmHosts))
	system_id := make([]string, len(vmHosts))
	names := make([]string, len(vmHosts))

	recommended := 0
	memory := int64(0)
	for i, vmHost := range vmHosts {
		numbers[i] = vmHost.ID
		system_id[i] = vmHost.Host.SystemID
		names[i] = vmHost.Name
		if vmHost.Available.Memory > memory {
			recommended = vmHost.ID
			memory = vmHost.Available.Memory
		}
	}

	tfState := map[string]interface{}{
		"id":          id,
		"system_id":   system_id,
		"no":          numbers,
		"name":        names,
		"recommended": recommended,
	}
	if err := setTerraformState(d, tfState); err != nil {
		return diag.FromErr(err)
	}
	return nil

}
