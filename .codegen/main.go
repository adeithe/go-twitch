package main

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"text/template"

	"github.com/adeithe/go-twitch/.codegen/data"
	"github.com/adeithe/go-twitch/.codegen/helpers"
)

const (
	TemplateAPIHeader   = "template_twitch_api_header.go.tmpl"
	TemplateAPIResource = "template_twitch_api_resource.go.tmpl"
	TemplateAPICall     = "template_twitch_api_call.go.tmpl"
)

func main() {
	var resourceCount, endpointCount int
	for _, resource := range data.Resources {
		tmplOutput := GetFileName(resource)
		out, err := os.Create(tmplOutput)
		if err != nil {
			slog.Error("failed to create output file",
				slog.String("file", tmplOutput),
				slog.Any("error", err),
			)
			_ = out.Close()
			_ = os.Remove(tmplOutput)
			continue
		}

		endpoints := data.ResourceMapping[resource]
		if err := Generate(out, TemplateAPIHeader, resource); err != nil {
			slog.Error("failed to create output file",
				slog.String("file", tmplOutput),
				slog.Any("error", err),
			)
			_ = out.Close()
			_ = os.Remove(tmplOutput)
			continue
		}

		numResources, numEndpoints, err := GenerateResource(out, resource, endpoints)
		if err != nil {
			slog.Error("failed to generate resource",
				slog.String("output", tmplOutput),
				slog.String("template", TemplateAPIResource),
				slog.Any("error", err),
			)
		}
		resourceCount += numResources
		endpointCount += numEndpoints
		_ = out.Close()
	}
	slog.Info("generated resources", slog.Int("resources", resourceCount), slog.Int("endpoints", endpointCount))
}

func GenerateResource(w io.Writer, resource *data.TwitchAPIResource, endpoints []*data.TwitchAPIEndpoint) (numResources, numEndpoints int, err error) {
	if err = Generate(w, TemplateAPIResource, resource); err != nil {
		return
	}

	for _, subresource := range resource.SubResources {
		subresource.Name = resource.Name + subresource.Name
		nr, ne, err := GenerateResource(w, subresource, data.ResourceMapping[subresource])
		if err != nil {
			return numResources, numEndpoints, err
		}
		numResources += nr
		numEndpoints += ne
	}

	for _, endpoint := range endpoints {
		if strings.HasPrefix(endpoint.DocsURL, "#") {
			endpoint.DocsURL = "https://dev.twitch.tv/docs/api/reference/" + endpoint.DocsURL
		}

		if err = Generate(w, TemplateAPICall, endpoint); err != nil {
			return
		}
		numEndpoints++
	}
	numResources++
	return
}

func Generate(w io.Writer, tmplFile string, data any) error {
	path := strings.Split(tmplFile, "/")
	tmplName := path[len(path)-1]
	tmpl, err := template.New(tmplName).Funcs(helpers.TemplateFuncs).ParseFiles(tmplFile)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(w, tmplName, data)
}

func GetFileName(r *data.TwitchAPIResource) string {
	return "twitch_" + strings.ToLower(helpers.ToSnakeCase(strings.ReplaceAll(r.Name, " ", ""))) + ".go"
}
