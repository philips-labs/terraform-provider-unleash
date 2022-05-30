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

resource "unleash_feature" "variants_example" {
  name       = "toggle_variants"
  project_id = data.unleash_project.example.project_id
  type       = data.unleash_feature_type.example.type_id
  variant {
    name = "Variant1"
  }
  variant {
    name = "Variant2"
    payload {
      type  = "string"
      value = "foo"
    }
    overrides {
      context_name = "appName"
      values       = ["bar", "xyz"]
    }
    overrides {
      context_name = "environment"
      values       = ["development"]
    }
  }
}
