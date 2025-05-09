---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "unleash_feature_type Data Source - terraform-provider-unleash"
subcategory: ""
description: |-
  Retrieve details of an existing feature type
---

# unleash_feature_type (Data Source)

Retrieve details of an existing feature type

## Example Usage

```terraform
data "unleash_feature_type" "example" {
  type_id = "kill-switch"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `type_id` (String) The id of the feature type

### Read-Only

- `description` (String) The description of the feature type
- `id` (String) The ID of this resource.
- `lifetime_days` (Number) The lifetime of the feature type in days
- `name` (String) Feature type name
