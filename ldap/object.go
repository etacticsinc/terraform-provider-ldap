package ldap

const (
	top = "top"
)

type Object interface {
	GetObjectClass() []string
	GetDN() string
	GetRelativeDN() string
	GetPath() string
	GetAttributes() Attributes
	SetAttributes(attributes Attributes)
	SetDN(dn string)
}
