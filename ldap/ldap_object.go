package ldap

type Object interface {
	Class() []string
	DN() string
	RelativeDN() string
	CN() string
	BaseDN() string
	Attributes() Attributes
	SetAttributes(attributes Attributes)
}
