package annotations

type wish struct {
	forStruct string
	wishes []string
	field  interface{}
}


type angel func(fileName, packagName, structName string, field interface{}) []byte

var angels map[string]angel


