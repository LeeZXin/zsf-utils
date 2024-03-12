package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"golang.org/x/net/http2"
	"io"
	"net/http"
	"time"
)

const (
	JsonContentType = "application/json;charset=utf-8"
)

type RetryableRoundTripper struct {
	Delegated http.RoundTripper
}

func (t *RetryableRoundTripper) RoundTrip(request *http.Request) (response *http.Response, err error) {
	buf := new(bytes.Buffer)
	hasBody := request.Body != nil
	if hasBody {
		_, err = io.Copy(buf, request.Body)
		request.Body.Close()
	}
	if err != nil {
		return
	}
	for i := 0; i < 3; i++ {
		if hasBody {
			request.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
		}
		response, err = t.Delegated.RoundTrip(request)
		if err == nil {
			break
		}
	}
	return
}

// NewRetryableHttpClient http client
func NewRetryableHttpClient(http2Enabled ...bool) *http.Client {
	http2EnabledRet := false
	if len(http2Enabled) > 0 {
		http2EnabledRet = http2Enabled[0]
	}
	return &http.Client{
		Transport: &RetryableRoundTripper{
			Delegated: newRoundTripper(http2EnabledRet),
		},
		Timeout: 30 * time.Second,
	}
}

func newRoundTripper(http2Enabled bool) http.RoundTripper {
	transport := &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
		MaxIdleConns:        100,
		IdleConnTimeout:     time.Minute,
		MaxConnsPerHost:     10,
	}
	if http2Enabled {
		ret, err := http2.ConfigureTransports(transport)
		if err != nil {
			panic(err)
		}
		ret.AllowHTTP = true
		return ret
	}
	return transport
}

// NewHttpClient http client
func NewHttpClient(http2Enabled ...bool) *http.Client {
	if len(http2Enabled) > 0 && http2Enabled[0] {
		return &http.Client{
			Transport: newRoundTripper(true),
			Timeout:   30 * time.Second,
		}
	}
	return &http.Client{
		Transport: newRoundTripper(false),
		Timeout:   30 * time.Second,
	}
}

func Post(ctx context.Context, client *http.Client, url string, header map[string]string, req, resp any) error {
	var (
		reqJson []byte
		err     error
	)
	if req != nil {
		reqJson, err = json.Marshal(req)
		if err != nil {
			return err
		}
	} else {
		reqJson = []byte{}
	}
	request, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqJson))
	if err != nil {
		return err
	}
	if header != nil {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}
	request.Header.Set("Content-Type", JsonContentType)
	post, err := client.Do(request)
	if err != nil {
		return err
	}
	defer post.Body.Close()
	respBody, err := io.ReadAll(post.Body)
	if err != nil {
		return err
	}
	if resp != nil {
		return json.Unmarshal(respBody, resp)
	}
	return nil
}

func Get(ctx context.Context, client *http.Client, url string, header map[string]string, resp any) error {
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	if header != nil {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}
	get, err := client.Do(request)
	if err != nil {
		return err
	}
	defer get.Body.Close()
	respBody, err := io.ReadAll(get.Body)
	if err != nil {
		return err
	}
	if resp != nil {
		return json.Unmarshal(respBody, resp)
	}
	return nil
}
