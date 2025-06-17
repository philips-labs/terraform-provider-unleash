package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/philips-labs/go-unleash-api/v2/api"
)

func resourceFeatureV2() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "(Experimental) Provides a resource for managing unleash features with variants and environment strategies all in a single resource.",

		CreateContext: resourceFeatureV2Create,
		ReadContext:   resourceFeatureV2Read,
		UpdateContext: resourceFeatureV2Update,
		DeleteContext: resourceFeatureV2Delete,

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
			"environment": {
				Description: "Use this to enable a feature in an environment and add strategies",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Environment name",
							Type:        schema.TypeString,
							Required:    true,
						},
						"enabled": {
							Description: "Whether the feature is on/off in the environment. Default is `true` (on)",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
						},
						"strategy": {
							Description: "Strategy to add in the environment",
							Type:        schema.TypeList,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Description: "Strategy unique name",
										Type:        schema.TypeString,
										Required:    true,
									},
									"parameters": {
										Description: "Strategy parameters. All the values need to informed as strings.",
										Type:        schema.TypeMap,
										Optional:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"variant": {
										Description: "Feature strategy variant",
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
													Description: "Variant payload. The type of the payload can be `string`, `json` or `csv` or `number`",
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
																Type:        schema.TypeString,
																Description: "Always a string value, independent of the type.",
																Required:    true,
															},
														},
													},
												},
											},
										},
									},
									"constraint": {
										Description: "Strategy constraint",
										Type:        schema.TypeList,
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"context_name": {
													Description:  "Constraint context. Can be `appName`, `currentTime`, `environment`, `sessionId` or `userId`",
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice([]string{"appName", "currentTime", "environment", "sessionId", "userId"}, false),
												},
												"operator": {
													Description:  "Constraint operator. Can be `IN`, `NOT_IN`, `STR_CONTAINS`, `STR_STARTS_WITH`, `STR_ENDS_WITH`, `NUM_EQ`, `NUM_GT`, `NUM_GTE`, `NUM_LT`, `NUM_LTE`, `SEMVER_EQ`, `SEMVER_GT` or `SEMVER_LT`",
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice([]string{"IN", "NOT_IN", "STR_CONTAINS", "STR_STARTS_WITH", "STR_ENDS_WITH", "NUM_EQ", "NUM_GT", "NUM_GTE", "NUM_LT", "NUM_LTE", "SEMVER_EQ", "SEMVER_GT", "SEMVER_LT"}, false),
												},
												"value": {
													Description: "Value to use in the evaluation of the constraint. Applies only to `DATE_`, `NUM_` and `SEMVER_` operators.",
													Type:        schema.TypeString,
													Optional:    true,
												},
												"values": {
													Description: "List of values to use in the evaluation of the constraint. Applies to all operators, except `DATE_`, `NUM_` and `SEMVER_`.",
													Type:        schema.TypeList,
													Optional:    true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"case_insensitive": {
													Description: "If operator is case-insensitive.",
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     false,
												},
												"inverted": {
													Description: "If constraint expressions will be negated, meaning that they get their opposite value.",
													Type:        schema.TypeBool,
													Optional:    true,
													Default:     false,
												},
											},
										},
									},
									"id": {
										Description: "Strategy ID",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"tag": {
				Description: "Tag to add to the feature",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Tag type. Default is `simple`.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "simple",
						},
						"value": {
							Description: "Tag value.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceFeatureV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if e, ok := d.GetOk("environment"); ok {
		tfEnvironments := e.([]interface{})
		for _, tfEnvironment := range tfEnvironments {
			environment := toFeatureEnvironment(tfEnvironment.(map[string]interface{}))

			for _, strategy := range environment.Strategies {
				_, resp, err := client.FeatureToggles.AddStrategyToFeature(feature.Project, feature.Name, environment.Name, strategy)
				if resp == nil || err != nil {
					client.FeatureToggles.ArchiveFeature(feature.Project, feature.Name)
					client.FeatureToggles.DeleteArchivedFeature(feature.Name)
					return diag.FromErr(err)
				}

			}
			ok, _, err := client.FeatureToggles.EnableFeatureOnEnvironment(feature.Project, feature.Name, environment.Name, environment.Enabled)
			if err != nil || !ok {
				client.FeatureToggles.ArchiveFeature(feature.Project, feature.Name)
				client.FeatureToggles.DeleteArchivedFeature(feature.Name)
				return diag.FromErr(err)
			}
		}
	}
	if t, ok := d.GetOk("tag"); ok {
		tfTags := t.([]interface{})
		for _, tfTag := range tfTags {
			tag := toFeatureTag(tfTag.(map[string]interface{}))
			_, resp, err := client.FeatureTags.CreateFeatureTags(feature.Name, tag)
			if resp == nil || err != nil {
				client.FeatureToggles.ArchiveFeature(feature.Project, feature.Name)
				client.FeatureToggles.DeleteArchivedFeature(feature.Name)
				return diag.FromErr(err)
			}
		}

	}

	d.SetId(createdFeature.Name)
	readDiags := resourceFeatureV2Read(ctx, d, meta)
	if readDiags != nil {
		diags = append(diags, readDiags...)
	}

	return diags
}

