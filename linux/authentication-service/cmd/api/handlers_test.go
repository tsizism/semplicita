package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type RoundTripFu func(req *http.Request) *http.Response

// type RoundTripper interface { RoundTrip(*Request) (*Response, error)
func (f RoundTripFu) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFu) *http.Client {
	return &http.Client{
		Transport: fn, // implements RoundTripper interface
	}
}

func Test_Authenticate(t *testing.T) {
	jsonToReturn := `{
		"error" : false,
		"message": "test msg",
	}`

	cfgTestApp.Client = NewTestClient(
		func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusConflict,
				Body:       io.NopCloser(bytes.NewBufferString(jsonToReturn)),
				Header:     make(http.Header),
			}
		})

	postBody := map[string]interface{}{
		"email":    "bla@bla.com",
		"password": "bla",
	}

	body, err := json.Marshal(postBody)
	if err != nil {
		t.Errorf("marshal fialed err=%v", err)
	}

	req, err := http.NewRequest("POST", "/authenticate", bytes.NewReader(body))
	if err != nil {
		t.Errorf("newRequest fialed err=%v", err)
	}
	rr := httptest.NewRecorder()

	fmt.Printf("req=%+v\n", *req)
	fmt.Printf("rr=%v\n", rr)

	appCtx := applicationContext{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		cfg:    cfgTestApp,
	}
	handler := http.HandlerFunc(appCtx.authenticateHandler)
	handler.ServeHTTP(rr, req)

	fmt.Printf("Code=%v", rr.Code)
	fmt.Printf("Body=%v", rr.Body)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected 202 and got %v", rr.Code)
	}

}
