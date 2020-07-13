package main

import (
	"fmt"
	"ldap/internal"
)

type OrganizationalUnit struct {
	City          string
	Country       string
	Description   string
	Name          string
	Path          string
	PostalCode    string
	State         string
	StreetAddress string
}

func (ou *OrganizationalUnit) Attributes() Attributes {
	return Attributes{map[string][]string{
		"objectClass":   ou.Class(),
		"l":             {ou.City},
		"c":             {ou.Country},
		"description":   {ou.Description},
		"ou":            {ou.Name},
		"postalCode":    {ou.PostalCode},
		"st":            {ou.State},
		"streetAddress": {ou.StreetAddress},
	}}
}

func (ou *OrganizationalUnit) SetAttributes(attributes Attributes) {
	ou.City = attributes.GetFirst("l")
	ou.Country = attributes.GetFirst("c")
	ou.Description = attributes.GetFirst("description")
	ou.Name = attributes.GetFirst("ou")
	ou.PostalCode = attributes.GetFirst("postalCode")
	ou.State = attributes.GetFirst("st")
	ou.StreetAddress = attributes.GetFirst("streetAddress")
}

func (ou *OrganizationalUnit) Class() []string {
	return []string{"top", "organizationalUnit"}
}

func (ou *OrganizationalUnit) DN() string {
	return fmt.Sprintf("%s,%s", ou.RelativeDN(), ou.Path)
}

func (ou *OrganizationalUnit) BaseDN() string {
	return internal.BaseDN(ou.Path)
}

func (ou *OrganizationalUnit) RelativeDN() string {
	return "ou=" + ou.Name
}
