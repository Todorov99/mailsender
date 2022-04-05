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

	senderController, err := controller.NewSenderController()
	if err != nil {
		return err
	}

	routes.HandleFunc("/api/mail/attachment/send", senderController.SendWithAttachment).Methods("POST")
	routes.HandleFunc("/api/mail/send", senderController.Send).Methods("POST")

	return http.ListenAndServe(fmt.Sprintf(":%s", port), routes)
}
