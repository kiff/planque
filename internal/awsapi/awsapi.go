// Package awsapi calls AWS's JSON APIs (DynamoDB, KMS) with Signature
// Version 4 and the standard library only. Credentials come from the
// environment, which is how Lambda provides them.
package awsapi

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// Credentials sign requests.
type Credentials struct {
	AccessKey, SecretKey, SessionToken string
}

// FromEnv reads AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and
// AWS_SESSION_TOKEN.
func FromEnv() (Credentials, error) {
	c := Credentials{os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_SESSION_TOKEN")}
	if c.AccessKey == "" || c.SecretKey == "" {
		return c, errors.New("awsapi: AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY are not set")
	}
	return c, nil
}

// Client calls one service in one region.
type Client struct {
	Service  string // "dynamodb", "kms"
	Region   string
	Target   string // the X-Amz-Target prefix: "DynamoDB_20120810", "TrentService"
	JSONType string // "application/x-amz-json-1.0" or "-1.1"
	Creds    func() (Credentials, error)
	HTTP     *http.Client
	Endpoint string // optional override, e.g. a local DynamoDB
}

// Error is an AWS error answer.
type Error struct {
	Status  int
	Type    string
	Message string
}

func (e *Error) Error() string { return fmt.Sprintf("%s (%d): %s", e.Type, e.Status, e.Message) }

// Is reports whether err is an AWS error of type typ, such as
// "ConditionalCheckFailedException".
func Is(err error, typ string) bool {
	var e *Error
	return errors.As(err, &e) && e.Type == typ
}

// Do calls action with in and decodes the answer into out.
func (c *Client) Do(ctx context.Context, action string, in, out any) error {
	body, err := json.Marshal(in)
	if err != nil {
		return err
	}
	endpoint := c.Endpoint
	if endpoint == "" {
		endpoint = "https://" + c.Service + "." + c.Region + ".amazonaws.com/"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", c.JSONType)
	req.Header.Set("X-Amz-Target", c.Target+"."+action)
	creds, err := c.Creds()
	if err != nil {
		return err
	}
	Sign(req, body, creds, c.Service, c.Region, time.Now())
	hc := c.HTTP
	if hc == nil {
		hc = &http.Client{Timeout: 10 * time.Second}
	}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		var e struct {
			Type    string `json:"__type"`
			Message string `json:"message"`
			Msg     string `json:"Message"`
		}
		_ = json.Unmarshal(b, &e)
		t := e.Type
		if i := strings.LastIndexByte(t, '#'); i >= 0 {
			t = t[i+1:]
		}
		return &Error{Status: resp.StatusCode, Type: t, Message: e.Message + e.Msg}
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(b, out)
}

// Sign adds SigV4 headers to req, whose body is body.
func Sign(req *http.Request, body []byte, c Credentials, service, region string, now time.Time) {
	now = now.UTC()
	amzDate := now.Format("20060102T150405Z")
	day := now.Format("20060102")
	payload := sha256.Sum256(body)
	req.Header.Set("X-Amz-Date", amzDate)
	if c.SessionToken != "" {
		req.Header.Set("X-Amz-Security-Token", c.SessionToken)
	}
	headers := map[string]string{"host": req.URL.Host}
	for k, v := range req.Header {
		headers[strings.ToLower(k)] = strings.TrimSpace(strings.Join(v, ","))
	}
	names := make([]string, 0, len(headers))
	for k := range headers {
		names = append(names, k)
	}
	sort.Strings(names)
	var canonHeaders strings.Builder
	for _, k := range names {
		canonHeaders.WriteString(k + ":" + headers[k] + "\n")
	}
	signed := strings.Join(names, ";")
	path := canonicalPath(req.URL.EscapedPath())
	canon := strings.Join([]string{req.Method, path, req.URL.RawQuery, canonHeaders.String(), signed,
		hex.EncodeToString(payload[:])}, "\n")
	scope := day + "/" + region + "/" + service + "/aws4_request"
	sum := sha256.Sum256([]byte(canon))
	toSign := "AWS4-HMAC-SHA256\n" + amzDate + "\n" + scope + "\n" + hex.EncodeToString(sum[:])
	k := mac([]byte("AWS4"+c.SecretKey), day)
	k = mac(k, region)
	k = mac(k, service)
	k = mac(k, "aws4_request")
	sig := hex.EncodeToString(mac(k, toSign))
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential="+c.AccessKey+"/"+scope+
		", SignedHeaders="+signed+", Signature="+sig)
}

// canonicalPath is SigV4's canonical URI for every service but S3: the
// path as it's sent, URI-encoded once more, every byte but A-Z a-z 0-9
// - _ . ~ and / as %XX. The API Gateway management API's path holds "@"
// (/@connections/) and connection ids ending in "=", which a path sends as
// they are and AWS signs as %40 and %3D: signed as sent, every post to a
// paired Chrome was refused (403, the signature doesn't match).
func canonicalPath(escaped string) string {
	if escaped == "" {
		return "/"
	}
	var b strings.Builder
	for i := 0; i < len(escaped); i++ {
		c := escaped[i]
		if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' || c == '-' || c == '_' || c == '.' || c == '~' || c == '/' {
			b.WriteByte(c)
			continue
		}
		fmt.Fprintf(&b, "%%%02X", c)
	}
	return b.String()
}

func mac(key []byte, s string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(s))
	return h.Sum(nil)
}
