package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-labs/go-unleash-api/api"
)

func dataSourceFeatureType() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Retrieve details of an existing feature type",

		ReadContext: dataSourceFeatureTypeRead,

		// This descriptions are used by the documentation generator and the language server.
		Schema: map[string]*schema.Schema{
			"type_id": {
				Description: "The id of the feature type",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Feature type name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"lifetime_days": {
				Description: "The lifetime of the feature type in days",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"description": {
				Description: "The description of the feature type",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceFeatureTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).PhilipsUnleashClient

	var diags diag.Diagnostics

	typeId := d.Get("type_id").(string)

	resp, _, err := client.FeatureTypes.GetAllFeatureTypes()
	if err != nil {
		return diag.FromErr(err)
	}
	types := resp.Types

	var foundFeatureType api.FeatureType
	for _, featureType := range types {
		if featureType.ID == typeId {
			foundFeatureType = featureType
			break
		}
	}

	if foundFeatureType.ID == "" {
		return diag.FromErr(api.ErrNotFound)
	}

	d.SetId(foundFeatureType.ID)
	_ = d.Set("type_id", foundFeatureType.ID)
	_ = d.Set("name", foundFeatureType.Name)
	_ = d.Set("description", foundFeatureType.Description)
	_ = d.Set("lifetime_days", foundFeatureType.LifetimeDays)

	return diags
}
