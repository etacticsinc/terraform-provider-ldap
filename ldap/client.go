package ldap

import (
	"errors"
	"fmt"
	"github.com/etacticsinc/terraform-provider-ldap/ldap/internal"
	"github.com/go-ldap/ldap/v3"
	"strings"
)

type Client struct {
	Server       string
	BindDN       string
	BindPassword string
}

func (c *Client) Add(obj Object) error {
	add := func(conn *ldap.Conn) error {
		dn := fmt.Sprintf("%s,%s", obj.GetRelativeDN(), obj.GetPath())
		attributes := obj.GetAttributes()
		request := ldap.NewAddRequest(dn, []ldap.Control{})
		attributes.ForEach(request.Attribute)
		if err := conn.Add(request); err != nil {
			return fmt.Errorf("%v\nattributes: %v", err, attributes.String())
		}
		return nil
	}
	return c.bindThen(add)
}

func (c *Client) Search(obj Object) error {
	search := func(conn *ldap.Conn) error {
		path := obj.GetPath()
		filter := internal.Filter(obj.GetRelativeDN(), obj.GetObjectClass())
		attributes := obj.GetAttributes()
		request := ldap.NewSearchRequest(path, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, attributes.Keys(), []ldap.Control{})
		result, err := conn.Search(request)
		if err != nil {
			return fmt.Errorf("%s\nserver: %s\nsearch base: %s\nfilter: %s", err, c.Server, path, filter)
		}
		entries := result.Entries
		if len(entries) == 0 { // Not found
			return errors.New(fmt.Sprintf("Resource not found.\nserver: %s\nsearch base: %s\nfilter: %s", c.Server, path, filter))
		} else if len(entries) > 1 { // Non-unique (shouldn't be possible)
			return errors.New(fmt.Sprintf("Non-unique search result.\nserver: %s\nsearch base: %s\nfilter: %s", c.Server, path, filter))
		}
		m := make(map[string][]string)
		entry := entries[0]
		obj.SetDN(entry.DN)
		for _, attr := range entry.Attributes {
			m[attr.Name] = attr.Values
		}
		obj.SetAttributes(Attributes{m})
		return nil
	}
	return c.bindThen(search)
}

func (c *Client) Delete(obj Object) error {
	delete := func(conn *ldap.Conn) error {
		dn := obj.GetDN()
		request := ldap.NewDelRequest(dn, []ldap.Control{})
		return conn.Del(request)
	}
	return c.bindThen(delete)
}

func (c *Client) Modify(old Object, new Object) error {
	modify := func(conn *ldap.Conn) error {
		if old.GetDN() != new.GetDN() {
			oldPath := strings.Replace(old.GetDN(), old.GetRelativeDN()+",", "", 1)
			newPath := strings.Replace(new.GetDN(), new.GetRelativeDN()+",", "", 1)
			if oldPath == newPath {
				newPath = ""
			}
			request := ldap.NewModifyDNRequest(old.GetDN(), new.GetRelativeDN(), true, newPath)
			if err := conn.ModifyDN(request); err != nil {
				return errors.New(fmt.Sprintf("%vdn: %v\nrdn: %v\npath: %v", err, old.GetDN(), new.GetRelativeDN(), newPath))
			}
		}
		oldAttributes := old.GetAttributes()
		newAttributes := new.GetAttributes()
		request := ldap.NewModifyRequest(new.GetDN(), []ldap.Control{})
		modified := false
		for _, key := range newAttributes.Keys() {
			if newAttributes.HasValue(key) {
				newAttribute := newAttributes.Get(key)
				oldAttribute := oldAttributes.Get(key)
				if oldAttribute == nil {
					request.Add(key, newAttribute)
				} else {
					for i, value := range newAttribute {
						if i >= len(oldAttribute) || oldAttribute[i] != value {
							request.Replace(key, newAttribute)
							modified = true
							break
						}
					}
				}
			}
		}
		for _, key := range oldAttributes.Keys() {
			if oldAttributes.HasValue(key) && !newAttributes.HasValue(key) {
				request.Delete(key, oldAttributes.Get(key))
				modified = true
			}
		}
		if modified {
			err := conn.Modify(request)
			if err != nil {
				return fmt.Errorf("%vattributes: %v", err, newAttributes.String())
			}
		}
		return nil
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
