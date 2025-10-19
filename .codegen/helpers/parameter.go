package helpers

import "strings"

// Parameter represents a parameter in a resource method.
type Parameter struct {
	Name      string
	FieldName string
	Type      string
	Required  bool
	Slice     bool
}

// Kind returns the underlying type of the parameter.
func (p Parameter) Kind() string {
	return strings.TrimPrefix(strings.TrimPrefix(p.Type, "[]"), "api.")
}

// AsVariadicParam returns the parameter type formatted as a variadic parameter.
func (p Parameter) AsVariadicParam() string {
	if !p.Slice {
		return p.Kind()
	}
	return "..." + p.Kind()
}

// ParameterName retrieves the parameters of the endpoint.
func (p Parameter) ParameterName() string {
	str := ToCamelCase(p.Name)
	if strings.HasSuffix(strings.ToLower(str), "id") {
		str = str[:len(str)-2] + "ID"
	}

	if strings.ToLower(str) == "id" || strings.ToLower(str) == "ids" {
		return strings.ToLower(str)
	}
	return str
}