func resourceFeatureV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if e, ok := d.GetOk("environment"); ok {
		toSave := []api.Environment{}
		for _, tfEnvironment := range e.([]interface{}) {
			for _, env := range feature.Environments {
				if tfEnvironment.(map[string]interface{})["name"] == env.Name {
					toSave = append(toSave, env)
				}
			}
		}
		_ = d.Set("environment", flattenEnvironments(toSave))
	}

	if t, ok := d.GetOk("tag"); ok {
		featureTags, _, err := client.FeatureTags.GetAllFeatureTags(feature.Name)
		if err != nil {
			return diag.FromErr(err)
		}
		toSave := []api.FeatureTag{}
		for _, tfTag := range t.([]interface{}) {
			for _, tag := range featureTags.Tags {
				if tfTag.(map[string]interface{})["value"] == tag.Value {
					toSave = append(toSave, tag)
				}
			}
		}

		_ = d.Set("tag", flattenTags(toSave))
	}

	return diags
}

func resourceFeatureV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if d.HasChange("tag") {
		o, a := d.GetChange("tag")
		old := o.([]interface{})
		new := a.([]interface{})

		toAdd := []api.FeatureTag{}
		toRemove := []api.FeatureTag{}

		for _, newTag := range new {
			newFeatureTag := toFeatureTag(newTag.(map[string]interface{}))
			if !isTagIn(newFeatureTag, old) {
				toAdd = append(toAdd, newFeatureTag)
			}
		}

		for _, oldTag := range old {
			oldFeatureTag := toFeatureTag(oldTag.(map[string]interface{}))
			if !isTagIn(oldFeatureTag, new) {
				toRemove = append(toRemove, oldFeatureTag)
			}
		}

		_, _, err := client.FeatureTags.UpdateFeatureTags(feature.Name, toAdd, toRemove)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("environment") {
		o, a := d.GetChange("environment")
		old := o.([]interface{})
		new := a.([]interface{})

		toAdd := []api.Environment{}
		toUpdate := []api.Environment{}
		toRemove := []api.Environment{}

		for _, newEnv := range new {
			newFeatureEnv := toFeatureEnvironment(newEnv.(map[string]interface{}))
			if isEnvIn(newFeatureEnv.Name, old) {
				toUpdate = append(toUpdate, newFeatureEnv)
			} else {
				toAdd = append(toAdd, newFeatureEnv)
			}
		}

		oldEnvs := []api.Environment{}
		for _, oldEnv := range old {
			oldFeatureEnv := toFeatureEnvironment(oldEnv.(map[string]interface{}))
			oldEnvs = append(oldEnvs, oldFeatureEnv)
			if !isEnvIn(oldFeatureEnv.Name, new) {
				toRemove = append(toRemove, oldFeatureEnv)
			}
		}

		for _, envToUpdate := range toUpdate {
			newStrats := envToUpdate.Strategies
			oldStrats := []api.FeatureStrategy{}
			for _, oldEnv := range oldEnvs {
				if envToUpdate.Name == oldEnv.Name {
					oldStrats = oldEnv.Strategies
				}
			}
			for _, newStrat := range newStrats {
				if isStratIn(newStrat.ID, oldStrats) {
					_, resp, err := client.FeatureToggles.UpdateFeatureStrategy(feature.Project, feature.Name, envToUpdate.Name, newStrat)
					if resp == nil {
						return diag.FromErr(fmt.Errorf("response is nil: %v", err))
					}
					if err != nil {
						return diag.FromErr(err)
					}
				} else {
					_, resp, err := client.FeatureToggles.AddStrategyToFeature(feature.Project, feature.Name, envToUpdate.Name, newStrat)
					if resp == nil {
						return diag.FromErr(fmt.Errorf("response is nil: %v", err))
					}
					if err != nil {
						return diag.FromErr(err)
					}
				}
			}

			for _, oldStrat := range oldStrats {
				if !isStratIn(oldStrat.ID, newStrats) {
					_, _, err = client.FeatureToggles.DeleteStrategyFromFeature(feature.Project, feature.Name, envToUpdate.Name, oldStrat.ID)
					if err != nil {
						return diag.FromErr(err)
					}
				}
			}

			ok, _, err := client.FeatureToggles.EnableFeatureOnEnvironment(feature.Project, feature.Name, envToUpdate.Name, envToUpdate.Enabled)
			if err != nil || !ok {
				return diag.FromErr(err)
			}
		}

		for _, envToRemove := range toRemove {
			for _, strategy := range envToRemove.Strategies {
				_, _, err = client.FeatureToggles.DeleteStrategyFromFeature(feature.Project, feature.Name, envToRemove.Name, strategy.ID)
				if err != nil {
					return diag.FromErr(err)
				}
			}
			ok, _, err := client.FeatureToggles.EnableFeatureOnEnvironment(feature.Project, feature.Name, envToRemove.Name, false)
			if err != nil || !ok {
				return diag.FromErr(err)
			}
		}

		for _, envToAdd := range toAdd {
			for _, strategy := range envToAdd.Strategies {
				_, resp, err := client.FeatureToggles.AddStrategyToFeature(feature.Project, feature.Name, envToAdd.Name, strategy)
				if resp == nil {
					return diag.FromErr(fmt.Errorf("response is nil: %v", err))
				}
				if err != nil {
					return diag.FromErr(err)
				}
			}
			ok, _, err := client.FeatureToggles.EnableFeatureOnEnvironment(feature.Project, feature.Name, envToAdd.Name, envToAdd.Enabled)
			if err != nil || !ok {
				return diag.FromErr(err)
			}
		}

		readDiags := resourceFeatureV2Read(ctx, d, meta)
		if readDiags != nil {
			diags = append(diags, readDiags...)
		}
	}

	return diags
}

