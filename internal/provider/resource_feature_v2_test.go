package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	openapiclient "github.com/Unleash/unleash-server-api-go/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-labs/terraform-provider-unleash/utils"
)

func TestAccResourceFeatureV2(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFeatureV2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("unleash_feature_v2.foo", "name", regexp.MustCompile("^my_nice_feature")),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "description", "manages my nice feature"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "type", "release"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "project_id", "default"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "archive_on_destroy", "false"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "impression_data", "false"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.0.name", "production"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.0.enabled", "false"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.name", "development"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.enabled", "true"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.0.name", "remoteAddress"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.0.parameters.IPs", "189.434.777.123,host.test.com"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.0.variant.0.name", "a"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.0.variant.1.name", "b"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.0.variant.1.weight", "500"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.0.variant.1.weight_type", "fix"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.0.variant.1.payload.0.type", "string"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.0.variant.1.payload.0.value", "foo"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.name", "flexibleRollout"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.parameters.rollout", "68"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.parameters.stickiness", "random"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.parameters.groupId", "toggle"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.constraint.0.context_name", "appName"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.constraint.0.operator", "NUM_EQ"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.constraint.0.case_insensitive", "false"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.constraint.0.inverted", "false"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.constraint.0.value", "1"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.constraint.1.context_name", "remoteAddress"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.constraint.1.operator", "IN"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.constraint.1.values.0", "dev.philips.com"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "environment.1.strategy.1.constraint.1.values.1", "customer.philips.com"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "tag.0.type", "simple"),
					resource.TestCheckResourceAttr("unleash_feature_v2.foo", "tag.0.value", "value"),
				),
			},
		},
	})
}

var testAccResourceFeatureV2 = fmt.Sprintf(`
resource "unleash_feature_v2" "foo" {
	name               = "my_nice_feature_%s"
	description        = "manages my nice feature"
	type               = "release"
	project_id         = "default"
	archive_on_destroy = false
	impression_data    = false

	environment {
	  name    = "production"
	  enabled = false
	}

	environment {
	  name    = "development"
	  enabled = true

	  strategy {
			name = "remoteAddress"
			parameters = {
				IPs = "189.434.777.123,host.test.com"
			}
			variant {
				name = "a"
			}
			variant {
				name = "b"
				weight = 500
				weight_type = "fix"
				payload {
					type  = "string"
					value = "foo"
				}
			}
	  }
	  strategy {
			name = "flexibleRollout"
			constraint {
				context_name = "appName"
				operator = "NUM_EQ"
				case_insensitive = false
				inverted = false
				value = "1"
			}
			constraint {
				context_name = "remoteAddress"
				operator = "IN"
				values = ["dev.philips.com", "customer.philips.com"]
			}
			parameters = {
				rollout    = "68"
				stickiness = "random"
				groupId    = "toggle"
			}
	  }
	}
	tag {
		type = "simple"
		value = "value"
	}
}
`, utils.RandomString(4))

func TestAccResourceFeatureV2_import(t *testing.T) {
	featureName := fmt.Sprintf("import_feature_%s", utils.RandomString(4))
	resourceName := "unleash_feature_v2.import_test"

	config := fmt.Sprintf(`
resource "unleash_feature_v2" "import_test" {
  name               = "%s"
  description        = "feature to import"
  type               = "release"
  project_id         = "default"
  archive_on_destroy = false

  environment {
    name    = "development"
    enabled = true

    strategy {
      name = "flexibleRollout"
      parameters = {
        rollout    = "100"
        stickiness = "default"
        groupId    = "toggle"
      }
    }
  }

  environment {
    name    = "production"
    enabled = false
  }

  tag {
    type  = "simple"
    value = "import-test"
  }
}`, featureName)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", featureName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           fmt.Sprintf("default/%s", featureName),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"archive_on_destroy"},
			},
		},
	})
}

