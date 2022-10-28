package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-labs/terraform-provider-unleash/utils"
)

func TestAccResourceApiToken(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApiToken,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("unleash_api_token.foo", "username", regexp.MustCompile("^bar")),
					resource.TestMatchResourceAttr("unleash_api_token.foo", "secret", regexp.MustCompile(`^\*:\*\.`)),
				),
			},
		},
	})
}

var testAccResourceApiToken = fmt.Sprintf(`
resource "unleash_api_token" "foo" {
  username    = "bar%s"
  type        = "admin"
  environment = "*"
  projects    = ["*"]
}
`, utils.RandomString(4))
