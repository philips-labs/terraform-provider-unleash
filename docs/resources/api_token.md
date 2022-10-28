---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "unleash_api_token Resource - terraform-provider-unleash"
subcategory: ""
description: |-
  Provides a resource for managing unleash api tokens.
---

# unleash_api_token (Resource)

Provides a resource for managing unleash api tokens.

## Example Usage

```terraform
resource "unleash_api_token" "my_token" {
  username    = "bobjoe"
  type        = "admin"
  expires_at  = "2024-10-19"
  environment = "*"
  projects    = ["*"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `type` (String) The type of the API token. Can be `client`, `admin` or `frontend`
- `username` (String)

### Optional

- `created_at` (String) The API token creation date.
- `environment` (String) The environment the token will have access to. Use `"*"` for all environments. By default, it will have access to all environments.
- `expires_at` (String) The API token expiration date.
- `projects` (Set of String) The project(s) the token will have access to. Use `["*"]` for all projects. By default, it will have access to all projects.
- `secret` (String, Sensitive) The API token secret.

### Read-Only

- `id` (String) The ID of this resource.

