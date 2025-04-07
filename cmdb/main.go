package main

import (
	"log"
	"net/http"
)

func main() {
	initShell()
	if shell[0] == "" || shell[1] == "" {
		log.Fatal("shell not found")
	}
	var svc = newCommMap(handel)
	go svc.handleConnections()
	http.HandleFunc("/ws", svc.ServeHTTP)
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handel(ws ws, code uint16, msg string) {
	switch code {
	case CLIENT_EXEC_COMMAND:
		GoSafe(func() {
			connExecute(ws, msg)
		})
	case CLIENT_CANCEL_EXEC_COMMAND:
		GoSafe(func() {
			connExecuteCancel(ws, msg)
		})
	}
}
