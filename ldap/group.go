package ldap

import (
	"fmt"
	"strconv"
)

const (
	DISTRIBUTION                  = "Distribution"
	SECURITY                      = "Security"
	DOMAIN_LOCAL                  = "DomainLocal"
	GLOBAL                        = "Global"
	UNIVERSAL                     = "Universal"
	POSIX_GROUP                   = "posixGroup"
	GROUP                         = "group"
	SAM_ALIAS_OBJECT              = "AliasObject"
	SAM_GROUP_OBJECT              = "GroupObject"
	SAM_NON_SECURITY_GROUP_OBJECT = "NonSecurityGroupObject"
)

type Group struct {
	CommonName     string
	Description    string
	DN             string
	DisplayName    string
	GidNumber      int
	GroupCategory  string
	GroupScope     string
	HomePage       string
	Members        []string
	MemberUids     []string
	Name           string
	ObjectClass    []string
	Path           string
	SamAccountName string
	SamAccountType string
}

func (g *Group) GetAttributes() Attributes {
	m := map[string][]string{
		"description":    {g.Description},
		"displayName":    {g.DisplayName},
		"gidNumber":      {""},
		"groupType":      {""},
		"member":         g.Members,
		"memberUid":      g.MemberUids,
		"name":           {g.Name},
		"objectClass":    g.ObjectClass,
		"sAMAccountName": {g.SamAccountName},
		"sAMAccountType": {""},
		"wWWHomePage":    {g.HomePage},
	}
	if g.GidNumber != 0 {
		m["gidNumber"] = []string{strconv.Itoa(g.GidNumber)}
	}
	if g.GroupCategory != "" && g.GroupScope != "" {
		// Group Category and Scope are stored as a single bitmask property 'groupType'
		// https://docs.microsoft.com/en-us/windows/win32/adschema/a-grouptype
		groupTypeMask := uint32(0)
		if categoryMask, ok := g.groupCategoryMasks()[g.GroupCategory]; ok {
			groupTypeMask |= categoryMask
		}
		if scopeMask, ok := g.groupScopeMasks()[g.GroupScope]; ok {
			groupTypeMask |= scopeMask
		}
		m["groupType"] = []string{fmt.Sprintf("%d", groupTypeMask)}
	}
	if g.SamAccountType == SAM_GROUP_OBJECT {
		m["sAMAccountType"] = []string{fmt.Sprintf("%d", 0x10000000)}
	} else if g.SamAccountType == SAM_NON_SECURITY_GROUP_OBJECT {
		m["sAMAccountType"] = []string{fmt.Sprintf("%d", 0x10000001)}
	} else if g.SamAccountType == SAM_ALIAS_OBJECT {
		m["sAMAccountType"] = []string{fmt.Sprintf("%d", 0x20000000)}
	}
	return Attributes{m}
}

func (g *Group) SetAttributes(attributes Attributes) {
	g.Description = attributes.GetFirst("description")
	g.DisplayName = attributes.GetFirst("displayName")
	if attributes.HasValue("gidNumber") {
		gidNumber, _ := strconv.Atoi(attributes.GetFirst("gidNumber"))
		g.GidNumber = gidNumber
	}
	if attributes.HasValue("groupType") {
		groupType := attributes.GetFirst("groupType")
		mask, err := strconv.Atoi(groupType)
		if err != nil {
			umask := uint32(mask)
			categoryMasks := g.groupCategoryMasks()
			if umask&categoryMasks[SECURITY] != 0 {
				g.GroupCategory = SECURITY
			} else {
				g.GroupCategory = DISTRIBUTION
			}
			scopeMasks := g.groupScopeMasks()
			if umask&scopeMasks[GLOBAL] != 0 {
				g.GroupScope = GLOBAL
			} else if umask&scopeMasks[DOMAIN_LOCAL] != 0 {
				g.GroupScope = DOMAIN_LOCAL
			} else if umask&scopeMasks[UNIVERSAL] != 0 {
				g.GroupScope = UNIVERSAL
			}
		}
	}
	g.HomePage = attributes.GetFirst("wWWHomePage")
	g.Members = attributes.Get("member")
	g.MemberUids = attributes.Get("memberUid")
	g.Name = attributes.GetFirst("name")
	g.ObjectClass = attributes.Get("objectClass")
	g.SamAccountName = attributes.GetFirst("sAMAccountName")
	if attributes.HasValue("sAMAccountType") {
		samAccountType, _ := strconv.Atoi(attributes.GetFirst("sAMAccountType"))
		if samAccountType == 0x10000000 {
			g.SamAccountType = SAM_GROUP_OBJECT
		} else if samAccountType == 0x10000001 {
			g.SamAccountType = SAM_NON_SECURITY_GROUP_OBJECT
		} else if samAccountType == 0x20000000 {
			g.SamAccountType = SAM_ALIAS_OBJECT
		}
	}
}

func (g *Group) GetObjectClass() []string {
	return g.ObjectClass
}

func (g *Group) GetDN() string {
	return g.DN
}

func (g *Group) GetPath() string {
	return g.Path
}

func (g *Group) GetRelativeDN() string {
	return "cn=" + g.CommonName
}

func (g *Group) SetDN(dn string) {
	g.DN = dn
}

func (g *Group) groupCategoryMasks() map[string]uint32 {
	return map[string]uint32{
		DISTRIBUTION: 0x00000000,
		SECURITY:     0x80000000,
	}
}

func (g *Group) groupScopeMasks() map[string]uint32 {
	return map[string]uint32{
		GLOBAL:       0x00000002,
		DOMAIN_LOCAL: 0x00000004,
		UNIVERSAL:    0x00000008,
	}
}
