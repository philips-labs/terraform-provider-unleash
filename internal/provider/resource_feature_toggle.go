package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/philips-labs/go-unleash-api/api"
)

func resourceFeatureToggle() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a resource for managing unleash features with variants and environment strategies.",

		CreateContext: resourceFeatureToggleCreate,
		ReadContext:   resourceFeatureToggleRead,
		UpdateContext: resourceFeatureToggleUpdate,
		DeleteContext: resourceFeatureToggleDelete,

		// The descriptions are used by the documentation generator and the language server.
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Feature name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"project_id": {
				Description: "The feature will be created in the given project",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Description: "Feature type",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Feature description",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"archive_on_destroy": {
				Description: "Whether to archive the feature toggle on destroy. Default is `true`. When `false`, it will permanently delete the feature toggle.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"environment": {
				Description: "Use this to enable a feature in an environment and add strategies",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Environment name",
							Type:        schema.TypeString,
							Required:    true,
						},
						"enabled": {
							Description: "Whether the feature is on/off in the environment",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
						},
						"strategy": {
							Description: "Strategy to add in the environment",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Description: "Strategy unique name",
										Type:        schema.TypeString,
										Required:    true,
									},
									"parameters": {
										Description: "Strategy parameters. All the values need to informed as strings.",
										Type:        schema.TypeMap,
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"id": {
										Description: "Strategy ID",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"variant": {
				Description: "Feature variant",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Variant name",
							Type:        schema.TypeString,
							Required:    true,
						},
						"stickiness": {
							Description: "Variant stickiness. Default is `default`.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "default",
						},
						"weight": {
							Description:  "Variant weight. Only considered when the `weight_type` is `fix`. It is calculated automatically if the `weight_type` is `variable`.",
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"weight_type": {
							Description: "Variant weight type. The weight type can be `fix` or `variable`. Default is `variable`.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "variable",
						},
						"payload": {
							Description: "Variant payload. The type of the payload can be `string`, `json` or `csv`",
							Type:        schema.TypeSet,
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"overrides": {
							Description: "Overrides existing context field values. Values are comma separated e.g `v1, v2, ...`)",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"context_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceFeatureToggleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	feature := &api.FeatureToggle{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Project:     d.Get("project_id").(string),
	}

	createdFeature, resp, err := client.FeatureToggles.CreateFeature(feature.Project, *feature)
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	if v, ok := d.GetOk("variant"); ok {
		tfVariants := v.([]interface{})
		variants := make([]api.Variant, 0, len(tfVariants))
		for _, tfVariant := range tfVariants {
			variants = append(variants, toFeatureVariant(tfVariant.(map[string]interface{})))
		}
		_, resp, err := client.Variants.AddVariantsForFeatureToggle(feature.Project, feature.Name, variants)
		if resp == nil {
			return diag.FromErr(fmt.Errorf("response is nil: %v", err))
		}
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if e, ok := d.GetOk("environment"); ok {
		tfEnvironments := e.([]interface{})
		for _, tfEnvironment := range tfEnvironments {
			environment := toFeatureEnvironment(tfEnvironment.(map[string]interface{}))

			for _, strategy := range environment.Strategies {
				_, resp, err := client.FeatureToggles.AddStrategyToFeature(feature.Project, feature.Name, environment.Name, strategy)
				if resp == nil {
					return diag.FromErr(fmt.Errorf("response is nil: %v", err))
				}
				if err != nil {
					return diag.FromErr(err)
				}

			}
			ok, _, err := client.FeatureToggles.EnableFeatureOnEnvironment(feature.Project, feature.Name, environment.Name, environment.Enabled)
			if err != nil || !ok {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(createdFeature.Name)
	readDiags := resourceFeatureToggleRead(ctx, d, meta)
	if readDiags != nil {
		diags = append(diags, readDiags...)
	}

	return diags
}

func resourceFeatureToggleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	featureName := d.Id()
	projectId := d.Get("project_id").(string)
	feature, _, err := client.FeatureToggles.GetFeatureByName(projectId, featureName)
	if err != nil {
		if err == api.ErrNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("name", feature.Name)
	_ = d.Set("description", feature.Description)
	_ = d.Set("type", feature.Type)
	_ = d.Set("project_id", feature.Project)
	_ = d.Set("variant", flattenVariants(feature.Variants))

	if e, ok := d.GetOk("environment"); ok {
		toSave := []api.Environment{}

		for _, env := range feature.Environments {
			for _, tfEnvironment := range e.([]interface{}) {
				if tfEnvironment.(map[string]interface{})["name"] == env.Name { // the api returns all envs, so we only add to the state whats defined in TF
					toSave = append(toSave, env)
				}
			}
		}

		_ = d.Set("environment", flattenEnvironments(toSave))
	}

	return diags
}

func resourceFeatureToggleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	feature := &api.FeatureToggle{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Project:     d.Get("project_id").(string),
	}

	_, resp, err := client.FeatureToggles.UpdateFeature(feature.Project, *feature)
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("variant") {
		tfVariants := d.Get("variant").([]interface{})
		variants := make([]api.Variant, 0, len(tfVariants))
		for _, tfVariant := range tfVariants {
			variants = append(variants, toFeatureVariant(tfVariant.(map[string]interface{})))
		}
		_, resp, err := client.Variants.AddVariantsForFeatureToggle(feature.Project, feature.Name, variants)
		if resp == nil {
			return diag.FromErr(fmt.Errorf("response is nil: %v", err))
		}
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// if d.HasChange("environment") {
	// 	o, n := d.GetChange("environment")
	// 	old := o.([]interface{})
	// 	new := n.([]interface{})

	// 	tfEnvironments := d.Get("environment").([]interface{})
	// 	for _, tfEnvironment := range tfEnvironments {
	// 		environment := toFeatureEnvironment(tfEnvironment.(map[string]interface{}))

	// 		for _, strategy := range environment.Strategies {
	// 			_, resp, err := client.FeatureToggles.AddStrategyToFeature(feature.Project, feature.Name, environment.Name, strategy)
	// 			if resp == nil {
	// 				return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	// 			}
	// 			if err != nil {
	// 				return diag.FromErr(err)
	// 			}
	// 		}
	// 	}
	// }

	return diags
}

// Archives a feature
func resourceFeatureToggleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	featureName := d.Id()
	projectId := d.Get("project_id").(string)
	_, _, err := client.FeatureToggles.ArchiveFeature(projectId, featureName)
	if err != nil {
		return diag.FromErr(err)
	}
	shouldArchive := d.Get("archive_on_destroy").(bool)
	if !shouldArchive {
		_, _, err := client.FeatureToggles.DeleteArchivedFeature(featureName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return diags
}

func toFeatureEnvironment(tfEnvironment map[string]interface{}) api.Environment {
	environment := api.Environment{}
	environment.Name = tfEnvironment["name"].(string)
	environment.Enabled = tfEnvironment["enabled"].(bool)

	if tfStrategies, ok := tfEnvironment["strategy"].(*schema.Set); ok && tfStrategies.Len() > 0 {
		strategiesList := tfStrategies.List()
		strategies := make([]api.FeatureStrategy, 0, len(strategiesList))
		for _, tfStrategy := range strategiesList {
			strategyMap := tfStrategy.(map[string]interface{})
			strategy := api.FeatureStrategy{
				Name: strategyMap["name"].(string),
			}
			if p, ok := strategyMap["parameters"]; ok {
				tfParams := p.(map[string]interface{})
				castedParameters := make(map[string]interface{})
				for k, v := range tfParams {
					castedParameters[k] = v.(string)
				}
				strategy.Parameters = castedParameters
			}
			strategies = append(strategies, strategy)
		}
		environment.Strategies = strategies
	}
	return environment
}

func flattenEnvironments(environments []api.Environment) []interface{} {
	if environments == nil {
		return []interface{}{}
	}

	tfEnvironments := []interface{}{}

	for _, env := range environments {
		tfEnvironment := map[string]interface{}{}
		tfEnvironment["name"] = env.Name
		tfEnvironment["enabled"] = env.Enabled

		if env.Strategies != nil {
			tfStrategies := []interface{}{}
			for _, strategy := range env.Strategies {
				tfStrategy := map[string]interface{}{}
				tfStrategy["id"] = strategy.ID
				tfStrategy["name"] = strategy.Name
				retrievedParams := strategy.Parameters.(map[string]interface{})
				castedParams := make(map[string]interface{})
				for k, v := range retrievedParams {
					castedParams[k] = v.(string)
				}
				tfStrategy["parameters"] = castedParams
				tfStrategies = append(tfStrategies, tfStrategy)
			}
			tfEnvironment["strategy"] = tfStrategies
		}

		tfEnvironments = append(tfEnvironments, tfEnvironment)
	}

	return tfEnvironments
}
