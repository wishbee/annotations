package wish

import (
	"bytes"
	"fmt"
	"go/ast"
	"strings"
)

type setter struct {
	multiField
}

func NewSetterWish(fileName, packageName, structName string, fields []*ast.Field) Wish {
	return &setter{multiField{
		fileName:    fileName,
		packageName: packageName,
		structName:  structName,
		fields:      fields,
	},
	}
}

func (r *setter) Name() string {
	return "setter"
}

func (r *setter) FullFill() []byte {
	set := func(fileName, packagName, structName string, field *ast.Field) []byte {
		fieldName := field.Names[0].Name
		fieldType := fieldType(field)
		s := fmt.Sprintf("func (s *%s)Set%s(v %s) {\n"+
			"	s.%s = v\n"+
			"}\n", structName, strings.Title(fieldName), fieldType, fieldName)
		return []byte(s)
	}
	var buf bytes.Buffer
	for _, field := range r.fields {
		buf.Write(set(r.fileName, r.packageName, r.structName, field))
	}
	return buf.Bytes()
}
