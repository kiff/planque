package awsapi

import (
	"net/http"
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
