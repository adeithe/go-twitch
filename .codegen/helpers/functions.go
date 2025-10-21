// Package helpers provides utility functions for use during code generation.
package helpers

import (
	"net/http"
	"reflect"
	"strings"
	"text/template"
)

// ResponseFields represents a field in an API response.
type ResponseFields struct {
	Name    string
	Type    string
	RawType string
}

// TemplateFuncs is a map of helper functions for use in templates at the global level.
var TemplateFuncs = template.FuncMap{
	"Type":        Type,
	"NoSpaces":    NoSpaces,
	"ToSnakeCase": ToSnakeCase,
	"ToCamelCase": ToCamelCase,
	"IsDefined":   IsDefined,
	"MethodName":  MethodName,
	"ToUpper":     strings.ToUpper,
	"TrimPrefix":  strings.TrimPrefix,
}

// Type returns the string representation of the type of the given value.
func Type(t any) string {
	return reflect.TypeOf(t).String()
}

// NoSpaces removes all spaces from the given string.
func NoSpaces(str string) string {
	return strings.Join(strings.Split(str, " "), "")
}

// ToSnakeCase converts a given string to snake_case.
func ToSnakeCase(str string) string {
	l := len(str)
	diff := 'a' - 'A'
	var b strings.Builder
	for i, v := range str {
		if v >= 'a' {
			b.WriteRune(v)
			continue
		}
		if (i != 0 || i == l-1) && ((i > 0 && rune(str[i-1]) >= 'a') || (i < l-1 && rune(str[i+1]) >= 'a')) {
			b.WriteRune('_')
		}
		b.WriteRune(v + diff)
	}
	return b.String()
}

// ToCamelCase converts a given string to camelCase.
func ToCamelCase(str string) string {
	if len(str) < 2 {
		return strings.ToLower(str)
	}
	return strings.ToLower(string(str[0])) + str[1:]
}

// IsDefined checks if a given value is defined (not nil and not zero value).
func IsDefined(v any) bool {
	if v == nil {
		return false
	}
	return !reflect.ValueOf(v).IsZero()
}

// MethodName returns a standard method name based on the HTTP method.
func MethodName(method string) string {
	switch method {
	case http.MethodGet:
		return "List"
	case http.MethodPost:
		return "Insert"
	case http.MethodPut:
		return "Update"
	case http.MethodPatch:
		return "Modify"
	case http.MethodDelete:
		return "Delete"
	default:
		return ""
	}
}