func TestAccResourceFeatureV2_importInvalidId(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "unleash_feature_v2" "import_test" {
  name               = "placeholder"
  type               = "release"
  project_id         = "default"
  archive_on_destroy = true
}`,
				ResourceName:  "unleash_feature_v2.import_test",
				ImportState:   true,
				ImportStateId: "missing-slash",
				ExpectError:   regexp.MustCompile(`unexpected format of ID`),
			},
		},
	})
}

func TestAccResourceFeatureV2_importWithConstraintsAndVariants(t *testing.T) {
	featureName := fmt.Sprintf("import_complex_%s", utils.RandomString(4))
	resourceName := "unleash_feature_v2.import_complex"

	config := fmt.Sprintf(`
resource "unleash_feature_v2" "import_complex" {
  name               = "%s"
  description        = "import test with constraints and variants"
  type               = "release"
  project_id         = "default"
  archive_on_destroy = false

  environment {
    name    = "development"
    enabled = true

    strategy {
      name = "flexibleRollout"
      parameters = {
        rollout    = "100"
        stickiness = "default"
        groupId    = "toggle"
      }
      constraint {
        context_name = "appName"
        operator     = "IN"
        values       = ["service-a", "service-b"]
      }
      constraint {
        context_name = "userId"
        operator     = "NOT_IN"
        values       = ["blocked-user"]
      }
      variant {
        name        = "config"
        weight_type = "variable"
        payload {
          type  = "json"
          value = "{\"key\":\"value\"}"
        }
      }
    }
  }

  environment {
    name    = "production"
    enabled = false
  }
}`, featureName)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "environment.0.strategy.0.constraint.0.context_name", "appName"),
					resource.TestCheckResourceAttr(resourceName, "environment.0.strategy.0.constraint.0.values.0", "service-a"),
					resource.TestCheckResourceAttr(resourceName, "environment.0.strategy.0.constraint.1.context_name", "userId"),
					resource.TestCheckResourceAttr(resourceName, "environment.0.strategy.0.variant.0.name", "config"),
					resource.TestCheckResourceAttr(resourceName, "environment.0.strategy.0.variant.0.payload.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           fmt.Sprintf("default/%s", featureName),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"archive_on_destroy"},
			},
		},
	})
}

func TestAccResourceFeatureV2_importWithDuplicateTagValues(t *testing.T) {
	featureName := fmt.Sprintf("import_duptag_%s", utils.RandomString(4))
	resourceName := "unleash_feature_v2.import_duptag"
	tagType := fmt.Sprintf("test%s", utils.RandomString(4))

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			createTagType(t, tagType)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "unleash_feature_v2" "import_duptag" {
  name               = "%s"
  description        = "import test with same-value different-type tags"
  type               = "release"
  project_id         = "default"
  archive_on_destroy = false

  environment {
    name    = "development"
    enabled = false
  }

  environment {
    name    = "production"
    enabled = false
  }

  tag {
    type  = "simple"
    value = "shared-value"
  }

  tag {
    type  = "%s"
    value = "shared-value"
  }
}`, featureName, tagType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tag.#", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           fmt.Sprintf("default/%s", featureName),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"archive_on_destroy"},
			},
		},
	})
}

