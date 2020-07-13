package main

import (
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"ldap/internal"
	"strings"
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
		baseDN := obj.BaseDN()
		filter := internal.Filter(obj.RelativeDN(), obj.Class())
		attributes := obj.Attributes()
		request := ldap.NewSearchRequest(baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, attributes.Keys(), []ldap.Control{})
		result, err := conn.Search(request)
		if err != nil {
			return fmt.Errorf("%v\nbase: %v\nfilter: %v", err, baseDN, filter)
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

func (c *Client) Delete(obj Object) error {
	delete := func(conn *ldap.Conn) error {
		dn := obj.DN()
		request := ldap.NewDelRequest(dn, []ldap.Control{})
		return conn.Del(request)
	}
	return c.bindThen(delete)
}

func (c *Client) Rename(obj Object, newName string) error {
	rename := func(conn *ldap.Conn) error {
		dn := obj.DN()
		rdn := obj.RelativeDN()
		rdnSplit := strings.Split(rdn, "=")
		rdnSplit[1] = newName
		rdn = strings.Join(rdnSplit, "=")
		request := ldap.NewModifyDNRequest(dn, rdn, true, "")
		return conn.ModifyDN(request)
	}
	return c.bindThen(rename)
}

func (c *Client) Modify(old Object, new Object) error {
	modify := func(conn *ldap.Conn) error {
		if old.DN() != new.DN() {
			oldPath := strings.Replace(old.DN(), old.RelativeDN()+",", "", 1)
			newPath := strings.Replace(new.DN(), new.RelativeDN()+",", "", 1)
			if oldPath == newPath {
				newPath = ""
			}
			request := ldap.NewModifyDNRequest(old.DN(), new.RelativeDN(), true, newPath)
			if err := conn.ModifyDN(request); err != nil {
				return err
			}
		}
		oldAttributes := old.Attributes()
		newAttributes := new.Attributes()
		request := ldap.NewModifyRequest(new.DN(), []ldap.Control{})
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
			return conn.Modify(request)
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