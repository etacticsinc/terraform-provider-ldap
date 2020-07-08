package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLdapOrganizationalUnit() *schema.Resource {
	return &schema.Resource{
		Create: resourceLdapOrganizationalUnitCreate,
		Read:   resourceLdapOrganizationalUnitRead,
		Update: resourceLdapOrganizationalUnitUpdate,
		Delete: resourceLdapOrganizationalUnitDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
		},
	}
}

func resourceLdapOrganizationalUnitCreate(d *schema.ResourceData, m interface{}) error {
	return resourceLdapOrganizationalUnitRead(d, m)
}

func resourceLdapOrganizationalUnitRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceLdapOrganizationalUnitUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceLdapOrganizationalUnitRead(d, m)
}

func resourceLdapOrganizationalUnitDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
