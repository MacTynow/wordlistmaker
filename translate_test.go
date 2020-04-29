package translator

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTranslate(t *testing.T) {
	tests := []struct {
		body string
		want string
	}{
		{body: `{"words": ["hello"]}`, want: "hello, 你好, nǐ hǎo\n"},
		{body: `{"words": ["hello", "my friend"]}`, want: "hello, 你好, nǐ hǎo\nmy friend, 我的朋友, wǒ de péng yǒu\n"},
		{body: `{"words": ["hello", "你好"]}`, want: "hello, 你好, nǐ hǎo\nhello there, 你好, nǐ hǎo\n"},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/", strings.NewReader(test.body))
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		Translate(rr, req)

		if got := rr.Body.String(); got != test.want {
			t.Errorf("Translate(%q) = %q, want %q", test.body, got, test.want)
		}
	}
}
