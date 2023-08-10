package provider

import (
	"github.com/Unleash/unleash-server-api-go/client"
	"github.com/philips-labs/go-unleash-api/api"
)

type ApiClients struct {
	PhilipsUnleashClient *api.ApiClient
	UnleashClient        *client.APIClient
}
