package hellosign_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/matryer/is"

	hellosign "github.com/milhamhidayat/go-hellosign-sdk"
	"github.com/milhamhidayat/go-hellosign-sdk/testdata"
)

func TestSignatureRequest_Get(t *testing.T) {
	is := is.New(t)

	signatureRequestJSON := testdata.GetGolden(t, "signature-request")

	signatureRequest := hellosign.SignatureRequest{}
	err := json.Unmarshal(signatureRequestJSON, &signatureRequest)
	is.NoErr(err)

	errNotFoundJSON := testdata.GetGolden(t, "signature-request-err-not-found")

	tests := map[string]struct {
		signatureRequestID       string
		signatureResponse        http.Response
		expectedSignatureRequest hellosign.SignatureRequest
		expectedError            error
	}{
		"success": {
			signatureRequestID: "fa5c8a0b0f492d768749333ad6fcc214c111e967",
			signatureResponse: http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewReader(signatureRequestJSON)),
				Header:     make(http.Header),
			},
			expectedSignatureRequest: signatureRequest,
			expectedError:            nil,
		},
		"not found": {
			signatureRequestID: "fa5c8a0b0f492d768749333ad6fcc214c111e967",
			signatureResponse: http.Response{
				StatusCode: http.StatusNotFound,
				Body:       ioutil.NopCloser(bytes.NewReader(errNotFoundJSON)),
				Header:     make(http.Header),
			},
			expectedSignatureRequest: hellosign.SignatureRequest{},
			expectedError:            errors.New("not_found: Not found"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			is := is.New(t)
			mockHTTPClient := testdata.NewClient(t, func(req *http.Request) *http.Response {
				return &test.signatureResponse
			})

			apiClient := hellosign.NewClient("123", mockHTTPClient)
			resp, err := apiClient.SignatureRequestAPI.Get(test.signatureRequestID)
			if err != nil {
				is.Equal(test.expectedError.Error(), err.Error())
				return
			}
			is.NoErr(err)
			is.Equal(test.expectedSignatureRequest, resp)
		})
	}
}
