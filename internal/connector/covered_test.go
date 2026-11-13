package connector

import "testing"

// A tool runs on a sign-in granted any one of the scopes it lists, or the
// full form of a read-only one; unknown either way, it's left to the vendor.
func TestAToolIsCoveredByAnyOneOfItsScopes(t *testing.T) {
	read := Tool{Scopes: []string{"https://www.googleapis.com/auth/drive.readonly", "https://www.googleapis.com/auth/drive.metadata.readonly"}}
	for name, c := range map[string]struct {
		tool    Tool
		granted []string
		want    bool
	}{
		"one of them":             {read, []string{"openid", "https://www.googleapis.com/auth/drive.metadata.readonly"}, true},
		"the full form":           {read, []string{"https://www.googleapis.com/auth/drive"}, true},
		"none of them":            {read, []string{"https://www.googleapis.com/auth/spreadsheets"}, false},
		"a narrower one":          {Tool{Scopes: []string{"https://www.googleapis.com/auth/drive"}}, []string{"https://www.googleapis.com/auth/drive.readonly"}, false},
		"granted unknown":         {read, nil, true},
		"written in another form": {read, []string{"drive.readonly"}, true},
		"the tool lists none":     {Tool{}, []string{"https://www.googleapis.com/auth/spreadsheets"}, true},
	} {
		if got := c.tool.Covered(c.granted); got != c.want {
			t.Errorf("%s: %v, want %v", name, got, c.want)
		}
	}
}
