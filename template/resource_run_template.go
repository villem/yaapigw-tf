package yaapigw_tf

import (
	"context"
	"net/http"
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
		UpdateContext: resourceTemplateUpdate,
		DeleteContext: resourceRunTemplateDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"manager_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"inputs": {
				ForceNew: false,
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
			"digest": {
				Computed: true,
				Type:     schema.TypeString,
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
		UpdateContext: resourceTemplateUpdate,
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
			"manager_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inputs": {
				ForceNew: false,
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
			"digest": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
	/*
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},*/
}

func convert_to_yaapigw_inputs(inputs_in interface{}) (
	*map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	var inputs map[string]interface{}
	var inputs_2 map[string]interface{}

	log.Trace().Msgf("[TRACE] inputs %v", inputs_in)

	switch input_type := inputs_in.(type) {

	case map[string]interface{}:
		log.Trace().Msgf("[TRACE] map.string branch")
		inputs = inputs_in.(map[string]interface{})

	case *schema.Set:
		log.Trace().Msgf("[TRACE] Set branch %v list=>>%v<<", input_type,
			input_type.List())
		inputs = make(map[string]interface{})

		for n, v := range input_type.List() {

			switch v_type := v.(type) {

			case map[string]interface{}:
				log.Trace().Msgf("[TRACE] 2 map string branch %v", v)
				// Check and converted extend_override key, value pairs input pairs
				inputs_2 = v.(map[string]interface{})
				log.Trace().Msgf("[TRACE] 2 map string branch %v ", v)

				all_empty := true
				// This is stupid, but GetChanges() gives data that gives values for
				// keys. This messes everything so let's prune it out
				for _, v2 := range inputs_2 {
					log.Trace().Msgf("[TRACE] empty map check v2=%v", v2)
					switch v2_type := v2.(type) {
					case map[string]interface{}:
						if len(v2.(map[string]interface{})) > 0 {
							all_empty = false
						}
					case string:
						if v2 != "" {
							all_empty = false
						}
					default:
						_ = v2_type
					}
				}
				if all_empty {
					log.Trace().Msgf("[TRACE] empty map..")
					break
				}

				for n2, v2 := range inputs_2 {

					switch v2_type := v2.(type) {

					case map[string]interface{}:
						log.Trace().Msgf("[TRACE] 3 map string branch %v", v2)

						inputs_3 := v2.(map[string]interface{})

						for n3, v3 := range inputs_3 {
							inputs[strings.ToUpper(n3)] = v3
						}

					default:
						log.Trace().Msgf("[TRACE] 3 default (string) branch %v", v2)

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
				return nil, diags

			}

			log.Trace().
				Msgf("[TRACE] Type Set input type %v, n=%v, v=%#v\n",
					*input_type, n, v)

		}

	default:
		log.Error().Msgf("[ERROR] unknown input type %v\n", input_type)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Non supported inputs data type",
			Detail:   "",
		})
		return nil, diags
	}

	return &inputs, nil
}

func RunTemplate(ctx context.Context, d *schema.ResourceData,
	m interface{}, inputs *map[string]interface{}) (
	*yc.TemplateResults, error) {

	c := m.(*yc.YaapiGWClient)
	ois := yc.TemplateRunArguments{}
	for name, item := range *inputs {
		log.Trace().Msgf("name %v item i %v\n", name, item)
	}
	ois.Name = d.Get("name").(string)
	ois.Inputs = *inputs

	log.Trace().Msgf("operating mode is %v\n", c.OperatingMode)
	if c.OperatingMode != "" {
		outputs := d.Get("outputs").(map[string]interface{})
		if val, ok := outputs["manager_id"]; ok {
			manager := val.(string)
			ois.Manger = &manager
		}
	}

	m_id := d.Get("manager_id")
	var m_id2 string
	if m_id != nil {
		m_id2 = m_id.(string)
		ois.Manger = &m_id2
	}
	log.Trace().Msgf("ois %#v m_id %v mid2 %#v", ois, m_id, m_id2)

	o, err := c.RunTemplate(&ois)
	log.Trace().Msgf("template run returned %v, err %v \n", o, err)
	return o, err
}

func resourceRunTemplateCreate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	log.Trace().Msgf("[TRACE] CREATE: START)")

	c := m.(*yc.YaapiGWClient)

	var diags diag.Diagnostics

	inputs_in := d.Get("inputs")
	inputs, err2 := convert_to_yaapigw_inputs(inputs_in)
	if err2 != nil {
		return err2
	}

	o, err := RunTemplate(ctx, d, m, inputs)

	if err != nil {
		return diag.FromErr(err)
	}

	if c.OperatingMode == "" {
		d.SetId(o.Id)
	}
	d.Set("outputs", o.Outputs)
	d.Set("digest", o.Outputs["digest"])
	return diags
}

func resourceRunTemplateRead(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	c := m.(*yc.YaapiGWClient)

	log.Info().Msgf("[INFO] READ: START)")
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	var outputs map[string]interface{}

	inputs, err2 := convert_to_yaapigw_inputs(d.Get("inputs"))
	if err2 != nil {
		return err2
	}

	template_args := yc.TemplateRunArguments{}

	template_args.Name = d.Get("name").(string)
	template_args.Inputs = *inputs

	outputs = d.Get("outputs").(map[string]interface{})
	log.Trace().Msgf("read outputs %v\n", outputs)
	log.Trace().Msgf("read inputs %v\n", *inputs)

	digest := d.Get("digest")
	if digest != "" {
		digest := digest.(string)
		template_args.Digest = &digest
		c.OperatingMode = "check"

	}

	if outputs != nil {
		if manager, ok := outputs["manager_id"]; ok {
			manager := manager.(string)
			template_args.Manger = &manager
		}

	}

	run_results, err := c.RunTemplate(&template_args)
	c.OperatingMode = ""
	log.Trace().Msgf("template run returned %v, err %v \n", run_results, err)

	if err != nil && run_results.ErrorCode == http.StatusUnprocessableEntity {
		log.Info().Msgf("[INFO] READ: Digest changed %v -> %v",
			digest, run_results.Outputs["digest"])

		d.Set("digest", run_results.Outputs["digest"])
		return diags
	} else if err != nil {
		return diag.FromErr(err)
	}
	if run_results.Outputs["initial_configuration"] == "" && outputs != nil &&
		outputs["initial_configuration"] != nil {
		run_results.Outputs["initial_configuration"] =
			outputs["initial_configuration"].(string)
	}

	d.Set("outputs", run_results.Outputs)
	d.Set("digest", run_results.Outputs["digest"])

	/* diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Template run results Ok!",
		Detail:   fmt.Sprintf("OH OH!  %v", run_results),
	}) */

	return diags
}

func resourceTemplateUpdate(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	c := m.(*yc.YaapiGWClient)
	var diags diag.Diagnostics

	log.Info().Msgf("[INFO] UPDATE: START, inputs change=%v digest change=%v",
		d.HasChange("inputs"), d.HasChange("digest"))

	changes := false
	if d.HasChange("inputs") {
		changes = true
		changes1, changes2 := d.GetChange("inputs")
		inputs, _ := convert_to_yaapigw_inputs(changes2)
		inputs_old, _ := convert_to_yaapigw_inputs(changes1)
		previous_ic :=
			d.Get("outputs").(map[string]interface{})["initial_configuration"]
		log.Trace().Msgf("[TRACE] UPDATE: %v",
			d.Get("outputs").(map[string]interface{}))
		log.Trace().Msgf("[TRACE] UPDATE: old=%#v new=%#v, changes1=%#v, "+
			"changes2=%#v",
			inputs, inputs_old, changes1, changes2)

		updated_inputs := make(map[string]interface{})
		for k, v := range *inputs {
			log.Trace().Msgf("[TRACE] UPDATE: old %v new %v",
				(*inputs_old)[k], v)
			if (*inputs_old)[k] != v {
				log.Trace().Msgf("[TRACE] UPDATE: read old %v is different from  %v",
					(*inputs_old)[k], v)
				new_key := k + "--UPDATE"
				if (*inputs_old)[k] != "" {
					updated_inputs[new_key] = (*inputs_old)[k]
				}
			}
			if v != "" {
				updated_inputs[k] = v
			}
		}

		log.Trace().Msgf("Updated_inputs: %#v", updated_inputs)
		c.OperatingMode = "update"
		o, err := RunTemplate(ctx, d, m, &updated_inputs)
		c.OperatingMode = ""
		if err != nil {
			return diag.FromErr(err)
		}
		if changes {
			o.Outputs["initial_configuration"] = previous_ic.(string)
			d.Set("outputs", o.Outputs)
			d.Set("digest", o.Outputs["digest"])
		}
	}

	return diags
}

func resourceRunTemplateDelete(ctx context.Context, d *schema.ResourceData,
	m interface{}) diag.Diagnostics {
	c := m.(*yc.YaapiGWClient)

	log.Trace().Msgf("[TRACE] DELETE: START)")
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c.OperatingMode = "delete"
	err := resourceRunTemplateCreate(ctx, d, m)
	c.OperatingMode = ""

	if err != nil {
		return err
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
