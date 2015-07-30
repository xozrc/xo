package rest

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"net/http"
	"reflect"
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

type RestResult interface{}

type RestReturnObj struct {
	ErrorCode int32      `json:"errorCode"`
	Result    RestResult `json:"result"`
}

func RestPostHandler() martini.Handler {
	return func(c martini.Context, rw http.ResponseWriter, req *http.Request) {

		defer func(c martini.Context, rw http.ResponseWriter, req *http.Request) {
			restErrorVal := c.Get(reflect.TypeOf((*RestError)(nil)))
			restResultVal := c.Get(reflect.TypeOf((*RestResult)(nil))) //c.Get(reflect.TypeOf((*RestResult)(nil)))
			fmt.Println(restErrorVal)
			fmt.Println(restResultVal)
			//no rest func
			if !restErrorVal.IsValid() && !restResultVal.IsValid() {

				return
			}

			restReturnObj := &RestReturnObj{}

			if restErrorVal.IsValid() {

				restErrorObj := restErrorVal.Interface().(*RestError)
				restReturnObj.ErrorCode = restErrorObj.ErrorCode
			} else {
				fmt.Println("asd")
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
