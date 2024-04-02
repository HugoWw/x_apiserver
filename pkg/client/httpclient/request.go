package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"github.com/HugoWw/x_apiserver/pkg/clog"
	"golang.org/x/net/http2"
	"io"
	"k8s.io/apimachinery/pkg/util/net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

type Request struct {
	c *RTClient

	timeout    time.Duration
	maxRetries int

	// generic components accessible via method setters
	scheme     string
	host       string
	verb       string
	pathPrefix string
	subpath    string
	params     url.Values
	headers    http.Header

	// output
	err  error
	body io.Reader
}

func newRequest(client *RTClient) *Request {
	r := &Request{
		c:          client,
		timeout:    client.timeout,
		maxRetries: client.maxRetries,
		host:       client.host,
		scheme:     client.scheme,
	}

	if client.maxRetries < 0 {
		r.maxRetries = 5
	}

	r.SetHeader("Content-Type", "application/json")
	return r
}

func (r *Request) Verb(verb string) *Request {
	r.verb = verb
	return r
}

func (r *Request) SetHeader(key string, values ...string) *Request {
	if r.headers == nil {
		r.headers = http.Header{}
	}
	r.headers.Del(key)
	for _, value := range values {
		r.headers.Add(key, value)
	}

	return r
}

// SetParams creates a query parameter with the given string value.
func (r *Request) SetParams(paramName, value string) *Request {
	if r.err != nil {
		return r
	}

	if r.params == nil {
		r.params = make(url.Values)
	}

	r.params[paramName] = append(r.params[paramName], value)
	return r
}

func (r *Request) Timeout(d time.Duration) *Request {
	if r.err != nil {
		return r
	}
	r.timeout = d
	return r
}

func (r *Request) MaxRetries(maxRetries int) *Request {
	if maxRetries < 0 {
		maxRetries = 0
	}

	r.maxRetries = maxRetries
	return r
}

// Body makes the request use obj as the body. Optional.
// If obj is a string, try to read a file of that name.
// If obj is a []byte, send it directly.
// If obj is an io.Reader, use it directly.
// If obj is a nil, return Request directly.
// Otherwise, set an error.
func (r *Request) Body(obj interface{}) *Request {
	if r.err != nil {
		return r
	}

	switch t := obj.(type) {
	case string:
		data, err := os.ReadFile(t)
		if err != nil {
			r.err = err
			return r
		}
		r.body = bytes.NewReader(data)
	case []byte:
		r.body = bytes.NewReader(t)
	case nil:
		return r
	case io.Reader:
		r.body = t
	default:
		r.err = fmt.Errorf("unknown type used for request body: %+v", obj)
	}

	return r
}

// Prefix set api url prefix
func (r *Request) Prefix(prefixPath ...string) *Request {
	if r.err != nil {
		return r
	}

	r.pathPrefix = path.Join(r.pathPrefix, path.Join(prefixPath...))
	return r
}

// SetPath set api url sub path
func (r *Request) SetPath(suffixPath ...string) *Request {
	if r.err != nil {
		return r
	}
	r.subpath = path.Join(r.subpath, path.Join(suffixPath...))
	return r
}

// RequestURL overwrites existing path and parameters with the value of the provided server relative
// for example:
// uri= http://www.example.com:8080/api
// overwrites info:
// r.path=/api;r.shceme=http;r.host=www.example.com:8080
func (r *Request) RequestURL(uri string) *Request {
	if r.err != nil {
		return r
	}

	urlInfo, err := url.Parse(uri)
	if err != nil {
		r.err = err
		return r
	}

	r.scheme = urlInfo.Scheme
	r.host = urlInfo.Host
	r.subpath = urlInfo.Path

	if len(urlInfo.Query()) > 0 {
		if r.params == nil {
			r.params = make(url.Values)
		}
		for k, v := range urlInfo.Query() {
			r.params[k] = v
		}
	}

	return r
}

// URL return current request url info
func (r *Request) URL() *url.URL {
	finalURL := &url.URL{}
	finalURL.Path = path.Clean(path.Join(r.pathPrefix, r.subpath))
	finalURL.Scheme = r.scheme
	finalURL.Host = r.host

	query := url.Values{}
	if r.params != nil {
		for key, values := range r.params {
			for _, value := range values {
				query.Add(key, value)
			}
		}
	}

	finalURL.RawQuery = query.Encode()

	return finalURL
}

