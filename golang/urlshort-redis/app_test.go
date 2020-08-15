package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
)

const (
	expTime       = 60
	longURL       = "https://www.example.com"
	shortLink     = "IFHzaO"
	shortLinkInfo = `{"url":"https://www.example.com", "created_at":"2017-06-09 16:59"`
)

type storageMock struct {
	mock.Mock
}

func (s *storageMock) Shorten(url string, exp int64) (string, error) {
	args := s.Called(url, exp)
	return args.String(0), args.Error(1)
}

func (s *storageMock) ShortlinkInfo(eid string) (interface{}, error) {
	args := s.Called(eid)
	return args.String(0), args.Error(1)
}

func (s *storageMock) Unshorten(eid string) (string, error) {
	args := s.Called(eid)
	return args.String(0), args.Error(1)
}

func init() {
	app = App{}
	mockR = new(storageMock)
	app.Initialize(&Env{S: mockR})
}

var app App
var mockR *storageMock

func TestCreateShortlink(t *testing.T) {
	var jsonStr = []byte(`{
"url":"https://www.example.com",
"expiration_in_minutes":60
}`)

	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal("Should be able to create a request", err)
	}

	req.Header.Set("Content-Type", "application/json")

	mockR.On("Shorten", longURL, int64(expTime)).Return(shortLink, nil).Once()
	rw := httptest.NewRecorder()
	app.Router.ServeHTTP(rw, req)

	if rw.Code != http.StatusCreated {
		t.Fatalf("Excepted receive %d. Got %d", http.StatusCreated, rw.Code)
	}

	resp := struct {
		Shortlink string `json:"shortlink"`
	}{}

	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("Should decode the response")
	}

	if resp.Shortlink != shortLink {
		t.Fatalf("Exprected receive %s. Got %s", shortLink, resp.Shortlink)
	}
}

func TestRedirect(t *testing.T) {
	r := fmt.Sprintf("/%s", shortLink)
	//req,err:= http.NewRequest("GET", r, nil)
	//if err != nil {
	//	t.Fatal("Should be able to create a request")
	//}
	req := httptest.NewRequest("GET", r, nil)

	mockR.On("Unshorten", shortLink).Return(longURL, nil).Once()

	rw := httptest.NewRecorder()
	app.Router.ServeHTTP(rw, req)

	if rw.Code != http.StatusTemporaryRedirect {
		t.Fatalf("Expected receive %d. Got %d", http.StatusTemporaryRedirect, rw.Code)
	}
}
