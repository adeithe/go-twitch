// Package data contains data structures representing the Twitch API resources and endpoints.
package data

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/adeithe/go-twitch/.codegen/helpers"
)

// TwitchAPIResource represents a Twitch API resource with potential sub-resources.
type TwitchAPIResource struct {
	Name         string
	SubResources []*TwitchAPIResource
}

// TwitchAPIEndpoint represents a Twitch API endpoint with its associated resource, parameters, and response.
type TwitchAPIEndpoint struct {
	Resource     *TwitchAPIResource
	Name         string
	Method, Path string
	DocsURL      string
	Comments     []string
	Params       any
	Response     any
}

// UsesPackage checks if the resource or any of its endpoints use a specific package.
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

// HasBody checks if the resource or any of its endpoints have body parameters.
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

// ResourceName returns the name of the resource with "Resource" suffix.
func (r TwitchAPIResource) ResourceName() string {
	return r.Name + "Resource"
}

// HasBody checks if the endpoint has body parameters.
func (e TwitchAPIEndpoint) HasBody() bool {
	return len(e.GetBodyParams()) > 0
}

// HasResponseField checks if the endpoint's response has a specific field.
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

// InitializationParameters returns a string of required parameters for initializing the endpoint.
func (e TwitchAPIEndpoint) InitializationParameters() string {
	var params []string
	for _, p := range append(e.GetQueryParams(), e.GetBodyParams()...) {
		if p.Required {
			params = append(params, fmt.Sprintf("%s %s", p.ParameterName(), p.Kind()))
		}
	}
	return strings.Join(params, ", ")
}

// GetQueryParams retrieves the query parameters for the endpoint.
func (e TwitchAPIEndpoint) GetQueryParams() (out []helpers.Parameter) {
	if e.Params == nil {
		return out
	}

	t := reflect.TypeOf(e.Params)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		structTag, ok := field.Tag.Lookup("query")
		if !ok {
			continue
		}

		fieldName := helpers.ToSnakeCase(field.Name)
		tags := strings.Split(structTag, ",")
		if len(tags) > 0 && len(tags[0]) > 0 && tags[0] != "-" {
			fieldName = tags[0]
		}

		out = append(out, helpers.Parameter{
			Name:      field.Name,
			FieldName: fieldName,
			Type:      field.Type.String(),
			Required:  slices.Contains(tags, "required"),
			Slice:     field.Type.Kind() == reflect.Slice,
		})
	}
	return out
}

// GetBodyParams retrieves the body parameters for the endpoint.
func (e TwitchAPIEndpoint) GetBodyParams() (out []helpers.Parameter) {
	if e.Params == nil {
		return
	}

	t := reflect.TypeOf(e.Params)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		structTag, ok := field.Tag.Lookup("body")
		if !ok {
			continue
		}

		fieldName := helpers.ToSnakeCase(field.Name)
		tags := strings.Split(structTag, ",")
		if len(tags) > 0 && len(tags[0]) > 0 && tags[0] != "-" {
			fieldName = tags[0]
		}

		out = append(out, helpers.Parameter{
			Name:      field.Name,
			FieldName: fieldName,
			Type:      field.Type.String(),
			Required:  slices.Contains(tags, "required"),
			Slice:     field.Type.Kind() == reflect.Slice,
		})
	}
	return
}

// GetResponseFields retrieves the fields of the endpoint's response.
func (e TwitchAPIEndpoint) GetResponseFields() (fields []helpers.ResponseFields) {
	if e.Response != nil {
		t := reflect.TypeOf(e.Response)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			rawFieldType := strings.TrimPrefix(field.Type.String(), "api.")
			fieldType := rawFieldType
			if field.Type.Kind() == reflect.Struct && field.Name == "Data" {
				fieldType = "[]" + fieldType
			}

			fields = append(fields, helpers.ResponseFields{
				Name:    field.Name,
				Type:    fieldType,
				RawType: rawFieldType,
			})
		}
	}
	return
}

// EndpointName returns the name of the endpoint method.
func (e TwitchAPIEndpoint) EndpointName() string {
	return e.Name + helpers.MethodName(e.Method) + "Call"
}

// ResponseType returns the type of the endpoint's response data.
func (e TwitchAPIEndpoint) ResponseType() string {
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
