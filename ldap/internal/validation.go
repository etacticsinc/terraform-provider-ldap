package internal

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func SetIntersection(values []interface{}, minLen int) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		s, ok := i.(schema.Set)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be schema.Set", k))
			return warnings, errors
		}
		size := 0
		for _, elem := range s.List() {
			for _, value := range values {
				if elem == value {
					size++
				}
			}
		}
		if size < minLen {
			if len(values) == minLen {
				errors = append(errors, fmt.Errorf("must contain %v", values))
			} else {
				errors = append(errors, fmt.Errorf("must contain at least %v of %v", minLen, values))
			}
		}
		return warnings, errors
	}
}