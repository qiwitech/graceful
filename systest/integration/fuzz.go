package integration

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"

	api "github.com/qiwitech/graceful/systest/api_v1"
)

var (
	methodNames []string
	APIbaseURL  string
	httpClient  = http.Client{}
)

func init() {
	// генерим список имен методов
	x := &api.APIHTTPClient{}
	t := reflect.TypeOf(x)
	methodNames = make([]string, t.NumMethod())
	for i := 0; i < t.NumMethod(); i++ {
		methodNames[i] = t.Method(i).Name
	}

	// запуск серверов
	addr, _, err := APIServiceChain(&PlutoMockProcessing{}, &PlutoMockStorage{})
	if err != nil {
		panic(err)
	}
	APIbaseURL = addr + "/"
}

func Fuzz(data []byte) int {
	result := 1 // 1 - good mutation, 0 - bad mutation
	req := bytes.NewBuffer(data)
	for _, methodName := range methodNames {
		// request
		resp, err := httpClient.Post(APIbaseURL+methodName, "application/octet-stream", req)
		if err != nil {
			result = 0
			break // can't send this
		}
		// response
		data, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			result = 0
		}
		if resp.StatusCode != http.StatusOK {
			result = 0
		}

	}
	return result
}
