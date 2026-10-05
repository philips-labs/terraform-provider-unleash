package provider

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) {
	// Ryuk (testcontainers' cleanup sidecar) cannot mount the Windows named pipe into a
	// Linux container. Disable it — TestMain handles explicit cleanup instead.
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	// Bypass: if UNLEASH_API_URL is already set (e.g. developer ran docker compose up -d),
	// skip container startup entirely.
	if os.Getenv("UNLEASH_API_URL") != "" {
		os.Exit(m.Run())
	}

	ctx := context.Background()

	nw, err := network.New(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create docker network: %s\n", err)
		os.Exit(1)
	}

	// Postgres — pre-configured module; alias "postgres" so unleash can reach it by hostname.
	pgContainer, err := tcpostgres.Run(
		ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("unleash"),
		tcpostgres.WithUsername("unleash_user"),
		tcpostgres.WithPassword("some_password"),
		network.WithNetwork([]string{"postgres"}, nw),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start postgres: %s\n", err)
		_ = nw.Remove(ctx)
		os.Exit(1)
	}

	// Unleash — no pre-built module exists; HTTP health check is the correct wait
	// strategy because the port opens before the DB migration finishes.
	unleashContainer, err := testcontainers.Run(
		ctx,
		"unleashorg/unleash-server:latest",
		testcontainers.WithExposedPorts("4242/tcp"),
		testcontainers.WithEnv(map[string]string{
			"DATABASE_HOST":     "postgres",
			"DATABASE_NAME":     "unleash",
			"DATABASE_USERNAME": "unleash_user",
			"DATABASE_PASSWORD": "some_password",
			"DATABASE_SSL":      "false",
			"AUTH_TYPE":         "NONE",
		}),
		network.WithNetwork([]string{"unleash"}, nw),
		testcontainers.WithWaitStrategy(
			wait.ForHTTP("/health").
				WithPort("4242/tcp").
				WithStartupTimeout(120*time.Second),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start unleash: %s\n", err)
		_ = testcontainers.TerminateContainer(pgContainer)
		_ = nw.Remove(ctx)
		os.Exit(1)
	}

	mappedPort, err := unleashContainer.MappedPort(ctx, "4242/tcp")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get unleash port: %s\n", err)
		_ = testcontainers.TerminateContainer(unleashContainer)
		_ = testcontainers.TerminateContainer(pgContainer)
		_ = nw.Remove(ctx)
		os.Exit(1)
	}
	host, err := unleashContainer.Host(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get unleash host: %s\n", err)
		_ = testcontainers.TerminateContainer(unleashContainer)
		_ = testcontainers.TerminateContainer(pgContainer)
		_ = nw.Remove(ctx)
		os.Exit(1)
	}

	os.Setenv("UNLEASH_API_URL", fmt.Sprintf("http://%s:%s/api/", host, mappedPort.Port()))
	os.Setenv("UNLEASH_AUTH_TOKEN", "token")

	code := m.Run()

	_ = testcontainers.TerminateContainer(unleashContainer)
	_ = testcontainers.TerminateContainer(pgContainer)
	_ = nw.Remove(ctx)

	os.Exit(code)
}

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"unleash": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
	if v := os.Getenv("UNLEASH_API_URL"); v == "" {
		t.Fatal("UNLEASH_API_URL must be set for acceptance tests")
	}
	if v := os.Getenv("UNLEASH_AUTH_TOKEN"); v == "" {
		t.Fatal("UNLEASH_AUTH_TOKEN must be set for acceptance tests")
	}
}
