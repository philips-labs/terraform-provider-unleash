package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-labs/go-unleash-api/api"
)

func dataSourceApiTokens() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Retrieves existing api tokens. Filters are optional.",

		ReadContext: dataSourceApiTokensRead,

		// This descriptions are used by the documentation generator and the language server.
		Schema: map[string]*schema.Schema{
			"username": {
				Description: "Filter tokens by username.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"projects": {
				Description: "Filter tokens by project(s).",
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"tokens": {
				Description: "List of api tokens.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
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

func dataSourceApiTokensRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	resp, _, err := client.ApiTokens.GetAllApiTokens()
	if err != nil {
		return diag.FromErr(err)
	}
	allTokens := resp.Tokens

	u, uOk := d.GetOk("username")
	p, pOk := d.GetOk("projects")

	var foundApiTokens []api.ApiToken
	if !uOk && !pOk {
		foundApiTokens = allTokens
		d.SetId(buildId("*", []string{"*"}))
	} else {
		username := u.(string)
		projects := p.(*schema.Set).List()
		for _, token := range allTokens {
			if (username == "" || token.Username == username) && subslice(toStringArr(projects), token.Projects) {
				foundApiTokens = append(foundApiTokens, token)
			}
		}
		d.SetId(buildId(username, toStringArr(projects)))
	}

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

func buildId(username string, projects []string) string {
	projectsStr := strings.Join(projects[:], ",")
	query := username + projectsStr
	return toMD5Str(query)
}

func toInterfaceArr(stringArr []string) []interface{} {
	tfList := make([]interface{}, 0, len(stringArr))
	for _, v := range stringArr {
		tfList = append(tfList, v)
	}
	return tfList
}

func subslice(s1 []string, s2 []string) bool {
	if len(s1) > len(s2) {
		return false
	}
	for _, e := range s1 {
		if !contains(s2, e) {
			return false
		}
	}
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
