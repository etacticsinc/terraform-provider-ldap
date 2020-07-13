package ldap

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"bind_dn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"bind_password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ldap_organizational_unit": resourceLdapOrganizationalUnit(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Server:       d.Get("server").(string),
		BindDN:       d.Get("bind_dn").(string),
		BindPassword: d.Get("bind_password").(string),
	}
	return config.Client()
}
