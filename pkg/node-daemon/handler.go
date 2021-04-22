package node_daemon

import "fmt"

type PrintHandler struct {}

func NewPrintHandler() *PrintHandler {
	return &PrintHandler{}
}

func (ph *PrintHandler) DoAdded(obj map[string]interface{}) {
	fmt.Println(obj)
}

func (ph *PrintHandler) DoModified(obj map[string]interface{}) {
	fmt.Println(obj)
}

func (ph *PrintHandler) DoDeleted(obj map[string]interface{}) {
	fmt.Println(obj)
}