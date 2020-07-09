package main

import (
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

type User struct {
	AccountExpirationDate string
	City                  string
	Company               string
	Country               string
	Department            string
	Description           string
	DisplayName           string
	Division              string
	EmailAddress          string
	EmployeeID            string
	EmployeeNumber        string
	Fax                   string
	GivenName             string
	HomeDirectory         string
	HomeDrive             string
	HomePage              string
	HomePhone             string
	Initials              string
	LogonWorkstations     string
	Manager               string
	MobilePhone           string
	Name                  string
	Office                string
	OfficePhone           string
	Organization          string
	OtherName             string
	POBox                 string
	Path                  string
	PostalCode            string
	ProfilePath           string
	SamAccountName        string
	ScriptPath            string
	ServicePrincipalNames string
	State                 string
	StreetAddress         string
	Surname               string
	Title                 string
	UserPrincipalName     string
}

func (u *User) DistinguishedName() (string, error) {
	if u.Name == "" {
		return "", errors.New("undefined user name")
	}

	if u.Path == "" {
		return "", errors.New("undefined user path")
	}

	return fmt.Sprintf("CN=%s,%s", ldap.EscapeFilter(u.Name), u.Path), nil
}

func (u *User) Attributes() (map[string][]string, error) {
	return map[string][]string{
		"objectClass":  {"top", "organizationalPerson", "user", "person"},
		"instanceType": {fmt.Sprintf("%d", 0x00000004)},
	}, nil
}
