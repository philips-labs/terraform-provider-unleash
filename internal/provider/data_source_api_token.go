package provider

import (
	"context"

	openapiclient "github.com/Unleash/unleash-server-api-go/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApiToken() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Retrieves a single api token based on provided filters. It raises an error if more than one token is returned.",

		ReadContext: dataSourceApiTokenRead,

		// This descriptions are used by the documentation generator and the language server.
		Schema: map[string]*schema.Schema{
			"token_name": {
				Description: "Filter token by the unique name of the token. This property replaced `username` in Unleash v5).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"projects": {
				Description: "Filter token by project(s).",
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"token": {
				Description: "API token",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token_name": {
							Description: "The unique name of the token. This property replaced `username` in Unleash v5).",
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
	client := meta.(*ApiClients).UnleashClient

	var diags diag.Diagnostics

	resp, _, err := client.APITokensAPI.GetAllApiTokens(ctx).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	allTokens := resp.Tokens

	tokenName := d.Get("token_name").(string)
	projects := d.Get("projects").(*schema.Set).List()
	var foundApiTokens []openapiclient.ApiTokenSchema
	for _, token := range allTokens {
		if token.TokenName == tokenName && subslice(toStringArr(projects), token.Projects) {
			foundApiTokens = append(foundApiTokens, token)
		}
	}

	if len(foundApiTokens) > 1 {
		return diag.FromErr(ErrMoreThanOneApiToken)
	}

	d.SetId(buildId(tokenName, toStringArr(projects)))

	tokens := []interface{}{}
	token := foundApiTokens[0]

	tfMap := map[string]interface{}{}
	tfMap["token_name"] = token.TokenName
	tfMap["type"] = token.Type
	tfMap["environment"] = token.Environment
	tfMap["projects"] = toInterfaceArr(token.Projects)
	tfMap["expires_at"] = token.ExpiresAt
	tfMap["created_at"] = token.CreatedAt
	tfMap["secret"] = token.Secret

	tokens = append(tokens, tfMap)
	_ = d.Set("token", tokens)

	return diags
}
