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
									"corpus": &openapi3.SchemaRef{
										Ref: "#/components/schemas/CorpusDraft",
									},
								},
							}),
						},
					},
				},
			},
			"/corpuses": &openapi3.PathItem{
				Post: &openapi3.Operation{
					Summary:     "Create a corpus",
					OperationID: "CreateCorpus",
					RequestBody: &openapi3.RequestBodyRef{
						Value: openapi3.NewRequestBody().WithJSONSchema(&openapi3.Schema{
							Type:     openapi3.TypeObject,
							Required: []string{"corpus"},
							Properties: openapi3.Schemas{
								"corpus": &openapi3.SchemaRef{
									Ref: "#/components/schemas/CorpusDraft",
								},
							},
						}),
					},
					Responses: openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Value: openapi3.NewResponse().WithDescription("Success").WithJSONSchema(&openapi3.Schema{
								Type:     openapi3.TypeObject,
								Required: []string{"id"},
								Properties: openapi3.Schemas{
									"id": mIdField,
								},
							}),
						},
					},
				},
			},
			"/corpus/{corpusId}/translations": &openapi3.PathItem{
				Post: &openapi3.Operation{
					Summary:     "Create a translation for a corpus",
					OperationID: "CreateCorpusTranslation",
					Parameters: openapi3.Parameters{
						&openapi3.ParameterRef{
							Value: &openapi3.Parameter{
								Name:     "corpusId",
								In:       openapi3.ParameterInPath,
								Required: true,
								Schema:   mIdField,
							},
						},
					},
					RequestBody: &openapi3.RequestBodyRef{
						Value: openapi3.NewRequestBody().WithJSONSchema(&openapi3.Schema{
							Type:     openapi3.TypeObject,
							Required: []string{"translation"},
							Properties: openapi3.Schemas{
								"translation": &openapi3.SchemaRef{
									Ref: "#/components/schemas/TranslationDraft",
								},
							},
						}),
					},
					Responses: openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Value: openapi3.NewResponse().WithDescription("Success").WithJSONSchema(&openapi3.Schema{
								Type:     openapi3.TypeObject,
								Required: []string{"id"},
								Properties: openapi3.Schemas{
									"id": mIdField,
								},
							}),
						},
					},
				},
			},
		},
		Components: openapi3.Components{
			Schemas: openapi3.Schemas{
				"CorpusDraft": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:     openapi3.TypeObject,
						Required: []string{"title", "original_language_iso_639_3"},
						Properties: openapi3.Schemas{
							"title":                       mStringField,
							"original_language_iso_639_3": mStringField,
							"translations": &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: openapi3.TypeArray,
									Items: &openapi3.SchemaRef{
										Ref: "#/components/schemas/TranslationDraft",
									},
								},
							},
						},
					},
				},
				"TranslationDraft": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:     openapi3.TypeObject,
						Required: []string{"language_iso_639_3", "title"},
						Properties: openapi3.Schemas{
							"title":              mStringField,
							"language_iso_639_3": mStringField,
							"blocks": &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: openapi3.TypeArray,
									Items: &openapi3.SchemaRef{
										Ref: "#/components/schemas/BlockDraft",
									},
								},
							},
						},
					},
				},
				"BlockDraft": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:     openapi3.TypeObject,
						Required: []string{"content", "rank", "uuid"},
						Properties: openapi3.Schemas{
							"content": mStringField,
							"rank":    mIntField,
							"uuid":    mStringField,
						},
					},
				},
			},
		},
	}
}

var mIdField = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.TypeString,
	},
}

var mStringField = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.TypeString,
	},
}

var mIntField = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: openapi3.TypeInteger,
	},
}
