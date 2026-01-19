package sys

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// FindPrefixedEnvVars finds prefixed environment variables.
func FindPrefixedEnvVars(environ []string, prefix string, element any) []string {
	prefixes := getRootPrefixes(element, prefix)

	var values []string
	for _, px := range prefixes {
		for _, value := range environ {
			if strings.HasPrefix(value, px) {
				values = append(values, value)
			}
		}
	}

	return values
}

func getRootPrefixes(element any, prefix string) []string {
	if element == nil {
		return nil
	}

	rootType := reflect.TypeOf(element)

	return getPrefixes(prefix, rootType)
}

func getPrefixes(prefix string, rootType reflect.Type) []string {
	var names []string

	if rootType.Kind() == reflect.Pointer {
		rootType = rootType.Elem()
	}

	if rootType.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < rootType.NumField(); i++ {
		field := rootType.Field(i)

		if !IsExported(field) {
			continue
		}

		if field.Anonymous &&
			(field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.Struct || field.Type.Kind() == reflect.Struct) {
			names = append(names, getPrefixes(prefix, field.Type)...)
			continue
		}

		names = append(names, prefix+strings.ToUpper(field.Name))
	}

	return names
}

// IsExported reports whether f is exported.
// https://golang.org/pkg/reflect/#StructField
func IsExported(f reflect.StructField) bool {
	return f.PkgPath == ""
}

func checkPrefix(prefix string) error {
	prefixPattern := `^[a-zA-Z0-9]+_$`
	matched, err := regexp.MatchString(prefixPattern, prefix)
	if err != nil {
		return err
	}

	if !matched {
		return fmt.Errorf("invalid prefix %q, the prefix pattern must match the following pattern: %s", prefix, prefixPattern)
	}

	return nil
}
