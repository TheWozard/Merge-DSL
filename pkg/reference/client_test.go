package reference_test

import (
	"bytes"
	"io"
	"merge-dsl/pkg/reference"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmbeddedClient(t *testing.T) {
	data, info, err := (&reference.EmbeddedClient{Format: "example"}).Import("data", reference.Info{Reference: "ref:other"})
	require.Nil(t, err)
	require.Equal(t, reference.Info{
		Type:      "",
		Format:    "example",
		Reference: "ref:other",
	}, info)
	require.Equal(t, []byte("data"), data)
}

type TestTripper struct {
	t        *testing.T
	url      string
	response string
}

func (tt *TestTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	require.Equal(tt.t, tt.url, req.URL.String())
	return &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"content-type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewBufferString(tt.response)),
	}, nil
}

func TestHTTPClient(t *testing.T) {
	client := reference.HTTPClient{
		Protocol: "https",
		Client:   &http.Client{Transport: &TestTripper{t: t, url: "https://basic.url", response: "{}"}},
		Headers: map[string]string{
			"user-agent": "example",
		},
	}
	raw, info, err := client.Import("basic.url", reference.Info{Reference: "the-reference"})
	require.Nil(t, err)
	require.Equal(t, []byte("{}"), raw)
	require.Equal(t, reference.Info{
		Type:      "",
		Format:    "json",
		Reference: "the-reference",
	}, info)
}
