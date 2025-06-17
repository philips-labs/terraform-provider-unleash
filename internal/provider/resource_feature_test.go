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
}
`, utils.RandomString(4))
