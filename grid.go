package main

import (
	"container/list"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"code.google.com/p/go.net/websocket"
)

const listenAddr = ":4000"

func main() {
	http.HandleFunc("/", rootHandler)
	http.Handle("/socket", websocket.Handler(socketHandler))
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// type subscriberQueue chan []byte
// var subscriberQueues = make([](chan *websocket.Conn), 100)

var clients = list.New()

// type subscriber chan []byte
// var subscribers = make([]subscriber, 100)

func socketHandler(ws *websocket.Conn) {
	fmt.Printf("readWriteServer %#v\n", ws.Config())

	clientListEl := clients.PushBack(ws)

	for {
		buf := make([]byte, 100)
		n, err := ws.Read(buf)
		if err != nil {
			fmt.Println(err)
			break
		}
		go handleMessage(buf[:n], ws)
	}

	clients.Remove(clientListEl)

	fmt.Println("readWriteServer finished")
}

var updateRegExp, _ = regexp.Compile("^UPDATE\\:([0-9]+)\\,([0-9]+)\\:(.)$")

func handleMessage(message []byte, origin *websocket.Conn) {
	fmt.Printf("recv:%q\n", message)

	submatches := updateRegExp.FindSubmatch(message)
	if submatches != nil {
		fmt.Printf("x:    %q\n", submatches[1])
		fmt.Printf("y:    %q\n", submatches[2])
		fmt.Printf("char: %q\n", submatches[3])
	}

	for e := clients.Front(); e != nil; e = e.Next() {
		// do something with e.Value
		conn := e.Value.(*websocket.Conn)

		// if conn != origin {
		websocket.Message.Send(conn, fmt.Sprintf("%s", message))
		// }
	}
}
