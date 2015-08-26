package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var print = fmt.Print

func PostUrlEncodeForm(postUrl string, val url.Values) (map[string]string, error) {

	resp, err1 := http.PostForm(postUrl, val)

	if err1 != nil {
		return nil, err1
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("connect remote server failed")
	}
	defer resp.Body.Close()
	content, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return nil, err2
	}

	fmt.Println(string(content))

	tempMap := make(map[string]string)
	err3 := json.Unmarshal(content, &tempMap)
	if err3 != nil {
		return nil, err3
	}
	return tempMap, nil
}
