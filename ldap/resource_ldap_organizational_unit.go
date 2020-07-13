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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the name of the object.",
				ForceNew:    true,
			},
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the X.500 path of the OU or container where the new object is created.",
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a description of the object.",
			},
			"street_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a street address.",
			},
			"city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the town or city.",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a state or province.",
			},
			"postal_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the postal code or zip code.",
			},
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the country or region code.",
			},
		},
	}
}

func NewOrganizationalUnit(d *schema.ResourceData) OrganizationalUnit {
	return OrganizationalUnit{
		Name:          d.Get("name").(string),
		Path:          d.Get("path").(string),
		Description:   d.Get("description").(string),
		StreetAddress: d.Get("street_address").(string),
		City:          d.Get("city").(string),
		State:         d.Get("state").(string),
		PostalCode:    d.Get("postal_code").(string),
		Country:       d.Get("country").(string),
	}
}

func resourceLdapOrganizationalUnitCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	ou := NewOrganizationalUnit(d)
	if err := client.Add(&ou); err != nil {
		return err
	}
	d.SetId(ou.DN())
	return resourceLdapOrganizationalUnitRead(d, m)
}

func resourceLdapOrganizationalUnitRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	ou := NewOrganizationalUnit(d)
	if err := client.Search(&ou); err != nil {
		return err
	}
	return nil
}

func resourceLdapOrganizationalUnitUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	old := NewOrganizationalUnit(d)
	new := NewOrganizationalUnit(d)
	if d.HasChange("name") {
		oldName, _ := d.GetChange("name")
		old.Name = oldName.(string)
	}
	if d.HasChange("path") {
		oldPath, _ := d.GetChange("path")
		old.Path = oldPath.(string)
	}
	if d.HasChange("description") {
		oldDescription, _ := d.GetChange("description")
		old.Description = oldDescription.(string)
	}
	if d.HasChange("street_address") {
		oldStreetAddress, _ := d.GetChange("street_address")
		old.StreetAddress = oldStreetAddress.(string)
	}
	if d.HasChange("city") {
		oldCity, _ := d.GetChange("city")
		old.City = oldCity.(string)
	}
	if d.HasChange("state") {
		oldState, _ := d.GetChange("state")
		old.State = oldState.(string)
	}
	if d.HasChange("postal_code") {
		oldPostalCode, _ := d.GetChange("postal_code")
		old.PostalCode = oldPostalCode.(string)
	}
	if d.HasChange("country") {
		oldCountry, _ := d.GetChange("country")
		old.Country = oldCountry.(string)
	}
	if err := client.Modify(&old, &new); err != nil {
		return err
	}
	return resourceLdapOrganizationalUnitRead(d, m)
}

func resourceLdapOrganizationalUnitDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	ou := NewOrganizationalUnit(d)
	if err := client.Delete(&ou); err != nil {
		return err
	}
	return nil
}
