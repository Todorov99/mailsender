package server

import (
	"fmt"
	"net/http"

	"github.com/Todorov99/mailsender/pkg/controller"
	"github.com/Todorov99/mailsender/pkg/global"
	"github.com/Todorov99/mailsender/pkg/server/config"
	"github.com/gorilla/mux"
)

// HandleRequest handles the supported REST request of the Web Server
func HandleRequest(port string) error {
	routes := mux.NewRouter().StrictSlash(true)

	senderController, err := controller.NewSenderController()
	if err != nil {
		return err
	}

	routes.HandleFunc("/api/mail/attachment/send", senderController.SendWithAttachments).Methods("POST")
	routes.HandleFunc("/api/mail/send", senderController.Send).Methods("POST")

	if config.GetTLSCfg() == nil {
		return http.ListenAndServe(fmt.Sprintf(":%s", port), routes)
	}

	return http.ListenAndServeTLS(fmt.Sprintf(":%s", port), fmt.Sprintf("%s/%s", global.CertificatesPath, config.GetTLSCfg().CertFile), fmt.Sprintf("%s/%s", global.CertificatesPath, config.GetTLSCfg().PrivateKey), routes)
}
