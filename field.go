package annotations

import (
	"fmt"
	"go/ast"
	"strings"
)
func init(){
	if angels == nil {
		angels = make(map[string]angel)
	}
	angels["set"] = set
	angels["get"] = get
}

func set(fileName, packagName, structName string, field interface{}) []byte {
	f, ok := field.(*ast.Field)
	if !ok {
		panic(fmt.Sprintf("expected 'field' argument as type '*ast.Field' for angel 'set' but got as '%T'",field))
	}
	fieldName := f.Names[0].Name
	fieldType := fieldType(f)
	s := fmt.Sprintf("func (s *%s)Set%s(v %s) {\n"+
		"	s.%s = %s\n"+
		"}\n", structName, strings.Title(fieldName), fieldType, fieldName, fieldName)
	//fmt.Println(s)
	return []byte(s)
}

func get(fileName, packagName, structName string, field interface{}) []byte {
	f, ok := field.(*ast.Field)
	if !ok {
		panic(fmt.Sprintf("expected 'field' argument as type '*ast.Field' for angel 'get' but got as '%T'",field))
	}
	fieldName := f.Names[0].Name
	fieldType := fieldType(f)
	s := fmt.Sprintf("func (s *%s)Get%s() %s {\n"+
		"	return s.%s\n"+
		"}\n", structName, strings.Title(fieldName), fieldType, fieldName)
	//fmt.Println(s)
	return []byte(s)
}

func fieldType(field *ast.Field) string {
	fieldType := "unknownFieldType"
	ident, ok := field.Type.(*ast.Ident)
	if ok {
		fieldType = ident.Name
	}
	starExpr, ok := field.Type.(*ast.StarExpr)
	if ok {
		ident, ok := starExpr.X.(*ast.Ident)
		if ok {
			fieldType = "*" + ident.Name
		}
	}
	return fieldType
}