package annotations

import (
	"bytes"
	"fmt"
	"github.com/wishbee/annotations/wish"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Process() {
	files := getSourceFiles(*Option.Dir)
	for _, file := range files {
		if *Option.Verbose {
			fmt.Printf("Reading source file: %s\n", file)
		}
		f, err := readSourceFile(file)
		if err != nil {
			if *Option.AbortOnError {
				fmt.Printf("Aborting due to error when reading source file: %s, error: %s", file, err)
				os.Exit(1)
			}
			fmt.Printf("Skipping file: %s due to error: %s", file, err)
			continue
		}
		wishes := collectWishes(file, f)
		if wishes != nil {
			fulfilWishes(wishes, file, f.Name.Name)
		}
	}
}

func getSourceFiles(dir string) []string {
	var files []string
	ignoreFilesWithPattern := []string{"_test.go", "_wish.go"}
	if *Option.IgnoreFiles != "" {
		userPatterns := strings.Split(*Option.IgnoreFiles, ",")
		for _, userPattern := range userPatterns {
			ignoreFilesWithPattern = append(ignoreFilesWithPattern, strings.TrimSpace(userPattern))
		}
	}
	dontIgnore := func(path string) bool {
		for _, pattern := range ignoreFilesWithPattern {
			if strings.HasSuffix(path, pattern) {
				return false
			}
		}
		return true
	}
	_ = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") && dontIgnore(path) {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func readSourceFile(file string) (f *ast.File, err error) {
	fSet := token.NewFileSet()
	return parser.ParseFile(fSet, file, nil, parser.ParseComments)
}

func collectWishes(file string, f *ast.File) []wish.Wish {
	// Loop through the declarations.
	var wishingWell []wish.Wish
	for _, d := range f.Decls {
		// If this declaration is GenDecl type then go further.
		if t, ok := d.(*ast.GenDecl); ok && t.Tok == token.TYPE {
			ts := t.Specs[0].(*ast.TypeSpec)
			st, ok := t.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType)
			if ok && !st.Incomplete {
				wishes := findWishesForStruct(file, f, t, ts, st)
				if wishes != nil {
					for _, v := range wishes {
						wishingWell = append(wishingWell, v)
					}
				}
			}
		}
	}
	return wishingWell
}

func findWishesForStruct(file string, f *ast.File, t *ast.GenDecl, ts *ast.TypeSpec, st *ast.StructType) []wish.Wish {
	var wishingWell []wish.Wish
	if t.Doc != nil {
		comment := t.Doc.Text()
		if hasAnyWish(comment) {
			wishes := getWishes(comment)
			if *Option.Verbose {
				fmt.Printf("Found Wishes: %v for struct: %s in file: %s\n", wishes, ts.Name, file)
			}
			var w wish.Wish
			for _, wishName := range wishes {
				if !wish.IsValidWishName(wishName) {
					errText := fmt.Sprintf("encountered invalid wish: %s in file: %s for struct: %s", wishName, file, ts.Name)
					if *Option.AbortOnError {
						panic(errText)
					}
					fmt.Println(errText)
				}
				if wishName == "setter" {
					w = wish.NewSetterWish(file, f.Name.Name, ts.Name.Name, st.Fields.List)
					wishingWell = append(wishingWell, w)
					continue
				}
				if wishName == "getter" {
					w = wish.NewGetterWish(file, f.Name.Name, ts.Name.Name, st.Fields.List)
					wishingWell = append(wishingWell, w)
					continue
				}
				if wishName == "fluent" {
					w = wish.NewFluentWish(file, f.Name.Name, ts.Name.Name, st.Fields.List)
					wishingWell = append(wishingWell, w)
					continue
				}
			}
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
						if *Option.Verbose {
							fmt.Printf("Found Wishes: %v for field: %s.%s in file: %s\n", wishes, ts.Name, field.Names[0].Name, file)
						}
						var w wish.Wish
						singleField := make([]*ast.Field, 1)
						singleField[0] = field
						for _, wishName := range wishes {
							if wishName == "setter" {
								w = wish.NewSetterWish(file, f.Name.Name, ts.Name.Name, singleField)
								wishingWell = append(wishingWell, w)
								continue
							}
							if wishName == "getter" {
								w = wish.NewGetterWish(file, f.Name.Name, ts.Name.Name, singleField)
								wishingWell = append(wishingWell, w)
								continue
							}
						}
					}
				}
			}
		}
	}
	return wishingWell
}

func hasAnyWish(comment string) bool {
	if idx := strings.Index(comment, "@wish:"); idx > 0 && len(comment)-idx > 1 {
		return true
	}
	return false
}

func getWishes(comment string) []string {
	var wishes []string
	idx := 0
	nextWish := func() bool {
		idx1 := strings.Index(comment[idx:], "@wish:")
		if idx1 > 0 {
			idx += idx1 + 6
			n := strings.IndexAny(comment[idx:], " \n\t")
			if n > 0 {
				wishes = append(wishes, comment[idx:idx+n])
				idx1 += n
				return true
			}
		}
		return false
	}

	for nextWish() {
	}
	return wishes
}

func fulfilWishes(wishingWell []wish.Wish, fileName, packageNAme string) {
	var output bytes.Buffer
	for _, v := range wishingWell {
		b := v.FullFill()
		if b != nil {
			output.Write(b)
		}
	}
	if output.Len() > 0 {
		topLine := fmt.Sprintf("// Code generated by WishGen. DO NOT EDIT.\n// Source: %s\npackage %s\n", fileName, packageNAme)
		if *Option.OutputToStdOutOnly {
			fmt.Print(topLine)
			fmt.Print(output.String())
			return
		}
		dir, file := filepath.Split(fileName)
		if dir == "" {
			dir = "."
		}
		newFile := strings.Split(file, ".")[0] + "_wish.go"
		newFile = dir + string(filepath.Separator) + newFile
		outFile, err := os.OpenFile(newFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			errText := fmt.Sprintf("unable to create outfile: %s. Error: %s", newFile, err)
			if *Option.AbortOnError {
				panic(errText)
			}
			fmt.Println(errText)
		}
		_, _ = outFile.WriteString(topLine)
		_, _ = outFile.WriteString(output.String())
		_ = outFile.Close()
		if *Option.Verbose {
			fmt.Printf("Wish completed for source file %s. See %s for output.", fileName, newFile)
		}
	}
}
