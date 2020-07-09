package main

type LdapObject interface {
	DistinguishedName() (string, error)
	Attributes() (map[string][]string, error)
}
