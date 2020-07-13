package internal

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"strings"
)

func DN(name string, path string) string {
	return fmt.Sprintf("CN=%s,%s", CN(name), path)
}

func BaseDN(path string) string {
	parts := strings.Split(path, ",")
	baseDN := ""
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if !strings.EqualFold("DC", strings.Split(part, "=")[0]) {
			return baseDN
		}
		baseDN = part + baseDN
	}
	return baseDN
}

func RelativeDN(name string) string {
	return "CN=" + CN(name)
}

func CN(name string) string {
	return ldap.EscapeFilter(name)
}

func Filter(relativeDN string, objectClass []string) string {
	filters := make([]string, len(objectClass)+1)
	filters[0] = relativeDN
	for i, class := range objectClass {
		filters[i+1] = "objectClass=" + ldap.EscapeFilter(class)
	}
	return "(&(" + strings.Join(filters, ")(") + "))"
}
