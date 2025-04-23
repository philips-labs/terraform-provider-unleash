package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-labs/terraform-provider-unleash/utils"
)

func TestAccResourceUser(t *testing.T) {
	randomSuffix := utils.RandomString(4)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// Initial creation
				Config: testAccResourceUserInitial(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("unleash_user.foo", "name", regexp.MustCompile("^bar")),
					resource.TestMatchResourceAttr("unleash_user.foo", "email", regexp.MustCompile("^foo.+@foo.com.br$")),
					resource.TestMatchResourceAttr("unleash_user.foo", "root_role", regexp.MustCompile("^Admin$")),
				),
			},
			{
				// Update configuration
				Config: testAccResourceUserUpdated(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify updates took effect
					resource.TestMatchResourceAttr("unleash_user.foo", "root_role", regexp.MustCompile("^Viewer$")),
					// Verify unchanged attributes
					resource.TestMatchResourceAttr("unleash_user.foo", "name", regexp.MustCompile("^bar")),
					resource.TestMatchResourceAttr("unleash_user.foo", "email", regexp.MustCompile("^foo.+@foo.com.br$")),
				),
			},
		},
	})
}

func testAccResourceUserInitial(suffix string) string {
	return fmt.Sprintf(`
resource "unleash_user" "foo" {
  name = "bar"
  email = "foo%s@foo.com.br"
  root_role = "Admin"
}`, suffix)
}

func testAccResourceUserUpdated(suffix string) string {
	return fmt.Sprintf(`
resource "unleash_user" "foo" {
  name = "bar"
  email = "foo%s@foo.com.br"
  root_role = "Viewer"
}`, suffix)
}
