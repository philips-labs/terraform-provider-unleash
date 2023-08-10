package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Retrieve details of an existing unleash project",

		ReadContext: dataSourceProjectRead,

		// This descriptions are used by the documentation generator and the language server.
		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "The project id of the unleash project",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Project name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "The date the unleash project was last updated",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "The description of the unleash project",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"environments": {
				Description: "The list of unleash environments in this project",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment": {
							Description: "The environment name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).PhilipsUnleashClient

	var diags diag.Diagnostics

	projectId := d.Get("project_id").(string)

	foundProject, _, err := client.Projects.GetProjectById(projectId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(foundProject.Name)
	_ = d.Set("name", foundProject.Name)
	_ = d.Set("description", foundProject.Description)
	_ = d.Set("updatedAt", foundProject.UpdatedAt)

	envs := []interface{}{}
	for _, env := range foundProject.Environments {
		tfMap := map[string]interface{}{}
		tfMap["environment"] = env.Environment
		envs = append(envs, tfMap)
	}
	_ = d.Set("environments", envs)

	return diags
}
