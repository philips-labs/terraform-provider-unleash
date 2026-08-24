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
	cache := c.contextFieldsCache
	c.contextFieldsMu.Unlock()

	if cache != nil {
		return cache, nil
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

	newCache := make(map[string]bool, len(builtins)+len(contextFields))
	for k, v := range builtins {
		newCache[k] = v
	}
	for _, cf := range contextFields {
		newCache[cf.Name] = true
	}

	c.contextFieldsMu.Lock()
	if c.contextFieldsCache == nil {
		c.contextFieldsCache = newCache
	}
	result := c.contextFieldsCache
	c.contextFieldsMu.Unlock()

	return result, nil
}
