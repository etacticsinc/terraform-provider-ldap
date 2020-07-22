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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the town or city.",
			},
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the country or region code.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a description of the object.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the object.",
				ForceNew:    true,
			},
			"object_class": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The list of classes from which this object is derived.",
			},
			"organizational_unit": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The OU name",
				ForceNew:    true,
			},
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the X.500 path of the OU or container where the new object is created.",
				ForceNew:    true,
			},
			"postal_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the postal code or zip code.",
			},
			"street_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a street address.",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a state or province.",
			},
		},
	}
}

func resourceLdapOrganizationalUnitCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	_, ou := resourceLdapOrganizationalUnitUnmarshal(d)
	if err := client.Add(ou); err != nil {
		return err
	}
	d.SetId(ou.GetDN())
	return resourceLdapOrganizationalUnitRead(d, m)
}

func resourceLdapOrganizationalUnitRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	_, ou := resourceLdapOrganizationalUnitUnmarshal(d)
	if err := client.Search(ou); err != nil {
		return err
	}
	return nil
}

func resourceLdapOrganizationalUnitUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Modify(resourceLdapOrganizationalUnitUnmarshal(d)); err != nil {
		return err
	}
	return resourceLdapOrganizationalUnitRead(d, m)
}

func resourceLdapOrganizationalUnitDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	_, ou := resourceLdapOrganizationalUnitUnmarshal(d)
	if err := client.Delete(ou); err != nil {
		return err
	}
	return nil
}

func resourceLdapOrganizationalUnitUnmarshal(d *schema.ResourceData) (oldOu *OrganizationalUnit, newOu *OrganizationalUnit) {
	properties := map[string]func(*OrganizationalUnit, interface{}){
		"city":        func(ou *OrganizationalUnit, v interface{}) { ou.City = v.(string) },
		"country":     func(ou *OrganizationalUnit, v interface{}) { ou.Country = v.(string) },
		"description": func(ou *OrganizationalUnit, v interface{}) { ou.Description = v.(string) },
		"name":        func(ou *OrganizationalUnit, v interface{}) { ou.Name = v.(string) },
		"object_class": func(ou *OrganizationalUnit, v interface{}) {
			set := v.(*schema.Set)
			if set.Len() > 0 {
				objectClass := make([]string, 0)
				for _, c := range set.List() {
					objectClass = append(objectClass, c.(string))
				}
				ou.ObjectClass = objectClass
			} else {
				ou.ObjectClass = []string{top, organizationalUnit}
				for _, objectClass := range ou.ObjectClass {
					set.Add(objectClass)
				}
			}
		},
		"organizational_unit": func(ou *OrganizationalUnit, v interface{}) { ou.OrganizationalUnit = v.(string) },
		"path":                func(ou *OrganizationalUnit, v interface{}) { ou.Path = v.(string) },
		"postal_code":         func(ou *OrganizationalUnit, v interface{}) { ou.PostalCode = v.(string) },
		"street_address":      func(ou *OrganizationalUnit, v interface{}) { ou.StreetAddress = v.(string) },
		"state":               func(ou *OrganizationalUnit, v interface{}) { ou.State = v.(string) },
	}
	newOu = new(OrganizationalUnit)
	oldOu = new(OrganizationalUnit)
	for property, fn := range properties {
		newVal := d.Get(property)
		fn(newOu, newVal)
		if d.HasChange(property) {
			oldVal, _ := d.GetChange(property)
			fn(oldOu, oldVal)
		} else {
			fn(oldOu, newVal)
		}
	}
	return
}
