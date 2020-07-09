package main

import (
	"errors"
	"fmt"
	"github.com/go-ldap/ldap"
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

func NewOrganizationalUnit(attributes map[string][]string) OrganizationalUnit {
	return OrganizationalUnit{
		City: attributes[city],
	}
}

func (ou *OrganizationalUnit) Attributes() (map[string][]string, error) {
	return map[string][]string{
		"objectClass": {"top", "organizationalUnit"},
		"l":           {ou.City},
		"c":           {ou.Country},
		"description": {ou.Description},
		"displayName": {ou.DisplayName},
		"managedBy":   {ou.ManagedBy},
		"name":        {ou.Name},
		"postalCode":  {ou.PostalCode},
		"st":          {ou.State},
		"street":      {ou.StreetAddress},
	}, nil
}

func (ou *OrganizationalUnit) DistinguishedName() (string, error) {
	if ou.Name == "" {
		return "", errors.New("undefined organizational unit name")
	}

	if ou.Path == "" {
		return "", errors.New("undefined organizational unit path")
	}

	return fmt.Sprintf("CN=%s,%s", ldap.EscapeFilter(ou.Name), ou.Path), nil
}
