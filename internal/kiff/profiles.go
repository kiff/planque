package kiff

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// ProfileApproval carries the exact bytes signed by the phone, base64
// encoded for the gateway. Only KIFF decides whether that evidence applies.
type ProfileApproval struct {
	Message   string `json:"message"`
	Body      string `json:"body"`
	Signature string `json:"signature"`
	DeviceID  string `json:"device_id"`
}

type ProfileRevision struct {
	Revision      int64  `json:"revision"`
	ProfileSHA256 string `json:"profile_sha256"`
	WordsSHA256   string `json:"words_sha256"`
}

func (c *Client) ApplyProfile(ctx context.Context, connector string, raw json.RawMessage, approval ProfileApproval) (ProfileRevision, error) {
	var out ProfileRevision
	status, err := c.do(ctx, http.MethodPost, c.Gateway+"/v1/profiles", connector, struct {
		Profile  json.RawMessage `json:"profile"`
		Approval ProfileApproval `json:"approval"`
	}{raw, approval}, &out)
	if err != nil {
		return out, err
	}
	if status != http.StatusOK || out.Revision <= 0 || out.ProfileSHA256 == "" || out.WordsSHA256 == "" {
		return ProfileRevision{}, errors.New("KIFF did not confirm applying the profile")
	}
	return out, nil
}
