package maas

import (
	"context"

	"github.com/canonical/gomaasclient/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMaasVMHost() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVMHostRead,
		Description: "Provides details about an existing MAAS VM hosts.",

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The VM host ID.",
			},
			"no": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The VM host internal ID (for create VM Instances).",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The new VM host name. This is computed if it's not set.",
			},
			"cores": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The VM host total number of CPU cores.",
			},
			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The VM host total RAM memory (in MB).",
			},
			"local_storage": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The VM host total local storage (in bytes).",
			},
		},
	}
}

func dataSourceVMHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	vmHost, err := getVMHost(m.(*client.Client), d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	tfState := map[string]interface{}{
		"id":            vmHost.Host.SystemID,
		"no":            vmHost.ID,
		"name":          vmHost.Name,
		"cores":         vmHost.Available.Cores,
		"memory":        vmHost.Available.Memory,
		"local_storage": vmHost.Available.LocalStorage,
	}
	if err := setTerraformState(d, tfState); err != nil {
		return diag.FromErr(err)
	}
	return nil

}
