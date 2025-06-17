package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/philips-labs/go-unleash-api/api"
)

func resourceStrategyAssignment() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a resource for add strategy to a feature toggle in the given environment.",

		CreateContext: resourceStrategyAssignmentCreate,
		ReadContext:   resourceStrategyAssignmentRead,
		UpdateContext: resourceStrategyAssignmentUpdate,
		DeleteContext: resourceStrategyAssignmentDelete,

		// The descriptions are used by the documentation generator and the language server.
		Schema: map[string]*schema.Schema{
			"feature_name": {
				Description: "Feature name to assign the strategy to",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"project_id": {
				Description: "The unleash project the feature is in",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"environment": {
				Description: "The environment where the strategy will take place",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"strategy_name": {
				Description: "Strategy unique name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"strategy_id": {
				Description: "Strategy id",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"parameters": {
				Description: "Strategy parameters. All the values need to informed as strings.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
							Description: "Variant payload. The type of the payload can be `string`, `json` or `csv` or `number`",
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
										Type:        schema.TypeString,
										Description: "Always a string value, independent of the type.",
										Required:    true,
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

func resourceStrategyAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).PhilipsUnleashClient

	var diags diag.Diagnostics

	featureStrategy := &api.FeatureStrategy{
		Name: d.Get("strategy_name").(string),
	}
	projectId := d.Get("project_id").(string)
	featureName := d.Get("feature_name").(string)
	environment := d.Get("environment").(string)

	if p, ok := d.GetOk("parameters"); ok {
		givenParams := p.(map[string]interface{})
		strategy, _, err := client.Strategies.GetStrategyByName(featureStrategy.Name)
		if err != nil {
			return diag.FromErr(err)
		}
		if strategy == nil {
			return diag.FromErr(api.ErrNotFound)
		}

		convertedParams := make(map[string]interface{})
		for _, param := range strategy.Parameters {
			if _, ok := givenParams[param.Name]; !ok && param.Required {
				return diag.FromErr(ErrStrategyParametersRequired)
			}
			convertedParams[param.Name] = givenParams[param.Name].(string)
		}

		featureStrategy.Parameters = convertedParams
	}

	if p, ok := d.GetOk("variant"); ok {
		tfVariants := p.([]interface{})
		variants := make([]api.Variant, 0, len(tfVariants))
		for _, tfVariant := range tfVariants {
			variants = append(variants, toFeatureVariant(tfVariant.(map[string]interface{})))
		}
		featureStrategy.Variants = variants
	}

	addedStrategy, resp, err := client.FeatureToggles.AddStrategyToFeature(projectId, featureName, environment, *featureStrategy)
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(addedStrategy.ID)
	readDiags := resourceStrategyAssignmentRead(ctx, d, meta)
	if readDiags != nil {
		diags = append(diags, readDiags...)
	}

	return diags
}

func resourceStrategyAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).PhilipsUnleashClient

	var diags diag.Diagnostics

	projectId := d.Get("project_id").(string)
	featureName := d.Get("feature_name").(string)
	variants := d.Get("variant").([]interface{})

	feature, _, err := client.FeatureToggles.GetFeatureByName(projectId, featureName)
	if err != nil {
		if err == api.ErrNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	strategyName := d.Get("strategy_name").(string)
	environment := d.Get("environment").(string)

	strategy, _, err := client.Strategies.GetStrategyByName(strategyName)
	if err != nil {
		return diag.FromErr(err)
	}
	if strategy == nil {
		return diag.FromErr(api.ErrNotFound)
	}

	for _, env := range feature.Environments {
		if env.Name == environment {
			for _, featureStrategy := range env.Strategies {
				if featureStrategy.Name == strategyName {
					_ = d.Set("strategy_name", featureStrategy.Name)
					convertedParams := make(map[string]interface{}) // convert actual param type back to string
					retrievedParams := featureStrategy.Parameters
					for _, param := range strategy.Parameters {
						convertedParams[param.Name] = retrievedParams.(map[string]interface{})[param.Name].(string)
					}
					_ = d.Set("parameters", convertedParams)
					_ = d.Set("variant", flattenVariants(variants, featureStrategy.Variants))
				}
			}
			break
		}
	}

	return diags
}

func resourceStrategyAssignmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).PhilipsUnleashClient

	var diags diag.Diagnostics

	strategy := &api.FeatureStrategy{
		ID:   d.Id(),
		Name: d.Get("strategy_name").(string),
	}
	projectId := d.Get("project_id").(string)
	featureName := d.Get("feature_name").(string)
	environment := d.Get("environment").(string)

	if v, ok := d.GetOk("parameters"); ok {
		vv := v.(map[string]interface{})
		found, _, err := client.Strategies.GetStrategyByName(strategy.Name)
		if err != nil {
			return diag.FromErr(err)
		}
		if found == nil {
			return diag.FromErr(api.ErrNotFound)
		}

		convertedParams := make(map[string]interface{})
		for _, param := range found.Parameters {
			if _, ok := vv[param.Name]; !ok && param.Required {
				return diag.FromErr(ErrStrategyParametersRequired)
			}
			convertedParams[param.Name] = vv[param.Name].(string)
		}

		strategy.Parameters = convertedParams
	}

	if d.HasChange("variant") {
		tfVariants := d.Get("variant").([]interface{})
		variants := make([]api.Variant, 0, len(tfVariants))
		for _, tfVariant := range tfVariants {
			variants = append(variants, toFeatureVariant(tfVariant.(map[string]interface{})))
		}
		strategy.Variants = variants
	}

	_, resp, err := client.FeatureToggles.UpdateFeatureStrategy(projectId, featureName, environment, *strategy)
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceStrategyAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).PhilipsUnleashClient

	var diags diag.Diagnostics

	strategyId := d.Id()
	featureName := d.Get("feature_name").(string)
	projectId := d.Get("project_id").(string)
	environment := d.Get("environment").(string)
	_, _, err := client.FeatureToggles.DeleteStrategyFromFeature(projectId, featureName, environment, strategyId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func toFeatureVariant(tfVariant map[string]interface{}) api.Variant {
	variant := api.Variant{}
	variant.Name = tfVariant["name"].(string)
	variant.Stickiness = tfVariant["stickiness"].(string)
	variant.Weight = tfVariant["weight"].(int)
	variant.WeightType = tfVariant["weight_type"].(string)

	if payloadSet, ok := tfVariant["payload"].(*schema.Set); ok && payloadSet.Len() > 0 {
		payloadList := payloadSet.List()
		payloadMap := payloadList[0].(map[string]interface{})
		variant.Payload = &api.VariantPayload{
			Type:  payloadMap["type"].(string),
			Value: payloadMap["value"].(string),
		}
	}
	return variant
}

// func flattenVariants(tfVariants []interface{}, variantsFromApi []api.Variant) []interface{} {
// 	if variantsFromApi == nil {
// 		return []interface{}{}
// 	}

// 	vVariants := []interface{}{}

// 	for _, variant := range variantsFromApi {
// 		mVariant := map[string]interface{}{}
// 		mVariant["name"] = variant.Name
// 		mVariant["weight"] = variant.Weight
// 		mVariant["weight_type"] = variant.WeightType
// 		mVariant["stickiness"] = variant.Stickiness

// 		if variant.Payload != nil {
// 			mPayloads := []interface{}{}
// 			mPayload := map[string]interface{}{}
// 			mPayload["type"] = variant.Payload.Type
// 			mPayload["value"] = variant.Payload.Value
// 			mPayloads = append(mPayloads, mPayload)
// 			mVariant["payload"] = mPayloads
// 		}

// 		vVariants = append(vVariants, mVariant)
// 	}

// 	return vVariants
// }

func flattenVariants(tfVariants []interface{}, variantsFromApi []api.Variant) []interface{} {
	apiVariantMap := make(map[string]api.Variant)
	for _, v := range variantsFromApi {
		apiVariantMap[v.Name] = v
	}

	result := make([]interface{}, 0, len(tfVariants))
	for _, tfV := range tfVariants {
		tfVariant := tfV.(map[string]interface{})
		name := tfVariant["name"].(string)

		if apiVariant, exists := apiVariantMap[name]; exists {
			flattened := map[string]interface{}{
				"name":        apiVariant.Name,
				"stickiness":  apiVariant.Stickiness,
				"weight":      apiVariant.Weight,
				"weight_type": apiVariant.WeightType,
			}

			if apiVariant.Payload != nil {
				payload := map[string]interface{}{
					"type":  apiVariant.Payload.Type,
					"value": apiVariant.Payload.Value,
				}
				flattened["payload"] = []interface{}{payload}
			}

			result = append(result, flattened)
		}
	}

	return result
}
