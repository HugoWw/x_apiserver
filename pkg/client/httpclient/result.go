package httpclient

import (
	"encoding/json"
	"fmt"
)

type Result struct {
	body       []byte
	err        error
	statusCode int
}

// Raw return response raw data and error
func (r Result) Raw() ([]byte, error) {
	return r.body, r.err
}

// StatusCode  return http response status code
func (r Result) StatusCode() int {
	return r.statusCode
}

// Into Unmarshal response data into struct obj,
// an error will be returned if the Result object has an error or a deserialization error
func (r Result) Into(obj any) error {
	if r.err != nil {
		return r.err
	}

	if len(r.body) == 0 {
		return nil
	}

	if err := json.Unmarshal(r.body, obj); err != nil {
		return fmt.Errorf("json unmarshal error: %s, response body info :%s\n", err, string(r.body))
	}

	return nil
}

// Error return executing the request error info
func (r Result) Error() error {
	return r.err
}
