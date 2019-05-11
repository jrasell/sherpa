package api

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest_setQueryOptions(t *testing.T) {
	testCases := []struct {
		testName        string
		request         *request
		queryOpts       *QueryOptions
		testGetKey      string
		expectedGetResp string
	}{
		{
			testName:        "Test case: count query parameter",
			request:         &request{params: make(map[string][]string)},
			queryOpts:       &QueryOptions{Params: map[string]string{"count": "100"}},
			testGetKey:      "count",
			expectedGetResp: "100",
		},
		{
			testName:        "Test case: empty query options",
			request:         &request{params: make(map[string][]string)},
			queryOpts:       nil,
			testGetKey:      "wontwork",
			expectedGetResp: "",
		},
	}

	for _, tc := range testCases {
		tc.request.setQueryOptions(tc.queryOpts)
		actualGetResp := tc.request.params.Get(tc.testGetKey)
		assert.Equal(t, tc.expectedGetResp, actualGetResp, tc.testName)
	}
}

func TestRequest_toHTTP(t *testing.T) {
	testCases := []struct {
		testName         string
		request          *request
		expectedHTTPResp *http.Request
		expectedError    error
	}{
		{
			testName:         "Test case: set PUT method with HTTP scheme",
			request:          &request{url: &url.URL{Scheme: "http", Host: "http://127.0.0.1:8000"}, method: "PUT"},
			expectedHTTPResp: generateRequest("PUT", &url.URL{Scheme: "http", Host: "http://127.0.0.1:8000"}),
			expectedError:    nil,
		},
		{
			testName:         "Test case: set POST method with HTTPS scheme",
			request:          &request{url: &url.URL{Scheme: "https", Host: "https://127.0.0.1:8000"}, method: "POST"},
			expectedHTTPResp: generateRequest("POST", &url.URL{Scheme: "https", Host: "https://127.0.0.1:8000"}),
			expectedError:    nil,
		},
		{
			testName:         "Test case: set GET method with HTTPS scheme on custom domain",
			request:          &request{url: &url.URL{Scheme: "https", Host: "https://sherpa.jrasell.com:8000"}, method: "GET"},
			expectedHTTPResp: generateRequest("GET", &url.URL{Scheme: "https", Host: "https://sherpa.jrasell.com:8000"}),
			expectedError:    nil,
		},
	}

	for _, tc := range testCases {
		actualHTTPResp, err := tc.request.toHTTP()
		assert.Equal(t, tc.expectedHTTPResp, actualHTTPResp, tc.testName)

		if tc.expectedError != nil {
			assert.EqualError(t, err, tc.expectedError.Error(), tc.testName)
		} else {
			assert.Nil(t, err, tc.testName)
		}
	}
}

func generateRequest(method string, url *url.URL) *http.Request {
	req, _ := http.NewRequest(method, url.RequestURI(), nil)
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Host = url.Host
	return req
}
