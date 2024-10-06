package maas

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/canonical/gomaasclient/client"
	"github.com/canonical/gomaasclient/entity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMaasVMInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVMInstanceCreate,
		ReadContext:   resourceInstanceRead,
		DeleteContext: resourceVMInstanceDelete,
		Description:   "Provides a resource to create and deploy VMs in MAAS, based on the specified parameters. If no parameters are given, a random machine will be created and deployed using the defaults.",
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				client := m.(*client.Client)
				machine, err := getMachine(client, d.Id())
				if err != nil {
					return nil, err
				}
				if machine.StatusName != "Deployed" {
					return nil, fmt.Errorf("machine '%s' needs to be already deployed to be imported as maas_instance resource", machine.Hostname)
				}
				d.SetId(machine.SystemID)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"cpu_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     1,
				Description: "The number of cores used to allocate the MAAS machine.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The Description.",
			},
			"ephemeral": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Deploy machine in memory",
			},
			"fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The deployed MAAS machine FQDN.",
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The hostname of the MAAS machine to be allocated.",
			},
			"ip_addresses": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The ip addresses.",
			},
			"kvm_no": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The number of the kvm host.",
			},
			"memory": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     2048,
				Description: "The RAM memory size (in MB) used to allocate the MAAS machine.",
			},
			"pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Default:     nil,
				Description: "The resource pool name of the MAAS machine to be allocated.",
			},
			"storage": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     8,
				Description: "The VM local storage (in GBytes).",
			},
			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A set of tag names that must be assigned on the MAAS machine to be allocated.",
			},
			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Cloud-init user data script that gets run on the machine once it has deployed. A good practice is to set this with `file.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Default:     nil,
				Description: "The zone name of the MAAS machine to be allocated.",
			},
		},
	}
}

func resourceVMInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*client.Client)

	no := d.Get("kvm_no").(int)
	params := &entity.VMHostMachineParams{
		Cores:    d.Get("cpu_count").(int),
		Memory:   int64(d.Get("memory").(int)),
		Storage:  strconv.Itoa(d.Get("storage").(int)),
		Hostname: d.Get("hostname").(string),
	}
	machine, err := client.VMHost.Compose(no, params)
	if err != nil {
		return diag.FromErr(err)
	}

	// Save system id
	d.SetId(machine.SystemID)

	// Wait for MAAS machine to be ready
	timeout := 10 * time.Minute
	_, err = waitForMachineStatus(ctx, client, machine.SystemID, []string{"Commissioning", "Testing"}, []string{"Ready"}, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update pool and zone
	mparams := &entity.MachineParams{
		Pool:        d.Get("pool").(string),
		Zone:        d.Get("zone").(string),
		Description: d.Get("description").(string),
	}
	mpower := make(map[string]interface{})
	mpower[machine.PowerType] = "virsh"

	_, err = client.Machine.Update(machine.SystemID, mparams, mpower)
	if err != nil {
		return diag.FromErr(err)
	}

	userdata := d.Get("user_data").(string)
	deploy := &entity.MachineDeployParams{
		UserData: userdata,
	}

	// Deploy
	_, err = client.Machine.Deploy(machine.SystemID, deploy)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for MAAS machine to be deployed
	timeout = 20 * time.Minute
	_, err = waitForMachineStatus(ctx, client, machine.SystemID, []string{"Deploying"}, []string{"Deployed"}, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Read MAAS machine info
	return resourceInstanceRead(ctx, d, meta)

}

func resourceVMInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	// Delete VM MAAS machine
	err := client.Machine.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
