package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-labs/terraform-provider-unleash/utils"
)

func TestAccResourceUser(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUser,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("unleash_user.foo", "name", regexp.MustCompile("^bar")),
					resource.TestMatchResourceAttr("unleash_user.foo", "email", regexp.MustCompile("^foo.+@foo.com.br$")),
					resource.TestMatchResourceAttr("unleash_user.foo", "username", regexp.MustCompile("^xyz")),
					resource.TestMatchResourceAttr("unleash_user.foo", "root_role", regexp.MustCompile("^Admin$")),
				),
			},
		},
	})
}

var testAccResourceUser = fmt.Sprintf(`
resource "unleash_user" "foo" {
  name = "bar"
  email = "foo%s@foo.com.br"
  username = "xyz%s"
  root_role = "Admin"
}
`, utils.RandomString(4), utils.RandomString(4))
