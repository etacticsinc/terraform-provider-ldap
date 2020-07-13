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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the name of the object.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the display name of the object.",
			},
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the X.500 path of the OU or container where the new object is created.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a description of the object.",
			},
			"managed_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user or group that manages the object.",
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
		DisplayName:   d.Get("display_name").(string),
		Path:          d.Get("path").(string),
		Description:   d.Get("description").(string),
		ManagedBy:     d.Get("managed_by").(string),
		StreetAddress: d.Get("street_address").(string),
		City:          d.Get("city").(string),
		State:         d.Get("state").(string),
		ObjectGUID:    d.Id(),
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
	d.SetId(ou.ObjectGUID)
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
	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")
		path := d.Get("path")
		if d.HasChange("path") {
			path, _ = d.GetChange("path")
		}
		if err := client.Rename(oldName.(string), path.(string), newName.(string)); err != nil {
			return err
		}
	}
	if d.HasChange("path") {
		oldPath, newPath := d.GetChange("path")
		name := d.Get("name")
		if err := client.Move(name.(string), oldPath.(string), newPath.(string)); err != nil {
			return err
		}
	}
	oldAttributes, newAttributes := make(map[string][]string), make(map[string][]string)
	for _, key := range []string{"display_name, description", "managed_by", "street_address", "city", "state", "postal_code", "country"} {
		if d.HasChange(key) {
			oldAttribute, newAttribute := d.GetChange(key)
			oldAttributes[key] = []string{oldAttribute.(string)}
			newAttributes[key] = []string{newAttribute.(string)}
		}
	}
	if len(newAttributes) > 0 {
		name := d.Get("name")
		path := d.Get("path")
		if err := client.Modify(name.(string), path.(string), oldAttributes, newAttributes); err != nil {
			return err
		}
	}
	return resourceLdapOrganizationalUnitRead(d, m)
}

func resourceLdapOrganizationalUnitDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	ou := NewOrganizationalUnit(d)
	if err := client.Delete(ou.Name, ou.Path); err != nil {
		return err
	}
	return nil
}
