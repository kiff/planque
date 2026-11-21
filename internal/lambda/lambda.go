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
	// QueryStringParameters is a WebSocket API's $connect query.
	QueryStringParameters map[string]string `json:"queryStringParameters"`
	RequestContext        struct {
		DomainName string `json:"domainName"`
		// RouteKey and ConnectionID are a WebSocket API's ($connect,
		// $default, $disconnect): only API Gateway sets them, never a
		// function URL request.
		RouteKey     string `json:"routeKey"`
		ConnectionID string `json:"connectionId"`
		Identity     struct {
			SourceIP string `json:"sourceIp"`
		} `json:"identity"`
		HTTP struct {
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

// WSEvent is one WebSocket API event: a connection opened ($connect, with
// its query), a message on it ($default), or it closed ($disconnect).
type WSEvent struct {
	Route      string
	Connection string
	Query      map[string]string
	Body       []byte
	SourceIP   string
}

// WebSocket handles a WebSocket API's events. Its answer is the status API
// Gateway gets: a $connect answered other than 2xx is refused.
type WebSocket interface {
	WebSocket(ctx context.Context, ev WSEvent) int
}

// Running reports whether this process is a Lambda runtime.
func Running() bool { return os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" }

// Start serves Lambda invocations with h forever, and a WebSocket API's
// events with ws (nil: none are taken).
func Start(h http.Handler, ws WebSocket) error {
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
		out, err := invoke(ctx, h, ws, raw)
		cancel()
		if err != nil {
			b, _ := json.Marshal(map[string]string{"errorMessage": err.Error(), "errorType": "PlanqueError"})
			post(hc, base+id+"/error", b)
			continue
		}
		post(hc, base+id+"/response", out)
	}
}

func invoke(ctx context.Context, h http.Handler, ws WebSocket, raw []byte) (b []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	var e Event
	if err := json.Unmarshal(raw, &e); err != nil {
		return nil, err
	}
	if e.RequestContext.ConnectionID != "" && e.RequestContext.RouteKey != "" && e.RequestContext.HTTP.Method == "" {
		// A WebSocket API's event never goes through the HTTP routes.
		status := http.StatusForbidden
		if ws != nil {
			body := []byte(e.Body)
			if e.IsBase64Encoded {
				if body, err = base64.StdEncoding.DecodeString(e.Body); err != nil {
					return nil, err
				}
			}
			status = ws.WebSocket(ctx, WSEvent{Route: e.RequestContext.RouteKey, Connection: e.RequestContext.ConnectionID,
				Query: e.QueryStringParameters, Body: body, SourceIP: e.RequestContext.Identity.SourceIP})
		}
		return json.Marshal(map[string]int{"statusCode": status})
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
