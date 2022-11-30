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
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.name", "flexibleRollout"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.parameters.rollout", "68"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.parameters.stickiness", "random"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.parameters.groupId", "toggle"),
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
`, utils.RandomString(4))
