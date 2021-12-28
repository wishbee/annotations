package wish

import (
	"bytes"
	"fmt"
	"go/ast"
	"strings"
)

// This file implements the fluent wish.
// Fluent wish generates the fluent pattern for a struct by creating
// SetXXX method which return the receiver as return value so that use can chain the set methods.

type fluent struct {
	multiField
}

func NewFluentWish(wishData *WishData, fileName, packageName, structName string, fields []*ast.Field) Wish {
	return &fluent{multiField{
		wishData:    wishData,
		fileName:    fileName,
		packageName: packageName,
		structName:  structName,
		fields:      fields,
	},
	}
}

func (r *fluent) FullFill() []byte {
	var buf bytes.Buffer
	set := func(fileName, packagName, structName string, field *ast.Field) []byte {
		fieldName := field.Names[0].Name
		fieldType := fieldType(field)
		s := fmt.Sprintf("func (s *%s)Set%s(v %s) *%s {\n"+
			"	s.%s = v\n"+
			"	return s\n"+
			"}\n", structName, strings.Title(fieldName), fieldType, structName, fieldName)
		return []byte(s)
	}

	for _, field := range r.fields {
		buf.Write(set(r.fileName, r.packageName, r.structName, field))
	}
	return buf.Bytes()
}
