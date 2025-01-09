package main

import (
	"net/http"
	"trace/data"

	shared "github.com/tsizism/semplicita/linux/shared"
)

type JSONPayload struct {
	Src string `json:"src"`
	Via string `json:"via"`
	Data string `json:"data"`
}

func (appCtx *applicationContext)writeTrace(w http.ResponseWriter, r *http.Request) {
	appCtx.logger.Printf("Hit trace writeTrace http.Request=%+v\n", r)

	var requestPaylod JSONPayload 
	shared.ReadJSON(w, r, &requestPaylod)

	event := data.TraceEntry {
		Src: requestPaylod.Src,
		Via: requestPaylod.Via,
		Data: requestPaylod.Data,
	}

	err := appCtx.models.Insert(event)

	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	resp := shared.JsonResponse {
		Error: false,
		Message: "inserted",
	}

	shared.WriteJSON(w, http.StatusAccepted, resp)
}