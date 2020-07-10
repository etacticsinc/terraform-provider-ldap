package main

import (
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"strings"
)

type LdapClient struct {
	Server       string
	BindDN       string
	BindPassword string
}

func (c *LdapClient) Add(obj LdapObject) error {

	attributes := obj.Attributes()

	request := ldap.NewAddRequest(obj.DistinguishedName(), []ldap.Control{})

	attributes.ForEach(request.Attribute)

	add := func(conn *ldap.Conn) error { return conn.Add(request) }

	return c.bindThen(add)
}

func (c *LdapClient) Search(obj LdapObject) error {

	attributes := obj.Attributes()

	search := func(conn *ldap.Conn) error {

		objectClass := obj.ObjectClass()

		filters := make([]string, len(objectClass) + 1)

		filters[0] = "CN=" + obj.CommonName()

		for i, class := range objectClass {
			filters[i + 1] = "objectClass=" + ldap.EscapeFilter(class)
		}

		filter := "(&(" + strings.Join(filters, ")(") + "))"

		request := ldap.NewSearchRequest(
			obj.BaseDN(),
			ldap.ScopeWholeSubtree,
			0,
			0,
			0,
			false,
			filter,
			attributes.Keys(),
			[]ldap.Control{})

		result, err := conn.Search(request)

		if err != nil {
			return err
		}

		entries := result.Entries

		if len(entries) == 0 { // Not found
			return errors.New(fmt.Sprintf("Not found. (filter=%s)", filter))
		} else if len(entries) > 1 { // Non-unique (shouldn't be possible)
			return errors.New(fmt.Sprintf("Non-unique result. (filter=%s)", filter))
		}

		m := make(map[string][]string)

		for _, attr  := range entries[0].Attributes {
			m[attr.Name] = attr.Values
		}

		obj.SetAttributes(Attributes{m})

		return nil
	}

	return c.bindThen(search)
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
