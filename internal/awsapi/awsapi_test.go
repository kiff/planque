package awsapi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// TestSignVanilla is the get-vanilla case of AWS's SigV4 test suite.
func TestSignVanilla(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.amazonaws.com/", nil)
	c := Credentials{AccessKey: "AKIDEXAMPLE", SecretKey: "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY"}
	Sign(req, nil, c, "service", "us-east-1", time.Date(2015, 8, 30, 12, 36, 0, 0, time.UTC))
	want := "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20150830/us-east-1/service/aws4_request, " +
		"SignedHeaders=host;x-amz-date, Signature=5fa00fa31553b73ebf1942676e86291e8372ff2a2260956d9b8aae1d763fbf31"
	if got := req.Header.Get("Authorization"); got != want {
		t.Fatalf("Authorization\n got %s\nwant %s", got, want)
	}
}

// SigV4's canonical URI for every service but S3 is the path as sent,
// URI-encoded once more. The API Gateway management API's path holds "@" and
// connection ids ending in "=": signed as sent, AWS refused every post to a
// paired Chrome (22 Nov 2026, "the signature we calculated does not match").
func TestTheCanonicalPathIsEncodedAsAWSSignsIt(t *testing.T) {
	for in, want := range map[string]string{
		"":                                  "/",
		"/":                                 "/",
		"/v1/@connections/Lx1ZbcMDoAMCJRQ=": "/v1/%40connections/Lx1ZbcMDoAMCJRQ%3D",
		"/a-b_c.d~e/f":                      "/a-b_c.d~e/f",
		"/example%20space/":                 "/example%2520space/",
		"/x+y:z":                            "/x%2By%3Az",
		// The AWS CLI's own canonical path for GetConnection on that API, read
		// from its --debug output: the id sent escaped, signed escaped again.
		"/v1/@connections/Zz9probeNotACo0%3D": "/v1/%40connections/Zz9probeNotACo0%253D",
	} {
		if got := canonicalPath(in); got != want {
			t.Errorf("canonicalPath(%q) = %q, want %q", in, got, want)
		}
	}
}

// A dropped connection (06:33Z, 26 Nov: "Post https://kms.eu-west-1.amazonaws.com/: EOF" failed a grants call) is sent once
// more, signed afresh, on a client whose calls are safe to send twice (KMS's). One whose writes may have landed (DynamoDB's)
// returns the drop as it was, and an AWS answer, even an error, is never sent again.
func TestADroppedCallIsSentOnceMoreOnlyWhereThatIsSafe(t *testing.T) {
	serve := func(dropFirst bool, status int) (*httptest.Server, *atomic.Int32) {
		var calls atomic.Int32 // the handler's goroutine adds, the test reads
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if n := calls.Add(1); dropFirst && n == 1 {
				conn, _, err := w.(http.Hijacker).Hijack()
				if err == nil {
					conn.Close()
				}
				return
			}
			if r.Header.Get("Authorization") == "" {
				t.Error("a call went unsigned")
			}
			w.WriteHeader(status)
			if status != http.StatusOK {
				fmt.Fprint(w, `{"__type":"com.amazonaws.kms#KMSInternalException","message":"try later"}`)
				return
			}
			fmt.Fprint(w, `{"Plaintext":"b2s="}`)
		}))
		t.Cleanup(srv.Close)
		return srv, &calls
	}
	creds := func() (Credentials, error) { return Credentials{AccessKey: "AKID", SecretKey: "secret"}, nil }
	call := func(srv *httptest.Server, retry bool) error {
		c := &Client{Service: "kms", Region: "eu-west-1", Target: "TrentService", JSONType: "application/x-amz-json-1.1",
			Creds: creds, Endpoint: srv.URL + "/", RetryDropped: retry}
		var out struct{ Plaintext string }
		return c.Do(context.Background(), "Decrypt", map[string]string{"CiphertextBlob": "eA=="}, &out)
	}

	srv, calls := serve(true, http.StatusOK)
	if err := call(srv, true); err != nil || calls.Load() != 2 {
		t.Fatalf("with the retry: err %v after %d calls, want success after 2", err, calls.Load())
	}
	srv, calls = serve(true, http.StatusOK)
	if err := call(srv, false); err == nil || !dropped(err) || calls.Load() != 1 {
		t.Fatalf("without it: err %v after %d calls, want the drop after 1", err, calls.Load())
	}
	srv, calls = serve(false, http.StatusInternalServerError)
	if err := call(srv, true); !Is(err, "KMSInternalException") || calls.Load() != 1 {
		t.Fatalf("an AWS error: err %v after %d calls, want it as answered after 1", err, calls.Load())
	}

	// The retry waits first; a caller whose time runs out in that wait gets the drop, and nothing is sent again.
	defer func(w time.Duration) { retryWait = w }(retryWait)
	retryWait = time.Second
	srv, calls = serve(true, http.StatusOK)
	c := &Client{Service: "kms", Region: "eu-west-1", Target: "TrentService", JSONType: "application/x-amz-json-1.1",
		Creds: creds, Endpoint: srv.URL + "/", RetryDropped: true}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	if err := c.Do(ctx, "Decrypt", map[string]string{"CiphertextBlob": "eA=="}, nil); err == nil || !dropped(err) || calls.Load() != 1 || time.Since(start) > 500*time.Millisecond {
		t.Fatalf("cancelled in the wait: err %v after %d calls in %v, want the drop after 1, at once", err, calls.Load(), time.Since(start))
	}
}
