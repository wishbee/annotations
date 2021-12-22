package annotations

import (
	"sync"
)

type wish struct {
	forFile string
	forPackage string
	forStruct string
	wishes []string
	field  interface{}
}

var wellLock sync.Mutex

var wishingWell []*wish

type angel func(fileName, packagName, structName string, field interface{}) []byte

var angels map[string]angel

func init() {
	wishingWell = make([]*wish,0)
}

func addWish(w *wish){
	wellLock.Lock()
	defer wellLock.Unlock()
	wishingWell = append(wishingWell, w)
}