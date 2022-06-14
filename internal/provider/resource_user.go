package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/philips-labs/go-unleash-api/api"
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
			"user_id": {
				Description: "The user's id.",
				Type:        schema.TypeInt,
				Computed:    true,
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
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	givenUserRole := d.Get("root_role").(string)
	roleId := getUserRoleId(givenUserRole)
	user := &api.User{
		Name:      d.Get("name").(string),
		Email:     d.Get("email").(string),
		Username:  d.Get("username").(string),
		RootRole:  roleId,
		SendEmail: d.Get("send_email").(bool),
	}

	createdUser, resp, err := client.Users.CreateUser(*user)
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("invite_link", createdUser.InviteLink)
	_ = d.Set("email_sent", createdUser.EmailSent)
	d.SetId(strconv.Itoa(createdUser.Id))
	readDiags := resourceUserRead(ctx, d, meta)
	if readDiags != nil {
		diags = append(diags, readDiags...)
	}

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	userId := d.Id()
	user, _, err := client.Users.GetUserById(userId)
	if err != nil {
		if err == api.ErrNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("user_id", user.Id)
	_ = d.Set("name", user.Name)
	_ = d.Set("username", user.Username)
	_ = d.Set("email", user.Email)
	_ = d.Set("root_role", getUserRole(user.RootRole))

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	givenUserRole := d.Get("root_role").(string)
	roleId := getUserRoleId(givenUserRole)
	user := &api.User{
		Name:      d.Get("name").(string),
		Email:     d.Get("email").(string),
		Username:  d.Get("username").(string),
		RootRole:  roleId,
		SendEmail: d.Get("send_email").(bool),
	}

	userId := d.Id()
	_, resp, err := client.Users.UpdateUser(userId, *user)
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.ApiClient)

	var diags diag.Diagnostics

	userId := d.Id()
	_, _, err := client.Users.DeleteUser(userId)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
