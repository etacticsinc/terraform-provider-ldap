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
)

type User struct {
	City          string
	Country       string
	Description   string
	GidNumber     int
	HomeDirectory string
	Name          string
	ObjectClass   []string
	Path          string
	PostalCode    string
	State         string
	StreetAddress string
	Surname       string
	Uid           string
	UidNumber     int
}

func (u *User) GetAttributes() Attributes {
	m := map[string][]string{
		"l":             {u.City},
		"c":             {u.Country},
		"description":   {u.Description},
		"homeDirectory": {u.HomeDirectory},
		"cn":            {u.Name},
		"objectClass":   u.ObjectClass,
		"postalCode":    {u.PostalCode},
		"st":            {u.State},
		"streetAddress": {u.StreetAddress},
		"sn":            {u.Surname},
		"uid":           {u.Uid},
	}
	if u.GidNumber != 0 {
		m["gidNumber"] = []string{strconv.Itoa(u.GidNumber)}
	}
	if u.UidNumber != 0 {
		m["uidNumber"] = []string{strconv.Itoa(u.UidNumber)}
	}
	return Attributes{m}
}

func (u *User) SetAttributes(attributes Attributes) {
	if attributes.HasValue("gidNumber") {
		gidNumber, _ := strconv.Atoi(attributes.GetFirst("gidNumber"))
		u.GidNumber = gidNumber
	}
	if attributes.HasValue("uidNumber") {
		uidNumber, _ := strconv.Atoi(attributes.GetFirst("uidNumber"))
		u.UidNumber = uidNumber
	}
	u.City = attributes.GetFirst("l")
	u.Country = attributes.GetFirst("c")
	u.Description = attributes.GetFirst("description")
	u.HomeDirectory = attributes.GetFirst("homeDirectory")
	u.Name = attributes.GetFirst("cn")
	u.ObjectClass = attributes.Get("objectClass")
	u.PostalCode = attributes.GetFirst("postalCode")
	u.State = attributes.GetFirst("st")
	u.StreetAddress = attributes.GetFirst("streetAddress")
	u.Surname = attributes.GetFirst("sn")
	u.Uid = attributes.GetFirst("uid")
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
	return "cn=" + u.Name
}