func TestAccResourceFeatureV2_builtinContextFallback(t *testing.T) {
	featureName := fmt.Sprintf("builtin_ctx_feature_%s", utils.RandomString(4))

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "unleash_feature_v2" "builtin_ctx" {
  name               = "%s"
  description        = "test builtin context fallback"
  type               = "release"
  project_id         = "default"
  archive_on_destroy = true

  environment {
    name    = "development"
    enabled = true

    strategy {
      name = "flexibleRollout"
      constraint {
        context_name = "remoteAddress"
        operator     = "IN"
        values       = ["192.168.1.1"]
      }
      parameters = {
        rollout    = "100"
        stickiness = "default"
        groupId    = "toggle"
      }
    }
  }
}`, featureName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("unleash_feature_v2.builtin_ctx", "environment.0.strategy.0.constraint.0.context_name", "remoteAddress"),
				),
			},
		},
	})
}

func TestAccResourceFeatureV2_validCustomContext(t *testing.T) {
	featureName := fmt.Sprintf("valid_custom_ctx_%s", utils.RandomString(4))
	contextName := fmt.Sprintf("testCtx%s", utils.RandomString(4))

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			createContextField(t, contextName)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "unleash_feature_v2" "valid_custom_ctx" {
  name               = "%s"
  description        = "test valid custom context"
  type               = "release"
  project_id         = "default"
  archive_on_destroy = true

  environment {
    name    = "development"
    enabled = true

    strategy {
      name = "flexibleRollout"
      constraint {
        context_name = "%s"
        operator     = "IN"
        values       = ["someValue"]
      }
      parameters = {
        rollout    = "100"
        stickiness = "default"
        groupId    = "toggle"
      }
    }
  }
}`, featureName, contextName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("unleash_feature_v2.valid_custom_ctx", "environment.0.strategy.0.constraint.0.context_name", contextName),
				),
			},
		},
	})
}

func TestAccResourceFeatureV2_invalidCustomContext(t *testing.T) {
	featureName := fmt.Sprintf("invalid_ctx_feature_%s", utils.RandomString(4))

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "unleash_feature_v2" "invalid_ctx" {
  name               = "%s"
  description        = "test invalid context validation"
  type               = "release"
  project_id         = "default"
  archive_on_destroy = true

  environment {
    name    = "development"
    enabled = true

    strategy {
      name = "flexibleRollout"
      constraint {
        context_name = "nonExistentContext"
        operator     = "IN"
        values       = ["someValue"]
      }
      parameters = {
        rollout    = "100"
        stickiness = "default"
        groupId    = "toggle"
      }
    }
  }
}`, featureName),
				ExpectError: regexp.MustCompile(`unknown context fields.*"nonExistentContext"`),
			},
		},
	})
}

func TestAccResourceFeatureV2_apiUnavailableBuiltinAccepted(t *testing.T) {
	featureName := fmt.Sprintf("api_down_builtin_%s", utils.RandomString(4))

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"unleash": func() (*schema.Provider, error) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if strings.Contains(r.URL.Path, "/api/admin/context") {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("{}"))
				}))
				t.Cleanup(srv.Close)

				p := New("test")()
				p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
					unleashConfig := openapiclient.NewConfiguration()
					unleashConfig.Servers = openapiclient.ServerConfigurations{
						openapiclient.ServerConfiguration{URL: srv.URL},
					}
					return &ApiClients{
						UnleashClient: openapiclient.NewAPIClient(unleashConfig),
					}, nil
				}
				return p, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "unleash_feature_v2" "builtin_ok" {
  name               = "%s"
  type               = "release"
  project_id         = "default"
  archive_on_destroy = true

  environment {
    name    = "development"
    enabled = true

    strategy {
      name = "flexibleRollout"
      constraint {
        context_name = "userId"
        operator     = "IN"
        values       = ["user1"]
      }
      parameters = {
        rollout    = "100"
        stickiness = "default"
        groupId    = "toggle"
      }
    }
  }
}`, featureName),
					PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				},
			},
		})
}

