package main

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"strings"
)

type OrganizationalUnit struct {
	City          string
	Country       string
	Description   string
	DisplayName   string
	ManagedBy     string
	Name          string
	Path          string
	PostalCode    string
	State         string
	StreetAddress string
}



func (ou *OrganizationalUnit) Attributes() Attributes {
	return Attributes{map[string][]string{
		"objectClass": ou.ObjectClass(),
		"l":           {ou.City},
		"c":           {ou.Country},
		"description": {ou.Description},
		"displayName": {ou.DisplayName},
		"managedBy":   {ou.ManagedBy},
		"name":        {ou.Name},
		"instanceType":{fmt.Sprintf("%d", 0x00000004)},
		"postalCode":  {ou.PostalCode},
		"st":          {ou.State},
		"street":      {ou.StreetAddress},
	}}
}

func (ou *OrganizationalUnit) SetAttributes(attributes Attributes) {
	ou.City = attributes.GetFirst("l")
	ou.Country = attributes.GetFirst("c")
	ou.Description = attributes.GetFirst("description")
	ou.DisplayName = attributes.GetFirst("displayName")
	ou.ManagedBy = attributes.GetFirst("managedBy")
	ou.Name = attributes.GetFirst("name")
	ou.PostalCode = attributes.GetFirst("postalCode")
	ou.State = attributes.GetFirst("st")
	ou.StreetAddress = attributes.GetFirst("street")
}

func (ou *OrganizationalUnit) ObjectClass() []string {
	return []string{"top", "organizationalUnit"}
}

func (ou *OrganizationalUnit) DistinguishedName() string {
	return fmt.Sprintf("CN=%s,%s", ou.CommonName(), ou.Path)
}

func (ou * OrganizationalUnit) BaseDN() string {

	parts := strings.Split(ou.Path, ",")

	baseDN := ""

	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if !strings.EqualFold("DC", strings.Split(part, "=")[0]) {
			return baseDN
		}
		baseDN = part + baseDN
	}

	return baseDN
}

func (ou *OrganizationalUnit) CommonName() string {
	return ldap.EscapeFilter(ou.Name)
}
