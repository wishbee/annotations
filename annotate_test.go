package annotations

import (
	"fmt"
	"go/ast"
	"strings"
	"testing"
)

type a struct{
	b int
}

func (r a)set(i int)  {
	r.b = i
}
func (r a)get()int  {
	return r.b
}

func TestProcess(t *testing.T){
	Process("/Users/lower/work/src/go/samples")
}


func Comments(f *ast.File) {
	for _, d := range f.Decls {
		if t,ok := d.(*ast.GenDecl); ok {
			fmt.Printf("%s, Comment: %s\n",TypeSpecs(t.Specs),t.Doc.Text())
		}
	}
}
func TypeSpecs(s []ast.Spec) string {
	r := ""
	for _,a := range s {
		if t,ok := a.(*ast.TypeSpec); ok {
			r += fmt.Sprintf("Name: %s, Type: %s", t.Name, Type(t))
		}
	}
	return r
}

func Type(t *ast.TypeSpec) string {
	var r string
	switch t.Type.(type) {
	case *ast.StructType:
		r += "struct"
	case *ast.InterfaceType:
		r += "interface"
	default:
		r += fmt.Sprintf("%T",t.Type)
	}
	return r
}

type class int

func (v class)Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	fmt.Printf("%s%T\n",strings.Repeat(" ",int(v)),n)
	switch d := n.(type) {
	case *ast.File:
		fmt.Printf("File (Package): %s\n",d.Name.Name)
	case *ast.TypeSpec:
		fmt.Printf("TypeSpec (Type): %s\n",d.Name.Name)
	case *ast.CommentGroup:
		fmt.Printf("CommentGroup (Type): %v\n",d.List)
	}
	return v + 1
}

