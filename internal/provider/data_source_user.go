package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Retrieve details of an existing user",

		ReadContext: dataSourceUserRead,

		// This descriptions are used by the documentation generator and the language server.
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Id used to search the user.",
				Type:        schema.TypeInt,
				Required:    true,
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
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).PhilipsUnleashClient

	var diags diag.Diagnostics

	id := d.Get("id").(int)
	stringId := strconv.Itoa(id)
	userDetails, _, err := client.Users.GetUserById(stringId)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(stringId)

	_ = d.Set("name", userDetails.Name)
	_ = d.Set("username", userDetails.Username)
	_ = d.Set("email", userDetails.Email)
	for k, v := range rolesLookup {
		if v == userDetails.RootRole {
			_ = d.Set("root_role", k)
		}
	}
	_ = d.Set("created_at", userDetails.Email)
	_ = d.Set("image_url", userDetails.ImageUrl)

	return diags
}
