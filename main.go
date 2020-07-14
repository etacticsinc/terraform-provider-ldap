package main

import (
	"github.com/etacticsinc/terraform-provider-ldap/ldap"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return ldap.Provider()
		},
	})
}
