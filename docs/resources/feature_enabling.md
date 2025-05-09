---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "unleash_feature_enabling Resource - terraform-provider-unleash"
subcategory: ""
description: |-
  Provides a resource for enabling a feature toggle in the given environment. This can be only done after the feature toggle has at least one strategy.
---

# unleash_feature_enabling (Resource)

Provides a resource for enabling a feature toggle in the given environment. This can be only done after the feature toggle has at least one strategy.

## Example Usage

```terraform
data "unleash_project" "example" {
  project_id = "default"
}

data "unleash_feature_type" "example" {
  type_id = "kill-switch"
}

resource "unleash_feature" "example" {
  name       = "toggle"
  project_id = data.unleash_project.example.project_id
  type       = data.unleash_feature_type.example.type_id
}

resource "unleash_strategy_assignment" "example" {
  feature_name  = unleash_feature.example.name
  project_id    = data.unleash_project.example.project_id
  environment   = "development"
  strategy_name = "flexibleRollout"
  parameters = {
    rollout    = "68"
    stickiness = "random"
    groupId    = "toggle"
  }
}

resource "unleash_feature_enabling" "example" {
  feature_name = unleash_feature.example.name
  project_id   = data.unleash_project.example.project_id
  environment  = "development"
  enabled      = true
  depends_on = [
    unleash_strategy_assignment.example // you can not enable the environment before it has strategies
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment` (String) The environment where the toggle will be enabled
- `feature_name` (String) Feature name to enabled
- `project_id` (String) The unleash project the feature is in

### Optional

- `enabled` (Boolean) Whether the feature is on/off in the provided environment. Default is `true` (on).

### Read-Only

- `id` (String) The ID of this resource.
