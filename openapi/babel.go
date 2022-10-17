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
				Get: &openapi3.Operation{
					Summary:     "List corpuses",
					OperationID: "ListCorpuses",
					Responses: openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Value: openapi3.NewResponse().WithDescription("Success").WithJSONSchema(&openapi3.Schema{
								Type:     openapi3.TypeObject,
								Required: []string{"corpuses"},
								Properties: openapi3.Schemas{
									"corpuses": &openapi3.SchemaRef{
										Value: &openapi3.Schema{
											Type: openapi3.TypeArray,
											Items: &openapi3.SchemaRef{
												Ref: "#/components/schemas/Corpus",
											},
										},
									},
								},
							}),
						},
					},
				},
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
			"/corpus/{corpusId}": &openapi3.PathItem{
				Get: &openapi3.Operation{
					Summary:     "Get for a corpus",
					OperationID: "GetCorpus",
					Parameters: openapi3.Parameters{
						&openapi3.ParameterRef{
							Value: mCorpusIdParameter,
						},
					},
					Responses: openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Value: openapi3.NewResponse().WithDescription("Success").WithJSONSchema(&openapi3.Schema{
								Type:     openapi3.TypeObject,
								Required: []string{"corpus"},
								Properties: openapi3.Schemas{
									"corpus": &openapi3.SchemaRef{
										Ref: "#/components/schemas/Corpus",
									},
								},
							}),
						},
					},
				},
			},
			"/corpus/{corpusId}/translations": &openapi3.PathItem{
				Get: &openapi3.Operation{
					Summary:     "List corpus translations",
					OperationID: "ListCorpusTranslations",
					Parameters: openapi3.Parameters{
						&openapi3.ParameterRef{
							Value: mCorpusIdParameter,
						},
					},
					Responses: openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Value: openapi3.NewResponse().WithDescription("Success").WithJSONSchema(&openapi3.Schema{
								Type:     openapi3.TypeObject,
								Required: []string{"translations"},
								Properties: openapi3.Schemas{
									"translations": &openapi3.SchemaRef{
										Value: &openapi3.Schema{
											Type: openapi3.TypeArray,
											Items: &openapi3.SchemaRef{
												Ref: "#/components/schemas/Translation",
											},
										},
									},
								},
							}),
						},
					},
				},
				Post: &openapi3.Operation{
					Summary:     "Create a translation for a corpus",
					OperationID: "CreateCorpusTranslation",
					Parameters: openapi3.Parameters{
						&openapi3.ParameterRef{
							Value: mCorpusIdParameter,
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
			"/translation/{translationId}/blocks/_search": &openapi3.PathItem{
				Post: &openapi3.Operation{
					Summary:     "Search translation blocks",
					OperationID: "SearchTranslationBlocks",
					Parameters: openapi3.Parameters{
						&openapi3.ParameterRef{
							Value: mTranslationIdParameter,
						},
					},
					RequestBody: &openapi3.RequestBodyRef{
						Value: openapi3.NewRequestBody().WithJSONSchema(&openapi3.Schema{
							Type:     openapi3.TypeObject,
							Required: []string{"pagination", "filter"},
							Properties: openapi3.Schemas{
								"pagination": &openapi3.SchemaRef{
									Ref: "#/components/schemas/Pagination",
								},
								"filter": &openapi3.SchemaRef{
									Ref: "#/components/schemas/BlockFilter",
								},
							},
						}),
					},
					Responses: openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Value: openapi3.NewResponse().WithDescription("Success").WithJSONSchema(&openapi3.Schema{
								Type:     openapi3.TypeObject,
								Required: []string{"blocks"},
								Properties: openapi3.Schemas{
									"blocks": &openapi3.SchemaRef{
										Value: &openapi3.Schema{
											Type: openapi3.TypeArray,
											Items: &openapi3.SchemaRef{
												Ref: "#/components/schemas/Block",
											},
										},
									},
								},
							}),
						},
					},
				},
			},
			"/translation/{translationId}/blocks/_count": &openapi3.PathItem{
				Post: &openapi3.Operation{
					Summary:     "Count translation blocks",
					OperationID: "CountTranslationBlocks",
					Parameters: openapi3.Parameters{
						&openapi3.ParameterRef{
							Value: mTranslationIdParameter,
						},
					},
					RequestBody: &openapi3.RequestBodyRef{
						Value: openapi3.NewRequestBody().WithJSONSchema(&openapi3.Schema{
							Type:     openapi3.TypeObject,
							Required: []string{"filter"},
							Properties: openapi3.Schemas{
								"filter": &openapi3.SchemaRef{
									Ref: "#/components/schemas/BlockFilter",
								},
							},
						}),
					},
					Responses: openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Value: openapi3.NewResponse().WithDescription("Success").WithJSONSchema(&openapi3.Schema{
								Type:     openapi3.TypeObject,
								Required: []string{"total_count"},
								Properties: openapi3.Schemas{
									"total_count": mIntField,
								},
							}),
						},
					},
				},
			},
			"/block/{blockId}/_translate": &openapi3.PathItem{
				Post: &openapi3.Operation{
					Summary:     "Translate a block",
					OperationID: "TranslateBlock",
					Parameters: openapi3.Parameters{
						&openapi3.ParameterRef{
							Value: mBlockIdParameter,
						},
					},
					RequestBody: &openapi3.RequestBodyRef{
						Value: openapi3.NewRequestBody().WithJSONSchema(&openapi3.Schema{
							Type:     openapi3.TypeObject,
							Required: []string{"translation_ids"},
							Properties: openapi3.Schemas{
								"translation_ids": &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type:  openapi3.TypeArray,
										Items: mIdField,
									},
								},
							},
						}),
					},
					Responses: openapi3.Responses{
						"200": &openapi3.ResponseRef{
							Value: openapi3.NewResponse().WithDescription("Success").WithJSONSchema(&openapi3.Schema{
								Type:     openapi3.TypeObject,
								Required: []string{"blocks"},
								Properties: openapi3.Schemas{
									"blocks": &openapi3.SchemaRef{
										Value: &openapi3.Schema{
											Type: openapi3.TypeArray,
											Items: &openapi3.SchemaRef{
												Ref: "#/components/schemas/Block",
											},
										},
									},
								},
							}),
						},
					},
				},
			},
		},
		Components: openapi3.Components{
			Schemas: openapi3.Schemas{
				"Pagination": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:     openapi3.TypeObject,
						Required: []string{"page", "page_size"},
						Properties: openapi3.Schemas{
							"page":      mIntField,
							"page_size": mIntField,
						},
					},
				},
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
				"Corpus": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:     openapi3.TypeObject,
						Required: []string{"id", "title", "original_language_iso_639_3"},
						Properties: openapi3.Schemas{
							"id":                          mIdField,
							"title":                       mStringField,
							"original_language_iso_639_3": mStringField,
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
				"Translation": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:     openapi3.TypeObject,
						Required: []string{"id", "corpus_id", "language_iso_639_3", "title"},
						Properties: openapi3.Schemas{
							"id":                 mIdField,
							"corpus_id":          mIdField,
							"title":              mStringField,
							"language_iso_639_3": mStringField,
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
				"Block": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:     openapi3.TypeObject,
						Required: []string{"id", "translation_id", "content", "rank", "uuid"},
						Properties: openapi3.Schemas{
							"id":             mIdField,
							"translation_id": mIdField,
							"content":        mStringField,
							"rank":           mIntField,
							"uuid":           mStringField,
						},
					},
				},
				"BlockFilter": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: openapi3.TypeObject,
						Properties: openapi3.Schemas{
							"parent_block_id": &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type:        mIdType,
									Description: "If provided, all sub-blocks, i.e. thoese blocks whose uuid has this block's uuid as prefix, at next rank will be returned, otherwise returns top-level blocks.",
								},
							},
							"content": mStringField,
						},
					},
				},
			},
		},
	}
}

const mIdType = openapi3.TypeString

var mIdField = &openapi3.SchemaRef{
	Value: &openapi3.Schema{
		Type: mIdType,
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

var mCorpusIdParameter = &openapi3.Parameter{
	Name:     "corpusId",
	In:       openapi3.ParameterInPath,
	Required: true,
	Schema:   mIdField,
}

var mTranslationIdParameter = &openapi3.Parameter{
	Name:     "translationId",
	In:       openapi3.ParameterInPath,
	Required: true,
	Schema:   mIdField,
}

var mBlockIdParameter = &openapi3.Parameter{
	Name:     "blockId",
	In:       openapi3.ParameterInPath,
	Required: true,
	Schema:   mIdField,
}
