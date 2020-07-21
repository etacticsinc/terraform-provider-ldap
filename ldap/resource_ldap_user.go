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
			"city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the town or city.",
			},
			"common_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name that represents the object. Used to perform searches",
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
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The display name for an object.",
			},
			"email_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user's e-mail address.",
			},
			"gid_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Contains an integer value that uniquely identifies a group in an administrative domain.",
			},
			"given_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Contains the given name (first name) of the user.",
			},
			"home_directory": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The home directory for the account.",
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
			"sam_account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The logon name.",
			},
			"sam_account_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Contains information about every account type object.",
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
			"user_principal_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies a user principal name (UPN) in the format <user>@<DNS-domain-name>.",
			},
		},
	}
}

func resourceLdapUserCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	_, u := resourceLdapUserUnmarshal(d)
	if err := client.Add(u); err != nil {
		return err
	}
	d.SetId(u.GetDN())
	return resourceLdapUserRead(d, m)
}

func resourceLdapUserRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	_, u := resourceLdapUserUnmarshal(d)
	if err := client.Search(u); err != nil {
		return err
	}
	return nil
}

func resourceLdapUserUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Modify(resourceLdapUserUnmarshal(d)); err != nil {
		return err
	}
	return resourceLdapUserRead(d, m)
}

func resourceLdapUserDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	_, u := resourceLdapUserUnmarshal(d)
	if err := client.Delete(u); err != nil {
		return err
	}
	return nil
}

func resourceLdapUserUnmarshal(d *schema.ResourceData) (oldUser *User, newUser *User) {
	properties := map[string]func(*User, interface{}){
		"city":           func(u *User, v interface{}) { u.City = v.(string) },
		"common_name":    func(u *User, v interface{}) { u.CommonName = v.(string) },
		"country":        func(u *User, v interface{}) { u.Country = v.(string) },
		"description":    func(u *User, v interface{}) { u.Description = v.(string) },
		"display_name":   func(u *User, v interface{}) { u.DisplayName = v.(string) },
		"email_address":  func(u *User, v interface{}) { u.EmailAddress = v.(string) },
		"gid_number":     func(u *User, v interface{}) { u.GidNumber = v.(int) },
		"given_name":     func(u *User, v interface{}) { u.GivenName = v.(string) },
		"home_directory": func(u *User, v interface{}) { u.HomeDirectory = v.(string) },
		"name":           func(u *User, v interface{}) { u.Name = v.(string) },
		"object_class": func(u *User, v interface{}) {
			if set := v.(*schema.Set); set != nil && set.Len() > 0 {
				objectClass := make([]string, 0)
				for _, c := range set.List() {
					objectClass = append(objectClass, c.(string))
				}
				u.ObjectClass = objectClass
			} else {
				u.ObjectClass = []string{top, person, organizationalPerson, user}
			}
		},
		"path":                func(u *User, v interface{}) { u.Path = v.(string) },
		"postal_code":         func(u *User, v interface{}) { u.PostalCode = v.(string) },
		"sam_account_name":    func(u *User, v interface{}) { u.SamAccountName = v.(string) },
		"sam_account_type":    func(u *User, v interface{}) { u.SamAccountType = v.(string) },
		"street_address":      func(u *User, v interface{}) { u.StreetAddress = v.(string) },
		"state":               func(u *User, v interface{}) { u.State = v.(string) },
		"surname":             func(u *User, v interface{}) { u.Surname = v.(string) },
		"uid":                 func(u *User, v interface{}) { u.Uid = v.(string) },
		"uid_number":          func(u *User, v interface{}) { u.UidNumber = v.(int) },
		"user_principal_name": func(u *User, v interface{}) { u.UserPrincipalName = v.(string) },
	}
	newUser = new(User)
	for property, fn := range properties {
		fn(newUser, d.Get(property))
		if d.HasChange(property) {
			if oldUser == nil {
				oldUser = new(User)
			}
			oldVal, _ := d.GetChange(property)
			fn(oldUser, oldVal)
		}
	}
	return
}
