package yaapigw_tf

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/zerolog/log"
	yc "github.com/villem/yaapigw-go-client"
)

func resourceRunTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRunTemplateCreate,
		ReadContext:   resourceRunTemplateRead,
		//UpdateContext: resourceOrderUpdate,
		DeleteContext: resourceRunTemplateDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"inputs": {
				ForceNew: true,
				Required: true,
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"outputs": {
				Computed: true,
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
	/*
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},*/
}

func resourceRunSingleDynamicFWTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRunTemplateCreate,
		ReadContext:   resourceRunTemplateRead,
		//UpdateContext: resourceOrderUpdate,
		DeleteContext: resourceRunTemplateDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "template_single_fw_for_public_cloud",
			},
			"inputs": {
				ForceNew: true,
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_server_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"single_fw__name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"extend_override": {
							Type: schema.TypeMap,
							Description: "A map containing key value that are automatically" +
								"converted as extra inputs. With this one can add new" +
								"key,value pairs that would be otherwise forbidden by schema" +
								"verification",
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"outputs": {
				Computed: true,
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
	/*
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},*/
}

func resourceRunTemplateCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	c := m.(*yc.YaapiGWClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var inputs map[string]interface{}
	var inputs_2 map[string]interface{}

	inputs_in := d.Get("inputs")
	switch input_type := inputs_in.(type) {
	case map[string]interface{}:
		inputs = inputs_in.(map[string]interface{})
	case *schema.Set:
		inputs = make(map[string]interface{})

		for n, v := range input_type.List() {
			switch v_type := v.(type) {
			case map[string]interface{}:
				// Check and converted extend_override key, value pairs input pairs
				inputs_2 = v.(map[string]interface{})
				for n2, v2 := range inputs_2 {
					switch v2_type := v2.(type) {
					case map[string]interface{}:
						inputs_3 := v2.(map[string]interface{})
						for n3, v3 := range inputs_3 {
							inputs[strings.ToUpper(n3)] = v3
						}

					default:
						_ = v2_type // shut up the compiler
						inputs[strings.ToUpper(n2)] = v2
					}
				}
			default:
				log.Trace().Msgf("[TRACE] Default branch v=%v type is %v\n", v, v_type)
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Non supported inputs member data type",
					Detail:   "",
				})
				return diags

			}

			log.Trace().Msgf("[TRACE] TypeSet input type %v, n=%v, v=%#v, inputs %v\n",
				*input_type, n, v, inputs)

		}

	default:
		//log.Printf("[ERROR] unknown input type %v\n", input_type)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Non supported inputs data type",
			Detail:   "",
		})
		return diags
	}

	ois := yc.TemplateRunArguments{}
	for name, item := range inputs {
		log.Trace().Msgf("name %v item i %v\n", name, item)
	}
	ois.Name = d.Get("name").(string)
	ois.Inputs = inputs
	o, err := c.RunTemplate(&ois)
	log.Trace().Msgf("template run returned %v, err %v \n", o, err)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(o.Id)
	d.Set("outputs", o.Outputs)
	//resourceOrderRead(ctx, d, m)
	return diags
}

/*
func flattenCoffee(coffee hc.Coffee) []interface{} {
	c := make(map[string]interface{})
	c["id"] = coffee.ID
	c["name"] = coffee.Name
	c["teaser"] = coffee.Teaser
	c["description"] = coffee.Description
	c["price"] = coffee.Price
	c["image"] = coffee.Image

	return []interface{}{c}
}

func flattenOrderItems(orderItems *[]hc.OrderItem) []interface{} {
	if orderItems != nil {
		ois := make([]interface{}, len(*orderItems), len(*orderItems))

		for i, orderItem := range *orderItems {
			oi := make(map[string]interface{})

			oi["coffee"] = flattenCoffee(orderItem.Coffee)
			oi["quantity"] = orderItem.Quantity
			ois[i] = oi

		}

		return ois
	}

	return make([]interface{}, 0)
}
*/
func resourceRunTemplateRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	//c := m.(*hc.YaapiGWClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	orderID := d.Id()

	/*order, err := c.GetOrder(orderID)
	if err != nil {
		return diag.FromErr(err)
	}
	*/
	//orderItems := flattenOrderItems(&order.Items)
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "UH!",
		Detail:   fmt.Sprintf("OH OH!  %s", orderID),
	})
	//if err := d.Set("items", orderItems); err != nil {
	//	return diag.FromErr(err)
	//}

	return diags
}

/*
func resourceOrderUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	c := m.(*hc.Client)

	orderID := d.Id()

	if d.HasChange("items") {
		items := d.Get("items").([]interface{})
		ois := []hc.OrderItem{}

		for _, item := range items {
			i := item.(map[string]interface{})

			co := i["coffee"].([]interface{})[0]
			coffee := co.(map[string]interface{})

			oi := hc.OrderItem{
				Coffee: hc.Coffee{
					ID: coffee["id"].(int),
				},
				Quantity: i["quantity"].(int),
			}
			ois = append(ois, oi)
		}

		_, err := c.UpdateOrder(orderID, ois)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceOrderRead(ctx, d, m)
}
*/
func resourceRunTemplateDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	//c := m.(*yc.YaapiGWClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//orderID := d.Id()

	/*err := c.DeleteOrder(orderID)
	if err != nil {
		return diag.FromErr(err)
	}
	*/
	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
