// Package lambda runs an http.Handler on AWS Lambda behind a function
// URL: the custom runtime API loop and the URL's event format (payload
// 2.0), with the standard library only.
package lambda

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Event is a function URL request.
type Event struct {
	RawPath         string            `json:"rawPath"`
	RawQueryString  string            `json:"rawQueryString"`
	Cookies         []string          `json:"cookies"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	RequestContext  struct {
		DomainName string `json:"domainName"`
		HTTP       struct {
			Method   string `json:"method"`
			SourceIP string `json:"sourceIp"`
		} `json:"http"`
	} `json:"requestContext"`
}

// Response is a function URL response.
type Response struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Cookies         []string          `json:"cookies,omitempty"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

// Serve converts one event, runs h, and converts the answer.
func Serve(ctx context.Context, h http.Handler, e Event) (Response, error) {
	body := []byte(e.Body)
	if e.IsBase64Encoded {
		b, err := base64.StdEncoding.DecodeString(e.Body)
		if err != nil {
			return Response{}, err
		}
		body = b
	}
	u := "https://" + e.RequestContext.DomainName + e.RawPath
	if e.RawQueryString != "" {
		u += "?" + e.RawQueryString
	}
	req, err := http.NewRequestWithContext(ctx, e.RequestContext.HTTP.Method, u, bytes.NewReader(body))
	if err != nil {
		return Response{}, err
	}
	for k, v := range e.Headers {
		req.Header.Set(k, v)
	}
	if len(e.Cookies) > 0 {
		req.Header.Set("Cookie", strings.Join(e.Cookies, "; "))
	}
	req.RemoteAddr = e.RequestContext.HTTP.SourceIP
	req.RequestURI = e.RawPath
	if e.RawQueryString != "" {
		req.RequestURI += "?" + e.RawQueryString
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	res := Response{StatusCode: rec.Code, Headers: map[string]string{}}
	for k, v := range rec.Header() {
		if strings.EqualFold(k, "Set-Cookie") {
			res.Cookies = append(res.Cookies, v...)
			continue
		}
		res.Headers[k] = strings.Join(v, ", ")
	}
	out := rec.Body.Bytes()
	if utf8.Valid(out) {
		res.Body = string(out)
	} else {
		res.Body, res.IsBase64Encoded = base64.StdEncoding.EncodeToString(out), true
	}
	return res, nil
}

// Running reports whether this process is a Lambda runtime.
func Running() bool { return os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" }

// Start serves Lambda invocations with h forever.
func Start(h http.Handler) error {
	base := "http://" + os.Getenv("AWS_LAMBDA_RUNTIME_API") + "/2018-06-01/runtime/invocation/"
	hc := &http.Client{} // next blocks until an invocation arrives
	for {
		resp, err := hc.Get(base + "next")
		if err != nil {
			return err
		}
		id := resp.Header.Get("Lambda-Runtime-Aws-Request-Id")
		deadline, _ := strconv.ParseInt(resp.Header.Get("Lambda-Runtime-Deadline-Ms"), 10, 64)
		raw, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}
		ctx, cancel := context.WithDeadline(context.Background(), time.UnixMilli(deadline))
		out, err := invoke(ctx, h, raw)
		cancel()
		if err != nil {
			b, _ := json.Marshal(map[string]string{"errorMessage": err.Error(), "errorType": "PlanqueError"})
			post(hc, base+id+"/error", b)
			continue
		}
		post(hc, base+id+"/response", out)
	}
}

func invoke(ctx context.Context, h http.Handler, raw []byte) (b []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	var e Event
	if err := json.Unmarshal(raw, &e); err != nil {
		return nil, err
	}
	res, err := Serve(ctx, h, e)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

func post(hc *http.Client, u string, b []byte) {
	resp, err := hc.Post(u, "application/json", bytes.NewReader(b))
	if err == nil {
		resp.Body.Close()
	}
}
