package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"github.com/Todorov99/mailsender/pkg/dto"
	"github.com/Todorov99/mailsender/pkg/service"
	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/sirupsen/logrus"
)

var controllerLogger = logger.NewLogrus("controller", os.Stdout)

var wg sync.WaitGroup

type senderController struct {
	mailSenderService service.MailSenderService
	mx                sync.RWMutex
	logger            *logrus.Entry
}

func NewSenderController() (senderController, error) {
	mailSenderService, err := service.NewMailSender()
	if err != nil {
		return senderController{}, err
	}

	return senderController{
		mailSenderService: mailSenderService,
		logger:            logger.NewLogrus("senderController", os.Stdout),
	}, nil
}

func (s *senderController) Send(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r.Body.Close()
	}()

	mailSenderDto := dto.MailSenderDto{}
	err := json.NewDecoder(r.Body).Decode(&mailSenderDto)
	if err != nil {
		s.logger.Error(err)
		response(w, err, nil, http.StatusBadRequest)
		return
	}

	err = s.mailSenderService.Send(mailSenderDto.Subject, mailSenderDto.Cc, mailSenderDto.To, mailSenderDto.Body)
	if err != nil {
		err = fmt.Errorf("failed sending email with subject: %q to: %q cc: %q: %w", mailSenderDto.Subject, mailSenderDto.To, mailSenderDto.Cc, err)
		s.logger.Error(err)
		response(w, err, nil, http.StatusBadRequest)
		return
	}

	response(w, err, mailSenderDto, http.StatusOK)
}

func (s *senderController) SendWithAttachments(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r.Body.Close()
	}()

	mailSenderDto := dto.MailSenderDto{}
	mailInfo := r.FormValue("mailInfo")

	err := json.Unmarshal([]byte(mailInfo), &mailSenderDto)
	if err != nil {
		s.logger.Error(err)
		response(w, err, nil, http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]

	filesToAttach := make(map[string][]byte)

	for _, f := range files {
		wg.Add(1)
		go func(mpf *multipart.FileHeader) {
			multiPartFile, err := mpf.Open()
			if err != nil {
				err = fmt.Errorf(fmt.Sprintf("failed reading file: %q", mpf.Filename))
				s.logger.Error(err)
				response(w, err, nil, http.StatusBadRequest)
				return
			}

			defer multiPartFile.Close()

			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, multiPartFile); err != nil {
				response(w, err, mailSenderDto, http.StatusBadRequest)
				return
			}

			s.mx.RLock()
			filesToAttach[mpf.Filename] = buf.Bytes()
			s.mx.RUnlock()
			wg.Done()
		}(f)

	}
	wg.Wait()

	err = s.mailSenderService.SendWithAttachments(mailSenderDto.Subject,
		mailSenderDto.Cc, mailSenderDto.To, mailSenderDto.Body, filesToAttach)

	if err != nil {
		err = fmt.Errorf("failed sending email with subject: %q to: %q cc: %q: %w", mailSenderDto.Subject, mailSenderDto.To, mailSenderDto.Cc, err)
		s.logger.Error(err)
		response(w, err, nil, http.StatusBadRequest)
		return
	}

	response(w, nil, mailSenderDto, http.StatusOK)
}

func response(w http.ResponseWriter, err error, model interface{}, statusCode int) {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		controllerLogger.Error(err)
		responseError := dto.Response{
			ErrorMessage: err.Error(),
		}
		json.NewEncoder(w).Encode(responseError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	responseError := dto.Response{
		MailInfo: model,
	}

	json.NewEncoder(w).Encode(responseError)
}
