package yaapigw_template

//func not_used() {
//	var code string = `
import (
	"context"

	yc "example.com/yaapigw_client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func providerConfigure(ctx context.Context,
	d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	/*diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Warning Message Summary",
		Detail:   "This is the detailed warning message from providerConfigure",
	})*/

	if (username != "") && (password != "") {
		c, err := yc.NewClient(nil, &username, &password)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create YAAPIGW client",
				Detail:   "Unable to auth user for authenticated YAAPIGW client",
			})

			return nil, diags

		}

		return c, diags
	}

	c, err := yc.NewClient(nil, nil, nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create YAAPIGW client",
			Detail:   "Unable to auth user for authenticated YAAPIGW client",
		})

		return nil, diags
	}

	return c, diags
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("YAAPIGW_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("YAAPIGW_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"yaapigw_run_template":            resourceRunTemplate(),
			"yaapigw_single_dynamic_firewall": resourceRunSingleDynamicFWTemplate(),
		},
		/*DataSourcesMap: map[string]*schema.Resource{
			"hashicups_coffees": dataSourceCoffees(),
			"hashicups_order":   dataSourceOrder(),
		},*/
		ConfigureContextFunc: providerConfigure,
	}
}

//`
//	fmt.Printf("%s", code)
//}