package rest

import (
	"encoding/json"
	"github.com/codegangsta/inject"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
)

///rest error
type RestError struct {
	ErrorCode int32  //error code
	ErrorDesc string //error desc
}

func (re *RestError) Error() string {
	return "rest error code[" + strconv.FormatInt(int64(re.ErrorCode), 10) + "]" + ",error desc[" + re.ErrorDesc + "]"
}

type RestResult interface {
	Result() string
}

type RestReturnObj struct {
	ErrorCode int32      `json:"errorCode"`
	Result    RestResult `json:"result"`
}

func RestPostHandler() martini.Handler {
	return func(c martini.Context, rw http.ResponseWriter, req *http.Request) {

		defer func(c martini.Context, rw http.ResponseWriter, req *http.Request) {
			restErrorVal := c.Get(inject.InterfaceOf((*error)(nil)))
			restResultVal := c.Get(inject.InterfaceOf((*RestResult)(nil)))

			//no rest func
			if !restErrorVal.IsValid() && !restResultVal.IsValid() {

				return
			}

			restReturnObj := &RestReturnObj{}

			if restErrorVal.IsValid() {

				restErrorObj := restErrorVal.Interface().(*RestError)
				restReturnObj.ErrorCode = restErrorObj.ErrorCode
			} else {

				restResultObj := restResultVal.Interface().(RestResult)
				restReturnObj.ErrorCode = 0
				restReturnObj.Result = restResultObj
			}

			rw.Header().Set("Content-Type", "application/json; charset=utf-8")
			content, err := json.Marshal(restReturnObj)

			if err != nil {
				panic("encode error")
			}

			rw.WriteHeader(http.StatusOK)
			rw.Write(content)
		}(c, rw, req)

		c.Next()
	}
}
