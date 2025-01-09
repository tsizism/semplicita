package main

import (
	"net/http"

	shared "github.com/tsizism/semplicita/linux/shared"
	// "golang.org/x/text/message"
)

func (appCtx applicationContext) sendMail(w http.ResponseWriter, r *http.Request) {
	appCtx.logger.Printf("Hit mail sendMail http.Request=%+v\n", r)

	type mailMsg struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Msg     string `json:"message"`
	}

	var requestPayload mailMsg

	err := shared.ReadJSON(w, r, &requestPayload)

	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}
	appCtx.logger.Printf("requestPayload=%+v", requestPayload)

	msg := msgProp{
		from:    requestPayload.From,
		to:      requestPayload.To,
		subject: requestPayload.Subject,
		data:    requestPayload.Msg,
	}

	appCtx.logger.Printf("Sending email %+v", msg)

	err = appCtx.sendSMTPMessage(msg)

	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	payload := shared.JsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	shared.WriteJSON(w, http.StatusAccepted, payload)
}
