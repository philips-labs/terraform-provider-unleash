# Unleash Terraform Provider

- Documentation on [registry.terraform.io](https://registry.terraform.io/providers/philips-labs/unleash/latest/docs)

## Overview

A Terraform provider to provision and manage Unleash admin resources - in early development.
To find out more about Unleash, please visit [getunleash.io/](https://www.getunleash.io/)

## Development Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

Clone repository somewhere *outside* your $GOPATH:

```sh
$ git clone git@github.com:philips-labs/terraform-provider-unleash
$ cd terraform-provider-unleash
$ go build .
```

Copy the resulting binary to the appropriate [plugin directory](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) e.g. `.terraform.d/plugins/darwin_amd64/terraform-provider-unleash`

## Acceptance tests

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, so it requires access to a Unleash server.

You can run the Unleash server locally using Docker, see [docs](https://docs.getunleash.io/deploy/getting_started#option-one---use-docker).

Then, you need to set the following environment variables:
- UNLEASH_API_URL: Usually the URL will be: `http:localhost:4242/api/`
- UNLEASH_AUTH_TOKEN: Login with the default Admin user (username: `admin` / password: `unleash4all`) and generate an `Admin` API token. See how to generate it [here](https://docs.getunleash.io/user_guide/api-token#create-api-token).

```sh
$ make testacc
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to the provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

**Terraform 0.13+**: To install this provider, copy and paste this code into your Terraform configuration. Then, run terraform init.

```terraform
terraform {
  required_providers {
    hsdp = {
      source = "philips-labs/unleash" // not published yet :/
      version = ">= 0.0.1"
    }
  }
}

provider "unleash" {
  api_url    = "http://unleash.api-url.com/api"
  auth_token = "auth-token"
}
```

## Documentation

To generate or update documentation, run `go generate`.

## LICENSE

License is MIT
