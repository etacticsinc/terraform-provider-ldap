package ldap

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceLdapGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceLdapGroupCreate,
		Read:   resourceLdapGroupRead,
		Update: resourceLdapGroupUpdate,
		Delete: resourceLdapGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"common_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name that represents the object. Used to perform searches",
				ForceNew:    true,
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
			"gid_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Contains an integer value that uniquely identifies a group in an administrative domain.",
			},
			"group_category": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  fmt.Sprintf("Specifies the category of the group. The acceptable values for this parameter are \"%s\" and \"%s\"", distribution, security),
				ValidateFunc: validation.StringInSlice([]string{distribution, security}, false),
			},
			"group_scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  fmt.Sprintf("Specifies the group scope of the group. The acceptable values for this parameter are \"%s,\" \"%s\" and \"%s\"", global, domainLocal, universal),
				ValidateFunc: validation.StringInSlice([]string{global, domainLocal, universal}, false),
			},
			"homepage": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the URL of the home page of the object.",
			},
			"members": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				ConflictsWith: []string{"member_uids"},
				Description:   "Specifies an array of user, group, and computer objects to add to the group.",
			},
			"member_uids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				ConflictsWith: []string{"members"},
				Description:   "Contains the login names of the members of a group.",
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
			"sam_account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the Security Account Manager (SAM) account name of the group.",
			},
			"sam_account_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the Security Account Manager (SAM) account type of the group.",
			},
		},
	}
}

func resourceLdapGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	_, g := resourceLdapGroupUnmarshal(d)
	if err := client.Add(g); err != nil {
		return err
	}
	d.SetId(g.GetDN())
	return resourceLdapGroupRead(d, m)
}

func resourceLdapGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	_, g := resourceLdapGroupUnmarshal(d)
	if err := client.Search(g); err != nil {
		return err
	}
	return nil
}

func resourceLdapGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	if err := client.Modify(resourceLdapGroupUnmarshal(d)); err != nil {
		return err
	}
	return resourceLdapGroupRead(d, m)
}

func resourceLdapGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	_, g := resourceLdapGroupUnmarshal(d)
	if err := client.Delete(g); err != nil {
		return err
	}
	return nil
}

func resourceLdapGroupUnmarshal(d *schema.ResourceData) (oldGroup *Group, newGroup *Group) {
	properties := map[string]func(*Group, interface{}){
		"common_name":    func(g *Group, v interface{}) { g.CommonName = v.(string) },
		"description":    func(g *Group, v interface{}) { g.Description = v.(string) },
		"display_name":   func(g *Group, v interface{}) { g.DisplayName = v.(string) },
		"gid_number":     func(g *Group, v interface{}) { g.GidNumber = v.(int) },
		"group_category": func(g *Group, v interface{}) { g.GroupCategory = v.(string) },
		"group_scope":    func(g *Group, v interface{}) { g.GroupScope = v.(string) },
		"homepage":       func(g *Group, v interface{}) { g.HomePage = v.(string) },
		"members": func(g *Group, v interface{}) {
			set := v.(*schema.Set)
			members := make([]string, 0)
			for _, m := range set.List() {
				members = append(members, m.(string))
			}
			g.Members = members
		},
		"member_uids": func(g *Group, v interface{}) {
			set := v.(*schema.Set)
			memberUids := make([]string, 0)
			for _, m := range set.List() {
				memberUids = append(memberUids, m.(string))
			}
			g.MemberUids = memberUids
		},
		"name": func(g *Group, v interface{}) { g.Name = v.(string) },
		"object_class": func(g *Group, v interface{}) {
			set := v.(*schema.Set)
			if set.Len() > 0 {
				objectClass := make([]string, 0)
				for _, c := range set.List() {
					objectClass = append(objectClass, c.(string))
				}
				g.ObjectClass = objectClass
			} else {
				g.ObjectClass = []string{top, group}
				for _, objectClass := range g.ObjectClass {
					set.Add(objectClass)
				}
			}
		},
		"path":             func(g *Group, v interface{}) { g.Path = v.(string) },
		"sam_account_name": func(g *Group, v interface{}) { g.SamAccountName = v.(string) },
		"sam_account_type": func(g *Group, v interface{}) { g.SamAccountType = v.(string) },
	}
	newGroup = new(Group)
	oldGroup = new(Group)
	for property, fn := range properties {
		newVal := d.Get(property)
		fn(newGroup, newVal)
		if d.HasChange(property) {
			oldVal, _ := d.GetChange(property)
			fn(oldGroup, oldVal)
		} else {
			fn(oldGroup, newVal)
		}
	}
	return
}