func isEnvIn(name string, envs []interface{}) bool {
	for _, env := range envs {
		fEnv := toFeatureEnvironment(env.(map[string]interface{}))
		if fEnv.Name == name {
			return true
		}
	}
	return false
}

func isTagIn(tag api.FeatureTag, tags []interface{}) bool {
	for _, t := range tags {
		tfTag := toFeatureTag(t.(map[string]interface{}))
		if tfTag.Type == tag.Type && tfTag.Value == tag.Value {
			return true
		}
	}
	return false
}

func isStratIn(id string, strats []api.FeatureStrategy) bool {
	for _, strat := range strats {
		if strat.ID == id {
			return true
		}
	}
	return false
}

// Archives a feature
func resourceFeatureV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func toFeatureEnvironment(tfEnvironment map[string]interface{}) api.Environment {
	environment := api.Environment{}
	environment.Name = tfEnvironment["name"].(string)
	environment.Enabled = tfEnvironment["enabled"].(bool)

	if tfStrategies, ok := tfEnvironment["strategy"].([]interface{}); ok && len(tfStrategies) > 0 {
		strategies := make([]api.FeatureStrategy, 0, len(tfStrategies))
		for _, tfStrategy := range tfStrategies {
			strategyMap := tfStrategy.(map[string]interface{})
			name := strategyMap["name"].(string)
			if len(name) > 0 {
				strategy := api.FeatureStrategy{
					Name: name,
				}
				id := strategyMap["id"].(string)
				if len(id) > 0 {
					strategy.ID = id
				}
				if p, ok := strategyMap["parameters"]; ok {
					tfParams := p.(map[string]interface{})
					castedParameters := make(map[string]interface{})
					for k, v := range tfParams {
						castedParameters[k] = v.(string)
					}
					strategy.Parameters = castedParameters
				}
				if tfConstraints, ok := strategyMap["constraint"].([]interface{}); ok && len(tfConstraints) > 0 {
					constraints := make([]api.StrategyConstraint, 0, len(tfConstraints))
					for _, tfConstraint := range tfConstraints {
						constraintMap := tfConstraint.(map[string]interface{})
						constraint := api.StrategyConstraint{
							ContextName:     constraintMap["context_name"].(string),
							Operator:        constraintMap["operator"].(string),
							Value:           constraintMap["value"].(string),
							Values:          toStringArray(constraintMap["values"].([]interface{})),
							Inverted:        constraintMap["inverted"].(bool),
							CaseInsensitive: constraintMap["case_insensitive"].(bool),
						}
						constraints = append(constraints, constraint)
					}
					strategy.Constraints = constraints
				}
				if tfVariants, ok := strategyMap["variant"].([]interface{}); ok && len(tfVariants) > 0 {
					variants := make([]api.Variant, 0, len(tfVariants))
					for _, tfVariant := range tfVariants {
						variants = append(variants, toFeatureVariant(tfVariant.(map[string]interface{})))
					}
					strategy.Variants = variants
				}
				strategies = append(strategies, strategy)
			}
		}
		environment.Strategies = strategies
	}
	return environment
}

