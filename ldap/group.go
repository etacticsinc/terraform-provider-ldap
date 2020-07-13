package main

import (
	"fmt"
	"ldap/internal"
	"strconv"
)

const (
	distribution = "Distribution"
	security     = "Security"
	domainLocal  = "DomainLocal"
	global       = "Global"
	universal    = "Universal"
	posixGroup   = "posixGroup"
	group        = "group"
)

type Group struct {
	Description    string
	GidNumber      string
	GroupCategory  string
	GroupScope     string
	HomePage       string
	Name           string
	ObjectClass    []string
	Path           string
	SamAccountName string
}

func (g *Group) GetAttributes() Attributes {
	m := map[string][]string{
		"objectClass":    g.ObjectClass,
		"description":    {g.Description},
		"gidNumber":      {g.GidNumber},
		"wWWHomePage":    {g.HomePage},
		"sAMAccountName": {g.SamAccountName},
	}
	if g.GroupCategory != "" && g.GroupScope != "" {
		// Group Category and Scope are stored as a single bitmask property 'groupType'
		// https://docs.microsoft.com/en-us/windows/win32/adschema/a-grouptype
		groupTypeMask := 0
		if categoryMask, ok := g.categoryMasks()[g.GroupCategory]; ok {
			groupTypeMask |= categoryMask
		}
		if scopeMask, ok := g.scopeMasks()[g.GroupScope]; ok {
			groupTypeMask |= scopeMask
		}
		m["groupType"] = []string{fmt.Sprintf("%d", groupTypeMask)}
	}
	return Attributes{m}
}

func (g *Group) SetAttributes(attributes Attributes) {
	groupType := attributes.GetFirst("groupType")
	if groupType != "" {
		mask, err := strconv.Atoi(groupType)
		if err != nil {
			categoryMasks := g.categoryMasks()
			securityMask := categoryMasks[security]
			if mask&securityMask == securityMask {
				g.GroupCategory = security
			} else {
				g.GroupCategory = distribution
			}
			scopeMasks := g.scopeMasks()
			globalMask := scopeMasks[global]
			domainLocalMask := scopeMasks[domainLocal]
			universalMask := scopeMasks[universal]
			if mask&globalMask == globalMask {
				g.GroupScope = global
			} else if mask&domainLocalMask == domainLocalMask {
				g.GroupScope = domainLocal
			} else if mask&universalMask == universalMask {
				g.GroupScope = universal
			}
		}
	}
	g.ObjectClass = attributes.Get("objectClass")
	g.Description = attributes.GetFirst("description")
	g.GidNumber = attributes.GetFirst("gidNumber")
	g.HomePage = attributes.GetFirst("wWWHomePage")
	g.SamAccountName = attributes.GetFirst("sAMAccountName")
}

func (g *Group) GetObjectClass() []string {
	return g.ObjectClass
}

func (g *Group) GetDN() string {
	return fmt.Sprintf("%s,%s", g.GetRelativeDN(), g.Path)
}

func (g *Group) GetBaseDN() string {
	return internal.BaseDN(g.Path)
}

func (g *Group) GetRelativeDN() string {
	return "cn=" + g.Name
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
