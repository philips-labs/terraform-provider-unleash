package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-labs/go-unleash-api/api"
)

var descriptions map[string]string

func init() {
	schema.DescriptionKind = schema.StringMarkdown

	descriptions = map[string]string{
		"api_url":    "URL of the unleash API",
		"auth_token": "Authentication token to authenticate to the Unleash API",
	}
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_url": {
					Type:        schema.TypeString,
					Required:    true,
					Description: descriptions["api_url"],
					DefaultFunc: schema.EnvDefaultFunc("UNLEASH_API_URL", nil),
				},
				"auth_token": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: descriptions["auth_token"],
					DefaultFunc: schema.EnvDefaultFunc("UNLEASH_AUTH_TOKEN", nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"unleash_feature":      dataSourceFeature(),
				"unleash_project":      dataSourceProject(),
				"unleash_feature_type": dataSourceFeatureType(),
				"unleash_users":        dataSourceUsers(),
				"unleash_user":         dataSourceUser(),
				"unleash_api_tokens":   dataSourceApiTokens(),
				"unleash_api_token":    dataSourceApiToken(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"unleash_feature":             resourceFeature(),
				"unleash_feature_v2":          resourceFeatureV2(),
				"unleash_strategy_assignment": resourceStrategyAssignment(),
				"unleash_feature_enabling":    resourceFeatureEnabling(),
				"unleash_user":                resourceUser(),
				"unleash_api_token":           resourceApiToken(),
				"unleash_project":             resourceProject(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics

		apiClient, err := api.NewClient(&http.Client{}, d.Get("api_url").(string), d.Get("auth_token").(string))
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return apiClient, diags
	}
}
