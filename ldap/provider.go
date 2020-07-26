package ldap

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net"
	url2 "net/url"
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
				Sensitive: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ldap_organizational_unit": resourceLdapOrganizationalUnit(),
			"ldap_user":                resourceLdapUser(),
			"ldap_group":               resourceLdapGroup(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	url, err := url2.Parse(d.Get("server").(string))
	if err != nil {
		return nil, err
	}
	if ip := net.ParseIP(url.Host); ip == nil {
		ips, err := net.LookupIP(url.Host)
		if err != nil {
			return nil, err
		}
		url.Host = ips[0].String()
	}
	config := Config{
		Server:       url.String(),
		BindDN:       d.Get("bind_dn").(string),
		BindPassword: d.Get("bind_password").(string),
	}
	return config.Client()
}
