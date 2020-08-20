package apidoc

import "github.com/getkin/kin-openapi/openapi3"

//Doc 接口文档
type Doc struct {
	openapi3.Swagger
}

var _doc = Doc{
	openapi3.Swagger{
		OpenAPI: "3.0",
		Info: &openapi3.Info{
			Title:          "yumi",
			Description:    "",
			TermsOfService: "",
			Contact: &openapi3.Contact{
				Name:  "liuxin",
				URL:   "",
				Email: "247274526@qq.com",
			},
			License: nil,
			Version: "v1",
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				URL: "http://localhost:8080",
			},
		},
		Paths:        openapi3.Paths{},
		Components:   openapi3.Components{},
		Tags:         openapi3.Tags{},
		Security:     openapi3.SecurityRequirements{},
		ExternalDocs: &openapi3.ExternalDocs{},
	},
}

//Component ...
type Component struct {
	Parameter   Parameter
	RequestBody RequestBody
	Response    Response
}
//Parameter ...
type Parameter interface{}
//RequestBody ...
type RequestBody interface{}
//Response ...
type Response interface{}

//RegisterComponent ...
func RegisterComponent(comp Component) {
	//TODO
}
