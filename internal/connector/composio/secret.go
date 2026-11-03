package composio

import (
	"context"
	"github.com/kiff/planque/internal/connector"
	"github.com/kiff/planque/internal/secretstore"
)

type JSONAPI = secretstore.JSONAPI

// SecretFromStore resolves one configured project credential, with a short cache.
func SecretFromStore(api JSONAPI, reference string) Secret {
	read := secretstore.Reader(api, reference, "api_key")
	return func(ctx context.Context) (string, error) {
		value, err := read(ctx)
		if err != nil {
			return "", connector.ErrUnavailable
		}
		return value, nil
	}
}
