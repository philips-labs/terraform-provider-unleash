package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-labs/go-unleash-api/api"
)

func dataSourceFeature() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Retrieve details of an existing feature",

		ReadContext: dataSourceFeatureRead,

		// This descriptions are used by the documentation generator and the language server.
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Feature name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project_id": {
				Description: "The project id of the feature toggle",
				Type:        schema.TypeString,
				Required:    true,
			},
			"archived": {
				Description: "Wether the feature toggle is archived or not",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"created_at": {
				Description: "The date the feature toggle was created",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "The description of the feature toggle",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"stale": {
				Description: "Wether the feature toggle is stale or not",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"type": {
				Description: "The type of the feature toggle",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"environments": {
				Description: "The environments of the feature toggle",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceFeatureRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	projectId := d.Get("project_id").(string)

	feature, _, err := client.FeatureToggles.GetFeatureByName(projectId, name)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(feature.Name)
	_ = d.Set("archived", feature.Archived)
	_ = d.Set("created_at", feature.CreatedAt)
	_ = d.Set("description", feature.Description)
	_ = d.Set("name", feature.Name)
	_ = d.Set("project_id", feature.Project)
	_ = d.Set("stale", feature.Stale)
	_ = d.Set("type", feature.Type)

	environments := []interface{}{}

	for _, env := range feature.Environments {
		tfMap := map[string]interface{}{}
		tfMap["name"] = env.Name
		tfMap["enabled"] = env.Enabled
		environments = append(environments, tfMap)
	}
	_ = d.Set("environments", environments)

	return diags
}
