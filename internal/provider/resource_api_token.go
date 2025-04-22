package provider

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	openapiclient "github.com/Unleash/unleash-server-api-go/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Description: "The name of the token. Used as `tokenName` in the API (username is deprecated).",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Description:  "The type of the API token. Can be `client`, `admin` or `frontend`",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"client", "admin", "frontend"}, false),
			},
			"environment": {
				Description: "The environment the token will have access to. By default, it has access to the `development` environment.",
				Type:        schema.TypeString,
				Default:     "development",
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
	client := meta.(*ApiClients).UnleashClient

	var diags diag.Diagnostics

	tokenName := d.Get("username").(string)
	tokenType := d.Get("type").(string)
	environment := d.Get("environment").(string)
	projects := toStringArr(d.Get("projects").(*schema.Set).List())
	expiresAt := d.Get("expires_at").(string)

	createApiTokenSchema := openapiclient.CreateApiTokenSchema{CreateApiTokenSchemaOneOf2: openapiclient.NewCreateApiTokenSchemaOneOf2(tokenType, tokenName)}
	createApiTokenSchema.CreateApiTokenSchemaOneOf2.Environment = &environment
	createApiTokenSchema.CreateApiTokenSchemaOneOf2.Projects = projects
	createApiTokenSchema.CreateApiTokenSchemaOneOf2.ExpiresAt = nil
	if expiresAt != "" {
		res, parseErr := time.Parse(time.RFC3339, expiresAt)
		if parseErr != nil {
			return diag.FromErr(parseErr)
		}
		createApiTokenSchema.CreateApiTokenSchemaOneOf2.ExpiresAt = &res
	}

	createdToken, resp, err := client.APITokensAPI.CreateApiToken(ctx).CreateApiTokenSchema(createApiTokenSchema).Execute()
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("secret", createdToken.Secret)
	_ = d.Set("created_at", createdToken.CreatedAt.Format(time.RFC3339))
	d.SetId(toMD5Str(createdToken.Secret))
	readDiags := resourceApiTokenRead(ctx, d, meta)
	if readDiags != nil {
		diags = append(diags, readDiags...)
	}

	return diags
}

func resourceApiTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).UnleashClient

	var diags diag.Diagnostics

	secret := d.Get("secret").(string)
	resp, _, err := client.APITokensAPI.GetAllApiTokens(ctx).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	tokens := resp.Tokens

	var foundApiToken openapiclient.ApiTokenSchema
	for _, token := range tokens {
		if token.Secret == secret {
			foundApiToken = token
			break
		}
	}

	_ = d.Set("username", foundApiToken.TokenName)
	_ = d.Set("type", foundApiToken.Type)
	_ = d.Set("environment", foundApiToken.Environment)
	_ = d.Set("created_at", foundApiToken.CreatedAt.Format(time.RFC3339))
	_ = d.Set("secret", foundApiToken.Secret)

	return diags
}

func resourceApiTokenUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).UnleashClient

	fmt.Print("Updating API token...")
	var diags diag.Diagnostics

	expiresAt := d.Get("expires_at").(string)
	parsedExpiresAt, parseErr := time.Parse(time.RFC3339, expiresAt)
	if parseErr != nil {
		return diag.FromErr(parseErr)
	}

	updateApiTokenSchema := *openapiclient.NewUpdateApiTokenSchemaWithDefaults()
	updateApiTokenSchema.ExpiresAt = parsedExpiresAt
	tokenSecret := d.Get("secret").(string)
	resp, err := client.APITokensAPI.UpdateApiToken(ctx, tokenSecret).UpdateApiTokenSchema(updateApiTokenSchema).Execute()
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceApiTokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).UnleashClient

	var diags diag.Diagnostics

	tokenSecret := d.Get("secret").(string)
	_, err := client.APITokensAPI.DeleteApiToken(ctx, tokenSecret).Execute()
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
