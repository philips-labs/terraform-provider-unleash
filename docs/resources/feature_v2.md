---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "unleash_feature_v2 Resource - terraform-provider-unleash"
subcategory: ""
description: |-
  (Experimental) Provides a resource for managing unleash features with variants and environment strategies all in a single resource.
---

# unleash_feature_v2 (Resource)

(Experimental) Provides a resource for managing unleash features with variants and environment strategies all in a single resource.

## Example Usage

```terraform
resource "unleash_feature_v2" "with_env_strategies" {
  name               = "my_nice_feature"
  description        = "manages my nice feature"
  type               = "release"
  project_id         = "default"
  archive_on_destroy = false

  environment {
    name    = "production"
    enabled = false
  }

  environment {
    name    = "development"
    enabled = true

    strategy {
      name = "remoteAddress"
      parameters = {
        IPs = "189.434.777.123,host.test.com"
      }
      constraint {
        context_name     = "appName"
        operator         = "SEMVER_EQ"
        case_insensitive = false
        inverted         = false
        value            = "1.0.0"
      }
      constraint {
        context_name = "appName"
        operator     = "IN"
        values       = ["foo", "bar"]
      }
    }
    strategy {
      name = "flexibleRollout"
      parameters = {
        rollout    = "68"
        stickiness = "random"
        groupId    = "toggle"
      }
      variant {
        name = "a" # if you see drifts with multiple variants, sort them by name.
        payload {
          type  = "string"
          value = "foo"
        }
      }
    }
  }

  tag {
    type  = "simple"
    value = "foo"
  }

  tag {
    type  = "simple"
    value = "bar"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Feature name
- `project_id` (String) The feature will be created in the given project
- `type` (String) Feature type

### Optional

- `archive_on_destroy` (Boolean) Whether to archive the feature toggle on destroy. Default is `true`. When `false`, it will permanently delete the feature toggle.
- `description` (String) Feature description
- `environment` (Block List) Use this to enable a feature in an environment and add strategies (see [below for nested schema](#nestedblock--environment))
- `tag` (Block List) Tag to add to the feature (see [below for nested schema](#nestedblock--tag))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--environment"></a>
### Nested Schema for `environment`

Required:

- `name` (String) Environment name

Optional:

- `enabled` (Boolean) Whether the feature is on/off in the environment. Default is `true` (on)
- `strategy` (Block List) Strategy to add in the environment (see [below for nested schema](#nestedblock--environment--strategy))

<a id="nestedblock--environment--strategy"></a>
### Nested Schema for `environment.strategy`

Required:

- `name` (String) Strategy unique name

Optional:

- `constraint` (Block List) Strategy constraint (see [below for nested schema](#nestedblock--environment--strategy--constraint))
- `parameters` (Map of String) Strategy parameters. All the values need to informed as strings.
- `variant` (Block List) Feature strategy variant. The api returns them sorted by name, so if you see drifts, sort them by name when declaring them in the configuration as well. (see [below for nested schema](#nestedblock--environment--strategy--variant))

Read-Only:

- `id` (String) Strategy ID

<a id="nestedblock--environment--strategy--constraint"></a>
### Nested Schema for `environment.strategy.constraint`

Required:

- `context_name` (String) Constraint context. Can be `appName`, `currentTime`, `environment`, `sessionId` or `userId`
- `operator` (String) Constraint operator. Can be `IN`, `NOT_IN`, `STR_CONTAINS`, `STR_STARTS_WITH`, `STR_ENDS_WITH`, `NUM_EQ`, `NUM_GT`, `NUM_GTE`, `NUM_LT`, `NUM_LTE`, `SEMVER_EQ`, `SEMVER_GT` or `SEMVER_LT`

Optional:

- `case_insensitive` (Boolean) If operator is case-insensitive.
- `inverted` (Boolean) If constraint expressions will be negated, meaning that they get their opposite value.
- `value` (String) Value to use in the evaluation of the constraint. Applies only to `DATE_`, `NUM_` and `SEMVER_` operators.
- `values` (List of String) List of values to use in the evaluation of the constraint. Applies to all operators, except `DATE_`, `NUM_` and `SEMVER_`.


<a id="nestedblock--environment--strategy--variant"></a>
### Nested Schema for `environment.strategy.variant`

Required:

- `name` (String) Variant name

Optional:

- `payload` (Block Set, Max: 1) Variant payload. The type of the payload can be `string`, `json` or `csv` or `number` (see [below for nested schema](#nestedblock--environment--strategy--variant--payload))
- `stickiness` (String) Variant stickiness. Default is `default`.
- `weight` (Number) Variant weight. Only considered when the `weight_type` is `fix`. It is calculated automatically if the `weight_type` is `variable`.
- `weight_type` (String) Variant weight type. The weight type can be `fix` or `variable`. Default is `variable`.

<a id="nestedblock--environment--strategy--variant--payload"></a>
### Nested Schema for `environment.strategy.variant.payload`

Required:

- `type` (String)
- `value` (String) Always a string value, independent of the type.





<a id="nestedblock--tag"></a>
### Nested Schema for `tag`

Required:

- `value` (String) Tag value.

Optional:

- `type` (String) Tag type. Default is `simple`.
