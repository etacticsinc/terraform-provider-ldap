package main

type LdapObject interface {
	ObjectClass() []string
	DistinguishedName() string
	CommonName() string
	BaseDN() string
	Attributes() Attributes
	SetAttributes(attributes Attributes)
}
