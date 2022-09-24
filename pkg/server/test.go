package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type Test struct {
	Name               string
	Path               string
	Method             string
	RequestBody        interface{}
	Setup              func() interface{}
	ExpectedStatusCode int
}

func (s *Server) RunTests(tests []*Test, t *testing.T) {
	for i := range tests {
		test := tests[i]
		t.Run(test.Name, func(tt *testing.T) {
			expectedBody := test.Setup()

			var requestBody io.Reader
			if test.RequestBody != http.NoBody {
				requestBodyBytes, _ := json.Marshal(test.RequestBody)
				requestBody = bytes.NewBuffer(requestBodyBytes)
			}
			request, _ := http.NewRequest(test.Method, test.Path, requestBody)

			if test.RequestBody != http.NoBody {
				request.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			}

			response, err := s.Test(request)
			assert.Nil(tt, err)
			assert.NotNil(tt, response)
			if response != nil {
				defer response.Body.Close()
				assert.Equal(tt, test.ExpectedStatusCode, response.StatusCode)

				if expectedBody != nil {
					var expectedBodyBytes []byte
					if expectedBody != http.NoBody {
						expectedBodyBytes, _ = json.Marshal(expectedBody)
					}
					responseBody, _ := io.ReadAll(response.Body)
					assert.Equal(tt, string(expectedBodyBytes), string(responseBody))
				}
			}
		})
	}
}
