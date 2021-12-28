package wish

import (
	"go/ast"
	"strings"
)

// Wish is an interface to hide the generation logic behind it.
type Wish interface {
	FullFill() []byte
	Name() string
	Data() string
	File() string
	Package() string
	Struct() string
	Fields() []*ast.Field
}

var supportedWishes = []string{
	"getter",
	"setter",
	"fluent",
}

type WishData struct {
	Name string
	Data string
}

func ParseWishData(s string) *WishData {
	w := &WishData{}
	t := strings.Split(s, ":")
	w.Name = t[0]
	if len(t) >= 1 {
		w.Data = t[0]
	}
	return w
}

// IsValidWishName checks if the wishName is valid or not.
func IsValidWishName(wishName string) bool {
	for _, name := range supportedWishes {
		if name == wishName {
			return true
		}
	}
	return false
}

type multiField struct {
	wishData    *WishData
	fileName    string
	packageName string
	structName  string
	fields      []*ast.Field
}

func (m *multiField) File() string {
	return m.fileName
}

func (m *multiField) Package() string {
	return m.packageName
}

func (m *multiField) Struct() string {
	return m.structName
}

func (m *multiField) Fields() []*ast.Field {
	return m.fields
}

func (r *multiField) Name() string {
	return r.wishData.Name
}

func (r *multiField) Data() string {
	return r.wishData.Data
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

func MakeAStructLevelWish(wishData *WishData, fileName, packageName, structName string, fields []*ast.Field) Wish {
	switch wishData.Name {
	case "setter":
		return NewSetterWish(wishData, fileName, packageName, structName, fields)
	case "getter":
		return NewGetterWish(wishData, fileName, packageName, structName, fields)
	case "fluent":
		return NewFluentWish(wishData, fileName, packageName, structName, fields)
	}
	return nil
}

func MakeAFieldLevelWish(wishData *WishData, fileName, packageName, structName string, field *ast.Field) Wish {
	singleField := make([]*ast.Field, 1)
	singleField[0] = field
	switch wishData.Name {
	case "setter":
		return NewSetterWish(wishData, fileName, packageName, structName, singleField)
	case "getter":
		return NewGetterWish(wishData, fileName, packageName, structName, singleField)
	}
	return nil
}
