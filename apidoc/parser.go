package apidoc

import (
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

var root = ".."

var _paths openapi3.Paths

func init() {
	_paths = make(openapi3.Paths)
}

// Parse ...
func Parse(filePath, prefix string) (err error) {
	fset := token.NewFileSet()

	filePath = strings.ReplaceAll(filePath, "\"", "")
	f, err := parser.ParseFile(fset, filePath+"/routes.go", nil, parser.Trace)
	if err != nil {
		return err
	}

	docs, err := ParseHandleDocs(filePath)
	if err != nil {
		return err
	}

	for _, obj := range f.Scope.Objects {
		if obj.Name == "Mount" {
			routes := make(map[string]string)
			for _, decl := range obj.Decl.(*ast.FuncDecl).Body.List {
				// 解析：ar := r.Group("api", DebugLog)
				if stmt, ok := decl.(*ast.AssignStmt); ok {
					expr := stmt.Rhs[0].(*ast.CallExpr)
					slct := expr.Fun.(*ast.SelectorExpr)
					if slct.Sel.Name != "Group" {
						continue
					}

					router := stmt.Lhs[0].(*ast.Ident).Name
					pattern := expr.Args[0].(*ast.BasicLit).Value
					routes[router] = joinPaths(prefix, pattern[1:len(pattern)-1])
				}

				if stmt, ok := decl.(*ast.ExprStmt); ok {
					expr := stmt.X.(*ast.CallExpr)
					slct := expr.Fun.(*ast.SelectorExpr)

					// 解析：admin.Mount(ar)
					if slct.Sel.Name == "Mount" {
						router := expr.Args[0].(*ast.Ident).Name
						if routes[router] != "" {
							subRouter := slct.X.(*ast.Ident).Name

							for _, imp := range f.Imports {
								name := ""
								if imp.Name != nil {
									name = imp.Name.Name
								}

								importPath := imp.Path.Value[1 : len(imp.Path.Value)-1]
								name = filepath.Base(importPath)

								if name == subRouter {
									index := strings.Index(importPath, "/")
									path := root + importPath[index:]
									err = Parse(path, routes[router])
									if err != nil {
										return err
									}
								}
							}
						}
					}

					// 解析：admin.Mount(ar)
					if strings.IndexAny("GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE", slct.Sel.Name) != -1 {
						funcName := expr.Args[1].(*ast.Ident).Name
						prms, body, resp, err := ParseDto(filePath, funcName)
						if err != nil {
							return err
						}
						oper := openapi3.Operation{}
						oper.Summary = docs[funcName]
						oper.OperationID = funcName
						oper.Tags = append(oper.Tags, filepath.Base(prefix))
						if prms != nil {
							oper.Parameters = *prms
						}
						if body != nil {
							oper.RequestBody = &openapi3.RequestBodyRef{Value: body}
						}
						oper.Responses = resp

						pi := openapi3.PathItem{}
						pi.Summary = docs[funcName]
						switch slct.Sel.Name {
						case http.MethodGet:
							pi.Get = &oper
						case http.MethodPost:
							pi.Post = &oper
						case http.MethodDelete:
							pi.Delete = &oper
						case http.MethodPut:
							pi.Put = &oper
						}

						pattern := expr.Args[0].(*ast.BasicLit).Value
						pattern = pattern[1 : len(pattern)-1]
						_paths[joinPaths(prefix, pattern)] = &pi
					}
				}
			}
		}
	}

	return err
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func joinPaths(absolutePath, relativePath string) string {
	if absolutePath[0] != '/' {
		absolutePath = "/" + absolutePath
	}
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}

// ParseHandleDocs ...
func ParseHandleDocs(filePath string) (map[string]string, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, filePath+"/handle.go", nil, parser.ParseComments)
	if err != nil {
		return nil, nil
	}

	docs := make(map[string]string)
	for _, obj := range f.Scope.Objects {
		fd, ok := obj.Decl.(*ast.FuncDecl)
		if ok {
			docText := fd.Doc.Text()
			docText = strings.TrimPrefix(docText, obj.Name)
			docText = strings.TrimPrefix(docText, " ")
			docs[obj.Name] = docText
		}
	}

	return docs, nil
}

// ParseDto ...
func ParseDto(filePath, funcName string) (
	prms *openapi3.Parameters, body *openapi3.RequestBody, resp openapi3.Responses, err error) {
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

	sche := getRespOject()

	respSche := parseSchema(f, funcName+"Response")
	if len(respSche.Properties) > 0 {
		sche.Properties["data"] = &openapi3.SchemaRef{Value: &respSche}
	}
	resp = make(openapi3.Responses)
	cnt := openapi3.Content{}
	cnt["application/json"] = &openapi3.MediaType{Schema: &openapi3.SchemaRef{Value: &sche}}
	resp["200"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Content: cnt,
		},
	}

	return
}

func getRespOject() openapi3.Schema {
	sche := openapi3.Schema{Type: "object"}
	sche.Description = "返回业务数据"

	sche.Properties = make(map[string]*openapi3.SchemaRef)
	sche.Properties["code"] = openapi3.NewSchemaRef("", &openapi3.Schema{Type: "integer", Description: "错误编码，0为成功"})
	sche.Properties["message"] = openapi3.NewSchemaRef("", &openapi3.Schema{Type: "string", Description: "错误消息，OK为成功"})
	detalis := &openapi3.Schema{Type: "array", Description: "错误详情"}
	detalis.Items = &openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}}
	sche.Properties["details"] = openapi3.NewSchemaRef("", detalis)

	return sche
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
				if tag.Get("json") == "" {
					continue
				}

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
