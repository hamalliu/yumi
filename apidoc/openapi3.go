package apidoc

import "github.com/getkin/kin-openapi/openapi3"

//Doc 接口文档
type Doc struct {
	openapi3.Swagger
}

var _doc = Doc{
	openapi3.Swagger{
		OpenAPI: "3.0.0",
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
	},
}
