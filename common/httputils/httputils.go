package httputils

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type httpUtils struct {
	req *http.Request
	Err error
}

func NewRequest(method string, url string, body ...io.Reader) *httpUtils {
	var body2 io.Reader = nil
	if len(body) > 0 {
		body2 = body[0]
	}

	r, err := http.NewRequest(method, url, body2)
	return &httpUtils{
		req: r,
		Err: err,
	}
}

func (h *httpUtils) SetBasicAuth(u, p string) *httpUtils {
	if h.Err != nil {
		return h
	}
	h.req.SetBasicAuth(u, p)
	return h
}

func (h *httpUtils) JSON(data interface{}) *httpUtils {
	if h.Err != nil {
		return h
	}

	resp, err := http.DefaultClient.Do(h.req)
	if err != nil {
		h.Err = err
		return h
	}
	defer resp.Body.Close()
	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.Err = err
		return h
	}

	/*
		if resp.StatusCode != http.StatusOK {
			h.Err = errors.New(string(buff))
			return h
		}
	*/

	log.Println("response:", string(buff))

	err = json.Unmarshal(buff, data)
	if err != nil {
		h.Err = errors.New(string(buff))
		return h
	}

	return h
}
