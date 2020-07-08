package main

import (
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"reflect"
)

type LdapClient struct {
	Server       string
	BindDN       string
	BindPassword string
}

//city string,
//country string,
//description string,
//displayName string,
//managedBy string,
//postalCode string,
//state string,
//streetAddress string,
func (c *LdapClient) OrganizationalUnitCreate(name string, path string, attributes map[string][]string) error {

	dn := "CN=" + name + "," + path

	organizationalUnit := ldap.NewAddRequest(dn, []ldap.Control{})

	for key, value := range attributes {
		if len(value) > 0 {
			organizationalUnit.Attribute(key, value)
		}
	}

	organizationalUnit.Attribute("objectClass", []string{"top", "organizationalUnit"})

	organizationalUnit.Attribute("name", []string{name})

	return c.bindThen(func(conn *ldap.Conn) error { return conn.Add(organizationalUnit) })
}

//description string,
//displayName string,
//groupCategory string,
//groupScope string,
//homePage string,
//managedBy string,
//name string,
//path string,
//samAccountName string
func (c *LdapClient) GroupCreate(name string, path string, category string, scope string, attributes map[string][]string) error {

	dn := "CN=" + name + "," + path

	group := ldap.NewAddRequest(dn, []ldap.Control{})

	for key, value := range attributes {
		if len(value) > 0 {
			group.Attribute(key, value)
		}
	}

	if category != "" || scope != "" {
		groupType, err := c.groupType(category, scope)

		if err != nil {
			return err
		}
		group.Attribute("groupType", []string{groupType})
	}

	group.Attribute("objectClass", []string{"top", "group"})

	group.Attribute("name", []string{name})

	// Allow writes on group object
	// https://docs.microsoft.com/en-us/windows/win32/adschema/a-instancetype

	group.Attribute("instanceType", []string{fmt.Sprintf("%d", 0x00000004)})

	return c.bindThen(func(conn *ldap.Conn) error { return conn.Add(group) })
}

func (c *LdapClient) UserCreate(name string, path string, accountExpirationDate string, attributes map[string][]string) error {

}

func (c *LdapClient) bindThen(fn func(*ldap.Conn) error) error {

	// Connect to LDAP server

	conn, err := ldap.DialURL(c.Server)

	if err != nil {
		return err
	}

	defer conn.Close()

	// Perform bind

	if c.BindPassword != "" {
		if err := conn.Bind(c.BindDN, c.BindPassword); err != nil {
			return err
		}
	} else {
		if err := conn.UnauthenticatedBind(c.BindDN); err != nil {
			return err
		}
	}

	return fn(conn)
}

func (c *LdapClient) groupType(category string, scope string) (string, error) {

	// Group Category and Scope are stored as a single bitmask property 'groupType'
	// https://docs.microsoft.com/en-us/windows/win32/adschema/a-grouptype

	categoryMasks := map[string]int{
		"Distribution": 0x00000000,
		"Security":     0x80000000,
	}

	categoryMask := categoryMasks[category]

	if categoryMask == 0 {
		return "", errors.New(fmt.Sprintf("invalid group category %q (expected one of %+q)", category, reflect.ValueOf(categoryMasks).MapKeys()))
	}

	scopeMasks := map[string]int{
		"DomainLocal": 0x00000004,
		"Global":      0x00000002,
		"Universal":   0x00000008,
	}

	scopeMask := scopeMasks[scope]

	if scopeMask == 0 {
		return "", errors.New(fmt.Sprintf("invalid group scope %q (expected one of %+q)", scope, reflect.ValueOf(scopeMasks).MapKeys()))
	}

	return fmt.Sprintf("%d", categoryMask|scopeMask), nil
}
