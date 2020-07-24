package ldap

import (
	"errors"
	"github.com/etacticsinc/terraform-provider-ldap/ldap/internal"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func resourceLdapUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceLdapUserCreate,
		Read:   resourceLdapUserRead,
		Update: resourceLdapUserUpdate,
		Delete: resourceLdapUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the town or city.",
			},
			"cn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name that represents the object. Used to perform searches",
				ForceNew:    true,
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
				Computed:    true,
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
				Computed:    true,
				Description: "Specifies the Security Account Manager (SAM) account name of the user.",
			},
			"sam_account_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the Security Account Manager (SAM) account type of the user.",
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
				Description: "Specifies a user principal name (UPN) in the format <USER>@<DNS-domain-name>.",
			},
		},
	}
}

func resourceLdapUserCreate(d *schema.ResourceData, m interface{}) error {
	_, u, err := resourceLdapUserUnmarshal(d)
	if err != nil {
		return err
	}
	client := m.(*Client)
	if err := client.Add(u); err != nil {
		return err
	}
	return resourceLdapUserRead(d, m)
}

func resourceLdapUserRead(d *schema.ResourceData, m interface{}) error {
	_, u, err := resourceLdapUserUnmarshal(d)
	if err != nil {
		return err
	}
	client := m.(*Client)
	if err := client.Search(u); err != nil {
		return err
	}
	return resourceLdapUserMarshal(u, d)
}

func resourceLdapUserUpdate(d *schema.ResourceData, m interface{}) error {
	oldUser, newUser, err := resourceLdapUserUnmarshal(d)
	if err != nil {
		return err
	}
	client := m.(*Client)
	if err := client.Modify(oldUser, newUser); err != nil {
		return err
	}
	return resourceLdapUserRead(d, m)
}

func resourceLdapUserDelete(d *schema.ResourceData, m interface{}) error {
	_, u, err := resourceLdapUserUnmarshal(d)
	if err != nil {
		return err
	}
	client := m.(*Client)
	if err := client.Delete(u); err != nil {
		return err
	}
	return nil
}

func resourceLdapUserMarshal(u *User, d *schema.ResourceData) error {
	if d.Id() != u.DN {
		d.SetId(u.DN)
	}
	d.Set("city", u.City)
	d.Set("cn", u.CommonName)
	d.Set("country", u.Country)
	d.Set("description", u.Description)
	d.Set("display_name", u.DisplayName)
	d.Set("email_address", u.EmailAddress)
	d.Set("gid_number", u.GidNumber)
	d.Set("given_name", u.GivenName)
	d.Set("home_directory", u.HomeDirectory)
	d.Set("name", u.Name)
	d.Set("object_class", u.ObjectClass)
	d.Set("path", u.Path)
	d.Set("postal_code", u.PostalCode)
	d.Set("sam_account_name", u.SamAccountName)
	d.Set("sam_account_type", u.SamAccountType)
	d.Set("street_address", u.StreetAddress)
	d.Set("state", u.State)
	d.Set("surname", u.Surname)
	d.Set("uid", u.Uid)
	d.Set("uid_number", u.UidNumber)
	d.Set("user_principal_name", u.UserPrincipalName)
	return nil
}

func resourceLdapUserUnmarshal(d *schema.ResourceData) (oldUser *User, newUser *User, err error) {
	newUser = &User{DN: d.Id()}
	oldUser = &User{DN: d.Id()}
	if _, ok := d.GetOk("path"); !ok { // Not present in import
		rdn, path, err := internal.ParseDN(d.Id())
		if err != nil {
			return oldUser, newUser, err
		}
		if !strings.HasPrefix(strings.ToLower(rdn), "cn=") {
			return nil, nil, errors.New("invalid distinguished name; expected prefix \"cn=\"")
		}
		newUser.CommonName = rdn[3:]
		newUser.Path = path
	} else {
		properties := map[string]func(*User, interface{}){
			"city":           func(u *User, v interface{}) { u.City = v.(string) },
			"cn":             func(u *User, v interface{}) { u.CommonName = v.(string) },
			"country":        func(u *User, v interface{}) { u.Country = v.(string) },
			"description":    func(u *User, v interface{}) { u.Description = v.(string) },
			"display_name":   func(u *User, v interface{}) { u.DisplayName = v.(string) },
			"email_address":  func(u *User, v interface{}) { u.EmailAddress = v.(string) },
			"gid_number":     func(u *User, v interface{}) { u.GidNumber = v.(int) },
			"given_name":     func(u *User, v interface{}) { u.GivenName = v.(string) },
			"home_directory": func(u *User, v interface{}) { u.HomeDirectory = v.(string) },
			"name":           func(u *User, v interface{}) { u.Name = v.(string) },
			"object_class": func(u *User, v interface{}) {
				set := v.(*schema.Set)
				if set.Len() > 0 {
					objectClass := make([]string, 0)
					for _, c := range set.List() {
						objectClass = append(objectClass, c.(string))
					}
					u.ObjectClass = objectClass
				} else {
					u.ObjectClass = []string{top, PERSON, ORGANIZATIONAL_PERSON, USER}
					for _, objectClass := range u.ObjectClass {
						set.Add(objectClass)
					}
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
		for property, fn := range properties {
			newVal := d.Get(property)
			fn(newUser, newVal)
			if d.HasChange(property) {
				oldVal, _ := d.GetChange(property)
				fn(oldUser, oldVal)
			} else {
				fn(oldUser, newVal)
			}
		}
	}
	return
}
