package ldap

import (
	"./internal"
	"fmt"
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
	ObjectGUID    string
	State         string
	StreetAddress string
}

func (ou *OrganizationalUnit) Attributes() Attributes {
	return Attributes{map[string][]string{
		"objectClass":  ou.Class(),
		"l":            {ou.City},
		"c":            {ou.Country},
		"description":  {ou.Description},
		"displayName":  {ou.DisplayName},
		"managedBy":    {ou.ManagedBy},
		"name":         {ou.Name},
		"instanceType": {fmt.Sprintf("%d", 0x00000004)},
		"objectGUID":   {ou.ObjectGUID},
		"postalCode":   {ou.PostalCode},
		"st":           {ou.State},
		"street":       {ou.StreetAddress},
	}}
}

func (ou *OrganizationalUnit) SetAttributes(attributes Attributes) {
	ou.City = attributes.GetFirst("l")
	ou.Country = attributes.GetFirst("c")
	ou.Description = attributes.GetFirst("description")
	ou.DisplayName = attributes.GetFirst("displayName")
	ou.ManagedBy = attributes.GetFirst("managedBy")
	ou.Name = attributes.GetFirst("name")
	ou.ObjectGUID = attributes.GetFirst("objectGUID")
	ou.PostalCode = attributes.GetFirst("postalCode")
	ou.State = attributes.GetFirst("st")
	ou.StreetAddress = attributes.GetFirst("street")
}

func (ou *OrganizationalUnit) Class() []string {
	return []string{"top", "organizationalUnit"}
}

func (ou *OrganizationalUnit) DN() string {
	return internal.DN(ou.Name, ou.Path)
}

func (ou *OrganizationalUnit) BaseDN() string {
	return internal.BaseDN(ou.Path)
}

func (ou *OrganizationalUnit) RelativeDN() string {
	return internal.RelativeDN(ou.Name)
}

func (ou *OrganizationalUnit) CN() string {
	return internal.CN(ou.Name)
}
