package openapi

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestSelectedOperationsKeepTheFullSourceDigest(t *testing.T) {
	var doc map[string]any
	if err := json.Unmarshal([]byte(desc), &doc); err != nil {
		t.Fatal(err)
	}
	paths := doc["paths"].(map[string]any)
	for i := range MaxOperations + 10 {
		paths[fmt.Sprintf("/items/%d", i)] = map[string]any{"get": map[string]any{"operationId": fmt.Sprintf("readItem%d", i)}}
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Parse(raw, key); err == nil {
		t.Fatal("an unrestricted snapshot exceeded its operation bound")
	}
	s, err := ParseOperations(raw, key, []string{"getOrder", "create_refund"})
	if err != nil || len(s.Operations) != 2 || len(s.Skipped) != 0 {
		t.Fatalf("snapshot=%+v err=%v", s, err)
	}
	h := sha256.Sum256(raw)
	if s.Digest != hex.EncodeToString(h[:]) {
		t.Fatal("subset lost the complete document's digest")
	}
	// A caller cannot silently receive a smaller selection than requested.
	for _, names := range [][]string{{"getOrder", "missing"}, {"getOrder", "getOrder"}, {}, {"getOrder", "upload"}} {
		if _, err := ParseOperations(raw, key, names); err == nil {
			t.Fatalf("accepted missing, duplicate or unsupported selection %v", names)
		}
	}
}

func TestDescriptionBetweenTwoAndSixteenMBIsAccepted(t *testing.T) {
	raw := []byte(strings.TrimSuffix(desc, "}") + strings.Repeat(" ", 3<<20) + "}")
	s, err := ParseOperations(raw, key, []string{"getOrder"})
	if err != nil || len(s.Operations) != 1 {
		t.Fatalf("a bounded large document was refused: %v", err)
	}
}

func TestOperationDigestIgnoresWordsButKeepsArgumentSemantics(t *testing.T) {
	operation := Operation{Name: "sendReply", Method: "POST", Path: "/tickets/{ticket}/replies", Body: "fields",
		InputSchema: json.RawMessage(`{"type":"object","description":"Old words","properties":{"description":{"type":"string","description":"Text"},"body":{"type":"string"}},"required":["body"]}`)}
	first, err := operation.Digest()
	if err != nil {
		t.Fatal(err)
	}
	operation.Summary, operation.Description = "New summary", "New words"
	operation.InputSchema = json.RawMessage(strings.ReplaceAll(string(operation.InputSchema), "Old words", "Better words"))
	second, err := operation.Digest()
	if err != nil || first != second {
		t.Fatalf("wording changed execution identity: %s %s %v", first, second, err)
	}
	operation.InputSchema = json.RawMessage(strings.Replace(string(operation.InputSchema), `"description":{"type":"string"`, `"description":{"type":"integer"`, 1))
	third, err := operation.Digest()
	if err != nil || first == third {
		t.Fatal("a property named description lost its type from the execution identity")
	}
	operation.Path = "/tickets/{ticket}/notes"
	fourth, err := operation.Digest()
	if err != nil || fourth == third {
		t.Fatal("a changed request path kept its identity")
	}
}
