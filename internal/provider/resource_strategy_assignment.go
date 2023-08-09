package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
