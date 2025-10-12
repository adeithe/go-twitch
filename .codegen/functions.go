package main

import (
	"net/http"
	"reflect"
	"slices"
	"strings"
	"text/template"
)

type QueryParam struct {
	Name     string
	Type     string
	Required bool
	Slice    bool
}

type BodyParam struct {
	Name     string
	Type     string
	Required bool
	Slice    bool
}

type ResponseFields struct {
	Name    string
	Type    string
	RawType string
}

var templateFuncs = template.FuncMap{
	"Type":              Type,
	"NoSpaces":          NoSpaces,
	"ToSnakeCase":       ToSnakeCase,
	"ToCamelCase":       ToCamelCase,
	"IsDefined":         IsDefined,
	"GetQueryParams":    GetQueryParams,
	"GetBodyParams":     GetBodyParams,
	"GetResponseFields": GetResponseFields,
	"MethodName":        MethodName,
	"ResourceName":      ResourceName,
	"EndpointName":      EndpointName,
	"ResponseType":      ResponseType,
	"ToUpper":           strings.ToUpper,
	"TrimPrefix":        strings.TrimPrefix,
}

func Type(t any) string {
	return reflect.TypeOf(t).String()
}

func NoSpaces(str string) string {
	return strings.Join(strings.Split(str, " "), "")
}

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

func ToCamelCase(str string) string {
	if len(str) < 2 {
		return strings.ToLower(str)
	}
	return strings.ToLower(string(str[0])) + str[1:]
}

func IsDefined(v any) bool {
	if v == nil {
		return false
	}
	return !reflect.ValueOf(v).IsZero()
}

func GetQueryParams(e TwitchAPIEndpoint) []QueryParam {
	var params []QueryParam
	if e.Params == nil {
		return params
	}

	t := reflect.TypeOf(e.Params)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		structTag, ok := field.Tag.Lookup("query")
		if !ok {
			continue
		}

		tags := strings.Split(structTag, ",")
		params = append(params, QueryParam{
			Name:     field.Name,
			Type:     field.Type.String(),
			Required: slices.Contains(tags, "required"),
			Slice:    field.Type.Kind() == reflect.Slice,
		})
	}
	return params
}

func GetBodyParams(e TwitchAPIEndpoint) []BodyParam {
	var params []BodyParam
	if e.Params == nil {
		return params
	}

	t := reflect.TypeOf(e.Params)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		structTag, ok := field.Tag.Lookup("body")
		if !ok {
			continue
		}

		tags := strings.Split(structTag, ",")
		params = append(params, BodyParam{
			Name:     field.Name,
			Type:     field.Type.String(),
			Required: slices.Contains(tags, "required"),
			Slice:    field.Type.Kind() == reflect.Slice,
		})
	}
	return params
}

func GetResponseFields(e TwitchAPIEndpoint) (fields []ResponseFields) {
	if e.Response != nil {
		t := reflect.TypeOf(e.Response)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			rawFieldType := strings.TrimPrefix(field.Type.String(), "api.")
			fieldType := rawFieldType
			if field.Type.Kind() == reflect.Struct && field.Name == "Data" {
				fieldType = "[]" + fieldType
			}

			fields = append(fields, ResponseFields{
				Name:    field.Name,
				Type:    fieldType,
				RawType: rawFieldType,
			})
		}
	}
	return
}

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

func ResourceName(r TwitchAPIResource) string {
	return r.Name + "Resource"
}

func EndpointName(e TwitchAPIEndpoint) string {
	return e.Name + MethodName(e.Method) + "Call"
}

func ResponseType(e TwitchAPIEndpoint) string {
	if e.Response != nil {
		t := reflect.TypeOf(e.Response)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Type.Kind() == reflect.Struct && field.Name == "Data" {
				return strings.TrimPrefix(field.Type.String(), "api.")
			}
		}
	}
	return "any"
}

func (r *TwitchAPIResource) UsesPackage(pkg string) bool {
	for _, s := range r.SubResources {
		if s.UsesPackage(pkg) {
			return true
		}
	}

	for _, e := range Endpoints {
		if e.Resource != r || e.Params == nil {
			continue
		}

		t := reflect.TypeOf(e.Params)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if strings.SplitN(field.Type.PkgPath(), ".", 2)[0] == pkg {
				return true
			}
		}
	}
	return false
}

func (r *TwitchAPIResource) HasBody() bool {
	for _, s := range r.SubResources {
		if s.HasBody() {
			return true
		}
	}

	for _, e := range Endpoints {
		if e.Resource != r {
			continue
		}

		if e.HasBody() {
			return true
		}
	}
	return false
}

func (e TwitchAPIEndpoint) HasBody() bool {
	return len(GetBodyParams(e)) > 0
}

func (e TwitchAPIEndpoint) HasResponseField(fieldName string) bool {
	if e.Response != nil {
		t := reflect.TypeOf(e.Response)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Name == fieldName {
				return true
			}
		}
	}
	return false
}

func (q QueryParam) Kind() string {
	return strings.TrimPrefix(q.Type, "[]")
}

func (q QueryParam) AsVariadicParam() string {
	if q.Slice {
		return "..." + q.Kind()
	}
	return q.Type
}

func (q BodyParam) Kind() string {
	return strings.TrimPrefix(strings.TrimPrefix(q.Type, "[]"), "api.")
}

func (q BodyParam) AsVariadicParam() string {
	return strings.Replace(q.Kind(), "[]", "...", 1)
}
