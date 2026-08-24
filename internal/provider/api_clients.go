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

	contextFieldsOnce  sync.Once
	contextFieldsCache map[string]bool
	contextFieldsErr   error
}

func (c *ApiClients) GetValidContextNames(ctx context.Context) (map[string]bool, error) {
	c.contextFieldsOnce.Do(func() {
		c.contextFieldsCache = map[string]bool{
			"appName":       true,
			"currentTime":   true,
			"environment":   true,
			"remoteAddress": true,
			"sessionId":     true,
			"userId":        true,
		}

		contextFields, _, err := c.UnleashClient.ContextAPI.GetContextFields(ctx).Execute()
		if err != nil {
			c.contextFieldsErr = fmt.Errorf("failed to fetch context fields from Unleash: %w", err)
			return
		}

		for _, cf := range contextFields {
			c.contextFieldsCache[cf.Name] = true
		}
	})

	return c.contextFieldsCache, c.contextFieldsErr
}
