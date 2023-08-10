package provider

import (
	"context"
	"fmt"

	openapiclient "github.com/Unleash/unleash-server-api-go/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/philips-labs/go-unleash-api/api"
)

func resourceFeature() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a resource for managing unleash features.",

		CreateContext: resourceFeatureCreate,
		ReadContext:   resourceFeatureRead,
		UpdateContext: resourceFeatureUpdate,
		DeleteContext: resourceFeatureDelete,

		// The descriptions are used by the documentation generator and the language server.
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Feature name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"project_id": {
				Description: "The feature will be created in the given project",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Description: "Feature type",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Feature description",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"archive_on_destroy": {
				Description: "Whether to archive the feature toggle on destroy. Default is `true`. When `false`, it will permanently delete the feature toggle.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"variant": {
				Description: "Feature variant",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Variant name",
							Type:        schema.TypeString,
							Required:    true,
						},
						"stickiness": {
							Description: "Variant stickiness. Default is `default`.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "default",
						},
						"weight": {
							Description:  "Variant weight. Only considered when the `weight_type` is `fix`. It is calculated automatically if the `weight_type` is `variable`.",
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"weight_type": {
							Description: "Variant weight type. The weight type can be `fix` or `variable`. Default is `variable`.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "variable",
						},
						"payload": {
							Description: "Variant payload. The type of the payload can be `string`, `json` or `csv`",
							Type:        schema.TypeSet,
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"overrides": {
							Description: "Overrides existing context field values. Values are comma separated e.g `v1, v2, ...`)",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"context_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceFeatureCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).PhilipsUnleashClient

	var diags diag.Diagnostics

	feature := &api.FeatureToggle{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Project:     d.Get("project_id").(string),
	}

	createdFeature, resp, err := client.FeatureToggles.CreateFeature(feature.Project, *feature)
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	if p, ok := d.GetOk("variant"); ok {
		tfVariants := p.([]interface{})
		variants := make([]api.Variant, 0, len(tfVariants))
		for _, tfVariant := range tfVariants {
			variants = append(variants, toFeatureVariant(tfVariant.(map[string]interface{})))
		}
		_, resp, err := client.Variants.AddVariantsForFeatureToggle(feature.Project, feature.Name, variants)
		if resp == nil {
			return diag.FromErr(fmt.Errorf("response is nil: %v", err))
		}
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(createdFeature.Name)
	readDiags := resourceFeatureRead(ctx, d, meta)
	if readDiags != nil {
		diags = append(diags, readDiags...)
	}

	return diags
}

func resourceFeatureRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).PhilipsUnleashClient

	var diags diag.Diagnostics

	featureName := d.Id()
	projectId := d.Get("project_id").(string)
	feature, _, err := client.FeatureToggles.GetFeatureByName(projectId, featureName)
	if err != nil {
		if err == api.ErrNotFound {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}
	_ = d.Set("name", feature.Name)
	_ = d.Set("description", feature.Description)
	_ = d.Set("type", feature.Type)
	_ = d.Set("project_id", feature.Project)
	_ = d.Set("variant", flattenVariants(feature.Variants))

	return diags
}

func resourceFeatureUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).PhilipsUnleashClient

	var diags diag.Diagnostics

	feature := &api.FeatureToggle{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Project:     d.Get("project_id").(string),
	}

	_, resp, err := client.FeatureToggles.UpdateFeature(feature.Project, *feature)
	if resp == nil {
		return diag.FromErr(fmt.Errorf("response is nil: %v", err))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("variant") {
		tfVariants := d.Get("variant").([]interface{})
		variants := make([]api.Variant, 0, len(tfVariants))
		for _, tfVariant := range tfVariants {
			variants = append(variants, toFeatureVariant(tfVariant.(map[string]interface{})))
		}
		_, resp, err := client.Variants.AddVariantsForFeatureToggle(feature.Project, feature.Name, variants)
		if resp == nil {
			return diag.FromErr(fmt.Errorf("response is nil: %v", err))
		}
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

// Archives a feature
func resourceFeatureDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ApiClients).PhilipsUnleashClient

	var diags diag.Diagnostics

	featureName := d.Id()
	projectId := d.Get("project_id").(string)
	_, _, err := client.FeatureToggles.ArchiveFeature(projectId, featureName)
	if err != nil {
		return diag.FromErr(err)
	}
	shouldArchive := d.Get("archive_on_destroy").(bool)
	if !shouldArchive {
		_, _, err := client.FeatureToggles.DeleteArchivedFeature(featureName)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return diags
}

func toStringArray(iArr []interface{}) []string {
	stringArr := make([]string, len(iArr))
	for i, v := range iArr {
		stringArr[i] = v.(string)
	}
	return stringArr
}

func toOpenApiFeatureVariant(tfVariant map[string]interface{}) openapiclient.VariantSchema {
	name := tfVariant["name"].(string)
	weight := float32(tfVariant["weight"].(int))
	stickiness := tfVariant["stickiness"].(string)
	weightType := tfVariant["weight_type"].(string)

	v := *openapiclient.NewVariantSchema(name, weight)
	v.Stickiness = &stickiness
	v.WeightType = &weightType

	if payloadSet, ok := tfVariant["payload"].(*schema.Set); ok && payloadSet.Len() > 0 {
		payloadList := payloadSet.List()
		payloadMap := payloadList[0].(map[string]interface{})
		v.Payload = openapiclient.NewVariantSchemaPayload(payloadMap["type"].(string), payloadMap["value"].(string))
	}

	if overridesSet, ok := tfVariant["overrides"].(*schema.Set); ok && overridesSet.Len() > 0 {
		overridesList := overridesSet.List()
		overrides := make([]openapiclient.OverrideSchema, 0, len(overridesList))
		for _, tfOverride := range overridesList {
			overrideMap := tfOverride.(map[string]interface{})
			override := *openapiclient.NewOverrideSchema(overrideMap["context_name"].(string), toStringArray(overrideMap["values"].([]interface{})))
			overrides = append(overrides, override)
		}
		v.Overrides = overrides
	}
	return v
}

func toFeatureVariant(tfVariant map[string]interface{}) api.Variant {
	variant := api.Variant{}
	variant.Name = tfVariant["name"].(string)
	variant.Stickiness = tfVariant["stickiness"].(string)
	variant.Weight = tfVariant["weight"].(int)
	variant.WeightType = tfVariant["weight_type"].(string)

	if payloadSet, ok := tfVariant["payload"].(*schema.Set); ok && payloadSet.Len() > 0 {
		payloadList := payloadSet.List()
		payloadMap := payloadList[0].(map[string]interface{})
		variant.Payload = &api.VariantPayload{
			Type:  payloadMap["type"].(string),
			Value: payloadMap["value"].(string),
		}
	}

	if overridesSet, ok := tfVariant["overrides"].(*schema.Set); ok && overridesSet.Len() > 0 {
		overridesList := overridesSet.List()
		overrides := make([]api.VariantOverride, 0, len(overridesList))
		for _, tfOverride := range overridesList {
			overrideMap := tfOverride.(map[string]interface{})
			override := api.VariantOverride{
				ContextName: overrideMap["context_name"].(string),
				Values:      toStringArray(overrideMap["values"].([]interface{})),
			}
			overrides = append(overrides, override)
		}
		variant.Overrides = overrides
	}
	return variant
}

func flattenVariants(variants []api.Variant) []interface{} {
	if variants == nil {
		return []interface{}{}
	}

	vVariants := []interface{}{}

	for _, variant := range variants {
		mVariant := map[string]interface{}{}
		mVariant["name"] = variant.Name
		mVariant["weight"] = variant.Weight
		mVariant["weight_type"] = variant.WeightType
		mVariant["stickiness"] = variant.Stickiness

		if variant.Payload != nil {
			mPayloads := []interface{}{}
			mPayload := map[string]interface{}{}
			mPayload["type"] = variant.Payload.Type
			mPayload["value"] = variant.Payload.Value
			mPayloads = append(mPayloads, mPayload)
			mVariant["payload"] = mPayloads
		}

		if variant.Overrides != nil {
			vOverrides := []interface{}{}
			for _, override := range variant.Overrides {
				mOverride := map[string]interface{}{}
				mOverride["context_name"] = override.ContextName
				mOverride["values"] = override.Values
				vOverrides = append(vOverrides, mOverride)
			}
			mVariant["overrides"] = vOverrides
		}

		vVariants = append(vVariants, mVariant)
	}

	return vVariants
}
