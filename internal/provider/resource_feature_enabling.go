package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-labs/go-unleash-api/api"
)

func resourceFeatureEnabling() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a resource for enabling a feature toggle in the given environment. This can be only done after the feature toggle has at least one strategy.",

		CreateContext: resourceFeatureEnablingCreate,
		ReadContext:   resourceFeatureEnablingRead,
		UpdateContext: resourceFeatureEnablingUpdate,
		DeleteContext: resourceFeatureEnablingDelete,

		// The descriptions are used by the documentation generator and the language server.
		Schema: map[string]*schema.Schema{
			"feature_name": {
				Description: "Feature name to enabled",
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
				Description: "The environment where the toggle will be enabled",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"enabled": {
				Description: "Whether the feature is on/off in the provided environment",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceFeatureEnablingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	projectId := d.Get("project_id").(string)
	featureName := d.Get("feature_name").(string)
	environment := d.Get("environment").(string)
	enabled := d.Get("enabled").(bool)

	ok, _, err := client.FeatureToggles.EnableFeatureOnEnvironment(projectId, featureName, environment, enabled)
	if err != nil || !ok {
		return diag.FromErr(err)
	}

	d.SetId(featureName + "/" + environment)
	readDiags := resourceFeatureEnablingRead(ctx, d, meta)
	if readDiags != nil {
		diags = append(diags, readDiags...)
	}

	return diags
}

func resourceFeatureEnablingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

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

	environment := d.Get("environment").(string)

	for _, env := range feature.Environments {
		if env.Name == environment {
			_ = d.Set("enabled", env.Enabled)
			break
		}
	}

	return diags
}

func resourceFeatureEnablingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	projectId := d.Get("project_id").(string)
	featureName := d.Get("feature_name").(string)
	environment := d.Get("environment").(string)
	enabled := d.Get("enabled").(bool)

	ok, _, err := client.FeatureToggles.EnableFeatureOnEnvironment(projectId, featureName, environment, enabled)
	if err != nil || !ok {
		return diag.FromErr(err)
	}

	return diags
}

func resourceFeatureEnablingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	projectId := d.Get("project_id").(string)
	featureName := d.Get("feature_name").(string)
	environment := d.Get("environment").(string)

	_, _, err := client.FeatureToggles.EnableFeatureOnEnvironment(projectId, featureName, environment, false)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
