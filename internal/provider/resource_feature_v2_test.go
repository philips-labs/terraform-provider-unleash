package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-labs/terraform-provider-unleash/utils"
)

func TestAccResourceFeatureV2(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFeatureV2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("unleash_feature_v2.foo", "name", regexp.MustCompile("^my_nice_feature")),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "description", "manages my nice feature"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "type", "release"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "project_id", "default"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "archive_on_destroy", "false"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.0.name", "production"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.0.enabled", "false"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.name", "development"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.enabled", "true"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.0.name", "remoteAddress"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.0.parameters.IPs", "189.434.777.123,host.test.com"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.0.variant.0.name", "Variant"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.name", "flexibleRollout"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.parameters.rollout", "68"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.parameters.stickiness", "random"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.parameters.groupId", "toggle"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.constraint.0.context_name", "appName"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.constraint.0.operator", "NUM_EQ"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.constraint.0.case_insensitive", "false"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.constraint.0.inverted", "false"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.constraint.0.value", "1"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "tag.0.type", "simple"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "tag.0.value", "value"),
				),
			},
		},
	})
}

var testAccResourceFeatureV2 = fmt.Sprintf(`
resource "unleash_feature_v2" "foo" {
	name               = "my_nice_feature_%s"
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
			variant {
				name = "Variant"
				payload {
					type  = "string"
					value = "foo"
				}
			}
	  }
	  strategy {
			name = "flexibleRollout"
			constraint {
				context_name = "appName"
				operator = "NUM_EQ"
				case_insensitive = false
				inverted = false
				value = "1"
			}
			parameters = {
				rollout    = "68"
				stickiness = "random"
				groupId    = "toggle"
			}
	  }
	}
	tag {
		type = "simple"
		value = "value"
	}
}
`, utils.RandomString(4))
