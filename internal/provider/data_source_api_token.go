package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-labs/go-unleash-api/api"
)

func dataSourceApiToken() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Retrieves details of a single api token.",

		ReadContext: dataSourceApiTokenRead,

		// This descriptions are used by the documentation generator and the language server.
		Schema: map[string]*schema.Schema{
			"username": {
				Description: "It will return the token defined for this username.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"projects": {
				Description: "It will return the token that have access to the projects defined here.",
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"tokens": {
				Description: "API token",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Description: "The user's username.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"type": {
							Description: "The type of the API token. Can be `client`, `admin` or `frontend`",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"environment": {
							Description: "The environment the token has access to. `\"*\"` means all environments.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"projects": {
							Description: "The project(s) the token will have access to. `[\"*\"]` means all projects.",
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
						},
						"expires_at": {
							Description: "The API token expiration date.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"created_at": {
							Description: "The API token creation date.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"secret": {
							Description: "The API token secret.",
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApiTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	resp, _, err := client.ApiTokens.GetAllApiTokens()
	if err != nil {
		return diag.FromErr(err)
	}
	allTokens := resp.Tokens

	username := d.Get("username").(string)
	projects := d.Get("projects").(*schema.Set).List()
	var foundApiTokens []api.ApiToken
	for _, token := range allTokens {
		if token.Username == username && subslice(toStringArr(projects), token.Projects) {
			foundApiTokens = append(foundApiTokens, token)
		}
	}

	if len(foundApiTokens) > 1 {
		return diag.FromErr(ErrMoreThanOneApiToken)
	}

	d.SetId(buildId(username, toStringArr(projects)))

	tokens := []interface{}{}
	for _, token := range foundApiTokens {
		tfMap := map[string]interface{}{}
		tfMap["username"] = token.Username
		tfMap["type"] = token.Type
		tfMap["environment"] = token.Environment
		tfMap["projects"] = toInterfaceArr(token.Projects)
		tfMap["expires_at"] = token.ExpiresAt
		tfMap["created_at"] = token.CreatedAt
		tfMap["secret"] = token.Secret
		tokens = append(tokens, tfMap)
	}
	_ = d.Set("tokens", tokens)

	return diags
}
