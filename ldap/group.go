package main

import (
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

const (
	distribution = "Distribution"
	security     = "Security"
	domainLocal  = "DomainLocal"
	global       = "Global"
	universal    = "Universal"
)

type Group struct {
	Description    string
	DisplayName    string
	GroupCategory  string
	GroupScope     string
	HomePage       string
	ManagedBy      string
	Name           string
	Path           string
	SamAccountName string
}

func (g *Group) DistinguishedName() (string, error) {
	if g.Name == "" {
		return "", errors.New("undefined group name")
	}

	if g.Path == "" {
		return "", errors.New("undefined group path")
	}

	return fmt.Sprintf("CN=%s,%s", ldap.EscapeFilter(g.Name), g.Path), nil
}

func (g *Group) Attributes() (map[string][]string, error) {

	// Group Category and Scope are stored as a single bitmask property 'groupType'
	// https://docs.microsoft.com/en-us/windows/win32/adschema/a-grouptype

	groupTypeMask := 0

	if categoryMask, ok := g.categoryMasks()[g.GroupCategory]; ok {
		groupTypeMask |= categoryMask
	} else if g.GroupCategory != "" {
		return nil, errors.New(fmt.Sprintf("invalid group category %q", g.GroupCategory))
	}

	if scopeMask, ok := g.scopeMasks()[g.GroupScope]; ok {
		groupTypeMask |= scopeMask
	} else if g.GroupScope != "" {
		return nil, errors.New(fmt.Sprintf("invalid group scope %q", g.GroupScope))
	}

	groupType := fmt.Sprintf("%d", groupTypeMask)

	return map[string][]string{
		"objectClass":    {"top", "group"},
		"description":    {g.Description},
		"displayName":    {g.DisplayName},
		"groupType":      {groupType},
		"wWWHomePage":    {g.HomePage},
		"managedBy":      {g.ManagedBy},
		"name":           {g.Name},
		"sAMAccountName": {g.SamAccountName},
	}, nil
}

func (g *Group) categoryMasks() map[string]int {
	return map[string]int{
		distribution: 0x00000000,
		security:     0x80000000,
	}
}

func (g *Group) scopeMasks() map[string]int {
	return map[string]int{
		global:      0x00000002,
		domainLocal: 0x00000004,
		universal:   0x00000008,
	}
}
