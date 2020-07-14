package ldap

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLdapUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceLdapUserCreate,
		Read:   resourceLdapUserRead,
		Update: resourceLdapUserUpdate,
		Delete: resourceLdapUserDelete,
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
			"surname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's last name or surname.",
			},
			"uid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A user ID.",
			},
			"uid_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Contains a number that uniquely identifies a user in an administrative domain.",
			},
			"gid_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Contains an integer value that uniquely identifies a group in an administrative domain.",
			},
			"home_directory": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The home directory for the account.",
			},
		},
	}
}

func NewUser(d *schema.ResourceData) User {
	ou := User{
		City:          d.Get("city").(string),
		Country:       d.Get("country").(string),
		Description:   d.Get("description").(string),
		GidNumber:     d.Get("gid_number").(int),
		HomeDirectory: d.Get("home_directory").(string),
		Name:          d.Get("name").(string),
		Path:          d.Get("path").(string),
		PostalCode:    d.Get("postal_code").(string),
		StreetAddress: d.Get("street_address").(string),
		State:         d.Get("state").(string),
		Surname:       d.Get("surname").(string),
		Uid:           d.Get("uid").(string),
		UidNumber:     d.Get("uid_number").(int),
	}
	if c := d.Get("object_class").(*schema.Set); c != nil && c.Len() > 0 {
		objectClass := make([]string, 0)
		for _, c := range c.List() {
			objectClass = append(objectClass, c.(string))
		}
		ou.ObjectClass = objectClass
	} else {
		ou.ObjectClass = []string{top, person, organizationalPerson, user}
	}
	return ou
}

func resourceLdapUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	ou := NewUser(d)
	if err := client.Add(&ou); err != nil {
		return err
	}
	d.SetId(ou.GetDN())
	return resourceLdapUserRead(d, m)
}

func resourceLdapUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	ou := NewUser(d)
	if err := client.Search(&ou); err != nil {
		return err
	}
	return nil
}

func resourceLdapUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	prev := NewUser(d)
	next := NewUser(d)
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
	if d.HasChange("surname") {
		prevSurname, _ := d.GetChange("surname")
		prev.Surname = prevSurname.(string)
	}
	if d.HasChange("uid") {
		prevUid, _ := d.GetChange("uid")
		prev.Uid = prevUid.(string)
	}
	if d.HasChange("uid_number") {
		prevUidNumber, _ := d.GetChange("uid_number")
		prev.UidNumber = prevUidNumber.(int)
	}
	if d.HasChange("gid_number") {
		prevGidNumber, _ := d.GetChange("gid_number")
		prev.GidNumber = prevGidNumber.(int)
	}
	if d.HasChange("home_directory") {
		prevHomeDirectory, _ := d.GetChange("home_directory")
		prev.HomeDirectory = prevHomeDirectory.(string)
	}
	if err := client.Modify(&prev, &next); err != nil {
		return err
	}
	return resourceLdapUserRead(d, m)
}

func resourceLdapUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	ou := NewUser(d)
	if err := client.Delete(&ou); err != nil {
		return err
	}
	return nil
}
