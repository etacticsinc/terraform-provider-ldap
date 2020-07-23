package internal

import (
	"errors"
	"strings"
)

func ParseDN(dn string) (rdn string, path string, err error) {
	rdnIndex := strings.IndexRune(dn, ',')
	if rdnIndex < 1 || rdnIndex == len(dn) - 1 {
		err = errors.New("invalid distinguished name")
		return
	}
	rdn = dn[:rdnIndex]
	path = dn[rdnIndex+1:]
	return
}

func Filter(relativeDN string, objectClass []string) string {
	filtersLen := len(objectClass) + 1
	if filtersLen < 2 {
		filtersLen = 2
	}
	filters := make([]string, filtersLen)
	filters[0] = relativeDN
	if objectClass == nil || len(objectClass) == 0 {
		objectClass = []string{"*"}
	}
	for i, class := range objectClass {
		filters[i + 1] = "objectClass=" + class
	}
	return "(&(" + strings.Join(filters, ")(") + "))"
}
