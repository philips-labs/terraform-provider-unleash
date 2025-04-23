package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-labs/terraform-provider-unleash/utils"
)

func TestAccResourceApiToken(t *testing.T) {
	randomSuffix := utils.RandomString(4)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				// Initial creation
				Config: testAccResourceApiTokenInitial(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("unleash_api_token.foo", "token_name", regexp.MustCompile("^bar")),
					resource.TestMatchResourceAttr("unleash_api_token.foo", "secret", regexp.MustCompile(`^\*:development\.`)),
					resource.TestCheckResourceAttr("unleash_api_token.foo", "environment", "development"),
					resource.TestCheckResourceAttr("unleash_api_token.foo", "projects.#", "1"),
				),
			},
			{
				// Update configuration
				Config: testAccResourceApiTokenUpdated(randomSuffix),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify updates took effect
					resource.TestMatchResourceAttr("unleash_api_token.foo", "expires_at", regexp.MustCompile("^2023-04-15T14:30:45Z$")),
					// Verify unchanged attributes
					resource.TestMatchResourceAttr("unleash_api_token.foo", "token_name", regexp.MustCompile("^bar")),
					resource.TestMatchResourceAttr("unleash_api_token.foo", "secret", regexp.MustCompile(`^\*:development\.`)),
					resource.TestCheckResourceAttr("unleash_api_token.foo", "environment", "development"),
					resource.TestCheckResourceAttr("unleash_api_token.foo", "projects.#", "1"),
				),
			},
		},
	})
}

func testAccResourceApiTokenInitial(suffix string) string {
	return fmt.Sprintf(`
resource "unleash_api_token" "foo" {
  token_name    = "bar%s"
  type        = "client"
  environment = "development"
  projects    = ["*"]
}`, suffix)
}

func testAccResourceApiTokenUpdated(suffix string) string {
	return fmt.Sprintf(`
resource "unleash_api_token" "foo" {
  token_name    = "bar%s"
  type        = "client"
  environment = "development"
  projects    = ["*"]
	expires_at = "2023-04-15T14:30:45Z"
}`, suffix)
}
