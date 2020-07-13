package main

const (
	top = "top"
)

type Object interface {
	GetObjectClass() []string
	GetDN() string
	GetRelativeDN() string
	GetBaseDN() string
	GetAttributes() Attributes
	SetAttributes(attributes Attributes)
}
