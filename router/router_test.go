package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/mirusky-dev/challenge-18/core"
	"github.com/mirusky-dev/challenge-18/core/background"
	"github.com/mirusky-dev/challenge-18/core/env"
	"github.com/mirusky-dev/challenge-18/models/dtos"
	"github.com/stretchr/testify/assert"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir) // HACK: Move to root
	if err != nil {
		panic(err)
	}
}

func TestRoutes(t *testing.T) {

	tokenUser1 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI3YWZmOWI2NC1kNTJhLTQ4NWMtYTczOC1mY2FlOWU1Y2VkZTAiLCJleHAiOjE3MjI1MDM1NTIsIm5iZiI6MTcyMTUwMzI1MiwiaWF0IjoxNzIxNTAzMjUyLCJqdGkiOiI3MTMxOTVmYy1kMGE2LTRiN2YtYmVjOC03ZTlmMmRiOGY5NTMiLCJyb2xlIjoidGVjaCJ9.sxvrrcXw98Qc6WKafIefg1HaF1lxNmoTVoITXHrHSZY"
	tokenUser2 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIwYTYyYzE3NS0xNTEyLTQzYjUtYWJkNC04YzViZjQ5ZjRhNDkiLCJleHAiOjE3MjI1MDM3MTAsIm5iZiI6MTcyMTUwMzQxMCwiaWF0IjoxNzIxNTAzNDEwLCJqdGkiOiI1M2NlMWQ5Yy0zNzFhLTQzMzYtOGJmMy1kMDY1NDlmNGUzYzkiLCJyb2xlIjoidGVjaCJ9.8oLMLXaOHtWIoFQm2uqh3k1sn4g83QRjk4ggDd5IYSk"
	tokenManager := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjQ3N2JjOS0zODMwLTQ4YWMtODJkNS1hOTFhNDliMGJkNzMiLCJleHAiOjE3MjI1MDM3NTUsIm5iZiI6MTcyMTUwMzQ1NSwiaWF0IjoxNzIxNTAzNDU1LCJqdGkiOiJhNDFlNDY5Zi1lMDc0LTRiOWMtOTZiZi1jMDY1YTQwODM5YTgiLCJyb2xlIjoibWFuYWdlciJ9.Y8-w8METGeXh3otgbMvpYK1uFdIt_USKFV985WPT20w"

	tests := []struct {
		description string

		// Test input
		route  string
		method string
		body   io.Reader
		header map[string][]string

		// Expected output
		expectedError bool
		expectedCode  int
		expectedBody  any
		ignoreBody    bool
	}{
		{
			description: "(Manager) Should list 2 tasks",
			route:       "/api/v1/tasks",
			method:      "GET",
			header: map[string][]string{
				"Authorization": {fmt.Sprintf("Bearer %s", tokenManager)},
			},
			expectedError: false,
			expectedCode:  200,
			expectedBody:  core.PagedResponse[dtos.Task]{},
		},
		{
			description: "(User1) Should list 1 task",
			route:       "/api/v1/tasks",
			method:      "GET",
			header: map[string][]string{
				"Authorization": {fmt.Sprintf("Bearer %s", tokenUser1)},
			},
			expectedError: false,
			expectedCode:  200,
			expectedBody:  core.PagedResponse[dtos.Task]{},
		},
		{
			description: "(User2) Should list 1 task",
			route:       "/api/v1/tasks",
			method:      "GET",
			header: map[string][]string{
				"Authorization": {fmt.Sprintf("Bearer %s", tokenUser2)},
			},
			expectedError: false,
			expectedCode:  200,
			expectedBody:  core.PagedResponse[dtos.Task]{},
		},
		{
			description: "(User2) Should perform 1 task",
			route:       "/api/v1/tasks/2f436752-0e23-4eb1-9963-ee8e9d04e972",
			method:      "PUT",
			header: map[string][]string{
				"Authorization": {fmt.Sprintf("Bearer %s", tokenUser2)},
				"Content-Type":  {"application/json"},
			},
			body:          bytes.NewBufferString(`{ "done": true }`),
			expectedError: false,
			expectedCode:  200,
			expectedBody:  dtos.Task{},
		},
		{
			description: "(User2) Attempt to read a task from another User",
			route:       "/api/v1/tasks/another-user-task-id",
			method:      "GET",
			header: map[string][]string{
				"Authorization": {fmt.Sprintf("Bearer %s", tokenUser2)},
			},
			expectedError: false,
			expectedCode:  404,
			expectedBody:  core.Exception{},
		},
		{
			description: "(User2) Should save a new task",
			route:       "/api/v1/tasks",
			method:      "POST",
			header: map[string][]string{
				"Authorization": {fmt.Sprintf("Bearer %s", tokenUser2)},
				"Content-Type":  {"application/json"},
			},
			body:          bytes.NewBufferString(`{ "summary": "my test task" }`),
			expectedError: false,
			expectedCode:  200,
			expectedBody:  dtos.Task{},
		},
	}

	config, err := env.Load()
	if err != nil {
		t.Error(err)
		return
	}

	backgroundJobClient, err := background.NewClient(config)
	if err != nil {
		t.Error(err)
		return
	}

	// Setup the app as it is done in the main function
	app := Setup(config, backgroundJobClient)

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route
		// from the test case

		req, _ := http.NewRequest(
			test.method,
			test.route,
			test.body,
		)

		for header, values := range test.header {
			for _, v := range values {
				req.Header.Add(header, v)
			}
		}

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := io.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		// Verify, that the reponse body equals the expected body

		err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&test.expectedBody)
		assert.Nilf(t, err, test.description)
	}
}
