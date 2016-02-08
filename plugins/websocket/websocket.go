package websocket

import (
	"fmt"
	"log"
	"net/http"

	"./../pluginbase"
	"github.com/gorilla/websocket"
	"github.com/robertkrimen/otto"
)

var upgrader = websocket.Upgrader{}

func InitPlugin(createVM func() (*otto.Otto, error)) *pluginbase.Plugin {
	p := pluginbase.Plugin{
		Name: "websocket",
		Init: func(vm *otto.Otto) {},
		HttpMapping: pluginbase.FuncMapping{
			"websocket": func(w http.ResponseWriter, r *http.Request) {
				doWebSocket(w, r, createVM)
			},
		},
	}

	return &p
}

func doWebSocket(w http.ResponseWriter, r *http.Request,
	createVM func() (*otto.Otto, error)) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	vm, err := createVM()

	fmt.Println("websocket start")

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		_, err = vm.Call("onWebSocketRequest", nil, string(message))
		if err != nil {
			log.Println("jserror", err)
		}

		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
	fmt.Println("websocket end")
}
