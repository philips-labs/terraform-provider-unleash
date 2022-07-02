package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-labs/terraform-provider-unleash/utils"
)

func TestAccResourceFeature(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFeature,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("unleash_feature.foo", "name", regexp.MustCompile("^bar")),
					resource.TestCheckResourceAttr("unleash_feature.foo", "project_id", "default"),
					resource.TestCheckResourceAttr("unleash_feature.foo", "type", "release"),
					resource.TestCheckResourceAttr("unleash_feature.foo", "variant.0.name", "Variant"),
					resource.TestCheckResourceAttr("unleash_feature.foo", "variant.0.payload.0.type", "string"),
					resource.TestCheckResourceAttr("unleash_feature.foo", "variant.0.payload.0.value", "foo"),
					resource.TestCheckResourceAttr("unleash_feature.foo", "variant.0.overrides.0.context_name", "appName"),
					resource.TestCheckResourceAttr("unleash_feature.foo", "variant.0.overrides.0.values.0", "bar"),
					resource.TestCheckResourceAttr("unleash_feature.foo", "variant.0.overrides.0.values.1", "xyz"),
					resource.TestCheckResourceAttr("unleash_feature.foo", "variant.0.overrides.1.context_name", "environment"),
					resource.TestCheckResourceAttr("unleash_feature.foo", "variant.0.overrides.1.values.0", "development"),
				),
			},
		},
	})
}

var testAccResourceFeature = fmt.Sprintf(`
resource "unleash_feature" "foo" {
  name = "bar%s"
  project_id = "default"
  type = "release"
  variant {
	name = "Variant"
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
`, utils.RandomString(4))