func flattenEnvironments(environments []api.Environment) []interface{} {
	if environments == nil {
		return []interface{}{}
	}

	tfEnvironments := []interface{}{}

	for _, env := range environments {
		tfEnvironment := map[string]interface{}{}
		tfEnvironment["name"] = env.Name
		tfEnvironment["enabled"] = env.Enabled

		if env.Strategies != nil {
			tfStrategies := []interface{}{}
			for _, strategy := range env.Strategies {
				tfStrategy := map[string]interface{}{}
				tfStrategy["id"] = strategy.ID
				tfStrategy["name"] = strategy.Name
				retrievedParams := strategy.Parameters.(map[string]interface{})
				castedParams := make(map[string]interface{})
				for k, v := range retrievedParams {
					castedParams[k] = v.(string)
				}
				tfStrategy["parameters"] = castedParams
				if strategy.Constraints != nil {
					tfConstraints := []interface{}{}
					for _, constraint := range strategy.Constraints {
						tfConstraint := map[string]interface{}{}
						tfConstraint["context_name"] = constraint.ContextName
						tfConstraint["operator"] = constraint.Operator
						tfConstraint["value"] = constraint.Value
						tfConstraint["values"] = constraint.Values
						tfConstraint["inverted"] = constraint.Inverted
						tfConstraint["case_insensitive"] = constraint.CaseInsensitive
						tfConstraints = append(tfConstraints, tfConstraint)
					}
					tfStrategy["constraint"] = tfConstraints
				}
				if strategy.Variants != nil {
					tfVariants := []interface{}{}
					for _, variant := range strategy.Variants {
						tfVariant := map[string]interface{}{}
						tfVariant["name"] = variant.Name
						tfVariant["stickiness"] = variant.Stickiness
						tfVariant["weight"] = variant.Weight
						tfVariant["weight_type"] = variant.WeightType
						if variant.Payload != nil {
							payload := map[string]interface{}{
								"type":  variant.Payload.Type,
								"value": variant.Payload.Value,
							}
							tfVariant["payload"] = []interface{}{payload}
						}
						tfVariants = append(tfVariants, tfVariant)
					}
					tfStrategy["variant"] = tfVariants
				}
				tfStrategies = append(tfStrategies, tfStrategy)
			}
			tfEnvironment["strategy"] = tfStrategies
		}

		tfEnvironments = append(tfEnvironments, tfEnvironment)
	}

	return tfEnvironments
}

func flattenTags(tags []api.FeatureTag) []interface{} {
	if tags == nil {
		return []interface{}{}
	}

	tfTags := []interface{}{}
	for _, tag := range tags {
		tfTag := map[string]interface{}{}
		tfTag["type"] = tag.Type
		tfTag["value"] = tag.Value

		tfTags = append(tfTags, tfTag)
	}
	return tfTags
}

func toFeatureTag(tfTag map[string]interface{}) api.FeatureTag {
	tag := api.FeatureTag{}
	tag.Type = tfTag["type"].(string)
	tag.Value = tfTag["value"].(string)
	return tag
}
