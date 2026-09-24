package strings_utils

import (
	"net/url"
	"reflect"
	"testing"
)

func TestEncodeURL(t *testing.T) {
	cases := map[string]string{
		"https://cdn.example/live/a b.m3u8?token=a+b/c=&name=é x": "https://cdn.example/live/a%20b.m3u8?name=%C3%A9+x&token=a+b%2Fc%3D",
		"https://user:pw@host:8443/p?q=1#frag":                    "https://user:pw@host:8443/p?q=1#frag",
		"https://host/path":                                       "https://host/path",
		"://bad":                                                  "",
	}
	for in, want := range cases {
		if got := EncodeURL(in); got != want {
			t.Errorf("EncodeURL(%q)\n got  %q\n want %q", in, got, want)
		}
	}
}

func TestEncodeURLValuesDecodeBack(t *testing.T) {
	// a value must survive one round of decoding unchanged (no double encoding)
	u, err := url.Parse(EncodeURL("https://h/p?v=a%20b%2Bc"))
	if err != nil {
		t.Fatal(err)
	}
	if got := u.Query().Get("v"); got != "a b+c" {
		t.Errorf("decoded value %q, want %q", got, "a b+c")
	}
}

func TestRemoveSliceDoesNotModifyInput(t *testing.T) {
	in := []string{"a", "b", "c"}
	if got := RemoveSlice(in, 1); !reflect.DeepEqual(got, []string{"a", "c"}) {
		t.Errorf("got %v", got)
	}
	if !reflect.DeepEqual(in, []string{"a", "b", "c"}) {
		t.Errorf("input modified: %v", in)
	}
	if got := RemoveSlice(in, 5); !reflect.DeepEqual(got, in) {
		t.Errorf("out of range: got %v", got)
	}
}

func TestCleanUrlPath(t *testing.T) {
	if got := CleanUrlPath("https://host//a///b/"); got != "https://host/a/b/" {
		t.Errorf("got %q", got)
	}
}