func TestAccResourceFeatureV2_apiUnavailableCustomRejected(t *testing.T) {
	featureName := fmt.Sprintf("api_down_custom_%s", utils.RandomString(4))

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"unleash": func() (*schema.Provider, error) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if strings.Contains(r.URL.Path, "/api/admin/context") {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("{}"))
				}))
				t.Cleanup(srv.Close)

				p := New("test")()
				p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
					unleashConfig := openapiclient.NewConfiguration()
					unleashConfig.Servers = openapiclient.ServerConfigurations{
						openapiclient.ServerConfiguration{URL: srv.URL},
					}
					return &ApiClients{
						UnleashClient: openapiclient.NewAPIClient(unleashConfig),
					}, nil
				}
				return p, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "unleash_feature_v2" "custom_fail" {
  name               = "%s"
  type               = "release"
  project_id         = "default"
  archive_on_destroy = true

  environment {
    name    = "development"
    enabled = true

    strategy {
      name = "flexibleRollout"
      constraint {
        context_name = "myCustomField"
        operator     = "IN"
        values       = ["val1"]
      }
      parameters = {
        rollout    = "100"
        stickiness = "default"
        groupId    = "toggle"
      }
    }
  }
}`, featureName),
				ExpectError: regexp.MustCompile(`cannot validate custom context fields.*"myCustomField".*failed to fetch context fields`),
			},
		},
	})
}

func TestResourceFeatureV2ImportState_parsesCompositeId(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceFeatureV2().Schema, map[string]interface{}{})
	d.SetId("my-project/my-feature")

	results, err := resourceFeatureV2ImportState(context.Background(), d, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Id() != "my-feature" {
		t.Errorf("expected ID 'my-feature', got %q", results[0].Id())
	}
	if results[0].Get("project_id") != "my-project" {
		t.Errorf("expected project_id 'my-project', got %q", results[0].Get("project_id"))
	}
}

func TestResourceFeatureV2ImportState_rejectsInvalidId(t *testing.T) {
	cases := []string{"", "no-slash", "/no-project", "no-feature/", "project/feature/extra"}
	for _, id := range cases {
		d := schema.TestResourceDataRaw(t, resourceFeatureV2().Schema, map[string]interface{}{})
		d.SetId(id)

		_, err := resourceFeatureV2ImportState(context.Background(), d, nil)
		if err == nil {
			t.Errorf("expected error for ID %q, got nil", id)
		}
	}
}

func createContextField(t *testing.T, name string) {
	t.Helper()

	apiURL := os.Getenv("UNLEASH_API_URL")
	authToken := os.Getenv("UNLEASH_AUTH_TOKEN")

	body := fmt.Sprintf(`{"name": "%s", "description": "test context field"}`, name)
	req, err := http.NewRequest("POST", strings.TrimRight(apiURL, "/")+"/admin/context", strings.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to create context field: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		t.Fatalf("failed to create context field %q: status %d", name, resp.StatusCode)
	}

	t.Cleanup(func() {
		delReq, err := http.NewRequest("DELETE", strings.TrimRight(apiURL, "/")+"/admin/context/"+name, nil)
		if err != nil {
			t.Logf("failed to build delete request for context field %q: %v", name, err)
			return
		}
		delReq.Header.Set("Authorization", authToken)

		delResp, err := http.DefaultClient.Do(delReq)
		if err != nil {
			t.Logf("failed to delete context field %q: %v", name, err)
			return
		}
		delResp.Body.Close()
	})
}

func createTagType(t *testing.T, name string) {
	t.Helper()

	apiURL := os.Getenv("UNLEASH_API_URL")
	authToken := os.Getenv("UNLEASH_AUTH_TOKEN")

	body := fmt.Sprintf(`{"name": "%s", "description": "test tag type"}`, name)
	req, err := http.NewRequest("POST", strings.TrimRight(apiURL, "/")+"/admin/tag-types", strings.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to create tag type: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		t.Fatalf("failed to create tag type %q: status %d", name, resp.StatusCode)
	}

	t.Cleanup(func() {
		delReq, err := http.NewRequest("DELETE", strings.TrimRight(apiURL, "/")+"/admin/tag-types/"+name, nil)
		if err != nil {
			t.Logf("failed to build delete request for tag type %q: %v", name, err)
			return
		}
		delReq.Header.Set("Authorization", authToken)

		delResp, err := http.DefaultClient.Do(delReq)
		if err != nil {
			t.Logf("failed to delete tag type %q: %v", name, err)
			return
		}
		delResp.Body.Close()
	})
}
