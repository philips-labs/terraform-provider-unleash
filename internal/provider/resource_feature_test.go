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
					resource.TestMatchResourceAttr("unleash_feature.foo", "project_id", regexp.MustCompile("^default$")),
					resource.TestMatchResourceAttr("unleash_feature.foo", "type", regexp.MustCompile("^release$")),
					resource.TestMatchResourceAttr("unleash_feature.foo", "variant.0.name", regexp.MustCompile("^Variant$")),
					resource.TestMatchResourceAttr("unleash_feature.foo", "variant.0.payload.0.type", regexp.MustCompile("^string$")),
					resource.TestMatchResourceAttr("unleash_feature.foo", "variant.0.payload.0.value", regexp.MustCompile("^foo$")),
					resource.TestMatchResourceAttr("unleash_feature.foo", "variant.0.overrides.0.context_name", regexp.MustCompile("^appName$")),
					resource.TestMatchResourceAttr("unleash_feature.foo", "variant.0.overrides.0.values.0", regexp.MustCompile("^bar$")),
					resource.TestMatchResourceAttr("unleash_feature.foo", "variant.0.overrides.0.values.1", regexp.MustCompile("^xyz$")),
					resource.TestMatchResourceAttr("unleash_feature.foo", "variant.0.overrides.1.context_name", regexp.MustCompile("^environment$")),
					resource.TestMatchResourceAttr("unleash_feature.foo", "variant.0.overrides.1.values.0", regexp.MustCompile("^development$")),
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
