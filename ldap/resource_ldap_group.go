package main

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
			"sam_account_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the Security Account Manager (SAM) account name of the user, group, computer, or service account.",
			},
		},
	}
}

func NewGroup(d *schema.ResourceData) Group {
	g := Group{
		Name:           d.Get("name").(string),
		Path:           d.Get("path").(string),
		Description:    d.Get("description").(string),
		GidNumber:      d.Get("gid_number").(int),
		GroupCategory:  d.Get("group_category").(string),
		GroupScope:     d.Get("group_scope").(string),
		HomePage:       d.Get("homepage").(string),
		SamAccountName: d.Get("sam_account_name").(string),
	}
	if m := d.Get("members").(*schema.Set); m != nil && m.Len() > 0 {
		members := make([]string, 0)
		for _, m := range m.List() {
			members = append(members, m.(string))
		}
		g.Members = members
	}
	if m := d.Get("member_uids").(*schema.Set); m != nil && m.Len() > 0 {
		memberUids := make([]string, 0)
		for _, m := range m.List() {
			memberUids = append(memberUids, m.(string))
		}
		g.MemberUids = memberUids
	}
	if c := d.Get("object_class").(*schema.Set); c != nil && c.Len() > 0 {
		objectClass := make([]string, 0)
		for _, c := range c.List() {
			objectClass = append(objectClass, c.(string))
		}
		g.ObjectClass = objectClass
	} else {
		g.ObjectClass = []string{top, group}
	}
	return g
}

func resourceLdapGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	g := NewGroup(d)
	if err := client.Add(&g); err != nil {
		return err
	}
	d.SetId(g.GetDN())
	return resourceLdapGroupRead(d, m)
}

func resourceLdapGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	g := NewGroup(d)
	if err := client.Search(&g); err != nil {
		return err
	}
	return nil
}

func resourceLdapGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	prev := NewGroup(d)
	next := NewGroup(d)
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
	if d.HasChange("members") {
		prevMembers := make([]string, 0)
		m, _ := d.GetChange("members")
		if m := m.(*schema.Set); m.Len() > 0 {
			for _, m := range m.List() {
				prevMembers = append(prevMembers, m.(string))
			}
		}
		prev.Members = prevMembers
	}
	if d.HasChange("member_uids") {
		prevMemberUids := make([]string, 0)
		m, _ := d.GetChange("member_uids")
		if m := m.(*schema.Set); m.Len() > 0 {
			for _, m := range m.List() {
				prevMemberUids = append(prevMemberUids, m.(string))
			}
		}
		prev.MemberUids = prevMemberUids
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
	if d.HasChange("gid_number") {
		prevGidNumber, _ := d.GetChange("gid_number")
		prev.GidNumber = prevGidNumber.(int)
	}
	if d.HasChange("group_category") {
		prevGroupCategory, _ := d.GetChange("group_category")
		prev.GroupCategory = prevGroupCategory.(string)
	}
	if d.HasChange("group_scope") {
		prevGroupScope, _ := d.GetChange("group_scope")
		prev.GroupScope = prevGroupScope.(string)
	}
	if d.HasChange("homepage") {
		prevHomePage, _ := d.GetChange("homepage")
		prev.HomePage = prevHomePage.(string)
	}
	if d.HasChange("sam_account_name") {
		prevSamAccountName, _ := d.GetChange("sam_account_name")
		prev.SamAccountName = prevSamAccountName.(string)
	}
	if err := client.Modify(&prev, &next); err != nil {
		return err
	}
	return resourceLdapGroupRead(d, m)
}

func resourceLdapGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	g := NewGroup(d)
	if err := client.Delete(&g); err != nil {
		return err
	}
	return nil
}
