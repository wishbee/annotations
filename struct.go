package annotations

import (
	"bytes"
	"fmt"
	"go/ast"
)

func init(){
	if angels == nil {
		angels = make(map[string]angel)
	}
	angels["setters"] = setter
	angels["getters"] = getter
}

func setter(fileName, packagName, structName string, fields interface{}) []byte {
	f, ok := fields.([]*ast.Field)
	if !ok {
		panic(fmt.Sprintf("expected 'fields' argument as type '[]*ast.Field' for angel 'setter' but got as '%T'",fields))
	}
	var buf bytes.Buffer
	for _, field := range f {
		buf.Write(set(fileName, packagName, structName, field))
	}
	return buf.Bytes()
}

func getter(fileName, packagName, structName string, fields interface{}) []byte {
	f, ok := fields.([]*ast.Field)
	if !ok {
		panic(fmt.Sprintf("expected 'fields' argument as type '[]*ast.Field' for angel 'getter' but got as '%T'",fields))
	}
	var buf bytes.Buffer
	for _, field := range f {
		buf.Write(get(fileName, packagName, structName, field))
	}
	return buf.Bytes()
}


