package provider

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/philips-labs/go-unleash-api/api"
)

func resourceApiToken() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a resource for managing unleash api tokens.",

		CreateContext: resourceApiTokenCreate,
		ReadContext:   resourceApiTokenRead,
		UpdateContext: resourceApiTokenUpdate,
		DeleteContext: resourceApiTokenDelete,

		// The descriptions are used by the documentation generator and the language server.
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Description:  "The type of the API token. Can be `client`, `admin` or `frontend`",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"client", "admin", "frontend"}, false),
			},
			"environment": {
				Description: "The environment the token will have access to. Use `\"*\"` for all environments. By default, it will have access to all environments.",
				Type:        schema.TypeString,
				Default:     "*",
				Optional:    true,
				ForceNew:    true,
			},
			"projects": {
				Description: "The project(s) the token will have access to. Use `[\"*\"]` for all projects. By default, it will have access to all projects.",
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				ForceNew:    true,
			},
			"expires_at": {
				Description: "The API token expiration date.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"created_at": {
				Description: "The API token creation date.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"secret": {
				Description: "The API token secret.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceApiTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).PhilipsUnleashClient

	var diags diag.Diagnostics

	apiToken := &api.ApiToken{
		Username:    d.Get("username").(string),
		Type:        d.Get("type").(string),
		Environment: d.Get("environment").(string),
		Projects:    toStringArr(d.Get("projects").(*schema.Set).List()),
		ExpiresAt:   d.Get("expires_at").(string),
	}

	createdToken, resp, err := client.ApiTokens.CreateApiToken(*apiToken)
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("secret", createdToken.Secret)
	_ = d.Set("created_at", createdToken.CreatedAt)
	d.SetId(toMD5Str(createdToken.Secret))
	readDiags := resourceApiTokenRead(ctx, d, meta)
	if readDiags != nil {
		diags = append(diags, readDiags...)
	}

	return diags
}

func resourceApiTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).PhilipsUnleashClient

	var diags diag.Diagnostics

	secret := d.Get("secret").(string)
	resp, _, err := client.ApiTokens.GetAllApiTokens()
	if err != nil {
		return diag.FromErr(err)
	}
	tokens := resp.Tokens

	var foundApiToken api.ApiToken
	for _, token := range tokens {
		if token.Secret == secret {
			foundApiToken = token
			break
		}
	}

	_ = d.Set("username", foundApiToken.Username)
	_ = d.Set("type", foundApiToken.Type)
	_ = d.Set("environment", foundApiToken.Environment)
	_ = d.Set("created_at", foundApiToken.CreatedAt)
	_ = d.Set("secret", foundApiToken.Secret)

	return diags
}

func resourceApiTokenUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).PhilipsUnleashClient

	var diags diag.Diagnostics

	apiToken := &api.ApiToken{
		Username:    d.Get("username").(string),
		Type:        d.Get("type").(string),
		Environment: d.Get("environment").(string),
		Projects:    toStringArr(d.Get("projects").(*schema.Set).List()),
		ExpiresAt:   d.Get("expires_at").(string),
	}

	tokenSecret := d.Get("secret").(string)
	_, resp, err := client.ApiTokens.UpdateApiToken(tokenSecret, *apiToken)
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceApiTokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).PhilipsUnleashClient

	var diags diag.Diagnostics

	tokenSecret := d.Get("secret").(string)
	_, _, err := client.ApiTokens.DeleteApiToken(tokenSecret)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func toStringArr(tfList []interface{}) []string {
	stringArr := make([]string, 0, len(tfList))
	for _, v := range tfList {
		stringArr = append(stringArr, v.(string))
	}
	return stringArr
}

func toMD5Str(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum([]byte{}))
}
