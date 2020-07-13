package ldap

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLdapGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceLdapGroupCreate,
		Read:   resourceLdapGroupRead,
		Update: resourceLdapGroupUpdate,
		Delete: resourceLdapGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
		},
	}
}

func resourceLdapGroupCreate(d *schema.ResourceData, m interface{}) error {
	return resourceLdapGroupRead(d, m)
}

func resourceLdapGroupRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceLdapGroupUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceLdapGroupRead(d, m)
}

func resourceLdapGroupDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
