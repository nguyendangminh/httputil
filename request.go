package httputil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/pkg/errors"
)

const REQUEST_TIMEOUT int = 3 // seconds

func Get(endpoint string) ([]byte, error) {
	return request("GET", endpoint, nil)
}

func Post(endpoint string, data []byte) ([]byte, error) {
	return request("POST", endpoint, data)
}

func Put(endpoint string, data []byte) ([]byte, error) {
	return request("PUT", endpoint, data)
}

func Delete(endpoint string) ([]byte, error) {
	return request("DELETE", endpoint, nil)
}

func request(method string, endpoint string, data []byte) ([]byte, error) {
	req, err := http.NewRequest(method, endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrapf(err, "NewRequest failed: %s %s\n", method, endpoint)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Duration(REQUEST_TIMEOUT) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "HTTP Do failed")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "ReadAll failed: body=%+v\n", resp.Body)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
    	return nil, fmt.Errorf("%d\n%s", resp.StatusCode, string(body))
    }

	return body, nil
}