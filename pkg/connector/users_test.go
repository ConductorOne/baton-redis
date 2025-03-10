package connector

import (
	"context"
	encoding "encoding/base64"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/conductorone/baton-redis/pkg/client"
	"github.com/conductorone/baton-redis/test"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

// Tests that the client can fetch users based on the documented API below.
// https://redis.io/docs/latest/operate/rs/references/rest-api/requests/users/#get-all-users
func TestRedisClient_GetUsers(t *testing.T) {
	// Create a mock response.
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(test.ReadFile("usersMock.json"))),
	}
	mockResponse.Header.Set("Content-Type", "application/json")

	// Create a test client with the mock response.
	testClient := test.NewTestClient(mockResponse, nil)

	// Call GetUsers
	ctx := context.Background()
	result, nextOptions, err := testClient.ListUsers(ctx)

	// Check for errors.
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the result.
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Check count.
	if len(result) != 2 {
		t.Errorf("Expected Count to be 2, got %d", len(result))
	}

	for index, user := range result {
		expectedUser := client.User{
			AuthMethod:        "regular",
			Email:             fmt.Sprintf("testuser%d@redislabs.com", index+1),
			EmailAlerts:       true,
			Name:              fmt.Sprintf("Test User %d", index+1),
			PasswordIssueDate: test.GetUniqueTime(),
			Role:              test.Roles[index],
			RoleUIDs:          test.RoleUIDs[index],
			Status:            "active",
			UID:               index + 1,
		}

		if !reflect.DeepEqual(user, expectedUser) {
			t.Errorf("Unexpected user: got %+v, want %+v", user, expectedUser)
		}
	}

	// Check next options.
	if nextOptions == nil {
		t.Fatal("Expected non-nil nextOptions")
	}
}

func TestRedisClient_GetUsers_RequestDetails(t *testing.T) {
	// Create a custom RoundTripper to capture the request.
	var capturedRequest *http.Request
	mockTransport := &test.MockRoundTripper{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`[]`)),
			Header:     make(http.Header),
		},
		Err: nil,
	}
	mockTransport.Response.Header.Set("Content-Type", "application/json")

	mockRoundTrip := func(req *http.Request) (*http.Response, error) {
		capturedRequest = req
		return mockTransport.Response, mockTransport.Err
	}
	mockTransport.SetRoundTrip(mockRoundTrip)

	// Create a test client with the mock transport.
	httpClient := &http.Client{Transport: mockTransport}
	baseHttpClient := uhttp.NewBaseHttpClient(httpClient)
	testClient := client.NewClient("username", "password", "http://localhost", "8080", baseHttpClient)

	// Call GetUsers.
	ctx := context.Background()
	_, _, err := testClient.ListUsers(ctx)

	// Check for errors.
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the request details.
	if capturedRequest == nil {
		t.Fatal("No request was captured")
	}

	// Check URL components.
	expectedURL := "http://localhost:8080/v1/users"
	if capturedRequest.URL.String() != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, capturedRequest.URL.String())
	}

	// Check headers.
	authorizationToken := encoding.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", "username", "password")))
	expectedHeaders := map[string]string{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Basic %s", authorizationToken),
	}

	for key, expectedValue := range expectedHeaders {
		if value := capturedRequest.Header.Get(key); value != expectedValue {
			t.Errorf("Expected header %s to be %s, got %s", key, expectedValue, value)
		}
	}
}
