package main

type Object interface {
	Class() []string
	DN() string
	RelativeDN() string
	BaseDN() string
	Attributes() Attributes
	SetAttributes(attributes Attributes)
}
