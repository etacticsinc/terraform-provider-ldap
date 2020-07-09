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
				Description: "Specifies the name of the object.",
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the display name of the object.",
			},
			"path": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the X.500 path of the OU or container where the new object is created.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a description of the object.",
			},
			"managed_by": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user or group that manages the object.",
			},
			"street_address": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a street address.",
			},
			"city": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the town or city.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a state or province.",
			},
			"postal_code": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the postal code or zip code.",
			},
			"country": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the country or region code.",
			},
		},
	}
}

func resourceLdapOrganizationalUnitCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*LdapClient)

	ou := OrganizationalUnit{
		Name:          d.Get("name").(string),
		DisplayName:   d.Get("display_name").(string),
		Path:          d.Get("path").(string),
		Description:   d.Get("description").(string),
		ManagedBy:     d.Get("managed_by").(string),
		StreetAddress: d.Get("street_address").(string),
		City:          d.Get("city").(string),
		State:         d.Get("state").(string),
		PostalCode:    d.Get("postal_code").(string),
		Country:       d.Get("country").(string),
	}

	if err := client.Add(&ou); err != nil {
		return err
	}

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
