package goimportsdeps

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strconv"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	goimportsdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
)

// Bind fills deps.Deps.Goimportsdeps, providing the capability to parse a Go
// source file into its package clause, imports and top-level declarations.
func Bind(deps *deps.Deps) {
	deps.Goimportsdeps = goimportsdeps.Lib{
		Parse:          parse,
		GetPackageName: getPackageName,
		GetImports:     getImports,
	}
}

// getImports fills goimportsdeps.Lib.GetImports, parsing only the import
// section of the given Go source and unquoting each import path.
func getImports(content string) ([]string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	imports := make([]string, 0, len(file.Imports))
	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return nil, err
		}
		imports = append(imports, path)
	}

	return imports, nil
}

// getPackageName fills goimportsdeps.Lib.GetPackageName, parsing only the
// package clause of the given Go source.
func getPackageName(content string) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.PackageClauseOnly)
	if err != nil {
		return "", err
	}
	return file.Name.Name, nil
}

// parse fills goimportsdeps.Lib.Parse, parsing the whole Go source (with
// comments) and flattening its declarations into a goimportsdeps.File.
func parse(content string) (*goimportsdeps.File, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	result := &goimportsdeps.File{
		Package:   file.Name.Name,
		Doc:       docText(file.Doc),
		Imports:   make([]goimportsdeps.Import, 0, len(file.Imports)),
		Functions: make([]goimportsdeps.Function, 0),
		Types:     make([]goimportsdeps.Type, 0),
		Constants: make([]goimportsdeps.Value, 0),
		Variables: make([]goimportsdeps.Value, 0),
	}

	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return nil, err
		}
		alias := ""
		if spec.Name != nil {
			alias = spec.Name.Name
		}
		result.Imports = append(result.Imports, goimportsdeps.Import{Alias: alias, Path: path})
	}

	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			result.Functions = append(result.Functions, funcFromDecl(fset, d))
		case *ast.GenDecl:
			switch d.Tok {
			case token.TYPE:
				for _, spec := range d.Specs {
					ts := spec.(*ast.TypeSpec)
					result.Types = append(result.Types, typeFromSpec(fset, d, ts))
				}
			case token.CONST:
				result.Constants = append(result.Constants, valuesFromDecl(fset, d)...)
			case token.VAR:
				result.Variables = append(result.Variables, valuesFromDecl(fset, d)...)
			}
		}
	}

	return result, nil
}

// docText trims a doc comment group to its plain text, "" when absent.
func docText(group *ast.CommentGroup) string {
	if group == nil {
		return ""
	}
	return strings.TrimSpace(group.Text())
}

// exprString renders an AST type expression back to Go source, "" when nil.
func exprString(fset *token.FileSet, expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, expr); err != nil {
		return ""
	}
	return buf.String()
}

// isExported reports whether a name begins with an upper-case letter.
func isExported(name string) bool {
	return name != "" && ast.IsExported(name)
}

// funcFromDecl converts a function or method declaration.
func funcFromDecl(fset *token.FileSet, decl *ast.FuncDecl) goimportsdeps.Function {
	fn := goimportsdeps.Function{
		Name:     decl.Name.Name,
		Doc:      docText(decl.Doc),
		Params:   fieldsToParams(fset, decl.Type.Params),
		Results:  fieldsToParams(fset, decl.Type.Results),
		Exported: isExported(decl.Name.Name),
	}
	if decl.Recv != nil && len(decl.Recv.List) > 0 {
		recv := decl.Recv.List[0].Type
		if star, ok := recv.(*ast.StarExpr); ok {
			fn.Pointer = true
			recv = star.X
		}
		fn.Receiver = exprString(fset, recv)
	}
	return fn
}

// fieldsToParams flattens a parameter or result list, one entry per name.
func fieldsToParams(fset *token.FileSet, list *ast.FieldList) []goimportsdeps.Param {
	params := make([]goimportsdeps.Param, 0)
	if list == nil {
		return params
	}
	for _, field := range list.List {
		typ := exprString(fset, field.Type)
		if len(field.Names) == 0 {
			params = append(params, goimportsdeps.Param{Type: typ})
			continue
		}
		for _, name := range field.Names {
			params = append(params, goimportsdeps.Param{Name: name.Name, Type: typ})
		}
	}
	return params
}