func (r *Request) request(ctx context.Context, fn func(r *http.Request, response *http.Response)) error {

	client := r.c.Client
	if client == nil {
		client = http.DefaultClient
	}

	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout*time.Second)
		defer cancel()
	}

	retries := 0
	for {
		url := r.URL().String()
		req, err := http.NewRequest(r.verb, url, r.body)
		if err != nil {
			return err
		}

		req = req.WithContext(ctx)
		req.Header = r.headers

		resp, err := client.Do(req)
		if err != nil {
			if r.verb != "GET" {
				return err
			}

			// for connection error and upstream server shutdown error retry
			if net.IsConnectionReset(err) || net.IsProbableEOF(err) {
				resp = &http.Response{
					StatusCode: http.StatusInternalServerError,
					Header:     http.Header{"Retry-After": []string{"1"}},
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}
			} else {
				return err
			}

		}

		done := func() bool {
			defer func() {
				if resp.ContentLength <= maxBodySlurpSize {
					io.Copy(io.Discard, &io.LimitedReader{R: resp.Body, N: maxBodySlurpSize})
				}
				resp.Body.Close()
			}()

			retries++

			if seconds, wait := checkWait(resp); wait && retries < r.maxRetries {
				if seeker, ok := r.body.(io.Seeker); ok && r.body != nil {
					_, err := seeker.Seek(0, 0)
					if err != nil {
						clog.Logger.Infof("Could not retry request, can't Seek() back to beginning of body for %T", r.body)
						fn(req, resp)
						return true
					}
				}

				clog.Logger.Infof("Got a Retry-After %ds response for attempt %d to %v", seconds, retries, url)
				time.Sleep(time.Duration(seconds) * time.Second)
				return false
			}

			fn(req, resp)
			return true
		}()

		if done {
			return nil
		}

	}
}

// checkWait returns true along with a number of seconds if the server instructed us to wait
// before retrying.
func checkWait(resp *http.Response) (int, bool) {
	switch r := resp.StatusCode; {
	// too many connection or shutdown error retry
	case r == http.StatusTooManyRequests, r >= 500:
	default:
		return 0, false

	}
	i, ok := retryAfterSeconds(resp)

	return i, ok
}

// retryAfterSeconds returns the value of the 'Retry-After' header and true, or 0 and false if
// the header was missing or not a valid number.
func retryAfterSeconds(resp *http.Response) (int, bool) {
	if h := resp.Header.Get("Retry-After"); len(h) > 0 {
		if i, err := strconv.Atoi(h); err == nil {
			return i, true
		}
	}

	return 0, false
}

// Do formats and executes the request. Returns a Result object
func (r *Request) Do(ctx context.Context) Result {
	var result Result
	err := r.request(ctx, func(req *http.Request, resp *http.Response) {
		result = r.transformResponse(resp, req)
	})

	if err != nil {
		return Result{err: err}
	}
	return result
}

// transformResponse only successful responses are converted to Result
func (r *Request) transformResponse(resp *http.Response, req *http.Request) Result {
	var body []byte
	if resp.Body != nil {
		data, err := io.ReadAll(resp.Body)
		switch err.(type) {
		case nil:
			body = data
		case http2.StreamError:
			clog.Logger.Errorf("Stream error %#v when reading response body, may be caused by closed connection.", err)
			streamError := fmt.Errorf("stream error when reading response body, may be caused by closed connection. Please retry. Original error: %v", err)
			return Result{err: streamError}
		default:
			clog.Logger.Errorf("Unexpected error when reading response body: %v", err)
			unexpectedErr := fmt.Errorf("unexpected error when reading response body. Please retry. Original error: %v", err)
			return Result{
				err: unexpectedErr,
			}
		}
	}

	switch {
	case resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusPartialContent:
		// the response was successful but the response was abnormal, example response 404|503..
		code_err := httpStatusCode2String(resp.StatusCode, r.verb, r.URL().String(), &body)

		return Result{
			body:       body,
			err:        code_err,
			statusCode: resp.StatusCode,
		}
	}

	//todo log body info

	return Result{
		body:       body,
		err:        nil,
		statusCode: resp.StatusCode,
	}
}

// Stream formats and executes the request, and offers streaming of the response.
// Returns io.ReadCloser which could be used for streaming of the response, or an error
func (r *Request) Stream(ctx context.Context) (io.ReadCloser, error) {
	if r.err != nil {
		return nil, r.err
	}

	client := r.c.Client
	if client == nil {
		client = http.DefaultClient
	}

	url := r.URL().String()
	req, err := http.NewRequest(r.verb, url, nil)
	if err != nil {
		return nil, err
	}

	if r.body != nil {
		req.Body = io.NopCloser(r.body)
	}

	req = req.WithContext(ctx)
	req.Header = r.headers

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	switch {
	case (resp.StatusCode >= 200) && (resp.StatusCode < 300):
		return resp.Body, nil

	default:
		// ensure we close the body before returning the error
		defer resp.Body.Close()

		result := r.transformResponse(resp, req)
		err := result.err
		if err == nil {
			err = fmt.Errorf("http stream response code %d while accessing %v,body info: %s", result.statusCode, url, string(result.body))
		}
		return nil, err
	}

}
