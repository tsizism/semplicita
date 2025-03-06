package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	shared "github.com/tsizism/semplicita/linux/shared"
)

func setupHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// http.ResponseWriter is interface, and existing types implementing this interface are pointers.
// Meaning there's no need to use a pointer to this interface, as it's already "backed" by a pointer

// http.Request is not interface, it's just struct, and since we want to change this struct
// and have web server see those changes, it has to be a pointer. If it was just a struct value,
// we would just modify a copy of it that the web server calling our code could not see

func (appCtx applicationContext) authenticateHandler(w http.ResponseWriter, r *http.Request) {
	appCtx.logger.Printf("Hit auth http.Request=%+v\n", r)
	setupHeader(w)

	/*{
		"email": "admin@example.com",
		"password": "verysecret"
	}*/
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := shared.ReadJSON(w, r, &requestPayload)

	if err != nil {
		appCtx.logger.Printf("ReadJSON error=%+v\n", err)
		shared.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	appCtx.logger.Printf("authenticateHandler requestPayload=%+v\n", requestPayload)

	user, err := appCtx.cfg.repo.GetByEmail(requestPayload.Email)

	if err != nil {
		appCtx.logger.Printf("Bad user err=%+v\n", err)
		shared.ErrorJSON(w, errors.New("invalid user credentials"), http.StatusBadRequest)
		return
	}

	// valid, err := user.PasswordMatches(requestPayload.Password)
	valid, err := appCtx.cfg.repo.PasswordMatches(requestPayload.Password, *user)

	if err != nil || !valid {
		appCtx.logger.Printf("Bad password err=%+v, valid=%+v\n", err, valid)
		shared.ErrorJSON(w, errors.New("invalid user credentials"), http.StatusBadRequest)
		return
	}

	err = appCtx.traceEvent("auth", "REST", fmt.Sprintf("%s logged in", user.Email))

	if err != nil {
		shared.ErrorJSON(w, err)
	}

	payload := shared.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	shared.WriteJSON(w, http.StatusAccepted, payload)
}



func (appCtx applicationContext) traceEvent(src, via, data string) error {
	var jsonPayload struct {
		Src string `json:"src"`
		Via  string `json:"via"` 
		Data string `json:"data"`
	}
		
	jsonPayload.Src = src
	jsonPayload.Via = via
	jsonPayload.Data = data


	jsonData , _ := json.MarshalIndent(jsonPayload, "", "\t")
	tarceServiceURL := "http://trace-service/trace"

	request, err := http.NewRequest("POST", tarceServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	// client := http.Client{}
	_, err = appCtx.cfg.Client.Do(request)
	
	if err != nil {
		return err
	}

	return nil
}
