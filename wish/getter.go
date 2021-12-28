package wish

import (
	"bytes"
	"fmt"
	"go/ast"
	"strings"
)

type getter struct {
	multiField
}

func NewGetterWish(wishData *WishData, fileName, packageName, structName string, fields []*ast.Field) Wish {
	return &getter{multiField{
		wishData:    wishData,
		fileName:    fileName,
		packageName: packageName,
		structName:  structName,
		fields:      fields,
	},
	}
}

func (r *getter) FullFill() []byte {
	get := func(fileName, packagName, structName string, field *ast.Field) []byte {
		fieldName := field.Names[0].Name
		fieldType := fieldType(field)
		s := fmt.Sprintf("func (s *%s)%s() %s {\n"+
			"	return s.%s\n"+
			"}\n", structName, strings.Title(fieldName), fieldType, fieldName)
		return []byte(s)
	}
	var buf bytes.Buffer
	for _, field := range r.fields {
		buf.Write(get(r.fileName, r.packageName, r.structName, field))
	}
	return buf.Bytes()
}
