package server

import (
	"fmt"
	"net/http"

	"github.com/Todorov99/mailsender/pkg/controller"
	"github.com/gorilla/mux"
)

// HandleRequest handles the supported REST request of the Web Server
func HandleRequest(port string) error {
	routes := mux.NewRouter().StrictSlash(true)

	routes.HandleFunc("/api/mail/send", controller.Send).Methods("Get")

	return http.ListenAndServe(fmt.Sprintf(":%s", port), routes)
}
