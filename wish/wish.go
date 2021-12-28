package wish

import "go/ast"

// Wish is an interface to hide the generation logic behind it.
type Wish interface {
	FullFill() []byte
	Name() string
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
