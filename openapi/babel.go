package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
)

func NewBabel() *openapi3.T {
	return &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "Babel API",
			Description: "REST APIs used for interacting with the Babel Service",
			Version:     "0.0.1",
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				URL: "/api",
			},
		},
		Paths: openapi3.Paths{
			"/metadata": &openapi3.PathItem{
				Get: &openapi3.Operation{
					Summary:     "Get server metadata",
					OperationID: "GetMetadata",
					Responses: openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Value: openapi3.NewResponse().WithDescription("Success").WithJSONSchema(&openapi3.Schema{
								Type:     openapi3.TypeObject,
								Required: []string{"commit_identifier", "commit_time", "version"},
								Properties: openapi3.Schemas{
									"commit_identifier": mStringField,
									"commit_time":       mStringField,
									"version":           mStringField,
								},
							}),
						},
					},
				},
			},
		},
	}
}

var mStringField = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.TypeString,
	},
}
