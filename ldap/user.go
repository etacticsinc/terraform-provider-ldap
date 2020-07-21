package ldap

import (
	"fmt"
	"github.com/etacticsinc/terraform-provider-ldap/ldap/internal"
	"strconv"
)

const (
	person               = "person"
	organizationalPerson = "organizationalPerson"
	inetOrgPerson        = "inetOrgPerson"
	posixAccount         = "posixAccount"
	shadowAccount        = "shadowAccount"
	user                 = "user"
	normalUserAccount    = "NormalUserAccount"
)

type User struct {
	City              string
	CommonName        string
	Country           string
	Description       string
	EmailAddress      string
	GidNumber         int
	GivenName         string
	HomeDirectory     string
	Name              string
	ObjectClass       []string
	Path              string
	PostalCode        string
	SamAccountName    string
	SamAccountType    string
	State             string
	StreetAddress     string
	Surname           string
	Uid               string
	UidNumber         int
	UserPrincipalName string
}

func (u *User) GetAttributes() Attributes {
	m := map[string][]string{
		"l":                 {u.City},
		"cn":                {u.CommonName},
		"c":                 {u.Country},
		"description":       {u.Description},
		"mail":              {u.EmailAddress},
		"givenName":         {u.GivenName},
		"homeDirectory":     {u.HomeDirectory},
		"name":              {u.Name},
		"objectClass":       u.ObjectClass,
		"postalCode":        {u.PostalCode},
		"sAMAccountName":    {u.SamAccountName},
		"st":                {u.State},
		"streetAddress":     {u.StreetAddress},
		"sn":                {u.Surname},
		"uid":               {u.Uid},
		"userPrincipalName": {u.UserPrincipalName},
	}
	if u.GidNumber != 0 {
		m["gidNumber"] = []string{strconv.Itoa(u.GidNumber)}
	}
	if u.SamAccountType == normalUserAccount {
		m["sAMAccountType"] = []string{fmt.Sprintf("%d", 0x30000000)}
	}
	if u.UidNumber != 0 {
		m["uidNumber"] = []string{strconv.Itoa(u.UidNumber)}
	}
	return Attributes{m}
}

func (u *User) SetAttributes(attributes Attributes) {
	u.City = attributes.GetFirst("l")
	u.CommonName = attributes.GetFirst("cn")
	u.Country = attributes.GetFirst("c")
	u.Description = attributes.GetFirst("description")
	if attributes.HasValue("gidNumber") {
		gidNumber, _ := strconv.Atoi(attributes.GetFirst("gidNumber"))
		u.GidNumber = gidNumber
	}
	u.EmailAddress = attributes.GetFirst("mail")
	u.GivenName = attributes.GetFirst("givenName")
	u.HomeDirectory = attributes.GetFirst("homeDirectory")
	u.Name = attributes.GetFirst("name")
	u.ObjectClass = attributes.Get("objectClass")
	u.PostalCode = attributes.GetFirst("postalCode")
	u.SamAccountName = attributes.GetFirst("sAMAccountName")
	if attributes.HasValue("sAMAccountType") {
		samAccountType, _ := strconv.Atoi(attributes.GetFirst("sAMAccountType"))
		if samAccountType == 0x30000000 {
			u.SamAccountType = normalUserAccount
		}
	}
	u.State = attributes.GetFirst("st")
	u.StreetAddress = attributes.GetFirst("streetAddress")
	u.Surname = attributes.GetFirst("sn")
	u.Uid = attributes.GetFirst("uid")
	if attributes.HasValue("uidNumber") {
		uidNumber, _ := strconv.Atoi(attributes.GetFirst("uidNumber"))
		u.UidNumber = uidNumber
	}
	u.UserPrincipalName = attributes.GetFirst("userPrincipalName")
}

func (u *User) GetObjectClass() []string {
	return u.ObjectClass
}

func (u *User) GetDN() string {
	return fmt.Sprintf("%s,%s", u.GetRelativeDN(), u.Path)
}

func (u *User) GetBaseDN() string {
	return internal.BaseDN(u.Path)
}

func (u *User) GetRelativeDN() string {
	return "cn=" + u.CommonName
}
