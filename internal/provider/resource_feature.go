package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-labs/go-unleash-api/api"
)

func resourceFeature() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a resource for managing unleash features.",

		CreateContext: resourceFeatureCreate,
		ReadContext:   resourceFeatureRead,
		UpdateContext: resourceFeatureUpdate,
		DeleteContext: resourceFeatureDelete,

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
		},
	}
}

func resourceFeatureCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	feature := &api.FeatureToggle{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Project:     d.Get("project_id").(string),
	}

	// TODO: check if provided project_id and feature type are valid

	createdFeature, resp, err := client.FeatureToggles.CreateFeature(feature.Project, *feature)
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(createdFeature.Name)
	readDiags := resourceFeatureRead(ctx, d, meta)
	if readDiags != nil {
		diags = append(diags, readDiags...)
	}

	return diags
}

func resourceFeatureRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return diags
}

func resourceFeatureUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	feature := &api.FeatureToggle{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Project:     d.Get("project_id").(string),
	}

	// TODO: check if provided project_id and feature type are valid

	_, resp, err := client.FeatureToggles.UpdateFeature(feature.Project, *feature)
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// Archives a feature
func resourceFeatureDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	featureName := d.Id()
	projectId := d.Get("project_id").(string)
	_, _, err := client.FeatureToggles.ArchiveFeature(projectId, featureName)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
