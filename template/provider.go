package yaapigw_tf

//func not_used() {
//	var code string = `
import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/zerolog"

	yc "github.com/villem/yaapigw-go-client"
)

func providerConfigure(ctx context.Context,
	d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	level := zerolog.TraceLevel
	log_level := &level
	if log_level != nil {
		zerolog.SetGlobalLevel(*log_level)
	} else {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}

	username := d.Get("user_name").(string)
	password := d.Get("password").(string)
	user_id := d.Get("user_id").(string)
	server := d.Get("server").(string)
	port := d.Get("port").(int)
	timeout := d.Get("timeout_in_seconds").(int)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	/*diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Warning Message Summary",
		Detail:   "This is the detailed warning message from providerConfigure",
	})*/

	if (username != "") && (password != "") {
		uri := fmt.Sprintf("https://%s:%d", server, port)
		c, err := yc.NewClient(&uri, &username, &password, &user_id, timeout)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create YAAPIGW client connection",
				Detail:   fmt.Sprintf("%s", err),
			})

			return nil, diags

		}

		return c, diags
	}

	c, err := yc.NewClient(nil, nil, nil, nil, timeout)
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
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("YAAPIGW_USER_NAME", nil),
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("YAAPIGW_USER_ID", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("YAAPIGW_PASSWORD", nil),
			},
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("YAAPIGW_SERVER", nil),
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("YAAPIGW_SERVER_PORT", 8848),
			},
			"timeout_in_seconds": {
				Type:      schema.TypeInt,
				Optional:  true,
				Sensitive: true,
				DefaultFunc: schema.EnvDefaultFunc("YAAPIGW_SERVER_TIMEOUT_IN_SECONDS",
					15),
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
