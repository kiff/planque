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
		"one of them":                               {read, []string{"openid", "https://www.googleapis.com/auth/drive.metadata.readonly"}, true},
		"the full form":                             {read, []string{"https://www.googleapis.com/auth/drive"}, true},
		"none of them":                              {read, []string{"https://www.googleapis.com/auth/spreadsheets"}, false},
		"a narrower one":                            {Tool{Scopes: []string{"https://www.googleapis.com/auth/drive"}}, []string{"https://www.googleapis.com/auth/drive.readonly"}, false},
		"granted unknown":                           {read, nil, true},
		"written in another form":                   {read, []string{"drive.readonly"}, true},
		"the tool lists none":                       {Tool{}, []string{"https://www.googleapis.com/auth/spreadsheets"}, true},
		"a write implies its read (shopify)":        {Tool{Toolkit: "shopify", Scopes: []string{"read_orders"}}, []string{"read_locations", "write_orders"}, true},
		"a read isn't its write":                    {Tool{Toolkit: "shopify", Scopes: []string{"write_orders"}}, []string{"read_orders"}, false},
		"another name":                              {Tool{Toolkit: "shopify", Scopes: []string{"read_orders"}}, []string{"read_locations", "write_products"}, false},
		"another vendor's write isn't a read":       {Tool{Toolkit: "intercom", Scopes: []string{"read_orders"}}, []string{"read_locations", "write_orders"}, false},
		"a full scope implies its read (pipedrive)": {Tool{Toolkit: "Pipedrive", Scopes: []string{"deals:read"}}, []string{"base", "deals:full"}, true},
		"a read scope isn't its full (pipedrive)":   {Tool{Toolkit: "pipedrive", Scopes: []string{"deals:full"}}, []string{"base", "deals:read"}, false},
		"another name's full (pipedrive)":           {Tool{Toolkit: "pipedrive", Scopes: []string{"deals:read"}}, []string{"base", "persons:full"}, false},
		"another vendor's full isn't a read":        {Tool{Toolkit: "intercom", Scopes: []string{"deals:read"}}, []string{"base", "deals:full"}, false},
	} {
		if got := c.tool.Covered(c.granted); got != c.want {
			t.Errorf("%s: %v, want %v", name, got, c.want)
		}
	}
}
