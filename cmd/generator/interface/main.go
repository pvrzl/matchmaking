package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"unicode"
)

func formatFieldList(fields *ast.FieldList) string {
	if fields == nil {
		return ""
	}
	var fieldStrings []string
	for _, field := range fields.List {
		var names []string
		for _, name := range field.Names {
			names = append(names, name.Name)
		}
		var types []string
		switch typ := field.Type.(type) {
		case *ast.Ident:
			types = append(types, typ.Name)
		case *ast.SelectorExpr:
			types = append(types, fmt.Sprintf("%s.%s", typ.X.(*ast.Ident).Name, typ.Sel.Name))
		case *ast.ArrayType:
			types = append(types, fmt.Sprintf("[]%s", formatFieldList(&ast.FieldList{List: []*ast.Field{{Type: typ.Elt}}})))
		case *ast.MapType:
			types = append(types, fmt.Sprintf("map[%s]%s", formatFieldList(&ast.FieldList{List: []*ast.Field{{Type: typ.Key}}}), formatFieldList(&ast.FieldList{List: []*ast.Field{{Type: typ.Value}}})))
		case *ast.FuncType:
			types = append(types, formatFieldList(typ.Params)+" -> "+formatFieldList(typ.Results))
		default:
			types = append(types, "unknown")
		}
		fieldStrings = append(fieldStrings, fmt.Sprintf("%s %s", strings.Join(names, ", "), strings.Join(types, ", ")))
	}
	return strings.Join(fieldStrings, ", ")
}

func formatFuncDecl(funcDecl *ast.FuncDecl) string {
	paramStr := formatFieldList(funcDecl.Type.Params)
	resultStr := formatFieldList(funcDecl.Type.Results)

	if resultStr != "" {
		return fmt.Sprintf("%s(%s) (%s)", funcDecl.Name.Name, paramStr, resultStr)
	}
	return fmt.Sprintf("%s(%s)", funcDecl.Name.Name, paramStr)
}

func getReceiverType(recv *ast.FieldList) string {
	if recv == nil || len(recv.List) == 0 {
		return ""
	}
	// Handle cases where receiver type is a pointer or struct
	switch expr := recv.List[0].Type.(type) {
	case *ast.Ident:
		return expr.Name // Direct struct name
	case *ast.StarExpr:
		if ident, ok := expr.X.(*ast.Ident); ok {
			return ident.Name // Dereferenced struct name
		}
	}
	return ""
}

func capitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func main() {
	fmt.Print("Enter the path to the Go source file: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("Error reading input.")
		os.Exit(1)
	}
	filePath := scanner.Text()

	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}
	sourceCode := string(file)

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", sourceCode, parser.AllErrors)
	if err != nil {
		fmt.Printf("Error parsing source code: %s\n", err)
		os.Exit(1)
	}

	methodsByStruct := make(map[string][]string)

	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				receiverType := getReceiverType(funcDecl.Recv)
				if receiverType != "" && ast.IsExported(funcDecl.Name.Name) {
					methodsByStruct[receiverType] = append(methodsByStruct[receiverType], formatFuncDecl(funcDecl))
				}
			}
		}
	}

	for structName, methods := range methodsByStruct {
		exportedStructName := capitalizeFirstLetter(structName)
		fmt.Printf("type %sInterface interface {\n", exportedStructName)
		for _, method := range methods {
			fmt.Printf("    %s\n", method)
		}
		fmt.Println("}")
	}

	if len(methodsByStruct) == 0 {
		fmt.Println("No exported struct methods found in the provided Go source file.")
	}
}
