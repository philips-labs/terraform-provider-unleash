package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-labs/terraform-provider-unleash/utils"
)

func TestAccResourceFeatureEnabling(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFeatureEnabling,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("unleash_feature_enabling.foo", "feature_name", regexp.MustCompile("^bar")),
					resource.TestCheckResourceAttr("unleash_feature_enabling.foo", "project_id", "default"),
					resource.TestCheckResourceAttr("unleash_feature_enabling.foo", "environment", "development"),
					resource.TestCheckResourceAttr("unleash_feature_enabling.foo", "enabled", "true"),
				),
			},
		},
	})
}

var testAccResourceFeatureEnabling = fmt.Sprintf(`
resource "unleash_feature" "foo" {
  name = "bar%s"
  project_id = "default"
  type = "release"
}
resource "unleash_strategy_assignment" "foo" {
	feature_name  = unleash_feature.foo.name
	project_id    = "default"
	environment   = "development"
	strategy_name = "userWithId"
	parameters = {
	  userIds    = "xyz,bar"
	}
}
resource "unleash_feature_enabling" "foo" {
	feature_name = unleash_feature.foo.name
	project_id   = "default"
	environment  = "development"
	depends_on = [
		unleash_strategy_assignment.foo
	]	
}
`, utils.RandomString(4))
