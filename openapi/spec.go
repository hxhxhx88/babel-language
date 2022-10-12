package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
)

func New() *openapi3.T {
	return &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "Babel API",
			Description: "REST APIs used for interacting with the Babel Service",
			Version:     "0.0.1",
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				Description: "Local development",
				URL:         "http://0.0.0.0:12346",
			},
		},
		Paths: openapi3.Paths{
			"/ping": &openapi3.PathItem{
				Get: &openapi3.Operation{
					Summary:     "Ping the server",
					OperationID: "ping",
					Responses: openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Value: openapi3.NewResponse().WithDescription("Success"),
						},
					},
				},
			},
		},
	}
}
