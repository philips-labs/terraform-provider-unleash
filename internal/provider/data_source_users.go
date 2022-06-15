package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/philips-labs/go-unleash-api/api"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Retrieve a collection of users that match the provided query.",

		ReadContext: dataSourceUsersRead,

		// This descriptions are used by the documentation generator and the language server.
		Schema: map[string]*schema.Schema{
			"query": {
				Description:  "Query used to search the user. It searches by `email`, `username` and `name` fields of users.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 255),
			},
			"users": {
				Description: "Collection of users that match the provided query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The user's id.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "The user's name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"email": {
							Description: "The user's email address.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"username": {
							Description: "The user's username.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"root_role": {
							Description: "The user's role.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"created_at": {
							Description: "The date of creation of the user.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"image_url": {
							Description: "The user's image URL.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	query := d.Get("query").(string)

	matchedUsers, _, err := client.Users.SearchUser(query)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(query)

	users := []interface{}{}
	for _, user := range *matchedUsers {
		userDetails, _, err := client.Users.GetUserById(strconv.Itoa(user.Id))
		if err != nil {
			return diag.FromErr(err)
		}
		tfMap := map[string]interface{}{}
		tfMap["id"] = userDetails.Id
		tfMap["name"] = userDetails.Name
		tfMap["username"] = userDetails.Username
		tfMap["email"] = userDetails.Email
		for k, v := range rolesLookup {
			if v == userDetails.RootRole {
				tfMap["root_role"] = k
			}
		}
		tfMap["created_at"] = userDetails.CreatedAt
		tfMap["image_url"] = userDetails.ImageUrl
		users = append(users, tfMap)
	}
	_ = d.Set("users", users)

	return diags
}
