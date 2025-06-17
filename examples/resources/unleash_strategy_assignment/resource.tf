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
  variant {
    name = "Variant1"
  }
  variant {
    name = "Variant2"
    payload {
      type  = "string"
      value = "foo"
    }
  }
}
