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
