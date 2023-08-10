package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	openapiclient "github.com/Unleash/unleash-server-api-go/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a resource for managing unleash users.",

		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,

		// The descriptions are used by the documentation generator and the language server.
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The user's name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"email": {
				Description:  "The user's email address.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}`), "must be a valid email with lowercase letters"),
			},
			"username": {
				Description: "The user's username.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"root_role": {
				Description:  "The role to assign to the user. Can be `Admin`, `Editor` or `Viewer`",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Admin", "Editor", "Viewer"}, false),
			},
			"send_email": {
				Description: "Whether to send a welcome email with a login link to the user or not. Defaults to `true`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"invite_link": {
				Description: "The link for the login link.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email_sent": {
				Description: "Whether the welcome email was successfully sent to the user.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).UnleashClient

	var diags diag.Diagnostics

	givenUserRole := d.Get("root_role").(string)
	roleId := int32(rolesLookup[givenUserRole])
	givenName := d.Get("name").(string)
	givenEmail := d.Get("email").(string)
	givenUsername := d.Get("username").(string)
	givenSendEmail := d.Get("send_email").(bool)

	createUserSchema := *openapiclient.NewCreateUserSchemaWithDefaults()
	createUserSchema.Name = &givenName
	createUserSchema.Email = &givenEmail
	createUserSchema.Username = &givenUsername
	createUserSchema.SendEmail = &givenSendEmail
	createUserSchema.RootRole = openapiclient.Int32AsCreateUserSchemaRootRole(&roleId)

	createdUser, resp, err := client.UsersApi.CreateUser(ctx).CreateUserSchema(createUserSchema).Execute()
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("invite_link", createdUser.InviteLink)
	_ = d.Set("email_sent", createdUser.EmailSent)
	d.SetId(strconv.Itoa(int(createdUser.Id)))
	readDiags := resourceUserRead(ctx, d, meta)
	if readDiags != nil {
		diags = append(diags, readDiags...)
	}

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).UnleashClient

	var diags diag.Diagnostics

	userId := d.Id()
	user, resp, err := client.UsersApi.GetUser(ctx, userId).Execute()
	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("name", user.Name)
	_ = d.Set("username", user.Username)
	_ = d.Set("email", user.Email)

	for k, v := range rolesLookup {
		if int32(v) == *user.RootRole {
			_ = d.Set("root_role", k)
		}
	}

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).UnleashClient

	var diags diag.Diagnostics

	givenUserRole := d.Get("root_role").(string)
	roleId := int32(rolesLookup[givenUserRole])
	givenName := d.Get("name").(string)
	givenEmail := d.Get("email").(string)
	rootRole := openapiclient.Int32AsCreateUserSchemaRootRole(&roleId)

	requestBody := map[string]interface{}{
		"name":     &givenName,
		"email":    &givenEmail,
		"rootRole": &rootRole,
	}

	userId := d.Id()
	_, resp, err := client.UsersApi.UpdateUser(ctx, userId).RequestBody(requestBody).Execute()
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).UnleashClient

	var diags diag.Diagnostics

	userId := d.Id()
	_, err := client.UsersApi.DeleteUser(ctx, userId).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
