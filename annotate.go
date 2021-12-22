package annotations

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)


func Process(dir string) {
	files := getSourceFiles(dir)
	for _, file := range files {
		f, err := readSourceFile(file)
		if err != nil {
			fmt.Printf("Skipping file: %s due to error: %s",file, err)
			continue
		}
		collectWishes(file,f)
	}
	fulfilWishes()
}



func getSourceFiles(dir string) []string {
	var files []string
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path,"_test.go") {
			fmt.Println("Go",path)
			files = append(files, path)
		}
		return nil
	})
	return files
}

func readSourceFile(file string) (f *ast.File, err error) {
	fs := token.NewFileSet()
	return parser.ParseFile(fs,file,nil, parser.ParseComments)
}

func collectWishes(file string, f *ast.File) {
	// Loop through the declarations.
	for _, d := range f.Decls {
		// If this declaration is GenDecl type then go further.
		if t,ok := d.(*ast.GenDecl); ok && t.Tok == token.TYPE {
			ts := t.Specs[0].(*ast.TypeSpec)
			st, ok := t.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType)
			if ok && !st.Incomplete {
				handleWishForStruct(file, f, t, ts, st)
			}
		}
	}
}

func handleWishForStruct(file string, f *ast.File, t *ast.GenDecl, ts *ast.TypeSpec, st *ast.StructType) {
	if t.Doc != nil {
		comment := t.Doc.Text()
		if hasAnyWish(comment) {
			wishes := getWishes(comment)
			fmt.Printf("Struct: %s, Wishes: %v\n", ts.Name, wishes)
			w := &wish{
				forFile:    file,
				forPackage: f.Name.Name,
				forStruct:  ts.Name.Name,
				wishes:     wishes,
				field:      st.Fields.List,
			}
			addWish(w)
		}
		// This declaration does not have wishes at struct level.
		// Check if any of the field has any wishes.
		if st.Fields != nil && len(st.Fields.List) > 0 {
			for _, field := range st.Fields.List {
				gc := field.Doc
				if gc != nil {
					comment := gc.Text()
					if hasAnyWish(comment) {
						wishes := getWishes(comment)
						fmt.Printf("Struct: %s, Field: %s, Wishes: %v\n", ts.Name, field.Names[0].Name, wishes)
						w := &wish{
							forFile:    file,
							forPackage: f.Name.Name,
							forStruct:  ts.Name.Name,
							wishes:     wishes,
							field:      field,
						}
						addWish(w)
					}
				}
			}
		}
	}
}

func hasAnyWish(comment string) bool {
	if idx := strings.Index(comment,"@wish:"); idx > 0 && len(comment)-idx > 1 {
		return true
	}
	return false
}

func getWishes(comment string) []string {
	var wishes []string
	idx := 0
	nextWish := func() bool {
		idx1 := strings.Index(comment[idx:],"@wish:")
		if idx1 > 0 {
			idx += idx1 + 6
			n := strings.IndexAny(comment[idx:]," \n\t");
			if n > 0 {
				wishes = append(wishes,comment[idx:idx+n])
				idx1 += n
				return true
			}
		}
		return false
	}

	for nextWish() {}
	return wishes
}

func fulfilWishes() {
	for _, wish := range wishingWell {
		for _, angelName := range wish.wishes {
			// Check if angel exists for this angelName.
			if angel, ok := angels[angelName]; ok {
				b := angel(wish.forFile, wish.forPackage,wish.forStruct,wish.field)
				fmt.Println(string(b))
			}else {
				panic(fmt.Sprintf("No angel for wish: %s\n",angelName))
			}
		}
	}
}