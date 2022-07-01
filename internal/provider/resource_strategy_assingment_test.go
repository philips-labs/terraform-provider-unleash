package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/philips-labs/terraform-provider-unleash/utils"
)

func TestAccResourceStrategyAssignment(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStrategyAssignment,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("unleash_strategy_assignment.foo", "feature_name", regexp.MustCompile("^bar")),
					resource.TestMatchResourceAttr("unleash_strategy_assignment.foo", "project_id", regexp.MustCompile("^default$")),
					resource.TestMatchResourceAttr("unleash_strategy_assignment.foo", "environment", regexp.MustCompile("^development$")),
					resource.TestMatchResourceAttr("unleash_strategy_assignment.foo", "strategy_name", regexp.MustCompile("^flexibleRollout$")),
					resource.TestMatchResourceAttr("unleash_strategy_assignment.foo", "parameters.rollout", regexp.MustCompile("^68$")),
					resource.TestMatchResourceAttr("unleash_strategy_assignment.foo", "parameters.stickiness", regexp.MustCompile("^random$")),
					resource.TestMatchResourceAttr("unleash_strategy_assignment.foo", "parameters.groupId", regexp.MustCompile("^toggle$")),
					resource.TestMatchResourceAttr("unleash_strategy_assignment.foo2", "strategy_name", regexp.MustCompile("^userWithId$")),
					resource.TestMatchResourceAttr("unleash_strategy_assignment.foo2", "parameters.userIds", regexp.MustCompile("xyz,bar")),
				),
			},
		},
	})
}

var testAccResourceStrategyAssignment = fmt.Sprintf(`
resource "unleash_feature" "foo" {
  name = "bar%s"
  project_id = "default"
  type = "release"
}
resource "unleash_strategy_assignment" "foo" {
	feature_name  = unleash_feature.foo.name
	project_id    = "default"
	environment   = "development"
	strategy_name = "flexibleRollout"
	parameters = {
	  rollout    = "68"
	  stickiness = "random"
	  groupId    = "toggle"
	}
}
resource "unleash_strategy_assignment" "foo2" {
	feature_name  = unleash_feature.foo.name
	project_id    = "default"
	environment   = "development"
	strategy_name = "userWithId"
	parameters = {
	  userIds    = "xyz,bar"
	}
}
`, utils.RandomString(4))