// typeFromSpec converts one type spec, classifying its kind.
func typeFromSpec(fset *token.FileSet, decl *ast.GenDecl, spec *ast.TypeSpec) goimportsdeps.Type {
	t := goimportsdeps.Type{
		Name:     spec.Name.Name,
		Doc:      typeDoc(decl, spec),
		Fields:   make([]goimportsdeps.Field, 0),
		Methods:  make([]goimportsdeps.Function, 0),
		Exported: isExported(spec.Name.Name),
	}
	switch underlying := spec.Type.(type) {
	case *ast.StructType:
		t.Kind = "struct"
		t.Fields = structFields(fset, underlying)
	case *ast.InterfaceType:
		t.Kind = "interface"
		t.Methods = interfaceMethods(fset, underlying)
	default:
		if spec.Assign.IsValid() {
			t.Kind = "alias"
		} else {
			t.Kind = "other"
		}
		t.Underlying = exprString(fset, spec.Type)
	}
	return t
}

// typeDoc picks the doc comment for a type spec, falling back to the GenDecl's
// own comment for a single-spec block.
func typeDoc(decl *ast.GenDecl, spec *ast.TypeSpec) string {
	if spec.Doc != nil {
		return docText(spec.Doc)
	}
	if len(decl.Specs) == 1 {
		return docText(decl.Doc)
	}
	return ""
}

// structFields flattens a struct's fields, one entry per name.
func structFields(fset *token.FileSet, st *ast.StructType) []goimportsdeps.Field {
	fields := make([]goimportsdeps.Field, 0)
	if st.Fields == nil {
		return fields
	}
	for _, field := range st.Fields.List {
		typ := exprString(fset, field.Type)
		tag := ""
		if field.Tag != nil {
			if unquoted, err := strconv.Unquote(field.Tag.Value); err == nil {
				tag = unquoted
			} else {
				tag = field.Tag.Value
			}
		}
		doc := docText(field.Doc)
		if len(field.Names) == 0 {
			fields = append(fields, goimportsdeps.Field{
				Type:     typ,
				Tag:      tag,
				Doc:      doc,
				Exported: isExported(embeddedName(typ)),
			})
			continue
		}
		for _, name := range field.Names {
			fields = append(fields, goimportsdeps.Field{
				Name:     name.Name,
				Type:     typ,
				Tag:      tag,
				Doc:      doc,
				Exported: isExported(name.Name),
			})
		}
	}
	return fields
}

// embeddedName reduces an embedded field's type expression to its bare name.
func embeddedName(typ string) string {
	typ = strings.TrimPrefix(typ, "*")
	if i := strings.LastIndex(typ, "."); i >= 0 {
		typ = typ[i+1:]
	}
	return typ
}

// interfaceMethods converts an interface's explicit methods; embedded
// interfaces are skipped.
func interfaceMethods(fset *token.FileSet, it *ast.InterfaceType) []goimportsdeps.Function {
	methods := make([]goimportsdeps.Function, 0)
	if it.Methods == nil {
		return methods
	}
	for _, member := range it.Methods.List {
		fnType, ok := member.Type.(*ast.FuncType)
		if !ok || len(member.Names) == 0 {
			continue
		}
		name := member.Names[0].Name
		methods = append(methods, goimportsdeps.Function{
			Name:     name,
			Doc:      docText(member.Doc),
			Params:   fieldsToParams(fset, fnType.Params),
			Results:  fieldsToParams(fset, fnType.Results),
			Exported: isExported(name),
		})
	}
	return methods
}

// valuesFromDecl flattens the names of a const or var block.
func valuesFromDecl(fset *token.FileSet, decl *ast.GenDecl) []goimportsdeps.Value {
	values := make([]goimportsdeps.Value, 0)
	for _, spec := range decl.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		typ := exprString(fset, vs.Type)
		doc := docText(vs.Doc)
		if doc == "" && len(decl.Specs) == 1 {
			doc = docText(decl.Doc)
		}
		for index, name := range vs.Names {
			value := ""
			if index < len(vs.Values) {
				value = exprString(fset, vs.Values[index])
			}
			values = append(values, goimportsdeps.Value{
				Name:     name.Name,
				Doc:      doc,
				Type:     typ,
				Value:    value,
				Exported: isExported(name.Name),
			})
		}
	}
	return values
}
