---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "unleash_feature Resource - terraform-provider-unleash"
subcategory: ""
description: |-
  Provides a resource for managing unleash features.
---

# unleash_feature (Resource)

Provides a resource for managing unleash features.

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

### Read-Only

- `id` (String) The ID of this resource.
