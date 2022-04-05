package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/Todorov99/mailsender/pkg/dto"
	"github.com/Todorov99/mailsender/pkg/service"
	"github.com/Todorov99/sensorcli/pkg/logger"
)

var controllerLogger = logger.NewLogrus("controller", os.Stdout)

type senderController struct {
	mailSenderService service.MailSenderService
}

func NewSenderController() (senderController, error) {
	mailSenderService, err := service.NewMailSender()
	if err != nil {
		return senderController{}, err
	}

	return senderController{
		mailSenderService: mailSenderService,
	}, nil
}

func (s *senderController) Send(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r.Body.Close()
	}()

	mailInfo := dto.MailSenderDto{}
	err := json.NewDecoder(r.Body).Decode(&mailInfo)
	if err != nil {
		response(w, "Failed sending mail", err, "failed sending meil", http.StatusBadRequest)
		return
	}

	err = s.mailSenderService.Send(mailInfo.Subject, mailInfo.Cc, mailInfo.To, mailInfo.Body)
	if err != nil {
		response(w, "Failed sending mail", err, "failed sending meil", http.StatusBadRequest)
		return
	}

	response(w, "Mail succesfully send", err, "mail successfully sent", http.StatusOK)
}
func (s *senderController) SendWithAttachment(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r.Body.Close()
	}()

	mailSenderDto := dto.MailSenderDto{}
	mailInfo := r.FormValue("mailInfo")

	err := json.Unmarshal([]byte(mailInfo), &mailSenderDto)
	if err != nil {
		response(w, "Failed sending mail with attachment", err, "failed sending meil", http.StatusBadRequest)
		return
	}

	f, _, err := r.FormFile("file")
	if err != nil {
		response(w, "Failed sending mail with attachment", err, "failed sending meil", http.StatusBadRequest)
		return
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		response(w, "Failed sending mail with attachment", err, "failed sending meil", http.StatusBadRequest)
		return
	}

	err = s.mailSenderService.SendWithAttachment(mailSenderDto.Subject,
		mailSenderDto.Cc, mailSenderDto.To, mailSenderDto.Body, mailSenderDto.Filename, buf.Bytes())

	if err != nil {
		response(w, "Failed sending mail with attachment", err, "failed sending meil", http.StatusBadRequest)
		return
	}

	response(w, "Mail succesfully send", err, "mail successfully sent", http.StatusOK)
}

func response(w http.ResponseWriter, loggMessagge string, err error, responseMsg string, failurStatusCode int) {
	controllerLogger.Info(loggMessagge)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(failurStatusCode)
		controllerLogger.Error(err)
		responseError := dto.Response{
			ErrorMessage:    err.Error(),
			ResponseMessage: responseMsg,
		}
		json.NewEncoder(w).Encode(responseError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseError := dto.Response{
		ResponseMessage: responseMsg,
	}

	json.NewEncoder(w).Encode(responseError)
}
