package internal

import (
	"github.com/go-ldap/ldap/v3"
	"strings"
)

func BaseDN(path string) string {
	parts := strings.Split(path, ",")
	baseDN := ""
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if !strings.EqualFold("dc", strings.Split(part, "=")[0]) {
			return baseDN
		}
		if baseDN != "" {
			part += ","
		}
		baseDN = part + baseDN
	}
	return baseDN
}

func Filter(relativeDN string, objectClass []string) string {
	filters := make([]string, len(objectClass)+1)
	filters[0] = relativeDN
	for i, class := range objectClass {
		filters[i+1] = "objectClass=" + ldap.EscapeFilter(class)
	}
	return "(&(" + strings.Join(filters, ")(") + "))"
}
