package openapi

import (
	"strings"
	"testing"
)

// The negatives for the connect batch's key places (qa-kiff): an extra fixed header never stands for a credential,
// the request's own headers or a method; a Basic username or password never breaks the header it goes in.
func TestAFixedHeaderIsNeverACredentialTheRequestsOwnOrAMethod(t *testing.T) {
	key := Auth{In: "header", Name: "X-API-Key"}
	basic := Auth{In: "basic", User: "me@example.com/token"}
	for _, c := range []struct {
		auth Auth
		name string
	}{
		{key, "Authorization"}, {key, "authorization"}, {basic, "AUTHORIZATION"},
		{key, "Host"}, {key, "Cookie"}, {key, "Proxy-Authorization"}, {key, "Proxy-Anything"},
		{key, "X-Forwarded-Host"}, {key, "Forwarded"}, {key, "Content-Type"}, {key, "Transfer-Encoding"},
		// The key's own header, in any case.
		{key, "x-api-key"},
		// A method override turns a read into a write at vendors that honour it, past what the Card calls a read.
		{key, "X-HTTP-Method-Override"}, {key, "X-HTTP-Method"}, {key, "X-Method-Override"},
		// Not a header name at all.
		{key, "Harvest Account"}, {key, "X-A\r\nHost"}, {key, ""},
	} {
		a := c.auth
		a.Fixed = &Fixed{Name: c.name, Value: "123"}
		if err := a.Check(); err == nil {
			t.Errorf("%q accepted as a fixed header (auth %s)", c.name, c.auth.In)
		}
	}
	for _, v := range []string{"", "1\r\nHost: evil.example", "1\nX", strings.Repeat("a", MaxBasicPart+1)} {
		a := key
		a.Fixed = &Fixed{Name: "Harvest-Account-Id", Value: v}
		if err := a.Check(); err == nil {
			t.Errorf("fixed value %.30q accepted", v)
		}
	}
	ok := key
	ok.Fixed = &Fixed{Name: "Harvest-Account-Id", Value: "123456"}
	if err := ok.Check(); err != nil {
		t.Errorf("an account id header refused: %v", err)
	}
}

func TestABasicUsernameOrPasswordNeverBreaksItsHeader(t *testing.T) {
	for _, a := range []Auth{
		{In: "basic", User: "me:you"},
		{In: "basic", User: "me\r\nHost: evil.example"},
		{In: "basic", Password: "x\ny"},
		{In: "basic", User: "me", Password: "also"},
		{In: "basic", Name: "Authorization"},
		{In: "basic", Prefix: "Basic "},
		{In: "header", Name: "X-API-Key", User: "me"},
		{In: "cookie", Name: "k"},
	} {
		if err := a.Check(); err == nil {
			t.Errorf("accepted: %+v", a)
		}
	}
}

// A description that names no scheme takes the key as asked; one operation that names a scheme the key isn't is
// still never offered.
func TestNoNamedSchemeStillRefusesAnOperationSecuredOtherwise(t *testing.T) {
	noScheme := strings.Replace(desc, `"securitySchemes": {"key": {"type": "apiKey", "in": "header", "name": "X-API-Key"}, "bearer": {"type": "http", "scheme": "bearer"}},`, ``, 1)
	noScheme = strings.Replace(noScheme, `"security": [{"key": []}],`, ``, 1)
	s, err := Parse([]byte(noScheme), Auth{In: "header", Name: "Authorization", Prefix: "token "})
	if err != nil {
		t.Fatalf("no named scheme, the key as asked: %v", err)
	}
	if len(s.Operations) == 0 {
		t.Fatal("no operation read")
	}
	secured := strings.Replace(noScheme, `"operationId": "getOrder",`, `"operationId": "getOrder", "security": [{"oauth": []}],`, 1)
	s, err = Parse([]byte(secured), Auth{In: "header", Name: "Authorization", Prefix: "token "})
	if err == nil {
		for _, o := range s.Operations {
			if o.Name == "getOrder" {
				t.Error("an operation secured by a scheme the description doesn't name is offered")
			}
		}
	}
}

// A form body is the arguments, each its own field however it's written: an "&" or "=" in a value never makes
// another field.
func TestAFormFieldNeverBecomesTwo(t *testing.T) {
	op := Operation{Name: "refund", Method: "POST", Path: "/refunds", Body: "form", InputSchema: []byte(`{"type":"object","properties":{"reason":{"type":"string"},"amount":{"type":"integer"}}}`)}
	_, _, body, err := op.Request("https://api.vendor.example/v1", map[string]any{"reason": "late&amount=999999&admin=true", "amount": float64(300)})
	if err != nil {
		t.Fatal(err)
	}
	if b := string(body); strings.Contains(b, "&amount=999999") || strings.Contains(b, "&admin=true") || strings.Count(b, "amount=") != 1 {
		t.Errorf("a value made another field: %s", b)
	}
}
