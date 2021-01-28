package apidoc

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Parse ...
func Parse(filePath, prefix string) (err error) {
	// TODO:
	return
}

func parseRouter(filePath, prefix string) (path string, pathItem *openapi3.PathItem) {
	// TODO:
	return 
}

// ParseDto ...
func ParseDto(filePath, funcName string) (
	prms *openapi3.Parameters, body *openapi3.RequestBody, resp *openapi3.Schema, err error) {
	fset := token.NewFileSet()

	filePath = strings.ReplaceAll(filePath, "\"", "")
	f, err := parser.ParseFile(fset, filePath+"/dto.go", nil, parser.Trace)
	if err != nil {
		return
	}

	ps := parseParameters(f, funcName+"Request")
	if len(ps) > 0 {
		prms = &ps
	}

	reqSche := parseSchema(f, funcName+"Request")
	if len(reqSche.Properties) > 0 {
		body.Content["application/json"] = &openapi3.MediaType{Schema: &openapi3.SchemaRef{Value: &reqSche}}
	}

	respSche := parseSchema(f, funcName+"Response")
	if len(respSche.Properties) > 0 {
		resp = &respSche
	}

	return
}

func parseParameters(f *ast.File, objName string) (ps openapi3.Parameters) {
	for _, obj := range f.Scope.Objects {
		if obj.Name == objName {
			fields := obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List
			for _, field := range fields {
				tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
				if tag.Get("query") != "" {
					p := openapi3.Parameter{
						Name:        tag.Get("query"),
						In:          "query",
						Description: tag.Get("desc"),
						Required:    tag.Get("binding") == "required",
						Schema:      &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
					}
					ps = append(ps, &openapi3.ParameterRef{Value: &p})
				}
			}
		}
	}
	return
}

func parseSchema(f *ast.File, objName string) (s openapi3.Schema) {
	s.Type = "object"
	s.Properties = make(map[string]*openapi3.SchemaRef)

	for _, obj := range f.Scope.Objects {
		if obj.Name == objName {
			fields := obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List

			for _, field := range fields {
				tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
				ss := openapi3.Schema{}
				ss.Description = tag.Get("desc")
				if tag.Get("binding") == "required" {
					ss.Required = append(ss.Required, tag.Get("json"))
				}

				if ident, ok := field.Type.(*ast.Ident); ok {
					ss.Type = ident.Name
				}
				if star, ok := field.Type.(*ast.StarExpr); ok {
					if ident, ok := star.X.(*ast.Ident); ok {
						ss.Type = ident.Name
					}
				}
				if array, ok := field.Type.(*ast.ArrayType); ok {
					if ident, ok := array.Elt.(*ast.Ident); ok {
						ss.Type = "array " + ident.Name
					}

					if star, ok := array.Elt.(*ast.StarExpr); ok {
						if ident, ok := star.X.(*ast.Ident); ok {
							ss.Type = "array " + ident.Name
						}
					}
				}

				switch ss.Type {
				case "string":
					ss.MinLength = 0
					ss.MaxLength = nil
					ss.Pattern = ""
					ss.Type = "string"
					ss.Format = ""
				case "uint8", "uint16", "uint32", "uint64",
					"int8", "int16", "int32", "int64", "uint", "int",
					"float32", "float64":
					ss.Min = nil
					ss.Max = nil
					ss.MultipleOf = nil
					ss.Type = "number"
				case "bool":
					ss.Type = "boolean"
				case "array string":
					item := openapi3.Schema{}
					item.MinLength = 0
					item.MaxLength = nil
					item.Pattern = ""
					item.Type = "string"

					ss.Type = "array"
					ss.Items = openapi3.NewSchemaRef("", &item)
				case "array uint8", "array uint16", "array uint32",
					"array uint64", "array int8", "array int16",
					"array int32", "array int64", "array uint",
					"array int", "array float32", "array float64":
					item := openapi3.Schema{}
					item.Min = nil
					item.Max = nil
					item.MultipleOf = nil
					item.Type = "number"

					ss.Type = "array"
					ss.Items = openapi3.NewSchemaRef("", &item)
				case "array bool":
					item := openapi3.Schema{}
					item.Type = "boolean"

					ss.Type = "array"
					ss.Items = openapi3.NewSchemaRef("", &item)
				default:
					if strings.HasPrefix(ss.Type, "array ") {
						ss.Type = strings.TrimPrefix(ss.Type, "array ")
						item := openapi3.Schema{}
						item = parseSchema(f, ss.Type)

						ss.Type = "array"
						ss.Items = openapi3.NewSchemaRef("", &item)
					} else {
						item := openapi3.Schema{}
						item = parseSchema(f, ss.Type)

						ss.Type = "object"
						ss.Properties = make(map[string]*openapi3.SchemaRef)
						ss.Properties[item.Title] = &openapi3.SchemaRef{Value: &item}
					}
				}

				s.Properties[tag.Get("json")] = &openapi3.SchemaRef{Value: &ss}
			}
		}
	}
	return
}
