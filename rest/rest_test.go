package rest_test

import (
	"github.com/go-martini/martini"
	"github.com/xo/rest"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRestError(t *testing.T) {
	err := &rest.RestError{-100, "test rest error"}
	if err.ErrorCode != -100 {
		t.Errorf("error code %d != -100", err.ErrorCode)
	}
}

func TestNormalHttp(t *testing.T) {
	m := martini.New()
	recorder := httptest.NewRecorder()
	m.Use(rest.RestPostHandler())
	m.ServeHTTP(recorder, (*http.Request)(nil))

	if recorder.Code != http.StatusOK {
		t.Error("failed status")
		return
	}

	if recorder.Body.Len() != 0 {
		t.Error("too much data")
	}
}

func TestRestErrorResult(t *testing.T) {
	m := martini.New()
	recorder := httptest.NewRecorder()
	m.Use(rest.RestPostHandler())
	m.Use(func(c martini.Context, res http.ResponseWriter, req *http.Request) {

		err := &rest.RestError{-1000, "12"}
		c.Map(err)
	})

	m.ServeHTTP(recorder, (*http.Request)(nil))

	if recorder.Code != http.StatusOK {
		t.Error("failed status")
		return
	}

	if recorder.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Error("failed content type")
		return
	}

	var returnObj rest.RestReturnObj

	if err := json.Unmarshal(recorder.Body.Bytes(), &returnObj); err != nil {
		t.Error("json decode failed")
		return
	}

	if returnObj.ErrorCode != -1000 {
		t.Error("error code failed")
		return
	}

}

type RestLoginResult struct {
	UserId int64 `json:"userId"`
}

func TestRestResult(t *testing.T) {

	m := martini.New()
	recorder := httptest.NewRecorder()
	m.Use(rest.RestPostHandler())
	m.Use(func(c martini.Context, res http.ResponseWriter, req *http.Request) {

		tempLogin := &RestLoginResult{}
		tempLogin.UserId = 10001
		c.MapTo(tempLogin, (*rest.RestResult)(nil))
	})

	m.ServeHTTP(recorder, (*http.Request)(nil))

	if recorder.Code != http.StatusOK {
		t.Error("failed status")
		return
	}

	var returnObj rest.RestReturnObj

	if err := json.Unmarshal(recorder.Body.Bytes(), &returnObj); err != nil {
		t.Error("json decode failed:" + err.Error())
		return
	}

	if returnObj.ErrorCode != 0 {
		t.Error("error code failed")
		return
	}
	result := returnObj.Result
	loginResult := result.(RestLoginResult)
	if loginResult.UserId != 10001 {
		t.Error("result failed")
		return
	}
}
