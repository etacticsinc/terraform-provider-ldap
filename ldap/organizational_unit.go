package ldap

import (
	"fmt"
	"github.com/etacticsinc/terraform-provider-ldap/ldap/internal"
)

const (
	organizationalUnit = "organizationalUnit"
)

type OrganizationalUnit struct {
	City               string
	Country            string
	Description        string
	Name               string
	ObjectClass        []string
	OrganizationalUnit string
	Path               string
	PostalCode         string
	State              string
	StreetAddress      string
}

func (ou *OrganizationalUnit) GetAttributes() Attributes {
	return Attributes{map[string][]string{
		"l":             {ou.City},
		"c":             {ou.Country},
		"description":   {ou.Description},
		"name":          {ou.Name},
		"objectClass":   ou.ObjectClass,
		"ou":            {ou.OrganizationalUnit},
		"postalCode":    {ou.PostalCode},
		"st":            {ou.State},
		"streetAddress": {ou.StreetAddress},
	}}
}

func (ou *OrganizationalUnit) SetAttributes(attributes Attributes) {
	ou.ObjectClass = attributes.Get("objectClass")
	ou.City = attributes.GetFirst("l")
	ou.Country = attributes.GetFirst("c")
	ou.Description = attributes.GetFirst("description")
	ou.Name = attributes.GetFirst("name")
	ou.OrganizationalUnit = attributes.GetFirst("ou")
	ou.PostalCode = attributes.GetFirst("postalCode")
	ou.State = attributes.GetFirst("st")
	ou.StreetAddress = attributes.GetFirst("streetAddress")
}

func (ou *OrganizationalUnit) GetObjectClass() []string {
	return ou.ObjectClass
}

func (ou *OrganizationalUnit) GetDN() string {
	return fmt.Sprintf("%s,%s", ou.GetRelativeDN(), ou.Path)
}

func (ou *OrganizationalUnit) GetBaseDN() string {
	return internal.BaseDN(ou.Path)
}

func (ou *OrganizationalUnit) GetRelativeDN() string {
	return "ou=" + ou.OrganizationalUnit
}
