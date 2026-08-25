package provider

import (
	"context"
	"fmt"
	"sync"

	"github.com/Unleash/unleash-server-api-go/client"
	"github.com/philips-labs/go-unleash-api/v2/api"
)

type ApiClients struct {
	PhilipsUnleashClient *api.ApiClient
	UnleashClient        *client.APIClient

	contextFieldsMu    sync.Mutex
	contextFieldsCache map[string]bool
}

func (c *ApiClients) GetValidContextNames(ctx context.Context) (map[string]bool, error) {
	c.contextFieldsMu.Lock()
	defer c.contextFieldsMu.Unlock()

	if c.contextFieldsCache != nil {
		return c.contextFieldsCache, nil
	}

	contextFields, _, err := c.UnleashClient.ContextAPI.GetContextFields(ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch context fields from Unleash: %w", err)
	}

	builtins := map[string]bool{
		"appName":       true,
		"currentTime":   true,
		"environment":   true,
		"remoteAddress": true,
		"sessionId":     true,
		"userId":        true,
	}

	c.contextFieldsCache = make(map[string]bool, len(builtins)+len(contextFields))
	for k, v := range builtins {
		c.contextFieldsCache[k] = v
	}
	for _, cf := range contextFields {
		c.contextFieldsCache[cf.Name] = true
	}

	return c.contextFieldsCache, nil
}
