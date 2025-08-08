# Disruptive Technologies Terraform Provider
This is a Terraform provider for Disruptive Technologies, allowing you to manage devices, data connectors, projects, and notification rules within the Disruptive Technologies platform. This provider is designed to work with the Disruptive Technologies API, enabling infrastructure as code for IoT solutions.

## Features
The provider currently supports the following resources and data sources:

- [x] Device data source
- [x] Data Connector resource
- [ ] Data Connector data source
- [ ] Labels Resource
- [ ] Labels Data Source
- [ ] Organization Data Source
- [x] Project data source
- [x] Project resource
- [x] Contact Group Resource
- [ ] Contact Group Data Source
- [x] Notification Rules Resource
- [ ] Notification Rules Data Source

## Usage

The provider requires a DT service account. Se how to setup a service account [here](https://disruptive.gitbook.io/docs/service-accounts/creating-a-service-account).

The provider requires the following variables to be set:
- `DT_API_KEY_ID` - The ID for the DT Service Account key
- `DT_API_KEY_SECRET` - The secret for the DT Service Account key
- `DT_OIDC_EMAIL` - The email for the DT Service Account

These variables are sensitive and should not be committed to version control.

Here is an example of how to configure the provider:

```hcl
terraform {
  required_providers {
    disruptive-technologies = {
      source = "registry.terraform.io/disruptive-technologies/dt"
    }
  }
}

provider "disruptive-technologies" {
  url            = "https://api.disruptive-technologies.com"
  token_endpoint = "https://identity.disruptive-technologies.com/oauth2/token"
}
```

See the [examples](examples) directory for example usage.
