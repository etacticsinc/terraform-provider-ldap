package ldap

import (
	"./internal"
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

type Client struct {
	Server       string
	BindDN       string
	BindPassword string
}

func (c *Client) Add(obj Object) error {
	add := func(conn *ldap.Conn) error {
		attributes := obj.Attributes()
		request := ldap.NewAddRequest(obj.DN(), []ldap.Control{})
		attributes.ForEach(request.Attribute)
		return conn.Add(request)
	}
	return c.bindThen(add)
}

func (c *Client) Search(obj Object) error {
	search := func(conn *ldap.Conn) error {
		filter := internal.Filter(obj.RelativeDN(), obj.Class())
		attributes := obj.Attributes()
		request := ldap.NewSearchRequest(obj.BaseDN(), ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, attributes.Keys(), []ldap.Control{})
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
		for _, attr := range entries[0].Attributes {
			m[attr.Name] = attr.Values
		}
		obj.SetAttributes(Attributes{m})
		return nil
	}
	return c.bindThen(search)
}

func (c *Client) Delete(name string, path string) error {
	delete := func(conn *ldap.Conn) error {
		dn := internal.DN(name, path)
		request := ldap.NewDelRequest(dn, []ldap.Control{})
		return conn.Del(request)
	}
	return c.bindThen(delete)
}

func (c *Client) Rename(name string, path string, newName string) error {
	rename := func(conn *ldap.Conn) error {
		dn := internal.DN(name, path)
		rdn := internal.RelativeDN(newName)
		request := ldap.NewModifyDNRequest(dn, rdn, true, "")
		return conn.ModifyDN(request)
	}
	return c.bindThen(rename)
}

func (c *Client) Move(name string, path string, newPath string) error {
	move := func(conn *ldap.Conn) error {
		dn := internal.DN(name, path)
		rdn := internal.RelativeDN(name)
		request := ldap.NewModifyDNRequest(dn, rdn, true, newPath)
		return conn.ModifyDN(request)
	}
	return c.bindThen(move)
}

func (c *Client) Modify(name string, path string, oldAttributes map[string][]string, newAttributes map[string][]string) error {
	modify := func(conn *ldap.Conn) error {
		dn := internal.DN(name, path)
		request := ldap.NewModifyRequest(dn, []ldap.Control{})
		for key, newAttribute := range newAttributes {
			if newAttribute == nil || len(newAttribute) == 0 || newAttribute[0] == "" {
				if oldAttribute, ok := oldAttributes[key]; ok && oldAttribute != nil {
					request.Delete(key, oldAttribute)
				} else {
					return errors.New(fmt.Sprintf("cannot delete attribute '%s,' no initial value specified", key))
				}
			} else {
				if oldAttribute, ok := oldAttributes[key]; ok && oldAttribute != nil {
					request.Replace(key, newAttribute)
				} else {
					request.Add(key, newAttribute)
				}
			}
		}
		return conn.Modify(request)
	}
	return c.bindThen(modify)
}

func (c *Client) bindThen(fn func(*ldap.Conn) error) error {
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
