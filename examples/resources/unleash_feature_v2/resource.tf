resource "unleash_feature_v2" "strategies_example" {
  name               = "toggle_strategies"
  description        = "description"
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
    }
    strategy {
      name = "flexibleRollout"
      parameters = {
        rollout    = "68"
        stickiness = "random"
        groupId    = "toggle"
      }
    }
  }
}
