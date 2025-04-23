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
					resource.TestMatchResourceAttr("unleash_user.foo", "username", regexp.MustCompile("^xyz")),
					resource.TestMatchResourceAttr("unleash_user.foo", "root_role", regexp.MustCompile("^Admin$")),
				),
			},
			{
				// Update configuration
				Config: testAccResourceUserUpdated(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify updates took effect
					resource.TestMatchResourceAttr("unleash_user.foo", "root_role", regexp.MustCompile("^Viewer$")),
					resource.TestMatchResourceAttr("unleash_user.foo", "name", regexp.MustCompile("^Joana")),
					resource.TestMatchResourceAttr("unleash_user.foo", "email", regexp.MustCompile("^foo.+@foozzz.com.br$")),
					// Verify unchanged attributes
					resource.TestMatchResourceAttr("unleash_user.foo", "username", regexp.MustCompile("^xyz")),
				),
			},
		},
	})
}

func testAccResourceUserInitial(suffix string) string {
	return fmt.Sprintf(`
resource "unleash_user" "foo" {
  name = "bar"
	username = "xyz%s"
  email = "foo%s@foo.com.br"
  root_role = "Admin"
}`, suffix, suffix)
}

func testAccResourceUserUpdated(suffix string) string {
	return fmt.Sprintf(`
resource "unleash_user" "foo" {
  name = "Joana Darc"
	username = "xyz%s"
  email = "foo%s@foozzz.com.br"
  root_role = "Viewer"
}`, suffix, suffix)
}
