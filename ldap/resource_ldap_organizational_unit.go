package ldap

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
			"object_class": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The list of classes from which this object is derived.",
			},
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
	ou := OrganizationalUnit{
		Name:          d.Get("name").(string),
		Path:          d.Get("path").(string),
		Description:   d.Get("description").(string),
		StreetAddress: d.Get("street_address").(string),
		City:          d.Get("city").(string),
		State:         d.Get("state").(string),
		PostalCode:    d.Get("postal_code").(string),
		Country:       d.Get("country").(string),
	}
	if c := d.Get("object_class").(*schema.Set); c != nil && c.Len() > 0 {
		objectClass := make([]string, 0)
		for _, c := range c.List() {
			objectClass = append(objectClass, c.(string))
		}
		ou.ObjectClass = objectClass
	} else {
		ou.ObjectClass = []string{top, organizationalUnit}
	}
	return ou
}

func resourceLdapOrganizationalUnitCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	ou := NewOrganizationalUnit(d)
	if err := client.Add(&ou); err != nil {
		return err
	}
	d.SetId(ou.GetDN())
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
	prev := NewOrganizationalUnit(d)
	next := NewOrganizationalUnit(d)
	if d.HasChange("object_class") {
		prevObjectClass := make([]string, 0)
		c, _ := d.GetChange("object_class")
		if c := c.(*schema.Set); c.Len() > 0 {
			for _, c := range c.List() {
				prevObjectClass = append(prevObjectClass, c.(string))
			}
		}
		prev.ObjectClass = prevObjectClass
	}
	if d.HasChange("name") {
		prevName, _ := d.GetChange("name")
		prev.Name = prevName.(string)
	}
	if d.HasChange("path") {
		prevPath, _ := d.GetChange("path")
		prev.Path = prevPath.(string)
	}
	if d.HasChange("description") {
		prevDescription, _ := d.GetChange("description")
		prev.Description = prevDescription.(string)
	}
	if d.HasChange("street_address") {
		prevStreetAddress, _ := d.GetChange("street_address")
		prev.StreetAddress = prevStreetAddress.(string)
	}
	if d.HasChange("city") {
		prevCity, _ := d.GetChange("city")
		prev.City = prevCity.(string)
	}
	if d.HasChange("state") {
		prevState, _ := d.GetChange("state")
		prev.State = prevState.(string)
	}
	if d.HasChange("postal_code") {
		prevPostalCode, _ := d.GetChange("postal_code")
		prev.PostalCode = prevPostalCode.(string)
	}
	if d.HasChange("country") {
		prevCountry, _ := d.GetChange("country")
		prev.Country = prevCountry.(string)
	}
	if err := client.Modify(&prev, &next); err != nil {
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
