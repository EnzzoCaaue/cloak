package web

import (
	"github.com/gorilla/context"
	"log"
	"net/http"
)

// Start starts the AAC http server
func Start(port, path string) {
	r := newRouter()
	r.ServeFiles("/public/*filepath", http.Dir(path+"/public/"))
	log.Println("Cloak AAC listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, context.ClearHandler(r)))
}
